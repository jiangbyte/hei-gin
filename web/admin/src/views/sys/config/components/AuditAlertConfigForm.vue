<!-- Author: Charlie -->

<script setup lang="ts">
import { configApi } from '@/api'
import { computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ConfigSectionLayout from './ConfigSectionLayout.vue'
import {
  loadByCategory,
  parseBool,
  parseNumber,
  saveByKeys,
  toBoolStr,
} from '../composables/useConfigForm'

const CATEGORY = 'AUDIT_ALERT'
const route = useRoute()
const router = useRouter()

const pushEngineLabel: Record<string, string> = {
  DINGTALK: '钉钉',
  LARK: '飞书',
  WECHAT_WORK: '企业微信',
}

const state = reactive({
  loading: false,
  saving: false,
  testingWebhook: false,
  testingPush: false,
  enabled: false,
  notifyEmail: true,
  notifyEmailTo: '',
  notifyPush: true,
  notifyCustomWebhook: false,
  webhookUrl: '',
  webhookSecret: '',
  analysisInterval: 60,
  cooldown: 1800,
  ruleBruteForce: true,
  bruteForceThreshold: 10,
  ruleUnusualHours: true,
  ruleSensitiveOps: true,
  ruleBulkDelete: true,
  ruleIpAnomaly: true,
  bulkDeleteThreshold: 20,
  ipAnomalyThreshold: 3,
  pushEngine: 'DINGTALK',
  pushConfigured: false,
  mailConfigured: false,
  snapshot: '',
})

const pushHint = computed(() => {
  const name = pushEngineLabel[state.pushEngine] || state.pushEngine
  return state.pushConfigured
    ? `复用「消息推送」中的默认引擎（当前：${name}）`
    : `复用「消息推送」配置（当前默认：${name}，尚未配置 Webhook）`
})

onMounted(() => {
  void reload()
})

function hasKey(map: Record<string, string>, key: string) {
  return Object.prototype.hasOwnProperty.call(map, key)
}

async function reload() {
  state.loading = true
  try {
    const [map, pushMap, mailMap, sysMap] = await Promise.all([
      loadByCategory(CATEGORY),
      loadByCategory('PUSH'),
      loadByCategory('MAIL'),
      loadByCategory('SYS'),
    ])
    state.enabled = parseBool(map.AUDIT_ALERT_ENABLED)
    state.notifyEmail = hasKey(map, 'AUDIT_ALERT_NOTIFY_EMAIL')
      ? parseBool(map.AUDIT_ALERT_NOTIFY_EMAIL)
      : true
    state.notifyEmailTo = map.AUDIT_ALERT_NOTIFY_EMAIL_TO || ''
    state.notifyPush = hasKey(map, 'AUDIT_ALERT_NOTIFY_PUSH')
      ? parseBool(map.AUDIT_ALERT_NOTIFY_PUSH)
      : true
    state.notifyCustomWebhook = hasKey(map, 'AUDIT_ALERT_NOTIFY_CUSTOM_WEBHOOK')
      ? parseBool(map.AUDIT_ALERT_NOTIFY_CUSTOM_WEBHOOK)
      : Boolean(map.AUDIT_ALERT_WEBHOOK_URL)
    state.webhookUrl = map.AUDIT_ALERT_WEBHOOK_URL || ''
    state.webhookSecret = map.AUDIT_ALERT_WEBHOOK_SECRET || ''
    state.analysisInterval = parseNumber(map.AUDIT_ALERT_ANALYSIS_INTERVAL_SECONDS, 60)
    state.cooldown = parseNumber(map.AUDIT_ALERT_ALERT_COOLDOWN_SECONDS, 1800)
    state.ruleBruteForce = hasKey(map, 'AUDIT_ALERT_RULE_BRUTE_FORCE')
      ? parseBool(map.AUDIT_ALERT_RULE_BRUTE_FORCE)
      : true
    state.bruteForceThreshold = parseNumber(map.AUDIT_ALERT_BRUTE_FORCE_THRESHOLD, 10)
    state.ruleUnusualHours = hasKey(map, 'AUDIT_ALERT_RULE_UNUSUAL_HOURS')
      ? parseBool(map.AUDIT_ALERT_RULE_UNUSUAL_HOURS)
      : true
    state.ruleSensitiveOps = hasKey(map, 'AUDIT_ALERT_RULE_SENSITIVE_OPS')
      ? parseBool(map.AUDIT_ALERT_RULE_SENSITIVE_OPS)
      : true
    state.ruleBulkDelete = hasKey(map, 'AUDIT_ALERT_RULE_BULK_DELETE')
      ? parseBool(map.AUDIT_ALERT_RULE_BULK_DELETE)
      : true
    state.ruleIpAnomaly = hasKey(map, 'AUDIT_ALERT_RULE_IP_ANOMALY')
      ? parseBool(map.AUDIT_ALERT_RULE_IP_ANOMALY)
      : true
    state.bulkDeleteThreshold = parseNumber(map.AUDIT_ALERT_BULK_DELETE_THRESHOLD, 20)
    state.ipAnomalyThreshold = parseNumber(map.AUDIT_ALERT_IP_ANOMALY_THRESHOLD, 3)

    state.pushEngine = (
      pushMap.DEFAULT_MESSAGE_PUSH_ENGINE ||
      sysMap.DEFAULT_MESSAGE_PUSH_ENGINE ||
      'DINGTALK'
    ).toUpperCase()
    state.pushConfigured = Boolean(
      pushMap.PUSH_DINGTALK_WEBHOOK ||
        pushMap.PUSH_LARK_WEBHOOK ||
        pushMap.PUSH_WECHAT_WORK_WEBHOOK,
    )
    state.mailConfigured = Boolean(mailMap.MAIL_LOCAL_HOST || mailMap.MAIL_LOCAL_FROM_EMAIL)

    state.snapshot = snapshotOf()
  } finally {
    state.loading = false
  }
}

function snapshotOf() {
  return JSON.stringify({
    enabled: state.enabled,
    notifyEmail: state.notifyEmail,
    notifyEmailTo: state.notifyEmailTo,
    notifyPush: state.notifyPush,
    notifyCustomWebhook: state.notifyCustomWebhook,
    webhookUrl: state.webhookUrl,
    webhookSecret: state.webhookSecret,
    analysisInterval: state.analysisInterval,
    cooldown: state.cooldown,
    ruleBruteForce: state.ruleBruteForce,
    bruteForceThreshold: state.bruteForceThreshold,
    ruleUnusualHours: state.ruleUnusualHours,
    ruleSensitiveOps: state.ruleSensitiveOps,
    ruleBulkDelete: state.ruleBulkDelete,
    ruleIpAnomaly: state.ruleIpAnomaly,
    bulkDeleteThreshold: state.bulkDeleteThreshold,
    ipAnomalyThreshold: state.ipAnomalyThreshold,
  })
}

function reset() {
  if (!state.snapshot) return
  Object.assign(state, JSON.parse(state.snapshot))
}

function goConfig(tab: string) {
  void router.replace({ query: { ...route.query, tab } })
}

async function save() {
  state.saving = true
  try {
    await saveByKeys([
      {
        config_key: 'AUDIT_ALERT_ENABLED',
        config_value: toBoolStr(state.enabled),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_NOTIFY_EMAIL',
        config_value: toBoolStr(state.notifyEmail),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_NOTIFY_EMAIL_TO',
        config_value: state.notifyEmailTo.trim(),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_NOTIFY_PUSH',
        config_value: toBoolStr(state.notifyPush),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_NOTIFY_CUSTOM_WEBHOOK',
        config_value: toBoolStr(state.notifyCustomWebhook),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_WEBHOOK_URL',
        config_value: state.webhookUrl,
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_WEBHOOK_SECRET',
        config_value: state.webhookSecret,
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_ANALYSIS_INTERVAL_SECONDS',
        config_value: String(state.analysisInterval),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_ALERT_COOLDOWN_SECONDS',
        config_value: String(state.cooldown),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_RULE_BRUTE_FORCE',
        config_value: toBoolStr(state.ruleBruteForce),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_BRUTE_FORCE_THRESHOLD',
        config_value: String(state.bruteForceThreshold),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_RULE_UNUSUAL_HOURS',
        config_value: toBoolStr(state.ruleUnusualHours),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_RULE_SENSITIVE_OPS',
        config_value: toBoolStr(state.ruleSensitiveOps),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_RULE_BULK_DELETE',
        config_value: toBoolStr(state.ruleBulkDelete),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_RULE_IP_ANOMALY',
        config_value: toBoolStr(state.ruleIpAnomaly),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_BULK_DELETE_THRESHOLD',
        config_value: String(state.bulkDeleteThreshold),
        category: CATEGORY,
      },
      {
        config_key: 'AUDIT_ALERT_IP_ANOMALY_THRESHOLD',
        config_value: String(state.ipAnomalyThreshold),
        category: CATEGORY,
      },
    ])
    window.$message.success('保存成功')
    state.snapshot = snapshotOf()
  } finally {
    state.saving = false
  }
}

async function testPush() {
  state.testingPush = true
  try {
    await configApi.testAuditAlertPush()
    window.$message.success('测试消息已通过默认推送引擎发送')
  } catch {
    window.$message.error('推送测试失败，请先在「消息推送」中配置 Webhook')
  } finally {
    state.testingPush = false
  }
}

async function testWebhook() {
  const url = state.webhookUrl.trim()
  if (!url) {
    window.$message.warning('请先填写自定义 Webhook 地址')
    return
  }
  state.testingWebhook = true
  try {
    await configApi.testAuditAlertWebhook({
      webhook_url: url,
      webhook_secret: state.webhookSecret,
    })
    window.$message.success('测试消息已发送，请检查 Webhook 接收端')
  } catch {
    window.$message.error('Webhook 测试失败，请检查 URL 和密钥')
  } finally {
    state.testingWebhook = false
  }
}
</script>

<template>
  <NSpin :show="state.loading">
    <ConfigSectionLayout
      description="配置审计告警开关、通知渠道与已实现规则。保存后下次任务调度执行即生效。"
      :saving="state.saving"
      @save="save"
      @reset="reset"
    >
      <NForm
        class="sys-config-form sys-config-form--wide"
        label-placement="top"
      >
        <NCard
          title="总开关与分析"
          size="small"
          :bordered="false"
        >
          <NGrid
            :cols="24"
            :x-gap="16"
          >
            <NGi :span="8">
              <NFormItem label="启用告警">
                <NSwitch v-model:value="state.enabled" />
              </NFormItem>
            </NGi>
            <NGi :span="8">
              <NFormItem label="分析周期（秒）">
                <NInputNumber
                  v-model:value="state.analysisInterval"
                  class="w-full"
                  :min="60"
                  :max="3600"
                />
              </NFormItem>
            </NGi>
            <NGi :span="8">
              <NFormItem label="告警冷却（秒）">
                <NInputNumber
                  v-model:value="state.cooldown"
                  class="w-full"
                  :min="60"
                  :max="86400"
                />
              </NFormItem>
            </NGi>
          </NGrid>
        </NCard>

        <NCard
          title="通知渠道"
          size="small"
          :bordered="false"
          class="mt-12px"
        >
          <div class="alert-rule">
            <div class="alert-rule__row alert-rule__row--channel">
              <div class="alert-rule__switch">
                <span class="alert-rule__label">邮件通知</span>
                <NSwitch v-model:value="state.notifyEmail" />
              </div>
              <div class="alert-rule__meta">
                <p class="sys-config__hint">
                  复用「邮件引擎」SMTP 配置
                  <template v-if="!state.mailConfigured"> （尚未配置） </template>
                </p>
                <NButton
                  text
                  type="primary"
                  @click="goConfig('MAIL')"
                >
                  去配置
                </NButton>
              </div>
            </div>

            <NFormItem
              v-if="state.notifyEmail"
              label="告警收件邮箱"
              class="mt-8px"
            >
              <NInput
                v-model:value="state.notifyEmailTo"
                placeholder="security@example.com"
              />
            </NFormItem>

            <div class="alert-rule__row alert-rule__row--channel">
              <div class="alert-rule__switch">
                <span class="alert-rule__label">消息推送</span>
                <NSwitch v-model:value="state.notifyPush" />
              </div>
              <div class="alert-rule__meta">
                <p class="sys-config__hint">
                  {{ pushHint }}
                </p>
                <NSpace>
                  <NButton
                    text
                    type="primary"
                    @click="goConfig('PUSH')"
                  >
                    去配置
                  </NButton>
                  <NButton
                    text
                    :loading="state.testingPush"
                    :disabled="!state.notifyPush"
                    @click="testPush"
                  >
                    测试推送
                  </NButton>
                </NSpace>
              </div>
            </div>

            <div class="alert-rule__row alert-rule__row--channel">
              <div class="alert-rule__switch">
                <span class="alert-rule__label">自定义 Webhook</span>
                <NSwitch v-model:value="state.notifyCustomWebhook" />
              </div>
              <p class="sys-config__hint">仅在需要独立接收地址时启用；一般优先用上方消息推送</p>
            </div>
          </div>

          <template v-if="state.notifyCustomWebhook">
            <NGrid
              :cols="24"
              :x-gap="16"
              class="mt-12px"
            >
              <NGi :span="16">
                <NFormItem label="Webhook 地址">
                  <div class="flex gap-8px w-full">
                    <NInput
                      v-model:value="state.webhookUrl"
                      placeholder="https://..."
                    />
                    <NButton
                      :loading="state.testingWebhook"
                      @click="testWebhook"
                    >
                      测试
                    </NButton>
                  </div>
                </NFormItem>
              </NGi>
              <NGi :span="8">
                <NFormItem label="签名密钥（留空表示不修改）">
                  <NInput
                    v-model:value="state.webhookSecret"
                    type="password"
                    show-password-on="click"
                  />
                </NFormItem>
              </NGi>
            </NGrid>
          </template>
        </NCard>

        <NCard
          title="告警规则"
          size="small"
          :bordered="false"
          class="mt-12px"
        >
          <div class="alert-rule">
            <div class="alert-rule__row">
              <div class="alert-rule__switch">
                <span class="alert-rule__label">暴力破解检测</span>
                <NSwitch v-model:value="state.ruleBruteForce" />
              </div>
              <NFormItem
                v-if="state.ruleBruteForce"
                label="窗口内审计条数阈值"
                class="alert-rule__extra"
                :show-feedback="false"
              >
                <NInputNumber
                  v-model:value="state.bruteForceThreshold"
                  class="w-full"
                  :min="1"
                  :max="100000"
                />
              </NFormItem>
            </div>

            <div class="alert-rule__row">
              <div class="alert-rule__switch">
                <span class="alert-rule__label">非常时段敏感操作</span>
                <NSwitch v-model:value="state.ruleUnusualHours" />
              </div>
              <p class="sys-config__hint">凌晨 0-6 点的角色/权限变更</p>
            </div>

            <div class="alert-rule__row">
              <div class="alert-rule__switch">
                <span class="alert-rule__label">敏感操作检测</span>
                <NSwitch v-model:value="state.ruleSensitiveOps" />
              </div>
              <p class="sys-config__hint">5 分钟内角色授权/权限变更次数</p>
            </div>

            <div class="alert-rule__row">
              <div class="alert-rule__switch">
                <span class="alert-rule__label">批量删除检测</span>
                <NSwitch v-model:value="state.ruleBulkDelete" />
              </div>
              <NFormItem
                v-if="state.ruleBulkDelete"
                label="5 分钟内删除次数阈值"
                class="alert-rule__extra"
                :show-feedback="false"
              >
                <NInputNumber
                  v-model:value="state.bulkDeleteThreshold"
                  class="w-full"
                  :min="1"
                  :max="100000"
                />
              </NFormItem>
            </div>

            <div class="alert-rule__row">
              <div class="alert-rule__switch">
                <span class="alert-rule__label">异地 IP 登录检测</span>
                <NSwitch v-model:value="state.ruleIpAnomaly" />
              </div>
              <NFormItem
                v-if="state.ruleIpAnomaly"
                label="15 分钟内不同 IP 数阈值"
                class="alert-rule__extra"
                :show-feedback="false"
              >
                <NInputNumber
                  v-model:value="state.ipAnomalyThreshold"
                  class="w-full"
                  :min="1"
                  :max="1000"
                />
              </NFormItem>
            </div>
          </div>
        </NCard>
      </NForm>
    </ConfigSectionLayout>
  </NSpin>
</template>

<style scoped>
.mt-12px {
  margin-top: 12px;
}

.mt-8px {
  margin-top: 8px;
}

.mb-12px {
  margin-bottom: 12px;
}

.flex {
  display: flex;
}

.gap-8px {
  gap: 8px;
}

.w-full {
  width: 100%;
}

.alert-rule {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.alert-rule__row {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr);
  gap: 24px;
  align-items: end;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.alert-rule__row:last-child {
  border-bottom: 0;
  padding-bottom: 0;
}

.alert-rule__row--channel {
  align-items: center;
}

.alert-rule__switch {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 34px;
  padding-bottom: 6px;
}

.alert-rule__row--channel .alert-rule__switch {
  padding-bottom: 0;
}

.alert-rule__label {
  font-size: 14px;
  color: var(--text-color-1);
}

.alert-rule__meta {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.alert-rule__extra {
  margin-bottom: 0;
}
</style>
