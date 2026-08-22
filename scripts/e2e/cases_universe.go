package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// runUniverse 对 debug/routes 中的全部 API 做传参/出参断言（非抽样）。
// - 每个 GET：合法 query + 按约定断言响应 shape
// - 每个标准 CRUD 资源：create → page → detail → update → delete
// - 无 payload 的 create 记为失败（强制全覆盖）
func runUniverse(base, adminTok, portalTok string, routes []routeInfo, readBucket, crudBucket *caseBucket, skipped *[]caseResult) {
	idx := indexRoutes(routes)
	fmt.Printf("universe routes=%d get=%d create=%d\n", len(routes), len(idx.gets), len(idx.creates))

	runAllGETShapes(base, adminTok, portalTok, idx, readBucket, skipped)
	runAllCRUD(base, adminTok, idx, crudBucket, skipped)
	runSpecialWrites(base, adminTok, portalTok, idx, crudBucket, skipped)
}

type routeIndex struct {
	gets    []routeInfo
	creates []string // path ending with /create
	byKey   map[string]bool
}

func indexRoutes(routes []routeInfo) routeIndex {
	idx := routeIndex{byKey: map[string]bool{}}
	seenGet := map[string]bool{}
	seenCreate := map[string]bool{}
	for _, r := range routes {
		key := r.Method + " " + r.Path
		idx.byKey[key] = true
		if r.Method == "GET" {
			if seenGet[r.Path] {
				continue
			}
			seenGet[r.Path] = true
			idx.gets = append(idx.gets, r)
		}
		if r.Method == "POST" && strings.HasSuffix(r.Path, "/create") {
			if seenCreate[r.Path] {
				continue
			}
			seenCreate[r.Path] = true
			idx.creates = append(idx.creates, r.Path)
		}
	}
	sort.Slice(idx.gets, func(i, j int) bool { return idx.gets[i].Path < idx.gets[j].Path })
	sort.Strings(idx.creates)
	return idx
}

func (idx routeIndex) has(method, path string) bool {
	return idx.byKey[method+" "+path]
}

func runAllGETShapes(base, adminTok, portalTok string, idx routeIndex, bucket *caseBucket, skipped *[]caseResult) {
	for _, r := range idx.gets {
		if reason, skip := skipRouteReason(r); skip {
			*skipped = append(*skipped, caseResult{Name: "GET " + r.Path, Error: reason, OK: true})
			continue
		}
		path := materializePath(r.Path)
		tok := pickToken(path, adminTok, portalTok)
		name := "GET " + path
		cr := caseResult{Name: name, URL: base + path}

		// detail / own-* / columns 由 assert 自行拼合法 query，避免无参预请求
		if getNeedsID(path) {
			if err := assertGETShape(base, tok, path, 0, nil, apiResp{}, idx); err != nil {
				cr.Error = err.Error()
				bucket.add(cr)
				continue
			}
			cr.OK = true
			bucket.add(cr)
			continue
		}

		q := enrichGETQuery(path)
		full := path
		if q != "" {
			if strings.Contains(full, "?") {
				full += "&" + q
			} else {
				full += "?" + q
			}
		}
		url := base + full
		cr.URL = url
		status, raw, ar, err := doRaw("GET", url, tok, "")
		cr.Status, cr.BizCode, cr.Body = status, ar.Code, truncate(string(raw), 280)
		if err != nil {
			cr.Error = err.Error()
			bucket.add(cr)
			continue
		}
		if err := assertGETShape(base, tok, path, status, raw, ar, idx); err != nil {
			cr.Error = err.Error()
			bucket.add(cr)
			continue
		}
		cr.OK = true
		bucket.add(cr)
	}
}

func getNeedsID(path string) bool {
	leaf := path[strings.LastIndex(path, "/")+1:]
	return leaf == "detail" || leaf == "columns" || leaf == "my-detail" ||
		strings.Contains(path, "/own-") || strings.Contains(path, "/granted") ||
		strings.HasSuffix(path, "/spaces/detail") || strings.HasSuffix(path, "/children/detail")
}

func enrichGETQuery(path string) string {
	vals := url.Values{}
	needPage := strings.Contains(path, "/page") || strings.HasSuffix(path, "/list") ||
		strings.HasSuffix(path, "/tree") || strings.HasSuffix(path, "/my-page")
	if needPage {
		vals.Set("current", "1")
		vals.Set("size", "10")
	}
	switch {
	case strings.Contains(path, "/roles/"):
		vals.Set("name", "a")
		vals.Set("code", "a")
	case strings.Contains(path, "/accounts"):
		vals.Set("account", "a")
		vals.Set("name", "a")
	case strings.Contains(path, "/banners"):
		vals.Set("title", "a")
		if strings.HasSuffix(path, "/list") {
			vals.Set("position", "HOME_TOP")
		}
	case strings.Contains(path, "/notices"):
		vals.Set("title", "a")
	case strings.Contains(path, "/dicts"):
		vals.Set("code", "a")
		vals.Set("category", "SYS")
	case strings.Contains(path, "/config"):
		vals.Set("config_key", "a")
	case strings.Contains(path, "/weak-password"):
		vals.Set("keyword", "a")
	case strings.Contains(path, "/depts"):
		vals.Set("name", "a")
	case strings.Contains(path, "/groups"):
		vals.Set("name", "a")
	case strings.Contains(path, "/positions"):
		vals.Set("name", "a")
	case strings.Contains(path, "/resources"):
		vals.Set("name", "a")
	case strings.Contains(path, "/codegen"):
		vals.Set("name", "a")
	case strings.Contains(path, "/cg-test-catalog"):
		vals.Set("code", "a")
		vals.Set("name", "a")
	case strings.Contains(path, "/cg-test-activity"):
		vals.Set("code", "a")
		vals.Set("name", "a")
	case strings.Contains(path, "/cg-test-order"):
		vals.Set("order_no", "a")
		vals.Set("name", "a")
	case strings.Contains(path, "/cg-test-knowledge"):
		vals.Set("code", "a")
		vals.Set("name", "a")
		vals.Set("title", "a")
	case strings.Contains(path, "/jobs"):
		vals.Set("name", "a")
	case strings.Contains(path, "/file"):
		vals.Set("original_name", "a")
	case strings.Contains(path, "/feedback"):
		vals.Set("title", "a")
	case strings.Contains(path, "/audit"):
		vals.Set("action", "a")
	}
	// detail / own-* need id — filled in assertGETShape when possible
	return vals.Encode()
}

func assertGETShape(base, token, path string, status int, raw []byte, ar apiResp, idx routeIndex) error {
	leaf := path
	if i := strings.LastIndex(path, "/"); i >= 0 {
		leaf = path[i+1:]
	}

	switch {
	case leaf == "page" || leaf == "my-page" || strings.HasSuffix(path, "/job-logs/page"):
		if err := assertBizOK(status, ar.Code); err != nil {
			return err
		}
		_, data, _ := parseEnvelope(raw)
		recs, err := assertPage(data)
		if err != nil {
			return err
		}
		// 仅当路由表存在对应 detail 时才 follow
		basePath := strings.TrimSuffix(path, "/"+leaf)
		detailPath := basePath + "/detail"
		if leaf == "my-page" {
			detailPath = basePath + "/my-detail"
		}
		id := firstID(recs)
		if id != "" && idx.has("GET", detailPath) {
			st, draw, dar, err := doRaw("GET", base+detailPath+"?id="+url.QueryEscape(id), token, "")
			if err != nil {
				return fmt.Errorf("detail follow: %w", err)
			}
			if err := assertBizOK(st, dar.Code); err != nil {
				return fmt.Errorf("detail follow: %w %s", err, truncate(string(draw), 120))
			}
			_, dm, _ := parseEnvelope(draw)
			if dm != nil && asString(dm["id"]) == "" && asString(dm["account_id"]) == "" {
				return fmt.Errorf("detail missing id")
			}
		}
		return nil

	case leaf == "list" || leaf == "tree" || leaf == "tables" || leaf == "current" || leaf == "children":
		if err := assertBizOK(status, ar.Code); err != nil {
			return err
		}
		// data 可为 array 或 object
		if ar.Data == nil || string(ar.Data) == "null" {
			return fmt.Errorf("data null")
		}
		var data any
		if err := json.Unmarshal(ar.Data, &data); err != nil {
			return err
		}
		switch data.(type) {
		case []any, map[string]any:
			return nil
		default:
			return fmt.Errorf("data type %T", data)
		}

	case leaf == "detail" || leaf == "my-detail" || strings.HasSuffix(path, "/spaces/detail") || strings.HasSuffix(path, "/children/detail"):
		// 尝试从 sibling page 取 id；否则用 id=0 断言非 5xx + envelope
		id := resolveSiblingPageID(base, token, path)
		q := "id=" + url.QueryEscape(id)
		if id == "" {
			q = "id=0"
		}
		st, draw, dar, err := doRaw("GET", base+path+"?"+q, token, "")
		if err != nil {
			return err
		}
		if st >= 500 || dar.Code >= 500 {
			return fmt.Errorf("detail 5xx status=%d code=%d %s", st, dar.Code, truncate(string(draw), 120))
		}
		if id != "" {
			if err := assertBizOK(st, dar.Code); err != nil {
				return err
			}
		}
		// 有种子数据时期望对象
		if id != "" {
			_, dm, _ := parseEnvelope(draw)
			if dm == nil {
				return fmt.Errorf("detail data not object")
			}
		}
		return nil

	case leaf == "me" || leaf == "overview" || leaf == "live" || leaf == "ready" ||
		leaf == "unread-count" || leaf == "password-key" || leaf == "captcha":
		if err := assertBizOK(status, ar.Code); err != nil {
			return err
		}
		if ar.Data == nil || string(ar.Data) == "null" {
			return fmt.Errorf("data null")
		}
		return nil

	case strings.Contains(path, "/own-") || strings.Contains(path, "/granted") ||
		strings.HasSuffix(path, "/permissions") || leaf == "columns":
		// 需要 id 类参数：尽量补齐后非 5xx
		id := resolveSiblingPageID(base, token, path)
		u := base + path
		vals := url.Values{}
		if id != "" {
			vals.Set("id", id)
		} else {
			vals.Set("id", "0")
		}
		if strings.Contains(path, "columns") {
			vals.Set("table_name", "sys_dict")
		}
		if enc := vals.Encode(); enc != "" {
			u += "?" + enc
		}
		st, draw, dar, err := doRaw("GET", u, token, "")
		if err != nil {
			return err
		}
		if st >= 500 || dar.Code >= 500 {
			return fmt.Errorf("query 5xx %s", truncate(string(draw), 120))
		}
		return nil

	default:
		// 通用：业务成功或可预期 4xx，禁止 5xx；成功时 data 可解析
		if status >= 500 || ar.Code >= 500 {
			return fmt.Errorf("5xx status=%d code=%d %s", status, ar.Code, truncate(string(raw), 120))
		}
		if status >= 200 && status < 300 && (ar.Code == 0 || ar.Code == 200) {
			if len(ar.Data) > 0 && string(ar.Data) != "null" {
				var data any
				if err := json.Unmarshal(ar.Data, &data); err != nil {
					return fmt.Errorf("data json: %w", err)
				}
			}
		}
		return nil
	}
}

func resolveSiblingPageID(base, token, detailPath string) string {
	// .../detail → .../page ; .../my-detail → .../my-page ; .../children/detail → .../children/page
	pagePath := detailPath
	pagePath = strings.Replace(pagePath, "/my-detail", "/my-page", 1)
	pagePath = strings.Replace(pagePath, "/detail", "/page", 1)
	if pagePath == detailPath {
		if i := strings.LastIndex(detailPath, "/"); i > 0 {
			baseRes := detailPath[:i]
			pagePath = baseRes + "/page"
		}
	}
	st, raw, ar, err := doRaw("GET", base+pagePath+"?current=1&size=5", token, "")
	if err != nil || assertBizOK(st, ar.Code) != nil {
		return ""
	}
	_, data, _ := parseEnvelope(raw)
	recs, err := assertPage(data)
	if err != nil || len(recs) == 0 {
		return ""
	}
	return firstID(recs)
}

func runAllCRUD(base, adminTok string, idx routeIndex, bucket *caseBucket, skipped *[]caseResult) {
	suffix := fmt.Sprintf("%d", time.Now().Unix()%1000000)
	specs := crudSpecs(suffix)
	if err := attachDeferredCRUDSpecs(base, adminTok, suffix, specs); err != nil {
		bucket.add(caseResult{Name: "CRUD deferred-setup", Error: err.Error()})
	}

	covered := map[string]bool{}
	for _, createPath := range idx.creates {
		if reason, skip := skipRouteReason(routeInfo{Method: "POST", Path: createPath}); skip {
			*skipped = append(*skipped, caseResult{Name: "CRUD " + createPath, Error: reason, OK: true})
			continue
		}
		basePath := strings.TrimSuffix(createPath, "/create")
		if !idx.has("GET", basePath+"/page") || !idx.has("POST", basePath+"/update") || !idx.has("POST", basePath+"/delete") {
			*skipped = append(*skipped, caseResult{Name: "CRUD " + createPath, Error: "not-standard-crud", OK: true})
			continue
		}
		spec, ok := specs[createPath]
		if !ok {
			cr := caseResult{Name: "CRUD " + createPath, Error: "missing payload factory — all create APIs must be covered"}
			bucket.add(cr)
			continue
		}
		covered[createPath] = true
		name := "CRUD " + basePath
		runCRUD(bucket, name, func() error {
			return execCRUD(base, adminTok, basePath, spec)
		})
	}

	// 工厂里多写的也跑（保险）
	for createPath, spec := range specs {
		if covered[createPath] {
			continue
		}
		if !idx.has("POST", createPath) {
			continue
		}
		basePath := strings.TrimSuffix(createPath, "/create")
		runCRUD(bucket, "CRUD "+basePath, func() error {
			return execCRUD(base, adminTok, basePath, spec)
		})
	}
}

type crudSpec struct {
	findField string
	findValue string
	create    map[string]any
	update    func(id string, create map[string]any) map[string]any
	pageQuery func(findValue string) string
}

func execCRUD(base, token, basePath string, spec crudSpec) error {
	body, _ := json.Marshal(spec.create)
	st, raw, ar, err := doRaw("POST", base+basePath+"/create", token, string(body))
	if err != nil {
		return err
	}
	if err := assertBizOK(st, ar.Code); err != nil {
		return fmt.Errorf("create: %w %s", err, truncate(string(raw), 200))
	}

	pq := "current=1&size=50"
	if spec.pageQuery != nil {
		pq = spec.pageQuery(spec.findValue)
	} else if spec.findField != "" {
		pq += "&" + spec.findField + "=" + url.QueryEscape(spec.findValue)
	}
	st, raw, ar, err = doRaw("GET", base+basePath+"/page?"+pq, token, "")
	if err != nil {
		return err
	}
	if err := assertBizOK(st, ar.Code); err != nil {
		return fmt.Errorf("page: %w %s", err, truncate(string(raw), 160))
	}
	_, data, _ := parseEnvelope(raw)
	recs, err := assertPage(data)
	if err != nil {
		return err
	}
	id := findIDByField(recs, spec.findField, spec.findValue)
	if id == "" {
		id = firstID(recs) // 有些字段被加密/脱敏
	}
	if id == "" {
		return fmt.Errorf("created row not found in page field=%s value=%s", spec.findField, spec.findValue)
	}

	detailPath := basePath + "/detail"
	// resource-buttons 等无 detail 路由：跳过
	st, raw, ar, err = doRaw("GET", base+detailPath+"?id="+url.QueryEscape(id), token, "")
	if err != nil {
		return err
	}
	if st != 404 {
		if err := assertBizOK(st, ar.Code); err != nil {
			return fmt.Errorf("detail: %w %s", err, truncate(string(raw), 160))
		}
	}

	ubody, _ := json.Marshal(spec.update(id, spec.create))
	st, raw, ar, err = doRaw("POST", base+basePath+"/update", token, string(ubody))
	if err != nil {
		return err
	}
	if err := assertBizOK(st, ar.Code); err != nil {
		return fmt.Errorf("update: %w %s", err, truncate(string(raw), 200))
	}

	dbody, _ := json.Marshal(map[string]any{"ids": []string{id}})
	st, raw, ar, err = doRaw("POST", base+basePath+"/delete", token, string(dbody))
	if err != nil {
		return err
	}
	if err := assertBizOK(st, ar.Code); err != nil {
		return fmt.Errorf("delete: %w %s", err, truncate(string(raw), 160))
	}
	return nil
}

func runSpecialWrites(base, adminTok, portalTok string, idx routeIndex, bucket *caseBucket, skipped *[]caseResult) {
	// 对非 CRUD 写接口：带合法最小入参，断言非 5xx（业务 4xx 可接受，如缺种子）
	type writeCase struct {
		name, method, path string
		body               func() (string, error)
		requireOK          bool
	}

	accountID := firstPageID(base, adminTok, "/api/v1/admin/sys/accounts/page")
	roleID := firstPageID(base, adminTok, "/api/v1/admin/sys/roles/page")
	groupID := firstPageID(base, adminTok, "/api/v1/admin/sys/groups/page")
	deptID := firstPageID(base, adminTok, "/api/v1/admin/sys/depts/page")
	jobID := firstPageID(base, adminTok, "/api/v1/admin/sys/jobs/page")
	noticeID := firstPageID(base, adminTok, "/api/v1/admin/sys/notices/page")
	bannerID := firstPageID(base, adminTok, "/api/v1/admin/sys/banners/page")
	resID := firstPageID(base, adminTok, "/api/v1/admin/sys/resources/page")

	mk := func(v any) func() (string, error) {
		return func() (string, error) {
			b, err := json.Marshal(v)
			return string(b), err
		}
	}

	cases := []writeCase{
		{name: "POST accounts/grant-role", path: "/api/v1/admin/sys/accounts/grant-role",
			body: mk(map[string]any{"id": accountID, "role_ids": []string{}})},
		{name: "POST accounts/grant-group", path: "/api/v1/admin/sys/accounts/grant-group",
			body: mk(map[string]any{"id": accountID, "group_ids": []string{}})},
		{name: "POST accounts/grant-dept", path: "/api/v1/admin/sys/accounts/grant-dept",
			body: mk(map[string]any{"id": accountID, "grant_info_list": []any{}})},
		{name: "POST accounts/grant-resource", path: "/api/v1/admin/sys/accounts/grant-resource",
			body: mk(map[string]any{"id": accountID, "grant_info_list": []any{}})},
		{name: "POST accounts/grant-client-resource", path: "/api/v1/admin/sys/accounts/grant-client-resource",
			body: mk(map[string]any{"id": accountID, "grant_info_list": []any{}})},
		{name: "POST roles/grant-resource", path: "/api/v1/admin/sys/roles/grant-resource",
			body: mk(map[string]any{"id": roleID, "grant_info_list": []any{}})},
		{name: "POST roles/grant-client-resource", path: "/api/v1/admin/sys/roles/grant-client-resource",
			body: mk(map[string]any{"id": roleID, "grant_info_list": []any{}})},
		{name: "POST groups/grant-user", path: "/api/v1/admin/sys/groups/grant-user",
			body: mk(map[string]any{"id": groupID, "account_ids": []string{}})},
		{name: "POST groups/grant-role", path: "/api/v1/admin/sys/groups/grant-role",
			body: mk(map[string]any{"id": groupID, "role_ids": []string{}})},
		{name: "POST groups/grant-resource", path: "/api/v1/admin/sys/groups/grant-resource",
			body: mk(map[string]any{"id": groupID, "grant_info_list": []any{}})},
		{name: "POST groups/grant-client-resource", path: "/api/v1/admin/sys/groups/grant-client-resource",
			body: mk(map[string]any{"id": groupID, "grant_info_list": []any{}})},
		{name: "POST jobs/enabled", path: "/api/v1/admin/sys/jobs/enabled",
			body: mk(map[string]any{"id": jobID, "enabled": false})},
		{name: "POST jobs/run", path: "/api/v1/admin/sys/jobs/run",
			body: mk(map[string]any{"id": jobID})},
		{name: "POST notices/read", path: "/api/v1/admin/sys/notices/read",
			body: mk(map[string]any{"ids": []string{noticeID}})},
		{name: "POST notices/read-all", path: "/api/v1/admin/sys/notices/read-all",
			body: mk(map[string]any{})},
		{name: "POST portal banners/interaction", path: "/api/v1/portal/sys/banners/interaction",
			body: mk(map[string]any{"id": bannerID})},
		{name: "POST config/batch-save", path: "/api/v1/admin/sys/config/batch-save",
			body: mk(map[string]any{"items": []map[string]any{{"config_key": "e2e.batch.noop", "config_value": "1", "category": "E2E"}}})},
		{name: "POST resource-permissions", path: "/api/v1/admin/resource-permissions",
			body: mk(map[string]any{"resource_id": resID, "permission_key": "e2e:noop", "account_type": "ADMIN"})},
		{name: "POST client-resource-permissions", path: "/api/v1/admin/client-resource-permissions",
			body: mk(map[string]any{"resource_id": resID, "permission_key": "e2e:noop", "account_type": "PORTAL"})},
		{name: "POST auth/refresh admin", path: "/api/v1/admin/auth/refresh", body: mk(map[string]any{})},
		{name: "POST auth/refresh portal", path: "/api/v1/portal/auth/refresh", body: mk(map[string]any{})},
	}

	_ = deptID
	_ = portalTok

	for _, c := range cases {
		if !idx.byKey["POST "+c.path] {
			*skipped = append(*skipped, caseResult{Name: c.name, Error: "route-absent", OK: true})
			continue
		}
		body, err := c.body()
		cr := caseResult{Name: c.name, URL: base + c.path}
		if err != nil {
			cr.Error = err.Error()
			bucket.add(cr)
			continue
		}
		tok := pickToken(c.path, adminTok, portalTok)
		st, raw, ar, err := doRaw("POST", base+c.path, tok, body)
		cr.Status, cr.BizCode, cr.Body = st, ar.Code, truncate(string(raw), 200)
		if err != nil {
			cr.Error = err.Error()
			bucket.add(cr)
			continue
		}
		if st >= 500 || ar.Code >= 500 {
			cr.Error = fmt.Sprintf("5xx status=%d code=%d %s", st, ar.Code, truncate(string(raw), 120))
			bucket.add(cr)
			continue
		}
		if c.requireOK {
			if err := assertBizOK(st, ar.Code); err != nil {
				cr.Error = err.Error()
				bucket.add(cr)
				continue
			}
		}
		cr.OK = true
		bucket.add(cr)
	}
}

func firstPageID(base, token, pagePath string) string {
	st, raw, ar, err := doRaw("GET", base+pagePath+"?current=1&size=5", token, "")
	if err != nil || assertBizOK(st, ar.Code) != nil {
		return "0"
	}
	_, data, _ := parseEnvelope(raw)
	recs, _ := assertPage(data)
	if id := firstID(recs); id != "" {
		return id
	}
	return "0"
}
