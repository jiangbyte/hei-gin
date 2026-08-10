<!-- Author: Charlie -->

<script setup lang="tsx">
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { deptApi } from '@/api'
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
import { dictList, dictTypeData, dictTypeColor, isDictLoaded, refreshDict } from '@/utils/dict'
import ModalDetail from './components/ModalDetail.vue'
import ModalForm from './components/ModalForm.vue'

const formModalRef = ref<any>(null)
const detailModalRef = ref<any>(null)
const state = reactive({
  depts: [] as any[],
  loading: false,
  searchValues: {} as any,
  checkedRowKeys: [] as string[],
})

const hasCheckedRows = computed(() => state.checkedRowKeys.length > 0)
const filteredDepts = computed(() => filterDeptTree(state.depts, state.searchValues))

const searchForm = createProSearchForm<any>({
  defaultCollapsed: true,
  onSubmit(values) {
    state.searchValues = normalizeSearchValues(values, {
      name: (value) => String(value).trim(),
      category: (value) => String(value).trim(),
      status: (value) => String(value).trim(),
    })
  },
  onReset() {
    state.searchValues = {}
  },
})

const searchColumns = computed<ProSearchFormColumns<any>>(() => [
  {
    title: '部门名称',
    path: 'name',
    field: 'input',
  },
  {
    title: '部门分类',
    path: 'category',
    field: 'select',
    fieldProps: {
      options: dictList('DEPT_CATEGORY'),
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

const tableColumns = computed<ProDataTableColumns<any>>(() => [
  {
    type: 'selection',
    fixed: 'left',
  },
  {
    title: '部门名称',
    path: 'name',
    width: 220,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '部门分类',
    path: 'category',
    width: 130,
    render: (row) => dictTypeData('DEPT_CATEGORY', row.category) || row.category,
  },
  {
    title: '负责人',
    width: 120,
    render: (row) => row.master_name || row.master_id || '-',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '副负责人',
    width: 120,
    render: (row) => row.deputy_master_name || row.deputy_master_id || '-',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '排序',
    path: 'sort',
    width: 80,
  },
  {
    title: '虚拟部门',
    path: 'is_virtual',
    width: 100,
    render: (row) => (row.is_virtual ? '是' : '否'),
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
    width: 180,
    render: (row) => formatDateTime(row.updated_at),
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    fixed: 'right',
    render: (row) => (
      <NFlex size={12}>
        {hasPermission('iam:dept:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => openDetailModal(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('iam:dept:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openEditModal(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('iam:dept:delete') ? (
          <NButton type="error" size="small" text={true} onClick={() => confirmDelete(row.id)}>
            {renderButtonIcon('icon-park-outline:delete')}
          </NButton>
        ) : null}
      </NFlex>
    ),
  },
])

onMounted(async () => {
  if (!isDictLoaded()) {
    await refreshDict()
  }
  fetchTree()
})

async function fetchTree() {
  state.loading = true
  try {
    const response = await deptApi.tree()
    state.depts = response.data ?? []
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
  if (!ids.length) return
  window.$dialog.warning({
    title: ids.length > 1 ? '批量删除' : '删除',
    draggable: true,
    maskClosable: false,
    content: ids.length > 1 ? `删除 ${ids.length} 个部门?` : '删除该部门?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteData(ids),
  })
}

async function deleteData(ids: string[]) {
  await deptApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))
  window.$message.success('删除成功')
  await fetchTree()
}

function filterDeptTree(items: any[], searchValues: any): any[] {
  const keyword = (searchValues.name || '').toLowerCase().trim()
  const codeKeyword = (searchValues.code || '').toLowerCase().trim()
  const categoryVal = searchValues.category || ''
  const statusVal = searchValues.status || ''
  if (!keyword && !codeKeyword && !categoryVal && !statusVal) return items
  return items
    .map((item) => {
      const match =
        (!keyword || (item.name || '').toLowerCase().includes(keyword)) &&
        (!codeKeyword || (item.code || '').toLowerCase().includes(codeKeyword)) &&
        (!categoryVal || item.category === categoryVal) &&
        (!statusVal || item.status === statusVal)
      const matchedChildren = item.children?.length
        ? filterDeptTree(item.children, searchValues)
        : []
      if (match || matchedChildren.length) {
        return { ...item, children: matchedChildren }
      }
      return null
    })
    .filter(Boolean)
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
      :title="'部门管理'"
      row-key="id"
      :scroll-x="1440"
      :columns="tableColumns"
      :data="filteredDepts"
      :loading="state.loading"
      :pagination="false"
      :checked-row-keys="state.checkedRowKeys"
      :on-update-checked-row-keys="handleCheckedRowKeys"
      default-expand-all
    >
      <template #toolbar>
        <NFlex>
          <NButton
            v-if="hasPermission('iam:dept:create')"
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
            @click="fetchTree"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:reload" />
              </NIcon>
            </template>
          </NButton>
          <NButton
            v-if="hasPermission('iam:dept:delete')"
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
      @saved="fetchTree"
    />
    <ModalDetail ref="detailModalRef" />
  </NFlex>
</template>

<style scoped></style>
