package codegen

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"
)

// Preview 按方案渲染预览文件。
func (s *Service) Preview(ctx context.Context, id string) (*PreviewResult, error) {
	plan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("codegen plan not found")
	}
	mainFields, err := s.repo.ListFields(ctx, plan.ID, "MAIN")
	if err != nil {
		return nil, err
	}
	if len(mainFields) == 0 {
		if err := s.syncReflectedFields(ctx, plan); err != nil {
			return nil, err
		}
		mainFields, err = s.repo.ListFields(ctx, plan.ID, "MAIN")
		if err != nil {
			return nil, err
		}
	}
	var subFields []Field
	if isRelationType(plan.GenType) && plan.SubTable != nil && *plan.SubTable != "" {
		subFields, err = s.repo.ListFields(ctx, plan.ID, "SUB")
		if err != nil {
			return nil, err
		}
		if len(subFields) == 0 {
			if err := s.syncReflectedFields(ctx, plan); err != nil {
				return nil, err
			}
			subFields, err = s.repo.ListFields(ctx, plan.ID, "SUB")
			if err != nil {
				return nil, err
			}
		}
	}
	files, err := renderPlan(plan, mainFields, subFields)
	if err != nil {
		return nil, err
	}
	return &PreviewResult{Files: files}, nil
}

// DownloadZip 按方案打包生成代码为 zip。
func (s *Service) DownloadZip(ctx context.Context, id string) ([]byte, string, error) {
	preview, err := s.Preview(ctx, id)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, f := range preview.Files {
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
	return buf.Bytes(), "codegen-" + id + ".zip", nil
}

// renderPlan 渲染全部生成文件（Go 后端 + Vue 前端 + SQL）。
func renderPlan(plan *Plan, mainFields, subFields []Field) ([]PreviewFileResult, error) {
	ctx := buildEmitContext(plan, mainFields, subFields)
	files := []struct {
		path     string
		language string
		tmpl     string
	}{}
	modRoot := "modules/" + firstSegment(ctx.ModulePath)

	// Go 后端（同包分文件：model/param/result/repo/service/handler/register）
	backend := []struct {
		file string
		tmpl string
	}{
		{ctx.Main.Package + "/model.go", goModelTmpl},
		{ctx.Main.Package + "/param.go", goParamTmpl},
		{ctx.Main.Package + "/result.go", goResultTmpl},
		{ctx.Main.Package + "/repo.go", goRepoTmpl},
		{ctx.Main.Package + "/service.go", goServiceTmpl},
		{ctx.Main.Package + "/handler.go", goHandlerTmpl},
		{ctx.Main.Package + "/register.go", goRegisterTmpl},
	}
	for _, b := range backend {
		files = append(files, struct {
			path     string
			language string
			tmpl     string
		}{modRoot + "/" + b.file, "go", b.tmpl})
	}
	if ctx.HasSub && ctx.Sub != nil {
		sub := []struct {
			file string
			tmpl string
		}{
			{ctx.Sub.Package + "/model.go", goModelTmpl},
			{ctx.Sub.Package + "/param.go", goParamTmpl},
			{ctx.Sub.Package + "/repo.go", goRepoTmpl},
			{ctx.Sub.Package + "/service.go", goServiceTmpl},
			{ctx.Sub.Package + "/handler.go", goHandlerTmpl},
		}
		for _, b := range sub {
			files = append(files, struct {
				path     string
				language string
				tmpl     string
			}{modRoot + "/" + b.file, "go", b.tmpl})
		}
	}

	// 前端
	files = append(files, struct {
		path     string
		language string
		tmpl     string
	}{ctx.ApiFile, "typescript", apiTsTmpl})
	files = append(files, struct {
		path     string
		language string
		tmpl     string
	}{ctx.ApiIndexAppend, "typescript", apiIndexAppendTmpl})
	files = append(files, struct {
		path     string
		language string
		tmpl     string
	}{ctx.ViewPath, "vue", indexVueTmpl})
	files = append(files, struct {
		path     string
		language string
		tmpl     string
	}{ctx.ViewComponentDir + "/ModalForm.vue", "vue", modalFormTmpl})
	files = append(files, struct {
		path     string
		language string
		tmpl     string
	}{ctx.ViewComponentDir + "/ModalDetail.vue", "vue", modalDetailTmpl})
	if ctx.HasSub {
		files = append(files, struct {
			path     string
			language string
			tmpl     string
		}{ctx.ViewComponentDir + "/children/ChildModalForm.vue", "vue", childModalFormTmpl})
		files = append(files, struct {
			path     string
			language string
			tmpl     string
		}{ctx.ViewComponentDir + "/children/ChildModalDetail.vue", "vue", childModalDetailTmpl})
	}

	// 菜单权限 SQL
	files = append(files, struct {
		path     string
		language string
		tmpl     string
	}{"scripts/" + toSnake(ctx.Main.EntityName) + "_menu_permission.sql", "sql", menuPermissionSqlTmpl})

	out := make([]PreviewFileResult, 0, len(files))
	funcs := template.FuncMap{
		"sq":           func(s string) string { return strings.ReplaceAll(s, "'", "''") },
		"replaceColon": func(s string) string { return strings.ReplaceAll(s, ":", "_") },
		"toSnake":      toSnake,
	}
	for _, f := range files {
		var b bytes.Buffer
		t, err := template.New(f.path).Funcs(funcs).Parse(f.tmpl)
		if err != nil {
			return nil, fmt.Errorf("parse template %s: %w", f.path, err)
		}
		if err := t.Execute(&b, ctx); err != nil {
			return nil, fmt.Errorf("render template %s: %w", f.path, err)
		}
		content := b.String()
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		out = append(out, PreviewFileResult{Path: f.path, Language: f.language, Content: content})
	}
	return out, nil
}

// emitCtx 模板上下文。
//
// Author: Charlie
type emitCtx struct {
	Author            string
	GeneratedAt       string
	GenType           string
	HasTree           bool
	HasSub            bool
	HasTreeParentForm bool
	PermissionPrefix  string
	APIPrefix         string
	ModulePath        string
	ModuleRoot        string
	BasePackage       string
	ApiFile           string
	ApiIndexAppend    string
	ApiExportName     string
	ViewPath          string
	ViewComponentDir  string
	Main              *entityCtx
	Sub               *entityCtx
	Plan              *Plan
	Menu              menuCtx
}

// entityCtx 主/子表上下文。
type entityCtx struct {
	Package      string
	EntityName   string
	VarName      string
	TableName    string
	PKName       string
	BusinessName string
	Fields       []fieldCtx
	FormFields   []fieldCtx
	QueryFields  []fieldCtx
	TableFields  []fieldCtx
	DetailFields []fieldCtx
	HasTime      bool
	HasJSON      bool
	HasBool      bool
	HasInt       bool
	HasFloat     bool
	HasDatetime  bool
	HasDict      bool
	IsSub        bool
}

// fieldCtx 单字段上下文。
type fieldCtx struct {
	Name           string // snake column
	Property       string // camel
	GoName         string // Pascal
	Label          string
	Comment        string
	DBType         string
	PythonType     string
	TypescriptType string
	GoType         string // *int64 等
	GoTypeBase     string
	IsPointer      bool
	FormWidget     string
	DictCode       string
	QueryOperator  string
	ShowInTable    bool
	ShowInForm     bool
	ShowInDetail   bool
	ShowInQuery    bool
	IsPrimaryKey   bool
	IsRequired     bool
	IsNullable     bool
	MaxLength      int
	Sort           int
	IsDatetime     bool
	IsJSON         bool
	IsBool         bool
	VueDefault     string
}

// menuCtx 菜单 SQL 上下文。
type menuCtx struct {
	MenuID  string
	Actions []menuAction
}

type menuAction struct {
	Key        string
	Label      string
	Sort       int
	ResourceID string
	RelationID string
}

// buildEmitContext 构建模板上下文。
func buildEmitContext(plan *Plan, mainFields, subFields []Field) *emitCtx {
	genType := plan.GenType
	if genType == "" {
		genType = "TABLE"
	}
	hasTree := genType == "TREE" || genType == "LEFT_TREE_TABLE"
	hasSub := isRelationType(genType) &&
		plan.SubEntityName != nil && *plan.SubEntityName != "" &&
		plan.SubTable != nil && *plan.SubTable != "" &&
		plan.SubPK != nil && *plan.SubPK != ""

	apiPrefix := plan.APIPrefix
	apiPrefix = strings.TrimPrefix(apiPrefix, "/api/v1/admin")
	if apiPrefix == "" || apiPrefix == "/" {
		apiPrefix = "/biz/" + toKebab(ctxVar(plan.MainBusinessName))
	}
	if !strings.HasPrefix(apiPrefix, "/") {
		apiPrefix = "/" + apiPrefix
	}

	main := buildEntityCtx(plan.MainEntityName, plan.MainTable, plan.MainPK, plan.MainBusinessName, mainFields, false, hasTree, plan.TreeParentField)
	var sub *entityCtx
	if hasSub && plan.SubEntityName != nil {
		sub = buildEntityCtx(*plan.SubEntityName, *plan.SubTable, *plan.SubPK, *plan.SubBusinessName, subFields, true, false, nil)
	}

	hasTreeParentForm := hasTree && plan.TreeParentField != nil && fieldIn(main.FormFields, *plan.TreeParentField)

	componentPath := strings.TrimSpace(plan.ComponentPath)
	viewPath := "web/admin/src/views/" + strings.TrimPrefix(componentPath, "/")
	viewDir := viewPath
	if idx := strings.LastIndex(viewPath, "/"); idx > 0 {
		viewDir = viewPath[:idx]
	}
	if viewDir == viewPath && strings.HasSuffix(viewPath, ".vue") {
		viewDir = "web/admin/src/views"
	}
	viewComponentDir := viewDir + "/components"
	apiFile := resolveApiFile(plan, componentPath)
	apiExportName := lowerFirst(plan.MainEntityName) + "Api"

	modPath := plan.MainModulePath
	if modPath == "" {
		modPath = "biz/" + ctxVar(plan.MainBusinessName)
	}
	modPath = strings.Trim(modPath, "/")
	moduleRoot := "modules/" + firstSegment(modPath)

	return &emitCtx{
		Author:            plan.Author,
		GeneratedAt:       time.Now().Format("2006-01-02 15:04:05"),
		GenType:           genType,
		HasTree:           hasTree,
		HasSub:            hasSub,
		HasTreeParentForm: hasTreeParentForm,
		PermissionPrefix:  plan.PermissionPrefix,
		APIPrefix:         apiPrefix,
		ModulePath:        modPath,
		ModuleRoot:        moduleRoot,
		BasePackage:       "hei-gin/" + moduleRoot + "/" + lastSegment(modPath),
		ApiFile:           apiFile,
		ApiIndexAppend:    "web/admin/src/api/index.ts.append",
		ApiExportName:     apiExportName,
		ViewPath:          viewPath,
		ViewComponentDir:  viewComponentDir,
		Main:              main,
		Sub:               sub,
		Plan:              plan,
		Menu:              buildMenuCtx(hasTree),
	}
}

func buildEntityCtx(entityName, tableName, pkName, businessName string, fields []Field, isSub, hasTree bool, treeParent *string) *entityCtx {
	ctx := &entityCtx{
		Package:      toSnake(entityName),
		EntityName:   entityName,
		VarName:      lowerFirst(entityName),
		TableName:    tableName,
		PKName:       pkName,
		BusinessName: businessName,
		IsSub:        isSub,
	}
	if pkName == "" {
		ctx.PKName = "id"
	}
	exclude := map[string]bool{}
	if hasTree && treeParent != nil && *treeParent != "" {
		exclude[*treeParent] = true
	}
	for _, f := range fields {
		fc := buildFieldCtx(f)
		if isAuditColumn(f.ColumnName) {
			continue
		}
		ctx.Fields = append(ctx.Fields, fc)
		if fc.ShowInForm && !fc.IsPrimaryKey {
			ctx.FormFields = append(ctx.FormFields, fc)
		}
		if fc.ShowInQuery && !fc.IsPrimaryKey {
			ctx.QueryFields = append(ctx.QueryFields, fc)
		}
		if fc.ShowInTable && !exclude[f.ColumnName] {
			ctx.TableFields = append(ctx.TableFields, fc)
		}
		if fc.ShowInDetail {
			ctx.DetailFields = append(ctx.DetailFields, fc)
		}
		if fc.GoTypeBase == "time.Time" {
			ctx.HasTime = true
		}
		if fc.GoTypeBase == "datatypes.JSON" {
			ctx.HasJSON = true
		}
		if fc.IsBool {
			ctx.HasBool = true
		}
		if fc.PythonType == "int" {
			ctx.HasInt = true
		}
		if fc.PythonType == "float" {
			ctx.HasFloat = true
		}
		if fc.IsDatetime {
			ctx.HasDatetime = true
		}
		if fc.DictCode != "" {
			ctx.HasDict = true
		}
	}
	return ctx
}

func buildFieldCtx(f Field) fieldCtx {
	py := def(f.PythonType, "str")
	goType, base, ptr := semanticToGoType(py, f.IsNullable)
	isDatetime := py == "datetime" || f.FormWidget == "datetime"
	isJSON := py == "dict" || strings.Contains(strings.ToLower(f.DBType), "json")
	isBool := py == "bool"
	label := ""
	if f.ColumnComment != nil {
		label = *f.ColumnComment
	}
	if label == "" {
		label = f.ColumnName
	}
	return fieldCtx{
		Name:           f.ColumnName,
		Property:       snakeToCamel(f.ColumnName),
		GoName:         toPascal(f.ColumnName),
		Label:          label,
		Comment:        label,
		DBType:         f.DBType,
		PythonType:     py,
		TypescriptType: def(f.TypescriptType, "string"),
		GoType:         goType,
		GoTypeBase:     base,
		IsPointer:      ptr,
		FormWidget:     def(f.FormWidget, "input"),
		DictCode:       deref(f.DictCode),
		QueryOperator:  deref(f.QueryOperator),
		ShowInTable:    f.ShowInTable,
		ShowInForm:     f.ShowInForm,
		ShowInDetail:   f.ShowInDetail,
		ShowInQuery:    f.ShowInQuery,
		IsPrimaryKey:   f.IsPrimaryKey,
		IsRequired:     f.IsRequired,
		IsNullable:     f.IsNullable,
		MaxLength:      derefInt(f.MaxLength),
		Sort:           f.Sort,
		IsDatetime:     isDatetime,
		IsJSON:         isJSON,
		IsBool:         isBool,
		VueDefault:     vueDefault(py, isDatetime, isJSON, isBool),
	}
}

func buildMenuCtx(hasTree bool) menuCtx {
	actions := []menuAction{
		{Key: "page", Label: "分页", Sort: 10},
		{Key: "create", Label: "新增", Sort: 20},
		{Key: "detail", Label: "详情", Sort: 30},
		{Key: "update", Label: "编辑", Sort: 40},
		{Key: "delete", Label: "删除", Sort: 50},
	}
	if hasTree {
		actions = append(actions, menuAction{Key: "list", Label: "树列表", Sort: 90})
	}
	for i := range actions {
		actions[i].ResourceID = snowflakeLikeID()
		actions[i].RelationID = snowflakeLikeID()
	}
	return menuCtx{MenuID: snowflakeLikeID(), Actions: actions}
}

func semanticToGoType(py string, nullable bool) (goType, base string, ptr bool) {
	switch py {
	case "int":
		base = "int64"
	case "float":
		base = "float64"
	case "bool":
		base = "bool"
	case "datetime":
		base = "time.Time"
	case "dict":
		base = "datatypes.JSON"
	default:
		base = "string"
	}
	if nullable && base != "datatypes.JSON" && base != "time.Time" {
		return "*" + base, base, true
	}
	if nullable && base == "time.Time" {
		return "*time.Time", "time.Time", true
	}
	return base, base, false
}

func vueDefault(py string, isDatetime, isJSON, isBool bool) string {
	if isDatetime {
		return "null"
	}
	if py == "int" || py == "float" {
		return "0"
	}
	if isBool {
		return "false"
	}
	if isJSON {
		return "'{}'"
	}
	return "''"
}

// ---- 命名工具 ----

func snakeToCamel(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == '_' || r == '-' || r == '.' })
	var b strings.Builder
	for i, p := range parts {
		if p == "" {
			continue
		}
		if i == 0 {
			b.WriteString(strings.ToLower(p))
		} else {
			b.WriteString(strings.ToUpper(p[:1]) + strings.ToLower(p[1:]))
		}
	}
	return b.String()
}

func toPascal(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == '_' || r == '-' || r == '.' || r == ' ' })
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

func toSnake(s string) string {
	if s == "" {
		return "entity"
	}
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1])) {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
		} else if r == '-' || r == ' ' {
			b.WriteByte('_')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func toKebab(s string) string {
	return strings.ReplaceAll(toSnake(s), "_", "-")
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func firstSegment(p string) string {
	p = strings.Trim(p, "/")
	for _, part := range strings.Split(p, "/") {
		if part != "" && part != "." {
			return strings.ToLower(part)
		}
	}
	return "biz"
}

func lastSegment(p string) string {
	p = strings.Trim(p, "/")
	parts := strings.Split(p, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && parts[i] != "." {
			return strings.ToLower(parts[i])
		}
	}
	return "entity"
}

func ctxVar(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "entity"
	}
	return v
}

func fieldIn(fields []fieldCtx, name string) bool {
	for _, f := range fields {
		if f.Name == name {
			return true
		}
	}
	return false
}

func resolveApiFile(plan *Plan, componentPath string) string {
	cleaned := strings.TrimPrefix(strings.TrimSpace(componentPath), "/")
	parts := strings.Split(cleaned, "/")
	if len(parts) >= 2 && strings.HasSuffix(parts[len(parts)-1], ".vue") {
		var b strings.Builder
		b.WriteString("web/admin/src/api")
		for i := 0; i < len(parts)-1; i++ {
			b.WriteString("/")
			b.WriteString(parts[i])
		}
		return b.String() + ".ts"
	}
	return "web/admin/src/api/" + toSnake(plan.MainEntityName) + ".ts"
}

func snowflakeLikeID() string {
	ts := time.Now().UnixMilli()
	return strconv.FormatInt(ts*10000+int64(time.Now().Nanosecond()%10000), 10)
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}
