<!--
  Author: Charlie

  消息新建/编辑页。
-->
<script setup lang="ts">
import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import { sysNoticeApi } from '@/api'
import { ACCOUNT_TYPE_OPTIONS } from '@/constants/account'
import { createRequiredArrayRule, createRequiredRule, dictList, formatDateTime } from '@/utils'
import { MdEditor, RichTextEditor } from '@/components/editor'
import UserSelector from '@/components/selector/UserSelector.vue'
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInst | null>(null)
const locationOptions = computed(() => dictList('NOTIFY_LOCATION'))
const targetScopeOptions = computed<SelectOption[]>(() =>
  dictList('TARGET_SCOPE').filter((item: SelectOption) =>
    ['ALL', 'ACCOUNT_TYPE', 'SPECIFIC'].includes(String(item.value)),
  ),
)
const kindOptions = [
  { label: '通知', value: 'NOTIFICATION' },
  { label: '公告', value: 'ANNOUNCEMENT' },
]

const defaultFormData: Record<string, any> = {
  kind: 'NOTIFICATION',
  title: '',
  content: '',
  content_type: 'text',
  category: 'SYSTEM',
  severity: 'INFO',
  target_scope: 'ALL',
  target_account_types: [],
  target_account_ids: [],
  target_dept_ids: [],
  target_role_ids: [],
  publish_locations: {},
  is_pinned: false,
  pinned_until: null,
  status: 'DRAFT',
  publish_at: null,
  expire_at: null,
}

const state = reactive({
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: normalizeFormData(),
  showUserSelector: false,
  userNames: [] as string[],
})

const isAnnouncement = computed(() => state.formModel.kind === 'ANNOUNCEMENT')
const pageTitle = computed(() => (state.dataId ? '编辑消息' : '新增消息'))
const listPath = '/sys/notice'

const rules = computed<FormRules>(() => {
  const scope = state.formModel.target_scope
  const base: FormRules = {
    kind: [createRequiredRule('类型', 'change')],
    title: [createRequiredRule('标题', 'input')],
    content: [createRequiredRule('内容', 'input')],
    content_type: [createRequiredRule('内容格式', 'change')],
    severity: [createRequiredRule('等级', 'change')],
    target_scope: [createRequiredRule('目标范围', 'change')],
    status: [createRequiredRule('状态', 'change')],
    target_account_types: [createRequiredArrayRule('目标账户类型')],
    target_account_ids: scope === 'SPECIFIC' ? [createRequiredArrayRule('目标用户')] : [],
  }
  if (!isAnnouncement.value) {
    base.category = [createRequiredRule('分类', 'change')]
  } else {
    base.publish_locations = [
      {
        required: true,
        trigger: 'change',
        validator: (_rule, value) => {
          const locs = value && typeof value === 'object' ? value : {}
          if (Object.values(locs).some(Boolean)) return true
          return new Error('请选择至少一个发布位置')
        },
      },
    ]
  }
  return base
})

function resolveQueryId() {
  const id = route.query.id
  return typeof id === 'string' && id ? id : null
}

async function initPage() {
  const id = resolveQueryId()
  state.dataId = id
  state.formModel = normalizeFormData()
  state.userNames = []
  if (id) {
    await fetchDetail(id)
  }
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const data = (await sysNoticeApi.detail({ id })).data ?? {}
    state.formModel = normalizeFormData(data)
    state.userNames = data.target_account_ids?.length
      ? data.target_account_ids.map(() => '已选中')
      : []
  } finally {
    state.loading = false
  }
}

function normalizeFormData(data: Record<string, any> = {}): Record<string, any> {
  return {
    ...defaultFormData,
    ...data,
    kind: data.kind || defaultFormData.kind,
    pinned_until: formatDateTime(data.pinned_until, '') || null,
    publish_at: formatDateTime(data.publish_at, '') || null,
    expire_at: formatDateTime(data.expire_at, '') || null,
  }
}

function normalizeSubmitData(data: Record<string, any>): Record<string, any> {
  const r = { ...data }
  const n = (v: unknown) => {
    const t = formatDateTime(v, '')
    if (!t) return null
    const d = new Date(t.replace(' ', 'T') + '+08:00')
    return Number.isNaN(d.getTime()) ? null : d.toISOString().replace(/\.\d{3}Z$/, 'Z')
  }
  r.publish_at = n(data.publish_at)
  if (r.kind === 'ANNOUNCEMENT') {
    r.pinned_until = n(data.pinned_until)
    r.expire_at = n(data.expire_at)
    if (typeof r.publish_locations === 'object' && !Array.isArray(r.publish_locations)) {
      // ok
    } else if (Array.isArray(r.publish_locations)) {
      const d: Record<string, boolean> = {}
      r.publish_locations.forEach((k: string) => (d[k] = true))
      r.publish_locations = d
    }
    r.category = r.category || 'SYSTEM'
  } else {
    r.publish_locations = {}
    r.is_pinned = false
    r.pinned_until = null
    r.expire_at = null
    r.view_count = 0
  }
  delete r.sender_account_type
  delete r.sender_account_id
  delete r.revoked_at
  delete r.source_type
  delete r.source_id
  r.extra = r.extra && typeof r.extra === 'object' ? r.extra : {}
  for (const k of [
    'target_account_types',
    'target_account_ids',
    'target_dept_ids',
    'target_role_ids',
  ]) {
    if (!Array.isArray(r[k])) r[k] = []
  }
  return r
}

function goBack() {
  router.push(listPath)
}

async function submitForm() {
  await formRef.value?.validate()
  state.submitLoading = true
  try {
    const payload = normalizeSubmitData(state.formModel)
    if (state.dataId) {
      await sysNoticeApi.update({ ...payload, id: state.dataId })
      window.$message.success('更新成功')
    } else {
      await sysNoticeApi.create(payload)
      window.$message.success('创建成功')
    }
    goBack()
  } finally {
    state.submitLoading = false
  }
}

function handleUserSelect(account: { id: string; name: string }) {
  state.formModel.target_account_ids = [account.id]
  state.userNames = [account.name]
  state.showUserSelector = false
}

function toggleLocation(key: string, checked: boolean) {
  state.formModel.publish_locations = { ...state.formModel.publish_locations, [key]: checked }
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
        <NForm
          ref="formRef"
          class="message-page-form"
          :model="state.formModel"
          :rules="rules"
          label-placement="left"
          label-width="108"
          :disabled="state.loading || state.submitLoading"
        >
          <div class="form-section-title">
            基本信息
          </div>
          <NGrid
            :cols="2"
            :x-gap="24"
            :y-gap="4"
          >
            <NGi :span="2">
              <NFormItem
                label="类型"
                path="kind"
              >
                <NRadioGroup
                  v-model:value="state.formModel.kind"
                  :disabled="!!state.dataId"
                >
                  <NSpace>
                    <NRadio
                      v-for="item in kindOptions"
                      :key="item.value"
                      :value="item.value"
                      :label="item.label"
                    />
                  </NSpace>
                </NRadioGroup>
              </NFormItem>
            </NGi>
            <NGi :span="2">
              <NFormItem
                label="标题"
                path="title"
              >
                <NInput
                  v-model:value="state.formModel.title"
                  placeholder="请输入消息标题"
                />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                label="内容格式"
                path="content_type"
              >
                <DictSelect
                  v-model="state.formModel.content_type"
                  dict-code="CONTENT_TYPE"
                />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                label="等级"
                path="severity"
              >
                <DictSelect
                  v-model="state.formModel.severity"
                  dict-code="NOTIFICATION_SEVERITY"
                />
              </NFormItem>
            </NGi>
            <NGi v-if="!isAnnouncement">
              <NFormItem
                label="分类"
                path="category"
              >
                <DictSelect
                  v-model="state.formModel.category"
                  dict-code="NOTIFICATION_CATEGORY"
                />
              </NFormItem>
            </NGi>
            <NGi :span="2">
              <NFormItem
                label="状态"
                path="status"
              >
                <DictSelect
                  v-model="state.formModel.status"
                  dict-code="PUBLISH_STATUS"
                  type="radio"
                />
              </NFormItem>
            </NGi>
          </NGrid>

          <div class="form-section-title">
            投放目标
          </div>
          <NGrid
            :cols="2"
            :x-gap="24"
            :y-gap="4"
          >
            <NGi :span="2">
              <NFormItem
                label="目标范围"
                path="target_scope"
              >
                <NRadioGroup v-model:value="state.formModel.target_scope">
                  <NSpace>
                    <NRadio
                      v-for="item in targetScopeOptions"
                      :key="String(item.value)"
                      :value="item.value"
                      :label="String(item.label)"
                    />
                  </NSpace>
                </NRadioGroup>
              </NFormItem>
            </NGi>
            <NGi :span="2">
              <NFormItem
                label="目标账户类型"
                path="target_account_types"
              >
                <NCheckboxGroup v-model:value="state.formModel.target_account_types">
                  <NSpace>
                    <NCheckbox
                      v-for="item in ACCOUNT_TYPE_OPTIONS"
                      :key="item.value"
                      :value="item.value"
                      :label="item.label"
                    />
                  </NSpace>
                </NCheckboxGroup>
              </NFormItem>
            </NGi>
            <NGi
              v-if="state.formModel.target_scope === 'SPECIFIC'"
              :span="2"
            >
              <NFormItem
                label="目标用户"
                path="target_account_ids"
              >
                <NInput
                  :value="state.userNames.join(', ') || ''"
                  readonly
                  placeholder="点击选择用户"
                >
                  <template #suffix>
                    <NButton
                      text
                      type="primary"
                      size="small"
                      @click="state.showUserSelector = true"
                    >
                      选择
                    </NButton>
                  </template>
                </NInput>
              </NFormItem>
            </NGi>
            <NGi
              v-if="isAnnouncement"
              :span="2"
            >
              <NFormItem
                label="发布位置"
                path="publish_locations"
              >
                <NFlex
                  :size="16"
                  :wrap="true"
                >
                  <NCheckbox
                    v-for="opt in locationOptions"
                    :key="opt.value"
                    :checked="!!state.formModel.publish_locations?.[opt.value]"
                    @update:checked="(v: boolean) => toggleLocation(opt.value, v)"
                  >
                    {{ opt.label }}
                  </NCheckbox>
                </NFlex>
              </NFormItem>
            </NGi>
          </NGrid>

          <div class="form-section-title">
            发布设置
          </div>
          <NGrid
            :cols="2"
            :x-gap="24"
            :y-gap="4"
          >
            <NGi>
              <NFormItem
                label="发布时间"
                path="publish_at"
              >
                <NDatePicker
                  v-model:formatted-value="state.formModel.publish_at"
                  type="datetime"
                  value-format="yyyy-MM-dd HH:mm:ss"
                  class="w-full"
                  clearable
                />
              </NFormItem>
            </NGi>
            <NGi v-if="isAnnouncement">
              <NFormItem
                label="过期时间"
                path="expire_at"
              >
                <NDatePicker
                  v-model:formatted-value="state.formModel.expire_at"
                  type="datetime"
                  value-format="yyyy-MM-dd HH:mm:ss"
                  class="w-full"
                  clearable
                />
              </NFormItem>
            </NGi>
            <NGi v-if="isAnnouncement">
              <NFormItem
                label="是否置顶"
                path="is_pinned"
              >
                <NSwitch v-model:value="state.formModel.is_pinned" />
              </NFormItem>
            </NGi>
            <NGi v-if="isAnnouncement && state.formModel.is_pinned">
              <NFormItem
                label="置顶截止"
                path="pinned_until"
              >
                <NDatePicker
                  v-model:formatted-value="state.formModel.pinned_until"
                  type="datetime"
                  value-format="yyyy-MM-dd HH:mm:ss"
                  class="w-full"
                  clearable
                />
              </NFormItem>
            </NGi>
          </NGrid>

          <div class="form-section-title">
            消息内容
          </div>
          <NFormItem
            label="内容"
            path="content"
            label-placement="top"
          >
            <div
              v-if="state.formModel.content_type === 'text'"
              class="w-full"
            >
              <NInput
                v-model:value="state.formModel.content"
                type="textarea"
                :autosize="{ minRows: 8, maxRows: 20 }"
                placeholder="请输入消息内容"
              />
            </div>
            <MdEditor
              v-else-if="state.formModel.content_type === 'markdown'"
              v-model:value="state.formModel.content"
              :height="420"
            />
            <RichTextEditor
              v-else
              v-model:value="state.formModel.content"
              :height="420"
            />
          </NFormItem>
        </NForm>
      </NSpin>
    </NCard>

    <UserSelector
      v-model:visible="state.showUserSelector"
      mode="single"
      @select="handleUserSelect"
    />
  </div>
</template>

<style scoped>
.message-page-form {
  max-width: 960px;
}

.form-section-title {
  margin: 4px 0 14px;
  padding-bottom: 8px;
  color: var(--n-text-color-2, #666);
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.02em;
  border-bottom: 1px solid var(--n-divider-color, rgba(0, 0, 0, 0.08));
}

.form-section-title:not(:first-child) {
  margin-top: 8px;
}
</style>
