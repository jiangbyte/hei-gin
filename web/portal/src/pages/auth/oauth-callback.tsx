/** Author: Charlie */

import { useEffect, useState } from 'react'
import { Result, Spin, message } from 'antd'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { oauthExchange } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import { clearToken, setToken } from '@/utils/session'
import { clearAuthStorage } from '@/utils/storage'
import { getSafeRedirect } from '@/utils/validate'
import { wireBool } from '@/utils/wire'
import { refreshDict, syncDictTree } from '@/utils/dict'
import { PortalAuthShell } from './PortalAuthShell'

/**
 * OAuth 前端回调页：用 oauth_code 兑换 token，再完成会话。
 */
export function OAuthCallbackPage() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const refreshUserInfo = useAuthStore((s) => s.refreshUserInfo)
  const resolveSecurityRedirect = useAuthStore((s) => s.resolveSecurityRedirect)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    void (async () => {
      const status = searchParams.get('oauth_status')
      const action = searchParams.get('oauth_action') || ''
      const rawMessage = searchParams.get('oauth_message')
      const oauthCode = searchParams.get('oauth_code')
      // 兼容旧版 URL token（过渡期）
      const legacyToken = searchParams.get('token')
      const redirect = searchParams.get('redirect')

      if (status !== 'ok') {
        let msg = '三方登录失败'
        if (rawMessage) {
          try {
            msg = decodeURIComponent(rawMessage)
          } catch {
            msg = rawMessage
          }
        }
        if (!cancelled) {
          setError(msg)
          message.error(msg)
          window.setTimeout(() => navigate('/auth/login', { replace: true }), 1600)
        }
        return
      }

      try {
        let passwordExpired = false
        let forceBindEmail = false
        let forceBindPhone = false

        if (oauthCode) {
          const { data } = await oauthExchange({ code: oauthCode })
          const token = data?.token
          if (!token) {
            throw new Error('兑换登录凭证失败')
          }
          clearToken()
          clearAuthStorage()
          setToken(String(token), true)
          passwordExpired = wireBool(data?.password_expired ?? false)
          forceBindEmail = wireBool(data?.force_bind_email ?? false)
          forceBindPhone = wireBool(data?.force_bind_phone ?? false)
        } else if (legacyToken) {
          clearToken()
          clearAuthStorage()
          setToken(legacyToken, true)
          passwordExpired = wireBool(searchParams.get('password_expired') ?? false)
          forceBindEmail = wireBool(searchParams.get('force_bind_email') ?? false)
          forceBindPhone = wireBool(searchParams.get('force_bind_phone') ?? false)
        }

        if (action === 'bound') {
          await refreshUserInfo()
          message.success('绑定成功')
          if (!cancelled) navigate('/usercenter?tab=oauth', { replace: true })
          return
        }

        if (passwordExpired) {
          message.warning('密码已过期，请先修改密码')
        } else if (forceBindEmail || forceBindPhone) {
          message.warning('请先完成账号安全绑定')
        } else {
          message.success('登录成功')
        }

        await refreshUserInfo()
        syncDictTree()
        await refreshDict()
        const next = resolveSecurityRedirect(getSafeRedirect(redirect))
        if (!cancelled) navigate(next, { replace: true })
      } catch (e: any) {
        const msg = e?.message || '登录会话建立失败'
        if (!cancelled) {
          setError(msg)
          message.error(msg)
          window.setTimeout(() => navigate('/auth/login', { replace: true }), 1600)
        }
      }
    })()
    return () => {
      cancelled = true
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <PortalAuthShell variant="center" title="三方登录" description={error || '正在完成登录…'}>
      {error ? (
        <Result status="error" title="三方登录失败" subTitle={error} />
      ) : (
        <div style={{ display: 'grid', placeItems: 'center', minHeight: 120 }}>
          <Spin size="large" tip="正在完成登录…" />
        </div>
      )}
    </PortalAuthShell>
  )
}
