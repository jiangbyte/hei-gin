/** Author: Charlie */

import { Typography } from 'antd'
import { useNavigate } from 'react-router-dom'

export function Logo() {
  const navigate = useNavigate()

  return (
    <Typography.Title level={4} className="!m-0 cursor-pointer" onClick={() => navigate('/')}>
      {import.meta.env.VITE_APP_TITLE || 'HEI'}
    </Typography.Title>
  )
}
