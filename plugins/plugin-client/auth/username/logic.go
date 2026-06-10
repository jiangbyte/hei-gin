package username

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/captcha"
	"hei-gin/sdk/db"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/log"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"
	cliUser "hei-gin/plugins/plugin-client/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DoLogin(c *gin.Context) {
	ctx := c.Request.Context()
	var param UsernameLoginParam
	if err := c.ShouldBindJSON(&param); err != nil {
		panic(exception.NewBusinessError("请求参数错误", 400))
	}

	

	if err := captcha.CCaptcha.CheckCaptcha(param.CaptchaID, param.CaptchaCode); err != nil {
		panic(exception.NewBusinessError(err.Error(), 400))
	}

	var user cliUser.ClientUser
	if err := db.DB.WithContext(ctx).Where("username = ?", param.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(exception.NewBusinessError("用户名或密码错误", 400))
		}
		panic(exception.NewBusinessError("系统异常", 500))
	}

	status := user.Status
	switch status {
	case string(enums.UserStatusLocked):
		panic(exception.NewBusinessError("账号已被锁定", 400))
	case string(enums.UserStatusInactive):
		panic(exception.NewBusinessError("账号已停用", 400))
	default:
		if status != string(enums.UserStatusActive) {
			panic(exception.NewBusinessError("账号状态异常", 400))
		}
	}

	rawPwd := utils.Decrypt(param.Password)
	if rawPwd == "" {
		panic(exception.NewBusinessError("用户名或密码错误", 400))
	}
	if user.Password == nil {
		panic(exception.NewBusinessError("用户名或密码错误", 400))
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(rawPwd)); err != nil {
		panic(exception.NewBusinessError("用户名或密码错误", 400))
	}

	ua := c.GetHeader("User-Agent")
	extra := map[string]any{
		"username":    safeStr(user.Username),
		"nickname":    safeStr(user.Nickname),
		"status":      status,
		"device_type": utils.GetBrowser(ua),
		"device_id":   param.DeviceID,
	}

	clientAuth := auth.Consumer
	token, err := clientAuth.Login(c, user.ID, extra)
	if err != nil {
		panic(exception.NewBusinessError("登录失败", 500))
	}

	now := time.Now()
	ip := utils.GetClientIP(c)
	db.DB.WithContext(ctx).Model(&cliUser.ClientUser{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ip,
		"login_count":   gorm.Expr("login_count + 1"),
	})

	username := safeStr(user.Username)
	log.RecordAuthLog(c, "登录", "LOGIN", "SUCCESS", "", username)

	result.Success(c, UsernameLoginResult{Token: token})
}

func DoRegister(c *gin.Context) {
	ctx := c.Request.Context()
	var param UsernameRegisterParam
	if err := c.ShouldBindJSON(&param); err != nil {
		panic(exception.NewBusinessError("请求参数错误", 400))
	}

	ctx = c.Request.Context()

	if err := captcha.CCaptcha.CheckCaptcha(param.CaptchaID, param.CaptchaCode); err != nil {
		panic(exception.NewBusinessError(err.Error(), 400))
	}

	var count int64
	db.DB.WithContext(ctx).Model(&cliUser.ClientUser{}).Where("username = ?", param.Username).Count(&count)
	if count > 0 {
		panic(exception.NewBusinessError("用户名已存在", 400))
	}

	rawPwd := utils.Decrypt(param.Password)
	if rawPwd == "" {
		panic(exception.NewBusinessError("密码解密失败", 400))
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(rawPwd), bcrypt.DefaultCost)
	if err != nil {
		panic(exception.NewBusinessError("密码加密失败", 500))
	}

	userID := utils.GenerateID()
	now := time.Now()
	hashedPwdStr := string(hashedPwd)

	entity := cliUser.ClientUser{
		ID: userID, Username: &param.Username, Password: &hashedPwdStr,
		Nickname: &param.Username, Status: string(enums.UserStatusActive),
		CreatedAt: &now, CreatedBy: &userID,
	}
	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil {
		panic(exception.NewBusinessError("注册失败", 500))
	}

	log.RecordAuthLog(c, "注册", "REGISTER", "SUCCESS", "", param.Username)

	result.Success(c, UsernameRegisterResult{Message: "注册成功"})
}

func DoLogout(c *gin.Context) {
	clientAuth := auth.Consumer
	userID := clientAuth.GetLoginIDDefaultNull(c)
	if userID != "" {
		var user cliUser.ClientUser
		if err := db.DB.First(&user, "id = ?", userID).Error; err == nil {
			username := safeStr(user.Username)
			log.RecordAuthLog(c, "登出", "LOGOUT", "SUCCESS", "", username)
		}
	}
	clientAuth.Logout(c)
	result.Success(c, UsernameLogoutResult{Message: "登出成功"})
}

func safeStr(s *string) string {
	if s == nil { return "" }
	return *s
}
