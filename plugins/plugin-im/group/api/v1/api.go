package v1

import (
	"time"
	"strconv"

	"hei-gin/sdk/middleware"
	"hei-gin/sdk/auth"
	authMW "hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/log"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	"hei-gin/plugins/plugin-im/group"

	"github.com/gin-gonic/gin"
)

func RegisterSysRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/sys/im/group").Use(authMW.HeiCheckLogin())
	{
		g.POST("/create", authMW.NoRepeat(3000), createHandler)
		g.GET("/my-groups", myGroupsHandler)
		g.GET("/detail", detailHandler)
		g.POST("/update", updateHandler)
		g.POST("/dissolve", log.SysLog("解散群"), dissolveHandler)

		g.POST("/invite", authMW.NoRepeat(3000), inviteHandler)
		g.POST("/join", joinHandler)
		g.GET("/pending-join-requests", pendingJoinRequestsHandler)
		g.POST("/handle-join-request", handleJoinRequestHandler)
		g.POST("/leave", leaveHandler)
		g.POST("/kick", log.SysLog("踢出成员"), kickHandler)
		g.POST("/set-role", log.SysLog("设置角色"), setRoleHandler)
		g.POST("/transfer-owner", log.SysLog("转让群"), transferOwnerHandler)
		g.POST("/set-nickname", setNicknameHandler)

		g.GET("/messages", messagesHandler)
		g.GET("/search", searchHandler)
		g.GET("/search-groups", searchGroupsHandler)
		g.POST("/send", middleware.RateLimiter("sys_group_send", 3, 20), authMW.NoRepeat(3000), sendHandler)
		g.POST("/recall", recallHandler)
		g.POST("/mark-read", markReadHandler)
		g.POST("/mute", log.SysLog("禁言"), muteHandler)
		g.POST("/unmute", log.SysLog("解禁"), unmuteHandler)
		g.GET("/members", membersHandler)
	}
}

func RegisterClientRoutes(r *gin.Engine) {
	g := r.Group("/api/v1/c/im/group").Use(authMW.HeiClientCheckLogin())
	{
		g.POST("/create", createHandler)
		g.GET("/my-groups", myGroupsHandler)
		g.GET("/detail", detailHandler)
		g.POST("/join", joinHandler)
		g.POST("/leave", leaveHandler)
		g.GET("/messages", messagesHandler)
		g.GET("/search", searchHandler)
		g.GET("/search-groups", searchGroupsHandler)
		g.POST("/send", middleware.RateLimiter("c_group_send", 3, 20), sendHandler)
		g.POST("/recall", recallHandler)
		g.POST("/mark-read", markReadHandler)
		g.GET("/members", membersHandler)
	}
}

// getLoginID returns the authenticated user based on route prefix.
func getLoginID(c *gin.Context) (string, string) {
	path := c.Request.URL.Path
	if len(path) > 8 && path[:8] == "/api/v1/c" {
		return auth.Consumer.GetLoginID(c), group.UserTypeConsumer
	}
	return auth.GetLoginID(c), group.UserTypeBusiness
}


func createHandler(c *gin.Context) {
	var p group.CreateParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	if userID == "" {
		result.Failure(c, "未登录", 401)
		return
	}
	g := group.Create(c.Request.Context(), userID, userType, &p)
	result.Success(c, g)
}

func updateHandler(c *gin.Context) {
	var p group.UpdateParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	group.Update(c.Request.Context(), userID, userType, &p)
	result.Success(c, nil)
}

func dissolveHandler(c *gin.Context) {
	var p struct{ GroupID string `json:"group_id"` }
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, _ := getLoginID(c)
	group.Dissolve(c.Request.Context(), userID, p.GroupID)
	result.Success(c, nil)
}

func detailHandler(c *gin.Context) {
	groupID := c.Query("group_id")
	vo := group.Detail(c.Request.Context(), groupID)
	result.Success(c, vo)
}

func myGroupsHandler(c *gin.Context) {
	userID, userType := getLoginID(c)
	list := group.MyGroups(c.Request.Context(), userID, userType)
	result.Success(c, list)
}

func inviteHandler(c *gin.Context) {
	var p group.InviteParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	group.Invite(c.Request.Context(), userID, userType, &p)
	result.Success(c, nil)
}

func joinHandler(c *gin.Context) {
	var p struct{ GroupID string `json:"group_id"` }
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	group.Join(c.Request.Context(), userID, userType, p.GroupID)
	result.Success(c, nil)
}

func pendingJoinRequestsHandler(c *gin.Context) {
	groupID := c.Query("group_id")
	requests := group.PendingJoinRequests(c.Request.Context(), groupID)
	result.Success(c, requests)
}

func handleJoinRequestHandler(c *gin.Context) {
	var p group.HandleJoinRequestParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	group.HandleJoinRequest(c.Request.Context(), userID, userType, &p)
	result.Success(c, nil)
}

func leaveHandler(c *gin.Context) {
	var p struct{ GroupID string `json:"group_id"` }
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	group.Leave(c.Request.Context(), userID, userType, p.GroupID)
	result.Success(c, nil)
}

func kickHandler(c *gin.Context) {
	var p group.KickParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	group.Kick(c.Request.Context(), userID, userType, &p)
	result.Success(c, nil)
}

func setRoleHandler(c *gin.Context) {
	var p group.SetRoleParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, _ := getLoginID(c)
	group.SetRole(c.Request.Context(), userID, &p)
	result.Success(c, nil)
}

func transferOwnerHandler(c *gin.Context) {
	var p group.TransferOwnerParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, _ := getLoginID(c)
	group.TransferOwner(c.Request.Context(), userID, &p)
	result.Success(c, nil)
}

func setNicknameHandler(c *gin.Context) {
	var p group.SetNicknameParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	group.SetMemberNickname(c.Request.Context(), userID, userType, &p)
	result.Success(c, nil)
}

func membersHandler(c *gin.Context) {
	groupID := c.Query("group_id")
	list := group.Members(c.Request.Context(), groupID)
	result.Success(c, list)
}

func messagesHandler(c *gin.Context) {
	groupID := c.Query("group_id")
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	msgs, hasMore := group.Messages(c.Request.Context(), groupID, cursor, size)
	result.Success(c, gin.H{"records": msgs, "has_more": hasMore})
}

func searchHandler(c *gin.Context) {
	groupID := c.Query("group_id")
	keyword := c.Query("keyword")
	cursor := c.Query("cursor")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	msgs, hasMore := group.SearchMessages(c.Request.Context(), groupID, keyword, cursor, size)
	result.Success(c, gin.H{"records": msgs, "has_more": hasMore})
}

func searchGroupsHandler(c *gin.Context) {
	keyword := c.Query("keyword")
	size := 20
	if s := c.Query("size"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			size = n
		}
	}
	list := group.SearchGroups(c.Request.Context(), keyword, size)
	result.Success(c, list)
}

func sendHandler(c *gin.Context) {
	var p group.SendMessageParam
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	vo := group.SendMessage(c.Request.Context(), userID, userType, &p)
	result.Success(c, vo)
}

func recallHandler(c *gin.Context) {
	var p struct {
		GroupID   string `json:"group_id"`
		MessageID string `json:"message_id"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	group.RecallMessage(c.Request.Context(), p.GroupID, p.MessageID, userID, userType)
	result.Success(c, nil)
}

func markReadHandler(c *gin.Context) {
	var p struct {
		GroupID   string `json:"group_id"`
		MessageID string `json:"message_id"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	group.MarkRead(c.Request.Context(), p.GroupID, userID, userType, p.MessageID)
	result.Success(c, nil)
}

func muteHandler(c *gin.Context) {
	var p struct {
		GroupID  string `json:"group_id"`
		UserID   string `json:"user_id"`
		UserType string `json:"user_type"`
		Duration int    `json:"duration"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	duration := 60
	if p.Duration > 0 {
		duration = p.Duration
	}
	kp := group.KickParam{GroupID: p.GroupID, UserID: p.UserID, UserType: p.UserType}
	group.MuteMember(c.Request.Context(), userID, userType, &kp, time.Duration(duration)*time.Minute)
	result.Success(c, nil)
}

func unmuteHandler(c *gin.Context) {
	var p struct {
		GroupID  string `json:"group_id"`
		UserID   string `json:"user_id"`
		UserType string `json:"user_type"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		result.Failure(c, "参数错误", 400)
		return
	}
	userID, userType := getLoginID(c)
	kp := group.KickParam{GroupID: p.GroupID, UserID: p.UserID, UserType: p.UserType}
	group.UnmuteMember(c.Request.Context(), userID, userType, &kp)
	result.Success(c, nil)
}

func init() {
	registry.RegisterRoute(RegisterSysRoutes)
	registry.RegisterRoute(RegisterClientRoutes)
}
