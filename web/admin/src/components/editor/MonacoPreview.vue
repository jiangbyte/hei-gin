<!-- Author: Charlie -->

<script setup lang="ts">
import type * as monaco from 'monaco-editor'
import MonacoEditor from './MonacoEditor.vue'

withDefaults(
  defineProps<{
    value?: string | null
    language?: string
    theme?: string
    height?: string | number
    options?: monaco.editor.IStandaloneEditorConstructionOptions
  }>(),
  {
    value: '',
    language: 'typescript',
    theme: 'vs',
    height: 360,
    options: () => ({}),
  },
)
</script>

<template>
  <MonacoEditor
    :value="value"
    :language="language"
    :theme="theme"
    :height="height"
    :options="options"
    read-only
  />
</template>
