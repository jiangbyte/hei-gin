/** Author: Charlie */

import { authApi } from '@/api'
import { encryptFieldsWithPublicKey } from '@/utils/rsaEncrypt'

export async function encryptPasswords<T extends Record<string, string | null | undefined>>(
  fields: T,
) {
  const response = await authApi.passwordKey()
  return {
    password_key_id: response.data.key_id,
    values: await encryptFieldsWithPublicKey(response.data.public_key, fields),
  }
}
