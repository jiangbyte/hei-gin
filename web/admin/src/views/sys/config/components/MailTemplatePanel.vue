<!-- Author: Charlie -->

<script setup lang="ts">
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import { computed, h, onMounted, reactive, ref } from 'vue'
import { NButton, NFlex, NTag } from 'naive-ui'
import { configApi } from '@/api'
import { MonacoEditor } from '@/components/editor'
import {
  MAIL_SCENE_LABELS,
  SCENE_PATTERN,
  normalizeScene,
  parseMailTemplate,
  stringifyMailTemplate,
  templateConfigKey,
} from '../composables/useTemplateConfig'

const CATEGORY = 'MAIL_TEMPLATE'

interface Row {
  id: string
  scene: string
  label: string
  subject: string
  body: string
  is_builtin: boolean
  sort_code: number
  config_key: string
}

const formRef = ref<FormInst | null>(null)

const state = reactive({
  loading: false,
  saving: false,
  rows: [] as Row[],
  drawerShow: false,
  isCreate: false,
  editingId: '' as string,
  draftScene: '',
  draftLabel: '',
  draftSubject: '',
  draftBody: '',
  draftBuiltin: false,
})

const rules = computed<FormRules>(() => ({
  draftScene: [
    { required: true, message: '请输入场景编码', trigger: ['blur', 'input'] },
    {
      validator: () => SCENE_PATTERN.test(normalizeScene(state.draftScene)),
      message: '场景编码须为 UPPER_SNAKE',
      trigger: ['blur', 'input'],
    },
  ],
  draftLabel: [{ required: true, message: '请输入名称', trigger: ['blur', 'input'] }],
}))

onMounted(() => {
  void reload()
})

async function reload() {
  state.loading = true
  try {
    const res = await configApi.list({ category: CATEGORY })
    state.rows = (res.data ?? []).map((item: any) => {
      const parsed = parseMailTemplate(item.config_value)
      return {
        id: String(item.id),
        scene: String(item.scene || ''),
        label: String(item.label || item.scene || item.config_key),
        subject: parsed.subject,
        body: parsed.body,
        is_builtin: !!item.is_builtin,
        sort_code: Number(item.sort_code || 0),
        config_key: String(item.config_key),
      }
    })
  } finally {
    state.loading = false
  }
}

const columns = computed<DataTableColumns<Row>>(() => [
  {
    title: '场景',
    key: 'scene',
    width: 200,
    ellipsis: { tooltip: true },
  },
  { title: '名称', key: 'label', width: 160, ellipsis: { tooltip: true } },
  { title: '主题', key: 'subject', ellipsis: { tooltip: true } },
  {
    title: '是否内置',
    key: 'is_builtin',
    width: 100,
    render: (row) =>
      h(
        NTag,
        {
          size: 'small',
          type: row.is_builtin ? 'info' : 'default',
          bordered: false,
        },
        () => (row.is_builtin ? '是' : '否'),
      ),
  },
  {
    title: '操作',
    key: 'actions',
    width: 140,
    render: (row) =>
      h(NFlex, { size: 12 }, () => [
        h(NButton, { text: true, type: 'primary', onClick: () => openEdit(row) }, () => '编辑'),
        h(
          NButton,
          {
            text: true,
            type: 'error',
            disabled: row.is_builtin,
            title: row.is_builtin ? '内置模板不可删除' : '删除',
            onClick: () => confirmDelete(row),
          },
          () => '删除',
        ),
      ]),
  },
])

function openCreate() {
  state.isCreate = true
  state.editingId = ''
  state.draftScene = ''
  state.draftLabel = ''
  state.draftSubject = ''
  state.draftBody = ''
  state.draftBuiltin = false
  state.drawerShow = true
}

function openEdit(row: Row) {
  state.isCreate = false
  state.editingId = row.id
  state.draftScene = row.scene
  state.draftLabel = row.label
  state.draftSubject = row.subject
  state.draftBody = row.body
  state.draftBuiltin = row.is_builtin
  state.drawerShow = true
}

function onSceneInput(value: string) {
  state.draftScene = value
  const scene = normalizeScene(value)
  if (!state.draftLabel && MAIL_SCENE_LABELS[scene]) {
    state.draftLabel = MAIL_SCENE_LABELS[scene]
  }
}

async function saveDraft() {
  await formRef.value?.validate()
  const scene = normalizeScene(state.draftScene)
  if (!scene) {
    window.$message.warning('请输入场景编码')
    return
  }
  if (state.isCreate && state.rows.some((r) => r.scene === scene)) {
    window.$message.warning('该场景已存在')
    return
  }
  state.saving = true
  try {
    const config_key = templateConfigKey('MAIL_TEMPLATE', scene)
    const config_value = stringifyMailTemplate({
      subject: state.draftSubject,
      body: state.draftBody,
    })
    if (state.isCreate) {
      await configApi.create({
        config_key,
        config_value,
        category: CATEGORY,
        remark: scene,
        sort_code: (state.rows.at(-1)?.sort_code || 0) + 1,
        value_type: 'JSON',
        label: state.draftLabel.trim(),
        scope: null,
        scene,
        is_builtin: false,
        ext_json: {},
      })
    } else {
      await configApi.update({
        id: state.editingId,
        config_key,
        config_value,
        category: CATEGORY,
        remark: scene,
        sort_code: state.rows.find((r) => r.id === state.editingId)?.sort_code || 0,
        value_type: 'JSON',
        label: state.draftLabel.trim(),
        scope: null,
        scene: state.draftBuiltin ? state.draftScene : scene,
        is_builtin: state.draftBuiltin,
        ext_json: {},
      })
    }
    window.$message.success('保存成功')
    state.drawerShow = false
    await reload()
  } finally {
    state.saving = false
  }
}

function confirmDelete(row: Row) {
  if (row.is_builtin) {
    window.$message.warning('内置模板不可删除')
    return
  }
  window.$dialog.warning({
    title: '删除模板',
    content: `确认删除场景「${row.scene}」？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await configApi.remove({ ids: [row.id] })
      window.$message.success('已删除')
      await reload()
    },
  })
}
</script>

<template>
  <NSpin :show="state.loading">
    <NFlex
      justify="end"
      class="mb-12px"
    >
      <NButton
        type="primary"
        @click="openCreate"
      >
        新建模板
      </NButton>
    </NFlex>

    <NDataTable
      :columns="columns"
      :data="state.rows"
      :bordered="false"
      :single-line="false"
    />

    <NDrawer
      v-model:show="state.drawerShow"
      :width="720"
      placement="right"
    >
      <NDrawerContent
        :title="state.isCreate ? '新建邮件模板' : '编辑邮件模板'"
        closable
      >
        <NForm
          ref="formRef"
          label-placement="top"
          :model="state"
          :rules="rules"
        >
          <NFormItem
            label="场景编码"
            path="draftScene"
          >
            <NInput
              :value="state.draftScene"
              :disabled="state.draftBuiltin"
              placeholder="如 LOGIN_CODE"
              @update:value="onSceneInput"
            />
          </NFormItem>
          <NFormItem
            label="名称"
            path="draftLabel"
          >
            <NInput
              v-model:value="state.draftLabel"
              placeholder="展示名称"
            />
          </NFormItem>
          <NFormItem label="主题">
            <NInput
              v-model:value="state.draftSubject"
              placeholder="支持 {{变量}} 占位"
            />
          </NFormItem>
          <NFormItem
            label="正文"
            :show-feedback="false"
          >
            <div class="template-editor-field">
              <div class="template-editor-field__monaco">
                <MonacoEditor
                  v-model:value="state.draftBody"
                  language="html"
                  theme="vs-dark"
                  height="420px"
                  :options="{ wordWrap: 'on', lineNumbers: 'on', tabSize: 2 }"
                />
              </div>
              <p class="sys-config__hint">
                支持 HTML 与 &#123;&#123;变量&#125;&#125; 占位，如
                &#123;&#123;validCode&#125;&#125;、&#123;&#123;sysName&#125;&#125;
              </p>
            </div>
          </NFormItem>
        </NForm>
        <template #footer>
          <NSpace>
            <NButton @click="state.drawerShow = false">
              取消
            </NButton>
            <NButton
              type="primary"
              :loading="state.saving"
              @click="saveDraft"
            >
              保存
            </NButton>
          </NSpace>
        </template>
      </NDrawerContent>
    </NDrawer>
  </NSpin>
</template>

<style scoped>
.mb-12px {
  margin-bottom: 12px;
}

.template-editor-field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
  min-width: 0;
}

.template-editor-field__monaco {
  width: 100%;
  border-radius: 6px;
  overflow: hidden;
}

.template-editor-field__monaco :deep(.monaco-editor) {
  border: 0;
}
</style>
