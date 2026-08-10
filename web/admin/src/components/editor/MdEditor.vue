<!-- Author: Charlie -->

<script setup lang="ts">
import type { Themes, UploadImgEvent } from 'md-editor-v3'
import type { EditorUploadedFile } from './shared'
import { MdEditor } from 'md-editor-v3'
import { computed } from 'vue'
import { toCssSize, uploadEditorFile } from './shared'
import 'md-editor-v3/lib/style.css'

const props = withDefaults(
  defineProps<{
    value?: string | null
    height?: string | number
    placeholder?: string
    theme?: Themes
    preview?: boolean
    readOnly?: boolean
    disabled?: boolean
    language?: string
    previewTheme?: string
    codeTheme?: string
    showCodeRowNumber?: boolean
    noUploadImg?: boolean
  }>(),
  {
    value: '',
    height: 360,
    placeholder: '请输入内容',
    theme: 'light',
    preview: true,
    readOnly: false,
    disabled: false,
    language: 'zh-CN',
    previewTheme: 'default',
    codeTheme: 'atom',
    showCodeRowNumber: true,
    noUploadImg: false,
  },
)

const emit = defineEmits<{
  'update:value': [value: string]
  change: [value: string]
  uploaded: [file: EditorUploadedFile]
  'upload-error': [error: unknown]
}>()

const editorId = `md-editor-${Math.random().toString(36).slice(2)}`
const modelValue = computed(() => props.value ?? '')
const editorStyle = computed(() => ({
  height: toCssSize(props.height),
}))

function updateValue(value: string) {
  emit('update:value', value)
  emit('change', value)
}

const handleUploadImg: UploadImgEvent = async (files, callback) => {
  const images: Array<{ url: string; alt: string; title: string }> = []

  for (const file of files) {
    try {
      const uploaded = await uploadEditorFile(file)
      if (uploaded.url) {
        images.push({
          url: uploaded.url,
          alt: uploaded.name,
          title: uploaded.name,
        })
      }
      emit('uploaded', uploaded)
    } catch (error) {
      emit('upload-error', error)
      window.$message?.error('上传失败')
    }
  }

  if (images.length) {
    callback(images)
  }
}
</script>

<template>
  <div class="editor-wrapper">
    <MdEditor
      :id="editorId"
      :model-value="modelValue"
      :style="editorStyle"
      :theme="theme"
      :preview="preview"
      :read-only="readOnly"
      :disabled="disabled"
      :language="language"
      :placeholder="placeholder"
      :preview-theme="previewTheme"
      :code-theme="codeTheme"
      :show-code-row-number="showCodeRowNumber"
      :no-upload-img="noUploadImg"
      @on-upload-img="handleUploadImg"
      @update:model-value="updateValue"
    />
  </div>
</template>

<style scoped>
.editor-wrapper {
  width: 100%;
  min-width: 0;
}
</style>
