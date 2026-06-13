package contracts

import "context"

type RealmID string

type PermissionAPI interface {
	GetPermissionList(ctx context.Context, realmID RealmID, userID string) ([]string, error)
	GetRoleList(ctx context.Context, realmID RealmID, userID string) ([]string, error)
	GetPermissionScopeMap(ctx context.Context, realmID RealmID, userID string) (map[string]ScopeInfo, error)
}

type ScopeInfo struct {
	GroupScope     string   `json:"group_scope"`
	OrgScope       string   `json:"org_scope"`
	CustomGroupIDs []string `json:"custom_group_ids"`
	CustomOrgIDs   []string `json:"custom_org_ids"`
}

type ScopeRow struct {
	PermissionCode string
	Scope          string
	CustomGroupIDs *string
	CustomOrgIDs   *string
}
