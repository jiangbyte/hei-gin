/** Author: Charlie */

import { useEffect, useState } from 'react'
import { Empty, Pagination, Skeleton, Tag } from 'antd'
import { Link } from 'react-router-dom'
import { useAuthStore } from '@/stores/auth'
import { formatDateTime } from '@/utils/time'
import { readPageMeta } from '@/utils/wire'
import { sysNoticeApi } from '@/api'

function announcementSummary(content: string, contentType: string) {
  const raw = content || ''
  const text =
    contentType === 'html' || contentType === 'markdown'
      ? raw
          .replace(/<[^>]+>/g, ' ')
          .replace(/\s+/g, ' ')
          .trim()
      : raw.replace(/\s+/g, ' ').trim()
  if (text.length <= 120) return text
  return `${text.slice(0, 120)}…`
}

export function AnnouncementListPage() {
  const isLogin = useAuthStore((s) => s.isLogin)
  const loggedIn = isLogin()

  const [loading, setLoading] = useState(true)
  const [records, setRecords] = useState<any[]>([])
  const [total, setTotal] = useState(0)
  const [current, setCurrent] = useState(1)
  const [size, setSize] = useState(10)

  useEffect(() => {
    let mounted = true

    async function load() {
      setLoading(true)
      try {
        const res = await sysNoticeApi.list({ current, size })
        if (!mounted) return
        setRecords(res.data.records ?? [])
        setTotal(readPageMeta(res.data).total)
      } catch {
        if (!mounted) return
        setRecords([])
        setTotal(0)
      } finally {
        if (mounted) setLoading(false)
      }
    }

    void load()
    return () => {
      mounted = false
    }
  }, [current, size, loggedIn])

  return (
    <div className="page-shell">
      <div className="mb-4">
        <h1 className="text-xl font-semibold">公告</h1>
        <p className="muted-text mt-1 text-sm">查看平台发布的最新通知与公告</p>
      </div>

      <section className="panel overflow-hidden">
        <Skeleton active loading={loading} paragraph={{ rows: 6 }}>
          {records.length ? (
            <div>
              {records.map((item) => (
                <Link
                  key={item.id}
                  to={`/announcements/${item.id}`}
                  className="flex flex-col border-b border-[var(--ant-color-border)] px-4 py-4 last:border-b-0 transition-colors hover:bg-[var(--ant-color-fill-quaternary)]"
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
                    <div className="muted-text mt-2 text-xs">{formatDateTime(item.publish_at)}</div>
                  ) : null}
                </Link>
              ))}
            </div>
          ) : !loading ? (
            <div className="py-16">
              <Empty description="暂无公告" />
            </div>
          ) : null}
        </Skeleton>
      </section>

      {total > 0 ? (
        <div className="mt-4 flex justify-end">
          <Pagination
            current={current}
            pageSize={size}
            total={total}
            showSizeChanger
            showTotal={(t) => `共 ${t} 条`}
            onChange={(page, pageSize) => {
              setCurrent(page)
              setSize(pageSize)
            }}
          />
        </div>
      ) : null}
    </div>
  )
}
