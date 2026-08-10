<!-- Author: Charlie -->

<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import {
  loadByCategory,
  parseBool,
  parseNumber,
  saveByKeys,
  toBoolStr,
} from '../composables/useConfigForm'
import {
  ACCOUNT_TYPES,
  ACCOUNT_TYPE_TABS,
  DEFAULT_ACCOUNT_TYPE,
  accountConfigKey,
  accountTypeLabel,
  createAccountTypeMap,
  mapAccountTypes,
  type AccountType,
} from '@/constants/account'

const CATEGORY = 'AUTH_LOGIN'
const PREFIX = 'AUTH_LOGIN'

interface ScopeForm {
  failureWindow: number
  maxFailures: number
  lockSeconds: number
  allowPhone: boolean
  phoneNoUserPolicy: string
  allowEmail: boolean
  emailNoUserPolicy: string
  allowOtp: boolean
}

const policyOptions = [
  { label: '不允许登录', value: 'DENY' },
  { label: '自动创建用户', value: 'AUTO_CREATE' },
]

function emptyScope(): ScopeForm {
  return {
    failureWindow: 300,
    maxFailures: 5,
    lockSeconds: 300,
    allowPhone: true,
    phoneNoUserPolicy: 'DENY',
    allowEmail: true,
    emailNoUserPolicy: 'DENY',
    allowOtp: true,
  }
}

const state = reactive({
  loading: false,
  saving: false,
  subTab: DEFAULT_ACCOUNT_TYPE as AccountType,
  byType: createAccountTypeMap(emptyScope),
  shared: {
    failureWindow: 900,
    accountMaxFailures: 5,
    ipMaxFailures: 30,
    lockSeconds: 900,
  },
  snapshot: '',
})

onMounted(() => {
  void reload()
})

function fillScope(map: Record<string, string>, type: AccountType): ScopeForm {
  return {
    failureWindow: parseNumber(map[accountConfigKey(PREFIX, type, 'FAILURE_WINDOW_SECONDS')], 300),
    maxFailures: parseNumber(map[accountConfigKey(PREFIX, type, 'MAX_FAILURES')], 5),
    lockSeconds: parseNumber(map[accountConfigKey(PREFIX, type, 'LOCK_SECONDS')], 300),
    allowPhone: parseBool(map[accountConfigKey(PREFIX, type, 'ALLOW_PHONE')]),
    phoneNoUserPolicy: map[accountConfigKey(PREFIX, type, 'PHONE_NO_USER_POLICY')] || 'DENY',
    allowEmail: parseBool(map[accountConfigKey(PREFIX, type, 'ALLOW_EMAIL')]),
    emailNoUserPolicy: map[accountConfigKey(PREFIX, type, 'EMAIL_NO_USER_POLICY')] || 'DENY',
    allowOtp: parseBool(map[accountConfigKey(PREFIX, type, 'ALLOW_OTP')]),
  }
}

async function reload() {
  state.loading = true
  try {
    const map = await loadByCategory(CATEGORY)
    for (const type of ACCOUNT_TYPES) {
      Object.assign(state.byType[type], fillScope(map, type))
    }
    state.shared.failureWindow = parseNumber(map.AUTH_LOGIN_FAILURE_WINDOW_SECONDS, 900)
    state.shared.accountMaxFailures = parseNumber(map.AUTH_LOGIN_ACCOUNT_MAX_FAILURES, 5)
    state.shared.ipMaxFailures = parseNumber(map.AUTH_LOGIN_IP_MAX_FAILURES, 30)
    state.shared.lockSeconds = parseNumber(map.AUTH_LOGIN_LOCK_SECONDS, 900)
    state.snapshot = JSON.stringify({ byType: state.byType, shared: state.shared })
  } finally {
    state.loading = false
  }
}

function reset() {
  if (!state.snapshot) return
  const data = JSON.parse(state.snapshot)
  for (const type of Object.keys(data.byType) as AccountType[]) {
    if (state.byType[type]) Object.assign(state.byType[type], data.byType[type])
  }
  Object.assign(state.shared, data.shared)
}

const current = computed(() => state.byType[state.subTab])

function scopeItems(type: AccountType, form: ScopeForm) {
  return [
    {
      config_key: accountConfigKey(PREFIX, type, 'FAILURE_WINDOW_SECONDS'),
      config_value: String(form.failureWindow),
      category: CATEGORY,
    },
    {
      config_key: accountConfigKey(PREFIX, type, 'MAX_FAILURES'),
      config_value: String(form.maxFailures),
      category: CATEGORY,
    },
    {
      config_key: accountConfigKey(PREFIX, type, 'LOCK_SECONDS'),
      config_value: String(form.lockSeconds),
      category: CATEGORY,
    },
    {
      config_key: accountConfigKey(PREFIX, type, 'ALLOW_PHONE'),
      config_value: toBoolStr(form.allowPhone),
      category: CATEGORY,
    },
    {
      config_key: accountConfigKey(PREFIX, type, 'PHONE_NO_USER_POLICY'),
      config_value: form.phoneNoUserPolicy,
      category: CATEGORY,
    },
    {
      config_key: accountConfigKey(PREFIX, type, 'ALLOW_EMAIL'),
      config_value: toBoolStr(form.allowEmail),
      category: CATEGORY,
    },
    {
      config_key: accountConfigKey(PREFIX, type, 'EMAIL_NO_USER_POLICY'),
      config_value: form.emailNoUserPolicy,
      category: CATEGORY,
    },
    {
      config_key: accountConfigKey(PREFIX, type, 'ALLOW_OTP'),
      config_value: toBoolStr(form.allowOtp),
      category: CATEGORY,
    },
  ]
}

async function save() {
  state.saving = true
  try {
    await saveByKeys([
      ...mapAccountTypes((type) => scopeItems(type, state.byType[type])).flat(),
      {
        config_key: 'AUTH_LOGIN_FAILURE_WINDOW_SECONDS',
        config_value: String(state.shared.failureWindow),
        category: CATEGORY,
      },
      {
        config_key: 'AUTH_LOGIN_ACCOUNT_MAX_FAILURES',
        config_value: String(state.shared.accountMaxFailures),
        category: CATEGORY,
      },
      {
        config_key: 'AUTH_LOGIN_IP_MAX_FAILURES',
        config_value: String(state.shared.ipMaxFailures),
        category: CATEGORY,
      },
      {
        config_key: 'AUTH_LOGIN_LOCK_SECONDS',
        config_value: String(state.shared.lockSeconds),
        category: CATEGORY,
      },
    ])
    window.$message.success('保存成功')
    state.snapshot = JSON.stringify({ byType: state.byType, shared: state.shared })
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
      :description="`配置「${accountTypeLabel(state.subTab)}」账户类型的登录策略。`"
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NForm
        class="sys-config-form sys-config-form--wide"
        label-placement="top"
      >
        <NCard
          title="登录锁定"
          size="small"
          :bordered="false"
        >
          <NGrid
            :cols="24"
            :x-gap="16"
          >
            <NGi :span="8">
              <NFormItem label="失败统计窗口（秒）">
                <NInputNumber
                  v-model:value="current.failureWindow"
                  class="w-full"
                  :min="0"
                />
              </NFormItem>
            </NGi>
            <NGi :span="8">
              <NFormItem label="最大失败次数">
                <NInputNumber
                  v-model:value="current.maxFailures"
                  class="w-full"
                  :min="0"
                />
              </NFormItem>
            </NGi>
            <NGi :span="8">
              <NFormItem label="锁定时间（秒）">
                <NInputNumber
                  v-model:value="current.lockSeconds"
                  class="w-full"
                  :min="0"
                />
              </NFormItem>
            </NGi>
          </NGrid>
        </NCard>

        <NCard
          title="登录方式"
          size="small"
          :bordered="false"
          class="mt-12px"
        >
          <div class="login-method">
            <div class="login-method__row">
              <div class="login-method__switch">
                <span class="login-method__label">手机号登录</span>
                <NSwitch v-model:value="current.allowPhone" />
              </div>
              <NFormItem
                v-if="current.allowPhone"
                label="无用户策略"
                class="login-method__policy"
                :show-feedback="false"
              >
                <NSelect
                  v-model:value="current.phoneNoUserPolicy"
                  :options="policyOptions"
                />
              </NFormItem>
            </div>

            <div class="login-method__row">
              <div class="login-method__switch">
                <span class="login-method__label">邮箱登录</span>
                <NSwitch v-model:value="current.allowEmail" />
              </div>
              <NFormItem
                v-if="current.allowEmail"
                label="无用户策略"
                class="login-method__policy"
                :show-feedback="false"
              >
                <NSelect
                  v-model:value="current.emailNoUserPolicy"
                  :options="policyOptions"
                />
              </NFormItem>
            </div>

            <div class="login-method__row login-method__row--single">
              <div class="login-method__switch">
                <span class="login-method__label">动态口令登录</span>
                <NSwitch v-model:value="current.allowOtp" />
              </div>
            </div>
          </div>
        </NCard>

        <NCard
          title="运行时共享"
          size="small"
          :bordered="false"
          class="mt-12px"
        >
          <p class="sys-config__hint mb-12px">
            各账户类型共用的兜底限流参数
          </p>
          <NGrid
            :cols="24"
            :x-gap="16"
          >
            <NGi :span="6">
              <NFormItem label="失败窗口（秒）">
                <NInputNumber
                  v-model:value="state.shared.failureWindow"
                  class="w-full"
                  :min="0"
                />
              </NFormItem>
            </NGi>
            <NGi :span="6">
              <NFormItem label="账号最大失败">
                <NInputNumber
                  v-model:value="state.shared.accountMaxFailures"
                  class="w-full"
                  :min="0"
                />
              </NFormItem>
            </NGi>
            <NGi :span="6">
              <NFormItem label="IP 最大失败">
                <NInputNumber
                  v-model:value="state.shared.ipMaxFailures"
                  class="w-full"
                  :min="0"
                />
              </NFormItem>
            </NGi>
            <NGi :span="6">
              <NFormItem label="锁定时间（秒）">
                <NInputNumber
                  v-model:value="state.shared.lockSeconds"
                  class="w-full"
                  :min="0"
                />
              </NFormItem>
            </NGi>
          </NGrid>
        </NCard>
      </NForm>
    </ConfigSectionLayout>
  </NSpin>
</template>

<style scoped>
.mt-12px {
  margin-top: 12px;
}

.login-method {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.login-method__row {
  display: grid;
  grid-template-columns: 200px minmax(0, 1fr);
  gap: 24px;
  align-items: end;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.login-method__row:last-child {
  border-bottom: 0;
  padding-bottom: 0;
}

.login-method__row--single {
  grid-template-columns: 200px;
}

.login-method__switch {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 34px;
  padding-bottom: 6px;
}

.login-method__label {
  font-size: 14px;
  color: var(--text-color-1);
  line-height: 1.4;
}

.login-method__policy {
  margin-bottom: 0 !important;
}

@media (max-width: 720px) {
  .login-method__row {
    grid-template-columns: 1fr;
    gap: 8px;
    align-items: stretch;
  }
}
</style>
