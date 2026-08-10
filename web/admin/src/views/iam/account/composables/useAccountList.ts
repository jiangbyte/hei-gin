/** Author: Charlie */

import { accountApi } from '@/api'
import type { AccountType } from '@/constants/account'
import { normalizeSearchValues } from '@/utils'
import { readPageMeta } from '@/utils/wire'
import { computed, onMounted, reactive } from 'vue'

/** 列表页公共分页 / 拉取 / 删除逻辑；各体系面板自行定义搜索与表格列。 */
export function useAccountList(accountType: AccountType) {
  const state = reactive({
    accounts: [] as any[],
    total: 0,
    loading: false,
    searchValues: {} as Record<string, unknown>,
    checkedRowKeys: [] as string[],
    page: 1,
    pageSize: 20,
  })

  const hasCheckedRows = computed(() => state.checkedRowKeys.length > 0)

  onMounted(() => {
    void fetchPage()
  })

  function applySearch(values: Record<string, unknown>) {
    state.searchValues = normalizeSearchValues(values, {
      account: (value) => String(value).trim(),
      name: (value) => String(value).trim(),
      phone: (value) => String(value).trim(),
      email: (value) => String(value).trim(),
    })
    state.page = 1
    void fetchPage()
  }

  function resetSearch() {
    state.searchValues = {}
    state.page = 1
    void fetchPage()
  }

  async function fetchPage() {
    state.loading = true
    try {
      const response = await accountApi.page({
        current: state.page,
        size: state.pageSize,
        account_type: accountType,
        ...state.searchValues,
      })
      const data = response.data ?? {}
      state.accounts = data.records ?? []
      const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
      state.total = pageMeta.total
      state.page = pageMeta.current
      state.pageSize = pageMeta.size
      state.checkedRowKeys = state.checkedRowKeys.filter((key) =>
        state.accounts.some((item) => item.id === key),
      )
    } finally {
      state.loading = false
    }
  }

  function handleCheckedRowKeys(keys: Array<string | number>) {
    state.checkedRowKeys = keys.map(String)
  }

  function confirmDelete(value: string | string[]) {
    const ids = Array.isArray(value) ? value : [value]
    if (!ids.length) return
    const isBatch = ids.length > 1
    window.$dialog.warning({
      title: isBatch ? '批量删除' : '删除',
      draggable: true,
      maskClosable: false,
      content: isBatch ? `删除 ${ids.length} 个账号?` : '删除该账号?',
      positiveText: '确认',
      negativeText: '取消',
      onPositiveClick: () => deleteData(ids),
    })
  }

  async function deleteData(ids: string[]) {
    await accountApi.remove({ ids })
    state.checkedRowKeys = state.checkedRowKeys.filter((key) => !ids.includes(key))
    window.$message.success('删除成功')
    await fetchPage()
    if (!state.accounts.length && state.total > 0 && state.page > 1) {
      state.page -= 1
      await fetchPage()
    }
  }

  return {
    state,
    hasCheckedRows,
    applySearch,
    resetSearch,
    fetchPage,
    handleCheckedRowKeys,
    confirmDelete,
  }
}

export function renderAccountAvatar(row: any) {
  const avatar = row.avatar || undefined
  const name = row.nickname || row.name || row.account || ''
  return { avatar, name }
}
