// Package app 应用装配与 HTTP 服务。
//
//	@title						HEI Gin API
//	@version					1.1.0-beta
//	@description				HEI Gin 后端 API。JSON 字段使用 snake_case；统一响应信封 code/message/data。
//	@termsOfService				https://github.com/jiangbyte/hei-gin
//	@contact.name				HEI
//	@license.name				Apache 2.0
//	@license.url				https://www.apache.org/licenses/LICENSE-2.0.html
//	@BasePath					/api
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				登录后返回的 token，形如 Bearer {token}
//
//go:generate go run ./scripts/gen-swag
package app
