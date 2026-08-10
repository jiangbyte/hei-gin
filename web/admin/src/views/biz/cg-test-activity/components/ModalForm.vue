<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-08 21:09:52
-->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { cgTestActivityApi } from '@/api'
import { wireBool, wireInt, wireFloat } from '@/utils/wire'
import { createRequiredRule, toApiDateTime, toFormDateTime } from '@/utils'
import { computed, reactive, ref } from 'vue'

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)
const defaultFormData: Record<string, any> = {
  code: '',
  name: '',
  category: '',
  type: '',
  status: '',
  cover_url: '',
  description: '',
  start_at: null,
  end_at: null,
  max_participants: 0,
  price: 0,
  is_public: false,
  need_approval: false,
  rule_config: '{}',
  extra: '{}',
}
const state = reactive({
  showModal: false,
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: normalizeFormData(),
})

const modalTitle = computed(() => (state.dataId ? '编辑CgTestActivity' : '新增CgTestActivity'))
const rules = computed<FormRules>(() => ({
  code: [createRequiredRule('活动编码', 'input')],
  name: [createRequiredRule('活动名称', 'input')],
  type: [createRequiredRule('活动类型', 'input')],
  status: [createRequiredRule('状态', 'input')],
  start_at: [createRequiredRule('开始时间', 'change')],
  max_participants: [
    {
      validator: () =>
        typeof state.formModel.max_participants === 'number' &&
        Number.isFinite(state.formModel.max_participants),
      message: '请输入最大参与人数',
      trigger: ['input', 'blur'],
    },
  ],
  price: [
    {
      validator: () =>
        typeof state.formModel.price === 'number' && Number.isFinite(state.formModel.price),
      message: '请输入报名费用',
      trigger: ['input', 'blur'],
    },
  ],
  is_public: [
    {
      validator: () => typeof state.formModel.is_public === 'boolean',
      message: '请选择是否公开',
      trigger: 'change',
    },
  ],
  need_approval: [
    {
      validator: () => typeof state.formModel.need_approval === 'boolean',
      message: '请选择是否需要审批',
      trigger: 'change',
    },
  ],
  rule_config: [
    createRequiredRule('规则配置', 'input'),
    {
      validator: () => isValidJsonValue(state.formModel.rule_config),
      message: '请输入合法 JSON 对象',
      trigger: ['input', 'blur'],
    },
  ],
  extra: [
    {
      validator: () => isValidJsonValue(state.formModel.extra),
      message: '请输入合法 JSON 对象',
      trigger: ['input', 'blur'],
    },
  ],
}))

async function openModal(id?: string, defaults: Partial<typeof defaultFormData> = {}) {
  state.dataId = id ?? null
  state.formModel = normalizeFormData(defaults)
  state.showModal = true
  if (id) {
    await fetchDetail(id)
  }
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await cgTestActivityApi.detail({ id })
    state.formModel = normalizeFormData(response.data ?? {})
  } finally {
    state.loading = false
  }
}

function normalizeFormData(data: Record<string, any> = {}): Record<string, any> {
  return {
    ...defaultFormData,
    ...data,
    is_public:
      data.is_public == null || data.is_public === ''
        ? defaultFormData.is_public
        : wireBool(String(data.is_public)),
    need_approval:
      data.need_approval == null || data.need_approval === ''
        ? defaultFormData.need_approval
        : wireBool(String(data.need_approval)),
    max_participants:
      data.max_participants == null || data.max_participants === ''
        ? defaultFormData.max_participants
        : wireInt(String(data.max_participants)),
    price:
      data.price == null || data.price === ''
        ? defaultFormData.price
        : wireFloat(String(data.price)),
    start_at: toFormDateTime(data.start_at),
    end_at: toFormDateTime(data.end_at),
    rule_config: stringifyJsonValue(data.rule_config),
    extra: stringifyJsonValue(data.extra),
  }
}

function normalizeSubmitData(data: Record<string, any>): Record<string, any> {
  return {
    ...data,
    start_at: toApiDateTime(data.start_at),
    end_at: toApiDateTime(data.end_at),
    rule_config: parseJsonValue(data.rule_config),
    extra: parseJsonValue(data.extra),
  }
}

function parseJsonValue(value: unknown) {
  const text = String(value ?? '').trim()
  if (!text) {
    return {}
  }
  const parsed = JSON.parse(text)
  if (Array.isArray(parsed) || typeof parsed !== 'object' || parsed === null) {
    throw new Error('JSON value must be an object')
  }
  return parsed
}

function isValidJsonValue(value: unknown) {
  try {
    parseJsonValue(value)
    return true
  } catch {
    return false
  }
}

function stringifyJsonValue(value: unknown) {
  if (value === undefined || value === null || value === '') {
    return '{}'
  }
  if (typeof value === 'string') {
    try {
      return JSON.stringify(JSON.parse(value), null, 2)
    } catch {
      return value
    }
  }
  return JSON.stringify(value, null, 2)
}

function closeModal() {
  state.showModal = false
  state.submitLoading = false
}

async function submitForm() {
  await formRef.value?.validate()
  state.submitLoading = true
  try {
    const payload = normalizeSubmitData(state.formModel)
    if (state.dataId) {
      await cgTestActivityApi.update({ ...payload, id: state.dataId })
      window.$message.success('更新成功')
    } else {
      await cgTestActivityApi.create(payload)
      window.$message.success('创建成功')
    }
    emit('saved')
    closeModal()
  } finally {
    state.submitLoading = false
  }
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
            label="活动编码"
            path="code"
          >
            <NInput v-model:value="state.formModel.code" />
          </NFormItem>
          <NFormItem
            label="活动名称"
            path="name"
          >
            <NInput v-model:value="state.formModel.name" />
          </NFormItem>
          <NFormItem
            label="活动分类"
            path="category"
          >
            <NInput v-model:value="state.formModel.category" />
          </NFormItem>
          <NFormItem
            label="活动类型"
            path="type"
          >
            <NInput v-model:value="state.formModel.type" />
          </NFormItem>
          <NFormItem
            label="状态"
            path="status"
          >
            <NInput v-model:value="state.formModel.status" />
          </NFormItem>
          <NFormItem
            label="封面地址"
            path="cover_url"
          >
            <NInput v-model:value="state.formModel.cover_url" />
          </NFormItem>
          <NFormItem
            label="活动描述"
            path="description"
          >
            <NInput
              v-model:value="state.formModel.description"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 8 }"
            />
          </NFormItem>
          <NFormItem
            label="开始时间"
            path="start_at"
          >
            <NDatePicker
              v-model:formatted-value="state.formModel.start_at"
              type="datetime"
              value-format="yyyy-MM-dd HH:mm:ss"
              class="w-full"
              clearable
            />
          </NFormItem>
          <NFormItem
            label="结束时间"
            path="end_at"
          >
            <NDatePicker
              v-model:formatted-value="state.formModel.end_at"
              type="datetime"
              value-format="yyyy-MM-dd HH:mm:ss"
              class="w-full"
              clearable
            />
          </NFormItem>
          <NFormItem
            label="最大参与人数"
            path="max_participants"
          >
            <NInputNumber
              v-model:value="state.formModel.max_participants"
              class="w-full"
            />
          </NFormItem>
          <NFormItem
            label="报名费用"
            path="price"
          >
            <NInputNumber
              v-model:value="state.formModel.price"
              class="w-full"
            />
          </NFormItem>
          <NFormItem
            label="是否公开"
            path="is_public"
          >
            <NSwitch v-model:value="state.formModel.is_public" />
          </NFormItem>
          <NFormItem
            label="是否需要审批"
            path="need_approval"
          >
            <NSwitch v-model:value="state.formModel.need_approval" />
          </NFormItem>
          <NFormItem
            label="规则配置"
            path="rule_config"
          >
            <NInput
              v-model:value="state.formModel.rule_config"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 8 }"
            />
          </NFormItem>
          <NFormItem
            label="扩展信息"
            path="extra"
          >
            <NInput
              v-model:value="state.formModel.extra"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 8 }"
            />
          </NFormItem>
        </NForm>
      </NScrollbar>
    </NSpin>

    <template #action>
      <NSpace justify="end">
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
