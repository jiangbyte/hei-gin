<!-- Author: Charlie -->

<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import { loadByCategory, parseNumber, saveByKeys } from '../composables/useConfigForm'

const CATEGORY = 'AUTH_TOKEN'

const state = reactive({
  loading: false,
  saving: false,
  tokenTtl: 2592000,
  resetTtl: 600,
  resetUrlAdmin: 'http://localhost:5173/auth/forgot-password',
  resetUrlPortal: 'http://localhost:5174/auth/forgot-password',
  snapshot: '',
})

onMounted(() => {
  void reload()
})

function snapshotOf() {
  return JSON.stringify({
    tokenTtl: state.tokenTtl,
    resetTtl: state.resetTtl,
    resetUrlAdmin: state.resetUrlAdmin,
    resetUrlPortal: state.resetUrlPortal,
  })
}

async function reload() {
  state.loading = true
  try {
    const map = await loadByCategory(CATEGORY)
    state.tokenTtl = parseNumber(map.AUTH_TOKEN_TTL_SECONDS, 2592000)
    state.resetTtl = parseNumber(map.AUTH_PASSWORD_RESET_TOKEN_TTL_SECONDS, 600)
    state.resetUrlAdmin =
      map.AUTH_PASSWORD_RESET_URL_ADMIN || 'http://localhost:5173/auth/forgot-password'
    state.resetUrlPortal =
      map.AUTH_PASSWORD_RESET_URL_PORTAL || 'http://localhost:5174/auth/forgot-password'
    state.snapshot = snapshotOf()
  } finally {
    state.loading = false
  }
}

function reset() {
  if (!state.snapshot) return
  const data = JSON.parse(state.snapshot)
  state.tokenTtl = data.tokenTtl
  state.resetTtl = data.resetTtl
  state.resetUrlAdmin = data.resetUrlAdmin
  state.resetUrlPortal = data.resetUrlPortal
}

async function save() {
  state.saving = true
  try {
    await saveByKeys([
      {
        config_key: 'AUTH_TOKEN_TTL_SECONDS',
        config_value: String(state.tokenTtl),
        category: CATEGORY,
      },
      {
        config_key: 'AUTH_PASSWORD_RESET_TOKEN_TTL_SECONDS',
        config_value: String(state.resetTtl),
        category: CATEGORY,
      },
      {
        config_key: 'AUTH_PASSWORD_RESET_URL_ADMIN',
        config_value: state.resetUrlAdmin.trim(),
        category: CATEGORY,
      },
      {
        config_key: 'AUTH_PASSWORD_RESET_URL_PORTAL',
        config_value: state.resetUrlPortal.trim(),
        category: CATEGORY,
      },
    ])
    window.$message.success('保存成功')
    state.snapshot = snapshotOf()
  } finally {
    state.saving = false
  }
}
</script>

<template>
  <NSpin :show="state.loading">
    <ConfigSectionLayout
      description="配置会话令牌、密码重置令牌有效期，以及各端密码重置页完整 URL。"
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NForm
        class="sys-config-form"
        label-placement="top"
      >
        <NFormItem label="会话令牌过期时间（秒）">
          <NInputNumber
            v-model:value="state.tokenTtl"
            class="w-full"
            :min="0"
          />
        </NFormItem>
        <NFormItem label="密码重置令牌有效期（秒）">
          <NInputNumber
            v-model:value="state.resetTtl"
            class="w-full"
            :min="0"
          />
        </NFormItem>
        <NFormItem label="ADMIN 密码重置页 URL">
          <NInput
            v-model:value="state.resetUrlAdmin"
            clearable
            placeholder="http://localhost:5173/auth/forgot-password"
          />
        </NFormItem>
        <NFormItem label="PORTAL 密码重置页 URL">
          <NInput
            v-model:value="state.resetUrlPortal"
            clearable
            placeholder="http://localhost:5174/auth/forgot-password"
          />
        </NFormItem>
      </NForm>
    </ConfigSectionLayout>
  </NSpin>
</template>
