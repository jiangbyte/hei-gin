<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormInst, FormItemRule, FormRules } from 'naive-ui'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores'
import { isValidEmail } from '@/utils'
import { encryptPasswords } from '@/utils/security'
import { computed, onMounted, reactive, ref } from 'vue'
import '../usercenter.css'
import BindConfirmModal from './BindConfirmModal.vue'

const authStore = useAuthStore()
const emailFormRef = ref<FormInst | null>(null)

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
    loading: false,
  },
})

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

async function refresh() {
  state.loading = true
  try {
    const data = await authStore.refreshUserInfo()
    const currentProfile = data?.profile ?? {}
    state.emailForm.email = currentProfile.email ?? ''
    state.emailForm.email_login_enabled = Boolean(currentProfile.email_login_enabled)
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
  state.bindConfirm.show = true
}

function validateEmailForm(_rule: FormItemRule, value: string) {
  const text = String(value ?? '').trim()
  if (!text) {
    return state.emailForm.email_login_enabled ? new Error('请输入邮箱') : true
  }
  if (!isValidEmail(text)) {
    return new Error('请输入有效邮箱')
  }
  return true
}

async function confirmBind() {
  if (!state.bindConfirm.password) {
    window.$message.warning('请输入当前密码')
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
      email_login_enabled: state.emailForm.email_login_enabled,
    })
    state.bindConfirm.show = false
    state.bindConfirm.password = ''
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
        <NSwitch v-model:value="state.emailForm.email_login_enabled" />
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
    title="确认更新邮箱"
    :loading="state.bindConfirm.loading"
    @confirm="confirmBind"
  />
</template>
