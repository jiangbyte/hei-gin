/** Author: Charlie */

import type { NScrollbar } from 'naive-ui'
import type { Ref } from 'vue'
import { nextTick, ref, watchEffect } from 'vue'

/**
 * 标签栏横向滚动控制。
 *
 * 负责两件事：
 * 1. 当前标签变化时，自动把激活标签滚动到可视区域；
 * 2. 用户在标签栏上滚动鼠标滚轮时，把纵向滚轮转换成横向滚动；
 * 3. 触屏设备上支持横向拖动滚动，同时避免误触发页签点击。
 */
export function useTabScroll(currentTabPath: Ref<string>) {
  // Naive UI NScrollbar 实例，TabBar.vue 会把 ref 绑定到 n-scrollbar 上。
  const scrollbar = ref<InstanceType<typeof NScrollbar>>()

  // 激活标签靠近左右边缘时预留的安全距离，避免标签贴边显示。
  const safeArea = ref(150)

  // 触屏横向拖动滚动状态。它用于样式反馈，并阻止拖动结束后的误点击。
  const touchScrolling = ref(false)

  // 滚轮节流定时器，避免高频 wheel 事件导致滚动过于抖动。
  let timer: number | undefined

  let touchPointerId: number | undefined
  let touchStartX = 0
  let touchStartY = 0
  let touchStartLeft = 0
  let isTouchDrag = false
  let suppressClickUntil = 0

  function getScrollWrapper() {
    return document.querySelector(
      '.tab-bar-scroller-wrapper .n-scrollbar-container',
    ) as HTMLElement | null
  }

  function getScrollContent() {
    return document.querySelector('.tab-bar-scroller-content') as HTMLElement | null
  }

  /**
   * 滚动到指定横向位置。
   *
   * 这里统一使用 smooth 行为，让点击标签、关闭标签后的自动定位更自然。
   */
  function handleTabSwitch(distance: number) {
    scrollbar.value?.scrollTo({
      left: distance,
      behavior: 'smooth',
    })
  }

  /**
   * 把当前激活标签滚动到可视区域。
   *
   * DOM 更新发生在路由和标签状态更新之后，所以需要 nextTick 后再读取元素位置。
   */
  function scrollToCurrentTab() {
    nextTick(() => {
      // 当前激活标签，TabBar.vue 会在普通标签上写入 data-tab-path。
      const currentTabElement = document.querySelector(
        `[data-tab-path="${currentTabPath.value}"]`,
      ) as HTMLElement | null

      // n-scrollbar 内部真正产生滚动的容器。
      const tabBarScrollWrapper = getScrollWrapper()

      // 标签内容容器，用于读取右侧 padding，避免计算滚动距离时忽略内边距。
      const tabBarScrollContent = getScrollContent()

      // DOM 不完整时直接跳过，避免初始化阶段或组件卸载阶段报错。
      if (!currentTabElement || !tabBarScrollWrapper || !tabBarScrollContent) {
        return
      }

      // 当前标签相对内容容器的左偏移。
      const tabLeft = currentTabElement.offsetLeft

      // 当前滚动容器已滚动的横向距离。
      const tabBarLeft = tabBarScrollWrapper.scrollLeft

      // 滚动容器可视宽度。
      const wrapperWidth = tabBarScrollWrapper.getBoundingClientRect().width

      // 当前标签自身宽度。
      const tabWidth = currentTabElement.getBoundingClientRect().width

      // 内容容器右侧 padding，确保最后一个标签滚动到右侧时不会被操作区遮挡。
      const containerPR = Number.parseFloat(
        window.getComputedStyle(tabBarScrollContent).paddingRight,
      )

      // 标签右侧超出可视区域时，向右滚动到标签完整可见。
      if (tabLeft + tabWidth + safeArea.value + containerPR > wrapperWidth + tabBarLeft) {
        handleTabSwitch(tabLeft + tabWidth + containerPR - wrapperWidth + safeArea.value)
      } else if (tabLeft - safeArea.value < tabBarLeft) {
        // 标签左侧超出可视区域时，向左滚动到标签完整可见。
        handleTabSwitch(tabLeft - safeArea.value)
      }
    })
  }

  /**
   * 鼠标滚轮横向滚动标签栏。
   *
   * 大多数鼠标只有纵向滚轮，这里把 deltaY 转换为横向滚动；如果本身是横向滚动手势则不处理。
   */
  function onWheel(e: WheelEvent) {
    e.preventDefault()
    if (Math.abs(e.deltaY) <= Math.abs(e.deltaX)) {
      return
    }

    window.clearTimeout(timer)
    timer = window.setTimeout(() => {
      scrollbar.value?.scrollBy({
        left: e.deltaY > 0 ? 400 : -400,
        behavior: 'smooth',
      })
    }, 40)
  }

  /**
   * 触屏横向拖动滚动标签栏。
   *
   * 页签本身支持拖拽排序；触屏上如果不区分“横向滑动”和“拖拽排序”，很容易滑动页签栏时
   * 误触发排序或点击。这里仅处理 touch pointer，桌面鼠标仍保留原有拖拽排序体验。
   */
  function onPointerDown(e: PointerEvent) {
    if (e.pointerType !== 'touch') {
      return
    }

    const wrapper = getScrollWrapper()
    if (!wrapper) {
      return
    }

    touchPointerId = e.pointerId
    touchStartX = e.clientX
    touchStartY = e.clientY
    touchStartLeft = wrapper.scrollLeft
    isTouchDrag = false
  }

  function onPointerMove(e: PointerEvent) {
    if (e.pointerType !== 'touch' || e.pointerId !== touchPointerId) {
      return
    }

    const wrapper = getScrollWrapper()
    if (!wrapper) {
      return
    }

    const deltaX = e.clientX - touchStartX
    const deltaY = e.clientY - touchStartY

    if (!isTouchDrag) {
      if (Math.abs(deltaX) < 8) {
        return
      }

      // 纵向移动更明显时交给页面正常滚动，避免标签栏拦截移动端页面上下滑动。
      if (Math.abs(deltaY) > Math.abs(deltaX)) {
        resetTouchScroll()
        return
      }

      isTouchDrag = true
      touchScrolling.value = true
    }

    e.preventDefault()
    wrapper.scrollLeft = touchStartLeft - deltaX
  }

  function onPointerEnd(e: PointerEvent) {
    if (e.pointerType !== 'touch' || e.pointerId !== touchPointerId) {
      return
    }

    if (isTouchDrag) {
      suppressClickUntil = window.performance.now() + 250
    }

    resetTouchScroll()
  }

  function onClickCapture(e: MouseEvent) {
    if (window.performance.now() > suppressClickUntil) {
      return
    }

    e.preventDefault()
    e.stopPropagation()
  }

  function resetTouchScroll() {
    touchPointerId = undefined
    isTouchDrag = false
    touchScrolling.value = false
  }

  // 当前标签变化时自动修正滚动位置。
  watchEffect(() => {
    if (currentTabPath.value) {
      scrollToCurrentTab()
    }
  })

  // 暴露给 TabBar.vue 使用的滚动实例、滚轮处理函数和可调安全距离。
  return {
    scrollbar,
    touchScrolling,
    onWheel,
    onPointerDown,
    onPointerMove,
    onPointerEnd,
    onClickCapture,
    safeArea,
    handleTabSwitch,
  }
}
