/** Author: Charlie */

import { useRef, useState } from 'react'
import { Button, ConfigProvider, Form, Input, message } from 'antd'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { authApi } from '@/api'
import { CaptchaInput, type CaptchaInputHandle } from '@/components/common/CaptchaInput'
import { PasswordStrength } from '@/components/common/PasswordStrength'
import { useAuthModalStore } from '@/stores/authModal'
import { encryptPasswords } from '@/utils/security'
import { isValidEmail } from '@/utils/validate'
import { AuthCenter } from './AuthSplit'

type FormValues = {
  email: string
  password?: string
  confirmPassword?: string
  captcha_id: string
  captcha_value: string
}

export function ForgotPasswordPage() {
  const [form] = Form.useForm<FormValues>()
  const [loading, setLoading] = useState(false)
  const captchaRef = useRef<CaptchaInputHandle>(null)
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const openAuthModal = useAuthModalStore((s) => s.open)
  const token = params.get('token') || ''
  const isResetMode = Boolean(token)
  const captchaId = Form.useWatch('captcha_id', form) || ''
  const captchaValue = Form.useWatch('captcha_value', form) || ''
  const password = Form.useWatch('password', form) || ''

  async function sendLink() {
    try {
      await form.validateFields(['email', 'captcha_value'])
    } catch {
      return
    }
    const values = form.getFieldsValue()
    if (!isValidEmail(values.email || '')) {
      message.warning('请输入有效邮箱')
      return
    }
    setLoading(true)
    try {
      await authApi.forgotPassword({
        email: values.email.trim(),
        captcha_id: values.captcha_id,
        captcha_value: values.captcha_value,
      })
      message.success('密码重置链接已发送')
      await captchaRef.current?.refresh()
    } catch {
      await captchaRef.current?.refresh()
    } finally {
      setLoading(false)
    }
  }

  async function resetPassword() {
    try {
      await form.validateFields(['password', 'confirmPassword', 'captcha_value'])
    } catch {
      return
    }
    const values = form.getFieldsValue()

    setLoading(true)
    try {
      const encrypted = await encryptPasswords({ password: values.password || '' })
      await authApi.resetPassword({
        token,
        password: encrypted.values.password || '',
        password_key_id: encrypted.password_key_id,
        captcha_id: values.captcha_id,
        captcha_value: values.captcha_value,
      })
      message.success('密码已重置，请重新登录')
      navigate('/')
      openAuthModal('login')
    } catch {
      await captchaRef.current?.refresh()
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthCenter
      title={isResetMode ? '重置密码' : '找回密码'}
      description={
        isResetMode
          ? '请设置新密码。重置链接在过期前仅可使用一次。'
          : '请输入已启用登录的邮箱，系统将发送密码重置链接。'
      }
    >
      <ConfigProvider componentSize="large">
        <Form
          form={form}
          layout="vertical"
          requiredMark={false}
          initialValues={{ captcha_id: '', captcha_value: '' }}
          onFinish={() => {
            void (isResetMode ? resetPassword() : sendLink())
          }}
        >
          {!isResetMode ? (
            <Form.Item name="email" rules={[{ required: true, message: '请输入登录邮箱' }]}>
              <Input placeholder="登录邮箱" allowClear />
            </Form.Item>
          ) : (
            <>
              <Form.Item name="password" rules={[{ required: true, message: '请输入新密码' }]}>
                <Input.Password placeholder="新密码（至少 8 位，含大小写、数字与特殊字符）" />
              </Form.Item>
              <PasswordStrength password={password} />
              <Form.Item
                name="confirmPassword"
                dependencies={['password']}
                rules={[
                  { required: true, message: '请确认密码' },
                  ({ getFieldValue }) => ({
                    validator(_, value) {
                      if (!value || getFieldValue('password') === value) {
                        return Promise.resolve()
                      }
                      return Promise.reject(new Error('两次输入的密码不一致'))
                    },
                  }),
                ]}
              >
                <Input.Password placeholder="确认新密码" />
              </Form.Item>
            </>
          )}

          <Form.Item name="captcha_id" hidden>
            <Input />
          </Form.Item>

          <Form.Item name="captcha_value" rules={[{ required: true, message: '请输入验证码' }]}>
            <CaptchaInput
              ref={captchaRef}
              size="large"
              captchaId={captchaId}
              captchaValue={captchaValue}
              onCaptchaIdChange={(v) => form.setFieldValue('captcha_id', v)}
              onCaptchaValueChange={(v) => form.setFieldValue('captcha_value', v)}
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading}>
              {isResetMode ? '重置密码' : '发送重置链接'}
            </Button>
          </Form.Item>

          <div className="auth-center__links">
            <Button
              type="link"
              className="!px-0"
              onClick={() => {
                navigate('/')
                openAuthModal('login')
              }}
            >
              返回登录
            </Button>
            {isResetMode ? (
              <Link to="/auth/forgot-password">重新申请链接</Link>
            ) : (
              <Button
                type="link"
                className="!px-0"
                onClick={() => {
                  navigate('/')
                  openAuthModal('register')
                }}
              >
                去注册
              </Button>
            )}
          </div>
        </Form>
      </ConfigProvider>
    </AuthCenter>
  )
}
