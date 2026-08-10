<!-- Author: Charlie -->

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { deptApi, roleApi } from '@/api'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import ModalGrantDept from '@/views/iam/components/ModalGrantDept.vue'
import ModalGrantRole from '@/views/iam/components/ModalGrantRole.vue'
import {
  ACCOUNT_TYPE_TABS,
  DEFAULT_ACCOUNT_TYPE,
  accountConfigKey,
  accountTypeLabel,
  createAccountTypeMap,
  mapAccountTypes,
  type AccountType,
} from '@/constants/account'
import { loadByCategory, parseBool, saveByKeys, toBoolStr } from '../composables/useConfigForm'

const CATEGORY = 'AUTH_REGISTER'
const PREFIX = 'AUTH_REGISTER'

type ScopeForm = {
  enabled: boolean
  requirePhone: boolean
  requireEmail: boolean
  defaultRoleId: string
  defaultRoleName: string
  defaultDeptId: string
  defaultDeptName: string
}

function emptyScope(): ScopeForm {
  return {
    enabled: false,
    requirePhone: false,
    requireEmail: false,
    defaultRoleId: '',
    defaultRoleName: '',
    defaultDeptId: '',
    defaultDeptName: '',
  }
}

const rolePickerRef = ref<InstanceType<typeof ModalGrantRole> | null>(null)
const deptPickerRef = ref<InstanceType<typeof ModalGrantDept> | null>(null)

const state = reactive({
  loading: false,
  saving: false,
  subTab: DEFAULT_ACCOUNT_TYPE as AccountType,
  byType: createAccountTypeMap(emptyScope),
  snapshot: '',
})

onMounted(() => {
  void reload()
})

async function resolveRoleName(id: string): Promise<string> {
  if (!id) return ''
  try {
    const res = await roleApi.detail({ id })
    return res?.data?.name || res?.data?.code || id
  } catch {
    return id
  }
}

async function resolveDeptName(id: string): Promise<string> {
  if (!id) return ''
  try {
    const res = await deptApi.detail({ id })
    return res?.data?.name || id
  } catch {
    return id
  }
}

async function fillScope(map: Record<string, string>, type: AccountType) {
  const target = state.byType[type]
  target.enabled = parseBool(map[accountConfigKey(PREFIX, type, 'ENABLED')])
  target.requirePhone = parseBool(map[accountConfigKey(PREFIX, type, 'REQUIRE_PHONE')])
  target.requireEmail = parseBool(map[accountConfigKey(PREFIX, type, 'REQUIRE_EMAIL')])
  target.defaultRoleId = map[accountConfigKey(PREFIX, type, 'DEFAULT_ROLE_ID')] || ''
  target.defaultDeptId = map[accountConfigKey(PREFIX, type, 'DEFAULT_DEPT_ID')] || ''
  const [roleName, deptName] = await Promise.all([
    resolveRoleName(target.defaultRoleId),
    resolveDeptName(target.defaultDeptId),
  ])
  target.defaultRoleName = roleName
  target.defaultDeptName = deptName
}

async function reload() {
  state.loading = true
  try {
    const map = await loadByCategory(CATEGORY)
    await Promise.all(mapAccountTypes((type) => fillScope(map, type)))
    state.snapshot = JSON.stringify(state.byType)
  } finally {
    state.loading = false
  }
}

function reset() {
  if (!state.snapshot) return
  const data = JSON.parse(state.snapshot) as Record<AccountType, ScopeForm>
  for (const type of Object.keys(data) as AccountType[]) {
    if (state.byType[type]) Object.assign(state.byType[type], data[type])
  }
}

const current = computed(() => state.byType[state.subTab])

function scopeItems(type: AccountType, form: ScopeForm) {
  return [
    {
      config_key: accountConfigKey(PREFIX, type, 'ENABLED'),
      config_value: toBoolStr(form.enabled),
      category: CATEGORY,
    },
    {
      config_key: accountConfigKey(PREFIX, type, 'REQUIRE_PHONE'),
      config_value: toBoolStr(form.requirePhone),
      category: CATEGORY,
    },
    {
      config_key: accountConfigKey(PREFIX, type, 'REQUIRE_EMAIL'),
      config_value: toBoolStr(form.requireEmail),
      category: CATEGORY,
    },
    {
      config_key: accountConfigKey(PREFIX, type, 'DEFAULT_ROLE_ID'),
      config_value: form.defaultRoleId,
      category: CATEGORY,
    },
    {
      config_key: accountConfigKey(PREFIX, type, 'DEFAULT_DEPT_ID'),
      config_value: form.defaultDeptId,
      category: CATEGORY,
    },
  ]
}

async function save() {
  state.saving = true
  try {
    await saveByKeys(mapAccountTypes((type) => scopeItems(type, state.byType[type])).flat())
    window.$message.success('保存成功')
    state.snapshot = JSON.stringify(state.byType)
  } finally {
    state.saving = false
  }
}

function clearRole() {
  current.value.defaultRoleId = ''
  current.value.defaultRoleName = ''
}

function clearDept() {
  current.value.defaultDeptId = ''
  current.value.defaultDeptName = ''
}

function openRolePicker() {
  void rolePickerRef.value?.openPicker({
    selectedIds: current.value.defaultRoleId ? [current.value.defaultRoleId] : [],
    multiple: false,
    title: '选择默认角色',
  })
}

function openDeptPicker() {
  void deptPickerRef.value?.openPicker({
    selectedIds: current.value.defaultDeptId ? [current.value.defaultDeptId] : [],
    multiple: false,
    title: '选择默认部门',
  })
}

function onRoleConfirm(items: Array<{ id: string; name: string }>) {
  const first = items[0]
  current.value.defaultRoleId = first?.id ?? ''
  current.value.defaultRoleName = first?.name ?? ''
}

function onDeptConfirm(items: Array<{ id: string; name: string }>) {
  const first = items[0]
  current.value.defaultDeptId = first?.id ?? ''
  current.value.defaultDeptName = first?.name ?? ''
}
</script>

<template>
  <NSpin :show="state.loading">
    <NTabs
      v-model:value="state.subTab"
      type="line"
      class="sys-config-subnav"
    >
      <NTab
        v-for="item in ACCOUNT_TYPE_TABS"
        :key="item.key"
        :name="item.key"
        :tab="item.label"
      />
    </NTabs>

    <ConfigSectionLayout
      :description="`配置「${accountTypeLabel(state.subTab)}」账户类型的注册策略。`"
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NForm
        class="sys-config-form"
        label-placement="top"
      >
        <NFormItem label="允许注册">
          <NSwitch v-model:value="current.enabled" />
        </NFormItem>
        <NFormItem label="注册后需绑定手机号">
          <NSwitch v-model:value="current.requirePhone" />
        </NFormItem>
        <NFormItem label="注册后需绑定邮箱">
          <NSwitch v-model:value="current.requireEmail" />
        </NFormItem>
        <NFormItem label="默认角色">
          <NInput
            :value="current.defaultRoleName"
            readonly
            clearable
            placeholder="可选，点击选择角色"
            @clear="clearRole"
          >
            <template #suffix>
              <NButton
                text
                type="primary"
                size="small"
                @click="openRolePicker"
              >
                选择
              </NButton>
            </template>
          </NInput>
        </NFormItem>
        <NFormItem label="默认部门">
          <NInput
            :value="current.defaultDeptName"
            readonly
            clearable
            placeholder="可选，点击选择部门"
            @clear="clearDept"
          >
            <template #suffix>
              <NButton
                text
                type="primary"
                size="small"
                @click="openDeptPicker"
              >
                选择
              </NButton>
            </template>
          </NInput>
        </NFormItem>
      </NForm>
    </ConfigSectionLayout>

    <ModalGrantRole
      ref="rolePickerRef"
      @confirm="onRoleConfirm"
    />
    <ModalGrantDept
      ref="deptPickerRef"
      @confirm="onDeptConfirm"
    />
  </NSpin>
</template>
