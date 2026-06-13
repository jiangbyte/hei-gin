package provider

import (
	"context"
	"encoding/json"
	"log"
	"sort"

	"gorm.io/gorm"

	roleModel "hei-gin/plugins/plugin-sys/role"
	userModel "hei-gin/plugins/plugin-sys/user"
	"hei-gin/sdk/auth"
	"hei-gin/sdk/constants"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/shared/contracts"
)

type PermissionProvider struct{}

// getRoleIDs returns the role IDs for a given login ID (single query).
func (p *PermissionProvider) getRoleIDs(ctx context.Context, loginID string) ([]string, error) {
	var entities []userModel.RelUserRole
	if err := db.DB.WithContext(ctx).Where("user_id = ?", loginID).Find(&entities).Error; err != nil {
		log.Printf("[Permission] Failed to query user roles: %v", err)
		return nil, err
	}
	if len(entities) == 0 {
		return nil, nil
	}
	roleIDs := make([]string, 0, len(entities))
	for _, e := range entities {
		roleIDs = append(roleIDs, e.RoleID)
	}
	return roleIDs, nil
}

// getRolesByIDs returns full role records for a batch of role IDs (single query).
func (p *PermissionProvider) getRolesByIDs(ctx context.Context, roleIDs []string) ([]roleModel.SysRole, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var roles []roleModel.SysRole
	if err := db.DB.WithContext(ctx).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (p *PermissionProvider) GetPermissionList(ctx context.Context, realmID contracts.RealmID, userID string) ([]string, error) {
	roleIDs, err := p.getRoleIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissionCodes := make(map[string]struct{})

	if realmID == contracts.RealmID(auth.BusinessID) || realmID == contracts.RealmID(auth.ConsumerID) {
		if len(roleIDs) > 0 {
			var entities []userModel.RelRolePermission
			if err := db.DB.WithContext(ctx).Where("role_id IN ?", roleIDs).Find(&entities).Error; err != nil {
				log.Printf("[Permission] Failed to query role permissions: %v", err)
			} else {
				for _, e := range entities {
					permissionCodes[e.PermissionCode] = struct{}{}
				}
			}
		}

		var entities []userModel.RelUserPermission
		if err := db.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&entities).Error; err != nil {
			log.Printf("[Permission] Failed to query user permissions: %v", err)
		} else {
			for _, e := range entities {
				permissionCodes[e.PermissionCode] = struct{}{}
			}
		}
	}

	result := make([]string, 0, len(permissionCodes))
	for code := range permissionCodes {
		result = append(result, code)
	}
	return result, nil
}

func (p *PermissionProvider) GetRoleList(ctx context.Context, realmID contracts.RealmID, userID string) ([]string, error) {
	_ = realmID
	roleIDs, err := p.getRoleIDs(ctx, userID)
	if err != nil || len(roleIDs) == 0 {
		return []string{}, err
	}

	roles, err := p.getRolesByIDs(ctx, roleIDs)
	if err != nil {
		log.Printf("[Permission] Failed to query role list: %v", err)
		return nil, err
	}

	roleCodes := make([]string, 0, len(roles))
	for _, r := range roles {
		roleCodes = append(roleCodes, r.Code)
	}
	return roleCodes, nil
}

func (p *PermissionProvider) GetPermissionScopeMap(ctx context.Context, realmID contracts.RealmID, userID string) (map[string]contracts.ScopeInfo, error) {
	if realmID != contracts.RealmID(auth.BusinessID) && realmID != contracts.RealmID(auth.ConsumerID) {
		return map[string]contracts.ScopeInfo{}, nil
	}

	roleIDs, err := p.getRoleIDs(ctx, userID)
	if err != nil {
		log.Printf("[Permission] Failed to query user roles: %v", err)
		roleIDs = nil
	}

	permScope := make(map[string]map[string]interface{})

	if len(roleIDs) > 0 {
		var entities []userModel.RelRolePermission
		if err := db.DB.WithContext(ctx).Where("role_id IN ?", roleIDs).Find(&entities).Error; err != nil {
			log.Printf("[Permission] Failed to query role permission scopes: %v", err)
		} else {
			scopeRows := make([]auth.ScopeRow, 0, len(entities))
			for _, e := range entities {
				scopeRows = append(scopeRows, auth.ScopeRow{
					PermissionCode: e.PermissionCode,
					Scope:          e.Scope,
					CustomGroupIDs: e.CustomScopeGroupIds,
					CustomOrgIDs:   e.CustomScopeOrgIds,
				})
			}
			auth.MergeScope(permScope, string(enums.PermissionPathUserRole), scopeRows)
		}
	}

	var entities []userModel.RelUserPermission
	if err := db.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&entities).Error; err != nil {
		log.Printf("[Permission] Failed to query user permission scopes: %v", err)
	} else {
		scopeRows := make([]auth.ScopeRow, 0, len(entities))
		for _, e := range entities {
			scopeRows = append(scopeRows, auth.ScopeRow{
				PermissionCode: e.PermissionCode,
				Scope:          e.Scope,
				CustomGroupIDs: e.CustomScopeGroupIds,
				CustomOrgIDs:   e.CustomScopeOrgIds,
			})
		}
		auth.MergeScope(permScope, string(enums.PermissionPathDirect), scopeRows)
	}

	result := make(map[string]contracts.ScopeInfo, len(permScope))
	for k, v := range permScope {
		result[k] = contracts.ScopeInfo{
			GroupScope:     safeString(v["group_scope"]),
			OrgScope:       safeString(v["org_scope"]),
			CustomGroupIDs: safeStringSlice(v["custom_group_ids"]),
			CustomOrgIDs:   safeStringSlice(v["custom_org_ids"]),
		}
	}
	return result, nil
}

func (p *PermissionProvider) getAllPermissionsFromRedis(ctx context.Context) ([]string, error) {
	val, err := db.Redis.Get(ctx, constants.PERMISSION_CACHE_KEY).Result()
	if err != nil {
		return nil, err
	}

	var perms []string
	if err := json.Unmarshal([]byte(val), &perms); err == nil {
		return perms, nil
	}

	var tree map[string]map[string]auth.PermissionEntry
	if err := json.Unmarshal([]byte(val), &tree); err != nil {
		return nil, err
	}

	codeSet := make(map[string]struct{})
	for _, modulePerms := range tree {
		for code := range modulePerms {
			codeSet[code] = struct{}{}
		}
	}

	perms = make([]string, 0, len(codeSet))
	for code := range codeSet {
		perms = append(perms, code)
	}
	sort.Strings(perms)
	return perms, nil
}

func safeString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func safeStringSlice(v interface{}) []string {
	if v == nil {
		return []string{}
	}
	if s, ok := v.([]string); ok {
		return s
	}
	return []string{}
}

var _ = &gorm.DB{}
