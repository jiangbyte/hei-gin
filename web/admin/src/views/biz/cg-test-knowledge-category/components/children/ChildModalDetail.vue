<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-09 21:39:43
-->

<script setup lang="ts">
import { cgTestKnowledgeCategoryApi } from '@/api'
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
    const response = await cgTestKnowledgeCategoryApi.childDetail({ id })
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
  <NModal v-model:show="state.showModal" preset="card" draggable :mask-closable="false" title="知识文档详情" style="width: 680px">
    <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions label-placement="left" bordered :column="1">
          <NDescriptionsItem label="id">
            {{ displayValue(state.detail.id) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="category_id">
            {{ displayValue(state.detail.category_id) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="code">
            {{ displayValue(state.detail.code) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="title">
            {{ displayValue(state.detail.title) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="type">
            {{ displayValue(state.detail.type) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="status">
            <NTag :color="createTagColor(dictTypeColor('COMMON_STATUS', state.detail.status))" :bordered="false">
              {{ dictTypeData('COMMON_STATUS', state.detail.status) || displayValue(state.detail.status) }}
            </NTag>
          </NDescriptionsItem>
          <NDescriptionsItem label="summary">
            {{ displayValue(state.detail.summary) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="content">
            {{ displayValue(state.detail.content) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="author">
            {{ displayValue(state.detail.author) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="published_at">
            {{ formatDateTime(state.detail.published_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="view_count">
            {{ displayValue(state.detail.view_count) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="sort">
            {{ displayValue(state.detail.sort) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="is_top">
            {{ displayValue(state.detail.is_top) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="settings">
            <NCode :code="formatJsonValue(state.detail.settings)" language="json" word-wrap />
          </NDescriptionsItem>
          <NDescriptionsItem label="extra">
            <NCode :code="formatJsonValue(state.detail.extra)" language="json" word-wrap />
          </NDescriptionsItem>
          <NDescriptionsItem label="创建时间">{{ formatDateTime(state.detail.created_at) }}</NDescriptionsItem>
          <NDescriptionsItem label="创建人">{{ displayValue(state.detail.created_name || state.detail.created_by) }}</NDescriptionsItem>
          <NDescriptionsItem label="更新时间">{{ formatDateTime(state.detail.updated_at) }}</NDescriptionsItem>
          <NDescriptionsItem label="更新人">{{ displayValue(state.detail.updated_name || state.detail.updated_by) }}</NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>
