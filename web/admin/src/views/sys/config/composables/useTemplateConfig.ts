/** Author: Charlie */

export const SCENE_PATTERN = /^[A-Z][A-Z0-9_]*$/

export const MAIL_SCENE_LABELS: Record<string, string> = {
  REGISTER_SUCCESS: '注册成功',
  LOGIN_CODE: '登录验证码',
  CHANGE_PASSWORD_CODE: '修改密码验证码',
  RESET_PASSWORD_CODE: '重置密码验证码',
  RESET_PASSWORD_SUCCESS: '重置密码成功',
  PASSWORD_EXPIRING: '密码即将过期',
  BIND_EMAIL_CODE: '绑定邮箱验证码',
  CHANGE_EMAIL_CODE: '修改邮箱验证码',
  ACCOUNT_CANCELLED: '账号注销确认',
  ACCOUNT_PURGED: '账号彻底删除',
}

export const SMS_SCENE_LABELS: Record<string, string> = {
  REGISTER_SUCCESS: '注册成功',
  LOGIN_CODE: '登录验证码',
  CHANGE_PASSWORD_CODE: '修改密码验证码',
  RESET_PASSWORD_CODE: '重置密码验证码',
  RESET_PASSWORD_SUCCESS: '重置密码成功',
  PASSWORD_EXPIRING: '密码即将过期',
  BIND_PHONE_CODE: '绑定手机验证码',
  CHANGE_PHONE_CODE: '修改手机验证码',
  ACCOUNT_CANCELLED: '账号注销确认',
  ACCOUNT_PURGED: '账号彻底删除',
}

/** 后端当前会发送的邮件场景（可编辑；标注「已接线」）。 */
export const MAIL_WIRED_SCENES = new Set([
  'REGISTER_SUCCESS',
  'LOGIN_CODE',
  'CHANGE_PASSWORD_CODE',
  'RESET_PASSWORD_CODE',
  'BIND_EMAIL_CODE',
  'ACCOUNT_CANCELLED',
  'ACCOUNT_PURGED',
])

/** 后端当前会发送的短信场景。 */
export const SMS_WIRED_SCENES = new Set([
  'LOGIN_CODE',
  'CHANGE_PASSWORD_CODE',
  'BIND_PHONE_CODE',
  'ACCOUNT_CANCELLED',
  'ACCOUNT_PURGED',
])

export function isMailSceneWired(scene: string): boolean {
  return MAIL_WIRED_SCENES.has(normalizeScene(scene))
}

export function isSmsSceneWired(scene: string): boolean {
  return SMS_WIRED_SCENES.has(normalizeScene(scene))
}

export interface MailTemplateValue {
  subject: string
  body: string
}

export interface SmsTemplateValue {
  code: string
  content: string
}

/** 全局模板键：`MAIL_TEMPLATE_{SCENE}` / `SMS_TEMPLATE_{SCENE}` */
export function templateConfigKey(prefix: 'MAIL_TEMPLATE' | 'SMS_TEMPLATE', scene: string): string {
  return `${prefix}_${scene}`
}

export function parseMailTemplate(raw: string | null | undefined): MailTemplateValue {
  if (!raw) return { subject: '', body: '' }
  try {
    const data = JSON.parse(raw)
    return {
      subject: String(data?.subject ?? ''),
      body: String(data?.body ?? ''),
    }
  } catch {
    return { subject: '', body: String(raw) }
  }
}

export function stringifyMailTemplate(value: MailTemplateValue): string {
  return JSON.stringify({
    subject: value.subject ?? '',
    body: value.body ?? '',
  })
}

export function parseSmsTemplate(raw: string | null | undefined): SmsTemplateValue {
  if (!raw) return { code: '', content: '' }
  try {
    const data = JSON.parse(raw)
    return {
      code: String(data?.code ?? ''),
      content: String(data?.content ?? ''),
    }
  } catch {
    return { code: '', content: String(raw) }
  }
}

export function stringifySmsTemplate(value: SmsTemplateValue): string {
  return JSON.stringify({
    code: value.code ?? '',
    content: value.content ?? '',
  })
}

export function normalizeScene(scene: string): string {
  return String(scene || '')
    .trim()
    .toUpperCase()
    .replace(/[^A-Z0-9_]/g, '_')
    .replace(/_+/g, '_')
    .replace(/^_|_$/g, '')
}
