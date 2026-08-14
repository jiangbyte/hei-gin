<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormInst } from 'naive-ui'
import ImageUpload from '@/components/upload/ImageUpload.vue'
import { accountApi } from '@/api'
import { toNullableString } from '@/utils'
import { computed, reactive, ref } from 'vue'
import {
  accountStatusOptions,
  buildBaseAccountPayload,
  createBaseFormRules,
} from '../../composables/useAccountForm'
import { wireBool } from '@/utils/wire'

/** 对应后端 AccountType.ADMIN + AdminProfileUpsertPayload */
const ACCOUNT_TYPE = 'ADMIN'

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)

/** SysAccount + identity + admin_user_profile 可写字段 */
const defaultFormData = {
  // SysAccount
  account_type: ACCOUNT_TYPE,
  account_status: 'ENABLED',
  password: '',
  // SysAccountIdentity（ACCOUNT / EMAIL / PHONE）
  account: '',
  email: '',
  phone: '',
  email_login_enabled: false,
  phone_login_enabled: false,
  // AdminUserProfile / AdminProfileUpsertPayload
  name: '',
  nickname: '',
  avatar: '',
  signature: '',
  remark: '',
}

const state = reactive({
  showModal: false,
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: { ...defaultFormData },
})

const statusOptions = computed(() => accountStatusOptions())
const modalTitle = computed(() => (state.dataId ? '编辑管理员' : '新增管理员'))
const rules = computed(() => createBaseFormRules(() => Boolean(state.dataId), state.formModel))

async function openModal(id?: string) {
  state.dataId = id ?? null
  state.formModel = { ...defaultFormData }
  state.showModal = true
  if (id) await fetchDetail(id)
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await accountApi.detail({ id })
    const data = response.data ?? {}
    state.formModel = {
      ...defaultFormData,
      account_type: ACCOUNT_TYPE,
      account_status: data.account_status ?? 'ENABLED',
      password: '',
      account: data.account ?? '',
      email: data.email ?? '',
      phone: data.phone ?? '',
      email_login_enabled: wireBool(data.email_login_enabled ?? false),
      phone_login_enabled: wireBool(data.phone_login_enabled ?? false),
      name: data.name ?? '',
      nickname: data.nickname ?? '',
      avatar: data.avatar ?? '',
      signature: data.signature ?? '',
      remark: data.remark ?? '',
    }
  } finally {
    state.loading = false
  }
}

function closeModal() {
  state.showModal = false
  state.submitLoading = false
}

async function submitForm() {
  await formRef.value?.validate()
  state.submitLoading = true
  try {
    // AdminProfileUpsertPayload: name/nickname/avatar/signature/phone/email/remark
    const payload = {
      ...(await buildBaseAccountPayload(state.formModel)),
      remark: toNullableString(state.formModel.remark),
    }

    if (state.dataId) {
      await accountApi.update({ ...payload, id: state.dataId })
      window.$message.success('更新成功')
    } else {
      await accountApi.create(payload)
      window.$message.success('创建成功')
    }

    closeModal()
    emit('saved')
  } finally {
    state.submitLoading = false
  }
}

defineExpose({ openModal })
</script>

<template>
  <NModal
    v-model:show="state.showModal"
    preset="card"
    draggable
    :mask-closable="false"
    :title="modalTitle"
    style="width: 720px"
    :segmented="{ content: true, action: true }"
  >
    <NSpin :show="state.loading">
      <NScrollbar class="h-[480px] pr-16px">
        <NForm
          ref="formRef"
          :model="state.formModel"
          :rules="rules"
          label-placement="left"
          label-width="110"
          :disabled="state.loading || state.submitLoading"
        >
          <NTabs
            type="line"
            animated
          >
            <NTabPane
              name="account"
              tab="账号"
            >
              <NFormItem
                label="密码"
                path="password"
              >
                <NInput
                  v-model:value="state.formModel.password"
                  type="password"
                  show-password-on="click"
                  :placeholder="state.dataId ? '留空则保持当前密码' : undefined"
                />
              </NFormItem>
              <PasswordStrengthBar :password="state.formModel.password" />
              <NFormItem
                label="账号状态"
                path="account_status"
              >
                <NRadioGroup v-model:value="state.formModel.account_status">
                  <NRadio
                    v-for="option in statusOptions"
                    :key="option.value"
                    :value="option.value"
                  >
                    {{ option.label }}
                  </NRadio>
                </NRadioGroup>
              </NFormItem>
            </NTabPane>

            <NTabPane
              name="identity"
              tab="登录身份"
            >
              <NFormItem
                label="账号"
                path="account"
              >
                <NInput v-model:value="state.formModel.account" />
              </NFormItem>
              <NFormItem
                label="邮箱"
                path="email"
              >
                <NInput v-model:value="state.formModel.email" />
              </NFormItem>
              <NFormItem label="启用邮箱登录">
                <NSwitch v-model:value="state.formModel.email_login_enabled" />
              </NFormItem>
              <NFormItem
                label="手机号"
                path="phone"
              >
                <NInput v-model:value="state.formModel.phone" />
              </NFormItem>
              <NFormItem label="启用手机号登录">
                <NSwitch v-model:value="state.formModel.phone_login_enabled" />
              </NFormItem>
            </NTabPane>

            <NTabPane
              name="profile"
              tab="资料"
            >
              <NFormItem
                label="姓名"
                path="name"
              >
                <NInput v-model:value="state.formModel.name" />
              </NFormItem>
              <NFormItem
                label="昵称"
                path="nickname"
              >
                <NInput v-model:value="state.formModel.nickname" />
              </NFormItem>
              <NFormItem
                label="头像"
                path="avatar"
              >
                <ImageUpload v-model:value="state.formModel.avatar" />
              </NFormItem>
              <NFormItem
                label="个性签名"
                path="signature"
              >
                <NInput
                  v-model:value="state.formModel.signature"
                  type="textarea"
                  :autosize="{ minRows: 3, maxRows: 5 }"
                />
              </NFormItem>
              <NFormItem
                label="备注"
                path="remark"
              >
                <NInput
                  v-model:value="state.formModel.remark"
                  type="textarea"
                  :autosize="{ minRows: 3, maxRows: 5 }"
                />
              </NFormItem>
            </NTabPane>
          </NTabs>
        </NForm>
      </NScrollbar>
    </NSpin>

    <template #action>
      <NSpace
        justify="end"
        align="center"
      >
        <NButton @click="closeModal">
          取消
        </NButton>
        <NButton
          type="primary"
          :loading="state.submitLoading"
          @click="submitForm"
        >
          确认
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
