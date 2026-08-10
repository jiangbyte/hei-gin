/** Author: Charlie */

import { MdEditor as MdEditorBase } from 'md-editor-rt'
import 'md-editor-rt/lib/style.css'
import type { Themes, UploadImgEvent } from 'md-editor-rt'
import { toCssSize } from './shared'

type Props = {
  value?: string | null
  onChange?: (value: string) => void
  height?: string | number
  placeholder?: string
  theme?: Themes
  preview?: boolean
  readOnly?: boolean
  disabled?: boolean
  language?: string
  previewTheme?: string
  codeTheme?: string
  showCodeRowNumber?: boolean
  noUploadImg?: boolean
  onUploadImg?: UploadImgEvent
  className?: string
}

export function MdEditor({
  value = '',
  onChange,
  height = 360,
  placeholder = '请输入内容',
  theme = 'light',
  preview = true,
  readOnly = false,
  disabled = false,
  language = 'zh-CN',
  previewTheme = 'github',
  codeTheme = 'atom',
  showCodeRowNumber = true,
  noUploadImg = false,
  onUploadImg,
  className,
}: Props) {
  return (
    <div className={className}>
      <MdEditorBase
        value={value ?? ''}
        style={{ height: toCssSize(height) }}
        placeholder={placeholder}
        theme={theme}
        preview={preview}
        readOnly={readOnly}
        disabled={disabled}
        language={language}
        previewTheme={previewTheme}
        codeTheme={codeTheme}
        showCodeRowNumber={showCodeRowNumber}
        noUploadImg={noUploadImg}
        onUploadImg={onUploadImg}
        onChange={onChange}
      />
    </div>
  )
}
