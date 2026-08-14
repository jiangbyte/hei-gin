// internal/modules/biz/cg_test_knowledge_category/result.go 出参定义。
//
// Author: Charlie

package cg_test_knowledge_category

// TreeNode 分类树节点。
//
// Author: Charlie
type TreeNode struct {
	Category
	Children []TreeNode `json:"children"`
}
