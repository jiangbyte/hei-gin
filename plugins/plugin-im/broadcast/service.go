package broadcast

import (
	"gorm.io/gorm"
	imModel "hei-gin/plugins/plugin-im/model"

	"hei-gin/plugins/plugin-im/ws"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *repository
}

// ==================== BroadcastSend ====================

func (s *Service) Send(c *gin.Context, p *SendBroadcastParam) {
	ctx := c.Request.Context()
	senderID := auth.Business.GetLoginID(c)

	if p.Scope == "" {
		p.Scope = "ALL"
	}

	e := imModel.Broadcast{
		ID:       utils.GenerateID(),
		Title:    p.Title,
		Content:  p.Content,
		Scope:    p.Scope,
		SenderID: senderID,
	}
	if err := s.repo.Create(ctx, &e); err != nil {
		result.WriteError(c, exception.NewBusinessError("发送通知失败: "+err.Error(), 500))
		return
	}

	// WS broadcast
	payload := map[string]interface{}{
		"title":   p.Title,
		"content": p.Content,
		"scope":   p.Scope,
		"action":  "broadcast",
	}
	msg := ws.Message{Type: "broadcast", Payload: payload}
	if runtime := ws.Runtime(); runtime != nil {
		switch p.Scope {
		case "ALL":
			runtime.BroadcastAll(msg)
		case string(enums.LoginTypeBusiness):
			runtime.BroadcastBusiness(msg)
		case string(enums.LoginTypeConsumer):
			runtime.BroadcastConsumers(msg)
		}
	}
}

// ==================== BroadcastPage (admin) ====================

func (s *Service) Page(c *gin.Context, cursor string, size int) ([]BroadcastVO, bool) {
	ctx := c.Request.Context()
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	records := s.repo.Page(ctx, cursor, size)

	hasMore := len(records) > size
	if hasMore {
		records = records[:size]
	}

	vos := make([]BroadcastVO, len(records))
	for i, b := range records {
		vos[i] = *ImBroadcastToBroadcastVO(&b)
	}
	return vos, hasMore
}

// ==================== BroadcastUnreadList ====================

func (s *Service) UnreadList(c *gin.Context, userType string) []BroadcastVO {
	ctx := c.Request.Context()
	var userID string
	if userType == string(enums.LoginTypeConsumer) {
		userID = auth.Consumer.GetLoginID(c)
	} else {
		userID = auth.Business.GetLoginID(c)
	}

	records := s.repo.ListLatest(ctx, 50)
	readRecords := s.repo.ListReads(ctx, userID, userType)
	readMap := make(map[string]bool)
	for _, r := range readRecords {
		readMap[r.BroadcastID] = true
	}

	vos := make([]BroadcastVO, 0, len(records))
	for _, b := range records {
		vo := *ImBroadcastToBroadcastVO(&b)
		if _, ok := readMap[b.ID]; ok {
			vo.Read = true
		}
		vos = append(vos, vo)
	}
	return vos
}

// ==================== BroadcastMarkRead ====================

func (s *Service) MarkRead(c *gin.Context, p *ReadParam) {
	ctx := c.Request.Context()
	userID := auth.Business.GetLoginID(c)

	s.repo.MarkRead(ctx, p.BroadcastID, userID, string(enums.LoginTypeBusiness))
}

// ==================== BroadcastMarkReadConsumer ====================

func (s *Service) MarkReadConsumer(c *gin.Context, p *ReadParam) {
	ctx := c.Request.Context()
	userID := auth.Consumer.GetLoginID(c)

	s.repo.MarkRead(ctx, p.BroadcastID, userID, string(enums.LoginTypeConsumer))
}

// ==================== BroadcastDetail ====================

func (s *Service) Detail(c *gin.Context, id string) *BroadcastVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("通知不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询通知失败: "+err.Error(), 500))
		return nil
	}
	return ImBroadcastToBroadcastVO(b)
}
