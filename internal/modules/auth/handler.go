package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
)

func (s *Service) registerRoutes(api *gin.RouterGroup) {
	rdb := s.repo.rdb

	api.GET("/v1/admin/captcha", middleware.RateLimit(rdb, "admin:captcha", 30, 60), s.captcha)
	api.GET("/v1/portal/captcha", middleware.RateLimit(rdb, "portal:captcha", 30, 60), s.captcha)
	api.GET("/v1/admin/password-key", s.passwordKey)
	api.GET("/v1/portal/password-key", s.passwordKey)

	api.POST("/v1/admin/login", middleware.RateLimit(rdb, "admin:login", 20, 60), s.login(security.AccountAdmin))
	api.POST("/v1/portal/login", middleware.RateLimit(rdb, "portal:login", 20, 60), s.login(security.AccountPortal))

	api.POST("/v1/admin/send-login-code", middleware.RateLimit(rdb, "admin:send-login-code", 10, 60), s.sendLoginCode(security.AccountAdmin))
	api.POST("/v1/portal/send-login-code", middleware.RateLimit(rdb, "portal:send-login-code", 10, 60), s.sendLoginCode(security.AccountPortal))

	api.POST("/v1/admin/forgot-password", middleware.RateLimit(rdb, "admin:forgot-password", 5, 60), s.forgotPassword(security.AccountAdmin))
	api.POST("/v1/portal/forgot-password", middleware.RateLimit(rdb, "portal:forgot-password", 5, 60), s.forgotPassword(security.AccountPortal))
	api.POST("/v1/admin/reset-password", middleware.RateLimit(rdb, "admin:reset-password", 10, 60), s.resetPassword(security.AccountAdmin))
	api.POST("/v1/portal/reset-password", middleware.RateLimit(rdb, "portal:reset-password", 10, 60), s.resetPassword(security.AccountPortal))

	api.POST("/v1/admin/logout", middleware.RequireAccountType(security.AccountAdmin), s.logout)
	api.POST("/v1/portal/logout", middleware.RequireAccountType(security.AccountPortal), s.logout)

	api.GET("/v1/admin/public/auth-options", s.authOptions(security.AccountAdmin))
	api.GET("/v1/portal/public/auth-options", s.authOptions(security.AccountPortal))

	api.POST("/v1/admin/auth/refresh", middleware.RequireAccountType(security.AccountAdmin), s.refresh(security.AccountAdmin))
	api.POST("/v1/portal/auth/refresh", middleware.RequireAccountType(security.AccountPortal), s.refresh(security.AccountPortal))

	api.POST("/v1/admin/cancel", middleware.RequireAccountType(security.AccountAdmin), s.cancel(security.AccountAdmin))
	api.POST("/v1/portal/cancel", middleware.RequireAccountType(security.AccountPortal), s.cancel(security.AccountPortal))

	api.POST("/v1/portal/register/send-code", middleware.RateLimit(rdb, "portal:register-send-code", 10, 60), s.registerSendCode)

	api.POST("/v1/portal/register", s.register)

	if s.oauth != nil {
		s.oauth.RegisterRoutes(api)
		s.oauth.RegisterBindingRoutes(api, s.perms)
	}
	s.registerSessionRoutes(api)
}

func (s *Service) captcha(c *gin.Context) {
	out, err := s.Captcha(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) passwordKey(c *gin.Context) {
	out, err := s.PasswordKey(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) login(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginParam
		if err := bind.JSON(c, &req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		out, err := s.Login(c.Request.Context(), accountType, req, c.ClientIP(), c.Request.UserAgent())
		if err != nil {
			switch err {
			case errInvalidCredentials, errInvalidOTP, errAccountLocked, errIPLocked:
				response.Fail(c, http.StatusUnauthorized, 401, err.Error())
			case errAccountFinder:
				response.Fail(c, http.StatusInternalServerError, 500, err.Error())
			default:
				response.Fail(c, http.StatusBadRequest, 400, err.Error())
			}
			return
		}
		ttlSec := s.cfg.Auth.TokenTTLSeconds
		if !req.RememberMe && s.cfg.Auth.TokenTTLShortSeconds > 0 {
			ttlSec = s.cfg.Auth.TokenTTLShortSeconds
		}
		if ttlSec <= 0 {
			ttlSec = 14400
		}
		s.SetSessionCookie(c, out.Token, accountType, time.Duration(ttlSec)*time.Second, req.RememberMe)
		response.OK(c, out)
	}
}

func (s *Service) logout(c *gin.Context) {
	token := s.ResolveLogoutToken(c)
	sess := contextx.Session(c.Request.Context())
	accountID, accountType := "", ""
	if sess != nil {
		accountID = sess.AccountID
		accountType = string(sess.AccountType)
	}
	_ = s.Logout(c.Request.Context(), token, accountID, accountType, c.ClientIP(), c.Request.UserAgent())
	s.ClearSessionCookie(c, contextx.AccountType(c.Request.Context()))
	response.OK(c, LogoutResult{Success: true})
}

func (s *Service) register(c *gin.Context) {
	var req RegisterParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := s.Register(c.Request.Context(), req)
	if err != nil {
		switch err {
		case errRegisterDisabled:
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
		case errPortalRegistrar:
			response.Fail(c, http.StatusNotImplemented, 501, err.Error())
		default:
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
		}
		return
	}
	response.OK(c, out)
}

func (s *Service) sendLoginCode(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SendLoginCodeParam
		if err := bind.JSON(c, &req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		if err := s.SendLoginCode(c.Request.Context(), accountType, req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, nil)
	}
}

func (s *Service) forgotPassword(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ForgotPasswordParam
		if err := bind.JSON(c, &req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		if err := s.ForgotPassword(c.Request.Context(), accountType, req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, nil)
	}
}

func (s *Service) resetPassword(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ResetPasswordParam
		if err := bind.JSON(c, &req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		if err := s.ResetPassword(c.Request.Context(), accountType, req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, nil)
	}
}

func (s *Service) authOptions(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		response.OK(c, s.AuthOptions(c.Request.Context(), accountType))
	}
}

func (s *Service) refresh(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		out, err := s.RefreshSession(c.Request.Context(), accountType, c.ClientIP(), c.Request.UserAgent())
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, 401, err.Error())
			return
		}
		response.OK(c, out)
	}
}

func (s *Service) cancel(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CancelParam
		_ = bind.JSON(c, &req)
		if err := s.CancelAccount(c.Request.Context(), accountType, c.ClientIP(), c.Request.UserAgent(), req.CancelReason); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		s.ClearSessionCookie(c, contextx.AccountType(c.Request.Context()))
		response.OK(c, LogoutResult{Success: true})
	}
}

func (s *Service) registerSendCode(c *gin.Context) {
	var req SendLoginCodeParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.sendRegisterCode(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}
