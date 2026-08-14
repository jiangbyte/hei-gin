package group

import (
	"hei-gin/internal/modules/iam/client"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/resource"
	"hei-gin/internal/modules/iam/role"
)

// OwnRoleResult ç”¨æˆ·ç»„å·²æ‹¥æœ‰è§’è‰²ç»“æžœã€‚
//
// Author: Charlie
type OwnRoleResult struct {
	ID      string      `json:"id"`
	Roles   []role.Role `json:"roles"`
	RoleIDs []string    `json:"role_ids"`
}

// OwnResourceResult ç”¨æˆ·ç»„å·²æ‹¥æœ‰ç®¡ç†ç«¯èµ„æºæŽˆæƒç»“æžœã€‚
//
// Author: Charlie
type OwnResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []resource.GrantModule       `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}

// OwnClientResourceResult ç”¨æˆ·ç»„å·²æ‹¥æœ‰å®¢æˆ·ç«¯èµ„æºæŽˆæƒç»“æžœã€‚
//
// Author: Charlie
type OwnClientResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []client.GrantModule         `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}
