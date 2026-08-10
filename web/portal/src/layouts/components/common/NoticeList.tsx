/** Author: Charlie */

import { Badge, Button, Empty, List, Spin, Tag, Typography } from 'antd'
import { NotificationOutlined, SoundOutlined } from '@ant-design/icons'

export interface BannerItem {
  avatar?: string | null
  id: string
  title: string
  icon?: 'announcement' | 'notification'
  tagTitle?: string
  tagType?: 'default' | 'error' | 'processing' | 'success' | 'warning'
  description?: string
  date: string
  isRead?: boolean
}

type Props = {
  list?: BannerItem[]
  loading?: boolean
  hasMore?: boolean
  onOpen?: (id: string) => void
  onLoadMore?: () => void
}

export function NoticeList({ list = [], loading, hasMore, onOpen, onLoadMore }: Props) {
  if (!loading && !list.length) {
    return (
      <div className="h-90 flex items-center justify-center">
        <Empty description="暂无消息" />
      </div>
    )
  }

  if (loading && !list.length) {
    return (
      <div className="h-90 flex items-center justify-center">
        <Spin size="small" />
      </div>
    )
  }

  return (
    <div className="h-90 overflow-y-auto">
      <List
        dataSource={list}
        renderItem={(item) => (
          <List.Item
            key={item.id}
            className="cursor-pointer !px-3 !py-2.5"
            onClick={() => onOpen?.(item.id)}
          >
            <div className="flex w-full min-w-0 items-start gap-3">
              <Badge dot={!item.isRead}>
                <div
                  className="flex h-8 w-8 items-center justify-center rounded-full"
                  style={{
                    background: 'var(--ant-color-fill-quaternary)',
                    color: item.isRead
                      ? 'var(--ant-color-text-tertiary)'
                      : 'var(--ant-color-primary)',
                  }}
                >
                  {item.icon === 'announcement' ? <SoundOutlined /> : <NotificationOutlined />}
                </div>
              </Badge>

              <div className="min-w-0 flex-1">
                <div className="flex min-w-0 items-center justify-between gap-2">
                  <Typography.Text
                    strong={!item.isRead}
                    type={item.isRead ? 'secondary' : undefined}
                    ellipsis
                    className="min-w-0 flex-1"
                  >
                    {item.title}
                  </Typography.Text>
                  {item.tagTitle ? (
                    <Tag
                      color={item.tagType === 'warning' ? 'warning' : 'processing'}
                      className="m-0 shrink-0"
                    >
                      {item.tagTitle}
                    </Tag>
                  ) : null}
                </div>
                {item.description ? (
                  <Typography.Text type="secondary" ellipsis className="mt-0.5 block text-xs">
                    {item.description
                      .replace(/<[^>]+>/g, ' ')
                      .replace(/\s+/g, ' ')
                      .trim()}
                  </Typography.Text>
                ) : null}
                <Typography.Text type="secondary" className="text-[11px]">
                  {item.date}
                </Typography.Text>
              </div>
            </div>
          </List.Item>
        )}
      />
      {hasMore ? (
        <div className="flex justify-center py-2">
          <Button type="link" size="small" loading={loading} onClick={() => onLoadMore?.()}>
            加载更多
          </Button>
        </div>
      ) : null}
    </div>
  )
}
