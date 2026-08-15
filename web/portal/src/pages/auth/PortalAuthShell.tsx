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
  return (
    <div className={variant === 'center' ? 'portal-auth portal-auth--center' : 'portal-auth'}>
      {/* 页面级波浪海洋背景（装饰层，不参与布局；静态三层波浪） */}
      <div className="portal-auth__waves" aria-hidden>
        <svg className="portal-auth__waves-svg portal-auth__waves-svg--far" viewBox="0 0 1440 320" preserveAspectRatio="none">
          <path d="M0,224L60,213.3C120,203,240,181,360,186.7C480,192,600,224,720,229.3C840,235,960,213,1080,197.3C1200,181,1320,171,1380,165.3L1440,160L1440,320L0,320Z" />
        </svg>
        <svg className="portal-auth__waves-svg portal-auth__waves-svg--mid" viewBox="0 0 1440 320" preserveAspectRatio="none">
          <path d="M0,256L80,245.3C160,235,320,213,480,213.3C640,213,800,235,960,245.3C1120,256,1280,256,1360,256L1440,256L1440,320L0,320Z" />
        </svg>
        <svg className="portal-auth__waves-svg portal-auth__waves-svg--near" viewBox="0 0 1440 320" preserveAspectRatio="none">
          <path d="M0,288L60,282.7C120,277,240,267,360,272C480,277,600,299,720,298.7C840,299,960,277,1080,266.7C1200,256,1320,256,1380,256L1440,256L1440,320L0,320Z" />
        </svg>
      </div>

      {variant === 'center' ? (
        <>
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
        </>
      ) : (
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
          <svg
            className="portal-auth__brand-waves"
            viewBox="0 0 1440 180"
            preserveAspectRatio="none"
            aria-hidden
          >
            <path
              className="portal-auth__brand-wave portal-auth__brand-wave--far"
              d="M0,96L60,88C120,80,240,64,360,69.3C480,75,600,101,720,106.7C840,112,960,96,1080,88C1200,80,1320,80,1380,80L1440,80L1440,180L0,180Z"
            />
            <path
              className="portal-auth__brand-wave portal-auth__brand-wave--near"
              d="M0,128L80,120C160,112,320,96,480,101.3C640,107,800,133,960,138.7C1120,144,1280,128,1360,124L1440,120L1440,180L0,180Z"
            />
          </svg>
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
      )}
    </div>
  )
}