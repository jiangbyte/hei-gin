/** Author: Charlie */

import { useState } from 'react'
import { Alert, Button, Form, Input, Modal, message } from 'antd'
import { useNavigate } from 'react-router-dom'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores/auth'
import '../usercenter.css'

export function CancelAccountPanel() {
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const resetSession = useAuthStore((s) => s.resetSession)
  const [submitting, setSubmitting] = useState(false)
  const confirmText = Form.useWatch('confirmText', form) || ''

  async function submitCancel() {
    const values = await form.validateFields()
    if (String(values.confirmText || '').trim() !== '注销') {
      message.warning('请输入「注销」以确认')
      return
    }
    Modal.confirm({
      title: '确认注销账号',
      content:
        '注销后将立即失效登录、释放手机号/邮箱等登录标识，并清理资料与授权。账号进入注销保留期；期间未再登录使用的，到期后将彻底删除并按绑定邮箱/短信通知。此操作不可撤销。',
      okText: '确认注销',
      okType: 'danger',
      cancelText: '再想想',
      onOk: () => doCancel(values.cancel_reason),
    })
  }

  async function doCancel(cancelReason?: string) {
    setSubmitting(true)
    try {
      await authApi.cancelAccount({
        cancel_reason: String(cancelReason || '').trim() || null,
      })
      message.success('账号已注销')
      resetSession()
      navigate('/auth/login', { replace: true })
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="user-center-form user-center-form--narrow w-full min-w-0">
      <Alert
        type="warning"
        showIcon
        className="mb-4"
        message="注销为不可逆操作：将清理个人资料、登录标识、角色/部门/用户组与资源授权，并强制下线全部会话。系统按保留天数（默认 15 天）暂存账号主记录；到期且未再登录使用后彻底删除，并通过已绑定邮箱通知（短信需在系统配置中填写模板编号）。"
      />
      <Form form={form} layout="vertical">
        <Form.Item name="cancel_reason" label="注销原因（可选）">
          <Input.TextArea rows={3} maxLength={500} showCount placeholder="例如：不再使用本系统" />
        </Form.Item>
        <Form.Item
          name="confirmText"
          label="请输入「注销」确认"
          rules={[{ required: true, message: '请输入「注销」' }]}
        >
          <Input placeholder="注销" />
        </Form.Item>
        <Form.Item>
          <Button
            danger
            type="primary"
            loading={submitting}
            disabled={String(confirmText).trim() !== '注销'}
            onClick={() => void submitCancel()}
          >
            注销账号
          </Button>
        </Form.Item>
      </Form>
    </div>
  )
}
