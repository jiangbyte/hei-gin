/** Author: Charlie */

import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import './auth-page.css'

const brandName = import.meta.env.VITE_APP_TITLE || 'HEI'

type ShellProps = {
  variant?: 'split' | 'center'
  title: string
  description?: string
  headerExtra?: ReactNode
  brandHeadline?: string
  brandLead?: string
  footerNote?: string
  children: ReactNode
}

/**
 * 门户全屏认证壳：表单左 / 品牌右（小屏品牌收为顶条）。
 */
export function PortalAuthShell({
  variant = 'split',
  title,
  description,
  headerExtra,
  brandHeadline = '登录门户，继续你的工作',
  brandLead = '个人中心、公告与反馈，开箱即用。',
  footerNote,
  children,
}: ShellProps) {
  if (variant === 'center') {
    return (
      <div className="portal-auth portal-auth--center">
        <div className="portal-auth__topbar">
          <Link to="/" className="portal-auth__brand-link">
            <span className="portal-auth__mark">{brandName.slice(0, 1).toUpperCase()}</span>
            <span className="portal-auth__name">{brandName}</span>
          </Link>
        </div>
        <main className="portal-auth__center-card portal-auth__enter">
          <h1 className="portal-auth__title">{title}</h1>
          {description ? <p className="portal-auth__desc">{description}</p> : null}
          <div className="portal-auth__body">{children}</div>
        </main>
      </div>
    )
  }

  return (
    <div className="portal-auth">
      <div className="portal-auth__stage portal-auth__enter">
        <section className="portal-auth__panel portal-auth__panel--form">
          <div className="portal-auth__mobile-brand">
            <Link to="/" className="portal-auth__brand-link">
              <span className="portal-auth__mark">{brandName.slice(0, 1).toUpperCase()}</span>
              <span className="portal-auth__name">{brandName}</span>
            </Link>
          </div>
          <header className="portal-auth__head">
            <h1 className="portal-auth__title">{title}</h1>
            {headerExtra ? <div className="portal-auth__head-extra">{headerExtra}</div> : null}
          </header>
          {description ? <p className="portal-auth__desc">{description}</p> : null}
          <div className="portal-auth__body">{children}</div>
          {footerNote ? <p className="portal-auth__legal">{footerNote}</p> : null}
        </section>

        <aside className="portal-auth__panel portal-auth__panel--brand" aria-hidden={false}>
          <div className="portal-auth__geo" aria-hidden />
          <div className="portal-auth__brand-inner">
            <Link to="/" className="portal-auth__brand-link portal-auth__brand-link--on-dark">
              <span className="portal-auth__mark">{brandName.slice(0, 1).toUpperCase()}</span>
              <span className="portal-auth__name">{brandName}</span>
            </Link>
            <div className="portal-auth__brand-copy">
              <p className="portal-auth__eyebrow">Portal</p>
              <h2 className="portal-auth__headline">{brandHeadline}</h2>
              <p className="portal-auth__lead">{brandLead}</p>
            </div>
            <Link to="/" className="portal-auth__home-cta">
              返回首页
            </Link>
          </div>
        </aside>
      </div>
    </div>
  )
}
