<!-- Author: Charlie -->

<script setup lang="tsx">
import type { FormRules, PaginationProps } from 'naive-ui'
import type { SelectMixedOption } from 'naive-ui/es/select/src/interface'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { codegenApi, resourceModuleApi } from '@/api'
import IconSelect from '@/components/common/IconSelect.vue'
import MonacoPreview from '@/components/editor/MonacoPreview.vue'
import {
  createRequiredRule,
  dictDataAll,
  formatDateTime,
  hasPermission,
  normalizeSearchValues,
  refreshDict,
  renderButtonIcon,
} from '@/utils'
import { NButton, NCheckbox, NFlex, NInput, NInputNumber, NSelect, NTreeSelect } from 'naive-ui'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'

const defaultForm = {
  name: '',
  gen_type: 'TABLE',
  author: '',
  description: '',
  main_table: '',
  main_pk: 'id',
  main_entity_name: '',
  main_module_path: '',
  main_business_name: '',
  api_prefix: '',
  permission_prefix: '',
  resource_module_id: null as string | null,
  parent_resource_id: null as string | null,
  menu_name: '',
  menu_path: '',
  component_path: '',
  icon: 'icon-park-outline:code',
  sort: 99,
  tree_parent_field: '',
  tree_label_field: '',
  sub_table: '',
  sub_pk: '',
  sub_foreign_key: '',
  sub_entity_name: '',
  sub_business_name: '',
}

const genTypeOptions = [
  { label: '普通表', value: 'TABLE' },
  { label: '树表', value: 'TREE' },
  { label: '左树右表', value: 'LEFT_TREE_TABLE' },
  { label: '主子表', value: 'MASTER_DETAIL' },
]

const widgetOptions = [
  { label: '输入框', value: 'input' },
  { label: '多行文本', value: 'textarea' },
  { label: '数字', value: 'number' },
  { label: '开关', value: 'switch' },
  { label: '字典', value: 'dict' },
  { label: '日期时间', value: 'datetime' },
]

const operatorOptions = [
  { label: '不查询', value: null },
  { label: '等于', value: 'EQ' },
  { label: '模糊', value: 'LIKE' },
] as SelectMixedOption[]

interface ColumnOption {
  label: string
  value: string
  isPrimaryKey: boolean
}

const formRef = ref<any>(null)
const state = reactive({
  rows: [] as any[],
  total: 0,
  loading: false,
  searchValues: {} as any,
  checkedRowKeys: [] as string[],
  page: 1,
  pageSize: 20,
  showForm: false,
  formLoading: false,
  submitLoading: false,
  editingId: null as string | null,
  form: { ...defaultForm },
  tableOptions: [] as any[],
  mainColumnOptions: [] as ColumnOption[],
  subColumnOptions: [] as ColumnOption[],
  moduleOptions: [] as any[],
  parentOptions: [] as any[],
  showFields: false,
  fieldPlanId: '',
  fieldRows: [] as any[],
  fieldLoading: false,
  fieldSaving: false,
  showPreview: false,
  previewFiles: [] as any[],
  previewPath: '',
  previewLoading: false,
  downloadingId: '',
})

const needsTree = computed(() => ['TREE', 'LEFT_TREE_TABLE'].includes(state.form.gen_type))
const needsSub = computed(() => ['LEFT_TREE_TABLE', 'MASTER_DETAIL'].includes(state.form.gen_type))
const previewEditorOptions = computed(() => ({
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  wordWrap: 'off' as const,
  automaticLayout: true,
}))
const hasCheckedRows = computed(() => state.checkedRowKeys.length > 0)
const dictCodeOptions = computed(() => toDictCodeTreeOptions(dictDataAll()))
const formRules = computed<FormRules>(() => ({
  name: createRequiredRule('方案名称', 'input'),
  author: createRequiredRule('作者', 'input'),
  gen_type: createRequiredRule('生成类型', 'change'),
  main_table: createRequiredRule('主表', 'change'),
  main_pk: createRequiredRule('主键', 'change'),
  main_entity_name: createRequiredRule('主实体类', 'input'),
  main_module_path: createRequiredRule('后端模块路径', 'input'),
  main_business_name: createRequiredRule('业务名称', 'input'),
  api_prefix: createRequiredRule('接口前缀', 'input'),
  permission_prefix: createRequiredRule('权限前缀', 'input'),
  menu_name: createRequiredRule('菜单名称', 'input'),
  menu_path: createRequiredRule('菜单路径', 'input'),
  component_path: createRequiredRule('组件路径', 'input'),
}))

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
  { title: '方案名称', path: 'name', field: 'input' },
  { title: '主表', path: 'main_table', field: 'input' },
  { title: '生成类型', path: 'gen_type', field: 'select', fieldProps: { options: genTypeOptions } },
])

const pagination = computed<PaginationProps>(() => ({
  page: state.page,
  pageSize: state.pageSize,
  itemCount: state.total,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: ({ itemCount }) => `${itemCount} 条`,
  onUpdatePage: (value) => {
    state.page = value
    fetchPage()
  },
  onUpdatePageSize: (value) => {
    state.pageSize = value
    state.page = 1
    fetchPage()
  },
}))

const tableColumns = computed<ProDataTableColumns<any>>(() => [
  { type: 'selection', fixed: 'left' },
  { title: '方案名称', path: 'name', width: 180, ellipsis: { tooltip: true } },
  {
    title: '生成类型',
    path: 'gen_type',
    width: 130,
    render: (row) =>
      genTypeOptions.find((item) => item.value === row.gen_type)?.label ?? row.gen_type,
  },
  { title: '作者', path: 'author', width: 120, ellipsis: { tooltip: true } },
  { title: '主表', path: 'main_table', width: 180, ellipsis: { tooltip: true } },
  { title: '子表', path: 'sub_table', width: 180, ellipsis: { tooltip: true } },
  { title: '模块路径', path: 'main_module_path', width: 220, ellipsis: { tooltip: true } },
  { title: '权限前缀', path: 'permission_prefix', width: 170, ellipsis: { tooltip: true } },
  {
    title: '更新时间',
    path: 'updated_at',
    width: 190,
    render: (row) => formatDateTime(row.updated_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 230,
    fixed: 'right',
    render: (row) => (
      <NFlex size={12}>
        {hasPermission('sys:codegen:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openForm(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('sys:codegen:update') ? (
          <NButton type="info" size="small" text={true} onClick={() => openFields(row.id)}>
            {renderButtonIcon('icon-park-outline:list-view')}
          </NButton>
        ) : null}
        {hasPermission('sys:codegen:preview') ? (
          <NButton type="info" size="small" text={true} onClick={() => openPreview(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('sys:codegen:download') ? (
          <NButton
            type="primary"
            size="small"
            text={true}
            loading={state.downloadingId === row.id}
            onClick={() => download(row.id)}
          >
            {renderButtonIcon('icon-park-outline:download')}
          </NButton>
        ) : null}
        {hasPermission('sys:codegen:delete') ? (
          <NButton type="error" size="small" text={true} onClick={() => confirmDelete(row.id)}>
            {renderButtonIcon('icon-park-outline:delete')}
          </NButton>
        ) : null}
      </NFlex>
    ),
  },
])

const fieldColumns = computed(() => [
  { title: '表', key: 'table_role', width: 80 },
  { title: '字段', key: 'column_name', width: 160, ellipsis: { tooltip: true } },
  {
    title: '注释',
    key: 'column_comment',
    width: 160,
    render: (row: any) => (
      <NInput
        value={row.column_comment}
        onUpdateValue={(value: string) => (row.column_comment = value)}
      />
    ),
  },
  { title: 'Python', key: 'python_type', width: 130 },
  { title: 'TS', key: 'typescript_type', width: 140 },
  {
    title: '控件',
    key: 'form_widget',
    width: 130,
    render: (row: any) => (
      <NSelect
        value={row.form_widget}
        options={widgetOptions}
        onUpdateValue={(value: string) => handleWidgetUpdate(row, value)}
      />
    ),
  },
  {
    title: '字典',
    key: 'dict_code',
    width: 220,
    render: (row: any) => (
      <NTreeSelect
        value={row.dict_code || null}
        options={dictCodeOptions.value}
        clearable={true}
        filterable={true}
        disabled={row.form_widget !== 'dict'}
        placeholder="选择字典"
        onUpdateValue={(value: string | number | Array<string | number> | null) =>
          handleDictCodeUpdate(row, value)
        }
      />
    ),
  },
  {
    title: '查询',
    key: 'query_operator',
    width: 120,
    render: (row: any) => (
      <NSelect
        value={row.query_operator}
        options={operatorOptions}
        onUpdateValue={(value: string | null) => (row.query_operator = value)}
      />
    ),
  },
  {
    title: '表格',
    key: 'show_in_table',
    width: 80,
    render: (row: any) => (
      <NCheckbox
        checked={row.show_in_table}
        onUpdateChecked={(value: boolean) => (row.show_in_table = value)}
      />
    ),
  },
  {
    title: '表单',
    key: 'show_in_form',
    width: 80,
    render: (row: any) => (
      <NCheckbox
        checked={row.show_in_form}
        onUpdateChecked={(value: boolean) => (row.show_in_form = value)}
      />
    ),
  },
  {
    title: '详情',
    key: 'show_in_detail',
    width: 80,
    render: (row: any) => (
      <NCheckbox
        checked={row.show_in_detail}
        onUpdateChecked={(value: boolean) => (row.show_in_detail = value)}
      />
    ),
  },
  {
    title: '检索',
    key: 'show_in_query',
    width: 80,
    render: (row: any) => (
      <NCheckbox
        checked={row.show_in_query}
        onUpdateChecked={(value: boolean) => (row.show_in_query = value)}
      />
    ),
  },
  {
    title: '必填',
    key: 'is_required',
    width: 80,
    render: (row: any) => (
      <NCheckbox
        checked={row.is_required}
        onUpdateChecked={(value: boolean) => (row.is_required = value)}
      />
    ),
  },
  {
    title: '排序',
    key: 'sort',
    width: 100,
    render: (row: any) => (
      <NInputNumber
        value={row.sort}
        min={0}
        onUpdateValue={(value: number | null) => (row.sort = value ?? 0)}
      />
    ),
  },
])

onMounted(async () => {
  await Promise.all([
    fetchPage(),
    fetchTables(),
    fetchModules(),
    fetchParentResources(),
    refreshDict(),
  ])
})

async function fetchPage() {
  state.loading = true
  try {
    const response = await codegenApi.page({
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

async function fetchTables() {
  const response = await codegenApi.tables()
  state.tableOptions = (response.data ?? []).map((item: any) => ({
    label: item.table_comment ? `${item.table_name} - ${item.table_comment}` : item.table_name,
    value: item.table_name,
  }))
}

async function fetchModules() {
  const response = await resourceModuleApi.selector({ client: 'ADMIN' })
  state.moduleOptions = (response.data ?? []).map((item: any) => ({
    label: item.name,
    value: String(item.id),
  }))
}

async function fetchParentResources(moduleId = state.form.resource_module_id) {
  const response = await codegenApi.parentResources({ module_id: moduleId || undefined })
  state.parentOptions = toTreeOptions(response.data ?? [])
}

async function fetchColumns(tableName: string, target: 'main' | 'sub') {
  if (!tableName) {
    if (target === 'main') {
      state.mainColumnOptions = []
    } else {
      state.subColumnOptions = []
    }
    return
  }
  const response = await codegenApi.tableColumns({ table_name: tableName })
  const options = (response.data ?? []).map((item: any) => ({
    label: `${item.column_name}${item.is_primary_key ? ' (主键)' : ''}${item.column_comment ? ` - ${item.column_comment}` : ''}`,
    value: item.column_name,
    isPrimaryKey: Boolean(item.is_primary_key),
  }))
  if (target === 'main') {
    state.mainColumnOptions = options
  } else {
    state.subColumnOptions = options
  }
}

function openCreateForm() {
  state.editingId = null
  state.form = { ...defaultForm }
  state.mainColumnOptions = []
  state.subColumnOptions = []
  state.showForm = true
}

async function openForm(id: string) {
  state.editingId = id
  state.showForm = true
  state.formLoading = true
  try {
    const response = await codegenApi.detail({ id })
    state.form = { ...defaultForm, ...response.data }
    await Promise.all([
      fetchColumns(state.form.main_table, 'main'),
      state.form.sub_table ? fetchColumns(state.form.sub_table, 'sub') : Promise.resolve(),
      fetchParentResources(state.form.resource_module_id),
    ])
  } finally {
    state.formLoading = false
  }
}

async function handleMainTableUpdate(value: string) {
  state.form.main_table = value
  await fetchColumns(value, 'main')
  state.form.main_pk = resolvePrimaryColumn('main', state.form.main_pk)
  if (!state.editingId) {
    const entity = toPascalCase(value)
    const backendPath = value.replaceAll('-', '_')
    const routePath = value.replaceAll('_', '-')
    state.form.main_entity_name = entity
    state.form.main_business_name = entity
    state.form.main_module_path = `biz/${backendPath}`
    state.form.api_prefix = `/biz/${routePath}`
    state.form.permission_prefix = `biz:${routePath.replaceAll('-', '')}`
    state.form.menu_name = entity
    state.form.menu_path = `/biz/${routePath}`
    state.form.component_path = `/biz/${routePath}/index.vue`
  }
}

async function handleSubTableUpdate(value: string) {
  state.form.sub_table = value
  await fetchColumns(value, 'sub')
  state.form.sub_pk = resolvePrimaryColumn('sub', state.form.sub_pk)
  if (!state.editingId) {
    state.form.sub_entity_name = toPascalCase(value)
    state.form.sub_business_name = toPascalCase(value)
  }
}

async function handleResourceModuleUpdate(value: string | null) {
  state.form.resource_module_id = value
  state.form.parent_resource_id = null
  await fetchParentResources(value)
}

async function submitForm() {
  await formRef.value?.validate()
  state.submitLoading = true
  try {
    const payload = { ...state.form }
    if (!needsTree.value) {
      payload.tree_parent_field = null as any
      payload.tree_label_field = null as any
    }
    if (!needsSub.value) {
      payload.sub_table = null as any
      payload.sub_pk = null as any
      payload.sub_foreign_key = null as any
      payload.sub_entity_name = null as any
      payload.sub_business_name = null as any
    }
    if (state.editingId) {
      await codegenApi.update({ ...payload, id: state.editingId })
      window.$message.success('更新成功')
    } else {
      await codegenApi.create(payload)
      window.$message.success('创建成功')
    }
    state.showForm = false
    await fetchPage()
  } finally {
    state.submitLoading = false
  }
}

async function openFields(id: string) {
  state.fieldPlanId = id
  state.showFields = true
  state.fieldLoading = true
  try {
    const response = await codegenApi.fields({ plan_id: id })
    state.fieldRows = response.data ?? []
  } finally {
    state.fieldLoading = false
  }
}

async function saveFields() {
  state.fieldSaving = true
  try {
    await codegenApi.updateFieldsBatch({ plan_id: state.fieldPlanId, fields: state.fieldRows })
    window.$message.success('字段配置已保存')
    state.showFields = false
  } finally {
    state.fieldSaving = false
  }
}

async function openPreview(id: string) {
  state.showPreview = true
  state.previewLoading = true
  state.previewFiles = []
  state.previewPath = ''
  try {
    const response = await codegenApi.preview({ id })
    state.previewFiles = response.data?.files ?? []
    state.previewPath = state.previewFiles[0]?.path ?? ''
  } finally {
    state.previewLoading = false
  }
}

async function download(id: string) {
  state.downloadingId = id
  try {
    await codegenApi.downloadZip(id)
  } finally {
    state.downloadingId = ''
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
    content: ids.length > 1 ? `删除 ${ids.length} 个生成方案?` : '删除该生成方案?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteRows(ids),
  })
}

async function deleteRows(ids: string[]) {
  await codegenApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))
  window.$message.success('删除成功')
  await fetchPage()
}

function findColumn(target: 'main' | 'sub', name: string) {
  const options = target === 'main' ? state.mainColumnOptions : state.subColumnOptions
  return options.find((item) => item.value === name)?.value
}

function resolvePrimaryColumn(target: 'main' | 'sub', currentValue?: string | null) {
  const options = target === 'main' ? state.mainColumnOptions : state.subColumnOptions
  const primaryOption = options.find((item) => item.isPrimaryKey)
  if (primaryOption) {
    return primaryOption.value
  }
  if (currentValue && options.some((item) => item.value === currentValue)) {
    return currentValue
  }
  return findColumn(target, 'id') || options[0]?.value || ''
}

function previewTabLabel(path: string) {
  const parts = path.split('/').filter(Boolean)
  return parts.slice(-2).join('/') || path
}

function toPascalCase(value: string) {
  return value
    .split(/[_\-\s]+/)
    .filter(Boolean)
    .map((item) => item.charAt(0).toUpperCase() + item.slice(1))
    .join('')
}

function toTreeOptions(items: any[]): any[] {
  return items
    .filter((item) => item.resource_type === 'CATALOG')
    .map((item) => ({
      label: item.name,
      key: String(item.id),
      children: item.children?.length ? toTreeOptions(item.children) : undefined,
    }))
}

function toDictCodeTreeOptions(items: any[]): any[] {
  return items.map((item) => ({
    label: toDictOptionLabel(item),
    key: item.code,
    disabled: item.status !== undefined && item.status !== null && item.status !== 'ENABLED',
    children: item.children?.length
      ? toDisabledDictItemOptions(item.children, item.code)
      : undefined,
  }))
}

function toDisabledDictItemOptions(items: any[], parentCode: string): any[] {
  return items.map((item) => ({
    label: toDictOptionLabel(item),
    key: `${parentCode}:${item.code}`,
    disabled: true,
    children: item.children?.length
      ? toDisabledDictItemOptions(item.children, `${parentCode}:${item.code}`)
      : undefined,
  }))
}

function toDictOptionLabel(item: any) {
  const label = item.label || item.name || item.code
  return label && label !== item.code ? `${label} (${item.code})` : item.code
}

function handleWidgetUpdate(row: any, value: string) {
  row.form_widget = value
  if (value !== 'dict') {
    row.dict_code = null
  }
}

function handleDictCodeUpdate(row: any, value: string | number | Array<string | number> | null) {
  row.dict_code = Array.isArray(value) ? null : value === null ? null : String(value)
}
</script>

<template>
  <NFlex
    class="h-full min-h-0"
    vertical
  >
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
      title="代码生成"
      row-key="id"
      :scroll-x="1600"
      :columns="tableColumns"
      :data="state.rows"
      :loading="state.loading"
      :pagination="pagination"
      :checked-row-keys="state.checkedRowKeys"
      :on-update-checked-row-keys="handleCheckedRowKeys"
    >
      <template #toolbar>
        <NFlex>
          <NButton
            v-if="hasPermission('sys:codegen:create')"
            type="primary"
            text
            title="新增"
            @click="openCreateForm"
          >
            <template #icon>
              <NIcon><Icon icon="icon-park-outline:plus" /></NIcon>
            </template>
          </NButton>
          <NButton
            text
            title="刷新"
            :loading="state.loading"
            @click="fetchPage"
          >
            <template #icon>
              <NIcon><Icon icon="icon-park-outline:reload" /></NIcon>
            </template>
          </NButton>
          <NButton
            v-if="hasPermission('sys:codegen:delete')"
            type="error"
            text
            title="批量删除"
            :disabled="!hasCheckedRows"
            @click="confirmDelete(state.checkedRowKeys)"
          >
            <template #icon>
              <NIcon><Icon icon="icon-park-outline:delete" /></NIcon>
            </template>
          </NButton>
        </NFlex>
      </template>
    </ProDataTable>

    <NModal
      v-model:show="state.showForm"
      preset="card"
      draggable
      :mask-closable="false"
      :title="state.editingId ? '编辑生成方案' : '新增生成方案'"
      style="width: min(980px, 96vw)"
      :segmented="{ content: true, action: true }"
    >
      <NSpin :show="state.formLoading">
        <NScrollbar class="max-h-[min(680px,calc(100vh-300px))] pr-16px">
          <NForm
            ref="formRef"
            :model="state.form"
            :rules="formRules"
            label-placement="left"
            label-width="116"
          >
            <NGrid
              :cols="2"
              :x-gap="16"
            >
              <NGi>
                <NFormItem
                  label="方案名称"
                  path="name"
                >
                  <NInput v-model:value="state.form.name" />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="作者"
                  path="author"
                >
                  <NInput v-model:value="state.form.author" />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="生成类型"
                  path="gen_type"
                >
                  <NSelect
                    v-model:value="state.form.gen_type"
                    :options="genTypeOptions"
                  />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="主表"
                  path="main_table"
                >
                  <NSelect
                    v-model:value="state.form.main_table"
                    filterable
                    :options="state.tableOptions"
                    @update:value="handleMainTableUpdate"
                  />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="主键"
                  path="main_pk"
                >
                  <NSelect
                    v-model:value="state.form.main_pk"
                    filterable
                    :options="state.mainColumnOptions"
                  />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="主实体类"
                  path="main_entity_name"
                >
                  <NInput v-model:value="state.form.main_entity_name" />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="业务名称"
                  path="main_business_name"
                >
                  <NInput v-model:value="state.form.main_business_name" />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="后端模块路径"
                  path="main_module_path"
                >
                  <NInput v-model:value="state.form.main_module_path" />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="接口前缀"
                  path="api_prefix"
                >
                  <NInput v-model:value="state.form.api_prefix" />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="权限前缀"
                  path="permission_prefix"
                >
                  <NInput v-model:value="state.form.permission_prefix" />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="菜单名称"
                  path="menu_name"
                >
                  <NInput v-model:value="state.form.menu_name" />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="菜单路径"
                  path="menu_path"
                >
                  <NInput v-model:value="state.form.menu_path" />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem
                  label="组件路径"
                  path="component_path"
                >
                  <NInput v-model:value="state.form.component_path" />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem label="资源模块">
                  <NSelect
                    v-model:value="state.form.resource_module_id"
                    clearable
                    :options="state.moduleOptions"
                    @update:value="handleResourceModuleUpdate"
                  />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem label="父级菜单">
                  <NTreeSelect
                    v-model:value="state.form.parent_resource_id"
                    clearable
                    filterable
                    :options="state.parentOptions"
                    key-field="key"
                    label-field="label"
                    children-field="children"
                  />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem label="图标">
                  <IconSelect v-model:value="state.form.icon" />
                </NFormItem>
              </NGi>
              <NGi>
                <NFormItem label="排序">
                  <NInputNumber
                    v-model:value="state.form.sort"
                    class="w-full"
                    :min="0"
                  />
                </NFormItem>
              </NGi>
              <NGi v-if="needsTree">
                <NFormItem label="父级字段">
                  <NSelect
                    v-model:value="state.form.tree_parent_field"
                    filterable
                    :options="state.mainColumnOptions"
                  />
                </NFormItem>
              </NGi>
              <NGi v-if="needsTree">
                <NFormItem label="展示字段">
                  <NSelect
                    v-model:value="state.form.tree_label_field"
                    filterable
                    :options="state.mainColumnOptions"
                  />
                </NFormItem>
              </NGi>
              <NGi v-if="needsSub">
                <NFormItem label="子表">
                  <NSelect
                    v-model:value="state.form.sub_table"
                    filterable
                    :options="state.tableOptions"
                    @update:value="handleSubTableUpdate"
                  />
                </NFormItem>
              </NGi>
              <NGi v-if="needsSub">
                <NFormItem label="子表主键">
                  <NSelect
                    v-model:value="state.form.sub_pk"
                    filterable
                    :options="state.subColumnOptions"
                  />
                </NFormItem>
              </NGi>
              <NGi v-if="needsSub">
                <NFormItem label="子表外键">
                  <NSelect
                    v-model:value="state.form.sub_foreign_key"
                    filterable
                    :options="state.subColumnOptions"
                  />
                </NFormItem>
              </NGi>
              <NGi v-if="needsSub">
                <NFormItem label="子实体类">
                  <NInput v-model:value="state.form.sub_entity_name" />
                </NFormItem>
              </NGi>
              <NGi v-if="needsSub">
                <NFormItem label="子业务名称">
                  <NInput v-model:value="state.form.sub_business_name" />
                </NFormItem>
              </NGi>
            </NGrid>
            <NFormItem label="描述">
              <NInput
                v-model:value="state.form.description"
                type="textarea"
              />
            </NFormItem>
          </NForm>
        </NScrollbar>
      </NSpin>
      <template #action>
        <NSpace justify="end">
          <NButton @click="state.showForm = false">
            取消
          </NButton>
          <NButton
            type="primary"
            :loading="state.submitLoading"
            @click="submitForm"
          >
            确认
          </NButton>
        </NSpace>
      </template>
    </NModal>

    <NDrawer
      v-model:show="state.showFields"
      width="min(1280px, 96vw)"
    >
      <NDrawerContent title="字段配置">
        <NDataTable
          :columns="fieldColumns"
          :data="state.fieldRows"
          :loading="state.fieldLoading"
          :scroll-x="1700"
          :pagination="false"
        />
        <template #footer>
          <NSpace justify="end">
            <NButton @click="state.showFields = false">
              取消
            </NButton>
            <NButton
              type="primary"
              :loading="state.fieldSaving"
              @click="saveFields"
            >
              保存
            </NButton>
          </NSpace>
        </template>
      </NDrawerContent>
    </NDrawer>

    <NDrawer
      v-model:show="state.showPreview"
      width="min(1320px, 96vw)"
    >
      <NDrawerContent
        title="代码预览"
        closable
        body-content-class="codegen-preview-drawer"
      >
        <NSpin :show="state.previewLoading">
          <div class="codegen-preview">
            <NTabs
              v-if="state.previewFiles.length"
              v-model:value="state.previewPath"
              type="card"
              size="small"
              animated
            >
              <NTabPane
                v-for="file in state.previewFiles"
                :key="file.path"
                :name="file.path"
                :tab="previewTabLabel(file.path)"
              >
                <div
                  v-if="state.previewPath === file.path"
                  class="codegen-preview__file"
                >
                  <div
                    class="codegen-preview__path"
                    :title="file.path"
                  >
                    {{ file.path }}
                  </div>
                  <MonacoPreview
                    :value="file.content"
                    :language="file.language"
                    height="calc(100vh - 176px)"
                    :options="previewEditorOptions"
                  />
                </div>
              </NTabPane>
            </NTabs>
            <NEmpty
              v-else
              description="暂无预览文件"
            />
          </div>
        </NSpin>
      </NDrawerContent>
    </NDrawer>
  </NFlex>
</template>

<style scoped>
.codegen-preview {
  min-height: calc(100vh - 120px);
}

:deep(.codegen-preview-drawer) {
  padding: 16px;
}

.codegen-preview :deep(.n-tabs) {
  height: 100%;
}

.codegen-preview :deep(.n-tabs-nav) {
  overflow: hidden;
}

.codegen-preview :deep(.n-tab-pane) {
  padding: 0;
}

.codegen-preview__file {
  min-width: 0;
}

.codegen-preview__path {
  min-width: 0;
  padding: 8px 12px;
  overflow: hidden;
  border-right: 1px solid var(--n-border-color);
  border-left: 1px solid var(--n-border-color);
  background: rgba(0, 0, 0, 0.04);
  color: var(--n-text-color);
  font-family: var(--n-font-family-mono);
  font-size: 12px;
  line-height: 18px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
