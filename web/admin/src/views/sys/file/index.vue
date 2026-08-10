<!-- Author: Charlie -->

<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { NButton, NFlex, NIcon, NImage, NTag } from 'naive-ui'
import { fileApi } from '@/api'
import {
  createTagColor,
  formatDateTime,
  formatFileSize,
  hasPermission,
  isImageFile,
  normalizeSearchValues,
  renderButtonIcon,
} from '@/utils'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'
import ModalDetail from './components/ModalDetail.vue'
import ModalForm from './components/ModalForm.vue'
import ModalUpload from './components/ModalUpload.vue'
import { STORAGE_PROVIDER_OPTIONS, storageProviderColor, storageProviderLabel } from './constants'
import { readPageMeta } from '@/utils/wire'

const formModalRef = ref<any>(null)
const detailModalRef = ref<any>(null)
const uploadModalRef = ref<any>(null)
const state = reactive({
  files: [] as any[],
  total: 0,
  loading: false,
  downloadingIds: [] as string[],
  searchValues: {} as any,
  checkedRowKeys: [] as string[],
  page: 1,
  pageSize: 20,
})

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
    title: '文件名',
    path: 'original_name',
    field: 'input',
  },
  {
    title: '对象路径',
    path: 'object_name',
    field: 'input',
  },
  {
    title: '存储提供商',
    path: 'storage_provider',
    field: 'select',
    fieldProps: {
      options: [...STORAGE_PROVIDER_OPTIONS],
    },
  },
  {
    title: '内容类型',
    path: 'content_type',
    field: 'input',
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
    title: '预览',
    key: 'preview',
    width: 120,
    render: (row) => renderPreview(row),
  },
  {
    title: '文件名',
    path: 'original_name',
    width: 220,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '对象路径',
    path: 'object_name',
    width: 320,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '存储提供商',
    path: 'storage_provider',
    width: 130,
    render: (row) => (
      <NTag
        size="small"
        color={createTagColor(storageProviderColor(row.storage_provider))}
        bordered={false}
      >
        {storageProviderLabel(row.storage_provider)}
      </NTag>
    ),
  },
  {
    title: '存储桶',
    path: 'bucket',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '内容类型',
    path: 'content_type',
    width: 180,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '文件大小',
    path: 'size',
    width: 120,
    render: (row) => formatFileSize(row.size),
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
    width: 150,
    fixed: 'right',
    render: (row) => (
      <NFlex size={12}>
        {hasPermission('sys:file:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => openDetailModal(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('sys:file:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openEditModal(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('sys:file:url') ? (
          <NButton
            type="primary"
            size="small"
            text={true}
            loading={isDownloading(row.id)}
            onClick={() => openFile(row)}
          >
            {renderButtonIcon('icon-park-outline:link')}
          </NButton>
        ) : null}
        {hasPermission('sys:file:delete') ? (
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
    const response = await fileApi.page({
      current: state.page,
      size: state.pageSize,
      ...state.searchValues,
    })
    const data = response.data ?? {}
    state.files = data.records ?? []
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
  } finally {
    state.loading = false
  }
}

function renderPreview(row: any) {
  const src = row.url || undefined
  if (!src || !isImageFile(row)) {
    return (
      <NTag size="small" bordered={false}>
        {row.content_type || '-'}
      </NTag>
    )
  }
  return (
    <NImage src={src} alt={row.original_name || '预览'} width={72} height={48} objectFit="cover" />
  )
}

function openDetailModal(id: string) {
  detailModalRef.value?.openModal(id)
}

function openEditModal(id: string) {
  formModalRef.value?.openModal(id)
}

function openUploadModal() {
  uploadModalRef.value?.openModal()
}

function handleCheckedRowKeys(keys: Array<string | number>) {
  state.checkedRowKeys = keys.map(String)
}

function isDownloading(id?: string | number | null) {
  return state.downloadingIds.includes(String(id ?? ''))
}

async function openFile(row: any) {
  const id = String(row.id || '')
  if (!id || isDownloading(id)) {
    return
  }
  state.downloadingIds.push(id)
  try {
    await fileApi.downloadFile(row, row.original_name || row.object_name || 'download')
  } finally {
    state.downloadingIds = state.downloadingIds.filter((item) => item !== id)
  }
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
    content: isBatch ? `删除 ${ids.length} 个文件?` : '删除该文件?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteData(ids),
  })
}

async function deleteData(ids: string[]) {
  await fileApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))

  window.$message.success('删除成功')
  await fetchPage()
  if (!state.files.length && state.total > 0 && state.page > 1) {
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
      :title="'文件管理'"
      row-key="id"
      :scroll-x="1800"
      :columns="tableColumns"
      :data="state.files"
      :loading="state.loading"
      :pagination="pagination"
      :checked-row-keys="state.checkedRowKeys"
      :on-update-checked-row-keys="handleCheckedRowKeys"
    >
      <template #toolbar>
        <NFlex align="center">
          <NButton
            v-if="hasPermission('sys:file:upload')"
            text
            :title="'上传文件'"
            :aria-label="'上传文件'"
            @click="openUploadModal"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:upload" />
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
            v-if="hasPermission('sys:file:delete')"
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
    <ModalUpload
      ref="uploadModalRef"
      @saved="fetchPage"
    />
  </NFlex>
</template>
