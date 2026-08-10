package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	contextx "hei-gin/framework/core/context"
	"hei-gin/framework/core/bind"
	"hei-gin/framework/core/response"
	"hei-gin/framework/core/security"
	"hei-gin/framework/middleware"
)

func (s *Service) registerRoutes(api *gin.RouterGroup) {
	api.GET("/v1/admin/captcha", s.captcha)
	api.GET("/v1/portal/captcha", s.captcha)
	api.GET("/v1/admin/password-key", s.passwordKey)
	api.GET("/v1/portal/password-key", s.passwordKey)

	api.POST("/v1/admin/login", s.login(security.AccountAdmin))
	api.POST("/v1/portal/login", s.login(security.AccountPortal))

	api.POST("/v1/admin/logout", middleware.RequireAccountType(security.AccountAdmin), s.logout)
	api.POST("/v1/portal/logout", middleware.RequireAccountType(security.AccountPortal), s.logout)

	api.POST("/v1/portal/register", s.register)
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
			case errInvalidCredentials:
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
	_ = s.Logout(c.Request.Context(), token)
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
