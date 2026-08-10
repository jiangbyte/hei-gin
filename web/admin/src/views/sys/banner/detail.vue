<!--
  Author: Charlie

  展示图详情页。
-->
<script setup lang="ts">
import { bannerApi } from '@/api'
import { accountTypeLabel } from '@/constants/account'
import { dictTypeData, displayValue, formatDateTime, hasPermission } from '@/utils'
import { computed, onMounted, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const listPath = '/sys/banner'

const state = reactive({
  loading: false,
  detail: {} as any,
})

const dataId = computed(() => {
  const id = route.query.id
  return typeof id === 'string' ? id : ''
})

const accountTypeLabels = computed(() =>
  (state.detail.target_account_types || []).map((t: string) => accountTypeLabel(t)),
)

async function fetchDetail(id: string) {
  if (!id) return
  state.loading = true
  try {
    const response = await bannerApi.detail({ id })
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
  router.push({ path: '/sys/banner/edit', query: { id: dataId.value } })
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
      title="展示图详情"
      :bordered="false"
    >
      <template #header-extra>
        <NSpace>
          <NButton @click="goBack">
            返回
          </NButton>
          <NButton
            v-if="hasPermission('sys:banner:update') && dataId"
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
            <NImage
              v-if="state.detail.image_url || state.detail.image"
              :src="state.detail.image_url || state.detail.image"
              :alt="state.detail.title || '图片'"
              width="320"
              height="180"
              object-fit="cover"
            />
          </header>

          <section class="meta-section">
            <h2 class="section-label">
              投放信息
            </h2>
            <div class="meta-grid">
              <div class="meta-item">
                <div class="meta-key">
                  目标账户类型
                </div>
                <div class="meta-value">
                  <span v-if="accountTypeLabels.length">{{ accountTypeLabels.join('、') }}</span>
                  <span v-else>{{ displayValue(null) }}</span>
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  分类
                </div>
                <div class="meta-value">
                  {{
                    dictTypeData('BANNER_CATEGORY', state.detail.category) ||
                      displayValue(state.detail.category)
                  }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  类型
                </div>
                <div class="meta-value">
                  {{
                    dictTypeData('BANNER_TYPE', state.detail.type) ||
                      displayValue(state.detail.type)
                  }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  位置
                </div>
                <div class="meta-value">
                  {{
                    dictTypeData('BANNER_POSITION', state.detail.position) ||
                      displayValue(state.detail.position)
                  }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  状态
                </div>
                <div class="meta-value">
                  {{
                    dictTypeData('COMMON_STATUS', state.detail.status) ||
                      displayValue(state.detail.status)
                  }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  排序
                </div>
                <div class="meta-value">
                  {{ displayValue(state.detail.sort) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  互动次数
                </div>
                <div class="meta-value">
                  {{ displayValue(state.detail.interaction_count) }}
                </div>
              </div>
            </div>
          </section>

          <section class="meta-section">
            <h2 class="section-label">
              跳转与时段
            </h2>
            <div class="meta-grid">
              <div class="meta-item">
                <div class="meta-key">
                  链接类型
                </div>
                <div class="meta-value">
                  {{
                    dictTypeData('BANNER_LINK_TYPE', state.detail.link_type) ||
                      displayValue(state.detail.link_type)
                  }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  目标 URL
                </div>
                <div class="meta-value">
                  {{ displayValue(state.detail.url) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  开始时间
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.start_at) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  结束时间
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.end_at) }}
                </div>
              </div>
            </div>
          </section>

          <section class="meta-section">
            <h2 class="section-label">
              文案
            </h2>
            <div class="meta-grid">
              <div class="meta-item meta-item--full">
                <div class="meta-key">
                  摘要
                </div>
                <div class="meta-value">
                  {{ displayValue(state.detail.summary) }}
                </div>
              </div>
              <div class="meta-item meta-item--full">
                <div class="meta-key">
                  描述
                </div>
                <div class="meta-value whitespace-pre-wrap">
                  {{ displayValue(state.detail.description) }}
                </div>
              </div>
            </div>
          </section>

          <section class="meta-section">
            <h2 class="section-label">
              审计信息
            </h2>
            <div class="meta-grid">
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
                  创建时间
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.created_at) }}
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
              <div class="meta-item">
                <div class="meta-key">
                  更新时间
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.updated_at) }}
                </div>
              </div>
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
  margin: 0 0 16px;
  color: var(--text-color-1, #1f1f1f);
  font-size: 22px;
  font-weight: 650;
  line-height: 1.35;
}

.meta-section {
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

.meta-item--full {
  grid-column: 1 / -1;
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
