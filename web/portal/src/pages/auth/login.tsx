/** Author: Charlie */

import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { authApi } from '@/api'
import { wireBool } from '@/utils/wire'
import { AuthLoginForm } from './AuthLoginForm'
import { PortalAuthShell } from './PortalAuthShell'

export function LoginPage() {
  const [registerEnabled, setRegisterEnabled] = useState(false)

  useEffect(() => {
    void authApi
      .authOptions()
      .then((res) => {
        setRegisterEnabled(wireBool(res?.data?.register_enabled ?? false))
      })
      .catch(() => undefined)
  }, [])

  return (
    <PortalAuthShell
      title="欢迎登录"
      brandHeadline="登录门户，继续你的工作"
      brandLead="个人中心、公告与反馈，开箱即用。"
      footerNote="登录即表示同意相关服务条款与隐私政策"
      headerExtra={
        registerEnabled ? (
          <Link to="/auth/register" className="linkish">
            没有账号？去注册
          </Link>
        ) : null
      }
    >
      <AuthLoginForm />
    </PortalAuthShell>
  )
}
