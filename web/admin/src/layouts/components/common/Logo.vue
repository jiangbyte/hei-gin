<!-- Author: Charlie -->

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores'

const router = useRouter()
const appStore = useAppStore()
const name = import.meta.env.VITE_APP_TITLE
const homePath = import.meta.env.VITE_HOME_PATH

const { sidebar = false } = defineProps<{
  sidebar?: boolean
}>()

// 桌面端跟随侧边栏折叠状态，移动端始终显示系统名称。
const hiddenLogoText = computed(() => sidebar && !appStore.isMobile && appStore.collapsed)
</script>

<template>
  <n-el
    tag="button"
    class="h-12 min-w-0 inline-flex items-center gap-3 border-0 bg-transparent p-0 text-left text-[var(--text-color-base)] cursor-pointer"
    :class="sidebar ? 'w-full items-center h-full justify-center px-3' : ''"
    type="button"
    @click="router.push(homePath)"
  >
    <span
      class="h-38px w-38px shrink-0 inline-flex items-center justify-center rounded-2 bg-[var(--primary-color)] text-white"
    >
      <NovaIcon
        icon="icon-park-outline:code-computer"
        :size="24"
      />
    </span>
    <span
      v-if="!hiddenLogoText"
      class="min-w-0 flex flex-col"
    >
      <span class="truncate text-base font-700 leading-5">{{ name }}</span>
    </span>
  </n-el>
</template>
