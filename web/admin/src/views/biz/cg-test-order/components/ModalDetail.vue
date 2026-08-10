<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-08 21:09:54
-->

<script setup lang="ts">
import { cgTestOrderApi } from '@/api'
import { displayValue, formatDateTime } from '@/utils'
import { reactive } from 'vue'

const state = reactive({
  showModal: false,
  loading: false,
  detail: {} as any,
})

async function openModal(id: string) {
  state.detail = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await cgTestOrderApi.detail({ id })
    state.detail = response.data ?? {}
  } finally {
    state.loading = false
  }
}

function formatJsonValue(value: unknown) {
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
    title="CgTestOrder详情"
    style="width: 680px"
  >
    <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions
          label-placement="left"
          bordered
          :column="1"
        >
          <NDescriptionsItem label="主键">
            {{ displayValue(state.detail.id) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="订单号">
            {{ displayValue(state.detail.order_no) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="订单名称">
            {{ displayValue(state.detail.name) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="客户名称">
            {{ displayValue(state.detail.customer_name) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="客户手机号">
            {{ displayValue(state.detail.customer_phone) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="状态">
            {{ displayValue(state.detail.status) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="订单类型">
            {{ displayValue(state.detail.type) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="下单时间">
            {{ formatDateTime(state.detail.ordered_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="支付时间">
            {{ formatDateTime(state.detail.paid_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="订单金额">
            {{ displayValue(state.detail.total_amount) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="商品数量">
            {{ displayValue(state.detail.item_count) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="是否开票">
            {{ displayValue(state.detail.need_invoice) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="发票配置">
            <NCode
              :code="formatJsonValue(state.detail.invoice_config)"
              language="json"
              word-wrap
            />
          </NDescriptionsItem>
          <NDescriptionsItem label="备注">
            {{ displayValue(state.detail.remark) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="扩展信息">
            <NCode
              :code="formatJsonValue(state.detail.extra)"
              language="json"
              word-wrap
            />
          </NDescriptionsItem>
          <NDescriptionsItem label="创建时间">
            {{ formatDateTime(state.detail.created_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="创建人">
            {{ displayValue(state.detail.created_name || state.detail.created_by) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="所属部门">
            {{ displayValue(state.detail.owner_dept_id) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="更新时间">
            {{ formatDateTime(state.detail.updated_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="更新人">
            {{ displayValue(state.detail.updated_name || state.detail.updated_by) }}
          </NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>
