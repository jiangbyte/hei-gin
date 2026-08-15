<!-- Author: Charlie -->

<script setup lang="ts">
import { sysFeedbackApi } from '@/api'
import {
  createTagColor,
  dictTypeColor,
  dictTypeData,
  displayValue,
  formatDateTime,
  formatFileSize,
  isImageFile,
} from '@/utils'
import { computed, onMounted, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

function renderDictTag(dictCode: string, value: string) {
  const label = dictTypeData(dictCode, value) || value
  const color = createTagColor(dictTypeColor(dictCode, value))
  return { label, color }
}

const route = useRoute()
const router = useRouter()
const state = reactive({
  loading: false,
  detail: null as any,
  notFound: false,
})

const dataId = computed(() => {
  const id = route.params.id
  return typeof id === 'string' ? id : ''
})

const attachments = computed(() => {
  const list = state.detail?.attachments
  if (Array.isArray(list) && list.length) return list
  return (state.detail?.attach_object_names || []).map((object_name: string) => ({
    object_name,
    original_name: object_name,
    url: null,
  }))
})

async function fetchDetail(id: string) {
  if (!id) {
    state.notFound = true
    state.detail = null
    return
  }
  state.loading = true
  state.notFound = false
  try {
    const response = await sysFeedbackApi.myDetail(id)
    state.detail = response.data ?? null
    if (!state.detail) state.notFound = true
  } catch {
    state.detail = null
    state.notFound = true
  } finally {
    state.loading = false
  }
}

function goBack() {
  void router.push('/feedback')
}

function attachmentLabel(item: any) {
  return item.original_name || item.object_name || '附件'
}

onMounted(() => {
  void fetchDetail(dataId.value)
})
watch(dataId, (id) => {
  void fetchDetail(id)
})
</script>

<template>
  <div class="h-full min-h-0 p-4">
    <div class="mb-4">
      <NButton
        text
        @click="goBack"
      >
        <template #icon>
          <NovaIcon icon="icon-park-outline:left" />
        </template>
        返回我的反馈
      </NButton>
    </div>

    <NCard
      :bordered="false"
      class="max-w-4xl"
    >
      <NSpin :show="state.loading">
        <NEmpty
          v-if="state.notFound"
          description="反馈不存在或无权查看"
        >
          <NButton
            type="primary"
            class="mt-3"
            @click="goBack"
          >
            返回列表
          </NButton>
        </NEmpty>
        <article v-else-if="state.detail">
          <h1 class="text-xl font-semibold md:text-2xl">
            {{ displayValue(state.detail.title) }}
          </h1>
          <div class="mt-3 flex flex-wrap items-center gap-2">
            <NTag
              size="small"
              :bordered="false"
              :color="renderDictTag('FEEDBACK_CATEGORY', state.detail.category).color"
            >
              {{ renderDictTag('FEEDBACK_CATEGORY', state.detail.category).label }}
            </NTag>
            <NTag
              size="small"
              :bordered="false"
              :color="renderDictTag('FEEDBACK_STATUS', state.detail.status).color"
            >
              {{ renderDictTag('FEEDBACK_STATUS', state.detail.status).label }}
            </NTag>
            <NText
              depth="3"
              class="text-xs"
            >
              提交于 {{ formatDateTime(state.detail.created_at) }}
            </NText>
          </div>

          <div class="mt-6 grid gap-3 text-sm md:grid-cols-2">
            <div>
              <div class="mb-1 text-[var(--text-color-3)]">
                联系方式
              </div>
              <div>{{ displayValue(state.detail.contact) }}</div>
            </div>
          </div>

          <div class="mt-6 border-t border-[var(--border-color)] pt-6">
            <div class="mb-2 text-sm font-medium">
              内容
            </div>
            <div class="whitespace-pre-wrap text-sm leading-7">
              {{ displayValue(state.detail.content) }}
            </div>
          </div>

          <div
            v-if="attachments.length"
            class="mt-6 border-t border-[var(--border-color)] pt-6"
          >
            <div class="mb-3 text-sm font-medium">
              附件
            </div>
            <div class="flex flex-wrap gap-3">
              <template
                v-for="item in attachments"
                :key="item.object_name || attachmentLabel(item)"
              >
                <div
                  v-if="item.url && isImageFile(item)"
                  class="w-28"
                >
                  <NImage
                    :src="item.url"
                    :alt="attachmentLabel(item)"
                    width="112"
                    height="112"
                    object-fit="cover"
                    style="border-radius: 8px"
                  />
                  <div class="mt-1 truncate text-xs">
                    {{ attachmentLabel(item) }}
                  </div>
                </div>
                <a
                  v-else
                  class="rounded border border-[var(--border-color)] px-3 py-2 text-sm"
                  :href="item.url || undefined"
                  target="_blank"
                  rel="noreferrer"
                >
                  {{ attachmentLabel(item) }}
                  <template v-if="item.size != null">（{{ formatFileSize(item.size) }}）</template>
                </a>
              </template>
            </div>
          </div>

          <div
            v-if="state.detail.reply"
            class="mt-6 rounded-lg border border-[var(--border-color)] bg-[var(--hover-color)] p-4"
          >
            <div class="mb-2 text-sm font-medium">
              官方回复
            </div>
            <div class="whitespace-pre-wrap text-sm leading-7">
              {{ state.detail.reply }}
            </div>
            <div
              v-if="state.detail.replied_at"
              class="mt-2 text-xs text-[var(--text-color-3)]"
            >
              回复于 {{ formatDateTime(state.detail.replied_at) }}
            </div>
          </div>
        </article>
      </NSpin>
    </NCard>
  </div>
</template>
