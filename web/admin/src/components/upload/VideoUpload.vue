<!-- Author: Charlie -->

<script setup lang="ts">
import FileUpload from './FileUpload.vue'

defineProps<{
  value?: string | null
}>()

defineEmits<{
  'update:value': [value: string]
  uploaded: [file: any]
}>()
</script>

<template>
  <FileUpload
    :value="value"
    accept="video/mp4,video/webm,video/quicktime,.mp4,.webm,.mov"
    preview="video"
    @update:value="$emit('update:value', $event)"
    @uploaded="$emit('uploaded', $event)"
  />
</template>
