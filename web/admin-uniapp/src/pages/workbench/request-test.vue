<!-- Author: Charlie -->

<template>
  <Layout title="请求测试" :back="true">
    <view class="px-4 py-4">
      <text class="text-base font-semibold text-gray-800 mb-3 block"
      >HTTP 请求测试
      </text
      >

      <!-- GET 请求测试 -->
      <view class="bg-white rounded p-3 mb-3">
        <view class="flex items-center justify-between mb-2">
          <text class="text-sm font-semibold text-gray-700"
          >GET {{ API_PREFIX }}/me
          </text
          >
          <u-button
              text="发送"
              size="mini"
              type="primary"
              :loading="loadingMe"
              @click="fetchMe"
          />
        </view>
        <view
            v-if="meResult"
            class="bg-gray-50 rounded p-2 text-xs font-mono text-gray-600 leading-relaxed max-h-32 overflow-y-auto"
        >
          <text>{{ meResult }}</text>
        </view>
        <view v-else class="bg-gray-50 rounded p-2 text-xs text-gray-400"
        >点击发送查看结果
        </view
        >
      </view>

      <!-- GET 请求测试 - 字典树 -->
      <view class="bg-white rounded p-3 mb-3">
        <view class="flex items-center justify-between mb-2">
          <text class="text-sm font-semibold text-gray-700"
          >GET {{ API_PREFIX }}/sys/dicts/tree
          </text
          >
          <u-button
              text="发送"
              size="mini"
              type="primary"
              :loading="loadingDictTree"
              @click="fetchDictTree"
          />
        </view>
        <view
            v-if="dictTreeResult"
            class="bg-gray-50 rounded p-2 text-xs font-mono text-gray-600 leading-relaxed max-h-40 overflow-y-auto"
        >
          <text>{{ dictTreeResult }}</text>
        </view>
        <view v-else class="bg-gray-50 rounded p-2 text-xs text-gray-400"
        >点击发送查看结果
        </view
        >
      </view>

      <!-- 错误请求测试 -->
      <view class="bg-white rounded p-3 mb-3">
        <view class="flex items-center justify-between mb-2">
          <text class="text-sm font-semibold text-red-600">错误请求测试</text>
          <u-button
              text="发送"
              size="mini"
              type="error"
              :loading="loadingError"
              @click="fetchError"
          />
        </view>
        <view
            v-if="errorResult"
            class="bg-red-50 rounded p-2 text-xs font-mono text-red-600 leading-relaxed max-h-32 overflow-y-auto"
        >
          <text>{{ errorResult }}</text>
        </view>
        <view v-else class="bg-gray-50 rounded p-2 text-xs text-gray-400"
        >点击发送测试错误请求响应
        </view
        >
      </view>

      <!-- 文件上传测试 -->
      <view class="bg-white rounded p-3 mb-3">
        <text class="text-sm font-semibold text-gray-700 mb-2 block"
        >文件上传测试 — POST {{ API_PREFIX }}/sys/file/upload
        </text
        >
        <view class="flex flex-col items-center gap-3 mb-3">
          <view
              class="w-full bg-gray-50 rounded p-3 flex flex-col items-center gap-2"
          >
            <u-icon name="file-text" size="40" color="#3c96f3"/>
            <text
                v-if="uploadFileInfo"
                class="text-xs text-gray-600 text-center"
            >{{ uploadFileInfo }}
            </text
            >
            <text v-else class="text-xs text-gray-400">尚未选择文件</text>
          </view>
          <view
              class="inline-block bg-blue-500 text-white text-xs px-4 py-1.5 rounded-full"
              @click="chooseUploadFile"
          >选择文件
          </view
          >
        </view>
        <view v-if="uploading" class="mb-2">
          <view class="bg-blue-100 rounded h-2 w-full relative">
            <view
                class="bg-blue-500 rounded h-2 absolute left-0 top-0"
                :style="{ width: uploadProgress + '%' }"
            />
          </view>
          <text class="text-xs text-gray-500 mt-1 block"
          >{{ uploadProgress }}%
          </text
          >
        </view>
        <u-button
            text="上传文件"
            type="primary"
            :disabled="!uploadFilePath"
            :loading="uploading"
            @click="doUploadFile"
        />
        <view
            v-if="uploadResult"
            class="bg-green-50 rounded p-2 mt-2 text-xs font-mono text-green-700 leading-relaxed max-h-48 overflow-y-auto"
        >
          <text>{{ uploadResult }}</text>
        </view>
        <view
            v-else-if="uploadError"
            class="bg-red-50 rounded p-2 mt-2 text-xs font-mono text-red-600 leading-relaxed max-h-32 overflow-y-auto"
        >
          <text>{{ uploadError }}</text>
        </view>
      </view>
    </view>
  </Layout>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import {API_PREFIX} from '@/constants/api'
import Layout from '@/layouts/index.vue'

const loadingMe = ref(false)
const loadingDictTree = ref(false)

const loadingError = ref(false)
const meResult = ref('')
const dictTreeResult = ref('')
const errorResult = ref('')

async function fetchMe() {
  loadingMe.value = true
  meResult.value = ''
  try {
    const {http} = await import('@/utils/request')
    const data = await http.get(`${API_PREFIX}/me`)
    meResult.value = JSON.stringify(data, null, 2)
  } catch (e: any) {
    meResult.value = `错误：${e.message || e}`
  } finally {
    loadingMe.value = false
  }
}

async function fetchDictTree() {
  loadingDictTree.value = true
  dictTreeResult.value = ''
  try {
    const {http} = await import('@/utils/request')
    const data = await http.get(`${API_PREFIX}/sys/dicts/tree`)
    const preview = (Array.isArray(data) ? data : [])
        .slice(0, 3)
        .map((n: any) => ({
          code: n.code,
          label: n.label,
          children: (n.children || []).map((c: any) => ({
            label: c.label,
            value: c.value,
            color: c.color,
          })),
        }))
    dictTreeResult.value = JSON.stringify(preview, null, 2)
  } catch (e: any) {
    dictTreeResult.value = `错误：${e.message || e}`
  } finally {
    loadingDictTree.value = false
  }
}

async function fetchError() {
  loadingError.value = true
  errorResult.value = ''
  try {
    const {http} = await import('@/utils/request')
    const data = await http.get(`${API_PREFIX}/non-existent-path`)
    errorResult.value = JSON.stringify(data, null, 2)
  } catch (e: any) {
    errorResult.value = `错误：${e.message || e}\ncode: ${e.code}\nstatusCode: ${e.statusCode}`
  } finally {
    loadingError.value = false
  }
}

const uploadFilePath = ref('')
const uploadFileInfo = ref('')
const uploadProgress = ref(0)
const uploading = ref(false)
const uploadResult = ref('')
const uploadError = ref('')

function chooseUploadFile() {
  uploadResult.value = ''
  uploadError.value = ''
  uni.chooseFile({
    count: 1,
    type: 'all',
    success: (res) => {
      uploadFilePath.value = res.tempFilePaths[0]
      const file = res.tempFiles?.[0]
      const name =
          file?.name ||
          uploadFilePath.value.split('/').pop() ||
          uploadFilePath.value
      const size = file?.size ? formatFileSize(file.size) : ''
      uploadFileInfo.value = [name, size].filter(Boolean).join(' ｜ ')
    },
  })
}

function formatFileSize(bytes: number) {
  if (bytes < 1024) return bytes + 'B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + 'KB'
  return (bytes / 1024 / 1024).toFixed(1) + 'MB'
}

async function doUploadFile() {
  if (!uploadFilePath.value) return
  uploading.value = true
  uploadProgress.value = 0
  uploadResult.value = ''
  uploadError.value = ''

  try {
    const {http} = await import('@/utils/request')
    uni.showLoading({title: '上传中...'})
    const result = await http.upload(
        `${API_PREFIX}/sys/file/upload`,
        uploadFilePath.value
    )
    uni.hideLoading()
    uploadProgress.value = 100
    uploadResult.value = JSON.stringify(result, null, 2)
  } catch (e: any) {
    uni.hideLoading()
    uploadError.value = `错误：${e.message || e}`
  } finally {
    uploading.value = false
  }
}
</script>

<style lang="scss" scoped>
page {
  background-color: #f5f6f7;
}
</style>
