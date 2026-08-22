#!/usr/bin/env python3
"""Convert hei-boot FreeMarker frontend codegen templates to Go text/template."""

from __future__ import annotations

import pathlib
import re
import sys

BOOT = pathlib.Path(__file__).resolve().parents[2] / "hei-boot/module/sys/src/main/resources/codegen/templates"
OUT = pathlib.Path(__file__).resolve().parents[1] / "internal/modules/sys/codegen/templates/frontend"
FRONT = [
    "index.vue.ftl",
    "ModalForm.vue.ftl",
    "ModalDetail.vue.ftl",
    "ChildModalForm.vue.ftl",
    "ChildModalDetail.vue.ftl",
]

ROOT_KEYS = {
    "generated_at": ".GeneratedAt",
    "author": ".Author",
    "api_export_name": ".APIExportName",
    "api_export": ".APIExport",
    "has_tree": ".HasTree",
    "has_sub": ".HasSub",
    "has_tree_parent_form": ".HasTreeParentForm",
    "permissionPrefix": ".PermissionPrefix",
    "permission_prefix": ".PermissionPrefix",
    "gen_type": ".GenType",
}

PREFIX_KEYS = {
    "plan": ".Plan",
    "main": ".Main",
    "sub": ".Sub",
    "target": ".Target",
}


def dot(path: str, list_var: str | None = None) -> str:
    path = path.strip()
    if path.startswith("("):
        return path
    if path in ROOT_KEYS:
        return ROOT_KEYS[path]
    parts = path.split(".")
    head = parts[0]
    if head == list_var and len(parts) > 1:
        return f'(index ${list_var} "{parts[1]}")'
    if head in PREFIX_KEYS:
        cur = PREFIX_KEYS[head]
        for p in parts[1:]:
            cur = f'(index {cur} "{p}")'
        return cur
    if head == list_var:
        return f"${list_var}"
    return f'.{head}' + "".join(f'["{p}"]' for p in parts[1:])


def expr(s: str, list_var: str | None = None) -> str:
    s = s.strip()
    s = re.sub(r'\$\{r"\$\{([^}]+)\}"\}', r"${\1}", s)
    if s in ("has_sub && (sub??)", "has_sub && sub"):
        return "and $.HasSub $.Sub"
    if s == "has_sub":
        return "$.HasSub"
    if s == "has_tree":
        return "$.HasTree"
    if s == "has_tree_parent_form":
        return "$.HasTreeParentForm"
    if s == "sub??":
        return "$.Sub"
    if s == "field?is_first":
        return "$isFirst"
    if s == "field?is_last":
        return "$isLast"
    if s in ("field?index gte 8", "(field?index gte 8)"):
        return "ge $idx 8"
    m = re.match(r"(\w+)\.name\s*==\s*plan\.(\w+)", s)
    if m and (list_var is None or m.group(1) == list_var):
        var = m.group(1)
        return f'eq (index ${var} "name") (index $.Plan "{m.group(2)}")'
    m = re.match(r"(.+?)\s*!=\s*\"([^\"]+)\"", s)
    if m:
        return f'ne (deref {dot(m.group(1), list_var)}) "{m.group(2)}"'
    m = re.match(r"(.+?)\s*==\s*\"([^\"]+)\"", s)
    if m:
        return f'eq (deref {dot(m.group(1), list_var)}) "{m.group(2)}"'
    m = re.match(r"(.+?)\s*==\s*(.+)", s)
    if m:
        return f"eq (deref {dot(m.group(1), list_var)}) (deref {dot(m.group(2), list_var)})"
    if re.fullmatch(r"\w+\.dict_code\?\? && \w+\.dict_code\?has_content", s):
        var = s.split(".", 1)[0]
        return f"hasDictCode ${var}"
    m = re.match(r"(.+?)\s*&&\s*(.+)", s)
    if m and "||" not in s.split("&&", 1)[0]:
        return f"and {expr(m.group(1), list_var)} {expr(m.group(2), list_var)}"
    if "||" in s:
        parts = [expr(p.strip(), list_var) for p in re.split(r"\s*\|\|\s*", s)]
        return "or " + " ".join(parts)
    if "&&" in s:
        parts = [expr(p.strip(), list_var) for p in re.split(r"\s*&&\s*", s)]
        return "and " + " ".join(parts)
    if s.endswith("??"):
        return dot(s[:-2], list_var)
    if s.endswith(".is_bool"):
        return f'(index ${s[:-8]} "is_bool")'
    if s.endswith(".is_datetime"):
        return f'(index ${s[:-12]} "is_datetime")'
    if s.endswith(".is_json"):
        return f'(index ${s[:-8]} "is_json")'
    return dot(s, list_var)


def vue_mustache(inner: str) -> str:
    return "{{" + inner.strip() + "}}"


def repl_var(match: re.Match[str], list_var: str | None) -> str:
    inner = match.group(1)
    inner = re.sub(r'\$\{r"\$\{([^}]+)\}"\}', r"${\1}", inner)
    if inner.startswith('"') or inner.startswith("'"):
        return "{{" + inner + "}}"
    if re.match(r'\([^)]+\)\?has_content\?then\(', inner):
        m = re.match(r'\(([^)]+)\)\?has_content\?then\(([^,]+),\s*"([^"]*)"\)', inner)
        if m:
            return (
                "{{if (deref " + dot(m.group(1), list_var) + ")}}"
                + "{{" + dot(m.group(1), list_var) + "}}"
                + "{{else}}" + m.group(3) + "{{end}}"
            )
    m = re.match(r'\(([^)]+)\)\?then\("([^"]*)",\s*"([^"]*)"\)', inner)
    if m:
        return (
            "{{if " + expr(m.group(1), list_var) + "}}"
            + m.group(2)
            + "{{else}}"
            + m.group(3)
            + "{{end}}"
        )
    if list_var and inner.startswith(list_var + "."):
        return "{{index $" + list_var + ' "' + inner.split(".", 1)[1] + '"}}'
    if inner.startswith("field."):
        return "{{index $field \"" + inner.split(".", 1)[1] + '"}}'
    if inner == "api_export_name":
        return "{{.APIExportName}}"
    if "{{" in inner:
        return inner
    return "{{" + dot(inner, list_var) + "}}"


def convert_inline(line: str, list_var: str | None) -> str:
    while True:
        m = re.search(r"<#if\s+([^>]+)>(.*?)<#else>(.*?)</#if>", line)
        if m:
            rep = "{{if " + expr(m.group(1), list_var) + "}}" + m.group(2) + "{{else}}" + m.group(3) + "{{end}}"
            line = line[: m.start()] + rep + line[m.end() :]
            continue
        m = re.search(r"<#if\s+([^>]+)>(.*?)</#if>", line)
        if m:
            rep = "{{if " + expr(m.group(1), list_var) + "}}" + m.group(2) + "{{end}}"
            line = line[: m.start()] + rep + line[m.end() :]
            continue
        break
    line = re.sub(r"<#if\s+([^>]+)>", lambda mo: "{{if " + expr(mo.group(1), list_var) + "}}", line)
    line = re.sub(r"<#elseif\s+([^>]+)>", lambda mo: "{{else if " + expr(mo.group(1), list_var) + "}}", line)
    line = line.replace("<#else>", "{{else}}")
    line = line.replace("</#if>", "{{end}}")
    line = re.sub(r'\$\{r"\$\{([^}]+)\}"\}', r"${\1}", line)
    line = re.sub(r'\$\{"\{\{"\}\s*([^$]+?)\s*\$\{"\}\}"\}', lambda mo: vue_mustache(mo.group(1)), line)
    line = re.sub(r"\$\{([^#][^}]*)\}", lambda mo: repl_var(mo, list_var), line)
    return line


def convert_assign_wire_block(lines: list[str], i: int) -> tuple[str | None, int]:
    if not re.match(r"\s*<#assign _wires = \[\]>", lines[i]):
        return None, i
    j = i + 1
    while j < len(lines) and re.match(r"\s*<#if target\.has_form_", lines[j]):
        j += 1
    if j >= len(lines) or not re.match(r"\s*<#if \(_wires\?size > 0\)>", lines[j]):
        return None, i
    j += 1
    if j >= len(lines):
        return None, i
    import_line = lines[j]
    j += 1
    if j < len(lines) and re.match(r"\s*</#if>", lines[j]):
        j += 1
    out = '{{if wireImports .Target}}import { {{wireImports .Target}} } from \'@/utils/wire\'{{end}}\n'
    return out, j


def preprocess(content: str) -> str:
    def list_break(match: re.Match[str]) -> str:
        coll = dot(match.group(1))
        var = match.group(2)
        limit = match.group(3)
        return f'{{{{range $idx, ${var} := {coll}}}}}\n{{{{if ge $idx {limit}}}}}{{{{break}}}}{{{{end}}}}'

    content = re.sub(
        r"<#list\s+(\S+)\s+as\s+(\w+)><#if\s+\(\w+\?index\s+gte\s+(\d+)\)><#break></#if>",
        list_break,
        content,
    )
    content = re.sub(
        r"<#list\s+(\S+)\s+as\s+(\w+)>\s*<#if\s+\(\w+\?index\s+gte\s+(\d+)\)>\s*<#break>\s*</#if>",
        list_break,
        content,
    )
    return content


def postprocess(text: str) -> str:
    text = re.sub(
        r'\{\{if and \$\.HasTree \(index \$field "name == plan"\)\}\}',
        r'{{if and $.HasTree (eq (index $field "name") (index $.Plan "tree_parent_field"))}}',
        text,
    )
    text = re.sub(
        r'\{\{if and \$\.HasTreeParentForm \(index \$field "name == plan"\)\}\}',
        r'{{if and $.HasTreeParentForm (eq (index $field "name") (index $.Plan "tree_parent_field"))}}',
        text,
    )
    text = re.sub(
        r'\{\{if and \(index \$field "dict_code"\) \(index \$field "dict_code\?has_content"\)\}\}',
        r"{{if hasDictCode $field}}",
        text,
    )
    text = re.sub(r'\{\{else if \.field\["dict_code"\] \.field\["dict_code\?has_content"\]\}\}', r"{{else if hasDictCode $field}}", text)
    text = re.sub(r'\{\{else if \.field\["is_bool"\]\}\}', r"{{else if (index $field \"is_bool\")}}", text)
    text = re.sub(r'\{\{else if \.field\["is_datetime"\]\}\}', r"{{else if (index $field \"is_datetime\")}}", text)
    text = re.sub(
        r'<#list\s+(\S+)\s+as\s+(\w+)>\{\{if \(field\?index gte (\d+)\)\}\}<#break>\{\{end\}\}',
        lambda m: f'{{{{range $idx, ${m.group(2)} := {dot(m.group(1))}}}}}\n{{{{if ge $idx {m.group(3)}}}}}{{{{break}}}}{{{{end}}}}',
        text,
    )
    return text


def convert(content: str) -> str:
    content = preprocess(content)
    out: list[str] = []
    list_var: str | None = None
    lines = content.splitlines()
    i = 0
    while i < len(lines):
        line = lines[i]
        if re.match(r"\s*<#--", line):
            out.append(line)
            i += 1
            continue
        wire_block, next_i = convert_assign_wire_block(lines, i)
        if wire_block is not None:
            out.append(wire_block.rstrip("\n"))
            i = next_i
            continue
        if re.match(r"\s*<#assign\b", line):
            i += 1
            continue
        if "<#if" in line or "</#if>" in line or "<#list" in line or "<#elseif" in line or "<#else>" in line:
            if not re.match(r"\s*<#(if|elseif|else|list|/if|/list)", line):
                line = convert_inline(line, list_var)
                out.append(line)
                i += 1
                continue
        m = re.match(r"(\s*)<#if\s+(.+?)>\s*$", line)
        if m:
            out.append(f"{m.group(1)}{{{{if {expr(m.group(2), list_var)}}}}}")
            i += 1
            continue
        m = re.match(r"(\s*)<#elseif\s+(.+?)>\s*$", line)
        if m:
            out.append(f"{m.group(1)}{{{{else if {expr(m.group(2), list_var)}}}}}")
            i += 1
            continue
        m = re.match(r"(\s*)<#else>\s*$", line)
        if m:
            out.append(f"{m.group(1)}{{{{else}}}}")
            i += 1
            continue
        m = re.match(r"(\s*)</#if>\s*$", line)
        if m:
            out.append(f"{m.group(1)}{{{{end}}}}")
            i += 1
            continue
        m = re.match(r"(\s*)<#list\s+(\S+)\s+as\s+(\w+)>\s*$", line)
        if m:
            coll = dot(m.group(2), list_var)
            list_var = m.group(3)
            out.append(f'{m.group(1)}{{{{range $idx, ${list_var} := {coll}}}}}')
            i += 1
            continue
        m = re.match(r"(\s*)</#list>\s*$", line)
        if m:
            out.append(f"{m.group(1)}{{{{end}}}}")
            list_var = None
            i += 1
            continue
        m = re.match(r"(\s*)<#break>\s*$", line)
        if m:
            out.append(f"{m.group(1)}{{{{break}}}}")
            i += 1
            continue
        line2 = convert_inline(line, list_var)
        out.append(line2)
        i += 1
    return postprocess("\n".join(out) + "\n")


def main() -> int:
    OUT.mkdir(parents=True, exist_ok=True)
    for name in FRONT:
        src = BOOT / name
        dst = OUT / name.replace(".ftl", ".tmpl")
        dst.write_text(convert(src.read_text(encoding="utf-8")), encoding="utf-8")
        print(f"wrote {dst.name} ({dst.stat().st_size} bytes)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
