/** Author: Charlie */

import { useEffect, useMemo, useState } from 'react'
import { Button, Empty, Image, Skeleton, Typography } from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import { Link, useParams } from 'react-router-dom'
import { feedbackApi } from '@/api'
import { DictTag } from '@/components/common/DictTag'
import { displayValue, formatDateTime, formatFileSize, isImageFile } from '@/utils'

export function FeedbackDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [loading, setLoading] = useState(Boolean(id))
  const [detail, setDetail] = useState<any>(null)
  const [notFound, setNotFound] = useState(!id)
  const [activeId, setActiveId] = useState(id)

  if (id !== activeId) {
    setActiveId(id)
    setLoading(Boolean(id))
    setNotFound(!id)
    setDetail(null)
  }

  useEffect(() => {
    if (!id) return

    let mounted = true

    void (async () => {
      try {
        const res = await feedbackApi.myDetail(id)
        if (!mounted) return
        setDetail(res.data ?? null)
        setNotFound(!res.data)
      } catch {
        if (!mounted) return
        setDetail(null)
        setNotFound(true)
      } finally {
        if (mounted) setLoading(false)
      }
    })()

    return () => {
      mounted = false
    }
  }, [id])

  const attachments = useMemo(() => {
    const list = detail?.attachments
    return Array.isArray(list) ? list : []
  }, [detail])

  if (loading) {
    return (
      <div className="page-shell">
        <Skeleton active paragraph={{ rows: 8 }} />
      </div>
    )
  }

  if (notFound || !detail) {
    return (
      <div className="page-shell">
        <Empty description="反馈不存在或无权查看" />
        <div className="mt-4 text-center">
          <Link to="/feedback">
            <Button type="primary">返回反馈列表</Button>
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="page-shell">
      <div className="mb-4">
        <Link to="/feedback">
          <Button type="text" icon={<ArrowLeftOutlined />} className="!px-0">
            返回反馈列表
          </Button>
        </Link>
      </div>

      <div className="rounded-lg border border-[var(--ant-color-border)] bg-[var(--ant-color-bg-container)] p-5">
        <div className="mb-3 flex flex-wrap items-center gap-2">
          <Typography.Title level={3} className="!mb-0 !text-xl">
            {displayValue(detail.title)}
          </Typography.Title>
          <DictTag dictCode="FEEDBACK_STATUS" value={detail.status} />
          <DictTag dictCode="FEEDBACK_TYPE" value={detail.type} />
        </div>
        <Typography.Paragraph type="secondary" className="!mb-4 text-xs">
          提交于 {formatDateTime(detail.created_at)}
        </Typography.Paragraph>
        <Typography.Paragraph className="whitespace-pre-wrap">
          {displayValue(detail.content)}
        </Typography.Paragraph>

        {attachments.length ? (
          <div className="mt-4">
            <Typography.Text strong>附件</Typography.Text>
            <div className="mt-2 flex flex-wrap gap-3">
              {attachments.map((file: any) => {
                const url = file.url || file.file_url
                const name = file.name || file.filename || '附件'
                if (isImageFile(name) || isImageFile(url)) {
                  return (
                    <Image
                      key={file.id || url || name}
                      src={url}
                      alt={name}
                      width={96}
                      height={96}
                      className="rounded object-cover"
                    />
                  )
                }
                return (
                  <a
                    key={file.id || url || name}
                    href={url}
                    target="_blank"
                    rel="noreferrer"
                    className="text-sm"
                  >
                    {name}
                    {file.size ? ` (${formatFileSize(file.size)})` : ''}
                  </a>
                )
              })}
            </div>
          </div>
        ) : null}

        {detail.reply_content ? (
          <div className="mt-6 rounded border border-[var(--ant-color-border-secondary)] bg-[var(--ant-color-fill-quaternary)] p-4">
            <Typography.Text strong>官方回复</Typography.Text>
            <Typography.Paragraph className="mt-2 mb-0 whitespace-pre-wrap">
              {displayValue(detail.reply_content)}
            </Typography.Paragraph>
            {detail.replied_at ? (
              <Typography.Text type="secondary" className="text-xs">
                回复于 {formatDateTime(detail.replied_at)}
              </Typography.Text>
            ) : null}
          </div>
        ) : null}
      </div>
    </div>
  )
}
