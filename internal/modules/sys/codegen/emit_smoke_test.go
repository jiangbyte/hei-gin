package codegen

import "testing"

func strptr(s string) *string { return &s }

func TestRenderPlanSmoke(t *testing.T) {
	main := []Field{
		{ColumnName: "id", DBType: "varchar", ValueType: "str", UIType: "string", PrimaryKey: true, InTable: true, InForm: false, InDetail: true, Sort: 1},
		{ColumnName: "name", Label: strptr("名称"), DBType: "varchar", ValueType: "str", UIType: "string", InTable: true, InForm: true, InDetail: true, Sort: 2},
		{ColumnName: "status", Label: strptr("状态"), DBType: "varchar", ValueType: "str", UIType: "string", Widget: "dict", DictCode: strptr("COMMON_STATUS"), InTable: true, InForm: true, InDetail: true, InQuery: true, Sort: 3},
		{ColumnName: "count", DBType: "int4", ValueType: "int", UIType: "number", Widget: "number", InTable: true, InForm: true, InDetail: true, Sort: 4},
		{ColumnName: "created_at", DBType: "timestamptz", ValueType: "datetime", UIType: "string", InTable: true, InForm: false, InDetail: true, Sort: 5},
	}
	plan := &Plan{
		Name: "smoke", GenType: "TABLE", Author: "Charlie", Table: "cg_test_smoke",
		PKColumn: "id", EntityName: "Smoke", ModulePath: "biz/smoke", BusinessName: "Smoke",
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
}

func TestRenderPlanPlain(t *testing.T) {
	main := []Field{
		{ColumnName: "id", DBType: "varchar", ValueType: "str", UIType: "string", PrimaryKey: true, InTable: true, InForm: false, InDetail: true, Sort: 1},
		{ColumnName: "name", DBType: "varchar", ValueType: "str", UIType: "string", InTable: true, InForm: true, InDetail: true, Sort: 2},
	}
	plan := &Plan{
		Name: "plain", GenType: "TABLE", Author: "Charlie", Table: "cg_test_plain",
		PKColumn: "id", EntityName: "Plain", ModulePath: "biz/plain", BusinessName: "Plain",
		APIPrefix: "/biz/plain", PermissionPrefix: "biz:plain", MenuName: "Plain",
		MenuPath: "/biz/plain", ComponentPath: "biz/plain/index.vue", Sort: 99,
	}
	if _, err := renderPlan(plan, main, nil); err != nil {
		t.Fatalf("renderPlan failed: %v", err)
	}
}
