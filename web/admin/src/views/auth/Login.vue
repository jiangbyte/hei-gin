<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormInst, FormItemRule, FormRules } from 'naive-ui'
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import CaptchaInput from '@/components/common/CaptchaInput.vue'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores'
import { isValidEmail } from '@/utils'
import { encryptPasswords } from '@/utils/security'
import { wireBool } from '@/utils/wire'
import AuthLayout from './AuthLayout.vue'

const OTP_COOLDOWN_SECONDS = 60

type LoginType = 'ACCOUNT' | 'EMAIL' | 'PHONE'

const route = useRoute()
const authStore = useAuthStore()
const formRef = ref<FormInst | null>(null)
const captchaRef = ref<InstanceType<typeof CaptchaInput> | null>(null)
const loading = ref(false)
const sendingCode = ref(false)
const otpCooldown = ref(0)
let otpCooldownTimer: ReturnType<typeof setInterval> | null = null
const activeType = ref<LoginType>('ACCOUNT')
const loginMode = ref<'PASSWORD' | 'OTP'>('PASSWORD')
const copyright = ref(import.meta.env.VITE_COPYRIGHT_INFO || '')
const copyrightUrl = ref('')

const options = reactive({
  allow_account: true,
  allow_email: true,
  allow_phone: true,
  allow_otp: true,
})

const sendCodeLabel = computed(() =>
  otpCooldown.value > 0 ? `${otpCooldown.value}s 后重发` : '发送验证码',
)

function startOtpCooldown() {
  otpCooldown.value = OTP_COOLDOWN_SECONDS
  if (otpCooldownTimer) clearInterval(otpCooldownTimer)
  otpCooldownTimer = setInterval(() => {
    otpCooldown.value -= 1
    if (otpCooldown.value <= 0 && otpCooldownTimer) {
      clearInterval(otpCooldownTimer)
      otpCooldownTimer = null
    }
  }, 1000)
}

const allLoginTypes: Array<{ key: LoginType; label: string; placeholder: string }> = [
  { key: 'ACCOUNT', label: '账号', placeholder: '请输入管理员账号' },
  { key: 'EMAIL', label: '邮箱', placeholder: '请输入登录邮箱' },
  { key: 'PHONE', label: '手机号', placeholder: '请输入登录手机号' },
]

const loginTypes = computed(() =>
  allLoginTypes.filter((item) => {
    if (item.key === 'ACCOUNT') return options.allow_account
    if (item.key === 'EMAIL') return options.allow_email
    return options.allow_phone
  }),
)

const form = reactive({
  account: '',
  email: '',
  phone: '',
  password: '',
  otp_code: '',
  captcha_id: '',
  captcha_value: '',
  remember: true,
})

const currentLogin = computed(
  () => loginTypes.value.find((item) => item.key === activeType.value) || loginTypes.value[0],
)
const activeField = computed(() => activeType.value.toLowerCase() as 'account' | 'email' | 'phone')
const otpAvailable = computed(
  () => options.allow_otp && (activeType.value === 'EMAIL' || activeType.value === 'PHONE'),
)

onMounted(async () => {
  try {
    const res = await authApi.authOptions()
    const data = res?.data || {}
    options.allow_account = wireBool(data.allow_account ?? true)
    options.allow_email = wireBool(data.allow_email ?? true)
    options.allow_phone = wireBool(data.allow_phone ?? true)
    options.allow_otp = wireBool(data.allow_otp ?? true)
    if (data.copyright_text) copyright.value = data.copyright_text
    copyrightUrl.value = data.copyright_url || ''
    if (!loginTypes.value.some((item) => item.key === activeType.value)) {
      activeType.value = loginTypes.value[0]?.key || 'ACCOUNT'
    }
  } catch {
    // 使用默认全开
  }
})

onUnmounted(() => {
  if (otpCooldownTimer) clearInterval(otpCooldownTimer)
})

watch(activeType, () => {
  if (!otpAvailable.value) loginMode.value = 'PASSWORD'
})

function validateLoginIdentity(_rule: FormItemRule, value: string) {
  const text = String(value ?? '').trim()
  if (!text) {
    return new Error(`请输入${currentLogin.value?.label || '账号'}`)
  }
  if (activeType.value === 'EMAIL' && !isValidEmail(text)) {
    return new Error('请输入有效邮箱')
  }
  return true
}

const rules = computed<FormRules>(() => {
  const next: FormRules = {
    [activeField.value]: [
      {
        validator: validateLoginIdentity,
        trigger: ['input', 'blur'],
      },
    ],
    captcha_value: [
      {
        required: true,
        message: '请输入验证码',
        trigger: ['input', 'blur'],
      },
    ],
  }
  if (loginMode.value === 'OTP') {
    next.otp_code = [{ required: true, message: '请输入登录验证码', trigger: ['input', 'blur'] }]
  } else {
    next.password = [{ required: true, message: '请输入密码', trigger: ['input', 'blur'] }]
  }
  return next
})

async function handleSendCode() {
  if (otpCooldown.value > 0 || sendingCode.value) return
  const target = form[activeField.value].trim()
  if (!target) {
    window.$message.warning(`请输入${currentLogin.value?.label}`)
    return
  }
  if (activeType.value === 'EMAIL' && !isValidEmail(target)) {
    window.$message.warning('请输入有效邮箱')
    return
  }
  if (!form.captcha_value.trim()) {
    window.$message.warning('请输入图形验证码')
    return
  }
  sendingCode.value = true
  try {
    await authApi.sendLoginCode({
      target,
      channel: activeType.value === 'EMAIL' ? 'EMAIL' : 'PHONE',
      captcha_id: form.captcha_id,
      captcha_value: form.captcha_value,
    })
    window.$message.success('验证码已发送，请查收后填写')
    startOtpCooldown()
    // 发送成功会消耗图形验证码，需刷新供登录提交使用
    await captchaRef.value?.refresh()
  } catch {
    await captchaRef.value?.refresh()
  } finally {
    sendingCode.value = false
  }
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : undefined
    const security: Record<string, string> = {
      captcha_id: form.captcha_id,
      captcha_value: form.captcha_value,
    }
    let password = ''
    if (loginMode.value === 'PASSWORD') {
      const encrypted = await encryptPasswords({ password: form.password })
      password = encrypted.values.password || ''
      security.password_key_id = encrypted.password_key_id
    }
    await authStore.login(
      form[activeField.value].trim(),
      password,
      redirect,
      form.remember,
      activeType.value,
      {
        ...security,
        login_mode: loginMode.value,
        ...(loginMode.value === 'OTP' && form.otp_code.trim()
          ? { otp_code: form.otp_code.trim() }
          : {}),
      } as any,
    )
    window.$message.success('登录成功')
  } catch {
    await captchaRef.value?.refresh()
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthLayout
    title="管理端登录"
    :copyright="copyright"
    :copyright-url="copyrightUrl"
  >
    <n-form
      ref="formRef"
      :model="form"
      :rules="rules"
      size="large"
      @submit.prevent="handleSubmit"
    >
      <n-tabs
        v-model:value="activeType"
        type="segment"
        size="large"
        class="auth-login-tabs"
      >
        <n-tab-pane
          v-for="item in loginTypes"
          :key="item.key"
          :name="item.key"
          :tab="item.label"
        />
      </n-tabs>

      <n-form-item :path="activeField">
        <n-input
          v-model:value="form[activeField]"
          size="large"
          :placeholder="currentLogin?.placeholder"
          clearable
        />
      </n-form-item>

      <div
        v-if="otpAvailable"
        class="auth-mode-row"
      >
        <n-radio-group
          v-model:value="loginMode"
          size="large"
        >
          <n-radio-button
            value="PASSWORD"
            label="密码登录"
          />
          <n-radio-button
            value="OTP"
            label="验证码登录"
          />
        </n-radio-group>
      </div>

      <n-form-item
        v-if="loginMode === 'PASSWORD'"
        path="password"
      >
        <n-input
          v-model:value="form.password"
          size="large"
          type="password"
          show-password-on="click"
          placeholder="请输入密码"
        />
      </n-form-item>
      <n-form-item
        v-else
        path="otp_code"
      >
        <div class="auth-otp-row">
          <n-input
            v-model:value="form.otp_code"
            size="large"
            placeholder="请输入登录验证码"
          />
          <n-button
            size="large"
            :loading="sendingCode"
            :disabled="otpCooldown > 0"
            @click="handleSendCode"
          >
            {{ sendCodeLabel }}
          </n-button>
        </div>
      </n-form-item>
      <n-form-item path="captcha_value">
        <CaptchaInput
          ref="captchaRef"
          v-model:captcha-id="form.captcha_id"
          v-model:captcha-value="form.captcha_value"
          size="large"
        />
      </n-form-item>
      <div class="auth-form-row">
        <n-checkbox
          v-model:checked="form.remember"
          size="large"
        >
          记住我
        </n-checkbox>
        <RouterLink to="/auth/forgot-password">
          忘记密码？
        </RouterLink>
      </div>
      <n-button
        class="auth-submit"
        type="primary"
        size="large"
        block
        attr-type="submit"
        :loading="loading"
      >
        登录
      </n-button>
    </n-form>
  </AuthLayout>
</template>

<style scoped>
.auth-login-tabs :deep(.n-tabs-pane-wrapper) {
  overflow: visible;
}

.auth-mode-row {
  margin: 0 0 12px;
}

.auth-otp-row {
  display: flex;
  gap: 8px;
  width: 100%;
}

.auth-form-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin: 0 0 16px;
  font-size: 14px;
}

.auth-form-row a {
  color: var(--n-primary-color, #1677ff);
  text-decoration: none;
}

@media (max-width: 420px) {
  .auth-form-row {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
