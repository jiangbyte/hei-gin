<!--
  Author: Charlie

  消息详情页。
-->
<script setup lang="ts">
import { msgNoticeApi } from '@/api'
import { accountTypeLabel } from '@/constants/account'
import { dictTypeData, displayValue, formatDateTime, hasPermission } from '@/utils'
import { MdPreview, RichTextPreview } from '@/components/editor'
import { computed, onMounted, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const listPath = '/message/notice'

const state = reactive({
  loading: false,
  detail: {} as any,
})

const dataId = computed(() => {
  const id = route.query.id
  return typeof id === 'string' ? id : ''
})

const locationTags = computed(() => {
  const locs = state.detail.publish_locations || {}
  return Object.keys(locs).filter((k) => locs[k])
})

const accountTypeLabels = computed(() =>
  (state.detail.target_account_types || []).map((t: string) => accountTypeLabel(t)),
)

async function fetchDetail(id: string) {
  if (!id) return
  state.loading = true
  try {
    const response = await msgNoticeApi.detail({ id })
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
  router.push({ path: '/message/notice/edit', query: { id: dataId.value } })
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
      title="消息详情"
      :bordered="false"
    >
      <template #header-extra>
        <NSpace>
          <NButton @click="goBack">
            返回
          </NButton>
          <NButton
            v-if="hasPermission('message:notice:update') && dataId"
            type="primary"
            @click="goEdit"
          >
            编辑
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
                  类型
                </div>
                <div class="meta-value">
                  {{ state.detail.kind === 'ANNOUNCEMENT' ? '公告' : '通知' }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  等级
                </div>
                <div class="meta-value">
                  {{
                    dictTypeData('NOTIFICATION_SEVERITY', state.detail.severity) ||
                      displayValue(state.detail.severity)
                  }}
                </div>
              </div>
              <div
                v-if="state.detail.kind !== 'ANNOUNCEMENT'"
                class="meta-item"
              >
                <div class="meta-key">
                  分类
                </div>
                <div class="meta-value">
                  {{
                    dictTypeData('NOTIFICATION_CATEGORY', state.detail.category) ||
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
                    dictTypeData('PUBLISH_STATUS', state.detail.status) ||
                      displayValue(state.detail.status)
                  }}
                </div>
              </div>
              <div
                v-if="state.detail.kind === 'ANNOUNCEMENT'"
                class="meta-item"
              >
                <div class="meta-key">
                  是否置顶
                </div>
                <div class="meta-value">
                  {{ state.detail.is_pinned ? '是' : '否' }}
                </div>
              </div>
            </div>
          </section>

          <section class="meta-section">
            <h2 class="section-label">
              投放信息
            </h2>
            <div class="meta-grid">
              <div class="meta-item">
                <div class="meta-key">
                  目标范围
                </div>
                <div class="meta-value">
                  {{
                    dictTypeData('TARGET_SCOPE', state.detail.target_scope) ||
                      displayValue(state.detail.target_scope)
                  }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  目标账户类型
                </div>
                <div class="meta-value">
                  <span v-if="accountTypeLabels.length">{{ accountTypeLabels.join('、') }}</span>
                  <span
                    v-else
                    class="text-muted"
                  >—</span>
                </div>
              </div>
              <div
                v-if="state.detail.kind === 'ANNOUNCEMENT'"
                class="meta-item"
              >
                <div class="meta-key">
                  发布位置
                </div>
                <div class="meta-value">
                  <span v-if="locationTags.length">{{
                    locationTags.map((k) => dictTypeData('NOTIFY_LOCATION', k) || k).join('、')
                  }}</span>
                  <span
                    v-else
                    class="text-muted"
                  >—</span>
                </div>
              </div>
              <div
                v-if="state.detail.kind === 'ANNOUNCEMENT'"
                class="meta-item"
              >
                <div class="meta-key">
                  查看次数
                </div>
                <div class="meta-value">
                  {{ displayValue(state.detail.view_count) }}
                </div>
              </div>
            </div>
          </section>

          <section class="meta-section">
            <h2 class="section-label">
              时间信息
            </h2>
            <div class="meta-grid">
              <div class="meta-item">
                <div class="meta-key">
                  发布时间
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.publish_at) }}
                </div>
              </div>
              <div
                v-if="state.detail.kind === 'ANNOUNCEMENT'"
                class="meta-item"
              >
                <div class="meta-key">
                  过期时间
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.expire_at) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  创建时间
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.created_at) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  创建人
                </div>
                <div class="meta-value">
                  {{ state.detail.created_name || displayValue(state.detail.created_by) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  更新时间
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.updated_at) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  更新人
                </div>
                <div class="meta-value">
                  {{ state.detail.updated_name || displayValue(state.detail.updated_by) }}
                </div>
              </div>
            </div>
          </section>

          <section class="content-section">
            <h2 class="section-label">
              消息内容
            </h2>
            <div class="detail-content">
              <div
                v-if="state.detail.content_type === 'text'"
                class="whitespace-pre-wrap"
              >
                {{ state.detail.content }}
              </div>
              <MdPreview
                v-else-if="state.detail.content_type === 'markdown'"
                :value="state.detail.content"
                :preview="true"
              />
              <RichTextPreview
                v-else
                :value="state.detail.content"
              />
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
  min-height: 120px;
  color: var(--text-color-1, #333);
  font-size: 15px;
  line-height: 1.75;
}

.whitespace-pre-wrap {
  white-space: pre-wrap;
  word-break: break-word;
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
