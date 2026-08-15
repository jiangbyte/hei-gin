<!-- Author: Charlie -->

<script setup lang="ts">
import { authApi } from '@/api'
import { wireBool } from '@/utils/wire'
import { computed, onMounted, reactive } from 'vue'
import '../profile.css'

type Binding = {
  provider: string
  label: string
  open_id_masked?: string
  nickname?: string
  avatar?: string
  bound_at?: string
}

type ProviderOption = {
  provider: string
  label: string
  enabled: boolean
  web_oauth: boolean
}

const state = reactive({
  loading: false,
  bindingProvider: null as string | null,
  bindings: [] as Binding[],
  providers: [] as ProviderOption[],
})

const boundSet = computed(() => new Set(state.bindings.map((item) => item.provider)))

onMounted(() => {
  void refresh()
})

async function refresh() {
  state.loading = true
  try {
    const [bindRes, optRes] = await Promise.all([authApi.oauthBindings(), authApi.authOptions()])
    state.bindings = Array.isArray(bindRes?.data) ? bindRes.data : []
    const list = Array.isArray(optRes?.data?.oauth_providers) ? optRes.data.oauth_providers : []
    state.providers = list
      .map((item: any) => ({
        provider: String(item.provider || ''),
        label: String(item.label || item.provider || ''),
        enabled: wireBool(item.enabled ?? false),
        web_oauth: wireBool(item.web_oauth ?? true),
      }))
      .filter((item: ProviderOption) => item.provider && item.enabled && item.web_oauth)
  } finally {
    state.loading = false
  }
}

async function bind(provider: string) {
  if (state.bindingProvider) return
  state.bindingProvider = provider
  try {
    const res = await authApi.oauthBindAuthorize(provider)
    const url = res?.data?.authorize_url
    if (!url) {
      window.$message.error('无法发起绑定')
      return
    }
    window.location.href = String(url)
  } catch {
    // 全局错误提示
  } finally {
    state.bindingProvider = null
  }
}

async function unbind(provider: string) {
  await authApi.oauthUnbind(provider)
  window.$message.success('已解绑')
  await refresh()
}
</script>

<template>
  <NSpin :show="state.loading">
    <div class="uc-panel">
      <p class="uc-panel__desc">
        绑定后可使用对应平台快速登录。管理端禁止三方自动开户，需先用密码登录再绑定。至少保留一种登录方式。
      </p>

      <NEmpty
        v-if="!state.bindings.length"
        description="暂无已绑定三方账号"
        class="mb-16px"
      />
      <NList
        v-else
        bordered
      >
        <NListItem
          v-for="item in state.bindings"
          :key="item.provider"
        >
          <template #prefix>
            <NAvatar :src="item.avatar || undefined">
              {{ (item.label || item.provider || '?').slice(0, 1) }}
            </NAvatar>
          </template>
          <NThing
            :title="`${item.label || item.provider}${item.nickname ? ` · ${item.nickname}` : ''}`"
            :description="`OpenID：${item.open_id_masked || '-'} · 绑定于 ${item.bound_at || '-'}`"
          />
          <template #suffix>
            <NPopconfirm @positive-click="() => unbind(item.provider)">
              <template #trigger>
                <NButton
                  type="error"
                  text
                >
                  解绑
                </NButton>
              </template>
              确认解绑该三方账号？
            </NPopconfirm>
          </template>
        </NListItem>
      </NList>

      <h3 class="uc-panel__subtitle">
        可绑定平台
      </h3>
      <NSpace>
        <NButton
          v-for="item in state.providers"
          :key="item.provider"
          :disabled="boundSet.has(item.provider) || Boolean(state.bindingProvider)"
          :loading="state.bindingProvider === item.provider"
          @click="bind(item.provider)"
        >
          {{ boundSet.has(item.provider) ? `已绑定 ${item.label}` : `绑定 ${item.label}` }}
        </NButton>
        <NText
          v-if="!state.providers.length"
          depth="3"
        >
          暂无可绑定的三方登录（未开启或未配置）
        </NText>
      </NSpace>
    </div>
  </NSpin>
</template>

<style scoped>
.uc-panel__desc {
  margin: 0 0 16px;
  color: var(--n-text-color-3, #6b7280);
  font-size: 13px;
  line-height: 1.6;
}

.uc-panel__subtitle {
  margin: 24px 0 12px;
  font-size: 15px;
  font-weight: 600;
}

.mb-16px {
  margin-bottom: 16px;
}
</style>
