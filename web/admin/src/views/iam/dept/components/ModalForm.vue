<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormInst, FormRules } from 'naive-ui'
import { deptApi } from '@/api'
import { createRequiredRule, toNullableString } from '@/utils'
import { computed, reactive, ref } from 'vue'
import UserSelector from '@/components/selector/UserSelector.vue'

const emit = defineEmits<{
  saved: []
}>()

const formRef = ref<FormInst | null>(null)
const defaultFormData = {
  name: '',
  category: null as string | null,
  parent_id: null as string | null,
  master_id: null as string | null,
  master_name: '',
  deputy_master_id: null as string | null,
  deputy_master_name: '',
  sort: 0,
  is_virtual: false,
  status: 'ENABLED',
  extra: {},
}
const state = reactive({
  showModal: false,
  loading: false,
  treeLoading: false,
  submitLoading: false,
  dataId: null as string | null,
  formModel: { ...defaultFormData },
  deptTree: [] as any[],
})

const modalTitle = computed(() => (state.dataId ? '编辑 部门' : '新增 部门'))

const rules = computed<FormRules>(() => ({
  name: createRequiredRule('部门名称', 'input'),
  category: createRequiredRule('部门分类', 'change'),
  status: createRequiredRule('状态', 'change'),
}))

const parentTreeOptions = computed(() => {
  const excludedIds = new Set(
    state.dataId ? collectDescendantIds(state.deptTree, state.dataId) : [],
  )
  return buildParentTreeOptions(state.deptTree, excludedIds)
})

/* ---- 账号用户选择器 ---- */

const userSelectorState = reactive({
  show: false,
  title: '',
  targetField: null as 'master_id' | 'deputy_master_id' | null,
})

function openUserSelector(field: 'master_id' | 'deputy_master_id') {
  userSelectorState.title = field === 'master_id' ? '选择负责人' : '选择副负责人'
  userSelectorState.targetField = field
  userSelectorState.show = true
}

function handleUserSelect(account: { id: string; name: string }) {
  if (userSelectorState.targetField === 'master_id') {
    state.formModel.master_id = account.id
    state.formModel.master_name = account.name
  } else if (userSelectorState.targetField === 'deputy_master_id') {
    state.formModel.deputy_master_id = account.id
    state.formModel.deputy_master_name = account.name
  }
  userSelectorState.show = false
}

async function openModal(id?: string) {
  state.dataId = id ?? null
  state.formModel = { ...defaultFormData }
  state.showModal = true

  await fetchDeptTree()

  if (id) {
    await fetchDetail(id)
  }
}

async function fetchDeptTree() {
  state.treeLoading = true
  try {
    const response = await deptApi.tree()
    state.deptTree = response.data ?? []
  } finally {
    state.treeLoading = false
  }
}

async function fetchDetail(id: string) {
  state.loading = true
  try {
    const response = await deptApi.detail({ id })
    state.formModel = Object.assign({}, defaultFormData, response.data, {
      parent_id: response.data?.parent_id ?? null,
      master_id: response.data?.master_id ?? null,
      deputy_master_id: response.data?.deputy_master_id ?? null,
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
      name: state.formModel.name.trim(),
      parent_id: toNullableString(state.formModel.parent_id),
      master_id: toNullableString(state.formModel.master_id),
      deputy_master_id: toNullableString(state.formModel.deputy_master_id),
      sort: Number(state.formModel.sort ?? 0),
      is_virtual: Boolean(state.formModel.is_virtual),
      extra: state.formModel.extra ?? {},
    }

    if (state.dataId) {
      await deptApi.update({ ...payload, id: state.dataId })
      window.$message.success('更新成功')
    } else {
      await deptApi.create(payload)
      window.$message.success('创建成功')
    }

    closeModal()
    emit('saved')
  } finally {
    state.submitLoading = false
  }
}

/* ---- 树形辅助 ---- */

function collectDescendantIds(nodes: any[], id: string): string[] {
  const ids: string[] = [id]
  for (const node of nodes) {
    if (node.id === id) {
      ids.push(...(node.children ? flattenIds(node.children) : []))
    } else if (node.children) {
      ids.push(...collectDescendantIds(node.children, id))
    }
  }
  return ids
}

function flattenIds(nodes: any[]): string[] {
  const ids: string[] = []
  for (const node of nodes) {
    ids.push(node.id)
    if (node.children) ids.push(...flattenIds(node.children))
  }
  return ids
}

function buildParentTreeOptions(nodes: any[], excludedIds: Set<string>): any[] {
  return nodes
    .filter((node) => !excludedIds.has(node.id))
    .map((node) => ({
      key: node.id,
      label: String(node.name ?? node.id ?? ''),
      children: node.children ? buildParentTreeOptions(node.children, excludedIds) : undefined,
    }))
}

defineExpose({ openModal })
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
          <NFormItem
            :label="'部门名称'"
            path="name"
          >
            <NInput v-model:value="state.formModel.name" />
          </NFormItem>
          <NFormItem
            :label="'部门分类'"
            path="category"
          >
            <DictSelect
              v-model="state.formModel.category"
              dict-code="DEPT_CATEGORY"
              :placeholder="'请选择部门分类'"
            />
          </NFormItem>
          <NFormItem
            :label="'父级部门'"
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
              children-field="children"
            />
          </NFormItem>
          <NFormItem
            :label="'负责人'"
            path="master_id"
          >
            <NInputGroup>
              <NInput
                :value="state.formModel.master_name || ''"
                readonly
                :placeholder="'点击右侧按钮选择负责人'"
              />
              <NButton
                type="primary"
                @click="openUserSelector('master_id')"
              >
                <template #icon>
                  <NovaIcon
                    icon="icon-park-outline:search"
                    :size="16"
                  />
                </template>
              </NButton>
              <NButton
                v-if="state.formModel.master_id"
                @click="
                  () => {
                    state.formModel.master_id = null
                    state.formModel.master_name = ''
                  }
                "
              >
                <template #icon>
                  <NovaIcon
                    icon="icon-park-outline:close"
                    :size="16"
                  />
                </template>
              </NButton>
            </NInputGroup>
          </NFormItem>
          <NFormItem
            :label="'副负责人'"
            path="deputy_master_id"
          >
            <NInputGroup>
              <NInput
                :value="state.formModel.deputy_master_name || ''"
                readonly
                :placeholder="'点击右侧按钮选择副负责人'"
              />
              <NButton
                type="primary"
                @click="openUserSelector('deputy_master_id')"
              >
                <template #icon>
                  <NovaIcon
                    icon="icon-park-outline:search"
                    :size="16"
                  />
                </template>
              </NButton>
              <NButton
                v-if="state.formModel.deputy_master_id"
                @click="
                  () => {
                    state.formModel.deputy_master_id = null
                    state.formModel.deputy_master_name = ''
                  }
                "
              >
                <template #icon>
                  <NovaIcon
                    icon="icon-park-outline:close"
                    :size="16"
                  />
                </template>
              </NButton>
            </NInputGroup>
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
            :label="'虚拟部门'"
            path="is_virtual"
          >
            <NSwitch v-model:value="state.formModel.is_virtual" />
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

  <!-- 用户选择器抽屉 -->
  <UserSelector
    v-model:visible="userSelectorState.show"
    title="选择用户"
    mode="single"
    @select="handleUserSelect"
  />
</template>
