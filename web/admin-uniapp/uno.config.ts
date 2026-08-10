import { presetUni } from '@uni-helper/unocss-preset-uni'
import {
  defineConfig,
  presetIcons,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'
import { presetWeapp } from 'unocss-preset-weapp'
import {
  transformerAttributify,
  transformerClass,
} from 'unocss-preset-weapp/transformer'

const isH5 = process.env.UNI_PLATFORM === 'h5'
export default defineConfig({
  presets: isH5
    ? [
        presetUni(),
        presetIcons({
          scale: 1.2,
          warn: true,
          extraProperties: {
            display: 'inline-block',
            'vertical-align': 'middle',
          },
        }),
      ]
    : [
        presetWeapp({
          platform: 'uniapp',
          // prefix: 'u-',
          // whRpx: true,
          // preflight: true,
          // arbitraryVariants: true,
        }),
      ],
  transformers: isH5
    ? [transformerDirectives(), transformerVariantGroup()]
    : [transformerAttributify(), transformerClass()],
  shortcuts: {
    'btn-primary': 'px-4 py-2 rounded bg-blue-500 text-white hover:bg-blue-600',
    card: 'bg-white rounded-lg shadow p-4',
  },
})
