/** Author: Charlie */

/**
 * 与后端 AccountType 对齐的唯一前端真相源。
 * 扩展新账户类型：在此追加一项，并同步后端枚举 + sys_config 种子键
 *（模式：`AUTH_REGISTER_{TYPE}_*` / `AUTH_LOGIN_{TYPE}_*` / `PASSWORD_{TYPE}_*`）。
 * 邮件/短信模板为全局 `MAIL_TEMPLATE_{SCENE}` / `SMS_TEMPLATE_{SCENE}`，不按类型拆分。
 */
export const ACCOUNT_TYPE = [
  { value: 'ADMIN', label: '管理员' },
  { value: 'PORTAL', label: '门户用户' },
] as const

export type AccountType = (typeof ACCOUNT_TYPE)[number]['value']

export const ACCOUNT_TYPES = ACCOUNT_TYPE.map((item) => item.value) as AccountType[]

export const DEFAULT_ACCOUNT_TYPE: AccountType = ACCOUNT_TYPE[0].value

/** NTab：name/tab */
export const ACCOUNT_TYPE_TABS = ACCOUNT_TYPE.map(({ value, label }) => ({
  key: value,
  label,
}))

/** NSelect / ProSearchForm select：label/value */
export const ACCOUNT_TYPE_OPTIONS = ACCOUNT_TYPE.map(({ value, label }) => ({
  label,
  value,
}))

export function accountTypeLabel(type: string): string {
  return ACCOUNT_TYPE.find((item) => item.value === type)?.label || type
}

/** 按账户类型初始化一份表单/状态映射，扩展类型时无需改业务代码。 */
export function createAccountTypeMap<T>(factory: () => T): Record<AccountType, T> {
  return Object.fromEntries(ACCOUNT_TYPE.map(({ value }) => [value, factory()])) as Record<
    AccountType,
    T
  >
}

export function mapAccountTypes<T>(fn: (type: AccountType) => T): T[] {
  return ACCOUNT_TYPE.map(({ value }) => fn(value))
}

/** 配置键：`PREFIX_ADMIN_FIELD` → 扩展类型自动拼出。 */
export function accountConfigKey(
  prefix: string,
  type: AccountType | string,
  field: string,
): string {
  return `${prefix}_${type}_${field}`
}
