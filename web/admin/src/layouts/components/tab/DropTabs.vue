<!-- Author: Charlie -->

<script setup lang="ts">
import type { DropdownOption } from 'naive-ui'
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTabStore } from '@/stores'
import { getRouteTitle } from '@/stores/route'
import { renderIcon } from '@/utils/icon'

const tabStore = useTabStore()
const router = useRouter()

// 下拉菜单展示所有可访问页签，包括固定页签和普通页签，作为横向页签溢出时的快速入口。
const options = computed<DropdownOption[]>(() =>
  tabStore.allTabs.map((route) => ({
    label: getRouteTitle(route),
    key: route.fullPath,
    icon: renderIcon(route.meta.icon ?? undefined),
  })),
)

// 选中下拉项后使用 fullPath 跳转，保留查询参数和哈希，确保与页签实际地址一致。
function handleDropTabs(key: string | number) {
  router.push(String(key))
}
</script>

<template>
  <n-dropdown
    :options="options"
    trigger="click"
    size="small"
    @select="handleDropTabs"
  >
    <CommonWrapper>
      <NovaIcon icon="icon-park-outline:application-menu" />
    </CommonWrapper>
  </n-dropdown>
</template>
