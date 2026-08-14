<!-- Author: Charlie -->

<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import { accountApi, authApi } from '@/api'
import { createTagColor, displayValue, formatDateTime, hasPermission } from '@/utils'
import { wireBool } from '@/utils/wire'
import { NButton } from 'naive-ui'
import { computed, reactive } from 'vue'
import { dictTypeColor, dictTypeData } from '@/utils/dict'

const state = reactive({
  showModal: false,
  loading: false,
  unbinding: '' as string,
  account: {} as any,
})

const avatarAlt = computed(() => state.account?.nickname || state.account?.name || '门户用户头像')
const avatarUrl = computed(() => state.account?.avatar || undefined)
const avatarImgProps = { referrerPolicy: 'no-referrer' } as any
const oauthBindings = computed(() =>
  Array.isArray(state.account?.oauth_bindings) ? state.account.oauth_bindings : [],
)

function maskOpenId(value?: string) {
  const text = String(value || '')
  if (text.length <= 8) return text ? '****' : '-'
  return `${text.slice(0, 4)}****${text.slice(-4)}`
}

function providerLabel(provider?: string) {
  return dictTypeData('OAUTH_PROVIDER', provider || '') || displayValue(provider)
}

const oauthColumns = computed<DataTableColumns<any>>(() => [
  {
    title: '提供商',
    key: 'provider',
    render: (row) => providerLabel(row.provider),
  },
  {
    title: 'OpenID',
    key: 'open_id',
    render: (row) => maskOpenId(row.open_id),
  },
  {
    title: '昵称',
    key: 'nickname',
    render: (row) => displayValue(row.nickname),
  },
  {
    title: '绑定时间',
    key: 'bound_at',
    render: (row) => formatDateTime(row.bound_at),
  },
  {
    title: '操作',
    key: 'actions',
    render: (row) =>
      hasPermission('iam:account:update') ? (
        <NButton
          text
          type="error"
          size="small"
          loading={state.unbinding === row.provider}
          onClick={() => void unbindOauth(row.provider)}
        >
          解绑
        </NButton>
      ) : null,
  },
])

async function openModal(id: string) {
  state.account = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await accountApi.detail({ id })
    state.account = response.data ?? {}
  } finally {
    state.loading = false
  }
}

async function unbindOauth(provider: string) {
  if (!state.account?.id || state.unbinding) return
  state.unbinding = provider
  try {
    await authApi.adminOauthUnbind({
      account_id: String(state.account.id),
      provider,
    })
    window.$message.success('已解绑')
    await fetchDetail(String(state.account.id))
  } finally {
    state.unbinding = ''
  }
}

defineExpose({ openModal })
</script>

<template>
  <NModal
    v-model:show="state.showModal"
    preset="card"
    draggable
    :mask-closable="false"
    title="门户用户详情"
    style="width: 680px"
  >
    <NScrollbar class="h-[480px] pr-16px">
      <NSpin :show="state.loading">
        <NTabs
          type="line"
          animated
        >
          <NTabPane
            name="account"
            tab="账号"
          >
            <NDescriptions
              label-placement="left"
              bordered
              :column="1"
            >
              <NDescriptionsItem label="账号 ID">
                {{ displayValue(state.account.id) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="账号状态">
                <NTag
                  :color="
                    createTagColor(dictTypeColor('ACCOUNT_STATUS', state.account.account_status))
                  "
                  :bordered="false"
                >
                  {{ dictTypeData('ACCOUNT_STATUS', state.account.account_status) }}
                </NTag>
              </NDescriptionsItem>
              <NDescriptionsItem label="注销时间">
                {{ formatDateTime(state.account.cancelled_at) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="注销人">
                {{ displayValue(state.account.cancelled_by) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="注销原因">
                {{ displayValue(state.account.cancel_reason) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="上次登录 IP">
                {{ displayValue(state.account.last_login_ip) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="上次登录地址">
                {{ displayValue(state.account.last_login_address) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="上次登录时间">
                {{ formatDateTime(state.account.last_login_time) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="上次登录设备">
                {{ displayValue(state.account.last_login_device) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="最近登录 IP">
                {{ displayValue(state.account.latest_login_ip) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="最近登录地址">
                {{ displayValue(state.account.latest_login_address) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="最近登录时间">
                {{ formatDateTime(state.account.latest_login_time) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="最近登录设备">
                {{ displayValue(state.account.latest_login_device) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="创建时间">
                {{ formatDateTime(state.account.created_at) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="更新时间">
                {{ formatDateTime(state.account.updated_at) }}
              </NDescriptionsItem>
            </NDescriptions>
          </NTabPane>

          <NTabPane
            name="identity"
            tab="登录身份"
          >
            <NDescriptions
              label-placement="left"
              bordered
              :column="1"
            >
              <NDescriptionsItem label="账号">
                {{ displayValue(state.account.account) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="邮箱身份">
                {{ displayValue(state.account.email_identity) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="启用邮箱登录">
                {{ wireBool(state.account.email_login_enabled ?? false) ? '是' : '否' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="邮箱已验证">
                {{ wireBool(state.account.email_identity_verified ?? false) ? '是' : '否' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="邮箱绑定状态">
                {{
                  dictTypeData(
                    'ACCOUNT_IDENTITY_BIND_STATUS',
                    state.account.email_identity_bind_status,
                  ) || displayValue(state.account.email_identity_bind_status)
                }}
              </NDescriptionsItem>
              <NDescriptionsItem label="手机号身份">
                {{ displayValue(state.account.phone_identity) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="启用手机号登录">
                {{ wireBool(state.account.phone_login_enabled ?? false) ? '是' : '否' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="手机号已验证">
                {{ wireBool(state.account.phone_identity_verified ?? false) ? '是' : '否' }}
              </NDescriptionsItem>
              <NDescriptionsItem label="手机号绑定状态">
                {{
                  dictTypeData(
                    'ACCOUNT_IDENTITY_BIND_STATUS',
                    state.account.phone_identity_bind_status,
                  ) || displayValue(state.account.phone_identity_bind_status)
                }}
              </NDescriptionsItem>
            </NDescriptions>

            <div class="mt-16px mb-8px font-600">
              三方绑定
            </div>
            <NEmpty
              v-if="!oauthBindings.length"
              description="暂无三方绑定"
            />
            <NDataTable
              v-else
              size="small"
              :bordered="false"
              :single-line="false"
              :columns="oauthColumns"
              :data="oauthBindings"
            />
          </NTabPane>

          <NTabPane
            name="profile"
            tab="资料"
          >
            <NDescriptions
              label-placement="left"
              bordered
              :column="1"
            >
              <NDescriptionsItem label="姓名">
                {{ displayValue(state.account.name) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="昵称">
                {{ displayValue(state.account.nickname) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="头像">
                <NAvatar
                  v-if="avatarUrl"
                  :src="avatarUrl"
                  :alt="avatarAlt"
                  :img-props="avatarImgProps"
                />
                <template v-else>
                  -
                </template>
              </NDescriptionsItem>
              <NDescriptionsItem label="个性签名">
                {{ displayValue(state.account.signature) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="手机号">
                {{ displayValue(state.account.phone) }}
              </NDescriptionsItem>
              <NDescriptionsItem label="邮箱">
                {{ displayValue(state.account.email) }}
              </NDescriptionsItem>
            </NDescriptions>
          </NTabPane>
        </NTabs>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>

<style scoped>
.mt-16px {
  margin-top: 16px;
}

.mb-8px {
  margin-bottom: 8px;
}

.font-600 {
  font-weight: 600;
}
</style>
