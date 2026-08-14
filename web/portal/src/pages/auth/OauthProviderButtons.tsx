/** Author: Charlie */

type OauthProvider = {
  provider: string
  label: string
}

type Props = {
  providers: OauthProvider[]
  loadingProvider?: string | null
  onSelect: (provider: string) => void
  hint?: string
}

function providerInitial(provider: string, label: string) {
  const key = provider.toUpperCase()
  if (key === 'GITHUB') return 'GH'
  if (key === 'GITEE') return 'GE'
  if (key === 'QQ') return 'QQ'
  if (key.startsWith('WECHAT')) return '微'
  return (label || provider).slice(0, 1).toUpperCase()
}

/** 三方登录圆形入口 */
export function OauthProviderButtons({ providers, loadingProvider, onSelect, hint }: Props) {
  if (!providers.length) return null
  return (
    <div className="auth-oauth">
      <div className="auth-oauth__divider">
        <span>其他登录方式</span>
      </div>
      <div className="auth-oauth__row">
        {providers.map((item) => {
          const busy = loadingProvider === item.provider
          return (
            <button
              key={item.provider}
              type="button"
              className="auth-oauth__icon"
              title={item.label}
              aria-label={item.label}
              disabled={Boolean(loadingProvider)}
              onClick={() => onSelect(item.provider)}
            >
              <span className={busy ? 'auth-oauth__icon-label is-busy' : 'auth-oauth__icon-label'}>
                {providerInitial(item.provider, item.label)}
              </span>
            </button>
          )
        })}
      </div>
      {hint ? <p className="auth-oauth__hint">{hint}</p> : null}
    </div>
  )
}
