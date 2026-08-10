/** Author: Charlie */

import { useEffect, useState } from 'react'
import { Button, Empty, Pagination, Skeleton } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { Link } from 'react-router-dom'
import { feedbackApi } from '@/api'
import { DictTag } from '@/components/common/DictTag'
import { formatDateTime } from '@/utils/time'
import { readPageMeta } from '@/utils/wire'

export function FeedbackListPage() {
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
        const res = await feedbackApi.myPage({ current, size })
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
  }, [current, size])

  return (
    <div className="page-shell">
      <div className="mb-4 flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-xl font-semibold">意见反馈</h1>
          <p className="muted-text mt-1 text-sm">提交问题与建议，并查看处理进度</p>
        </div>
        <Link to="/feedback/new">
          <Button type="primary" icon={<PlusOutlined />}>
            提交反馈
          </Button>
        </Link>
      </div>

      <section className="panel overflow-hidden">
        <Skeleton active loading={loading} paragraph={{ rows: 6 }}>
          {records.length ? (
            <div>
              {records.map((item) => (
                <Link
                  key={item.id}
                  to={`/feedback/${item.id}`}
                  className="flex flex-col border-b border-[var(--ant-color-border)] px-4 py-4 last:border-b-0 transition-colors hover:bg-[var(--ant-color-fill-quaternary)]"
                >
                  <div className="flex flex-wrap items-center gap-2">
                    <div className="min-w-0 flex-1 truncate text-[15px] font-medium">
                      {item.title}
                    </div>
                    <DictTag dictCode="FEEDBACK_CATEGORY" value={item.category} />
                    <DictTag dictCode="FEEDBACK_STATUS" value={item.status} />
                  </div>
                  <div className="muted-text mt-1 line-clamp-2 text-sm">
                    {String(item.content || '')
                      .replace(/\s+/g, ' ')
                      .trim()
                      .slice(0, 120) || '点击查看详情'}
                  </div>
                  <div className="muted-text mt-2 text-xs">{formatDateTime(item.created_at)}</div>
                </Link>
              ))}
            </div>
          ) : !loading ? (
            <div className="py-16">
              <Empty description="暂无反馈">
                <Link to="/feedback/new">
                  <Button type="primary">提交第一条反馈</Button>
                </Link>
              </Empty>
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
