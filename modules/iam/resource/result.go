package resource

// TreeNode 资源树节点。
//
// Author: Charlie
type TreeNode struct {
	Resource
	Children []TreeNode `json:"children"`
}

// ModuleOption 资源模块选择项。
//
// Author: Charlie
type ModuleOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
