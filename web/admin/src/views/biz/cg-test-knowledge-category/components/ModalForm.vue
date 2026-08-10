<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-08 21:09:55
-->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { cgTestKnowledgeCategoryApi } from '@/api'
import { wireBool, wireInt } from '@/utils/wire'
import { createRequiredRule } from '@/utils'
import { computed, reactive, ref } from 'vue'

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)
const defaultFormData: Record<string, any> = {
  parent_id: null,
  code: '',
  name: '',
  status: '',
  sort: 0,
  is_visible: false,
  description: '',
  extra: '{}',
}
const state = reactive({
  showModal: false,
  loading: false,
  treeLoading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: normalizeFormData(),
  treeRows: [] as any[],
})

const modalTitle = computed(() =>
  state.dataId ? '编辑CgTestKnowledgeCategory' : '新增CgTestKnowledgeCategory',
)
const parentTreeOptions = computed(() => buildParentTreeOptions(state.treeRows, state.dataId))
const rules = computed<FormRules>(() => ({
  code: [createRequiredRule('分类编码', 'input')],
  name: [createRequiredRule('分类名称', 'input')],
  status: [createRequiredRule('状态', 'input')],
  sort: [
    {
      validator: () =>
        typeof state.formModel.sort === 'number' && Number.isFinite(state.formModel.sort),
      message: '请输入排序',
      trigger: ['input', 'blur'],
    },
  ],
  is_visible: [
    {
      validator: () => typeof state.formModel.is_visible === 'boolean',
      message: '请选择是否显示',
      trigger: 'change',
    },
  ],
  extra: [
    createRequiredRule('扩展信息', 'input'),
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
  await fetchTreeRows()
  if (id) {
    await fetchDetail(id)
  }
}

async function fetchTreeRows() {
  state.treeLoading = true
  try {
    const response = await cgTestKnowledgeCategoryApi.tree()
    state.treeRows = response.data ?? []
  } finally {
    state.treeLoading = false
  }
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await cgTestKnowledgeCategoryApi.detail({ id })
    state.formModel = normalizeFormData(response.data ?? {})
  } finally {
    state.loading = false
  }
}

function normalizeFormData(data: Record<string, any> = {}): Record<string, any> {
  return {
    ...defaultFormData,
    ...data,
    is_visible:
      data.is_visible == null || data.is_visible === ''
        ? defaultFormData.is_visible
        : wireBool(String(data.is_visible)),
    sort: data.sort == null || data.sort === '' ? defaultFormData.sort : wireInt(String(data.sort)),
    extra: stringifyJsonValue(data.extra),
  }
}

function normalizeSubmitData(data: Record<string, any>): Record<string, any> {
  return {
    ...data,
    parent_id: data.parent_id === '' ? null : (data.parent_id ?? null),
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

function buildParentTreeOptions(
  items: any[],
  editingId: string | null,
  disabledParent = false,
): any[] {
  return items.map((item) => {
    const itemId = String(item.id ?? '')
    const disabled = disabledParent || (editingId !== null && itemId === editingId)
    return {
      key: item.id,
      label: String(item.name ?? item.id ?? ''),
      disabled,
      children: buildParentTreeOptions(item.children ?? [], editingId, disabled),
    }
  })
}

async function submitForm() {
  await formRef.value?.validate()
  state.submitLoading = true
  try {
    const payload = normalizeSubmitData(state.formModel)
    if (state.dataId) {
      await cgTestKnowledgeCategoryApi.update({ ...payload, id: state.dataId })
      window.$message.success('更新成功')
    } else {
      await cgTestKnowledgeCategoryApi.create(payload)
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
    <NSpin :show="state.loading || state.treeLoading">
      <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
        <NForm
          ref="formRef"
          :model="state.formModel"
          :rules="rules"
          label-placement="left"
          label-width="110"
          :disabled="state.loading || state.treeLoading || state.submitLoading"
        >
          <NFormItem
            label="父级"
            path="parent_id"
          >
            <NTreeSelect
              v-model:value="state.formModel.parent_id"
              clearable
              filterable
              :options="parentTreeOptions"
              :loading="state.treeLoading"
              key-field="key"
              label-field="label"
              children-field="children"
              class="w-full"
            />
          </NFormItem>
          <NFormItem
            label="分类编码"
            path="code"
          >
            <NInput v-model:value="state.formModel.code" />
          </NFormItem>
          <NFormItem
            label="分类名称"
            path="name"
          >
            <NInput v-model:value="state.formModel.name" />
          </NFormItem>
          <NFormItem
            label="状态"
            path="status"
          >
            <NInput v-model:value="state.formModel.status" />
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
            label="是否显示"
            path="is_visible"
          >
            <NSwitch v-model:value="state.formModel.is_visible" />
          </NFormItem>
          <NFormItem
            label="描述"
            path="description"
          >
            <NInput
              v-model:value="state.formModel.description"
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
