import { defineConfig, presetUno } from 'unocss'

export default defineConfig({
  presets: [presetUno()],
  shortcuts: {
    'wh-full': 'w-full h-full',
    'flex-center': 'flex items-center justify-center',
    'flex-y-center': 'flex items-center',
  },
})
