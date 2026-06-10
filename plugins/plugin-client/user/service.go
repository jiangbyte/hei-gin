package user

import (
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ClientUserPage(c *gin.Context, param *ClientUserPageParam) {
	ctx := c.Request.Context()
	if param.Current < 1 {
		param.Current = 1
	}
	if param.Size < 1 || param.Size > 100 {
		param.Size = 10
	}

	q := db.DB.WithContext(ctx).Model(&ClientUser{})
	if param.Keyword != "" {
		like := "%" + param.Keyword + "%"
		q = q.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?", like, like, like, like)
	}
	if param.Status != "" {
		q = q.Where("status = ?", param.Status)
	}

	var total int64
	q.Count(&total)

	var records []ClientUser
	q.Order("created_at DESC").Limit(param.Size).Offset((param.Current - 1) * param.Size).Find(&records)

	vos := make([]ClientUserVO, len(records))
	for i, r := range records {
		vos[i] = toVO(&r)
	}

	result.PageDataResult(c, vos, total, param.Current, param.Size)
}

func ClientUserDetail(c *gin.Context, id string) *ClientUserVO {
	if id == "" {
		return nil
	}
	ctx := c.Request.Context()
	var entity ClientUser
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		panic(exception.NewBusinessError("查询用户详情失败: "+err.Error(), 500))
	}
	rec := toVO(&entity)
	return &rec
}

func ClientUserCreate(c *gin.Context, vo *ClientUserVO, userID string) {
	ctx := c.Request.Context()
	now := time.Now()

	entity := ClientUser{
		ID:        utils.GenerateID(),
		Status:    string(enums.UserStatusActive),
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	if vo.Username != nil {
		var count int64
		db.DB.WithContext(ctx).Model(&ClientUser{}).Where("username = ?", *vo.Username).Count(&count)
		if count > 0 {
			panic(exception.NewBusinessError("帐号已存在", 400))
		}
		entity.Username = vo.Username
	}
	if vo.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*vo.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(exception.NewBusinessError("密码加密失败", 500))
		}
		s := string(hashed)
		entity.Password = &s
	}
	if vo.Nickname != nil {
		entity.Nickname = vo.Nickname
	}
	if vo.Avatar != nil {
		entity.Avatar = vo.Avatar
	}
	if vo.Email != nil {
		entity.Email = vo.Email
	}
	if vo.Phone != nil {
		entity.Phone = vo.Phone
	}
	if userID != "" {
		entity.CreatedBy = &userID
		entity.UpdatedBy = &userID
	}

	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil {
		panic(exception.NewBusinessError("添加用户失败: "+err.Error(), 500))
	}
}

func ClientUserModify(c *gin.Context, vo *ClientUserVO, userID string) {
	ctx := c.Request.Context()
	var entity ClientUser
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", vo.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(exception.NewBusinessError("数据不存在", 400))
		}
		panic(exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
	}

	updates := map[string]interface{}{"updated_at": time.Now()}
	if vo.Nickname != nil {
		updates["nickname"] = *vo.Nickname
	}
	if vo.Avatar != nil {
		updates["avatar"] = *vo.Avatar
	}
	if vo.Email != nil {
		updates["email"] = *vo.Email
	}
	if vo.Phone != nil {
		updates["phone"] = *vo.Phone
	}
	if vo.Status != "" {
		updates["status"] = vo.Status
	}
	if userID != "" {
		updates["updated_by"] = userID
	}

	if err := db.DB.WithContext(ctx).Model(&ClientUser{}).Where("id = ?", vo.ID).Updates(updates).Error; err != nil {
		panic(exception.NewBusinessError("编辑用户失败: "+err.Error(), 500))
	}
}

func ClientUserRemove(c *gin.Context, ids []string) {
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&ClientUser{})
}

func Current(c *gin.Context, userID string) *ClientUserVO {
	if userID == "" {
		return nil
	}
	return ClientUserDetail(c, userID)
}

func UpdateProfile(c *gin.Context, userID string, param *UpdateProfileParam) {
	if userID == "" {
		return
	}
	ctx := c.Request.Context()

	if param.Username != nil && *param.Username != "" {
		var count int64
		db.DB.WithContext(ctx).Model(&ClientUser{}).Where("username = ? AND id != ?", *param.Username, userID).Count(&count)
		if count > 0 {
			panic(exception.NewBusinessError("用户名已存在", 400))
		}
	}

	updates := map[string]interface{}{"updated_at": time.Now()}
	if param.Nickname != nil {
		updates["nickname"] = *param.Nickname
	}
	if param.Username != nil {
		updates["username"] = *param.Username
	}
	if param.Avatar != nil {
		updates["avatar"] = *param.Avatar
	}
	if param.Email != nil {
		updates["email"] = *param.Email
	}
	if param.Phone != nil {
		updates["phone"] = *param.Phone
	}
	db.DB.WithContext(ctx).Model(&ClientUser{}).Where("id = ?", userID).Updates(updates)
}

func UpdateAvatar(c *gin.Context, userID string, param *UpdateAvatarParam) {
	if userID == "" {
		panic(exception.NewBusinessError("用户未登录", 401))
	}
	if param.Avatar == "" {
		panic(exception.NewBusinessError("头像不能为空", 400))
	}

	avatar := utils.CompressBase64Image(param.Avatar, 512, 512, 80)

	ctx := c.Request.Context()
	var entity ClientUser
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(exception.NewBusinessError("用户不存在", 404))
		}
		panic(exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
	}
	if err := db.DB.WithContext(ctx).Model(&entity).Update("avatar", avatar).Error; err != nil {
		panic(exception.NewBusinessError("保存头像失败: "+err.Error(), 500))
	}
}

func UpdatePassword(c *gin.Context, userID string, param *UpdatePasswordParam) {
	if userID == "" {
		panic(exception.NewBusinessError("用户未登录", 401))
	}
	ctx := c.Request.Context()
	var entity ClientUser
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(exception.NewBusinessError("用户不存在", 404))
		}
		panic(exception.NewBusinessError("查询用户失败: "+err.Error(), 500))
	}
	if entity.Password == nil || *entity.Password == "" {
		panic(exception.NewBusinessError("未设置密码，无法修改", 400))
	}
	if bcrypt.CompareHashAndPassword([]byte(*entity.Password), []byte(utils.Decrypt(param.CurrentPassword))) != nil {
		panic(exception.NewBusinessError("当前密码不正确", 400))
	}
	h, _ := bcrypt.GenerateFromPassword([]byte(utils.Decrypt(param.NewPassword)), bcrypt.DefaultCost)
	db.DB.WithContext(ctx).Model(&ClientUser{}).Where("id = ?", userID).Update("password", string(h))
}
