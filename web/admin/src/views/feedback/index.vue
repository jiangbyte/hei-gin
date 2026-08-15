<!-- Author: Charlie -->

<script setup lang="ts">
import { sysFeedbackApi } from '@/api'
import { createTagColor, dictTypeColor, dictTypeData, formatDateTime } from '@/utils'
import { readPageMeta } from '@/utils/wire'
import { onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'

function renderDictTag(dictCode: string, value: string) {
  const label = dictTypeData(dictCode, value) || value
  const color = createTagColor(dictTypeColor(dictCode, value))
  return { label, color }
}

const router = useRouter()
const state = reactive({
  loading: false,
  rows: [] as any[],
  total: 0,
  page: 1,
  pageSize: 10,
})

onMounted(() => {
  void fetchPage()
})

async function fetchPage() {
  state.loading = true
  try {
    const response = await sysFeedbackApi.myPage({
      current: state.page,
      size: state.pageSize,
    })
    const data = response.data ?? {}
    state.rows = data.records ?? []
    const pageMeta = readPageMeta(data, { current: state.page, size: state.pageSize })
    state.total = pageMeta.total
    state.page = pageMeta.current
    state.pageSize = pageMeta.size
  } catch {
    state.rows = []
    state.total = 0
  } finally {
    state.loading = false
  }
}

function excerpt(content: string) {
  const text = String(content || '')
    .replace(/\s+/g, ' ')
    .trim()
  if (!text) return '点击查看详情'
  return text.length > 120 ? `${text.slice(0, 120)}…` : text
}

function goNew() {
  void router.push('/feedback/new')
}

function goDetail(id: string) {
  void router.push(`/feedback/${id}`)
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
  <div class="h-full min-h-0 p-4">
    <div class="mb-4 flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 class="text-xl font-semibold">
          意见反馈
        </h1>
        <p class="mt-1 text-sm text-[var(--text-color-3)]">
          提交问题与建议，并查看处理进度
        </p>
      </div>
      <NButton
        type="primary"
        @click="goNew"
      >
        <template #icon>
          <NovaIcon icon="icon-park-outline:plus" />
        </template>
        提交反馈
      </NButton>
    </div>

    <NCard
      :bordered="false"
      class="min-h-0"
    >
      <NSpin :show="state.loading">
        <NEmpty
          v-if="!state.loading && !state.rows.length"
          description="暂无反馈"
        >
          <NButton
            type="primary"
            class="mt-3"
            @click="goNew"
          >
            提交第一条反馈
          </NButton>
        </NEmpty>
        <NList
          v-else
          clickable
          show-divider
        >
          <NListItem
            v-for="item in state.rows"
            :key="item.id"
            class="cursor-pointer"
            @click="goDetail(item.id)"
          >
            <div class="flex w-full min-w-0 flex-col gap-1 py-1">
              <div class="flex flex-wrap items-center gap-2">
                <NText
                  strong
                  class="min-w-0 flex-1 truncate text-[15px]"
                >
                  {{ item.title }}
                </NText>
                <NTag
                  size="small"
                  :bordered="false"
                  :color="renderDictTag('FEEDBACK_CATEGORY', item.category).color"
                >
                  {{ renderDictTag('FEEDBACK_CATEGORY', item.category).label }}
                </NTag>
                <NTag
                  size="small"
                  :bordered="false"
                  :color="renderDictTag('FEEDBACK_STATUS', item.status).color"
                >
                  {{ renderDictTag('FEEDBACK_STATUS', item.status).label }}
                </NTag>
              </div>
              <NText
                depth="3"
                class="line-clamp-2 text-sm"
              >
                {{ excerpt(item.content) }}
              </NText>
              <NText
                depth="3"
                class="text-xs"
              >
                {{ formatDateTime(item.created_at) }}
              </NText>
            </div>
          </NListItem>
        </NList>
      </NSpin>
    </NCard>

    <div
      v-if="state.total > 0"
      class="mt-4 flex justify-end"
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
    </div>
  </div>
</template>
