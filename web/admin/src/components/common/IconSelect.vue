<!-- Author: Charlie -->

<script setup lang="ts">
import { icons as iconParkOutline } from '@iconify-json/icon-park-outline'
import { computed, onBeforeUnmount, ref, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    value?: string | null
    placeholder?: string
    clearable?: boolean
    disabled?: boolean
    size?: 'tiny' | 'small' | 'medium' | 'large'
    prefix?: string
    columns?: number
    pageSize?: number
    modalWidth?: string | number
  }>(),
  {
    value: '',
    placeholder: '请选择图标',
    clearable: true,
    disabled: false,
    size: 'medium',
    prefix: 'icon-park-outline',
    columns: 8,
    pageSize: 64,
    modalWidth: 720,
  },
)

const emit = defineEmits<{
  'update:value': [value: string]
  change: [value: string]
}>()

const iconNames = Object.keys(iconParkOutline.icons).sort((a, b) => a.localeCompare(b))
const showPanel = ref(false)
const searchKey = ref('')
const deferredSearchKey = ref('')
const page = ref(1)
let searchTimer: number | null = null
const currentIcon = computed(() => normalizeIconValue(props.value))
const columnCount = computed(() => Math.max(1, Math.floor(Number(props.columns) || 1)))
const normalizedPageSize = computed(() => Math.max(1, Math.floor(Number(props.pageSize) || 1)))
const modalStyle = computed(() => ({
  width: toCssSize(props.modalWidth),
  '--icon-select-columns': String(columnCount.value),
}))
const filteredIcons = computed(() => {
  const keyword = deferredSearchKey.value.trim().toLowerCase()
  if (!keyword) {
    return iconNames
  }

  return iconNames.filter((name) => {
    const fullName = toFullIconName(name)
    return name.includes(keyword) || fullName.includes(keyword)
  })
})
const pageCount = computed(() =>
  Math.max(1, Math.ceil(filteredIcons.value.length / normalizedPageSize.value)),
)
const pagedIcons = computed(() => {
  const currentPage = Math.min(page.value, pageCount.value)
  const start = (currentPage - 1) * normalizedPageSize.value
  return filteredIcons.value.slice(start, start + normalizedPageSize.value)
})
const resultText = computed(() => {
  const total = filteredIcons.value.length
  if (!total) {
    return '0 个图标'
  }
  const start = (Math.min(page.value, pageCount.value) - 1) * normalizedPageSize.value + 1
  const end = Math.min(start + pagedIcons.value.length - 1, total)
  return `${start}-${end} / ${total}`
})

function toFullIconName(name: string) {
  return name.includes(':') ? name : `${props.prefix}:${name}`
}

function normalizeIconValue(value?: string | null) {
  const rawValue = String(value ?? '').trim()
  if (!rawValue) {
    return ''
  }
  return toFullIconName(rawValue)
}

function handleUpdateValue(value: string | null) {
  const icon = normalizeIconValue(value)
  emit('update:value', icon)
  emit('change', icon)
}

function selectIcon(name: string) {
  handleUpdateValue(toFullIconName(name))
  showPanel.value = false
}

function clearIcon() {
  if (props.disabled || !props.clearable) {
    return
  }
  handleUpdateValue('')
}

function openPanel() {
  if (!props.disabled) {
    showPanel.value = true
  }
}

function handleUpdatePage(value: number) {
  page.value = value
}

function toCssSize(value: string | number) {
  return typeof value === 'number' ? `${value}px` : value
}

watch(searchKey, (value) => {
  if (searchTimer) {
    window.clearTimeout(searchTimer)
  }
  searchTimer = window.setTimeout(() => {
    deferredSearchKey.value = value
    page.value = 1
  }, 120)
})

watch(showPanel, (value) => {
  if (value) {
    deferredSearchKey.value = searchKey.value
    page.value = 1
  }
})

onBeforeUnmount(() => {
  if (searchTimer) {
    window.clearTimeout(searchTimer)
  }
})
</script>

<template>
  <div class="icon-select">
    <span class="icon-select__preview">
      <NovaIcon
        v-if="currentIcon"
        :icon="currentIcon"
        :size="18"
      />
    </span>
    <NInput
      class="icon-select__input"
      :value="currentIcon"
      :placeholder="placeholder"
      :disabled="disabled"
      :size="size"
      :clearable="clearable"
      @update:value="handleUpdateValue"
      @clear="clearIcon"
    />
    <NButton
      class="icon-select__button"
      :size="size"
      :disabled="disabled"
      title="选择图标"
      aria-label="选择图标"
      @click="openPanel"
    >
      <template #icon>
        <NovaIcon icon="icon-park-outline:all-application" />
      </template>
    </NButton>
  </div>

  <NModal
    v-model:show="showPanel"
    preset="card"
    draggable
    :mask-closable="true"
    :title="'选择图标'"
    :style="modalStyle"
    :segmented="{ content: true }"
  >
    <div class="icon-select-panel">
      <div class="icon-select-panel__header">
        <NInput
          v-model:value="searchKey"
          size="small"
          placeholder="搜索图标名称"
          clearable
        />
        <div class="icon-select-panel__meta">
          <NText depth="3">
            {{ resultText }}
          </NText>
          <NPagination
            v-if="pageCount > 1"
            :page="page"
            :page-count="pageCount"
            size="small"
            simple
            @update:page="handleUpdatePage"
          />
        </div>
      </div>
      <div class="icon-select-panel__scroll">
        <div class="icon-select-panel__grid">
          <button
            v-for="name in pagedIcons"
            :key="name"
            class="icon-select-panel__item"
            :class="{ 'icon-select-panel__item--active': currentIcon === toFullIconName(name) }"
            type="button"
            :title="toFullIconName(name)"
            @click="selectIcon(name)"
          >
            <NovaIcon
              :icon="toFullIconName(name)"
              :size="20"
            />
            <span>{{ name }}</span>
          </button>
        </div>
        <NEmpty
          v-if="!filteredIcons.length"
          size="small"
          description="未找到图标"
        />
      </div>
    </div>
  </NModal>
</template>

<style scoped>
.icon-select {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  min-width: 0;
}

.icon-select__preview {
  display: inline-flex;
  flex: 0 0 32px;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-color-2);
}

.icon-select__input {
  flex: 1;
  min-width: 0;
}

.icon-select__button {
  flex: 0 0 auto;
}

.icon-select-panel {
  display: grid;
  grid-template-rows: auto minmax(0, 1fr);
  gap: 10px;
  min-width: 280px;
  min-height: 0;
  max-height: min(640px, calc(100vh - 180px));
  max-width: 100%;
  overflow-x: hidden;
}

.icon-select-panel__header {
  display: grid;
  gap: 10px;
  min-width: 0;
}

.icon-select-panel__grid {
  display: grid;
  grid-template-columns: repeat(var(--icon-select-columns), minmax(0, 1fr));
  gap: 6px;
  min-width: 0;
  padding-right: 4px;
}

.icon-select-panel__scroll {
  min-height: 0;
  max-height: min(520px, calc(100vh - 300px));
  overflow-y: auto;
  overflow-x: hidden;
  overscroll-behavior: contain;
}

.icon-select-panel__meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-width: 0;
}

.icon-select-panel__item {
  display: grid;
  grid-template-rows: 24px 16px;
  gap: 4px;
  align-items: center;
  justify-items: center;
  min-width: 0;
  height: 54px;
  padding: 5px 4px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--body-color);
  color: var(--text-color-2);
  cursor: pointer;
}

.icon-select-panel__item:hover,
.icon-select-panel__item--active {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.icon-select-panel__item span {
  max-width: 100%;
  overflow: hidden;
  font-size: 11px;
  line-height: 16px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:global(.n-modal .icon-select-panel) {
  max-height: calc(100vh - 180px);
}
</style>
