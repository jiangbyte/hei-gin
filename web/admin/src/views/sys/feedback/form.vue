<!--
  Author: Charlie

  反馈处理页 — 仅更新状态与回复。
-->
<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { sysFeedbackApi } from '@/api'
import { createRequiredRule, dictTypeData, displayValue, formatDateTime } from '@/utils'
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInst | null>(null)
const listPath = '/sys/feedback'

const defaultFormData: Record<string, any> = {
  status: 'PENDING',
  reply: null,
}

const state = reactive({
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: { ...defaultFormData },
  detail: {} as any,
})

const pageTitle = computed(() => '处理反馈')
const rules = computed<FormRules>(() => ({
  status: [createRequiredRule('状态', 'change')],
}))

function resolveQueryId() {
  const id = route.query.id
  return typeof id === 'string' && id ? id : null
}

async function initPage() {
  const id = resolveQueryId()
  state.dataId = id
  state.formModel = { ...defaultFormData }
  state.detail = {}
  if (id) {
    await fetchDetail(id)
  }
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const data = (await sysFeedbackApi.detail({ id })).data ?? {}
    state.detail = data
    state.formModel = {
      status: data.status || 'PENDING',
      reply: data.reply ?? null,
    }
  } finally {
    state.loading = false
  }
}

function goBack() {
  router.push(listPath)
}

async function submitForm() {
  await formRef.value?.validate()
  if (!state.dataId) return
  state.submitLoading = true
  try {
    await sysFeedbackApi.update({
      id: state.dataId,
      status: state.formModel.status,
      reply: state.formModel.reply || null,
    })
    window.$message.success('更新成功')
    goBack()
  } finally {
    state.submitLoading = false
  }
}

onMounted(initPage)
watch(
  () => route.query.id,
  () => {
    void initPage()
  },
)
</script>

<template>
  <div class="h-full min-h-0">
    <NCard
      class="h-full min-h-0 overflow-auto"
      :title="pageTitle"
      :bordered="false"
    >
      <template #header-extra>
        <NSpace>
          <NButton @click="goBack">
            返回
          </NButton>
          <NButton
            type="primary"
            :loading="state.submitLoading"
            @click="submitForm"
          >
            保存
          </NButton>
        </NSpace>
      </template>
      <NSpin :show="state.loading">
        <div class="message-page-form">
          <div class="form-section-title">
            反馈摘要
          </div>
          <div class="summary-grid">
            <div class="summary-item summary-item--full">
              <div class="summary-key">
                标题
              </div>
              <div class="summary-value">
                {{ displayValue(state.detail.title) }}
              </div>
            </div>
            <div class="summary-item">
              <div class="summary-key">
                分类
              </div>
              <div class="summary-value">
                {{
                  dictTypeData('FEEDBACK_CATEGORY', state.detail.category) ||
                    displayValue(state.detail.category)
                }}
              </div>
            </div>
            <div class="summary-item">
              <div class="summary-key">
                提交时间
              </div>
              <div class="summary-value">
                {{ formatDateTime(state.detail.created_at) }}
              </div>
            </div>
            <div class="summary-item summary-item--full">
              <div class="summary-key">
                反馈内容
              </div>
              <div class="summary-value whitespace-pre-wrap">
                {{ displayValue(state.detail.content) }}
              </div>
            </div>
          </div>

          <div class="form-section-title">
            处理
          </div>
          <NForm
            ref="formRef"
            :model="state.formModel"
            :rules="rules"
            label-placement="left"
            label-width="108"
            :disabled="state.loading || state.submitLoading"
          >
            <NFormItem
              label="状态"
              path="status"
            >
              <DictSelect
                v-model="state.formModel.status"
                dict-code="FEEDBACK_STATUS"
                type="radio"
              />
            </NFormItem>
            <NFormItem
              label="管理员回复"
              path="reply"
            >
              <NInput
                v-model:value="state.formModel.reply"
                type="textarea"
                :autosize="{ minRows: 4, maxRows: 12 }"
                placeholder="输入回复内容（可选）"
              />
            </NFormItem>
          </NForm>
        </div>
      </NSpin>
    </NCard>
  </div>
</template>

<style scoped>
.message-page-form {
  max-width: 960px;
}

.form-section-title {
  margin: 4px 0 14px;
  padding-bottom: 8px;
  color: var(--text-color-2, #666);
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.02em;
  border-bottom: 1px solid var(--n-divider-color, rgba(0, 0, 0, 0.08));
}

.form-section-title:not(:first-child) {
  margin-top: 20px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px 24px;
  margin-bottom: 8px;
}

.summary-item--full {
  grid-column: 1 / -1;
}

.summary-key {
  margin-bottom: 4px;
  color: var(--text-color-3, #999);
  font-size: 12px;
}

.summary-value {
  color: var(--text-color-1, #333);
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
}

.whitespace-pre-wrap {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
