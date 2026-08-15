<!-- Author: Charlie -->

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { myNoticeApi } from '@/api'
import MessageDetailModal from '@/components/sys/MessageDetailModal.vue'
import { useMessageUnreadStore } from '@/stores'
import { formatDateTime, wireBool } from '@/utils'
import NoticeList, { type BannerItem } from '../common/NoticeList.vue'
import { readPageMeta } from '@/utils/wire'

const pageSize = 8
const router = useRouter()

type LoadMode = 'replace' | 'merge' | 'append'
type NoticeKind = 'NOTIFICATION' | 'ANNOUNCEMENT'

interface NoticeSource {
  id: string
  title: string
  icon: string
  tagTitle?: string
  tagType?: BannerItem['tagType']
  description?: string
  date: string
  sourceType: NoticeKind
  sourceId: string
  isRead: boolean
}

const listState = reactive({
  records: [] as NoticeSource[],
  current: 0,
  size: pageSize,
  total: 0,
  loading: false,
  loaded: false,
})
const unreadStore = useMessageUnreadStore()
const showPopover = ref(false)
const detailModalRef = ref<InstanceType<typeof MessageDetailModal> | null>(null)

const list = computed(() => listState.records.map(toNoticeItem))
const hasMore = computed(() => listState.records.length < listState.total)
const unreadTotal = computed(() => unreadStore.unreadTotal)

onMounted(() => {
  void unreadStore.refresh()
})

async function refresh() {
  await Promise.all([unreadStore.refresh(), loadList(1, listState.loaded ? 'merge' : 'replace')])
}

async function loadMore() {
  if (listState.loading || listState.records.length >= listState.total) return
  await loadList(listState.current + 1, 'append')
}

async function loadList(page = 1, mode: LoadMode = 'replace') {
  if (listState.loading) return
  listState.loading = true
  try {
    const response = await myNoticeApi.myPage({ current: page, size: listState.size })
    const data = response.data ?? {}
    const incoming = (data.records ?? []).map((item: any) => mapHistoryItem(item))
    listState.records = mergeNoticeRecords(listState.records, incoming, mode)
    const pageMeta = readPageMeta(data, { current: page, size: listState.size })
    listState.total = pageMeta.total || listState.records.length
    listState.current = pageMeta.current
    listState.size = pageMeta.size
    listState.loaded = true
  } finally {
    listState.loading = false
  }
}

async function handleOpen(id: string) {
  const item = listState.records.find((notice) => notice.id === id)
  if (!item) return
  await detailModalRef.value?.open({
    ...item,
    id: item.sourceId,
    sourceType: item.sourceType,
    is_read: item.isRead,
  })
}

function handleDetailChanged(payload: { type: string; id: string }) {
  const item = listState.records.find((notice) => notice.id === `${payload.type}:${payload.id}`)
  if (item) item.isRead = true
}

async function markAllRead() {
  try {
    await myNoticeApi.readAll()
    listState.records.forEach((item) => {
      item.isRead = true
    })
    unreadStore.notifyReadAll()
  } catch {
    /* 忽略 */
  }
}

function goMore() {
  showPopover.value = false
  void router.push({ path: '/profile', query: { tab: 'my_messages' } })
}

function mergeNoticeRecords(
  current: NoticeSource[],
  incoming: NoticeSource[],
  mode: LoadMode,
): NoticeSource[] {
  if (mode === 'replace') return incoming
  const currentMap = new Map(current.map((item) => [item.id, item]))
  const result = current.map((item) => ({
    ...item,
    ...(incoming.find((i) => i.id === item.id) ?? {}),
  }))
  incoming.forEach((item) => {
    if (!currentMap.has(item.id)) result.push(item)
  })
  return result
}

function mapHistoryItem(item: any): NoticeSource {
  const kind: NoticeKind = item.kind === 'ANNOUNCEMENT' ? 'ANNOUNCEMENT' : 'NOTIFICATION'
  return {
    id: `${kind}:${item.id}`,
    title: item.title,
    icon:
      kind === 'ANNOUNCEMENT' ? 'icon-park-outline:volume-notice' : 'icon-park-outline:tips-one',
    tagTitle: kind === 'ANNOUNCEMENT' ? '公告' : '通知',
    tagType: kind === 'ANNOUNCEMENT' ? 'warning' : 'info',
    description: item.content,
    date: formatDateTime(item.publish_at || item.created_at),
    sourceType: kind,
    sourceId: item.id,
    isRead: wireBool(item.is_read ?? false),
  }
}

function toNoticeItem(item: NoticeSource): BannerItem {
  return { ...item, isRead: item.isRead } as BannerItem
}
</script>

<template>
  <n-popover
    v-model:show="showPopover"
    placement="bottom"
    trigger="click"
    arrow-point-to-center
    class="!p-0"
    @update:show="(show: boolean) => show && refresh()"
  >
    <template #trigger>
      <n-tooltip
        placement="bottom"
        trigger="hover"
      >
        <template #trigger>
          <CommonWrapper>
            <n-badge
              :value="unreadTotal"
              :max="99"
              style="color: unset"
            >
              <NovaIcon icon="icon-park-outline:remind" />
            </n-badge>
          </CommonWrapper>
        </template>
        消息
      </n-tooltip>
    </template>
    <NSpace
      vertical
      :size="0"
      style="width: 390px"
    >
      <div style="padding: 12px 12px 0; font-weight: 600">
        我的消息
      </div>
      <NoticeList
        :list="list"
        :loading="listState.loading"
        :has-more="hasMore"
        @open="handleOpen"
        @load-more="loadMore"
      />
      <NDivider style="margin: 0" />
      <div class="flex gap-2 p-2">
        <NButton
          class="flex-1"
          tertiary
          size="small"
          :disabled="unreadTotal <= 0"
          @click="markAllRead"
        >
          全部已读
        </NButton>
        <NButton
          class="flex-1"
          tertiary
          size="small"
          @click="goMore"
        >
          查看更多
        </NButton>
      </div>
    </NSpace>
  </n-popover>
  <MessageDetailModal
    ref="detailModalRef"
    @changed="handleDetailChanged"
  />
</template>
