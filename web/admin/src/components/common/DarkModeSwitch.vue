<!-- Author: Charlie -->

<script setup lang="ts">
import { Icon } from '@iconify/vue/offline'
import { NFlex } from 'naive-ui'
import { computed, h } from 'vue'
import { useAppStore } from '@/stores'

const appStore = useAppStore()

const options = computed(() => [
  {
    label: '浅色',
    value: 'light',
    icon: 'icon-park-outline:sun-one',
  },
  {
    label: '深色',
    value: 'dark',
    icon: 'icon-park-outline:moon',
  },
  {
    label: '系统',
    value: 'auto',
    icon: 'icon-park-outline:laptop-computer',
  },
])

function renderLabel(option: { label: string; icon: string }) {
  return h(
    NFlex,
    { align: 'center', wrap: false },
    {
      default: () => [h(Icon, { icon: option.icon }), option.label],
    },
  )
}
</script>

<template>
  <n-popselect
    :value="appStore.storeColorMode"
    :render-label="renderLabel"
    :options="options"
    trigger="click"
    @update:value="appStore.setColorMode"
  >
    <CommonWrapper>
      <NovaIcon
        :icon="
          appStore.storeColorMode === 'dark'
            ? 'icon-park-outline:moon'
            : appStore.storeColorMode === 'light'
              ? 'icon-park-outline:sun-one'
              : 'icon-park-outline:laptop-computer'
        "
      />
    </CommonWrapper>
  </n-popselect>
</template>
