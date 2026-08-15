<!-- Author: Charlie -->

<template>
  <view class="mine-container">
    <!-- 顶部个人信息栏 -->
    <view class="header-section" @click="goToProfile">
      <view class="header-content">
        <u-avatar
          :src="authStore.userInfo?.avatar || ''"
          size="120"
          icon="account-fill"
        />
        <view class="user-info">
          <text class="user-name">{{ displayName }}</text>
          <text class="user-account">{{ authStore.userInfo?.account }}</text>
        </view>
        <view class="header-right">
          <text class="header-right-text">个人资料</text>
          <u-icon name="arrow-right" color="rgba(255,255,255,0.6)" size="14" />
        </view>
      </view>
    </view>

    <view class="content-section">
      <!-- 统计指标 -->
      <view class="action-card">
        <u-grid :col="4" :border="false">
          <u-grid-item>
            <view class="stat-item">
              <text class="stat-number">1,286</text>
              <text class="stat-label">账号总数</text>
            </view>
          </u-grid-item>
          <u-grid-item>
            <view class="stat-item">
              <text class="stat-number">86</text>
              <text class="stat-label">今日新增</text>
            </view>
          </u-grid-item>
          <u-grid-item>
            <view class="stat-item">
              <text class="stat-number">43</text>
              <text class="stat-label">在线设备</text>
            </view>
          </u-grid-item>
          <u-grid-item>
            <view class="stat-item">
              <text class="stat-number">12</text>
              <text class="stat-label">待办任务</text>
            </view>
          </u-grid-item>
        </u-grid>
      </view>
      <view class="menu-card">
        <u-cell-group :border="false">
          <u-cell-item title="退出登录" :arrow="false" @click="confirmLogout" />
        </u-cell-group>
      </view>
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

function goToProfile() {
  uni.navigateTo({ url: '/pages/profile/profile' })
}

function confirmLogout() {
  uni.showModal({
    title: '提示',
    content: '确定要退出登录吗？',
    success: (res) => {
      if (res.confirm) {
        authStore.logout()
      }
    },
  })
}
</script>

<style lang="scss" scoped>
page {
  background-color: #f5f6f7;
}

.mine-container {
  width: 100%;
  min-height: 100vh;
}

.header-section {
  padding: 30px 15px 55px 15px;
  background-color: #3c96f3;
  color: white;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-arrow {
  margin-left: auto;
}

.header-right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 4px;
}

.header-right-text {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
}

.user-info {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.user-account {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.85);
}

.content-section {
  position: relative;
  top: -40px;
  padding: 0 15px;
}

.action-card {
  padding: 8px 0;
  background-color: #fff;
  margin-bottom: 15px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-number {
  font-size: 16px;
  font-weight: 700;
  color: #333;
}

.stat-label {
  font-size: 11px;
  color: #999;
  margin-top: 4px;
}

.menu-card {
  background-color: #fff;
  overflow: hidden;
  margin-bottom: 15px;
}
</style>
