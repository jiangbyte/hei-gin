<!-- Author: Charlie -->

<template>
  <view class="login-page">
    <view class="login-header">
      <image
        class="login-header__logo"
        src="/static/logo.png"
        mode="aspectFit"
      ></image>
      <text class="login-header__title">管理端控制台</text>
      <text class="login-header__subtitle">登录后继续</text>
    </view>

    <u-card class="login-card" :show-head="false">
      <u-tabs
        :list="loginTabs"
        :current="activeIndex"
        @change="changeLoginType"
      ></u-tabs>

      <u-form :model="form">
        <u-form-item :border-bottom="false">
          <view class="form-field">
            <text class="form-field__label">{{ currentLogin.label }}</text>
            <u-input
              v-model="form[activeField]"
              :placeholder="currentLogin.placeholder"
              border
              clearable
            ></u-input>
          </view>
        </u-form-item>

        <u-form-item :border-bottom="false">
          <view class="form-field">
            <text class="form-field__label">密码</text>
            <u-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              border
            ></u-input>
          </view>
        </u-form-item>

        <u-form-item :border-bottom="false">
          <view class="form-field">
            <text class="form-field__label">验证码</text>
            <view class="captcha-row">
              <u-input
                v-model="form.captcha_value"
                placeholder="请输入验证码"
                border
              ></u-input>
              <image
                v-if="captchaImage"
                :src="captchaImage"
                mode="aspectFit"
                class="captcha-image"
                @click="loadCaptcha"
              ></image>
            </view>
          </view>
        </u-form-item>
      </u-form>

      <view class="login-submit">
        <u-button
          text="登录"
          type="primary"
          :loading="loading"
          @click="submit"
        ></u-button>
      </view>
    </u-card>
  </view>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores/auth'

type LoginType = 'ACCOUNT' | 'EMAIL' | 'PHONE'

const authStore = useAuthStore()
const loading = ref(false)
const captchaImage = ref('')
const activeType = ref<LoginType>('ACCOUNT')

const loginTypes: Array<{
  key: LoginType
  name: string
  label: string
  placeholder: string
}> = [
  {
    key: 'ACCOUNT',
    name: '账号',
    label: '账号',
    placeholder: '请输入管理员账号',
  },
  {
    key: 'EMAIL',
    name: '邮箱',
    label: '邮箱',
    placeholder: '请输入登录邮箱',
  },
  {
    key: 'PHONE',
    name: '手机号',
    label: '手机号',
    placeholder: '请输入登录手机号',
  },
]

const form = reactive({
  account: '',
  email: '',
  phone: '',
  password: '',
  captcha_id: '',
  captcha_value: '',
})

const currentLogin = computed(
  () =>
    loginTypes.find((item) => item.key === activeType.value) || loginTypes[0]
)
const activeField = computed(
  () => activeType.value.toLowerCase() as 'account' | 'email' | 'phone'
)
const loginTabs = computed(() =>
  loginTypes.map((item) => ({ name: item.name, key: item.key }))
)
const activeIndex = computed(() =>
  loginTypes.findIndex((item) => item.key === activeType.value)
)

function changeLoginType(event: any) {
  const index = typeof event === 'number' ? event : event.index
  activeType.value = loginTypes[index]?.key || 'ACCOUNT'
}

onLoad(() => {
  loadCaptcha()
})

async function loadCaptcha() {
  const captcha = await authApi.captcha({ format: 'png' })
  form.captcha_id = captcha.captcha_id
  captchaImage.value = `data:${captcha.image_type || 'image/png'};base64,${captcha.image_base64}`
}

async function submit() {
  const identity = form[activeField.value]
  if (!identity || !form.password || !form.captcha_value) {
    uni.showToast({ title: '请完整填写登录信息', icon: 'none' })
    return
  }
  loading.value = true
  try {
    await authStore.login({
      account: identity.trim(),
      password: form.password,
      captcha_id: form.captcha_id,
      captcha_value: form.captcha_value,
      identity_type: activeType.value,
    })
    uni.switchTab({ url: '/pages/dashboard/index' })
  } catch {
    form.captcha_value = ''
    await loadCaptcha()
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16) var(--space-8);
  background-color: var(--color-neutral-50);
}

.login-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
}

.login-header__logo {
  width: 64px;
  height: 64px;
}

.login-header__title {
  font-size: var(--text-xl);
  font-weight: 600;
  line-height: 1.25;
  color: var(--color-neutral-900);
}

.login-header__subtitle {
  font-size: var(--text-sm);
  color: var(--color-neutral-500);
}

.login-card {
  width: 100%;
}

.form-field {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.form-field__label {
  font-size: var(--text-sm);
  color: var(--color-neutral-500);
}

.captcha-row {
  width: 100%;
  display: flex;
  gap: var(--space-3);
  align-items: center;
}

.captcha-image {
  width: 80px;
  height: 32px;
}

.login-submit {
  margin-top: var(--space-8);
}
</style>
