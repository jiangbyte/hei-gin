<!--
  由 HEI 代码生成器生成。
  Author: Charlie
  生成时间：2026-08-08 21:09:52
-->

<script setup lang="ts">
import { cgTestActivityApi } from '@/api'
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
  <NModal
    v-model:show="state.showModal"
    preset="card"
    draggable
    :mask-closable="false"
    title="CgTestActivity详情"
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
          <NDescriptionsItem label="活动编码">
            {{ displayValue(state.detail.code) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="活动名称">
            {{ displayValue(state.detail.name) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="活动分类">
            {{ displayValue(state.detail.category) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="活动类型">
            {{ displayValue(state.detail.type) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="状态">
            {{ displayValue(state.detail.status) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="封面地址">
            {{ displayValue(state.detail.cover_url) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="活动描述">
            {{ displayValue(state.detail.description) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="开始时间">
            {{ formatDateTime(state.detail.start_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="结束时间">
            {{ formatDateTime(state.detail.end_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="最大参与人数">
            {{ displayValue(state.detail.max_participants) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="报名费用">
            {{ displayValue(state.detail.price) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="是否公开">
            {{ displayValue(state.detail.is_public) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="是否需要审批">
            {{ displayValue(state.detail.need_approval) }}
          </NDescriptionsItem>
          <NDescriptionsItem label="规则配置">
            <NCode
              :code="formatJsonValue(state.detail.rule_config)"
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
