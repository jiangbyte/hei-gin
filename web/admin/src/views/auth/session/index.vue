<!--
  Author: Charlie

  在线会话管理：台账式总览 + 分端列表。
-->
<script setup lang="tsx">
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { sessionApi } from '@/api'
import {
  ACCOUNT_TYPE_TABS,
  DEFAULT_ACCOUNT_TYPE,
  accountTypeLabel,
  type AccountType,
} from '@/constants/account'
import { formatDateTime, hasPermission, normalizeSearchValues, renderButtonIcon } from '@/utils'
import { NButton, NDataTable, NFlex, NIcon } from 'naive-ui'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, onMounted, reactive } from 'vue'
import { readPageMeta } from '@/utils/wire'

const state = reactive({
  rows: [] as any[],
  tokens: [] as any[],
  analysis: {} as any,
  total: 0,
  loading: false,
  tokenModalShow: false,
  searchValues: {} as any,
  accountType: DEFAULT_ACCOUNT_TYPE as AccountType,
  page: 1,
  pageSize: 20,
})

const searchForm = createProSearchForm<any>({
  defaultCollapsed: true,
  onSubmit(values) {
    state.searchValues = normalizeSearchValues(values, {
      account: (value) => String(value).trim(),
      ip: (value) => String(value).trim(),
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

const analysisCards = computed(() => [
  {
    key: 'online_account_count',
    label: '在线账号',
    note: '当前有会话的账号',
  },
  {
    key: 'online_token_count',
    label: '在线设备',
    note: '全部有效令牌',
  },
  {
    key: 'admin_account_count',
    label: '管理端',
    note: '管理员在线账号',
  },
  {
    key: 'portal_account_count',
    label: '门户端',
    note: '门户用户在线账号',
  },
  {
    key: 'one_hour_new_count',
    label: '近 1 小时',
    note: '新增登录次数',
  },
  {
    key: 'max_token_count',
    label: '单账号峰值',
    note: '最大同时设备数',
    danger: true,
  },
])

const searchColumns = computed<ProSearchFormColumns<any>>(() => [
  { title: '账号', path: 'account', field: 'input' },
  { title: '客户端 IP', path: 'ip', field: 'input' },
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
  { title: '账号 ID', path: 'account_id', width: 170, ellipsis: { tooltip: true } },
  { title: '账号', path: 'account', width: 160, ellipsis: { tooltip: true } },
  { title: '名称', path: 'name', width: 160, ellipsis: { tooltip: true } },
  { title: '设备数', path: 'token_count', width: 110 },
  {
    title: '客户端 IP',
    key: 'client_ip',
    width: 150,
    render: (row) => row.tokens?.[0]?.client_ip || row.latest_login_ip || '-',
  },
  {
    title: '设备',
    key: 'device',
    width: 140,
    render: (row) => row.tokens?.[0]?.device_label || '-',
  },
  {
    title: '最近活跃时间',
    path: 'latest_active_at',
    width: 190,
    ellipsis: { tooltip: true },
    render: (row) => formatDateTime(row.latest_active_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 130,
    fixed: 'right',
    render: (row) => (
      <NFlex size={12}>
        {hasPermission('auth:session:tokenlist') ? (
          <NButton type="info" size="small" text={true} onClick={() => openTokens(row)}>
            {renderButtonIcon('icon-park-outline:preview-open')}
          </NButton>
        ) : null}
        {hasPermission('auth:session:exit') ? (
          <NButton type="error" size="small" text={true} onClick={() => confirmExitAccount(row)}>
            {renderButtonIcon('icon-park-outline:logout')}
          </NButton>
        ) : null}
      </NFlex>
    ),
  },
])

const tokenColumns = computed<DataTableColumns<any>>(() => [
  { title: '令牌', key: 'token', width: 220, ellipsis: { tooltip: true } },
  { title: '设备', key: 'device_label', width: 110 },
  { title: '客户端 IP', key: 'client_ip', width: 140 },
  {
    title: '登录时间',
    key: 'login_at',
    width: 180,
    ellipsis: { tooltip: true },
    render: (row) => formatDateTime(row.login_at),
  },
  {
    title: '上次活跃时间',
    key: 'last_active_at',
    width: 180,
    ellipsis: { tooltip: true },
    render: (row) => formatDateTime(row.last_active_at),
  },
  {
    title: '过期时间',
    key: 'expires_at',
    width: 180,
    ellipsis: { tooltip: true },
    render: (row) => formatDateTime(row.expires_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 90,
    fixed: 'right',
    render: (row) =>
      hasPermission('auth:session:tokenexit') ? (
        <NButton type="error" size="small" text={true} onClick={() => confirmExitToken(row.token)}>
          {renderButtonIcon('icon-park-outline:logout')}
        </NButton>
      ) : null,
  },
])

const tableTitle = computed(() => {
  const label = accountTypeLabel(state.accountType) || state.accountType
  return label ? `${label}会话` : '在线会话'
})

onMounted(() => {
  void fetchAll()
})

function handleAccountTypeChange(value: string | number) {
  state.accountType = String(value) as AccountType
  state.page = 1
  void fetchPage()
}

async function fetchAll() {
  await Promise.all([fetchAnalysis(), fetchPage()])
}

async function fetchAnalysis() {
  const response = await sessionApi.analysis()
  state.analysis = response.data ?? {}
}

async function fetchPage() {
  if (!state.accountType) return
  state.loading = true
  try {
    const response = await sessionApi.page({
      current: state.page,
      size: state.pageSize,
      account_type: state.accountType,
      ...state.searchValues,
    })
    const data = response.data ?? {}
    state.rows = data.records ?? []
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
  } finally {
    state.loading = false
  }
}

function openTokens(row: any) {
  state.tokens = row.tokens ?? []
  state.tokenModalShow = true
}

function confirmExitAccount(row: any) {
  window.$dialog.warning({
    title: '强制下线',
    draggable: true,
    maskClosable: false,
    content: '强制下线该账号的所有在线设备?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: async () => {
      await sessionApi.exit({
        targets: [{ account_type: row.account_type, account_id: row.account_id }],
      })
      window.$message.success('强制下线成功')
      await fetchAll()
    },
  })
}

function confirmExitToken(token: string) {
  window.$dialog.warning({
    title: '强制下线',
    draggable: true,
    maskClosable: false,
    content: '强制下线该在线设备?',
    positiveText: '确认',
    negativeText: '取消',
    onPositiveClick: async () => {
      await sessionApi.tokenExit({ tokens: [token] })
      window.$message.success('强制下线成功')
      state.tokens = state.tokens.filter((item) => item.token !== token)
      await fetchAll()
    },
  })
}
</script>

<template>
  <NFlex
    class="session-page h-full min-h-0"
    vertical
    :size="12"
  >
    <header class="session-head">
      <div class="session-head__copy">
        <div class="session-head__title">
          在线会话
        </div>
        <div class="session-head__sub">
          查看各端在线账号与设备，必要时可强制下线
        </div>
      </div>
      <NButton
        text
        :loading="state.loading"
        @click="fetchAll"
      >
        <template #icon>
          <NIcon>
            <Icon icon="icon-park-outline:reload" />
          </NIcon>
        </template>
        刷新总览
      </NButton>
    </header>

    <section class="session-ledger">
      <div
        v-for="item in analysisCards"
        :key="item.key"
        class="session-ledger__cell"
        :class="{ 'session-ledger__cell--danger': item.danger }"
      >
        <span class="session-ledger__label">{{ item.label }}</span>
        <span class="session-ledger__value">{{ state.analysis[item.key] ?? 0 }}</span>
        <span class="session-ledger__note">{{ item.note }}</span>
      </div>
    </section>

    <NTabs
      class="session-account-tabs"
      :value="state.accountType"
      type="line"
      animated
      @update:value="handleAccountTypeChange"
    >
      <NTabPane
        v-for="item in ACCOUNT_TYPE_TABS"
        :key="item.key"
        :name="item.key"
        :tab="item.label"
      />
    </NTabs>

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
      :title="tableTitle"
      row-key="account_id"
      :scroll-x="1210"
      :columns="tableColumns"
      :data="state.rows"
      :loading="state.loading"
      :pagination="pagination"
    >
      <template #toolbar>
        <NButton
          text
          :title="'刷新'"
          :aria-label="'刷新'"
          :loading="state.loading"
          @click="fetchAll"
        >
          <template #icon>
            <NIcon>
              <Icon icon="icon-park-outline:reload" />
            </NIcon>
          </template>
        </NButton>
      </template>
    </ProDataTable>

    <NModal
      v-model:show="state.tokenModalShow"
      preset="card"
      draggable
      :title="'设备详情'"
      style="width: min(960px, calc(100vw - 32px))"
    >
      <NScrollbar style="max-height: min(540px, 70vh)">
        <NDataTable
          :row-key="(row) => row.token"
          :scroll-x="1170"
          :columns="tokenColumns"
          :data="state.tokens"
          :pagination="false"
        />
      </NScrollbar>
    </NModal>
  </NFlex>
</template>

<style scoped>
.session-page {
  min-width: 0;
}

.session-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 14px;
  background: var(--card-color, #fff);
  border: 1px solid var(--border-color, #eef2f7);
}

.session-head__title {
  color: var(--text-color-1, #1f1f1f);
  font-size: 17px;
  font-weight: 650;
  line-height: 1.3;
}

.session-head__sub {
  margin-top: 4px;
  color: var(--text-color-3, #999);
  font-size: 12px;
  line-height: 1.4;
}

.session-ledger {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  background: var(--card-color, #fff);
  border: 1px solid var(--border-color, #eef2f7);
}

.session-ledger__cell {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  min-width: 0;
  padding: 12px 10px;
  border-right: 1px solid var(--border-color, #eef2f7);
}

.session-ledger__cell:last-child {
  border-right: 0;
}

.session-ledger__label {
  color: var(--text-color-3, #999);
  font-size: 12px;
}

.session-ledger__value {
  color: var(--text-color-1, #1f1f1f);
  font-size: 22px;
  font-weight: 720;
  font-variant-numeric: tabular-nums;
  line-height: 1;
}

.session-ledger__cell--danger .session-ledger__value {
  color: #cf1322;
}

.session-ledger__note {
  color: var(--text-color-3, #999);
  font-size: 11px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.session-account-tabs :deep(.n-tabs-pane-wrapper) {
  display: none;
}

@media (max-width: 1280px) {
  .session-ledger {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .session-ledger__cell {
    border-bottom: 1px solid var(--border-color, #eef2f7);
  }

  .session-ledger__cell:nth-child(3n) {
    border-right: 0;
  }

  .session-ledger__cell:nth-last-child(-n + 3) {
    border-bottom: 0;
  }
}

@media (max-width: 900px) {
  .session-head {
    flex-wrap: wrap;
  }

  .session-ledger__note {
    display: none;
  }

  .session-ledger__value {
    font-size: 20px;
  }
}

@media (max-width: 720px) {
  .session-ledger {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .session-ledger__cell:nth-child(3n) {
    border-right: 1px solid var(--border-color, #eef2f7);
  }

  .session-ledger__cell:nth-child(2n) {
    border-right: 0;
  }

  .session-ledger__cell:nth-last-child(-n + 3) {
    border-bottom: 1px solid var(--border-color, #eef2f7);
  }

  .session-ledger__cell:nth-last-child(-n + 2) {
    border-bottom: 0;
  }
}
</style>
