<!-- Author: Charlie -->

<script setup lang="ts">
import type { FormItemInst, FormItemRule } from 'naive-ui'
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
const identityItemRef = ref<FormItemInst | null>(null)
const passwordItemRef = ref<FormItemInst | null>(null)
const otpItemRef = ref<FormItemInst | null>(null)
const captchaItemRef = ref<FormItemInst | null>(null)
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

type OauthProviderOption = {
  provider: string
  label: string
  enabled: boolean
  web_oauth: boolean
}

const oauthProviders = ref<OauthProviderOption[]>([])
const oauthLoading = ref<string | null>(null)

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

const identityRule: FormItemRule = {
  trigger: ['input', 'blur'],
  validator() {
    const text = form[activeField.value].trim()
    if (!text) {
      return new Error(`请输入${currentLogin.value?.label || '账号'}`)
    }
    if (activeType.value === 'EMAIL' && !isValidEmail(text)) {
      return new Error('请输入有效邮箱')
    }
    return true
  },
}

const passwordRule: FormItemRule = {
  trigger: ['input', 'blur'],
  validator() {
    if (!form.password) return new Error('请输入密码')
    return true
  },
}

const otpRule: FormItemRule = {
  trigger: ['input', 'blur'],
  validator() {
    if (!form.otp_code.trim()) return new Error('请输入登录验证码')
    return true
  },
}

const captchaRule: FormItemRule = {
  trigger: ['input', 'blur'],
  validator() {
    if (!form.captcha_value.trim()) return new Error('请输入验证码')
    return true
  },
}

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
    const providers = Array.isArray(data.oauth_providers) ? data.oauth_providers : []
    oauthProviders.value = providers
      .map((item: any) => ({
        provider: String(item.provider || ''),
        label: String(item.label || item.provider || ''),
        enabled: wireBool(item.enabled ?? false),
        web_oauth: wireBool(item.web_oauth ?? true),
      }))
      .filter((item: OauthProviderOption) => item.provider && item.enabled && item.web_oauth)
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
  identityItemRef.value?.restoreValidation()
})

watch(loginMode, () => {
  passwordItemRef.value?.restoreValidation()
  otpItemRef.value?.restoreValidation()
})

async function validateFields(items: Array<FormItemInst | null | undefined>) {
  try {
    await Promise.all(items.filter(Boolean).map((item) => item!.validate()))
    return true
  } catch {
    return false
  }
}

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
    await captchaRef.value?.refresh()
  } catch {
    await captchaRef.value?.refresh()
  } finally {
    sendingCode.value = false
  }
}

async function handleOauthLogin(provider: string) {
  if (oauthLoading.value) return
  oauthLoading.value = provider
  try {
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : undefined
    const res = await authApi.oauthAuthorize(provider, {
      intent: 'LOGIN',
      redirect,
    })
    const url = res?.data?.authorize_url
    if (!url) {
      window.$message.error('无法发起三方登录')
      return
    }
    window.location.href = String(url)
  } catch {
    // 全局错误提示；管理端未绑定时会在回调页提示
  } finally {
    oauthLoading.value = null
  }
}

function providerInitial(provider: string, label: string) {
  const key = provider.toUpperCase()
  if (key === 'GITHUB') return 'GH'
  if (key === 'GITEE') return 'GE'
  if (key === 'QQ') return 'QQ'
  if (key.startsWith('WECHAT')) return '微'
  return (label || provider).slice(0, 1).toUpperCase()
}

async function handleSubmit() {
  const credentialItem = loginMode.value === 'PASSWORD' ? passwordItemRef.value : otpItemRef.value
  const ok = await validateFields([identityItemRef.value, credentialItem, captchaItemRef.value])
  if (!ok) return

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
    <form
      class="auth-form"
      @submit.prevent="handleSubmit"
    >
      <n-tabs
        v-model:value="activeType"
        type="line"
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

      <n-form-item
        ref="identityItemRef"
        :show-label="false"
        :rule="identityRule"
      >
        <n-input
          v-model:value="form[activeField]"
          size="large"
          :placeholder="currentLogin?.placeholder"
          clearable
        />
      </n-form-item>

      <n-form-item
        v-if="loginMode === 'PASSWORD'"
        ref="passwordItemRef"
        :show-label="false"
        :rule="passwordRule"
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
        ref="otpItemRef"
        :show-label="false"
        :rule="otpRule"
      >
        <div class="auth-otp-row">
          <n-input
            v-model:value="form.otp_code"
            class="auth-otp-row__input"
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

      <n-form-item
        ref="captchaItemRef"
        :show-label="false"
        :rule="captchaRule"
      >
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
        <div class="auth-form-row__links">
          <button
            v-if="otpAvailable"
            type="button"
            class="auth-mode-link"
            @click="loginMode = loginMode === 'PASSWORD' ? 'OTP' : 'PASSWORD'"
          >
            {{ loginMode === 'PASSWORD' ? '验证码登录' : '密码登录' }}
          </button>
          <span
            v-if="otpAvailable"
            class="auth-form-row__sep"
          >·</span>
          <RouterLink to="/auth/forgot-password">
            忘记密码？
          </RouterLink>
        </div>
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

      <div
        v-if="oauthProviders.length"
        class="auth-oauth"
      >
        <div class="auth-oauth__divider">
          <span>其他登录方式</span>
        </div>
        <div class="auth-oauth__row">
          <button
            v-for="item in oauthProviders"
            :key="item.provider"
            type="button"
            class="auth-oauth__icon"
            :title="item.label"
            :aria-label="item.label"
            :disabled="Boolean(oauthLoading)"
            @click="handleOauthLogin(item.provider)"
          >
            {{ providerInitial(item.provider, item.label) }}
          </button>
        </div>
        <p class="auth-oauth__hint">
          管理端需先用密码登录并在用户中心绑定三方账号
        </p>
      </div>

      <p class="auth-form-foot">
        登录遇到问题？请联系系统管理员
      </p>
    </form>
  </AuthLayout>
</template>
