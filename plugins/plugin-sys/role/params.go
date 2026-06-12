package role

import "hei-gin/sdk/utils"

// RoleVO is the view object for a role, used for create/modify requests and API responses.
type RoleVO struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
	SortCode    int     `json:"sort_code"`
	Extra       *string `json:"extra"`
	CreatedAt   string  `json:"created_at"`
	CreatedBy   *string `json:"created_by"`
	UpdatedAt   string  `json:"updated_at"`
	UpdatedBy   *string `json:"updated_by"`
}

// RolePageParam holds pagination parameters for the role page query.
type RolePageParam struct {
	Current  int    `json:"current" form:"current"`
	Size     int    `json:"size" form:"size"`
	Keyword  string `json:"keyword" form:"keyword"`
	Category string `json:"category" form:"category"`
}

// PermissionItem represents a permission to be granted to a role.
type PermissionItem struct {
	PermissionCode      string  `json:"permission_code"`
	Scope               string  `json:"scope"`
	CustomScopeGroupIds *string `json:"custom_scope_group_ids"`
	CustomScopeOrgIds   *string `json:"custom_scope_org_ids"`
}

// GrantPermissionParam holds the parameters for granting permissions to a role.
type GrantPermissionParam struct {
	RoleID      string           `json:"role_id"`
	Permissions []PermissionItem `json:"permissions"`
}

// GrantResourceParam holds the parameters for granting resources to a role.
type GrantResourceParam struct {
	RoleID      string           `json:"role_id"`
	ResourceIDs []string         `json:"resource_ids"`
	Permissions []PermissionItem `json:"permissions"`
}

func SysRoleToRoleVO(src *SysRole) *RoleVO {
	if src == nil {
		return nil
	}

	dst := &RoleVO{}
	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.Description = src.Description
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.Extra = src.Extra
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)
	return dst
}

func RoleVOToSysRole(src *RoleVO) *SysRole {
	if src == nil {
		return nil
	}

	dst := &SysRole{}
	dst.ID = src.ID
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Category = src.Category
	dst.Description = src.Description
	dst.Status = src.Status
	dst.SortCode = src.SortCode
	dst.Extra = src.Extra
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)
	return dst
}
