/** Author: Charlie */

import { Carousel } from 'antd'
import { Link } from 'react-router-dom'

export type PromoSlide = {
  key: string
  title: string
  desc: string
  tag?: string
  /** 内部 SPA 路由 */
  to?: string
  /** 外部 / 绝对 URL */
  href?: string
  /** Tailwind 渐变类，如 from-[...] to-[...] */
  tone?: string
  cta?: string
  /** 可选封面图；按固定幻灯片高度裁剪 */
  imageUrl?: string
  onClick?: () => void
}

type Props = {
  slides: PromoSlide[]
  /** 固定幻灯片高度 — 后续加图时保持布局稳定 */
  height?: number
  className?: string
  autoplay?: boolean
}

const DEFAULT_TONES = [
  'from-[var(--ant-color-primary)] to-[var(--ant-color-primary-active)]',
  'from-[var(--ant-color-info)] to-[var(--ant-color-primary)]',
  'from-[var(--ant-color-success)] to-[var(--ant-color-info)]',
  'from-[var(--ant-color-warning)] to-[var(--ant-color-error)]',
]

export function PromoCarousel({ slides, height = 240, className, autoplay = true }: Props) {
  if (!slides.length) return null

  return (
    <div className={`panel overflow-hidden ${className ?? ''}`}>
      <Carousel
        autoplay={autoplay}
        autoplaySpeed={4500}
        dots
        effect="scrollx"
        className="promo-carousel"
      >
        {slides.map((slide, index) => {
          const tone = slide.tone ?? DEFAULT_TONES[index % DEFAULT_TONES.length]
          const body = (
            <div
              className={`relative flex w-full items-center overflow-hidden bg-gradient-to-br ${tone} px-7 py-8 text-white md:px-10`}
              style={{ height }}
            >
              {slide.imageUrl ? (
                <img
                  src={slide.imageUrl}
                  alt=""
                  className="pointer-events-none absolute inset-0 h-full w-full object-cover"
                />
              ) : null}
              <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_85%_20%,rgba(255,255,255,0.22),transparent_45%)]" />
              {slide.imageUrl ? (
                <div className="pointer-events-none absolute inset-0 bg-black/25" />
              ) : null}
              <div className="relative z-[1] max-w-2xl">
                {slide.tag ? (
                  <span className="mb-3 inline-flex rounded-full bg-white/15 px-3 py-1 text-xs tracking-wide">
                    {slide.tag}
                  </span>
                ) : null}
                <div className="line-clamp-2 text-2xl font-semibold leading-snug md:text-[28px]">
                  {slide.title}
                </div>
                {slide.desc ? (
                  <div className="mt-3 line-clamp-2 max-w-xl text-sm leading-relaxed text-white/88 md:text-base">
                    {slide.desc}
                  </div>
                ) : null}
                {slide.cta ? (
                  <span className="mt-5 inline-flex rounded-lg bg-white/15 px-4 py-2 text-sm ring-1 ring-white/25">
                    {slide.cta}
                  </span>
                ) : null}
              </div>
            </div>
          )

          if (slide.to) {
            return (
              <Link key={slide.key} to={slide.to} className="block" onClick={slide.onClick}>
                {body}
              </Link>
            )
          }

          if (slide.href) {
            return (
              <a
                key={slide.key}
                href={slide.href}
                className="block"
                target="_blank"
                rel="noreferrer"
                onClick={slide.onClick}
              >
                {body}
              </a>
            )
          }

          return (
            <div
              key={slide.key}
              onClick={slide.onClick}
              role={slide.onClick ? 'button' : undefined}
            >
              {body}
            </div>
          )
        })}
      </Carousel>
    </div>
  )
}
