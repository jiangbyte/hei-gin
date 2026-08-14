// internal/modules/sys/dict/result.go 出参定义。
//
// Author: Charlie

package dict

// TreeNode 字典树节点。
//
// Author: Charlie
type TreeNode struct {
	Dict
	Children []TreeNode `json:"children"`
}
