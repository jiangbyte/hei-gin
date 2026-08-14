<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-09 21:39:38
-->

<script setup lang="ts">
import { cgTestActivityApi } from '@/api'
import { createTagColor, dictTypeColor, dictTypeData, displayValue, formatDateTime } from '@/utils'
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
    const response = await cgTestActivityApi.detail({ id })
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
  <NModal v-model:show="state.showModal" preset="card" draggable :mask-closable="false" title="Activity详情" style="width: 680px">
    <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions label-placement="left" bordered :column="1">
          <NDescriptionsItem label="id">
            {{ displayValue(state.detail.id) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="code">
            {{ displayValue(state.detail.code) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="name">
            {{ displayValue(state.detail.name) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="category">
            {{ displayValue(state.detail.category) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="type">
            {{ displayValue(state.detail.type) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="status">
            <NTag :color="createTagColor(dictTypeColor('COMMON_STATUS', state.detail.status))" :bordered="false">
              {{ dictTypeData('COMMON_STATUS', state.detail.status) || displayValue(state.detail.status) }}
            </NTag>
          </NDescriptionsItem>
          <NDescriptionsItem label="cover_url">
            {{ displayValue(state.detail.cover_url) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="description">
            {{ displayValue(state.detail.description) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="start_at">
            {{ formatDateTime(state.detail.start_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="end_at">
            {{ formatDateTime(state.detail.end_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="max_participants">
            {{ displayValue(state.detail.max_participants) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="price">
            {{ displayValue(state.detail.price) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="is_public">
            {{ displayValue(state.detail.is_public) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="need_approval">
            {{ displayValue(state.detail.need_approval) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="rule_config">
            <NCode :code="formatJsonValue(state.detail.rule_config)" language="json" word-wrap />
          </NDescriptionsItem>
          <NDescriptionsItem label="extra">
            <NCode :code="formatJsonValue(state.detail.extra)" language="json" word-wrap />
          </NDescriptionsItem>
          <NDescriptionsItem label="owner_dept_id">
            {{ displayValue(state.detail.owner_dept_id) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="创建时间">{{ formatDateTime(state.detail.created_at) }}</NDescriptionsItem>
          <NDescriptionsItem label="创建人">{{ displayValue(state.detail.created_name || state.detail.created_by) }}</NDescriptionsItem>
          <NDescriptionsItem label="所属部门">{{ displayValue(state.detail.owner_dept_id) }}</NDescriptionsItem>
          <NDescriptionsItem label="更新时间">{{ formatDateTime(state.detail.updated_at) }}</NDescriptionsItem>
          <NDescriptionsItem label="更新人">{{ displayValue(state.detail.updated_name || state.detail.updated_by) }}</NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>
