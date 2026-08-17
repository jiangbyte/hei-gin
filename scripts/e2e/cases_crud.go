package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func runCRUDCases(base, adminTok string, bucket *caseBucket) {
	suffix := fmt.Sprintf("%d", time.Now().Unix()%1000000)

	// ---- dict ----
	runCRUD(bucket, "crud_dict", func() error {
		code := "E2E_DICT_" + suffix
		cat := "BIZ"
		label := "e2e-dict"
		body, _ := json.Marshal(map[string]any{
			"code": code, "label": label, "value": code, "category": cat, "status": "ENABLED", "sort": 99,
		})
		st, raw, ar, err := doRaw("POST", base+"/api/v1/admin/sys/dicts/create", adminTok, string(body))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("create: %w %s", err, truncate(string(raw), 160))
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/sys/dicts/page?current=1&size=20&code="+code, adminTok, "")
		if err != nil {
			return err
		}
		_, data, _ := parseEnvelope(raw)
		recs, err := assertPage(data)
		if err != nil {
			return err
		}
		id := findIDByField(recs, "code", code)
		if id == "" {
			return fmt.Errorf("created dict not found in page")
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/sys/dicts/detail?id="+id, adminTok, "")
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("detail: %w", err)
		}
		_, dm, _ := parseEnvelope(raw)
		if asString(dm["code"]) != code {
			return fmt.Errorf("detail code want %s got %s", code, asString(dm["code"]))
		}
		newLabel := "e2e-dict-upd"
		ubody, _ := json.Marshal(map[string]any{
			"id": id, "code": code, "label": newLabel, "value": code, "category": cat, "status": "ENABLED", "sort": 98,
		})
		st, raw, ar, err = doRaw("POST", base+"/api/v1/admin/sys/dicts/update", adminTok, string(ubody))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("update: %w %s", err, truncate(string(raw), 160))
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/sys/dicts/detail?id="+id, adminTok, "")
		if err != nil {
			return err
		}
		_, dm, _ = parseEnvelope(raw)
		if asString(dm["label"]) != newLabel {
			return fmt.Errorf("update not reflected label=%v", dm["label"])
		}
		dbody, _ := json.Marshal(map[string]any{"ids": []string{id}})
		st, raw, ar, err = doRaw("POST", base+"/api/v1/admin/sys/dicts/delete", adminTok, string(dbody))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("delete: %w %s", err, truncate(string(raw), 160))
		}
		return nil
	})

	// ---- weak password ----
	runCRUD(bucket, "crud_weak_password", func() error {
		pwd := "e2e_weak_" + suffix
		body, _ := json.Marshal(map[string]any{"password": pwd})
		st, raw, ar, err := doRaw("POST", base+"/api/v1/admin/sys/weak-password/create", adminTok, string(body))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("create: %w %s", err, truncate(string(raw), 160))
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/sys/weak-password/page?current=1&size=50&keyword="+pwd, adminTok, "")
		if err != nil {
			return err
		}
		_, data, _ := parseEnvelope(raw)
		recs, err := assertPage(data)
		if err != nil {
			return err
		}
		id := findIDByField(recs, "password", pwd)
		if id == "" {
			return fmt.Errorf("weak password not found")
		}
		newPwd := pwd + "_u"
		ubody, _ := json.Marshal(map[string]any{"id": id, "password": newPwd})
		st, raw, ar, err = doRaw("POST", base+"/api/v1/admin/sys/weak-password/update", adminTok, string(ubody))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("update: %w %s", err, truncate(string(raw), 160))
		}
		dbody, _ := json.Marshal(map[string]any{"ids": []string{id}})
		st, raw, ar, err = doRaw("POST", base+"/api/v1/admin/sys/weak-password/delete", adminTok, string(dbody))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("delete: %w %s", err, truncate(string(raw), 160))
		}
		return nil
	})

	// ---- banner ----
	runCRUD(bucket, "crud_banner", func() error {
		title := "E2E Banner " + suffix
		body, _ := json.Marshal(map[string]any{
			"title": title, "image": "https://example.com/e2e.png", "link_type": "NONE",
			"category": "HOME", "type": "CAROUSEL", "position": "HOME_TOP",
			"sort": 99, "status": "DISABLED",
		})
		st, raw, ar, err := doRaw("POST", base+"/api/v1/admin/sys/banners/create", adminTok, string(body))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("create: %w %s", err, truncate(string(raw), 200))
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/sys/banners/page?current=1&size=50&title=E2E", adminTok, "")
		if err != nil {
			return err
		}
		_, data, _ := parseEnvelope(raw)
		recs, err := assertPage(data)
		if err != nil {
			return err
		}
		id := findIDByField(recs, "title", title)
		if id == "" {
			return fmt.Errorf("banner not found")
		}
		if err := assertKeys(recs[0], "target_account_types", "position"); err != nil {
			return err
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/sys/banners/detail?id="+id, adminTok, "")
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("detail: %w", err)
		}
		_, dm, _ := parseEnvelope(raw)
		if asString(dm["title"]) != title {
			return fmt.Errorf("detail title mismatch")
		}
		newTitle := title + " upd"
		ubody, _ := json.Marshal(map[string]any{
			"id": id, "title": newTitle, "image": "https://example.com/e2e2.png", "link_type": "NONE",
			"category": "HOME", "type": "CAROUSEL", "position": "HOME_TOP",
			"sort": 98, "status": "DISABLED",
		})
		st, raw, ar, err = doRaw("POST", base+"/api/v1/admin/sys/banners/update", adminTok, string(ubody))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("update: %w %s", err, truncate(string(raw), 200))
		}
		dbody, _ := json.Marshal(map[string]any{"ids": []string{id}})
		st, raw, ar, err = doRaw("POST", base+"/api/v1/admin/sys/banners/delete", adminTok, string(dbody))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("delete: %w %s", err, truncate(string(raw), 160))
		}
		return nil
	})

	// ---- notice ----
	runCRUD(bucket, "crud_notice", func() error {
		title := "E2E Notice " + suffix
		cat := "SYSTEM"
		body, _ := json.Marshal(map[string]any{
			"kind": "NOTIFICATION", "title": title, "content": "e2e content", "content_type": "text",
			"category": cat, "severity": "INFO", "target_scope": "ALL",
			"target_account_types": []string{"ADMIN"}, "status": "DRAFT",
			"publish_locations": map[string]any{"center": true},
		})
		st, raw, ar, err := doRaw("POST", base+"/api/v1/admin/sys/notices/create", adminTok, string(body))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("create: %w %s", err, truncate(string(raw), 200))
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/sys/notices/page?current=1&size=50&title=E2E", adminTok, "")
		if err != nil {
			return err
		}
		_, data, _ := parseEnvelope(raw)
		recs, err := assertPage(data)
		if err != nil {
			return err
		}
		id := findIDByField(recs, "title", title)
		if id == "" {
			return fmt.Errorf("notice not found")
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/sys/notices/detail?id="+id, adminTok, "")
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("detail: %w %s", err, truncate(string(raw), 160))
		}
		_, dm, _ := parseEnvelope(raw)
		if err := assertKeys(dm, "id", "title", "kind", "target_scope", "target_account_types"); err != nil {
			return err
		}
		newTitle := title + " upd"
		ubody, _ := json.Marshal(map[string]any{
			"id": id, "kind": "NOTIFICATION", "title": newTitle, "content": "e2e content2", "content_type": "text",
			"category": cat, "severity": "INFO", "target_scope": "ALL",
			"target_account_types": []string{"ADMIN"}, "status": "DRAFT",
			"publish_locations": map[string]any{"center": true},
		})
		st, raw, ar, err = doRaw("POST", base+"/api/v1/admin/sys/notices/update", adminTok, string(ubody))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("update: %w %s", err, truncate(string(raw), 200))
		}
		dbody, _ := json.Marshal(map[string]any{"ids": []string{id}})
		st, raw, ar, err = doRaw("POST", base+"/api/v1/admin/sys/notices/delete", adminTok, string(dbody))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("delete: %w %s", err, truncate(string(raw), 160))
		}
		return nil
	})

	// ---- cg_test_catalog ----
	runCRUD(bucket, "crud_cg_catalog", func() error {
		code := "E2E_CAT_" + suffix
		name := "E2E Catalog " + suffix
		body, _ := json.Marshal(map[string]any{
			"code": code, "name": name, "category": "BUSINESS", "status": "ENABLED",
			"sort": 99, "is_visible": true, "extra": map[string]any{"e2e": true},
		})
		st, raw, ar, err := doRaw("POST", base+"/api/v1/admin/biz/cg-test-catalog/create", adminTok, string(body))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("create: %w %s", err, truncate(string(raw), 200))
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/biz/cg-test-catalog/page?current=1&size=50&code="+code, adminTok, "")
		if err != nil {
			return err
		}
		_, data, _ := parseEnvelope(raw)
		recs, err := assertPage(data)
		if err != nil {
			return err
		}
		id := findIDByField(recs, "code", code)
		if id == "" {
			return fmt.Errorf("catalog not found")
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/biz/cg-test-catalog/detail?id="+id, adminTok, "")
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("detail: %w", err)
		}
		newName := name + " upd"
		ubody, _ := json.Marshal(map[string]any{
			"id": id, "code": code, "name": newName, "category": "BUSINESS", "status": "ENABLED",
			"sort": 98, "is_visible": false, "extra": map[string]any{"e2e": true},
		})
		st, raw, ar, err = doRaw("POST", base+"/api/v1/admin/biz/cg-test-catalog/update", adminTok, string(ubody))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("update: %w %s", err, truncate(string(raw), 200))
		}
		st, raw, ar, err = doRaw("GET", base+"/api/v1/admin/biz/cg-test-catalog/tree", adminTok, "")
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("tree: %w", err)
		}
		dbody, _ := json.Marshal(map[string]any{"ids": []string{id}})
		st, raw, ar, err = doRaw("POST", base+"/api/v1/admin/biz/cg-test-catalog/delete", adminTok, string(dbody))
		if err != nil {
			return err
		}
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("delete: %w %s", err, truncate(string(raw), 160))
		}
		return nil
	})
}

func runCRUD(bucket *caseBucket, name string, fn func() error) {
	cr := caseResult{Name: name}
	if err := fn(); err != nil {
		cr.Error = err.Error()
		bucket.add(cr)
		return
	}
	cr.OK = true
	bucket.add(cr)
}
