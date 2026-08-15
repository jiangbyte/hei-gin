// internal/modules/profile/admin.go 管理端用户中心路由与处理器（对齐 hei-boot AdminProfileController）。
//
// Author: Charlie

package profile

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/bind"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
)

// adminRoutes 管理端 /v1/admin/* 用户中心路由。
func (s *Service) adminRoutes(api *gin.RouterGroup) {
	g := api.Group("/v1/admin", middleware.RequireAccountType(security.AccountAdmin))
	g.GET("/me", s.adminMe)
	profile := g.Group("/profile")
	profile.POST("/update", s.adminUpdateProfile)
	profile.POST("/avatar/upload", s.adminUploadAvatar)
	profile.POST("/password/send-code", s.adminSendPasswordCode)
	profile.POST("/password/update", s.adminUpdatePassword)
	profile.POST("/phone/send-code", s.adminSendPhoneCode)
	profile.POST("/phone/update", s.adminUpdatePhone)
	profile.POST("/email/send-code", s.adminSendEmailCode)
	profile.POST("/email/update", s.adminUpdateEmail)
	profile.GET("/org-info", s.adminOrgInfo)
}

func (s *Service) adminMe(c *gin.Context) {
	out, err := s.Me(c.Request.Context(), contextx.Session(c.Request.Context()))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) adminUpdateProfile(c *gin.Context) {
	var req ProfileUpdateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdateProfile(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) adminUploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, "file required")
		return
	}
	const maxSize = 2 << 20
	if file.Size > maxSize {
		response.Fail(c, http.StatusBadRequest, 400, "avatar too large")
		return
	}
	f, err := file.Open()
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	defer f.Close()
	url, err := s.UploadAvatar(c.Request.Context(), contextx.AccountID(c.Request.Context()),
		file.Filename, f, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.OK(c, AvatarResult{Avatar: url})
}

func (s *Service) adminSendPasswordCode(c *gin.Context) {
	if err := s.SendPasswordChangeCode(c.Request.Context(), contextx.AccountID(c.Request.Context())); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) adminUpdatePassword(c *gin.Context) {
	var req PasswordUpdateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	err := s.UpdatePassword(c.Request.Context(), contextx.AccountID(c.Request.Context()), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, 404, "account not found")
			return
		}
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) adminSendPhoneCode(c *gin.Context) {
	var req SendCodeParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.SendBindCode(c.Request.Context(), contextx.AccountID(c.Request.Context()), "PHONE", req.Target); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) adminUpdatePhone(c *gin.Context) {
	var req PhoneUpdateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdatePhone(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) adminSendEmailCode(c *gin.Context) {
	var req SendCodeParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.SendBindCode(c.Request.Context(), contextx.AccountID(c.Request.Context()), "EMAIL", req.Target); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) adminUpdateEmail(c *gin.Context) {
	var req EmailUpdateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdateEmail(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) adminOrgInfo(c *gin.Context) {
	response.OK(c, s.OrgInfo(contextx.Session(c.Request.Context())))
}
