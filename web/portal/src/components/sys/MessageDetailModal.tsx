/** Author: Charlie */

import { useEffect, useMemo, useState } from 'react'
import { Button, Modal, Space, Spin, Tag, Typography } from 'antd'
import { myNoticeApi } from '@/api'
import { useMessageUnreadStore } from '@/stores/messageUnread'
import { displayValue, formatDateTime, wireBool } from '@/utils'
import { dictTypeData } from '@/utils/dict'

type MessageKind = 'NOTIFICATION' | 'ANNOUNCEMENT'

export type MessageDetailSource = {
  id: string
  sourceType?: MessageKind | string
  type?: string
  title?: string
  content?: string
  publish_at?: string
  created_at?: string
  severity?: string
  is_read?: boolean | string
  kind?: string
}

type Props = {
  open: boolean
  source: MessageDetailSource | null
  onClose: () => void
  onChanged?: (payload: { type: string; id: string }) => void
}

function resolveKind(raw: unknown): MessageKind {
  return String(raw || 'NOTIFICATION').toUpperCase() === 'ANNOUNCEMENT'
    ? 'ANNOUNCEMENT'
    : 'NOTIFICATION'
}

function asReadFlag(value: unknown): boolean {
  if (typeof value === 'boolean' || typeof value === 'string') {
    return wireBool(value)
  }
  return false
}

function MessageDetailBody({
  source,
  onChanged,
}: {
  source: MessageDetailSource
  onChanged?: (payload: { type: string; id: string }) => void
}) {
  const notifyRead = useMessageUnreadStore((s) => s.notifyRead)
  const refreshUnread = useMessageUnreadStore((s) => s.refresh)
  const [loading, setLoading] = useState(true)
  const [actionLoading, setActionLoading] = useState(false)
  const [detail, setDetail] = useState<any>({})
  const [readLocally, setReadLocally] = useState(false)

  useEffect(() => {
    let mounted = true
    const wasUnread = !asReadFlag(source.is_read)
    const sourceId = source.id
    const sourceKind = source.sourceType

    void (async () => {
      try {
        // myDetail 后端会顺带标记已读
        const response = await myNoticeApi.myDetail(sourceId)
        if (!mounted) return
        const data = response.data ?? {}
        setDetail(data)

        const id = data.id || sourceId
        const kind = resolveKind(data.kind || sourceKind)
        if (id && wasUnread) {
          setDetail((prev: any) => ({ ...prev, is_read: true }))
          setReadLocally(true)
          notifyRead()
          onChanged?.({ type: kind, id })
          void refreshUnread()
        }
      } finally {
        if (mounted) setLoading(false)
      }
    })()

    return () => {
      mounted = false
    }
    // MessageDetailBody 通过 key=source.id 重挂载
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const messageKind: MessageKind = useMemo(
    () => resolveKind(detail.kind || source.sourceType || source.type),
    [detail.kind, source],
  )

  const kindLabel = messageKind === 'ANNOUNCEMENT' ? '公告' : '通知'
  const titleText = displayValue(detail.title || source.title)
  const contentText = displayValue(detail.content || source.content)
  const publishText = formatDateTime(detail.publish_at || source.publish_at || detail.created_at)
  const severity = detail.severity || source.severity
  const severityLabel = severity ? dictTypeData('NOTIFICATION_SEVERITY', severity) || severity : ''
  const isRead = readLocally || asReadFlag(detail.is_read) || asReadFlag(source.is_read)

  async function markRead() {
    const id = detail.id || source.id
    if (!id || isRead) return
    const kind = resolveKind(detail.kind || messageKind)
    setActionLoading(true)
    try {
      await myNoticeApi.read({ ids: [id] })
      setDetail((prev: any) => ({ ...prev, is_read: true }))
      setReadLocally(true)
      notifyRead()
      onChanged?.({ type: kind, id })
      void refreshUnread()
    } finally {
      setActionLoading(false)
    }
  }

  return (
    <>
      <div className="mb-4">
        <Space size={8}>
          <Tag color={messageKind === 'ANNOUNCEMENT' ? 'warning' : 'processing'}>{kindLabel}</Tag>
          <Tag color={isRead ? 'success' : 'warning'}>{isRead ? '已读' : '未读'}</Tag>
        </Space>
      </div>
      <Spin spinning={loading}>
        <Space direction="vertical" size={16} className="w-full">
          <div>
            <Typography.Title level={4} style={{ margin: 0, fontSize: 18 }}>
              {titleText}
            </Typography.Title>
            <Space size={8} className="mt-2" wrap>
              {severityLabel ? <Tag>{severityLabel}</Tag> : null}
              <Typography.Text type="secondary" className="text-xs">
                {publishText}
              </Typography.Text>
            </Space>
          </div>
          <div className="border-t border-[var(--ant-color-border)] pt-4">
            <div
              className="max-h-[min(420px,calc(100vh-280px))] overflow-y-auto text-sm leading-7 whitespace-pre-wrap"
              style={{ color: 'var(--ant-color-text-secondary)' }}
            >
              {contentText}
            </div>
          </div>
          {!isRead ? (
            <div className="flex justify-end">
              <Button type="link" loading={actionLoading} onClick={() => void markRead()}>
                标记已读
              </Button>
            </div>
          ) : null}
        </Space>
      </Spin>
    </>
  )
}

export function MessageDetailModal({ open, source, onClose, onChanged }: Props) {
  return (
    <Modal
      open={open}
      onCancel={onClose}
      footer={null}
      width={560}
      destroyOnHidden
      title="消息详情"
    >
      {open && source?.id ? (
        <MessageDetailBody key={source.id} source={source} onChanged={onChanged} />
      ) : null}
    </Modal>
  )
}
