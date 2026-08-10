/** Author: Charlie
 *
 * 浏览器 WebCrypto RSA-OAEP 辅助。 */
export async function importSpkiPublicKey(publicKeyB64OrPem: string) {
  const binary = globalThis.atob(
    publicKeyB64OrPem
      .replace('-----BEGIN PUBLIC KEY-----', '')
      .replace('-----END PUBLIC KEY-----', '')
      .replace(/\s/g, ''),
  )
  const bytes = new Uint8Array(binary.length)
  for (let index = 0; index < binary.length; index += 1) {
    bytes[index] = binary.charCodeAt(index)
  }
  return globalThis.crypto.subtle.importKey(
    'spki',
    bytes.buffer,
    { name: 'RSA-OAEP', hash: 'SHA-256' },
    false,
    ['encrypt'],
  )
}

export async function encryptTextRsaOaep(publicKey: CryptoKey, value: string) {
  const encrypted = await globalThis.crypto.subtle.encrypt(
    { name: 'RSA-OAEP' },
    publicKey,
    new TextEncoder().encode(value),
  )
  return arrayBufferToBase64(encrypted)
}

function arrayBufferToBase64(buffer: ArrayBuffer) {
  const bytes = new Uint8Array(buffer)
  let binary = ''
  for (let index = 0; index < bytes.byteLength; index += 1) {
    binary += String.fromCharCode(bytes[index])
  }
  return globalThis.btoa(binary)
}

export async function encryptFieldsWithPublicKey<
  T extends Record<string, string | null | undefined>,
>(publicKeyB64OrPem: string, fields: T) {
  const publicKey = await importSpkiPublicKey(publicKeyB64OrPem)
  const result: Record<string, string | null> = {}
  for (const [key, value] of Object.entries(fields)) {
    result[key] = value ? await encryptTextRsaOaep(publicKey, value) : null
  }
  return result as { [K in keyof T]: string | null }
}
