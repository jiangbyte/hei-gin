package dict

// TreeNode 字典树节点。
//
// Author: Charlie
type TreeNode struct {
	Dict
	Children []TreeNode `json:"children"`
}
