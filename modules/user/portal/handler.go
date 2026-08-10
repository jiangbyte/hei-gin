package portal

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"hei-gin/framework/core/bind"
	"hei-gin/framework/core/response"
	"hei-gin/framework/core/security"
	"hei-gin/framework/middleware"
)

func (s *Service) registerRoutes(api *gin.RouterGroup) {
	g := api.Group("/v1/portal", middleware.RequireAccountType(security.AccountPortal))
	g.GET("/me", s.me)
	uc := g.Group("/user-center")
	uc.POST("/profile/update", s.updateProfile)
	uc.POST("/avatar/upload", s.uploadAvatar)
	uc.POST("/password/send-code", s.sendPasswordCode)
	uc.POST("/password/update", s.updatePassword)
	uc.POST("/phone/update", s.updatePhone)
	uc.POST("/email/update", s.updateEmail)
}

func (s *Service) me(c *gin.Context) {
	out, err := s.Me(c.Request.Context(), SessionFromContext(c.Request.Context()))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) updateProfile(c *gin.Context) {
	var req ProfileUpdateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdateProfile(c.Request.Context(), AccountIDFromContext(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) uploadAvatar(c *gin.Context) {
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
	url, err := s.UploadAvatar(c.Request.Context(), AccountIDFromContext(c.Request.Context()),
		file.Filename, f, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.OK(c, AvatarResult{Avatar: url})
}

func (s *Service) sendPasswordCode(c *gin.Context) {
	response.OK(c, nil)
}

func (s *Service) updatePassword(c *gin.Context) {
	var req PasswordUpdateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	err := s.UpdatePassword(c.Request.Context(), AccountIDFromContext(c.Request.Context()), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, 404, "account not found")
			return
		}
		if err == errOldPassword {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) updatePhone(c *gin.Context) {
	var req PhoneUpdateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdatePhone(c.Request.Context(), AccountIDFromContext(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) updateEmail(c *gin.Context) {
	var req EmailUpdateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.UpdateEmail(c.Request.Context(), AccountIDFromContext(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}
