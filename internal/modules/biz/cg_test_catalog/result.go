// internal/modules/biz/cg_test_catalog/result.go 出参定义。
//
// Author: Charlie

package cg_test_catalog

// TreeNode 目录树节点。
//
// Author: Charlie
type TreeNode struct {
	Catalog
	Children []TreeNode `json:"children"`
}
