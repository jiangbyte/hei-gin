import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import UnoCSS from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'

const rootDir = path.dirname(fileURLToPath(import.meta.url))

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const port = Number(env.VITE_PORT || 5174)
  // Cookie 会话需同源；开发时由 Vite 反代到后端，避免跨域丢 Authorization。
  const proxyTarget = env.VITE_API_PROXY_TARGET || 'http://127.0.0.1:8000'
  const apiProxy = {
    '/api': {
      target: proxyTarget,
      changeOrigin: true,
    },
  }

  return {
    plugins: [
      react(),
      UnoCSS(),
      AutoImport({
        dts: 'src/typing/auto-imports.d.ts',
        imports: [
          'react',
          'react-router-dom',
          {
            zustand: ['create'],
            'zustand/middleware': ['persist'],
            antd: ['message', 'App'],
          },
        ],
        eslintrc: {
          enabled: true,
          filepath: './.eslintrc-auto-import.json',
        },
      }),
    ],
    resolve: {
      alias: {
        '@': path.resolve(rootDir, 'src'),
      },
    },
    server: {
      host: '0.0.0.0',
      port,
      strictPort: false,
      proxy: apiProxy,
    },
    preview: {
      host: '0.0.0.0',
      port,
      strictPort: false,
      proxy: apiProxy,
    },
  }
})
