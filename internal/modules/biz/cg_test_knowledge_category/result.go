package cg_test_knowledge_category

// TreeNode 分类树节点。
//
// Author: Charlie
type TreeNode struct {
	Category
	Children []TreeNode `json:"children"`
}
