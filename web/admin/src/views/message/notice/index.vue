<!--
  Author: Charlie

  消息管理页。
-->
<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { msgNoticeApi } from '@/api'
import { accountTypeLabel } from '@/constants/account'
import {
  createTagColor,
  dictList,
  dictTypeColor,
  dictTypeData,
  formatDateTime,
  hasPermission,
  normalizeSearchValues,
  renderButtonIcon,
} from '@/utils'
import { NButton, NFlex, NIcon, NTag } from 'naive-ui'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { readPageMeta } from '@/utils/wire'

const router = useRouter()
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

const kindOptions = [
  { label: '通知', value: 'NOTIFICATION' },
  { label: '公告', value: 'ANNOUNCEMENT' },
]

const searchColumns = computed<ProSearchFormColumns<any>>(() => [
  { title: '标题', path: 'title', field: 'input' },
  {
    title: '类型',
    path: 'kind',
    field: 'select',
    fieldProps: { options: kindOptions },
  },
  {
    title: '状态',
    path: 'status',
    field: 'select',
    fieldProps: { options: dictList('PUBLISH_STATUS') },
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

const tableColumns = computed<ProDataTableColumns<any>>(() => [
  { type: 'selection', fixed: 'left' },
  { title: '标题', path: 'title', width: 200, ellipsis: { tooltip: true } },
  {
    title: '类型',
    path: 'kind',
    width: 80,
    render: (row) => (
      <NTag size="small" bordered={false} type={row.kind === 'ANNOUNCEMENT' ? 'warning' : 'info'}>
        {row.kind === 'ANNOUNCEMENT' ? '公告' : '通知'}
      </NTag>
    ),
  },
  {
    title: '等级',
    path: 'severity',
    width: 80,
    render: (row) => {
      const color = createTagColor(dictTypeColor('NOTIFICATION_SEVERITY', row.severity))
      const label = dictTypeData('NOTIFICATION_SEVERITY', row.severity)
      return (
        <NTag size="small" color={color} bordered={false}>
          {label || row.severity}
        </NTag>
      )
    },
  },
  {
    title: '目标范围',
    path: 'target_scope',
    width: 90,
    render: (row) => {
      const label = dictTypeData('TARGET_SCOPE', row.target_scope)
      return <span>{label || row.target_scope}</span>
    },
  },
  {
    title: '目标账户类型',
    path: 'target_account_types',
    width: 160,
    render: (row) => {
      const types = Array.isArray(row.target_account_types) ? row.target_account_types : []
      if (!types.length) return <span />
      return (
        <NFlex size={4} wrap>
          {types.map((t: string) => (
            <NTag key={t} size="small" bordered={false} type="info">
              {accountTypeLabel(t)}
            </NTag>
          ))}
        </NFlex>
      )
    },
  },
  {
    title: '置顶',
    path: 'is_pinned',
    width: 70,
    render: (row) =>
      row.kind === 'ANNOUNCEMENT' ? (
        <NTag size="small" bordered={false} type={row.is_pinned ? 'warning' : 'default'}>
          {row.is_pinned ? '是' : '否'}
        </NTag>
      ) : (
        <span>—</span>
      ),
  },
  {
    title: '状态',
    path: 'status',
    width: 80,
    render: (row) => {
      const color = createTagColor(dictTypeColor('PUBLISH_STATUS', row.status))
      const label = dictTypeData('PUBLISH_STATUS', row.status)
      return (
        <NTag size="small" color={color} bordered={false}>
          {label || row.status}
        </NTag>
      )
    },
  },
  {
    title: '查看次数',
    path: 'view_count',
    width: 80,
    align: 'right',
    render: (row) => (row.kind === 'ANNOUNCEMENT' ? row.view_count : '—'),
  },
  {
    title: '发布时间',
    path: 'publish_at',
    width: 170,
    render: (row) => formatDateTime(row.publish_at),
  },
  {
    title: '过期时间',
    path: 'expire_at',
    width: 170,
    render: (row) => (row.kind === 'ANNOUNCEMENT' ? formatDateTime(row.expire_at) : '—'),
  },
  {
    title: '更新时间',
    path: 'updated_at',
    width: 170,
    render: (row) => formatDateTime(row.updated_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 130,
    fixed: 'right',
    render: (row) => (
      <NFlex size={12}>
        {hasPermission('message:notice:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => openDetailPage(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('message:notice:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openEditPage(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('message:notice:delete') ? (
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
    const response = await msgNoticeApi.page({
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
    state.checkedRowKeys = state.checkedRowKeys.filter((key) =>
      state.rows.some((item) => item.id === key),
    )
  } finally {
    state.loading = false
  }
}

function openDetailPage(id: string) {
  router.push({ path: '/message/notice/detail', query: { id } })
}

function openCreatePage() {
  router.push('/message/notice/create')
}

function openEditPage(id: string) {
  router.push({ path: '/message/notice/edit', query: { id } })
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
    content: ids.length > 1 ? `删除 ${ids.length} 条记录?` : '删除该记录?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteRows(ids),
  })
}

async function deleteRows(ids: string[]) {
  await msgNoticeApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))
  window.$message.success('删除成功')
  await fetchPage()
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
      />
    </ProCard>

    <ProDataTable
      class="min-h-0 flex-1"
      remote
      title="消息管理"
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
          <NButton
            v-if="hasPermission('message:notice:create')"
            type="primary"
            text
            @click="openCreatePage"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:plus" />
              </NIcon>
            </template>
          </NButton>
          <NButton
            text
            :loading="state.loading"
            @click="fetchPage"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:refresh" />
              </NIcon>
            </template>
          </NButton>
          <NButton
            v-if="hasPermission('message:notice:delete')"
            type="error"
            text
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
