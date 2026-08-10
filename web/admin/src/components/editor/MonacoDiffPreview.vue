<!-- Author: Charlie -->

<script setup lang="ts">
import type * as monaco from 'monaco-editor'
import { computed, onBeforeUnmount, onMounted, ref, shallowRef, watch } from 'vue'
import { toCssSize } from './shared'
import './monacoWorkers'

const props = withDefaults(
  defineProps<{
    original?: string | null
    modified?: string | null
    language?: string
    theme?: string
    height?: string | number
    options?: monaco.editor.IDiffEditorConstructionOptions
  }>(),
  {
    original: '',
    modified: '',
    language: 'typescript',
    theme: 'vs',
    height: 420,
    options: () => ({}),
  },
)

const containerRef = ref<HTMLElement | null>(null)
const monacoRef = shallowRef<typeof monaco | null>(null)
const diffEditorRef = shallowRef<monaco.editor.IStandaloneDiffEditor | null>(null)
const originalModelRef = shallowRef<monaco.editor.ITextModel | null>(null)
const modifiedModelRef = shallowRef<monaco.editor.ITextModel | null>(null)
let disposed = false
const containerStyle = computed(() => ({
  height: toCssSize(props.height),
}))

async function loadMonaco() {
  if (!monacoRef.value) {
    monacoRef.value = await import('monaco-editor')
  }
  return monacoRef.value
}

function setModelValue(model: monaco.editor.ITextModel | null, value: string) {
  if (model && model.getValue() !== value) {
    model.setValue(value)
  }
}

onMounted(async () => {
  if (!containerRef.value) {
    return
  }

  const monacoInstance = await loadMonaco()
  if (disposed || !containerRef.value) {
    return
  }
  originalModelRef.value = monacoInstance.editor.createModel(props.original ?? '', props.language)
  modifiedModelRef.value = monacoInstance.editor.createModel(props.modified ?? '', props.language)
  diffEditorRef.value = monacoInstance.editor.createDiffEditor(containerRef.value, {
    ...props.options,
    theme: props.theme,
    readOnly: true,
    automaticLayout: true,
    renderSideBySide: true,
    originalEditable: false,
  })
  diffEditorRef.value.setModel({
    original: originalModelRef.value,
    modified: modifiedModelRef.value,
  })
})

watch(
  () => props.original,
  (value) => {
    setModelValue(originalModelRef.value, value ?? '')
  },
)

watch(
  () => props.modified,
  (value) => {
    setModelValue(modifiedModelRef.value, value ?? '')
  },
)

watch(
  () => props.language,
  (language) => {
    const monacoInstance = monacoRef.value
    if (!monacoInstance) {
      return
    }
    if (originalModelRef.value) {
      monacoInstance.editor.setModelLanguage(originalModelRef.value, language)
    }
    if (modifiedModelRef.value) {
      monacoInstance.editor.setModelLanguage(modifiedModelRef.value, language)
    }
  },
)

watch(
  () => props.theme,
  (theme) => {
    monacoRef.value?.editor.setTheme(theme)
  },
)

watch(
  () => props.options,
  (options) => {
    diffEditorRef.value?.updateOptions(options)
  },
)

onBeforeUnmount(() => {
  disposed = true
  diffEditorRef.value?.dispose()
  originalModelRef.value?.dispose()
  modifiedModelRef.value?.dispose()
  diffEditorRef.value = null
  originalModelRef.value = null
  modifiedModelRef.value = null
})
</script>

<template>
  <div
    ref="containerRef"
    class="monaco-diff-preview"
    :style="containerStyle"
  />
</template>

<style scoped>
.monaco-diff-preview {
  width: 100%;
  min-width: 0;
  border: 1px solid var(--border-color);
}
</style>
