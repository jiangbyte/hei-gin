/** Author: Charlie */

import { Tag } from 'antd'
import { useDict } from '@/hooks/useDict'
import { dictTypeColor, dictTypeData } from '@/utils/dict'

type Props = {
  dictCode: string
  value?: string | number | null
}

/** 对齐 admin：展示 dictTypeData + dictTypeColor。 */
export function DictTag({ dictCode, value }: Props) {
  useDict()
  const label = dictTypeData(dictCode, value)
  if (!label) {
    return <span>{value == null || value === '' ? '-' : String(value)}</span>
  }
  const color = dictTypeColor(dictCode, value)
  return <Tag color={color || undefined}>{label}</Tag>
}
