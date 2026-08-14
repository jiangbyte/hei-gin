<!-- Author: Charlie -->

<script setup lang="ts">
import { authApi } from '@/api'
import { computed, onMounted, ref } from 'vue'

withDefaults(
  defineProps<{
    size?: 'tiny' | 'small' | 'medium' | 'large'
  }>(),
  {
    size: 'large',
  },
)

const captchaId = defineModel<string>('captchaId', { required: true })
const captchaValue = defineModel<string>('captchaValue', { required: true })
const loading = ref(false)
const imageBase64 = ref('')

const imageSrc = computed(() =>
  imageBase64.value ? `data:image/svg+xml;base64,${imageBase64.value}` : '',
)

async function refresh() {
  loading.value = true
  try {
    const response = await authApi.captcha()
    captchaId.value = response.data.captcha_id
    captchaValue.value = ''
    imageBase64.value = response.data.image_base64
  } finally {
    loading.value = false
  }
}

onMounted(refresh)

defineExpose({ refresh })
</script>

<template>
  <div
    class="captcha-input"
    :class="`captcha-input--${size}`"
  >
    <NInput
      v-model:value="captchaValue"
      class="captcha-input__field"
      :size="size"
      :placeholder="'请输入验证码'"
      clearable
    />
    <button
      class="captcha-image"
      type="button"
      :disabled="loading"
      aria-label="刷新验证码"
      @click="refresh"
    >
      <NSpin
        :show="loading"
        size="small"
      >
        <img
          v-if="imageSrc"
          :src="imageSrc"
          alt="验证码"
        >
      </NSpin>
    </button>
  </div>
</template>

<style scoped>
/* 以输入框真实高度为准，图片绝对定位铺满，避免 SVG 撑高 */
.captcha-input {
  display: flex;
  align-items: stretch;
  gap: 10px;
  width: 100%;
}

.captcha-input__field {
  flex: 1;
  min-width: 0;
}

.captcha-image {
  position: relative;
  box-sizing: border-box;
  flex: 0 0 140px;
  width: 140px;
  padding: 0;
  overflow: hidden;
  cursor: pointer;
  line-height: 0;
  background: var(--n-color, #f8fafc);
  border: 1px solid var(--n-border-color);
  border-radius: var(--border-radius, var(--n-border-radius, 3px));
}

.captcha-image:disabled {
  cursor: wait;
}

.captcha-image :deep(.n-spin-nested-loading),
.captcha-image :deep(.n-spin-container) {
  position: absolute;
  inset: 0;
}

.captcha-image img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
