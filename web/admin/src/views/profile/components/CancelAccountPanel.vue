<!--
  Author: Charlie

  账号注销：调用 /cancel，清理本地会话并跳转登录。
-->
<script setup lang="ts">
import { authApi } from '@/api'
import { useAuthStore } from '@/stores'
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import '../profile.css'

const authStore = useAuthStore()
const router = useRouter()

const state = reactive({
  confirmText: '',
  cancelReason: '',
  submitting: false,
})

async function submitCancel() {
  if (state.confirmText.trim() !== '注销') {
    window.$message.warning('请输入「注销」以确认')
    return
  }
  window.$dialog.warning({
    title: '确认注销账号',
    content:
      '注销后将立即失效登录、释放手机号/邮箱等登录标识，并清理资料与授权。账号进入注销保留期；期间未再登录使用的，到期后将彻底删除并按绑定邮箱/短信通知。此操作不可撤销。',
    positiveText: '确认注销',
    negativeText: '再想想',
    draggable: true,
    maskClosable: false,
    onPositiveClick: () => doCancel(),
  })
}

async function doCancel() {
  state.submitting = true
  try {
    await authApi.cancelAccount({
      cancel_reason: state.cancelReason.trim() || null,
    })
    window.$message.success('账号已注销')
    authStore.resetSession()
    await router.push({ path: '/auth/login' })
  } finally {
    state.submitting = false
  }
}
</script>

<template>
  <div class="profile-form profile-form--narrow w-full min-w-0">
    <NAlert
      type="warning"
      :show-icon="true"
      class="mb-16px"
    >
      注销为不可逆操作：将清理个人资料、登录标识、角色/部门/用户组与资源授权，并强制下线全部会话。系统按保留天数（默认
      15
      天）暂存账号主记录；到期且未再登录使用后彻底删除，并通过已绑定邮箱通知（短信需在系统配置中填写模板编号）。
    </NAlert>

    <NForm label-placement="top">
      <NFormItem label="注销原因（可选）">
        <NInput
          v-model:value="state.cancelReason"
          type="textarea"
          :rows="3"
          maxlength="500"
          show-count
          placeholder="例如：不再使用本系统"
        />
      </NFormItem>
      <NFormItem label="请输入「注销」确认">
        <NInput
          v-model:value="state.confirmText"
          placeholder="注销"
        />
      </NFormItem>
      <NFormItem>
        <NButton
          type="error"
          :loading="state.submitting"
          :disabled="state.confirmText.trim() !== '注销'"
          @click="submitCancel"
        >
          注销账号
        </NButton>
      </NFormItem>
    </NForm>
  </div>
</template>

<style scoped>
.mb-16px {
  margin-bottom: 16px;
}
</style>
