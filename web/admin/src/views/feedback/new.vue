<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormInst, FormRules, UploadCustomRequestOptions, UploadFileInfo } from 'naive-ui'
import { fileApi, msgFeedbackApi } from '@/api'
import { createRequiredRule, normalizeUploadedFile } from '@/utils'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const formRef = ref<FormInst | null>(null)
const objectNameByUid = new Map<string, string>()
const state = reactive({
  submitting: false,
  fileList: [] as UploadFileInfo[],
  formModel: {
    title: '',
    category: null as string | null,
    content: '',
    contact: '',
  },
})

const rules: FormRules = {
  title: [
    createRequiredRule('标题', 'input'),
    { max: 255, message: '标题最多 255 个字符', trigger: 'blur' },
  ],
  category: [createRequiredRule('分类', 'change')],
  content: [createRequiredRule('内容', 'input')],
}

async function customRequest(options: UploadCustomRequestOptions) {
  const { file, onFinish, onError } = options
  try {
    const raw = file.file
    if (!raw) throw new Error('empty file')
    const res = await fileApi.upload(raw)
    const normalized = normalizeUploadedFile(res.data, raw, 'object_name')
    file.url = normalized.url
    file.name = normalized.name || file.name
    objectNameByUid.set(file.id, normalized.objectName)
    onFinish()
  } catch (error) {
    onError()
    throw error
  }
}

function handleFileListUpdate(list: UploadFileInfo[]) {
  const keep = new Set(list.map((item) => item.id))
  for (const uid of [...objectNameByUid.keys()]) {
    if (!keep.has(uid)) objectNameByUid.delete(uid)
  }
  state.fileList = list
}

function resolveObjectName(item: UploadFileInfo) {
  return String(objectNameByUid.get(item.id) || '').trim()
}

async function handleSubmit() {
  await formRef.value?.validate()
  if (state.fileList.some((item) => item.status === 'uploading')) {
    window.$message.warning('请等待附件上传完成')
    return
  }
  if (state.fileList.some((item) => item.status === 'error')) {
    window.$message.warning('请移除上传失败的附件后再提交')
    return
  }

  const attachObjectNames = state.fileList.map(resolveObjectName).filter(Boolean)

  state.submitting = true
  try {
    await msgFeedbackApi.submit({
      title: state.formModel.title.trim(),
      content: state.formModel.content.trim(),
      category: String(state.formModel.category || ''),
      contact: state.formModel.contact.trim() || null,
      attach_object_names: attachObjectNames,
    })
    window.$message.success('反馈已提交')
    void router.push('/feedback')
  } finally {
    state.submitting = false
  }
}

function goBack() {
  void router.push('/feedback')
}
</script>

<template>
  <div class="h-full min-h-0 p-4">
    <div class="mb-4">
      <NButton
        text
        @click="goBack"
      >
        <template #icon>
          <NovaIcon icon="icon-park-outline:left" />
        </template>
        返回我的反馈
      </NButton>
    </div>

    <NCard
      title="提交反馈"
      :bordered="false"
      class="max-w-3xl"
    >
      <NForm
        ref="formRef"
        :model="state.formModel"
        :rules="rules"
        label-placement="top"
        require-mark-placement="right-hanging"
      >
        <NFormItem
          label="标题"
          path="title"
        >
          <NInput
            v-model:value="state.formModel.title"
            placeholder="简要概括问题或建议"
            clearable
          />
        </NFormItem>
        <NFormItem
          label="分类"
          path="category"
        >
          <DictSelect
            v-model="state.formModel.category"
            dict-code="FEEDBACK_CATEGORY"
            placeholder="请选择分类"
          />
        </NFormItem>
        <NFormItem
          label="内容"
          path="content"
        >
          <NInput
            v-model:value="state.formModel.content"
            type="textarea"
            :rows="6"
            placeholder="请尽量描述清楚场景、期望与复现步骤"
            clearable
          />
        </NFormItem>
        <NFormItem
          label="联系方式"
          path="contact"
        >
          <NInput
            v-model:value="state.formModel.contact"
            placeholder="可选，便于我们联系你"
            clearable
          />
        </NFormItem>
        <NFormItem label="附件">
          <NUpload
            multiple
            directory-dnd
            :file-list="state.fileList"
            :custom-request="customRequest"
            @update:file-list="handleFileListUpdate"
          >
            <NUploadDragger>
              <div class="flex flex-col items-center gap-2 py-4">
                <NovaIcon
                  icon="icon-park-outline:upload-one"
                  :size="28"
                />
                <NText>点击或拖拽文件到此处上传</NText>
                <NText
                  depth="3"
                  class="text-xs"
                >
                  支持多文件上传
                </NText>
              </div>
            </NUploadDragger>
          </NUpload>
        </NFormItem>
        <NFormItem>
          <NButton
            type="primary"
            :loading="state.submitting"
            @click="handleSubmit"
          >
            提交
          </NButton>
        </NFormItem>
      </NForm>
    </NCard>
  </div>
</template>
