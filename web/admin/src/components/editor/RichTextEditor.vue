<!-- Author: Charlie -->

<script setup lang="ts">
import type { IDomEditor, IEditorConfig, IToolbarConfig } from '@wangeditor/editor'
import type { EditorUploadedFile } from './shared'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { computed, onBeforeUnmount, shallowRef, watch } from 'vue'
import { toCssSize, uploadEditorFile } from './shared'
import '@wangeditor/editor/dist/css/style.css'

type InsertFileFn = (...args: string[]) => void
type RichTextEditorConfig = Partial<IEditorConfig> & {
  MENU_CONF?: Record<string, unknown>
}

const props = withDefaults(
  defineProps<{
    value?: string | null
    height?: string | number
    placeholder?: string
    mode?: 'default' | 'simple'
    readOnly?: boolean
    disabled?: boolean
    editorConfig?: RichTextEditorConfig
    toolbarConfig?: Partial<IToolbarConfig>
  }>(),
  {
    value: '',
    height: 360,
    placeholder: '请输入内容',
    mode: 'default',
    readOnly: false,
    disabled: false,
    editorConfig: () => ({}),
    toolbarConfig: () => ({}),
  },
)

const emit = defineEmits<{
  'update:value': [value: string]
  change: [value: string]
  uploaded: [file: EditorUploadedFile]
  'upload-error': [error: unknown]
}>()

const editorRef = shallowRef<IDomEditor | null>(null)
const modelValue = computed(() => props.value ?? '')
const editorStyle = computed(() => ({
  height: toCssSize(props.height),
  overflowY: 'hidden',
}))

const resolvedToolbarConfig = computed<Partial<IToolbarConfig>>(() => ({
  ...props.toolbarConfig,
}))

const resolvedEditorConfig = computed<RichTextEditorConfig>(() => {
  const menuConf = {
    ...((props.editorConfig.MENU_CONF ?? {}) as Record<string, unknown>),
    uploadImage: {
      ...(((props.editorConfig.MENU_CONF as Record<string, unknown> | undefined)?.uploadImage ??
        {}) as Record<string, unknown>),
      customUpload: async (file: File, insertFn: InsertFileFn) => {
        await uploadAndInsert(file, insertFn)
      },
    },
    uploadVideo: {
      ...(((props.editorConfig.MENU_CONF as Record<string, unknown> | undefined)?.uploadVideo ??
        {}) as Record<string, unknown>),
      customUpload: async (file: File, insertFn: InsertFileFn) => {
        await uploadAndInsert(file, insertFn)
      },
    },
  }

  return {
    ...props.editorConfig,
    placeholder: props.placeholder,
    MENU_CONF: menuConf,
  }
})

async function uploadAndInsert(file: File, insertFn: InsertFileFn) {
  try {
    const uploaded = await uploadEditorFile(file)
    if (uploaded.url) {
      insertFn(uploaded.url, uploaded.name, uploaded.url)
    }
    emit('uploaded', uploaded)
  } catch (error) {
    emit('upload-error', error)
    window.$message?.error('上传失败')
  }
}

function syncDisabledState() {
  const editor = editorRef.value
  if (!editor) {
    return
  }
  if (props.disabled || props.readOnly) {
    editor.disable()
  } else {
    editor.enable()
  }
}

function handleCreated(editor: IDomEditor) {
  editorRef.value = editor
  syncDisabledState()
}

function updateValue(value: string) {
  emit('update:value', value)
  emit('change', value)
}

watch(() => [props.disabled, props.readOnly], syncDisabledState)

onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor) {
    editor.destroy()
  }
  editorRef.value = null
})
</script>

<template>
  <div class="rich-text-editor">
    <Toolbar
      class="rich-text-editor__toolbar"
      :editor="editorRef"
      :default-config="resolvedToolbarConfig"
      :mode="mode"
    />
    <Editor
      class="rich-text-editor__body"
      :model-value="modelValue"
      :default-config="resolvedEditorConfig"
      :mode="mode"
      :style="editorStyle"
      @on-created="handleCreated"
      @update:model-value="updateValue"
    />
  </div>
</template>

<style scoped>
.rich-text-editor {
  width: 100%;
  min-width: 0;
  border: 1px solid var(--border-color);
}

.rich-text-editor__toolbar {
  border-bottom: 1px solid var(--border-color);
}

.rich-text-editor__body {
  min-height: 0;
}
</style>
