<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-08 21:09:55
-->

<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { cgTestKnowledgeCategoryApi } from '@/api'
import { readPageMeta } from '@/utils/wire'
import { formatDateTime, hasPermission, normalizeSearchValues, renderButtonIcon } from '@/utils'
import { NButton, NFlex, NIcon, NInput, NInputGroup } from 'naive-ui'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'
import ChildModalDetail from './components/children/ChildModalDetail.vue'
import ChildModalForm from './components/children/ChildModalForm.vue'

const childFormModalRef = ref<any>(null)
const childDetailModalRef = ref<any>(null)
const state = reactive({
  treeRows: [] as any[],
  treeLoading: false,
  selectedTreeKeys: [] as string[],
  treeKeyword: '',
  childRows: [] as any[],
  childTotal: 0,
  childLoading: false,
  childSearchValues: {} as any,
  childCheckedRowKeys: [] as string[],
  childPage: 1,
  childPageSize: 20,
  selectedMasterId: null as string | null,
})

const treeData = computed(() => buildTreeNodes(state.treeRows))
const hasChildCheckedRows = computed(() => state.childCheckedRowKeys.length > 0)
const canCreateChild = computed(() => Boolean(state.selectedMasterId))

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
  { title: '分类ID', path: 'category_id', field: 'input' },
  { title: '文档编码', path: 'code', field: 'input' },
  { title: '文档标题', path: 'title', field: 'input' },
  { title: '文档类型', path: 'type', field: 'input' },
  { title: '状态', path: 'status', field: 'input' },
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

const childColumns = computed<ProDataTableColumns<any>>(() => [
  { type: 'selection', fixed: 'left' },
  { title: '主键', path: 'id', width: 150, ellipsis: { tooltip: true } },
  { title: '分类ID', path: 'category_id', width: 150, ellipsis: { tooltip: true } },
  { title: '文档编码', path: 'code', width: 150, ellipsis: { tooltip: true } },
  { title: '文档标题', path: 'title', width: 150, ellipsis: { tooltip: true } },
  { title: '文档类型', path: 'type', width: 150, ellipsis: { tooltip: true } },
  { title: '状态', path: 'status', width: 150, ellipsis: { tooltip: true } },
  { title: '摘要', path: 'summary', width: 150, ellipsis: { tooltip: true } },
  { title: '正文内容', path: 'content', width: 150, ellipsis: { tooltip: true } },
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
        {hasPermission('biz:cgtestknowledgecategory:detail') ? (
          <NButton
            type="info"
            size="small"
            text={true}
            onClick={() => openChildDetailModal(row.id)}
          >
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('biz:cgtestknowledgecategory:update') ? (
          <NButton
            type="primary"
            size="small"
            text={true}
            onClick={() => openChildEditModal(row.id)}
          >
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('biz:cgtestknowledgecategory:delete') ? (
          <NButton type="error" size="small" text={true} onClick={() => confirmChildDelete(row.id)}>
            {renderButtonIcon('icon-park-outline:delete')}
          </NButton>
        ) : null}
      </NFlex>
    ),
  },
])

onMounted(() => {
  fetchTree()
  fetchChildPage()
})

async function fetchTree() {
  state.treeLoading = true
  try {
    const response = await cgTestKnowledgeCategoryApi.tree({
      keyword: state.treeKeyword || undefined,
    })
    state.treeRows = response.data ?? []
  } finally {
    state.treeLoading = false
  }
}

async function searchTree() {
  state.selectedTreeKeys = []
  state.selectedMasterId = null
  state.childPage = 1
  await Promise.all([fetchTree(), fetchChildPage()])
}

async function resetTreeSearch() {
  state.treeKeyword = ''
  await searchTree()
}

function handleTreeSelect(keys: Array<string | number>) {
  state.selectedTreeKeys = keys.map(String)
  state.selectedMasterId = state.selectedTreeKeys[0] ?? null
  state.childPage = 1
  fetchChildPage()
}

function buildTreeNodes(items: any[]): any[] {
  return items.map((item) => ({
    key: item.id,
    label: item.name,
    children: item.children?.length ? buildTreeNodes(item.children) : undefined,
  }))
}

async function fetchChildPage() {
  state.childLoading = true
  try {
    const response = await cgTestKnowledgeCategoryApi.childPage({
      current: state.childPage,
      size: state.childPageSize,
      category_id: state.selectedMasterId,
      ...state.childSearchValues,
    })
    const data = response.data ?? {}
    state.childRows = data.records ?? []
    const childPageMeta = readPageMeta(data, {
      current: state.childPage,
      size: state.childPageSize,
    })
    state.childTotal = childPageMeta.total
    state.childPage = childPageMeta.current
    state.childPageSize = childPageMeta.size
    state.childCheckedRowKeys = state.childCheckedRowKeys.filter((key) =>
      state.childRows.some((item) => item.id === key),
    )
  } finally {
    state.childLoading = false
  }
}

function openChildDetailModal(id: string) {
  childDetailModalRef.value?.openModal(id)
}

function openChildCreateModal() {
  childFormModalRef.value?.openModal(undefined, { category_id: state.selectedMasterId })
}

function openChildEditModal(id: string) {
  childFormModalRef.value?.openModal(id)
}

function handleChildCheckedRowKeys(keys: Array<string | number>) {
  state.childCheckedRowKeys = keys.map(String)
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
  await cgTestKnowledgeCategoryApi.childRemove({ ids })
  state.childCheckedRowKeys = state.childCheckedRowKeys.filter((key) => !ids.includes(key))
  window.$message.success('删除成功')
  await fetchChildPage()
}
</script>

<template>
  <div class="generated-left-tree-table">
    <ProCard
      class="generated-tree"
      content-class="h-full min-h-0 overflow-hidden"
    >
      <NFlex
        class="generated-tree-layout"
        vertical
        :size="12"
      >
        <NInputGroup>
          <NInput
            v-model:value="state.treeKeyword"
            clearable
            placeholder="搜索CgTestKnowledgeCategory"
            @keyup.enter="searchTree"
          />
          <NButton
            type="primary"
            :loading="state.treeLoading"
            title="搜索"
            @click="searchTree"
          >
            <template #icon>
              <NIcon><Icon icon="icon-park-outline:search" /></NIcon>
            </template>
          </NButton>
          <NButton
            :disabled="!state.treeKeyword"
            title="重置"
            @click="resetTreeSearch"
          >
            <template #icon>
              <NIcon><Icon icon="icon-park-outline:refresh" /></NIcon>
            </template>
          </NButton>
        </NInputGroup>
        <div class="generated-tree-body">
          <NSpin
            :show="state.treeLoading"
            class="generated-tree-spin"
            content-class="generated-tree-spin-content"
          >
            <NScrollbar
              class="generated-tree-scroll"
              content-class="generated-tree-scroll-content"
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
          :form="childSearchForm"
          :columns="childSearchColumns"
          :reset-button-props="{ content: '重置' }"
          :search-button-props="{ content: '搜索' }"
        />
      </ProCard>
      <ProDataTable
        class="min-h-0 flex-1"
        remote
        title="CgTestKnowledgeDoc"
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
            <NButton
              v-if="hasPermission('biz:cgtestknowledgecategory:create')"
              type="primary"
              text
              :disabled="!canCreateChild"
              @click="openChildCreateModal"
            >
              <template #icon>
                <NIcon><Icon icon="icon-park-outline:plus" /></NIcon>
              </template>
            </NButton>
            <NButton
              text
              :loading="state.childLoading"
              @click="fetchChildPage"
            >
              <template #icon>
                <NIcon><Icon icon="icon-park-outline:refresh" /></NIcon>
              </template>
            </NButton>
            <NButton
              v-if="hasPermission('biz:cgtestknowledgecategory:delete')"
              type="error"
              text
              :disabled="!hasChildCheckedRows"
              @click="confirmChildDelete(state.childCheckedRowKeys)"
            >
              <template #icon>
                <NIcon><Icon icon="icon-park-outline:delete" /></NIcon>
              </template>
            </NButton>
          </NFlex>
        </template>
      </ProDataTable>
    </NFlex>

    <ChildModalDetail ref="childDetailModalRef" />
    <ChildModalForm
      ref="childFormModalRef"
      @saved="fetchChildPage"
    />
  </div>
</template>

<style scoped>
.generated-left-tree-table {
  display: grid;
  grid-template-columns: 300px minmax(0, 1fr);
  gap: 12px;
  height: 100%;
  min-height: 0;
}

.generated-tree {
  min-height: 0;
}

.generated-tree-layout {
  height: 100%;
  min-height: 0;
}

.generated-tree-body {
  min-height: 0;
  flex: 1;
}

.generated-tree-spin,
.generated-tree-spin :deep(.generated-tree-spin-content),
.generated-tree-scroll {
  height: 100%;
  min-height: 0;
}

.generated-tree-scroll :deep(.generated-tree-scroll-content) {
  min-width: max-content;
  padding-right: 8px;
}

@media (max-width: 900px) {
  .generated-left-tree-table {
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: minmax(260px, 34vh) minmax(0, 1fr);
  }
}
</style>
