#!/usr/bin/env python3
"""Full API contract diff: hei-boot OpenAPI vs hei-gin debug routes + sampled schemas."""

from __future__ import annotations

import json
import re
import sys
import urllib.error
import urllib.request
from collections import defaultdict
from pathlib import Path
from typing import Any

BOOT = "http://127.0.0.1:8000"
GIN = "http://127.0.0.1:8001"
OPENAPI_FILE = Path(__file__).with_name("boot-openapi.json")

SKIP_BOOT_PREFIXES = (
    "/actuator/",
    "/doc.html",
    "/swagger-ui",
    "/v3/api-docs",
    "/webjars/",
)
SKIP_BOOT_PATHS = {
    "/api/v1/admin/dashboard/overview",  # removed, replaced by workspace
}
GIN_ONLY_OK = {
    "/api/v1/internal/debug/routes",
    "/metrics",
}


def http_json(method: str, url: str, token: str = "", body: dict | None = None) -> tuple[int, Any]:
    data = None
    headers = {"Accept": "application/json"}
    if body is not None:
        data = json.dumps(body).encode()
        headers["Content-Type"] = "application/json"
    if token:
        headers["Authorization"] = token
    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=60) as resp:
            raw = resp.read().decode()
            return resp.status, json.loads(raw) if raw else None
    except urllib.error.HTTPError as e:
        raw = e.read().decode()
        try:
            return e.code, json.loads(raw)
        except Exception:
            return e.code, raw


def login(base: str) -> str:
    _, resp = http_json("POST", f"{base}/api/v1/admin/login", body={
        "account": "superadmin",
        "password": "123456",
        "login_mode": "PASSWORD",
        "captcha_id": "",
        "captcha_code": "",
        "password_encrypted": False,
    })
    tok = (resp or {}).get("data", {}).get("token", "")
    if not tok:
        raise RuntimeError(f"login failed: {resp}")
    return tok


def norm_path(p: str) -> str:
    p = re.sub(r"\{[^}]+\}", ":id", p)
    if not p.startswith("/"):
        p = "/" + p
    return p


def boot_routes(oa: dict) -> set[tuple[str, str]]:
    out: set[tuple[str, str]] = set()
    for path, methods in (oa.get("paths") or {}).items():
        if any(path.startswith(p) for p in SKIP_BOOT_PREFIXES):
            continue
        if path in SKIP_BOOT_PATHS:
            continue
        for m in methods:
            if m.upper() in {"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}:
                out.add((m.upper(), norm_path(path)))
    return out


def gin_routes(token: str) -> set[tuple[str, str]]:
    _, resp = http_json("GET", f"{GIN}/api/v1/internal/debug/routes", token)
    out: set[tuple[str, str]] = set()
    for r in (resp or {}).get("data") or []:
        m = (r.get("method") or "").upper()
        p = norm_path(r.get("path") or "")
        if m and p:
            out.add((m, p))
    return out


def ref_resolve(oa: dict, node: Any) -> Any:
    if not isinstance(node, dict):
        return node
    if "$ref" in node:
        ref = node["$ref"]
        if ref.startswith("#/components/schemas/"):
            name = ref.split("/")[-1]
            return ref_resolve(oa, (oa.get("components", {}).get("schemas") or {}).get(name, {}))
    out = {}
    for k, v in node.items():
        out[k] = ref_resolve(oa, v)
    return out


def schema_props(oa: dict, schema: Any) -> set[str]:
    schema = ref_resolve(oa, schema or {})
    if not isinstance(schema, dict):
        return set()
    t = schema.get("type")
    if t == "array":
        return schema_props(oa, schema.get("items"))
    props = set(schema.get("properties") or {})
    for key in ("allOf", "oneOf", "anyOf"):
        for sub in schema.get(key) or []:
            props |= schema_props(oa, sub)
    return props


def endpoint_io(oa: dict, method: str, path: str) -> dict[str, set[str]]:
    raw_path = path
    # reverse norm for openapi lookup: try exact and pattern match
    paths = oa.get("paths") or {}
    op = None
    for p, ms in paths.items():
        if norm_path(p) == raw_path and method.lower() in ms:
            op = ms[method.lower()]
            break
    if not op:
        return {"req": set(), "res": set()}
    req_props: set[str] = set()
    for content in (op.get("requestBody") or {}).get("content", {}).values():
        req_props |= schema_props(oa, content.get("schema"))
    res_props: set[str] = set()
    for content in (op.get("responses") or {}).get("200", {}).get("content", {}).values():
        res_props |= schema_props(oa, content.get("schema"))
    return {"req": req_props, "res": res_props}


def infer_keys_from_sample(data: Any) -> set[str]:
    if isinstance(data, dict):
        if "records" in data and isinstance(data["records"], list) and data["records"]:
            return set(data["records"][0].keys()) if isinstance(data["records"][0], dict) else set(data.keys())
        return set(data.keys())
    return set()


def sample_call(base: str, token: str, method: str, path: str) -> tuple[int, Any]:
    # materialize :id placeholders with known seed ids
    p = path.replace(":id", "1")
    if "jobs/detail" in p:
        p = p.replace("id=1", "id=7541000000000000001")
    if "detail" in p and "id=" not in p:
        sep = "&" if "?" in p else "?"
        if "job" in p:
            p += f"{sep}id=7541000000000000001"
        else:
            p += f"{sep}id=1"
    if "page" in p and "current=" not in p:
        sep = "&" if "?" in p else "?"
        p += f"{sep}current=1&size=5"
    return http_json(method, base + p, token)


def main() -> int:
    if not OPENAPI_FILE.exists():
        print("missing", OPENAPI_FILE, "- run curl boot /v3/api-docs first")
        return 2
    oa = json.loads(OPENAPI_FILE.read_text(encoding="utf-8"))
    boot_tok = login(BOOT)
    gin_tok = login(GIN)

    br = boot_routes(oa)
    gr = gin_routes(gin_tok)
    print(f"boot routes: {len(br)}")
    print(f"gin routes:  {len(gr)}")

    only_boot = sorted(br - gr)
    only_gin = sorted(r for r in (gr - br) if r[1] not in GIN_ONLY_OK)

    print(f"\n=== PATH COVERAGE ===")
    print(f"only in boot ({len(only_boot)}):")
    for m, p in only_boot[:80]:
        print(f"  {m} {p}")
    if len(only_boot) > 80:
        print(f"  ... and {len(only_boot)-80} more")

    print(f"\nonly in gin ({len(only_gin)}):")
    for m, p in only_gin[:40]:
        print(f"  {m} {p}")

    common = sorted(br & gr)
    print(f"\ncommon routes: {len(common)}")

    shape_diffs: list[str] = []
    checked = 0
    for method, path in common:
        if method not in {"GET"}:
            continue
        if any(x in path for x in ("/download", "/upload", "/captcha", "/oauth/", "/preview", "/generate")):
            continue
        bs, brsp = sample_call(BOOT, boot_tok, method, path)
        gs, grsp = sample_call(GIN, gin_tok, method, path)
        if not isinstance(brsp, dict) or not isinstance(grsp, dict):
            continue
        bd = brsp.get("data")
        gd = grsp.get("data")
        if bd is None and gd is None:
            continue
        bk = infer_keys_from_sample(bd)
        gk = infer_keys_from_sample(gd)
        if bk != gk and (bk or gk):
            only_b = sorted(bk - gk)
            only_g = sorted(gk - bk)
            if only_b or only_g:
                shape_diffs.append(f"{method} {path}: boot_only={only_b} gin_only={only_g} status boot={bs} gin={gs}")
        checked += 1

    print(f"\n=== RESPONSE SHAPE DIFFS (GET sampled {checked}) ===")
    for d in shape_diffs[:100]:
        print(" -", d)
    if len(shape_diffs) > 100:
        print(f" ... and {len(shape_diffs)-100} more")

    out = {
        "boot_count": len(br),
        "gin_count": len(gr),
        "only_boot": [f"{m} {p}" for m, p in only_boot],
        "only_gin": [f"{m} {p}" for m, p in only_gin],
        "shape_diffs": shape_diffs,
    }
    out_file = Path(__file__).with_name("full-api-diff.json")
    out_file.write_text(json.dumps(out, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"\nwritten {out_file}")

    if only_boot or only_gin or shape_diffs:
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
