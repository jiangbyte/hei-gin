<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { deptApi, roleApi } from '@/api'
import { createRequiredRule, toNullableString } from '@/utils'
import { computed, reactive, ref } from 'vue'

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)
const defaultFormData = {
  code: '',
  name: '',
  category: null as string | null,
  scope_type: null as string | null,
  owner_dept_id: '',
  sort: 0,
  status: 'ENABLED',
  is_builtin: false,
  description: '',
  extra: {},
}
const state = reactive({
  showModal: false,
  loading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: { ...defaultFormData },
  deptTree: [] as any[],
})

const modalTitle = computed(() => (state.dataId ? '编辑角色' : '新增角色'))

const rules = computed<FormRules>(() => ({
  code: createRequiredRule('角色编码', 'input'),
  name: createRequiredRule('角色名称', 'input'),
  category: createRequiredRule('角色分类', 'change'),
  scope_type: createRequiredRule('范围类型', 'change'),
  status: createRequiredRule('状态', 'change'),
}))

function buildDeptTreeOptions(nodes: any[]): any[] {
  return nodes.map((node: any) => ({
    key: node.id,
    label: node.name,
    children: node.children ? buildDeptTreeOptions(node.children) : undefined,
  }))
}

const deptTreeOptions = computed(() => buildDeptTreeOptions(state.deptTree))

async function openModal(id?: string) {
  state.dataId = id ?? null
  state.formModel = { ...defaultFormData }
  state.showModal = true
  fetchDeptTree()

  if (id) {
    await fetchDetail(id)
  }
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await roleApi.detail({ id })
    state.formModel = Object.assign({}, defaultFormData, response.data, {
      owner_dept_id: response.data?.owner_dept_id ?? '',
      description: response.data?.description ?? '',
      extra: response.data?.extra ?? {},
    })
  } finally {
    state.loading = false
  }
}

async function fetchDeptTree() {
  try {
    const resp = await deptApi.tree()
    state.deptTree = resp.data ?? []
  } catch {
    state.deptTree = []
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
      owner_dept_id: toNullableString(state.formModel.owner_dept_id),
      sort: Number(state.formModel.sort ?? 0),
      is_builtin: Boolean(state.formModel.is_builtin),
      description: toNullableString(state.formModel.description),
      extra: state.formModel.extra ?? {},
    }

    if (state.dataId) {
      await roleApi.update({
        ...payload,
        id: state.dataId,
      })
      window.$message.success('更新成功')
    } else {
      await roleApi.create(payload)
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
</script>

<template>
  <NModal
    v-model:show="state.showModal"
    preset="card"
    draggable
    :mask-closable="false"
    :title="modalTitle"
    style="width: 720px"
    :segmented="{ content: true, action: true }"
  >
    <NSpin :show="state.loading">
      <NScrollbar class="max-h-[min(620px,calc(100vh-300px))] pr-16px">
        <NForm
          ref="formRef"
          :model="state.formModel"
          :rules="rules"
          label-placement="left"
          label-width="110"
          :disabled="state.loading || state.submitLoading"
        >
          <NFormItem
            :label="'角色名称'"
            path="name"
          >
            <NInput v-model:value="state.formModel.name" />
          </NFormItem>
          <NFormItem
            :label="'角色编码'"
            path="code"
          >
            <NInput v-model:value="state.formModel.code" />
          </NFormItem>
          <NFormItem
            :label="'角色分类'"
            path="category"
          >
            <DictSelect
              v-model="state.formModel.category"
              dict-code="SYS_BIZ_CATEGORY"
            />
          </NFormItem>
          <NFormItem
            :label="'范围类型'"
            path="scope_type"
          >
            <DictSelect
              v-model="state.formModel.scope_type"
              dict-code="ROLE_SCOPE_TYPE"
            />
          </NFormItem>
          <NFormItem
            :label="'所属部门'"
            path="owner_dept_id"
          >
            <NTreeSelect
              v-model:value="state.formModel.owner_dept_id"
              clearable
              filterable
              :options="deptTreeOptions"
              :placeholder="'请选择所属部门'"
              key-field="key"
              label-field="label"
              children-field="children"
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
            :label="'内置角色'"
            path="is_builtin"
          >
            <NSwitch v-model:value="state.formModel.is_builtin" />
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
