/** Author: Charlie */

import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import './auth-page.css'

type SplitProps = {
  title: string
  headerExtra?: ReactNode
  copyright?: string
  copyrightUrl?: string
  children: ReactNode
}

const brandName = import.meta.env.VITE_APP_TITLE || 'HEI'

function resolveCopyrightHref(url?: string) {
  const value = (url || '').trim()
  if (!value) return ''
  if (/^https?:\/\//i.test(value)) return value
  return `https://${value}`
}

/** 登录 / 注册：左右分栏卡片 */
export function AuthSplit({ title, headerExtra, copyright, copyrightUrl, children }: SplitProps) {
  const href = resolveCopyrightHref(copyrightUrl)
  return (
    <div className="auth-page">
      <div className="auth-card">
        <aside className="auth-card__brand">
          <div className="auth-card__brand-deco" aria-hidden />
          <div className="auth-card__brand-inner">
            <Link to="/" className="auth-card__logo">
              <span className="auth-card__logo-mark">{brandName.slice(0, 1).toUpperCase()}</span>
              <span className="auth-card__logo-text">{brandName}</span>
            </Link>
            <p className="auth-card__eyebrow">Portal</p>
            <h2 className="auth-card__headline">门户服务台</h2>
            <p className="auth-card__lead">登录注册、个人中心与公告，开箱即用。</p>
            <div className="auth-card__brand-foot">
              {copyright ? (
                <div style={{ marginBottom: 8 }}>
                  {href ? (
                    <a
                      className="auth-card__copyright-link"
                      href={href}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      {copyright}
                    </a>
                  ) : (
                    copyright
                  )}
                </div>
              ) : null}
              <Link to="/" className="auth-card__brand-link">
                进入门户首页
              </Link>
            </div>
          </div>
        </aside>

        <div className="auth-card__form">
          <div className="auth-card__form-head">
            <h1 className="auth-card__title">{title}</h1>
            {headerExtra ? <div className="auth-card__form-extra">{headerExtra}</div> : null}
          </div>
          {children}
        </div>
      </div>
    </div>
  )
}

type CenterProps = {
  title: string
  description?: string
  children: ReactNode
}

/** 找回 / 重置密码：居中简版 */
export function AuthCenter({ title, description, children }: CenterProps) {
  return (
    <div className="auth-page auth-page--center">
      <div className="auth-center">
        <Link to="/" className="auth-center__logo">
          <span className="auth-center__logo-mark">{brandName.slice(0, 1).toUpperCase()}</span>
        </Link>
        <h1 className="auth-center__title">{title}</h1>
        {description ? <p className="auth-center__desc">{description}</p> : null}
        <div className="auth-center__body">{children}</div>
      </div>
    </div>
  )
}
