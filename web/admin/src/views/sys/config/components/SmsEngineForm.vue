<!-- Author: Charlie -->

<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import { loadByCategory, saveByKeys } from '../composables/useConfigForm'

const CATEGORY = 'SMS'
type Engine = 'ALIYUN' | 'TENCENT'

const engineOptions = [
  { label: '阿里云短信', value: 'ALIYUN' },
  { label: '腾讯云短信', value: 'TENCENT' },
]

const state = reactive({
  loading: false,
  saving: false,
  subTab: 'ALIYUN' as Engine,
  defaultEngine: 'ALIYUN' as Engine,
  aliyun: {
    accessKeyId: '',
    accessKeySecret: '',
    signName: '',
  },
  tencent: {
    secretId: '',
    secretKey: '',
    sdkAppId: '',
    signName: '',
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
    state.defaultEngine = (map.DEFAULT_SMS_ENGINE || sys.DEFAULT_SMS_ENGINE || 'ALIYUN') as Engine
    state.aliyun.accessKeyId = map.SMS_ALIYUN_ACCESS_KEY_ID || ''
    state.aliyun.accessKeySecret = map.SMS_ALIYUN_ACCESS_KEY_SECRET || ''
    state.aliyun.signName = map.SMS_ALIYUN_SIGN_NAME || ''
    state.tencent.secretId = map.SMS_TENCENT_SECRET_ID || ''
    state.tencent.secretKey = map.SMS_TENCENT_SECRET_KEY || ''
    state.tencent.sdkAppId = map.SMS_TENCENT_SDK_APP_ID || ''
    state.tencent.signName = map.SMS_TENCENT_SIGN_NAME || ''
    state.tencent.region = map.SMS_TENCENT_REGION || ''
    state.snapshot = JSON.stringify({
      defaultEngine: state.defaultEngine,
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
  Object.assign(state.aliyun, data.aliyun)
  Object.assign(state.tencent, data.tencent)
}

async function save() {
  state.saving = true
  try {
    await saveByKeys([
      {
        config_key: 'DEFAULT_SMS_ENGINE',
        config_value: state.defaultEngine,
        category: CATEGORY,
      },
      {
        config_key: 'SMS_ALIYUN_ACCESS_KEY_ID',
        config_value: state.aliyun.accessKeyId,
        category: CATEGORY,
      },
      {
        config_key: 'SMS_ALIYUN_ACCESS_KEY_SECRET',
        config_value: state.aliyun.accessKeySecret,
        category: CATEGORY,
      },
      {
        config_key: 'SMS_ALIYUN_SIGN_NAME',
        config_value: state.aliyun.signName,
        category: CATEGORY,
      },
      {
        config_key: 'SMS_TENCENT_SECRET_ID',
        config_value: state.tencent.secretId,
        category: CATEGORY,
      },
      {
        config_key: 'SMS_TENCENT_SECRET_KEY',
        config_value: state.tencent.secretKey,
        category: CATEGORY,
      },
      {
        config_key: 'SMS_TENCENT_SDK_APP_ID',
        config_value: state.tencent.sdkAppId,
        category: CATEGORY,
      },
      {
        config_key: 'SMS_TENCENT_SIGN_NAME',
        config_value: state.tencent.signName,
        category: CATEGORY,
      },
      {
        config_key: 'SMS_TENCENT_REGION',
        config_value: state.tencent.region,
        category: CATEGORY,
      },
    ])
    window.$message.success('保存成功')
    state.snapshot = JSON.stringify({
      defaultEngine: state.defaultEngine,
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
      <NFormItem label="默认短信引擎">
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
        name="ALIYUN"
        tab="阿里云短信"
      />
      <NTab
        name="TENCENT"
        tab="腾讯云短信"
      />
    </NTabs>

    <ConfigSectionLayout
      :description="state.subTab === 'ALIYUN' ? '配置阿里云短信凭证。' : '配置腾讯云短信凭证。'"
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NForm
        v-if="state.subTab === 'ALIYUN'"
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
            <NFormItem label="短信签名">
              <NInput v-model:value="state.aliyun.signName" />
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
            <NFormItem label="应用编号">
              <NInput v-model:value="state.tencent.sdkAppId" />
            </NFormItem>
          </NGi>
          <NGi :span="12">
            <NFormItem label="短信签名">
              <NInput v-model:value="state.tencent.signName" />
            </NFormItem>
          </NGi>
          <NGi :span="24">
            <NFormItem label="区域（留空默认 ap-guangzhou）">
              <NInput v-model:value="state.tencent.region" placeholder="ap-guangzhou" />
            </NFormItem>
          </NGi>
        </NGrid>
      </NForm>
    </ConfigSectionLayout>
  </NSpin>
</template>
