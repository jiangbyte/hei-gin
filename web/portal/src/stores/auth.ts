/** Author: Charlie */

import { create } from 'zustand'
import { message } from 'antd'
import { authApi } from '@/api'
import { clearDict, refreshDict, syncDictTree } from '@/utils/dict'
import { clearToken, setToken } from '@/utils/session'
import { clearAuthStorage, getStoredUserInfo, setStoredUserInfo } from '@/utils/storage'
import { getSafeRedirect } from '@/utils/validate'
import { wireBool } from '@/utils/wire'

const userCenterPasswordPath = '/usercenter?tab=password'
const userCenterEmailPath = '/usercenter?tab=email'
const userCenterPhonePath = '/usercenter?tab=phone'

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
  roleIdNames?: { id: string; name: string }[]
  deptIdNames?: { id: string; name: string }[]
  groupIdNames?: { id: string; name: string }[]
  profile?: Record<string, unknown> | null
  passwordExpired?: boolean
  forceBindEmail?: boolean
  forceBindPhone?: boolean
  loginAt: number
}

interface AuthState {
  userInfo: AuthUserInfo | null
  sessionChecked: boolean
  /** 主动退出中：屏蔽并发请求的 401 提示 */
  loggingOut: boolean
  isLogin: () => boolean
  ensureSession: () => Promise<boolean>
  login: (
    account: string,
    password: string,
    redirect?: string,
    rememberMe?: boolean,
    identityType?: string,
    security?: {
      password_key_id?: string
      captcha_id: string
      captcha_value: string
      login_mode?: 'PASSWORD' | 'OTP'
      otp_code?: string
    },
  ) => Promise<string>
  refreshUserInfo: () => Promise<any>
  logout: (redirect?: string) => Promise<void>
  resetSession: () => void
  resolveSecurityRedirect: (fallback?: string) => string
}

function mapMe(data: any, loginAt = Date.now()): AuthUserInfo {
  return {
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
    profile: data.profile ?? null,
    passwordExpired: wireBool(data.password_expired ?? false),
    forceBindEmail: wireBool(data.force_bind_email ?? false),
    forceBindPhone: wireBool(data.force_bind_phone ?? false),
    loginAt,
  }
}

export function resolveSecurityWallPath(user: AuthUserInfo | null | undefined): string | null {
  if (!user) return null
  if (user.passwordExpired) return userCenterPasswordPath
  if (user.forceBindEmail) return userCenterEmailPath
  if (user.forceBindPhone) return userCenterPhonePath
  return null
}

export function isAllowedUnderSecurityWall(
  pathname: string,
  search: string,
  user: AuthUserInfo | null | undefined,
): boolean {
  if (!user) return true
  if (user.passwordExpired) {
    return pathname.startsWith('/usercenter') && new URLSearchParams(search).get('tab') === 'password'
  }
  if (user.forceBindEmail || user.forceBindPhone) {
    if (!pathname.startsWith('/usercenter')) return false
    const tab = new URLSearchParams(search).get('tab')
    if (user.forceBindEmail && tab === 'email') return true
    if (user.forceBindPhone && tab === 'phone') return true
    return false
  }
  return true
}

export const useAuthStore = create<AuthState>((set, get) => ({
  userInfo: getStoredUserInfo<AuthUserInfo>(),
  sessionChecked: false,
  loggingOut: false,

  isLogin: () => Boolean(get().userInfo?.accountId),

  resolveSecurityRedirect: (fallback) => {
    return resolveSecurityWallPath(get().userInfo) || getSafeRedirect(fallback)
  },

  ensureSession: async () => {
    if (get().sessionChecked) {
      return get().isLogin()
    }
    set({ sessionChecked: true })
    try {
      // 静默探测会话（Cookie 和/或本地 Authorization）；未登录时的 401 属预期
      const meResponse = await authApi.me({ probe: true })
      const userInfo = mapMe(meResponse.data, get().userInfo?.loginAt ?? Date.now())
      setStoredUserInfo(userInfo)
      set({ userInfo })
      return true
    } catch {
      clearToken()
      clearAuthStorage()
      set({ userInfo: null })
      return false
    }
  },

  login: async (
    account,
    password,
    redirect,
    rememberMe = true,
    identityType = 'ACCOUNT',
    security,
  ) => {
    const response = await authApi.login({
      account,
      password: password || undefined,
      identity_type: identityType,
      remember_me: rememberMe,
      password_key_id: security?.password_key_id,
      captcha_id: security?.captcha_id || '',
      captcha_value: security?.captcha_value || '',
      login_mode: security?.login_mode || 'PASSWORD',
      ...(security?.otp_code ? { otp_code: security.otp_code } : {}),
    })

    // Cookie 与 Header 双通道：本地持久化 opaque token，供无 Cookie 时鉴权。
    clearToken()
    clearAuthStorage()
    if (response.data.token) {
      setToken(String(response.data.token), rememberMe)
    }
    set({ sessionChecked: true })

    // WireBool 序列化为 "true"/"false" 字符串，不能直接当 JS 真值用
    const passwordExpired = wireBool(response.data.password_expired ?? false)
    const forceBindEmail = wireBool(response.data.force_bind_email ?? false)
    const forceBindPhone = wireBool(response.data.force_bind_phone ?? false)
    const warningDays = response.data.password_expiry_warning_days
    if (passwordExpired) {
      message.warning('密码已过期，请先修改密码')
    } else if (forceBindEmail || forceBindPhone) {
      message.warning('请先完成账号安全绑定')
    } else if (typeof warningDays === 'number' && warningDays > 0) {
      message.warning(`密码将在 ${warningDays} 天后过期，请及时修改`)
    }

    await get().refreshUserInfo()

    syncDictTree()
    await refreshDict()

    return get().resolveSecurityRedirect(redirect)
  },

  refreshUserInfo: async () => {
    const meResponse = await authApi.me()
    const userInfo = mapMe(meResponse.data, get().userInfo?.loginAt ?? Date.now())
    setStoredUserInfo(userInfo)
    set({ userInfo })
    return meResponse.data
  },

  resetSession: () => {
    clearToken()
    clearAuthStorage()
    clearDict()
    set({ userInfo: null, sessionChecked: true })
  },

  logout: async (redirectTo) => {
    set({ loggingOut: true })
    try {
      await authApi.logout()
    } catch {
      // 忽略
    } finally {
      get().resetSession()
      set({ loggingOut: false })
    }

    const query =
      redirectTo && !redirectTo.startsWith('/auth')
        ? `?redirect=${encodeURIComponent(redirectTo)}`
        : ''
    window.location.assign(`/auth/login${query}`)
  },
}))
