<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-08 21:09:55
-->

<script setup lang="ts">
import { cgTestKnowledgeCategoryApi } from '@/api'
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
  <NModal
    v-model:show="state.showModal"
    preset="card"
    draggable
    :mask-closable="false"
    title="CgTestKnowledgeDoc详情"
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
          <NDescriptionsItem label="分类ID">
            {{ displayValue(state.detail.category_id) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="文档编码">
            {{ displayValue(state.detail.code) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="文档标题">
            {{ displayValue(state.detail.title) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="文档类型">
            {{ displayValue(state.detail.type) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="状态">
            {{ displayValue(state.detail.status) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="摘要">
            {{ displayValue(state.detail.summary) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="正文内容">
            {{ displayValue(state.detail.content) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="作者">
            {{ displayValue(state.detail.author) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="发布时间">
            {{ formatDateTime(state.detail.published_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="浏览次数">
            {{ displayValue(state.detail.view_count) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="排序">
            {{ displayValue(state.detail.sort) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="是否置顶">
            {{ displayValue(state.detail.is_top) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="展示设置">
            <NCode
              :code="formatJsonValue(state.detail.settings)"
              language="json"
              word-wrap
            />
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
