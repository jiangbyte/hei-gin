/** Author: Charlie */

import * as authApi from '@/api/auth'

type RsaPublicKey = {
  modulus: Uint8Array
  exponent: Uint8Array
}

const bigZero = BigInt(0)
const bigOne = BigInt(1)
const bigEight = BigInt(8)
const bigByteMask = BigInt(0xff)

export async function encryptPasswords<
  T extends Record<string, string | null | undefined>,
>(fields: T) {
  const key = await authApi.passwordKey()
  const publicKey = parsePublicKey(base64ToBytes(key.public_key))
  const result: Record<string, string | null> = {}

  for (const [field, value] of Object.entries(fields)) {
    result[field] = value ? await encryptText(publicKey, value) : null
  }

  return {
    password_key_id: key.key_id,
    values: result as { [K in keyof T]: string | null },
  }
}

async function encryptText(publicKey: RsaPublicKey, value: string) {
  const message = utf8Encode(value)
  const encoded = await oaepEncode(message, publicKey.modulus.length)
  const encrypted = modPow(
    bytesToBigInt(encoded),
    bytesToBigInt(publicKey.exponent),
    bytesToBigInt(publicKey.modulus)
  )
  return bytesToBase64(bigIntToBytes(encrypted, publicKey.modulus.length))
}

async function oaepEncode(message: Uint8Array, keyLength: number) {
  const hashLength = 32
  if (message.length > keyLength - 2 * hashLength - 2) {
    throw new Error('密码过长，无法加密')
  }

  const labelHash = sha256(new Uint8Array())
  const padding = new Uint8Array(
    keyLength - message.length - 2 * hashLength - 2
  )
  const dataBlock = concatBytes(
    labelHash,
    padding,
    new Uint8Array([1]),
    message
  )
  const seed = await randomBytes(hashLength)
  const dataMask = mgf1(seed, keyLength - hashLength - 1)
  const maskedData = xorBytes(dataBlock, dataMask)
  const seedMask = mgf1(maskedData, hashLength)
  const maskedSeed = xorBytes(seed, seedMask)

  return concatBytes(new Uint8Array([0]), maskedSeed, maskedData)
}

function parsePublicKey(der: Uint8Array): RsaPublicKey {
  const root = readDerElement(der, 0, 0x30)
  const algorithm = readDerElement(der, root.valueStart, 0x30)
  const bitString = readDerElement(der, algorithm.end, 0x03)
  const rsaKey = readDerElement(der, bitString.valueStart + 1, 0x30)
  const modulus = readDerElement(der, rsaKey.valueStart, 0x02)
  const exponent = readDerElement(der, modulus.end, 0x02)

  return {
    modulus: trimPositiveInteger(
      der.slice(modulus.valueStart, modulus.valueEnd)
    ),
    exponent: trimPositiveInteger(
      der.slice(exponent.valueStart, exponent.valueEnd)
    ),
  }
}

function readDerElement(
  bytes: Uint8Array,
  offset: number,
  expectedTag: number
) {
  const tag = bytes[offset]
  if (tag !== expectedTag) {
    throw new Error('公钥无效')
  }
  const lengthByte = bytes[offset + 1]
  let length = lengthByte
  let valueStart = offset + 2

  if (lengthByte & 0x80) {
    const lengthBytes = lengthByte & 0x7f
    length = 0
    valueStart += lengthBytes
    for (let index = 0; index < lengthBytes; index += 1) {
      length = (length << 8) | bytes[offset + 2 + index]
    }
  }

  return {
    valueStart,
    valueEnd: valueStart + length,
    end: valueStart + length,
  }
}

function trimPositiveInteger(bytes: Uint8Array) {
  let offset = 0
  while (offset < bytes.length - 1 && bytes[offset] === 0) {
    offset += 1
  }
  return bytes.slice(offset)
}

function mgf1(seed: Uint8Array, length: number) {
  const output = new Uint8Array(length)
  let offset = 0
  let counter = 0

  while (offset < length) {
    const digest = sha256(concatBytes(seed, i2osp(counter, 4)))
    output.set(
      digest.slice(0, Math.min(digest.length, length - offset)),
      offset
    )
    offset += digest.length
    counter += 1
  }

  return output
}

function modPow(base: bigint, exponent: bigint, modulus: bigint) {
  let result = bigOne
  let nextBase = base % modulus
  let nextExponent = exponent

  while (nextExponent > bigZero) {
    if (nextExponent & bigOne) {
      result = (result * nextBase) % modulus
    }
    nextExponent >>= bigOne
    nextBase = (nextBase * nextBase) % modulus
  }

  return result
}

function bytesToBigInt(bytes: Uint8Array) {
  let value = bigZero
  for (const byte of bytes) {
    value = (value << bigEight) + BigInt(byte)
  }
  return value
}

function bigIntToBytes(value: bigint, length: number) {
  const bytes = new Uint8Array(length)
  let nextValue = value

  for (let index = length - 1; index >= 0; index -= 1) {
    bytes[index] = Number(nextValue & bigByteMask)
    nextValue >>= bigEight
  }

  return bytes
}

function i2osp(value: number, length: number) {
  const bytes = new Uint8Array(length)
  for (let index = length - 1; index >= 0; index -= 1) {
    bytes[index] = value & 0xff
    value >>>= 8
  }
  return bytes
}

function sha256(message: Uint8Array) {
  const constants = [
    0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1,
    0x923f82a4, 0xab1c5ed5, 0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3,
    0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174, 0xe49b69c1, 0xefbe4786,
    0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
    0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147,
    0x06ca6351, 0x14292967, 0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13,
    0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85, 0xa2bfe8a1, 0xa81a664b,
    0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
    0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a,
    0x5b9cca4f, 0x682e6ff3, 0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208,
    0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
  ]
  const words = new Uint32Array(64)
  const paddedLength = ((message.length + 9 + 63) >> 6) << 6
  const padded = new Uint8Array(paddedLength)
  padded.set(message)
  padded[message.length] = 0x80

  const bitLength = message.length * 8
  for (let index = 0; index < 8; index += 1) {
    padded[padded.length - 1 - index] =
      Math.floor(bitLength / 2 ** (index * 8)) & 0xff
  }

  let h0 = 0x6a09e667
  let h1 = 0xbb67ae85
  let h2 = 0x3c6ef372
  let h3 = 0xa54ff53a
  let h4 = 0x510e527f
  let h5 = 0x9b05688c
  let h6 = 0x1f83d9ab
  let h7 = 0x5be0cd19

  for (let chunk = 0; chunk < padded.length; chunk += 64) {
    for (let index = 0; index < 16; index += 1) {
      const offset = chunk + index * 4
      words[index] =
        ((padded[offset] << 24) |
          (padded[offset + 1] << 16) |
          (padded[offset + 2] << 8) |
          padded[offset + 3]) >>>
        0
    }
    for (let index = 16; index < 64; index += 1) {
      const s0 =
        rotateRight(words[index - 15], 7) ^
        rotateRight(words[index - 15], 18) ^
        (words[index - 15] >>> 3)
      const s1 =
        rotateRight(words[index - 2], 17) ^
        rotateRight(words[index - 2], 19) ^
        (words[index - 2] >>> 10)
      words[index] = (words[index - 16] + s0 + words[index - 7] + s1) >>> 0
    }

    let a = h0
    let b = h1
    let c = h2
    let d = h3
    let e = h4
    let f = h5
    let g = h6
    let h = h7

    for (let index = 0; index < 64; index += 1) {
      const s1 = rotateRight(e, 6) ^ rotateRight(e, 11) ^ rotateRight(e, 25)
      const ch = (e & f) ^ (~e & g)
      const temp1 = (h + s1 + ch + constants[index] + words[index]) >>> 0
      const s0 = rotateRight(a, 2) ^ rotateRight(a, 13) ^ rotateRight(a, 22)
      const maj = (a & b) ^ (a & c) ^ (b & c)
      const temp2 = (s0 + maj) >>> 0

      h = g
      g = f
      f = e
      e = (d + temp1) >>> 0
      d = c
      c = b
      b = a
      a = (temp1 + temp2) >>> 0
    }

    h0 = (h0 + a) >>> 0
    h1 = (h1 + b) >>> 0
    h2 = (h2 + c) >>> 0
    h3 = (h3 + d) >>> 0
    h4 = (h4 + e) >>> 0
    h5 = (h5 + f) >>> 0
    h6 = (h6 + g) >>> 0
    h7 = (h7 + h) >>> 0
  }

  return concatBytes(
    i2osp(h0, 4),
    i2osp(h1, 4),
    i2osp(h2, 4),
    i2osp(h3, 4),
    i2osp(h4, 4),
    i2osp(h5, 4),
    i2osp(h6, 4),
    i2osp(h7, 4)
  )
}

function rotateRight(value: number, bits: number) {
  return (value >>> bits) | (value << (32 - bits))
}

function concatBytes(...arrays: Uint8Array[]) {
  const length = arrays.reduce((total, array) => total + array.length, 0)
  const result = new Uint8Array(length)
  let offset = 0

  arrays.forEach((array) => {
    result.set(array, offset)
    offset += array.length
  })

  return result
}

function xorBytes(left: Uint8Array, right: Uint8Array) {
  const result = new Uint8Array(left.length)
  for (let index = 0; index < left.length; index += 1) {
    result[index] = left[index] ^ right[index]
  }
  return result
}

function utf8Encode(value: string) {
  const bytes: number[] = []
  for (let index = 0; index < value.length; index += 1) {
    let codePoint = value.charCodeAt(index)
    if (
      codePoint >= 0xd800 &&
      codePoint <= 0xdbff &&
      index + 1 < value.length
    ) {
      const next = value.charCodeAt(index + 1)
      if (next >= 0xdc00 && next <= 0xdfff) {
        codePoint = 0x10000 + ((codePoint - 0xd800) << 10) + (next - 0xdc00)
        index += 1
      }
    }
    if (codePoint < 0x80) {
      bytes.push(codePoint)
    } else if (codePoint < 0x800) {
      bytes.push(0xc0 | (codePoint >> 6), 0x80 | (codePoint & 0x3f))
    } else if (codePoint < 0x10000) {
      bytes.push(
        0xe0 | (codePoint >> 12),
        0x80 | ((codePoint >> 6) & 0x3f),
        0x80 | (codePoint & 0x3f)
      )
    } else {
      bytes.push(
        0xf0 | (codePoint >> 18),
        0x80 | ((codePoint >> 12) & 0x3f),
        0x80 | ((codePoint >> 6) & 0x3f),
        0x80 | (codePoint & 0x3f)
      )
    }
  }
  return new Uint8Array(bytes)
}

async function randomBytes(length: number) {
  const crypto = (globalThis as any).crypto
  if (crypto?.getRandomValues) {
    const bytes = new Uint8Array(length)
    crypto.getRandomValues(bytes)
    return bytes
  }

  const wxApi = (globalThis as any).wx
  if (wxApi?.getRandomValues) {
    return new Promise<Uint8Array>((resolve, reject) => {
      wxApi.getRandomValues({
        length,
        success: (res: { randomValues: ArrayBuffer }) =>
          resolve(new Uint8Array(res.randomValues)),
        fail: reject,
      })
    })
  }

  const uniApi = (globalThis as any).uni
  if (uniApi?.getRandomValues) {
    return new Promise<Uint8Array>((resolve, reject) => {
      uniApi.getRandomValues({
        length,
        success: (res: { randomValues: ArrayBuffer }) =>
          resolve(new Uint8Array(res.randomValues)),
        fail: reject,
      })
    })
  }

  throw new Error('安全随机数生成器不可用')
}

function base64ToBytes(value: string) {
  const chars =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'
  const clean = value.replace(/\s/g, '')
  const bytes: number[] = []

  for (let index = 0; index < clean.length; index += 4) {
    const first = chars.indexOf(clean[index])
    const second = chars.indexOf(clean[index + 1])
    const third =
      clean[index + 2] === '=' ? -1 : chars.indexOf(clean[index + 2])
    const fourth =
      clean[index + 3] === '=' ? -1 : chars.indexOf(clean[index + 3])
    const triplet =
      (first << 18) | (second << 12) | ((third & 0x3f) << 6) | (fourth & 0x3f)

    bytes.push((triplet >> 16) & 0xff)
    if (third >= 0) {
      bytes.push((triplet >> 8) & 0xff)
    }
    if (fourth >= 0) {
      bytes.push(triplet & 0xff)
    }
  }

  return new Uint8Array(bytes)
}

function bytesToBase64(bytes: Uint8Array) {
  const chars =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'
  let result = ''

  for (let index = 0; index < bytes.length; index += 3) {
    const first = bytes[index]
    const second = bytes[index + 1]
    const third = bytes[index + 2]
    const triplet = (first << 16) | ((second || 0) << 8) | (third || 0)

    result += chars[(triplet >> 18) & 0x3f]
    result += chars[(triplet >> 12) & 0x3f]
    result += index + 1 < bytes.length ? chars[(triplet >> 6) & 0x3f] : '='
    result += index + 2 < bytes.length ? chars[triplet & 0x3f] : '='
  }

  return result
}
