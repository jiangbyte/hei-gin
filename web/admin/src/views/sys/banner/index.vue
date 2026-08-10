<!-- Author: Charlie -->

<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { NButton, NFlex, NIcon, NImage, NTag } from 'naive-ui'
import { bannerApi } from '@/api'
import { ACCOUNT_TYPE_OPTIONS, accountTypeLabel } from '@/constants/account'
import {
  createTagColor,
  formatDateTime,
  hasPermission,
  normalizeSearchValues,
  renderButtonIcon,
} from '@/utils'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { dictList, dictTypeData, dictTypeColor } from '@/utils/dict'
import { readPageMeta } from '@/utils/wire'

const router = useRouter()
const state = reactive({
  banners: [] as any[],
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
    title: '目标账户类型',
    path: 'target_account_type',
    field: 'select',
    fieldProps: {
      options: ACCOUNT_TYPE_OPTIONS,
    },
  },
  {
    title: '分类',
    path: 'category',
    field: 'select',
    fieldProps: {
      options: dictList('BANNER_CATEGORY'),
    },
  },
  {
    title: '类型',
    path: 'type',
    field: 'select',
    fieldProps: {
      options: dictList('BANNER_TYPE'),
    },
  },
  {
    title: '位置',
    path: 'position',
    field: 'select',
    fieldProps: {
      options: dictList('BANNER_POSITION'),
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
    title: '标题',
    path: 'title',
    width: 180,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '图片',
    key: 'image',
    width: 90,
    render: (row) => renderImage(row),
  },
  {
    title: '目标账户类型',
    path: 'target_account_types',
    width: 160,
    render: (row) => {
      const types = Array.isArray(row.target_account_types) ? row.target_account_types : []
      if (!types.length) return <span />
      return (
        <NFlex size={4} wrap>
          {types.map((t: string) => (
            <NTag key={t} size="small" bordered={false} type="info">
              {accountTypeLabel(t)}
            </NTag>
          ))}
        </NFlex>
      )
    },
  },
  {
    title: '分类',
    path: 'category',
    width: 150,
    render: (row) => (
      <NTag
        size="small"
        color={createTagColor(dictTypeColor('BANNER_CATEGORY', row.category))}
        bordered={false}
      >
        {dictTypeData('BANNER_CATEGORY', row.category) || row.category}
      </NTag>
    ),
  },
  {
    title: '类型',
    path: 'type',
    width: 120,
    render: (row) => (
      <NTag
        size="small"
        color={createTagColor(dictTypeColor('BANNER_TYPE', row.type))}
        bordered={false}
      >
        {dictTypeData('BANNER_TYPE', row.type) || row.type}
      </NTag>
    ),
  },
  {
    title: '位置',
    path: 'position',
    width: 160,
    render: (row) => (
      <NTag
        size="small"
        color={createTagColor(dictTypeColor('BANNER_POSITION', row.position))}
        bordered={false}
      >
        {dictTypeData('BANNER_POSITION', row.position) || row.position}
      </NTag>
    ),
  },
  {
    title: '链接类型',
    path: 'link_type',
    width: 110,
    render: (row) => (
      <NTag
        size="small"
        color={createTagColor(dictTypeColor('BANNER_LINK_TYPE', row.link_type))}
        bordered={false}
      >
        {dictTypeData('BANNER_LINK_TYPE', row.link_type) || row.link_type}
      </NTag>
    ),
  },
  {
    title: '排序',
    path: 'sort',
    width: 90,
  },
  {
    title: '互动次数',
    path: 'interaction_count',
    width: 120,
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
    title: '开始时间',
    path: 'start_at',
    width: 190,
    ellipsis: { tooltip: true },
    render: (row) => formatDateTime(row.start_at),
  },
  {
    title: '结束时间',
    path: 'end_at',
    width: 190,
    ellipsis: { tooltip: true },
    render: (row) => formatDateTime(row.end_at),
  },
  {
    title: '更新时间',
    path: 'updated_at',
    width: 190,
    ellipsis: { tooltip: true },
    render: (row) => formatDateTime(row.updated_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    fixed: 'right',
    render: (row) => (
      <NFlex size={12}>
        {hasPermission('sys:banner:detail') ? (
          <NButton type="info" size="small" text={true} onClick={() => openDetail(row.id)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('sys:banner:update') ? (
          <NButton type="primary" size="small" text={true} onClick={() => openEdit(row.id)}>
            {renderButtonIcon('icon-park-outline:edit')}
          </NButton>
        ) : null}
        {hasPermission('sys:banner:delete') ? (
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
    const response = await bannerApi.page({
      current: state.page,
      size: state.pageSize,
      ...state.searchValues,
    })
    const data = response.data ?? {}
    state.banners = data.records ?? []
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
  } finally {
    state.loading = false
  }
}

function renderImage(row: any) {
  const src = row.image_url || row.image || undefined
  if (!src) {
    return <span>-</span>
  }
  return <NImage src={src} alt={row.title || '图片'} width={72} height={48} objectFit="cover" />
}

function openDetail(id: string) {
  router.push({ path: '/sys/banner/detail', query: { id } })
}

function openCreate() {
  router.push('/sys/banner/create')
}

function openEdit(id: string) {
  router.push({ path: '/sys/banner/edit', query: { id } })
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
    content: isBatch ? `删除 ${ids.length} 张展示图?` : '删除该展示图?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: () => deleteData(ids),
  })
}

async function deleteData(ids: string[]) {
  await bannerApi.remove({ ids })
  state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))

  window.$message.success('删除成功')
  await fetchPage()
  if (!state.banners.length && state.total > 0 && state.page > 1) {
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
      :title="'展示图管理'"
      row-key="id"
      :scroll-x="1780"
      :columns="tableColumns"
      :data="state.banners"
      :loading="state.loading"
      :pagination="pagination"
      :checked-row-keys="state.checkedRowKeys"
      :on-update-checked-row-keys="handleCheckedRowKeys"
    >
      <template #toolbar>
        <NFlex>
          <NButton
            v-if="hasPermission('sys:banner:create')"
            type="primary"
            text
            :title="'新增'"
            :aria-label="'新增'"
            @click="openCreate"
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
            v-if="hasPermission('sys:banner:delete')"
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
</template>
