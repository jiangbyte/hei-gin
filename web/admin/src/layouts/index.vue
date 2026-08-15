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
import SidebarMenuProvider from './components/common/SidebarMenuProvider.vue'

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
          <ModuleSwitch />
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
      <SidebarMenuProvider>
        <n-scrollbar class="sidebar-menu-scrollbar">
          <n-menu
            v-bind="layout.verticalMenuProps"
            :collapsed-width="64"
          />
        </n-scrollbar>
      </SidebarMenuProvider>
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

/* 滚动容器：SidebarMenuProvider 已撑满 sidebar，这里让 scrollbar 占满并滚动 */
.sidebar-menu-scrollbar {
  min-height: 0;
  flex: 1 1 0;
}

/*
 * 深色侧边栏：仅作用于侧边栏容器（aside），不影响顶栏/内容区。
 * --pro-layout-color 同时用于顶栏背景，因此在此作用域内单独覆盖。
 * 配色引用 style.css 中的模式感知变量：浅色模式深蓝灰、暗黑模式更深一档并与内容协调。
 */
:deep(.n-pro-layout__aside) {
  --pro-layout-color: var(--sidebar-bg);
  --pro-layout-border-color: var(--sidebar-border);
  /* 侧边栏内滚动条浅色滑块（深底） */
  --app-scrollbar-thumb-color: var(--sidebar-scrollbar-thumb);
  --app-scrollbar-thumb-color-hover: var(--sidebar-scrollbar-thumb-hover);
}
</style>
