package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// crudSpecs 覆盖全部标准 /create 资源；缺工厂会在 runAllCRUD 中判失败。
func crudSpecs(suffix string) map[string]crudSpec {
	now := time.Now().UTC().Format(time.RFC3339)
	later := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
	enabled := false

	cloneUpdate := func(mutate func(m map[string]any)) func(string, map[string]any) map[string]any {
		return func(id string, create map[string]any) map[string]any {
			m := map[string]any{"id": id}
			for k, v := range create {
				m[k] = v
			}
			mutate(m)
			return m
		}
	}

	return map[string]crudSpec{
		"/api/v1/admin/sys/dicts/create": {
			findField: "code", findValue: "E2E_DICT_" + suffix,
			create: map[string]any{
				"code": "E2E_DICT_" + suffix, "label": "e2e-dict", "value": "E2E_DICT_" + suffix,
				"category": "BIZ", "status": "ENABLED", "sort": 99,
			},
			update: cloneUpdate(func(m map[string]any) { m["label"] = "e2e-dict-upd"; m["sort"] = 98 }),
		},
		"/api/v1/admin/sys/weak-password/create": {
			findField: "password", findValue: "e2e_weak_" + suffix,
			create: map[string]any{"password": "e2e_weak_" + suffix},
			update: cloneUpdate(func(m map[string]any) { m["password"] = "e2e_weak_" + suffix + "_u" }),
			pageQuery: func(v string) string {
				return "current=1&size=50&keyword=" + url.QueryEscape(v)
			},
		},
		"/api/v1/admin/sys/banners/create": {
			findField: "title", findValue: "E2E Banner " + suffix,
			create: map[string]any{
				"title": "E2E Banner " + suffix, "image": "https://example.com/e2e.png", "link_type": "NONE",
				"category": "HOME", "type": "CAROUSEL", "position": "HOME_TOP",
				"sort": 99, "status": "DISABLED",
			},
			update: cloneUpdate(func(m map[string]any) {
				m["title"] = "E2E Banner " + suffix + " upd"
				m["image"] = "https://example.com/e2e2.png"
				m["sort"] = 98
			}),
			pageQuery: func(string) string { return "current=1&size=50&title=E2E" },
		},
		"/api/v1/admin/sys/notices/create": {
			findField: "title", findValue: "E2E Notice " + suffix,
			create: map[string]any{
				"kind": "NOTIFICATION", "title": "E2E Notice " + suffix, "content": "e2e content", "content_type": "text",
				"category": "SYSTEM", "severity": "INFO", "target_scope": "ALL",
				"target_account_types": []string{"ADMIN"}, "status": "DRAFT",
				"publish_locations": map[string]any{"center": true},
			},
			update: cloneUpdate(func(m map[string]any) {
				m["title"] = "E2E Notice " + suffix + " upd"
				m["content"] = "e2e content2"
			}),
			pageQuery: func(string) string { return "current=1&size=50&title=E2E" },
		},
		"/api/v1/admin/sys/config/create": {
			findField: "config_key", findValue: "e2e.config." + suffix,
			create: map[string]any{
				"config_key": "e2e.config." + suffix, "config_value": "v1", "category": "E2E",
				"value_type": "STRING", "label": "e2e", "sort_code": 99,
			},
			update: cloneUpdate(func(m map[string]any) { m["config_value"] = "v2"; m["sort_code"] = 98 }),
		},
		"/api/v1/admin/sys/roles/create": {
			findField: "code", findValue: "E2E_ROLE_" + suffix,
			create: map[string]any{
				"code": "E2E_ROLE_" + suffix, "name": "E2E Role " + suffix, "category": "CUSTOM",
				"scope_type": "ALL", "sort": 99, "status": "ENABLED",
			},
			update: cloneUpdate(func(m map[string]any) { m["name"] = "E2E Role " + suffix + " upd"; m["sort"] = 98 }),
		},
		"/api/v1/admin/sys/positions/create": {
			findField: "name", findValue: "E2E Pos " + suffix,
			create: map[string]any{
				"name": "E2E Pos " + suffix, "category": "STAFF", "sort": 99, "status": "ENABLED",
			},
			update: cloneUpdate(func(m map[string]any) { m["name"] = "E2E Pos " + suffix + " upd"; m["sort"] = 98 }),
			pageQuery: func(v string) string { return "current=1&size=50&name=" + url.QueryEscape(v) },
		},
		"/api/v1/admin/sys/groups/create": {
			findField: "name", findValue: "E2E Group " + suffix,
			create: map[string]any{"name": "E2E Group " + suffix, "status": "ENABLED"},
			update: cloneUpdate(func(m map[string]any) { m["name"] = "E2E Group " + suffix + " upd" }),
			pageQuery: func(v string) string { return "current=1&size=50&name=" + url.QueryEscape(v) },
		},
		"/api/v1/admin/sys/depts/create": {
			findField: "name", findValue: "E2E Dept " + suffix,
			create: map[string]any{
				"name": "E2E Dept " + suffix, "category": "DEPT", "sort": 99, "status": "ENABLED",
			},
			update: cloneUpdate(func(m map[string]any) { m["name"] = "E2E Dept " + suffix + " upd"; m["sort"] = 98 }),
			pageQuery: func(v string) string { return "current=1&size=50&name=" + url.QueryEscape(v) },
		},
		"/api/v1/admin/sys/resource-modules/create": {
			findField: "code", findValue: "E2E_RMOD_" + suffix,
			create: map[string]any{
				"name": "E2E ResMod " + suffix, "code": "E2E_RMOD_" + suffix, "client": "ADMIN",
				"sort": 99, "status": "ENABLED",
			},
			update: cloneUpdate(func(m map[string]any) { m["name"] = "E2E ResMod " + suffix + " upd"; m["sort"] = 98 }),
		},
		"/api/v1/admin/sys/client-modules/create": {
			findField: "code", findValue: "E2E_CMOD_" + suffix,
			create: map[string]any{
				"name": "E2E CliMod " + suffix, "code": "E2E_CMOD_" + suffix, "account_type": "PORTAL",
				"sort": 99, "status": "ENABLED",
			},
			update: cloneUpdate(func(m map[string]any) { m["name"] = "E2E CliMod " + suffix + " upd"; m["sort"] = 98 }),
		},
		"/api/v1/admin/sys/resources/create": {
			findField: "code", findValue: "E2E_RES_" + suffix,
			create: map[string]any{
				"code": "E2E_RES_" + suffix, "name": "E2E Res " + suffix, "resource_type": "MENU",
				"path": "/e2e/" + suffix, "component": "e2e/Page", "sort": 99, "status": "ENABLED",
				"is_visible": true,
			},
			update: cloneUpdate(func(m map[string]any) { m["name"] = "E2E Res " + suffix + " upd"; m["sort"] = 98 }),
		},
		"/api/v1/admin/sys/client-resources/create": {
			findField: "code", findValue: "E2E_CRES_" + suffix,
			create: map[string]any{
				"code": "E2E_CRES_" + suffix, "name": "E2E CliRes " + suffix, "resource_type": "MENU",
				"path": "/e2e-c/" + suffix, "component": "e2e/CPage", "sort": 99, "status": "ENABLED",
				"is_visible": true,
			},
			update: cloneUpdate(func(m map[string]any) { m["name"] = "E2E CliRes " + suffix + " upd"; m["sort"] = 98 }),
		},
		"/api/v1/admin/sys/jobs/create": {
			findField: "name", findValue: "E2E Job " + suffix,
			create: map[string]any{
				"name": "E2E Job " + suffix, "handler": "sys_job_sample",
				"trigger_type": "FIXED", "trigger_config": "3600",
				"params": map[string]any{}, "sort": 99, "enabled": &enabled,
			},
			update: cloneUpdate(func(m map[string]any) {
				m["name"] = "E2E Job " + suffix + " upd"
				m["trigger_config"] = "7200"
				m["sort"] = 98
			}),
			pageQuery: func(v string) string { return "current=1&size=50&name=" + url.QueryEscape(v) },
		},
		"/api/v1/admin/biz/cg-test-catalog/create": {
			findField: "code", findValue: "E2E_CAT_" + suffix,
			create: map[string]any{
				"code": "E2E_CAT_" + suffix, "name": "E2E Catalog " + suffix, "category": "BUSINESS",
				"status": "ENABLED", "sort": 99, "is_visible": true, "extra": map[string]any{"e2e": true},
			},
			update: cloneUpdate(func(m map[string]any) {
				m["name"] = "E2E Catalog " + suffix + " upd"
				m["sort"] = 98
				m["is_visible"] = false
			}),
		},
		"/api/v1/admin/biz/cg-test-activity/create": {
			findField: "code", findValue: "E2E_ACT_" + suffix,
			create: map[string]any{
				"code": "E2E_ACT_" + suffix, "name": "E2E Act " + suffix, "type": "ONLINE",
				"status": "DRAFT", "start_at": now, "end_at": later, "is_public": true,
				"need_approval": false, "price": 0, "max_participants": 10,
			},
			update: cloneUpdate(func(m map[string]any) { m["name"] = "E2E Act " + suffix + " upd" }),
		},
		"/api/v1/admin/biz/cg-test-order/create": {
			findField: "order_no", findValue: "E2E_ORD_" + suffix,
			create: map[string]any{
				"order_no": "E2E_ORD_" + suffix, "name": "E2E Order " + suffix,
				"customer_name": "e2e", "status": "CREATED", "type": "NORMAL",
				"ordered_at": now, "total_amount": 1.23, "item_count": 1, "need_invoice": false,
			},
			update: cloneUpdate(func(m map[string]any) { m["name"] = "E2E Order " + suffix + " upd"; m["item_count"] = 2 }),
		},
		"/api/v1/admin/biz/cg-test-knowledge-category/create": {
			findField: "code", findValue: "E2E_KC_" + suffix,
			create: map[string]any{
				"code": "E2E_KC_" + suffix, "name": "E2E KC " + suffix, "status": "ENABLED",
				"sort": 99, "is_visible": true,
			},
			update: cloneUpdate(func(m map[string]any) { m["name"] = "E2E KC " + suffix + " upd"; m["sort"] = 98 }),
		},
	}
}

// attachDeferredCRUDSpecs 依赖运行时种子（parent_id / password key / table）的 create。
func attachDeferredCRUDSpecs(base, adminTok, suffix string, specs map[string]crudSpec) error {
	// resource-buttons
	parentID := firstPageID(base, adminTok, "/api/v1/admin/sys/resources/page")
	specs["/api/v1/admin/sys/resource-buttons/create"] = crudSpec{
		findField: "code", findValue: "E2E_BTN_" + suffix,
		create: map[string]any{
			"parent_id": parentID, "code": "E2E_BTN_" + suffix, "name": "E2E Btn " + suffix,
			"permission_key": "e2e:btn:" + suffix, "sort": 99, "status": "ENABLED",
		},
		update: func(id string, create map[string]any) map[string]any {
			m := map[string]any{"id": id}
			for k, v := range create {
				m[k] = v
			}
			m["name"] = "E2E Btn " + suffix + " upd"
			return m
		},
	}

	// order children（明细）依赖订单；若库空则先建临时订单
	orderID := firstPageID(base, adminTok, "/api/v1/admin/biz/cg-test-order/page")
	if orderID == "" || orderID == "0" {
		ono := "E2E_ORD_TMP_" + suffix
		body, _ := json.Marshal(map[string]any{
			"order_no": ono, "name": "tmp", "customer_name": "e2e", "status": "CREATED", "type": "NORMAL",
			"ordered_at": time.Now().UTC().Format(time.RFC3339), "total_amount": 1, "item_count": 1,
		})
		st, raw, ar, err := doRaw("POST", base+"/api/v1/admin/biz/cg-test-order/create", adminTok, string(body))
		if err != nil || assertBizOK(st, ar.Code) != nil {
			return fmt.Errorf("seed order for children: %v %s", err, truncate(string(raw), 120))
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/biz/cg-test-order/page?current=1&size=20&order_no="+url.QueryEscape(ono), adminTok, "")
		if err == nil && assertBizOK(st, ar.Code) == nil {
			_, data, _ := parseEnvelope(raw)
			recs, _ := assertPage(data)
			orderID = findIDByField(recs, "order_no", ono)
		}
		if orderID == "" {
			orderID = firstPageID(base, adminTok, "/api/v1/admin/biz/cg-test-order/page")
		}
	}
	specs["/api/v1/admin/biz/cg-test-order/children/create"] = crudSpec{
		findField: "sku_code", findValue: "E2E_SKU_" + suffix,
		create: map[string]any{
			"order_id": orderID, "sku_code": "E2E_SKU_" + suffix, "name": "E2E Item " + suffix,
			"status": "CREATED", "quantity": 1, "unit_price": 9.9, "is_gift": false,
		},
		update: func(id string, create map[string]any) map[string]any {
			m := map[string]any{"id": id}
			for k, v := range create {
				m[k] = v
			}
			m["name"] = "E2E Item " + suffix + " upd"
			m["quantity"] = 2
			return m
		},
		pageQuery: func(string) string {
			return "current=1&size=50&order_id=" + url.QueryEscape(orderID)
		},
	}

	// knowledge children（文档）依赖分类
	catID := firstPageID(base, adminTok, "/api/v1/admin/biz/cg-test-knowledge-category/page")
	specs["/api/v1/admin/biz/cg-test-knowledge-category/children/create"] = crudSpec{
		findField: "code", findValue: "E2E_DOC_" + suffix,
		create: map[string]any{
			"category_id": catID, "code": "E2E_DOC_" + suffix, "title": "E2E Doc " + suffix,
			"type": "ARTICLE", "status": "DRAFT", "content": "e2e", "sort": 99,
		},
		update: func(id string, create map[string]any) map[string]any {
			m := map[string]any{"id": id}
			for k, v := range create {
				m[k] = v
			}
			m["title"] = "E2E Doc " + suffix + " upd"
			return m
		},
		pageQuery: func(v string) string {
			return "current=1&size=50&code=" + url.QueryEscape(v)
		},
	}

	// accounts：RSA 加密默认密码
	enc, keyID, err := encryptPassword(base, "/api/v1/admin", "123456")
	if err != nil {
		return fmt.Errorf("account password encrypt: %w", err)
	}
	acc := "e2e_u_" + suffix
	specs["/api/v1/admin/sys/accounts/create"] = crudSpec{
		findField: "account", findValue: acc,
		create: map[string]any{
			"account": acc, "password": enc, "password_key_id": keyID,
			"account_type": "ADMIN", "account_status": "ENABLED", "name": "E2E User " + suffix,
		},
		update: func(id string, create map[string]any) map[string]any {
			return map[string]any{
				"id": id, "account": acc, "account_type": "ADMIN", "account_status": "ENABLED",
				"name": "E2E User " + suffix + " upd",
			}
		},
		pageQuery: func(v string) string {
			return "current=1&size=50&account=" + url.QueryEscape(v)
		},
	}

	// codegen：用真实表名
	table := firstTableName(base, adminTok)
	if table == "" {
		table = "sys_dict"
	}
	specs["/api/v1/admin/sys/codegen/create"] = crudSpec{
		findField: "name", findValue: "E2E Gen " + suffix,
		create: map[string]any{
			"name": "E2E Gen " + suffix, "gen_type": "SINGLE", "author": "e2e",
			"table_name": table, "pk_column": "id", "entity_name": "E2eEntity" + suffix,
			"module_path": "biz/e2e_" + suffix, "business_name": "e2e" + suffix,
			"api_prefix": "/v1/admin/biz/e2e-" + suffix, "permission_prefix": "biz:e2e" + suffix,
			"menu_name": "E2E Gen", "menu_path": "/e2e-gen/" + suffix, "component_path": "e2e/Gen",
			"sort": 99,
		},
		update: func(id string, create map[string]any) map[string]any {
			m := map[string]any{"id": id}
			for k, v := range create {
				m[k] = v
			}
			m["name"] = "E2E Gen " + suffix + " upd"
			return m
		},
		pageQuery: func(v string) string { return "current=1&size=50&name=" + url.QueryEscape(v) },
	}
	return nil
}

func encryptPassword(base, prefix, plain string) (cipherB64, keyID string, err error) {
	var pk struct {
		Data struct {
			KeyID     string `json:"key_id"`
			PublicKey string `json:"public_key"`
		} `json:"data"`
	}
	if err = getJSON(base+prefix+"/password-key", &pk); err != nil {
		return "", "", err
	}
	der, err := base64.StdEncoding.DecodeString(pk.Data.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubAny, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return "", "", err
	}
	enc, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubAny.(*rsa.PublicKey), []byte(plain), nil)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(enc), pk.Data.KeyID, nil
}

func firstTableName(base, token string) string {
	st, raw, ar, err := doRaw("GET", base+"/api/v1/admin/sys/codegen/tables", token, "")
	if err != nil || assertBizOK(st, ar.Code) != nil {
		return ""
	}
	arr, err := parseDataArrayMaps(ar)
	_ = raw
	if err != nil || len(arr) == 0 {
		return ""
	}
	return asString(arr[0]["table_name"])
}
