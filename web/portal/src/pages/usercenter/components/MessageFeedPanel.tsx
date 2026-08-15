/** Author: Charlie */

import { useCallback, useEffect, useState } from 'react'
import { Badge, Button, Empty, List, Pagination, Space, Spin, Tag, Typography, message } from 'antd'
import { NotificationOutlined, SoundOutlined } from '@ant-design/icons'
import { myNoticeApi } from '@/api'
import {
  MessageDetailModal,
  type MessageDetailSource,
} from '@/components/sys/MessageDetailModal'
import { useMessageUnreadStore } from '@/stores/messageUnread'
import { formatDateTime, wireBool } from '@/utils'
import { dictTypeData } from '@/utils/dict'
import { readPageMeta } from '@/utils/wire'

export function MessageFeedPanel() {
  const notifyReadAll = useMessageUnreadStore((s) => s.notifyReadAll)
  const [rows, setRows] = useState<any[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [detailOpen, setDetailOpen] = useState(false)
  const [detailSource, setDetailSource] = useState<MessageDetailSource | null>(null)
  const requestKey = `${page}|${pageSize}`
  const [activeKey, setActiveKey] = useState(requestKey)

  if (activeKey !== requestKey) {
    setActiveKey(requestKey)
    setLoading(true)
  }

  const applyPage = useCallback((data: any, current: number, size: number) => {
    setRows(
      (data.records ?? []).map((row: any) => ({
        ...row,
        is_read: wireBool(row.is_read ?? false),
      })),
    )
    const pageMeta = readPageMeta(data, { current, size })
    setTotal(pageMeta.total)
    if (pageMeta.current !== current) setPage(pageMeta.current)
    if (pageMeta.size !== size) setPageSize(pageMeta.size)
  }, [])

  const fetchPage = useCallback(async () => {
    setLoading(true)
    try {
      const response = await myNoticeApi.myPage({ current: page, size: pageSize })
      applyPage(response.data ?? {}, page, pageSize)
    } finally {
      setLoading(false)
    }
  }, [applyPage, page, pageSize])

  useEffect(() => {
    let cancelled = false
    void (async () => {
      try {
        const response = await myNoticeApi.myPage({ current: page, size: pageSize })
        if (!cancelled) applyPage(response.data ?? {}, page, pageSize)
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [applyPage, page, pageSize])

  function kindLabel(row: any) {
    return row.kind === 'ANNOUNCEMENT' ? '公告' : '通知'
  }

  function severityLabel(row: any) {
    if (!row.severity) return ''
    return dictTypeData('NOTIFICATION_SEVERITY', row.severity) || row.severity
  }

  function excerpt(row: any) {
    const text = String(row.content || '')
      .replace(/<[^>]+>/g, ' ')
      .replace(/\s+/g, ' ')
      .trim()
    if (!text) return ''
    return text.length > 96 ? `${text.slice(0, 96)}…` : text
  }

  function openDetail(row: any) {
    setDetailSource({
      id: row.id,
      sourceType: row.kind === 'ANNOUNCEMENT' ? 'ANNOUNCEMENT' : 'NOTIFICATION',
      title: row.title,
      is_read: row.is_read,
      publish_at: row.publish_at,
      content: row.content,
      severity: row.severity,
    })
    setDetailOpen(true)
  }

  function handleDetailChanged(payload: { type: string; id: string }) {
    setRows((prev) =>
      prev.map((item) => (item.id === payload.id ? { ...item, is_read: true } : item)),
    )
  }

  async function markAllRead() {
    await myNoticeApi.readAll()
    setRows((prev) => prev.map((row) => ({ ...row, is_read: true })))
    notifyReadAll()
    message.success('已全部标记为已读')
  }

  return (
    <Space direction="vertical" size={12} className="w-full min-w-0">
      <div className="flex justify-end gap-4">
        <Button type="link" loading={loading} onClick={() => void fetchPage()}>
          刷新
        </Button>
        <Button type="link" onClick={() => void markAllRead()}>
          全部已读
        </Button>
      </div>

      <Spin spinning={loading}>
        {!loading && !rows.length ? (
          <Empty description="暂无消息" />
        ) : (
          <List
            className="w-full min-w-0"
            dataSource={rows}
            renderItem={(row) => (
              <List.Item key={row.id} style={{ cursor: 'pointer' }} onClick={() => openDetail(row)}>
                <List.Item.Meta
                  avatar={
                    <Badge dot={!row.is_read}>
                      <div
                        className="flex h-9 w-9 items-center justify-center rounded-full"
                        style={{
                          background: 'var(--ant-color-fill-quaternary)',
                          color: row.is_read
                            ? 'var(--ant-color-text-tertiary)'
                            : 'var(--ant-color-primary)',
                        }}
                      >
                        {row.kind === 'ANNOUNCEMENT' ? <SoundOutlined /> : <NotificationOutlined />}
                      </div>
                    </Badge>
                  }
                  title={
                    <div className="flex flex-wrap items-center gap-2">
                      <Tag
                        color={row.kind === 'ANNOUNCEMENT' ? 'warning' : 'processing'}
                        className="m-0"
                      >
                        {kindLabel(row)}
                      </Tag>
                      {severityLabel(row) ? <Tag className="m-0">{severityLabel(row)}</Tag> : null}
                      <Typography.Text
                        strong={!row.is_read}
                        type={row.is_read ? 'secondary' : undefined}
                      >
                        {row.title}
                      </Typography.Text>
                      <Typography.Text type="secondary" className="text-xs whitespace-nowrap">
                        {formatDateTime(row.publish_at || row.created_at)}
                      </Typography.Text>
                    </div>
                  }
                  description={
                    excerpt(row) ? (
                      <Typography.Text type="secondary" className="text-[13px]">
                        {excerpt(row)}
                      </Typography.Text>
                    ) : null
                  }
                />
              </List.Item>
            )}
          />
        )}
      </Spin>

      {total > 0 ? (
        <div className="flex justify-end">
          <Pagination
            current={page}
            pageSize={pageSize}
            total={total}
            showSizeChanger
            pageSizeOptions={[10, 20, 30]}
            onChange={(p, size) => {
              setPage(p)
              setPageSize(size)
            }}
          />
        </div>
      ) : null}

      <MessageDetailModal
        open={detailOpen}
        source={detailSource}
        onClose={() => {
          setDetailOpen(false)
          setDetailSource(null)
        }}
        onChanged={handleDetailChanged}
      />
    </Space>
  )
}
