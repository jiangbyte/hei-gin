<!-- Author: Charlie -->

<script setup lang="ts">
import Logo from './Logo.vue'

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
      :native-scrollbar="false"
      :body-content-style="{ padding: '0' }"
    >
      <template #header>
        <Logo />
      </template>
      <n-el
        tag="div"
        class="min-h-full text-[var(--text-color-base)]"
      >
        <slot />
      </n-el>
    </n-drawer-content>
  </n-drawer>
</template>
