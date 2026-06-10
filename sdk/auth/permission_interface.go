package auth

import (
	"encoding/json"
	"log"

	"hei-gin/api"
	"hei-gin/sdk/enums"
)

var PermissionDelegate api.PermissionAPI

func RegisterInterface(impl api.PermissionAPI) {
	PermissionDelegate = impl
	log.Println("[auth] PermissionInterface registered via delegate")
}

// ScopeRow is an alias for api.ScopeRow to avoid import conflicts.
type ScopeRow = api.ScopeRow

// ScopeInfo is an alias for api.ScopeInfo.
type ScopeInfo = api.ScopeInfo

func MergeScope(permScope map[string]map[string]interface{}, path string, rows []ScopeRow) {
	priority := map[string]int{
		string(enums.PermissionPathDirect):   0,
		string(enums.PermissionPathUserRole): 1,
	}
	currentPrio, _ := priority[path]
	for _, row := range rows {
		existing, exists := permScope[row.PermissionCode]
		if !exists {
			permScope[row.PermissionCode] = map[string]interface{}{
				"group_scope":      row.Scope,
				"org_scope":        row.Scope,
				"custom_group_ids": parseCSV(row.CustomGroupIDs),
				"custom_org_ids":   parseCSV(row.CustomOrgIDs),
			}
			continue
		}
		permScope[row.PermissionCode] = mergeScopeEntry(existing, row, currentPrio)
	}
}

func mergeScopeEntry(existing map[string]interface{}, row ScopeRow, currentPrio int) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range existing {
		result[k] = v
	}
	result["group_scope"] = enums.MostRestrictive(
		safeString(result["group_scope"]),
		row.Scope,
	)
	prevCustomGroups := safeStringSlice(result["custom_group_ids"])
	prevCustomOrgs := safeStringSlice(result["custom_org_ids"])
	newCustomGroups := parseCSV(row.CustomGroupIDs)
	newCustomOrgs := parseCSV(row.CustomOrgIDs)
	result["custom_group_ids"] = mergeUnique(prevCustomGroups, newCustomGroups)
	result["custom_org_ids"] = mergeUnique(prevCustomOrgs, newCustomOrgs)
	result["org_scope"] = enums.MostRestrictive(
		safeString(result["org_scope"]),
		row.Scope,
	)
	return result
}

func parseCSV(s *string) []string {
	if s == nil || *s == "" {
		return nil
	}
	var result []string
	if err := json.Unmarshal([]byte(*s), &result); err != nil {
		return []string{*s}
	}
	return result
}

func mergeUnique(a, b []string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	result := make([]string, 0, len(a)+len(b))
	for _, s := range a {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	for _, s := range b {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
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

func safeStrPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

