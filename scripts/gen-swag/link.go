package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type routeLink struct {
	Method  string
	Path    string
	Handler string
	Summary string
	ModDir  string
	PkgName string
}

type handlerInfo struct {
	BindType   string // e.g. AddParam, schema.IDQuery
	BindLoc    string // body, query, multipart, none
	IsPage     bool
	ResultType string
}

var (
	reGroupAPI    = regexp.MustCompile(`(\w+)\s*:=\s*api\.Group\("([^"]+)"`)
	reGroupNested = regexp.MustCompile(`(\w+)\s*:=\s*(\w+)\.Group\("([^"]+)"`)
	reDirectAPI   = regexp.MustCompile(`api\.(GET|POST|PUT|DELETE|PATCH)\("([^"]+)"`)
	reGroupMethod = regexp.MustCompile(`(\w+)\.(GET|POST|PUT|DELETE|PATCH)\("([^"]+)"`)
	reMountGroup  = regexp.MustCompile(`(\w+)\s*:=\s*r\.Group\("([^"]+)"`)
	reMountMethod = regexp.MustCompile(`(\w+)\.(GET|POST|PUT|DELETE|PATCH)\("([^"]+)"`)
	reHandlerS    = regexp.MustCompile(`s\.(\w+)(?:\([^)]*\))?`)
	reHandlerFn   = regexp.MustCompile(`,\s*(\w+)\s*\)`)
	rePermission  = regexp.MustCompile(`RequirePermission\([^,]+,[^,]+,\s*"([^"]+)"\)`)
	reAuditLabel  = regexp.MustCompile(`OperationAudit\([^,]+,[^,]+,\s*"([^"]+)"\)`)
)

func buildRouteIndex(root string) map[string]routeLink {
	index := make(map[string]routeLink)
	_ = filepath.Walk(filepath.Join(root, "internal"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil || !looksLikeRouteFile(content) {
			return nil
		}
		modDir := filepath.Dir(path)
		pkgName := packageNameFromFile(path)
		var links []routeLink
		if strings.Contains(string(content), "r.Group(") {
			links = parseMountRoutes(string(content), modDir, pkgName)
		} else {
			links = parseRegisterRoutes(normalizeRouteSource(string(content)), modDir, pkgName)
		}
		for _, lk := range links {
			key := lk.Method + " " + lk.Path
			index[key] = lk
		}
		return nil
	})
	return index
}

func looksLikeRouteFile(content []byte) bool {
	s := string(content)
	return strings.Contains(s, ".Group(") ||
		strings.Contains(s, "api.GET(") ||
		strings.Contains(s, "api.POST(") ||
		strings.Contains(s, "api.PUT(") ||
		strings.Contains(s, "api.DELETE(") ||
		strings.Contains(s, ".GET(\"/") ||
		strings.Contains(s, ".POST(\"/")
}

func normalizeRouteSource(content string) string {
	re := regexp.MustCompile(`(\w+)\.(GET|POST|PUT|DELETE|PATCH)\(\s*\n\s*"`)
	return re.ReplaceAllString(content, `${1}.${2}("`)
}

func packageNameFromFile(path string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.PackageClauseOnly)
	if err != nil {
		return ""
	}
	return f.Name.Name
}

func parseRegisterRoutes(content, modDir, pkgName string) []routeLink {
	groups := map[string]string{"api": ""}
	var out []routeLink
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if m := reGroupAPI.FindStringSubmatch(line); m != nil {
			groups[m[1]] = joinPath(groups["api"], m[2])
			continue
		}
		if m := reGroupNested.FindStringSubmatch(line); m != nil {
			parent := groups[m[2]]
			groups[m[1]] = joinPath(parent, m[3])
			continue
		}
		var method, path string
		switch {
		case reDirectAPI.MatchString(line):
			m := reDirectAPI.FindStringSubmatch(line)
			method, path = m[1], joinPath(groups["api"], m[2])
		case reGroupMethod.MatchString(line):
			m := reGroupMethod.FindStringSubmatch(line)
			method, path = m[2], joinPath(groups[m[1]], m[3])
		default:
			continue
		}
		handler := lastHandlerName(line)
		summary := extractSummary(line)
		out = append(out, routeLink{
			Method: strings.ToUpper(method), Path: path, Handler: handler,
			Summary: summary, ModDir: modDir, PkgName: pkgName,
		})
	}
	return out
}

func parseMountRoutes(content, modDir, pkgName string) []routeLink {
	groups := map[string]string{"r": ""}
	var out []routeLink
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if m := reMountGroup.FindStringSubmatch(line); m != nil {
			groups[m[1]] = joinPath("", m[2])
			continue
		}
		if !reMountMethod.MatchString(line) {
			continue
		}
		m := reMountMethod.FindStringSubmatch(line)
		method, path := m[2], joinPath(groups[m[1]], m[3])
		handler := ""
		if hm := reHandlerFn.FindStringSubmatch(line); len(hm) > 1 {
			handler = hm[1]
		}
		out = append(out, routeLink{
			Method: strings.ToUpper(method), Path: path, Handler: handler,
			ModDir: modDir, PkgName: pkgName,
		})
	}
	return out
}

func lastHandlerName(line string) string {
	matches := reHandlerS.FindAllStringSubmatch(line, -1)
	if len(matches) == 0 {
		return ""
	}
	return matches[len(matches)-1][1]
}

func extractSummary(line string) string {
	if m := rePermission.FindStringSubmatch(line); len(m) > 1 {
		return m[1]
	}
	if m := reAuditLabel.FindStringSubmatch(line); len(m) > 1 {
		return m[1]
	}
	return ""
}

func joinPath(prefix, rel string) string {
	if rel == "" {
		return prefix
	}
	if strings.HasPrefix(rel, "/") {
		if prefix == "" {
			return rel
		}
		return strings.TrimSuffix(prefix, "/") + rel
	}
	if prefix == "" {
		return "/" + rel
	}
	return strings.TrimSuffix(prefix, "/") + "/" + rel
}

func buildHandlerIndex(root string) map[string]handlerInfo {
	out := make(map[string]handlerInfo)
	_ = filepath.Walk(filepath.Join(root, "internal"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		base := filepath.Base(path)
		if base == "register.go" || base == "service.go" || base == "repo.go" || base == "model.go" {
			return nil
		}
		modDir := filepath.Dir(path)
		mergeHandlerIndex(out, modDir, path)
		return nil
	})
	return out
}

func handlerKey(modDir, name string) string {
	return modDir + ":" + name
}

func mergeHandlerIndex(out map[string]handlerInfo, modDir, path string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return
	}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || fn.Name == nil {
			continue
		}
		if fn.Name.Name == "registerRoutes" || fn.Name.Name == "RegisterRoutes" || fn.Name.Name == "Mount" ||
			strings.HasPrefix(fn.Name.Name, "Register") && strings.HasSuffix(fn.Name.Name, "Routes") {
			continue
		}
		body := fn.Body
		if body == nil {
			continue
		}
		// closure: func (s *Service) login(...) gin.HandlerFunc { return func(c *gin.Context) { ... } }
		if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
			if inner := findClosureBody(body); inner != nil {
				out[handlerKey(modDir, fn.Name.Name)] = analyzeHandlerBody(inner)
				continue
			}
		}
		out[handlerKey(modDir, fn.Name.Name)] = analyzeHandlerBody(body)
	}
}

func findClosureBody(body *ast.BlockStmt) *ast.BlockStmt {
	for _, stmt := range body.List {
		ret, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(ret.Results) == 0 {
			continue
		}
		fnLit, ok := ret.Results[0].(*ast.FuncLit)
		if !ok {
			continue
		}
		return fnLit.Body
	}
	return nil
}

func checkCall(call *ast.CallExpr, info *handlerInfo) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	pkg := exprName(sel.X)
	fn := sel.Sel.Name
	switch {
	case pkg == "bind" && fn == "JSON":
		info.BindLoc = "body"
	case pkg == "bind" && fn == "Query":
		info.BindLoc = "query"
	case fn == "FormFile":
		info.BindLoc = "multipart"
	case pkg == "response" && fn == "Page":
		info.IsPage = true
	case pkg == "response" && fn == "OK":
		if len(call.Args) >= 2 {
			if arg := call.Args[1]; !isNil(arg) {
				info.ResultType = "_non_nil"
			}
		}
	case fn == "ShouldBindQuery" || fn == "ShouldBindJSON" || fn == "ShouldBindUri":
		if fn == "ShouldBindQuery" {
			info.BindLoc = "query"
		} else if fn == "ShouldBindJSON" {
			info.BindLoc = "body"
		} else {
			info.BindLoc = "path"
		}
	}
}

func analyzeHandlerBody(body *ast.BlockStmt) handlerInfo {
	info := handlerInfo{}
	for _, stmt := range body.List {
		switch x := stmt.(type) {
		case *ast.DeclStmt:
			gen, ok := x.Decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gen.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range vs.Names {
					t := typeExprString(vs.Type)
					if t == "" {
						continue
					}
					n := name.Name
					if strings.HasPrefix(n, "req") || n == "body" || strings.HasPrefix(n, "q") {
						info.BindType = t
					}
				}
			}
		case *ast.AssignStmt:
			for i, lhs := range x.Lhs {
				if _, ok := lhs.(*ast.Ident); !ok || i >= len(x.Rhs) {
					continue
				}
				if call, ok := x.Rhs[i].(*ast.CallExpr); ok {
					if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "FormFile" {
						info.BindLoc = "multipart"
					}
				}
			}
		}
	}
	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			checkCall(call, &info)
		}
		return true
	})
	return info
}

func typeExprString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprName(t.X) + "." + t.Sel.Name
	default:
		return ""
	}
}

func exprName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		return ""
	}
}

func isNil(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

// suffix inference when handler link missing
func inferFromPathSuffix(path, method string) (bindType, bindLoc string) {
	if !isStandardCRUDPath(path) {
		return "", ""
	}
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, "/create"), strings.HasSuffix(lower, "/createitem"):
		return "AddParam", "body"
	case strings.HasSuffix(lower, "/update"), strings.HasSuffix(lower, "/updateitem"):
		return "EditParam", "body"
	case strings.HasSuffix(lower, "/delete"), strings.HasSuffix(lower, "/deleteitem"):
		return "IDsParam", "body"
	case strings.HasSuffix(lower, "/detail"):
		return "schema.IDQuery", "query"
	case strings.HasSuffix(lower, "/page"), strings.HasSuffix(lower, "/my-page"), strings.HasSuffix(lower, "/page-admin"):
		return "PageParam", "query"
	case strings.HasSuffix(lower, "/tree"):
		return "", "query"
	case method == "GET" && strings.Contains(lower, "/own-"):
		return "schema.IDQuery", "query"
	case strings.Contains(lower, "/grant-"):
		return "", "body"
	}
	return "", ""
}

func isStandardCRUDPath(path string) bool {
	return standardCRUDPath.MatchString(path)
}

var standardCRUDPath = regexp.MustCompile(`/v1/(admin|portal)/(sys|biz)/[^/]+(?:/[^/]+)?/(create|update|delete|detail|page)(?:item|doc)?$`)

func inferBindFromHandler(handler string) (bindType, bindLoc string) {
	switch handler {
	case "grantRole":
		return "GrantRoleParam", "body"
	case "grantGroup":
		return "GrantGroupParam", "body"
	case "grantDept":
		return "GrantDeptParam", "body"
	case "grantResource", "grantClientResource":
		return "GrantResourceParam", "body"
	case "ownResource":
		return "OwnResourceQuery", "query"
	case "login", "register", "sendLoginCode", "forgotPassword", "forgotPasswordByPhone",
		"resetPassword", "resetPasswordByPhone", "refresh", "cancel", "registerSendCode":
		return "", "body"
	case "captcha", "passwordKey", "siteFooter", "logout":
		return "", "none"
	default:
		if strings.HasPrefix(handler, "own") {
			return "schema.IDQuery", "query"
		}
		if strings.HasPrefix(handler, "grant") {
			return "", "body"
		}
	}
	return "", ""
}
