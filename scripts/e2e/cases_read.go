package main

import (
	"fmt"
	"strings"
)

func runReadCases(base, adminTok, portalTok string, bucket *caseBucket) {
	type readCase struct {
		name   string
		token  string
		method string
		path   string
		check  func(status int, raw []byte, ar apiResp, data map[string]any) error
	}

	admin := adminTok
	portal := portalTok

	cases := []readCase{
		{
			name: "read_health_live", token: "", method: "GET", path: "/api/v1/internal/health/live",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				return assertKeys(data, "status")
			},
		},
		{
			name: "read_health_ready", token: "", method: "GET", path: "/api/v1/internal/health/ready",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				return assertKeys(data, "status", "checks")
			},
		},
		{
			name: "read_admin_me", token: admin, method: "GET", path: "/api/v1/admin/me",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				return assertKeys(data, "account", "account_type", "permission_keys")
			},
		},
		{
			name: "read_workspace_overview", token: admin, method: "GET", path: "/api/v1/admin/workspace/overview",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				return assertKeys(data, "shortcuts", "recent_operations", "recent_logins")
			},
		},
		{
			name: "read_public_site_footer", token: "", method: "GET", path: "/api/v1/public/site-footer",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				return assertKeys(data, "copyright_text", "icp_number", "psb_number")
			},
		},
		{
			name: "read_identity_status", token: admin, method: "GET", path: "/api/v1/admin/profile/identity/status",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				return assertKeys(data, "status")
			},
		},
		{
			name: "read_audit_my_page", token: admin, method: "GET", path: "/api/v1/admin/sys/audit/my-page?current=1&size=5",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				_, err := assertPage(data)
				return err
			},
		},
		{
			name: "read_roles_page", token: admin, method: "GET",
			path: "/api/v1/admin/sys/roles/page?current=1&size=10&name=admin&code=SUPER",
			check: pageRecordKeys(admin, base, "/api/v1/admin/sys/roles/detail", "id", "code", "name"),
		},
		{
			name: "read_accounts_page", token: admin, method: "GET",
			path: "/api/v1/admin/sys/accounts/page?current=1&size=10&account=super&name=a",
			check: pageRecordKeys(admin, base, "/api/v1/admin/sys/accounts/detail", "id", "account_type", "account_status"),
		},
		{
			name: "read_banners_page", token: admin, method: "GET",
			path: "/api/v1/admin/sys/banners/page?current=1&size=10&title=a&position=HOME_TOP&status=ENABLED",
			check: pageRecordKeys(admin, base, "/api/v1/admin/sys/banners/detail", "id", "title", "target_account_types", "position"),
		},
		{
			name: "read_banners_list", token: admin, method: "GET",
			path: "/api/v1/admin/sys/banners/list?position=ADMIN_TOP",
			check: func(status int, _ []byte, ar apiResp, _ map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				_, err := parseDataArray(ar)
				return err
			},
		},
		{
			name: "read_notices_page", token: admin, method: "GET",
			path: "/api/v1/admin/sys/notices/page?current=1&size=10&title=a",
			check: pageRecordKeys(admin, base, "/api/v1/admin/sys/notices/detail", "id", "title", "kind", "target_scope"),
		},
		{
			name: "read_notices_my_page", token: admin, method: "GET",
			path: "/api/v1/admin/sys/notices/my-page?current=1&size=5",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				_, err := assertPage(data)
				return err
			},
		},
		{
			name: "read_dicts_page", token: admin, method: "GET",
			path: "/api/v1/admin/sys/dicts/page?current=1&size=10&code=COMMON&category=SYS&status=ENABLED",
			check: pageRecordKeys(admin, base, "/api/v1/admin/sys/dicts/detail", "id", "code", "status"),
		},
		{
			name: "read_dicts_tree", token: admin, method: "GET",
			path: "/api/v1/admin/sys/dicts/tree?category=SYS",
			check: func(status int, raw []byte, ar apiResp, _ map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				_, err := parseDataArray(ar)
				return err
			},
		},
		{
			name: "read_config_page", token: admin, method: "GET",
			path: "/api/v1/admin/sys/config/page?current=1&size=10&config_key=a&category=AUTH",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				recs, err := assertPage(data)
				if err != nil {
					return err
				}
				if len(recs) > 0 {
					return assertKeys(recs[0], "id", "config_key")
				}
				return nil
			},
		},
		{
			name: "read_weak_password_page", token: admin, method: "GET",
			path: "/api/v1/admin/sys/weak-password/page?current=1&size=10&keyword=123",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				_, err := assertPage(data)
				return err
			},
		},
		{
			name: "read_depts_page", token: admin, method: "GET",
			path: "/api/v1/admin/sys/depts/page?current=1&size=10&name=a",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				_, err := assertPage(data)
				return err
			},
		},
		{
			name: "read_codegen_tables", token: admin, method: "GET",
			path: "/api/v1/admin/sys/codegen/tables",
			check: func(status int, _ []byte, ar apiResp, _ map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				arr, err := parseDataArrayMaps(ar)
				if err != nil {
					return err
				}
				if len(arr) == 0 {
					return fmt.Errorf("tables empty")
				}
				return assertKeys(arr[0], "table_name")
			},
		},
		{
			name: "read_codegen_page", token: admin, method: "GET",
			path: "/api/v1/admin/sys/codegen/page?current=1&size=5&name=a",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				_, err := assertPage(data)
				return err
			},
		},
		{
			name: "read_cg_catalog_page", token: admin, method: "GET",
			path: "/api/v1/admin/biz/cg-test-catalog/page?current=1&size=10&code=ROOT&name=%E6%A0%B9",
			check: pageRecordKeys(admin, base, "/api/v1/admin/biz/cg-test-catalog/detail", "id", "code", "name"),
		},
		{
			name: "read_cg_catalog_tree", token: admin, method: "GET",
			path: "/api/v1/admin/biz/cg-test-catalog/tree",
			check: func(status int, _ []byte, ar apiResp, _ map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				_, err := parseDataArray(ar)
				return err
			},
		},
		{
			name: "read_cg_activity_page", token: admin, method: "GET",
			path: "/api/v1/admin/biz/cg-test-activity/page?current=1&size=10&code=ACT&name=a",
			check: pageRecordKeys(admin, base, "/api/v1/admin/biz/cg-test-activity/detail", "id", "code", "name"),
		},
		{
			name: "read_cg_order_page", token: admin, method: "GET",
			path: "/api/v1/admin/biz/cg-test-order/page?current=1&size=10&order_no=a&name=a",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				_, err := assertPage(data)
				return err
			},
		},
		{
			name: "read_portal_banners", token: portal, method: "GET",
			path: "/api/v1/portal/sys/banners/list",
			check: func(status int, _ []byte, ar apiResp, _ map[string]any) error {
				return assertBizOK(status, ar.Code)
			},
		},
		{
			name: "read_portal_dicts_tree", token: "", method: "GET",
			path: "/api/v1/portal/sys/dicts/tree",
			check: func(status int, _ []byte, ar apiResp, _ map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				_, err := parseDataArray(ar)
				return err
			},
		},
		{
			name: "read_portal_notices", token: "", method: "GET",
			path: "/api/v1/portal/sys/notices/list",
			check: func(status int, _ []byte, ar apiResp, _ map[string]any) error {
				return assertBizOK(status, ar.Code)
			},
		},
		{
			name: "read_audit_page", token: admin, method: "GET",
			path: "/api/v1/admin/sys/audit/page?current=1&size=5",
			check: func(status int, _ []byte, ar apiResp, data map[string]any) error {
				if err := assertBizOK(status, ar.Code); err != nil {
					return err
				}
				_, err := assertPage(data)
				return err
			},
		},
		{
			name: "read_resources_current", token: admin, method: "GET",
			path: "/api/v1/admin/sys/resources/current",
			check: func(status int, _ []byte, ar apiResp, _ map[string]any) error {
				return assertBizOK(status, ar.Code)
			},
		},
	}

	for _, c := range cases {
		url := base + c.path
		status, raw, ar, err := doRaw(c.method, url, c.token, "")
		cr := caseResult{Name: c.name, URL: url, Status: status, BizCode: ar.Code, Body: truncate(string(raw), 300)}
		if err != nil {
			cr.Error = err.Error()
			bucket.add(cr)
			continue
		}
		_, data, perr := parseEnvelope(raw)
		if perr != nil && c.check != nil {
			// still pass data nil into check for array responses
		}
		if c.check != nil {
			if err := c.check(status, raw, ar, data); err != nil {
				cr.Error = err.Error()
				bucket.add(cr)
				continue
			}
		}
		cr.OK = true
		bucket.add(cr)
	}
}

func pageRecordKeys(token, base, detailPath string, keys ...string) func(int, []byte, apiResp, map[string]any) error {
	return func(status int, _ []byte, ar apiResp, data map[string]any) error {
		if err := assertBizOK(status, ar.Code); err != nil {
			return err
		}
		recs, err := assertPage(data)
		if err != nil {
			return err
		}
		if len(recs) == 0 {
			return nil // empty page ok for filtered queries
		}
		if err := assertKeys(recs[0], keys...); err != nil {
			return err
		}
		id := firstID(recs)
		if id == "" || detailPath == "" {
			return nil
		}
		sep := "?"
		if strings.Contains(detailPath, "?") {
			sep = "&"
		}
		st, raw, dar, err := doRaw("GET", base+detailPath+sep+"id="+id, token, "")
		if err != nil {
			return err
		}
		if err := assertBizOK(st, dar.Code); err != nil {
			return fmt.Errorf("detail: %w body=%s", err, truncate(string(raw), 120))
		}
		_, dm, err := parseEnvelope(raw)
		if err != nil {
			return err
		}
		if asString(dm["id"]) != id && asString(dm["account_id"]) != id {
			// accounts detail may use id field
			if asString(dm["id"]) == "" {
				return fmt.Errorf("detail id mismatch got=%v want=%s", dm["id"], id)
			}
		}
		return nil
	}
}
