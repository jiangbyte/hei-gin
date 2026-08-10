<!-- Author: Charlie -->

<script setup lang="ts">
import type { Themes } from 'md-editor-v3'
import { MdPreview } from 'md-editor-v3'
import { computed } from 'vue'
import 'md-editor-v3/lib/preview.css'

const props = withDefaults(
  defineProps<{
    value?: string | null
    theme?: Themes
    language?: string
    previewTheme?: string
    codeTheme?: string
    showCodeRowNumber?: boolean
  }>(),
  {
    value: '',
    theme: 'light',
    language: 'zh-CN',
    previewTheme: 'default',
    codeTheme: 'atom',
    showCodeRowNumber: true,
  },
)

const modelValue = computed(() => props.value ?? '')
</script>

<template>
  <div class="editor-preview">
    <MdPreview
      :model-value="modelValue"
      :theme="theme"
      :language="language"
      :preview-theme="previewTheme"
      :code-theme="codeTheme"
      :show-code-row-number="showCodeRowNumber"
    />
  </div>
</template>

<style scoped>
.editor-preview {
  width: 100%;
  min-width: 0;
}
</style>
