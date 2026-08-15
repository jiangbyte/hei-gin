<!-- Author: Charlie -->
<!--
  模块选择器：下拉菜单（n-dropdown）形式，常驻顶栏（折叠/抽屉按钮之后）。
  跟随全局 light/dark 主题，无需深色变体。
-->
<script setup lang="ts">
import type { DropdownOption } from 'naive-ui'
import { computed } from 'vue'
import { useAppStore, useRouteStore } from '@/stores'
import { renderIcon } from '@/utils/icon'

const { forceLabel = false, block = false } = defineProps<{
  forceLabel?: boolean
  block?: boolean
}>()

const appStore = useAppStore()
const routeStore = useRouteStore()

const currentModule = computed(() =>
  routeStore.resourceModules.find((item) => item.id === routeStore.activeModuleId),
)

const options = computed<DropdownOption[]>(() =>
  routeStore.resourceModules.map((item) => ({
    label: item.name,
    key: item.id,
    icon: renderIcon(item.icon || 'icon-park-outline:blocks-and-arrows'),
  })),
)

const showLabel = computed(() => forceLabel || (!appStore.collapsed && !appStore.isMobile))

function handleSelect(key: string | number) {
  routeStore.setActiveModule(String(key))
}
</script>

<template>
  <n-dropdown
    v-if="options.length"
    trigger="click"
    :options="options"
    @select="handleSelect"
  >
    <button
      class="module-switch"
      :class="{ 'module-switch--block': block }"
      type="button"
      :title="currentModule?.name || 'Switch module'"
    >
      <NovaIcon :icon="currentModule?.icon || 'icon-park-outline:blocks-and-arrows'" />
      <span
        v-if="showLabel"
        class="module-switch__label"
      >
        {{ currentModule?.name || 'Modules' }}
      </span>
      <NovaIcon
        v-if="showLabel"
        icon="icon-park-outline:down"
        class="module-switch__arrow"
      />
    </button>
  </n-dropdown>
</template>

<style scoped>
.module-switch {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 32px;
  max-width: min(220px, 32vw);
  height: 32px;
  gap: 6px;
  padding: 0 8px;
  color: var(--text-color-2);
  white-space: nowrap;
  cursor: pointer;
  background: transparent;
  border: 0;
  border-radius: var(--border-radius);
  transition:
    background-color 0.2s var(--cubic-bezier-ease-in-out),
    color 0.2s var(--cubic-bezier-ease-in-out);
}

.module-switch:hover {
  color: var(--primary-color);
  background-color: var(--button-color-2-hover);
}

.module-switch__label {
  min-width: 0;
  overflow: hidden;
  font-size: 14px;
  text-overflow: ellipsis;
}

.module-switch__arrow {
  flex: 0 0 auto;
  font-size: 12px;
}

.module-switch--block {
  width: 100%;
  max-width: none;
  justify-content: flex-start;
  padding: 0 12px;
}

.module-switch--block .module-switch__label {
  flex: 1 1 auto;
  text-align: left;
}
</style>
