<!-- Author: Charlie -->

<template>
  <view>
    <u-navbar :title="title" :is-back="back" :placeholder="true" :border="true">
      <template #right>
        <slot name="right"></slot>
      </template>
    </u-navbar>
    <slot></slot>
  </view>
</template>

<script setup lang="ts">
defineProps({
  title: {
    type: String,
    required: true,
  },
  back: {
    type: Boolean,
    default: false,
  },
})
</script>

<style lang="scss" scoped></style>
