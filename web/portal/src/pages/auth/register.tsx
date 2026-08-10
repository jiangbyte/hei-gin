/** Author: Charlie */

import { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useAuthModalStore } from '@/stores/authModal'

/** 兼容深链：打开注册弹窗并回到首页作为背景。 */
export function RegisterPage() {
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const open = useAuthModalStore((s) => s.open)

  useEffect(() => {
    open('register', params.get('redirect'))
    navigate('/', { replace: true })
  }, [navigate, open, params])

  return null
}
