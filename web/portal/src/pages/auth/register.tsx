/** Author: Charlie */

import { Link } from 'react-router-dom'
import { AuthRegisterForm } from './AuthRegisterForm'
import { PortalAuthShell } from './PortalAuthShell'

export function RegisterPage() {
  return (
    <PortalAuthShell
      title="注册账号"
      brandHeadline="加入门户，开启更多能力"
      brandLead="注册后即可使用个人中心、公告与反馈。"
      footerNote="注册即表示同意相关服务条款与隐私政策"
      headerExtra={
        <Link to="/auth/login" className="linkish">
          已有账号？去登录
        </Link>
      }
    >
      <AuthRegisterForm />
    </PortalAuthShell>
  )
}
