<!-- Author: Charlie -->

<script setup lang="ts">
import { usePreferredDark } from '@vueuse/core'
import { computed } from 'vue'
import DarkModeSwitch from '@/components/common/DarkModeSwitch.vue'
import { useAppStore } from '@/stores'
import './auth-page.css'

const props = withDefaults(
  defineProps<{
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
    class="admin-auth"
    :class="{
      'admin-auth--center': variant === 'center',
      'admin-auth--dark': isDarkTheme,
    }"
  >
    <div class="admin-auth__tools">
      <DarkModeSwitch />
    </div>

    <template v-if="variant === 'split'">
      <div class="admin-auth__stage admin-auth__enter">
        <section class="admin-auth__panel admin-auth__panel--form">
          <div class="admin-auth__mobile-brand">
            <RouterLink
              class="admin-auth__brand-link"
              to="/auth/login"
            >
              <span class="admin-auth__mark">{{ brandMark }}</span>
              <span class="admin-auth__name">{{ appTitle }}</span>
            </RouterLink>
          </div>
          <header class="admin-auth__head">
            <h1 class="admin-auth__title">
              {{ title }}
            </h1>
            <div
              v-if="$slots.headerExtra"
              class="admin-auth__head-extra"
            >
              <slot name="headerExtra" />
            </div>
          </header>
          <p
            v-if="description"
            class="admin-auth__desc"
          >
            {{ description }}
          </p>
          <div class="admin-auth__body">
            <slot />
          </div>
        </section>

        <aside class="admin-auth__panel admin-auth__panel--brand">
          <div
            class="admin-auth__scan"
            aria-hidden="true"
          />
          <div class="admin-auth__brand-inner">
            <RouterLink
              class="admin-auth__brand-link admin-auth__brand-link--on-dark"
              to="/auth/login"
            >
              <span class="admin-auth__mark">{{ brandMark }}</span>
              <span class="admin-auth__name">{{ appTitle }}</span>
            </RouterLink>
            <div class="admin-auth__brand-copy">
              <p class="admin-auth__eyebrow">
                Administration
              </p>
              <h2 class="admin-auth__headline">
                管理端控制台
              </h2>
              <p class="admin-auth__lead">
                统一管理组织、权限、消息与系统配置。
              </p>
            </div>
            <div
              v-if="copyright"
              class="admin-auth__foot"
            >
              <a
                v-if="copyrightHref"
                class="admin-auth__copyright"
                :href="copyrightHref"
                target="_blank"
                rel="noopener noreferrer"
              >{{ copyright }}</a>
              <template v-else>
                {{ copyright }}
              </template>
            </div>
            <div
              v-else
              class="admin-auth__foot admin-auth__foot--spacer"
              aria-hidden="true"
            />
          </div>
        </aside>
      </div>
    </template>

    <template v-else>
      <div class="admin-auth__center admin-auth__enter">
        <RouterLink
          class="admin-auth__brand-link"
          to="/auth/login"
        >
          <span class="admin-auth__mark">{{ brandMark }}</span>
          <span class="admin-auth__name">{{ appTitle }}</span>
        </RouterLink>
        <h1 class="admin-auth__title">
          {{ title }}
        </h1>
        <p
          v-if="description"
          class="admin-auth__desc"
        >
          {{ description }}
        </p>
        <div class="admin-auth__body">
          <slot />
        </div>
      </div>
    </template>
  </div>
</template>
