<!-- Author: Charlie -->

<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { resourceModuleApi } from '@/api'
import NovaIcon from '@/components/common/NovaIcon.vue'
import { ACCOUNT_TYPE_OPTIONS, accountTypeLabel } from '@/constants/account'
import {
  createTagColor,
  formatDateTime,
  hasPermission,
  normalizeSearchValues,
  renderButtonIcon,
} from '@/utils'
import { NButton, NFlex, NIcon, NTag } from 'naive-ui'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'
import { dictList, dictTypeColor, dictTypeData } from '@/utils/dict'
import ModalDetail from './components/ModalDetail.vue'
import ModalForm from './components/ModalForm.vue'
import { readPageMeta } from '@/utils/wire'

const formModalRef = ref<any>(null)
const detailModalRef = ref<any>(null)
const state = reactive({
  modules: [] as any[],
  total: 0,
  loading: false,
  searchValues: {} as any,
  checkedRowKeys: [] as string[],
  page: 1,
  pageSize: 20,
})

const searchForm = createProSearchForm<any>({
  defaultCollapsed: true,
  onSubmit(values) {
    state.searchValues = normalizeSearchValues(values, {
      name: (value) => String(value).trim(),
      code: (value) => String(value).trim(),
      client: (value) => String(value).trim(),
    })
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
    title: '模块名称',
    path: 'name',
    field: 'input',
  },
  {
    title: '模块编码',
    path: 'code',
    field: 'input',
  },
  {
    title: '客户端',
    path: 'client',
    field: 'select',
    fieldProps: {
      options: ACCOUNT_TYPE_OPTIONS,
    },
  },
  {
    title: '状态',
    path: 'status',
    field: 'select',
    fieldProps: {
      options: dictList('COMMON_STATUS'),
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

const tableColumns = computed<ProDataTableColumns<any>>(() => [
  {
    type: 'selection',
    fixed: 'left',
  },
  {
    title: '模块名称',
    path: 'name',
    width: 160,
    render: (row) => (
      <span
        style={{
          display: 'inline-flex',
          alignItems: 'center',
          gap: '8px',
          maxWidth: '100%',
          minWidth: 0,
          verticalAlign: 'middle',
        }}
      >
        {row.icon ? (
          <span
            title={row.icon}
            style={{
              display: 'inline-flex',
              flex: '0 0 auto',
              alignItems: 'center',
              justifyContent: 'center',
              width: '20px',
              height: '20px',
              lineHeight: 1,
            }}
          >
            <NovaIcon icon={row.icon} size={18} />
          </span>
        ) : null}
        <span
          style={{
            minWidth: 0,
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
            lineHeight: '20px',
          }}
        >
          {row.name}
        </span>
      </span>
    ),
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '模块编码',
    path: 'code',
    width: 160,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '客户端',
    path: 'client',
    width: 120,
    render: (row) => accountTypeLabel(row.client) || row.client,
  },
  {
    title: '颜色',
    path: 'color',
    width: 110,
    render: (row) =>
      row.color ? (
        <NTag size="small" color={createTagColor(row.color)} bordered={false}>
          {row.color}
        </NTag>
      ) : (
        '-'
      ),
  },
  {
    title: '排序',
    path: 'sort',
    width: 90,
  },
  {
    title: '状态',
    path: 'status',
    width: 110,
    render: (row) => (
      <NTag
        size="small"
        color={createTagColor(dictTypeColor('COMMON_STATUS', row.status))}
        bordered={false}
      >
        {dictTypeData('COMMON_STATUS', row.status) || row.status}
      </NTag>
    ),
  },
  {
    title: '更新时间',
    path: 'updated_at',
    width: 190,
    ellipsis: {
      tooltip: true,
    },
    render: (row) => formatDateTime(row.updated_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    fixed: 'right',
    render: (row) => (
      <NFlex size={12}>
        {hasPermission('iam:resourcemodule:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => openDetailModal(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('iam:resourcemodule:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openEditModal(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('iam:resourcemodule:delete') ? (
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
    const response = await resourceModuleApi.page({
      current: state.page,
      size: state.pageSize,
      ...state.searchValues,
    })
    const data = response.data ?? {}
    state.modules = data.records ?? []
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
    state.checkedRowKeys = state.checkedRowKeys.filter((key) =>
      state.modules.some((item) => item.id === key),
    )
  } finally {
    state.loading = false
  }
}

function openDetailModal(id: string) {
  detailModalRef.value?.openModal(id)
}

function openCreateModal() {
  formModalRef.value?.openModal()
}

function openEditModal(id: string) {
  formModalRef.value?.openModal(id)
}

function handleCheckedRowKeys(keys: Array<string | number>) {
  state.checkedRowKeys = keys.map(String)
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
    content: isBatch ? `删除 ${ids.length} 个资源模块?` : '删除该资源模块?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteData(ids),
  })
}

async function deleteData(ids: string[]) {
  await resourceModuleApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))

  window.$message.success('删除成功')
  await fetchPage()
  if (!state.modules.length && state.total > 0 && state.page > 1) {
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
      :title="'资源模块管理'"
      row-key="id"
      :scroll-x="1420"
      :columns="tableColumns"
      :data="state.modules"
      :loading="state.loading"
      :pagination="pagination"
      :checked-row-keys="state.checkedRowKeys"
      :on-update-checked-row-keys="handleCheckedRowKeys"
    >
      <template #toolbar>
        <NFlex>
          <NButton
            v-if="hasPermission('iam:resourcemodule:create')"
            type="primary"
            text
            :title="'新增'"
            :aria-label="'新增'"
            @click="openCreateModal"
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
            v-if="hasPermission('iam:resourcemodule:delete')"
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

    <ModalForm
      ref="formModalRef"
      @saved="fetchPage"
    />
    <ModalDetail ref="detailModalRef" />
  </NFlex>
</template>

<style scoped></style>
