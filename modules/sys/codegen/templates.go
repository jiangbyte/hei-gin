package codegen

// 生成 model.go（GORM 实体）。
const goModelTmpl = `// Package {{.Main.Package}} 由 HEI 代码生成器生成。
package {{.Main.Package}}

import (
{{- if .Main.HasTime }}
	"time"
{{- end }}
{{- if .Main.HasJSON }}
	"gorm.io/datatypes"
{{- end }}
)

// {{.Main.EntityName}} 对应表 {{.Main.TableName}}。
//
// Author: {{.Author}}
type {{.Main.EntityName}} struct {
{{- range .Main.Fields }}
	{{.GoName}} {{.GoType}} ` + "`" + ` + "` + "`" + `" + ` + "`" + `gorm:"column:{{.Name}};{{if .IsPrimaryKey}}primaryKey;size:64{{else if .MaxLength}}size:{{.MaxLength}}{{end}}{{if .IsJSON}};type:jsonb{{end}}{{if .IsDatetime}};type:timestamp{{end}}" json:"{{.Name}}"` + "`" + ` + "` + "`" + `" + ` + "`" + `
{{- end }}
}

// TableName 返回 {{.Main.EntityName}} 对应的数据库表名。
func ({{.Main.EntityName}}) TableName() string { return "{{.Main.TableName}}" }
`

// 生成 param.go（Add/Edit/IDs/Page）。
const goParamTmpl = `package {{.Main.Package}}

import (
{{- if .Main.HasTime }}
	"time"
{{- end }}
	"hei-gin/framework/core/schema"
)

// AddParam 创建{{.Main.BusinessName}}入参。
//
// Author: {{.Author}}
type AddParam struct {
{{- range .Main.FormFields }}
{{- if eq .PythonType "int" }}
	{{.GoName}} {{if .IsRequired}}int64{{else}}*int64{{end}} ` + "`" + ` + "` + "`" + `" + ` + "`" + `json:"{{.Name}}"{{if .IsRequired}} binding:"required"{{end}}` + "`" + ` + "` + "`" + `" + ` + "`" + `
{{- else if eq .PythonType "float" }}
	{{.GoName}} {{if .IsRequired}}float64{{else}}*float64{{end}} ` + "`" + ` + "` + "`" + `" + ` + "`" + `json:"{{.Name}}"{{if .IsRequired}} binding:"required"{{end}}` + "`" + ` + "` + "`" + `" + ` + "`" + `
{{- else if eq .PythonType "bool" }}
	{{.GoName}} bool ` + "`" + ` + "` + "`" + `" + ` + "`" + `json:"{{.Name}}"` + "`" + ` + "` + "`" + `" + ` + "`" + `
{{- else if eq .PythonType "datetime" }}
	{{.GoName}} {{if .IsRequired}}time.Time{{else}}*time.Time{{end}} ` + "`" + ` + "` + "`" + `" + ` + "`" + `json:"{{.Name}}"{{if .IsRequired}} binding:"required"{{end}}` + "`" + ` + "` + "`" + `" + ` + "`" + `
{{- else if eq .PythonType "dict" }}
	{{.GoName}} map[string]any ` + "`" + ` + "` + "`" + `" + ` + "`" + `json:"{{.Name}}"` + "`" + ` + "` + "`" + `" + ` + "`" + `
{{- else }}
	{{.GoName}} {{if .IsRequired}}string{{else}}*string{{end}} ` + "`" + ` + "` + "`" + `" + ` + "`" + `json:"{{.Name}}"{{if .IsRequired}} binding:"required"{{end}}` + "`" + ` + "` + "`" + `" + ` + "`" + `
{{- end }}
{{- end }}
}

// EditParam 更新{{.Main.BusinessName}}入参。
//
// Author: {{.Author}}
type EditParam struct {
	ID string ` + "`" + ` + "` + "`" + `" + ` + "`" + `json:"id" binding:"required"` + "`" + ` + "` + "`" + `" + ` + "`" + `
	AddParam
}

// IDsParam 批量 ID 入参。
//
// Author: {{.Author}}
type IDsParam struct {
	IDs []string ` + "`" + ` + "` + "`" + `" + ` + "`" + `json:"ids" binding:"required"` + "`" + ` + "` + "`" + `" + ` + "`" + `
}

// PageParam {{.Main.BusinessName}}分页查询。
//
// Author: {{.Author}}
type PageParam struct {
	schema.PageQuery
{{- range .Main.QueryFields }}
	{{.GoName}} string ` + "`" + ` + "` + "`" + `" + ` + "`" + `form:"{{.Name}}"` + "`" + ` + "` + "`" + `" + ` + "`" + `
{{- end }}
}
`

// 生成 result.go。
const goResultTmpl = `package {{.Main.Package}}

// DetailResult {{.Main.BusinessName}}详情（与实体一致）。
type DetailResult = {{.Main.EntityName}}
`

// 生成 repo.go。
const goRepoTmpl = `package {{.Main.Package}}

import (
	"context"

	"gorm.io/gorm"
)

// Repo {{.Main.BusinessName}}持久化。
//
// Author: {{.Author}}
type Repo struct{ db *gorm.DB }

// NewRepo 构造 Repo。
func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) with(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建{{.Main.BusinessName}}。
func (r *Repo) Create(ctx context.Context, row *{{.Main.EntityName}}) error {
	return r.with(ctx).Create(row).Error
}

// Update 更新{{.Main.BusinessName}}。
func (r *Repo) Update(ctx context.Context, id string, updates map[string]any) error {
	return r.with(ctx).Model(&{{.Main.EntityName}}{}).Where("{{.Main.PKName}} = ?", id).Updates(updates).Error
}

// DeleteByIDs 批量删除。
func (r *Repo) DeleteByIDs(ctx context.Context, ids []string) error {
	return r.with(ctx).Where("{{.Main.PKName}} IN ?", ids).Delete(&{{.Main.EntityName}}{}).Error
}

// GetByID 按主键查询。
func (r *Repo) GetByID(ctx context.Context, id string) (*{{.Main.EntityName}}, error) {
	var row {{.Main.EntityName}}
	if err := r.with(ctx).First(&row, "{{.Main.PKName}} = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Page 分页查询。
func (r *Repo) Page(ctx context.Context, p PageParam) (rows []{{.Main.EntityName}}, total int64, err error) {
	cur, size := p.Normalize()
	db := r.with(ctx).Model(&{{.Main.EntityName}}{})
{{- range .Main.QueryFields }}
{{- if eq .PythonType "str" }}
	if p.{{.GoName}} != "" {
		db = db.Where("{{.Name}} ILIKE ?", "%"+p.{{.GoName}}+"%")
	}
{{- end }}
{{- end }}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Order("created_at DESC").Offset((cur - 1) * size).Limit(size).Find(&rows).Error
	return rows, total, err
}
`

// 生成 service.go。
const goServiceTmpl = `package {{.Main.Package}}

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/framework/core/security"
	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service {{.Main.BusinessName}}服务。
//
// Author: {{.Author}}
type Service struct {
	repo *Repo
}

// NewService 构造服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 biz.{{.Main.Package}} 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "{{.ModulePath}}",
		Models: []any{&{{.Main.EntityName}}{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建{{.Main.BusinessName}}。
func (s *Service) Create(ctx context.Context, accountID string, req AddParam) error {
	row := fromAddParam(req)
	row.ID = idgen.Next()
	row.CreatedBy = &accountID
	row.UpdatedBy = &accountID
	return s.repo.Create(ctx, &row)
}

// Update 更新{{.Main.BusinessName}}。
func (s *Service) Update(ctx context.Context, accountID string, req EditParam) error {
	row := fromAddParam(req.AddParam)
	return s.repo.Update(ctx, req.ID, map[string]any{
{{- range .Main.FormFields }}
		"{{.Name}}": row.{{.GoName}},
{{- end }}
		"updated_by": accountID,
	})
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail {{.Main.BusinessName}}详情。
func (s *Service) Detail(ctx context.Context, id string) (*{{.Main.EntityName}}, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, p PageParam, sess *security.SessionPayload) (rows []{{.Main.EntityName}}, total int64, current, size int, err error) {
	current, size = p.Normalize()
	rows, total, err = s.repo.Page(ctx, p)
	return rows, total, current, size, err
}

func fromAddParam(req AddParam) {{.Main.EntityName}} {
	return {{.Main.EntityName}}{
{{- range .Main.FormFields }}
		{{.GoName}}: req.{{.GoName}},
{{- end }}
	}
}

func mustJSON(v any) datatypes.JSON {
	if v == nil {
		return datatypes.JSON([]byte("{}"))
	}
	b, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON([]byte("{}"))
	}
	return b
}
`

// 生成 handler.go。
const goHandlerTmpl = `package {{.Main.Package}}

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contextx "hei-gin/framework/core/context"
	"hei-gin/framework/core/bind"
	"hei-gin/framework/core/response"
	"hei-gin/framework/core/schema"
	"hei-gin/framework/core/security"
	"hei-gin/framework/middleware"
	"hei-gin/modules/shared"
)

func (s *Service) registerRoutes(d *shared.Deps) func(*gin.RouterGroup) {
	return func(api *gin.RouterGroup) {
		g := api.Group("{{.APIPrefix}}", middleware.RequireAccountType(security.AccountAdmin))
		g.POST("/create", middleware.RequirePermission(d.Perms, "{{.PermissionPrefix}}:create", "创建{{.Main.BusinessName}}"), s.create)
		g.POST("/update", middleware.RequirePermission(d.Perms, "{{.PermissionPrefix}}:update", "更新{{.Main.BusinessName}}"), s.update)
		g.POST("/delete", middleware.RequirePermission(d.Perms, "{{.PermissionPrefix}}:delete", "删除{{.Main.BusinessName}}"), s.delete)
		g.GET("/detail", middleware.RequirePermission(d.Perms, "{{.PermissionPrefix}}:detail", "{{.Main.BusinessName}}详情"), s.detail)
		g.GET("/page", middleware.RequirePermission(d.Perms, "{{.PermissionPrefix}}:page", "{{.Main.BusinessName}}分页"), s.page)
	}
}

func (s *Service) create(c *gin.Context) {
	var req AddParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Create(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) update(c *gin.Context) {
	var req EditParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Update(c.Request.Context(), contextx.AccountID(c.Request.Context()), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) delete(c *gin.Context) {
	var req IDsParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.Delete(c.Request.Context(), req.IDs); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}

func (s *Service) detail(c *gin.Context) {
	var q schema.IDQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	row, err := s.Detail(c.Request.Context(), q.ID)
	if err != nil {
		response.Fail(c, http.StatusNotFound, 404, "not found")
		return
	}
	response.OK(c, row)
}

func (s *Service) page(c *gin.Context) {
	var q PageParam
	_ = c.ShouldBindQuery(&q)
	rows, total, cur, size, err := s.Page(c.Request.Context(), q, contextx.Session(c.Request.Context()))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.Page(c, int64(cur), int64(size), total, rows)
}
`

// 生成 register.go。
const goRegisterTmpl = `package {{.Main.Package}}

import (
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// init 自注册 {{.ModulePath}} 模块。
func init() {
	module.Register("{{.ModulePath}}", 90, func(d *module.Deps) module.Module {
		return New(shared.FromModule(d))
	})
}
`

// 生成前端 api.ts。
const apiTsTmpl = `/**
 * 由 HEI 代码生成器生成。
 * Author: {{.Author}}
 * 生成时间：{{.GeneratedAt}}
 */
import { API_PREFIX } from '@/constants/api'
import { http } from '@/utils'

const prefix = ` + "`" + `${API_PREFIX}{{.APIPrefix}}` + "`" + `

export function page(params: any) {
	return http.get<any>(` + "`" + `${prefix}` + "`" + ` + '/page', { params })
}

export function detail(params: any) {
	return http.get<any>(` + "`" + `${prefix}` + "`" + ` + '/detail', { params })
}

export function create(data: any) {
	return http.post<any>(` + "`" + `${prefix}` + "`" + ` + '/create', data)
}

export function update(data: any) {
	return http.post<any>(` + "`" + `${prefix}` + "`" + ` + '/update', data)
}

export function remove(data: any) {
	return http.post<any>(` + "`" + `${prefix}` + "`" + ` + '/delete', data)
}
{{- if .HasTree }}

export function tree(params?: any) {
	return http.get<any>(` + "`" + `${prefix}` + "`" + ` + '/tree', { params })
}
{{- end }}
{{- if .HasSub }}

export function childPage(params: any) {
	return http.get<any>(` + "`" + `${prefix}` + "`" + ` + '/children/page', { params })
}

export function childDetail(params: any) {
	return http.get<any>(` + "`" + `${prefix}` + "`" + ` + '/children/detail', { params })
}

export function childCreate(data: any) {
	return http.post<any>(` + "`" + `${prefix}` + "`" + ` + '/children/create', data)
}

export function childUpdate(data: any) {
	return http.post<any>(` + "`" + `${prefix}` + "`" + ` + '/children/update', data)
}

export function childRemove(data: any) {
	return http.post<any>(` + "`" + `${prefix}` + "`" + ` + '/children/delete', data)
}
{{- end }}`

// 生成 api/index.ts.append 追加行。
const apiIndexAppendTmpl = `{{.ApiExportName}}
`

// 生成菜单权限 SQL。
const menuPermissionSqlTmpl = `-- 由 HEI 代码生成器生成。
-- Author: {{.Author}}
-- 生成时间：{{.GeneratedAt}}
-- 执行前请按需调整 module_id/parent_id。
BEGIN;

INSERT INTO sys_resource (id, parent_id, code, name, resource_type, module_id, path, component, icon, sort, is_visible, is_cache, is_affix, status, description, extra)
VALUES (
  '{{.Menu.MenuID}}',
{{- if .Plan.ParentResourceID }}
  '{{.Plan.ParentResourceID | sq}}',
{{- else }}
  NULL,
{{- end }}
  '{{.PermissionPrefix | replaceColon}}',
  '{{.Plan.MenuName | sq}}',
  'MENU',
{{- if .Plan.ResourceModuleID }}
  '{{.Plan.ResourceModuleID | sq}}',
{{- else }}
  NULL,
{{- end }}
  '{{.Plan.MenuPath | sq}}',
  '{{.Plan.ComponentPath | sq}}',
{{- if .Plan.Icon }}
  '{{.Plan.Icon | sq}}',
{{- else }}
  NULL,
{{- end }}
  {{.Plan.Sort}},
  true,
  false,
  false,
  'ENABLED',
{{- if .Plan.Description }}
  '{{.Plan.Description | sq}}',
{{- else }}
  NULL,
{{- end }}
  '{}'
)
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name, path = EXCLUDED.path, component = EXCLUDED.component, updated_at = now();

{{- range .Menu.Actions }}

INSERT INTO sys_resource (id, parent_id, code, name, resource_type, module_id, sort, is_visible, is_cache, is_affix, status, extra)
VALUES (
  '{{.ResourceID}}',
  '{{$.Menu.MenuID}}',
  '{{$.PermissionPrefix | replaceColon}}_{{.Key}}',
  '{{.Label}}{{$.Main.BusinessName | sq}}',
  'BUTTON',
{{- if $.Plan.ResourceModuleID }}
  '{{$.Plan.ResourceModuleID | sq}}',
{{- else }}
  NULL,
{{- end }}
  {{.Sort}},
  false,
  false,
  false,
  'ENABLED',
  '{}'
)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, updated_at = now();

INSERT INTO sys_iam_relation (id, subject_type, subject_id, relation_type, target_type, target_id, target_key, grant_mode, data_scope, custom_scope_dept_ids, is_primary, sort, status, description, extra)
VALUES (
  '{{.RelationID}}',
  'RESOURCE',
  '{{.ResourceID}}',
  'RESOURCE_PERMISSION',
  'PERMISSION',
  '',
  '{{$.PermissionPrefix}}:{{.Key}}',
  'CASCADE',
  'ALL',
  '[]',
  false,
  {{.Sort}},
  'ENABLED',
  '{{.Label}}{{$.Main.BusinessName | sq}}',
  '{}'
)
ON CONFLICT (id)
DO UPDATE SET description = EXCLUDED.description, updated_at = now();

{{- end }}
COMMIT;
`

// 生成新增/编辑弹窗。
const modalFormTmpl = `<!--
  由 HEI 代码生成器生成。
  Author: {{.Author}}
  生成时间：{{.GeneratedAt}}
-->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { {{.ApiExportName}} } from '@/api'
{{- if .Main.HasBool }}
import { wireBool } from '@/utils/wire'
{{- end }}
{{- if .Main.HasInt }}
import { wireInt } from '@/utils/wire'
{{- end }}
{{- if .Main.HasFloat }}
import { wireFloat } from '@/utils/wire'
{{- end }}
import { createRequiredRule{{- if .Main.HasDatetime }}, toApiDateTime, toFormDateTime{{- end }} } from '@/utils'
import { computed, reactive, ref } from 'vue'

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)
const defaultFormData: Record<string, any> = {
{{- range .Main.FormFields }}
  {{.Name}}: {{.VueDefault}},
{{- end }}
}
const state = reactive({
  showModal: false,
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: normalizeFormData(),
})

const modalTitle = computed(() => state.dataId ? '编辑{{.Main.BusinessName}}' : '新增{{.Main.BusinessName}}')
const rules = computed<FormRules>(() => ({
{{- range .Main.FormFields }}
{{- if .IsRequired }}
  {{.Name}}: [
{{- if .IsBool }}
    {
      validator: () => typeof state.formModel.{{.Name}} === 'boolean',
      message: '请选择{{.Label}}',
      trigger: 'change',
    },
{{- else if or (eq .PythonType "int") (eq .PythonType "float") }}
    {
      validator: () => typeof state.formModel.{{.Name}} === 'number' && Number.isFinite(state.formModel.{{.Name}}),
      message: '请输入{{.Label}}',
      trigger: ['input', 'blur'],
    },
{{- else }}
    createRequiredRule('{{.Label}}', {{if or (eq .FormWidget "dict") .IsBool .IsDatetime}}'change'{{else}}'input'{{end}}),
{{- end }}
  ],
{{- end }}
{{- end }}
}))

async function openModal(id?: string, defaults: Partial<typeof defaultFormData> = {}) {
  state.dataId = id ?? null
  state.formModel = normalizeFormData(defaults)
  state.showModal = true
  if (id) {
    await fetchDetail(id)
  }
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await {{.ApiExportName}}.detail({ id })
    state.formModel = normalizeFormData(response.data ?? {})
  } finally {
    state.loading = false
  }
}

function normalizeFormData(data: Record<string, any> = {}): Record<string, any> {
  return {
    ...defaultFormData,
    ...data,
{{- range .Main.FormFields }}
{{- if .IsBool }}
    {{.Name}}: data.{{.Name}} == null || data.{{.Name}} === '' ? defaultFormData.{{.Name}} : wireBool(String(data.{{.Name}})),
{{- end }}
{{- end }}
{{- range .Main.FormFields }}
{{- if eq .PythonType "int" }}
    {{.Name}}: data.{{.Name}} == null || data.{{.Name}} === '' ? defaultFormData.{{.Name}} : wireInt(String(data.{{.Name}})),
{{- end }}
{{- end }}
{{- range .Main.FormFields }}
{{- if eq .PythonType "float" }}
    {{.Name}}: data.{{.Name}} == null || data.{{.Name}} === '' ? defaultFormData.{{.Name}} : wireFloat(String(data.{{.Name}})),
{{- end }}
{{- end }}
{{- range .Main.FormFields }}
{{- if .IsDatetime }}
    {{.Name}}: toFormDateTime(data.{{.Name}}),
{{- end }}
{{- end }}
{{- range .Main.FormFields }}
{{- if .IsJSON }}
    {{.Name}}: stringifyJsonValue(data.{{.Name}}),
{{- end }}
{{- end }}
  }
}
{{- if .Main.HasDatetime }}

function normalizeSubmitData(data: Record<string, any>): Record<string, any> {
  return {
    ...data,
{{- range .Main.FormFields }}
{{- if .IsDatetime }}
    {{.Name}}: toApiDateTime(data.{{.Name}}),
{{- end }}
{{- end }}
{{- range .Main.FormFields }}
{{- if .IsJSON }}
    {{.Name}}: parseJsonValue(data.{{.Name}}),
{{- end }}
{{- end }}
  }
}
{{- end }}
{{- if .Main.HasJSON }}

function parseJsonValue(value: unknown) {
  const text = String(value ?? '').trim()
  if (!text) {
    return {}
  }
  const parsed = JSON.parse(text)
  if (Array.isArray(parsed) || typeof parsed !== 'object' || parsed === null) {
    throw new Error('JSON value must be an object')
  }
  return parsed
}

function isValidJsonValue(value: unknown) {
  try {
    parseJsonValue(value)
    return true
  } catch {
    return false
  }
}

function stringifyJsonValue(value: unknown) {
  if (value === undefined || value === null || value === '') {
    return '{}'
  }
  if (typeof value === 'string') {
    try {
      return JSON.stringify(JSON.parse(value), null, 2)
    } catch {
      return value
    }
  }
  return JSON.stringify(value, null, 2)
}
{{- end }}

function closeModal() {
  state.showModal = false
  state.submitLoading = false
}

async function submitForm() {
  await formRef.value?.validate()
  state.submitLoading = true
  try {
{{- if .Main.HasDatetime }}
    const payload = normalizeSubmitData(state.formModel)
{{- else }}
    const payload = state.formModel
{{- end }}
    if (state.dataId) {
      await {{.ApiExportName}}.update({ ...payload, id: state.dataId })
      window.$message.success('更新成功')
    } else {
      await {{.ApiExportName}}.create(payload)
      window.$message.success('创建成功')
    }
    emit('saved')
    closeModal()
  } finally {
    state.submitLoading = false
  }
}

defineExpose({
  openModal,
})
</script>

<template>
  <NModal
    v-model:show="state.showModal"
    preset="card"
    draggable
    :mask-closable="false"
    :title="modalTitle"
    style="width: 720px"
    :segmented="{ content: true, action: true }"
  >
    <NSpin :show="state.loading">
      <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
        <NForm ref="formRef" :model="state.formModel" :rules="rules" label-placement="left" label-width="110" :disabled="state.loading || state.submitLoading">
{{- range .Main.FormFields }}
          <NFormItem label="{{.Label}}" path="{{.Name}}">
{{- if eq .PythonType "int" }}
            <NInputNumber v-model:value="state.formModel.{{.Name}}" class="w-full" />
{{- else if eq .PythonType "float" }}
            <NInputNumber v-model:value="state.formModel.{{.Name}}" class="w-full" />
{{- else if eq .PythonType "bool" }}
            <NSwitch v-model:value="state.formModel.{{.Name}}" />
{{- else if .IsJSON }}
            <NInput v-model:value="state.formModel.{{.Name}}" type="textarea" :autosize="{ minRows: 4, maxRows: 12 }" />
{{- else if .IsDatetime }}
            <NDatePicker v-model:formatted-value="state.formModel.{{.Name}}" type="datetime" value-format="yyyy-MM-dd HH:mm:ss" class="w-full" clearable />
{{- else if eq .FormWidget "textarea" }}
            <NInput v-model:value="state.formModel.{{.Name}}" type="textarea" :autosize="{ minRows: 3, maxRows: 8 }" />
{{- else if eq .FormWidget "dict" }}
            <DictSelect v-model="state.formModel.{{.Name}}" dict-code="{{.DictCode}}" />
{{- else }}
            <NInput v-model:value="state.formModel.{{.Name}}" />
{{- end }}
          </NFormItem>
{{- end }}
        </NForm>
      </NScrollbar>
    </NSpin>

    <template #action>
      <NSpace justify="end">
        <NButton @click="closeModal">取消</NButton>
        <NButton type="primary" :loading="state.submitLoading" @click="submitForm">确认</NButton>
      </NSpace>
    </template>
  </NModal>
</template>`

// 生成前端列表页 index.vue。
const indexVueTmpl = `<!--
  由 HEI 代码生成器生成。
  Author: {{.Author}}
  生成时间：{{.GeneratedAt}}
-->

<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { {{.ApiExportName}} } from '@/api'
import { formatDateTime, hasPermission, normalizeSearchValues, renderButtonIcon } from '@/utils'
import { NButton, NFlex, NIcon } from 'naive-ui'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'
import ModalDetail from './components/ModalDetail.vue'
import ModalForm from './components/ModalForm.vue'

const formModalRef = ref<any>(null)
const detailModalRef = ref<any>(null)
const state = reactive({
  rows: [] as any[],
  total: 0,
  loading: false,
  searchValues: {} as any,
  checkedRowKeys: [] as string[],
  page: 1,
  pageSize: 20,
})

const hasCheckedRows = computed(() => state.checkedRowKeys.length > 0)

const searchForm = createProSearchForm<any>({
  defaultCollapsed: true,
  onSubmit(values) {
    state.searchValues = normalizeSearchValues(values)
    state.page = 1
    fetchPage()
  },
  onReset() {
    state.searchValues = {}
    state.page = 1
    fetchPage()
  },
})

const searchColumns = computed<ProSearchFormColumns<any>>(() => [
{{- range .Main.QueryFields }}
  { title: '{{.Label}}', path: '{{.Name}}', field: 'input' },
{{- end }}
])

const pagination = computed<PaginationProps>(() => ({
  page: state.page,
  pageSize: state.pageSize,
  itemCount: state.total,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: ({ itemCount }) => ` + "`" + `${itemCount} 条` + "`" + `,
  onUpdatePage: (value) => {
    state.page = value
    fetchPage()
  },
  onUpdatePageSize: (value) => {
    state.pageSize = value
    state.page = 1
    fetchPage()
  },
})

const tableColumns = computed<ProDataTableColumns<any>>(() => [
  { type: 'selection', fixed: 'left' },
{{- range .Main.TableFields }}
{{- if .DictCode }}
  {
    title: '{{.Label}}',
    path: '{{.Name}}',
    width: 150,
    ellipsis: { tooltip: true },
    render: row => (
      <span>{row.{{.Name}}}</span>
    ),
  },
{{- else if .IsDatetime }}
  { title: '{{.Label}}', path: '{{.Name}}', width: 190, render: row => formatDateTime(row.{{.Name}}) },
{{- else }}
  { title: '{{.Label}}', path: '{{.Name}}', width: 150, ellipsis: { tooltip: true } },
{{- end }}
{{- end }}
  {
    title: '创建时间',
    path: 'created_at',
    width: 190,
    render: (row) => formatDateTime(row.created_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    fixed: 'right',
    render: (row) => (
      <NFlex size={12}>
        {hasPermission('{{.PermissionPrefix}}:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => formModalRef.value?.openModal(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('{{.PermissionPrefix}}:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => detailModalRef.value?.openModal(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('{{.PermissionPrefix}}:delete') ? (
          <NButton type="error" size="small" text={true} onClick={() => confirmDelete(row.id)}>
            {renderButtonIcon('icon-park-outline:delete')}
          </NButton>
        ) : null}
      </NFlex>
    ),
  },
])

onMounted(() => {
  fetchPage()
})

async function fetchPage() {
  state.loading = true
  try {
    const response = await {{.ApiExportName}}.page({
      current: state.page,
      size: state.pageSize,
      ...state.searchValues,
    })
    state.rows = response.data?.records ?? []
    state.total = response.data?.total ?? 0
  } finally {
    state.loading = false
  }
}

function handleCheckedRowKeys(keys: Array<string | number>) {
  state.checkedRowKeys = keys.map(String)
}

function confirmDelete(value: string | string[]) {
  const ids = Array.isArray(value) ? value : [value]
  if (!ids.length) {
    return
  }
  window.$dialog.warning({
    title: ids.length > 1 ? '批量删除' : '删除',
    content: ids.length > 1 ? ` + "`" + `${ids.length} 条记录?` + "`" + ` : '删除该记录?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteRows(ids),
  })
}

async function deleteRows(ids: string[]) {
  await {{.ApiExportName}}.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))
  window.$message.success('删除成功')
  await fetchPage()
}
</script>

<template>
  <NFlex class="h-full min-h-0" vertical>
    <ProCard content-class="pb-0!">
      <ProSearchForm
        :form="searchForm"
        :columns="searchColumns"
        :reset-button-props="{ content: '重置' }"
        :search-button-props="{ content: '搜索' }"
        :collapse-button-props="{ content: searchForm.collapsed.value ? '展开' : '收起' }"
      />
    </ProCard>

    <ProDataTable
      class="min-h-0 flex-1"
      remote
      title="{{.Main.BusinessName}}"
      row-key="id"
      :scroll-x="1400"
      :columns="tableColumns"
      :data="state.rows"
      :loading="state.loading"
      :pagination="pagination"
      :checked-row-keys="state.checkedRowKeys"
      :on-update-checked-row-keys="handleCheckedRowKeys"
    >
      <template #toolbar>
        <NFlex>
          <NButton v-if="hasPermission('{{.PermissionPrefix}}:create')" type="primary" text title="新增" @click="formModalRef.value?.openModal()">
            <template #icon><NIcon><Icon icon="icon-park-outline:plus" /></NIcon></template>
          </NButton>
          <NButton text title="刷新" :loading="state.loading" @click="fetchPage">
            <template #icon><NIcon><Icon icon="icon-park-outline:reload" /></NIcon></template>
          </NButton>
          <NButton v-if="hasPermission('{{.PermissionPrefix}}:delete')" type="error" text title="批量删除" :disabled="!hasCheckedRows" @click="confirmDelete(state.checkedRowKeys)">
            <template #icon><NIcon><Icon icon="icon-park-outline:delete" /></NIcon></template>
          </NButton>
        </NFlex>
      </template>
    </ProDataTable>

    <ModalForm ref="formModalRef" @saved="fetchPage" />
    <ModalDetail ref="detailModalRef" />
  </NFlex>
</template>`

// 生成详情弹窗。
const modalDetailTmpl = `<!--
  由 HEI 代码生成器生成。
  Author: {{.Author}}
  生成时间：{{.GeneratedAt}}
-->

<script setup lang="ts">
import { {{.ApiExportName}} } from '@/api'
import { displayValue, formatDateTime } from '@/utils'
import { reactive } from 'vue'

const state = reactive({
  showModal: false,
  loading: false,
  detail: {} as any,
})

async function openModal(id: string) {
  state.detail = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await {{.ApiExportName}}.detail({ id })
    state.detail = response.data ?? {}
  } finally {
    state.loading = false
  }
}
{{- if .Main.HasJSON }}

function formatJsonValue(value: unknown) {
  if (value === undefined || value === null || value === '') {
    return '{}'
  }
  if (typeof value === 'string') {
    try {
      return JSON.stringify(JSON.parse(value), null, 2)
    } catch {
      return value
    }
  }
  return JSON.stringify(value, null, 2)
}
{{- end }}

defineExpose({
  openModal,
})
</script>

<template>
  <NModal v-model:show="state.showModal" preset="card" draggable :mask-closable="false" title="{{.Main.BusinessName}}详情" style="width: 680px">
    <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions label-placement="left" bordered :column="1">
{{- range .Main.DetailFields }}
          <NDescriptionsItem label="{{.Label}}">
{{- if .IsDatetime }}
            {{"{{ formatDateTime(state.detail.{{.Name}}) }}"}}
{{- else if .IsJSON }}
            <NCode :code="formatJsonValue(state.detail.{{.Name}})" language="json" word-wrap />
{{- else }}
            {{"{{ displayValue(state.detail.{{.Name}}) }}"}}
{{- end }}
          </NDescriptionsItem>
{{- end }}
          <NDescriptionsItem label="创建时间">{{"{{ formatDateTime(state.detail.created_at) }}"}}</NDescriptionsItem>
          <NDescriptionsItem label="创建人">{{"{{ displayValue(state.detail.created_by) }}"}}</NDescriptionsItem>
          <NDescriptionsItem label="更新时间">{{"{{ formatDateTime(state.detail.updated_at) }}"}}</NDescriptionsItem>
          <NDescriptionsItem label="更新人">{{"{{ displayValue(state.detail.updated_by) }}"}}</NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>`

// 生成子表新增/编辑弹窗。
const childModalFormTmpl = `<!--
  由 HEI 代码生成器生成。
  Author: {{.Author}}
  生成时间：{{.GeneratedAt}}
-->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { {{.ApiExportName}} } from '@/api'
{{- if .Sub.HasBool }}
import { wireBool } from '@/utils/wire'
{{- end }}
{{- if .Sub.HasInt }}
import { wireInt } from '@/utils/wire'
{{- end }}
{{- if .Sub.HasFloat }}
import { wireFloat } from '@/utils/wire'
{{- end }}
import { createRequiredRule{{- if .Sub.HasDatetime }}, toApiDateTime, toFormDateTime{{- end }} } from '@/utils'
import { computed, reactive, ref } from 'vue'

const props = defineProps<{
  masterId?: string | null
}>()
const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)
const defaultFormData: Record<string, any> = {
{{- range .Sub.FormFields }}
  {{.Name}}: {{.VueDefault}},
{{- end }}
}
const state = reactive({
  showModal: false,
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: normalizeFormData(),
})

const modalTitle = computed(() => state.dataId ? '编辑{{.Sub.BusinessName}}' : '新增{{.Sub.BusinessName}}')
const rules = computed<FormRules>(() => ({
{{- range .Sub.FormFields }}
{{- if .IsRequired }}
  {{.Name}}: [
{{- if .IsBool }}
    {
      validator: () => typeof state.formModel.{{.Name}} === 'boolean',
      message: '请选择{{.Label}}',
      trigger: 'change',
    },
{{- else if or (eq .PythonType "int") (eq .PythonType "float") }}
    {
      validator: () => typeof state.formModel.{{.Name}} === 'number' && Number.isFinite(state.formModel.{{.Name}}),
      message: '请输入{{.Label}}',
      trigger: ['input', 'blur'],
    },
{{- else }}
    createRequiredRule('{{.Label}}', {{if or (eq .FormWidget "dict") .IsBool .IsDatetime}}'change'{{else}}'input'{{end}}),
{{- end }}
  ],
{{- end }}
{{- end }}
}))

async function openModal(id?: string, defaults: Partial<typeof defaultFormData> = {}) {
  state.dataId = id ?? null
  state.formModel = normalizeFormData(defaults)
  state.showModal = true
  if (id) {
    await fetchDetail(id)
  }
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await {{.ApiExportName}}.childDetail({ id })
    state.formModel = normalizeFormData(response.data ?? {})
  } finally {
    state.loading = false
  }
}

function normalizeFormData(data: Record<string, any> = {}): Record<string, any> {
  return {
    ...defaultFormData,
    ...data,
{{- range .Sub.FormFields }}
{{- if .IsBool }}
    {{.Name}}: data.{{.Name}} == null || data.{{.Name}} === '' ? defaultFormData.{{.Name}} : wireBool(String(data.{{.Name}})),
{{- end }}
{{- end }}
{{- range .Sub.FormFields }}
{{- if eq .PythonType "int" }}
    {{.Name}}: data.{{.Name}} == null || data.{{.Name}} === '' ? defaultFormData.{{.Name}} : wireInt(String(data.{{.Name}})),
{{- end }}
{{- end }}
{{- range .Sub.FormFields }}
{{- if eq .PythonType "float" }}
    {{.Name}}: data.{{.Name}} == null || data.{{.Name}} === '' ? defaultFormData.{{.Name}} : wireFloat(String(data.{{.Name}})),
{{- end }}
{{- end }}
{{- range .Sub.FormFields }}
{{- if .IsDatetime }}
    {{.Name}}: toFormDateTime(data.{{.Name}}),
{{- end }}
{{- end }}
  }
}

function closeModal() {
  state.showModal = false
  state.submitLoading = false
}

async function submitForm() {
  await formRef.value?.validate()
  state.submitLoading = true
  try {
    const payload: Record<string, any> = { ...state.formModel }
{{- if .Sub.HasDatetime }}
{{- range .Sub.FormFields }}
{{- if .IsDatetime }}
    payload.{{.Name}} = toApiDateTime(payload.{{.Name}})
{{- end }}
{{- end }}
{{- end }}
    if (props.masterId) {
      payload.{{.Plan.SubForeignKey | toSnake}} = props.masterId
    }
    if (state.dataId) {
      await {{.ApiExportName}}.childUpdate({ ...payload, id: state.dataId })
      window.$message.success('更新成功')
    } else {
      await {{.ApiExportName}}.childCreate(payload)
      window.$message.success('创建成功')
    }
    emit('saved')
    closeModal()
  } finally {
    state.submitLoading = false
  }
}

defineExpose({
  openModal,
})
</script>

<template>
  <NModal
    v-model:show="state.showModal"
    preset="card"
    draggable
    :mask-closable="false"
    :title="modalTitle"
    style="width: 720px"
    :segmented="{ content: true, action: true }"
  >
    <NSpin :show="state.loading">
      <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
        <NForm ref="formRef" :model="state.formModel" :rules="rules" label-placement="left" label-width="110" :disabled="state.loading || state.submitLoading">
{{- range .Sub.FormFields }}
          <NFormItem label="{{.Label}}" path="{{.Name}}">
{{- if eq .PythonType "int" }}
            <NInputNumber v-model:value="state.formModel.{{.Name}}" class="w-full" />
{{- else if eq .PythonType "float" }}
            <NInputNumber v-model:value="state.formModel.{{.Name}}" class="w-full" />
{{- else if eq .PythonType "bool" }}
            <NSwitch v-model:value="state.formModel.{{.Name}}" />
{{- else if .IsJSON }}
            <NInput v-model:value="state.formModel.{{.Name}}" type="textarea" :autosize="{ minRows: 4, maxRows: 12 }" />
{{- else if .IsDatetime }}
            <NDatePicker v-model:formatted-value="state.formModel.{{.Name}}" type="datetime" value-format="yyyy-MM-dd HH:mm:ss" class="w-full" clearable />
{{- else if eq .FormWidget "textarea" }}
            <NInput v-model:value="state.formModel.{{.Name}}" type="textarea" :autosize="{ minRows: 3, maxRows: 8 }" />
{{- else if eq .FormWidget "dict" }}
            <DictSelect v-model="state.formModel.{{.Name}}" dict-code="{{.DictCode}}" />
{{- else }}
            <NInput v-model:value="state.formModel.{{.Name}}" />
{{- end }}
          </NFormItem>
{{- end }}
        </NForm>
      </NScrollbar>
    </NSpin>

    <template #action>
      <NSpace justify="end">
        <NButton @click="closeModal">取消</NButton>
        <NButton type="primary" :loading="state.submitLoading" @click="submitForm">确认</NButton>
      </NSpace>
    </template>
  </NModal>
</template>`

// 生成子表详情弹窗。
const childModalDetailTmpl = `<!--
  由 HEI 代码生成器生成。
  Author: {{.Author}}
  生成时间：{{.GeneratedAt}}
-->

<script setup lang="ts">
import { {{.ApiExportName}} } from '@/api'
import { displayValue, formatDateTime } from '@/utils'
import { reactive } from 'vue'

const state = reactive({
  showModal: false,
  loading: false,
  detail: {} as any,
})

async function openModal(id: string) {
  state.detail = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await {{.ApiExportName}}.childDetail({ id })
    state.detail = response.data ?? {}
  } finally {
    state.loading = false
  }
}
{{- if .Sub.HasJSON }}

function formatJsonValue(value: unknown) {
  if (value === undefined || value === null || value === '') {
    return '{}'
  }
  if (typeof value === 'string') {
    try {
      return JSON.stringify(JSON.parse(value), null, 2)
    } catch {
      return value
    }
  }
  return JSON.stringify(value, null, 2)
}
{{- end }}

defineExpose({
  openModal,
})
</script>

<template>
  <NModal v-model:show="state.showModal" preset="card" draggable :mask-closable="false" title="{{.Sub.BusinessName}}详情" style="width: 680px">
    <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions label-placement="left" bordered :column="1">
{{- range .Sub.DetailFields }}
          <NDescriptionsItem label="{{.Label}}">
{{- if .IsDatetime }}
            {{"{{ formatDateTime(state.detail.{{.Name}}) }}"}}
{{- else if .IsJSON }}
            <NCode :code="formatJsonValue(state.detail.{{.Name}})" language="json" word-wrap />
{{- else }}
            {{"{{ displayValue(state.detail.{{.Name}}) }}"}}
{{- end }}
          </NDescriptionsItem>
{{- end }}
          <NDescriptionsItem label="创建时间">{{"{{ formatDateTime(state.detail.created_at) }}"}}</NDescriptionsItem>
          <NDescriptionsItem label="创建人">{{"{{ displayValue(state.detail.created_by) }}"}}</NDescriptionsItem>
          <NDescriptionsItem label="更新时间">{{"{{ formatDateTime(state.detail.updated_at) }}"}}</NDescriptionsItem>
          <NDescriptionsItem label="更新人">{{"{{ displayValue(state.detail.updated_by) }}"}}</NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>`
