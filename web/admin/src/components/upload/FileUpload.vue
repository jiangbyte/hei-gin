<!-- Author: Charlie -->

<script setup lang="ts">
import type { UploadFileInfo } from 'naive-ui'
import { Icon } from '@iconify/vue/offline'
import { fileApi } from '@/api'
import { buildPublicFileUrl, formatFileSize, normalizeUploadedFile } from '@/utils'
import type { UploadedFileValueType } from '@/utils'
import { computed, reactive, ref, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    value?: string | null
    file?: File | null
    accept?: string
    buttonText?: string
    icon?: string
    mode?: 'button' | 'icon' | 'upload'
    uploadVariant?: 'default' | 'dragger'
    autoUpload?: boolean
    storageProvider?: string | null
    preview?: 'image' | 'video' | 'file'
    compact?: boolean
    valueType?: UploadedFileValueType
  }>(),
  {
    value: '',
    accept: '',
    buttonText: '',
    icon: '',
    mode: 'button',
    uploadVariant: 'dragger',
    autoUpload: true,
    storageProvider: null,
    preview: 'file',
    compact: false,
    valueType: 'auto',
    file: null,
  },
)

const emit = defineEmits<{
  'update:value': [value: string]
  'update:file': [file: File | null]
  selected: [file: File]
  cleared: []
  uploaded: [file: any]
}>()

const inputRef = ref<HTMLInputElement | null>(null)
const state = reactive({
  loading: false,
  fileName: '',
  fileUrl: '',
  fileSize: null as number | null,
  contentType: null as string | null,
  selectedFile: null as File | null,
  uploadFileList: [] as UploadFileInfo[],
})

const currentUrl = computed(() => {
  if (state.fileUrl) {
    return state.fileUrl
  }
  const value = String(props.value || '').trim()
  if (!value) {
    return undefined
  }
  // 已是可访问地址，或上传后暂存的 url
  if (/^(https?:|data:|blob:)/i.test(value) || value.startsWith('/')) {
    return value
  }
  // 业务字段存的是 object_name 时，走统一公开访问路径预览
  return buildPublicFileUrl(value)
})
const currentName = computed(
  () => state.fileName || props.file?.name || props.value || '未选择文件',
)
const uploadText = computed(() => props.buttonText || '上传')
const actionIcon = computed(() => props.icon || 'icon-park-outline:upload')

watch(
  () => props.file,
  (file) => {
    if (file === state.selectedFile) {
      return
    }
    setSelectedFile(file ?? null, false)
  },
  { immediate: true },
)

function triggerUpload() {
  inputRef.value?.click()
}

function clearValue() {
  clearSelectedFile()
  emit('update:value', '')
}

function clearSelectedFile() {
  state.fileName = ''
  state.fileUrl = ''
  state.fileSize = null
  state.contentType = null
  state.selectedFile = null
  state.uploadFileList = []
  emit('update:file', null)
  emit('cleared')
}

async function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) {
    return
  }
  await selectFile(file)
}

async function handleUploadFileListChange(fileList: UploadFileInfo[]) {
  const fileInfo = [...fileList].reverse().find((item) => item.file)
  if (!fileInfo?.file) {
    clearSelectedFile()
    return
  }
  await selectFile(fileInfo.file, fileInfo)
}

async function selectFile(file: File, fileInfo?: UploadFileInfo) {
  setSelectedFile(file, true, fileInfo)
  if (props.autoUpload) {
    await uploadSelectedFile()
  }
}

function setSelectedFile(file: File | null, shouldEmit: boolean, fileInfo?: UploadFileInfo) {
  state.selectedFile = file
  state.fileUrl = ''
  state.fileName = file?.name || ''
  state.fileSize = file?.size ?? null
  state.contentType = file?.type || null
  state.uploadFileList = file ? [fileInfo ?? createUploadFileInfo(file)] : []
  if (shouldEmit) {
    emit('update:file', file)
    if (file) {
      emit('selected', file)
    }
  }
}

async function uploadSelectedFile() {
  const file = state.selectedFile
  if (!file) {
    return undefined
  }
  state.loading = true
  updateUploadStatus('uploading')
  try {
    const response = await fileApi.upload(file, {
      storage_provider: props.storageProvider,
    })
    const uploaded = response.data ?? {}
    const normalized = normalizeUploadedFile(uploaded, file, props.valueType)
    state.fileName = normalized.name
    state.fileUrl = normalized.url
    state.fileSize = normalized.size
    state.contentType = normalized.contentType
    emit('update:value', normalized.value)
    emit('uploaded', { ...normalized, ...uploaded })
    window.$message.success('上传成功')
    updateUploadStatus('finished')
    return uploaded
  } catch (error) {
    updateUploadStatus('error')
    throw error
  } finally {
    state.loading = false
  }
}

function updateUploadStatus(status: UploadFileInfo['status']) {
  state.uploadFileList = state.uploadFileList.map((item) => ({
    ...item,
    status,
    percentage: status === 'finished' ? 100 : item.percentage,
  }))
}

function createUploadFileInfo(file: File): UploadFileInfo {
  return {
    id: `${Date.now()}-${file.name}`,
    name: file.name,
    status: 'pending',
    file,
    type: file.type,
  }
}

defineExpose({
  clear: clearValue,
  upload: uploadSelectedFile,
})
</script>

<template>
  <div
    class="file-upload"
    :class="{ 'file-upload--compact': compact && mode !== 'upload' }"
  >
    <input
      ref="inputRef"
      class="file-upload__input"
      type="file"
      :accept="accept"
      @change="handleFileChange"
    >

    <div
      v-if="mode !== 'upload' && !compact && preview === 'image' && currentUrl"
      class="file-upload__image"
    >
      <NImage
        :src="currentUrl"
        object-fit="cover"
        :alt="currentName"
        width="160"
        height="90"
      />
    </div>
    <video
      v-else-if="mode !== 'upload' && !compact && preview === 'video' && currentUrl"
      class="file-upload__video"
      controls
      :src="currentUrl"
    />
    <NEllipsis
      v-else-if="mode !== 'upload' && !compact"
      class="file-upload__name"
    >
      {{ currentName }}
    </NEllipsis>
    <div
      v-if="mode !== 'upload' && !compact && (state.fileSize !== null || state.contentType)"
      class="file-upload__meta"
    >
      <span v-if="state.fileSize !== null">{{ formatFileSize(state.fileSize) }}</span>
      <span v-if="state.contentType">{{ state.contentType }}</span>
    </div>

    <NUpload
      v-if="mode === 'upload'"
      class="file-upload__native"
      :accept="accept"
      :default-upload="false"
      :disabled="state.loading"
      :file-list="state.uploadFileList"
      :max="1"
      :show-cancel-button="false"
      :show-retry-button="false"
      @update:file-list="handleUploadFileListChange"
    >
      <NUploadDragger v-if="uploadVariant === 'dragger'">
        <div class="file-upload__dragger">
          <NIcon size="28">
            <Icon :icon="actionIcon" />
          </NIcon>
          <div>{{ uploadText || '选择文件' }}</div>
        </div>
      </NUploadDragger>
      <NButton
        v-else
        :loading="state.loading"
      >
        <template #icon>
          <NIcon>
            <Icon :icon="actionIcon" />
          </NIcon>
        </template>
        {{ uploadText || '选择文件' }}
      </NButton>
    </NUpload>

    <div
      v-else
      class="file-upload__actions"
      :class="{ 'file-upload__actions--compact': compact }"
    >
      <NButton
        v-if="mode === 'icon'"
        text
        :loading="state.loading"
        :title="uploadText"
        :aria-label="uploadText"
        @click="triggerUpload"
      >
        <template #icon>
          <NIcon>
            <Icon :icon="actionIcon" />
          </NIcon>
        </template>
      </NButton>
      <NButton
        v-else
        size="small"
        :loading="state.loading"
        @click="triggerUpload"
      >
        <template
          v-if="icon"
          #icon
        >
          <NIcon>
            <Icon :icon="icon" />
          </NIcon>
        </template>
        {{ uploadText }}
      </NButton>
      <NButton
        v-if="value"
        size="small"
        text
        type="error"
        @click="clearValue"
      >
        清除
      </NButton>
    </div>
  </div>
</template>

<style scoped>
.file-upload {
  display: grid;
  gap: 8px;
  width: 100%;
  min-width: 0;
}

.file-upload--compact {
  display: inline-flex;
  align-items: center;
  width: auto;
}

.file-upload__input {
  display: none;
}

.file-upload__native {
  width: 100%;
  min-width: 0;
}

.file-upload__dragger {
  display: grid;
  justify-items: center;
  gap: 6px;
  padding: 6px 0;
  color: var(--text-color-3);
}

.file-upload__actions {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.file-upload__actions--compact {
  gap: 0;
}

.file-upload__image,
.file-upload__video {
  width: 160px;
  max-width: 100%;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  overflow: hidden;
  background: var(--body-color);
}

.file-upload__video {
  aspect-ratio: 16 / 9;
}

.file-upload__name {
  max-width: 100%;
  color: var(--text-color-3);
}

.file-upload__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  color: var(--text-color-3);
  font-size: 12px;
  line-height: 1.3;
}
</style>
