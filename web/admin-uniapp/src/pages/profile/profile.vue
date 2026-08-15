<!-- Author: Charlie -->

<template>
  <Layout title="个人资料" :back="true">
    <view class="profile-container">
      <view class="avatar-section" @click="chooseAvatar">
        <u-avatar :src="form.avatar || ''" size="200" icon="account-fill" />
      </view>
      <view class="form-card">
        <u-form :model="form">
          <u-form-item label="账号">
            <text class="field-disabled">{{ userInfo?.account }}</text>
          </u-form-item>
          <u-form-item label="昵称">
            <u-input v-model="form.nickname" placeholder="请输入昵称" border />
          </u-form-item>
          <u-form-item label="姓名">
            <u-input v-model="form.name" placeholder="请输入姓名" border />
          </u-form-item>
        </u-form>
      </view>

      <view class="submit-btn">
        <u-button
          text="保存"
          type="primary"
          :loading="saving"
          @click="submit"
        />
      </view>
    </view>
  </Layout>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import Layout from '@/layouts/index.vue'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api'

const authStore = useAuthStore()
const userInfo = authStore.userInfo
const saving = ref(false)

const form = reactive({
  avatar: userInfo?.avatar || '',
  nickname: userInfo?.nickname || '',
  name: userInfo?.name || '',
})

onShow(() => {
  if (!authStore.isLogin) {
    uni.reLaunch({ url: '/pages/auth/login/login' })
    return
  }
  // 从 store 中刷新 form 数据（裁剪头像后返回需要更新）
  const ua = authStore.userInfo
  if (ua) {
    form.avatar = ua.avatar || form.avatar
    form.nickname = ua.nickname || form.nickname
    form.name = ua.name || form.name
  }
})

function chooseAvatar() {
  uni.chooseImage({
    count: 1,
    success: (res) => {
      const path = encodeURIComponent(res.tempFilePaths[0])
      uni.navigateTo({ url: `/pages/profile/avatar-crop?imagePath=${path}` })
    },
  })
}

async function submit() {
  saving.value = true
  try {
    await authApi.updateProfile({
      nickname: form.nickname,
      name: form.name,
    })
    await authStore.refreshUserInfo()
    uni.showToast({ title: '保存成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 1500)
  } catch {
    uni.showToast({ title: '保存失败', icon: 'none' })
  } finally {
    saving.value = false
  }
}
</script>

<style lang="scss" scoped>
.profile-container {
  padding: 15px;
}

.form-card {
  background-color: #fff;
  padding: 0 15px;
}

.avatar-section {
  text-align: center;
  padding: 16px 0;
}

.field-disabled {
  color: #999;
  font-size: 14px;
}

.submit-btn {
  margin-top: 30px;
}
</style>
