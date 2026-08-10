/** Author: Charlie */

import { MdPreview } from '@/components/editor/MdPreview'

type Props = {
  content: string
  className?: string
}

export function Markdown({ content, className }: Props) {
  return <MdPreview value={content} className={className} />
}
