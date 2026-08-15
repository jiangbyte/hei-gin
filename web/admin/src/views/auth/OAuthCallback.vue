<!-- Author: Charlie -->

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { oauthExchange } from '@/api/auth'
import { useAuthStore } from '@/stores'
import { clearToken, setToken } from '@/utils/session'
import { wireBool } from '@/utils/wire'
import { refreshDict, syncDictTree } from '@/utils/dict'
import AuthLayout from './AuthLayout.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const error = ref('')
const tip = ref('正在完成登录…')

onMounted(() => {
  void handleCallback()
})

async function handleCallback() {
  const status = String(route.query.oauth_status || '')
  const action = String(route.query.oauth_action || '')
  const rawMessage = typeof route.query.oauth_message === 'string' ? route.query.oauth_message : ''
  const oauthCode = typeof route.query.oauth_code === 'string' ? route.query.oauth_code : ''
  // 兼容旧版 URL token（过渡期）
  const legacyToken = typeof route.query.token === 'string' ? route.query.token : ''
  const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : undefined

  if (status !== 'ok') {
    let msg = '三方登录失败'
    if (rawMessage) {
      try {
        msg = decodeURIComponent(rawMessage)
      } catch {
        msg = rawMessage
      }
    }
    error.value = msg
    tip.value = msg
    window.$message?.error?.(msg)
    window.setTimeout(() => {
      void router.replace({
        path: '/auth/login',
        query: redirect ? { redirect } : undefined,
      })
    }, 1600)
    return
  }

  try {
    let passwordExpired = false
    let forceBindEmail = false
    let forceBindPhone = false

    if (oauthCode) {
      const { data } = await oauthExchange({ code: oauthCode })
      const token = data?.token
      if (!token) {
        throw new Error('兑换登录凭证失败')
      }
      clearToken()
      setToken(String(token), true)
      authStore.sessionChecked = true
      passwordExpired = wireBool(data?.password_expired ?? false)
      forceBindEmail = wireBool(data?.force_bind_email ?? false)
      forceBindPhone = wireBool(data?.force_bind_phone ?? false)
    } else if (legacyToken) {
      clearToken()
      setToken(legacyToken, true)
      authStore.sessionChecked = true
      passwordExpired = wireBool(String(route.query.password_expired ?? false))
      forceBindEmail = wireBool(String(route.query.force_bind_email ?? false))
      forceBindPhone = wireBool(String(route.query.force_bind_phone ?? false))
    }

    if (action === 'bound') {
      tip.value = '绑定成功，正在跳转…'
      await authStore.refreshUserInfo()
      window.$message?.success?.('绑定成功')
      await router.replace('/profile?tab=oauth')
      return
    }

    if (passwordExpired) {
      window.$message?.warning?.('密码已过期，请先修改密码')
    } else if (forceBindEmail || forceBindPhone) {
      window.$message?.warning?.('请先完成账号安全绑定')
    } else {
      window.$message?.success?.('登录成功')
    }

    tip.value = '登录成功，正在进入系统…'
    await authStore.finishLogin(redirect)
    syncDictTree()
    await refreshDict()
  } catch (e: any) {
    const msg = e?.message || '登录会话建立失败'
    error.value = msg
    tip.value = msg
    window.$message?.error?.(msg)
    window.setTimeout(() => {
      void router.replace('/auth/login')
    }, 1600)
  }
}
</script>

<template>
  <AuthLayout
    variant="center"
    title="三方登录"
    copyright=""
  >
    <div class="oauth-callback">
      <NSpin
        v-if="!error"
        :show="true"
      >
        <div class="oauth-callback__tip">
          {{ tip }}
        </div>
      </NSpin>
      <NResult
        v-else
        status="error"
        title="三方登录失败"
        :description="error"
      />
    </div>
  </AuthLayout>
</template>
