<!-- Author: Charlie -->

<script setup lang="ts">
defineOptions({
  inheritAttrs: false,
})

const modelValue = defineModel<string | null>({
  default: null,
})

const props = withDefaults(
  defineProps<{
    presetColors?: string[]
    disabled?: boolean
    size?: 'small' | 'medium' | 'large'
  }>(),
  {
    presetColors: () => [
      '#18a058',
      '#36ad6a',
      '#0f766e',
      '#14b8a6',
      '#2080f0',
      '#4098fc',
      '#2563eb',
      '#06b6d4',
      '#f0a020',
      '#f59e0b',
      '#f97316',
      '#d03050',
      '#ef4444',
      '#dc2626',
      '#7c3aed',
      '#8b5cf6',
      '#a855f7',
      '#ec4899',
      '#64748b',
      '#475569',
      '#111827',
      '#000000',
      '#ffffff',
    ],
    disabled: false,
    size: 'medium',
  },
)

function selectPresetColor(color: string) {
  if (!props.disabled) {
    modelValue.value = color
  }
}
</script>

<template>
  <div class="common-color-picker">
    <NColorPicker
      v-model:value="modelValue"
      :modes="['hex']"
      :show-alpha="false"
      :actions="['clear']"
      :swatches="presetColors"
      :disabled="disabled"
      :size="size"
      v-bind="$attrs"
    />

    <div class="common-color-picker__swatches">
      <button
        v-for="color in presetColors"
        :key="color"
        class="common-color-picker__swatch"
        :class="{ 'common-color-picker__swatch--active': modelValue === color }"
        type="button"
        :aria-label="color"
        :title="color"
        :disabled="disabled"
        :style="{ backgroundColor: color }"
        @click="selectPresetColor(color)"
      />
    </div>
  </div>
</template>

<style scoped>
.common-color-picker {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.common-color-picker__swatches {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.common-color-picker__swatch {
  width: 24px;
  height: 24px;
  padding: 0;
  border: 1px solid var(--n-border-color, #d9d9d9);
  border-radius: 4px;
  cursor: pointer;
}

.common-color-picker__swatch:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.common-color-picker__swatch--active {
  outline: 2px solid var(--n-primary-color, #1677ff);
  outline-offset: 2px;
}
</style>
