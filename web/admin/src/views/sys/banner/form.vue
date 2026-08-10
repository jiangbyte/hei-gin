<!--
  Author: Charlie

  展示图新建/编辑页。
-->
<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import ImageUpload from '@/components/upload/ImageUpload.vue'
import { bannerApi } from '@/api'
import { ACCOUNT_TYPE_OPTIONS } from '@/constants/account'
import {
  createRequiredArrayRule,
  createRequiredRule,
  formatDateTime,
  toApiDateTime,
  toNullableString,
} from '@/utils'
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInst | null>(null)
const listPath = '/sys/banner'

const defaultFormData: Record<string, any> = {
  title: '',
  image: '',
  url: '',
  link_type: 'URL',
  summary: '',
  description: '',
  category: 'HOME',
  type: 'CAROUSEL',
  position: 'HOME_TOP',
  target_account_types: [],
  sort: 0,
  status: 'ENABLED',
  dateRange: null,
}

const state = reactive({
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: { ...defaultFormData },
})

const pageTitle = computed(() => (state.dataId ? '编辑展示图' : '新增展示图'))

const rules = computed<FormRules>(() => ({
  title: createRequiredRule('标题', 'input'),
  image: createRequiredRule('图片', 'input'),
  link_type: createRequiredRule('链接类型', 'change'),
  category: createRequiredRule('分类', 'change'),
  type: createRequiredRule('类型', 'change'),
  position: createRequiredRule('位置', 'change'),
  target_account_types: [createRequiredArrayRule('目标账户类型')],
  status: createRequiredRule('状态', 'change'),
}))

function resolveQueryId() {
  const id = route.query.id
  return typeof id === 'string' && id ? id : null
}

async function initPage() {
  const id = resolveQueryId()
  state.dataId = id
  state.formModel = { ...defaultFormData }
  if (id) {
    await fetchDetail(id)
  }
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await bannerApi.detail({ id })
    const data = response.data ?? {}
    const startStr = data.start_at ? formatDateTime(data.start_at, '') : null
    const endStr = data.end_at ? formatDateTime(data.end_at, '') : null
    state.formModel = {
      ...defaultFormData,
      ...data,
      dateRange: startStr && endStr ? [startStr, endStr] : null,
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
  state.submitLoading = true
  try {
    const payload: Record<string, any> = {
      ...state.formModel,
      url: toNullableString(state.formModel.url),
      summary: toNullableString(state.formModel.summary),
      description: toNullableString(state.formModel.description),
      sort: Number(state.formModel.sort ?? 0),
      start_at: state.formModel.dateRange?.[0] ? toApiDateTime(state.formModel.dateRange[0]) : null,
      end_at: state.formModel.dateRange?.[1] ? toApiDateTime(state.formModel.dateRange[1]) : null,
    }
    delete payload.dateRange
    delete payload.image_url
    delete payload.interaction_count
    delete payload.created_at
    delete payload.created_by
    delete payload.created_name
    delete payload.updated_at
    delete payload.updated_by
    delete payload.updated_name

    if (state.dataId) {
      await bannerApi.update({ ...payload, id: state.dataId })
      window.$message.success('更新成功')
    } else {
      await bannerApi.create(payload)
      window.$message.success('创建成功')
    }
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
        <NForm
          ref="formRef"
          class="banner-page-form"
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
                label="标题"
                path="title"
              >
                <NInput
                  v-model:value="state.formModel.title"
                  placeholder="请输入标题"
                />
              </NFormItem>
            </NGi>
            <NGi :span="2">
              <NFormItem
                label="图片"
                path="image"
              >
                <ImageUpload
                  v-model:value="state.formModel.image"
                  value-type="object_name"
                />
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
            <NGi>
              <NFormItem
                label="状态"
                path="status"
              >
                <DictSelect
                  v-model="state.formModel.status"
                  dict-code="COMMON_STATUS"
                  type="radio"
                />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                label="分类"
                path="category"
              >
                <DictSelect
                  v-model="state.formModel.category"
                  dict-code="BANNER_CATEGORY"
                />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                label="类型"
                path="type"
              >
                <DictSelect
                  v-model="state.formModel.type"
                  dict-code="BANNER_TYPE"
                />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                label="位置"
                path="position"
              >
                <DictSelect
                  v-model="state.formModel.position"
                  dict-code="BANNER_POSITION"
                />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                label="排序"
                path="sort"
              >
                <NInputNumber
                  v-model:value="state.formModel.sort"
                  class="w-full"
                  :min="0"
                />
              </NFormItem>
            </NGi>
          </NGrid>

          <div class="form-section-title">
            跳转与时段
          </div>
          <NGrid
            :cols="2"
            :x-gap="24"
            :y-gap="4"
          >
            <NGi :span="2">
              <NFormItem
                label="链接类型"
                path="link_type"
              >
                <DictSelect
                  v-model="state.formModel.link_type"
                  dict-code="BANNER_LINK_TYPE"
                  type="radio"
                />
              </NFormItem>
            </NGi>
            <NGi :span="2">
              <NFormItem
                label="目标 URL"
                path="url"
              >
                <NInput
                  v-model:value="state.formModel.url"
                  placeholder="可选"
                />
              </NFormItem>
            </NGi>
            <NGi :span="2">
              <NFormItem
                label="展示时间"
                path="dateRange"
              >
                <NDatePicker
                  v-model:formatted-value="state.formModel.dateRange"
                  type="datetimerange"
                  clearable
                  value-format="yyyy-MM-dd HH:mm:ss"
                  class="w-full"
                />
              </NFormItem>
            </NGi>
          </NGrid>

          <div class="form-section-title">
            文案
          </div>
          <NGrid
            :cols="2"
            :x-gap="24"
            :y-gap="4"
          >
            <NGi :span="2">
              <NFormItem
                label="摘要"
                path="summary"
              >
                <NInput
                  v-model:value="state.formModel.summary"
                  placeholder="可选"
                />
              </NFormItem>
            </NGi>
            <NGi :span="2">
              <NFormItem
                label="描述"
                path="description"
              >
                <NInput
                  v-model:value="state.formModel.description"
                  type="textarea"
                  :autosize="{ minRows: 3, maxRows: 8 }"
                  placeholder="可选"
                />
              </NFormItem>
            </NGi>
          </NGrid>
        </NForm>
      </NSpin>
    </NCard>
  </div>
</template>

<style scoped>
.banner-page-form {
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
  margin-top: 12px;
}
</style>
