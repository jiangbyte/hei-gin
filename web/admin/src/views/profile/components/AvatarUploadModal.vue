<!-- Author: Charlie -->

<script setup lang="ts">
import { authApi } from '@/api'
import { computed, reactive, ref, watch } from 'vue'
import { Cropper } from 'vue-advanced-cropper'
import 'vue-advanced-cropper/dist/style.css'

const props = defineProps<{
  show: boolean
  avatar?: string | null
}>()

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void
  (e: 'uploaded'): void
}>()

const fileInputRef = ref<HTMLInputElement | null>(null)
const cropperRef = ref<any>(null)

const state = reactive({
  source: '',
  previewUrl: '',
  fileName: '',
  uploading: false,
})

const modalShow = computed({
  get: () => props.show,
  set: (value: boolean) => emit('update:show', value),
})

watch(
  () => props.show,
  (show) => {
    if (!show) {
      resetSource()
    }
  },
)

function openFilePicker() {
  fileInputRef.value?.click()
}

function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) {
    return
  }
  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
    window.$message.warning('仅支持 JPG、PNG 和 WebP 图片，裁剪后上传')
    return
  }
  resetSource()
  state.fileName = file.name
  state.source = URL.createObjectURL(file)
  updatePreview()
}

async function uploadAvatar() {
  const canvas = cropperRef.value?.getResult?.().canvas
  if (!canvas) {
    window.$message.warning('请先选择头像图片')
    return
  }

  state.uploading = true
  try {
    const blob = await canvasToBlob(canvas)
    await authApi.uploadAvatar(new File([blob], 'avatar.png', { type: 'image/png' }))
    window.$message.success('头像已更新')
    emit('uploaded')
    modalShow.value = false
  } finally {
    state.uploading = false
  }
}

function onCropperChange() {
  updatePreview()
}

function updatePreview() {
  window.requestAnimationFrame(() => {
    const canvas = cropperRef.value?.getResult?.().canvas
    state.previewUrl = canvas ? canvas.toDataURL('image/png') : ''
  })
}

function canvasToBlob(canvas: HTMLCanvasElement) {
  return new Promise<Blob>((resolve, reject) => {
    canvas.toBlob(
      (blob) => {
        if (blob) {
          resolve(blob)
          return
        }
        reject(new Error('头像图片导出失败'))
      },
      'image/png',
      0.92,
    )
  })
}

function resetSource() {
  if (state.source) {
    URL.revokeObjectURL(state.source)
  }
  state.source = ''
  state.previewUrl = ''
  state.fileName = ''
}
</script>

<template>
  <NModal
    v-model:show="modalShow"
    preset="card"
    :title="'上传头像'"
    class="max-w-150"
    :bordered="false"
    :mask-closable="!state.uploading"
  >
    <input
      ref="fileInputRef"
      class="hidden"
      type="file"
      accept="image/jpeg,image/png,image/webp"
      @change="onFileChange"
    >

    <div class="avatar-upload-modal">
      <div
        v-if="state.source"
        class="avatar-editor"
      >
        <div class="avatar-cropper-wrap">
          <Cropper
            ref="cropperRef"
            class="avatar-cropper"
            :src="state.source"
            :stencil-props="{ aspectRatio: 1 }"
            :canvas="{ width: 320, height: 320 }"
            @change="onCropperChange"
          />
        </div>
        <div class="avatar-preview-panel">
          <div class="avatar-preview">
            <img
              v-if="state.previewUrl"
              :src="state.previewUrl"
              alt=""
            >
          </div>
          <div class="mt-3 text-sm text-[var(--text-color-3)]">
            实时预览
          </div>
        </div>
      </div>
      <div
        v-else
        class="avatar-empty"
      >
        <NAvatar
          v-if="avatar"
          round
          :size="96"
          :src="avatar"
        />
        <NAvatar
          v-else
          round
          :size="96"
        >
          <NovaIcon
            icon="icon-park-outline:user"
            :size="40"
          />
        </NAvatar>
        <NButton
          secondary
          type="primary"
          class="mt-4"
          @click="openFilePicker"
        >
          <template #icon>
            <NovaIcon icon="icon-park-outline:upload-picture" />
          </template>
          选择图片
        </NButton>
      </div>
    </div>

    <template #footer>
      <NSpace
        justify="space-between"
        align="center"
      >
        <div class="max-w-70 truncate text-sm text-[var(--text-color-3)]">
          {{ state.fileName || '仅支持 JPG、PNG 和 WebP 图片，裁剪后上传' }}
        </div>
        <NSpace>
          <NButton @click="modalShow = false">
            取消
          </NButton>
          <NButton
            v-if="state.source"
            secondary
            @click="openFilePicker"
          >
            重新选择
          </NButton>
          <NButton
            type="primary"
            :loading="state.uploading"
            @click="uploadAvatar"
          >
            裁剪并上传
          </NButton>
        </NSpace>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped>
.avatar-upload-modal {
  min-height: 360px;
}

.avatar-editor {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 180px;
  gap: 16px;
  align-items: stretch;
}

.avatar-cropper-wrap {
  height: min(56vh, 420px);
  overflow: hidden;
  border: 1px solid var(--border-color);
  border-radius: 8px;
}

.avatar-cropper {
  height: 100%;
  width: 100%;
  background: var(--body-color);
}

.avatar-preview-panel {
  min-height: 360px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 16px;
}

.avatar-preview {
  width: 128px;
  height: 128px;
  overflow: hidden;
  border-radius: 999px;
  background: var(--body-color);
  box-shadow: inset 0 0 0 1px var(--border-color);
}

.avatar-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.avatar-empty {
  min-height: 360px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 1px dashed var(--border-color);
  border-radius: 8px;
}

@media (max-width: 640px) {
  .avatar-editor {
    grid-template-columns: 1fr;
  }

  .avatar-preview-panel {
    min-height: 190px;
  }
}
</style>
