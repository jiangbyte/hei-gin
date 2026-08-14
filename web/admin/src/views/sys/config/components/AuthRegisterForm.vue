<!-- Author: Charlie -->

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { deptApi, roleApi } from '@/api'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import ModalGrantDept from '@/views/iam/components/ModalGrantDept.vue'
import ModalGrantRole from '@/views/iam/components/ModalGrantRole.vue'
import {
  ACCOUNT_TYPE_TABS,
  accountConfigKey,
  accountTypeLabel,
  createAccountTypeMap,
  mapAccountTypes,
  type AccountType,
} from '@/constants/account'
import { loadByCategory, parseBool, saveByKeys, toBoolStr } from '../composables/useConfigForm'

const CATEGORY = 'AUTH_REGISTER'
const PREFIX = 'AUTH_REGISTER'
const FORCE_CATEGORY = 'AUTH_FORCE_BIND'
const FORCE_PREFIX = 'AUTH_FORCE_BIND'

type ScopeForm = {
  enabled: boolean
  allowAccount: boolean
  allowEmail: boolean
  allowPhone: boolean
  forceBindEmail: boolean
  forceBindPhone: boolean
  defaultRoleId: string
  defaultRoleName: string
  defaultDeptId: string
  defaultDeptName: string
}

function emptyScope(): ScopeForm {
  return {
    enabled: false,
    allowAccount: true,
    allowEmail: true,
    allowPhone: false,
    forceBindEmail: false,
    forceBindPhone: false,
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
  subTab: 'PORTAL' as AccountType,
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

async function fillScope(registerMap: Record<string, string>, forceMap: Record<string, string>, type: AccountType) {
  const target = state.byType[type]
  target.enabled = parseBool(registerMap[accountConfigKey(PREFIX, type, 'ENABLED')])
  target.allowAccount = parseBool(registerMap[accountConfigKey(PREFIX, type, 'ALLOW_ACCOUNT')] ?? 'TRUE')
  target.allowEmail = parseBool(registerMap[accountConfigKey(PREFIX, type, 'ALLOW_EMAIL')] ?? 'TRUE')
  target.allowPhone = parseBool(registerMap[accountConfigKey(PREFIX, type, 'ALLOW_PHONE')] ?? 'FALSE')
  target.forceBindEmail = parseBool(forceMap[accountConfigKey(FORCE_PREFIX, type, 'EMAIL')])
  target.forceBindPhone = parseBool(forceMap[accountConfigKey(FORCE_PREFIX, type, 'PHONE')])
  target.defaultRoleId = registerMap[accountConfigKey(PREFIX, type, 'DEFAULT_ROLE_ID')] || ''
  target.defaultDeptId = registerMap[accountConfigKey(PREFIX, type, 'DEFAULT_DEPT_ID')] || ''
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
    const [registerMap, forceMap] = await Promise.all([
      loadByCategory(CATEGORY),
      loadByCategory(FORCE_CATEGORY),
    ])
    await Promise.all(mapAccountTypes((type) => fillScope(registerMap, forceMap, type)))
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
const isPortal = computed(() => state.subTab === 'PORTAL')

function saveItems(type: AccountType, form: ScopeForm) {
  const items = [
    {
      config_key: accountConfigKey(FORCE_PREFIX, type, 'EMAIL'),
      config_value: toBoolStr(form.forceBindEmail),
      category: FORCE_CATEGORY,
    },
    {
      config_key: accountConfigKey(FORCE_PREFIX, type, 'PHONE'),
      config_value: toBoolStr(form.forceBindPhone),
      category: FORCE_CATEGORY,
    },
  ]
  if (type === 'PORTAL') {
    items.unshift(
      {
        config_key: accountConfigKey(PREFIX, 'PORTAL', 'ENABLED'),
        config_value: toBoolStr(form.enabled),
        category: CATEGORY,
      },
      {
        config_key: accountConfigKey(PREFIX, 'PORTAL', 'ALLOW_ACCOUNT'),
        config_value: toBoolStr(form.allowAccount),
        category: CATEGORY,
      },
      {
        config_key: accountConfigKey(PREFIX, 'PORTAL', 'ALLOW_EMAIL'),
        config_value: toBoolStr(form.allowEmail),
        category: CATEGORY,
      },
      {
        config_key: accountConfigKey(PREFIX, 'PORTAL', 'ALLOW_PHONE'),
        config_value: toBoolStr(form.allowPhone),
        category: CATEGORY,
      },
      {
        config_key: accountConfigKey(PREFIX, 'PORTAL', 'DEFAULT_ROLE_ID'),
        config_value: form.defaultRoleId,
        category: CATEGORY,
      },
      {
        config_key: accountConfigKey(PREFIX, 'PORTAL', 'DEFAULT_DEPT_ID'),
        config_value: form.defaultDeptId,
        category: CATEGORY,
      },
    )
  }
  return items
}

async function save() {
  state.saving = true
  try {
    await saveByKeys(saveItems(state.subTab, state.byType[state.subTab]))
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
      :description="
        isPortal
          ? '门户自助注册通道与强制绑定；管理员仅配置强制绑定。'
          : `「${accountTypeLabel(state.subTab)}」无自助注册，可配置登录后强制绑定邮箱/手机。`
      "
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NAlert
        v-if="!isPortal"
        type="info"
        :bordered="false"
        title="管理员账户仅支持后台开户"
        class="mb-12px"
      >
        不提供自助注册；请在「账号管理」中创建管理员。下方可配置强制绑定。
      </NAlert>

      <NForm
        class="sys-config-form"
        label-placement="top"
      >
        <template v-if="isPortal">
          <NFormItem label="允许门户注册">
            <NSwitch v-model:value="current.enabled" />
          </NFormItem>
          <NFormItem label="允许用户名注册">
            <NSwitch v-model:value="current.allowAccount" />
          </NFormItem>
          <NFormItem label="允许邮箱注册">
            <NSwitch v-model:value="current.allowEmail" />
          </NFormItem>
          <NFormItem label="允许手机注册">
            <NSwitch v-model:value="current.allowPhone" />
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
        </template>

        <NDivider title-placement="left">
          强制绑定
        </NDivider>
        <NAlert
          type="warning"
          :bordered="false"
          class="mb-12px"
        >
          开启后，未绑定对应身份的用户登录将被硬拦截至用户中心绑定页（优先级低于密码过期）。
        </NAlert>
        <NFormItem label="强制绑定邮箱">
          <NSwitch v-model:value="current.forceBindEmail" />
        </NFormItem>
        <NFormItem label="强制绑定手机">
          <NSwitch v-model:value="current.forceBindPhone" />
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

<style scoped>
.mb-12px {
  margin-bottom: 12px;
}
</style>
