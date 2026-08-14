/** Author: Charlie */

import { useCallback, useEffect, useMemo, useState } from 'react'
import { Avatar, Button, List, Popconfirm, Space, Spin, Typography, message } from 'antd'
import { authApi } from '@/api'
import { wireBool } from '@/utils/wire'
import '../usercenter.css'

type Binding = {
  provider: string
  label: string
  open_id_masked?: string
  nickname?: string
  avatar?: string
  bound_at?: string
}

type ProviderOption = {
  provider: string
  label: string
  enabled: boolean
  web_oauth: boolean
}

export function OauthPanel() {
  const [loading, setLoading] = useState(true)
  const [bindingProvider, setBindingProvider] = useState<string | null>(null)
  const [bindings, setBindings] = useState<Binding[]>([])
  const [providers, setProviders] = useState<ProviderOption[]>([])

  const refresh = useCallback(async () => {
    setLoading(true)
    try {
      const [bindRes, optRes] = await Promise.all([authApi.oauthBindings(), authApi.authOptions()])
      setBindings(Array.isArray(bindRes?.data) ? bindRes.data : [])
      const list = Array.isArray(optRes?.data?.oauth_providers) ? optRes.data.oauth_providers : []
      setProviders(
        list
          .map((item: any) => ({
            provider: String(item.provider || ''),
            label: String(item.label || item.provider || ''),
            enabled: wireBool(item.enabled ?? false),
            web_oauth: wireBool(item.web_oauth ?? true),
          }))
          .filter((item: ProviderOption) => item.provider && item.enabled && item.web_oauth),
      )
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void refresh()
  }, [refresh])

  const boundSet = useMemo(() => new Set(bindings.map((item) => item.provider)), [bindings])

  async function bind(provider: string) {
    if (bindingProvider) return
    setBindingProvider(provider)
    try {
      const res = await authApi.oauthBindAuthorize(provider)
      const url = res?.data?.authorize_url
      if (!url) {
        message.error('无法发起绑定')
        return
      }
      window.location.href = String(url)
    } catch {
      // 全局错误提示
    } finally {
      setBindingProvider(null)
    }
  }

  async function unbind(provider: string) {
    await authApi.oauthUnbind(provider)
    message.success('已解绑')
    await refresh()
  }

  return (
    <Spin spinning={loading}>
      <div className="uc-panel">
        <Typography.Paragraph type="secondary">
          绑定后可使用对应平台快速登录。至少保留一种登录方式（密码/邮箱/手机/其它三方）。
        </Typography.Paragraph>

        <List
          locale={{ emptyText: '暂无已绑定三方账号' }}
          dataSource={bindings}
          renderItem={(item) => (
            <List.Item
              actions={[
                <Popconfirm
                  key="unbind"
                  title="确认解绑该三方账号？"
                  onConfirm={() => void unbind(item.provider)}
                >
                  <Button type="link" danger>
                    解绑
                  </Button>
                </Popconfirm>,
              ]}
            >
              <List.Item.Meta
                avatar={<Avatar src={item.avatar || undefined}>{item.label?.[0]}</Avatar>}
                title={`${item.label || item.provider}${item.nickname ? ` · ${item.nickname}` : ''}`}
                description={`OpenID：${item.open_id_masked || '-'} · 绑定于 ${item.bound_at || '-'}`}
              />
            </List.Item>
          )}
        />

        <Typography.Title level={5} style={{ marginTop: 24 }}>
          可绑定平台
        </Typography.Title>
        <Space wrap>
          {providers.map((item) => {
            const bound = boundSet.has(item.provider)
            return (
              <Button
                key={item.provider}
                disabled={bound || Boolean(bindingProvider)}
                loading={bindingProvider === item.provider}
                onClick={() => void bind(item.provider)}
              >
                {bound ? `已绑定 ${item.label}` : `绑定 ${item.label}`}
              </Button>
            )
          })}
          {!providers.length ? (
            <Typography.Text type="secondary">暂无可绑定的三方登录（未开启或未配置）</Typography.Text>
          ) : null}
        </Space>
      </div>
    </Spin>
  )
}
