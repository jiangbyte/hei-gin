package health

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/internal/framework/core/response"
	"hei-gin/internal/framework/core/stringly"
)

func (s *Service) registerRoutes(api *gin.RouterGroup) {
	api.GET("/v1/internal/health/live", s.live)
	api.GET("/v1/internal/health/ready", s.ready)
}

func (s *Service) live(c *gin.Context) {
	response.OK(c, s.Live())
}

func (s *Service) ready(c *gin.Context) {
	out, ok := s.Ready(c.Request.Context())
	if ok {
		response.OK(c, out)
		return
	}
	b, err := stringly.Marshal(response.ApiResponse{Code: 503, Message: "not ready", Data: out})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, "marshal error")
		return
	}
	c.Data(http.StatusServiceUnavailable, "application/json; charset=utf-8", b)
}
