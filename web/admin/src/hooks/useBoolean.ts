/** Author: Charlie */

import { ref } from 'vue'

/**
 * 管理布尔状态的通用组合式函数。
 *
 * 适合弹窗显隐、抽屉展开、开关状态等简单布尔场景，避免组件里重复声明 setTrue/setFalse/toggle。
 */
export function useBoolean(initValue = false) {
  // 当前布尔状态。使用 ref 方便在模板和 watch 中保持响应式。
  const bool = ref(initValue)

  // 显式设置布尔值，供需要根据外部条件赋值的场景使用。
  function setBool(value: boolean) {
    bool.value = value
  }

  // 快捷设置为 true，常用于打开弹窗或启用状态。
  function setTrue() {
    setBool(true)
  }

  // 快捷设置为 false，常用于关闭弹窗或禁用状态。
  function setFalse() {
    setBool(false)
  }

  // 在 true/false 之间切换。
  function toggle() {
    setBool(!bool.value)
  }

  // 返回状态和操作方法，调用方可按需解构。
  return {
    bool,
    setBool,
    setTrue,
    setFalse,
    toggle,
  }
}
