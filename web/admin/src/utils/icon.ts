/** Author: Charlie */

import { Icon } from '@iconify/vue/offline'
import { NIcon } from 'naive-ui'
import { h } from 'vue'

/**
 * 生成 Naive UI 菜单、下拉菜单等组件可使用的图标渲染函数。
 *
 * Iconify 图标名为空时返回 undefined，让调用方自然地不渲染图标插槽。
 */
export function renderIcon(icon?: string, size = 18) {
  if (!icon) {
    return undefined
  }

  // 返回函数而不是 VNode，符合 Naive UI option.icon 的懒渲染约定。
  return () => h(Icon, { icon, width: size, height: size })
}

export function renderButtonIcon(icon: string, size = 16) {
  return {
    icon: () =>
      h(NIcon, null, {
        default: () => h(Icon, { icon, width: size, height: size }),
      }),
  }
}
