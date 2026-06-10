package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: go run cmd/codegen/main.go <command>

Commands:
  list                    List all discovered plugins
  scaffold <name>         Create a new plugin scaffold (e.g. plugin-xxx)
  gen-imports             Regenerate blank imports in main.go
`)
	}
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	repoRoot := findRepoRoot()

	switch args[0] {
	case "list":
		cmdList(repoRoot)
	case "scaffold":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: go run cmd/codegen/main.go scaffold <name>")
			os.Exit(1)
		}
		cmdScaffold(repoRoot, args[1])
	case "gen-imports":
		cmdGenImports(repoRoot)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func findRepoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		fatal("failed to getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			fatal("go.mod not found")
		}
		dir = parent
	}
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

// ── list ────────────────────────────────────────────────────────────

func cmdList(repoRoot string) {
	pluginsDir := filepath.Join(repoRoot, "plugins")
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		fatal("failed to read plugins dir: %v", err)
	}
	fmt.Println("Discovered plugins:")
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		pluginFile := filepath.Join(pluginsDir, e.Name(), "plugin.go")
		hasPlugin := false
		if _, err := os.Stat(pluginFile); err == nil {
			hasPlugin = true
		}
		status := "[ ]"
		if hasPlugin {
			status = "[✓]"
		}
		fmt.Printf("  %s %s\n", status, e.Name())
	}
}

// ── scaffold ────────────────────────────────────────────────────────

func cmdScaffold(repoRoot, name string) {
	if !strings.HasPrefix(name, "plugin-") {
		name = "plugin-" + name
	}
	pluginsDir := filepath.Join(repoRoot, "plugins")
	target := filepath.Join(pluginsDir, name)
	if _, err := os.Stat(target); err == nil {
		fatal("plugin '%s' already exists", name)
	}

	// Compute naming variants
	pkgName := strings.ReplaceAll(name, "-", "_")
	parts := strings.Split(strings.TrimPrefix(name, "plugin-"), "-")
	pascalName := ""
	for _, p := range parts {
		if len(p) > 0 {
			pascalName += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	if pascalName == "" {
		pascalName = "Sample"
	}
	pluginName := strings.TrimPrefix(name, "plugin-")

	// Create directories
	os.MkdirAll(filepath.Join(target, "api", "v1"), 0755)

	writeTmpl := func(tmpl, path string) {
		r := strings.NewReplacer("@PKG@", pkgName, "@PASCAL@", pascalName, "@NAME@", pluginName, "@FULL@", name)
		result := r.Replace(tmpl)
		if err := os.WriteFile(path, []byte(result), 0644); err != nil {
			fatal("failed to write %s: %v", path, err)
		}
	}

	// ── plugin.go ───────────────────────────────────────────────────
	writeTmpl(pluginGoTmpl, filepath.Join(target, "plugin.go"))
	fmt.Printf("  created %s/plugin.go\n", name)

	// ── model.go ────────────────────────────────────────────────────
	writeTmpl(modelGoTmpl, filepath.Join(target, "model.go"))
	fmt.Printf("  created %s/model.go\n", name)

	// ── migrate.go ──────────────────────────────────────────────────
	writeTmpl(migrateGoTmpl, filepath.Join(target, "migrate.go"))
	fmt.Printf("  created %s/migrate.go\n", name)

	// ── params.go ───────────────────────────────────────────────────
	writeTmpl(paramsGoTmpl, filepath.Join(target, "params.go"))
	fmt.Printf("  created %s/params.go\n", name)

	// ── service.go ──────────────────────────────────────────────────
	writeTmpl(serviceGoTmpl, filepath.Join(target, "service.go"))
	fmt.Printf("  created %s/service.go\n", name)

	// ── api/v1/api.go ───────────────────────────────────────────────
	writeTmpl(apiV1GoTmpl, filepath.Join(target, "api", "v1", "api.go"))
	fmt.Printf("  created %s/api/v1/api.go\n", name)

	// ── Register in main.go ─────────────────────────────────────────
	registerInMainGo(repoRoot, name)

	fmt.Printf("\n✓ Created plugin scaffold: %s\n", name)
	fmt.Println("  Next steps:")
	fmt.Printf("  1. Edit plugins/%s/model.go — define your GORM models\n", name)
	fmt.Printf("  2. Edit plugins/%s/params.go — define request/response types\n", name)
	fmt.Printf("  3. Edit plugins/%s/service.go — implement business logic\n", name)
	fmt.Printf("  4. Edit plugins/%s/api/v1/api.go — implement route handlers\n", name)
	fmt.Printf("  5. go run main.go\n")
}

// ── Templates (@PKG@, @PASCAL@, @NAME@, @FULL@) ────────────────────

var pluginGoTmpl = `package @PKG@

import (
	"hei-gin/sdk/plugin"
	"hei-gin/sdk/plugin"
)

type @PASCAL@Plugin struct{}

func (p *@PASCAL@Plugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name:        "@NAME@",
		Version:     "1.0.0",
		Description: "@NAME@ plugin",
	}
}
func (p *@PASCAL@Plugin) Name() string { return "@FULL@" }
func (p *@PASCAL@Plugin) Init() error  { return nil }
func (p *@PASCAL@Plugin) Start() error { return nil }
func (p *@PASCAL@Plugin) Stop() error  { return nil }

func init() {
	plugin.Register(&@PASCAL@Plugin{})
}
`

var modelGoTmpl = `package @PKG@

import (
	"time"

	"gorm.io/gorm"
)

type Sample struct {
	ID        string         ` + "`" + `gorm:"primaryKey;size:32" json:"id"` + "`" + `
	Name      string         ` + "`" + `gorm:"size:100" json:"name"` + "`" + `
	CreatedAt *time.Time     ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt *time.Time     ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

func (Sample) TableName() string { return "@NAME@s" }
`

var migrateGoTmpl = `package @PKG@

import "hei-gin/sdk/db"

func init() {
	db.RegisterModel(&Sample{})
}
`

var paramsGoTmpl = `package @PKG@

type SampleVO struct {
	ID   string ` + "`" + `json:"id"` + "`" + `
	Name string ` + "`" + `json:"name"` + "`" + `
}
`

var serviceGoTmpl = `package @PKG@

import (
	"context"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/utils"
)

func Create(ctx context.Context, name string) (*Sample, error) {
	entity := &Sample{ID: utils.GenerateID(), Name: name}
	if err := db.DB.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, exception.NewBusinessError("创建失败", 500)
	}
	return entity, nil
}
`

var apiV1GoTmpl = `package v1

import (
	"github.com/gin-gonic/gin"

	"hei-gin/sdk/auth/middleware"
	"hei-gin/sdk/registry"
)

func init() {
	registry.RegisterRoute(func(r *gin.Engine) {
		g := r.Group("/api/v1/sys/@NAME@").Use(middleware.HeiCheckLogin())
		{
			g.GET("/page", pageHandler)
		}
	})
}

func pageHandler(c *gin.Context) {
	// TODO: implement
	c.JSON(200, gin.H{"code": 0, "data": nil, "success": true})
}
`

func registerInMainGo(repoRoot, pluginName string) {
	mainPath := filepath.Join(repoRoot, "main.go")
	data, err := os.ReadFile(mainPath)
	if err != nil {
		fatal("failed to read main.go: %v", err)
	}
	content := string(data)

	importLine := fmt.Sprintf("\t_ \"hei-gin/plugins/%s\"", pluginName)
	if strings.Contains(content, importLine) {
		return
	}

	marker := "// Plugin route/permission self-registration"
	idx := strings.Index(content, marker)
	if idx < 0 {
		fatal("cannot find plugin import marker in main.go")
	}
	closing := strings.Index(content[idx:], ")")
	if closing < 0 {
		fatal("cannot find import closing bracket")
	}
	insertPos := idx + closing

	newContent := content[:insertPos] + "\t" + importLine + "\n" + content[insertPos:]
	if err := os.WriteFile(mainPath, []byte(newContent), 0644); err != nil {
		fatal("failed to write main.go: %v", err)
	}
	fmt.Printf("  registered %s in main.go\n", pluginName)
}

// ── gen-imports ─────────────────────────────────────────────────────

func cmdGenImports(repoRoot string) {
	pluginsDir := filepath.Join(repoRoot, "plugins")
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		fatal("failed to read plugins dir: %v", err)
	}

	var imports []string
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		pluginFile := filepath.Join(pluginsDir, e.Name(), "plugin.go")
		if _, err := os.Stat(pluginFile); err == nil {
			imports = append(imports, fmt.Sprintf("\t_ \"hei-gin/plugins/%s\"", e.Name()))
		}
	}
	sort.Strings(imports)

	mainPath := filepath.Join(repoRoot, "main.go")
	data, err := os.ReadFile(mainPath)
	if err != nil {
		fatal("failed to read main.go: %v", err)
	}
	content := string(data)

	marker := "// Plugin route/permission self-registration"
	idx := strings.Index(content, marker)
	if idx < 0 {
		fatal("cannot find plugin import marker in main.go")
	}

	before := content[:idx]
	after := content[idx:]
	closing := strings.Index(after, ")")
	if closing < 0 {
		fatal("cannot find import closing bracket")
	}
	after = after[closing:]

	newContent := before + marker + "\n" + strings.Join(imports, "\n") + "\n" + after
	if err := os.WriteFile(mainPath, []byte(newContent), 0644); err != nil {
		fatal("failed to write main.go: %v", err)
	}

	fmt.Printf("Generated %d imports in main.go\n", len(imports))
}
