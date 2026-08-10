/** Author: Charlie */

import { useEffect, useState } from 'react'
import { Avatar, Empty, Spin } from 'antd'
import { HomeOutlined, UserOutlined } from '@ant-design/icons'
import { Link, useSearchParams } from 'react-router-dom'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api'

export function ProfilePage() {
  const [params] = useSearchParams()
  const userInfo = useAuthStore((s) => s.userInfo)
  const accountId = params.get('account_id') || userInfo?.accountId || ''
  const isSelf = Boolean(userInfo?.accountId && accountId === userInfo.accountId)

  const [loading, setLoading] = useState(Boolean(accountId))
  const [profile, setProfile] = useState<any>(null)
  const [activeAccountId, setActiveAccountId] = useState(accountId)

  if (accountId !== activeAccountId) {
    setActiveAccountId(accountId)
    setLoading(Boolean(accountId))
    setProfile(null)
  }

  useEffect(() => {
    if (!accountId) return
    let cancelled = false
    void (async () => {
      try {
        const spaceRes = await authApi.getPublicSpace(accountId)
        if (!cancelled) setProfile(spaceRes.data)
      } catch {
        if (!cancelled) setProfile(null)
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [accountId])

  const displayName =
    profile?.nickname ||
    profile?.name ||
    (isSelf ? userInfo?.nickname || userInfo?.name || userInfo?.account : null) ||
    '未命名用户'
  const avatarSrc = profile?.avatar || (isSelf ? userInfo?.avatar : undefined) || undefined
  const signature = String(profile?.signature || '').trim()

  if (!accountId) {
    return (
      <div className="flex min-h-[360px] items-center justify-center px-4 py-12">
        <Empty description="请先登录查看个人主页" />
      </div>
    )
  }

  return (
    <div className="min-h-[calc(100vh-64px-72px)] w-full bg-[var(--ant-color-bg-layout)] text-[var(--ant-color-text)]">
      <Spin spinning={loading}>
        {profile || isSelf ? (
          <>
            <section className="relative overflow-hidden text-[#f8fafc]" aria-label="个人资料">
              <div
                className="pointer-events-none absolute inset-0 bg-[linear-gradient(135deg,#0b1f3a_0%,#12325c_42%,#1a4a7a_100%)]"
                aria-hidden
              >
                <div className="absolute inset-0 bg-[radial-gradient(1200px_280px_at_12%_20%,rgba(22,119,255,0.55),transparent_60%)]" />
                <div className="absolute inset-0 bg-[radial-gradient(900px_260px_at_88%_10%,rgba(14,165,233,0.35),transparent_55%)]" />
                <div className="absolute inset-0 bg-[linear-gradient(to_bottom,rgba(15,23,42,0.08),rgba(15,23,42,0.45))]" />
              </div>

              <div className="relative z-[1] flex w-full items-end justify-between gap-4 px-4 py-10 md:px-6 md:pb-7 md:pt-14">
                <div className="flex min-w-0 items-end gap-4">
                  <Avatar
                    size={88}
                    src={avatarSrc || undefined}
                    icon={<UserOutlined />}
                    className="shrink-0 border-[3px] border-solid border-white/90 bg-[color-mix(in_srgb,var(--ant-color-primary)_18%,#fff)] shadow-[0_8px_24px_rgba(15,23,42,0.28)]"
                  />
                  <div className="min-w-0 pb-1">
                    <h1 className="m-0 text-[clamp(1.5rem,2.4vw,2rem)] font-bold leading-tight tracking-tight [text-shadow:0_1px_12px_rgba(15,23,42,0.35)]">
                      {displayName}
                    </h1>
                    {signature ? (
                      <p className="mt-2 max-w-2xl text-[0.925rem] leading-relaxed text-white/85">
                        {signature}
                      </p>
                    ) : null}
                  </div>
                </div>
                {isSelf ? (
                  <Link
                    to="/usercenter"
                    className="mb-1.5 inline-flex shrink-0 items-center border border-white/35 bg-white/12 px-3.5 py-2 text-[0.875rem] text-white backdrop-blur-sm transition-colors hover:border-white/55 hover:bg-white/22"
                  >
                    我的资料
                  </Link>
                ) : null}
              </div>
            </section>

            <nav
              className="sticky top-16 z-20 border-b border-[color-mix(in_srgb,var(--ant-color-border)_70%,transparent)] bg-[color-mix(in_srgb,var(--ant-color-bg-container)_94%,transparent)] backdrop-blur-md"
              aria-label="个人主页导航"
            >
              <div className="flex min-h-12 w-full items-stretch gap-1 px-4 md:px-6">
                <span className="inline-flex items-center gap-1.5 border-b-2 border-[var(--ant-color-primary)] px-3.5 text-[0.925rem] font-semibold text-[var(--ant-color-primary)]">
                  <HomeOutlined />
                  主页
                </span>
              </div>
            </nav>

            <div className="w-full px-4 py-5 md:px-6 md:pb-10">
              <section className="border border-[color-mix(in_srgb,var(--ant-color-border)_65%,transparent)] bg-[var(--ant-color-bg-container)] px-5 py-[18px]">
                <h2 className="mb-3 mt-0 text-base font-semibold">简介</h2>
                {signature ? (
                  <p className="m-0 whitespace-pre-wrap text-[0.925rem] leading-relaxed">
                    {signature}
                  </p>
                ) : (
                  <p className="m-0 text-[0.875rem] text-[var(--ant-color-text-secondary)]">
                    暂未填写签名
                  </p>
                )}
              </section>
            </div>
          </>
        ) : (
          <div className="flex min-h-[360px] items-center justify-center px-4 py-12">
            <Empty description="用户不存在或资料未公开" />
          </div>
        )}
      </Spin>
    </div>
  )
}
