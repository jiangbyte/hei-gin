<!-- Author: Charlie -->

<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { roleApi } from '@/api'
import {
  createTagColor,
  formatDateTime,
  hasPermission,
  normalizeSearchValues,
  renderButtonIcon,
} from '@/utils'
import { NButton, NDropdown, NFlex, NIcon, NTag } from 'naive-ui'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'
import { dictList, dictTypeData, dictTypeColor } from '@/utils/dict'
import ModalDetail from './components/ModalDetail.vue'
import ModalForm from './components/ModalForm.vue'
import ModalGrantClientResource from './components/ModalGrantClientResource.vue'
import ModalGrantResource from './components/ModalGrantResource.vue'
import ModalGrantUser from '../components/ModalGrantUser.vue'
import { readPageMeta } from '@/utils/wire'

const formModalRef = ref<any>(null)
const detailModalRef = ref<any>(null)
const grantResourceModalRef = ref<any>(null)
const grantClientResourceModalRef = ref<any>(null)
const grantUserModalRef = ref<any>(null)
const state = reactive({
  roles: [] as any[],
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
      code: (value) => String(value).trim(),
      name: (value) => String(value).trim(),
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
    title: '角色名称',
    path: 'name',
    field: 'input',
  },
  {
    title: '角色编码',
    path: 'code',
    field: 'input',
  },
  {
    title: '角色分类',
    path: 'category',
    field: 'select',
    fieldProps: {
      options: dictList('SYS_BIZ_CATEGORY'),
    },
  },
  {
    title: '范围类型',
    path: 'scope_type',
    field: 'select',
    fieldProps: {
      options: dictList('ROLE_SCOPE_TYPE'),
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
    title: '角色名称',
    path: 'name',
    width: 160,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '角色编码',
    path: 'code',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '角色分类',
    path: 'category',
    width: 130,
    render: (row) => dictTypeData('SYS_BIZ_CATEGORY', row.category) || row.category,
  },
  {
    title: '范围类型',
    path: 'scope_type',
    width: 130,
    render: (row) => dictTypeData('ROLE_SCOPE_TYPE', row.scope_type) || row.scope_type,
  },
  {
    title: '所属部门',
    path: 'owner_dept_name',
    width: 150,
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
    title: '内置角色',
    path: 'is_builtin',
    width: 110,
    render: (row) => (row.is_builtin ? '是' : '否'),
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
    width: 150,
    fixed: 'right',
    render: (row) => {
      const options = grantOptions.value
      return (
        <NFlex size={12}>
          {hasPermission('iam:role:detail') ? (
            <NButton type="info" size="small" text={true} onClick={() => openDetailModal(row.id)}>
              {renderButtonIcon('icon-park-outline:preview-open')}
            </NButton>
          ) : null}
          {hasPermission('iam:role:update') ? (
            <NButton type="primary" size="small" text={true} onClick={() => openEditModal(row.id)}>
              {renderButtonIcon('icon-park-outline:edit')}
            </NButton>
          ) : null}
          {options.length ? (
            <NDropdown
              trigger="click"
              options={options}
              onSelect={(key) => openGrantModal(String(key), row)}
            >
              <NButton type="warning" size="small" text={true}>
                {renderButtonIcon('icon-park-outline:permissions')}
              </NButton>
            </NDropdown>
          ) : null}
          {hasPermission('iam:role:delete') ? (
            <NButton type="error" size="small" text={true} onClick={() => confirmDelete(row.id)}>
              {renderButtonIcon('icon-park-outline:delete')}
            </NButton>
          ) : null}
        </NFlex>
      )
    },
  },
])

const grantOptions = computed(() =>
  [
    {
      label: '分配资源',
      key: 'resource',
      permission: 'iam:role:grantresource',
    },
    {
      label: '分配客户端资源',
      key: 'client-resource',
      permission: 'iam:role:grantclientresource',
    },
    {
      label: '分配用户',
      key: 'user',
      permission: 'iam:role:grantuser',
    },
  ].filter((item) => hasPermission(item.permission)),
)
const hasCheckedRows = computed(() => state.checkedRowKeys.length > 0)

onMounted(() => {
  fetchPage()
})

async function fetchPage() {
  state.loading = true
  try {
    const response = await roleApi.page({
      current: state.page,
      size: state.pageSize,
      ...state.searchValues,
    })
    const data = response.data ?? {}
    state.roles = data.records ?? []
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
    state.checkedRowKeys = state.checkedRowKeys.filter((key) =>
      state.roles.some((item) => item.id === key),
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

function openGrantModal(type: string, row: any) {
  const role = {
    id: row.id,
    code: row.code,
    name: row.name,
  }
  if (type === 'resource') {
    grantResourceModalRef.value?.openModal(role)
  } else if (type === 'client-resource') {
    grantClientResourceModalRef.value?.openModal(role, roleApi, '分配客户端资源')
  } else if (type === 'user') {
    grantUserModalRef.value?.openModal(role)
  }
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
    content: isBatch ? `删除 ${ids.length} 个角色?` : '删除该角色?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteData(ids),
  })
}

async function deleteData(ids: string[]) {
  await roleApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))

  window.$message.success('删除成功')
  await fetchPage()
  if (!state.roles.length && state.total > 0 && state.page > 1) {
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
      :title="'角色管理'"
      row-key="id"
      :scroll-x="1520"
      :columns="tableColumns"
      :data="state.roles"
      :loading="state.loading"
      :pagination="pagination"
      :checked-row-keys="state.checkedRowKeys"
      :on-update-checked-row-keys="handleCheckedRowKeys"
    >
      <template #toolbar>
        <NFlex>
          <NButton
            v-if="hasPermission('iam:role:create')"
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
            v-if="hasPermission('iam:role:delete')"
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
    <ModalGrantResource
      ref="grantResourceModalRef"
      @saved="fetchPage"
    />
    <ModalGrantClientResource
      ref="grantClientResourceModalRef"
      @saved="fetchPage"
    />
    <ModalGrantUser
      ref="grantUserModalRef"
      @saved="fetchPage"
    />
  </NFlex>
</template>

<style scoped></style>
