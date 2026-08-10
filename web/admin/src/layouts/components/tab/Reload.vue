<!-- Author: Charlie -->

<script setup lang="ts">
import { ref } from 'vue'
import { useAppStore } from '@/stores'

const appStore = useAppStore()
const loading = ref(false)

// 刷新当前路由视图时给图标一个最短 800ms 的旋转反馈，避免页面很快重载时用户感知不到点击结果。
function handleReload() {
  loading.value = true
  appStore.reloadPage()
  window.setTimeout(() => {
    loading.value = false
  }, 800)
}
</script>

<template>
  <n-tooltip
    placement="bottom"
    trigger="hover"
  >
    <template #trigger>
      <CommonWrapper @click="handleReload">
        <NovaIcon
          icon="icon-park-outline:refresh"
          :class="{ 'animate-spin': loading }"
        />
      </CommonWrapper>
    </template>
    刷新
  </n-tooltip>
</template>
