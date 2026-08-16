<!-- Author: Charlie -->

<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import {
  NButton,
  NFlex,
  NIcon,
  NInput,
  NModal,
  NSelect,
  NSwitch,
  NTag,
  NForm,
  NFormItem,
} from 'naive-ui'
import { jobApi } from '@/api'
import {
  createTagColor,
  formatDateTime,
  hasPermission,
  normalizeSearchValues,
  renderButtonIcon,
} from '@/utils'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'
import { readPageMeta } from '@/utils/wire'

const state = reactive({
  rows: [] as any[],
  total: 0,
  loading: false,
  searchValues: {} as any,
  checkedRowKeys: [] as string[],
  page: 1,
  pageSize: 20,
  handlers: [] as any[],
})

// ---------- 表单弹窗 ----------
const formVisible = ref(false)
const formSaving = ref(false)
const editingId = ref<string | null>(null)
const form = reactive({
  handler_key: '',
  name: '',
  cron_expr: '',
  params: '',
  status: 'ENABLED',
  description: '',
})
const formRules = {
  handler_key: { required: true, message: '请选择处理器', trigger: ['blur', 'change'] },
  name: { required: true, message: '请输入名称', trigger: ['blur', 'input'] },
  cron_expr: {
    required: true,
    message: '请输入 cron 表达式（6 段含秒）',
    trigger: ['blur', 'input'],
  },
}

// ---------- 日志弹窗 ----------
const logVisible = ref(false)
const logLoading = ref(false)
const logRows = ref([] as any[])
const logTotal = ref(0)
const logPage = ref(1)
const logPageSize = ref(10)
const logJobId = ref('')
const logJobName = ref('')

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
  {
    title: '关键字',
    path: 'keyword',
    field: 'input',
    fieldProps: { placeholder: '名称 / 处理器', clearable: true },
  },
  {
    title: '状态',
    path: 'status',
    field: 'select',
    fieldProps: {
      options: [
        { label: '启用', value: 'ENABLED' },
        { label: '停用', value: 'DISABLED' },
      ],
    },
  },
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

const logPagination = computed<PaginationProps>(() => ({
  page: logPage.value,
  pageSize: logPageSize.value,
  itemCount: logTotal.value,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  prefix: ({ itemCount }) => `${itemCount} 条`,
  onUpdatePage: (value) => {
    logPage.value = value
    fetchLogs()
  },
  onUpdatePageSize: (value) => {
    logPageSize.value = value
    logPage.value = 1
    fetchLogs()
  },
}))

const logColumns: ProDataTableColumns<any> = [
  {
    title: '开始时间',
    path: 'started_at',
    width: 190,
    render: (row) => formatDateTime(row.started_at),
  },
  { title: '状态', key: 'status', width: 100, render: (row) => renderLogStatus(row.status) },
  { title: '耗时(ms)', path: 'duration_ms', width: 120 },
  { title: '消息', path: 'message', minWidth: 240, ellipsis: { tooltip: true } },
  {
    title: '完成时间',
    path: 'finished_at',
    width: 190,
    render: (row) => formatDateTime(row.finished_at),
  },
]

const tableColumns = computed<ProDataTableColumns<any>>(() => [
  {
    type: 'selection',
    fixed: 'left',
  },
  { title: '名称', path: 'name', width: 180, ellipsis: { tooltip: true } },
  {
    title: '处理器',
    path: 'handler_key',
    width: 200,
    render: (row) => (
      <NTag size="small" bordered={false} type="info">
        {row.handler_key}
      </NTag>
    ),
  },
  {
    title: 'cron 表达式',
    path: 'cron_expr',
    width: 140,
    render: (row) => <code>{row.cron_expr}</code>,
  },
  { title: '参数', path: 'params', width: 140, ellipsis: { tooltip: true } },
  {
    title: '状态',
    path: 'status',
    width: 100,
    render: (row) => (
      <NTag
        size="small"
        color={createTagColor(row.status === 'ENABLED' ? '#18a058' : '#d03050')}
        bordered={false}
      >
        {row.status === 'ENABLED' ? '启用' : '停用'}
      </NTag>
    ),
  },
  {
    title: '上次执行',
    path: 'last_run_at',
    width: 190,
    render: (row) => formatDateTime(row.last_run_at),
  },
  {
    title: '下次执行',
    path: 'next_run_at',
    width: 190,
    render: (row) => formatDateTime(row.next_run_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 240,
    fixed: 'right',
    render: (row) => (
      <NFlex size={8}>
        {hasPermission('sys:job:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openEdit(row)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('sys:job:update') ? (
          <NButton size="small" text={true} onClick={() => toggleStatus(row)}>
            {renderButtonIcon(
              row.status === 'ENABLED' ? 'icon-park-outline:pause' : 'icon-park-outline:play',
            )}
          </NButton>
        ) : null}
        {hasPermission('sys:job:update') ? (
          <NButton type="warning" size="small" text={true} onClick={() => confirmTrigger(row)}>
            {renderButtonIcon('icon-park-outline:flash-payment')}
          </NButton>
        ) : null}
        {hasPermission('sys:job:page') ? (
          <NButton type="info" size="small" text={true} onClick={() => openLogs(row)}>
            {renderButtonIcon('icon-park-outline:file-tips')}
          </NButton>
        ) : null}
        {hasPermission('sys:job:delete') ? (
          <NButton type="error" size="small" text={true} onClick={() => confirmDelete(row.id)}>
            {renderButtonIcon('icon-park-outline:delete')}
          </NButton>
        ) : null}
      </NFlex>
    ),
  },
])

const hasCheckedRows = computed(() => state.checkedRowKeys.length > 0)

onMounted(() => {
  fetchHandlers()
  fetchPage()
})

async function fetchHandlers() {
  try {
    const response = await jobApi.handlers()
    state.handlers = response.data ?? []
  } catch {
    state.handlers = []
  }
}

async function fetchPage() {
  state.loading = true
  try {
    const response = await jobApi.page({
      current: state.page,
      size: state.pageSize,
      ...state.searchValues,
    })
    const data = response.data ?? {}
    state.rows = data.records ?? []
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
  } finally {
    state.loading = false
  }
}

async function fetchLogs() {
  logLoading.value = true
  try {
    const response = await jobApi.logs({
      job_id: logJobId.value,
      current: logPage.value,
      size: logPageSize.value,
    })
    const data = response.data ?? {}
    logRows.value = data.records ?? []
    const pageMeta = readPageMeta(data, { current: logPage.value, size: logPageSize.value })
    logTotal.value = pageMeta.total
    logPage.value = pageMeta.current
    logPageSize.value = pageMeta.size
  } finally {
    logLoading.value = false
  }
}

function openCreate() {
  editingId.value = null
  Object.assign(form, {
    handler_key: '',
    name: '',
    cron_expr: '',
    params: '',
    status: 'ENABLED',
    description: '',
  })
  formVisible.value = true
}

async function openEdit(row: any) {
  editingId.value = row.id
  Object.assign(form, {
    handler_key: row.handler_key ?? '',
    name: row.name ?? '',
    cron_expr: row.cron_expr ?? '',
    params: row.params ?? '',
    status: row.status ?? 'ENABLED',
    description: row.description ?? '',
  })
  formVisible.value = true
}

async function submitForm() {
  formSaving.value = true
  try {
    const payload = {
      handler_key: form.handler_key,
      name: form.name,
      cron_expr: form.cron_expr,
      params: form.params,
      status: form.status,
      description: form.description || null,
    }
    if (editingId.value) {
      await jobApi.update({ id: editingId.value, ...payload })
      window.$message.success('更新成功')
    } else {
      await jobApi.create(payload)
      window.$message.success('创建成功')
    }
    formVisible.value = false
    await fetchPage()
  } finally {
    formSaving.value = false
  }
}

async function toggleStatus(row: any) {
  const next = row.status === 'ENABLED' ? 'DISABLED' : 'ENABLED'
  await jobApi.setStatus({ id: row.id, status: next })
  window.$message.success(next === 'ENABLED' ? '已启用' : '已停用')
  await fetchPage()
}

function confirmTrigger(row: any) {
  window.$dialog.warning({
    title: '立即触发',
    draggable: true,
    maskClosable: false,
    content: `立即执行「${row.name}」？`,
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: async () => {
      await jobApi.trigger({ id: row.id })
      window.$message.success('已触发')
    },
  })
}

function confirmDelete(value: string | string[]) {
  const ids = Array.isArray(value) ? value : [value]
  if (!ids.length) {
    return
  }
  const isBatch = ids.length > 1
  window.$dialog.warning({
    title: isBatch ? '批量删除' : '删除',
    draggable: true,
    maskClosable: false,
    content: isBatch ? `删除 ${ids.length} 个任务?` : '删除该任务?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: async () => {
      await jobApi.remove({ ids })
      state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))
      window.$message.success('删除成功')
      await fetchPage()
      if (!state.rows.length && state.total > 0 && state.page > 1) {
        state.page -= 1
        await fetchPage()
      }
    },
  })
}

function openLogs(row: any) {
  logJobId.value = row.id
  logJobName.value = row.name
  logRows.value = []
  logTotal.value = 0
  logPage.value = 1
  logVisible.value = true
  fetchLogs()
}

function handleCheckedRowKeys(keys: Array<string | number>) {
  state.checkedRowKeys = keys.map(String)
}

function renderLogStatus(status: string) {
  const map: Record<string, { color: string; text: string }> = {
    SUCCESS: { color: '#18a058', text: '成功' },
    FAIL: { color: '#d03050', text: '失败' },
    CANCEL: { color: '#f0a020', text: '取消' },
    TIMEOUT: { color: '#f0a020', text: '超时' },
  }
  const item = map[status] ?? { color: '#909399', text: status }
  return (
    <NTag size="small" color={createTagColor(item.color)} bordered={false}>
      {item.text}
    </NTag>
  )
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
      :title="'任务管理'"
      row-key="id"
      :scroll-x="1680"
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
            v-if="hasPermission('sys:job:create')"
            type="primary"
            text
            :title="'新增'"
            @click="openCreate"
          >
            <template #icon
              ><NIcon><Icon icon="icon-park-outline:plus" /></NIcon
            ></template>
          </NButton>
          <NButton text :title="'刷新'" :loading="state.loading" @click="fetchPage">
            <template #icon
              ><NIcon><Icon icon="icon-park-outline:reload" /></NIcon
            ></template>
          </NButton>
          <NButton
            v-if="hasPermission('sys:job:delete')"
            type="error"
            text
            :title="'批量删除'"
            :disabled="!hasCheckedRows"
            @click="confirmDelete(state.checkedRowKeys)"
          >
            <template #icon
              ><NIcon><Icon icon="icon-park-outline:delete" /></NIcon
            ></template>
          </NButton>
        </NFlex>
      </template>
    </ProDataTable>

    <!-- 新建 / 编辑弹窗 -->
    <NModal
      v-model:show="formVisible"
      preset="card"
      :title="editingId ? '编辑任务' : '新建任务'"
      style="width: 560px"
    >
      <NForm
        :model="form"
        :rules="formRules"
        label-placement="left"
        label-width="110"
        require-mark-placement="right"
      >
        <NFormItem label="处理器" path="handler_key">
          <NSelect
            v-model:value="form.handler_key"
            :options="state.handlers.map((h: any) => ({ label: `${h.name} (${h.key})`, value: h.key }))"
            filterable
            placeholder="选择已注册的处理器"
          />
        </NFormItem>
        <NFormItem label="名称" path="name">
          <NInput v-model:value="form.name" placeholder="任务显示名称" />
        </NFormItem>
        <NFormItem label="cron 表达式" path="cron_expr">
          <NInput
            v-model:value="form.cron_expr"
            placeholder="如 0 0 2 * * *（秒 分 时 日 月 周）"
          />
        </NFormItem>
        <NFormItem label="参数">
          <NInput v-model:value="form.params" placeholder="传给处理器的参数（可选）" />
        </NFormItem>
        <NFormItem label="状态">
          <NSwitch
            v-model:value="form.status"
            :checked-value="'ENABLED'"
            :unchecked-value="'DISABLED'"
          />
        </NFormItem>
        <NFormItem label="描述">
          <NInput
            v-model:value="form.description"
            type="textarea"
            :rows="2"
            placeholder="任务说明（可选）"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NFlex justify="end">
          <NButton @click="formVisible = false">取消</NButton>
          <NButton type="primary" :loading="formSaving" @click="submitForm">保存</NButton>
        </NFlex>
      </template>
    </NModal>

    <!-- 执行日志弹窗：固定高度容器 + flex-height，表格始终可见、超出时 body 内部滚动（对齐其他模块弹窗定高策略） -->
    <NModal
      v-model:show="logVisible"
      preset="card"
      :title="`执行日志 - ${logJobName}`"
      style="width: 960px"
    >
      <div class="h-[min(560px,calc(100vh-300px))]">
        <ProDataTable
          remote
          flex-height
          row-key="id"
          :columns="logColumns"
          :data="logRows"
          :loading="logLoading"
          :pagination="logPagination"
          :scroll-x="980"
          class="h-full"
        />
      </div>
    </NModal>
  </NFlex>
</template>
