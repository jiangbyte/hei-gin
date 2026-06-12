package user

import (
	"gorm.io/gorm"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *repository
}

func (s *Service) Page(c *gin.Context, p *ClientUserPageParam) {
	ctx := c.Request.Context()
	if p.Current < 1 {
		p.Current = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}

	rows, total := s.repo.Page(ctx, p)

	vos := make([]*ClientUserVO, len(rows))
	for i, r := range rows {
		vos[i] = ClientUserToClientUserVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func (s *Service) Detail(c *gin.Context, id string) *ClientUserVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询用户详情失败: "+err.Error(), 500))
		return nil
	}
	return ClientUserToClientUserVO(e)
}

func (s *Service) Create(c *gin.Context, v *ClientUserVO) {
	ctx := c.Request.Context()

	if v.Username != nil {
		cnt := s.repo.CountByUsername(ctx, *v.Username, "")
		if cnt > 0 {
			result.WriteError(c, exception.NewBusinessError("帐号已存在", 400))
			return
		}
	}

	e := ClientUserVOToClientUser(v)
	e.Status = string(enums.UserStatusActive)

	if v.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*v.Password), bcrypt.DefaultCost)
		if err != nil {
			result.WriteError(c, exception.NewBusinessError("密码加密失败", 500))
			return
		}
		sv := string(hashed)
		e.Password = &sv
	}

	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加用户失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Modify(c *gin.Context, v *ClientUserVO) {
	ctx := c.Request.Context()
	if v.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	_, err := s.repo.FindByID(ctx, v.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
		return
	}

	up := map[string]interface{}{}
	if v.Nickname != nil {
		up["nickname"] = *v.Nickname
	}
	if v.Avatar != nil {
		up["avatar"] = *v.Avatar
	}
	if v.Email != nil {
		up["email"] = *v.Email
	}
	if v.Phone != nil {
		up["phone"] = *v.Phone
	}
	if v.Status != "" {
		up["status"] = v.Status
	}
	if len(up) == 0 {
		return
	}
	if err := s.repo.UpdateByID(ctx, v.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑用户失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Remove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	if err := s.repo.DeleteByIDs(c.Request.Context(), ids); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除用户失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Current(c *gin.Context) *ClientUserVO {
	userID := auth.Consumer.GetLoginIDDefaultNull(c)
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("用户未登录", 401))
		return nil
	}

	e, err := s.repo.FindByID(c.Request.Context(), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("用户不存在", 404))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
		return nil
	}
	return ClientUserToClientUserVO(e)
}

func (s *Service) UpdateProfile(c *gin.Context, param *UpdateProfileParam) {
	userID := auth.Consumer.GetLoginIDDefaultNull(c)
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("用户未登录", 401))
		return
	}
	ctx := c.Request.Context()

	if param.Username != nil && *param.Username != "" {
		count := s.repo.CountByUsername(ctx, *param.Username, userID)
		if count > 0 {
			result.WriteError(c, exception.NewBusinessError("用户名已存在", 400))
			return
		}
	}

	up := map[string]interface{}{}
	if param.Nickname != nil {
		up["nickname"] = *param.Nickname
	}
	if param.Username != nil {
		up["username"] = *param.Username
	}
	if param.Avatar != nil {
		up["avatar"] = *param.Avatar
	}
	if param.Email != nil {
		up["email"] = *param.Email
	}
	if param.Phone != nil {
		up["phone"] = *param.Phone
	}
	if len(up) == 0 {
		return
	}
	if err := s.repo.UpdateByID(ctx, userID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("更新个人信息失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) UpdateAvatar(c *gin.Context, param *UpdateAvatarParam) {
	userID := auth.Consumer.GetLoginIDDefaultNull(c)
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("用户未登录", 401))
		return
	}
	if param.Avatar == "" {
		result.WriteError(c, exception.NewBusinessError("头像不能为空", 400))
		return
	}

	avatar := utils.CompressBase64Image(param.Avatar, 512, 512, 80)

	ctx := c.Request.Context()
	e, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("用户不存在", 404))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
		return
	}
	if err := s.repo.UpdateAvatar(ctx, e, avatar); err != nil {
		result.WriteError(c, exception.NewBusinessError("保存头像失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) UpdatePassword(c *gin.Context, param *UpdatePasswordParam) {
	userID := auth.Consumer.GetLoginIDDefaultNull(c)
	if userID == "" {
		result.WriteError(c, exception.NewBusinessError("用户未登录", 401))
		return
	}
	ctx := c.Request.Context()
	e, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("用户不存在", 404))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
		return
	}
	if e.Password == nil || *e.Password == "" {
		result.WriteError(c, exception.NewBusinessError("未设置密码，无法修改", 400))
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(*e.Password), []byte(utils.Decrypt(param.CurrentPassword))) != nil {
		result.WriteError(c, exception.NewBusinessError("当前密码不正确", 400))
		return
	}
	newPwd := utils.Decrypt(param.NewPassword)
	if newPwd == "" {
		result.WriteError(c, exception.NewBusinessError("新密码解密失败", 400))
		return
	}
	h, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("密码加密失败", 500))
		return
	}
	if err := s.repo.UpdatePassword(ctx, userID, string(h)); err != nil {
		result.WriteError(c, exception.NewBusinessError("修改密码失败: "+err.Error(), 500))
		return
	}
}
