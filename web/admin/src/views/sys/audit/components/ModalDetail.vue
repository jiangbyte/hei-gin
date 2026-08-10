<!--
  Author: Charlie

  操作审计详情。
-->
<script setup lang="ts">
import { auditApi } from '@/api'
import { accountTypeLabel } from '@/constants/account'
import { displayValue, formatDateTime } from '@/utils'
import { wireBool } from '@/utils/wire'
import { computed, reactive } from 'vue'

const state = reactive({
  showModal: false,
  loading: false,
  record: {} as any,
})

const successText = computed(() => {
  if (state.record?.success === undefined || state.record?.success === null) {
    return '-'
  }
  return wireBool(state.record.success) ? '成功' : '失败'
})

function formatJson(value: unknown) {
  if (value === undefined || value === null || value === '') {
    return '-'
  }
  if (typeof value === 'string') {
    try {
      return JSON.stringify(JSON.parse(value), null, 2)
    } catch {
      return value
    }
  }
  try {
    return JSON.stringify(value, null, 2)
  } catch {
    return String(value)
  }
}

async function openModal(id: string) {
  state.record = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await auditApi.detail({ id })
    state.record = response.data ?? {}
  } finally {
    state.loading = false
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
    title="审计详情"
    style="width: 760px"
  >
    <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions
          label-placement="left"
          bordered
          :column="1"
          label-style="min-width: 120px"
        >
          <NDescriptionsItem label="时间">
            {{ formatDateTime(state.record.created_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="模块">
            {{ displayValue(state.record.module) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="动作">
            {{ displayValue(state.record.action) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="摘要">
            {{ displayValue(state.record.summary) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="结果">
            {{ successText }}
          </NDescriptionsItem>
          <NDescriptionsItem label="账号 ID">
            {{ displayValue(state.record.account_id) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="账号类型">
            {{ state.record.account_type ? accountTypeLabel(state.record.account_type) : '-' }}
          </NDescriptionsItem>
          <NDescriptionsItem label="资源类型">
            {{ displayValue(state.record.resource_type) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="资源 ID">
            {{ displayValue(state.record.resource_id) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="请求 ID">
            {{ displayValue(state.record.request_id) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="IP">
            {{ displayValue(state.record.ip) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="User-Agent">
            <pre class="audit-pre">{{ displayValue(state.record.user_agent) }}</pre>
          </NDescriptionsItem>
          <NDescriptionsItem label="错误信息">
            <pre class="audit-pre">{{ displayValue(state.record.error_message) }}</pre>
          </NDescriptionsItem>
          <NDescriptionsItem label="变更前">
            <pre class="audit-pre">{{ formatJson(state.record.before_data) }}</pre>
          </NDescriptionsItem>
          <NDescriptionsItem label="变更后">
            <pre class="audit-pre">{{ formatJson(state.record.after_data) }}</pre>
          </NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>

<style scoped>
.audit-pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  line-height: 1.5;
}
</style>
