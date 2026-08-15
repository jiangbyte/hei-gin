<!-- Author: Charlie -->

<script setup lang="ts">
import { authApi } from '@/api'
import { useAuthStore } from '@/stores'
import { isValidPhone } from '@/utils'
import { wireBool } from '@/utils/wire'
import { encryptPasswords } from '@/utils/security'
import { computed, onMounted, onUnmounted, reactive } from 'vue'
import '../profile.css'
import BindConfirmModal from './BindConfirmModal.vue'

const OTP_COOLDOWN_SECONDS = 60
const authStore = useAuthStore()
const forceBindPhone = computed(() => Boolean(authStore.userInfo?.forceBindPhone))

const state = reactive({
  loading: false,
  savingPhone: false,
  phoneForm: {
    phone: '',
    phone_login_enabled: false,
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
    state.phoneForm.phone = currentProfile.phone ?? ''
    state.phoneForm.phone_login_enabled = forceBindPhone.value
      ? true
      : wireBool(currentProfile.phone_login_enabled ?? false)
  } finally {
    state.loading = false
  }
}

function savePhone() {
  const phone = state.phoneForm.phone.trim()
  if (phone && !isValidPhone(phone)) {
    window.$message.warning('请输入有效手机号')
    return
  }
  if (!phone && (state.phoneForm.phone_login_enabled || forceBindPhone.value)) {
    window.$message.warning('请输入手机号')
    return
  }
  state.bindConfirm.password = ''
  state.bindConfirm.otpCode = ''
  state.bindConfirm.show = true
}

async function sendBindCode() {
  if (state.bindConfirm.otpCooldown > 0 || state.bindConfirm.sendingCode) return
  const phone = state.phoneForm.phone.trim()
  if (!phone || !isValidPhone(phone)) {
    window.$message.warning('请先填写有效手机号')
    return
  }
  state.bindConfirm.sendingCode = true
  try {
    await authApi.sendBindPhoneCode({ target: phone })
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
  if (state.phoneForm.phone.trim() && !state.bindConfirm.otpCode.trim()) {
    window.$message.warning('请输入手机验证码')
    return
  }
  state.bindConfirm.loading = true
  state.savingPhone = true
  try {
    const encrypted = await encryptPasswords({ password: state.bindConfirm.password })
    await authApi.updatePhone({
      password: encrypted.values.password,
      password_key_id: encrypted.password_key_id,
      phone: state.phoneForm.phone.trim() || null,
      phone_login_enabled: forceBindPhone.value ? true : state.phoneForm.phone_login_enabled,
      otp_code: state.bindConfirm.otpCode.trim() || undefined,
    })
    state.bindConfirm.show = false
    state.bindConfirm.password = ''
    state.bindConfirm.otpCode = ''
    await refresh()
    window.$message.success('绑定已更新')
  } finally {
    state.bindConfirm.loading = false
    state.savingPhone = false
  }
}

defineExpose({ refresh })
</script>

<template>
  <NSpin :show="state.loading">
    <NAlert
      v-if="forceBindPhone"
      type="warning"
      :bordered="false"
      class="mb-12px"
      title="请先绑定手机号后才能继续使用系统"
    />
    <NForm
      class="profile-form profile-form--narrow w-full min-w-0"
      label-placement="top"
    >
      <NFormItem label="手机号">
        <NInput v-model:value="state.phoneForm.phone" />
      </NFormItem>
      <NFormItem label="启用手机号登录">
        <NSwitch
          v-model:value="state.phoneForm.phone_login_enabled"
          :disabled="forceBindPhone"
        />
      </NFormItem>
      <NFormItem :show-label="false">
        <NButton
          type="primary"
          :loading="state.savingPhone"
          @click="savePhone"
        >
          更新手机号
        </NButton>
      </NFormItem>
    </NForm>
  </NSpin>

  <BindConfirmModal
    v-model:show="state.bindConfirm.show"
    v-model:password="state.bindConfirm.password"
    v-model:otp-code="state.bindConfirm.otpCode"
    title="确认更新手机号"
    otp-label="手机验证码"
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
