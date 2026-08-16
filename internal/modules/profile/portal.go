// internal/modules/profile/portal.go 门户端用户中心路由与处理器（对齐 hei-boot PortalProfileController）。
//
// Author: Charlie

package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/audit"
)

// registerRoutes 门户端 /v1/portal/* 用户中心路由。
func (s *Service) registerRoutes(api *gin.RouterGroup) {
	// 操作审计登记（对齐 hei-boot @OperationAudit：profile_center 门户端）
	s.auditReg.RegisterSpecs(
		audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/portal/profile/update", ResourceType: "profile_center", Action: "update_profile"},
		audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/portal/profile/avatar/upload", ResourceType: "profile_center", Action: "upload_avatar"},
		audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/portal/profile/password/send-code", ResourceType: "profile_center", Action: "send_password_code"},
		audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/portal/profile/password/update", ResourceType: "profile_center", Action: "update_password"},
		audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/portal/profile/phone/send-code", ResourceType: "profile_center", Action: "send_bind_phone_code"},
		audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/portal/profile/phone/update", ResourceType: "profile_center", Action: "update_phone"},
		audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/portal/profile/email/send-code", ResourceType: "profile_center", Action: "send_bind_email_code"},
		audit.AuditSpec{Method: "POST", PathPattern: "/api/v1/portal/profile/email/update", ResourceType: "profile_center", Action: "update_email"},
	)
	g := api.Group("/v1/portal", middleware.RequireAccountType(security.AccountPortal))
	g.GET("/me", s.portalMe)
	profile := g.Group("/profile")
	profile.POST("/update", s.portalUpdateProfile)
	profile.POST("/avatar/upload", s.portalUploadAvatar)
	profile.POST("/password/send-code", s.portalSendPasswordCode)
	profile.POST("/password/update", s.portalUpdatePassword)
	profile.POST("/phone/send-code", s.portalSendPhoneCode)
	profile.POST("/phone/update", s.portalUpdatePhone)
	profile.POST("/email/send-code", s.portalSendEmailCode)
	profile.POST("/email/update", s.portalUpdateEmail)
	g.GET("/spaces/detail", s.portalSpacesDetail)
}

func (s *Service) portalMe(c *gin.Context) {
	out, err := s.Me(c.Request.Context(), contextx.Session(c.Request.Context()))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) portalUpdateProfile(c *gin.Context) {
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

func (s *Service) portalUploadAvatar(c *gin.Context) {
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

func (s *Service) portalSendPasswordCode(c *gin.Context) {
	if err := s.SendPasswordChangeCode(c.Request.Context(), contextx.AccountID(c.Request.Context())); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) portalUpdatePassword(c *gin.Context) {
	var req PasswordUpdateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdatePassword(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) portalSendPhoneCode(c *gin.Context) {
	var req SendCodeParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.SendBindCode(c.Request.Context(), contextx.AccountID(c.Request.Context()), "BIND_PHONE_CODE", "PHONE", req.Target); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) portalUpdatePhone(c *gin.Context) {
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

func (s *Service) portalSendEmailCode(c *gin.Context) {
	var req SendCodeParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.SendBindCode(c.Request.Context(), contextx.AccountID(c.Request.Context()), "BIND_EMAIL_CODE", "EMAIL", req.Target); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) portalUpdateEmail(c *gin.Context) {
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

// portalSpacesDetail 按账号 ID 查询门户公开资料（空间详情）。
func (s *Service) portalSpacesDetail(c *gin.Context) {
	accountID := c.Query("account_id")
	if accountID == "" {
		response.Fail(c, http.StatusBadRequest, 400, "account_id required")
		return
	}
	p, err := s.repo.GetProfile(c.Request.Context(), accountID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "profile not found")
		return
	}
	response.OK(c, gin.H{
		"account_id": accountID,
		"name":       p.Name,
		"nickname":   p.Nickname,
		"avatar":     p.Avatar,
		"signature":  p.Signature,
	})
}
