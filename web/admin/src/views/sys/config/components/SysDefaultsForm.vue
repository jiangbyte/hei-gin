<!-- Author: Charlie -->

<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import { loadByCategory, saveByKeys } from '../composables/useConfigForm'

const CATEGORY = 'SYS'

const state = reactive({
  loading: false,
  saving: false,
  copyrightText: '',
  copyrightUrl: '',
  snapshot: '',
})

onMounted(() => {
  void reload()
})

async function reload() {
  state.loading = true
  try {
    const map = await loadByCategory(CATEGORY)
    state.copyrightText = map.COPYRIGHT_TEXT || ''
    state.copyrightUrl = map.COPYRIGHT_URL || ''
    state.snapshot = JSON.stringify({
      copyrightText: state.copyrightText,
      copyrightUrl: state.copyrightUrl,
    })
  } finally {
    state.loading = false
  }
}

function reset() {
  if (!state.snapshot) return
  const data = JSON.parse(state.snapshot)
  state.copyrightText = data.copyrightText
  state.copyrightUrl = data.copyrightUrl
}

async function save() {
  state.saving = true
  try {
    await saveByKeys([
      {
        config_key: 'COPYRIGHT_TEXT',
        config_value: state.copyrightText,
        category: CATEGORY,
      },
      {
        config_key: 'COPYRIGHT_URL',
        config_value: state.copyrightUrl,
        category: CATEGORY,
      },
    ])
    window.$message.success('保存成功')
    state.snapshot = JSON.stringify({
      copyrightText: state.copyrightText,
      copyrightUrl: state.copyrightUrl,
    })
  } finally {
    state.saving = false
  }
}
</script>

<template>
  <NSpin :show="state.loading">
    <ConfigSectionLayout
      description="配置站点版权文案与链接。"
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NForm
        class="sys-config-form"
        label-placement="top"
      >
        <NFormItem label="版权文案">
          <NInput v-model:value="state.copyrightText" />
        </NFormItem>
        <NFormItem label="版权链接">
          <NInput
            v-model:value="state.copyrightUrl"
            placeholder="https://"
          />
        </NFormItem>
      </NForm>
    </ConfigSectionLayout>
  </NSpin>
</template>
