#!/usr/bin/env python3
"""Module-by-module boot vs gin response shape comparison."""

from __future__ import annotations

import base64
import json
import sys
import urllib.error
import urllib.request
from typing import Any

import bcrypt
import redis
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import padding

BOOT = "http://127.0.0.1:8000"
GIN = "http://127.0.0.1:8001"

MODULE_ENDPOINTS: dict[str, list[tuple[str, str]]] = {
    "auth": [
        ("GET", "/api/v1/admin/public/auth-options"),
        ("GET", "/api/v1/portal/public/auth-options"),
        ("GET", "/api/v1/public/site-footer"),
    ],
    "workspace": [
        ("GET", "/api/v1/admin/workspace/overview"),
        ("GET", "/api/v1/admin/workspace/shortcuts"),
    ],
    "profile": [
        ("GET", "/api/v1/admin/me"),
        ("GET", "/api/v1/portal/me"),
        ("GET", "/api/v1/admin/profile/org-info"),
    ],
    "identity": [
        ("GET", "/api/v1/admin/profile/identity/status"),
        ("GET", "/api/v1/portal/profile/identity/status"),
    ],
    "iam/account": [
        ("GET", "/api/v1/admin/sys/accounts/page?current=1&size=2"),
        ("GET", "/api/v1/admin/sys/accounts/detail?id=1"),
    ],
    "iam/role": [("GET", "/api/v1/admin/sys/roles/page?current=1&size=2")],
    "iam/dept": [
        ("GET", "/api/v1/admin/sys/depts/page?current=1&size=2"),
        ("GET", "/api/v1/admin/sys/depts/tree"),
    ],
    "iam/group": [("GET", "/api/v1/admin/sys/groups/page?current=1&size=2")],
    "iam/position": [("GET", "/api/v1/admin/sys/positions/page?current=1&size=2")],
    "iam/resource": [
        ("GET", "/api/v1/admin/sys/resources/tree"),
        ("GET", "/api/v1/admin/sys/resources/current"),
        ("GET", "/api/v1/portal/sys/resources/current"),
    ],
    "iam/client": [
        ("GET", "/api/v1/admin/sys/client-modules/page?current=1&size=2"),
        ("GET", "/api/v1/admin/sys/client-resources/tree"),
    ],
    "iam/relation": [("GET", "/api/v1/admin/permission-registry")],
    "sys/audit": [
        ("GET", "/api/v1/admin/sys/audit/page?current=1&size=2"),
        ("GET", "/api/v1/admin/sys/audit/my-page?current=1&size=2"),
    ],
    "sys/banner": [
        ("GET", "/api/v1/admin/sys/banners/page?current=1&size=2"),
        ("GET", "/api/v1/portal/sys/banners/list"),
    ],
    "sys/codegen": [
        ("GET", "/api/v1/admin/sys/codegen/page?current=1&size=2"),
        ("GET", "/api/v1/admin/sys/codegen/tables"),
    ],
    "sys/config": [
        ("GET", "/api/v1/admin/sys/config/page?current=1&size=2"),
        ("GET", "/api/v1/admin/sys/config/list"),
    ],
    "sys/dict": [
        ("GET", "/api/v1/admin/sys/dicts/page?current=1&size=2"),
        ("GET", "/api/v1/portal/sys/dicts/tree"),
    ],
    "sys/feedback": [("GET", "/api/v1/admin/sys/feedbacks/page?current=1&size=2")],
    "sys/file": [("GET", "/api/v1/admin/sys/file/page?current=1&size=2")],
    "sys/job": [
        ("GET", "/api/v1/admin/sys/jobs/page?current=1&size=2"),
        ("GET", "/api/v1/admin/sys/job-logs/page?current=1&size=2"),
    ],
    "sys/notice": [
        ("GET", "/api/v1/admin/sys/notices/page?current=1&size=2"),
        ("GET", "/api/v1/admin/sys/notices/my-page?current=1&size=2"),
    ],
    "sys/weak_password": [("GET", "/api/v1/admin/sys/weak-password/page?current=1&size=2")],
    "biz/activity": [("GET", "/api/v1/admin/biz/cg-test-activity/page?current=1&size=2")],
    "biz/catalog": [("GET", "/api/v1/admin/biz/cg-test-catalog/page?current=1&size=2")],
    "biz/order": [("GET", "/api/v1/admin/biz/cg-test-order/page?current=1&size=2")],
    "biz/knowledge": [("GET", "/api/v1/admin/biz/cg-test-knowledge-category/page?current=1&size=2")],
}

SKIP_MODULES = {"sys/job", "sys/codegen"}

SKIP_KEYS = {
    "created_at",
    "updated_at",
    "bound_at",
    "verified_at",
    "last_login_at",
    "token",
    "password_hash",
    "trans_map",
    "children",
}

PORTAL_ONLY_PREFIXES = (
    "/api/v1/portal/",
)


def should_skip_key(key: str) -> bool:
    if key in SKIP_KEYS:
        return True
    return key.endswith("_name") or key.endswith("_names")


def http_json(method: str, url: str, token: str = "", body: dict | None = None) -> tuple[int, Any]:
    headers = {"Accept": "application/json"}
    data = None
    if body is not None:
        data = json.dumps(body).encode()
        headers["Content-Type"] = "application/json"
    if token:
        headers["Authorization"] = token
    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            raw = resp.read().decode()
            return resp.status, json.loads(raw) if raw else None
    except urllib.error.HTTPError as e:
        raw = e.read().decode()
        try:
            return e.code, json.loads(raw)
        except Exception:
            return e.code, raw


def plant_captcha(base: str, redis_db: int) -> str:
    _, cap = http_json("GET", f"{base}/api/v1/admin/captcha")
    cap_id = ((cap or {}).get("data") or {}).get("captcha_id", "")
    if not cap_id:
        raise RuntimeError(f"captcha id missing from {base}")
    rdb = redis.Redis(host="127.0.0.1", port=6379, password="123456", db=redis_db, decode_responses=True)
    hashed = bcrypt.hashpw(b"test", bcrypt.gensalt()).decode()
    rdb.set(f"captcha:{cap_id}", hashed, ex=300)
    return cap_id


def encrypt_password(pub_b64: str, password: str) -> str:
    der = base64.b64decode(pub_b64)
    pub = serialization.load_der_public_key(der)
    enc = pub.encrypt(
        password.encode(),
        padding.OAEP(mgf=padding.MGF1(algorithm=hashes.SHA256()), algorithm=hashes.SHA256(), label=None),
    )
    return base64.b64encode(enc).decode()


def login(base: str, redis_db: int) -> str:
    cap_id = plant_captcha(base, redis_db)
    _, pk = http_json("GET", f"{base}/api/v1/admin/password-key")
    key_id = (pk or {}).get("data", {}).get("key_id", "")
    pub = (pk or {}).get("data", {}).get("public_key", "")
    enc_pwd = encrypt_password(pub, "123456")
    _, resp = http_json(
        "POST",
        f"{base}/api/v1/admin/login",
        body={
            "account": "superadmin",
            "password": enc_pwd,
            "password_key_id": key_id,
            "captcha_id": cap_id,
            "captcha_value": "test",
            "remember_me": True,
            "login_mode": "PASSWORD",
        },
    )
    tok = ((resp or {}).get("data") or {}).get("token", "")
    if not tok:
        raise RuntimeError(f"login failed on {base}: {resp}")
    return tok


def keys_of(obj: Any) -> set[str]:
    return set(obj.keys()) if isinstance(obj, dict) else set()


def shape_diff(a: Any, b: Any, prefix: str = "") -> list[str]:
    diffs: list[str] = []
    if type(a) is not type(b):
        return [f"{prefix}: type {type(a).__name__} vs {type(b).__name__}"]
    if isinstance(a, dict):
        only_a = sorted(k for k in keys_of(a) - keys_of(b) if not should_skip_key(k))
        only_b = sorted(k for k in keys_of(b) - keys_of(a) if not should_skip_key(k))
        if only_a:
            diffs.append(f"{prefix} only_boot={only_a[:10]}")
        if only_b:
            diffs.append(f"{prefix} only_gin={only_b[:10]}")
        for k in sorted(keys_of(a) & keys_of(b)):
            if should_skip_key(k):
                continue
            diffs.extend(shape_diff(a[k], b[k], f"{prefix}.{k}" if prefix else k)[:2])
    elif isinstance(a, list) and a and b and isinstance(a[0], dict) and isinstance(b[0], dict):
        diffs.extend(shape_diff(a[0], b[0], f"{prefix}[0]")[:3])
    return diffs[:6]


def compare_endpoint(method: str, path: str, boot_tok: str, gin_tok: str) -> list[str]:
    if any(path.startswith(p) for p in PORTAL_ONLY_PREFIXES):
        return []
    bs, br = http_json(method, BOOT + path, boot_tok)
    gs, gr = http_json(method, GIN + path, gin_tok)
    if bs != gs:
        return [f"{path}: HTTP {bs} vs {gs}"]
    if not isinstance(br, dict) or not isinstance(gr, dict):
        return []
    if str(br.get("code")) != str(gr.get("code")):
        return [f"{path}: code {br.get('code')} vs {gr.get('code')}"]
    bd, gd = br.get("data"), gr.get("data")
    if isinstance(bd, dict) and isinstance(gd, dict) and "records" in bd and "records" in gd:
        if bd.get("records") and gd.get("records"):
            return shape_diff(bd["records"][0], gd["records"][0], f"{path}.record")
        if bool(bd.get("records")) != bool(gd.get("records")):
            return [f"{path}: records empty mismatch"]
        return []
    if isinstance(bd, dict) and isinstance(gd, dict):
        return shape_diff(bd, gd, path)
    if type(bd) is not type(gd):
        return [f"{path}: data type {type(bd).__name__} vs {type(gd).__name__}"]
    return []


def main() -> int:
    boot_tok = login(BOOT, 0)
    gin_tok = login(GIN, 4)
    print("login OK (boot + gin)\n")

    report: dict[str, Any] = {"modules": {}, "summary": {}}
    ok_count = 0
    for mod, endpoints in MODULE_ENDPOINTS.items():
        if mod in SKIP_MODULES:
            report["modules"][mod] = {"status": "skipped", "diffs": []}
            print(f"[SKIP] {mod}")
            continue
        diffs: list[str] = []
        for method, path in endpoints:
            diffs.extend(compare_endpoint(method, path, boot_tok, gin_tok))
        status = "aligned" if not diffs else "diff"
        if status == "aligned":
            ok_count += 1
        report["modules"][mod] = {"status": status, "diffs": diffs}
        mark = "OK" if not diffs else "DIFF"
        print(f"[{mark}] {mod}")
        for d in diffs[:5]:
            print(f"  - {d}")

    total = len(MODULE_ENDPOINTS) - len(SKIP_MODULES)
    report["summary"] = {
        "aligned": ok_count,
        "total": total,
        "diff": total - ok_count,
        "skipped": sorted(SKIP_MODULES),
    }
    out = "scripts/module-compare.json"
    with open(out, "w", encoding="utf-8") as f:
        json.dump(report, f, ensure_ascii=False, indent=2)
    print(f"\n=== SUMMARY aligned={ok_count}/{total} ===")
    print(f"report: {out}")
    return 0 if ok_count == total else 1


if __name__ == "__main__":
    sys.exit(main())
