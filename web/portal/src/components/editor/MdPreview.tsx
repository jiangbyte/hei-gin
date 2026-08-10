/** Author: Charlie */

import { MdPreview as MdPreviewBase } from 'md-editor-rt'
import { useAppStore } from '@/stores/app'
import 'md-editor-rt/lib/preview.css'
import './md-preview.css'

type PreviewTheme = 'light' | 'dark'

type Props = {
  value?: string | null
  className?: string
  theme?: PreviewTheme
  previewTheme?: string
  codeTheme?: string
  showCodeRowNumber?: boolean
}

export function MdPreview({
  value = '',
  className,
  theme,
  previewTheme = 'github',
  codeTheme = 'atom',
  showCodeRowNumber = true,
}: Props) {
  const resolvedTheme = useAppStore((s) => s.resolvedTheme)
  const previewThemeMode = theme ?? resolvedTheme

  return (
    <div className={`md-preview ${className ?? ''}`}>
      <MdPreviewBase
        value={value ?? ''}
        theme={previewThemeMode}
        previewTheme={previewTheme}
        codeTheme={codeTheme}
        showCodeRowNumber={showCodeRowNumber}
        style={{ padding: 0 }}
      />
    </div>
  )
}
