<!-- Author: Charlie -->

<script setup lang="ts">
import Logo from './Logo.vue'
import SidebarMenuProvider from './SidebarMenuProvider.vue'

// 与父组件通过 v-model:show 双向绑定抽屉显隐状态，移动端点击菜单入口后由父组件打开。
const showDrawer = defineModel<boolean>('show', { default: false })
</script>

<template>
  <n-drawer
    v-model:show="showDrawer"
    :width="280"
    placement="left"
    :mask-closable="true"
    :close-on-esc="true"
  >
    <n-drawer-content
      class="dark-sidebar-drawer"
      :native-scrollbar="false"
      :body-content-style="{ padding: '0' }"
    >
      <template #header>
        <Logo />
      </template>
      <n-el
        tag="div"
        class="min-h-full text-[var(--sidebar-text)]"
      >
        <!-- 移动端抽屉即侧边栏：同样使用深色菜单主题，与桌面端深色侧边栏一致 -->
        <SidebarMenuProvider>
          <slot />
        </SidebarMenuProvider>
      </n-el>
    </n-drawer-content>
  </n-drawer>
</template>

<!--
  注意：n-drawer 会把内容 teleport 到 body，scoped/:deep 选择器无法命中抽屉内部 DOM，
  因此抽屉的深色样式必须用全局选择器；配色引用 style.css 的 --sidebar-* 变量（:root / html.dark），
  teleport 后依然从 html 继承，与桌面端 aside 完全一致。
-->
<style>
.dark-sidebar-drawer {
  background-color: var(--sidebar-bg);
  color: var(--sidebar-text);
}

.dark-sidebar-drawer .n-drawer-header {
  background-color: var(--sidebar-bg);
  /* 去掉头部底部白色分隔条（深色抽屉下 n-divider-color 为浅色） */
  border-bottom: none !important;
}

.dark-sidebar-drawer .n-drawer-body {
  background-color: var(--sidebar-bg);
}

/* 抽屉内滚动条滑块浅色（深底） */
.dark-sidebar-drawer .n-scrollbar-rail {
  background-color: transparent;
}
</style>
