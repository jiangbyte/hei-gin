<!-- Author: Charlie -->

<script setup lang="ts">
import { useDialog, useLoadingBar, useMessage, useNotification } from 'naive-ui'
import { defineComponent, h } from 'vue'

const NaiveProviderContent = defineComponent({
  name: 'NaiveProviderContent',
  setup() {
    window.$loadingBar = useLoadingBar()
    window.$dialog = useDialog()
    window.$message = useMessage()
    window.$notification = useNotification()
  },
  render() {
    return h('div')
  },
})
</script>

<template>
  <n-loading-bar-provider>
    <n-dialog-provider>
      <n-notification-provider>
        <n-message-provider>
          <slot />
          <NaiveProviderContent />
        </n-message-provider>
      </n-notification-provider>
    </n-dialog-provider>
  </n-loading-bar-provider>
</template>

<style scoped></style>
