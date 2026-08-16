<!--
  Author: Charlie

  任务执行记录页（按 job_id 分页查询）。
-->
<script setup lang="tsx">
import type { PaginationProps } from 'naive-ui'
import type { ProDataTableColumns } from 'pro-naive-ui'
import { Icon } from '@iconify/vue/offline'
import { NButton, NFlex, NIcon, NTag } from 'naive-ui'
import { jobApi } from '@/api'
import { formatDateTime } from '@/utils'
import { ProDataTable } from 'pro-naive-ui'
import { computed, onMounted, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { readPageMeta } from '@/utils/wire'

const route = useRoute()
const router = useRouter()
const listPath = '/sys/job'

const state = reactive({
  jobName: '',
  logs: [] as any[],
  total: 0,
  loading: false,
  page: 1,
  pageSize: 20,
})

const jobId = computed(() => {
  const id = route.query.id
  return typeof id === 'string' ? id : ''
})

const tableTitle = computed(() => (state.jobName ? `执行记录 · ${state.jobName}` : '执行记录'))

const pagination = computed<PaginationProps>(() => ({
  page: state.page,
  pageSize: state.pageSize,
  itemCount: state.total,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  prefix: ({ itemCount }) => `${itemCount} 条`,
  onUpdatePage: (value) => {
    state.page = value
    fetchLogs()
  },
  onUpdatePageSize: (value) => {
    state.pageSize = value
    state.page = 1
    fetchLogs()
  },
}))

const tableColumns = computed<ProDataTableColumns<any>>(() => [
  {
    title: '执行时间',
    path: 'execute_time',
    width: 190,
    ellipsis: { tooltip: true },
    render: (row) => formatDateTime(row.execute_time),
  },
  {
    title: '执行人',
    path: 'executor',
    width: 100,
    ellipsis: { tooltip: true },
    render: (row) => row.executor || '-',
  },
  {
    title: '执行用时',
    path: 'execute_duration_ms',
    width: 110,
    render: (row) => (row.execute_duration_ms == null ? '-' : `${row.execute_duration_ms} ms`),
  },
  {
    title: '结果',
    path: 'success',
    width: 90,
    render: (row) =>
      row.success ? (
        <NTag size="small" type="success" bordered={false}>
          成功
        </NTag>
      ) : (
        <NTag size="small" type="error" bordered={false}>
          失败
        </NTag>
      ),
  },
  {
    title: '执行结果',
    path: 'execute_result',
    width: 260,
    ellipsis: { tooltip: true },
  },
  {
    title: 'IP',
    path: 'ip',
    width: 140,
    ellipsis: { tooltip: true },
  },
  {
    title: '进程ID',
    path: 'process_id',
    width: 110,
  },
  {
    title: '程序目录',
    path: 'app_dir',
    width: 240,
    ellipsis: { tooltip: true },
  },
])

async function fetchJobName(id: string) {
  if (!id) {
    state.jobName = ''
    return
  }
  try {
    const response = await jobApi.detail({ id })
    state.jobName = response.data?.job_name ?? ''
  } catch {
    state.jobName = ''
  }
}

async function fetchLogs() {
  if (!jobId.value) {
    state.logs = []
    state.total = 0
    return
  }
  state.loading = true
  try {
    const response = await jobApi.logPage({
      current: state.page,
      size: state.pageSize,
      job_id: jobId.value,
    })
    const data = response.data ?? {}
    state.logs = data.records ?? []
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
  } finally {
    state.loading = false
  }
}

function goBack() {
  router.push(listPath)
}

onMounted(() => {
  void fetchJobName(jobId.value)
  void fetchLogs()
})
watch(jobId, (id) => {
  state.page = 1
  void fetchJobName(id)
  void fetchLogs()
})
</script>

<template>
  <NFlex
    class="h-full min-h-0"
    vertical
  >
    <ProDataTable
      class="min-h-0 flex-1"
      remote
      :title="tableTitle"
      row-key="id"
      :scroll-x="1300"
      :columns="tableColumns"
      :data="state.logs"
      :loading="state.loading"
      :pagination="pagination"
    >
      <template #toolbar>
        <NFlex align="center">
          <NButton
            text
            title="返回"
            aria-label="返回"
            @click="goBack"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:return" />
              </NIcon>
            </template>
          </NButton>
          <NButton
            text
            title="刷新"
            aria-label="刷新"
            :loading="state.loading"
            @click="fetchLogs"
          >
            <template #icon>
              <NIcon>
                <Icon icon="icon-park-outline:reload" />
              </NIcon>
            </template>
          </NButton>
        </NFlex>
      </template>
    </ProDataTable>
  </NFlex>
</template>
