<!-- Author: Charlie -->

<script setup lang="ts">
import { Icon } from '@iconify/vue/offline'

withDefaults(
  defineProps<{
    icon?: string
    color?: string
    size?: number
    depth?: 1 | 2 | 3 | 4 | 5
  }>(),
  {
    icon: '',
    color: '',
    size: 18,
    depth: undefined,
  },
)
</script>

<template>
  <n-icon
    v-if="icon"
    :size="size"
    :depth="depth"
    :color="color"
  >
    <Icon :icon="icon" />
  </n-icon>
</template>
