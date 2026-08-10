<!-- Author: Charlie -->

<script setup lang="ts">
import { resourceApi } from '@/api'
import { createTagColor, displayValue, formatDateTime } from '@/utils'
import { reactive } from 'vue'
import { dictTypeData, dictTypeColor } from '@/utils/dict'

const state = reactive({
  showModal: false,
  loading: false,
  resource: {} as any,
})

async function openModal(id: string) {
  state.resource = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await resourceApi.detail({ id })
    state.resource = response.data ?? {}
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
    :title="'资源详情'"
    style="width: 720px"
  >
    <NScrollbar class="max-h-[min(640px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions
          label-placement="left"
          bordered
          :column="1"
        >
          <NDescriptionsItem :label="'资源ID'">
            {{ displayValue(state.resource.id) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'资源名称'">
            {{ displayValue(state.resource.name) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'资源编码'">
            {{ displayValue(state.resource.code) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'资源类型'">
            {{
              dictTypeData('RESOURCE_TYPE', state.resource.resource_type) ||
                displayValue(state.resource.resource_type)
            }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'父级资源ID'">
            {{ displayValue(state.resource.parent_id) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'资源模块'">
            {{ displayValue(state.resource.module_id_name || state.resource.module_id) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'路由路径'">
            {{ displayValue(state.resource.path) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'组件'">
            {{ displayValue(state.resource.component) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'重定向'">
            {{ displayValue(state.resource.redirect) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'图标'">
            <span
              v-if="state.resource.icon"
              class="icon-detail-preview"
              :title="state.resource.icon"
            >
              <NovaIcon
                :icon="state.resource.icon"
                :size="22"
              />
            </span>
            <template v-else>
              -
            </template>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'颜色'">
            <NTag
              v-if="state.resource.color"
              :color="createTagColor(state.resource.color)"
              :bordered="false"
            >
              {{ state.resource.color }}
            </NTag>
            <span v-else>{{ displayValue(state.resource.color) }}</span>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'外链'">
            {{ displayValue(state.resource.href) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'排序'">
            {{ displayValue(state.resource.sort) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'可见'">
            {{ state.resource.is_visible ? '是' : '否' }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'缓存'">
            {{ state.resource.is_cache ? '是' : '否' }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'固定标签'">
            {{ state.resource.is_affix ? '是' : '否' }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'状态'">
            <NTag
              :color="createTagColor(dictTypeColor('COMMON_STATUS', state.resource.status))"
              :bordered="false"
            >
              {{
                dictTypeData('COMMON_STATUS', state.resource.status) ||
                  displayValue(state.resource.status)
              }}
            </NTag>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'描述'">
            {{ displayValue(state.resource.description) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'布局'">
            {{ displayValue(state.resource.layout) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建时间'">
            {{ formatDateTime(state.resource.created_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建人'">
            {{ displayValue(state.resource.created_by) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新时间'">
            {{ formatDateTime(state.resource.updated_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新人'">
            {{ displayValue(state.resource.updated_by) }}
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
