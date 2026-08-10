<!-- Author: Charlie -->

<script setup lang="ts">
import type { DropdownOption } from 'naive-ui'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useDraggable } from 'vue-draggable-plus'
import { useTabScroll } from '@/hooks/useTabScroll'
import { useAppStore, useTabStore } from '@/stores'
import { getRouteTitle } from '@/stores/route'
import type { AppTab } from '@/stores/tab'
import { renderIcon } from '@/utils/icon'
import DropTabs from './DropTabs.vue'
import Reload from './Reload.vue'

const tabStore = useTabStore()
const appStore = useAppStore()
const router = useRouter()
// useTabScroll 将鼠标滚轮转换成横向滚动，并在当前页签变化时滚动到可视区域。
const {
  scrollbar,
  touchScrolling,
  onWheel,
  onPointerDown,
  onPointerMove,
  onPointerEnd,
  onClickCapture,
} = useTabScroll(computed(() => tabStore.currentTabPath))
void scrollbar

const el = ref<HTMLElement>()

// 普通页签支持拖拽排序。固定页签单独渲染，不参与拖拽，避免用户改变固定入口位置。
useDraggable(
  el,
  computed({
    // draggable-plus 会直接修改数组顺序，这里通过 computed setter 回写 Pinia，保持 store 为唯一数据源。
    get: () => tabStore.tabs,
    set: (value) => {
      tabStore.tabs = value
    },
  }),
  {
    animation: 150,
    ghostClass: 'ghost',
    delay: 180,
    delayOnTouchOnly: true,
    touchStartThreshold: 8,
    fallbackTolerance: 6,
  },
)

// 当前激活页签可能来自固定页签或普通页签，因此需要从 allTabs 中查找。
const currentTab = computed(() =>
  tabStore.allTabs.find((item) => item.fullPath === tabStore.currentTabPath),
)

// 固定页签不可关闭，也不能作为“关闭其他/左侧/右侧”的操作基准。
const isCurrentAffixTab = computed(() => Boolean(currentTab.value?.meta.is_affix))

// 页签右侧操作菜单。disabled 规则集中在这里，模板只负责展示，避免交互规则散落在视图层。
const options = computed<DropdownOption[]>(() => {
  const disabledCurrent = !currentTab.value || isCurrentAffixTab.value
  const disabledNormal = !tabStore.tabs.length

  return [
    {
      label: '刷新',
      key: 'reload',
      icon: renderIcon('icon-park-outline:redo'),
    },
    {
      label: '关闭当前',
      key: 'closeCurrent',
      icon: renderIcon('icon-park-outline:close'),
      disabled: disabledCurrent,
    },
    {
      label: '关闭其他',
      key: 'closeOther',
      icon: renderIcon('icon-park-outline:delete-four'),
      disabled: disabledCurrent || disabledNormal,
    },
    {
      label: '关闭左侧',
      key: 'closeLeft',
      icon: renderIcon('icon-park-outline:to-left'),
      disabled: disabledCurrent || disabledNormal,
    },
    {
      label: '关闭右侧',
      key: 'closeRight',
      icon: renderIcon('icon-park-outline:to-right'),
      disabled: disabledCurrent || disabledNormal,
    },
    {
      label: '关闭全部',
      key: 'closeAll',
      icon: renderIcon('icon-park-outline:fullwidth'),
      disabled: disabledNormal,
    },
  ]
})

// 点击页签时按 fullPath 跳转，保留 query/hash，保证页签切换后回到打开时的完整地址。
function handleTab(route: AppTab) {
  router.push(route.fullPath)
}

// 关闭按钮位于页签内部，需要阻止冒泡，否则关闭时会先触发页签点击导致路由跳转。
function handleCloseTab(e: MouseEvent, fullPath: string) {
  e.stopPropagation()
  tabStore.closeTab(fullPath)
}

// 操作菜单基于当前激活页签执行。具体关闭策略由 tabStore 统一处理，组件只做命令分发。
function handleSelect(key: string | number) {
  const path = currentTab.value?.fullPath
  if (!path) {
    return
  }

  const handleFn: Record<string, () => void> = {
    reload: () => appStore.reloadPage(),
    closeCurrent: () => tabStore.closeTab(path),
    closeOther: () => tabStore.closeOtherTabs(path),
    closeLeft: () => tabStore.closeLeftTabs(path),
    closeRight: () => tabStore.closeRightTabs(path),
    closeAll: () => tabStore.closeAllTabs(),
  }
  handleFn[String(key)]?.()
}
</script>

<template>
  <div class="relative flex h-full w-full min-w-0 overflow-hidden px-2">
    <n-scrollbar
      ref="scrollbar"
      class="relative h-full flex-1 min-w-0 tab-bar-scroller-wrapper"
      content-class="h-full tab-bar-scroller-content"
      :x-scrollable="true"
      @wheel="onWheel"
    >
      <div
        class="p-l-2 inline-flex h-full min-w-full items-center gap-1 relative tab-bar-scroll-surface"
        :class="{ 'is-touch-scrolling': touchScrolling }"
        @pointerdown="onPointerDown"
        @pointermove="onPointerMove"
        @pointerup="onPointerEnd"
        @pointercancel="onPointerEnd"
        @click.capture="onClickCapture"
      >
        <div class="flex items-center gap-1">
          <n-tag
            v-for="item in tabStore.affixTabs"
            :key="item.fullPath"
            class="tab-item"
            :type="tabStore.currentTabPath === item.fullPath ? 'primary' : 'default'"
            @click="handleTab(item)"
          >
            <template #icon>
              <NovaIcon
                v-if="item.meta.icon"
                :icon="item.meta.icon"
              />
            </template>
            {{ getRouteTitle(item) }}
          </n-tag>
        </div>
        <div
          ref="el"
          class="flex items-center gap-1 flex-1"
        >
          <n-tag
            v-for="item in tabStore.tabs"
            :key="item.fullPath"
            closable
            class="tab-item"
            :data-tab-path="item.fullPath"
            :type="tabStore.currentTabPath === item.fullPath ? 'primary' : 'default'"
            @close="handleCloseTab($event, item.fullPath)"
            @click="handleTab(item)"
          >
            <template #icon>
              <NovaIcon
                v-if="item.meta.icon"
                :icon="item.meta.icon"
              />
            </template>
            {{ getRouteTitle(item) }}
          </n-tag>
        </div>
      </div>
    </n-scrollbar>
    <n-el class="flex h-full shrink-0 items-center gap-1 bg-[var(--card-color)]">
      <Reload />
      <n-dropdown
        :options="options"
        trigger="click"
        placement="bottom-start"
        @select="handleSelect"
      >
        <CommonWrapper>
          <NovaIcon icon="icon-park-outline:setting-two" />
        </CommonWrapper>
      </n-dropdown>
      <DropTabs />
    </n-el>
  </div>
</template>

<style scoped>
.tab-item {
  cursor: pointer;
  white-space: nowrap;
}

.tab-bar-scroll-surface {
  touch-action: pan-y;
  overscroll-behavior-inline: contain;
}

.tab-bar-scroll-surface.is-touch-scrolling {
  user-select: none;
}
</style>
