<!-- Author: Charlie -->

<script setup lang="ts">
defineProps<{
  title?: string
  description?: string
  saving?: boolean
  hideActions?: boolean
}>()

const emit = defineEmits<{
  save: []
  reset: []
}>()
</script>

<template>
  <div class="sys-config-section">
    <div
      v-if="title || description"
      class="sys-config-section__header"
    >
      <h3
        v-if="title"
        class="sys-config-section__title"
      >
        {{ title }}
      </h3>
      <p
        v-if="description"
        class="sys-config-section__desc"
      >
        {{ description }}
      </p>
    </div>

    <div class="sys-config-section__body">
      <slot />
    </div>

    <div
      v-if="!hideActions"
      class="sys-config-section__footer"
    >
      <NButton
        type="primary"
        :loading="saving"
        @click="emit('save')"
      >
        保存
      </NButton>
      <NButton
        :disabled="saving"
        @click="emit('reset')"
      >
        重置
      </NButton>
    </div>
  </div>
</template>

<style scoped>
.sys-config-section__header {
  margin-bottom: 16px;
}

.sys-config-section__title {
  margin: 0 0 6px;
  font-size: 16px;
  font-weight: 600;
  line-height: 1.35;
}

.sys-config-section__desc {
  margin: 0;
  font-size: 13px;
  color: var(--text-color-3);
  line-height: 1.5;
}

.sys-config-section__footer {
  display: flex;
  gap: 12px;
  margin-top: 20px;
}
</style>
