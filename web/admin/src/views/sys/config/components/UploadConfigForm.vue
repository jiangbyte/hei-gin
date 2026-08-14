<!-- Author: Charlie -->

<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import {
  loadByCategory,
  parseNumber,
  parseTags,
  saveByKeys,
  tagsToJson,
} from '../composables/useConfigForm'

const CATEGORY = 'UPLOAD'

const state = reactive({
  loading: false,
  saving: false,
  maxBytes: 10485760,
  presignExpire: 3600,
  categoryMaxLength: 64,
  allowedContentTypes: [] as string[],
  allowedExtensions: [] as string[],
  deniedExtensions: [] as string[],
  snapshot: '',
})

onMounted(() => {
  void reload()
})

async function reload() {
  state.loading = true
  try {
    const map = await loadByCategory(CATEGORY)
    state.maxBytes = parseNumber(map.STORAGE_UPLOAD_MAX_BYTES, 10485760)
    state.presignExpire = parseNumber(map.STORAGE_PRESIGN_EXPIRE_SECONDS, 3600)
    state.categoryMaxLength = parseNumber(map.STORAGE_UPLOAD_CATEGORY_MAX_LENGTH, 64)
    state.allowedContentTypes = parseTags(map.STORAGE_UPLOAD_ALLOWED_CONTENT_TYPES)
    state.allowedExtensions = parseTags(map.STORAGE_UPLOAD_ALLOWED_EXTENSIONS)
    state.deniedExtensions = parseTags(map.STORAGE_UPLOAD_DENIED_EXTENSIONS)
    state.snapshot = JSON.stringify({
      maxBytes: state.maxBytes,
      presignExpire: state.presignExpire,
      categoryMaxLength: state.categoryMaxLength,
      allowedContentTypes: state.allowedContentTypes,
      allowedExtensions: state.allowedExtensions,
      deniedExtensions: state.deniedExtensions,
    })
  } finally {
    state.loading = false
  }
}

function reset() {
  if (!state.snapshot) return
  const data = JSON.parse(state.snapshot)
  state.maxBytes = data.maxBytes
  state.presignExpire = data.presignExpire
  state.categoryMaxLength = data.categoryMaxLength
  state.allowedContentTypes = data.allowedContentTypes
  state.allowedExtensions = data.allowedExtensions
  state.deniedExtensions = data.deniedExtensions
}

async function save() {
  state.saving = true
  try {
    await saveByKeys([
      {
        config_key: 'STORAGE_UPLOAD_MAX_BYTES',
        config_value: String(state.maxBytes),
        category: CATEGORY,
      },
      {
        config_key: 'STORAGE_PRESIGN_EXPIRE_SECONDS',
        config_value: String(state.presignExpire),
        category: CATEGORY,
      },
      {
        config_key: 'STORAGE_UPLOAD_CATEGORY_MAX_LENGTH',
        config_value: String(state.categoryMaxLength),
        category: CATEGORY,
      },
      {
        config_key: 'STORAGE_UPLOAD_ALLOWED_CONTENT_TYPES',
        config_value: tagsToJson(state.allowedContentTypes),
        category: CATEGORY,
      },
      {
        config_key: 'STORAGE_UPLOAD_ALLOWED_EXTENSIONS',
        config_value: tagsToJson(state.allowedExtensions),
        category: CATEGORY,
      },
      {
        config_key: 'STORAGE_UPLOAD_DENIED_EXTENSIONS',
        config_value: tagsToJson(state.deniedExtensions),
        category: CATEGORY,
      },
    ])
    window.$message.success('保存成功')
    state.snapshot = JSON.stringify({
      maxBytes: state.maxBytes,
      presignExpire: state.presignExpire,
      categoryMaxLength: state.categoryMaxLength,
      allowedContentTypes: state.allowedContentTypes,
      allowedExtensions: state.allowedExtensions,
      deniedExtensions: state.deniedExtensions,
    })
  } finally {
    state.saving = false
  }
}
</script>

<template>
  <NSpin :show="state.loading">
    <ConfigSectionLayout
      description="配置上传大小、预签名有效期、分类名长度与扩展名限制。分类名长度作用于对象路径前缀。"
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NForm
        class="sys-config-form sys-config-form--wide"
        label-placement="top"
      >
        <NGrid
          :cols="24"
          :x-gap="16"
        >
          <NGi :span="8">
            <NFormItem label="上传大小上限（字节）">
              <NInputNumber
                v-model:value="state.maxBytes"
                class="w-full"
                :min="0"
              />
            </NFormItem>
          </NGi>
          <NGi :span="8">
            <NFormItem label="预签名有效期（秒）">
              <NInputNumber
                v-model:value="state.presignExpire"
                class="w-full"
                :min="0"
              />
            </NFormItem>
          </NGi>
          <NGi :span="8">
            <NFormItem label="分类名最大长度">
              <NInputNumber
                v-model:value="state.categoryMaxLength"
                class="w-full"
                :min="1"
              />
            </NFormItem>
          </NGi>
          <NGi :span="24">
            <NFormItem label="允许的 MIME 类型">
              <NDynamicTags
                v-model:value="state.allowedContentTypes"
                :max="100"
              />
            </NFormItem>
          </NGi>
          <NGi :span="24">
            <NFormItem label="允许的扩展名">
              <NDynamicTags
                v-model:value="state.allowedExtensions"
                :max="100"
              />
            </NFormItem>
          </NGi>
          <NGi :span="24">
            <NFormItem label="禁止的扩展名">
              <NDynamicTags
                v-model:value="state.deniedExtensions"
                :max="100"
              />
            </NFormItem>
          </NGi>
        </NGrid>
      </NForm>
    </ConfigSectionLayout>
  </NSpin>
</template>
