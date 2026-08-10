<!--
  Author: Charlie

  反馈详情页。
-->
<script setup lang="ts">
import { msgFeedbackApi } from '@/api'
import { accountTypeLabel } from '@/constants/account'
import {
  dictTypeData,
  displayValue,
  formatDateTime,
  formatFileSize,
  hasPermission,
  isImageFile,
} from '@/utils'
import { computed, onMounted, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const listPath = '/message/feedback'

const state = reactive({
  loading: false,
  detail: {} as any,
})

const dataId = computed(() => {
  const id = route.query.id
  return typeof id === 'string' ? id : ''
})

const attachments = computed(() => {
  const list = state.detail.attachments
  if (Array.isArray(list) && list.length) return list
  return (state.detail.attach_object_names || []).map((object_name: string) => ({
    object_name,
    original_name: object_name,
    url: null,
  }))
})

async function fetchDetail(id: string) {
  if (!id) return
  state.loading = true
  try {
    const response = await msgFeedbackApi.detail({ id })
    state.detail = response.data ?? {}
  } finally {
    state.loading = false
  }
}

function goBack() {
  router.push(listPath)
}

function goEdit() {
  if (!dataId.value) return
  router.push({ path: '/message/feedback/edit', query: { id: dataId.value } })
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
  <div class="h-full min-h-0">
    <NCard
      class="h-full min-h-0 overflow-auto"
      title="反馈详情"
      :bordered="false"
    >
      <template #header-extra>
        <NSpace>
          <NButton @click="goBack">
            返回
          </NButton>
          <NButton
            v-if="hasPermission('message:feedback:update') && dataId"
            type="primary"
            @click="goEdit"
          >
            处理
          </NButton>
        </NSpace>
      </template>
      <NSpin :show="state.loading">
        <div class="detail-page">
          <header class="detail-header">
            <h1 class="detail-title">
              {{ displayValue(state.detail.title) }}
            </h1>
          </header>

          <section class="meta-section">
            <h2 class="section-label">
              基础信息
            </h2>
            <div class="meta-grid">
              <div class="meta-item">
                <div class="meta-key">
                  分类
                </div>
                <div class="meta-value">
                  {{
                    dictTypeData('FEEDBACK_CATEGORY', state.detail.category) ||
                      displayValue(state.detail.category)
                  }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  状态
                </div>
                <div class="meta-value">
                  {{
                    dictTypeData('FEEDBACK_STATUS', state.detail.status) ||
                      displayValue(state.detail.status)
                  }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  联系方式
                </div>
                <div class="meta-value">
                  {{ displayValue(state.detail.contact) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  账号类型
                </div>
                <div class="meta-value">
                  {{
                    accountTypeLabel(state.detail.submitter_account_type) ||
                      displayValue(state.detail.submitter_account_type)
                  }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  提交账号
                </div>
                <div class="meta-value">
                  <NFlex
                    align="center"
                    :size="8"
                  >
                    <NAvatar
                      v-if="state.detail.submitter_avatar"
                      :src="state.detail.submitter_avatar"
                      :size="24"
                      round
                    />
                    <NAvatar
                      v-else
                      :size="24"
                      round
                      color="#d9d9d9"
                    >
                      {{
                        (state.detail.submitter_nickname ||
                          state.detail.submitter_account_id ||
                          '?')[0]?.toUpperCase()
                      }}
                    </NAvatar>
                    <span>{{
                      state.detail.submitter_nickname || state.detail.submitter_account_id
                    }}</span>
                  </NFlex>
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  提交时间
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.created_at) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  回复时间
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.replied_at) }}
                </div>
              </div>
            </div>
          </section>

          <section class="content-section">
            <h2 class="section-label">
              反馈内容
            </h2>
            <div class="detail-content whitespace-pre-wrap">
              {{ displayValue(state.detail.content) }}
            </div>
          </section>

          <section class="meta-section">
            <h2 class="section-label">
              附件
            </h2>
            <div
              v-if="attachments.length"
              class="attach-list"
            >
              <div
                v-for="(item, idx) in attachments"
                :key="item.object_name || idx"
                class="attach-item"
              >
                <a
                  v-if="item.url"
                  :href="item.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="attach-link"
                >
                  <NImage
                    v-if="isImageFile(item)"
                    :src="item.url"
                    :alt="attachmentLabel(item)"
                    width="72"
                    height="72"
                    object-fit="cover"
                    class="attach-thumb"
                  />
                  <span v-else>{{ attachmentLabel(item) }}</span>
                </a>
                <span v-else>{{ attachmentLabel(item) }}</span>
                <span
                  v-if="item.size != null"
                  class="attach-meta"
                >{{
                  formatFileSize(item.size)
                }}</span>
              </div>
            </div>
            <div
              v-else
              class="text-muted"
            >
              无附件
            </div>
          </section>

          <section
            v-if="state.detail.reply"
            class="content-section"
          >
            <h2 class="section-label">
              管理员回复
            </h2>
            <div class="detail-content whitespace-pre-wrap">
              {{ state.detail.reply }}
            </div>
          </section>
        </div>
      </NSpin>
    </NCard>
  </div>
</template>

<style scoped>
.detail-page {
  max-width: 880px;
}

.detail-header {
  margin-bottom: 28px;
}

.detail-title {
  margin: 0 0 14px;
  color: var(--text-color-1, #1f1f1f);
  font-size: 22px;
  font-weight: 650;
  line-height: 1.35;
}

.meta-section,
.content-section {
  margin-bottom: 28px;
}

.section-label {
  margin: 0 0 14px;
  color: var(--text-color-2, #666);
  font-size: 13px;
  font-weight: 600;
}

.meta-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px 28px;
}

.meta-item {
  min-width: 0;
}

.meta-key {
  margin-bottom: 4px;
  color: var(--text-color-3, #999);
  font-size: 12px;
  line-height: 1.4;
}

.meta-value {
  color: var(--text-color-1, #333);
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
}

.text-muted {
  color: var(--text-color-3, #999);
}

.detail-content {
  min-height: 80px;
  color: var(--text-color-1, #333);
  font-size: 15px;
  line-height: 1.75;
}

.whitespace-pre-wrap {
  white-space: pre-wrap;
  word-break: break-word;
}

.attach-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 20px;
}

.attach-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.attach-link {
  color: var(--text-color-1, #333);
  text-decoration: none;
}

.attach-link:hover {
  color: var(--primary-color, #18a058);
}

.attach-thumb {
  border-radius: 6px;
  overflow: hidden;
}

.attach-meta {
  color: var(--text-color-3, #999);
  font-size: 12px;
}

@media (max-width: 960px) {
  .meta-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .meta-grid {
    grid-template-columns: 1fr;
  }
}
</style>
