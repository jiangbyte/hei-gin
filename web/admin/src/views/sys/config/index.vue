<!-- Author: Charlie -->

<script setup lang="ts">
import type { MenuOption } from 'naive-ui'
import type { Component } from 'vue'
import { renderIcon } from '@/utils/icon'
import { computed, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import './config.css'
import SysDefaultsForm from './components/SysDefaultsForm.vue'
import AuthRegisterForm from './components/AuthRegisterForm.vue'
import AuthLoginForm from './components/AuthLoginForm.vue'
import AuthPasswordForm from './components/AuthPasswordForm.vue'
import AuthTokenForm from './components/AuthTokenForm.vue'
import AuthOauthForm from './components/AuthOauthForm.vue'
import MailEngineForm from './components/MailEngineForm.vue'
import MailTemplatePanel from './components/MailTemplatePanel.vue'
import SmsEngineForm from './components/SmsEngineForm.vue'
import SmsTemplatePanel from './components/SmsTemplatePanel.vue'
import StorageConfigForm from './components/StorageConfigForm.vue'
import UploadConfigForm from './components/UploadConfigForm.vue'
import PushConfigForm from './components/PushConfigForm.vue'
import AuditAlertConfigForm from './components/AuditAlertConfigForm.vue'
import OtherConfigPanel from './components/OtherConfigPanel.vue'

const route = useRoute()
const router = useRouter()

const navItems: Array<{ key: string; label: string }> = [
  { key: 'SYS', label: '版权信息' },
  { key: 'AUTH_REGISTER', label: '注册配置' },
  { key: 'AUTH_LOGIN', label: '登录配置' },
  { key: 'AUTH_PASSWORD', label: '密码配置' },
  { key: 'AUTH_TOKEN', label: '令牌配置' },
  { key: 'AUTH_OAUTH', label: '三方登录' },
  { key: 'MAIL', label: '邮件引擎' },
  { key: 'MAIL_TEMPLATE', label: '邮件模板' },
  { key: 'SMS', label: '短信引擎' },
  { key: 'SMS_TEMPLATE', label: '短信模板' },
  { key: 'STORAGE', label: '文件存储' },
  { key: 'UPLOAD', label: '上传限制' },
  { key: 'PUSH', label: '消息推送' },
  { key: 'AUDIT_ALERT', label: '审计告警' },
  { key: 'OTHER', label: '其他配置' },
]

const panelMap: Record<string, Component> = {
  SYS: SysDefaultsForm,
  AUTH_REGISTER: AuthRegisterForm,
  AUTH_LOGIN: AuthLoginForm,
  AUTH_PASSWORD: AuthPasswordForm,
  AUTH_TOKEN: AuthTokenForm,
  AUTH_OAUTH: AuthOauthForm,
  MAIL: MailEngineForm,
  MAIL_TEMPLATE: MailTemplatePanel,
  SMS: SmsEngineForm,
  SMS_TEMPLATE: SmsTemplatePanel,
  STORAGE: StorageConfigForm,
  UPLOAD: UploadConfigForm,
  PUSH: PushConfigForm,
  AUDIT_ALERT: AuditAlertConfigForm,
  OTHER: OtherConfigPanel,
}

const state = reactive({
  activeTab: resolveInitialTab(),
})

const menuOptions: MenuOption[] = [
  {
    type: 'group',
    label: '认证',
    key: 'g-auth',
    children: [
      {
        key: 'AUTH_REGISTER',
        label: '注册配置',
        icon: renderIcon('icon-park-outline:user', 16),
      },
      {
        key: 'AUTH_LOGIN',
        label: '登录配置',
        icon: renderIcon('icon-park-outline:login', 16),
      },
      {
        key: 'AUTH_PASSWORD',
        label: '密码配置',
        icon: renderIcon('icon-park-outline:lock', 16),
      },
      {
        key: 'AUTH_TOKEN',
        label: '令牌配置',
        icon: renderIcon('icon-park-outline:key', 16),
      },
      {
        key: 'AUTH_OAUTH',
        label: '三方登录',
        icon: renderIcon('icon-park-outline:people', 16),
      },
    ],
  },
  {
    type: 'group',
    label: '通知',
    key: 'g-notify',
    children: [
      { key: 'MAIL', label: '邮件引擎', icon: renderIcon('icon-park-outline:mail', 16) },
      {
        key: 'MAIL_TEMPLATE',
        label: '邮件模板',
        icon: renderIcon('icon-park-outline:email-lock', 16),
      },
      {
        key: 'SMS',
        label: '短信引擎',
        icon: renderIcon('icon-park-outline:message', 16),
      },
      {
        key: 'SMS_TEMPLATE',
        label: '短信模板',
        icon: renderIcon('icon-park-outline:message-emoji', 16),
      },
      {
        key: 'PUSH',
        label: '消息推送',
        icon: renderIcon('icon-park-outline:send', 16),
      },
    ],
  },
  {
    type: 'group',
    label: '存储',
    key: 'g-storage',
    children: [
      {
        key: 'STORAGE',
        label: '文件存储',
        icon: renderIcon('icon-park-outline:cloud-storage', 16),
      },
      {
        key: 'UPLOAD',
        label: '上传限制',
        icon: renderIcon('icon-park-outline:upload', 16),
      },
    ],
  },
  {
    type: 'group',
    label: '安全与集成',
    key: 'g-security',
    children: [
      {
        key: 'AUDIT_ALERT',
        label: '审计告警',
        icon: renderIcon('icon-park-outline:alarm', 16),
      },
      {
        key: 'SYS',
        label: '版权信息',
        icon: renderIcon('icon-park-outline:copyright', 16),
      },
      {
        key: 'OTHER',
        label: '其他配置',
        icon: renderIcon('icon-park-outline:list', 16),
      },
    ],
  },
]

const activeNav = computed(
  () => navItems.find((item) => item.key === state.activeTab) ?? navItems[0],
)
const activePanel = computed(() => panelMap[state.activeTab] ?? AuthRegisterForm)

function resolveInitialTab() {
  const tab = typeof route.query.tab === 'string' ? route.query.tab : ''
  if (tab && navItems.some((item) => item.key === tab)) {
    return tab
  }
  return 'AUTH_REGISTER'
}

function selectTab(key: string) {
  if (!key || state.activeTab === key) return
  if (!navItems.some((item) => item.key === key)) return
  state.activeTab = key
  void router.replace({ query: { ...route.query, tab: key } })
}

watch(
  () => route.query.tab,
  (tab) => {
    if (typeof tab !== 'string' || !navItems.some((item) => item.key === tab)) return
    if (state.activeTab !== tab) state.activeTab = tab
  },
)
</script>

<template>
  <div class="sys-config w-full min-w-0">
    <div class="sys-config__body">
      <aside class="sys-config__sidebar">
        <NMenu
          :value="state.activeTab"
          :options="menuOptions"
          :root-indent="12"
          :indent="18"
          @update:value="selectTab"
        />
      </aside>

      <section class="sys-config__content">
        <div class="sys-config__panel">
          <h2 class="sys-config__panel-title">
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
  </div>
</template>
