<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormItemInst, FormItemRule } from 'naive-ui'
import { authApi } from '@/api'
import CaptchaInput from '@/components/common/CaptchaInput.vue'
import { isValidEmail } from '@/utils'
import { encryptPasswords } from '@/utils/security'
import { computed, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AuthLayout from './AuthLayout.vue'

const route = useRoute()
const router = useRouter()
const emailItemRef = ref<FormItemInst | null>(null)
const passwordItemRef = ref<FormItemInst | null>(null)
const confirmItemRef = ref<FormItemInst | null>(null)
const captchaItemRef = ref<FormItemInst | null>(null)
const captchaRef = ref<InstanceType<typeof CaptchaInput> | null>(null)
const loading = ref(false)

const form = reactive({
  email: '',
  token: typeof route.query.token === 'string' ? route.query.token : '',
  password: '',
  confirmPassword: '',
  captcha_id: '',
  captcha_value: '',
})

const isResetMode = computed(() => Boolean(form.token))

const emailRule: FormItemRule = {
  trigger: ['input', 'blur'],
  validator() {
    const text = form.email.trim()
    if (!text) return new Error('请输入登录邮箱')
    if (!isValidEmail(text)) return new Error('请输入有效邮箱')
    return true
  },
}

const passwordRule: FormItemRule = {
  trigger: ['input', 'blur'],
  validator() {
    if (!form.password) return new Error('请输入新密码')
    if (form.password.length < 8) return new Error('密码至少 8 个字符')
    return true
  },
}

const confirmPasswordRule: FormItemRule = {
  trigger: ['input', 'blur'],
  validator() {
    if (!form.confirmPassword) return new Error('请确认密码')
    if (form.confirmPassword !== form.password) return new Error('两次输入的密码不一致')
    return true
  },
}

const captchaRule: FormItemRule = {
  trigger: ['input', 'blur'],
  validator() {
    if (!form.captcha_value.trim()) return new Error('请输入验证码')
    return true
  },
}

async function validateFields(items: Array<FormItemInst | null | undefined>) {
  try {
    await Promise.all(items.filter(Boolean).map((item) => item!.validate()))
    return true
  } catch {
    return false
  }
}

async function sendLink() {
  const ok = await validateFields([emailItemRef.value, captchaItemRef.value])
  if (!ok) return
  loading.value = true
  try {
    await authApi.forgotPassword({
      email: form.email.trim(),
      captcha_id: form.captcha_id,
      captcha_value: form.captcha_value,
    })
    window.$message.success('密码重置链接已发送')
    await captchaRef.value?.refresh()
  } catch {
    await captchaRef.value?.refresh()
  } finally {
    loading.value = false
  }
}

async function resetPassword() {
  const ok = await validateFields([
    passwordItemRef.value,
    confirmItemRef.value,
    captchaItemRef.value,
  ])
  if (!ok) return
  loading.value = true
  try {
    const encrypted = await encryptPasswords({ password: form.password })
    await authApi.resetPassword({
      token: form.token,
      password: encrypted.values.password,
      password_key_id: encrypted.password_key_id,
      captcha_id: form.captcha_id,
      captcha_value: form.captcha_value,
    })
    window.$message.success('密码已重置，请重新登录')
    router.push('/auth/login')
  } catch {
    await captchaRef.value?.refresh()
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthLayout
    :title="isResetMode ? '重置密码' : '找回密码'"
    :description="
      isResetMode
        ? '请设置新密码。重置链接在过期前仅可使用一次。'
        : '请输入已启用管理端登录的邮箱，系统将发送密码重置链接。'
    "
  >
    <template #headerExtra>
      <RouterLink
        class="linkish"
        to="/auth/login"
      >
        返回登录
      </RouterLink>
    </template>

    <form
      class="auth-form"
      @submit.prevent="isResetMode ? resetPassword() : sendLink()"
    >
      <n-form-item
        v-if="!isResetMode"
        ref="emailItemRef"
        :show-label="false"
        :rule="emailRule"
      >
        <n-input
          v-model:value="form.email"
          size="large"
          clearable
          placeholder="登录邮箱"
        />
      </n-form-item>

      <template v-if="isResetMode">
        <n-form-item
          ref="passwordItemRef"
          :show-label="false"
          :rule="passwordRule"
        >
          <n-input
            v-model:value="form.password"
            size="large"
            type="password"
            show-password-on="click"
            placeholder="新密码（至少 8 位）"
          />
        </n-form-item>
        <PasswordStrengthBar :password="form.password" />
        <n-form-item
          ref="confirmItemRef"
          :show-label="false"
          :rule="confirmPasswordRule"
        >
          <n-input
            v-model:value="form.confirmPassword"
            size="large"
            type="password"
            show-password-on="click"
            placeholder="确认新密码"
          />
        </n-form-item>
      </template>

      <n-form-item
        ref="captchaItemRef"
        :show-label="false"
        :rule="captchaRule"
      >
        <CaptchaInput
          ref="captchaRef"
          v-model:captcha-id="form.captcha_id"
          v-model:captcha-value="form.captcha_value"
          size="large"
        />
      </n-form-item>

      <n-button
        type="primary"
        size="large"
        block
        attr-type="submit"
        :loading="loading"
      >
        {{ isResetMode ? '重置密码' : '发送重置链接' }}
      </n-button>

      <div
        v-if="isResetMode"
        class="auth-form-row"
        style="margin-top: 12px"
      >
        <RouterLink to="/auth/forgot-password">
          重新申请链接
        </RouterLink>
      </div>
    </form>
  </AuthLayout>
</template>
