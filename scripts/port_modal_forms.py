#!/usr/bin/env python3
"""Port ModalForm/ChildModalForm FTL to Go templates."""

from __future__ import annotations

import pathlib
import re

BOOT = pathlib.Path(__file__).resolve().parents[2] / "hei-boot/module/sys/src/main/resources/codegen/templates"
OUT = pathlib.Path(__file__).resolve().parents[1] / "internal/modules/sys/codegen/templates/frontend"


def port(content: str, child: bool = False) -> str:
    content = content.replace("<#-- Author: Charlie -->", "")
    content = re.sub(r"\$\{plan\.author\}", '{{index .Plan "author"}}', content)
    content = re.sub(r"\$\{generated_at\}", "{{.GeneratedAt}}", content)
    content = re.sub(r"\$\{api_export_name\}", "{{.APIExportName}}", content)
    content = re.sub(r"\$\{plan\.main_business_name\}", '{{index .Plan "main_business_name"}}', content)
    content = re.sub(r"\$\{plan\.tree_parent_field\}", '{{index .Plan "tree_parent_field"}}', content)
    content = re.sub(r"\$\{plan\.tree_label_field\}", '{{index .Plan "tree_label_field"}}', content)
    content = re.sub(r"\$\{main\.pk_name\}", '{{index .Main "pk_name"}}', content)
    content = re.sub(r"\$\{target\.pk_name\}", '{{index .Target "pk_name"}}', content)
    content = re.sub(
        r"\$\{\(plan\.sub_business_name\)\?has_content\?then\(plan\.sub_business_name, \"明细\"\)\}",
        '{{bootThen (index .Plan "sub_business_name") "明细"}}',
        content,
    )
    content = re.sub(
        r"\$\{\(target\.needs_submit_normalize \|\| has_tree_parent_form\)\?then\(\"normalizeSubmitData\(state\.formModel\)\", \"state\.formModel\"\)\}",
        '{{if or (index .Target "needs_submit_normalize") $.HasTreeParentForm}}normalizeSubmitData(state.formModel){{else}}state.formModel{{end}}',
        content,
    )
    content = re.sub(
        r"\$\{\(target\.needs_submit_normalize\)\?then\(\"normalizeSubmitData\(state\.formModel\)\", \"state\.formModel\"\)\}",
        '{{if (index .Target "needs_submit_normalize")}}normalizeSubmitData(state.formModel){{else}}state.formModel{{end}}',
        content,
    )
    content = re.sub(r"\$\{\"\{\{\"\}\s*([^$]+?)\s*\$\{\"\}\}\"\}", r"{{ \1 }}", content)

    # wire assign block
    content = re.sub(
        r"<#assign _wires = \[\]>\s*"
        r"<#if target\.has_form_bool><#assign _wires = _wires \+ \[\"wireBool\"\]></#if>\s*"
        r"<#if target\.has_form_int><#assign _wires = _wires \+ \[\"wireInt\"\]></#if>\s*"
        r"<#if target\.has_form_float><#assign _wires = _wires \+ \[\"wireFloat\"\]></#if>\s*"
        r"<#if \(_wires\?size > 0\)>\s*"
        r"import \{ \$\{_wires\?join\(\", \"\)\} \} from '@/utils/wire'\s*"
        r"</#if>",
        "{{if wireImports .Target}}import { {{wireImports .Target}} } from '@/utils/wire'{{end}}",
        content,
        flags=re.S,
    )

    content = re.sub(
        r"import \{ createRequiredRule<#if target\.has_form_datetime>, toApiDateTime, toFormDateTime</#if> \}",
        'import { createRequiredRule{{if (index .Target "has_form_datetime")}}, toApiDateTime, toFormDateTime{{end}} }',
        content,
    )
    content = re.sub(r"<#if target\.has_form_icon>", '{{if (index .Target "has_form_icon")}}', content)
    content = re.sub(r"<#if target\.has_form_editor>", '{{if (index .Target "has_form_editor")}}', content)
    content = re.sub(r"<#if target\.has_form_richtext>", '{{if (index .Target "has_form_richtext")}}', content)
    content = re.sub(r"<#if target\.has_form_markdown>", '{{if (index .Target "has_form_markdown")}}', content)
    content = re.sub(r"<#if target\.has_form_code>", '{{if (index .Target "has_form_code")}}', content)
    content = re.sub(r"<#if has_tree_parent_form>", "{{if $.HasTreeParentForm}}", content)
    content = re.sub(r"<#if target\.needs_submit_normalize \|\| has_tree_parent_form>", '{{if or (index .Target "needs_submit_normalize") $.HasTreeParentForm}}', content)
    content = re.sub(r"<#if target\.needs_submit_normalize>", '{{if (index .Target "needs_submit_normalize")}}', content)
    content = re.sub(r"<#if target\.has_form_json>", '{{if (index .Target "has_form_json")}}', content)

    # inline list+if on one line
    content = re.sub(
        r"<#list target\.form_fields as field><#if field\.required \|\| field\.is_json>",
        '{{range $idx, $field := (index .Target "form_fields")}}\n{{if or (index $field "required") (index $field "is_json")}}',
        content,
    )
    content = re.sub(
        r"</#if></#list>",
        "{{end}}\n{{end}}",
        content,
    )
    content = re.sub(
        r"<#list target\.form_fields as field><#if field\.is_bool>",
        '{{range $idx, $field := (index .Target "form_fields")}}\n{{if (index $field "is_bool")}}',
        content,
    )
    content = re.sub(
        r"<#list target\.form_fields as field><#if field\.value_type == \"int\">",
        '{{range $idx, $field := (index .Target "form_fields")}}\n{{if eq (index $field "value_type") "int"}}',
        content,
    )
    content = re.sub(
        r"<#list target\.form_fields as field><#if field\.value_type == \"float\">",
        '{{range $idx, $field := (index .Target "form_fields")}}\n{{if eq (index $field "value_type") "float"}}',
        content,
    )
    content = re.sub(
        r"<#list target\.form_fields as field><#if field\.is_datetime>",
        '{{range $idx, $field := (index .Target "form_fields")}}\n{{if (index $field "is_datetime")}}',
        content,
    )
    content = re.sub(
        r"<#list target\.form_fields as field><#if field\.is_json>",
        '{{range $idx, $field := (index .Target "form_fields")}}\n{{if (index $field "is_json")}}',
        content,
    )

    content = re.sub(r"<#list target\.form_fields as field>", '{{range $idx, $field := (index .Target "form_fields")}}', content)
    content = re.sub(r"</#list>", "{{end}}", content)
    content = re.sub(r"<#if has_tree_parent_form && field\.name == plan\.tree_parent_field>", "{{if and $.HasTreeParentForm (eq (index $field \"name\") (index $.Plan \"tree_parent_field\"))}}", content)
    content = re.sub(r"<#elseif \(field\.value_type == \"int\" \|\| field\.value_type == \"float\"\)>", '{{else if or (eq (index $field "value_type") "int") (eq (index $field "value_type") "float")}}', content)
    content = re.sub(r"<#elseif field\.widget == \"number\" \|\| \(field\.value_type == \"int\" \|\| field\.value_type == \"float\"\)>", '{{else if or (eq (index $field "widget") "number") (eq (index $field "value_type") "int") (eq (index $field "value_type") "float")}}', content)
    content = re.sub(r"<#elseif field\.widget == \"richtext\">", '{{else if eq (index $field "widget") "richtext"}}', content)
    content = re.sub(r"<#elseif field\.widget == \"markdown\">", '{{else if eq (index $field "widget") "markdown"}}', content)
    content = re.sub(r"<#elseif field\.widget == \"code\">", '{{else if eq (index $field "widget") "code"}}', content)
    content = re.sub(r"<#elseif field\.widget == \"icon\">", '{{else if eq (index $field "widget") "icon"}}', content)
    content = re.sub(r"<#elseif field\.widget == \"textarea\">", '{{else if eq (index $field "widget") "textarea"}}', content)
    content = re.sub(r"<#elseif field\.widget == \"dict\" && field\.dict_code>", '{{else if and (eq (index $field "widget") "dict") (index $field "dict_code")}}', content)
    content = re.sub(r"<#elseif field\.is_json>", '{{else if (index $field "is_json")}}', content)
    content = re.sub(r"<#elseif field\.value_type == \"bool\">", '{{else if eq (index $field "value_type") "bool"}}', content)
    content = re.sub(r"<#elseif field\.widget == \"datetime\" \|\| field\.value_type == \"datetime\">", '{{else if or (eq (index $field "widget") "datetime") (eq (index $field "value_type") "datetime")}}', content)
    content = re.sub(r"<#if field\.required>", "{{if (index $field \"required\")}}", content)
    content = re.sub(r"<#if field\.value_type == \"bool\">", '{{if eq (index $field "value_type") "bool"}}', content)
    content = re.sub(r"<#if field\.is_json>", "{{if (index $field \"is_json\")}}", content)
    content = re.sub(r"<#if field\.widget == \"dict\" \|\| field\.value_type == \"bool\" \|\| field\.is_datetime \|\| field\.widget == \"icon\">'change'<#else>'input'</#if>", "{{if or (eq (index $field \"widget\") \"dict\") (eq (index $field \"value_type\") \"bool\") (index $field \"is_datetime\") (eq (index $field \"widget\") \"icon\")}}'change'{{else}}'input'{{end}}", content)
    content = re.sub(r"<#if target\.has_form_editor>960px<#else>720px</#if>", '{{if (index .Target "has_form_editor")}}960px{{else}}720px{{end}}', content)
    content = re.sub(r"state\.loading<#if has_tree_parent_form> \|\| state\.treeLoading</#if>", "state.loading{{if $.HasTreeParentForm}} || state.treeLoading{{end}}", content)

    content = re.sub(r"<#if\s+([^>]+)>", lambda m: "{{if " + m.group(1).replace("has_tree_parent_form", "$.HasTreeParentForm").replace("field.required", '(index $field "required")').replace("field.is_json", '(index $field "is_json")').replace('field.value_type == "bool"', 'eq (index $field "value_type") "bool"') + "}}", content)
    content = re.sub(r"<#elseif\s+([^>]+)>", "{{else if \\1}}", content)
    content = content.replace("<#else>", "{{else}}")
    content = content.replace("</#if>", "{{end}}")

    content = re.sub(r"\$\{field\.(\w+)\}", r'{{index $field "\1"}}', content)
    return content


def main() -> int:
    for name in ("ModalForm.vue.ftl", "ChildModalForm.vue.ftl"):
        src = (BOOT / name).read_text(encoding="utf-8")
        dst = OUT / name.replace(".ftl", ".tmpl")
        dst.write_text(port(src, child=name.startswith("Child")), encoding="utf-8")
        print(f"wrote {dst.name}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
