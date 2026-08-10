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
  order_no: '',
  name: '',
  customer_name: '',
  customer_phone: '',
  status: '',
  type: '',
  ordered_at: null,
  paid_at: null,
  total_amount: 0,
  item_count: 0,
  need_invoice: false,
  invoice_config: '{}',
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

const modalTitle = computed(() => (state.dataId ? '编辑CgTestOrder' : '新增CgTestOrder'))
const rules = computed<FormRules>(() => ({
  order_no: [createRequiredRule('订单号', 'input')],
  name: [createRequiredRule('订单名称', 'input')],
  customer_name: [createRequiredRule('客户名称', 'input')],
  status: [createRequiredRule('状态', 'input')],
  type: [createRequiredRule('订单类型', 'input')],
  ordered_at: [createRequiredRule('下单时间', 'change')],
  total_amount: [
    {
      validator: () =>
        typeof state.formModel.total_amount === 'number' &&
        Number.isFinite(state.formModel.total_amount),
      message: '请输入订单金额',
      trigger: ['input', 'blur'],
    },
  ],
  item_count: [
    {
      validator: () =>
        typeof state.formModel.item_count === 'number' &&
        Number.isFinite(state.formModel.item_count),
      message: '请输入商品数量',
      trigger: ['input', 'blur'],
    },
  ],
  need_invoice: [
    {
      validator: () => typeof state.formModel.need_invoice === 'boolean',
      message: '请选择是否开票',
      trigger: 'change',
    },
  ],
  invoice_config: [
    createRequiredRule('发票配置', 'input'),
    {
      validator: () => isValidJsonValue(state.formModel.invoice_config),
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
    const response = await cgTestOrderApi.detail({ id })
    state.formModel = normalizeFormData(response.data ?? {})
  } finally {
    state.loading = false
  }
}

function normalizeFormData(data: Record<string, any> = {}): Record<string, any> {
  return {
    ...defaultFormData,
    ...data,
    need_invoice:
      data.need_invoice == null || data.need_invoice === ''
        ? defaultFormData.need_invoice
        : wireBool(String(data.need_invoice)),
    item_count:
      data.item_count == null || data.item_count === ''
        ? defaultFormData.item_count
        : wireInt(String(data.item_count)),
    total_amount:
      data.total_amount == null || data.total_amount === ''
        ? defaultFormData.total_amount
        : wireFloat(String(data.total_amount)),
    ordered_at: toFormDateTime(data.ordered_at),
    paid_at: toFormDateTime(data.paid_at),
    invoice_config: stringifyJsonValue(data.invoice_config),
    extra: stringifyJsonValue(data.extra),
  }
}

function normalizeSubmitData(data: Record<string, any>): Record<string, any> {
  return {
    ...data,
    ordered_at: toApiDateTime(data.ordered_at),
    paid_at: toApiDateTime(data.paid_at),
    invoice_config: parseJsonValue(data.invoice_config),
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
      await cgTestOrderApi.update({ ...payload, id: state.dataId })
      window.$message.success('更新成功')
    } else {
      await cgTestOrderApi.create(payload)
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
            label="订单号"
            path="order_no"
          >
            <NInput v-model:value="state.formModel.order_no" />
          </NFormItem>
          <NFormItem
            label="订单名称"
            path="name"
          >
            <NInput v-model:value="state.formModel.name" />
          </NFormItem>
          <NFormItem
            label="客户名称"
            path="customer_name"
          >
            <NInput v-model:value="state.formModel.customer_name" />
          </NFormItem>
          <NFormItem
            label="客户手机号"
            path="customer_phone"
          >
            <NInput v-model:value="state.formModel.customer_phone" />
          </NFormItem>
          <NFormItem
            label="状态"
            path="status"
          >
            <NInput v-model:value="state.formModel.status" />
          </NFormItem>
          <NFormItem
            label="订单类型"
            path="type"
          >
            <NInput v-model:value="state.formModel.type" />
          </NFormItem>
          <NFormItem
            label="下单时间"
            path="ordered_at"
          >
            <NDatePicker
              v-model:formatted-value="state.formModel.ordered_at"
              type="datetime"
              value-format="yyyy-MM-dd HH:mm:ss"
              class="w-full"
              clearable
            />
          </NFormItem>
          <NFormItem
            label="支付时间"
            path="paid_at"
          >
            <NDatePicker
              v-model:formatted-value="state.formModel.paid_at"
              type="datetime"
              value-format="yyyy-MM-dd HH:mm:ss"
              class="w-full"
              clearable
            />
          </NFormItem>
          <NFormItem
            label="订单金额"
            path="total_amount"
          >
            <NInputNumber
              v-model:value="state.formModel.total_amount"
              class="w-full"
            />
          </NFormItem>
          <NFormItem
            label="商品数量"
            path="item_count"
          >
            <NInputNumber
              v-model:value="state.formModel.item_count"
              class="w-full"
            />
          </NFormItem>
          <NFormItem
            label="是否开票"
            path="need_invoice"
          >
            <NSwitch v-model:value="state.formModel.need_invoice" />
          </NFormItem>
          <NFormItem
            label="发票配置"
            path="invoice_config"
          >
            <NInput
              v-model:value="state.formModel.invoice_config"
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
