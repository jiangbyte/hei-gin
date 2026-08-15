<!-- Author: Charlie -->

<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import {
  ACCOUNT_TYPE_TABS,
  accountConfigKey,
  accountTypeLabel,
  createAccountTypeMap,
  mapAccountTypes,
  type AccountType,
} from '@/constants/account'
import { loadByCategory, parseBool, saveByKeys, toBoolStr } from '../composables/useConfigForm'

const CATEGORY = 'AUTH_OAUTH'
const PREFIX = 'AUTH_OAUTH'

const PROVIDERS = [
  { key: 'GITHUB', label: 'GitHub', fields: ['CLIENT_ID', 'CLIENT_SECRET', 'REDIRECT_URI'] as const },
  { key: 'GITEE', label: 'Gitee', fields: ['CLIENT_ID', 'CLIENT_SECRET', 'REDIRECT_URI'] as const },
  { key: 'QQ', label: 'QQ', fields: ['CLIENT_ID', 'CLIENT_SECRET', 'REDIRECT_URI'] as const },
  {
    key: 'WECHAT_OPEN',
    label: '微信开放平台',
    fields: ['CLIENT_ID', 'CLIENT_SECRET', 'REDIRECT_URI'] as const,
  },
  { key: 'WECHAT_MP', label: '微信小程序', fields: ['APP_ID', 'APP_SECRET'] as const, portalOnly: true },
] as const

type ProviderKey = (typeof PROVIDERS)[number]['key']

type ProviderForm = {
  enabled: boolean
  clientId: string
  clientSecret: string
  redirectUri: string
  appId: string
  appSecret: string
}

type ScopeForm = Record<ProviderKey, ProviderForm>

function emptyProvider(): ProviderForm {
  return {
    enabled: false,
    clientId: '',
    clientSecret: '',
    redirectUri: '',
    appId: '',
    appSecret: '',
  }
}

function emptyScope(): ScopeForm {
  return {
    GITHUB: emptyProvider(),
    GITEE: emptyProvider(),
    QQ: emptyProvider(),
    WECHAT_OPEN: emptyProvider(),
    WECHAT_MP: emptyProvider(),
  }
}

const state = reactive({
  loading: false,
  saving: false,
  subTab: 'PORTAL' as AccountType,
  providerTab: 'GITHUB' as ProviderKey,
  byType: createAccountTypeMap(emptyScope),
  frontendPortal: '',
  frontendAdmin: '',
  snapshot: '',
})

onMounted(() => {
  void reload()
})

function providerConfigKey(type: AccountType, provider: string, field: string) {
  return accountConfigKey(PREFIX, type, `${provider}_${field}`)
}

function fillProvider(map: Record<string, string>, type: AccountType, provider: ProviderKey): ProviderForm {
  return {
    enabled: parseBool(map[providerConfigKey(type, provider, 'ENABLED')]),
    clientId: map[providerConfigKey(type, provider, 'CLIENT_ID')] || '',
    clientSecret: map[providerConfigKey(type, provider, 'CLIENT_SECRET')] || '',
    redirectUri: map[providerConfigKey(type, provider, 'REDIRECT_URI')] || '',
    appId: map[providerConfigKey(type, provider, 'APP_ID')] || '',
    appSecret: map[providerConfigKey(type, provider, 'APP_SECRET')] || '',
  }
}

async function reload() {
  state.loading = true
  try {
    const map = await loadByCategory(CATEGORY)
    for (const type of mapAccountTypes((t) => t)) {
      const scope = emptyScope()
      for (const provider of PROVIDERS) {
        scope[provider.key] = fillProvider(map, type, provider.key)
      }
      state.byType[type] = scope
    }
    state.frontendPortal = map.AUTH_OAUTH_FRONTEND_CALLBACK_PORTAL || ''
    state.frontendAdmin = map.AUTH_OAUTH_FRONTEND_CALLBACK_ADMIN || ''
    state.snapshot = JSON.stringify({
      byType: state.byType,
      frontendPortal: state.frontendPortal,
      frontendAdmin: state.frontendAdmin,
    })
  } finally {
    state.loading = false
  }
}

function reset() {
  if (!state.snapshot) return
  const data = JSON.parse(state.snapshot) as {
    byType: Record<AccountType, ScopeForm>
    frontendPortal: string
    frontendAdmin: string
  }
  for (const type of Object.keys(data.byType) as AccountType[]) {
    Object.assign(state.byType[type], data.byType[type])
  }
  state.frontendPortal = data.frontendPortal
  state.frontendAdmin = data.frontendAdmin
}

const current = computed(() => state.byType[state.subTab])
const visibleProviders = computed(() =>
  PROVIDERS.filter((item) => !('portalOnly' in item && item.portalOnly && state.subTab === 'ADMIN')),
)

async function save() {
  state.saving = true
  try {
    const items: Array<{ config_key: string; config_value: string; category: string }> = [
      {
        config_key: 'AUTH_OAUTH_FRONTEND_CALLBACK_PORTAL',
        config_value: state.frontendPortal,
        category: CATEGORY,
      },
      {
        config_key: 'AUTH_OAUTH_FRONTEND_CALLBACK_ADMIN',
        config_value: state.frontendAdmin,
        category: CATEGORY,
      },
    ]
    for (const type of mapAccountTypes((t) => t)) {
      for (const provider of PROVIDERS) {
        if ('portalOnly' in provider && provider.portalOnly && type === 'ADMIN') continue
        const form = state.byType[type][provider.key]
        items.push({
          config_key: providerConfigKey(type, provider.key, 'ENABLED'),
          config_value: toBoolStr(form.enabled),
          category: CATEGORY,
        })
        for (const field of provider.fields) {
          let value = ''
          if (field === 'CLIENT_ID') value = form.clientId
          if (field === 'CLIENT_SECRET') value = form.clientSecret
          if (field === 'REDIRECT_URI') value = form.redirectUri
          if (field === 'APP_ID') value = form.appId
          if (field === 'APP_SECRET') value = form.appSecret
          items.push({
            config_key: providerConfigKey(type, provider.key, field),
            config_value: value,
            category: CATEGORY,
          })
        }
      }
    }
    await saveByKeys(items)
    window.$message.success('保存成功')
    state.snapshot = JSON.stringify({
      byType: state.byType,
      frontendPortal: state.frontendPortal,
      frontendAdmin: state.frontendAdmin,
    })
  } finally {
    state.saving = false
  }
}
</script>

<template>
  <NSpin :show="state.loading">
    <NTabs
      v-model:value="state.subTab"
      type="line"
      class="sys-config-subnav"
    >
      <NTab
        v-for="item in ACCOUNT_TYPE_TABS"
        :key="item.key"
        :name="item.key"
        :tab="item.label"
      />
    </NTabs>

    <ConfigSectionLayout
      :description="`配置「${accountTypeLabel(state.subTab)}」三方登录开关与应用凭据。管理端禁止三方自动开户。`"
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NForm
        class="sys-config-form"
        label-placement="top"
      >
        <NFormItem label="前端 OAuth 回调页（门户）">
          <NInput
            v-model:value="state.frontendPortal"
            placeholder="如 https://portal.example.com/auth/oauth/callback"
          />
        </NFormItem>
        <NFormItem label="前端 OAuth 回调页（管理端）">
          <NInput
            v-model:value="state.frontendAdmin"
            placeholder="如 https://admin.example.com/auth/oauth/callback"
          />
        </NFormItem>

        <NTabs
          v-model:value="state.providerTab"
          type="segment"
          class="mb-12px"
        >
          <NTab
            v-for="item in visibleProviders"
            :key="item.key"
            :name="item.key"
            :tab="item.label"
          />
        </NTabs>

        <template
          v-for="item in visibleProviders"
          :key="item.key"
        >
          <div v-show="state.providerTab === item.key">
            <NFormItem :label="`启用 ${item.label}`">
              <NSwitch v-model:value="current[item.key].enabled" />
            </NFormItem>
            <template v-if="item.key === 'WECHAT_MP'">
              <NFormItem label="AppId">
                <NInput v-model:value="current[item.key].appId" />
              </NFormItem>
              <NFormItem label="AppSecret">
                <NInput
                  v-model:value="current[item.key].appSecret"
                  type="password"
                  show-password-on="click"
                />
              </NFormItem>
            </template>
            <template v-else>
              <NFormItem :label="item.key.startsWith('WECHAT') ? 'AppId' : 'Client ID'">
                <NInput v-model:value="current[item.key].clientId" />
              </NFormItem>
              <NFormItem :label="item.key.startsWith('WECHAT') ? 'AppSecret' : 'Client Secret'">
                <NInput
                  v-model:value="current[item.key].clientSecret"
                  type="password"
                  show-password-on="click"
                />
              </NFormItem>
              <NFormItem label="Redirect URI（后端回调）">
                <NInput
                  v-model:value="current[item.key].redirectUri"
                  :placeholder="`/api/v1/${state.subTab.toLowerCase()}/oauth/${item.key.toLowerCase()}/callback`"
                />
              </NFormItem>
            </template>
          </div>
        </template>
      </NForm>
    </ConfigSectionLayout>
  </NSpin>
</template>

<style scoped>
.mb-12px {
  margin-bottom: 12px;
}
</style>
