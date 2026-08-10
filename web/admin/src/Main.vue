<!-- Author: Charlie -->

<script setup lang="ts">
import type { GlobalThemeOverrides } from 'naive-ui'
import type { ProConfigProviderProps } from 'pro-naive-ui'
import { darkTheme, dateZhCN, zhCN } from 'naive-ui'
import { ProConfigProvider, zhCN as proZhCN } from 'pro-naive-ui'
import { computed } from 'vue'
import { hljs } from './plugins/hljs'
import { useAppStore } from './stores'
import themeConfig from './stores/app/theme.json'

const appStore = useAppStore()
const theme = themeConfig as GlobalThemeOverrides
const chinaTimeZone = 'Asia/Shanghai'
const proConfigProviderProps = computed<ProConfigProviderProps>(() => ({
  abstract: true,
  locale: proZhCN,
  dateLocale: dateZhCN,
  propOverrides: {
    ProButton: {
      focusable: false,
    },
    ProDataTable: {
      size: 'small',
      flexHeight: !appStore.isMobile,
      pagination: {
        pageSlot: appStore.isMobile ? 6 : undefined,
      },
    },
    ProModalForm: {
      preset: 'card',
      labelPlacement: 'left',
      labelWidth: '100',
    },
    ProTime: {
      fieldProps: {
        timeZone: chinaTimeZone,
      },
    },
  },
}))
</script>

<template>
  <n-config-provider
    class="wh-full"
    inline-theme-disabled
    :theme="appStore.naiveTheme === 'dark' ? darkTheme : null"
    :theme-overrides="theme"
    :locale="zhCN"
    :date-locale="dateZhCN"
    :hljs="hljs"
  >
    <naive-provider>
      <ProConfigProvider v-bind="proConfigProviderProps">
        <router-view />
        <watermark
          show
          text="HEI"
        />
      </ProConfigProvider>
    </naive-provider>
  </n-config-provider>
</template>
