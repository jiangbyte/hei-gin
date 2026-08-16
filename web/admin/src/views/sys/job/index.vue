<!-- Author: Charlie -->

<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { NButton, NFlex, NIcon, NTag } from 'naive-ui'
import { jobApi } from '@/api'
import { formatDateTime, hasPermission, normalizeSearchValues, renderButtonIcon } from '@/utils'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { readPageMeta, wireBool } from '@/utils/wire'

const router = useRouter()
const state = reactive({
  jobs: [] as any[],
  total: 0,
  loading: false,
  searchValues: {} as any,
  checkedRowKeys: [] as string[],
  page: 1,
  pageSize: 20,
})

const EXECUTE_TYPE_OPTIONS = [
  { label: 'CRON 表达式', value: 'CRON' },
  { label: '固定间隔', value: 'FIXED' },
]
const ENABLED_OPTIONS = [
  { label: '启用', value: 'true' },
  { label: '停用', value: 'false' },
]

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
    title: '任务名称',
    path: 'job_name',
    field: 'input',
    fieldProps: { placeholder: '请输入任务名称' },
  },
  {
    title: '触发类型',
    path: 'execute_type',
    field: 'select',
    fieldProps: {
      options: EXECUTE_TYPE_OPTIONS,
    },
  },
  {
    title: '状态',
    path: 'enabled',
    field: 'select',
    fieldProps: {
      options: ENABLED_OPTIONS,
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

function renderEnabled(row: any) {
  return wireBool(row.enabled) ? (
    <NTag size="small" type="success" bordered={false}>
      启用
    </NTag>
  ) : (
    <NTag size="small" type="error" bordered={false}>
      停用
    </NTag>
  )
}

const tableColumns = computed<ProDataTableColumns<any>>(() => [
  {
    type: 'selection',
    fixed: 'left',
  },
  {
    title: '任务名称',
    path: 'job_name',
    width: 180,
    ellipsis: { tooltip: true },
  },
  {
    title: '执行类',
    path: 'execute_class',
    width: 280,
    ellipsis: { tooltip: true },
  },
  {
    title: '触发类型',
    path: 'execute_type',
    width: 110,
    render: (row) => (
      <NTag size="small" type={row.execute_type === 'CRON' ? 'info' : 'warning'} bordered={false}>
        {row.execute_type === 'CRON' ? 'CRON 表达式' : '固定间隔'}
      </NTag>
    ),
  },
  {
    title: '触发配置',
    path: 'trigger_config',
    width: 140,
    ellipsis: { tooltip: true },
  },
  {
    title: '下次执行',
    path: 'next_run_time',
    width: 190,
    ellipsis: { tooltip: true },
    render: (row) => formatDateTime(row.next_run_time),
  },
  {
    title: '上次执行',
    path: 'last_run_time',
    width: 190,
    ellipsis: { tooltip: true },
    render: (row) => formatDateTime(row.last_run_time),
  },
  {
    title: '上次结果',
    path: 'last_execute_result',
    width: 200,
    ellipsis: { tooltip: true },
  },
  {
    title: '状态',
    path: 'enabled',
    width: 90,
    render: (row) => renderEnabled(row),
  },
  {
    title: '排序',
    path: 'sort',
    width: 70,
  },
  {
    title: '更新时间',
    path: 'updated_at',
    width: 190,
    ellipsis: { tooltip: true },
    render: (row) => formatDateTime(row.updated_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    fixed: 'right',
    render: (row) => (
      <NFlex size={12}>
        {hasPermission('sys:job:run') ? (
          <NButton type="warning" size="small" text={true} onClick={() => confirmRun(row)}>
            {renderButtonIcon('icon-park-outline:play')}
          </NButton>
        ) : null}
        {hasPermission('sys:job:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => openDetail(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('sys:joblog:page') ? (
          <NButton type="info" size="small" text={true} onClick={() => openLog(row.id)}>
            {renderButtonIcon('icon-park-outline:log')}
          </NButton>
        ) : null}
        {hasPermission('sys:job:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openEdit(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('sys:job:update') ? (
          <NButton size="small" text={true} onClick={() => toggleEnabled(row)}>
            {renderButtonIcon(wireBool(row.enabled) ? 'icon-park-outline:pause' : 'icon-park-outline:play')}
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
  fetchPage()
})

async function fetchPage() {
  state.loading = true
  try {
    const response = await jobApi.page({
      current: state.page,
      size: state.pageSize,
      ...state.searchValues,
    })
    const data = response.data ?? {}
    state.jobs = data.records ?? []
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
  } finally {
    state.loading = false
  }
}

function openDetail(id: string) {
  router.push({ path: '/sys/job/detail', query: { id } })
}

function openLog(id: string) {
  router.push({ path: '/sys/job/log', query: { id } })
}

function openCreate() {
  router.push('/sys/job/create')
}

function openEdit(id: string) {
  router.push({ path: '/sys/job/edit', query: { id } })
}

function handleCheckedRowKeys(keys: Array<string | number>) {
  state.checkedRowKeys = keys.map(String)
}

function confirmRun(row: any) {
  window.$dialog.info({
    title: '立即执行',
    draggable: true,
    maskClosable: false,
    content: `确认立即执行任务「${row.job_name}」？`,
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: async () => {
      await jobApi.run({ id: row.id })
      window.$message.success('已触发执行，结果见执行日志')
    },
  })
}

function toggleEnabled(row: any) {
  const target = !wireBool(row.enabled)
  window.$dialog.info({
    title: target ? '启用任务' : '停用任务',
    draggable: true,
    maskClosable: false,
    content: target ? `确认启用任务「${row.job_name}」？` : `确认停用任务「${row.job_name}」？`,
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: async () => {
      await jobApi.enabled({ id: row.id, enabled: target })
      window.$message.success(target ? '已启用' : '已停用')
      await fetchPage()
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
    onPositiveClick: () => deleteData(ids),
  })
}

async function deleteData(ids: string[]) {
  await jobApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))

  window.$message.success('删除成功')
  await fetchPage()
  if (!state.jobs.length && state.total > 0 && state.page > 1) {
    state.page -= 1
    await fetchPage()
  }
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
        :collapse-button-props="{
          content: searchForm.collapsed.value ? '展开' : '收起',
        }"
      />
    </ProCard>

    <ProDataTable
      class="min-h-0 flex-1"
      remote
      :title="'任务管理'"
      row-key="id"
      :scroll-x="1720"
      :columns="tableColumns"
      :data="state.jobs"
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
            :aria-label="'新增'"
            @click="openCreate"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:plus" />
              </NIcon>
            </template>
          </NButton>
          <NButton
            text
            :title="'刷新'"
            :aria-label="'刷新'"
            :loading="state.loading"
            @click="fetchPage"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:reload" />
              </NIcon>
            </template>
          </NButton>
          <NButton
            v-if="hasPermission('sys:job:delete')"
            type="error"
            text
            :title="'批量删除'"
            :aria-label="'批量删除'"
            :disabled="!hasCheckedRows"
            @click="confirmDelete(state.checkedRowKeys)"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:delete" />
              </NIcon>
            </template>
          </NButton>
        </NFlex>
      </template>
    </ProDataTable>
  </NFlex>
</template>
