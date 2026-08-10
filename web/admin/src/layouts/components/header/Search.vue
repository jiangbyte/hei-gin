<!-- Author: Charlie -->

<script setup lang="ts">
import { computed, nextTick, ref, watchEffect } from 'vue'
import { useMagicKeys } from '@vueuse/core'
import { useRouter } from 'vue-router'
import { useBoolean } from '@/hooks'
import { useAppStore, useRouteStore } from '@/stores'

const appStore = useAppStore()
const routeStore = useRouteStore()
const router = useRouter()

// 搜索弹窗中的可跳转项，value 固定使用路由 path，code 用于辅助识别权限/菜单编码。
interface SearchOption {
  label: string
  value: string
  code: string
  icon?: string
}

const searchValue = ref('')
const selectedIndex = ref(0)
const scrollbarRef = ref()
// 控制全局搜索弹窗显隐，统一通过 useBoolean 暴露语义化操作，减少模板中的状态细节。
const {
  bool: showModal,
  setTrue: openModal,
  setFalse: closeModal,
  toggle: toggleModal,
} = useBoolean(false)
// 标记当前高亮是否由键盘触发：键盘导航时忽略 mouseenter，避免鼠标停留位置抢回高亮。
const {
  bool: keyboardFlag,
  setTrue: setKeyboardTrue,
  setFalse: setKeyboardFalse,
} = useBoolean(false)

const { ctrl_k, arrowup, arrowdown, enter, escape } = useMagicKeys({
  passive: false,
  onEventFired(e) {
    if (e.ctrlKey && e.key.toLowerCase() === 'k' && e.type === 'keydown') {
      e.preventDefault()
    }
  },
})

// Ctrl + K 是全局搜索入口；使用 watchEffect 监听组合键状态，保证键盘按下时立即切换弹窗。
watchEffect(() => {
  if (ctrl_k.value) {
    toggleModal()
  }
})

// 根据输入实时从动态路由表中过滤菜单/页面。
// 仅允许启用、可见、可导航的 MENU/PAGE 进入搜索结果，避免跳转到按钮权限或隐藏资源。
const options = computed<SearchOption[]>(() => {
  const keyword = searchValue.value.trim().toLowerCase()
  if (!keyword) {
    return []
  }

  return routeStore.rowRoutes
    .map((item) => ({
      item,
      label: item.name,
    }))
    .filter(({ item, label }) => {
      if (
        item.status !== 'ENABLED' ||
        !item.is_visible ||
        (item.resource_type !== 'MENU' && item.resource_type !== 'PAGE') ||
        !item.path
      ) {
        return false
      }

      return [item.name, label, item.path, item.code].some((value) =>
        value.toLowerCase().includes(keyword),
      )
    })
    .map(({ item, label }) => ({
      label,
      value: item.path!,
      code: item.code,
      icon: item.icon ?? undefined,
    }))
})

// 关闭弹窗时重置输入和选中项，确保下次打开不会继承上一次搜索上下文。
function handleClose() {
  searchValue.value = ''
  selectedIndex.value = 0
  closeModal()
}

// 输入变化后把高亮归位到第一项，避免旧索引超出新结果列表长度。
function handleInputChange() {
  selectedIndex.value = 0
}

// 选择结果后立即关闭弹窗并跳转；nextTick 后再次清空输入，兼容关闭动画期间的输入回显。
function handleSelect(value: string) {
  handleClose()
  router.push(value)
  nextTick(() => {
    searchValue.value = ''
  })
}

// 弹窗打开后统一处理 Esc、上下方向键、Enter。
// options 为空时不处理导航键，避免对空列表产生无意义索引变更。
watchEffect(() => {
  if (!showModal.value) {
    return
  }

  if (escape.value) {
    handleClose()
    return
  }

  if (!options.value.length) {
    return
  }

  setKeyboardTrue()
  if (arrowup.value) {
    handleArrowup()
  }
  if (arrowdown.value) {
    handleArrowdown()
  }
  if (enter.value) {
    handleEnter()
  }
})

// 向上导航支持首尾循环，便于用户连续按键快速浏览结果。
function handleArrowup() {
  selectedIndex.value =
    selectedIndex.value === 0 ? options.value.length - 1 : selectedIndex.value - 1
  handleScroll(selectedIndex.value)
}

// 向下导航支持末尾回到第一项，与常见命令面板交互保持一致。
function handleArrowdown() {
  selectedIndex.value =
    selectedIndex.value === options.value.length - 1 ? 0 : selectedIndex.value + 1
  handleScroll(selectedIndex.value)
}

// 键盘移动高亮时同步滚动容器，使当前项保持在可视区域内。
// keepIndex 表示允许高亮项位于列表前 5 项范围内，超过后再滚动，减少频繁滚动造成的抖动。
function handleScroll(currentIndex: number) {
  const keepIndex = 5
  const optionHeight = 70
  const distance =
    currentIndex * optionHeight > keepIndex * optionHeight
      ? currentIndex * optionHeight - keepIndex * optionHeight
      : 0
  scrollbarRef.value?.scrollTo({ top: distance })
}

// Enter 只在当前索引命中结果时触发跳转，防止快速输入导致结果变化后访问空项。
function handleEnter() {
  const target = options.value[selectedIndex.value]
  if (target) {
    handleSelect(target.value)
  }
}

// 非键盘导航状态下才允许鼠标悬停更新高亮，兼顾鼠标浏览和键盘连续操作两种场景。
function handleMouseEnter(index: number) {
  if (!keyboardFlag.value) {
    selectedIndex.value = index
  }
}
</script>

<template>
  <CommonWrapper
    class="px-2"
    @click="openModal"
  >
    <NovaIcon icon="icon-park-outline:search" />
    <n-tag
      v-if="!appStore.isMobile"
      round
      size="small"
      class="font-mono cursor-pointer"
    >
      CtrlK
    </n-tag>
  </CommonWrapper>

  <n-modal
    v-model:show="showModal"
    class="w-560px fixed top-60px inset-x-0 max-w-full"
    size="small"
    preset="card"
    :segmented="{ content: true, footer: true }"
    :closable="false"
    @after-leave="handleClose"
  >
    <template #header>
      <n-input
        v-model:value="searchValue"
        :placeholder="'搜索菜单、路径或编码'"
        clearable
        size="large"
        @input="handleInputChange"
      >
        <template #prefix>
          <NovaIcon icon="icon-park-outline:search" />
        </template>
      </n-input>
    </template>

    <n-scrollbar
      ref="scrollbarRef"
      class="h-450px"
    >
      <ul
        v-if="options.length"
        class="flex flex-col gap-8px p-8px p-r-3"
      >
        <n-el
          v-for="(option, index) in options"
          :key="option.value"
          tag="li"
          role="option"
          class="cursor-pointer shadow h-62px transition-colors"
          :class="{
            'text-[var(--base-color)] bg-[var(--primary-color-hover)]': index === selectedIndex,
          }"
          @click="handleSelect(option.value)"
          @mouseenter="handleMouseEnter(index)"
          @mousemove="setKeyboardFalse"
        >
          <div class="grid grid-rows-2 grid-cols-[40px_1fr_30px] h-full p-2">
            <div class="row-span-2 place-self-center">
              <NovaIcon :icon="option.icon" />
            </div>
            <span>{{ option.label }}</span>
            <NovaIcon
              icon="icon-park-outline:right"
              class="row-span-2 place-self-center"
            />
            <span class="op-70">{{ option.value }} / {{ option.code }}</span>
          </div>
        </n-el>
      </ul>

      <n-empty
        v-else
        size="large"
        class="h-450px flex-center"
        :description="'暂无结果'"
      />
    </n-scrollbar>

    <template #footer>
      <n-flex class="items-center">
        <span class="flex-y-center gap-1">
          <n-tag
            size="small"
            round
          >Enter</n-tag>
          <span>选择</span>
        </span>
        <span class="flex-y-center gap-1">
          <n-tag
            size="small"
            round
          >↑</n-tag>
          <n-tag
            size="small"
            round
          >↓</n-tag>
          <span>导航</span>
        </span>
        <span class="flex-y-center gap-1">
          <n-tag
            size="small"
            round
          >Esc</n-tag>
          <span>关闭</span>
        </span>
      </n-flex>
    </template>
  </n-modal>
</template>
