<!-- Author: Charlie -->

<script setup lang="ts">
import { messageApi } from '@/api'
import MessageDetailModal from '@/components/message/MessageDetailModal.vue'
import { useMessageUnreadStore } from '@/stores'
import { dictTypeData, formatDateTime, wireBool } from '@/utils'
import { onMounted, reactive, ref } from 'vue'
import { readPageMeta } from '@/utils/wire'

const detailModalRef = ref<InstanceType<typeof MessageDetailModal> | null>(null)
const unreadStore = useMessageUnreadStore()
const state = reactive({
  rows: [] as any[],
  total: 0,
  loading: false,
  page: 1,
  pageSize: 10,
})

onMounted(() => {
  void fetchPage()
})

async function fetchPage() {
  state.loading = true
  try {
    const response = await messageApi.myPage({
      current: state.page,
      size: state.pageSize,
    })
    const data = response.data ?? {}
    state.rows = (data.records ?? []).map((row: any) => ({
      ...row,
      is_read: wireBool(row.is_read ?? false),
    }))
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
  } finally {
    state.loading = false
  }
}

function kindLabel(row: any) {
  return row.kind === 'ANNOUNCEMENT' ? '公告' : '通知'
}

function severityLabel(row: any) {
  if (!row.severity) return ''
  return dictTypeData('NOTIFICATION_SEVERITY', row.severity) || row.severity
}

function excerpt(row: any) {
  const text = String(row.content || '')
    .replace(/\s+/g, ' ')
    .trim()
  if (!text) return ''
  return text.length > 96 ? `${text.slice(0, 96)}…` : text
}

async function openDetail(row: any) {
  await detailModalRef.value?.open({
    id: row.id,
    sourceType: row.kind === 'ANNOUNCEMENT' ? 'ANNOUNCEMENT' : 'NOTIFICATION',
    title: row.title,
    is_read: row.is_read,
    publish_at: row.publish_at,
    content: row.content,
    severity: row.severity,
  })
}

async function handleDetailChanged(payload: { type: string; id: string }) {
  const row = state.rows.find((item) => item.id === payload.id)
  if (row) row.is_read = true
}

async function markAllRead() {
  await messageApi.readAll()
  state.rows.forEach((row) => {
    row.is_read = true
  })
  unreadStore.notifyReadAll()
  window.$message.success('已全部标记为已读')
}

function handlePageChange(page: number) {
  state.page = page
  void fetchPage()
}

function handlePageSizeChange(size: number) {
  state.pageSize = size
  state.page = 1
  void fetchPage()
}
</script>

<template>
  <NSpace
    vertical
    :size="12"
    class="w-full min-w-0"
  >
    <div class="flex justify-end gap-1">
      <NButton
        text
        :loading="state.loading"
        @click="fetchPage"
      >
        刷新
      </NButton>
      <NButton
        text
        @click="markAllRead"
      >
        全部已读
      </NButton>
    </div>

    <NSpin
      :show="state.loading"
      class="w-full min-w-0"
    >
      <NEmpty
        v-if="!state.loading && !state.rows.length"
        description="暂无消息"
      />
      <NList
        v-else
        clickable
        show-divider
        class="w-full min-w-0"
      >
        <NListItem
          v-for="row in state.rows"
          :key="row.id"
          @click="openDetail(row)"
        >
          <NThing class="w-full min-w-0">
            <template #avatar>
              <NBadge
                :dot="!row.is_read"
                type="info"
              >
                <NAvatar
                  round
                  :size="36"
                >
                  <NovaIcon
                    :icon="
                      row.kind === 'ANNOUNCEMENT'
                        ? 'icon-park-outline:volume-notice'
                        : 'icon-park-outline:tips-one'
                    "
                    :size="18"
                    :style="{ color: row.is_read ? 'var(--text-color-3)' : 'var(--primary-color)' }"
                  />
                </NAvatar>
              </NBadge>
            </template>
            <template #header>
              <NSpace
                :size="8"
                align="center"
                :wrap="true"
              >
                <NTag
                  size="small"
                  :bordered="false"
                  :type="row.kind === 'ANNOUNCEMENT' ? 'warning' : 'info'"
                >
                  {{ kindLabel(row) }}
                </NTag>
                <NTag
                  v-if="severityLabel(row)"
                  size="small"
                  :bordered="false"
                >
                  {{ severityLabel(row) }}
                </NTag>
                <NText
                  :depth="row.is_read ? 3 : 1"
                  :strong="!row.is_read"
                >
                  {{ row.title }}
                </NText>
              </NSpace>
            </template>
            <template #header-extra>
              <NText
                depth="3"
                style="font-size: 12px; white-space: nowrap"
              >
                {{ formatDateTime(row.publish_at || row.created_at) }}
              </NText>
            </template>
            <template
              v-if="excerpt(row)"
              #description
            >
              <NEllipsis
                :line-clamp="2"
                :tooltip="false"
              >
                <NText
                  depth="3"
                  style="font-size: 13px"
                >
                  {{ excerpt(row) }}
                </NText>
              </NEllipsis>
            </template>
          </NThing>
        </NListItem>
      </NList>
    </NSpin>

    <NSpace
      v-if="state.total > 0"
      justify="end"
    >
      <NPagination
        :page="state.page"
        :page-size="state.pageSize"
        :item-count="state.total"
        :page-sizes="[10, 20, 30]"
        show-size-picker
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
      />
    </NSpace>

    <MessageDetailModal
      ref="detailModalRef"
      @changed="handleDetailChanged"
    />
  </NSpace>
</template>
