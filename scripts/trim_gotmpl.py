#!/usr/bin/env python3
"""Add Go template whitespace trim markers to frontend templates."""

from __future__ import annotations

import pathlib
import re

ROOT = pathlib.Path(__file__).resolve().parents[1] / "internal/modules/sys/codegen/templates/frontend"


def trim_actions(text: str) -> str:
    def repl(match: re.Match[str]) -> str:
        tag = match.group(1)
        if tag.startswith('"'):
            return match.group(0)
        if tag.startswith("-"):
            return match.group(0)
        return "{{-" + tag

    return re.sub(r"\{\{((?:\"|\.)?[^}]+)\}\}", repl, text)


def main() -> int:
    for path in ROOT.glob("*.tmpl"):
        if path.name.startswith("api"):
            continue
        text = path.read_text(encoding="utf-8")
        text = trim_actions(text)
        text = text.replace("{{- end }}", "{{- end -}}")
        path.write_text(text, encoding="utf-8")
        print(path.name)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
