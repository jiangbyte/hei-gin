package codegen

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"unicode"
)

// TableInfo information_schema 表摘要。
//
// Author: Charlie
type TableInfo struct {
	TableName string `json:"table_name" gorm:"column:table_name"`
	TableType string `json:"table_type" gorm:"column:table_type"`
}

// ColumnInfo information_schema 列摘要。
//
// Author: Charlie
type ColumnInfo struct {
	ColumnName    string  `json:"column_name" gorm:"column:column_name"`
	DataType      string  `json:"data_type" gorm:"column:data_type"`
	IsNullable    string  `json:"is_nullable" gorm:"column:is_nullable"`
	ColumnDefault *string `json:"column_default" gorm:"column:column_default"`
	UDTName       string  `json:"udt_name" gorm:"column:udt_name"`
}

// EmitRequest 预览 / 下载入参。
//
// Author: Charlie
type EmitRequest struct {
	TableName    string `json:"table_name" form:"table_name" binding:"required"`
	ModulePath   string `json:"module_path" form:"module_path"`
	EntityName   string `json:"entity_name" form:"entity_name"`
	BusinessName string `json:"business_name" form:"business_name"`
	Author       string `json:"author" form:"author"`
	APIPrefix    string `json:"api_prefix" form:"api_prefix"`
	PermPrefix   string `json:"permission_prefix" form:"permission_prefix"`
}

// GeneratedFile 生成文件。
//
// Author: Charlie
type GeneratedFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// ListTables 列出 public schema 下的基表。
func (s *Service) ListTables(ctx context.Context) ([]TableInfo, error) {
	var rows []TableInfo
	err := s.repo.with(ctx).Raw(`
SELECT table_name, table_type
FROM information_schema.tables
WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
ORDER BY table_name`).Scan(&rows).Error
	return rows, err
}

// Preview 预览生成的 Go 模块桩代码。
func (s *Service) Preview(ctx context.Context, req EmitRequest) ([]GeneratedFile, error) {
	return s.generate(ctx, req)
}

// DownloadZip 打包生成代码为 zip。
func (s *Service) DownloadZip(ctx context.Context, req EmitRequest) ([]byte, string, error) {
	files, err := s.generate(ctx, req)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, f := range files {
		w, err := zw.Create(f.Path)
		if err != nil {
			_ = zw.Close()
			return nil, "", err
		}
		if _, err := w.Write([]byte(f.Content)); err != nil {
			_ = zw.Close()
			return nil, "", err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, "", err
	}
	name := sanitizeIdent(req.TableName) + "_module.zip"
	return buf.Bytes(), name, nil
}

type emitData struct {
	Package      string
	TableName    string
	EntityName   string
	BusinessName string
	Author       string
	APIPrefix    string
	PermPrefix   string
	ModulePath   string
	Fields       []emitField
}

type emitField struct {
	Name       string
	JSONName   string
	ColumnName string
	GoType     string
	GormTag    string
}

func (s *Service) generate(ctx context.Context, req EmitRequest) ([]GeneratedFile, error) {
	table := strings.TrimSpace(req.TableName)
	if table == "" {
		return nil, fmt.Errorf("table_name required")
	}
	cols, err := s.listColumns(ctx, table)
	if err != nil {
		return nil, err
	}
	pkg := sanitizeIdent(table)
	entity := req.EntityName
	if entity == "" {
		entity = toPascal(table)
	}
	biz := req.BusinessName
	if biz == "" {
		biz = strings.ToLower(entity)
	}
	author := req.Author
	if author == "" {
		author = "Charlie"
	}
	modPath := req.ModulePath
	if modPath == "" {
		modPath = "modules/biz/" + pkg
	}
	api := req.APIPrefix
	if api == "" {
		api = "/v1/admin/biz/" + strings.ReplaceAll(pkg, "_", "-")
	}
	perm := req.PermPrefix
	if perm == "" {
		perm = "biz:" + strings.ReplaceAll(pkg, "_", "") + ":"
	}
	data := emitData{
		Package: pkg, TableName: table, EntityName: entity, BusinessName: biz,
		Author: author, APIPrefix: api, PermPrefix: perm, ModulePath: modPath,
		Fields: mapColumns(cols),
	}
	files := []struct {
		name string
		tmpl string
	}{
		{"model.go", modelTmpl},
		{"service.go", serviceTmpl},
		{"handler.go", handlerTmpl},
		{"register.go", registerTmpl},
	}
	out := make([]GeneratedFile, 0, len(files))
	for _, f := range files {
		t, err := template.New(f.name).Parse(f.tmpl)
		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		if err := t.Execute(&b, data); err != nil {
			return nil, err
		}
		out = append(out, GeneratedFile{Path: modPath + "/" + f.name, Content: b.String()})
	}
	return out, nil
}

func (s *Service) listColumns(ctx context.Context, table string) ([]ColumnInfo, error) {
	var cols []ColumnInfo
	err := s.repo.with(ctx).Raw(`
SELECT column_name, data_type, is_nullable, column_default, udt_name
FROM information_schema.columns
WHERE table_schema = 'public' AND table_name = ?
ORDER BY ordinal_position`, table).Scan(&cols).Error
	return cols, err
}

func mapColumns(cols []ColumnInfo) []emitField {
	out := make([]emitField, 0, len(cols))
	for _, c := range cols {
		gt := sqlToGo(c.UDTName, c.DataType, c.IsNullable == "YES")
		name := toPascal(c.ColumnName)
		out = append(out, emitField{
			Name: name, JSONName: c.ColumnName, ColumnName: c.ColumnName, GoType: gt,
			GormTag: fmt.Sprintf(`gorm:"column:%s" json:"%s"`, c.ColumnName, c.ColumnName),
		})
	}
	return out
}

func sqlToGo(udt, dataType string, nullable bool) string {
	t := strings.ToLower(udt)
	if t == "" {
		t = strings.ToLower(dataType)
	}
	var base string
	switch {
	case strings.Contains(t, "int"):
		base = "int64"
	case strings.Contains(t, "bool"):
		base = "bool"
	case strings.Contains(t, "float") || strings.Contains(t, "numeric") || strings.Contains(t, "double") || strings.Contains(t, "real"):
		base = "float64"
	case strings.Contains(t, "timestamp") || t == "date":
		base = "time.Time"
	case t == "jsonb" || t == "json":
		base = "datatypes.JSON"
	default:
		base = "string"
	}
	if nullable && base != "datatypes.JSON" {
		return "*" + base
	}
	return base
}

func sanitizeIdent(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	out := b.String()
	if out == "" {
		return "entity"
	}
	return out
}

func toPascal(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == '.' || r == ' '
	})
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]) + strings.ToLower(p[1:]))
	}
	out := b.String()
	if out == "" {
		return "Entity"
	}
	return out
}

const modelTmpl = `// Package {{.Package}} 由代码生成预览产生的业务模块桩。
package {{.Package}}

import (
	"time"

	"gorm.io/datatypes"
)

// {{.EntityName}} 对应表 {{.TableName}}。
//
// Author: {{.Author}}
type {{.EntityName}} struct {
{{- range .Fields}}
	{{.Name}} {{.GoType}} ` + "`" + `{{.GormTag}}` + "`" + `
{{- end}}
}

// TableName 返回表名。
func ({{.EntityName}}) TableName() string { return "{{.TableName}}" }
`

const serviceTmpl = `package {{.Package}}

import (
	"context"

	"gorm.io/gorm"

	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service {{.BusinessName}} 业务服务。
//
// Author: {{.Author}}
type Service struct{ db *gorm.DB }

// NewService 构造服务。
func NewService(db *gorm.DB) *Service { return &Service{db: db} }

// New 构建模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "biz.{{.Package}}",
		Models: []any{&{{.EntityName}}{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Detail 按主键查询。
func (s *Service) Detail(ctx context.Context, id string) (*{{.EntityName}}, error) {
	var row {{.EntityName}}
	if err := s.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Create 创建（示例：仅写入 ID）。
func (s *Service) Create(ctx context.Context, row *{{.EntityName}}) error {
	if row.ID == "" {
		row.ID = idgen.Next()
	}
	return s.db.WithContext(ctx).Create(row).Error
}
`

const handlerTmpl = `package {{.Package}}

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hei-gin/framework/core/response"
	"hei-gin/framework/core/schema"
	"hei-gin/framework/core/security"
	"hei-gin/framework/middleware"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

func (s *Service) registerRoutes(d *shared.Deps) module.RouteRegistrar {
	return func(api *gin.RouterGroup) {
		admin := middleware.RequireAccountType(security.AccountAdmin)
		g := api.Group("{{.APIPrefix}}", admin)
		g.GET("/detail", middleware.RequirePermission(d.Perms, "{{.PermPrefix}}detail", "{{.EntityName}} detail"), s.detail)
	}
}

func (s *Service) detail(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := s.Detail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, row)
}
`

const registerTmpl = `package {{.Package}}

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册 biz.{{.Package}} 模块。
func init() {
	module.Register("biz.{{.Package}}", 90, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
`
