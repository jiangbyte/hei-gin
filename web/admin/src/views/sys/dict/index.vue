<!-- Author: Charlie -->

<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { dictApi } from '@/api'
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
import { dictList, dictTypeData, dictTypeColor, getDictLabel, refreshDict } from '@/utils/dict'
import ModalDetail from './components/ModalDetail.vue'
import ModalForm from './components/ModalForm.vue'
import { readPageMeta } from '@/utils/wire'

const detailModalRef = ref<any>(null)
const formModalRef = ref<any>(null)
const state = reactive({
  dicts: [] as any[],
  dictTree: [] as any[],
  total: 0,
  loading: false,
  treeLoading: false,
  category: 'SYS',
  treeSearchKey: '',
  selectedTreeKeys: [] as string[],
  checkedRowKeys: [] as string[],
  searchValues: {} as any,
  page: 1,
  pageSize: 20,
})

const selectedParentId = computed(() => state.selectedTreeKeys[0] ?? null)
const hasCheckedRows = computed(() => state.checkedRowKeys.length > 0)
const flatTreeDicts = computed(() => flattenDictTree(state.dictTree))

const searchForm = createProSearchForm<any>({
  defaultCollapsed: true,
  onSubmit(values) {
    state.searchValues = normalizeSearchValues(values, {
      code: (value) => String(value).trim().toUpperCase(),
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
    title: '编码',
    path: 'code',
    field: 'input',
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

const treeData = computed(() => {
  const keyword = state.treeSearchKey.trim().toLowerCase()
  return buildTreeNodes(state.dictTree, keyword)
})

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
    title: '编码',
    path: 'code',
    width: 190,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '标签',
    path: 'label',
    width: 150,
    render: (row) => getDictLabel(row),
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '值',
    path: 'value',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '颜色',
    path: 'color',
    width: 120,
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
    title: '分类',
    path: 'category',
    width: 120,
    render: (row) => dictTypeData('SYS_BIZ_CATEGORY', row.category),
  },
  {
    title: '父级字典',
    path: 'parent_id_name',
    width: 180,
    ellipsis: {
      tooltip: true,
    },
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
        {dictTypeData('COMMON_STATUS', row.status)}
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
        {hasPermission('sys:dict:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => openDetailModal(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('sys:dict:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openEditModal(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('sys:dict:delete') ? (
          <NButton type="error" size="small" text={true} onClick={() => confirmDelete(row.id)}>
            {renderButtonIcon('icon-park-outline:delete')}
          </NButton>
        ) : null}
      </NFlex>
    ),
  },
])

onMounted(() => {
  refreshData()
})

function handleCategoryUpdate(value: string) {
  state.category = value
  state.searchValues = {}
  state.treeSearchKey = ''
  state.selectedTreeKeys = []
  state.checkedRowKeys = []
  state.page = 1
  refreshData()
}

function handleTreeSelect(keys: Array<string | number>) {
  state.selectedTreeKeys = keys.map(String)
  state.checkedRowKeys = []
  state.page = 1
  fetchPage()
}

function handleCheckedRowKeys(keys: Array<string | number>) {
  state.checkedRowKeys = keys.map(String)
}

function openCreateModal() {
  formModalRef.value?.openModal(undefined, {
    category: state.category,
    parentId: selectedParentId.value,
  })
}

function openEditModal(id: string) {
  formModalRef.value?.openModal(id)
}

function openDetailModal(id: string) {
  detailModalRef.value?.openModal(id)
}

async function handleSaved() {
  state.page = 1
  await refreshData()
  await refreshDict()
}

async function refreshData() {
  await fetchTree()
  await fetchPage()
}

async function fetchTree() {
  state.treeLoading = true
  try {
    const response = await dictApi.tree({ category: state.category })
    state.dictTree = response.data ?? []
    if (
      selectedParentId.value &&
      !flatTreeDicts.value.some((item) => item.id === selectedParentId.value)
    ) {
      state.selectedTreeKeys = []
    }
  } finally {
    state.treeLoading = false
  }
}

async function fetchPage() {
  state.loading = true
  try {
    const params: any = {
      current: state.page,
      size: state.pageSize,
      category: state.category,
      ...state.searchValues,
    }
    if (selectedParentId.value) {
      params.parent_id = selectedParentId.value
    }

    const response = await dictApi.page(params)
    const data = response.data ?? {}
    state.dicts = data.records ?? []
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
    state.checkedRowKeys = state.checkedRowKeys.filter((key) =>
      state.dicts.some((item) => item.id === key),
    )
  } finally {
    state.loading = false
  }
}

function confirmDelete(value: string | string[]) {
  const ids = Array.isArray(value) ? value : [value]
  if (!ids.length) {
    return
  }
  const deleteIds = collectDeleteIds(ids)
  const isBatch = ids.length > 1

  window.$dialog.warning({
    title: isBatch ? '批量删除' : '删除',
    draggable: true,
    maskClosable: false,
    content: isBatch ? `删除 ${deleteIds.length} 个字典?` : '删除该字典及其子级?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteData(deleteIds),
  })
}

async function deleteData(ids: string[]) {
  await dictApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))
  if (selectedParentId.value && ids.includes(selectedParentId.value)) {
    state.selectedTreeKeys = []
  }
  window.$message.success('删除成功')

  await refreshData()
  await refreshDict()
  if (!state.dicts.length && state.total > 0 && state.page > 1) {
    state.page -= 1
    await fetchPage()
  }
}

function collectDeleteIds(ids: string[]) {
  const result = new Set(ids)
  const walk = (parentId: string) => {
    flatTreeDicts.value
      .filter((item) => item.parent_id === parentId)
      .forEach((item) => {
        result.add(item.id)
        walk(item.id)
      })
  }
  ids.forEach(walk)
  return Array.from(result)
}

function buildTreeNodes(items: any[], keyword: string): any[] {
  const nodes = items.map((item) => ({
    key: item.id,
    label: getDictLabel(item),
    children: buildTreeNodes(item.children ?? [], keyword),
    sort: item.sort ?? 0,
    id: item.id,
    raw: item,
  }))

  return sortAndFilterTree(nodes, keyword)
}

function sortAndFilterTree(nodes: any[], keyword: string): any[] {
  return nodes
    .sort(sortTreeNodes)
    .map((node) => ({
      ...node,
      children: sortAndFilterTree(node.children, keyword),
    }))
    .filter((node) => {
      if (!keyword) {
        return true
      }
      const raw = node.raw
      if (!raw) return false
      return (
        (raw.code ?? '').toLowerCase().includes(keyword) ||
        String(raw.label ?? '')
          .toLowerCase()
          .includes(keyword) ||
        node.children.length > 0
      )
    })
    .map((node) => ({
      key: node.key,
      label: node.label,
      children: node.children,
      sort: node.sort,
      id: node.id,
    }))
}

function sortTreeNodes(a: any, b: any) {
  return a.sort === b.sort ? String(b.id).localeCompare(String(a.id)) : a.sort - b.sort
}

function flattenDictTree(items: any[]) {
  const result: any[] = []
  const walk = (nodes: any[]) => {
    nodes.forEach((node) => {
      result.push(node)
      walk(node.children ?? [])
    })
  }
  walk(items)
  return result
}
</script>

<template>
  <div class="dict-page">
    <ProCard
      class="dict-tree-card"
      content-class="h-full min-h-0 overflow-hidden"
    >
      <NFlex
        class="dict-tree-layout"
        vertical
        :size="12"
      >
        <NInput
          v-model:value="state.treeSearchKey"
          clearable
          :placeholder="'搜索 dict'"
        >
          <template #prefix>
            <NIcon>
              <Icon icon="icon-park-outline:search" />
            </NIcon>
          </template>
        </NInput>
        <NTabs
          :value="state.category"
          type="line"
          animated
          justify-content="space-evenly"
          @update:value="handleCategoryUpdate"
        >
          <NTabPane
            name="SYS"
            :tab="'系统'"
          />
          <NTabPane
            name="BIZ"
            :tab="'业务'"
          />
        </NTabs>
        <div class="dict-tree-body">
          <NSpin
            :show="state.treeLoading"
            class="dict-tree-spin"
            content-class="dict-tree-spin-content"
          >
            <NScrollbar
              class="dict-tree-scroll"
              content-class="dict-tree-scroll-content"
            >
              <NTree
                block-line
                block-node
                show-line
                :data="treeData"
                :selected-keys="state.selectedTreeKeys"
                key-field="key"
                label-field="label"
                children-field="children"
                @update:selected-keys="handleTreeSelect"
              />
            </NScrollbar>
          </NSpin>
        </div>
      </NFlex>
    </ProCard>

    <NFlex
      class="min-w-0 min-h-0 h-full"
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
        :title="'字典管理'"
        row-key="id"
        :scroll-x="1510"
        :columns="tableColumns"
        :data="state.dicts"
        :loading="state.loading"
        :pagination="pagination"
        :checked-row-keys="state.checkedRowKeys"
        :on-update-checked-row-keys="handleCheckedRowKeys"
      >
        <template #toolbar>
          <NFlex>
            <NButton
              v-if="hasPermission('sys:dict:create')"
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
              :loading="state.loading || state.treeLoading"
              @click="refreshData"
            >
              <template #icon>
                <NIcon>
                  <Icon icon="icon-park-outline:reload" />
                </NIcon>
              </template>
            </NButton>
            <NButton
              v-if="hasPermission('sys:dict:delete')"
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

    <ModalForm
      ref="formModalRef"
      :dicts="flatTreeDicts"
      @saved="handleSaved"
    />
    <ModalDetail ref="detailModalRef" />
  </div>
</template>

<style scoped>
.dict-page {
  display: grid;
  grid-template-columns: minmax(260px, 320px) minmax(0, 1fr);
  grid-template-rows: minmax(0, 1fr);
  gap: 16px;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.dict-page > * {
  min-height: 0;
  min-width: 0;
}

.dict-tree-card {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.dict-tree-card :deep(.n-card__content),
.dict-tree-card :deep(.n-card-content) {
  display: flex;
  flex: 1 1 0;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

.dict-tree-card :deep(.n-pro-collapse-transition) {
  display: flex;
  flex: 1 1 0;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

.dict-tree-layout {
  flex: 1 1 0;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.dict-tree-layout > :not(.dict-tree-body) {
  flex-shrink: 0;
}

.dict-tree-body {
  flex: 1 1 0;
  min-height: 0;
  overflow: hidden;
}

.dict-tree-spin,
.dict-tree-body :deep(.dict-tree-spin-content),
.dict-tree-body :deep(.n-spin-container),
.dict-tree-body :deep(.n-spin-content) {
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.dict-tree-scroll {
  flex: 1 1 0;
  height: 100%;
  min-height: 0;
}

.dict-tree-scroll :deep(.n-scrollbar-container),
.dict-tree-scroll :deep(.dict-tree-scroll-content) {
  min-height: 0;
}

@media (max-width: 900px) {
  .dict-page {
    grid-template-columns: 1fr;
    grid-template-rows: minmax(280px, 36vh) minmax(0, 1fr);
  }
}
</style>
