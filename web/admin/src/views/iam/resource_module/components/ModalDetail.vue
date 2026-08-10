<!-- Author: Charlie -->

<script setup lang="ts">
import { resourceModuleApi } from '@/api'
import { accountTypeLabel } from '@/constants/account'
import { createTagColor, displayValue, formatDateTime } from '@/utils'
import { reactive } from 'vue'
import { dictTypeColor, dictTypeData } from '@/utils/dict'

const state = reactive({
  showModal: false,
  loading: false,
  module: {} as any,
})

async function openModal(id: string) {
  state.module = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await resourceModuleApi.detail({ id })
    state.module = response.data ?? {}
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
    :title="'资源模块详情'"
    style="width: 620px"
  >
    <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions
          label-placement="left"
          bordered
          :column="1"
        >
          <NDescriptionsItem :label="'资源模块 ID'">
            {{ displayValue(state.module.id) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'模块名称'">
            {{ displayValue(state.module.name) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'模块编码'">
            {{ displayValue(state.module.code) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'客户端'">
            {{ accountTypeLabel(state.module.client) || displayValue(state.module.client) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'图标'">
            <span
              v-if="state.module.icon"
              class="icon-detail-preview"
              :title="state.module.icon"
            >
              <NovaIcon
                :icon="state.module.icon"
                :size="22"
              />
            </span>
            <template v-else>
              -
            </template>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'颜色'">
            <NTag
              v-if="state.module.color"
              :color="createTagColor(state.module.color)"
              :bordered="false"
            >
              {{ state.module.color }}
            </NTag>
            <template v-else>
              -
            </template>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'排序'">
            {{ displayValue(state.module.sort) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'状态'">
            <NTag
              :color="createTagColor(dictTypeColor('COMMON_STATUS', state.module.status))"
              :bordered="false"
            >
              {{
                dictTypeData('COMMON_STATUS', state.module.status) ||
                  displayValue(state.module.status)
              }}
            </NTag>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'描述'">
            {{ displayValue(state.module.description) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建时间'">
            {{ formatDateTime(state.module.created_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建人'">
            {{ displayValue(state.module.created_by) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新时间'">
            {{ formatDateTime(state.module.updated_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新人'">
            {{ displayValue(state.module.updated_by) }}
          </NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>

<style scoped>
.icon-detail-preview {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
}
</style>
