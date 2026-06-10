package role

import (
	"hei-gin/sdk/utils"
)

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

// ButtonPermissionScope represents a button permission with data scope.
type ButtonPermissionScope struct {
	PermissionCode      string  `json:"permission_code"`
	Scope               string  `json:"scope"`
	CustomScopeGroupIds *string `json:"custom_scope_group_ids"`
	CustomScopeOrgIds   *string `json:"custom_scope_org_ids"`
}

// GrantResourceParam holds the parameters for granting resources to a role.
type GrantResourceParam struct {
	RoleID      string                  `json:"role_id"`
	ResourceIDs []string                `json:"resource_ids"`
	Permissions []ButtonPermissionScope `json:"permissions"`
}


func toVO(entity *SysRole) *RoleVO {
	if entity == nil { return nil }
	return &RoleVO{
		ID: entity.ID, Code: entity.Code, Name: entity.Name, Category: entity.Category,
		Description: entity.Description, Status: entity.Status, SortCode: entity.SortCode,
		Extra: entity.Extra, CreatedAt: utils.FormatDateTimePtr(entity.CreatedAt),
		CreatedBy: entity.CreatedBy, UpdatedAt: utils.FormatDateTimePtr(entity.UpdatedAt),
		UpdatedBy: entity.UpdatedBy,
	}
}


func (p *RolePageParam) GetCurrent() int { return p.Current }
func (p *RolePageParam) GetSize() int    { return p.Size }
