/** Author: Charlie */

import { forwardRef, useEffect, useImperativeHandle, useRef, useState } from 'react'
import { Input, Spin } from 'antd'
import { authApi } from '@/api'

export type CaptchaInputHandle = {
  refresh: () => Promise<void>
}

type Props = {
  captchaId: string
  captchaValue: string
  onCaptchaIdChange: (value: string) => void
  onCaptchaValueChange: (value: string) => void
  size?: 'middle' | 'large'
}

export const CaptchaInput = forwardRef<CaptchaInputHandle, Props>(function CaptchaInput(
  { captchaValue, onCaptchaIdChange, onCaptchaValueChange, size = 'middle' },
  ref,
) {
  const [loading, setLoading] = useState(false)
  const [imageBase64, setImageBase64] = useState('')
  const idChangeRef = useRef(onCaptchaIdChange)
  const valueChangeRef = useRef(onCaptchaValueChange)
  idChangeRef.current = onCaptchaIdChange
  valueChangeRef.current = onCaptchaValueChange

  async function refresh() {
    setLoading(true)
    try {
      const response = await authApi.captcha('svg')
      idChangeRef.current(response.data.captcha_id)
      valueChangeRef.current('')
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
  const imageHeight = size === 'large' ? 'h-10' : 'h-8'

  return (
    <div className="grid grid-cols-[minmax(0,1fr)_140px] gap-2.5 items-center">
      <Input
        size={size}
        value={captchaValue}
        placeholder="请输入验证码"
        allowClear
        onChange={(e) => onCaptchaValueChange(e.target.value)}
      />
      <button
        type="button"
        className={`${imageHeight} w-140px overflow-hidden rounded-md bg-[var(--ant-color-fill-quaternary)] p-0 cursor-pointer disabled:cursor-wait`}
        disabled={loading}
        onClick={() => void refresh()}
      >
        <Spin spinning={loading}>
          {imageSrc ? (
            <img src={imageSrc} alt="验证码" className={`block ${imageHeight} w-140px`} />
          ) : null}
        </Spin>
      </button>
    </div>
  )
})
