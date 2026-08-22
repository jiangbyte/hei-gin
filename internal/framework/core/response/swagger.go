// internal/framework/core/response/swagger.go Swagger 响应类型（供 swag 引用）。
//
// Author: Charlie

package response

// SwaggerApiResponse 通用成功响应信封（data 为动态类型，请在路由注解中用 {data=xxx} 指定）。
type SwaggerApiResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"success"`
	Data    any    `json:"data" swaggertype:"object"`
}

// SwaggerApiVoidResponse 写操作成功响应（无 data 字段）。
type SwaggerApiVoidResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"success"`
}

// SwaggerPageData 分页 data 结构。
type SwaggerPageData struct {
	Size    int64 `json:"size" example:"20"`
	Current int64 `json:"current" example:"1"`
	Total   int64 `json:"total" example:"100"`
	Pages   int64 `json:"pages" example:"5"`
	Records any   `json:"records" swaggertype:"array,object"`
}

// SwaggerApiPageResponse 分页成功响应。
type SwaggerApiPageResponse struct {
	Code    int               `json:"code" example:"200"`
	Message string            `json:"message" example:"success"`
	Data    SwaggerPageData   `json:"data"`
}
