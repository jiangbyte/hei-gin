<!-- Author: Charlie -->

<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import {
  loadByCategory,
  parseBool,
  parseNumber,
  saveByKeys,
  toBoolStr,
} from '../composables/useConfigForm'

const CATEGORY = 'MAIL'
type Engine = 'LOCAL' | 'ALIYUN' | 'TENCENT'

const engineOptions = [
  { label: '本地 SMTP', value: 'LOCAL' },
  { label: '阿里云邮件', value: 'ALIYUN' },
  { label: '腾讯云邮件', value: 'TENCENT' },
]

const state = reactive({
  loading: false,
  saving: false,
  subTab: 'LOCAL' as Engine,
  defaultEngine: 'LOCAL' as Engine,
  local: {
    host: '',
    port: 1025,
    username: '',
    password: '',
    fromEmail: '',
    fromName: '',
    authRequired: true,
    useSsl: false,
    useStarttls: false,
  },
  aliyun: {
    accessKeyId: '',
    accessKeySecret: '',
    accountName: '',
  },
  tencent: {
    secretId: '',
    secretKey: '',
    fromEmail: '',
    region: '',
  },
  snapshot: '',
})

onMounted(() => {
  void reload()
})

async function reload() {
  state.loading = true
  try {
    const [map, sys] = await Promise.all([loadByCategory(CATEGORY), loadByCategory('SYS')])
    state.defaultEngine = (map.DEFAULT_EMAIL_ENGINE ||
      sys.DEFAULT_EMAIL_ENGINE ||
      'LOCAL') as Engine
    state.local.host = map.MAIL_LOCAL_HOST || ''
    state.local.port = parseNumber(map.MAIL_LOCAL_PORT, 1025)
    state.local.username = map.MAIL_LOCAL_USERNAME || ''
    state.local.password = map.MAIL_LOCAL_PASSWORD || ''
    state.local.fromEmail = map.MAIL_LOCAL_FROM_EMAIL || ''
    state.local.fromName = map.MAIL_LOCAL_FROM_NAME || ''
    state.local.authRequired = parseBool(map.MAIL_LOCAL_AUTH_REQUIRED)
    state.local.useSsl = parseBool(map.MAIL_LOCAL_USE_SSL)
    state.local.useStarttls = parseBool(map.MAIL_LOCAL_USE_STARTTLS)
    state.aliyun.accessKeyId = map.MAIL_ALIYUN_ACCESS_KEY_ID || ''
    state.aliyun.accessKeySecret = map.MAIL_ALIYUN_ACCESS_KEY_SECRET || ''
    state.aliyun.accountName = map.MAIL_ALIYUN_ACCOUNT_NAME || ''
    state.tencent.secretId = map.MAIL_TENCENT_SECRET_ID || ''
    state.tencent.secretKey = map.MAIL_TENCENT_SECRET_KEY || ''
    state.tencent.fromEmail = map.MAIL_TENCENT_FROM_EMAIL || ''
    state.tencent.region = map.MAIL_TENCENT_REGION || ''
    state.snapshot = JSON.stringify({
      defaultEngine: state.defaultEngine,
      local: state.local,
      aliyun: state.aliyun,
      tencent: state.tencent,
    })
  } finally {
    state.loading = false
  }
}

function reset() {
  if (!state.snapshot) return
  const data = JSON.parse(state.snapshot)
  state.defaultEngine = data.defaultEngine
  Object.assign(state.local, data.local)
  Object.assign(state.aliyun, data.aliyun)
  Object.assign(state.tencent, data.tencent)
}

async function save() {
  state.saving = true
  try {
    await saveByKeys([
      {
        config_key: 'DEFAULT_EMAIL_ENGINE',
        config_value: state.defaultEngine,
        category: CATEGORY,
      },
      { config_key: 'MAIL_LOCAL_HOST', config_value: state.local.host, category: CATEGORY },
      {
        config_key: 'MAIL_LOCAL_PORT',
        config_value: String(state.local.port),
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_LOCAL_USERNAME',
        config_value: state.local.username,
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_LOCAL_PASSWORD',
        config_value: state.local.password,
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_LOCAL_FROM_EMAIL',
        config_value: state.local.fromEmail,
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_LOCAL_FROM_NAME',
        config_value: state.local.fromName,
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_LOCAL_AUTH_REQUIRED',
        config_value: toBoolStr(state.local.authRequired),
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_LOCAL_USE_SSL',
        config_value: toBoolStr(state.local.useSsl),
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_LOCAL_USE_STARTTLS',
        config_value: toBoolStr(state.local.useStarttls),
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_ALIYUN_ACCESS_KEY_ID',
        config_value: state.aliyun.accessKeyId,
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_ALIYUN_ACCESS_KEY_SECRET',
        config_value: state.aliyun.accessKeySecret,
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_ALIYUN_ACCOUNT_NAME',
        config_value: state.aliyun.accountName,
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_TENCENT_SECRET_ID',
        config_value: state.tencent.secretId,
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_TENCENT_SECRET_KEY',
        config_value: state.tencent.secretKey,
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_TENCENT_FROM_EMAIL',
        config_value: state.tencent.fromEmail,
        category: CATEGORY,
      },
      {
        config_key: 'MAIL_TENCENT_REGION',
        config_value: state.tencent.region,
        category: CATEGORY,
      },
    ])
    window.$message.success('保存成功')
    state.snapshot = JSON.stringify({
      defaultEngine: state.defaultEngine,
      local: state.local,
      aliyun: state.aliyun,
      tencent: state.tencent,
    })
  } finally {
    state.saving = false
  }
}
</script>

<template>
  <NSpin :show="state.loading">
    <NForm
      class="sys-config-form mb-12px"
      label-placement="top"
    >
      <NFormItem label="默认邮件引擎">
        <NRadioGroup v-model:value="state.defaultEngine">
          <NSpace>
            <NRadio
              v-for="opt in engineOptions"
              :key="opt.value"
              :value="opt.value"
              :label="opt.label"
            />
          </NSpace>
        </NRadioGroup>
      </NFormItem>
    </NForm>

    <NTabs
      v-model:value="state.subTab"
      type="line"
      class="sys-config-subnav"
    >
      <NTab
        name="LOCAL"
        tab="本地邮件"
      />
      <NTab
        name="ALIYUN"
        tab="阿里云邮件"
      />
      <NTab
        name="TENCENT"
        tab="腾讯云邮件"
      />
    </NTabs>

    <ConfigSectionLayout
      :description="
        state.subTab === 'LOCAL'
          ? '配置本地 SMTP 发信。'
          : state.subTab === 'ALIYUN'
            ? '配置阿里云邮件推送凭证。'
            : '配置腾讯云邮件推送凭证。'
      "
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NForm
        v-if="state.subTab === 'LOCAL'"
        class="sys-config-form sys-config-form--wide"
        label-placement="top"
      >
        <NGrid
          :cols="24"
          :x-gap="16"
        >
          <NGi :span="16">
            <NFormItem label="SMTP 服务器">
              <NInput
                v-model:value="state.local.host"
                placeholder="smtp.example.com"
              />
            </NFormItem>
          </NGi>
          <NGi :span="8">
            <NFormItem label="端口">
              <NInputNumber
                v-model:value="state.local.port"
                class="w-full"
                :min="1"
              />
            </NFormItem>
          </NGi>
          <NGi :span="12">
            <NFormItem label="发送账号">
              <NInput v-model:value="state.local.username" />
            </NFormItem>
          </NGi>
          <NGi :span="12">
            <NFormItem label="邮箱密钥">
              <NInput
                v-model:value="state.local.password"
                type="password"
                show-password-on="click"
                placeholder="留空不修改"
              />
            </NFormItem>
          </NGi>
          <NGi :span="12">
            <NFormItem label="发件人邮箱">
              <NInput v-model:value="state.local.fromEmail" />
            </NFormItem>
          </NGi>
          <NGi :span="12">
            <NFormItem label="发件人名称">
              <NInput v-model:value="state.local.fromName" />
            </NFormItem>
          </NGi>
          <NGi :span="6">
            <NFormItem label="账号密码验证">
              <NSwitch v-model:value="state.local.authRequired" />
            </NFormItem>
          </NGi>
          <NGi :span="6">
            <NFormItem label="SSL">
              <NSwitch v-model:value="state.local.useSsl" />
            </NFormItem>
          </NGi>
          <NGi :span="6">
            <NFormItem label="STARTTLS">
              <NSwitch v-model:value="state.local.useStarttls" />
            </NFormItem>
          </NGi>
        </NGrid>
      </NForm>

      <NForm
        v-else-if="state.subTab === 'ALIYUN'"
        class="sys-config-form sys-config-form--wide"
        label-placement="top"
      >
        <NGrid
          :cols="24"
          :x-gap="16"
        >
          <NGi :span="12">
            <NFormItem label="密钥 ID">
              <NInput v-model:value="state.aliyun.accessKeyId" />
            </NFormItem>
          </NGi>
          <NGi :span="12">
            <NFormItem label="密钥">
              <NInput
                v-model:value="state.aliyun.accessKeySecret"
                type="password"
                show-password-on="click"
                placeholder="留空不修改"
              />
            </NFormItem>
          </NGi>
          <NGi :span="24">
            <NFormItem label="发信地址">
              <NInput v-model:value="state.aliyun.accountName" />
            </NFormItem>
          </NGi>
        </NGrid>
      </NForm>

      <NForm
        v-else
        class="sys-config-form sys-config-form--wide"
        label-placement="top"
      >
        <NGrid
          :cols="24"
          :x-gap="16"
        >
          <NGi :span="12">
            <NFormItem label="密钥 ID">
              <NInput v-model:value="state.tencent.secretId" />
            </NFormItem>
          </NGi>
          <NGi :span="12">
            <NFormItem label="密钥">
              <NInput
                v-model:value="state.tencent.secretKey"
                type="password"
                show-password-on="click"
                placeholder="留空不修改"
              />
            </NFormItem>
          </NGi>
          <NGi :span="12">
            <NFormItem label="发件邮箱">
              <NInput v-model:value="state.tencent.fromEmail" />
            </NFormItem>
          </NGi>
          <NGi :span="12">
            <NFormItem label="区域（留空默认 ap-guangzhou）">
              <NInput v-model:value="state.tencent.region" placeholder="ap-guangzhou" />
            </NFormItem>
          </NGi>
        </NGrid>
      </NForm>
    </ConfigSectionLayout>
  </NSpin>
</template>
