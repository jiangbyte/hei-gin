<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormInst, FormItemRule, FormRules } from 'naive-ui'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores'
import { isValidEmail } from '@/utils'
import { wireBool } from '@/utils/wire'
import { encryptPasswords } from '@/utils/security'
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import '../usercenter.css'
import BindConfirmModal from './BindConfirmModal.vue'

const OTP_COOLDOWN_SECONDS = 60
const authStore = useAuthStore()
const emailFormRef = ref<FormInst | null>(null)
const forceBindEmail = computed(() => Boolean(authStore.userInfo?.forceBindEmail))

const state = reactive({
  loading: false,
  savingEmail: false,
  emailForm: {
    email: '',
    email_login_enabled: false,
  },
  bindConfirm: {
    show: false,
    password: '',
    otpCode: '',
    loading: false,
    sendingCode: false,
    otpCooldown: 0,
  },
})

let cooldownTimer: number | undefined

const emailRules = computed<FormRules>(() => ({
  email: [
    {
      validator: validateEmailForm,
      trigger: ['input', 'blur'],
    },
  ],
}))

onMounted(async () => {
  await refresh()
})

onUnmounted(() => {
  if (cooldownTimer) window.clearInterval(cooldownTimer)
})

async function refresh() {
  state.loading = true
  try {
    const data = await authStore.refreshUserInfo()
    const currentProfile = data?.profile ?? {}
    state.emailForm.email = currentProfile.email ?? ''
    state.emailForm.email_login_enabled = forceBindEmail.value
      ? true
      : wireBool(currentProfile.email_login_enabled ?? false)
  } finally {
    state.loading = false
  }
}

async function saveEmail() {
  try {
    await emailFormRef.value?.validate()
  } catch {
    return
  }
  state.bindConfirm.password = ''
  state.bindConfirm.otpCode = ''
  state.bindConfirm.show = true
}

function validateEmailForm(_rule: FormItemRule, value: string) {
  const text = String(value ?? '').trim()
  if (!text) {
    return state.emailForm.email_login_enabled || forceBindEmail.value
      ? new Error('请输入邮箱')
      : true
  }
  if (!isValidEmail(text)) {
    return new Error('请输入有效邮箱')
  }
  return true
}

async function sendBindCode() {
  if (state.bindConfirm.otpCooldown > 0 || state.bindConfirm.sendingCode) return
  const email = state.emailForm.email.trim()
  if (!email || !isValidEmail(email)) {
    window.$message.warning('请先填写有效邮箱')
    return
  }
  state.bindConfirm.sendingCode = true
  try {
    await authApi.sendBindEmailCode({ target: email })
    window.$message.success('验证码已发送')
    state.bindConfirm.otpCooldown = OTP_COOLDOWN_SECONDS
    if (cooldownTimer) window.clearInterval(cooldownTimer)
    cooldownTimer = window.setInterval(() => {
      if (state.bindConfirm.otpCooldown <= 1) {
        state.bindConfirm.otpCooldown = 0
        if (cooldownTimer) window.clearInterval(cooldownTimer)
        return
      }
      state.bindConfirm.otpCooldown -= 1
    }, 1000)
  } finally {
    state.bindConfirm.sendingCode = false
  }
}

async function confirmBind() {
  if (!state.bindConfirm.password) {
    window.$message.warning('请输入当前密码')
    return
  }
  if (state.emailForm.email.trim() && !state.bindConfirm.otpCode.trim()) {
    window.$message.warning('请输入邮箱验证码')
    return
  }
  state.bindConfirm.loading = true
  state.savingEmail = true
  try {
    const encrypted = await encryptPasswords({ password: state.bindConfirm.password })
    await authApi.updateUserCenterEmail({
      password: encrypted.values.password,
      password_key_id: encrypted.password_key_id,
      email: state.emailForm.email.trim() || null,
      email_login_enabled: forceBindEmail.value ? true : state.emailForm.email_login_enabled,
      otp_code: state.bindConfirm.otpCode.trim() || undefined,
    })
    state.bindConfirm.show = false
    state.bindConfirm.password = ''
    state.bindConfirm.otpCode = ''
    await refresh()
    window.$message.success('绑定已更新')
  } finally {
    state.bindConfirm.loading = false
    state.savingEmail = false
  }
}

defineExpose({ refresh })
</script>

<template>
  <NSpin :show="state.loading">
    <NAlert
      v-if="forceBindEmail"
      type="warning"
      :bordered="false"
      class="mb-12px"
      title="请先绑定邮箱后才能继续使用系统"
    />
    <NForm
      ref="emailFormRef"
      class="user-center-form user-center-form--narrow w-full min-w-0"
      :model="state.emailForm"
      :rules="emailRules"
      label-placement="top"
    >
      <NFormItem
        label="邮箱"
        path="email"
      >
        <NInput v-model:value="state.emailForm.email" />
      </NFormItem>
      <NFormItem label="启用邮箱登录">
        <NSwitch
          v-model:value="state.emailForm.email_login_enabled"
          :disabled="forceBindEmail"
        />
      </NFormItem>
      <NFormItem :show-label="false">
        <NButton
          type="primary"
          :loading="state.savingEmail"
          @click="saveEmail"
        >
          更新邮箱
        </NButton>
      </NFormItem>
    </NForm>
  </NSpin>

  <BindConfirmModal
    v-model:show="state.bindConfirm.show"
    v-model:password="state.bindConfirm.password"
    v-model:otp-code="state.bindConfirm.otpCode"
    title="确认更新邮箱"
    otp-label="邮箱验证码"
    :loading="state.bindConfirm.loading"
    :sending-code="state.bindConfirm.sendingCode"
    :otp-cooldown="state.bindConfirm.otpCooldown"
    @confirm="confirmBind"
    @send-code="sendBindCode"
  />
</template>

<style scoped>
.mb-12px {
  margin-bottom: 12px;
}
</style>
