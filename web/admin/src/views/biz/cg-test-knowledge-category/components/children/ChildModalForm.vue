<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-08 21:09:55
-->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { cgTestKnowledgeCategoryApi } from '@/api'
import { wireBool, wireInt } from '@/utils/wire'
import { createRequiredRule, toApiDateTime, toFormDateTime } from '@/utils'
import { computed, reactive, ref } from 'vue'

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)
const defaultFormData: Record<string, any> = {
  category_id: '',
  code: '',
  title: '',
  type: '',
  status: '',
  summary: '',
  content: '',
  author: '',
  published_at: null,
  view_count: 0,
  sort: 0,
  is_top: false,
  settings: '{}',
  extra: '{}',
}
const state = reactive({
  showModal: false,
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: normalizeFormData(),
})

const modalTitle = computed(() =>
  state.dataId ? '编辑CgTestKnowledgeDoc' : '新增CgTestKnowledgeDoc',
)
const rules = computed<FormRules>(() => ({
  category_id: [createRequiredRule('分类ID', 'input')],
  code: [createRequiredRule('文档编码', 'input')],
  title: [createRequiredRule('文档标题', 'input')],
  type: [createRequiredRule('文档类型', 'input')],
  status: [createRequiredRule('状态', 'input')],
  view_count: [
    {
      validator: () =>
        typeof state.formModel.view_count === 'number' &&
        Number.isFinite(state.formModel.view_count),
      message: '请输入浏览次数',
      trigger: ['input', 'blur'],
    },
  ],
  sort: [
    {
      validator: () =>
        typeof state.formModel.sort === 'number' && Number.isFinite(state.formModel.sort),
      message: '请输入排序',
      trigger: ['input', 'blur'],
    },
  ],
  is_top: [
    {
      validator: () => typeof state.formModel.is_top === 'boolean',
      message: '请选择是否置顶',
      trigger: 'change',
    },
  ],
  settings: [
    createRequiredRule('展示设置', 'input'),
    {
      validator: () => isValidJsonValue(state.formModel.settings),
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
    const response = await cgTestKnowledgeCategoryApi.childDetail({ id })
    state.formModel = normalizeFormData(response.data ?? {})
  } finally {
    state.loading = false
  }
}

function normalizeFormData(data: Record<string, any> = {}): Record<string, any> {
  return {
    ...defaultFormData,
    ...data,
    is_top:
      data.is_top == null || data.is_top === ''
        ? defaultFormData.is_top
        : wireBool(String(data.is_top)),
    view_count:
      data.view_count == null || data.view_count === ''
        ? defaultFormData.view_count
        : wireInt(String(data.view_count)),
    sort: data.sort == null || data.sort === '' ? defaultFormData.sort : wireInt(String(data.sort)),
    published_at: toFormDateTime(data.published_at),
    settings: stringifyJsonValue(data.settings),
    extra: stringifyJsonValue(data.extra),
  }
}

function normalizeSubmitData(data: Record<string, any>): Record<string, any> {
  return {
    ...data,
    published_at: toApiDateTime(data.published_at),
    settings: parseJsonValue(data.settings),
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
      await cgTestKnowledgeCategoryApi.childUpdate({ ...payload, id: state.dataId })
      window.$message.success('更新成功')
    } else {
      await cgTestKnowledgeCategoryApi.childCreate(payload)
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
            label="分类ID"
            path="category_id"
          >
            <NInput v-model:value="state.formModel.category_id" />
          </NFormItem>
          <NFormItem
            label="文档编码"
            path="code"
          >
            <NInput v-model:value="state.formModel.code" />
          </NFormItem>
          <NFormItem
            label="文档标题"
            path="title"
          >
            <NInput v-model:value="state.formModel.title" />
          </NFormItem>
          <NFormItem
            label="文档类型"
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
            label="摘要"
            path="summary"
          >
            <NInput v-model:value="state.formModel.summary" />
          </NFormItem>
          <NFormItem
            label="正文内容"
            path="content"
          >
            <NInput
              v-model:value="state.formModel.content"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 8 }"
            />
          </NFormItem>
          <NFormItem
            label="作者"
            path="author"
          >
            <NInput v-model:value="state.formModel.author" />
          </NFormItem>
          <NFormItem
            label="发布时间"
            path="published_at"
          >
            <NDatePicker
              v-model:formatted-value="state.formModel.published_at"
              type="datetime"
              value-format="yyyy-MM-dd HH:mm:ss"
              class="w-full"
              clearable
            />
          </NFormItem>
          <NFormItem
            label="浏览次数"
            path="view_count"
          >
            <NInputNumber
              v-model:value="state.formModel.view_count"
              class="w-full"
            />
          </NFormItem>
          <NFormItem
            label="排序"
            path="sort"
          >
            <NInputNumber
              v-model:value="state.formModel.sort"
              class="w-full"
            />
          </NFormItem>
          <NFormItem
            label="是否置顶"
            path="is_top"
          >
            <NSwitch v-model:value="state.formModel.is_top" />
          </NFormItem>
          <NFormItem
            label="展示设置"
            path="settings"
          >
            <NInput
              v-model:value="state.formModel.settings"
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
