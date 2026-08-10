<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-08 21:09:54
-->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { cgTestOrderApi } from '@/api'
import { wireBool, wireInt, wireFloat } from '@/utils/wire'
import { createRequiredRule, toApiDateTime, toFormDateTime } from '@/utils'
import { computed, reactive, ref } from 'vue'

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)
const defaultFormData: Record<string, any> = {
  order_id: '',
  sku_code: '',
  name: '',
  category: '',
  status: '',
  quantity: 0,
  unit_price: 0,
  shipped_at: null,
  is_gift: false,
  item_config: '{}',
  remark: '',
  extra: '{}',
}
const state = reactive({
  showModal: false,
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: normalizeFormData(),
})

const modalTitle = computed(() => (state.dataId ? '编辑CgTestOrderItem' : '新增CgTestOrderItem'))
const rules = computed<FormRules>(() => ({
  order_id: [createRequiredRule('订单ID', 'input')],
  sku_code: [createRequiredRule('SKU编码', 'input')],
  name: [createRequiredRule('商品名称', 'input')],
  status: [createRequiredRule('状态', 'input')],
  quantity: [
    {
      validator: () =>
        typeof state.formModel.quantity === 'number' && Number.isFinite(state.formModel.quantity),
      message: '请输入数量',
      trigger: ['input', 'blur'],
    },
  ],
  unit_price: [
    {
      validator: () =>
        typeof state.formModel.unit_price === 'number' &&
        Number.isFinite(state.formModel.unit_price),
      message: '请输入单价',
      trigger: ['input', 'blur'],
    },
  ],
  is_gift: [
    {
      validator: () => typeof state.formModel.is_gift === 'boolean',
      message: '请选择是否赠品',
      trigger: 'change',
    },
  ],
  item_config: [
    createRequiredRule('明细配置', 'input'),
    {
      validator: () => isValidJsonValue(state.formModel.item_config),
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
    const response = await cgTestOrderApi.childDetail({ id })
    state.formModel = normalizeFormData(response.data ?? {})
  } finally {
    state.loading = false
  }
}

function normalizeFormData(data: Record<string, any> = {}): Record<string, any> {
  return {
    ...defaultFormData,
    ...data,
    is_gift:
      data.is_gift == null || data.is_gift === ''
        ? defaultFormData.is_gift
        : wireBool(String(data.is_gift)),
    quantity:
      data.quantity == null || data.quantity === ''
        ? defaultFormData.quantity
        : wireInt(String(data.quantity)),
    unit_price:
      data.unit_price == null || data.unit_price === ''
        ? defaultFormData.unit_price
        : wireFloat(String(data.unit_price)),
    shipped_at: toFormDateTime(data.shipped_at),
    item_config: stringifyJsonValue(data.item_config),
    extra: stringifyJsonValue(data.extra),
  }
}

function normalizeSubmitData(data: Record<string, any>): Record<string, any> {
  return {
    ...data,
    shipped_at: toApiDateTime(data.shipped_at),
    item_config: parseJsonValue(data.item_config),
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
      await cgTestOrderApi.childUpdate({ ...payload, id: state.dataId })
      window.$message.success('更新成功')
    } else {
      await cgTestOrderApi.childCreate(payload)
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
            label="订单ID"
            path="order_id"
          >
            <NInput v-model:value="state.formModel.order_id" />
          </NFormItem>
          <NFormItem
            label="SKU编码"
            path="sku_code"
          >
            <NInput v-model:value="state.formModel.sku_code" />
          </NFormItem>
          <NFormItem
            label="商品名称"
            path="name"
          >
            <NInput v-model:value="state.formModel.name" />
          </NFormItem>
          <NFormItem
            label="商品分类"
            path="category"
          >
            <NInput v-model:value="state.formModel.category" />
          </NFormItem>
          <NFormItem
            label="状态"
            path="status"
          >
            <NInput v-model:value="state.formModel.status" />
          </NFormItem>
          <NFormItem
            label="数量"
            path="quantity"
          >
            <NInputNumber
              v-model:value="state.formModel.quantity"
              class="w-full"
            />
          </NFormItem>
          <NFormItem
            label="单价"
            path="unit_price"
          >
            <NInputNumber
              v-model:value="state.formModel.unit_price"
              class="w-full"
            />
          </NFormItem>
          <NFormItem
            label="发货时间"
            path="shipped_at"
          >
            <NDatePicker
              v-model:formatted-value="state.formModel.shipped_at"
              type="datetime"
              value-format="yyyy-MM-dd HH:mm:ss"
              class="w-full"
              clearable
            />
          </NFormItem>
          <NFormItem
            label="是否赠品"
            path="is_gift"
          >
            <NSwitch v-model:value="state.formModel.is_gift" />
          </NFormItem>
          <NFormItem
            label="明细配置"
            path="item_config"
          >
            <NInput
              v-model:value="state.formModel.item_config"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 8 }"
            />
          </NFormItem>
          <NFormItem
            label="备注"
            path="remark"
          >
            <NInput
              v-model:value="state.formModel.remark"
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
