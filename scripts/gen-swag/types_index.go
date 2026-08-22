package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type fieldMeta struct {
	Name     string
	JSONName string
	FormName string
	Swagger  string
	Required bool
}

type structIndex map[string][]fieldMeta // type name -> fields

type resultIndex map[string][]string // module dir -> result type names

func buildStructIndex(root string) map[string]structIndex {
	out := make(map[string]structIndex)
	_ = filepath.Walk(filepath.Join(root, "internal"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if base != "param.go" && base != "schema.go" {
			return nil
		}
		modDir := filepath.Dir(path)
		if out[modDir] == nil {
			out[modDir] = make(structIndex)
		}
		mergeStructIndex(out[modDir], path)
		return nil
	})
	// schema package
	schemaPath := filepath.Join(root, "internal", "framework", "core", "schema", "schema.go")
	if _, err := os.Stat(schemaPath); err == nil {
		key := "schema"
		if out[key] == nil {
			out[key] = make(structIndex)
		}
		mergeStructIndex(out[key], schemaPath)
	}
	return out
}

func mergeStructIndex(idx structIndex, path string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return
	}
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			idx[ts.Name.Name] = parseStructFields(st)
		}
	}
}

func parseStructFields(st *ast.StructType) []fieldMeta {
	var fields []fieldMeta
	for _, field := range st.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		name := field.Names[0].Name
		fm := fieldMeta{
			Name:    name,
			Swagger: goTypeToSwagger(field.Type),
		}
		if field.Tag != nil {
			tag := field.Tag.Value
			fm.JSONName = structTag(tag, "json")
			fm.FormName = structTag(tag, "form")
			fm.Required = strings.Contains(tag, `binding:"required`) || strings.Contains(tag, `binding:"required,`)
		}
		if fm.JSONName == "" {
			fm.JSONName = name
		}
		if fm.FormName == "" {
			fm.FormName = fm.JSONName
		}
		if fm.JSONName == "-" {
			continue
		}
		fields = append(fields, fm)
	}
	return fields
}

func structTag(tag, key string) string {
	tag = strings.Trim(tag, "`")
	for _, part := range strings.Split(tag, " ") {
		if strings.HasPrefix(part, key+":") {
			v := strings.TrimPrefix(part, key+":")
			v = strings.Trim(v, `"`)
			if idx := strings.Index(v, ","); idx >= 0 {
				v = v[:idx]
			}
			return v
		}
	}
	return ""
}

func goTypeToSwagger(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return identToSwagger(t.Name)
	case *ast.StarExpr:
		return identToSwagger(typeString(t.X))
	case *ast.ArrayType:
		return "array"
	case *ast.MapType:
		return "object"
	case *ast.SelectorExpr:
		return identToSwagger(t.Sel.Name)
	default:
		return "string"
	}
}

func identToSwagger(name string) string {
	switch name {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "integer"
	case "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	default:
		return "string"
	}
}

func typeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		return "string"
	}
}

func buildResultIndex(root string) map[string][]string {
	out := make(map[string][]string)
	_ = filepath.Walk(filepath.Join(root, "internal"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Base(path) != "result.go" {
			return nil
		}
		modDir := filepath.Dir(path)
		out[modDir] = parseResultTypes(path)
		return nil
	})
	return out
}

func parseResultTypes(path string) []string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil
	}
	var names []string
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, ok := ts.Type.(*ast.StructType); ok {
				names = append(names, ts.Name.Name)
			}
		}
	}
	return names
}

func inferResultType(handler string, types []string, pkgName string) string {
	if len(types) == 0 {
		return ""
	}
	if strings.HasPrefix(handler, "own") && len(handler) > 3 {
		rest := handler[3:]
		cand := "Own" + strings.ToUpper(rest[:1]) + rest[1:] + "Result"
		for _, t := range types {
			if t == cand {
				return t
			}
		}
	}
	if handler == "tree" || handler == "portalTree" {
		for _, t := range types {
			if t == "TreeNode" {
				return t
			}
		}
	}
	if isDetailHandler(handler) {
		return inferDetailRecordType(handler, pkgName, types)
	}
	if isPageHandler(handler) {
		return inferPageRecordType(handler, pkgName, types)
	}
	// auth helpers
	switch handler {
	case "captcha":
		for _, t := range types {
			if strings.Contains(t, "Captcha") {
				return t
			}
		}
	case "login", "refresh":
		for _, t := range types {
			if strings.Contains(t, "Login") || strings.Contains(t, "Token") {
				return t
			}
		}
	case "passwordKey":
		for _, t := range types {
			if strings.Contains(t, "Password") || strings.Contains(t, "Key") {
				return t
			}
		}
	}
	return ""
}

func isDetailHandler(handler string) bool {
	return handler == "detail" || strings.HasPrefix(handler, "detail")
}

func isPageHandler(handler string) bool {
	return handler == "page" || strings.HasPrefix(handler, "page") ||
		strings.HasSuffix(handler, "Page") || handler == "reviewPage" || handler == "identityPage"
}

func inferPageRecordType(handler, pkgName string, types []string) string {
	switch handler {
	case "pageItem", "pageDoc":
		for _, cand := range []string{"ItemDetailResult", "DocDetailResult"} {
			for _, t := range types {
				if t == cand {
					return t
				}
			}
		}
	}
	return inferDetailRecordType(handler, pkgName, types)
}

func inferDetailRecordType(handler, pkgName string, types []string) string {
	switch handler {
	case "detailItem":
		for _, t := range types {
			if t == "ItemDetailResult" {
				return t
			}
		}
	case "detailDoc":
		for _, t := range types {
			if t == "DocDetailResult" {
				return t
			}
		}
	}
	for _, t := range types {
		if t == "DetailResult" {
			return t
		}
	}
	if pkgName != "" {
		cand := pkgTitle(pkgName) + "Result"
		for _, t := range types {
			if t == cand {
				return t
			}
		}
		entity := pkgTitle(pkgName)
		for _, t := range types {
			if t == entity {
				return t
			}
		}
	}
	for _, t := range types {
		if strings.HasSuffix(t, "Result") && !strings.HasPrefix(t, "Own") &&
			!strings.HasPrefix(t, "Identity") && !strings.HasPrefix(t, "OAuth") {
			return t
		}
	}
	for _, t := range types {
		if t != "TreeNode" && !strings.HasPrefix(t, "Identity") && !strings.HasPrefix(t, "OAuth") {
			return t
		}
	}
	return ""
}

func pkgTitle(pkgName string) string {
	if pkgName == "" {
		return ""
	}
	return strings.ToUpper(pkgName[:1]) + pkgName[1:]
}
