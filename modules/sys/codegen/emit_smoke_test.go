package codegen

import (
	"context"
	"testing"
)

// TestRenderPlanSmoke 校验模板可解析并渲染出 Go 文件。
func TestRenderPlanSmoke(t *testing.T) {
	main := []Field{
		{ColumnName: "id", DBType: "varchar", PythonType: "str", TypescriptType: "string", IsPrimaryKey: true, ShowInTable: true, ShowInForm: false, ShowInDetail: true, Sort: 1},
		{ColumnName: "name", ColumnComment: strptr("名称"), DBType: "varchar", PythonType: "str", TypescriptType: "string", ShowInTable: true, ShowInForm: true, ShowInDetail: true, Sort: 2},
		{ColumnName: "status", ColumnComment: strptr("状态"), DBType: "varchar", PythonType: "str", TypescriptType: "string", FormWidget: "dict", DictCode: strptr("COMMON_STATUS"), ShowInTable: true, ShowInForm: true, ShowInDetail: true, ShowInQuery: true, Sort: 3},
		{ColumnName: "count", DBType: "int4", PythonType: "int", TypescriptType: "number", FormWidget: "number", ShowInTable: true, ShowInForm: true, ShowInDetail: true, Sort: 4},
		{ColumnName: "created_at", DBType: "timestamptz", PythonType: "datetime", TypescriptType: "string", ShowInTable: true, ShowInForm: false, ShowInDetail: true, Sort: 5},
	}
	plan := &Plan{
		Name: "smoke", GenType: "TABLE", Author: "Charlie", MainTable: "cg_test_smoke",
		MainPK: "id", MainEntityName: "Smoke", MainModulePath: "biz/smoke", MainBusinessName: "Smoke",
		APIPrefix: "/biz/smoke", PermissionPrefix: "biz:smoke", MenuName: "Smoke",
		MenuPath: "/biz/smoke", ComponentPath: "biz/smoke/index.vue", Sort: 99,
	}
	files, err := renderPlan(plan, main, nil)
	if err != nil {
		t.Fatalf("renderPlan failed: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no files rendered")
	}
	paths := map[string]bool{}
	for _, f := range files {
		paths[f.Path] = true
		if f.Content == "" {
			t.Errorf("empty content for %s", f.Path)
		}
	}
	for _, want := range []string{
		"modules/biz/smoke/model.go",
		"modules/biz/smoke/param.go",
		"modules/biz/smoke/result.go",
		"modules/biz/smoke/repo.go",
		"modules/biz/smoke/service.go",
		"modules/biz/smoke/handler.go",
		"modules/biz/smoke/register.go",
		"scripts/smoke_menu_permission.sql",
	} {
		if !paths[want] {
			t.Errorf("missing expected file %s (have %d files)", want, len(files))
		}
	}
	_ = context.Background()
}

func strptr(s string) *string { return &s }

// TestRenderPlainTable 校验无 dict/json 字段的表也能渲染（不产生未使用 import）。
func TestRenderPlainTable(t *testing.T) {
	main := []Field{
		{ColumnName: "id", DBType: "varchar", PythonType: "str", TypescriptType: "string", IsPrimaryKey: true, ShowInTable: true, ShowInForm: false, ShowInDetail: true, Sort: 1},
		{ColumnName: "name", DBType: "varchar", PythonType: "str", TypescriptType: "string", ShowInTable: true, ShowInForm: true, ShowInDetail: true, Sort: 2},
	}
	plan := &Plan{
		Name: "plain", GenType: "TABLE", Author: "Charlie", MainTable: "cg_test_plain",
		MainPK: "id", MainEntityName: "Plain", MainModulePath: "biz/plain", MainBusinessName: "Plain",
		APIPrefix: "/biz/plain", PermissionPrefix: "biz:plain", MenuName: "Plain",
		MenuPath: "/biz/plain", ComponentPath: "biz/plain/index.vue", Sort: 99,
	}
	files, err := renderPlan(plan, main, nil)
	if err != nil {
		t.Fatalf("renderPlan failed: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no files rendered")
	}
}
