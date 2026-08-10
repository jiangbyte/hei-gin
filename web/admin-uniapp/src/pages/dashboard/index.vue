<!-- Author: Charlie -->

<template>
  <view class="content">
    <image class="logo" src="/static/logo.png"></image>
    <view class="text-area">
      <text class="title">欢迎回来</text>
      <text class="subtitle">{{ displayName }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const displayName = computed(
  () =>
    authStore.userInfo?.nickname ||
    authStore.userInfo?.name ||
    authStore.userInfo?.account ||
    '管理员'
)

onShow(() => {
  if (!authStore.isLogin) {
    uni.reLaunch({ url: '/pages/auth/login/login' })
  }
})
</script>

<style lang="scss" scoped>
page {
  background-color: #f5f6f7;
}

.content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
}

.logo {
  height: 200rpx;
  width: 200rpx;
  margin-top: -100rpx;
  margin-bottom: 40rpx;
}

.text-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.title {
  font-size: 36rpx;
  color: #333;
  font-weight: 600;
}

.subtitle {
  font-size: 28rpx;
  color: #999;
}
</style>
