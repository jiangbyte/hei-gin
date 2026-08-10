<!-- Author: Charlie -->

<script setup lang="ts">
import { messageApi } from '@/api'
import { useMessageUnreadStore } from '@/stores'
import { displayValue, formatDateTime, wireBool } from '@/utils'
import { dictTypeData } from '@/utils/dict'
import { computed, reactive } from 'vue'

type MessageKind = 'NOTIFICATION' | 'ANNOUNCEMENT'

const emit = defineEmits<{
  changed: [payload: { type: string; id: string }]
}>()

const unreadStore = useMessageUnreadStore()

const state = reactive({
  show: false,
  loading: false,
  actionLoading: false,
  source: {} as any,
  detail: {} as any,
  readLocally: false,
})

function resolveKind(raw: unknown): MessageKind {
  return String(raw || 'NOTIFICATION').toUpperCase() === 'ANNOUNCEMENT'
    ? 'ANNOUNCEMENT'
    : 'NOTIFICATION'
}

function asReadFlag(value: unknown): boolean {
  if (typeof value === 'boolean' || typeof value === 'string') {
    return wireBool(value)
  }
  return false
}

const messageKind = computed<MessageKind>(() =>
  resolveKind(state.detail.kind || state.source.sourceType || state.source.type),
)

const kindLabel = computed(() => (messageKind.value === 'ANNOUNCEMENT' ? '公告' : '通知'))

const titleText = computed(() => displayValue(state.detail.title || state.source.title))

const contentText = computed(() => displayValue(state.detail.content || state.source.content))

const publishText = computed(() =>
  formatDateTime(state.detail.publish_at || state.source.publish_at || state.detail.created_at),
)

const severityLabel = computed(() => {
  const severity = state.detail.severity || state.source.severity
  if (!severity) return ''
  return dictTypeData('NOTIFICATION_SEVERITY', severity) || severity
})

const isRead = computed(
  () => state.readLocally || asReadFlag(state.detail.is_read) || asReadFlag(state.source.is_read),
)

async function open(source: any) {
  const wasUnread = !asReadFlag(source?.is_read)
  state.source = source ?? {}
  state.detail = {}
  state.readLocally = false
  state.show = true
  state.loading = true
  try {
    // myDetail 后端会顺带标记已读
    const response = await messageApi.myDetail(state.source.id)
    state.detail = response.data ?? {}
  } finally {
    state.loading = false
  }

  const id = state.detail.id || state.source.id
  const kind = resolveKind(state.detail.kind || messageKind.value)
  if (id && wasUnread) {
    state.detail.is_read = true
    state.source.is_read = true
    state.readLocally = true
    unreadStore.notifyRead()
    emit('changed', { type: kind, id })
    void unreadStore.refresh()
  }
}

async function markRead() {
  const id = state.detail.id || state.source.id
  if (!id || isRead.value) return
  const kind = resolveKind(state.detail.kind || messageKind.value)
  state.actionLoading = true
  try {
    await messageApi.read({ ids: [id] })
    state.detail.is_read = true
    state.source.is_read = true
    state.readLocally = true
    unreadStore.notifyRead()
    emit('changed', { type: kind, id })
    void unreadStore.refresh()
  } finally {
    state.actionLoading = false
  }
}

defineExpose({ open })
</script>

<template>
  <NModal
    v-model:show="state.show"
    preset="card"
    :mask-closable="true"
    :bordered="false"
    size="medium"
    style="width: 560px; max-width: calc(100vw - 32px)"
  >
    <template #header>
      <NSpace
        align="center"
        :size="8"
      >
        <NTag
          size="small"
          :bordered="false"
          :type="messageKind === 'ANNOUNCEMENT' ? 'warning' : 'info'"
        >
          {{ kindLabel }}
        </NTag>
        <NTag
          size="small"
          :bordered="false"
          :type="isRead ? 'success' : 'warning'"
        >
          {{ isRead ? '已读' : '未读' }}
        </NTag>
      </NSpace>
    </template>

    <NSpin :show="state.loading">
      <NSpace
        vertical
        :size="16"
      >
        <NThing>
          <template #header>
            <NText
              strong
              style="font-size: 18px; line-height: 1.4"
            >
              {{ titleText }}
            </NText>
          </template>
          <template #description>
            <NSpace
              :size="8"
              align="center"
            >
              <NTag
                v-if="severityLabel"
                size="small"
                :bordered="false"
              >
                {{ severityLabel }}
              </NTag>
              <NText
                depth="3"
                style="font-size: 12px"
              >
                {{ publishText }}
              </NText>
            </NSpace>
          </template>
        </NThing>

        <NDivider style="margin: 0" />

        <NScrollbar style="max-height: min(420px, calc(100vh - 280px))">
          <NText
            depth="2"
            style="font-size: 14px; line-height: 1.75; white-space: pre-wrap"
          >
            {{ contentText }}
          </NText>
        </NScrollbar>

        <NSpace
          v-if="!isRead"
          justify="end"
        >
          <NButton
            text
            type="primary"
            :loading="state.actionLoading"
            @click="markRead"
          >
            标记已读
          </NButton>
        </NSpace>
      </NSpace>
    </NSpin>
  </NModal>
</template>
