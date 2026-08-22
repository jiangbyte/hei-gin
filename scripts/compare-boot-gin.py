#!/usr/bin/env python3
"""Compare hei-boot (8000) vs hei-gin (8001) API request/response shapes."""

from __future__ import annotations

import json
import sys
import urllib.error
import urllib.parse
import urllib.request
from typing import Any

BOOT = "http://127.0.0.1:8000"
GIN = "http://127.0.0.1:8001"


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
        with urllib.request.urlopen(req, timeout=30) as resp:
            raw = resp.read().decode()
            return resp.status, json.loads(raw) if raw else None
    except urllib.error.HTTPError as e:
        raw = e.read().decode()
        try:
            return e.code, json.loads(raw)
        except Exception:
            return e.code, raw


def login(base: str, account: str = "superadmin", password: str = "123456") -> str:
    _, pk = http_json("GET", f"{base}/api/v1/admin/password-key")
    pub = pk.get("data", {}).get("public_key", "")
    # boot/gin both accept plaintext password when crypto not enforced in dev
    _, resp = http_json("POST", f"{base}/api/v1/admin/login", body={
        "account": account,
        "password": password,
        "login_mode": "PASSWORD",
        "captcha_id": "",
        "captcha_code": "",
        "password_encrypted": False,
        "public_key": pub,
    })
    tok = resp.get("data", {}).get("token", "")
    if not tok:
        raise RuntimeError(f"login failed on {base}: {resp}")
    return tok


def keys_of(obj: Any) -> set[str]:
    if isinstance(obj, dict):
        return set(obj.keys())
    return set()


def compare_record(boot_rec: dict, gin_rec: dict, label: str) -> list[str]:
    diffs: list[str] = []
    bk, gk = keys_of(boot_rec), keys_of(gin_rec)
    only_boot = sorted(bk - gk)
    only_gin = sorted(gk - bk)
    if only_boot:
        diffs.append(f"{label}: only in boot keys={only_boot}")
    if only_gin:
        diffs.append(f"{label}: only in gin keys={only_gin}")
    return diffs


def main() -> int:
    print("=== Login ===")
    boot_tok = login(BOOT)
    gin_tok = login(GIN)
    print("boot login OK, gin login OK")

    endpoints = [
        ("GET", "/api/v1/admin/me", None, "me"),
        ("GET", "/api/v1/admin/workspace/overview", None, "workspace"),
        ("GET", "/api/v1/admin/sys/jobs/page?current=1&size=5", None, "jobs_page"),
        ("GET", "/api/v1/admin/sys/jobs/detail?id=7541000000000000001", None, "job_detail"),
        ("GET", "/api/v1/admin/sys/job-logs/page?current=1&size=3&job_id=7541000000000000001", None, "job_logs"),
        ("GET", "/api/v1/admin/sys/roles/page?current=1&size=3", None, "roles_page"),
        ("GET", "/api/v1/public/site-footer", None, "site_footer"),
    ]

    all_diffs: list[str] = []
    print("\n=== Endpoint shape compare ===")
    for method, path, body, name in endpoints:
        bs, br = http_json(method, BOOT + path, boot_tok, body)
        gs, gr = http_json(method, GIN + path, gin_tok, body)
        print(f"\n[{name}] boot={bs} gin={gs}")
        if bs != gs:
            all_diffs.append(f"{name}: status boot={bs} gin={gs}")
            continue
        if not isinstance(br, dict) or not isinstance(gr, dict):
            continue
        if br.get("code") != gr.get("code"):
            all_diffs.append(f"{name}: code boot={br.get('code')} gin={gr.get('code')}")
        bd, gd = br.get("data"), gr.get("data")
        if isinstance(bd, dict) and isinstance(gd, dict) and "records" in bd and "records" in gd:
            if bd.get("records") and gd.get("records"):
                all_diffs.extend(compare_record(bd["records"][0], gd["records"][0], f"{name}.record[0]"))
        elif isinstance(bd, dict) and isinstance(gd, dict):
            all_diffs.extend(compare_record(bd, gd, name))

    print("\n=== Job detail field values (boot vs gin) ===")
    _, bj = http_json("GET", f"{BOOT}/api/v1/admin/sys/jobs/detail?id=7541000000000000001", boot_tok)
    _, gj = http_json("GET", f"{GIN}/api/v1/admin/sys/jobs/detail?id=7541000000000000001", gin_tok)
    for key in ["id", "name", "handler", "trigger_type", "trigger_config", "enabled", "sort"]:
        bv = (bj.get("data") or {}).get(key)
        gv = (gj.get("data") or {}).get(key)
        match = "OK" if bv == gv else "DIFF"
        print(f"  {key}: boot={bv!r} gin={gv!r} [{match}]")

    print("\n=== OpenAPI paths (boot) vs debug routes (gin) ===")
    try:
        _, oa = http_json("GET", f"{BOOT}/v3/api-docs")
        boot_paths = set((oa.get("paths") or {}).keys())
    except Exception as e:
        boot_paths = set()
        print("boot openapi fetch failed:", e)
    try:
        _, dr = http_json("GET", f"{GIN}/api/v1/internal/debug/routes", gin_tok)
        gin_paths = {r.get("path") for r in (dr.get("data") or []) if r.get("path")}
    except Exception as e:
        gin_paths = set()
        print("gin routes fetch failed:", e)

  # normalize parameterized paths
    def norm(p: str) -> str:
        return p.replace("{", ":").replace("}", "")

    boot_norm = {norm(p) for p in boot_paths}
    only_boot = sorted(boot_norm - gin_paths)[:30]
    only_gin = sorted(gin_paths - boot_norm)[:30]
    print(f"boot paths={len(boot_norm)} gin paths={len(gin_paths)}")
    if only_boot:
        print("sample only boot:", only_boot[:15])
    if only_gin:
        print("sample only gin:", only_gin[:15])

    print("\n=== SUMMARY ===")
    if all_diffs:
        print(f"DIFFERENCES ({len(all_diffs)}):")
        for d in all_diffs:
            print(" -", d)
        return 1
    print("No field-key differences on sampled endpoints.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
