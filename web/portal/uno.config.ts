import { defineConfig, presetUno } from 'unocss'

export default defineConfig({
  preflights: [
    {
      getCSS: () => `
        *,
        *::before,
        *::after {
          border-width: 0;
          border-style: solid;
          border-color: color-mix(in srgb, var(--ant-color-border, #e5e7eb) 70%, transparent);
        }
      `,
    },
  ],
  presets: [presetUno()],
  shortcuts: {
    'wh-full': 'w-full h-full',
    'flex-center': 'flex items-center justify-center',
    'flex-y-center': 'flex items-center',
    'page-shell': 'text-[var(--ant-color-text)]',
    workspace: 'bg-[var(--ant-color-bg-layout)] text-[var(--ant-color-text)]',
    // 校园练习台：轻描边 + 柔和阴影，避免厚边塑料感
    panel:
      'overflow-hidden rounded-xl border border-[color-mix(in_srgb,var(--ant-color-border)_65%,transparent)] bg-[var(--ant-color-bg-container)] shadow-[0_2px_8px_rgba(22,119,255,0.04)]',
    'panel-header':
      'flex min-h-10 items-center gap-2 border-b border-[color-mix(in_srgb,var(--ant-color-border)_55%,transparent)] bg-[var(--ant-color-bg-container)] px-3',
    'panel-body': 'bg-[var(--ant-color-bg-container)]',
    toolbar:
      'border-b border-[color-mix(in_srgb,var(--ant-color-border)_55%,transparent)] bg-[var(--ant-color-bg-container)]',
    'muted-box':
      'border border-dashed border-[color-mix(in_srgb,var(--ant-color-border)_70%,transparent)] bg-[var(--ant-color-fill-quaternary)] text-[var(--ant-color-text-secondary)]',
    'result-box':
      'border border-[color-mix(in_srgb,var(--ant-color-border)_55%,transparent)] bg-[var(--ant-color-fill-quaternary)]',
    'case-row':
      'border border-[color-mix(in_srgb,var(--ant-color-border)_55%,transparent)] bg-[var(--ant-color-bg-container)]',
    'editor-shell':
      'overflow-hidden rounded-lg border border-[color-mix(in_srgb,var(--ant-color-border)_55%,transparent)] bg-[var(--ant-color-bg-container)]',
    'muted-text': 'text-[var(--ant-color-text-secondary)]',
    'list-row': 'transition-colors hover:bg-[var(--ant-color-fill-quaternary)]',
    chip: 'inline-flex items-center rounded-full px-3 py-1 text-sm transition-colors',
    'error-box': 'bg-[var(--ant-color-error-bg)] text-[var(--ant-color-error)]',
    'tabs-shell': 'text-[var(--ant-color-text)]',
    'tabs-bar':
      'border-b border-[color-mix(in_srgb,var(--ant-color-border)_55%,transparent)] bg-[var(--ant-color-bg-container)]',
    'tabs-btn':
      'text-[var(--ant-color-text-secondary)] hover:bg-[var(--ant-color-fill-quaternary)] hover:text-[var(--ant-color-text)]',
    'tabs-btn-active': 'bg-transparent text-[var(--ant-color-text)]',
  },
})
