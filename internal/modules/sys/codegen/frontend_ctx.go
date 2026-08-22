// internal/modules/sys/codegen/frontend_ctx.go boot 兼容前端模板上下文。
//
// Author: Charlie

package codegen

import (
	"strings"
	"time"
)

// frontendCtx boot CodegenTemplateEngine.buildContext 兼容结构（map 供 text/template 使用）。
type frontendCtx struct {
	Author            string
	GeneratedAt       string
	GenType           string
	HasTree           bool
	HasSub            bool
	HasTreeParentForm bool
	PermissionPrefix  string
	APIExportName     string
	APIExport         string
	Plan              map[string]any
	Main              map[string]any
	Sub               map[string]any
	Target            map[string]any
}

func buildFrontendContext(plan *Plan, mainFields, subFields []Field) *frontendCtx {
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
		apiPrefix = "/biz/" + toKebab(ctxVar(plan.BusinessName))
	}
	if !strings.HasPrefix(apiPrefix, "/") {
		apiPrefix = "/" + apiPrefix
	}

	mainExclude := map[string]bool{}
	if genType == "TREE" && plan.TreeParentField != nil && *plan.TreeParentField != "" {
		mainExclude[*plan.TreeParentField] = true
	}
	main := buildBootEntityContext(plan.EntityName, plan.Table, plan.PKColumn, mainFields, mainExclude)

	var sub map[string]any
	if hasSub && plan.SubEntityName != nil {
		sub = buildBootEntityContext(*plan.SubEntityName, *plan.SubTable, *plan.SubPK, subFields, map[string]bool{})
	}

	hasTreeParentForm := hasTree && plan.TreeParentField != nil &&
		bootFieldIn(main, "form_fields", *plan.TreeParentField)

	apiFile := strings.TrimPrefix(resolveApiFile(plan, plan.ComponentPath), frontendRoot+"/")
	apiExportName := lowerFirst(plan.EntityName) + "Api"
	apiRel := strings.TrimPrefix(apiFile, "src/api/")
	apiRel = strings.TrimSuffix(apiRel, ".ts")
	apiExport := "export * as " + apiExportName + " from './" + apiRel + "'"

	planMap := map[string]any{
		"author":              plan.Author,
		"gen_type":              genType,
		"description":           deref(plan.Description),
		"table_name":            plan.Table,
		"pk_column":             plan.PKColumn,
		"entity_name":           plan.EntityName,
		"module_path":           plan.ModulePath,
		"business_name":         plan.BusinessName,
		"api_prefix":            apiPrefix,
		"permission_prefix":     plan.PermissionPrefix,
		"resource_module_id":    deref(plan.ResourceModuleID),
		"parent_resource_id":    deref(plan.ParentResourceID),
		"menu_name":             plan.MenuName,
		"menu_path":             plan.MenuPath,
		"component_path":        plan.ComponentPath,
		"icon":                  deref(plan.Icon),
		"sort":                  plan.Sort,
		"tree_parent_field":     deref(plan.TreeParentField),
		"tree_label_field":      deref(plan.TreeLabelField),
		"sub_table":             deref(plan.SubTable),
		"sub_pk":                deref(plan.SubPK),
		"sub_foreign_key":       deref(plan.SubForeignKey),
		"sub_entity_name":       deref(plan.SubEntityName),
		"sub_business_name":     deref(plan.SubBusinessName),
		"main_business_name":    "",
	}

	return &frontendCtx{
		Author:            plan.Author,
		GeneratedAt:       time.Now().Format("2006-01-02 15:04:05"),
		GenType:           genType,
		HasTree:           hasTree,
		HasSub:            hasSub,
		HasTreeParentForm: hasTreeParentForm,
		PermissionPrefix:  normalizePermissionPrefix(plan.PermissionPrefix),
		APIExportName:     apiExportName,
		APIExport:         apiExport,
		Plan:              planMap,
		Main:              main,
		Sub:               sub,
		Target:            main,
	}
}

func withFrontendTarget(ctx *frontendCtx, child bool) *frontendCtx {
	copy := *ctx
	if child && ctx.Sub != nil {
		copy.Target = ctx.Sub
	} else {
		copy.Target = ctx.Main
	}
	if child {
		copy.HasTreeParentForm = false
	} else {
		copy.HasTreeParentForm = ctx.HasTreeParentForm
	}
	return &copy
}

func buildBootEntityContext(entityName, tableName, pkName string, fields []Field, tableExclude map[string]bool) map[string]any {
	if pkName == "" {
		pkName = "id"
	}
	var modelFields, formFields, queryFields, tableFields, detailFields []map[string]any
	for _, f := range fields {
		fc := bootFieldContext(f)
		if !isAuditColumn(f.ColumnName) {
			modelFields = append(modelFields, fc)
		}
		if isBootFormField(f) {
			formFields = append(formFields, fc)
		}
		if f.InQuery && !f.PrimaryKey {
			queryFields = append(queryFields, fc)
		}
		if f.InTable && !isAuditColumn(f.ColumnName) && !tableExclude[f.ColumnName] {
			tableFields = append(tableFields, fc)
		}
		if f.InDetail && !isAuditColumn(f.ColumnName) {
			detailFields = append(detailFields, fc)
		}
	}
	entity := map[string]any{
		"entity_name": entityName,
		"var_name":    lowerFirst(entityName),
		"table_name":  tableName,
		"pk_name":     pkName,
		"form_fields": formFields,
		"query_fields": queryFields,
		"table_fields": tableFields,
		"detail_fields": detailFields,
		"has_form_datetime": bootAny(formFields, "is_datetime"),
		"has_form_json":     bootAny(formFields, "is_json"),
		"has_form_bool":     bootAny(formFields, "is_bool"),
		"has_form_richtext": bootAnyWidget(formFields, "richtext"),
		"has_form_markdown": bootAnyWidget(formFields, "markdown"),
		"has_form_code":     bootAnyWidget(formFields, "code"),
		"has_form_icon":     bootAnyWidget(formFields, "icon"),
		"has_form_editor":   bootAnyEditor(formFields),
		"has_form_int":      bootAnyValueType(formFields, "int"),
		"has_form_float":    bootAnyValueType(formFields, "float"),
		"has_detail_json":   bootAny(detailFields, "is_json"),
		"has_detail_richtext": bootAnyWidget(detailFields, "richtext"),
		"has_detail_markdown": bootAnyWidget(detailFields, "markdown"),
		"has_detail_code":     bootAnyWidget(detailFields, "code"),
		"has_detail_icon":     bootAnyWidget(detailFields, "icon"),
		"has_detail_editor":   bootAnyEditor(detailFields),
		"has_detail_dict":     bootAnyDict(detailFields),
		"has_query_dict":      bootAnyDict(queryFields),
		"has_table_dict":      bootAnyDict(tableFields),
		"has_table_bool":      bootAny(tableFields, "is_bool"),
		"has_table_tag":       bootAnyTag(tableFields),
		"needs_submit_normalize": bootNeedsNormalize(formFields),
	}
	_ = modelFields
	return entity
}

func bootFieldContext(f Field) map[string]any {
	py := def(f.ValueType, "str")
	isDatetime := f.Widget == "datetime" || py == "datetime"
	isJSON := py == "dict" || strings.Contains(strings.ToLower(f.DBType), "json")
	isBool := py == "bool"
	label := deref(f.Label)
	if label == "" {
		label = f.ColumnName
	}
	qop := deref(f.QueryOperator)
	if qop == "" {
		qop = "LIKE"
	}
	return map[string]any{
		"name":            f.ColumnName,
		"label":           label,
		"db_type":         f.DBType,
		"value_type":      py,
		"widget":          def(f.Widget, "input"),
		"dict_code":       deref(f.DictCode),
		"code_language":   codeLanguage(f),
		"query_operator":  qop,
		"required":        f.Required,
		"vue_default":     vueDefault(py, isDatetime, isJSON, isBool),
		"is_datetime":     isDatetime,
		"is_json":         isJSON,
		"is_bool":         isBool,
	}
}

func isBootFormField(f Field) bool {
	return f.InForm && !f.PrimaryKey && !isAuditColumn(f.ColumnName)
}

func bootFieldIn(entity map[string]any, listKey, name string) bool {
	raw, _ := entity[listKey].([]map[string]any)
	for _, f := range raw {
		if f["name"] == name {
			return true
		}
	}
	return false
}

func bootAny(fields []map[string]any, key string) bool {
	for _, f := range fields {
		if v, ok := f[key].(bool); ok && v {
			return true
		}
	}
	return false
}

func bootAnyWidget(fields []map[string]any, widget string) bool {
	for _, f := range fields {
		if f["widget"] == widget {
			return true
		}
	}
	return false
}

func bootAnyEditor(fields []map[string]any) bool {
	for _, w := range []string{"richtext", "markdown", "code"} {
		if bootAnyWidget(fields, w) {
			return true
		}
	}
	return false
}

func bootAnyValueType(fields []map[string]any, vt string) bool {
	for _, f := range fields {
		if f["value_type"] == vt {
			return true
		}
	}
	return false
}

func bootAnyDict(fields []map[string]any) bool {
	for _, f := range fields {
		if dc, ok := f["dict_code"].(string); ok && strings.TrimSpace(dc) != "" {
			return true
		}
	}
	return false
}

func bootAnyTag(fields []map[string]any) bool {
	return bootAnyDict(fields) || bootAny(fields, "is_bool")
}

func bootNeedsNormalize(fields []map[string]any) bool {
	return bootAny(fields, "is_datetime") || bootAny(fields, "is_json")
}

func codeLanguage(f Field) string {
	if f.ValueType == "dict" || strings.Contains(strings.ToLower(f.DBType), "json") {
		return "json"
	}
	col := strings.ToLower(f.ColumnName)
	if strings.Contains(col, "sql") {
		return "sql"
	}
	if strings.Contains(col, "script") || strings.Contains(col, "code") || strings.Contains(col, "template") {
		return "typescript"
	}
	return "plaintext"
}
