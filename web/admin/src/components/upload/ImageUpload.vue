<!-- Author: Charlie -->

<script setup lang="ts">
import FileUpload from './FileUpload.vue'
import type { UploadedFileValueType } from '@/utils'

withDefaults(
  defineProps<{
    value?: string | null
    valueType?: UploadedFileValueType
  }>(),
  {
    value: '',
    valueType: 'object_name',
  },
)

defineEmits<{
  'update:value': [value: string]
  uploaded: [file: any]
}>()
</script>

<template>
  <FileUpload
    :value="value"
    :value-type="valueType"
    accept="image/jpeg,image/png,image/webp,image/gif"
    preview="image"
    @update:value="$emit('update:value', $event)"
    @uploaded="$emit('uploaded', $event)"
  />
</template>
