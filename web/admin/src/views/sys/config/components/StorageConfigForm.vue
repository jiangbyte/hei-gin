<!-- Author: Charlie -->

<script setup lang="ts">
import { useMessage } from 'naive-ui'
import { onMounted, reactive } from 'vue'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import { FILES_PUBLIC_PATH } from '@/constants/api'
import { loadByCategory, parseBool, saveByKeys, toBoolStr } from '../composables/useConfigForm'

const CATEGORY = 'STORAGE'
type Engine = 'LOCAL' | 'ALIYUN' | 'TENCENT' | 'MINIO' | 'RUSTFS'

const engineOptions = [
  { label: '本地文件', value: 'LOCAL' as Engine },
  { label: '阿里云 OSS', value: 'ALIYUN' as Engine },
  { label: '腾讯云 COS', value: 'TENCENT' as Engine },
  { label: 'MinIO', value: 'MINIO' as Engine },
  { label: 'RustFS', value: 'RUSTFS' as Engine },
]

type CloudForm = {
  accessKey: string
  secretKey: string
  accessKeySet: boolean
  secretKeySet: boolean
  endpoint: string
  bucket: string
  region: string
  useSsl: boolean
  baseUrl: string
  publicPath: string
}

function emptyCloud(defaults: Partial<CloudForm> = {}): CloudForm {
  return {
    accessKey: '',
    secretKey: '',
    accessKeySet: false,
    secretKeySet: false,
    endpoint: '',
    bucket: 'defaultbucket',
    region: '',
    useSsl: false,
    baseUrl: '',
    publicPath: FILES_PUBLIC_PATH,
    ...defaults,
  }
}

const emit = defineEmits<{ saved: [] }>()
const message = useMessage()

const state = reactive({
  loading: false,
  saving: false,
  subTab: 'LOCAL' as Engine,
  defaultEngine: 'MINIO' as Engine,
  local: {
    localRoot: '/defaultUploadFolder',
    windowsRoot: 'D:/defaultUploadFolder',
    publicPath: FILES_PUBLIC_PATH,
    baseUrl: '',
  },
  aliyun: emptyCloud({
    endpoint: 'oss-cn-hangzhou.aliyuncs.com',
    region: 'cn-hangzhou',
    useSsl: true,
  }),
  tencent: emptyCloud({
    region: 'ap-beijing',
    useSsl: true,
  }),
  minio: emptyCloud({
    endpoint: 'https://play.min.io',
  }),
  rustfs: emptyCloud({
    endpoint: 'http://127.0.0.1:9002',
    region: 'us-east-1',
    useSsl: false,
  }),
  snapshot: '',
})

onMounted(() => {
  void reload()
})

function loadCloud(
  target: CloudForm,
  map: Record<string, string>,
  prefix: string,
  defaults: CloudForm,
) {
  target.accessKey = ''
  target.secretKey = ''
  target.accessKeySet = parseBool(map[`${prefix}_ACCESS_KEY_SET`])
  target.secretKeySet = parseBool(map[`${prefix}_SECRET_KEY_SET`])
  target.endpoint = map[`${prefix}_ENDPOINT`] || defaults.endpoint
  target.bucket = map[`${prefix}_BUCKET`] || defaults.bucket
  target.region = map[`${prefix}_REGION`] || defaults.region
  target.useSsl = map[`${prefix}_USE_SSL`] ? parseBool(map[`${prefix}_USE_SSL`]) : defaults.useSsl
  target.baseUrl = map[`${prefix}_BASE_URL`] || ''
  target.publicPath = map[`${prefix}_PUBLIC_PATH`] || defaults.publicPath
}

function cloudKeys(prefix: string, form: CloudForm) {
  return [
    { config_key: `${prefix}_ACCESS_KEY`, config_value: form.accessKey, category: CATEGORY },
    { config_key: `${prefix}_SECRET_KEY`, config_value: form.secretKey, category: CATEGORY },
    { config_key: `${prefix}_ENDPOINT`, config_value: form.endpoint, category: CATEGORY },
    { config_key: `${prefix}_BUCKET`, config_value: form.bucket, category: CATEGORY },
    { config_key: `${prefix}_REGION`, config_value: form.region, category: CATEGORY },
    { config_key: `${prefix}_USE_SSL`, config_value: toBoolStr(form.useSsl), category: CATEGORY },
    { config_key: `${prefix}_BASE_URL`, config_value: form.baseUrl, category: CATEGORY },
    { config_key: `${prefix}_PUBLIC_PATH`, config_value: form.publicPath, category: CATEGORY },
  ]
}

async function reload() {
  state.loading = true
  try {
    const map = await loadByCategory(CATEGORY)
    state.defaultEngine = (map.DEFAULT_FILE_ENGINE || 'MINIO') as Engine
    state.local.localRoot = map.STORAGE_LOCAL_LOCAL_ROOT || state.local.localRoot
    state.local.windowsRoot = map.STORAGE_LOCAL_WINDOWS_ROOT || state.local.windowsRoot
    state.local.publicPath = map.STORAGE_LOCAL_PUBLIC_PATH || state.local.publicPath
    state.local.baseUrl = map.STORAGE_LOCAL_BASE_URL || ''

    loadCloud(
      state.aliyun,
      map,
      'STORAGE_ALIYUN',
      emptyCloud({
        endpoint: 'oss-cn-hangzhou.aliyuncs.com',
        region: 'cn-hangzhou',
        useSsl: true,
      }),
    )
    loadCloud(
      state.tencent,
      map,
      'STORAGE_TENCENT',
      emptyCloud({
        region: 'ap-beijing',
        useSsl: true,
      }),
    )
    loadCloud(
      state.minio,
      map,
      'STORAGE_MINIO',
      emptyCloud({
        endpoint: 'https://play.min.io',
      }),
    )
    loadCloud(
      state.rustfs,
      map,
      'STORAGE_RUSTFS',
      emptyCloud({
        endpoint: 'http://127.0.0.1:9002',
        region: 'us-east-1',
        useSsl: false,
      }),
    )

    if (engineOptions.some((o) => o.value === state.defaultEngine)) {
      state.subTab = state.defaultEngine
    }
    state.snapshot = JSON.stringify({
      defaultEngine: state.defaultEngine,
      local: state.local,
      aliyun: state.aliyun,
      tencent: state.tencent,
      minio: state.minio,
      rustfs: state.rustfs,
    })
  } finally {
    state.loading = false
  }
}

function reset() {
  if (!state.snapshot) return
  const data = JSON.parse(state.snapshot)
  state.defaultEngine = data.defaultEngine
  Object.assign(state.local, data.local)
  Object.assign(state.aliyun, data.aliyun)
  Object.assign(state.tencent, data.tencent)
  Object.assign(state.minio, data.minio)
  Object.assign(state.rustfs, data.rustfs)
}

async function save() {
  state.saving = true
  try {
    await saveByKeys([
      {
        config_key: 'DEFAULT_FILE_ENGINE',
        config_value: state.defaultEngine,
        category: CATEGORY,
      },
      {
        config_key: 'STORAGE_LOCAL_LOCAL_ROOT',
        config_value: state.local.localRoot,
        category: CATEGORY,
      },
      {
        config_key: 'STORAGE_LOCAL_WINDOWS_ROOT',
        config_value: state.local.windowsRoot,
        category: CATEGORY,
      },
      {
        config_key: 'STORAGE_LOCAL_PUBLIC_PATH',
        config_value: state.local.publicPath,
        category: CATEGORY,
      },
      {
        config_key: 'STORAGE_LOCAL_BASE_URL',
        config_value: state.local.baseUrl,
        category: CATEGORY,
      },
      ...cloudKeys('STORAGE_ALIYUN', state.aliyun),
      ...cloudKeys('STORAGE_TENCENT', state.tencent),
      ...cloudKeys('STORAGE_MINIO', state.minio),
      ...cloudKeys('STORAGE_RUSTFS', state.rustfs),
    ])
    message.success('保存成功')
    await reload()
    emit('saved')
  } finally {
    state.saving = false
  }
}
</script>

<template>
  <NSpin :show="state.loading">
    <NForm
      class="sys-config-form mb-12px"
      label-placement="top"
    >
      <NFormItem label="默认文件存储">
        <NRadioGroup v-model:value="state.defaultEngine">
          <NSpace>
            <NRadio
              v-for="opt in engineOptions"
              :key="opt.value"
              :value="opt.value"
              :label="opt.label"
            />
          </NSpace>
        </NRadioGroup>
      </NFormItem>
    </NForm>

    <NTabs
      v-model:value="state.subTab"
      type="line"
      class="sys-config-subnav"
    >
      <NTab
        v-for="opt in engineOptions"
        :key="opt.value"
        :name="opt.value"
        :tab="opt.label"
      />
    </NTabs>

    <ConfigSectionLayout
      description="配置各文件存储引擎参数。上方单选切换默认引擎（互斥）；保存后热重载生效。RustFS 为 S3 兼容存储，默认 path-style，Region 建议 us-east-1。"
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NForm
        class="sys-config-form"
        label-placement="left"
        label-width="140"
        require-mark-placement="left"
      >
        <template v-if="state.subTab === 'LOCAL'">
          <NFormItem
            label="WINDOWS存储位置"
            required
          >
            <NInput
              v-model:value="state.local.windowsRoot"
              placeholder="D:/defaultUploadFolder"
            />
          </NFormItem>
          <NFormItem
            label="LINUX存储位置"
            required
          >
            <NInput
              v-model:value="state.local.localRoot"
              placeholder="/defaultUploadFolder"
            />
          </NFormItem>
        </template>

        <template v-else-if="state.subTab === 'ALIYUN'">
          <NFormItem
            label="阿里云密钥ID"
            required
          >
            <NInput
              v-model:value="state.aliyun.accessKey"
              :placeholder="
                state.aliyun.accessKeySet ? '已配置，留空不修改' : '阿里云文件 AccessKeyId'
              "
            />
          </NFormItem>
          <NFormItem
            label="阿里云密钥SECRET"
            required
          >
            <NInput
              v-model:value="state.aliyun.secretKey"
              type="password"
              show-password-on="click"
              :placeholder="
                state.aliyun.secretKeySet ? '已配置，留空不修改' : '阿里云文件 AccessKeySecret'
              "
            />
          </NFormItem>
          <NFormItem
            label="阿里云文件端点"
            required
          >
            <NInput
              v-model:value="state.aliyun.endpoint"
              placeholder="oss-cn-hangzhou.aliyuncs.com"
            />
          </NFormItem>
          <NFormItem
            label="阿里云默认储存桶"
            required
          >
            <NInput
              v-model:value="state.aliyun.bucket"
              placeholder="defaultbucket"
            />
          </NFormItem>
        </template>

        <template v-else-if="state.subTab === 'TENCENT'">
          <NFormItem
            label="腾讯云密钥ID"
            required
          >
            <NInput
              v-model:value="state.tencent.accessKey"
              :placeholder="
                state.tencent.accessKeySet ? '已配置，留空不修改' : '腾讯云文件 SecretId'
              "
            />
          </NFormItem>
          <NFormItem
            label="腾讯云密钥SECRET"
            required
          >
            <NInput
              v-model:value="state.tencent.secretKey"
              type="password"
              show-password-on="click"
              :placeholder="
                state.tencent.secretKeySet ? '已配置，留空不修改' : '腾讯云文件 SecretKey'
              "
            />
          </NFormItem>
          <NFormItem
            label="腾讯云区域ID"
            required
          >
            <NInput
              v-model:value="state.tencent.region"
              placeholder="ap-beijing"
            />
          </NFormItem>
          <NFormItem
            label="腾讯云存储桶"
            required
          >
            <NInput
              v-model:value="state.tencent.bucket"
              placeholder="defaultbucket"
            />
          </NFormItem>
        </template>

        <template v-else-if="state.subTab === 'MINIO'">
          <NFormItem
            label="MINIO通道KEY"
            required
          >
            <NInput
              v-model:value="state.minio.accessKey"
              :placeholder="state.minio.accessKeySet ? '已配置，留空不修改' : 'MINIO Access Key'"
            />
          </NFormItem>
          <NFormItem
            label="MINIO密钥KEY"
            required
          >
            <NInput
              v-model:value="state.minio.secretKey"
              type="password"
              show-password-on="click"
              :placeholder="state.minio.secretKeySet ? '已配置，留空不修改' : 'MINIO Secret Key'"
            />
          </NFormItem>
          <NFormItem
            label="MINIO端点"
            required
          >
            <NInput
              v-model:value="state.minio.endpoint"
              placeholder="https://play.min.io"
            />
          </NFormItem>
          <NFormItem
            label="MINIO存储桶"
            required
          >
            <NInput
              v-model:value="state.minio.bucket"
              placeholder="defaultbucket"
            />
          </NFormItem>
        </template>

        <template v-else-if="state.subTab === 'RUSTFS'">
          <NFormItem
            label="RustFS Access Key"
            required
          >
            <NInput
              v-model:value="state.rustfs.accessKey"
              :placeholder="state.rustfs.accessKeySet ? '已配置，留空不修改' : 'admin'"
            />
          </NFormItem>
          <NFormItem
            label="RustFS Secret Key"
            required
          >
            <NInput
              v-model:value="state.rustfs.secretKey"
              type="password"
              show-password-on="click"
              :placeholder="state.rustfs.secretKeySet ? '已配置，留空不修改' : '123456789'"
            />
          </NFormItem>
          <NFormItem
            label="RustFS 端点"
            required
          >
            <NInput
              v-model:value="state.rustfs.endpoint"
              placeholder="http://127.0.0.1:9002"
            />
          </NFormItem>
          <NFormItem
            label="RustFS Region"
            required
          >
            <NInput
              v-model:value="state.rustfs.region"
              placeholder="us-east-1"
            />
          </NFormItem>
          <NFormItem
            label="RustFS 存储桶"
            required
          >
            <NInput
              v-model:value="state.rustfs.bucket"
              placeholder="defaultbucket"
            />
          </NFormItem>
          <NFormItem label="使用 SSL">
            <NSwitch v-model:value="state.rustfs.useSsl" />
          </NFormItem>
          <NFormItem label="自定义基础 URL">
            <NInput
              v-model:value="state.rustfs.baseUrl"
              placeholder="可选，公网访问前缀；留空则用预签名 URL"
            />
          </NFormItem>
        </template>
      </NForm>
    </ConfigSectionLayout>
  </NSpin>
</template>

<style scoped>
.mb-12px {
  margin-bottom: 12px;
}
</style>
