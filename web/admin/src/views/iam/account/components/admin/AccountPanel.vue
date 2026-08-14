<!-- Author: Charlie -->

<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns, ProSearchFormColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { accountApi } from '@/api'
import { createTagColor, formatDateTime, hasPermission, renderButtonIcon } from '@/utils'
import { NAvatar, NButton, NDropdown, NFlex, NIcon, NTag } from 'naive-ui'
import { createProSearchForm, ProCard, ProDataTable, ProSearchForm } from 'pro-naive-ui'
import { computed, ref } from 'vue'
import { dictList, dictTypeColor, dictTypeData } from '@/utils/dict'
import { renderAccountAvatar, useAccountList } from '../../composables/useAccountList'
import ModalDetail from './ModalDetail.vue'
import ModalForm from './ModalForm.vue'
import ModalGrantDept from '../../../components/ModalGrantDept.vue'
import ModalGrantClientResource from '../../../role/components/ModalGrantClientResource.vue'
import ModalGrantResource from '../../../role/components/ModalGrantResource.vue'
import ModalGrantGroup from '../../../components/ModalGrantGroup.vue'
import ModalGrantRole from '../../../components/ModalGrantRole.vue'

const formRef = ref<{ openModal: (id?: string) => void } | null>(null)
const detailRef = ref<{ openModal: (id?: string) => void } | null>(null)
const grantRoleModalRef = ref<any>(null)
const grantGroupModalRef = ref<any>(null)
const grantDeptModalRef = ref<any>(null)
const grantResourceModalRef = ref<any>(null)
const grantClientResourceModalRef = ref<any>(null)

const {
  state,
  hasCheckedRows,
  applySearch,
  resetSearch,
  fetchPage,
  handleCheckedRowKeys,
  confirmDelete,
} = useAccountList('ADMIN')

const searchForm = createProSearchForm<any>({
  defaultCollapsed: true,
  onSubmit: applySearch,
  onReset: resetSearch,
})

/** 管理员列表搜索：对接 AccountAdminPageQuery */
const searchColumns = computed<ProSearchFormColumns<any>>(() => [
  { title: '账号', path: 'account', field: 'input' },
  { title: '姓名', path: 'name', field: 'input' },
  { title: '手机号', path: 'phone', field: 'input' },
  { title: '邮箱', path: 'email', field: 'input' },
  {
    title: '账号状态',
    path: 'account_status',
    field: 'select',
    fieldProps: { options: dictList('ACCOUNT_STATUS') },
  },
])

const pagination = computed<PaginationProps>(() => ({
  page: state.page,
  pageSize: state.pageSize,
  itemCount: state.total,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: ({ itemCount }) => `${itemCount} 条`,
  onUpdatePage: (value) => {
    state.page = value
    void fetchPage()
  },
  onUpdatePageSize: (value) => {
    state.pageSize = value
    state.page = 1
    void fetchPage()
  },
}))

const grantOptions = computed(() =>
  [
    { label: '分配角色', key: 'role', permission: 'iam:account:grantrole' },
    { label: '分配用户组', key: 'group', permission: 'iam:account:grantgroup' },
    { label: '分配部门', key: 'dept', permission: 'iam:account:grantdept' },
    { label: '分配资源', key: 'resource', permission: 'iam:account:grantresource' },
    {
      label: '分配客户端资源',
      key: 'client-resource',
      permission: 'iam:account:grantclientresource',
    },
  ].filter((item) => hasPermission(item.permission)),
)

const avatarImgProps = { referrerPolicy: 'no-referrer' } as any

/** 管理员表格列：含 admin_user_profile.remark */
const tableColumns = computed<ProDataTableColumns<any>>(() => [
  { type: 'selection', fixed: 'left' },
  {
    title: '头像',
    key: 'avatar',
    width: 80,
    render: (row) => {
      const { avatar, name } = renderAccountAvatar(row)
      return (
        <NAvatar round size={32} src={avatar} imgProps={avatarImgProps}>
          {avatar
            ? undefined
            : String(name || '-')
                .slice(0, 1)
                .toUpperCase()}
        </NAvatar>
      )
    },
  },
  { title: '账号', path: 'account', width: 140, ellipsis: { tooltip: true } },
  { title: '姓名', path: 'name', width: 130, ellipsis: { tooltip: true } },
  { title: '昵称', path: 'nickname', width: 130, ellipsis: { tooltip: true } },
  {
    title: '账号状态',
    path: 'account_status',
    width: 120,
    render: (row) => (
      <NTag
        size="small"
        color={createTagColor(dictTypeColor('ACCOUNT_STATUS', row.account_status))}
        bordered={false}
      >
        {dictTypeData('ACCOUNT_STATUS', row.account_status)}
      </NTag>
    ),
  },
  { title: '手机号', path: 'phone', width: 150 },
  { title: '邮箱', path: 'email', width: 220, ellipsis: { tooltip: true } },
  {
    title: '三方',
    key: 'oauth_bindings',
    width: 160,
    render: (row) => {
      const list = Array.isArray(row.oauth_bindings) ? row.oauth_bindings : []
      if (!list.length) return '-'
      return (
        <NFlex size={4} wrap>
          {list.map((item: any) => (
            <NTag key={item.provider} size="small" bordered={false}>
              {dictTypeData('OAUTH_PROVIDER', item.provider) || item.provider}
            </NTag>
          ))}
        </NFlex>
      )
    },
  },
  { title: '备注', path: 'remark', width: 180, ellipsis: { tooltip: true } },
  {
    title: '最近登录时间',
    path: 'latest_login_time',
    width: 190,
    render: (row) => formatDateTime(row.latest_login_time),
  },
  {
    title: '更新时间',
    path: 'updated_at',
    width: 190,
    render: (row) => formatDateTime(row.updated_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    fixed: 'right',
    render: (row) => {
      const options = grantOptions.value
      return (
        <NFlex size={12}>
          {hasPermission('iam:account:detail') ? (
            <NButton
              type="info"
              size="small"
              text
              onClick={() => detailRef.value?.openModal(row.id)}
            >
              {renderButtonIcon('icon-park-outline:preview-open')}
            </NButton>
          ) : null}
          {hasPermission('iam:account:update') ? (
            <NButton
              type="primary"
              size="small"
              text
              onClick={() => formRef.value?.openModal(row.id)}
            >
              {renderButtonIcon('icon-park-outline:edit')}
            </NButton>
          ) : null}
          {options.length ? (
            <NDropdown
              trigger="click"
              options={options}
              onSelect={(key) => openGrantModal(String(key), row)}
            >
              <NButton type="warning" size="small" text>
                {renderButtonIcon('icon-park-outline:permissions')}
              </NButton>
            </NDropdown>
          ) : null}
          {hasPermission('iam:account:delete') ? (
            <NButton type="error" size="small" text onClick={() => confirmDelete(row.id)}>
              {renderButtonIcon('icon-park-outline:delete')}
            </NButton>
          ) : null}
        </NFlex>
      )
    },
  },
])

function openGrantModal(type: string, row: any) {
  const account = {
    id: row.id,
    code: row.account,
    name: row.nickname || row.name || '-',
  }
  if (type === 'role') grantRoleModalRef.value?.openModal(account)
  else if (type === 'group') grantGroupModalRef.value?.openModal(account)
  else if (type === 'dept') grantDeptModalRef.value?.openModal(account)
  else if (type === 'resource') {
    grantResourceModalRef.value?.openModal(account, accountApi, '分配资源', {
      accountType: 'ADMIN',
      lockAccountType: true,
    })
  } else if (type === 'client-resource') {
    grantClientResourceModalRef.value?.openModal(account, accountApi, '分配客户端资源', {
      accountType: 'ADMIN',
      lockAccountType: true,
    })
  }
}
</script>

<template>
  <NFlex
    class="h-full min-h-0"
    vertical
  >
    <ProCard content-class="pb-0!">
      <ProSearchForm
        :form="searchForm"
        :columns="searchColumns"
        :reset-button-props="{ content: '重置' }"
        :search-button-props="{ content: '搜索' }"
        :collapse-button-props="{
          content: searchForm.collapsed.value ? '展开' : '收起',
        }"
      />
    </ProCard>

    <ProDataTable
      class="min-h-0 flex-1"
      remote
      title="管理员管理"
      row-key="id"
      :scroll-x="1980"
      :columns="tableColumns"
      :data="state.accounts"
      :loading="state.loading"
      :pagination="pagination"
      :checked-row-keys="state.checkedRowKeys"
      :on-update-checked-row-keys="handleCheckedRowKeys"
    >
      <template #toolbar>
        <NFlex>
          <NButton
            v-if="hasPermission('iam:account:create')"
            type="primary"
            text
            title="新增"
            aria-label="新增"
            @click="formRef?.openModal()"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:plus" />
              </NIcon>
            </template>
          </NButton>
          <NButton
            text
            title="刷新"
            aria-label="刷新"
            :loading="state.loading"
            @click="fetchPage"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:reload" />
              </NIcon>
            </template>
          </NButton>
          <NButton
            v-if="hasPermission('iam:account:delete')"
            type="error"
            text
            title="批量删除"
            aria-label="批量删除"
            :disabled="!hasCheckedRows"
            @click="confirmDelete(state.checkedRowKeys)"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:delete" />
              </NIcon>
            </template>
          </NButton>
        </NFlex>
      </template>
    </ProDataTable>

    <ModalForm
      ref="formRef"
      @saved="fetchPage"
    />
    <ModalDetail ref="detailRef" />
    <ModalGrantRole
      ref="grantRoleModalRef"
      @saved="fetchPage"
    />
    <ModalGrantGroup
      ref="grantGroupModalRef"
      @saved="fetchPage"
    />
    <ModalGrantDept
      ref="grantDeptModalRef"
      @saved="fetchPage"
    />
    <ModalGrantResource
      ref="grantResourceModalRef"
      @saved="fetchPage"
    />
    <ModalGrantClientResource
      ref="grantClientResourceModalRef"
      @saved="fetchPage"
    />
  </NFlex>
</template>
