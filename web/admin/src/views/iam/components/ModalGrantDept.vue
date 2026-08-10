<!--
  Author: Charlie

  分配部门：已选在上、候选在下。
-->
<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import { accountApi, deptApi } from '@/api'
import { renderButtonIcon } from '@/utils'
import { NButton, NTag } from 'naive-ui'
import { computed, reactive } from 'vue'

export type DeptPickItem = { id: string; name: string }

interface FlatDept {
  id: string
  name: string
  depth: number
}

const emit = defineEmits<{
  saved: []
  confirm: [items: DeptPickItem[]]
}>()

const state = reactive({
  showModal: false,
  mode: 'grant' as 'grant' | 'pick',
  multiple: true,
  titleOverride: '',
  loading: false,
  submitLoading: false,
  searchKey: '',
  account: {} as any,
  flatDepts: [] as FlatDept[],
  selectedData: [] as FlatDept[],
  primaryId: null as string | null,
  page: 1,
  pageSize: 20,
})

const isPick = computed(() => state.mode === 'pick')

const modalTitle = computed(() => {
  if (state.titleOverride) return state.titleOverride
  if (isPick.value) return '选择部门'
  return state.account?.name ? `分配部门 - ${state.account.name}` : '分配部门'
})

const filteredDepts = computed(() => {
  const k = state.searchKey.trim().toLowerCase()
  return k ? state.flatDepts.filter((d) => d.name.toLowerCase().includes(k)) : state.flatDepts
})

const tableDepts = computed(() => {
  const s = (state.page - 1) * state.pageSize
  return filteredDepts.value.slice(s, s + state.pageSize)
})

const selectedIds = computed(() => new Set(state.selectedData.map((d) => d.id)))

const listColumns = computed<DataTableColumns<FlatDept>>(() => [
  {
    title: '操作',
    key: 'action',
    align: 'center',
    width: 56,
    render: (row) => (
      <NButton
        text
        type="primary"
        disabled={selectedIds.value.has(row.id)}
        onClick={() => addRecord(row)}
      >
        {renderButtonIcon('icon-park-outline:plus')}
      </NButton>
    ),
  },
  {
    title: '部门名称',
    key: 'name',
    minWidth: 160,
    render: (row) => <span style={{ paddingLeft: `${row.depth * 16}px` }}>{row.name}</span>,
  },
])

const selectedColumns = computed<DataTableColumns<FlatDept>>(() => {
  const cols: DataTableColumns<FlatDept> = [
    {
      title: '操作',
      key: 'action',
      align: 'center',
      width: 70,
      render: (row) => (
        <NButton text type="error" onClick={() => delRecord(row)}>
          {renderButtonIcon('icon-park-outline:delete')}
        </NButton>
      ),
    },
    { title: '部门名称', key: 'name', minWidth: 120 },
  ]
  if (!isPick.value) {
    cols.push({
      title: '主部门',
      key: 'primary',
      width: 110,
      render: (row) =>
        state.primaryId === row.id ? (
          <NTag type="success" bordered={false}>
            主部门
          </NTag>
        ) : (
          <NButton text type="primary" onClick={() => (state.primaryId = row.id)}>
            设为主部门
          </NButton>
        ),
    })
  }
  return cols
})

async function openModal(account: any) {
  state.mode = 'grant'
  state.multiple = true
  state.titleOverride = ''
  state.account = account ?? {}
  state.searchKey = ''
  state.flatDepts = []
  state.selectedData = []
  state.primaryId = null
  state.page = 1
  state.showModal = true
  await fetchGrantData()
}

async function openPicker(options?: {
  selectedIds?: string[]
  multiple?: boolean
  title?: string
}) {
  state.mode = 'pick'
  state.multiple = options?.multiple ?? false
  state.titleOverride = options?.title || ''
  state.account = {}
  state.searchKey = ''
  state.flatDepts = []
  state.selectedData = []
  state.primaryId = null
  state.page = 1
  state.showModal = true
  await fetchPickData(options?.selectedIds ?? [])
}

async function fetchGrantData() {
  if (!state.account?.id) return
  state.loading = true
  try {
    const [deptRes, grantRes] = await Promise.all([
      deptApi.tree().catch(() => ({ data: [] })),
      accountApi.ownDepts(state.account.id),
    ])
    state.flatDepts = flattenTree(deptRes.data ?? [])
    const infoList = grantRes.data?.grant_info_list ?? []
    const selIds = new Set(infoList.map((i: any) => String(i.dept_id)))
    state.selectedData = state.flatDepts.filter((d) => selIds.has(d.id))
    const primary = infoList.find((i: any) => i.is_primary)?.dept_id
    state.primaryId = primary ? String(primary) : (state.selectedData[0]?.id ?? null)
  } finally {
    state.loading = false
  }
}

async function fetchPickData(selectedIds: string[]) {
  state.loading = true
  try {
    const deptRes = await deptApi.tree().catch(() => ({ data: [] }))
    state.flatDepts = flattenTree(deptRes.data ?? [])
    const sel = new Set(selectedIds.map(String))
    state.selectedData = state.flatDepts.filter((d) => sel.has(d.id))
  } finally {
    state.loading = false
  }
}

async function submit() {
  if (isPick.value) {
    emit(
      'confirm',
      state.selectedData.map((d) => ({ id: d.id, name: d.name })),
    )
    closeModal()
    return
  }

  state.submitLoading = true
  try {
    const ids = state.selectedData.map((d) => d.id)
    const primary = state.primaryId && ids.includes(state.primaryId) ? state.primaryId : ids[0]
    await accountApi.grantDepts({
      id: state.account.id,
      grant_info_list: ids.map((id) => ({ dept_id: id, is_primary: id === primary })),
    })
    window.$message.success('授权保存成功')
    closeModal()
    emit('saved')
  } finally {
    state.submitLoading = false
  }
}

function flattenTree(nodes: any[], depth = 0): FlatDept[] {
  const result: FlatDept[] = []
  for (const n of nodes) {
    result.push({ id: String(n.id), name: n.name, depth })
    if (n.children) result.push(...flattenTree(n.children, depth + 1))
  }
  return result
}

function closeModal() {
  state.flatDepts = []
  state.selectedData = []
  state.primaryId = null
  state.showModal = false
  state.submitLoading = false
}

function addRecord(r: FlatDept) {
  if (isPick.value && !state.multiple) {
    state.selectedData = [r]
    return
  }
  if (!selectedIds.value.has(r.id)) state.selectedData.push(r)
}

function addAllPageRecord() {
  if (isPick.value && !state.multiple) return
  tableDepts.value.forEach(addRecord)
}

function delRecord(r: FlatDept) {
  state.selectedData = state.selectedData.filter((d) => d.id !== r.id)
  if (state.primaryId === r.id) state.primaryId = state.selectedData[0]?.id ?? null
}

function delAllRecord() {
  state.selectedData = []
  state.primaryId = null
}

function resetSearch() {
  state.searchKey = ''
  state.page = 1
}

defineExpose({ openModal, openPicker })
</script>

<template>
  <NDrawer
    v-model:show="state.showModal"
    :default-width="760"
    placement="right"
    resizable
    :mask-closable="false"
  >
    <NDrawerContent
      :title="modalTitle"
      closable
      :native-scrollbar="false"
    >
      <div class="grant-pick">
        <section class="grant-pick__panel grant-pick__panel--selected">
          <NFlex
            class="grant-pick__bar"
            justify="space-between"
            align="center"
          >
            <div class="grant-pick__meta">
              <span class="grant-pick__title">已选部门</span>
              <NText depth="3">
                {{ state.selectedData.length }} 个
              </NText>
            </div>
            <NButton
              dashed
              type="error"
              :disabled="!state.selectedData.length"
              @click="delAllRecord"
            >
              全部移除
            </NButton>
          </NFlex>
          <NDataTable
            size="small"
            :row-key="(r: FlatDept) => r.id"
            :columns="selectedColumns"
            :data="state.selectedData"
            :bordered="true"
            :single-line="false"
            max-height="220px"
          />
        </section>

        <section class="grant-pick__panel">
          <NFlex
            class="grant-pick__bar"
            vertical
            :size="10"
          >
            <NInputGroup>
              <NInput
                v-model:value="state.searchKey"
                clearable
                placeholder="搜索部门"
                @keyup.enter="state.page = 1"
                @clear="resetSearch"
              />
              <NButton
                type="primary"
                @click="state.page = 1"
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
              <div class="grant-pick__meta">
                <span class="grant-pick__title">候选部门</span>
                <NText depth="3">
                  共 {{ filteredDepts.length }} 个
                </NText>
              </div>
              <NButton
                v-if="!(isPick && !state.multiple)"
                dashed
                @click="addAllPageRecord"
              >
                新增当前页
              </NButton>
            </NFlex>
          </NFlex>
          <NDataTable
            size="small"
            :row-key="(r: FlatDept) => r.id"
            :columns="listColumns"
            :data="tableDepts"
            :loading="state.loading"
            :bordered="true"
            :single-line="false"
            max-height="calc(100vh - 520px)"
          />
          <NFlex
            class="grant-pick__pager"
            justify="end"
          >
            <NPagination
              v-model:page="state.page"
              v-model:page-size="state.pageSize"
              show-size-picker
              :item-count="filteredDepts.length"
              :page-sizes="[10, 20, 50, 100]"
            />
          </NFlex>
        </section>
      </div>

      <template #footer>
        <NSpace
          justify="end"
          align="center"
        >
          <NButton @click="closeModal">
            关闭
          </NButton>
          <NButton
            type="primary"
            :loading="state.submitLoading"
            @click="submit"
          >
            {{ isPick ? '确认' : '保存' }}
          </NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped>
.grant-pick {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.grant-pick__panel {
  min-width: 0;
}

.grant-pick__panel--selected {
  padding-bottom: 14px;
  border-bottom: 1px solid var(--border-color, #eef2f7);
}

.grant-pick__bar {
  margin-bottom: 10px;
}

.grant-pick__meta {
  display: flex;
  align-items: baseline;
  gap: 8px;
  min-width: 0;
}

.grant-pick__title {
  color: var(--text-color-1, #1f1f1f);
  font-size: 14px;
  font-weight: 600;
}

.grant-pick__pager {
  margin-top: 10px;
}
</style>
