<!-- Author: Charlie -->

<script setup lang="ts">
import { useMagicKeys } from '@vueuse/core'
import { useAppStore } from '@/stores'

const appStore = useAppStore()

// 接管浏览器默认 F11 行为，统一走 appStore.toggleFullScreen，保证按钮和快捷键状态来源一致。
useMagicKeys({
  passive: false,
  onEventFired(e) {
    if (e.key === 'F11' && e.type === 'keydown') {
      e.preventDefault()
      appStore.toggleFullScreen()
    }
  },
})
</script>

<template>
  <n-tooltip
    placement="bottom"
    trigger="hover"
  >
    <template #trigger>
      <CommonWrapper @click="appStore.toggleFullScreen">
        <NovaIcon
          :icon="
            appStore.fullScreen ? 'icon-park-outline:off-screen' : 'icon-park-outline:full-screen'
          "
        />
      </CommonWrapper>
    </template>
    全屏
  </n-tooltip>
</template>
