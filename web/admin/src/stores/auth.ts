/** Author: Charlie */

import { defineStore } from 'pinia'
import { router } from '@/router'
import { authApi } from '@/api'
import { clearDict, refreshDict, syncDictTree } from '@/utils/dict'
import { clearToken, setToken } from '@/utils/session'
import { wireBool } from '@/utils/wire'
import { useRouteStore } from './route'
import { useTabStore } from './tab'

interface AuthUserInfo {
  accountId: string
  account: string
  accountType: string
  name?: string | null
  nickname?: string | null
  avatar?: string | null
  roleIds: string[]
  deptIds: string[]
  groupIds: string[]
  roleIdNames?: { id: string; name: string }[]
  deptIdNames?: { id: string; name: string }[]
  groupIdNames?: { id: string; name: string }[]
  permissionKeys: string[]
  profile?: Record<string, unknown> | null
  passwordExpired?: boolean
  forceBindEmail?: boolean
  forceBindPhone?: boolean
  loginAt: number
}

interface AuthState {
  userInfo: AuthUserInfo | null
  sessionChecked: boolean
}

const userInfoKey = 'user_info'
const loginPath = '/auth/login'
const userCenterPasswordPath = '/profile?tab=password'
const userCenterEmailPath = '/profile?tab=email'
const userCenterPhonePath = '/profile?tab=phone'

function getStoredUserInfo() {
  const raw = localStorage.getItem(userInfoKey)
  if (!raw) {
    return null
  }

  try {
    return JSON.parse(raw) as AuthUserInfo
  } catch {
    localStorage.removeItem(userInfoKey)
    return null
  }
}

function getSafeRedirect(redirect?: string) {
  if (!redirect || redirect.startsWith('/auth')) {
    return import.meta.env.VITE_HOME_PATH
  }
  return redirect
}

export function resolveSecurityWallPath(user: AuthUserInfo | null | undefined): string | null {
  if (!user) return null
  if (user.passwordExpired) return userCenterPasswordPath
  if (user.forceBindEmail) return userCenterEmailPath
  if (user.forceBindPhone) return userCenterPhonePath
  return null
}

export function isAllowedUnderSecurityWall(
  path: string,
  queryTab: string | undefined | null,
  user: AuthUserInfo | null | undefined,
): boolean {
  if (!user) return true
  if (user.passwordExpired) {
    return path.startsWith('/profile') && queryTab === 'password'
  }
  if (user.forceBindEmail || user.forceBindPhone) {
    if (!path.startsWith('/profile')) return false
    if (user.forceBindEmail && queryTab === 'email') return true
    if (user.forceBindPhone && queryTab === 'phone') return true
    return false
  }
  return true
}

export const useAuthStore = defineStore('auth-store', {
  state: (): AuthState => ({
    userInfo: getStoredUserInfo(),
    sessionChecked: false,
  }),
  getters: {
    isLogin: (state) => Boolean(state.userInfo?.accountId),
  },
  actions: {
    async ensureSession() {
      if (this.sessionChecked) {
        return this.isLogin
      }
      this.sessionChecked = true
      try {
        await this.refreshUserInfo()
        return true
      } catch {
        this.clearAuthStorage()
        return false
      }
    },

    async login(
      account: string,
      password: string,
      redirect?: string,
      rememberMe?: boolean,
      identityType = 'ACCOUNT',
      security?: {
        password_key_id?: string
        captcha_id: string
        captcha_value: string
        login_mode?: 'PASSWORD' | 'OTP'
        otp_code?: string
      },
    ) {
      const persist = rememberMe ?? true
      const response = await authApi.login({
        account,
        password: password || undefined,
        identity_type: identityType,
        remember_me: persist,
        password_key_id: security?.password_key_id,
        captcha_id: security?.captcha_id,
        captcha_value: security?.captcha_value,
        login_mode: security?.login_mode || 'PASSWORD',
        ...(security?.otp_code ? { otp_code: security.otp_code } : {}),
      })

      // Cookie 与 Header 双通道：本地持久化 opaque token，供无 Cookie 时鉴权。
      clearToken()
      if (response.data.token) {
        setToken(String(response.data.token), persist)
      }

      this.sessionChecked = true

      // WireBool 序列化为 "true"/"false" 字符串，不能直接当 JS 真值用
      const passwordExpired = wireBool(response.data.password_expired ?? false)
      const forceBindEmail = wireBool(response.data.force_bind_email ?? false)
      const forceBindPhone = wireBool(response.data.force_bind_phone ?? false)
      const warningDays = response.data.password_expiry_warning_days
      if (passwordExpired) {
        window.$message?.warning?.('密码已过期，请先修改密码')
        await this.finishLogin(userCenterPasswordPath)
        return
      }
      if (forceBindEmail || forceBindPhone) {
        window.$message?.warning?.('请先完成账号安全绑定')
        await this.finishLogin(forceBindEmail ? userCenterEmailPath : userCenterPhonePath)
        return
      }
      if (typeof warningDays === 'number' && warningDays > 0) {
        window.$message?.warning?.(`密码将在 ${warningDays} 天后过期，请及时修改`)
      }

      await this.finishLogin(redirect)
    },

    async finishLogin(redirect?: string) {
      await this.refreshUserInfo()

      const routeStore = useRouteStore()
      await routeStore.initAuthRoute()
      syncDictTree()
      await refreshDict()
      const wall = resolveSecurityWallPath(this.userInfo)
      await router.push(getSafeRedirect(wall || redirect))
    },

    async refreshUserInfo() {
      const meResponse = await authApi.me()
      const userInfo: AuthUserInfo = {
        ...(this.userInfo ?? { loginAt: Date.now() }),
        accountId: meResponse.data.account_id,
        account: meResponse.data.account,
        accountType: meResponse.data.account_type,
        name: meResponse.data.name,
        nickname: meResponse.data.nickname,
        avatar: meResponse.data.avatar,
        roleIds: meResponse.data.role_ids ?? [],
        deptIds: meResponse.data.dept_ids ?? [],
        groupIds: meResponse.data.group_ids ?? [],
        roleIdNames: meResponse.data.role_id_names ?? [],
        deptIdNames: meResponse.data.dept_id_names ?? [],
        groupIdNames: meResponse.data.group_id_names ?? [],
        permissionKeys: meResponse.data.permission_keys ?? [],
        profile: meResponse.data.profile ?? null,
        passwordExpired: wireBool(meResponse.data.password_expired ?? false),
        forceBindEmail: wireBool(meResponse.data.force_bind_email ?? false),
        forceBindPhone: wireBool(meResponse.data.force_bind_phone ?? false),
        loginAt: this.userInfo?.loginAt ?? Date.now(),
      }

      localStorage.setItem(userInfoKey, JSON.stringify(userInfo))
      this.userInfo = userInfo
      return meResponse.data
    },

    hasPermission(permissionKey: string) {
      const keys = this.userInfo?.permissionKeys ?? []
      return keys.includes('*:*:*') || keys.includes(permissionKey)
    },

    clearAuthStorage() {
      clearToken()
      localStorage.removeItem(userInfoKey)
      this.userInfo = null
    },

    resetSession() {
      this.clearAuthStorage()
      this.sessionChecked = true

      const routeStore = useRouteStore()
      routeStore.resetRouteStore()

      const tabStore = useTabStore()
      tabStore.clearAllTabs()

      clearDict()
    },

    async logout(redirect?: string) {
      const currentRoute = router.currentRoute.value
      const finalRedirect = redirect ?? currentRoute.fullPath

      try {
        await authApi.logout()
      } catch {
        // 后端登出失败不阻塞本地会话清理。
      } finally {
        this.resetSession()
      }

      await router.push({
        path: loginPath,
        query: finalRedirect.startsWith('/auth') ? undefined : { redirect: finalRedirect },
      })
    },
  },
})
