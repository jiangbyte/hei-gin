<!-- Author: Charlie -->

<script setup lang="ts">
import { roleApi } from '@/api'
import { createTagColor, displayValue, formatDateTime } from '@/utils'
import { reactive } from 'vue'
import { dictTypeData, dictTypeColor } from '@/utils/dict'

const state = reactive({
  showModal: false,
  loading: false,
  role: {} as any,
})

async function openModal(id: string) {
  state.role = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await roleApi.detail({ id })
    state.role = response.data ?? {}
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
    :title="'角色详情'"
    style="width: 680px"
  >
    <NScrollbar class="max-h-[min(640px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions
          label-placement="left"
          bordered
          :column="1"
        >
          <NDescriptionsItem :label="'角色ID'">
            {{ displayValue(state.role.id) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'角色名称'">
            {{ displayValue(state.role.name) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'角色编码'">
            {{ displayValue(state.role.code) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'角色分类'">
            {{
              dictTypeData('SYS_BIZ_CATEGORY', state.role.category) ||
                displayValue(state.role.category)
            }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'范围类型'">
            {{
              dictTypeData('ROLE_SCOPE_TYPE', state.role.scope_type) ||
                displayValue(state.role.scope_type)
            }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'所属部门'">
            {{ displayValue(state.role.owner_dept_name) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'排序'">
            {{ displayValue(state.role.sort) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'内置角色'">
            {{ state.role.is_builtin ? '是' : '否' }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'状态'">
            <NTag
              :color="createTagColor(dictTypeColor('COMMON_STATUS', state.role.status))"
              :bordered="false"
            >
              {{
                dictTypeData('COMMON_STATUS', state.role.status) || displayValue(state.role.status)
              }}
            </NTag>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'描述'">
            {{ displayValue(state.role.description) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建时间'">
            {{ formatDateTime(state.role.created_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建人'">
            {{ displayValue(state.role.created_name) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新时间'">
            {{ formatDateTime(state.role.updated_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新人'">
            {{ displayValue(state.role.updated_name) }}
          </NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>
