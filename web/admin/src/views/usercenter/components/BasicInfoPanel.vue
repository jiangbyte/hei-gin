<!-- Author: Charlie -->

<script setup lang="ts">
import { authApi } from '@/api'
import { useAuthStore } from '@/stores'
import { computed, onMounted, reactive } from 'vue'
import { displayValue, mapNames } from '../composables/useUserCenterProfile'
import '../usercenter.css'
import AvatarUploadModal from './AvatarUploadModal.vue'

const authStore = useAuthStore()
const avatarImgProps = { referrerPolicy: 'no-referrer' } as any

const state = reactive({
  loading: false,
  savingProfile: false,
  basicInfoTab: 'avatar' as 'profile' | 'avatar',
  avatarModalShow: false,
  me: null as any,
  profileForm: {
    name: '',
    nickname: '',
    avatar: '',
    signature: '',
    remark: '',
  },
})

const avatarUrl = computed(() => state.profileForm.avatar || undefined)
const displayName = computed(() => {
  const nickname = String(state.me?.nickname ?? '').trim()
  const name = String(state.me?.name ?? '').trim()
  if (nickname && name && nickname !== name) {
    return `${nickname}（${name}）`
  }
  return nickname || name || '-'
})
const deptText = computed(() => mapNames(state.me?.dept_id_names))
const roleText = computed(() => mapNames(state.me?.role_id_names))
const groupText = computed(() => mapNames(state.me?.group_id_names))
const contactText = computed(() => {
  const profile = state.me?.profile ?? {}
  const parts = [profile.phone, profile.email].filter(Boolean)
  return parts.length ? parts.join(' / ') : ''
})

onMounted(async () => {
  await refresh()
})

async function refresh() {
  state.loading = true
  try {
    const data = await authStore.refreshUserInfo()
    state.me = data
    syncForms(data)
  } finally {
    state.loading = false
  }
}

function syncForms(data: any) {
  const currentProfile = data?.profile ?? {}
  state.profileForm.name = data?.name ?? currentProfile.name ?? ''
  state.profileForm.nickname = data?.nickname ?? currentProfile.nickname ?? ''
  state.profileForm.avatar = data?.avatar ?? currentProfile.avatar ?? ''
  state.profileForm.signature = currentProfile.signature ?? ''
  state.profileForm.remark = currentProfile.remark ?? ''
}

async function saveProfile() {
  state.savingProfile = true
  try {
    await authApi.updateUserCenterProfile({
      name: state.profileForm.name || null,
      nickname: state.profileForm.nickname || null,
      signature: state.profileForm.signature || null,
      remark: state.profileForm.remark || null,
    })
    await refresh()
    window.$message.success('保存成功')
  } finally {
    state.savingProfile = false
  }
}

defineExpose({ refresh })
</script>

<template>
  <NSpin :show="state.loading">
    <NTabs
      v-model:value="state.basicInfoTab"
      type="line"
      animated
      class="user-center__subtabs"
    >
      <NTabPane
        name="avatar"
        tab="头像"
      >
        <div class="user-center__avatar-card">
          <button
            class="user-center__avatar-edit"
            type="button"
            title="更换头像"
            @click="state.avatarModalShow = true"
          >
            <NAvatar
              v-if="avatarUrl"
              round
              :size="160"
              :src="avatarUrl"
              :img-props="avatarImgProps"
            />
            <NAvatar
              v-else
              round
              :size="160"
            >
              <NovaIcon
                icon="icon-park-outline:user"
                :size="64"
              />
            </NAvatar>
            <span class="user-center__avatar-badge">
              <NovaIcon
                icon="icon-park-outline:edit"
                :size="14"
              />
              编辑
            </span>
          </button>
          <div class="user-center__avatar-name">
            {{ displayName }}
          </div>
          <div class="user-center__avatar-account">
            {{ state.me?.account || '-' }}
          </div>
          <NDescriptions
            class="user-center__avatar-desc"
            :column="1"
            label-placement="left"
            size="small"
          >
            <NDescriptionsItem label="部门">
              {{ displayValue(deptText) }}
            </NDescriptionsItem>
            <NDescriptionsItem label="角色">
              {{ displayValue(roleText) }}
            </NDescriptionsItem>
            <NDescriptionsItem label="用户组">
              {{ displayValue(groupText) }}
            </NDescriptionsItem>
            <NDescriptionsItem label="联系方式">
              {{ displayValue(contactText) }}
            </NDescriptionsItem>
          </NDescriptions>
        </div>
      </NTabPane>

      <NTabPane
        name="profile"
        tab="基本信息"
      >
        <NForm
          class="user-center-form user-center-form--narrow w-full min-w-0"
          label-placement="top"
        >
          <NFormItem label="账号">
            <NInput
              :value="state.me?.account"
              disabled
            />
            <template #feedback>
              <span class="user-center__hint">登录账号不可修改。</span>
            </template>
          </NFormItem>
          <NFormItem label="姓名">
            <NInput v-model:value="state.profileForm.name" />
            <template #feedback>
              <span class="user-center__hint">姓名可能出现在审批、审计等场景中。</span>
            </template>
          </NFormItem>
          <NFormItem label="昵称">
            <NInput v-model:value="state.profileForm.nickname" />
          </NFormItem>
          <NFormItem label="个性签名">
            <NInput
              v-model:value="state.profileForm.signature"
              type="textarea"
              :rows="3"
              placeholder="一句话介绍自己"
            />
          </NFormItem>
          <NFormItem label="备注">
            <NInput
              v-model:value="state.profileForm.remark"
              type="textarea"
              :rows="3"
            />
          </NFormItem>
          <NFormItem :show-label="false">
            <NButton
              type="primary"
              :loading="state.savingProfile"
              @click="saveProfile"
            >
              更新资料
            </NButton>
          </NFormItem>
        </NForm>
      </NTabPane>
    </NTabs>
  </NSpin>

  <AvatarUploadModal
    v-model:show="state.avatarModalShow"
    :avatar="avatarUrl"
    @uploaded="refresh"
  />
</template>
