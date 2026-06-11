package username

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	userModel "hei-gin/plugins/plugin-sys/user"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/captcha"
	"hei-gin/sdk/config"
	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/log"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"
)

// DoLogin handles username/password login.
func DoLogin(c *gin.Context) {
	var param UsernameLoginParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.WriteError(c, exception.NewBusinessError("请求参数错误", 400))
		return
	}
	ctx := c.Request.Context()
	if err := captcha.BCaptcha.CheckCaptcha(param.CaptchaID, param.CaptchaCode); err != nil {
		result.WriteError(c, exception.NewBusinessError(err.Error(), 400))
		return
	}
	var user userModel.SysUser
	if err := db.DB.WithContext(ctx).Where("username = ?", param.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("用户名或密码错误", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("系统异常", 500))
		return
	}
	switch user.Status {
	case string(enums.UserStatusLocked):
		result.WriteError(c, exception.NewBusinessError("账号已被锁定", 400))
		return
	case string(enums.UserStatusInactive):
		result.WriteError(c, exception.NewBusinessError("账号已停用", 400))
		return
	default:
		if user.Status != string(enums.UserStatusActive) {
			result.WriteError(c, exception.NewBusinessError("账号状态异常", 400))
			return
		}
	}
	rawPwd := utils.Decrypt(param.Password)
	if rawPwd == "" {
		result.WriteError(c, exception.NewBusinessError("用户名或密码错误", 400))
		return
	}
	if user.Password == nil {
		result.WriteError(c, exception.NewBusinessError("用户名或密码错误", 400))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(rawPwd)); err != nil {
		result.WriteError(c, exception.NewBusinessError("用户名或密码错误", 400))
		return
	}
	ua := c.GetHeader("User-Agent")
	extra := map[string]any{
		"username":    utils.SafeStrPtr(user.Username),
		"nickname":    utils.SafeStrPtr(user.Nickname),
		"status":      user.Status,
		"device_type": utils.GetBrowser(ua),
		"device_id":   param.DeviceID,
	}
	tokenStr, err := auth.Login(c, user.ID, extra)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("登录失败", 500))
		return
	}
	ip := utils.GetClientIP(c)
	db.DB.WithContext(ctx).Model(&userModel.SysUser{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"last_login_ip": ip,
		"login_count":   gorm.Expr("login_count + 1"),
	})
	username := utils.SafeStrPtr(user.Username)
	log.RecordAuthLog(c, "登录", "LOGIN", "SUCCESS", "", username)
	result.Success(c, UsernameLoginResult{Token: tokenStr})
}

// DoRegister handles user registration.
func DoRegister(c *gin.Context) {
	if !config.C.Auth.BusinessRegisterEnabled {
		result.WriteError(c, exception.NewBusinessError("后台注册未开放", 403))
		return
	}
	var param UsernameRegisterParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.WriteError(c, exception.NewBusinessError("请求参数错误", 400))
		return
	}
	ctx := c.Request.Context()
	if err := captcha.BCaptcha.CheckCaptcha(param.CaptchaID, param.CaptchaCode); err != nil {
		result.WriteError(c, exception.NewBusinessError(err.Error(), 400))
		return
	}
	var count int64
	db.DB.WithContext(ctx).Model(&userModel.SysUser{}).Where("username = ?", param.Username).Count(&count)
	if count > 0 {
		result.WriteError(c, exception.NewBusinessError("用户名已存在", 400))
		return
	}
	rawPwd := utils.Decrypt(param.Password)
	if rawPwd == "" {
		result.WriteError(c, exception.NewBusinessError("密码解密失败", 400))
		return
	}
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(rawPwd), bcrypt.DefaultCost)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("密码加密失败", 500))
		return
	}
	hashedPwdStr := string(hashedPwd)
	entity := userModel.SysUser{
		Username: &param.Username, Password: &hashedPwdStr,
		Nickname: &param.Username, Status: string(enums.UserStatusActive),
	}
	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("注册失败", 500))
		return
	}
	log.RecordAuthLog(c, "注册", "REGISTER", "SUCCESS", "", param.Username)
	result.Success(c, UsernameRegisterResult{Message: "注册成功"})
}

// DoLogout handles user logout.
func DoLogout(c *gin.Context) {
	userID := auth.GetLoginIDDefaultNull(c)
	if userID != "" {
		var user userModel.SysUser
		if err := db.DB.First(&user, "id = ?", userID).Error; err == nil {
			username := utils.SafeStrPtr(user.Username)
			log.RecordAuthLog(c, "登出", "LOGOUT", "SUCCESS", "", username)
		}
	}
	auth.Logout(c)
	result.Success(c, UsernameLogoutResult{Message: "登出成功"})
}
