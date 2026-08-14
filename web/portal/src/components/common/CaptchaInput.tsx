/** Author: Charlie */

import { forwardRef, useEffect, useImperativeHandle, useRef, useState } from 'react'
import { Input, Spin } from 'antd'
import type { InputProps } from 'antd'
import { authApi } from '@/api'
import './captcha-input.css'

export type CaptchaInputHandle = {
  refresh: () => Promise<void>
}

type Props = {
  /** Form.Item 注入 */
  value?: string
  onChange?: (value: string) => void
  status?: InputProps['status']
  onCaptchaIdChange?: (value: string) => void
  size?: 'middle' | 'large'
}

export const CaptchaInput = forwardRef<CaptchaInputHandle, Props>(function CaptchaInput(
  { value = '', onChange, status, onCaptchaIdChange, size = 'middle' },
  ref,
) {
  const [loading, setLoading] = useState(false)
  const [imageBase64, setImageBase64] = useState('')
  const idChangeRef = useRef(onCaptchaIdChange)
  const valueChangeRef = useRef(onChange)
  idChangeRef.current = onCaptchaIdChange
  valueChangeRef.current = onChange

  async function refresh() {
    setLoading(true)
    try {
      const response = await authApi.captcha('svg')
      idChangeRef.current?.(response.data.captcha_id)
      valueChangeRef.current?.('')
      setImageBase64(response.data.image_base64)
    } finally {
      setLoading(false)
    }
  }

  useImperativeHandle(ref, () => ({ refresh }))

  useEffect(() => {
    void refresh()
  }, [])

  const imageSrc = imageBase64 ? `data:image/svg+xml;base64,${imageBase64}` : ''

  return (
    <div className={`captcha-input${size === 'large' ? ' captcha-input--large' : ''}`}>
      <Input
        size={size}
        status={status}
        value={value}
        placeholder="请输入验证码"
        allowClear
        onChange={(e) => onChange?.(e.target.value)}
      />
      <button
        type="button"
        className="captcha-input__image"
        disabled={loading}
        onClick={() => void refresh()}
        aria-label="刷新验证码"
      >
        <Spin spinning={loading}>
          {imageSrc ? <img src={imageSrc} alt="验证码" /> : null}
        </Spin>
      </button>
    </div>
  )
})
