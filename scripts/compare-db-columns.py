#!/usr/bin/env python3
"""Compare MySQL column sets: hei_boot (source) vs target copy."""

import json
import subprocess

TABLES = [
    "sys_job", "sys_job_log", "sys_codegen_plan", "sys_codegen_field",
    "profile_identity", "real_name_case", "sys_workspace_shortcut",
    "iam_account", "sys_audit_log", "profile_user_admin",
]

TARGET_DB = "hei_gin"
SOURCE_DB = "hei_boot"


def cols(db: str, table: str) -> set[str]:
    sql = (
        "SELECT column_name FROM information_schema.columns "
        f"WHERE table_schema='{db}' AND table_name='{table}' ORDER BY 1"
    )
    out = subprocess.check_output(
        [
            "wsl",
            "docker",
            "exec",
            "dev-mysql",
            "mysql",
            "-uroot",
            "-p123456",
            "-N",
            "-e",
            sql,
        ],
        text=True,
    )
    return {line.strip() for line in out.splitlines() if line.strip()}


def main():
    report = {}
    for t in TABLES:
        boot = cols(SOURCE_DB, t)
        gin = cols(TARGET_DB, t)
        report[t] = {
            "boot_only": sorted(boot - gin),
            "gin_only": sorted(gin - boot),
            "same": boot == gin,
        }
    print(json.dumps(report, ensure_ascii=False, indent=2))
    diff_tables = [t for t, v in report.items() if not v["same"]]
    print(f"\nDIFF tables: {diff_tables}")


if __name__ == "__main__":
    main()
