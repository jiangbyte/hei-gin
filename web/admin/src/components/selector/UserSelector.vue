<!-- Author: Charlie -->

<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import { accountApi } from '@/api'
import { renderButtonIcon } from '@/utils'
import { NAvatar, NButton } from 'naive-ui'
import { computed, reactive, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    visible: boolean
    mode?: 'single' | 'multiple'
    title?: string
    selected?: any[]
  }>(),
  {
    mode: 'single',
    title: '选择用户',
    selected: () => [],
  },
)

const emit = defineEmits<{
  'update:visible': [value: boolean]
  select: [account: { id: string; name: string; account: string }]
  'update:selected': [value: any[]]
  confirm: [accounts: any[]]
}>()

const avatarImgProps = { referrerPolicy: 'no-referrer' } as any

const state = reactive({
  loading: false,
  searchKey: '',
  options: [] as any[],
  total: 0,
  selectedData: [] as any[],
  page: 1,
  pageSize: 10,
})

const selectedIds = computed(() => new Set(state.selectedData.map((item) => String(item.id))))

const singleColumns = computed<DataTableColumns<any>>(() => [
  {
    title: '头像',
    key: 'avatar',
    width: 60,
    render: (row) => {
      const url = row.avatar || undefined
      if (url) {
        return <NAvatar round size={36} src={url} imgProps={avatarImgProps} />
      }
      return (
        <NAvatar round size={36}>
          {row.name?.[0] || '?'}
        </NAvatar>
      )
    },
  },
  {
    title: '名称',
    key: 'name',
    minWidth: 120,
    ellipsis: { tooltip: true },
  },
  {
    title: '账号',
    key: 'account',
    width: 120,
    ellipsis: { tooltip: true },
  },
  {
    title: '操作',
    key: 'action',
    width: 60,
    align: 'center',
    render: (row) => (
      <NButton text type="primary" size="small" onClick={() => handleSelect(row)}>
        选择
      </NButton>
    ),
  },
])

const multipleLeftColumns = computed<DataTableColumns<any>>(() => [
  {
    title: '操作',
    key: 'action',
    align: 'center',
    width: 56,
    render: (row) => (
      <NButton
        text
        type="primary"
        size="small"
        disabled={selectedIds.value.has(String(row.id))}
        onClick={() => addRecord(row)}
      >
        {renderButtonIcon('icon-park-outline:plus')}
      </NButton>
    ),
  },
  {
    title: '头像',
    key: 'avatar',
    width: 56,
    render: (row) => {
      const url = row.avatar || undefined
      if (url) {
        return <NAvatar size="small" src={url} imgProps={avatarImgProps} />
      }
      return row.name ? <NAvatar size="small">{row.name?.slice(0, 1)}</NAvatar> : null
    },
  },
  {
    title: '名称',
    key: 'name',
    minWidth: 120,
    ellipsis: { tooltip: true },
  },
  {
    title: '账号',
    key: 'account',
    minWidth: 120,
    ellipsis: { tooltip: true },
  },
])

const multipleRightColumns = computed<DataTableColumns<any>>(() => [
  {
    title: '操作',
    key: 'action',
    align: 'center',
    width: 70,
    render: (row) => (
      <NButton text type="error" size="small" onClick={() => delRecord(row)}>
        {renderButtonIcon('icon-park-outline:delete')}
      </NButton>
    ),
  },
  {
    title: '名称',
    key: 'name',
    minWidth: 120,
    ellipsis: { tooltip: true },
  },
])

watch(
  () => props.visible,
  (val) => {
    if (val) {
      state.selectedData = [...props.selected]
      state.searchKey = ''
      state.page = 1
      loadOptions()
    } else {
      state.selectedData = []
    }
  },
)

watch(
  () => state.page,
  () => {
    if (props.visible) loadOptions()
  },
)

watch(
  () => state.pageSize,
  () => {
    if (props.visible) {
      state.page = 1
      loadOptions()
    }
  },
)

async function loadOptions() {
  state.loading = true
  try {
    const params: any = { current: state.page, size: state.pageSize }
    const keyword = state.searchKey.trim()
    if (keyword) {
      params.name = keyword
      params.account = keyword
    }
    const res = await accountApi.page(params)
    const records = res?.data?.records ?? []
    state.options = records.map((item: any) => ({
      id: item.id,
      name: item.name || item.account || '',
      account: item.account || '',
      avatar: item.avatar || null,
    }))
    state.total = res?.data?.total ?? 0
  } catch {
    state.options = []
    state.total = 0
  } finally {
    state.loading = false
  }
}

function handleSelect(account: any) {
  const result = { id: account.id, name: account.name, account: account.account }
  emit('select', result)
  close()
}

function doSearch() {
  state.page = 1
  loadOptions()
}

function addRecord(record: any) {
  if (!selectedIds.value.has(String(record.id))) {
    state.selectedData.push(record)
  }
}

function addAllPage() {
  state.options.forEach(addRecord)
}

function delRecord(record: any) {
  state.selectedData = state.selectedData.filter((item) => String(item.id) !== String(record.id))
}

function delAll() {
  state.selectedData = []
}

function handleConfirm() {
  emit('update:selected', [...state.selectedData])
  emit('confirm', [...state.selectedData])
  close()
}

function close() {
  emit('update:visible', false)
}

function resetSearch() {
  state.searchKey = ''
  state.page = 1
  loadOptions()
}
</script>

<template>
  <NDrawer
    :show="visible"
    :default-width="mode === 'multiple' ? 1000 : 500"
    placement="right"
    :mask-closable="false"
    @update:show="(val: boolean) => emit('update:visible', val)"
  >
    <NDrawerContent
      :title="title"
      closable
      :native-scrollbar="false"
    >
      <!-- 单选模式 -->
      <template v-if="mode === 'single'">
        <NSpace vertical>
          <NInput
            v-model:value="state.searchKey"
            clearable
            :placeholder="'搜索用户名/账号'"
            @keyup.enter="doSearch"
            @clear="resetSearch"
          />
          <NDataTable
            :row-key="(row: any) => row.id"
            :columns="singleColumns"
            :data="state.options"
            :loading="state.loading"
            :bordered="true"
            :single-line="false"
            max-height="calc(100vh - 290px)"
          />
          <NPagination
            v-model:page="state.page"
            v-model:page-size="state.pageSize"
            show-size-picker
            :item-count="state.total"
            :page-sizes="[10, 20, 50, 100]"
          />
        </NSpace>
      </template>

      <!-- 多选模式 -->
      <template v-else>
        <NGrid
          :cols="24"
          :x-gap="10"
        >
          <NGi :span="16">
            <NSpace vertical>
              <NInputGroup>
                <NInput
                  v-model:value="state.searchKey"
                  clearable
                  :placeholder="'请输入用户名'"
                  @keyup.enter="doSearch"
                  @clear="resetSearch"
                />
                <NButton
                  type="primary"
                  @click="doSearch"
                >
                  搜索
                </NButton>
                <NButton @click="resetSearch">
                  重置
                </NButton>
              </NInputGroup>
              <NFlex
                justify="space-between"
                align="center"
              >
                <NText>{{ `待处理列表: ${state.total}` }}</NText>
                <NButton
                  dashed
                  size="small"
                  @click="addAllPage"
                >
                  新增当前页
                </NButton>
              </NFlex>
              <NDataTable
                :row-key="(row: any) => row.id"
                :columns="multipleLeftColumns"
                :data="state.options"
                :loading="state.loading"
                :bordered="true"
                :single-line="false"
                max-height="calc(100vh - 340px)"
              />
              <NPagination
                v-model:page="state.page"
                v-model:page-size="state.pageSize"
                show-size-picker
                :item-count="state.total"
                :page-sizes="[10, 20, 50, 100]"
              />
            </NSpace>
          </NGi>
          <NGi :span="8">
            <NSpace vertical>
              <NFlex
                justify="space-between"
                align="center"
              >
                <NText>已选择：{{ state.selectedData.length }}</NText>
                <NButton
                  dashed
                  type="error"
                  size="small"
                  @click="delAll"
                >
                  全部移除
                </NButton>
              </NFlex>
              <NDataTable
                :row-key="(row: any) => row.id"
                :columns="multipleRightColumns"
                :data="state.selectedData"
                :bordered="true"
                :single-line="false"
                max-height="calc(100vh - 280px)"
              />
            </NSpace>
          </NGi>
        </NGrid>
      </template>

      <template #footer>
        <NSpace
          justify="end"
          align="center"
        >
          <NButton @click="close">
            关闭
          </NButton>
          <NButton
            v-if="mode === 'multiple'"
            type="primary"
            @click="handleConfirm"
          >
            确认
          </NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>
