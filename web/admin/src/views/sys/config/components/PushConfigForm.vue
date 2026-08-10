<!-- Author: Charlie -->

<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import { loadByCategory, saveByKeys } from '../composables/useConfigForm'

const CATEGORY = 'PUSH'
type Channel = 'DINGTALK' | 'LARK' | 'WECHAT_WORK'

const engineOptions = [
  { label: '钉钉', value: 'DINGTALK' },
  { label: '飞书', value: 'LARK' },
  { label: '企业微信', value: 'WECHAT_WORK' },
]

const state = reactive({
  loading: false,
  saving: false,
  subTab: 'DINGTALK' as Channel,
  defaultEngine: 'DINGTALK' as Channel,
  dingtalk: { webhook: '', secret: '' },
  lark: { webhook: '', secret: '' },
  wechatWork: { webhook: '' },
  snapshot: '',
})

onMounted(() => {
  void reload()
})

async function reload() {
  state.loading = true
  try {
    const [map, sys] = await Promise.all([loadByCategory(CATEGORY), loadByCategory('SYS')])
    state.defaultEngine = (map.DEFAULT_MESSAGE_PUSH_ENGINE ||
      sys.DEFAULT_MESSAGE_PUSH_ENGINE ||
      'DINGTALK') as Channel
    state.dingtalk.webhook = map.PUSH_DINGTALK_WEBHOOK || ''
    state.dingtalk.secret = map.PUSH_DINGTALK_SECRET || ''
    state.lark.webhook = map.PUSH_LARK_WEBHOOK || ''
    state.lark.secret = map.PUSH_LARK_SECRET || ''
    state.wechatWork.webhook = map.PUSH_WECHAT_WORK_WEBHOOK || ''
    state.snapshot = JSON.stringify({
      defaultEngine: state.defaultEngine,
      dingtalk: state.dingtalk,
      lark: state.lark,
      wechatWork: state.wechatWork,
    })
  } finally {
    state.loading = false
  }
}

function reset() {
  if (!state.snapshot) return
  const data = JSON.parse(state.snapshot)
  state.defaultEngine = data.defaultEngine
  Object.assign(state.dingtalk, data.dingtalk)
  Object.assign(state.lark, data.lark)
  Object.assign(state.wechatWork, data.wechatWork)
}

async function save() {
  state.saving = true
  try {
    await saveByKeys([
      {
        config_key: 'DEFAULT_MESSAGE_PUSH_ENGINE',
        config_value: state.defaultEngine,
        category: CATEGORY,
      },
      {
        config_key: 'PUSH_DINGTALK_WEBHOOK',
        config_value: state.dingtalk.webhook,
        category: CATEGORY,
      },
      {
        config_key: 'PUSH_DINGTALK_SECRET',
        config_value: state.dingtalk.secret,
        category: CATEGORY,
      },
      {
        config_key: 'PUSH_LARK_WEBHOOK',
        config_value: state.lark.webhook,
        category: CATEGORY,
      },
      {
        config_key: 'PUSH_LARK_SECRET',
        config_value: state.lark.secret,
        category: CATEGORY,
      },
      {
        config_key: 'PUSH_WECHAT_WORK_WEBHOOK',
        config_value: state.wechatWork.webhook,
        category: CATEGORY,
      },
    ])
    window.$message.success('保存成功')
    state.snapshot = JSON.stringify({
      defaultEngine: state.defaultEngine,
      dingtalk: state.dingtalk,
      lark: state.lark,
      wechatWork: state.wechatWork,
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
      <NFormItem label="默认消息推送引擎">
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
        name="DINGTALK"
        tab="钉钉"
      />
      <NTab
        name="LARK"
        tab="飞书"
      />
      <NTab
        name="WECHAT_WORK"
        tab="企业微信"
      />
    </NTabs>

    <ConfigSectionLayout
      :description="
        state.subTab === 'DINGTALK'
          ? '配置钉钉机器人 Webhook。'
          : state.subTab === 'LARK'
            ? '配置飞书机器人 Webhook。'
            : '配置企业微信机器人 Webhook。'
      "
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NForm
        v-if="state.subTab === 'DINGTALK'"
        class="sys-config-form"
        label-placement="top"
      >
        <NFormItem label="Webhook 地址">
          <NInput v-model:value="state.dingtalk.webhook" />
        </NFormItem>
        <NFormItem label="加签密钥（空值表示不修改）">
          <NInput
            v-model:value="state.dingtalk.secret"
            type="password"
            show-password-on="click"
          />
        </NFormItem>
      </NForm>

      <NForm
        v-else-if="state.subTab === 'LARK'"
        class="sys-config-form"
        label-placement="top"
      >
        <NFormItem label="Webhook 地址">
          <NInput v-model:value="state.lark.webhook" />
        </NFormItem>
        <NFormItem label="密钥（空值表示不修改）">
          <NInput
            v-model:value="state.lark.secret"
            type="password"
            show-password-on="click"
          />
        </NFormItem>
      </NForm>

      <NForm
        v-else
        class="sys-config-form"
        label-placement="top"
      >
        <NFormItem label="Webhook 地址">
          <NInput v-model:value="state.wechatWork.webhook" />
        </NFormItem>
      </NForm>
    </ConfigSectionLayout>
  </NSpin>
</template>
