<!-- Author: Charlie -->

<script setup lang="ts">
import { authApi } from '@/api'
import { useAuthStore } from '@/stores'
import { encryptPasswords } from '@/utils/security'
import { onMounted, reactive } from 'vue'
import '../usercenter.css'
import BindConfirmModal from './BindConfirmModal.vue'

const authStore = useAuthStore()

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
    loading: false,
  },
})

onMounted(async () => {
  await refresh()
})

async function refresh() {
  state.loading = true
  try {
    const data = await authStore.refreshUserInfo()
    const currentProfile = data?.profile ?? {}
    state.phoneForm.phone = currentProfile.phone ?? ''
    state.phoneForm.phone_login_enabled = Boolean(currentProfile.phone_login_enabled)
  } finally {
    state.loading = false
  }
}

function savePhone() {
  state.bindConfirm.password = ''
  state.bindConfirm.show = true
}

async function confirmBind() {
  if (!state.bindConfirm.password) {
    window.$message.warning('请输入当前密码')
    return
  }
  state.bindConfirm.loading = true
  state.savingPhone = true
  try {
    const encrypted = await encryptPasswords({ password: state.bindConfirm.password })
    await authApi.updateUserCenterPhone({
      password: encrypted.values.password,
      password_key_id: encrypted.password_key_id,
      phone: state.phoneForm.phone || null,
      phone_login_enabled: state.phoneForm.phone_login_enabled,
    })
    state.bindConfirm.show = false
    state.bindConfirm.password = ''
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
    <NForm
      class="user-center-form user-center-form--narrow w-full min-w-0"
      label-placement="top"
    >
      <NFormItem label="手机号">
        <NInput v-model:value="state.phoneForm.phone" />
      </NFormItem>
      <NFormItem label="启用手机号登录">
        <NSwitch v-model:value="state.phoneForm.phone_login_enabled" />
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
    title="确认更新手机号"
    :loading="state.bindConfirm.loading"
    @confirm="confirmBind"
  />
</template>
