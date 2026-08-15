/** Author: Charlie */

import { useEffect, useState } from 'react'
import { Empty, Skeleton, Tag } from 'antd'
import { NotificationOutlined, RightOutlined } from '@ant-design/icons'
import { Link } from 'react-router-dom'
import { PromoCarousel } from '@/components/common/PromoCarousel'
import { useBannerSlides } from '@/hooks/useBannerSlides'
import { useAuthStore } from '@/stores/auth'
import { formatDateTime } from '@/utils/time'
import { sysNoticeApi } from '@/api'

const HOME_BANNER_QUERY = { position: 'HOME_TOP' } as const

function announcementSummary(content: string, contentType: string) {
  const raw = content || ''
  const text =
    contentType === 'html' || contentType === 'markdown'
      ? raw
          .replace(/<[^>]+>/g, ' ')
          .replace(/\s+/g, ' ')
          .trim()
      : raw.replace(/\s+/g, ' ').trim()
  if (text.length <= 80) return text
  return `${text.slice(0, 80)}…`
}

export function HomePage() {
  const userInfo = useAuthStore((s) => s.userInfo)
  const isLogin = useAuthStore((s) => s.isLogin)
  const ensureSession = useAuthStore((s) => s.ensureSession)
  const loggedIn = isLogin()

  const [announcements, setAnnouncements] = useState<any[]>([])
  const [announceLoading, setAnnounceLoading] = useState(true)
  const { slides: bannerSlides, loading: bannerLoading } = useBannerSlides(HOME_BANNER_QUERY)

  useEffect(() => {
    void ensureSession()
  }, [ensureSession])

  useEffect(() => {
    let mounted = true

    async function loadAnnouncements() {
      setAnnounceLoading(true)
      try {
        const res = await sysNoticeApi.list({ current: 1, size: 5 })
        if (!mounted) return
        setAnnouncements(res.data.records ?? [])
      } catch {
        if (!mounted) return
        setAnnouncements([])
      } finally {
        if (mounted) setAnnounceLoading(false)
      }
    }

    void loadAnnouncements()
    return () => {
      mounted = false
    }
  }, [loggedIn])

  const displayName = userInfo?.nickname || userInfo?.account || '用户'
  const brand = import.meta.env.VITE_APP_TITLE || 'HEI'
  const showBannerColumn = bannerLoading || bannerSlides.length > 0

  return (
    <div className="page-shell flex w-full flex-col gap-5">
      <section
        className={
          showBannerColumn
            ? 'grid grid-cols-1 gap-5 lg:grid-cols-3 lg:items-stretch'
            : 'grid grid-cols-1'
        }
      >
        {showBannerColumn ? (
          <div className="min-w-0 lg:col-span-2">
            {bannerLoading ? (
              <div className="panel h-full overflow-hidden">
                <Skeleton.Node
                  active
                  style={{ width: '100%', height: 280, borderRadius: 0 }}
                  className="!flex !w-full"
                >
                  <div className="h-[280px] w-full" />
                </Skeleton.Node>
              </div>
            ) : (
              <PromoCarousel slides={bannerSlides} height={280} className="h-full" />
            )}
          </div>
        ) : null}

        <section
          className={
            showBannerColumn
              ? 'panel flex h-full min-h-[280px] flex-col justify-center rounded-xl px-5 py-6 lg:col-span-1 lg:px-6 lg:py-8'
              : 'panel flex flex-col justify-center rounded-xl px-6 py-10 md:px-10'
          }
        >
          <div className="text-sm text-[var(--ant-color-text-secondary)]">{brand}</div>
          <h1
            className={
              showBannerColumn
                ? 'mt-2 text-xl font-semibold leading-snug md:text-2xl'
                : 'mt-2 text-2xl font-semibold md:text-3xl'
            }
          >
            {loggedIn ? `${displayName}，欢迎回来` : `欢迎使用 ${brand}`}
          </h1>
          <p className="muted-text mt-3 text-sm leading-6 md:text-base">
            这是 HEI FastAPI 门户脚手架。账号认证、个人中心与公告已就绪，可在此基础上扩展业务模块。
          </p>
          <div className="mt-6 flex flex-wrap gap-3">
            <Link
              to="/announcements"
              className="inline-flex items-center gap-1 rounded-lg bg-[var(--ant-color-primary)] px-4 py-2 text-sm font-medium text-white hover:opacity-90"
            >
              查看公告 <RightOutlined />
            </Link>
            {loggedIn ? (
              <Link
                to="/usercenter"
                className="inline-flex items-center rounded-lg bg-[var(--ant-color-fill-quaternary)] px-4 py-2 text-sm font-medium hover:bg-[var(--ant-color-fill-secondary)]"
              >
                账号设置
              </Link>
            ) : (
              <Link
                to="/auth/login"
                className="inline-flex items-center rounded-lg bg-[var(--ant-color-fill-quaternary)] px-4 py-2 text-sm font-medium hover:bg-[var(--ant-color-fill-secondary)]"
              >
                登录
              </Link>
            )}
          </div>
        </section>
      </section>

      <section className="panel overflow-hidden">
        <div className="flex items-center justify-between border-b border-[var(--ant-color-border)] px-4 py-3">
          <div className="flex items-center gap-2">
            <NotificationOutlined className="text-[var(--ant-color-text-secondary)]" />
            <h2 className="text-base font-semibold">最新公告</h2>
          </div>
          <Link
            to="/announcements"
            className="text-sm text-[var(--ant-color-text-secondary)] hover:text-[var(--ant-color-text)]"
          >
            全部 <RightOutlined />
          </Link>
        </div>
        <Skeleton active loading={announceLoading} paragraph={{ rows: 3 }}>
          {announcements.length ? (
            <div>
              {announcements.map((item) => (
                <Link
                  key={item.id}
                  to={`/announcements/${item.id}`}
                  className="flex w-full flex-col border-b border-[var(--ant-color-border)] px-4 py-3.5 text-left last:border-b-0 transition-colors hover:bg-[var(--ant-color-fill-quaternary)]"
                >
                  <div className="flex items-center gap-2">
                    <div className="min-w-0 flex-1 truncate text-[15px] font-medium">
                      {item.title}
                    </div>
                    {item.is_pinned ? (
                      <Tag color="warning" className="m-0 shrink-0">
                        置顶
                      </Tag>
                    ) : null}
                    {isLogin() && !item.is_read ? (
                      <span className="h-2 w-2 shrink-0 rounded-full bg-[var(--ant-color-primary)]" />
                    ) : null}
                  </div>
                  <div className="muted-text mt-1 line-clamp-2 text-sm">
                    {announcementSummary(item.content, item.content_type) || '点击查看详情'}
                  </div>
                  {item.publish_at ? (
                    <div className="muted-text mt-1 text-xs">{formatDateTime(item.publish_at)}</div>
                  ) : null}
                </Link>
              ))}
            </div>
          ) : !announceLoading ? (
            <div className="py-10">
              <Empty description="暂无公告" />
            </div>
          ) : null}
        </Skeleton>
      </section>
    </div>
  )
}
