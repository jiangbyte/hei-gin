<!-- Author: Charlie -->

<template>
  <view class="card-head">
    <text class="card-head__title">{{ title }}</text>
    <text v-if="subTitle" class="card-head__subtitle">{{ subTitle }}</text>
  </view>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  subTitle?: string
}>()
</script>

<style lang="scss" scoped>
.card-head {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.card-head__title {
  font-weight: 600;
  line-height: 1.25;
}
</style>
