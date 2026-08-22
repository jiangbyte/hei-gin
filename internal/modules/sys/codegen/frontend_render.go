// internal/modules/sys/codegen/frontend_render.go boot 对齐前端模板渲染。
//
// Author: Charlie

package codegen

import (
	"bytes"
	"embed"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

//go:embed templates/frontend/*.tmpl
var bootFrontendTemplates embed.FS

type bootFrontendSpec struct {
	path     string
	tmpl     string
	language string
	child    bool
}

// bootFrontendFiles 与 hei-boot 一致的前端产物。
var bootFrontendFiles = []bootFrontendSpec{
	{"api.ts", "api.ts.tmpl", "typescript", false},
	{"api_index.ts.append", "api_index_export.ts.tmpl", "typescript", false},
	{"index.vue", "index.vue.tmpl", "vue", false},
	{"ModalForm.vue", "ModalForm.vue.tmpl", "vue", false},
	{"ModalDetail.vue", "ModalDetail.vue.tmpl", "vue", false},
	{"ChildModalForm.vue", "ChildModalForm.vue.tmpl", "vue", true},
	{"ChildModalDetail.vue", "ChildModalDetail.vue.tmpl", "vue", true},
}

func renderBootFrontend(plan *Plan, mainFields, subFields []Field) ([]PreviewFileResult, error) {
	base := buildFrontendContext(plan, mainFields, subFields)
	funcs := bootTemplateFuncs()
	out := make([]PreviewFileResult, 0, len(bootFrontendFiles))
	for _, spec := range bootFrontendFiles {
		ctx := base
		if spec.child {
			if !base.HasSub || base.Sub == nil {
				continue
			}
			ctx = withFrontendTarget(base, true)
		}
		content, err := renderBootTemplate(spec.tmpl, ctx, funcs)
		if err != nil {
			return nil, err
		}
		out = append(out, PreviewFileResult{
			Path:     bootFrontendPath(plan, spec.path),
			Language: spec.language,
			Content:  content,
		})
	}
	return out, nil
}

func renderBootTemplate(name string, ctx *frontendCtx, funcs template.FuncMap) (string, error) {
	raw, err := bootFrontendTemplates.ReadFile("templates/frontend/" + name)
	if err != nil {
		return "", err
	}
	t, err := template.New(name).Funcs(funcs).Parse(string(raw))
	if err != nil {
		return "", fmt.Errorf("parse frontend template %s: %w", name, err)
	}
	var b bytes.Buffer
	if err := t.Execute(&b, ctx); err != nil {
		return "", fmt.Errorf("render frontend template %s: %w", name, err)
	}
	content := b.String()
	content = normalizeFrontendOutput(content)
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content, nil
}

func bootTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"deref":       bootDeref,
		"or":          bootOr,
		"and":         bootAnd,
		"ne":          bootNe,
		"eq":          bootEq,
		"ge":          bootGe,
		"hasDictCode": bootHasDictCode,
		"wireImports": bootWireImports,
		"bootThen":    bootThen,
		"bootIsFirst": bootIsFirst,
		"bootIsLast":  bootIsLast,
	}
}

func bootFrontendPath(plan *Plan, name string) string {
	componentPath := strings.TrimSpace(plan.ComponentPath)
	switch name {
	case "api.ts":
		return resolveApiFile(plan, componentPath)
	case "api_index.ts.append":
		return frontendRoot + "/src/api/index.ts.append"
	case "index.vue":
		return frontendRoot + "/src/views/" + strings.TrimPrefix(componentPath, "/")
	default:
		viewDir := frontendRoot + "/src/views/" + strings.TrimPrefix(componentPath, "/")
		if idx := strings.LastIndex(viewDir, "/"); idx > 0 {
			viewDir = viewDir[:idx]
		}
		if strings.HasPrefix(name, "Child") {
			return viewDir + "/components/children/" + name
		}
		return viewDir + "/components/" + name
	}
}

func bootHasDictCode(field any) bool {
	m, ok := field.(map[string]any)
	if !ok {
		return false
	}
	dc, _ := m["dict_code"].(string)
	return strings.TrimSpace(dc) != ""
}

func bootWireImports(target any) string {
	m, ok := target.(map[string]any)
	if !ok {
		return ""
	}
	fields, _ := m["form_fields"].([]map[string]any)
	var parts []string
	if bootAny(fields, "is_bool") {
		parts = append(parts, "wireBool")
	}
	if bootAnyValueType(fields, "int") {
		parts = append(parts, "wireInt")
	}
	if bootAnyValueType(fields, "float") {
		parts = append(parts, "wireFloat")
	}
	return strings.Join(parts, ", ")
}

func bootThen(val any, def string) string {
	if bootTruthy(val) {
		return fmt.Sprint(val)
	}
	return def
}

func bootIsFirst(idx int) bool { return idx == 0 }

func bootIsLast(idx int, list any) bool {
	switch v := list.(type) {
	case []map[string]any:
		return len(v) > 0 && idx == len(v)-1
	default:
		return false
	}
}

func bootDeref(v any) any {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return nil
		}
		return rv.Elem().Interface()
	case reflect.Map:
		return v
	default:
		return v
	}
}

func bootOr(args ...any) bool {
	for _, a := range args {
		if bootTruthy(a) {
			return true
		}
	}
	return false
}
func bootAnd(args ...any) bool {
	for _, a := range args {
		if !bootTruthy(a) {
			return false
		}
	}
	return true
}

func bootTruthy(v any) bool {
	if v == nil {
		return false
	}
	switch x := v.(type) {
	case bool:
		return x
	case string:
		return strings.TrimSpace(x) != ""
	case map[string]any:
		return len(x) > 0
	default:
		return true
	}
}

func bootNe(a any, b string) bool {
	return fmt.Sprint(bootDeref(a)) != b
}
func bootEq(a, b any) bool {
	return fmt.Sprint(bootDeref(a)) == fmt.Sprint(bootDeref(b))
}
func bootGe(a, b int) bool { return a >= b }

func normalizeFrontendOutput(content string) string {
	repl := []struct{ old, new string }{
		{"import {displayValue", "import { displayValue"},
		{"width:680px", "width: 680px"},
		{"width:720px", "width: 720px"},
		{"width:960px", "width: 960px"},
	}
	for _, r := range repl {
		content = strings.ReplaceAll(content, r.old, r.new)
	}
	for strings.Contains(content, "\n\n\n") {
		content = strings.ReplaceAll(content, "\n\n\n", "\n\n")
	}
	return strings.TrimRight(content, "\n") + "\n"
}
