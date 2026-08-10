<!-- Author: Charlie -->

<script setup lang="tsx">
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { resourceApi, resourceModuleApi } from '@/api'
import NovaIcon from '@/components/common/NovaIcon.vue'
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
import { dictList, dictTypeData, dictTypeColor } from '@/utils/dict'
import ModalButtonPermission from './components/ModalButtonPermission.vue'
import ModalDetail from './components/ModalDetail.vue'
import ModalForm from './components/ModalForm.vue'

const formModalRef = ref<any>(null)
const detailModalRef = ref<any>(null)
const buttonPermissionModalRef = ref<any>(null)
const state = reactive({
  resources: [] as any[],
  modules: [] as any[],
  activeModuleId: null as string | null,
  moduleOptions: [] as any[],
  loading: false,
  searchValues: {} as any,
  checkedRowKeys: [] as string[],
})

const hasCheckedRows = computed(() => state.checkedRowKeys.length > 0)
const filteredResources = computed(() => filterResourceTree(state.resources, state.searchValues))

const searchForm = createProSearchForm<any>({
  defaultCollapsed: true,
  onSubmit(values) {
    state.searchValues = normalizeSearchValues(values, {
      code: (value) => String(value).trim(),
      name: (value) => String(value).trim(),
      module_id: (value) => String(value).trim(),
      parent_id: (value) => String(value).trim(),
    })
  },
  onReset() {
    state.searchValues = {}
  },
})

const searchColumns = computed<ProSearchFormColumns<any>>(() => [
  {
    title: '资源编码',
    path: 'code',
    field: 'input',
  },
  {
    title: '资源名称',
    path: 'name',
    field: 'input',
  },
  {
    title: '资源类型',
    path: 'resource_type',
    field: 'select',
    fieldProps: {
      options: dictList('RESOURCE_TYPE'),
    },
  },
  {
    title: '资源模块',
    path: 'module_id',
    field: 'select',
    fieldProps: {
      options: state.moduleOptions,
    },
  },
  {
    title: '父级资源ID',
    path: 'parent_id',
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

const tableColumns = computed<ProDataTableColumns<any>>(() => [
  {
    type: 'selection',
    fixed: 'left',
  },
  {
    title: '资源名称',
    path: 'name',
    width: 220,
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
    title: '资源编码',
    path: 'code',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '资源类型',
    path: 'resource_type',
    width: 100,
    render: (row) => dictTypeData('RESOURCE_TYPE', row.resource_type) || row.resource_type,
  },
  {
    title: '资源模块',
    path: 'module_id_name',
    width: 130,
    ellipsis: {
      tooltip: true,
    },
    render: (row) => row.module_id_name || row.module_id || '-',
  },
  {
    title: '路由路径',
    path: 'path',
    width: 100,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '组件',
    path: 'component',
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
    title: '排序',
    path: 'sort',
    width: 90,
  },
  {
    title: '可见',
    path: 'is_visible',
    width: 90,
    render: (row) => (row.is_visible ? '是' : '否'),
  },
  {
    title: '状态',
    path: 'status',
    width: 100,
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
    width: 150,
    fixed: 'right',
    render: (row) => {
      return (
        <NFlex size={12}>
          {hasPermission('iam:resource:detail') ? (
            <NButton type="info" size="small" text={true} onClick={() => openDetailModal(row.id)}>
              {renderButtonIcon('icon-park-outline:preview-open')}
            </NButton>
          ) : null}
          {hasPermission('iam:resource:update') ? (
            <NButton type="primary" size="small" text={true} onClick={() => openEditModal(row.id)}>
              {renderButtonIcon('icon-park-outline:edit')}
            </NButton>
          ) : null}
          {row.resource_type !== 'BUTTON' &&
          row.resource_type !== 'ACTION' &&
          hasPermission('iam:resource:create') ? (
            <NButton
              type="warning"
              size="small"
              text={true}
              onClick={() => openCreateModal(row.id)}
            >
              {renderButtonIcon('icon-park-outline:plus')}
            </NButton>
          ) : null}
          {row.resource_type !== 'BUTTON' &&
          row.resource_type !== 'ACTION' &&
          hasPermission('iam:resource:list') ? (
            <NButton
              type="warning"
              size="small"
              text={true}
              onClick={() => openButtonPermissionModal(row)}
            >
              {renderButtonIcon('icon-park-outline:config')}
            </NButton>
          ) : null}
          {hasPermission('iam:resource:delete') ? (
            <NButton type="error" size="small" text={true} onClick={() => confirmDelete(row.id)}>
              {renderButtonIcon('icon-park-outline:delete')}
            </NButton>
          ) : null}
        </NFlex>
      )
    },
  },
])

onMounted(() => {
  initialize()
})

async function initialize() {
  await fetchModules()
  await fetchTree()
}

async function fetchTree() {
  if (!state.activeModuleId) {
    state.resources = []
    state.checkedRowKeys = []
    return
  }
  state.loading = true
  try {
    const response = await resourceApi.tree({
      module_id: state.activeModuleId,
    })
    state.resources = response.data ?? []
    const existingIds = new Set(flattenResourceTree(state.resources).map((item) => item.id))
    state.checkedRowKeys = state.checkedRowKeys.filter((key) => existingIds.has(key))
  } finally {
    state.loading = false
  }
}

async function fetchModules() {
  const response = await resourceModuleApi.selector()
  state.modules = response.data ?? []
  state.moduleOptions = state.modules.map((item: any) => ({
    label: item.name,
    value: item.id,
  }))
  if (!state.modules.some((item) => item.id === state.activeModuleId)) {
    state.activeModuleId = state.modules[0]?.id ?? null
  }
}

async function handleModuleChange(moduleId: string | number) {
  state.activeModuleId = String(moduleId)
  state.checkedRowKeys = []
  await fetchTree()
}

function openDetailModal(id: string) {
  detailModalRef.value?.openModal(id)
}

function openCreateModal(parentId?: string) {
  formModalRef.value?.openModal(undefined, parentId, state.activeModuleId)
}

function openEditModal(id: string) {
  formModalRef.value?.openModal(id, undefined, state.activeModuleId)
}

function openButtonPermissionModal(row: any) {
  buttonPermissionModalRef.value?.openModal(row)
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
    content: isBatch ? `删除 ${ids.length} 个资源?` : '删除该资源?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteData(ids),
  })
}

async function deleteData(ids: string[]) {
  await resourceApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))
  window.$message.success('删除成功')
  await fetchTree()
}

function filterResourceTree(items: any[], searchValues: any): any[] {
  return items
    .map((item) => {
      const children = filterResourceTree(item.children ?? [], searchValues)
      if (matchesResource(item, searchValues) || children.length) {
        return {
          ...item,
          children,
        }
      }
      return null
    })
    .filter(Boolean)
}

function matchesResource(item: any, searchValues: any) {
  return (
    containsValue(item.code, searchValues.code) &&
    containsValue(item.name, searchValues.name) &&
    equalsValue(item.module_id, searchValues.module_id) &&
    containsValue(item.parent_id, searchValues.parent_id) &&
    equalsValue(item.resource_type, searchValues.resource_type) &&
    equalsValue(item.status, searchValues.status)
  )
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

function flattenResourceTree(items: any[]) {
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

    <ProCard
      v-if="state.moduleOptions.length"
      content-class="py-12px!"
    >
      <div class="max-w-full overflow-x-auto">
        <NButtonGroup>
          <NButton
            v-for="option in state.moduleOptions"
            :key="option.value"
            :type="state.activeModuleId === option.value ? 'primary' : 'default'"
            :secondary="state.activeModuleId !== option.value"
            :focusable="false"
            @click="handleModuleChange(option.value)"
          >
            {{ option.label }}
          </NButton>
        </NButtonGroup>
      </div>
    </ProCard>

    <ProDataTable
      class="min-h-0 flex-1"
      :title="'资源管理'"
      row-key="id"
      :scroll-x="1800"
      :columns="tableColumns"
      :data="filteredResources"
      :loading="state.loading"
      :pagination="false"
      :checked-row-keys="state.checkedRowKeys"
      :on-update-checked-row-keys="handleCheckedRowKeys"
      default-expand-all
    >
      <template #toolbar>
        <NFlex>
          <NButton
            v-if="hasPermission('iam:resource:create')"
            type="primary"
            text
            :title="'新增'"
            :aria-label="'新增'"
            @click="openCreateModal()"
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
            v-if="hasPermission('iam:resource:delete')"
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
    <ModalButtonPermission ref="buttonPermissionModalRef" />
  </NFlex>
</template>

<style scoped></style>
