<!-- Author: Charlie -->

<script setup lang="ts">
import { configApi } from '@/api'
import { displayValue, formatDateTime } from '@/utils'
import { computed, reactive } from 'vue'

const state = reactive({
  showModal: false,
  loading: false,
  config: {} as any,
})

const extJsonText = computed(() => {
  const value = state.config.ext_json
  if (!value || typeof value !== 'object') {
    return '{}'
  }
  return JSON.stringify(value, null, 2)
})

async function openModal(id: string) {
  state.config = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await configApi.detail({ id })
    state.config = response.data ?? {}
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
    :title="'系统配置详情'"
    style="width: 680px"
  >
    <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions
          label-placement="left"
          bordered
          :column="1"
        >
          <NDescriptionsItem :label="'配置ID'">
            {{ displayValue(state.config.id) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'配置键'">
            {{ displayValue(state.config.config_key) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'配置值'">
            <NCode
              :code="displayValue(state.config.config_value)"
              word-wrap
            />
          </NDescriptionsItem>
          <NDescriptionsItem :label="'分类'">
            {{ displayValue(state.config.category) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'备注'">
            {{ displayValue(state.config.remark) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'排序码'">
            {{ displayValue(state.config.sort_code) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'扩展信息'">
            <NCode
              :code="extJsonText"
              language="json"
              word-wrap
            />
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建时间'">
            {{ formatDateTime(state.config.created_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建人'">
            {{ displayValue(state.config.created_by) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新时间'">
            {{ formatDateTime(state.config.updated_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新人'">
            {{ displayValue(state.config.updated_by) }}
          </NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>
