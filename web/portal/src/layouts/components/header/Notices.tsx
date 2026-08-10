/** Author: Charlie */

import { useCallback, useEffect, useMemo, useState } from 'react'
import { Badge, Button, Divider, Popover, Tooltip } from 'antd'
import { BellOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { messageApi } from '@/api'
import { formatDateTime, wireBool } from '@/utils'
import { readPageMeta } from '@/utils/wire'
import { useMessageUnreadStore } from '@/stores/messageUnread'
import {
  MessageDetailModal,
  type MessageDetailSource,
} from '@/components/message/MessageDetailModal'
import { NoticeList, type BannerItem } from '../common/NoticeList'

const pageSize = 8

type LoadMode = 'replace' | 'merge' | 'append'
type NoticeKind = 'NOTIFICATION' | 'ANNOUNCEMENT'

interface NoticeSource {
  id: string
  title: string
  icon: 'announcement' | 'notification'
  tagTitle?: string
  tagType?: BannerItem['tagType']
  description?: string
  date: string
  sourceType: NoticeKind
  sourceId: string
  isRead: boolean
}

function mapHistoryItem(item: any): NoticeSource {
  const kind: NoticeKind = item.kind === 'ANNOUNCEMENT' ? 'ANNOUNCEMENT' : 'NOTIFICATION'
  return {
    id: `${kind}:${item.id}`,
    title: item.title,
    icon: kind === 'ANNOUNCEMENT' ? 'announcement' : 'notification',
    tagTitle: kind === 'ANNOUNCEMENT' ? '公告' : '通知',
    tagType: kind === 'ANNOUNCEMENT' ? 'warning' : 'processing',
    description: item.content,
    date: formatDateTime(item.publish_at || item.created_at),
    sourceType: kind,
    sourceId: item.id,
    isRead: wireBool(item.is_read ?? false),
  }
}

function mergeNoticeRecords(
  current: NoticeSource[],
  incoming: NoticeSource[],
  mode: LoadMode,
): NoticeSource[] {
  if (mode === 'replace') return incoming
  const currentMap = new Map(current.map((item) => [item.id, item]))
  const result = current.map((item) => ({
    ...item,
    ...(incoming.find((i) => i.id === item.id) ?? {}),
  }))
  incoming.forEach((item) => {
    if (!currentMap.has(item.id)) result.push(item)
  })
  return result
}

export function Notices() {
  const navigate = useNavigate()
  const unreadTotal = useMessageUnreadStore((s) => s.unreadTotal)
  const refreshUnreadCount = useMessageUnreadStore((s) => s.refresh)
  const notifyReadAll = useMessageUnreadStore((s) => s.notifyReadAll)
  const [open, setOpen] = useState(false)
  const [records, setRecords] = useState<NoticeSource[]>([])
  const [current, setCurrent] = useState(0)
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [loaded, setLoaded] = useState(false)
  const [detailOpen, setDetailOpen] = useState(false)
  const [detailSource, setDetailSource] = useState<MessageDetailSource | null>(null)

  const list = useMemo(() => records.map((item) => ({ ...item }) as BannerItem), [records])
  const hasMore = records.length < total

  const loadList = useCallback(async (page = 1, mode: LoadMode = 'replace') => {
    setLoading(true)
    try {
      const response = await messageApi.myPage({ current: page, size: pageSize })
      const data = response.data ?? {}
      const incoming = (data.records ?? []).map((item: any) => mapHistoryItem(item))
      setRecords((prev) => mergeNoticeRecords(prev, incoming, mode))
      const pageMeta = readPageMeta(data, { current: page, size: pageSize })
      setTotal(pageMeta.total || incoming.length)
      setCurrent(pageMeta.current)
      setLoaded(true)
    } finally {
      setLoading(false)
    }
  }, [])

  const refresh = useCallback(async () => {
    await Promise.all([refreshUnreadCount(), loadList(1, loaded ? 'merge' : 'replace')])
  }, [loadList, loaded, refreshUnreadCount])

  useEffect(() => {
    void refreshUnreadCount()
  }, [refreshUnreadCount])

  async function loadMore() {
    if (loading || records.length >= total) return
    await loadList(current + 1, 'append')
  }

  function handleOpen(id: string) {
    const item = records.find((notice) => notice.id === id)
    if (!item) return
    setDetailSource({
      id: item.sourceId,
      sourceType: item.sourceType,
      title: item.title,
      content: item.description,
      is_read: item.isRead,
    })
    setDetailOpen(true)
  }

  function handleDetailChanged(payload: { type: string; id: string }) {
    setRecords((prev) =>
      prev.map((notice) =>
        notice.id === `${payload.type}:${payload.id}` ? { ...notice, isRead: true } : notice,
      ),
    )
  }

  async function markAllRead() {
    try {
      await messageApi.readAll()
      setRecords((prev) => prev.map((item) => ({ ...item, isRead: true })))
      notifyReadAll()
    } catch {
      /* ignore */
    }
  }

  function goMore() {
    setOpen(false)
    void navigate('/usercenter?tab=my_messages')
  }

  const content = (
    <div style={{ width: 390 }}>
      <div style={{ padding: '12px 12px 0', fontWeight: 600 }}>我的消息</div>
      <NoticeList
        list={list}
        loading={loading}
        hasMore={hasMore}
        onOpen={handleOpen}
        onLoadMore={() => void loadMore()}
      />
      <Divider style={{ margin: 0 }} />
      <div className="flex gap-2 p-2">
        <Button
          className="flex-1"
          size="small"
          type="text"
          disabled={unreadTotal <= 0}
          onClick={() => void markAllRead()}
        >
          全部已读
        </Button>
        <Button className="flex-1" size="small" type="text" onClick={goMore}>
          查看更多
        </Button>
      </div>
    </div>
  )

  return (
    <>
      <Popover
        content={content}
        trigger="click"
        open={open}
        onOpenChange={(next) => {
          setOpen(next)
          if (next) void refresh()
        }}
        placement="bottomRight"
        arrow={{ pointAtCenter: true }}
      >
        <Tooltip title="消息" placement="bottom">
          <button type="button" className="header-icon-btn" aria-label="消息">
            <Badge count={unreadTotal} overflowCount={99} size="small">
              <BellOutlined style={{ fontSize: 18 }} />
            </Badge>
          </button>
        </Tooltip>
      </Popover>

      <MessageDetailModal
        open={detailOpen}
        source={detailSource}
        onClose={() => {
          setDetailOpen(false)
          setDetailSource(null)
        }}
        onChanged={handleDetailChanged}
      />
    </>
  )
}
