// internal/modules/iam/dept/result.go 出参定义。
//
// Author: Charlie

package dept

// TreeNode 部门树节点。
//
// Author: Charlie
type TreeNode struct {
	Dept
	Children []TreeNode `json:"children"`
}
