package client

// TreeNode 客户端资源树节点。
//
// Author: Charlie
type TreeNode struct {
	ClientResource
	Children []TreeNode `json:"children"`
}

// ModuleOption 客户端模块选择项。
//
// Author: Charlie
type ModuleOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
