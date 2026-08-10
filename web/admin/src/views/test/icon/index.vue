<!-- Author: Charlie -->

<script setup lang="ts">
import { computed, ref } from 'vue'
import IconSelect from '@/components/common/IconSelect.vue'

const icon = ref('icon-park-outline:people')
const columns = ref(8)
const exampleIcons = [
  'icon-park-outline:analysis',
  'icon-park-outline:setting-two',
  'icon-park-outline:file-code',
  'icon-park-outline:experiment-one',
]
const currentName = computed(() => icon.value.split(':')[1] || '')

function setExampleIcon(value: string) {
  icon.value = value
}

function clearIcon() {
  icon.value = ''
}
</script>

<template>
  <n-el class="icon-test-page">
    <div class="icon-test-header">
      <div class="min-w-0">
        <h1>图标选择器测试</h1>
        <p>验证离线 Iconify 图标搜索、选择、回显和清空。</p>
      </div>
    </div>

    <NGrid
      cols="1 l:24"
      responsive="screen"
      :x-gap="16"
      :y-gap="16"
    >
      <NGridItem span="1 l:14">
        <NCard
          title="选择器"
          :bordered="false"
          class="icon-test-card"
        >
          <NForm
            label-placement="left"
            label-width="100"
          >
            <NFormItem label="图标">
              <IconSelect
                v-model:value="icon"
                :columns="columns"
              />
            </NFormItem>
            <NFormItem label="每行数量">
              <NInputNumber
                v-model:value="columns"
                class="w-full"
                :min="4"
                :max="12"
              />
            </NFormItem>
            <NFormItem label="示例">
              <NSpace>
                <NButton
                  v-for="item in exampleIcons"
                  :key="item"
                  size="small"
                  @click="setExampleIcon(item)"
                >
                  <template #icon>
                    <NovaIcon :icon="item" />
                  </template>
                  {{ item.split(':')[1] }}
                </NButton>
                <NButton
                  size="small"
                  tertiary
                  @click="clearIcon"
                >
                  清空
                </NButton>
              </NSpace>
            </NFormItem>
          </NForm>
        </NCard>
      </NGridItem>

      <NGridItem span="1 l:10">
        <NCard
          title="当前值"
          :bordered="false"
          class="icon-test-card"
        >
          <div class="icon-preview">
            <div class="icon-preview__box">
              <NovaIcon
                v-if="icon"
                :icon="icon"
                :size="40"
              />
            </div>
            <div class="icon-preview__meta">
              <NText depth="3">
                完整值
              </NText>
              <NCode :code="icon || '-'" />
              <NText depth="3">
                图标名
              </NText>
              <NCode :code="currentName || '-'" />
            </div>
          </div>
        </NCard>
      </NGridItem>
    </NGrid>
  </n-el>
</template>

<style scoped>
.icon-test-page {
  min-width: 0;
  min-height: 100%;
}

.icon-test-header {
  margin-bottom: 16px;
}

.icon-test-header h1 {
  margin: 0;
  font-size: 26px;
  line-height: 1.25;
}

.icon-test-header p {
  margin: 8px 0 0;
  color: var(--text-color-3);
}

.icon-test-card {
  height: 100%;
}

.icon-test-card :deep(.n-card__content) {
  min-width: 0;
}

.icon-preview {
  display: grid;
  gap: 16px;
}

.icon-preview__box {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 96px;
  height: 96px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-color-1);
}

.icon-preview__meta {
  display: grid;
  gap: 8px;
  min-width: 0;
}
</style>
