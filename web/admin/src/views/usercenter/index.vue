<!-- Author: Charlie -->

<script setup lang="ts">
import type { MenuOption } from 'naive-ui'
import type { Component } from 'vue'
import { renderIcon } from '@/utils/icon'
import { computed, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import './usercenter.css'
import BasicInfoPanel from './components/BasicInfoPanel.vue'
import CancelAccountPanel from './components/CancelAccountPanel.vue'
import EmailPanel from './components/EmailPanel.vue'
import MyMessagesPanel from './components/MyMessagesPanel.vue'
import PasswordPanel from './components/PasswordPanel.vue'
import PhonePanel from './components/PhonePanel.vue'

const route = useRoute()
const router = useRouter()

const navItems: Array<{ key: string; label: string; icon: string }> = [
  { key: 'basic_info', label: '公开资料', icon: 'icon-park-outline:people' },
  { key: 'my_messages', label: '我的消息', icon: 'icon-park-outline:message' },
  { key: 'password', label: '密码', icon: 'icon-park-outline:lock' },
  { key: 'phone', label: '手机号', icon: 'icon-park-outline:phone-telephone' },
  { key: 'email', label: '邮箱', icon: 'icon-park-outline:mail' },
  { key: 'cancel_account', label: '账号注销', icon: 'icon-park-outline:logout' },
]

const panelMap: Record<string, Component> = {
  basic_info: BasicInfoPanel,
  my_messages: MyMessagesPanel,
  password: PasswordPanel,
  phone: PhonePanel,
  email: EmailPanel,
  cancel_account: CancelAccountPanel,
}

const state = reactive({
  loading: false,
  activeTab: resolveInitialTab(),
})

const menuOptions: MenuOption[] = [
  {
    key: 'basic_info',
    label: '公开资料',
    icon: renderIcon('icon-park-outline:people', 16),
  },
  {
    type: 'group',
    label: '消息',
    key: 'messages',
    children: [
      {
        key: 'my_messages',
        label: '我的消息',
        icon: renderIcon('icon-park-outline:message', 16),
      },
    ],
  },
  {
    type: 'group',
    label: '访问与安全',
    key: 'access',
    children: [
      {
        key: 'password',
        label: '密码',
        icon: renderIcon('icon-park-outline:lock', 16),
      },
      {
        key: 'phone',
        label: '手机号',
        icon: renderIcon('icon-park-outline:phone-telephone', 16),
      },
      {
        key: 'email',
        label: '邮箱',
        icon: renderIcon('icon-park-outline:mail', 16),
      },
      {
        key: 'cancel_account',
        label: '账号注销',
        icon: renderIcon('icon-park-outline:logout', 16),
      },
    ],
  },
]

const activeNav = computed(
  () => navItems.find((item) => item.key === state.activeTab) ?? navItems[0],
)
const activePanel = computed(() => panelMap[state.activeTab] ?? BasicInfoPanel)

function resolveInitialTab() {
  const tab = typeof route.query.tab === 'string' ? route.query.tab : ''
  if (tab && navItems.some((item) => item.key === tab)) {
    return tab
  }
  return 'basic_info'
}

function selectTab(key: string) {
  if (!key || state.activeTab === key) {
    return
  }
  if (!navItems.some((item) => item.key === key)) {
    return
  }
  state.activeTab = key
  void router.replace({ query: { ...route.query, tab: key } })
}

watch(
  () => route.query.tab,
  (tab) => {
    if (
      typeof tab === 'string' &&
      navItems.some((item) => item.key === tab) &&
      state.activeTab !== tab
    ) {
      state.activeTab = tab
    }
  },
)
</script>

<template>
  <div class="user-center w-full min-w-0">
    <NSpin :show="state.loading">
      <div class="user-center__body">
        <aside class="user-center__sidebar">
          <NMenu
            :value="state.activeTab"
            :options="menuOptions"
            :root-indent="12"
            :indent="18"
            @update:value="selectTab"
          />
        </aside>

        <section class="user-center__content">
          <div class="user-center__panel">
            <h2
              class="user-center__panel-title"
              :class="{ 'user-center__panel-title--with-tabs': state.activeTab === 'basic_info' }"
            >
              {{ activeNav.label }}
            </h2>

            <KeepAlive>
              <component
                :is="activePanel"
                :key="state.activeTab"
              />
            </KeepAlive>
          </div>
        </section>
      </div>
    </NSpin>
  </div>
</template>
