/** Author: Charlie */

import { Link } from 'react-router-dom'
import { MailOutlined, PhoneOutlined } from '@ant-design/icons'
import './app-footer.css'

const brand = import.meta.env.VITE_APP_TITLE || 'HEI'
const copyright = import.meta.env.VITE_COPYRIGHT_INFO || `© ${new Date().getFullYear()} ${brand}`
const icp = import.meta.env.VITE_ICP_FILING || ''
const supportEmail = import.meta.env.VITE_SUPPORT_EMAIL || 'support@example.com'
const supportPhone = import.meta.env.VITE_SUPPORT_PHONE || ''

const linkGroups = [
  {
    title: '产品服务',
    items: ['功能介绍', '使用指南', '更新日志'],
  },
  {
    title: '帮助支持',
    items: ['常见问题', '意见反馈', '技术支持'],
  },
  {
    title: '账号相关',
    items: ['账号安全', '隐私设置', '权限说明'],
  },
]

export function AppFooter() {
  return (
    <footer className="app-footer">
      <div className="app-footer__main">
        <div className="app-footer__inner">
          <div className="app-footer__brand">
            <Link to="/" className="app-footer__logo">
              <span className="app-footer__logo-mark">{brand.slice(0, 1).toUpperCase()}</span>
              <span className="app-footer__logo-text">{brand}</span>
            </Link>
            <p className="app-footer__tagline">登录注册、个人中心与公告，开箱即用。</p>
            <p className="app-footer__copy">{copyright}</p>
          </div>

          <div className="app-footer__cols">
            {linkGroups.map((group) => (
              <div key={group.title} className="app-footer__col">
                <h3 className="app-footer__col-title">{group.title}</h3>
                <ul className="app-footer__list">
                  {group.items.map((label) => (
                    <li key={label}>
                      <span className="app-footer__placeholder">{label}</span>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>

          <div className="app-footer__contact">
            <h3 className="app-footer__col-title">联系我们</h3>
            <ul className="app-footer__contact-list">
              <li>
                <MailOutlined aria-hidden />
                <a href={`mailto:${supportEmail}`}>{supportEmail}</a>
              </li>
              {supportPhone ? (
                <li>
                  <PhoneOutlined aria-hidden />
                  <a href={`tel:${supportPhone.replace(/\s+/g, '')}`}>{supportPhone}</a>
                </li>
              ) : (
                <li>
                  <PhoneOutlined aria-hidden />
                  <span>工作日 9:00 – 18:00</span>
                </li>
              )}
            </ul>
          </div>
        </div>
      </div>

      <div className="app-footer__bar">
        <div className="app-footer__inner app-footer__bar-inner">
          <p className="app-footer__bar-left">{brand} Portal · 统一门户服务</p>
          <p className="app-footer__bar-center">
            <span>服务条款</span>
            <span className="app-footer__dot">·</span>
            <span>隐私政策</span>
            {icp ? (
              <>
                <span className="app-footer__dot">·</span>
                <span>{icp}</span>
              </>
            ) : null}
          </p>
        </div>
      </div>
    </footer>
  )
}
