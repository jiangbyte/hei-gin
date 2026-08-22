// internal/modules/profile/identity/handler.go HTTP 处理器。
//
// Author: Charlie
package identity

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/bind"
	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/schema"
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/middleware"
	"hei-gin/internal/framework/platform/module"
)

func (s *Service) registerRoutes(d *module.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		adminUser := api.Group("/v1/admin", middleware.RequireAccountType(security.AccountAdmin))
		adminUser.GET("/profile/identity/status", s.adminIdentityStatus)
		adminUser.GET("/real-name/case/options", s.adminOptions)
		adminUser.POST("/real-name/case/submit", middleware.OperationAudit(d.Audit, "real_name_case", "submit"), s.adminSubmit)
		adminUser.POST("/real-name/case/init-third-party", middleware.OperationAudit(d.Audit, "real_name_case", "init_third_party"), s.adminInitThirdParty)
		adminUser.POST("/real-name/case/callback", s.adminCallback)
		adminUser.GET("/real-name/case/my-page", s.adminMyPage)

		portalUser := api.Group("/v1/portal", middleware.RequireAccountType(security.AccountPortal))
		portalUser.GET("/profile/identity/status", s.portalIdentityStatus)
		portalUser.GET("/real-name/case/options", s.portalOptions)
		portalUser.POST("/real-name/case/submit", middleware.OperationAudit(d.Audit, "real_name_case", "submit"), s.portalSubmit)
		portalUser.POST("/real-name/case/init-third-party", middleware.OperationAudit(d.Audit, "real_name_case", "init_third_party"), s.portalInitThirdParty)
		portalUser.POST("/real-name/case/callback", s.portalCallback)
		portalUser.GET("/real-name/case/my-page", s.portalMyPage)

		manage := api.Group("/v1/admin/sys", middleware.RequireAccountType(security.AccountAdmin))
		manage.GET("/real-name-case/review-page", middleware.RequirePermission(d.Perms, "sys:realname:review:verify", "实名审核分页"), s.reviewPage)
		manage.GET("/real-name-case/detail", middleware.RequirePermission(d.Perms, "sys:realname:review:verify", "实名工单详情"), s.detail)
		manage.POST("/real-name-case/approve", middleware.RequirePermission(d.Perms, "sys:realname:review:verify", "实名审核通过"), middleware.OperationAudit(d.Audit, "real_name_case", "approve"), s.approve)
		manage.POST("/real-name-case/reject", middleware.RequirePermission(d.Perms, "sys:realname:review:verify", "实名审核驳回"), middleware.OperationAudit(d.Audit, "real_name_case", "reject"), s.reject)
		manage.GET("/identity/page", middleware.RequirePermission(d.Perms, "sys:realname:identity:revoke", "实名快照分页"), s.identityPage)
		manage.POST("/identity/revoke", middleware.RequirePermission(d.Perms, "sys:realname:identity:revoke", "撤销实名"), middleware.OperationAudit(d.Audit, "profile_identity", "revoke"), s.revoke)
	}
}

func (s *Service) adminIdentityStatus(c *gin.Context)   { s.identityStatus(c) }
func (s *Service) portalIdentityStatus(c *gin.Context)  { s.identityStatus(c) }
func (s *Service) adminOptions(c *gin.Context)          { s.options(c) }
func (s *Service) portalOptions(c *gin.Context)         { s.options(c) }
func (s *Service) adminSubmit(c *gin.Context)           { s.submit(c) }
func (s *Service) portalSubmit(c *gin.Context)          { s.submit(c) }
func (s *Service) adminInitThirdParty(c *gin.Context)   { s.initThirdParty(c) }
func (s *Service) portalInitThirdParty(c *gin.Context)  { s.initThirdParty(c) }
func (s *Service) adminCallback(c *gin.Context)         { s.callback(c) }
func (s *Service) portalCallback(c *gin.Context)        { s.callback(c) }
func (s *Service) adminMyPage(c *gin.Context)           { s.myPage(c) }
func (s *Service) portalMyPage(c *gin.Context)          { s.myPage(c) }

func (s *Service) identityStatus(c *gin.Context) {
	accountID := contextx.AccountID(c.Request.Context())
	out, err := s.GetUserStatus(c.Request.Context(), accountID)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, out)
}

func (s *Service) options(c *gin.Context) {
	response.OK(c, s.Options(c.Request.Context()))
}

func (s *Service) submit(c *gin.Context) {
	var req RealNameCaseSubmitParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	accountID := contextx.AccountID(c.Request.Context())
	if err := s.Submit(c.Request.Context(), accountID, req); err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *Service) initThirdParty(c *gin.Context) {
	var req RealNameCaseInitThirdPartyParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	accountID := contextx.AccountID(c.Request.Context())
	out, err := s.InitThirdParty(c.Request.Context(), accountID, req)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, out)
}

func (s *Service) callback(c *gin.Context) {
	var req RealNameCaseCallbackParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Callback(c.Request.Context(), req); err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *Service) myPage(c *gin.Context) {
	var q RealNameCaseMyPageParam
	_ = c.ShouldBindQuery(&q)
	accountID := contextx.AccountID(c.Request.Context())
	rows, total, current, size, err := s.MyPage(c.Request.Context(), accountID, q)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (s *Service) reviewPage(c *gin.Context) {
	var q RealNameCaseReviewPageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, current, size, err := s.ReviewPage(c.Request.Context(), q)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (s *Service) detail(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := s.Detail(c.Request.Context(), q.ID)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, out)
}

func (s *Service) approve(c *gin.Context) {
	var req RealNameCaseApproveParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	reviewerID := contextx.AccountID(c.Request.Context())
	if err := s.Approve(c.Request.Context(), reviewerID, req); err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *Service) reject(c *gin.Context) {
	var req RealNameCaseRejectParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	reviewerID := contextx.AccountID(c.Request.Context())
	if err := s.Reject(c.Request.Context(), reviewerID, req); err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, nil)
}

func (s *Service) identityPage(c *gin.Context) {
	var q IdentityPageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, current, size, err := s.IdentityPage(c.Request.Context(), q)
	if err != nil {
		failBiz(c, err)
		return
	}
	response.Page(c, int64(current), int64(size), total, rows)
}

func (s *Service) revoke(c *gin.Context) {
	var req IdentityRevokeParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	operatorID := contextx.AccountID(c.Request.Context())
	if err := s.Revoke(c.Request.Context(), operatorID, req); err != nil {
		failBiz(c, err)
		return
	}
	response.OK(c, nil)
}

func failBiz(c *gin.Context, err error) {
	var be *BizError
	if errors.As(err, &be) {
		status := be.HTTPStatus
		if status == 0 {
			status = http.StatusBadRequest
		}
		code := be.Code
		if code == 0 {
			code = 400
		}
		response.Fail(c, status, code, be.Message)
		return
	}
	response.Fail(c, http.StatusBadRequest, 400, err.Error())
}
