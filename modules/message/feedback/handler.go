package feedback

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contextx "hei-gin/framework/core/context"
	"hei-gin/framework/core/bind"
	"hei-gin/framework/core/response"
	"hei-gin/framework/core/schema"
	"hei-gin/framework/core/security"
	"hei-gin/framework/middleware"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

func (s *Service) registerRoutes(d *shared.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := api.Group("/v1/admin/message/feedbacks", middleware.RequireAccountType(security.AccountAdmin))
		admin.GET("/page", middleware.RequirePermission(d.Perms, "message:feedback:page", "Feedback page"), s.pageAdmin)
		admin.GET("/detail", middleware.RequirePermission(d.Perms, "message:feedback:detail", "Feedback detail"), s.detail)
		admin.POST("/update", middleware.RequirePermission(d.Perms, "message:feedback:update", "Reply feedback"), s.update)
		admin.POST("/delete", middleware.RequirePermission(d.Perms, "message:feedback:delete", "Delete feedback"), s.delete)
		admin.POST("/submit", s.submit)
		admin.GET("/my-page", s.myPage)
		admin.GET("/my-detail", s.myDetail)

		portal := api.Group("/v1/portal/message/feedbacks", middleware.RequireAccountType(security.AccountPortal))
		portal.POST("/submit", s.submit)
		portal.GET("/my-page", s.myPage)
		portal.GET("/my-detail", s.myDetail)
	}
}

func (s *Service) submit(c *gin.Context) {
	var req CreateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	sess := contextx.Session(c.Request.Context())
	if sess == nil {
		response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}
	meta := SubmitMeta{
		AccountType: string(sess.AccountType), AccountID: sess.AccountID, CreatedBy: sess.AccountID,
	}
	if err := s.Submit(c.Request.Context(), req, meta); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) update(c *gin.Context) {
	var req UpdateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	aid := contextx.AccountID(c.Request.Context())
	if err := s.Update(c.Request.Context(), req, ReplyMeta{RepliedBy: aid, UpdatedBy: aid}); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) delete(c *gin.Context) {
	var req IDsParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Delete(c.Request.Context(), req.IDs); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) detail(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := s.Detail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "feedback not found")
		return
	}
	response.OK(c, row)
}

func (s *Service) pageAdmin(c *gin.Context) {
	var q PageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, current, size, err := s.PageAdmin(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (s *Service) myPage(c *gin.Context) {
	var q schema.PageQuery
	_ = c.ShouldBindQuery(&q)
	sess := contextx.Session(c.Request.Context())
	rows, total, current, size, err := s.MyPage(c.Request.Context(), q, sess.AccountID, string(sess.AccountType))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (s *Service) myDetail(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	sess := contextx.Session(c.Request.Context())
	row, err := s.MyDetail(c.Request.Context(), q.ID, sess.AccountID, string(sess.AccountType))
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "feedback not found")
		return
	}
	response.OK(c, row)
}
