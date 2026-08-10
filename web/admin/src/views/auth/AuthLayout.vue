<!-- Author: Charlie -->

<script setup lang="ts">
import { usePreferredDark } from '@vueuse/core'
import { computed } from 'vue'
import DarkModeSwitch from '@/components/common/DarkModeSwitch.vue'
import { useAppStore } from '@/stores'
import './auth-page.css'

const props = withDefaults(
  defineProps<{
    /** split：登录分栏；center：找回/重置居中 */
    variant?: 'split' | 'center'
    title: string
    description?: string
    copyright?: string
    copyrightUrl?: string
  }>(),
  {
    variant: 'split',
    description: '',
    copyright: '',
    copyrightUrl: '',
  },
)

const appStore = useAppStore()
const prefersDark = usePreferredDark()
const appTitle = import.meta.env.VITE_APP_TITLE || 'Admin'
const copyright = computed(() => props.copyright || import.meta.env.VITE_COPYRIGHT_INFO || '')
const copyrightHref = computed(() => {
  const url = (props.copyrightUrl || '').trim()
  if (!url) return ''
  if (/^https?:\/\//i.test(url)) return url
  return `https://${url}`
})

const isDarkTheme = computed(
  () =>
    appStore.storeColorMode === 'dark' || (appStore.storeColorMode === 'auto' && prefersDark.value),
)

const brandMark = computed(() => String(appTitle).slice(0, 1).toUpperCase())
</script>

<template>
  <div
    class="auth-page"
    :class="{
      'auth-page--center': variant === 'center',
      'auth-page--dark': isDarkTheme,
    }"
  >
    <div class="auth-page__tools">
      <DarkModeSwitch />
    </div>

    <template v-if="variant === 'split'">
      <div class="auth-card">
        <aside class="auth-card__brand">
          <div
            class="auth-card__brand-deco"
            aria-hidden="true"
          />
          <div class="auth-card__brand-inner">
            <RouterLink
              class="auth-card__logo"
              to="/auth/login"
            >
              <span class="auth-card__logo-mark">{{ brandMark }}</span>
              <span class="auth-card__logo-text">{{ appTitle }}</span>
            </RouterLink>
            <p class="auth-card__eyebrow">
              Administration Console
            </p>
            <h2 class="auth-card__headline">
              管理端控制台
            </h2>
            <p class="auth-card__lead">
              统一管理组织、权限、消息与系统配置。
            </p>
            <div
              v-if="copyright"
              class="auth-card__brand-foot"
            >
              <a
                v-if="copyrightHref"
                class="auth-card__copyright-link"
                :href="copyrightHref"
                target="_blank"
                rel="noopener noreferrer"
              >{{ copyright }}</a>
              <template v-else>
                {{ copyright }}
              </template>
            </div>
          </div>
        </aside>

        <div class="auth-card__form">
          <div class="auth-card__form-head">
            <h1 class="auth-card__title">
              {{ title }}
            </h1>
            <div
              v-if="$slots.headerExtra"
              class="auth-card__form-extra"
            >
              <slot name="headerExtra" />
            </div>
          </div>
          <slot />
        </div>
      </div>
    </template>

    <template v-else>
      <div class="auth-center">
        <RouterLink
          class="auth-center__logo"
          to="/auth/login"
        >
          <span class="auth-center__logo-mark">{{ brandMark }}</span>
        </RouterLink>
        <h1 class="auth-center__title">
          {{ title }}
        </h1>
        <p
          v-if="description"
          class="auth-center__desc"
        >
          {{ description }}
        </p>
        <div class="auth-center__body">
          <slot />
        </div>
      </div>
    </template>
  </div>
</template>
