/** Author: Charlie */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '@/api/auth'
import { getStorage, setStorage } from '@/utils/storage'
import { clearDict } from '@/utils/dict'
import { encryptPasswords } from '@/utils/security'
import {
  clearSessionStorage,
  getToken,
  onSessionCleared,
  setToken,
  userInfoKey,
} from '@/utils/session'

export interface AuthUserInfo {
  accountId: string
  account: string
  accountType: string
  name?: string | null
  nickname?: string | null
  avatar?: string | null
  roleIds: string[]
  deptIds: string[]
  groupIds: string[]
  roleIdNames: any[]
  deptIdNames: any[]
  groupIdNames: any[]
  permissionKeys: string[]
  profile?: Record<string, any> | null
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(getToken())
  const userInfo = ref<AuthUserInfo | null>(
    getStorage<AuthUserInfo | null>(userInfoKey, null)
  )

  const isLogin = computed(() => Boolean(token.value))

  onSessionCleared(() => {
    token.value = ''
    userInfo.value = null
  })

  async function login(payload: {
    account: string
    password: string
    captcha_id: string
    captcha_value: string
    identity_type?: string
  }) {
    const security = await encryptPasswords({ password: payload.password })
    const response = await authApi.login({
      account: payload.account,
      password: security.values.password,
      identity_type: payload.identity_type ?? 'ACCOUNT',
      password_key_id: security.password_key_id,
      captcha_id: payload.captcha_id,
      captcha_value: payload.captcha_value,
    })
    if (!response?.token) {
      throw new Error('登录失败：未返回会话令牌')
    }
    token.value = response.token
    setToken(response.token)
    await refreshUserInfo()
  }

  async function refreshUserInfo() {
    const data = await authApi.me()
    const nextUser: AuthUserInfo = {
      accountId: data.account_id,
      account: data.account,
      accountType: data.account_type,
      name: data.name,
      nickname: data.nickname,
      avatar: data.avatar,
      roleIds: data.role_ids ?? [],
      deptIds: data.dept_ids ?? [],
      groupIds: data.group_ids ?? [],
      roleIdNames: data.role_id_names ?? [],
      deptIdNames: data.dept_id_names ?? [],
      groupIdNames: data.group_id_names ?? [],
      permissionKeys: data.permission_keys ?? [],
      profile: data.profile ?? null,
    }
    userInfo.value = nextUser
    setStorage(userInfoKey, nextUser)
    return nextUser
  }

  function hasPermission(permissionKey?: string) {
    if (!permissionKey) {
      return true
    }
    const keys = userInfo.value?.permissionKeys ?? []
    return keys.includes('*:*:*') || keys.includes(permissionKey)
  }

  function resetSession() {
    clearSessionStorage()
    clearDict()
  }

  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // 本地会话清理优先。
    } finally {
      resetSession()
      uni.reLaunch({ url: '/pages/auth/login/login' })
    }
  }

  return {
    token,
    userInfo,
    isLogin,
    login,
    refreshUserInfo,
    hasPermission,
    resetSession,
    logout,
  }
})
