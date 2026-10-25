/**
 * 模型广场展示计价工具：块网格与价格弹窗共用。
 * 实付口径与后端计费一致（渠道/分组定价 × 生效倍率），官方参考价不乘倍率。
 */

import { computed } from 'vue'
import { formatScaled, resolveIntervalPrices } from '@/utils/pricing'
import {
  BILLING_MODE_TOKEN,
  BILLING_MODE_IMAGE,
  BILLING_MODE_VIDEO,
  REASONING_EFFORT_LEVELS,
  type BillingMode
} from '@/constants/channel'
import type { PlazaModel, PlazaTimePricingPeriod } from '@/api/modelPlaza'
import type { UserPricingInterval } from '@/api/channels'

const PER_MILLION = 1_000_000
const MIN_DECIMALS = 2

export interface PlazaRateContext {
  rateMultiplier: number
  userRateMultiplier?: number | null
  imageRateIndependent?: boolean
  imageRateMultiplier?: number | null
  videoRateIndependent?: boolean
  videoRateMultiplier?: number | null
  peakWindow?: string
  peakRateMultiplier?: number | null
}

export function usePlazaPricing(props: PlazaRateContext) {
  const effectiveRate = computed(() => props.userRateMultiplier ?? props.rateMultiplier)
  const hasCustomRate = computed(
    () => props.userRateMultiplier != null && props.userRateMultiplier !== props.rateMultiplier
  )

  function billingMode(m: PlazaModel): BillingMode {
    return (m.pricing?.billing_mode || BILLING_MODE_TOKEN) as BillingMode
  }

  function usesIndependentMediaRate(m: PlazaModel): boolean {
    return (billingMode(m) === BILLING_MODE_IMAGE && props.imageRateIndependent === true) ||
      (billingMode(m) === BILLING_MODE_VIDEO && props.videoRateIndependent === true)
  }

  /** 非 token 行的生效倍率，图片和视频可分别覆盖通用倍率。 */
  function requestRate(m: PlazaModel): number {
    if (billingMode(m) === BILLING_MODE_VIDEO && props.videoRateIndependent === true) {
      return Math.max(0, props.videoRateMultiplier ?? 1)
    }
    return billingMode(m) === BILLING_MODE_IMAGE && props.imageRateIndependent === true
      ? (props.imageRateMultiplier ?? 1)
      : effectiveRate.value
  }

  /** 时段行的生效倍率 = 生效倍率 × 时段倍率（去掉浮点噪声）。 */
  function periodRate(period: PlazaTimePricingPeriod): number {
    return Math.round(effectiveRate.value * period.multiplier * 1000) / 1000
  }

  /** 实付价 = 渠道单价 × 生效倍率（时段行再乘时段倍率），按 $/1M token 展示。 */
  function paidPerMillion(value: number | null | undefined, period: PlazaTimePricingPeriod | null = null): string {
    if (value == null) return '-'
    const rate = period ? periodRate(period) : effectiveRate.value
    return formatScaled(value * rate, PER_MILLION, MIN_DECIMALS)
  }

  /** 非 token 单价乘该行生效倍率，不换算为每百万 token。 */
  function paidRequestPrice(m: PlazaModel, value: number | null | undefined): string {
    if (value == null) return '-'
    return formatScaled(value * requestRate(m), 1, MIN_DECIMALS)
  }

  /** 官方参考价不乘倍率。 */
  function official(value: number | null | undefined): string {
    if (value == null) return '-'
    return formatScaled(value, PER_MILLION, MIN_DECIMALS)
  }

  return {
    effectiveRate,
    hasCustomRate,
    billingMode,
    usesIndependentMediaRate,
    requestRate,
    periodRate,
    paidPerMillion,
    paidRequestPrice,
    official
  }
}

export function plazaBillingMode(m: PlazaModel): BillingMode {
  return (m.pricing?.billing_mode || BILLING_MODE_TOKEN) as BillingMode
}

/** 上下文档位按下限升序展示（后端已升序，此处兜底）。 */
export function sortByContext(intervals: UserPricingInterval[]): UserPricingInterval[] {
  return [...intervals].sort((a, b) => a.min_tokens - b.min_tokens)
}

/** token 模式的阶梯定价（内联进输入/输出/缓存列）。 */
export function tokenIntervals(m: PlazaModel): UserPricingInterval[] {
  return sortByContext(m.pricing?.intervals ?? []).map(iv => resolveIntervalPrices(iv, m.pricing!))
}

/** 官方阶梯（后端按目录规则合成，不受分组开关影响）。 */
export function officialIntervals(m: PlazaModel): UserPricingInterval[] {
  return sortByContext(m.official_pricing?.intervals ?? [])
}

/** 按次/按图模式的阶梯定价（仅保留配了按次价的档位）。 */
export function requestIntervals(m: PlazaModel): UserPricingInterval[] {
  return (m.pricing?.intervals ?? []).filter((iv) => iv.per_request_price != null)
}

/** 任一档带缓存价才按档渲染缓存列；否则沿用平价的写入/读取两行。 */
export function hasTierCachePricing(intervals: UserPricingInterval[]): boolean {
  return intervals.some((iv) =>
    iv.cache_write_price != null || iv.cache_write_1h_price != null || iv.cache_read_price != null ||
    iv.cache_write_multiplier != null || iv.cache_read_multiplier != null
  )
}

export function hasCachePricing(m: PlazaModel): boolean {
  return m.pricing?.cache_write_price != null || m.pricing?.cache_write_1h_price != null || m.pricing?.cache_read_price != null
}

export function hasOfficialCache(o: NonNullable<PlazaModel['official_pricing']>): boolean {
  return o.cache_write_price != null || o.cache_read_price != null || o.cache_write_1h_price != null
}

/**
 * 档位标签：优先后端/管理员给出的 tier_label，否则按区间生成统一形态——
 * 有上限为「≤上限」，末档为「>下限」；档位升序排列。
 */
export function tierLabel(iv: UserPricingInterval): string {
  if (iv.tier_label) return iv.tier_label
  const { min_tokens: min, max_tokens: max } = iv
  return max == null ? `>${formatTokenCount(min)}` : `≤${formatTokenCount(max)}`
}

/** 「00:30–08:30」；整分钟的 HH:mm:ss 省略秒。 */
export function formatTimeWindow(p: PlazaTimePricingPeriod): string {
  const clock = (v: string) => v.replace(/^(\d{2}:\d{2}):00$/, '$1')
  return `${clock(p.start_time)}–${clock(p.end_time)}`
}

export function formatTokenCount(n: number): string {
  if (n >= 1_000_000) return `${trimZero(n / 1_000_000)}M`
  if (n >= 1_000) return `${trimZero(n / 1_000)}K`
  return String(n)
}

function trimZero(n: number): string {
  return String(Math.round(n * 100) / 100)
}

/** 思考等级倍率徽章数据（保持展示顺序）。 */
export function reasoningEffortMultipliers(model: PlazaModel): [string, number][] {
  const multipliers = model.pricing?.reasoning_effort_multipliers
  return REASONING_EFFORT_LEVELS.flatMap(effort => {
    const multiplier = multipliers?.[effort]
    return typeof multiplier === 'number' && Number.isFinite(multiplier) && multiplier > 0
      ? [[effort, multiplier] as [string, number]]
      : []
  })
}
