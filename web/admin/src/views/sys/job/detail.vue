<!--
  Author: Charlie

  任务详情页。
-->
<script setup lang="ts">
import { jobApi } from '@/api'
import { displayValue, formatDateTime, hasPermission, wireBool } from '@/utils'
import { computed, onMounted, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const listPath = '/sys/job'

const state = reactive({
  loading: false,
  detail: {} as any,
})

const dataId = computed(() => {
  const id = route.query.id
  return typeof id === 'string' ? id : ''
})

async function fetchDetail(id: string) {
  if (!id) return
  state.loading = true
  try {
    const response = await jobApi.detail({ id })
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
  router.push({ path: '/sys/job/edit', query: { id: dataId.value } })
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
      title="任务详情"
      :bordered="false"
    >
      <template #header-extra>
        <NSpace>
          <NButton @click="goBack">
            返回
          </NButton>
          <NButton
            v-if="hasPermission('sys:job:update') && dataId"
            type="primary"
            @click="goEdit"
          >
            编辑
          </NButton>
        </NSpace>
      </template>
      <NSpin :show="state.loading">
        <div class="detail-page">
          <section class="meta-section">
            <h2 class="section-label">
              基本信息
            </h2>
            <div class="meta-grid">
              <div class="meta-item">
                <div class="meta-key">
                  任务名称
                </div>
                <div class="meta-value">
                  {{ displayValue(state.detail.job_name) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  启用状态
                </div>
                <div class="meta-value">
                  {{ wireBool(state.detail.enabled) ? '启用' : '停用' }}
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
              <div class="meta-item meta-item--full">
                <div class="meta-key">
                  任务描述
                </div>
                <div class="meta-value">
                  {{ displayValue(state.detail.description) }}
                </div>
              </div>
            </div>
          </section>

          <section class="meta-section">
            <h2 class="section-label">
              调度配置
            </h2>
            <div class="meta-grid">
              <div class="meta-item meta-item--full">
                <div class="meta-key">
                  执行类
                </div>
                <div class="meta-value">
                  {{ displayValue(state.detail.execute_class) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  触发类型
                </div>
                <div class="meta-value">
                  {{ state.detail.execute_type === 'CRON' ? 'CRON 表达式' : '固定间隔' }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  触发配置
                </div>
                <div class="meta-value">
                  {{ displayValue(state.detail.trigger_config) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  上次执行
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.last_run_time) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  下次执行
                </div>
                <div class="meta-value">
                  {{ formatDateTime(state.detail.next_run_time) }}
                </div>
              </div>
              <div class="meta-item">
                <div class="meta-key">
                  任务参数
                </div>
                <div class="meta-value">
                  {{
                    state.detail.execute_param && typeof state.detail.execute_param === 'object'
                      ? JSON.stringify(state.detail.execute_param)
                      : displayValue(state.detail.execute_param)
                  }}
                </div>
              </div>
              <div class="meta-item meta-item--full">
                <div class="meta-key">
                  上次执行结果
                </div>
                <div class="meta-value">
                  {{ displayValue(state.detail.last_execute_result) }}
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
  max-width: 1180px;
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
