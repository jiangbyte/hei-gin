<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { dictApi } from '@/api'
import CommonColorPicker from '@/components/common/CommonColorPicker.vue'
import { createRequiredRule, isHexColor, toNullableString } from '@/utils'
import { computed, reactive, ref, watch } from 'vue'

const props = defineProps<{
  dicts: any[]
}>()

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)
const defaultFormData = {
  code: '',
  label: '',
  value: '',
  color: '',
  category: 'SYS',
  parent_id: null as string | null,
  status: 'ENABLED',
  sort: 99,
}
const state = reactive({
  showModal: false,
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: { ...defaultFormData },
})

const modalTitle = computed(() => (state.dataId ? '编辑字典' : '新增字典'))
const parentTreeOptions = computed(() =>
  buildTreeOptions(
    props.dicts.filter((item) => item.category === state.formModel.category),
    state.dataId,
  ),
)

const rules = computed<FormRules>(() => ({
  code: [
    createRequiredRule('编码', 'input'),
    {
      pattern: /^[A-Z0-9_]+$/,
      message: '编码只能包含大写字母、数字和下划线',
      trigger: ['input', 'blur'],
    },
  ],
  label: createRequiredRule('标签', 'input'),
  color: [
    {
      validator: (_rule, value) => isHexColor(value),
      message: '请选择有效的十六进制颜色',
      trigger: ['change', 'blur'],
    },
  ],
  category: createRequiredRule('分类', 'change'),
  status: createRequiredRule('状态', 'change'),
}))

// 切换分类时重置父级字典
watch(
  () => state.formModel.category,
  () => {
    state.formModel.parent_id = null
  },
)
async function openModal(id?: string, options?: { category?: string; parentId?: string | null }) {
  state.dataId = id ?? null
  state.formModel = {
    ...defaultFormData,
    category: options?.category ?? 'SYS',
    parent_id: options?.parentId ?? null,
  }
  state.showModal = true

  if (id) {
    await fetchDetail(id)
  }
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await dictApi.detail({ id })
    const data = response.data ?? {}
    state.formModel = Object.assign({}, defaultFormData, data, {
      value: data.value ?? '',
      color: data.color ?? '',
      category: data.category ?? 'SYS',
      parent_id: data.parent_id ?? null,
      status: data.status ?? 'ENABLED',
      sort: data.sort ?? 0,
    })
  } finally {
    state.loading = false
  }
}

function closeModal() {
  state.showModal = false
  state.submitLoading = false
}

async function submitForm() {
  await formRef.value?.validate()
  const payload = {
    ...state.formModel,
    code: state.formModel.code.trim().toUpperCase(),
    label: String(state.formModel.label ?? '').trim(),
    value: toNullableString(state.formModel.value),
    color: toNullableString(state.formModel.color),
    parent_id: state.formModel.parent_id ?? null,
    sort: Number(state.formModel.sort ?? 0),
  }

  state.submitLoading = true
  try {
    if (state.dataId) {
      await dictApi.update({
        ...payload,
        id: state.dataId,
      })
      window.$message.success('更新成功')
    } else {
      await dictApi.create(payload)
      window.$message.success('创建成功')
    }

    emit('saved')
    closeModal()
  } finally {
    state.submitLoading = false
  }
}

function updateCode(value: string) {
  state.formModel.code = value.toUpperCase()
}

function buildTreeOptions(items: any[], excludeId?: string | null) {
  const excludeIds = excludeId ? collectChildIds(items, excludeId) : new Set<string>()
  if (excludeId) {
    excludeIds.add(excludeId)
  }

  const nodeMap = new Map(
    items
      .filter((item) => !excludeIds.has(item.id))
      .map((item) => [
        item.id,
        {
          key: item.id,
          label: `${item.label || item.code} (${item.code})`,
          children: [] as any[],
          raw: item,
        },
      ]),
  )
  const roots: any[] = []

  nodeMap.forEach((node) => {
    const parentId = node.raw.parent_id
    const parent = parentId ? nodeMap.get(parentId) : null
    if (parent) {
      parent.children.push(node)
    } else {
      roots.push(node)
    }
  })

  return normalizeTreeOptions(roots)
}

function collectChildIds(items: any[], parentId: string) {
  const result = new Set<string>()
  const walk = (id: string) => {
    items
      .filter((item) => item.parent_id === id)
      .forEach((item) => {
        result.add(item.id)
        walk(item.id)
      })
  }
  walk(parentId)
  return result
}

function sortTreeOptions(a: any, b: any) {
  return (a.raw?.sort ?? 0) - (b.raw?.sort ?? 0)
}

function normalizeTreeOptions(nodes: any[]): any[] {
  return nodes.sort(sortTreeOptions).map((node) => ({
    key: node.key,
    label: node.label,
    children: normalizeTreeOptions(node.children),
  }))
}

defineExpose({
  openModal,
})
</script>

<template>
  <NModal
    v-model:show="state.showModal"
    preset="card"
    draggable
    :mask-closable="false"
    :title="modalTitle"
    style="width: 640px"
    :segmented="{ content: true, action: true }"
  >
    <NSpin :show="state.loading">
      <NForm
        ref="formRef"
        :model="state.formModel"
        :rules="rules"
        label-placement="left"
        label-width="100"
        :disabled="state.loading || state.submitLoading"
      >
        <NFormItem
          :label="'分类'"
          path="category"
        >
          <DictSelect
            v-model="state.formModel.category"
            dict-code="SYS_BIZ_CATEGORY"
            type="radio"
          />
        </NFormItem>
        <NFormItem
          :label="'父级字典'"
          path="parent_id"
        >
          <NTreeSelect
            v-model:value="state.formModel.parent_id"
            clearable
            filterable
            :options="parentTreeOptions"
            :placeholder="'Top 等级'"
            key-field="key"
            label-field="label"
          />
        </NFormItem>
        <NFormItem
          :label="'编码'"
          path="code"
        >
          <NInput
            :value="state.formModel.code"
            @update:value="updateCode"
          />
        </NFormItem>
        <NFormItem
          :label="'标签'"
          path="label"
        >
          <NInput v-model:value="state.formModel.label" />
        </NFormItem>
        <NFormItem
          :label="'值'"
          path="value"
        >
          <NInput v-model:value="state.formModel.value" />
        </NFormItem>
        <NFormItem
          :label="'颜色'"
          path="color"
        >
          <CommonColorPicker
            v-model="state.formModel.color"
            :disabled="state.loading || state.submitLoading"
          />
        </NFormItem>
        <NFormItem
          :label="'排序'"
          path="sort"
        >
          <NInputNumber
            v-model:value="state.formModel.sort"
            class="w-full"
            :min="0"
          />
        </NFormItem>
        <NFormItem
          :label="'状态'"
          path="status"
        >
          <DictSelect
            v-model="state.formModel.status"
            dict-code="COMMON_STATUS"
            type="radio"
          />
        </NFormItem>
      </NForm>
    </NSpin>

    <template #action>
      <NSpace
        justify="end"
        align="center"
      >
        <NButton @click="closeModal">
          取消
        </NButton>
        <NButton
          type="primary"
          :loading="state.submitLoading"
          @click="submitForm"
        >
          确认
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
