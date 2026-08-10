<!-- Author: Charlie -->

<script setup lang="ts">
import type { EditorUploadedFile } from '@/components/editor'
import { computed, reactive } from 'vue'
import {
  MdEditor,
  MdPreview,
  MonacoDiffPreview,
  MonacoEditor,
  MonacoPreview,
  RichTextEditor,
  RichTextPreview,
} from '@/components/editor'

const defaultMarkdown = `# 编辑器测试

这里用于验证 Markdown 编辑、预览和图片上传。

- 支持基础 Markdown 语法
- 上传图片后会插入文件服务返回的 URL
- 右侧实时预览

\`\`\`ts
const message: string = 'hello editor'
console.log(message)
\`\`\`
`

const defaultRichText = `<h2>富文本测试</h2><p>这里用于验证 wangEditor 5 编辑、预览和图片/视频上传。</p><ul><li>工具栏可编辑内容</li><li>上传文件会走现有文件接口</li></ul>`

const defaultCode = `interface EditorSample {
  id: string
  title: string
  enabled: boolean
}

const sample: EditorSample = {
  id: 'editor-test',
  title: '编辑器测试',
  enabled: true,
}

export default sample
`

const defaultOriginal = `export function sum(a: number, b: number) {
  return a + b
}
`

const defaultModified = `export function sum(a: number, b: number) {
  const result = a + b
  return result
}
`

const state = reactive({
  activeTab: 'markdown',
  markdown: defaultMarkdown,
  richText: defaultRichText,
  code: defaultCode,
  diffOriginal: defaultOriginal,
  diffModified: defaultModified,
  uploadLogs: [] as EditorUploadedFile[],
})

const uploadLogItems = computed(() => [...state.uploadLogs].reverse())

function appendUploadLog(file: EditorUploadedFile) {
  state.uploadLogs.push(file)
}

function clearUploadLogs() {
  state.uploadLogs = []
}

function resetMarkdown() {
  state.markdown = defaultMarkdown
}

function resetRichText() {
  state.richText = defaultRichText
}

function resetCode() {
  state.code = defaultCode
}

function resetDiff() {
  state.diffOriginal = defaultOriginal
  state.diffModified = defaultModified
}
</script>

<template>
  <n-el class="editor-test-page">
    <div class="editor-test-header">
      <div class="min-w-0">
        <h1>编辑器测试</h1>
        <p>验证 Markdown、富文本、Monaco 代码编辑与预览组件。</p>
      </div>
      <NButton
        text
        :disabled="!state.uploadLogs.length"
        title="清空上传记录"
        aria-label="清空上传记录"
        @click="clearUploadLogs"
      >
        <template #icon>
          <NovaIcon icon="icon-park-outline:clear" />
        </template>
      </NButton>
    </div>

    <NGrid
      cols="1 xl:24"
      responsive="screen"
      :x-gap="16"
      :y-gap="16"
    >
      <NGridItem span="1 xl:18">
        <NCard
          :bordered="false"
          class="editor-test-card"
        >
          <NTabs
            v-model:value="state.activeTab"
            type="segment"
            animated
          >
            <NTabPane
              name="markdown"
              tab="Markdown"
            >
              <div class="pane-toolbar">
                <NButton
                  size="small"
                  @click="resetMarkdown"
                >
                  <template #icon>
                    <NovaIcon icon="icon-park-outline:refresh" />
                  </template>
                  重置
                </NButton>
              </div>
              <NGrid
                cols="1 l:2"
                responsive="screen"
                :x-gap="16"
                :y-gap="16"
              >
                <NGridItem>
                  <MdEditor
                    v-model:value="state.markdown"
                    height="520px"
                    @uploaded="appendUploadLog"
                  />
                </NGridItem>
                <NGridItem>
                  <div class="preview-panel">
                    <MdPreview :value="state.markdown" />
                  </div>
                </NGridItem>
              </NGrid>
            </NTabPane>

            <NTabPane
              name="richText"
              tab="富文本"
            >
              <div class="pane-toolbar">
                <NButton
                  size="small"
                  @click="resetRichText"
                >
                  <template #icon>
                    <NovaIcon icon="icon-park-outline:refresh" />
                  </template>
                  重置
                </NButton>
              </div>
              <NGrid
                cols="1 l:2"
                responsive="screen"
                :x-gap="16"
                :y-gap="16"
              >
                <NGridItem>
                  <RichTextEditor
                    v-model:value="state.richText"
                    height="460px"
                    @uploaded="appendUploadLog"
                  />
                </NGridItem>
                <NGridItem>
                  <div class="preview-panel">
                    <RichTextPreview :value="state.richText" />
                  </div>
                </NGridItem>
              </NGrid>
            </NTabPane>

            <NTabPane
              name="monaco"
              tab="Monaco"
            >
              <div class="pane-toolbar">
                <NButton
                  size="small"
                  @click="resetCode"
                >
                  <template #icon>
                    <NovaIcon icon="icon-park-outline:refresh" />
                  </template>
                  重置
                </NButton>
              </div>
              <NGrid
                cols="1 l:2"
                responsive="screen"
                :x-gap="16"
                :y-gap="16"
              >
                <NGridItem>
                  <MonacoEditor
                    v-model:value="state.code"
                    language="typescript"
                    height="520px"
                  />
                </NGridItem>
                <NGridItem>
                  <MonacoPreview
                    :value="state.code"
                    language="typescript"
                    height="520px"
                  />
                </NGridItem>
              </NGrid>
            </NTabPane>

            <NTabPane
              name="diff"
              tab="Diff"
            >
              <div class="pane-toolbar">
                <NButton
                  size="small"
                  @click="resetDiff"
                >
                  <template #icon>
                    <NovaIcon icon="icon-park-outline:refresh" />
                  </template>
                  重置
                </NButton>
              </div>
              <NGrid
                class="mb-16px"
                cols="1 l:2"
                responsive="screen"
                :x-gap="16"
                :y-gap="16"
              >
                <NGridItem>
                  <NInput
                    v-model:value="state.diffOriginal"
                    type="textarea"
                    :autosize="{ minRows: 8, maxRows: 12 }"
                  />
                </NGridItem>
                <NGridItem>
                  <NInput
                    v-model:value="state.diffModified"
                    type="textarea"
                    :autosize="{ minRows: 8, maxRows: 12 }"
                  />
                </NGridItem>
              </NGrid>
              <MonacoDiffPreview
                :original="state.diffOriginal"
                :modified="state.diffModified"
                language="typescript"
                height="520px"
              />
            </NTabPane>
          </NTabs>
        </NCard>
      </NGridItem>

      <NGridItem span="1 xl:6">
        <NCard
          title="上传记录"
          :bordered="false"
          class="editor-test-card upload-card"
        >
          <NEmpty
            v-if="!uploadLogItems.length"
            description="暂无上传记录"
          />
          <NList v-else>
            <NListItem
              v-for="(item, index) in uploadLogItems"
              :key="`${item.value}-${index}`"
            >
              <NThing :title="item.name">
                <template #description>
                  <div class="upload-detail">
                    <span v-if="item.contentType">{{ item.contentType }}</span>
                    <NEllipsis class="upload-url">
                      {{ item.url }}
                    </NEllipsis>
                  </div>
                </template>
              </NThing>
            </NListItem>
          </NList>
        </NCard>
      </NGridItem>
    </NGrid>
  </n-el>
</template>

<style scoped>
.editor-test-page {
  min-width: 0;
  min-height: 100%;
}

.editor-test-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.editor-test-header h1 {
  margin: 0;
  font-size: 26px;
  line-height: 1.25;
}

.editor-test-header p {
  margin: 8px 0 0;
  color: var(--text-color-3);
}

.editor-test-card {
  height: 100%;
}

.editor-test-card :deep(.n-card__content) {
  min-width: 0;
}

.pane-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 12px;
}

.preview-panel {
  width: 100%;
  min-width: 0;
  height: 520px;
  overflow: auto;
  border: 1px solid var(--border-color);
}

.upload-card :deep(.n-card__content) {
  max-height: calc(100vh - 220px);
  overflow: auto;
}

.upload-url {
  max-width: 100%;
  color: var(--text-color-3);
  font-size: 12px;
}

.upload-detail {
  display: grid;
  gap: 4px;
  min-width: 0;
}
</style>
