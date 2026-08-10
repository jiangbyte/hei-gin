<!-- Author: Charlie -->

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getRouteTitle } from '@/stores/route'

const router = useRouter()
const route = useRoute()

// 只展示配置了 meta.name 的匹配路由，过滤掉纯布局层或无标题的中间路由。
const routes = computed(() => route.matched.filter((item) => item.meta.name))
</script>

<template>
  <TransitionGroup
    name="list"
    tag="ul"
    class="flex items-center gap-2"
  >
    <template
      v-for="(item, index) in routes"
      :key="item.path"
    >
      <li
        v-if="index > 0"
        class="breadcrumb-separator"
      >
        /
      </li>
      <n-el
        tag="li"
        class="flex items-center gap-1 breadcrumb-item cursor-pointer"
        @click="router.push(item.path || '/')"
      >
        <NovaIcon :icon="item.meta.icon ?? undefined" />
        <span class="whitespace-nowrap">
          {{ getRouteTitle(item) }}
        </span>
      </n-el>
    </template>
  </TransitionGroup>
</template>

<style scoped>
.breadcrumb-item:hover {
  color: var(--primary-color);
}

.breadcrumb-separator {
  color: var(--text-color-3);
  user-select: none;
}
</style>
