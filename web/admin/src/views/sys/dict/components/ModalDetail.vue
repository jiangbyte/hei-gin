<!-- Author: Charlie -->

<script setup lang="ts">
import { dictApi } from '@/api'
import { createTagColor, displayValue, formatDateTime } from '@/utils'
import { reactive } from 'vue'
import { dictTypeData, dictTypeColor } from '@/utils/dict'

const state = reactive({
  showModal: false,
  loading: false,
  data: {} as any,
})

async function openModal(id: string) {
  state.data = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await dictApi.detail({ id })
    state.data = response.data ?? {}
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
    :title="'字典详情'"
    style="width: 620px"
  >
    <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions
          label-placement="left"
          bordered
          :column="1"
        >
          <NDescriptionsItem :label="'ID'">
            {{ displayValue(state.data.id) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'编码'">
            {{ displayValue(state.data.code) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'标签'">
            {{ displayValue(state.data.label) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'值'">
            {{ displayValue(state.data.value) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'颜色'">
            <NTag
              v-if="state.data.color"
              size="small"
              :color="createTagColor(state.data.color)"
              :bordered="false"
            >
              {{ state.data.color }}
            </NTag>
            <template v-else>
              -
            </template>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'分类'">
            {{ dictTypeData('SYS_BIZ_CATEGORY', state.data.category) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'父级字典'">
            {{ displayValue(state.data.parent_id_name) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'排序'">
            {{ displayValue(state.data.sort) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'状态'">
            <NTag
              size="small"
              :color="createTagColor(dictTypeColor('COMMON_STATUS', state.data.status))"
              :bordered="false"
            >
              {{ dictTypeData('COMMON_STATUS', state.data.status) }}
            </NTag>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建时间'">
            {{ formatDateTime(state.data.created_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建人'">
            {{ displayValue(state.data.created_name) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新时间'">
            {{ formatDateTime(state.data.updated_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新人'">
            {{ displayValue(state.data.updated_name) }}
          </NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>
