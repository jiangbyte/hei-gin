<!-- Author: Charlie -->

<script setup lang="ts">
import { authApi } from '@/api'
import { encryptPasswords } from '@/utils/security'
import { computed, onMounted, reactive } from 'vue'
import '../usercenter.css'

type VerifyMethod = 'OLD_PASSWORD' | 'EMAIL_CODE' | 'PHONE_CODE'

const state = reactive({
  savingPassword: false,
  sendingCode: false,
  verifyMethod: 'OLD_PASSWORD' as VerifyMethod,
  passwordForm: {
    old_password: '',
    otp_code: '',
    new_password: '',
    confirm_password: '',
  },
})

const verifyHint = computed(() => {
  if (state.verifyMethod === 'EMAIL_CODE') return '将向已绑定邮箱发送验证码'
  if (state.verifyMethod === 'PHONE_CODE') return '将向已绑定手机号发送验证码'
  return ''
})

onMounted(async () => {
  try {
    const res = await authApi.authOptions()
    const method = String(res?.data?.password_change_verify_method || 'OLD_PASSWORD').toUpperCase()
    if (method === 'EMAIL_CODE' || method === 'PHONE_CODE' || method === 'OLD_PASSWORD') {
      state.verifyMethod = method
    }
  } catch {
    state.verifyMethod = 'OLD_PASSWORD'
  }
})

async function sendCode() {
  state.sendingCode = true
  try {
    await authApi.sendPasswordChangeCode()
    window.$message.success('验证码已发送')
  } finally {
    state.sendingCode = false
  }
}

async function savePassword() {
  if (state.passwordForm.new_password !== state.passwordForm.confirm_password) {
    window.$message.warning('两次输入的新密码不一致')
    return
  }
  if (state.verifyMethod === 'OLD_PASSWORD' && !state.passwordForm.old_password) {
    window.$message.warning('请输入旧密码')
    return
  }
  if (state.verifyMethod !== 'OLD_PASSWORD' && !state.passwordForm.otp_code.trim()) {
    window.$message.warning('请输入验证码')
    return
  }
  state.savingPassword = true
  try {
    const payload: Record<string, string> = {
      new_password: state.passwordForm.new_password,
    }
    if (state.verifyMethod === 'OLD_PASSWORD') {
      payload.old_password = state.passwordForm.old_password
    }
    const encrypted = await encryptPasswords(payload)
    await authApi.updateUserCenterPassword({
      old_password: encrypted.values.old_password,
      new_password: encrypted.values.new_password,
      password_key_id: encrypted.password_key_id,
      otp_code:
        state.verifyMethod === 'OLD_PASSWORD' ? undefined : state.passwordForm.otp_code.trim(),
    })
    state.passwordForm.old_password = ''
    state.passwordForm.otp_code = ''
    state.passwordForm.new_password = ''
    state.passwordForm.confirm_password = ''
    window.$message.success('密码已更新')
  } finally {
    state.savingPassword = false
  }
}
</script>

<template>
  <NForm
    class="user-center-form user-center-form--narrow w-full min-w-0"
    label-placement="top"
  >
    <NFormItem
      v-if="state.verifyMethod === 'OLD_PASSWORD'"
      label="旧密码"
    >
      <NInput
        v-model:value="state.passwordForm.old_password"
        type="password"
        show-password-on="click"
      />
    </NFormItem>
    <NFormItem
      v-else
      :label="state.verifyMethod === 'EMAIL_CODE' ? '邮箱验证码' : '手机验证码'"
    >
      <div class="password-otp-row">
        <NInput
          v-model:value="state.passwordForm.otp_code"
          :placeholder="verifyHint"
        />
        <NButton
          :loading="state.sendingCode"
          @click="sendCode"
        >
          发送验证码
        </NButton>
      </div>
    </NFormItem>
    <NFormItem label="新密码">
      <NInput
        v-model:value="state.passwordForm.new_password"
        type="password"
        show-password-on="click"
      />
    </NFormItem>
    <NFormItem label="确认密码">
      <NInput
        v-model:value="state.passwordForm.confirm_password"
        type="password"
        show-password-on="click"
      />
    </NFormItem>
    <NFormItem :show-label="false">
      <NButton
        type="primary"
        :loading="state.savingPassword"
        @click="savePassword"
      >
        更新密码
      </NButton>
    </NFormItem>
  </NForm>
</template>

<style scoped>
.password-otp-row {
  display: flex;
  gap: 8px;
  width: 100%;
}
</style>
