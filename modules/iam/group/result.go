package group

import (
	"hei-gin/modules/iam/client"
	"hei-gin/modules/iam/relation"
	"hei-gin/modules/iam/resource"
	"hei-gin/modules/iam/role"
)

// OwnRoleResult 用户组已拥有角色结果。
//
// Author: Charlie
type OwnRoleResult struct {
	ID      string      `json:"id"`
	Roles   []role.Role `json:"roles"`
	RoleIDs []string    `json:"role_ids"`
}

// OwnResourceResult 用户组已拥有管理端资源授权结果。
//
// Author: Charlie
type OwnResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []resource.GrantModule       `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}

// OwnClientResourceResult 用户组已拥有客户端资源授权结果。
//
// Author: Charlie
type OwnClientResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []client.GrantModule         `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}
