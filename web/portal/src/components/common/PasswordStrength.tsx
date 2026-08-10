/** Author: Charlie */

import { useMemo } from 'react'
import { theme } from 'antd'

type Props = {
  password: string
}

interface StrengthLevel {
  label: string
  color: string
  percent: number
}

const policyItems = [
  { label: '至少 8 个字符', met: (pwd: string) => pwd.length >= 8 },
  { label: '包含大写字母', met: (pwd: string) => /[A-Z]/.test(pwd) },
  { label: '包含小写字母', met: (pwd: string) => /[a-z]/.test(pwd) },
  { label: '包含数字', met: (pwd: string) => /[0-9]/.test(pwd) },
  { label: '包含特殊字符', met: (pwd: string) => /[^A-Za-z0-9]/.test(pwd) },
]

export function PasswordStrength({ password }: Props) {
  const { token } = theme.useToken()
  const pwd = password ?? ''

  const levels: StrengthLevel[] = useMemo(
    () => [
      { label: '弱', color: token.colorError, percent: 25 },
      { label: '较弱', color: token.colorWarning, percent: 50 },
      { label: '中等', color: token.colorWarningHover, percent: 75 },
      { label: '强', color: token.colorSuccess, percent: 100 },
    ],
    [token.colorError, token.colorWarning, token.colorWarningHover, token.colorSuccess],
  )

  const strength = useMemo(() => {
    if (!pwd) {
      return { label: '', color: token.colorBorder, percent: 0 }
    }
    let score = 0
    if (pwd.length >= 8) score += 1
    if (pwd.length >= 12) score += 1
    if (/[a-z]/.test(pwd) && /[A-Z]/.test(pwd)) score += 1
    if (/[0-9]/.test(pwd) && /[^A-Za-z0-9]/.test(pwd)) score += 1
    return levels[Math.min(score, levels.length - 1)]
  }, [pwd, levels, token.colorBorder])

  if (!pwd) {
    return null
  }

  return (
    <div className="mt-1">
      <div className="mb-1 h-1 overflow-hidden rounded bg-[var(--ant-color-fill-quaternary)]">
        <div
          className="h-full rounded transition-all duration-300"
          style={{ width: `${strength.percent}%`, backgroundColor: strength.color }}
        />
      </div>
      <span className="mb-1 inline-block text-xs font-semibold" style={{ color: strength.color }}>
        {strength.label}
      </span>
      <div className="flex flex-wrap gap-x-3 gap-y-0.5">
        {policyItems.map((item) => {
          const met = item.met(pwd)
          return (
            <span
              key={item.label}
              className={`inline-flex items-center gap-0.5 text-[11px] ${
                met ? 'text-[var(--ant-color-text)]' : 'text-[var(--ant-color-text-quaternary)]'
              }`}
            >
              <span
                className={`text-[10px] ${
                  met
                    ? 'text-[var(--ant-color-success)]'
                    : 'text-[var(--ant-color-text-quaternary)]'
                }`}
              >
                {met ? '✓' : '○'}
              </span>
              {item.label}
            </span>
          )
        })}
      </div>
    </div>
  )
}
