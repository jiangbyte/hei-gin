<!--
  Author: Charlie

  反馈管理页。
-->
<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { sysFeedbackApi } from '@/api'
import {
  ACCOUNT_TYPE_TABS,
  DEFAULT_ACCOUNT_TYPE,
  accountTypeLabel,
  type AccountType,
} from '@/constants/account'
import {
  createTagColor,
  dictTypeColor,
  dictTypeData,
  dictList,
  formatDateTime,
  hasPermission,
  normalizeSearchValues,
  renderButtonIcon,
} from '@/utils'
import { NAvatar, NButton, NFlex, NIcon, NTag } from 'naive-ui'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { readPageMeta } from '@/utils/wire'

const router = useRouter()
const state = reactive({
  accountType: DEFAULT_ACCOUNT_TYPE as AccountType,
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
  { title: '标题', path: 'title', field: 'input' },
  {
    title: '反馈分类',
    path: 'category',
    field: 'select',
    fieldProps: { options: dictList('FEEDBACK_CATEGORY') },
  },
  {
    title: '状态',
    path: 'status',
    field: 'select',
    fieldProps: { options: dictList('FEEDBACK_STATUS') },
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
  { title: '标题', path: 'title', width: 220, ellipsis: { tooltip: true } },
  {
    title: '分类',
    path: 'category',
    width: 90,
    render: (row) => {
      const color = createTagColor(dictTypeColor('FEEDBACK_CATEGORY', row.category))
      const label = dictTypeData('FEEDBACK_CATEGORY', row.category)
      return (
        <NTag size="small" color={color} bordered={false}>
          {label || row.category}
        </NTag>
      )
    },
  },
  {
    title: '提交人',
    key: 'submitter',
    width: 160,
    render: (row) => {
      const name = row.submitter_nickname || row.submitter_account_id || '-'
      const initial = String(name)[0]?.toUpperCase() || '?'
      return (
        <NFlex align="center" size={8}>
          {row.submitter_avatar ? (
            <NAvatar src={row.submitter_avatar} size={24} round />
          ) : (
            <NAvatar size={24} round color="#d9d9d9">
              {initial}
            </NAvatar>
          )}
          <span>{name}</span>
        </NFlex>
      )
    },
  },
  { title: '联系方式', path: 'contact', width: 140, ellipsis: { tooltip: true } },
  {
    title: '状态',
    path: 'status',
    width: 80,
    render: (row) => {
      const color = createTagColor(dictTypeColor('FEEDBACK_STATUS', row.status))
      const label = dictTypeData('FEEDBACK_STATUS', row.status)
      return (
        <NTag size="small" color={color} bordered={false}>
          {label || row.status}
        </NTag>
      )
    },
  },
  {
    title: '提交时间',
    path: 'created_at',
    width: 170,
    render: (row) => formatDateTime(row.created_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 130,
    fixed: 'right',
    render: (row) => (
      <NFlex size={12}>
        {hasPermission('sys:feedback:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => openDetail(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('sys:feedback:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openEdit(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('sys:feedback:delete') ? (
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

function handleAccountTypeChange(value: string | number) {
  state.accountType = String(value) as AccountType
  state.page = 1
  state.checkedRowKeys = []
  void fetchPage()
}

async function fetchPage() {
  state.loading = true
  try {
    const response = await sysFeedbackApi.page({
      current: state.page,
      size: state.pageSize,
      submitter_account_type: state.accountType,
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

function openDetail(id: string) {
  router.push({ path: '/sys/feedback/detail', query: { id } })
}

function openEdit(id: string) {
  router.push({ path: '/sys/feedback/edit', query: { id } })
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
  await sysFeedbackApi.remove({ ids })
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
    <NTabs
      class="account-type-tabs"
      :value="state.accountType"
      type="line"
      animated
      @update:value="handleAccountTypeChange"
    >
      <NTabPane
        v-for="item in ACCOUNT_TYPE_TABS"
        :key="item.key"
        :name="item.key"
        :tab="item.label"
      />
    </NTabs>

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
      :title="`意见反馈 · ${accountTypeLabel(state.accountType)}`"
      row-key="id"
      :scroll-x="1100"
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
            text
            :loading="state.loading"
            @click="fetchPage"
          >
            <template #icon>
              <NIcon><Icon icon="icon-park-outline:refresh" /></NIcon>
            </template>
          </NButton>
          <NButton
            v-if="hasPermission('sys:feedback:delete')"
            type="error"
            text
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
  </NFlex>
</template>

<style scoped>
.account-type-tabs :deep(.n-tabs-pane-wrapper) {
  display: none;
}
</style>
