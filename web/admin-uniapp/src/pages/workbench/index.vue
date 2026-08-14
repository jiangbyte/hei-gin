<!-- Author: Charlie -->

<template>
  <Layout title="工作台">
    <view class="work-container">
      <!-- 快捷入口 -->
      <view class="mx-4 mt-5">
        <text class="section-title">系统管理</text>
      </view>
      <view class="mx-4 mt-2 bg-white">
        <u-grid :col="4" :border="false">
          <u-grid-item
            v-for="(item, index) in gridItems"
            :key="index"
            @click="onGridClick(index)"
          >
            <view class="grid-item-box">
              <u-icon :name="item.icon" size="30" />
              <text class="grid-text">{{ item.name }}</text>
            </view>
          </u-grid-item>
        </u-grid>
      </view>
    </view>
  </Layout>
</template>

<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app'
import Layout from '@/layouts/index.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const gridItems = [
  {name: '字典测试', icon: 'grid'},
  {name: '请求测试', icon: 'reload'},
]

onShow(() => {
  if (!authStore.isLogin) {
    uni.reLaunch({ url: '/pages/auth/login/login' })
  }
})

function onGridClick(index: number) {
  const item = gridItems[index].name
  if (item === '字典测试') {
    uni.navigateTo({url: '/pages/workbench/dict-test'})
  } else if (item === '请求测试') {
    uni.navigateTo({url: '/pages/workbench/request-test'})
  }
}
</script>

<style lang="scss" scoped>
.work-container {
  min-height: 100vh;
  padding-bottom: 20px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.grid-item-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 15px 0;
}

.grid-text {
  text-align: center;
  font-size: 26rpx;
  margin-top: 10rpx;
  color: #555;
}
</style>
