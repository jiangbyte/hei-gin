/** Author: Charlie */

import { useState } from 'react'
import { Button, Form, Input, Upload, message } from 'antd'
import { ArrowLeftOutlined, InboxOutlined } from '@ant-design/icons'
import type { UploadFile, UploadProps } from 'antd'
import { Link, useNavigate } from 'react-router-dom'
import { feedbackApi, fileApi } from '@/api'
import { DictSelect } from '@/components/common/DictSelect'
import { normalizeUploadedFile } from '@/utils/file'

type FormValues = {
  title: string
  content: string
  category: string
  contact?: string
}

export function FeedbackNewPage() {
  const navigate = useNavigate()
  const [form] = Form.useForm<FormValues>()
  const [submitting, setSubmitting] = useState(false)
  const [fileList, setFileList] = useState<UploadFile[]>([])

  const customRequest: UploadProps['customRequest'] = async (options) => {
    const { file, onSuccess, onError } = options
    try {
      const res = await fileApi.uploadFile(file as File)
      const normalized = normalizeUploadedFile(res.data, file as File, 'object_name')
      onSuccess?.(
        {
          ...res.data,
          object_name: normalized.objectName,
          url: normalized.url,
          name: normalized.name,
        },
        file as any,
      )
    } catch (error) {
      onError?.(error as Error)
    }
  }

  async function handleSubmit() {
    const values = await form.validateFields()
    const attachObjectNames = fileList
      .map((item) => {
        const response = item.response as { object_name?: string } | undefined
        return response?.object_name || ''
      })
      .filter(Boolean)

    if (fileList.some((item) => item.status === 'uploading')) {
      message.warning('请等待附件上传完成')
      return
    }
    if (fileList.some((item) => item.status === 'error')) {
      message.warning('请移除上传失败的附件后再提交')
      return
    }

    setSubmitting(true)
    try {
      await feedbackApi.submit({
        title: values.title.trim(),
        content: values.content.trim(),
        category: values.category,
        contact: values.contact?.trim() || null,
        attach_object_names: attachObjectNames,
      })
      message.success('反馈已提交')
      void navigate('/feedback')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="page-shell">
      <div className="mb-4">
        <Link to="/feedback">
          <Button type="text" icon={<ArrowLeftOutlined />} className="!px-0">
            返回我的反馈
          </Button>
        </Link>
      </div>

      <section className="panel rounded-xl p-5 md:p-8">
        <h1 className="mb-6 text-xl font-semibold">提交反馈</h1>
        <Form
          form={form}
          layout="vertical"
          className="max-w-2xl"
          onFinish={() => void handleSubmit()}
        >
          <Form.Item
            name="title"
            label="标题"
            rules={[
              { required: true, message: '请输入标题' },
              { max: 255, message: '标题最多 255 个字符' },
            ]}
          >
            <Input placeholder="简要概括问题或建议" allowClear />
          </Form.Item>
          <Form.Item
            name="category"
            label="分类"
            rules={[{ required: true, message: '请选择分类' }]}
          >
            <DictSelect dictCode="FEEDBACK_CATEGORY" placeholder="请选择分类" />
          </Form.Item>
          <Form.Item
            name="content"
            label="内容"
            rules={[{ required: true, message: '请输入内容' }]}
          >
            <Input.TextArea rows={6} placeholder="请尽量描述清楚场景、期望与复现步骤" allowClear />
          </Form.Item>
          <Form.Item name="contact" label="联系方式">
            <Input placeholder="可选，便于我们联系你" allowClear />
          </Form.Item>
          <Form.Item label="附件">
            <Upload.Dragger
              multiple
              fileList={fileList}
              customRequest={customRequest}
              onChange={({ fileList: next }) => setFileList(next)}
              onRemove={(file) => {
                setFileList((prev) => prev.filter((item) => item.uid !== file.uid))
                return true
              }}
            >
              <p className="ant-upload-drag-icon">
                <InboxOutlined />
              </p>
              <p className="ant-upload-text">点击或拖拽文件到此处上传</p>
              <p className="ant-upload-hint">支持多文件上传</p>
            </Upload.Dragger>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting}>
              提交
            </Button>
          </Form.Item>
        </Form>
      </section>
    </div>
  )
}
