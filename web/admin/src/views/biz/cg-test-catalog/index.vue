<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-08 21:09:53
-->

<script setup lang="tsx">
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { cgTestCatalogApi } from '@/api'
import { formatDateTime, hasPermission, normalizeSearchValues, renderButtonIcon } from '@/utils'
import { NButton, NFlex, NIcon, NTag } from 'naive-ui'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'
import ModalDetail from './components/ModalDetail.vue'
import ModalForm from './components/ModalForm.vue'

const formModalRef = ref<any>(null)
const detailModalRef = ref<any>(null)
const state = reactive({
  searchValues: {} as any,
  checkedRowKeys: [] as string[],
  treeRows: [] as any[],
  treeLoading: false,
})

const hasCheckedRows = computed(() => state.checkedRowKeys.length > 0)
const filteredTreeRows = computed(() => filterTreeRows(state.treeRows, state.searchValues))

const searchForm = createProSearchForm<any>({
  defaultCollapsed: true,
  onSubmit(values) {
    state.searchValues = normalizeSearchValues(values)
  },
  onReset() {
    state.searchValues = {}
  },
})

const searchColumns = computed<ProSearchFormColumns<any>>(() => [
  { title: '目录编码', path: 'code', field: 'input' },
  { title: '目录名称', path: 'name', field: 'input' },
  { title: '目录分类', path: 'category', field: 'input' },
  { title: '状态', path: 'status', field: 'input' },
])

const tableColumns = computed<ProDataTableColumns<any>>(() => [
  { type: 'selection', fixed: 'left' },
  { title: '主键', path: 'id', width: 150, ellipsis: { tooltip: true } },
  { title: '目录编码', path: 'code', width: 150, ellipsis: { tooltip: true } },
  { title: '目录名称', path: 'name', width: 150, ellipsis: { tooltip: true } },
  { title: '目录分类', path: 'category', width: 150, ellipsis: { tooltip: true } },
  { title: '状态', path: 'status', width: 150, ellipsis: { tooltip: true } },
  { title: '排序', path: 'sort', width: 150, ellipsis: { tooltip: true } },
  {
    title: '是否显示',
    path: 'is_visible',
    width: 120,
    render: (row) => (
      <NTag type={row.is_visible === 'true' ? 'success' : 'default'} bordered={false}>
        {row.is_visible === 'true' ? '是' : row.is_visible === 'false' ? '否' : '-'}
      </NTag>
    ),
  },
  { title: '图标', path: 'icon', width: 150, ellipsis: { tooltip: true } },
  {
    title: '更新时间',
    path: 'updated_at',
    width: 190,
    render: (row) => formatDateTime(row.updated_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 130,
    fixed: 'right',
    render: (row) => (
      <NFlex size={12}>
        {hasPermission('biz:cgtestcatalog:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => openDetailModal(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('biz:cgtestcatalog:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openEditModal(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('biz:cgtestcatalog:delete') ? (
          <NButton type="error" size="small" text={true} onClick={() => confirmDelete(row.id)}>
            {renderButtonIcon('icon-park-outline:delete')}
          </NButton>
        ) : null}
      </NFlex>
    ),
  },
])

onMounted(() => {
  fetchTree()
})

async function fetchTree() {
  state.treeLoading = true
  try {
    const response = await cgTestCatalogApi.tree()
    state.treeRows = response.data ?? []
  } finally {
    state.treeLoading = false
  }
}

function filterTreeRows(items: any[], searchValues: any): any[] {
  return items
    .map((item) => {
      const children = filterTreeRows(item.children ?? [], searchValues)
      if (matchesTreeRow(item, searchValues) || children.length) {
        return { ...item, children }
      }
      return null
    })
    .filter(Boolean)
}

function matchesTreeRow(item: any, searchValues: any) {
  const conditions = [
    containsValue(item.code, searchValues.code),
    containsValue(item.name, searchValues.name),
    containsValue(item.category, searchValues.category),
    equalsValue(item.status, searchValues.status),
  ]
  return conditions.every(Boolean)
}

function containsValue(source: unknown, target: unknown) {
  if (target === undefined || target === null || target === '') {
    return true
  }
  return String(source ?? '')
    .toLowerCase()
    .includes(String(target).toLowerCase())
}

function equalsValue(source: unknown, target: unknown) {
  if (target === undefined || target === null || target === '') {
    return true
  }
  return String(source ?? '') === String(target)
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
  window.$dialog.warning({
    title: ids.length > 1 ? '批量删除' : '删除',
    content: ids.length > 1 ? `删除 ${ids.length} 条记录?` : '删除该记录?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteRows(ids),
  })
}

async function deleteRows(ids: string[]) {
  await cgTestCatalogApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))
  window.$message.success('删除成功')
  await fetchTree()
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
      title="CgTestCatalog"
      row-key="id"
      :scroll-x="1300"
      :columns="tableColumns"
      :data="filteredTreeRows"
      :loading="state.treeLoading"
      :pagination="false"
      default-expand-all
      :checked-row-keys="state.checkedRowKeys"
      :on-update-checked-row-keys="handleCheckedRowKeys"
    >
      <template #toolbar>
        <NFlex>
          <NButton
            v-if="hasPermission('biz:cgtestcatalog:create')"
            type="primary"
            text
            @click="openCreateModal"
          >
            <template #icon>
              <NIcon><Icon icon="icon-park-outline:plus" /></NIcon>
            </template>
          </NButton>
          <NButton
            text
            :loading="state.treeLoading"
            @click="fetchTree"
          >
            <template #icon>
              <NIcon><Icon icon="icon-park-outline:refresh" /></NIcon>
            </template>
          </NButton>
          <NButton
            v-if="hasPermission('biz:cgtestcatalog:delete')"
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

    <ModalDetail ref="detailModalRef" />
    <ModalForm
      ref="formModalRef"
      @saved="fetchTree"
    />
  </NFlex>
</template>
