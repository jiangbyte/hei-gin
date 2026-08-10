<!-- Author: Charlie -->

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ProLayout, useLayoutMenu } from 'pro-naive-ui'
import { useAppStore, useRouteStore } from '@/stores'
import {
  BackTop,
  Breadcrumb,
  CollapaseButton,
  FullScreen,
  Logo,
  MobileDrawer,
  ModuleSwitch,
  Notices,
  Search,
  TabBar,
  UserCenter,
} from './components'
import Content from './Content.vue'

const appStore = useAppStore()
const routeStore = useRouteStore()

const menus = computed(() => {
  return routeStore.menus
})

const { layout, activeKey } = useLayoutMenu({
  mode: 'vertical',
  accordion: true,
  menus,
} as never)

watch(
  () => routeStore.currentMenuPath,
  (currentMenuPath) => {
    activeKey.value = currentMenuPath
  },
  { immediate: true },
)

const showMobileDrawer = ref(false)
</script>

<template>
  <ProLayout
    v-model:collapsed="appStore.collapsed"
    mode="vertical"
    :is-mobile="appStore.isMobile"
    :show-logo="!appStore.isMobile"
    :show-footer="false"
    show-tabbar
    nav-fixed
    show-nav
    show-sidebar
    :nav-height="60"
    :tabbar-height="45"
    :sidebar-width="240"
    :sidebar-collapsed-width="64"
  >
    <template #logo>
      <Logo sidebar />
    </template>

    <template #nav-left>
      <template v-if="appStore.isMobile">
        <div class="h-full flex-y-center gap-3 p-x-sm">
          <CommonWrapper @click="showMobileDrawer = true">
            <NovaIcon icon="icon-park-outline:hamburger-button" />
          </CommonWrapper>
        </div>
      </template>
      <template v-else>
        <div class="h-full flex-y-center gap-1 p-x-sm">
          <CollapaseButton />
          <ModuleSwitch />
          <Breadcrumb />
        </div>
      </template>
    </template>

    <template #nav-right>
      <div class="h-full flex-y-center gap-1 p-x-xl">
        <template v-if="appStore.isMobile">
          <Search />
          <Notices />
          <DarkModeSwitch />
          <UserCenter />
        </template>
        <template v-else>
          <Search />
          <Notices />
          <FullScreen />
          <DarkModeSwitch />
          <UserCenter />
        </template>
      </div>
    </template>

    <template #sidebar>
      <n-scrollbar class="sidebar-menu-scrollbar">
        <n-menu
          v-bind="layout.verticalMenuProps"
          :collapsed-width="64"
        />
      </n-scrollbar>
    </template>

    <template #sidebar-extra>
      <n-scrollbar class="flex-[1_0_0]">
        <n-menu
          v-bind="layout.verticalExtraMenuProps"
          :collapsed-width="64"
        />
      </n-scrollbar>
    </template>

    <template #tabbar>
      <TabBar />
    </template>

    <Content />
    <BackTop class="z-999" />

    <MobileDrawer v-model:show="showMobileDrawer">
      <div class="mt-3 mb-2 px-3">
        <ModuleSwitch
          force-label
          block
        />
      </div>
      <n-menu
        v-bind="layout.verticalMenuProps"
        :collapsed="false"
      />
    </MobileDrawer>
  </ProLayout>
</template>

<style scoped>
:deep(.n-pro-layout__sidebar) {
  min-height: 0;
  overflow: hidden;
}

.sidebar-menu-scrollbar {
  min-height: 0;
  flex: 1 1 0;
}
</style>
