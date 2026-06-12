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
  scaffold <name>         Create a new plugin with demo sub-module (e.g. plugin-xxx)
  add-module <plugin> <module>  Add a sub-module to an existing plugin (e.g. plugin-xxx modname)
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
	case "add-module":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: go run cmd/codegen/main.go add-module <plugin> <module>")
			os.Exit(1)
		}
		cmdAddModule(repoRoot, args[1], args[2])
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

	// ── imports.go ─────────────────────────────────────────────────
	writeTmpl(importsGoTmpl, filepath.Join(target, "imports.go"))
	fmt.Printf("  created %s/imports.go\n", name)

	// ── demo/model.go ───────────────────────────────────────────────
	writeTmpl(modelGoTmpl, filepath.Join(target, moduleName, "model.go"))
	fmt.Printf("  created %s/%s/model.go\n", name, moduleName)

	// ── demo/params.go ──────────────────────────────────────────────
	writeTmpl(paramsGoTmpl, filepath.Join(target, moduleName, "params.go"))
	fmt.Printf("  created %s/%s/params.go\n", name, moduleName)

	// ── demo/mapper.go ──────────────────────────────────────────────
	writeTmpl(mapperGoTmpl, filepath.Join(target, moduleName, "mapper.go"))
	fmt.Printf("  created %s/%s/mapper.go\n", name, moduleName)

	// ── demo/service.go ─────────────────────────────────────────────
	writeTmpl(serviceGoTmpl, filepath.Join(target, moduleName, "service.go"))
	fmt.Printf("  created %s/%s/service.go\n", name, moduleName)

	// ── demo/migrate.go ─────────────────────────────────────────────
	writeTmpl(migrateGoTmpl, filepath.Join(target, moduleName, "migrate.go"))
	fmt.Printf("  created %s/%s/migrate.go\n", name, moduleName)

	// ── demo/api/v1/api.go ──────────────────────────────────────────
	writeTmpl(apiV1GoTmpl, filepath.Join(target, moduleName, "api", "v1", "api.go"))
	fmt.Printf("  created %s/%s/api/v1/api.go\n", name, moduleName)

	// ── Register in main.go ─────────────────────────────────────────
	registerInMainGo(repoRoot, name)

	fmt.Printf("\n✓ Created plugin scaffold: %s\n", name)
	fmt.Printf("  Next steps:\n")
	fmt.Printf("  1. Edit plugins/%s/%s/model.go — define your GORM models\n", name, moduleName)
	fmt.Printf("  2. Edit plugins/%s/%s/params.go — define request/response types\n", name, moduleName)
	fmt.Printf("  3. Edit plugins/%s/%s/service.go — implement business logic\n", name, moduleName)
	fmt.Printf("  4. Edit plugins/%s/%s/api/v1/api.go — implement route handlers\n", name, moduleName)
	fmt.Printf("  5. go run main.go\n")
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

	writeTmpl := func(tmpl, path string) {
		r := strings.NewReplacer(
			"@PKG@", pkgName,
			"@PASCAL@", pluginPascal,
			"@NAME@", pluginShortName,
			"@FULL@", pluginName,
			"@MODULE@", moduleName,
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

	// ── mapper.go ──────────────────────────────────────────────────
	writeTmpl(mapperGoTmpl, filepath.Join(pluginDir, moduleName, "mapper.go"))
	fmt.Printf("  created %s/%s/mapper.go\n", pluginName, moduleName)

	// ── service.go ─────────────────────────────────────────────────
	writeTmpl(serviceGoTmpl, filepath.Join(pluginDir, moduleName, "service.go"))
	fmt.Printf("  created %s/%s/service.go\n", pluginName, moduleName)

	// ── migrate.go ─────────────────────────────────────────────────
	writeTmpl(migrateGoTmpl, filepath.Join(pluginDir, moduleName, "migrate.go"))
	fmt.Printf("  created %s/%s/migrate.go\n", pluginName, moduleName)

	// ── api/v1/api.go ──────────────────────────────────────────────
	writeTmpl(apiV1GoTmpl, filepath.Join(pluginDir, moduleName, "api", "v1", "api.go"))
	fmt.Printf("  created %s/%s/api/v1/api.go\n", pluginName, moduleName)

	// Update imports.go
	updateImportsGo(pluginDir, pkgName, pluginName, moduleName)

	fmt.Printf("\n✓ Added sub-module %s to %s\n", moduleName, pluginName)
}

func updateImportsGo(pluginDir, pkgName, pluginName, moduleName string) {
	importsPath := filepath.Join(pluginDir, "imports.go")
	data, err := os.ReadFile(importsPath)
	if err != nil {
		fatal("failed to read imports.go: %v", err)
	}
	content := string(data)

	modelImport := fmt.Sprintf("\t_ \"hei-gin/plugins/%s/%s\"", pluginName, moduleName)
	routeImport := fmt.Sprintf("\t_ \"hei-gin/plugins/%s/%s/api/v1\"", pluginName, moduleName)

	if strings.Contains(content, modelImport) {
		fmt.Printf("  imports.go: %s already registered, skipping\n", modelImport)
		return
	}

	// Find "// Route registrations" section and add before it
	routeMarker := "\t// Route registrations"

	// Add model import after model registrations section
	modelSectionEnd := strings.Index(content, routeMarker)
	if modelSectionEnd < 0 {
		// Fallback: add at end of import block
		closing := strings.LastIndex(content, ")")
		if closing < 0 {
			fatal("cannot find import closing bracket in imports.go")
		}
		content = content[:closing] + modelImport + "\n" + routeImport + "\n" + content[closing:]
	} else {
		// Add model import before route marker
		content = content[:modelSectionEnd] + modelImport + "\n\n" + content[modelSectionEnd:]
		// Add route import in route section
		closing := strings.LastIndex(content, ")")
		if closing >= 0 {
			content = content[:closing] + routeImport + "\n" + content[closing:]
		}
	}

	if err := os.WriteFile(importsPath, []byte(content), 0644); err != nil {
		fatal("failed to write imports.go: %v", err)
	}
	fmt.Printf("  updated %s/imports.go\n", pluginName)
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
	"hei-gin/sdk/plugin"
)

type @PASCAL@Plugin struct {
	plugin.NoopPlugin
}

func (p *@PASCAL@Plugin) Info() plugin.PluginInfo {
	return plugin.PluginInfo{
		Name:        "@NAME@",
		Version:     "1.0.0",
		Description: "@NAME@ plugin",
	}
}
func (p *@PASCAL@Plugin) Name() string { return "@FULL@" }

func init() {
	plugin.Register(&@PASCAL@Plugin{})
}
`

var importsGoTmpl = `package @PKG@

import (
	// Model registrations (migrate.go)
	_ "hei-gin/plugins/@FULL@/@MODULE@"

	// Route registrations (api/v1/api.go)
	_ "hei-gin/plugins/@FULL@/@MODULE@/api/v1"
)
`

var modelGoTmpl = `package @MODULE@

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

var migrateGoTmpl = `package @MODULE@

import "hei-gin/sdk/db"

func init() {
	db.RegisterModel(&@MODULE_PASCAL@{})
}
`

var paramsGoTmpl = `package @MODULE@

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
`

var mapperGoTmpl = `package @MODULE@

import "hei-gin/sdk/utils"

// Sys@MODULE_PASCAL@To@MODULE_PASCAL@VO 将 entity 映射到 VO
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

var serviceGoTmpl = `package @MODULE@

import (
	"gorm.io/gorm"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

// @MODULE_PASCAL@Page 分页
func @MODULE_PASCAL@Page(c *gin.Context, p *@MODULE_PASCAL@PageParam) {
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

	q := db.DB.WithContext(ctx).Model(&@MODULE_PASCAL@{})
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("name LIKE ?", like)
	}

	var total int64
	q.Count(&total)

	var rows []@MODULE_PASCAL@
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

	vos := make([]*@MODULE_PASCAL@VO, len(rows))
	for i, r := range rows {
		vos[i] = Sys@MODULE_PASCAL@To@MODULE_PASCAL@VO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

// @MODULE_PASCAL@Detail 详情
func @MODULE_PASCAL@Detail(c *gin.Context, id string) *@MODULE_PASCAL@VO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e @MODULE_PASCAL@
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询失败: "+err.Error(), 500))
		return nil
	}
	return Sys@MODULE_PASCAL@To@MODULE_PASCAL@VO(&e)
}

// @MODULE_PASCAL@Create 创建
func @MODULE_PASCAL@Create(c *gin.Context, vo *@MODULE_PASCAL@VO) {
	ctx := c.Request.Context()

	e := @MODULE_PASCAL@VOToSys@MODULE_PASCAL@(vo)
	e.Status = string(enums.StatusEnabled)
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("创建失败: "+err.Error(), 500))
		return
	}
}

// @MODULE_PASCAL@Remove 删除
func @MODULE_PASCAL@Remove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	if err := db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&@MODULE_PASCAL@{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除失败: "+err.Error(), 500))
		return
	}
}

// @MODULE_PASCAL@Modify 编辑
func @MODULE_PASCAL@Modify(c *gin.Context, vo *@MODULE_PASCAL@VO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	var e @MODULE_PASCAL@
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", vo.ID).Error; err != nil {
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
	if err := db.DB.WithContext(ctx).Model(&@MODULE_PASCAL@{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑失败: "+err.Error(), 500))
		return
	}
}
`

var apiV1GoTmpl = `package v1

import (
	"hei-gin/sdk/utils"
	"hei-gin/sdk/result"
	"hei-gin/sdk/registry"
	"hei-gin/sdk/log"
	@MODULE@ "hei-gin/plugins/@FULL@/@MODULE@"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// GET /api/v1/sys/@MODULE@/page
	r.GET("/api/v1/sys/@MODULE@/page",
		registry.Perm("sys:@MODULE@:page", "@MODULE_PASCAL@分页"),
		pageHandler,
	)

	// POST /api/v1/sys/@MODULE@/create
	r.POST("/api/v1/sys/@MODULE@/create",
		registry.Perm("sys:@MODULE@:create", "添加@MODULE_PASCAL@"),
		log.SysLog("添加@MODULE_PASCAL@"),
		createHandler,
	)

	// GET /api/v1/sys/@MODULE@/detail
	r.GET("/api/v1/sys/@MODULE@/detail",
		registry.Perm("sys:@MODULE@:detail", "@MODULE_PASCAL@详情"),
		detailHandler,
	)

	// POST /api/v1/sys/@MODULE@/modify
	r.POST("/api/v1/sys/@MODULE@/modify",
		registry.Perm("sys:@MODULE@:modify", "编辑@MODULE_PASCAL@"),
		log.SysLog("编辑@MODULE_PASCAL@"),
		modifyHandler,
	)

	// POST /api/v1/sys/@MODULE@/remove
	r.POST("/api/v1/sys/@MODULE@/remove",
		registry.Perm("sys:@MODULE@:remove", "删除@MODULE_PASCAL@"),
		log.SysLog("删除@MODULE_PASCAL@"),
		removeHandler,
	)
}

func init() {
	registry.RegisterRoute(RegisterRoutes)
}

// pageHandler handles GET /api/v1/sys/@MODULE@/page
func pageHandler(c *gin.Context) {
	var param @MODULE@.@MODULE_PASCAL@PageParam
	if err := c.ShouldBindQuery(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	@MODULE@.@MODULE_PASCAL@Page(c, &param)
}

// createHandler handles POST /api/v1/sys/@MODULE@/create
func createHandler(c *gin.Context) {
	var vo @MODULE@.@MODULE_PASCAL@VO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	@MODULE@.@MODULE_PASCAL@Create(c, &vo)
	result.Success(c, nil)
}

// detailHandler handles GET /api/v1/sys/@MODULE@/detail
func detailHandler(c *gin.Context) {
	vo := @MODULE@.@MODULE_PASCAL@Detail(c, c.Query("id"))
	result.Success(c, vo)
}

// modifyHandler handles POST /api/v1/sys/@MODULE@/modify
func modifyHandler(c *gin.Context) {
	var vo @MODULE@.@MODULE_PASCAL@VO
	if err := c.ShouldBindJSON(&vo); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	@MODULE@.@MODULE_PASCAL@Modify(c, &vo)
	result.Success(c, nil)
}

// removeHandler handles POST /api/v1/sys/@MODULE@/remove
func removeHandler(c *gin.Context) {
	var param utils.IdsParam
	if err := c.ShouldBindJSON(&param); err != nil {
		result.Failure(c, "参数错误: "+err.Error(), 400)
		return
	}
	@MODULE@.@MODULE_PASCAL@Remove(c, &param)
	result.Success(c, nil)
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
