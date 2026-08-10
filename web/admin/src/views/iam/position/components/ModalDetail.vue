<!-- Author: Charlie -->

<script setup lang="ts">
import { positionApi } from '@/api'
import { createTagColor, displayValue, formatDateTime } from '@/utils'
import { reactive } from 'vue'
import { dictTypeData, dictTypeColor } from '@/utils/dict'

const state = reactive({
  showModal: false,
  loading: false,
  position: {} as any,
})

async function openModal(id: string) {
  state.position = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await positionApi.detail({ id })
    state.position = response.data ?? {}
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
    :title="'岗位详情'"
    style="width: 680px"
  >
    <NScrollbar class="max-h-[min(640px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions
          label-placement="left"
          bordered
          :column="1"
        >
          <NDescriptionsItem :label="'岗位 ID'">
            {{ displayValue(state.position.id) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'岗位名称'">
            {{ displayValue(state.position.name) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'岗位分类'">
            {{
              dictTypeData('POSITION_CATEGORY', state.position.category) ||
                displayValue(state.position.category)
            }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'排序'">
            {{ displayValue(state.position.sort) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'虚拟岗位'">
            {{ state.position.is_virtual ? '是' : '否' }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'状态'">
            <NTag
              :color="createTagColor(dictTypeColor('COMMON_STATUS', state.position.status))"
              :bordered="false"
            >
              {{
                dictTypeData('COMMON_STATUS', state.position.status) ||
                  displayValue(state.position.status)
              }}
            </NTag>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'描述'">
            {{ displayValue(state.position.description) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建时间'">
            {{ formatDateTime(state.position.created_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建人'">
            {{ displayValue(state.position.created_name) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新时间'">
            {{ formatDateTime(state.position.updated_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新人'">
            {{ displayValue(state.position.updated_name) }}
          </NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>
