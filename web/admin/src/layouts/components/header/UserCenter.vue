<!-- Author: Charlie -->

<script setup lang="ts">
import type { DropdownOption } from 'naive-ui'
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { renderIcon } from '@/utils/icon'
import { useAuthStore } from '@/stores'

const router = useRouter()
const authStore = useAuthStore()
const homePath = import.meta.env.VITE_HOME_PATH

const displayName = computed(() => authStore.userInfo?.nickname || '-')

const avatar = computed(() => authStore.userInfo?.avatar || undefined)
const avatarImgProps = { referrerPolicy: 'no-referrer' } as any

// 桌面端用户菜单项。项目首页使用环境配置路径。
const options = computed<DropdownOption[]>(() => [
  {
    label: '个人中心',
    key: 'userCenter',
    icon: renderIcon('icon-park-outline:user'),
  },
  {
    label: '我的反馈',
    key: 'feedback',
    icon: renderIcon('icon-park-outline:message'),
  },
  {
    type: 'divider',
    key: 'divider-1',
  },
  {
    label: '项目首页',
    key: 'home',
    icon: renderIcon('icon-park-outline:home'),
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: renderIcon('icon-park-outline:logout'),
  },
])

// 根据下拉菜单 key 分发行为；保留弹窗确认可以避免后续接入真实登出时误触。
function handleSelect(key: string | number) {
  if (key === 'userCenter') {
    router.push('/usercenter')
  }
  if (key === 'feedback') {
    router.push('/feedback')
  }
  if (key === 'home') {
    router.push(homePath)
  }
  if (key === 'logout') {
    window.$dialog.info({
      title: '退出登录',
      content: '确定退出当前账号？',
      positiveText: '确认',
      negativeText: '取消',
      onPositiveClick: async () => {
        await authStore.logout()
        window.$message.success('已退出登录')
      },
    })
  }
}
</script>

<template>
  <n-dropdown
    trigger="click"
    :options="options"
    @select="handleSelect"
  >
    <div class="flex items-center gap-2">
      <n-avatar
        v-if="avatar"
        round
        class="cursor-pointer"
        :src="avatar"
        :img-props="avatarImgProps"
      />
      <n-avatar
        v-else
        round
        class="cursor-pointer"
      >
        <NovaIcon icon="icon-park-outline:user" />
      </n-avatar>
      <span class="hidden md:inline">{{ displayName }}</span>
    </div>
  </n-dropdown>
</template>
