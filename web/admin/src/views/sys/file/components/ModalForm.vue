<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { fileApi } from '@/api'
import { createRequiredRule } from '@/utils'
import { computed, reactive, ref } from 'vue'

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)
const defaultFormData = {
  original_name: '',
}
const state = reactive({
  showModal: false,
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: { ...defaultFormData },
})

const rules = computed<FormRules>(() => ({
  original_name: createRequiredRule('文件名', 'input'),
}))

async function openModal(id: string) {
  state.dataId = id
  state.formModel = { ...defaultFormData }
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await fileApi.detail({ id })
    state.formModel = {
      original_name: response.data?.original_name ?? '',
    }
  } finally {
    state.loading = false
  }
}

function closeModal() {
  state.showModal = false
  state.submitLoading = false
}

async function submitForm() {
  await formRef.value?.validate()
  if (!state.dataId) {
    return
  }
  state.submitLoading = true
  try {
    await fileApi.update({
      id: state.dataId,
      original_name: state.formModel.original_name,
    })
    window.$message.success('更新成功')
    closeModal()
    emit('saved')
  } finally {
    state.submitLoading = false
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
    :title="'编辑 File'"
    style="width: 560px"
    :segmented="{ content: true, action: true }"
  >
    <NSpin :show="state.loading">
      <NForm
        ref="formRef"
        :model="state.formModel"
        :rules="rules"
        label-placement="left"
        label-width="110"
        :disabled="state.loading || state.submitLoading"
      >
        <NFormItem
          :label="'文件名'"
          path="original_name"
        >
          <NInput v-model:value="state.formModel.original_name" />
        </NFormItem>
      </NForm>
    </NSpin>

    <template #action>
      <NSpace
        justify="end"
        align="center"
      >
        <NButton @click="closeModal">
          取消
        </NButton>
        <NButton
          type="primary"
          :loading="state.submitLoading"
          @click="submitForm"
        >
          确认
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
