<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { resourceApi, resourceModuleApi } from '@/api'
import CommonColorPicker from '@/components/common/CommonColorPicker.vue'
import IconSelect from '@/components/common/IconSelect.vue'
import { createRequiredRule, isHexColor, toNullableString } from '@/utils'
import { computed, reactive, ref } from 'vue'

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)
const defaultFormData = {
  code: '',
  name: '',
  resource_type: null as string | null,
  parent_id: null as string | null,
  module_id: null as string | null,
  path: '',
  component: '',
  redirect: '',
  icon: '',
  color: '',
  href: '',
  sort: 0,
  is_visible: true,
  is_cache: false,
  is_affix: false,
  status: 'ENABLED',
  description: '',
  layout: '' as string | null,
  extra: {},
}
const state = reactive({
  showModal: false,
  loading: false,
  treeLoading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: { ...defaultFormData },
  resourceTree: [] as any[],
  moduleOptions: [] as any[],
})

const modalTitle = computed(() => (state.dataId ? '编辑资源' : '新增资源'))

const rules = computed<FormRules>(() => ({
  code: createRequiredRule('资源编码', 'input'),
  name: createRequiredRule('资源名称', 'input'),
  resource_type: createRequiredRule('资源类型', 'change'),
  module_id: createRequiredRule('资源模块', 'change'),
  color: [
    {
      validator: (_rule, value) => isHexColor(value),
      message: '请输入十六进制颜色，例如 #1677ff',
      trigger: ['change', 'blur'],
    },
  ],
  status: createRequiredRule('状态', 'change'),
}))

const parentTreeOptions = computed(() => {
  const excludedIds = new Set(
    state.dataId ? collectDescendantIds(state.resourceTree, state.dataId) : [],
  )
  return buildParentTreeOptions(state.resourceTree, excludedIds)
})

async function openModal(id?: string, parentId?: string, moduleId?: string | null) {
  state.dataId = id ?? null
  state.formModel = {
    ...defaultFormData,
    parent_id: parentId ?? null,
    module_id: moduleId ?? null,
  }
  state.showModal = true
  await Promise.all([fetchResourceTree(moduleId), fetchModules()])

  if (id) {
    await fetchDetail(id)
  } else if (parentId) {
    const parent = findResourceNode(state.resourceTree, parentId)
    state.formModel.module_id = parent?.module_id ?? null
  }
}

async function fetchResourceTree(moduleId?: string | null) {
  state.treeLoading = true
  try {
    const response = await resourceApi.tree({
      module_id: moduleId ?? undefined,
    })
    state.resourceTree = response.data ?? []
  } finally {
    state.treeLoading = false
  }
}

async function fetchModules() {
  const response = await resourceModuleApi.selector()
  state.moduleOptions = (response.data ?? []).map((item: any) => ({
    label: item.name,
    value: item.id,
  }))
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await resourceApi.detail({ id })
    state.formModel = Object.assign({}, defaultFormData, response.data, {
      parent_id: response.data?.parent_id ?? null,
      module_id: response.data?.module_id ?? null,
      path: response.data?.path ?? '',
      component: response.data?.component ?? '',
      redirect: response.data?.redirect ?? '',
      icon: response.data?.icon ?? '',
      color: response.data?.color ?? '',
      href: response.data?.href ?? '',
      description: response.data?.description ?? '',
      layout: response.data?.layout ?? null,
      extra: response.data?.extra ?? {},
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
  state.submitLoading = true
  try {
    const payload = {
      ...state.formModel,
      code: state.formModel.code.trim(),
      name: state.formModel.name.trim(),
      parent_id: toNullableString(state.formModel.parent_id),
      module_id: toNullableString(state.formModel.module_id),
      path: toNullableString(state.formModel.path),
      component: toNullableString(state.formModel.component),
      redirect: toNullableString(state.formModel.redirect),
      icon: toNullableString(state.formModel.icon),
      color: toNullableString(state.formModel.color),
      href: toNullableString(state.formModel.href),
      sort: Number(state.formModel.sort ?? 0),
      is_visible: Boolean(state.formModel.is_visible),
      is_cache: Boolean(state.formModel.is_cache),
      is_affix: Boolean(state.formModel.is_affix),
      description: toNullableString(state.formModel.description),
      layout: toNullableString(state.formModel.layout),
      extra: state.formModel.extra ?? {},
    }

    if (state.dataId) {
      await resourceApi.update({
        ...payload,
        id: state.dataId,
      })
      window.$message.success('更新成功')
    } else {
      await resourceApi.create(payload)
      window.$message.success('创建成功')
    }

    closeModal()
    emit('saved')
  } finally {
    state.submitLoading = false
  }
}

defineExpose({
  openModal,
})

function buildParentTreeOptions(items: any[], excludedIds: Set<string>): any[] {
  return items
    .filter((item) => !excludedIds.has(item.id))
    .map((item) => ({
      id: item.id,
      name: `${item.name} (${item.code})`,
      children: buildParentTreeOptions(item.children ?? [], excludedIds),
    }))
}

function collectDescendantIds(items: any[], targetId: string) {
  const result = new Set<string>([targetId])
  const target = findResourceNode(items, targetId)
  const walk = (nodes: any[]) => {
    nodes.forEach((node) => {
      result.add(node.id)
      walk(node.children ?? [])
    })
  }
  if (target) {
    walk(target.children ?? [])
  }
  return Array.from(result)
}

function findResourceNode(items: any[], id: string): any | null {
  for (const item of items) {
    if (item.id === id) {
      return item
    }
    const child = findResourceNode(item.children ?? [], id)
    if (child) {
      return child
    }
  }
  return null
}
</script>

<template>
  <NModal
    v-model:show="state.showModal"
    preset="card"
    draggable
    :mask-closable="false"
    :title="modalTitle"
    style="width: 760px"
    :segmented="{ content: true, action: true }"
  >
    <NSpin :show="state.loading || state.treeLoading">
      <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
        <NForm
          ref="formRef"
          :model="state.formModel"
          :rules="rules"
          label-placement="left"
          label-width="110"
          :disabled="state.loading || state.treeLoading || state.submitLoading"
        >
          <NGrid
            :cols="2"
            :x-gap="16"
          >
            <NGi>
              <NFormItem
                :label="'资源名称'"
                path="name"
              >
                <NInput v-model:value="state.formModel.name" />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                :label="'资源编码'"
                path="code"
              >
                <NInput v-model:value="state.formModel.code" />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                :label="'资源类型'"
                path="resource_type"
              >
                <DictSelect
                  v-model="state.formModel.resource_type"
                  dict-code="RESOURCE_TYPE"
                  :placeholder="'请选择资源类型'"
                />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                :label="'资源模块'"
                path="module_id"
              >
                <NSelect
                  v-model:value="state.formModel.module_id"
                  filterable
                  clearable
                  :options="state.moduleOptions"
                />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                :label="'父级资源ID'"
                path="parent_id"
              >
                <NTreeSelect
                  v-model:value="state.formModel.parent_id"
                  clearable
                  filterable
                  :options="parentTreeOptions"
                  :placeholder="'父级资源ID'"
                  key-field="id"
                  label-field="name"
                  children-field="children"
                />
              </NFormItem>
            </NGi>
            <NGi>
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
            </NGi>
            <NGi>
              <NFormItem
                :label="'图标'"
                path="icon"
              >
                <IconSelect v-model:value="state.formModel.icon" />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                :label="'布局'"
                path="layout"
              >
                <NInput
                  v-model:value="state.formModel.layout"
                  placeholder="default / fullscreen"
                />
              </NFormItem>
            </NGi>
            <NGi :span="2">
              <NFormItem
                :label="'颜色'"
                path="color"
              >
                <CommonColorPicker v-model="state.formModel.color" />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                :label="'路由路径'"
                path="path"
              >
                <NInput v-model:value="state.formModel.path" />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                :label="'组件'"
                path="component"
              >
                <NInput v-model:value="state.formModel.component" />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                :label="'重定向'"
                path="redirect"
              >
                <NInput v-model:value="state.formModel.redirect" />
              </NFormItem>
            </NGi>
            <NGi>
              <NFormItem
                :label="'外链'"
                path="href"
              >
                <NInput v-model:value="state.formModel.href" />
              </NFormItem>
            </NGi>
            <NGi :span="2">
              <NFlex :size="24">
                <NFormItem
                  :label="'可见'"
                  path="is_visible"
                >
                  <NSwitch v-model:value="state.formModel.is_visible" />
                </NFormItem>
                <NFormItem
                  :label="'缓存'"
                  path="is_cache"
                >
                  <NSwitch v-model:value="state.formModel.is_cache" />
                </NFormItem>
                <NFormItem
                  :label="'固定标签'"
                  path="is_affix"
                >
                  <NSwitch v-model:value="state.formModel.is_affix" />
                </NFormItem>
              </NFlex>
            </NGi>
            <NGi :span="2">
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
            </NGi>
          </NGrid>
          <NFormItem
            :label="'描述'"
            path="description"
          >
            <NInput
              v-model:value="state.formModel.description"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
            />
          </NFormItem>
        </NForm>
      </NScrollbar>
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
