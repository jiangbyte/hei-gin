<!-- Author: Charlie -->

<template>
  <Layout title="字典测试" :back="true">
    <view class="px-4 py-4">
      <!-- 刷新按钮 -->
      <view class="flex items-center justify-between mb-4">
        <text class="text-base font-semibold text-gray-800">字典工具测试</text>
        <view class="flex gap-2">
          <u-button
              text="刷新字典"
              size="mini"
              type="primary"
              :loading="loading"
              @click="reloadDict"
          />
        </view>
      </view>

      <!-- 状态提示 -->
      <view class="bg-blue-50 rounded px-3 py-2 mb-4 text-xs text-blue-600">
        当前字典条目数：{{ dictCount }} | 以下使用 sys_dict.sql 中的字典 KEY
        测试
      </view>

      <!-- 1. dictTypeList 测试 - 列出字典项 -->
      <view class="bg-white rounded p-3 mb-3">
        <text class="text-sm font-semibold text-gray-700 mb-2 block"
        >1. dictTypeList('COMMON_STATUS')
        </text
        >
        <view class="flex flex-wrap gap-2">
          <view
              v-for="item in commonStatusItems"
              :key="item.value"
              class="rounded px-2.5 py-1 text-xs"
              :style="{
              backgroundColor: item.color + '20',
              color: item.color,
              border: '1px solid ' + item.color,
            }"
          >
            {{ item.label }} ({{ item.value }})
          </view>
        </view>
      </view>

      <!-- 2. dictTypeList 测试 - RESOURCE_TYPE -->
      <view class="bg-white rounded p-3 mb-3">
        <text class="text-sm font-semibold text-gray-700 mb-2 block"
        >2. dictTypeList('RESOURCE_TYPE')
        </text
        >
        <u-table border="1" :show-summary="false">
          <u-tr>
            <u-th>标签</u-th>
            <u-th>值</u-th>
            <u-th>颜色</u-th>
          </u-tr>
          <u-tr v-for="item in resourceTypeItems" :key="item.value">
            <u-td>{{ item.label }}</u-td>
            <u-td
            >
              <text class="text-gray-500">{{ item.value }}</text>
            </u-td
            >
            <u-td>
              <view
                  class="inline-block rounded px-2 py-0.5 text-xs"
                  :style="{
                  backgroundColor: item.color + '20',
                  color: item.color,
                }"
              >
                {{ item.color }}
              </view>
            </u-td>
          </u-tr>
        </u-table>
      </view>

      <!-- 3. dictTypeData 测试 - 值转标签 -->
      <view class="bg-white rounded p-3 mb-3">
        <text class="text-sm font-semibold text-gray-700 mb-2 block"
        >3. dictTypeData('ACCOUNT_STATUS', value) — 值转标签
        </text
        >
        <view class="space-y-2">
          <view
              v-for="item in accountStatusItems"
              :key="item.value"
              class="flex items-center gap-3 text-xs"
          >
            <text class="w-20 text-gray-500">{{ item.value }}</text>
            <u-icon name="arrow-right" size="12" color="#ccc"/>
            <text class="font-medium" :style="{ color: item.color }">{{
                item.label
              }}
            </text>
          </view>
        </view>
      </view>

      <!-- 4. dictTypeColor 测试 - 获取颜色 -->
      <view class="bg-white rounded p-3 mb-3">
        <text class="text-sm font-semibold text-gray-700 mb-2 block"
        >4. dictTypeColor('NOTIFICATION_SEVERITY', value) — 获取颜色
        </text
        >
        <view class="flex flex-wrap gap-3">
          <view
              v-for="item in severityItems"
              :key="item.value"
              class="rounded px-3 py-1.5 text-xs text-white"
              :style="{ backgroundColor: item.color }"
          >
            {{ item.label }}
          </view>
        </view>
      </view>

      <!-- 5. dictList 测试 - 下拉选择框 -->
      <view class="bg-white rounded p-3 mb-3">
        <text class="text-sm font-semibold text-gray-700 mb-2 block"
        >5. dictList('DATA_SCOPE') — 下拉选择
        </text
        >
        <u-select
            v-model="showSelect"
            :list="dataScopeList"
            @confirm="onSelectConfirm"
        />
        <u-button
            text="选择数据范围"
            size="mini"
            type="info"
            @click="showSelect = true"
        />
        <text class="text-xs text-gray-500 ml-2"
        >已选：{{ selectedScope || '未选择' }}
        </text
        >
      </view>

      <!-- 6. 原始字典树结构 -->
      <view class="bg-white rounded p-3 mb-3">
        <text class="text-sm font-semibold text-gray-700 mb-2 block"
        >6. dictDataAll() — 原始字典树
        </text
        >
        <view class="bg-gray-50 rounded p-2 max-h-40 overflow-y-auto">
          <text class="text-xs text-gray-400 font-mono leading-relaxed">{{
              dictTreePreview
            }}
          </text>
        </view>
      </view>
    </view>
  </Layout>
</template>

<script setup lang="ts">
import {computed, ref} from 'vue'
import {onShow} from '@dcloudio/uni-app'
import Layout from '@/layouts/index.vue'
import {
  refreshDict,
  dictDataAll,
  dictTypeList,
  dictTypeData,
  dictTypeColor,
  dictList,
  isDictLoaded,
} from '@/utils/dict'

const loading = ref(false)
const showSelect = ref(false)
const selectedScope = ref('')

const dictCount = computed(() => dictDataAll().length)

const commonStatusItems = computed(() =>
    dictTypeList('COMMON_STATUS').map((item: any) => ({
      label: item.label || item.code,
      value: item.value || item.code,
      color: item.color || '#999',
    }))
)

const resourceTypeItems = computed(() =>
    dictTypeList('RESOURCE_TYPE').map((item: any) => ({
      label: item.label || item.code,
      value: item.value || item.code,
      color: item.color || '#999',
    }))
)

const accountStatusItems = computed(() =>
    dictTypeList('ACCOUNT_STATUS').map((item: any) => ({
      value: item.value || item.code,
      label:
          dictTypeData('ACCOUNT_STATUS', item.value || item.code) ||
          item.label ||
          item.code,
      color:
          dictTypeColor('ACCOUNT_STATUS', item.value || item.code) ||
          item.color ||
          '#999',
    }))
)

const severityItems = computed(() =>
    dictTypeList('NOTIFICATION_SEVERITY').map((item: any) => ({
      label: item.label || item.code,
      value: item.value || item.code,
      color:
          dictTypeColor('NOTIFICATION_SEVERITY', item.value || item.code) ||
          item.color ||
          '#999',
    }))
)

const dataScopeList = computed(() => dictList('DATA_SCOPE'))

const dictTreePreview = computed(() => {
  const tree = dictDataAll()
  if (!tree.length) return '（空）'
  return JSON.stringify(
      tree.slice(0, 5).map((n: any) => ({
        code: n.code,
        label: n.label,
        childrenCount: n.children?.length || 0,
      })),
      null,
      2
  )
})

onShow(async () => {
  if (!isDictLoaded()) {
    await doRefresh()
  }
})

async function reloadDict() {
  loading.value = true
  try {
    await doRefresh()
    uni.showToast({title: '字典已刷新', icon: 'success'})
  } catch {
    uni.showToast({title: '字典刷新失败', icon: 'none'})
  } finally {
    loading.value = false
  }
}

async function doRefresh() {
  await refreshDict()
}

function onSelectConfirm(e: any) {
  selectedScope.value = e?.[0]?.label || e?.label || ''
}
</script>

<style lang="scss" scoped>
page {
  background-color: #f5f6f7;
}
</style>
