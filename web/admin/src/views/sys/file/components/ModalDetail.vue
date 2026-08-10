<!-- Author: Charlie -->

<script setup lang="ts">
import { fileApi } from '@/api'
import { displayValue, formatDateTime, formatFileSize, isImageFile } from '@/utils'
import { computed, reactive } from 'vue'
import { storageProviderLabel } from '../constants'

const state = reactive({
  showModal: false,
  loading: false,
  downloading: false,
  file: {} as any,
})

const fileUrl = computed(() => state.file?.url || undefined)
const imageAlt = computed(() => state.file?.original_name ?? '预览')
const isImage = computed(() => isImageFile(state.file))

async function openModal(id: string) {
  state.file = {}
  state.showModal = true
  await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await fileApi.detail({ id })
    state.file = response.data
  } finally {
    state.loading = false
  }
}

function openFile() {
  if (!fileUrl.value) {
    return
  }
  window.open(fileUrl.value, '_blank', 'noopener,noreferrer')
}

async function downloadFile() {
  const id = String(state.file?.id || '')
  if (!id || state.downloading) {
    return
  }
  state.downloading = true
  try {
    await fileApi.downloadFile(
      state.file,
      state.file.original_name || state.file.object_name || 'download',
    )
  } finally {
    state.downloading = false
  }
}

async function copyText(value?: string | null) {
  const text = String(value || '')
  if (!text) {
    return
  }
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text)
  } else {
    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
  }
  window.$message.success('Copied successfully')
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
    :title="'File 详情'"
    style="width: 700px"
  >
    <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
      <NSpin :show="state.loading">
        <NDescriptions
          label-placement="left"
          bordered
          :column="1"
          label-style="min-width: 120px"
          class="file-detail-descriptions"
        >
          <NDescriptionsItem :label="'预览'">
            <NImage
              v-if="isImage && fileUrl"
              class="file-detail-image"
              :src="fileUrl"
              :alt="imageAlt"
              :width="220"
              :height="140"
              object-fit="cover"
            />
            <NButton
              v-else-if="state.file.id"
              type="primary"
              text
              :loading="state.downloading"
              @click="downloadFile"
            >
              下载
            </NButton>
            <template v-else>
              -
            </template>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'文件ID'">
            {{ displayValue(state.file.id) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'文件名'">
            {{ displayValue(state.file.original_name) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'对象路径'">
            <NFlex
              align="center"
              :size="8"
            >
              <span class="file-detail-text">{{ displayValue(state.file.object_name) }}</span>
              <NButton
                size="small"
                text
                type="primary"
                @click="copyText(state.file.object_name)"
              >
                复制
              </NButton>
            </NFlex>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'访问URL'">
            <NFlex
              align="center"
              :size="8"
            >
              <span class="file-detail-text">{{ displayValue(fileUrl) }}</span>
              <NButton
                v-if="fileUrl"
                size="small"
                text
                type="primary"
                @click="openFile"
              >
                打开
              </NButton>
              <NButton
                v-if="state.file.id"
                size="small"
                text
                type="primary"
                :loading="state.downloading"
                @click="downloadFile"
              >
                下载
              </NButton>
              <NButton
                v-if="fileUrl"
                size="small"
                text
                type="primary"
                @click="copyText(fileUrl)"
              >
                复制
              </NButton>
            </NFlex>
          </NDescriptionsItem>
          <NDescriptionsItem :label="'存储提供商'">
            {{
              storageProviderLabel(state.file.storage_provider) ||
                displayValue(state.file.storage_provider)
            }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'存储桶'">
            {{ displayValue(state.file.bucket) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'内容类型'">
            {{ displayValue(state.file.content_type) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'文件大小'">
            {{ formatFileSize(state.file.size) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建时间'">
            {{ formatDateTime(state.file.created_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'创建人'">
            {{ displayValue(state.file.created_name) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新时间'">
            {{ formatDateTime(state.file.updated_at) }}
          </NDescriptionsItem>
          <NDescriptionsItem :label="'更新人'">
            {{ displayValue(state.file.updated_name) }}
          </NDescriptionsItem>
        </NDescriptions>
      </NSpin>
    </NScrollbar>
  </NModal>
</template>

<style scoped>
.file-detail-image {
  border-radius: 6px;
}

.file-detail-text {
  word-break: break-all;
}
</style>
