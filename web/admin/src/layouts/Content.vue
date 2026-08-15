<!-- Author: Charlie -->

<script setup lang="ts">
import { useAppStore, useRouteStore } from '@/stores'

const appStore = useAppStore()
const routeStore = useRouteStore()
</script>

<template>
  <n-el
    class="h-full min-h-0 overflow-hidden p-8px"
  >
    <router-view v-slot="{ Component, route }">
      <transition
        name="fade-slide"
        mode="out-in"
      >
        <keep-alive :include="routeStore.cacheRoutes">
          <component
            :is="Component"
            v-if="appStore.loadFlag"
            :key="`${route.path}::${String(route.query.id ?? route.params.id ?? '')}`"
          />
        </keep-alive>
      </transition>
    </router-view>
  </n-el>
</template>
