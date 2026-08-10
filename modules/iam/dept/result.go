package dept

// TreeNode 部门树节点。
//
// Author: Charlie
type TreeNode struct {
	Dept
	Children []TreeNode `json:"children"`
}
