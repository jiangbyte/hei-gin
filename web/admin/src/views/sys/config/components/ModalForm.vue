<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { configApi } from '@/api'
import { createRequiredRule, toNullableString } from '@/utils'
import { computed, reactive, ref } from 'vue'

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)

const defaultFormData = {
  config_key: '',
  config_value: '',
  category: 'OTHER',
  remark: '',
  sort_code: 0,
  ext_json: '{}',
}
const state = reactive({
  showModal: false,
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: { ...defaultFormData },
})

const modalTitle = computed(() => (state.dataId ? '编辑系统配置' : '新增系统配置'))

const KEY_PATTERN = /^[A-Z][A-Z0-9_]*$/

const rules = computed<FormRules>(() => ({
  config_key: [
    createRequiredRule('配置键', 'input'),
    {
      validator: (_rule, value: string) => KEY_PATTERN.test(String(value || '').trim()),
      message: '配置键须匹配 ^[A-Z][A-Z0-9_]*$',
      trigger: ['input', 'blur'],
    },
  ],
  ext_json: [
    {
      validator: isValidExtJson,
      message: '请输入合法 JSON 对象',
      trigger: ['input', 'blur'],
    },
  ],
}))

async function openModal(id?: string) {
  state.dataId = id ?? null
  state.formModel = { ...defaultFormData }
  state.showModal = true

  if (id) {
    await fetchDetail(id)
  }
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await configApi.detail({ id })
    const data = response.data ?? {}
    state.formModel = Object.assign({}, defaultFormData, data, {
      config_value: data.config_value ?? '',
      category: data.category ?? 'OTHER',
      remark: data.remark ?? '',
      sort_code: data.sort_code ?? 0,
      ext_json: stringifyExtJson(data.ext_json),
    })
  } finally {
    state.loading = false
  }
}

function closeModal() {
  state.showModal = false
  state.submitLoading = false
}

async function submitForm() {
  await formRef.value?.validate()
  state.submitLoading = true
  try {
    const payload = {
      ...state.formModel,
      config_key: state.formModel.config_key.trim(),
      config_value: toNullableString(state.formModel.config_value),
      category: toNullableString(state.formModel.category),
      remark: toNullableString(state.formModel.remark),
      sort_code: Number(state.formModel.sort_code ?? 0),
      ext_json: parseExtJson(),
    }

    if (state.dataId) {
      await configApi.update({
        ...payload,
        id: state.dataId,
      })
      window.$message.success('更新成功')
    } else {
      await configApi.create(payload)
      window.$message.success('创建成功')
    }

    closeModal()
    emit('saved')
  } finally {
    state.submitLoading = false
  }
}

function parseExtJson() {
  const text = String(state.formModel.ext_json || '').trim()
  if (!text) {
    return {}
  }
  const value = JSON.parse(text)
  if (Array.isArray(value) || typeof value !== 'object' || value === null) {
    throw new Error('ext_json must be an object')
  }
  return value
}

function isValidExtJson() {
  try {
    parseExtJson()
    return true
  } catch {
    return false
  }
}

function stringifyExtJson(value: unknown) {
  if (!value || typeof value !== 'object') {
    return '{}'
  }
  return JSON.stringify(value, null, 2)
}

defineExpose({
  openModal,
})
</script>

<template>
  <NModal
    v-model:show="state.showModal"
    preset="card"
    draggable
    :mask-closable="false"
    :title="modalTitle"
    style="width: 720px"
    :segmented="{ content: true, action: true }"
  >
    <NSpin :show="state.loading">
      <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
        <NForm
          ref="formRef"
          :model="state.formModel"
          :rules="rules"
          label-placement="left"
          label-width="110"
          :disabled="state.loading || state.submitLoading"
        >
          <NFormItem
            :label="'配置键'"
            path="config_key"
          >
            <NInput
              v-model:value="state.formModel.config_key"
              placeholder="如 CUSTOM_FEATURE_FLAG"
              :disabled="!!state.dataId"
            />
          </NFormItem>
          <NFormItem
            :label="'配置值'"
            path="config_value"
          >
            <NInput
              v-model:value="state.formModel.config_value"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 8 }"
            />
          </NFormItem>

          <NFormItem
            :label="'备注'"
            path="remark"
          >
            <NInput v-model:value="state.formModel.remark" />
          </NFormItem>
          <NFormItem
            :label="'排序码'"
            path="sort_code"
          >
            <NInputNumber
              v-model:value="state.formModel.sort_code"
              class="w-full"
              :min="0"
            />
          </NFormItem>
          <NFormItem
            :label="'扩展信息'"
            path="ext_json"
          >
            <NInput
              v-model:value="state.formModel.ext_json"
              type="textarea"
              :autosize="{ minRows: 4, maxRows: 10 }"
            />
          </NFormItem>
        </NForm>
      </NScrollbar>
    </NSpin>

    <template #action>
      <NSpace
        justify="end"
        align="center"
      >
        <NButton @click="closeModal">
          取消
        </NButton>
        <NButton
          type="primary"
          :loading="state.submitLoading"
          @click="submitForm"
        >
          确认
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
