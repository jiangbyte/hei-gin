<!--
  Author: Charlie

  运营工作台：状态条 + 台账式指标 + 主图侧栏，数据来自 dashboard overview / banners list。
-->
<script setup lang="ts">
import { Chart } from '@antv/g2'
import { Icon } from '@iconify/vue/offline'
import { NIcon } from 'naive-ui'
import { bannerApi, dashboardApi } from '@/api'
import { useAuthStore } from '@/stores'
import { accountTypeLabel } from '@/constants/account'
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

type ChartInstance = InstanceType<typeof Chart>

const authStore = useAuthStore()
const router = useRouter()
const trendChartRef = ref<HTMLDivElement | null>(null)
const filePieChartRef = ref<HTMLDivElement | null>(null)
let trendChart: ChartInstance | null = null
let filePieChart: ChartInstance | null = null
const clockTimer = ref<number | null>(null)
const avatarImgProps = { referrerPolicy: 'no-referrer' } as any
const appTitle = import.meta.env.VITE_APP_TITLE || 'HEI Admin'

const state = reactive({
  loading: false,
  chartLoadError: false,
  now: Date.now(),
  banners: [] as any[],
  overview: {
    summary: {
      account_total: 0,
      online_sessions: 0,
      file_total: 0,
      storage_bytes: 0,
    },
    accounts: {
      enabled: 0,
      disabled: 0,
      today_new: 0,
      by_type: [] as Array<{ name: string; value: number }>,
    },
    iam: {
      role_count: 0,
      dept_count: 0,
      group_count: 0,
      menu_count: 0,
    },
    ops_today: {
      audit_total: 0,
      audit_failed: 0,
      feedback_pending: 0,
    },
    trends: {
      account_trend: [] as any[],
      audit_trend: [] as any[],
    },
    files: {
      by_content_type: [] as Array<{ name: string; value: number }>,
    },
  },
})

const displayName = computed(() => {
  const user = authStore.userInfo
  const nickname = String(user?.nickname ?? '').trim()
  const name = String(user?.name ?? '').trim()
  if (nickname && name && nickname !== name) {
    return `${nickname}（${name}）`
  }
  return nickname || name || user?.account || '-'
})

const avatarUrl = computed(() => authStore.userInfo?.avatar || undefined)
const deptText = computed(() => mapNames(authStore.userInfo?.deptIdNames) || '未分配部门')

const greeting = computed(() => {
  const hour = new Date(state.now).getHours()
  if (hour < 6) return '凌晨好'
  if (hour < 11) return '早上好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

const clockText = computed(() => {
  const date = new Date(state.now)
  const weekdays = ['日', '一', '二', '三', '四', '五', '六']
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  const ss = String(date.getSeconds()).padStart(2, '0')
  return `${y}.${m}.${d} 周${weekdays[date.getDay()]} ${hh}:${mm}:${ss}`
})

const pulseItems = computed(() => [
  {
    key: 'audit',
    label: '今日审计',
    value: Number(state.overview.ops_today.audit_total ?? 0),
  },
  {
    key: 'failed',
    label: '失败',
    value: Number(state.overview.ops_today.audit_failed ?? 0),
    danger: true,
  },
  {
    key: 'online',
    label: '在线会话',
    value: Number(state.overview.summary.online_sessions ?? 0),
    path: '/auth/session',
  },
  {
    key: 'feedback',
    label: '待处理反馈',
    value: Number(state.overview.ops_today.feedback_pending ?? 0),
    path: '/sys/feedback',
  },
])

const ledgerItems = computed(() => [
  {
    key: 'account',
    label: '账号',
    value: Number(state.overview.summary.account_total ?? 0),
    note: `今日 +${Number(state.overview.accounts.today_new ?? 0)}`,
    path: '/iam/account',
  },
  {
    key: 'enabled',
    label: '启用',
    value: Number(state.overview.accounts.enabled ?? 0),
    note: `禁用 ${Number(state.overview.accounts.disabled ?? 0)}`,
    path: '/iam/account',
  },
  {
    key: 'dept',
    label: '部门',
    value: Number(state.overview.iam.dept_count ?? 0),
    note: '组织节点',
    path: '/iam/dept',
  },
  {
    key: 'role',
    label: '角色',
    value: Number(state.overview.iam.role_count ?? 0),
    note: '授权主体',
    path: '/iam/role',
  },
  {
    key: 'group',
    label: '用户组',
    value: Number(state.overview.iam.group_count ?? 0),
    note: '批量授权',
    path: '/iam/group',
  },
  {
    key: 'menu',
    label: '菜单',
    value: Number(state.overview.iam.menu_count ?? 0),
    note: '可见入口',
    path: '/iam/resource',
  },
  {
    key: 'file',
    label: '文件',
    value: Number(state.overview.summary.file_total ?? 0),
    note: formatFileSize(state.overview.summary.storage_bytes),
    path: '/sys/file',
  },
])

const accountBars = computed(() => {
  const rows = (state.overview.accounts.by_type ?? []).map((item) => ({
    name: accountTypeLabel(item.name) || item.name,
    value: Number(item.value ?? 0),
  }))
  const max = Math.max(...rows.map((row) => row.value), 1)
  return rows.map((row) => ({ ...row, ratio: Math.round((row.value / max) * 100) }))
})

const filePieData = computed(() =>
  (state.overview.files.by_content_type ?? [])
    .slice(0, 6)
    .map((item) => ({
      name: simplifyContentType(item.name),
      value: Number(item.value ?? 0),
    }))
    .filter((item) => item.value > 0),
)

const trendData = computed(() => [
  ...state.overview.trends.account_trend.map((item) => ({
    date: item.date,
    value: Number(item.value ?? 0),
    type: '新增账号',
  })),
  ...state.overview.trends.audit_trend.map((item) => ({
    date: item.date,
    value: Number(item.value ?? 0),
    type: '审计量',
  })),
])

onMounted(() => {
  void fetchOverview()
  void loadBanners()
  clockTimer.value = window.setInterval(() => {
    state.now = Date.now()
  }, 1000)
})

onBeforeUnmount(() => {
  destroyCharts()
  if (clockTimer.value != null) {
    window.clearInterval(clockTimer.value)
    clockTimer.value = null
  }
})

watch(
  () => trendData.value,
  async () => {
    await nextTick()
    await renderTrendChart()
  },
)

watch(
  () => filePieData.value,
  async () => {
    await nextTick()
    await renderFilePieChart()
  },
)

async function fetchOverview() {
  state.loading = true
  try {
    const response = await dashboardApi.overview()
    state.overview = Object.assign(state.overview, response.data ?? {})
    await nextTick()
    await Promise.all([renderTrendChart(), renderFilePieChart()])
  } finally {
    state.loading = false
  }
}

async function loadBanners() {
  try {
    const response = await bannerApi.list({ position: 'ADMIN_TOP' })
    state.banners = Array.isArray(response.data) ? response.data : []
  } catch {
    state.banners = []
  }
}

function mapNames(items?: Array<{ id?: string; name?: string }>) {
  return (items ?? [])
    .map((item) => item.name)
    .filter(Boolean)
    .join(' / ')
}

function formatFileSize(size?: number | string | null) {
  const value = Number(size ?? 0)
  if (!Number.isFinite(value) || value <= 0) {
    return '0 B'
  }
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let current = value
  let unitIndex = 0
  while (current >= 1024 && unitIndex < units.length - 1) {
    current /= 1024
    unitIndex += 1
  }
  return `${current.toFixed(unitIndex === 0 ? 0 : 2)} ${units[unitIndex]}`
}

function simplifyContentType(value: string) {
  const text = String(value || 'unknown')
  if (text.includes('/')) {
    return text.split('/').pop() || text
  }
  return text
}

async function renderTrendChart() {
  destroyTrendChart()
  state.chartLoadError = false
  if (!trendChartRef.value) {
    return
  }
  try {
    const chart = new Chart({
      container: trendChartRef.value,
      autoFit: true,
      height: 300,
    })
    chart.options({
      type: 'view',
      autoFit: true,
      data: trendData.value,
      children: [
        {
          type: 'area',
          encode: { x: 'date', y: 'value', color: 'type' },
          scale: { color: { range: ['rgba(22,119,255,0.18)', 'rgba(19,194,194,0.16)'] } },
          style: { fillOpacity: 1 },
          legend: false,
          axis: false,
        },
        {
          type: 'line',
          encode: { x: 'date', y: 'value', color: 'type' },
          scale: { color: { range: ['#1677FF', '#13C2C2'] } },
          style: { lineWidth: 2.2 },
          axis: { x: { title: false }, y: { title: false, grid: true } },
          legend: { color: { position: 'top' } },
        },
      ],
    })
    trendChart = chart
    await chart.render()
  } catch {
    state.chartLoadError = true
  }
}

async function renderFilePieChart() {
  destroyFilePieChart()
  if (!filePieChartRef.value || !filePieData.value.length) {
    return
  }
  try {
    const chart = new Chart({
      container: filePieChartRef.value,
      autoFit: true,
      height: 220,
    })
    chart.options({
      type: 'interval',
      autoFit: true,
      data: filePieData.value,
      encode: { y: 'value', color: 'name' },
      transform: [{ type: 'stackY' }],
      coordinate: { type: 'theta', outerRadius: 0.85 },
      scale: {
        color: {
          range: ['#1677FF', '#13C2C2', '#722ED1', '#FA8C16', '#52C41A', '#EB2F96'],
        },
      },
      legend: {
        color: {
          position: 'bottom',
          layout: { justifyContent: 'center' },
        },
      },
      tooltip: {
        title: 'name',
        items: [{ channel: 'y', name: '数量' }],
      },
      style: { stroke: '#fff', lineWidth: 1 },
      animate: false,
    })
    filePieChart = chart
    await chart.render()
  } catch {
    // 侧栏饼图失败时保留空态即可，不阻断主图
  }
}

function destroyTrendChart() {
  trendChart?.destroy()
  trendChart = null
}

function destroyFilePieChart() {
  filePieChart?.destroy()
  filePieChart = null
}

function destroyCharts() {
  destroyTrendChart()
  destroyFilePieChart()
}

function go(path?: string) {
  if (!path) return
  router.push(path)
}

function openBanner(banner: any) {
  const link = String(banner?.url || '').trim()
  if (!link || banner?.link_type === 'NONE') return
  if (banner.link_type === 'ROUTE' || link.startsWith('/')) {
    router.push(link.startsWith('/') ? link : `/${link}`)
    return
  }
  window.open(link, '_blank', 'noopener,noreferrer')
}
</script>

<template>
  <NSpin :show="state.loading">
    <n-el class="board">
      <!-- 欢迎区：左欢迎卡 + 右轮播（轮播仅占右侧栏） -->
      <section
        class="welcome"
        :class="{ 'welcome--solo': !state.banners.length }"
      >
        <header class="rail">
          <div class="rail__top">
            <div class="rail__identity">
              <NAvatar
                v-if="avatarUrl"
                round
                :size="44"
                :src="avatarUrl"
                :img-props="avatarImgProps"
              />
              <NAvatar
                v-else
                round
                :size="44"
              >
                <NIcon :size="22">
                  <Icon icon="icon-park-outline:user" />
                </NIcon>
              </NAvatar>
              <div class="rail__copy">
                <div class="rail__hello">
                  {{ greeting }}，{{ displayName }}
                </div>
                <div class="rail__meta">
                  {{ deptText }} · {{ appTitle }}
                </div>
              </div>
            </div>
            <div class="rail__clock">
              {{ clockText }}
            </div>
          </div>

          <div class="rail__pulse">
            <button
              v-for="item in pulseItems"
              :key="item.key"
              type="button"
              class="pulse"
              :class="{ 'pulse--danger': item.danger, 'pulse--link': Boolean(item.path) }"
              @click="go(item.path)"
            >
              <span class="pulse__label">{{ item.label }}</span>
              <span class="pulse__value">{{ item.value }}</span>
            </button>
          </div>
        </header>

        <div
          v-if="state.banners.length"
          class="promo"
        >
          <NCarousel
            autoplay
            :interval="5000"
            show-dots
          >
            <button
              v-for="banner in state.banners"
              :key="banner.id"
              type="button"
              class="promo__slide"
              @click="openBanner(banner)"
            >
              <img
                v-if="banner.image_url || banner.image"
                :src="banner.image_url || banner.image"
                :alt="banner.title"
                class="promo__image"
              >
              <div class="promo__veil" />
              <div class="promo__text">
                <strong>{{ banner.title }}</strong>
                <span v-if="banner.summary">{{ banner.summary }}</span>
              </div>
            </button>
          </NCarousel>
        </div>
      </section>

      <!-- 台账式指标：横排分割 -->
      <section class="ledger">
        <button
          v-for="item in ledgerItems"
          :key="item.key"
          type="button"
          class="ledger__cell"
          @click="go(item.path)"
        >
          <span class="ledger__label">{{ item.label }}</span>
          <span class="ledger__value">{{ item.value }}</span>
          <span class="ledger__note">{{ item.note }}</span>
        </button>
      </section>

      <!-- 主区：趋势图 + 侧栏占比 -->
      <section class="main">
        <div class="panel trend">
          <div class="panel__head">
            <div>
              <div class="panel__title">
                近 7 日动态
              </div>
              <div class="panel__sub">
                新增账号与审计量
              </div>
            </div>
            <NButton
              text
              :loading="state.loading"
              @click="fetchOverview"
            >
              <template #icon>
                <NIcon>
                  <Icon icon="icon-park-outline:reload" />
                </NIcon>
              </template>
              刷新
            </NButton>
          </div>
          <div
            ref="trendChartRef"
            class="trend__chart"
          />
          <NAlert
            v-if="state.chartLoadError"
            class="mt-3"
            type="warning"
            :show-icon="false"
          >
            图表加载失败，请刷新后重试。
          </NAlert>
        </div>

        <aside class="side">
          <div class="panel">
            <div class="panel__head">
              <div class="panel__title">
                账号构成
              </div>
              <NButton
                text
                type="primary"
                @click="go('/iam/account')"
              >
                管理
              </NButton>
            </div>
            <div
              v-if="accountBars.length"
              class="bars"
            >
              <div
                v-for="row in accountBars"
                :key="row.name"
                class="bar-row"
              >
                <div class="bar-row__meta">
                  <span>{{ row.name }}</span>
                  <strong>{{ row.value }}</strong>
                </div>
                <div class="bar-row__track">
                  <div
                    class="bar-row__fill bar-row__fill--blue"
                    :style="{ width: `${row.ratio}%` }"
                  />
                </div>
              </div>
            </div>
            <NEmpty
              v-else
              description="暂无账号数据"
              size="small"
            />
          </div>

          <div class="panel">
            <div class="panel__head">
              <div class="panel__title">
                文件类型占比
              </div>
              <NButton
                text
                type="primary"
                @click="go('/sys/file')"
              >
                管理
              </NButton>
            </div>
            <div
              v-show="filePieData.length"
              ref="filePieChartRef"
              class="file-pie"
            />
            <NEmpty
              v-if="!filePieData.length"
              description="暂无文件数据"
              size="small"
            />
          </div>
        </aside>
      </section>
    </n-el>
  </NSpin>
</template>

<style scoped>
.board {
  display: flex;
  flex-direction: column;
  gap: 14px;
  min-width: 0;
}

.welcome {
  display: grid;
  grid-template-columns: minmax(0, 1.7fr) minmax(260px, 0.9fr);
  gap: 14px;
  align-items: stretch;
}

.welcome--solo {
  grid-template-columns: 1fr;
}

.rail {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 18px;
  min-width: 0;
  padding: 16px 18px;
  background: var(--card-color, #fff);
  border: 1px solid var(--border-color, #eef2f7);
}

.rail__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.rail__identity {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.rail__hello {
  font-size: 18px;
  font-weight: 650;
  line-height: 1.3;
  color: var(--text-color-1, #1f1f1f);
}

.rail__meta {
  margin-top: 4px;
  color: var(--text-color-3, #999);
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rail__pulse {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.pulse {
  display: inline-flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
  min-width: 84px;
  padding: 8px 10px;
  border: 1px solid var(--border-color, #eef2f7);
  background: var(--body-color, #f7f9fc);
  color: inherit;
  cursor: default;
  text-align: left;
}

.pulse--link {
  cursor: pointer;
}

.pulse--link:hover {
  border-color: color-mix(in srgb, var(--primary-color, #1677ff) 45%, transparent);
}

.pulse__label {
  color: var(--text-color-3, #999);
  font-size: 11px;
  letter-spacing: 0.02em;
}

.pulse__value {
  font-size: 20px;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
  line-height: 1.1;
  color: var(--text-color-1, #1f1f1f);
}

.pulse--danger .pulse__value {
  color: #cf1322;
}

.rail__clock {
  color: var(--text-color-2, #666);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.promo {
  min-width: 0;
  overflow: hidden;
  border: 1px solid var(--border-color, #eef2f7);
  background: #0f172a;
}

.promo,
.promo :deep(.n-carousel),
.promo :deep(.n-carousel__slide) {
  height: 100%;
  min-height: 168px;
}

.promo__slide {
  position: relative;
  display: block;
  width: 100%;
  height: 100%;
  min-height: 168px;
  padding: 0;
  border: 0;
  cursor: pointer;
  overflow: hidden;
  background: linear-gradient(145deg, #0f172a, #1d4ed8);
  text-align: left;
}

.promo__image {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.promo__veil {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.15), rgba(15, 23, 42, 0.72));
}

.promo__text {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  gap: 4px;
  height: 100%;
  padding: 14px;
  color: #fff;
}

.promo__text strong {
  font-size: 15px;
  font-weight: 650;
  line-height: 1.35;
}

.promo__text span {
  font-size: 12px;
  line-height: 1.4;
  opacity: 0.88;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.ledger {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  background: var(--card-color, #fff);
  border: 1px solid var(--border-color, #eef2f7);
}

.ledger__cell {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 6px;
  min-width: 0;
  padding: 14px 12px;
  border: 0;
  border-right: 1px solid var(--border-color, #eef2f7);
  background: transparent;
  color: inherit;
  cursor: pointer;
  text-align: left;
}

.ledger__cell:last-child {
  border-right: 0;
}

.ledger__cell:hover {
  background: color-mix(in srgb, var(--primary-color, #1677ff) 6%, transparent);
}

.ledger__label {
  color: var(--text-color-3, #999);
  font-size: 12px;
}

.ledger__value {
  color: var(--text-color-1, #1f1f1f);
  font-size: 24px;
  font-weight: 720;
  font-variant-numeric: tabular-nums;
  line-height: 1;
}

.ledger__note {
  color: var(--text-color-3, #999);
  font-size: 11px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.main {
  display: grid;
  grid-template-columns: minmax(0, 1.7fr) minmax(260px, 1fr);
  gap: 14px;
  align-items: stretch;
}

.panel {
  min-width: 0;
  padding: 14px 16px 16px;
  background: var(--card-color, #fff);
  border: 1px solid var(--border-color, #eef2f7);
}

.panel__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
}

.panel__title {
  color: var(--text-color-1, #1f1f1f);
  font-size: 15px;
  font-weight: 650;
}

.panel__sub {
  margin-top: 2px;
  color: var(--text-color-3, #999);
  font-size: 12px;
}

.trend__chart {
  width: 100%;
  height: 300px;
}

.side {
  display: flex;
  flex-direction: column;
  gap: 14px;
  min-width: 0;
}

.bars {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.bar-row__meta {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 6px;
  color: var(--text-color-2, #666);
  font-size: 12px;
}

.bar-row__meta strong {
  color: var(--text-color-1, #1f1f1f);
  font-variant-numeric: tabular-nums;
}

.bar-row__track {
  height: 6px;
  overflow: hidden;
  background: var(--body-color, #f0f3f8);
}

.bar-row__fill {
  height: 100%;
}

.bar-row__fill--blue {
  background: #1677ff;
}

.file-pie {
  width: 100%;
  height: 220px;
}

@media (max-width: 1100px) {
  .ledger {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .ledger__cell:nth-child(4n) {
    border-right: 0;
  }

  .ledger__cell {
    border-bottom: 1px solid var(--border-color, #eef2f7);
  }
}

@media (max-width: 900px) {
  .welcome,
  .main {
    grid-template-columns: 1fr;
  }

  .promo,
  .promo :deep(.n-carousel),
  .promo :deep(.n-carousel__slide),
  .promo__slide {
    min-height: 148px;
  }
}

@media (max-width: 640px) {
  .ledger {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .ledger__cell:nth-child(2n) {
    border-right: 0;
  }
}
</style>
