/** Author: Charlie */

import { useSyncExternalStore } from 'react'
import { dictTreeState, subscribeDict } from '@/utils/dict'

/** 订阅 dictTreeState，字典刷新后触发重渲染（对齐 admin 对 shallowRef 的依赖）。 */
export function useDict() {
  return useSyncExternalStore(
    subscribeDict,
    () => dictTreeState.value,
    () => dictTreeState.value,
  )
}
