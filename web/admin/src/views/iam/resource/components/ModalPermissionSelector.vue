<!-- Author: Charlie -->

<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import { resourceApi } from '@/api'
import { NButton, NInput, NInputGroup, NTag } from 'naive-ui'
import { computed, reactive } from 'vue'

const emit = defineEmits<{
  selected: [permission: any]
}>()

const state = reactive({
  showModal: false,
  loading: false,
  searchText: '',
  currentKey: '',
  rows: [] as any[],
  page: 1,
  pageSize: 10,
})

const filteredRows = computed(() => {
  const keyword = state.searchText.trim().toUpperCase()
  if (!keyword) {
    return state.rows
  }
  return state.rows.filter((item) => item.api.toUpperCase().includes(keyword))
})
const tableRows = computed(() => {
  const start = (state.page - 1) * state.pageSize
  return filteredRows.value.slice(start, start + state.pageSize)
})
const firstShowMap = computed<Record<string, number[]>>(() => {
  const map: Record<string, number[]> = {}
  tableRows.value.forEach((item: any, index: number) => {
    if (map[item.prefix]) {
      map[item.prefix].push(index)
    } else {
      map[item.prefix] = [index]
    }
  })
  return map
})

const columns = computed<DataTableColumns<any>>(() => [
  {
    title: '权限分组',
    key: 'prefix',
    fixed: 'left',
    width: 220,
    rowSpan: (row, rowIndex) => {
      const indexArr = firstShowMap.value[row.prefix] ?? []
      return rowIndex === indexArr[0] ? indexArr.length : 0
    },
    render: (row) => row.prefix,
  },
  {
    title: '权限标识',
    key: 'suffix',
    minWidth: 360,
    filter: true,
    renderFilterMenu: ({ hide }) => (
      <div class="permission-selector-filter">
        <NInput
          value={state.searchText}
          placeholder={'输入权限标识或名称'}
          onUpdateValue={(value) => {
            state.searchText = value
          }}
          onKeyup={(event: KeyboardEvent) => {
            if (event.key === 'Enter') {
              handleSearch()
              hide()
            }
          }}
        />
        <NInputGroup>
          <NButton
            type="primary"
            size="small"
            onClick={() => {
              handleSearch()
              hide()
            }}
          >
            {'搜索'}
          </NButton>
          <NButton
            size="small"
            onClick={() => {
              resetSearch()
              hide()
            }}
          >
            {'重置'}
          </NButton>
        </NInputGroup>
      </div>
    ),
    render: (row) => (
      <div class="permission-selector-key">
        <span>{row.suffix}</span>
        {row.permission_key === state.currentKey ? (
          <NTag size="small" type="success" bordered={false}>
            {'已选择'}
          </NTag>
        ) : null}
      </div>
    ),
  },
  {
    title: '权限名称',
    key: 'name',
    minWidth: 200,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    fixed: 'right',
    render: (row) => (
      <NButton type="primary" size="small" text={true} onClick={() => selectPermission(row)}>
        {'选择权限'}
      </NButton>
    ),
  },
])

async function openModal(currentKey = '') {
  state.currentKey = currentKey
  state.searchText = ''
  state.page = 1
  state.showModal = true
  if (!state.rows.length) {
    await fetchPermissions()
  }
}

async function fetchPermissions() {
  state.loading = true
  try {
    const response = await resourceApi.permissionRegistry()
    state.rows = (response.data ?? []).map((item: any) => {
      const permissionKey = item.permission_key ?? item
      const [prefix, suffix] = splitByPermissionKey(permissionKey)
      return {
        ...item,
        permission_key: permissionKey,
        name: item.name ?? permissionKey,
        api: `${permissionKey}[${item.name ?? permissionKey}]`,
        prefix,
        suffix,
      }
    })
  } finally {
    state.loading = false
  }
}

function closeModal() {
  state.showModal = false
}

function selectPermission(row: any) {
  emit('selected', row)
  closeModal()
}

function handleSearch() {
  state.page = 1
}

function resetSearch() {
  state.searchText = ''
  handleSearch()
}

function splitByPermissionKey(permissionKey: string) {
  const parts = permissionKey.split(':')
  return [parts.slice(0, 2).join(':') || permissionKey, parts.slice(2).join(':') || permissionKey]
}

defineExpose({
  openModal,
})
</script>

<template>
  <NModal
    v-model:show="state.showModal"
    preset="card"
    draggable
    :mask-closable="false"
    :title="'选择权限'"
    style="width: 980px"
    :segmented="{ content: true, action: true }"
  >
    <NSpin :show="state.loading">
      <NDataTable
        size="medium"
        :row-key="(row: any) => row.permission_key"
        :columns="columns"
        :data="tableRows"
        :bordered="true"
        :single-line="false"
        :scroll-x="900"
        max-height="calc(100vh - 340px)"
      />
      <NFlex
        justify="end"
        class="mt-10px"
      >
        <NPagination
          v-model:page="state.page"
          v-model:page-size="state.pageSize"
          show-size-picker
          size="small"
          :item-count="filteredRows.length"
          :page-sizes="[10, 20, 50, 100]"
        />
      </NFlex>
    </NSpin>

    <template #action>
      <NSpace
        justify="end"
        align="center"
      >
        <NButton @click="closeModal">
          关闭
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped>
.permission-selector-filter {
  width: 300px;
  padding: 12px;
}

.permission-selector-key {
  display: flex;
  gap: 8px;
  align-items: center;
}
</style>
