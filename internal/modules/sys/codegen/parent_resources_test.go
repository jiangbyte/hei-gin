package codegen

import "testing"

func TestNormalizeResourceParents(t *testing.T) {
	root := "root"
	orphan := "missing"
	rows := []ResourceNode{
		{ID: "root", ParentID: nil},
		{ID: "child", ParentID: &root},
		{ID: "orphan", ParentID: &orphan},
	}
	normalizeResourceParents(rows)
	if rows[1].ParentID == nil || *rows[1].ParentID != root {
		t.Fatalf("valid parent should remain: %+v", rows[1])
	}
	if rows[2].ParentID != nil {
		t.Fatalf("orphan parent should be cleared: %+v", rows[2])
	}
}

func TestBuildResourceTree(t *testing.T) {
	pRoot := "root"
	rows := []ResourceNode{
		{ID: "root", Name: "Root"},
		{ID: "a", ParentID: &pRoot, Name: "A"},
		{ID: "b", ParentID: &pRoot, Name: "B"},
	}
	tree := buildResourceTree(rows, nil)
	if len(tree) != 1 || tree[0].ID != "root" {
		t.Fatalf("expected single root, got %+v", tree)
	}
	if len(tree[0].Children) != 2 {
		t.Fatalf("expected 2 children, got %+v", tree[0].Children)
	}
}
