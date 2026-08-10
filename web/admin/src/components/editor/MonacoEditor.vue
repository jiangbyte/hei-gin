<!-- Author: Charlie -->

<script setup lang="ts">
import type * as monaco from 'monaco-editor'
import { computed, onBeforeUnmount, onMounted, ref, shallowRef, watch } from 'vue'
import { toCssSize } from './shared'
import './monacoWorkers'

const props = withDefaults(
  defineProps<{
    value?: string | null
    language?: string
    theme?: string
    height?: string | number
    readOnly?: boolean
    options?: monaco.editor.IStandaloneEditorConstructionOptions
  }>(),
  {
    value: '',
    language: 'typescript',
    theme: 'vs',
    height: 360,
    readOnly: false,
    options: () => ({}),
  },
)

const emit = defineEmits<{
  'update:value': [value: string]
  change: [value: string]
}>()

const containerRef = ref<HTMLElement | null>(null)
const monacoRef = shallowRef<typeof monaco | null>(null)
const editorRef = shallowRef<monaco.editor.IStandaloneCodeEditor | null>(null)
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

function syncValue(value: string) {
  const editor = editorRef.value
  if (!editor || editor.getValue() === value) {
    return
  }
  editor.setValue(value)
}

onMounted(async () => {
  if (!containerRef.value) {
    return
  }

  const monacoInstance = await loadMonaco()
  if (disposed || !containerRef.value) {
    return
  }
  editorRef.value = monacoInstance.editor.create(containerRef.value, {
    ...props.options,
    value: props.value ?? '',
    language: props.language,
    theme: props.theme,
    readOnly: props.readOnly,
    automaticLayout: true,
    minimap: {
      enabled: false,
    },
    scrollBeyondLastLine: false,
  })

  editorRef.value.onDidChangeModelContent(() => {
    const value = editorRef.value?.getValue() ?? ''
    emit('update:value', value)
    emit('change', value)
  })
})

watch(
  () => props.value,
  (value) => {
    syncValue(value ?? '')
  },
)

watch(
  () => props.language,
  (language) => {
    const monacoInstance = monacoRef.value
    const model = editorRef.value?.getModel()
    if (monacoInstance && model) {
      monacoInstance.editor.setModelLanguage(model, language)
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
  () => props.readOnly,
  (readOnly) => {
    editorRef.value?.updateOptions({ readOnly })
  },
)

watch(
  () => props.options,
  (options) => {
    editorRef.value?.updateOptions(options)
  },
)

onBeforeUnmount(() => {
  disposed = true
  const model = editorRef.value?.getModel()
  editorRef.value?.dispose()
  model?.dispose()
  editorRef.value = null
})
</script>

<template>
  <div
    ref="containerRef"
    class="monaco-editor"
    :style="containerStyle"
  />
</template>

<style scoped>
.monaco-editor {
  width: 100%;
  min-width: 0;
  border: 1px solid var(--border-color);
}
</style>
