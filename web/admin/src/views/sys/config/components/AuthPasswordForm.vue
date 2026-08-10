<!-- Author: Charlie -->

<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import WeakPasswordPanel from './WeakPasswordPanel.vue'
import {
  loadByCategory,
  parseBool,
  parseNumber,
  saveByKeys,
  toBoolStr,
} from '../composables/useConfigForm'

const CATEGORY = 'AUTH_PASSWORD'

const verifyOptions = [
  { label: '旧密码', value: 'OLD_PASSWORD' },
  { label: '邮箱验证码', value: 'EMAIL_CODE' },
  { label: '手机验证码', value: 'PHONE_CODE' },
]

const complexityOptions = [
  { label: '无限制', value: 'NO_LIMIT' },
  { label: '必须包含数字和字母', value: 'DIGITS_AND_LETTERS' },
  { label: '必须包含数字和大写字母', value: 'DIGITS_AND_UPPERCASE' },
  {
    label: '必须包含数字、大写、小写和特殊字符',
    value: 'DIGITS_UPPER_LOWER_SPECIAL',
  },
  { label: '字母/数字/特殊字符至少两类', value: 'TWO_OF_THREE' },
  { label: '大写/小写/数字/特殊字符至少三类', value: 'THREE_OF_FOUR' },
]

const state = reactive({
  loading: false,
  saving: false,
  subTab: 'POLICY' as 'POLICY' | 'WEAK',
  defaultPassword: '',
  changeVerifyMethod: 'OLD_PASSWORD',
  minLength: 6,
  maxLength: 20,
  complexity: 'NO_LIMIT',
  maxConsecutiveChars: 3,
  forbidUserInfo: true,
  forbidHistorical: true,
  historyCheckCount: 3,
  forbidWeakList: true,
  validityDays: 30,
  expiryWarningDays: 3,
  snapshot: '',
})

onMounted(() => {
  void reload()
})

async function reload() {
  state.loading = true
  try {
    const map = await loadByCategory(CATEGORY)
    state.defaultPassword = map.AUTH_DEFAULT_PASSWORD || ''
    state.changeVerifyMethod = map.PASSWORD_CHANGE_VERIFY_METHOD || 'OLD_PASSWORD'
    state.minLength = parseNumber(map.PASSWORD_MIN_LENGTH, 6)
    state.maxLength = parseNumber(map.PASSWORD_MAX_LENGTH, 20)
    state.complexity = map.PASSWORD_COMPLEXITY || 'NO_LIMIT'
    state.maxConsecutiveChars = parseNumber(map.PASSWORD_MAX_CONSECUTIVE_CHARS, 3)
    state.forbidUserInfo = parseBool(map.PASSWORD_FORBID_USER_INFO)
    state.forbidHistorical = parseBool(map.PASSWORD_FORBID_HISTORICAL)
    state.historyCheckCount = parseNumber(map.PASSWORD_HISTORY_CHECK_COUNT, 3)
    state.forbidWeakList = parseBool(map.PASSWORD_FORBID_WEAK_LIST)
    state.validityDays = parseNumber(map.PASSWORD_VALIDITY_DAYS, 30)
    state.expiryWarningDays = parseNumber(map.PASSWORD_EXPIRY_WARNING_DAYS, 3)
    state.snapshot = JSON.stringify({
      defaultPassword: state.defaultPassword,
      changeVerifyMethod: state.changeVerifyMethod,
      minLength: state.minLength,
      maxLength: state.maxLength,
      complexity: state.complexity,
      maxConsecutiveChars: state.maxConsecutiveChars,
      forbidUserInfo: state.forbidUserInfo,
      forbidHistorical: state.forbidHistorical,
      historyCheckCount: state.historyCheckCount,
      forbidWeakList: state.forbidWeakList,
      validityDays: state.validityDays,
      expiryWarningDays: state.expiryWarningDays,
    })
  } finally {
    state.loading = false
  }
}

function reset() {
  if (!state.snapshot) return
  Object.assign(state, JSON.parse(state.snapshot))
}

async function save() {
  state.saving = true
  try {
    await saveByKeys([
      {
        config_key: 'AUTH_DEFAULT_PASSWORD',
        config_value: state.defaultPassword,
        category: CATEGORY,
      },
      {
        config_key: 'PASSWORD_CHANGE_VERIFY_METHOD',
        config_value: state.changeVerifyMethod,
        category: CATEGORY,
      },
      {
        config_key: 'PASSWORD_MIN_LENGTH',
        config_value: String(state.minLength),
        category: CATEGORY,
      },
      {
        config_key: 'PASSWORD_MAX_LENGTH',
        config_value: String(state.maxLength),
        category: CATEGORY,
      },
      {
        config_key: 'PASSWORD_COMPLEXITY',
        config_value: state.complexity,
        category: CATEGORY,
      },
      {
        config_key: 'PASSWORD_MAX_CONSECUTIVE_CHARS',
        config_value: String(state.maxConsecutiveChars),
        category: CATEGORY,
      },
      {
        config_key: 'PASSWORD_FORBID_USER_INFO',
        config_value: toBoolStr(state.forbidUserInfo),
        category: CATEGORY,
      },
      {
        config_key: 'PASSWORD_FORBID_HISTORICAL',
        config_value: toBoolStr(state.forbidHistorical),
        category: CATEGORY,
      },
      {
        config_key: 'PASSWORD_HISTORY_CHECK_COUNT',
        config_value: String(state.historyCheckCount),
        category: CATEGORY,
      },
      {
        config_key: 'PASSWORD_FORBID_WEAK_LIST',
        config_value: toBoolStr(state.forbidWeakList),
        category: CATEGORY,
      },
      {
        config_key: 'PASSWORD_VALIDITY_DAYS',
        config_value: String(state.validityDays),
        category: CATEGORY,
      },
      {
        config_key: 'PASSWORD_EXPIRY_WARNING_DAYS',
        config_value: String(state.expiryWarningDays),
        category: CATEGORY,
      },
    ])
    window.$message.success('保存成功')
    await reload()
  } finally {
    state.saving = false
  }
}
</script>

<template>
  <div>
    <NTabs
      v-model:value="state.subTab"
      type="line"
      class="sys-config-subnav"
    >
      <NTab
        name="POLICY"
        tab="密码策略"
      />
      <NTab
        name="WEAK"
        tab="弱密码库"
      />
    </NTabs>

    <WeakPasswordPanel v-if="state.subTab === 'WEAK'" />

    <NSpin
      v-else
      :show="state.loading"
    >
      <ConfigSectionLayout
        description="全局密码策略，管理端与门户共用；保存后热重载生效。"
        :saving="state.saving"
        @save="save"
        @reset="reset"
      >
        <NForm
          class="sys-config-form sys-config-form--wide"
          label-placement="top"
        >
          <NCard
            title="基础"
            size="small"
            :bordered="false"
          >
            <NGrid
              :cols="24"
              :x-gap="16"
            >
              <NGi :span="12">
                <NFormItem label="默认用户密码">
                  <NInput
                    v-model:value="state.defaultPassword"
                    type="password"
                    show-password-on="click"
                    placeholder="留空不修改"
                  />
                </NFormItem>
                <p class="sys-config__hint">
                  敏感项；用于新建账户等场景
                </p>
              </NGi>
              <NGi :span="12">
                <NFormItem label="修改密码验证方式">
                  <NSelect
                    v-model:value="state.changeVerifyMethod"
                    :options="verifyOptions"
                  />
                </NFormItem>
              </NGi>
            </NGrid>
          </NCard>

          <NCard
            title="强度规则"
            size="small"
            :bordered="false"
            class="mt-12px"
          >
            <NGrid
              :cols="24"
              :x-gap="16"
            >
              <NGi :span="24">
                <NFormItem label="复杂度">
                  <NSelect
                    v-model:value="state.complexity"
                    :options="complexityOptions"
                  />
                </NFormItem>
              </NGi>
              <NGi :span="8">
                <NFormItem label="最小长度">
                  <NInputNumber
                    v-model:value="state.minLength"
                    class="w-full"
                    :min="1"
                  />
                </NFormItem>
              </NGi>
              <NGi :span="8">
                <NFormItem label="最大长度">
                  <NInputNumber
                    v-model:value="state.maxLength"
                    class="w-full"
                    :min="1"
                  />
                </NFormItem>
              </NGi>
              <NGi :span="8">
                <NFormItem label="连续相同字符上限">
                  <NInputNumber
                    v-model:value="state.maxConsecutiveChars"
                    class="w-full"
                    :min="1"
                  />
                </NFormItem>
              </NGi>
            </NGrid>
          </NCard>

          <NCard
            title="禁止项"
            size="small"
            :bordered="false"
            class="mt-12px"
          >
            <div class="pwd-rule">
              <div class="pwd-rule__row pwd-rule__row--single">
                <div class="pwd-rule__switch">
                  <span class="pwd-rule__label">禁止包含用户信息</span>
                  <NSwitch v-model:value="state.forbidUserInfo" />
                </div>
              </div>

              <div class="pwd-rule__row pwd-rule__row--single">
                <div class="pwd-rule__switch">
                  <span class="pwd-rule__label">禁止弱密码库</span>
                  <NSwitch v-model:value="state.forbidWeakList" />
                </div>
              </div>

              <div class="pwd-rule__row">
                <div class="pwd-rule__switch">
                  <span class="pwd-rule__label">禁止历史密码</span>
                  <NSwitch v-model:value="state.forbidHistorical" />
                </div>
                <NFormItem
                  v-if="state.forbidHistorical"
                  label="检查最近个数"
                  class="pwd-rule__extra"
                  :show-feedback="false"
                >
                  <NInputNumber
                    v-model:value="state.historyCheckCount"
                    class="w-full"
                    :min="0"
                  />
                </NFormItem>
              </div>
            </div>
          </NCard>

          <NCard
            title="有效期"
            size="small"
            :bordered="false"
            class="mt-12px"
          >
            <NGrid
              :cols="24"
              :x-gap="16"
            >
              <NGi :span="12">
                <NFormItem label="有效期（天）">
                  <NInputNumber
                    v-model:value="state.validityDays"
                    class="w-full"
                    :min="0"
                  />
                </NFormItem>
                <p class="sys-config__hint">
                  0 表示不过期
                </p>
              </NGi>
              <NGi :span="12">
                <NFormItem label="过期提醒（天）">
                  <NInputNumber
                    v-model:value="state.expiryWarningDays"
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
  </div>
</template>

<style scoped>
.mt-12px {
  margin-top: 12px;
}

.pwd-rule {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.pwd-rule__row {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr);
  gap: 24px;
  align-items: end;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.pwd-rule__row:last-child {
  border-bottom: 0;
  padding-bottom: 0;
}

.pwd-rule__row--single {
  grid-template-columns: 220px;
}

.pwd-rule__switch {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 34px;
  padding-bottom: 6px;
}

.pwd-rule__label {
  font-size: 14px;
  color: var(--text-color-1);
  line-height: 1.4;
}

.pwd-rule__extra {
  margin-bottom: 0 !important;
}

@media (max-width: 720px) {
  .pwd-rule__row {
    grid-template-columns: 1fr;
    gap: 8px;
    align-items: stretch;
  }
}
</style>
