/** Author: Charlie */

import type { FormItemRule, FormRules } from 'naive-ui'
import { createRequiredRule, isValidEmail, toNullableString } from '@/utils'
import { dictList } from '@/utils/dict'
import { encryptPasswords } from '@/utils/security'

export type AccountBaseForm = {
  account: string
  password: string
  account_type: string
  account_status: string
  name: string
  nickname: string
  avatar: string
  signature: string
  phone: string
  email: string
  email_login_enabled: boolean
  phone_login_enabled: boolean
}

export function accountStatusOptions() {
  return dictList('ACCOUNT_STATUS').filter((o: any) => !String(o.value).includes('CANCELLED'))
}

export function createBaseFormRules(
  isEdit: () => boolean,
  form: { email_login_enabled: boolean },
): FormRules {
  return {
    account: createRequiredRule('账号', 'input'),
    password: [
      {
        validator: (_rule, value) => {
          if (!isEdit() && !String(value ?? '').trim()) {
            return new Error('请输入密码')
          }
          return true
        },
        trigger: ['input', 'blur'],
      },
    ],
    account_status: createRequiredRule('账号状态', 'change'),
    email: [
      {
        validator: (_rule: FormItemRule, value: string) => {
          const text = String(value ?? '').trim()
          if (!text) {
            return form.email_login_enabled ? new Error('请输入邮箱') : true
          }
          if (!isValidEmail(text)) {
            return new Error('请输入有效邮箱')
          }
          return true
        },
        trigger: ['input', 'blur'],
      },
    ],
  }
}

/**
 * 组装 AccountCreate/Update 公共字段（不含各体系专属资料字段）。
 * - SysAccount: account_type / account_status / password
 * - SysAccountIdentity: account / email|phone + login flags
 * - 资料表共有: name / nickname / avatar / signature / phone / email
 */
export async function buildBaseAccountPayload(form: AccountBaseForm) {
  const payload: Record<string, unknown> = {
    account: form.account.trim(),
    password: toNullableString(form.password),
    account_type: form.account_type,
    account_status: form.account_status,
    name: form.name.trim(),
    nickname: toNullableString(form.nickname),
    avatar: toNullableString(form.avatar),
    signature: toNullableString(form.signature),
    phone: toNullableString(form.phone),
    email: toNullableString(form.email),
    email_login_enabled: Boolean(form.email_login_enabled),
    phone_login_enabled: Boolean(form.phone_login_enabled),
    email_identity: form.email_login_enabled ? toNullableString(form.email) : null,
    phone_identity: form.phone_login_enabled ? toNullableString(form.phone) : null,
    email_identity_verified: Boolean(form.email_login_enabled),
    phone_identity_verified: Boolean(form.phone_login_enabled),
    email_identity_bind_status: 'BOUND',
    phone_identity_bind_status: 'BOUND',
  }

  if (payload.password) {
    const encrypted = await encryptPasswords({ password: String(payload.password) })
    payload.password = encrypted.values.password
    payload.password_key_id = encrypted.password_key_id
  } else {
    payload.password_key_id = null
  }

  return payload
}
