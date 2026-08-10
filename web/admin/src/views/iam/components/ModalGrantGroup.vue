<!--
  Author: Charlie

  分配用户组：已选在上、候选在下。
-->
<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import { accountApi, groupApi } from '@/api'
import { createTagColor, renderButtonIcon } from '@/utils'
import { dictTypeColor, dictTypeData } from '@/utils/dict'
import { NButton, NTag } from 'naive-ui'
import { computed, reactive } from 'vue'

const emit = defineEmits<{
  saved: []
}>()

const state = reactive({
  showModal: false,
  loading: false,
  submitLoading: false,
  searchKey: '',
  account: {} as any,
  items: [] as any[],
  selectedData: [] as any[],
  page: 1,
  pageSize: 10,
  total: 0,
})

const modalTitle = computed(() =>
  state.account?.name ? `分配用户组 - ${state.account.name}` : '分配用户组',
)

const selectedIds = computed(() => new Set(state.selectedData.map((item) => String(item.id))))

const listColumns = computed<DataTableColumns<any>>(() => [
  {
    title: '操作',
    key: 'action',
    align: 'center',
    width: 56,
    render: (row) => (
      <NButton
        text
        type="primary"
        disabled={selectedIds.value.has(String(row.id))}
        onClick={() => addRecord(row)}
      >
        {renderButtonIcon('icon-park-outline:plus')}
      </NButton>
    ),
  },
  {
    title: '名称',
    key: 'name',
    minWidth: 160,
    ellipsis: { tooltip: true },
  },
  {
    title: '状态',
    key: 'status',
    width: 110,
    render: (row) => (
      <NTag color={createTagColor(dictTypeColor('COMMON_STATUS', row.status))} bordered={false}>
        {dictTypeData('COMMON_STATUS', row.status) || row.status}
      </NTag>
    ),
  },
])

const selectedColumns = computed<DataTableColumns<any>>(() => [
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
  {
    title: '名称',
    key: 'name',
    minWidth: 120,
    ellipsis: { tooltip: true },
  },
])

async function openModal(account: any) {
  state.account = account ?? {}
  state.searchKey = ''
  state.items = []
  state.selectedData = []
  state.page = 1
  state.total = 0
  state.showModal = true
  await Promise.all([fetchGrantedSelected(), fetchCandidatePage()])
}

async function fetchGrantedSelected() {
  if (!state.account?.id) return
  const response = await accountApi.ownGroups(state.account.id)
  state.selectedData = response.data?.groups ?? []
}

async function fetchCandidatePage() {
  state.loading = true
  try {
    const keyword = state.searchKey.trim()
    const response = await groupApi.page({
      current: state.page,
      size: state.pageSize,
      ...(keyword ? { name: keyword } : {}),
    })
    state.items = response.data?.records ?? []
    state.total = Number(response.data?.total ?? 0)
  } finally {
    state.loading = false
  }
}

async function submitGrant() {
  state.submitLoading = true
  try {
    await accountApi.grantGroups({
      id: state.account.id,
      group_ids: state.selectedData.map((item) => item.id),
    })
    window.$message.success('授权保存成功')
    closeModal()
    emit('saved')
  } finally {
    state.submitLoading = false
  }
}

function closeModal() {
  state.items = []
  state.selectedData = []
  state.total = 0
  state.showModal = false
  state.submitLoading = false
}

function addRecord(record: any) {
  if (!selectedIds.value.has(String(record.id))) {
    state.selectedData.push(record)
  }
}

function addAllPageRecord() {
  state.items.forEach(addRecord)
}

function delRecord(record: any) {
  state.selectedData = state.selectedData.filter((item) => String(item.id) !== String(record.id))
}

function delAllRecord() {
  state.selectedData = []
}

function doSearch() {
  state.page = 1
  void fetchCandidatePage()
}

function resetSearch() {
  state.searchKey = ''
  state.page = 1
  void fetchCandidatePage()
}

function onPageChange(page: number) {
  state.page = page
  void fetchCandidatePage()
}

function onPageSizeChange(size: number) {
  state.pageSize = size
  state.page = 1
  void fetchCandidatePage()
}

defineExpose({ openModal })
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
              <span class="grant-pick__title">已选用户组</span>
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
            :row-key="(row: any) => row.id"
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
                placeholder="请输入用户组名称"
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
              <div class="grant-pick__meta">
                <span class="grant-pick__title">候选用户组</span>
                <NText depth="3">
                  共 {{ state.total }} 个
                </NText>
              </div>
              <NButton
                dashed
                @click="addAllPageRecord"
              >
                新增当前页
              </NButton>
            </NFlex>
          </NFlex>
          <NDataTable
            size="small"
            :row-key="(row: any) => row.id"
            :columns="listColumns"
            :data="state.items"
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
              :page="state.page"
              :page-size="state.pageSize"
              show-size-picker
              :item-count="state.total"
              :page-sizes="[10, 20, 50, 100]"
              @update:page="onPageChange"
              @update:page-size="onPageSizeChange"
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
            @click="submitGrant"
          >
            保存
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
