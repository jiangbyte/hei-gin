package username

import (
	"golang.org/x/crypto/bcrypt"

	cliUser "hei-gin/plugins/plugin-client/user"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/captcha"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/log"
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DoLogin handles consumer username/password login.
// @Summary      C端认证登录
// @Description  访问 /api/v1/public/c/login，C端认证登录
// @Tags         C端认证
// @Accept       json
// @Produce      json
// @Param        body  body  UsernameLoginParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/public/c/login [post]
func DoLogin(c *gin.Context) {
	ctx := c.Request.Context()
	var param UsernameLoginParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.WriteError(c, exception.NewBusinessError("请求参数错误", 400))
		return
	}

	if err := captcha.CCaptcha.CheckCaptcha(param.CaptchaID, param.CaptchaCode); err != nil {
		result.WriteError(c, exception.NewBusinessError(err.Error(), 400))
		return
	}

	var user cliUser.ClientUser
	if err := db.DB.WithContext(ctx).Where("username = ?", param.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("用户名或密码错误", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("系统异常", 500))
		return
	}

	status := user.Status
	switch status {
	case string(enums.UserStatusLocked):
		result.WriteError(c, exception.NewBusinessError("账号已被锁定", 400))
		return
	case string(enums.UserStatusInactive):
		result.WriteError(c, exception.NewBusinessError("账号已停用", 400))
		return
	default:
		if status != string(enums.UserStatusActive) {
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
		"status":      status,
		"device_type": utils.GetBrowser(ua),
		"device_id":   param.DeviceID,
	}

	clientAuth := auth.Consumer
	token, err := clientAuth.Login(c, user.ID, extra)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("登录失败", 500))
		return
	}

	ip := utils.GetClientIP(c)
	db.DB.WithContext(ctx).Model(&cliUser.ClientUser{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"last_login_ip": ip,
		"login_count":   gorm.Expr("login_count + 1"),
	})

	username := utils.SafeStrPtr(user.Username)
	log.RecordAuthLog(c, "登录", "LOGIN", "SUCCESS", "", username)

	result.Success(c, UsernameLoginResult{Token: token})
}

// DoRegister handles consumer user registration.
// @Summary      C端认证注册
// @Description  访问 /api/v1/public/c/register，C端认证注册
// @Tags         C端认证
// @Accept       json
// @Produce      json
// @Param        body  body  UsernameRegisterParam  true  "请求体"
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/public/c/register [post]
func DoRegister(c *gin.Context) {
	ctx := c.Request.Context()
	var param UsernameRegisterParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.WriteError(c, exception.NewBusinessError("请求参数错误", 400))
		return
	}

	if err := captcha.CCaptcha.CheckCaptcha(param.CaptchaID, param.CaptchaCode); err != nil {
		result.WriteError(c, exception.NewBusinessError(err.Error(), 400))
		return
	}

	var count int64
	db.DB.WithContext(ctx).Model(&cliUser.ClientUser{}).Where("username = ?", param.Username).Count(&count)
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
	entity := cliUser.ClientUser{
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

// DoLogout handles consumer user logout.
// @Summary      C端认证登出
// @Description  访问 /api/v1/c/logout，C端认证登出
// @Tags         C端认证
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any  "成功响应"
// @Router       /api/v1/c/logout [post]
func DoLogout(c *gin.Context) {
	clientAuth := auth.Consumer
	userID := clientAuth.GetLoginIDDefaultNull(c)
	if userID != "" {
		var user cliUser.ClientUser
		if err := db.DB.First(&user, "id = ?", userID).Error; err == nil {
			username := utils.SafeStrPtr(user.Username)
			log.RecordAuthLog(c, "登出", "LOGOUT", "SUCCESS", "", username)
		}
	}
	clientAuth.Logout(c)
	result.Success(c, UsernameLogoutResult{Message: "登出成功"})
}
