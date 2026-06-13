package auth

import (
	"testing"
)

func TestMergeScopeCombinesRoleAndDirectScope(t *testing.T) {
	roleGroups := `["g1"]`
	roleOrgs := `["o1"]`
	directGroups := `["g2"]`
	directOrgs := `["o2"]`

	permScope := make(map[string]map[string]interface{})
	MergeScope(permScope, permissionPathUserRole, []ScopeRow{
		{
			PermissionCode: "sys:user:view",
			Scope:          "ORG",
			CustomGroupIDs: &roleGroups,
			CustomOrgIDs:   &roleOrgs,
		},
	})
	MergeScope(permScope, permissionPathDirect, []ScopeRow{
		{
			PermissionCode: "sys:user:view",
			Scope:          "SELF",
			CustomGroupIDs: &directGroups,
			CustomOrgIDs:   &directOrgs,
		},
	})

	got := permScope["sys:user:view"]
	if got["group_scope"] != "SELF" {
		t.Fatalf("group_scope = %v", got["group_scope"])
	}
	if got["org_scope"] != "SELF" {
		t.Fatalf("org_scope = %v", got["org_scope"])
	}

	groups := safeStringSlice(got["custom_group_ids"])
	orgs := safeStringSlice(got["custom_org_ids"])
	if len(groups) != 2 || groups[0] != "g1" || groups[1] != "g2" {
		t.Fatalf("custom_group_ids = %#v", groups)
	}
	if len(orgs) != 2 || orgs[0] != "o1" || orgs[1] != "o2" {
		t.Fatalf("custom_org_ids = %#v", orgs)
	}
}

func TestParseCSVFallbacksToRawString(t *testing.T) {
	raw := "dept-1"
	result := parseCSV(&raw)
	if len(result) != 1 || result[0] != "dept-1" {
		t.Fatalf("result = %#v", result)
	}
}
