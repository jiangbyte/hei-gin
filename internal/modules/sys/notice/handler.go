// internal/modules/sys/notice/handler.go HTTP 处理器。
//
// Author: Charlie

package notice

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

func (s *Service) registerRoutes(d *shared.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := api.Group("/v1/admin/sys/notices", middleware.RequireAccountType(security.AccountAdmin))
		admin.POST("/create", middleware.RequirePermission(d.Perms, "sys:notice:create", "创建 notice"), s.create)
		admin.POST("/update", middleware.RequirePermission(d.Perms, "sys:notice:update", "Update notice"), s.update)
		admin.POST("/delete", middleware.RequirePermission(d.Perms, "sys:notice:delete", "Delete notice"), s.delete)
		admin.GET("/detail", middleware.RequirePermission(d.Perms, "sys:notice:detail", "Notice detail"), s.detail)
		admin.GET("/page", middleware.RequirePermission(d.Perms, "sys:notice:page", "Notice page"), s.pageAdmin)
		admin.POST("/publish", middleware.RequirePermission(d.Perms, "sys:notice:publish", "Publish notice"), s.publish)
		admin.POST("/revoke", middleware.RequirePermission(d.Perms, "sys:notice:revoke", "Revoke notice"), s.revoke)
		admin.POST("/pin", middleware.RequirePermission(d.Perms, "sys:notice:pin", "Pin notice"), s.pin)
		s.registerMyRoutes(admin)

		portal := api.Group("/v1/portal/sys/notices")
		portal.GET("/list", s.portalList)
		portalAuth := portal.Group("", middleware.RequireAccountType(security.AccountPortal))
		s.registerMyRoutes(portalAuth)
	}
}

func (s *Service) registerMyRoutes(g *gin.RouterGroup) {
	g.GET("/my-page", s.myPage)
	g.GET("/my-detail", s.myDetail)
	g.GET("/unread-count", s.unreadCount)
	g.POST("/read", s.markRead)
	g.POST("/read-all", s.markAllRead)
}

func (s *Service) create(c *gin.Context) {
	var req CreateParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	aid := contextx.AccountID(c.Request.Context())
	if err := s.Create(c.Request.Context(), req, &aid, &aid); err != nil {
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
	if err := s.Update(c.Request.Context(), req, &aid); err != nil {
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
		response.Fail(c, http.StatusNotFound, 404, "notice not found")
		return
	}
	response.OK(c, row)
}

func (s *Service) pageAdmin(c *gin.Context) {
	var q PageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, current, size, err := s.PageAdmin(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (s *Service) publish(c *gin.Context) {
	var req IDsParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	now := time.Now().UTC()
	aid := contextx.AccountID(c.Request.Context())
	atype := string(contextx.AccountType(c.Request.Context()))
	if err := s.Publish(c.Request.Context(), req.IDs, PublishParam{
		Status: "PUBLISHED", PublishAt: now,
		SenderAccountID: aid, SenderAccountType: atype, UpdatedBy: aid,
	}); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) revoke(c *gin.Context) {
	var req IDsParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Revoke(c.Request.Context(), req.IDs, RevokeParam{
		Status: "REVOKED", RevokedAt: time.Now().UTC(),
	}); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) pin(c *gin.Context) {
	var req PinParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Pin(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) portalList(c *gin.Context) {
	var q PageParam
	_ = c.ShouldBindQuery(&q)
	accountType, accountID := sessionScope(c)
	rows, total, current, size, err := s.PagePublished(c.Request.Context(), q, accountType, accountID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (s *Service) myPage(c *gin.Context) {
	var q PageParam
	_ = c.ShouldBindQuery(&q)
	accountType, accountID := sessionScope(c)
	rows, total, current, size, err := s.PagePublished(c.Request.Context(), q, accountType, accountID)
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
	accountType, accountID := sessionScope(c)
	row, err := s.MyDetail(c.Request.Context(), q.ID, accountType, accountID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "notice not found")
		return
	}
	response.OK(c, row)
}

// sessionScope 从上下文取账户类型与 ID（匿名访问时类型取 PORTAL、ID 为空）。
func sessionScope(c *gin.Context) (string, string) {
	sess := contextx.Session(c.Request.Context())
	if sess == nil {
		return string(security.AccountPortal), ""
	}
	return string(sess.AccountType), sess.AccountID
}

func (s *Service) unreadCount(c *gin.Context) {
	sess := contextx.Session(c.Request.Context())
	if sess == nil {
		response.OK(c, 0)
		return
	}
	total, err := s.UnreadCount(c.Request.Context(), string(sess.AccountType), sess.AccountID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, total)
}

func (s *Service) markRead(c *gin.Context) {
	var req ReadParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	sess := contextx.Session(c.Request.Context())
	if sess == nil {
		response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}
	now := time.Now().UTC()
	for _, id := range req.IDs {
		_ = s.MarkRead(c.Request.Context(), ReadRecord{
			NoticeID: id, AccountType: string(sess.AccountType), AccountID: sess.AccountID, ReadAt: now,
		})
	}
	response.OK(c, nil)
}

func (s *Service) markAllRead(c *gin.Context) {
	sess := contextx.Session(c.Request.Context())
	if sess == nil {
		response.Fail(c, http.StatusUnauthorized, 401, "unauthorized")
		return
	}
	if err := s.MarkAllRead(c.Request.Context(), string(sess.AccountType), sess.AccountID, time.Now().UTC()); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}
