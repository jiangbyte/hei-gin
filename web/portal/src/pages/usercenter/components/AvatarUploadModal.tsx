/** Author: Charlie */

import { useState } from 'react'
import { Avatar, Button, Modal, Slider, Space, Upload, message } from 'antd'
import { UploadOutlined, UserOutlined } from '@ant-design/icons'
import Cropper, { type Area } from 'react-easy-crop'
import { authApi } from '@/api'
import './avatar-upload.css'

type Props = {
  open: boolean
  avatar?: string | null
  onClose: () => void
  onUploaded: () => void
}

const ACCEPTED = ['image/jpeg', 'image/png', 'image/webp']
const MAX_SIZE = 2 * 1024 * 1024
const CROP_SIZE = 320

export function AvatarUploadModal({ open, avatar, onClose, onUploaded }: Props) {
  const [source, setSource] = useState('')
  const [fileName, setFileName] = useState('')
  const [crop, setCrop] = useState({ x: 0, y: 0 })
  const [zoom, setZoom] = useState(1)
  const [croppedAreaPixels, setCroppedAreaPixels] = useState<Area | null>(null)
  const [previewUrl, setPreviewUrl] = useState('')
  const [uploading, setUploading] = useState(false)

  function resetSource() {
    if (source) {
      URL.revokeObjectURL(source)
    }
    setSource('')
    setFileName('')
    setPreviewUrl('')
    setCrop({ x: 0, y: 0 })
    setZoom(1)
    setCroppedAreaPixels(null)
  }

  async function onFileChange(file: File) {
    if (!ACCEPTED.includes(file.type)) {
      message.warning('仅支持 JPG、PNG 和 WebP 图片')
      return
    }
    if (file.size > MAX_SIZE) {
      message.warning('图片大小不能超过 2MB')
      return
    }
    if (source) {
      URL.revokeObjectURL(source)
    }
    setFileName(file.name)
    setSource(URL.createObjectURL(file))
    setPreviewUrl('')
  }

  function onCropComplete(_area: Area, areaPixels: Area) {
    setCroppedAreaPixels(areaPixels)
    void updatePreview(areaPixels)
  }

  async function updatePreview(areaPixels: Area) {
    if (!source) return
    const blob = await cropToBlob(source, areaPixels)
    if (blob) {
      setPreviewUrl(URL.createObjectURL(blob))
    }
  }

  async function uploadAvatar() {
    if (!source || !croppedAreaPixels) {
      message.warning('请先选择头像图片')
      return
    }
    const blob = await cropToBlob(source, croppedAreaPixels)
    if (!blob) {
      message.warning('头像图片导出失败')
      return
    }

    setUploading(true)
    try {
      await authApi.uploadAvatar(new File([blob], 'avatar.png', { type: 'image/png' }))
      message.success('头像已更新')
      onUploaded()
      onClose()
    } catch {
      // 错误提示由拦截器统一处理
    } finally {
      setUploading(false)
    }
  }

  const hasSource = Boolean(source)

  return (
    <Modal
      open={open}
      title="上传头像"
      footer={
        <Space>
          <span className="avatar-upload__hint">{fileName || '仅支持 JPG、PNG 和 WebP 图片'}</span>
          <Button onClick={onClose}>取消</Button>
          {hasSource ? (
            <Button onClick={() => document.getElementById('avatar-file-input')?.click()}>
              重新选择
            </Button>
          ) : null}
          <Button
            type="primary"
            icon={<UploadOutlined />}
            loading={uploading}
            disabled={!hasSource}
            onClick={() => void uploadAvatar()}
          >
            裁剪并上传
          </Button>
        </Space>
      }
      onCancel={onClose}
      afterOpenChange={(next) => {
        if (!next) resetSource()
      }}
      destroyOnHidden
    >
      <div className="avatar-upload-modal">
        {hasSource ? (
          <div className="avatar-upload-modal__editor">
            <div className="avatar-upload-modal__cropper">
              <Cropper
                image={source}
                crop={crop}
                zoom={zoom}
                aspect={1}
                cropShape="round"
                showGrid={false}
                onCropChange={setCrop}
                onZoomChange={setZoom}
                onCropComplete={onCropComplete}
              />
            </div>
            <div className="avatar-upload-modal__preview-panel">
              <div className="avatar-upload-modal__preview">
                {previewUrl ? <img src={previewUrl} alt="裁剪预览" /> : null}
              </div>
              <div className="mt-3 text-center text-sm text-[var(--ant-color-text-secondary)]">
                实时预览
              </div>
            </div>
            <div className="avatar-upload-modal__zoom">
              <Slider min={1} max={4} step={0.01} value={zoom} onChange={setZoom} />
            </div>
          </div>
        ) : (
          <div className="avatar-upload-modal__empty">
            <Avatar size={96} src={avatar || undefined} icon={<UserOutlined />} />
            <Upload accept={ACCEPTED.join(',')} showUploadList={false} beforeUpload={onFileChange}>
              <Button type="primary" icon={<UploadOutlined />} className="mt-4">
                选择图片
              </Button>
            </Upload>
          </div>
        )}
        <input
          id="avatar-file-input"
          type="file"
          accept={ACCEPTED.join(',')}
          className="hidden"
          onChange={(e) => {
            const file = e.target.files?.[0]
            e.target.value = ''
            if (file) {
              void onFileChange(file)
            }
          }}
        />
      </div>
    </Modal>
  )
}

function createImage(src: string) {
  return new Promise<HTMLImageElement>((resolve, reject) => {
    const image = new Image()
    image.onload = () => resolve(image)
    image.onerror = () => reject(new Error('图片加载失败'))
    image.src = src
  })
}

async function cropToBlob(imageSrc: string, area: Area) {
  try {
    const image = await createImage(imageSrc)
    const canvas = document.createElement('canvas')
    canvas.width = CROP_SIZE
    canvas.height = CROP_SIZE
    const ctx = canvas.getContext('2d')
    if (!ctx) return null

    ctx.drawImage(image, area.x, area.y, area.width, area.height, 0, 0, CROP_SIZE, CROP_SIZE)

    return new Promise<Blob | null>((resolve) => {
      canvas.toBlob((blob) => resolve(blob), 'image/png', 0.92)
    })
  } catch {
    return null
  }
}
