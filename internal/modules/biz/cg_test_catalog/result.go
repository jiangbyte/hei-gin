package cg_test_catalog

// TreeNode 目录树节点。
//
// Author: Charlie
type TreeNode struct {
	Catalog
	Children []TreeNode `json:"children"`
}
