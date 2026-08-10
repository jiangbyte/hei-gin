<!-- Author: Charlie -->

<script setup lang="ts">
import { computed, reactive, type Component } from 'vue'
import { ACCOUNT_TYPE_TABS, DEFAULT_ACCOUNT_TYPE, type AccountType } from '@/constants/account'
import AdminAccountPanel from './components/admin/AccountPanel.vue'
import PortalAccountPanel from './components/portal/AccountPanel.vue'

const panelMap: Record<AccountType, Component> = {
  ADMIN: AdminAccountPanel,
  PORTAL: PortalAccountPanel,
}

const state = reactive({
  accountType: DEFAULT_ACCOUNT_TYPE as AccountType,
})

const activePanel = computed(() => panelMap[state.accountType])

function handleAccountTypeChange(value: string | number) {
  state.accountType = String(value) as AccountType
}
</script>

<template>
  <NFlex
    class="h-full min-h-0"
    vertical
  >
    <NTabs
      class="account-type-tabs"
      :value="state.accountType"
      type="line"
      animated
      @update:value="handleAccountTypeChange"
    >
      <NTabPane
        v-for="item in ACCOUNT_TYPE_TABS"
        :key="item.key"
        :name="item.key"
        :tab="item.label"
      />
    </NTabs>

    <div class="account-panel min-h-0 flex-1">
      <KeepAlive>
        <component
          :is="activePanel"
          :key="state.accountType"
        />
      </KeepAlive>
    </div>
  </NFlex>
</template>

<style scoped>
.account-type-tabs :deep(.n-tabs-pane-wrapper) {
  display: none;
}

.account-panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.account-panel > * {
  flex: 1;
  min-height: 0;
}
</style>
