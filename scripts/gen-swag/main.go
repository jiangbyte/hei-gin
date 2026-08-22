// scripts/gen-swag 从 Gin 已注册路由生成 swag 注释桩并执行 swag init。
//
// 用法：go run ./scripts/gen-swag [-config config-8001.yaml]
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	_ "hei-gin/internal/modules/all"
	"hei-gin/internal/app"
)

type route struct {
	Method string
	Path   string
}

func main() {
	cfgPath := flag.String("config", "config-8001.yaml", "config file path")
	flag.Parse()

	cfg := app.LoadOrDie(*cfgPath)
	deps, err := app.OpenInfra(cfg)
	if err != nil {
		fatalf("open infra: %v", err)
	}
	app.AttachRegisteredModules(deps)
	gin.SetMode(gin.ReleaseMode)
	api := app.NewAPI(deps)

	routes := collectRoutes(api.Engine.Routes())
	if len(routes) == 0 {
		fatalf("no api routes collected")
	}

	root, err := findModuleRoot()
	if err != nil {
		fatalf("module root: %v", err)
	}

	out := filepath.Join(root, "internal", "app", "routes_swag_gen.go")
	ann := newAnnotator(root, routes)
	if err := ann.writeStubs(out); err != nil {
		fatalf("write stubs: %v", err)
	}
	fmt.Printf("wrote %d route stubs to %s\n", len(routes), out)

	swag, err := exec.LookPath("swag")
	if err != nil {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			out, _ := exec.Command("go", "env", "GOPATH").Output()
			gopath = strings.TrimSpace(string(out))
		}
		candidate := filepath.Join(gopath, "bin", "swag")
		if runtime.GOOS == "windows" {
			candidate += ".exe"
		}
		if _, statErr := os.Stat(candidate); statErr == nil {
			swag = candidate
			err = nil
		}
	}
	if err != nil {
		fmt.Println("swag CLI not found; run: go install github.com/swaggo/swag/cmd/swag@latest")
		return
	}
	docOut := filepath.Join(root, "internal", "app", "openapi")
	cmd := exec.Command(swag, "init", "-g", "doc.go", "-d", filepath.Join(root, "internal", "app"), "-o", docOut, "--parseDependency", "--parseInternal", "--outputTypes", "go,json")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatalf("swag init: %v", err)
	}
	fmt.Println("swag init ok ->", docOut)
}

func collectRoutes(in []gin.RouteInfo) []route {
	seen := make(map[string]bool)
	out := make([]route, 0, len(in))
	for _, rt := range in {
		method := strings.ToUpper(rt.Method)
		if method == "HEAD" || method == "OPTIONS" {
			continue
		}
		path := rt.Path
		var rel string
		switch {
		case strings.HasPrefix(path, "/api/"):
			if strings.HasPrefix(path, "/api/v1/internal/") {
				continue
			}
			rel = strings.TrimPrefix(path, "/api")
		default:
			continue
		}
		if rel == "" {
			continue
		}
		key := method + " " + rel
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, route{Method: method, Path: rel})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path == out[j].Path {
			return out[i].Method < out[j].Method
		}
		return out[i].Path < out[j].Path
	})
	return out
}

func swaggerTag(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == "v1" {
		if parts[1] == "admin" || parts[1] == "portal" {
			if len(parts) > 2 {
				return parts[1] + "-" + parts[2]
			}
			return parts[1]
		}
	}
	return "api"
}

func paramColonToBrace(path string) string {
	var b strings.Builder
	for i := 0; i < len(path); i++ {
		if path[i] == ':' {
			j := i + 1
			for j < len(path) && path[j] != '/' {
				j++
			}
			b.WriteByte('{')
			b.WriteString(path[i+1 : j])
			b.WriteByte('}')
			i = j - 1
			continue
		}
		b.WriteByte(path[i])
	}
	return b.String()
}

func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
