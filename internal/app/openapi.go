// internal/app/openapi.go 提供 /v3/api-docs 与 Swagger UI（swaggo 生成）。
//
// Author: Charlie
package app

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed openapi/swagger.json
var openAPISpec []byte

// MountOpenAPI 挂载 OpenAPI 文档与 Swagger UI。
func MountOpenAPI(r *gin.Engine, baseURL, appName string) {
	spec := patchOpenAPISpec(openAPISpec, baseURL, appName)

	r.GET("/v3/api-docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", spec)
	})
	r.GET("/v3/api-docs/*any", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", spec)
	})
	r.GET("/doc.html", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger-ui/index.html")
	})
	r.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/v3/api-docs")))
}

func patchOpenAPISpec(raw []byte, baseURL, appName string) []byte {
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return raw
	}
	info, ok := doc["info"].(map[string]any)
	if !ok {
		info = map[string]any{}
		doc["info"] = info
	}
	info["title"] = openAPITitle(appName)
	info["description"] = "HEI Gin backend API. JSON fields use snake_case."
	if baseURL != "" {
		if u, err := url.Parse(baseURL); err == nil {
			doc["servers"] = []map[string]string{{
				"url":         strings.TrimRight(baseURL, "/"),
				"description": "Current server",
			}}
			_ = u
		}
	}
	out, err := json.Marshal(doc)
	if err != nil {
		return raw
	}
	return out
}

func openAPITitle(appName string) string {
	switch strings.ToLower(strings.TrimSpace(appName)) {
	case "hei-gin", "":
		return "HEI Gin API"
	default:
		if len(appName) == 0 {
			return "API"
		}
		return strings.ToUpper(appName[:1]) + appName[1:] + " API"
	}
}
