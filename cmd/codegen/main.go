package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: go run cmd/codegen/main.go <command>

Commands:
  list                    List all discovered plugins
  scaffold <name>         Create a new plugin with demo sub-module (e.g. plugin-xxx)
  add-module <plugin> <module>  Add a sub-module to an existing plugin (e.g. plugin-xxx modname)
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
	case "add-module":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: go run cmd/codegen/main.go add-module <plugin> <module>")
			os.Exit(1)
		}
		cmdAddModule(repoRoot, args[1], args[2])
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

	moduleName := "demo"
	modulePascal := "Demo"

	// Create directories
	os.MkdirAll(filepath.Join(target, moduleName, "api", "v1"), 0755)

	writeTmpl := func(tmpl, path string) {
		r := strings.NewReplacer(
			"@PKG@", pkgName,
			"@PASCAL@", pascalName,
			"@NAME@", pluginName,
			"@FULL@", name,
			"@MODULE@", moduleName,
			"@MODULE_PKG@", moduleImportAlias(moduleName),
			"@MODULE_PASCAL@", modulePascal,
		)
		result := r.Replace(tmpl)
		if err := os.WriteFile(path, []byte(result), 0644); err != nil {
			fatal("failed to write %s: %v", path, err)
		}
	}

	// ── plugin.go ───────────────────────────────────────────────────
	writeTmpl(pluginGoTmpl, filepath.Join(target, "plugin.go"))
	fmt.Printf("  created %s/plugin.go\n", name)

	// ── go.mod ──────────────────────────────────────────────────────
	writeTmpl(goModTmpl, filepath.Join(target, "go.mod"))
	fmt.Printf("  created %s/go.mod\n", name)

	// ── migrate.go ──────────────────────────────────────────────────
	writeTmpl(migrateGoTmpl, filepath.Join(target, "migrate.go"))
	fmt.Printf("  created %s/migrate.go\n", name)

	// ── demo/model.go ───────────────────────────────────────────────
	writeTmpl(modelGoTmpl, filepath.Join(target, moduleName, "model.go"))
	fmt.Printf("  created %s/%s/model.go\n", name, moduleName)

	// ── demo/params.go ──────────────────────────────────────────────
	writeTmpl(paramsGoTmpl, filepath.Join(target, moduleName, "params.go"))
	fmt.Printf("  created %s/%s/params.go\n", name, moduleName)

	// ── demo/repository.go ──────────────────────────────────────────
	writeTmpl(repositoryGoTmpl, filepath.Join(target, moduleName, "repository.go"))
	fmt.Printf("  created %s/%s/repository.go\n", name, moduleName)

	// ── demo/service.go ─────────────────────────────────────────────
	writeTmpl(serviceGoTmpl, filepath.Join(target, moduleName, "service.go"))
	fmt.Printf("  created %s/%s/service.go\n", name, moduleName)

	// ── demo/module.go ──────────────────────────────────────────────
	writeTmpl(moduleGoTmpl, filepath.Join(target, moduleName, "module.go"))
	fmt.Printf("  created %s/%s/module.go\n", name, moduleName)

	// ── demo/api/v1/api.go ──────────────────────────────────────────
	writeTmpl(apiV1GoTmpl, filepath.Join(target, moduleName, "api", "v1", "api.go"))
	fmt.Printf("  created %s/%s/api/v1/api.go\n", name, moduleName)

	gofmtTree(target)
	runGoModTidy(target)

	fmt.Printf("\n✓ Created plugin scaffold: %s\n", name)
	fmt.Printf("  Next steps:\n")
	fmt.Printf("  1. Edit plugins/%s/%s/model.go — define your GORM models\n", name, moduleName)
	fmt.Printf("  2. Edit plugins/%s/%s/params.go — define request/response types\n", name, moduleName)
	fmt.Printf("  3. Edit plugins/%s/%s/repository.go — implement data access\n", name, moduleName)
	fmt.Printf("  4. Edit plugins/%s/%s/service.go — implement business logic\n", name, moduleName)
	fmt.Printf("  5. Edit plugins/%s/%s/api/v1/api.go — implement route handlers\n", name, moduleName)
	fmt.Printf("  6. Add %s.RegisterPlugin/RegisterRoutes/RegisterMigrations in main.go and cmd/migrate/main.go\n", pkgName)
	fmt.Printf("  7. go run main.go\n")
}

// ── add-module ─────────────────────────────────────────────────────

func cmdAddModule(repoRoot, pluginName, moduleName string) {
	if !strings.HasPrefix(pluginName, "plugin-") {
		pluginName = "plugin-" + pluginName
	}
	pluginDir := filepath.Join(repoRoot, "plugins", pluginName)
	if _, err := os.Stat(pluginDir); err != nil {
		fatal("plugin %s not found", pluginName)
	}

	modulePascal := ""
	for _, p := range strings.Split(moduleName, "-") {
		if len(p) > 0 {
			modulePascal += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	if modulePascal == "" {
		modulePascal = "Demo"
	}

	// Create sub-module directories
	os.MkdirAll(filepath.Join(pluginDir, moduleName, "api", "v1"), 0755)

	// Compute main plugin naming variants
	pkgName := strings.ReplaceAll(pluginName, "-", "_")
	pluginShortName := strings.TrimPrefix(pluginName, "plugin-")
	parts := strings.Split(pluginShortName, "-")
	pluginPascal := ""
	for _, p := range parts {
		if len(p) > 0 {
			pluginPascal += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	if pluginPascal == "" {
		pluginPascal = "Sample"
	}

	if _, err := os.Stat(filepath.Join(pluginDir, "plugin.go")); err != nil {
		fatal("plugin.go not found in %s", pluginName)
	}
	if _, err := os.Stat(filepath.Join(pluginDir, "migrate.go")); err != nil {
		fatal("migrate.go not found in %s", pluginName)
	}

	writeTmpl := func(tmpl, path string) {
		r := strings.NewReplacer(
			"@PKG@", pkgName,
			"@PASCAL@", pluginPascal,
			"@NAME@", pluginShortName,
			"@FULL@", pluginName,
			"@MODULE@", moduleName,
			"@MODULE_PKG@", moduleImportAlias(moduleName),
			"@MODULE_PASCAL@", modulePascal,
		)
		result := r.Replace(tmpl)
		if err := os.WriteFile(path, []byte(result), 0644); err != nil {
			fatal("failed to write %s: %v", path, err)
		}
	}

	// ── model.go ───────────────────────────────────────────────────
	writeTmpl(modelGoTmpl, filepath.Join(pluginDir, moduleName, "model.go"))
	fmt.Printf("  created %s/%s/model.go\n", pluginName, moduleName)

	// ── params.go ──────────────────────────────────────────────────
	writeTmpl(paramsGoTmpl, filepath.Join(pluginDir, moduleName, "params.go"))
	fmt.Printf("  created %s/%s/params.go\n", pluginName, moduleName)

	// ── repository.go ──────────────────────────────────────────────
	writeTmpl(repositoryGoTmpl, filepath.Join(pluginDir, moduleName, "repository.go"))
	fmt.Printf("  created %s/%s/repository.go\n", pluginName, moduleName)

	// ── service.go ─────────────────────────────────────────────────
	writeTmpl(serviceGoTmpl, filepath.Join(pluginDir, moduleName, "service.go"))
	fmt.Printf("  created %s/%s/service.go\n", pluginName, moduleName)

	// ── module.go ──────────────────────────────────────────────────
	writeTmpl(moduleGoTmpl, filepath.Join(pluginDir, moduleName, "module.go"))
	fmt.Printf("  created %s/%s/module.go\n", pluginName, moduleName)

	// ── api/v1/api.go ──────────────────────────────────────────────
	writeTmpl(apiV1GoTmpl, filepath.Join(pluginDir, moduleName, "api", "v1", "api.go"))
	fmt.Printf("  created %s/%s/api/v1/api.go\n", pluginName, moduleName)

	updatePluginRegistrations(pluginDir, pluginName, moduleName)
	gofmtTree(pluginDir)

	fmt.Printf("\n✓ Added sub-module %s to %s\n", moduleName, pluginName)
}

func updatePluginRegistrations(pluginDir, pluginName, moduleName string) {
	updatePluginFile(pluginDir, pluginName, moduleName)
	updateMigrateFile(pluginDir, pluginName, moduleName)
}

func updatePluginFile(pluginDir, pluginName, moduleName string) {
	pluginPath := filepath.Join(pluginDir, "plugin.go")
	data, err := os.ReadFile(pluginPath)
	if err != nil {
		fatal("failed to read plugin.go: %v", err)
	}
	content := string(data)

	alias := moduleImportAlias(moduleName)
	importLine := fmt.Sprintf("\t%sv1 \"hei-gin/plugins/%s/%s/api/v1\"\n", alias, pluginName, moduleName)
	callLine := fmt.Sprintf("\t%sv1.Register()\n", alias)

	if strings.Contains(content, importLine) && strings.Contains(content, callLine) {
		fmt.Printf("  plugin.go: %s already registered, skipping\n", moduleName)
		return
	}

	importBlockEnd := strings.Index(content, "\n)\n")
	if importBlockEnd < 0 {
		fatal("cannot find import block in plugin.go")
	}
	content = content[:importBlockEnd] + "\n" + importLine + content[importBlockEnd:]

	marker := "func RegisterRoutes() {\n"
	idx := strings.Index(content, marker)
	if idx < 0 {
		fatal("cannot find RegisterRoutes in plugin.go")
	}
	insertPos := idx + len(marker)
	content = content[:insertPos] + callLine + content[insertPos:]

	if err := os.WriteFile(pluginPath, []byte(content), 0644); err != nil {
		fatal("failed to write plugin.go: %v", err)
	}
	fmt.Printf("  updated %s/plugin.go\n", pluginName)
}

func updateMigrateFile(pluginDir, pluginName, moduleName string) {
	migratePath := filepath.Join(pluginDir, "migrate.go")
	data, err := os.ReadFile(migratePath)
	if err != nil {
		fatal("failed to read migrate.go: %v", err)
	}
	content := string(data)

	alias := moduleImportAlias(moduleName)
	importLine := fmt.Sprintf("\t%s \"hei-gin/plugins/%s/%s\"\n", alias, pluginName, moduleName)
	callLine := fmt.Sprintf("\tdb.RegisterModel(&%s.%s{})\n", alias, modulePascalName(moduleName))

	if strings.Contains(content, importLine) && strings.Contains(content, callLine) {
		fmt.Printf("  migrate.go: %s already registered, skipping\n", moduleName)
		return
	}

	importBlockEnd := strings.Index(content, "\n)\n")
	if importBlockEnd < 0 {
		fatal("cannot find import block in migrate.go")
	}
	content = content[:importBlockEnd] + "\n" + importLine + content[importBlockEnd:]

	marker := "func RegisterMigrations() {\n"
	idx := strings.Index(content, marker)
	if idx < 0 {
		fatal("cannot find RegisterMigrations in migrate.go")
	}
	insertPos := idx + len(marker)
	content = content[:insertPos] + callLine + content[insertPos:]

	if err := os.WriteFile(migratePath, []byte(content), 0644); err != nil {
		fatal("failed to write migrate.go: %v", err)
	}
	fmt.Printf("  updated %s/migrate.go\n", pluginName)
}

func moduleImportAlias(moduleName string) string {
	return strings.ReplaceAll(moduleName, "-", "_")
}

func modulePascalName(moduleName string) string {
	parts := strings.Split(moduleName, "-")
	pascal := ""
	for _, p := range parts {
		if p == "" {
			continue
		}
		pascal += strings.ToUpper(p[:1]) + p[1:]
	}
	if pascal == "" {
		return "Demo"
	}
	return pascal
}

func gofmtTree(root string) {
	cmd := exec.Command("gofmt", "-w", root)
	if err := cmd.Run(); err == nil {
		return
	}

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		if runErr := exec.Command("gofmt", "-w", path).Run(); runErr != nil {
			return runErr
		}
		return nil
	})
}

func runGoModTidy(dir string) {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	_ = cmd.Run()
}

// ── Templates (@PKG@, @PASCAL@, @NAME@, @FULL@) ────────────────────

var goModTmpl = `module hei-gin/plugins/@FULL@

go 1.25.10

require (
	hei-gin/sdk v0.0.0
	github.com/gin-gonic/gin v1.12.0
	gorm.io/gorm v1.25.12
)

replace hei-gin/sdk => ../../sdk
`

var pluginGoTmpl = `package @PKG@

import (
	"sync"

		@MODULE_PKG@v1 "hei-gin/plugins/@FULL@/@MODULE@/api/v1"
	"hei-gin/sdk/kernel/plugin"
)

type @PASCAL@Plugin struct {
	plugin.NoopPlugin
}

var registerOnce sync.Once

func (p *@PASCAL@Plugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name:        "@FULL@",
		Version:     "1.0.0",
		Description: "@NAME@ plugin",
	}
}
func (p *@PASCAL@Plugin) Name() string { return "@FULL@" }

func RegisterPlugin() {
	registerOnce.Do(func() {
		plugin.Register(&@PASCAL@Plugin{})
	})
}

func RegisterRoutes() {
	@MODULE_PKG@v1.Register()
}
`

var migrateGoTmpl = `package @PKG@

import (
	@MODULE_PKG@ "hei-gin/plugins/@FULL@/@MODULE@"
	"hei-gin/sdk/infra/db"
)

func RegisterMigrations() {
	db.RegisterModel(&@MODULE_PKG@.@MODULE_PASCAL@{})
}
`

var modelGoTmpl = `package @MODULE_PKG@

import "time"

type @MODULE_PASCAL@ struct {
	ID        string     ` + "`" + `gorm:"primaryKey;size:32" json:"id"` + "`" + `
	Name      string     ` + "`" + `gorm:"size:100" json:"name"` + "`" + `
	Status    string     ` + "`" + `gorm:"size:16;default:ENABLED" json:"status"` + "`" + `
	SortCode  int        ` + "`" + `gorm:"default:0" json:"sort_code"` + "`" + `
	CreatedAt *time.Time ` + "`" + `json:"created_at"` + "`" + `
	CreatedBy *string    ` + "`" + `gorm:"size:32" json:"created_by"` + "`" + `
	UpdatedAt *time.Time ` + "`" + `json:"updated_at"` + "`" + `
	UpdatedBy *string    ` + "`" + `gorm:"size:32" json:"updated_by"` + "`" + `
}

func (@MODULE_PASCAL@) TableName() string { return "sys_@MODULE@_template" }
`

var paramsGoTmpl = `package @MODULE_PKG@

import "hei-gin/sdk/utils"

// @MODULE_PASCAL@VO 视图对象
type @MODULE_PASCAL@VO struct {
	ID        string  ` + "`" + `json:"id"` + "`" + `
	Name      string  ` + "`" + `json:"name"` + "`" + `
	Status    string  ` + "`" + `json:"status"` + "`" + `
	SortCode  int     ` + "`" + `json:"sort_code"` + "`" + `
	CreatedAt string  ` + "`" + `json:"created_at"` + "`" + `
	CreatedBy *string ` + "`" + `json:"created_by"` + "`" + `
	UpdatedAt string  ` + "`" + `json:"updated_at"` + "`" + `
	UpdatedBy *string ` + "`" + `json:"updated_by"` + "`" + `
}

// @MODULE_PASCAL@PageParam 分页参数
type @MODULE_PASCAL@PageParam struct {
	Current int    ` + "`" + `json:"current" form:"current"` + "`" + `
	Size    int    ` + "`" + `json:"size" form:"size"` + "`" + `
	Keyword string ` + "`" + `json:"keyword" form:"keyword"` + "`" + `
}

func Sys@MODULE_PASCAL@To@MODULE_PASCAL@VO(src *@MODULE_PASCAL@) *@MODULE_PASCAL@VO {
	if src == nil {
		return nil
	}
	dst := &@MODULE_PASCAL@VO{}
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	// *time.Time → string
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)

	return dst
}

// @MODULE_PASCAL@VOToSys@MODULE_PASCAL@ 将 VO 映射到 entity
func @MODULE_PASCAL@VOToSys@MODULE_PASCAL@(src *@MODULE_PASCAL@VO) *@MODULE_PASCAL@ {
	if src == nil {
		return nil
	}
	dst := &@MODULE_PASCAL@{}
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	// string → *time.Time
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)

	return dst
}
`

var repositoryGoTmpl = `package @MODULE_PKG@

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) Page(ctx context.Context, p *@MODULE_PASCAL@PageParam) ([]@MODULE_PASCAL@, int64) {
	q := r.db.WithContext(ctx).Model(&@MODULE_PASCAL@{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("name LIKE ?", like)
	}

	var total int64
	q.Count(&total)

	var rows []@MODULE_PASCAL@
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)
	return rows, total
}

func (r *repository) FindByID(ctx context.Context, id string) (*@MODULE_PASCAL@, error) {
	var e @MODULE_PASCAL@
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) Create(ctx context.Context, entity *@MODULE_PASCAL@) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *repository) UpdateByID(ctx context.Context, id string, up map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&@MODULE_PASCAL@{}).Where("id = ?", id).Updates(up).Error
}

func (r *repository) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&@MODULE_PASCAL@{}).Error
}
`

var serviceGoTmpl = `package @MODULE_PKG@

import (
	"gorm.io/gorm"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/web/exception"
	"hei-gin/sdk/web/result"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *repository
}

func (s *Service) Page(c *gin.Context, p *@MODULE_PASCAL@PageParam) {
	ctx := c.Request.Context()
	if p.Current < 1 {
		p.Current = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}

	rows, total := s.repo.Page(ctx, p)

	vos := make([]*@MODULE_PASCAL@VO, len(rows))
	for i, r := range rows {
		vos[i] = Sys@MODULE_PASCAL@To@MODULE_PASCAL@VO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func (s *Service) Detail(c *gin.Context, id string) *@MODULE_PASCAL@VO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询失败: "+err.Error(), 500))
		return nil
	}
	return Sys@MODULE_PASCAL@To@MODULE_PASCAL@VO(e)
}

func (s *Service) Create(c *gin.Context, vo *@MODULE_PASCAL@VO) {
	ctx := c.Request.Context()
	e := @MODULE_PASCAL@VOToSys@MODULE_PASCAL@(vo)
	e.Status = string(enums.StatusEnabled)
	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("创建失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Remove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	if err := s.repo.DeleteByIDs(c.Request.Context(), ids); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除失败: "+err.Error(), 500))
		return
	}
}

func (s *Service) Modify(c *gin.Context, vo *@MODULE_PASCAL@VO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	_, err := s.repo.FindByID(ctx, vo.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询失败: "+err.Error(), 500))
		return
	}

	up := map[string]interface{}{
		"name": vo.Name, "sort_code": vo.SortCode,
	}
	if vo.Status != "" {
		up["status"] = vo.Status
	}
	if err := s.repo.UpdateByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑失败: "+err.Error(), 500))
		return
	}
}
`

var moduleGoTmpl = `package @MODULE_PKG@

import "hei-gin/sdk/infra/db"

type Module struct {
	service *Service
}

var DefaultModule = NewModule()

func NewModule() *Module {
	repo := &repository{db: db.DB}
	svc := &Service{repo: repo}
	return &Module{service: svc}
}

func (m *Module) Service() *Service {
	return m.service
}
`

var apiV1GoTmpl = `package v1

import (
	"hei-gin/sdk/utils"
	"hei-gin/sdk/web/result"
	"hei-gin/sdk/kernel/registry"
	"hei-gin/sdk/log"
	@MODULE_PKG@ "hei-gin/plugins/@FULL@/@MODULE@"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *@MODULE_PKG@.Service
}

var defaultHandler = newHandler(@MODULE_PKG@.DefaultModule)

func newHandler(module *@MODULE_PKG@.Module) *handler {
	return &handler{service: module.Service()}
}

func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/@MODULE@/page
	r.GET("/api/v1/sys/@MODULE@/page",
		registry.Perm("sys:@MODULE@:page", "@MODULE_PASCAL@分页"),
		defaultHandler.page,
	)

	// POST /api/v1/sys/@MODULE@/create
	r.POST("/api/v1/sys/@MODULE@/create",
		registry.Perm("sys:@MODULE@:create", "添加@MODULE_PASCAL@"),
		log.SysLog("添加@MODULE_PASCAL@"),
		defaultHandler.create,
	)

	// GET /api/v1/sys/@MODULE@/detail
	r.GET("/api/v1/sys/@MODULE@/detail",
		registry.Perm("sys:@MODULE@:detail", "@MODULE_PASCAL@详情"),
		defaultHandler.detail,
	)

	// POST /api/v1/sys/@MODULE@/modify
	r.POST("/api/v1/sys/@MODULE@/modify",
		registry.Perm("sys:@MODULE@:modify", "编辑@MODULE_PASCAL@"),
		log.SysLog("编辑@MODULE_PASCAL@"),
		defaultHandler.modify,
	)

	// POST /api/v1/sys/@MODULE@/remove
	r.POST("/api/v1/sys/@MODULE@/remove",
		registry.Perm("sys:@MODULE@:remove", "删除@MODULE_PASCAL@"),
		log.SysLog("删除@MODULE_PASCAL@"),
		defaultHandler.remove,
	)
}

func Register() {
	registry.RegisterRoute(RegisterRoutes)
}

func (h *handler) page(c *gin.Context) {
	var param @MODULE_PKG@.@MODULE_PASCAL@PageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Page(c, &param)
}

func (h *handler) create(c *gin.Context) {
	var vo @MODULE_PKG@.@MODULE_PASCAL@VO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Create(c, &vo)
	result.Success(c, nil)
}

func (h *handler) detail(c *gin.Context) {
	vo := h.service.Detail(c, c.Query("id"))
	result.Success(c, vo)
}

func (h *handler) modify(c *gin.Context) {
	var vo @MODULE_PKG@.@MODULE_PASCAL@VO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Modify(c, &vo)
	result.Success(c, nil)
}

func (h *handler) remove(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	h.service.Remove(c, &param)
	result.Success(c, nil)
}
`
