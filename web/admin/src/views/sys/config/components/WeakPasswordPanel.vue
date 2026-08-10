<!-- Author: Charlie -->

<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { weakPasswordApi } from '@/api'
import { formatDateTime } from '@/utils'
import { computed, h, onMounted, reactive } from 'vue'
import { loadByCategory, saveByKeys } from '../composables/useConfigForm'

const CATEGORY = 'AUTH_PASSWORD'

const state = reactive({
  loading: false,
  savingWords: false,
  customWeakWords: '',
  snapshotWords: '',
  rows: [] as any[],
  total: 0,
  page: 1,
  pageSize: 10,
  checkedRowKeys: [] as string[],
  modalShow: false,
  modalSaving: false,
  editingId: '' as string,
  formPassword: '',
})

const columns = computed<DataTableColumns<any>>(() => [
  { type: 'selection' },
  { title: '弱密码', key: 'password', ellipsis: { tooltip: true } },
  {
    title: '创建时间',
    key: 'created_at',
    width: 180,
    render: (row) => formatDateTime(row.created_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 140,
    render: (row) =>
      h('div', { class: 'flex gap-12px' }, [
        h(
          'a',
          {
            class: 'text-primary cursor-pointer',
            onClick: () => openEdit(row),
          },
          '编辑',
        ),
        h(
          'a',
          {
            class: 'text-error cursor-pointer',
            onClick: () => void removeOne(row.id),
          },
          '删除',
        ),
      ]),
  },
])

onMounted(() => {
  void reloadWords()
  void fetchPage()
})

async function reloadWords() {
  const map = await loadByCategory(CATEGORY)
  state.customWeakWords = map.PASSWORD_CUSTOM_WEAK_WORDS || ''
  state.snapshotWords = state.customWeakWords
}

async function fetchPage() {
  state.loading = true
  try {
    const res = await weakPasswordApi.page({
      page: state.page,
      size: state.pageSize,
    })
    state.rows = res.data?.records ?? res.data?.items ?? []
    state.total = res.data?.total ?? state.rows.length
  } finally {
    state.loading = false
  }
}

async function saveWords() {
  state.savingWords = true
  try {
    await saveByKeys([
      {
        config_key: 'PASSWORD_CUSTOM_WEAK_WORDS',
        config_value: state.customWeakWords,
        category: CATEGORY,
      },
    ])
    state.snapshotWords = state.customWeakWords
    window.$message.success('自定义弱密码已保存')
  } finally {
    state.savingWords = false
  }
}

function openCreate() {
  state.editingId = ''
  state.formPassword = ''
  state.modalShow = true
}

function openEdit(row: any) {
  state.editingId = row.id
  state.formPassword = row.password
  state.modalShow = true
}

async function submitModal() {
  const password = state.formPassword.trim()
  if (!password) {
    window.$message.warning('请输入弱密码')
    return
  }
  state.modalSaving = true
  try {
    if (state.editingId) {
      await weakPasswordApi.update({ id: state.editingId, password })
    } else {
      await weakPasswordApi.create({ password })
    }
    state.modalShow = false
    window.$message.success('保存成功')
    await fetchPage()
  } finally {
    state.modalSaving = false
  }
}

async function removeOne(id: string) {
  await weakPasswordApi.remove({ ids: [id] })
  window.$message.success('已删除')
  await fetchPage()
}

async function removeBatch() {
  if (!state.checkedRowKeys.length) {
    window.$message.warning('请先选择')
    return
  }
  await weakPasswordApi.remove({ ids: state.checkedRowKeys })
  state.checkedRowKeys = []
  window.$message.success('已删除')
  await fetchPage()
}
</script>

<template>
  <div class="weak-password-panel">
    <NForm
      class="sys-config-form mb-20px"
      label-placement="top"
    >
      <NFormItem label="自定义额外弱密码（逗号分隔）">
        <NInput
          v-model:value="state.customWeakWords"
          type="textarea"
          :autosize="{ minRows: 2, maxRows: 6 }"
          placeholder="password,admin123"
        />
      </NFormItem>
      <NSpace>
        <NButton
          type="primary"
          :loading="state.savingWords"
          @click="saveWords"
        >
          保存词表
        </NButton>
        <NButton @click="state.customWeakWords = state.snapshotWords">
          重置
        </NButton>
      </NSpace>
    </NForm>

    <NSpace class="mb-12px">
      <NButton
        type="primary"
        @click="openCreate"
      >
        新增
      </NButton>
      <NButton
        type="error"
        ghost
        @click="removeBatch"
      >
        批量删除
      </NButton>
      <NButton
        quaternary
        @click="fetchPage"
      >
        刷新
      </NButton>
    </NSpace>

    <NDataTable
      v-model:checked-row-keys="state.checkedRowKeys"
      :loading="state.loading"
      :columns="columns"
      :data="state.rows"
      :row-key="(row: any) => row.id"
      :pagination="{
        page: state.page,
        pageSize: state.pageSize,
        itemCount: state.total,
        showSizePicker: true,
        pageSizes: [10, 20, 50],
        onUpdatePage: (p: number) => {
          state.page = p
          fetchPage()
        },
        onUpdatePageSize: (s: number) => {
          state.pageSize = s
          state.page = 1
          fetchPage()
        },
      }"
    />

    <NModal
      v-model:show="state.modalShow"
      preset="card"
      :title="state.editingId ? '编辑弱密码' : '新增弱密码'"
      class="max-w-120"
      :mask-closable="false"
    >
      <NForm label-placement="top">
        <NFormItem
          label="弱密码"
          required
        >
          <NInput
            v-model:value="state.formPassword"
            placeholder="例如 admin"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="state.modalShow = false">
            取消
          </NButton>
          <NButton
            type="primary"
            :loading="state.modalSaving"
            @click="submitModal"
          >
            保存
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
.mb-20px {
  margin-bottom: 20px;
}
.mb-12px {
  margin-bottom: 12px;
}
</style>
