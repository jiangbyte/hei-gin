<!-- Author: Charlie -->

<script setup lang="ts">
interface Props {
  show: boolean
  text?: string
}
const { text = '水印' } = defineProps<Props>()
</script>

<template>
  <n-watermark
    v-if="show"
    :content="text"
    cross
    fullscreen
    :font-size="16"
    :line-height="16"
    :width="384"
    :height="384"
    :x-offset="12"
    :y-offset="60"
    :rotate="-15"
  />
</template>
