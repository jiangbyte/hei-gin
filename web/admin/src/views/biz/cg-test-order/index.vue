<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-09 21:39:42
-->

<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { cgTestOrderApi } from '@/api'
import { readPageMeta } from '@/utils/wire'
import { createTagColor, dictTypeColor, dictTypeData, displayValue, formatDateTime, hasPermission, normalizeSearchValues, renderButtonIcon } from '@/utils'
import { NButton, NFlex, NIcon, NTag } from 'naive-ui'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'
import ChildModalDetail from './components/children/ChildModalDetail.vue'
import ChildModalForm from './components/children/ChildModalForm.vue'
import ModalDetail from './components/ModalDetail.vue'
import ModalForm from './components/ModalForm.vue'

const formModalRef = ref<any>(null)
const detailModalRef = ref<any>(null)
const childFormModalRef = ref<any>(null)
const childDetailModalRef = ref<any>(null)
const state = reactive({
  rows: [] as any[],
  total: 0,
  loading: false,
  searchValues: {} as any,
  checkedRowKeys: [] as string[],
  page: 1,
  pageSize: 20,
  childRows: [] as any[],
  childTotal: 0,
  childLoading: false,
  childDrawerVisible: false,
  childSearchValues: {} as any,
  childCheckedRowKeys: [] as string[],
  childPage: 1,
  childPageSize: 20,
  selectedMasterId: null as string | null,
})

const hasCheckedRows = computed(() => state.checkedRowKeys.length > 0)
const hasChildCheckedRows = computed(() => state.childCheckedRowKeys.length > 0)
const canCreateChild = computed(() => Boolean(state.selectedMasterId))

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
  { title: 'order_no', path: 'order_no', field: 'input' },
  { title: 'name', path: 'name', field: 'input' },
  { title: 'customer_name', path: 'customer_name', field: 'input' },
  { title: 'status', path: 'status', field: 'input' },
  { title: 'type', path: 'type', field: 'input' },
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

const childSearchForm = createProSearchForm<any>({
  defaultCollapsed: true,
  onSubmit(values) {
    state.childSearchValues = normalizeSearchValues(values)
    state.childPage = 1
    fetchChildPage()
  },
  onReset() {
    state.childSearchValues = {}
    state.childPage = 1
    fetchChildPage()
  },
})

const childSearchColumns = computed<ProSearchFormColumns<any>>(() => [
  { title: 'order_id', path: 'order_id', field: 'input' },
  { title: 'sku_code', path: 'sku_code', field: 'input' },
  { title: 'name', path: 'name', field: 'input' },
  { title: 'status', path: 'status', field: 'input' },
])

const childPagination = computed<PaginationProps>(() => ({
  page: state.childPage,
  pageSize: state.childPageSize,
  itemCount: state.childTotal,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: ({ itemCount }) => `${itemCount} 条`,
  onUpdatePage: (value) => {
    state.childPage = value
    fetchChildPage()
  },
  onUpdatePageSize: (value) => {
    state.childPageSize = value
    state.childPage = 1
    fetchChildPage()
  },
}))

const tableColumns = computed<ProDataTableColumns<any>>(() => [
  { type: 'selection', fixed: 'left' },
  { title: 'order_no', path: 'order_no', width: 150, ellipsis: { tooltip: true } },
  { title: 'name', path: 'name', width: 150, ellipsis: { tooltip: true } },
  { title: 'customer_name', path: 'customer_name', width: 150, ellipsis: { tooltip: true } },
  { title: 'customer_phone', path: 'customer_phone', width: 150, ellipsis: { tooltip: true } },
  {
    title: 'status',
    path: 'status',
    width: 150,
    ellipsis: { tooltip: true },
    render: row => (
      <NTag color={createTagColor(dictTypeColor('COMMON_STATUS', row.status))} bordered={false}>
        {dictTypeData('COMMON_STATUS', row.status) || displayValue(row.status)}
      </NTag>
    ),
  },
  { title: 'type', path: 'type', width: 150, ellipsis: { tooltip: true } },
  { title: 'ordered_at', path: 'ordered_at', width: 190, render: row => formatDateTime(row.ordered_at) },
  { title: 'paid_at', path: 'paid_at', width: 190, render: row => formatDateTime(row.paid_at) },
  { title: '更新时间', path: 'updated_at', width: 190, render: row => formatDateTime(row.updated_at) },
  {
    title: '操作',
    key: 'actions',
    width: 170,
    fixed: 'right',
    render: row => (
      <NFlex size={12}>
        {hasPermission('biz:cgtestorder:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => openDetailModal(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('biz:cgtestorder:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openEditModal(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        <NButton type="info" size="small" text={true} onClick={() => selectMaster(row.id)}>
          {renderButtonIcon('icon-park-outline:list-view')}
        </NButton>
        {hasPermission('biz:cgtestorder:delete') ? (
          <NButton type="error" size="small" text={true} onClick={() => confirmDelete(row.id)}>
            {renderButtonIcon('icon-park-outline:delete')}
          </NButton>
        ) : null}
      </NFlex>
    ),
  },
])

const childColumns = computed<ProDataTableColumns<any>>(() => [
  { type: 'selection', fixed: 'left' },
  { title: 'order_id', path: 'order_id', width: 150, ellipsis: { tooltip: true } },
  { title: 'sku_code', path: 'sku_code', width: 150, ellipsis: { tooltip: true } },
  { title: 'name', path: 'name', width: 150, ellipsis: { tooltip: true } },
  { title: 'category', path: 'category', width: 150, ellipsis: { tooltip: true } },
  {
    title: 'status',
    path: 'status',
    width: 150,
    ellipsis: { tooltip: true },
    render: row => (
      <NTag color={createTagColor(dictTypeColor('COMMON_STATUS', row.status))} bordered={false}>
        {dictTypeData('COMMON_STATUS', row.status) || displayValue(row.status)}
      </NTag>
    ),
  },
  { title: 'quantity', path: 'quantity', width: 150, ellipsis: { tooltip: true } },
  { title: 'unit_price', path: 'unit_price', width: 150, ellipsis: { tooltip: true } },
  { title: 'shipped_at', path: 'shipped_at', width: 190, render: row => formatDateTime(row.shipped_at) },
  { title: '更新时间', path: 'updated_at', width: 190, render: row => formatDateTime(row.updated_at) },
  {
    title: '操作',
    key: 'actions',
    width: 130,
    fixed: 'right',
    render: row => (
      <NFlex size={12}>
        {hasPermission('biz:cgtestorder:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => openChildDetailModal(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('biz:cgtestorder:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openChildEditModal(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('biz:cgtestorder:delete') ? (
          <NButton type="error" size="small" text={true} onClick={() => confirmChildDelete(row.id)}>
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
    const response = await cgTestOrderApi.page({ current: state.page, size: state.pageSize, ...state.searchValues })
    const data = response.data ?? {}
    state.rows = data.records ?? []
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
    state.checkedRowKeys = state.checkedRowKeys.filter(key => state.rows.some(item => item.id === key))
  } finally {
    state.loading = false
  }
}

async function selectMaster(id: string) {
  state.selectedMasterId = id
  state.childDrawerVisible = true
  state.childCheckedRowKeys = []
  state.childPage = 1
  await fetchChildPage()
}

async function fetchChildPage() {
  state.childLoading = true
  try {
    const response = await cgTestOrderApi.childPage({
      current: state.childPage,
      size: state.childPageSize,
      order_id: state.selectedMasterId,
      ...state.childSearchValues,
    })
    const data = response.data ?? {}
    state.childRows = data.records ?? []
    const childPageMeta = readPageMeta(data, { current: state.childPage, size: state.childPageSize })
    state.childTotal = childPageMeta.total
    state.childPage = childPageMeta.current
    state.childPageSize = childPageMeta.size
    state.childCheckedRowKeys = state.childCheckedRowKeys.filter(key => state.childRows.some(item => item.id === key))
  } finally {
    state.childLoading = false
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

function openChildDetailModal(id: string) {
  childDetailModalRef.value?.openModal(id)
}

function openChildCreateModal() {
  childFormModalRef.value?.openModal(undefined, { order_id: state.selectedMasterId })
}

function openChildEditModal(id: string) {
  childFormModalRef.value?.openModal(id)
}

function handleCheckedRowKeys(keys: Array<string | number>) {
  state.checkedRowKeys = keys.map(String)
}

function handleChildCheckedRowKeys(keys: Array<string | number>) {
  state.childCheckedRowKeys = keys.map(String)
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
  await cgTestOrderApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter(key => !ids.includes(key))
  if (state.selectedMasterId && ids.includes(state.selectedMasterId)) {
    state.selectedMasterId = null
    state.childDrawerVisible = false
    state.childRows = []
    state.childTotal = 0
    state.childCheckedRowKeys = []
  }
  window.$message.success('删除成功')
  await fetchPage()
}

function confirmChildDelete(value: string | string[]) {
  const ids = Array.isArray(value) ? value : [value]
  if (!ids.length) {
    return
  }
  window.$dialog.warning({
    title: ids.length > 1 ? '批量删除' : '删除',
    content: ids.length > 1 ? `删除 ${ids.length} 条明细?` : '删除该明细?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteChildRows(ids),
  })
}

async function deleteChildRows(ids: string[]) {
  await cgTestOrderApi.childRemove({ ids })
  state.childCheckedRowKeys = state.childCheckedRowKeys.filter(key => !ids.includes(key))
  window.$message.success('删除成功')
  await fetchChildPage()
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
      />
    </ProCard>

    <ProDataTable
      class="min-h-0 flex-1"
      remote
      title="CgTestOrder"
      row-key="id"
      :scroll-x="1300"
      :columns="tableColumns"
      :data="state.rows"
      :loading="state.loading"
      :pagination="pagination"
      :checked-row-keys="state.checkedRowKeys"
      :on-update-checked-row-keys="handleCheckedRowKeys"
    >
      <template #toolbar>
        <NFlex>
          <NButton v-if="hasPermission('biz:cgtestorder:create')" type="primary" text @click="openCreateModal">
            <template #icon><NIcon><Icon icon="icon-park-outline:plus" /></NIcon></template>
          </NButton>
          <NButton text :loading="state.loading" @click="fetchPage">
            <template #icon><NIcon><Icon icon="icon-park-outline:refresh" /></NIcon></template>
          </NButton>
          <NButton v-if="hasPermission('biz:cgtestorder:delete')" type="error" text :disabled="!hasCheckedRows" @click="confirmDelete(state.checkedRowKeys)">
            <template #icon><NIcon><Icon icon="icon-park-outline:delete" /></NIcon></template>
          </NButton>
        </NFlex>
      </template>
    </ProDataTable>

    <NDrawer v-model:show="state.childDrawerVisible" :width="960" placement="right">
      <NDrawerContent title="订单明细管理" closable>
        <NFlex style="height: calc(100vh - 110px)" vertical>
          <ProCard content-class="pb-0!">
            <ProSearchForm
              :form="childSearchForm"
              :columns="childSearchColumns"
              :reset-button-props="{ content: '重置' }"
              :search-button-props="{ content: '搜索' }"
            />
          </ProCard>
          <ProDataTable
            class="min-h-0 flex-1"
            remote
            title="订单明细"
            row-key="id"
            :scroll-x="1300"
            :columns="childColumns"
            :data="state.childRows"
            :loading="state.childLoading"
            :pagination="childPagination"
            :checked-row-keys="state.childCheckedRowKeys"
            :on-update-checked-row-keys="handleChildCheckedRowKeys"
          >
            <template #toolbar>
              <NFlex>
                <NButton v-if="hasPermission('biz:cgtestorder:create')" type="primary" text :disabled="!canCreateChild" @click="openChildCreateModal">
                  <template #icon><NIcon><Icon icon="icon-park-outline:plus" /></NIcon></template>
                </NButton>
                <NButton text :loading="state.childLoading" @click="fetchChildPage">
                  <template #icon><NIcon><Icon icon="icon-park-outline:refresh" /></NIcon></template>
                </NButton>
                <NButton v-if="hasPermission('biz:cgtestorder:delete')" type="error" text :disabled="!hasChildCheckedRows" @click="confirmChildDelete(state.childCheckedRowKeys)">
                  <template #icon><NIcon><Icon icon="icon-park-outline:delete" /></NIcon></template>
                </NButton>
              </NFlex>
            </template>
          </ProDataTable>
        </NFlex>
      </NDrawerContent>
    </NDrawer>

    <ModalDetail ref="detailModalRef" />
    <ModalForm ref="formModalRef" @saved="fetchPage" />
    <ChildModalDetail ref="childDetailModalRef" />
    <ChildModalForm ref="childFormModalRef" @saved="fetchChildPage" />
  </NFlex>
</template>
