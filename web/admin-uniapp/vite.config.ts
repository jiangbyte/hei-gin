import { defineConfig, loadEnv } from 'vite'
import { fileURLToPath, URL } from 'node:url'
import uni from '@dcloudio/vite-plugin-uni'

export default defineConfig(async ({ mode }) => {
  const env = loadEnv(mode, __dirname, '')
  const port = Number(env.VITE_PORT || 5173)

  const UnoCSS = (await import('unocss/vite')).default
  return {
    server: {
      host: '0.0.0.0',
      port,
    },
    plugins: [UnoCSS(), uni()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
  }
})
