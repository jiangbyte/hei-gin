<!-- Author: Charlie -->

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  show: boolean
  title: string
  password: string
  otpCode?: string
  otpLabel?: string
  loading?: boolean
  sendingCode?: boolean
  otpCooldown?: number
}>()

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void
  (e: 'update:password', value: string): void
  (e: 'update:otpCode', value: string): void
  (e: 'confirm'): void
  (e: 'sendCode'): void
}>()

const modalShow = computed({
  get: () => props.show,
  set: (value: boolean) => emit('update:show', value),
})

const passwordModel = computed({
  get: () => props.password,
  set: (value: string) => emit('update:password', value),
})

const otpModel = computed({
  get: () => props.otpCode ?? '',
  set: (value: string) => emit('update:otpCode', value),
})
</script>

<template>
  <NModal
    v-model:show="modalShow"
    preset="card"
    :title="title"
    class="max-w-120"
    :bordered="false"
    :mask-closable="false"
  >
    <NForm label-placement="top">
      <NFormItem :label="otpLabel || '验证码'">
        <NInputGroup>
          <NInput
            v-model:value="otpModel"
            placeholder="请输入验证码"
          />
          <NButton
            :loading="sendingCode"
            :disabled="(otpCooldown ?? 0) > 0"
            @click="emit('sendCode')"
          >
            {{ (otpCooldown ?? 0) > 0 ? `${otpCooldown}s` : '发送验证码' }}
          </NButton>
        </NInputGroup>
      </NFormItem>
      <NFormItem label="当前密码">
        <NInput
          v-model:value="passwordModel"
          type="password"
          show-password-on="click"
          placeholder="请输入当前密码"
          @keydown.enter="emit('confirm')"
        />
      </NFormItem>
    </NForm>
    <template #footer>
      <NSpace justify="end">
        <NButton @click="modalShow = false">
          取消
        </NButton>
        <NButton
          type="primary"
          :loading="loading"
          @click="emit('confirm')"
        >
          确认
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
