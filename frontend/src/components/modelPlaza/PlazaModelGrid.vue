<template>
  <div class="grid grid-cols-1 gap-2.5 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
    <button
      v-for="m in sortedModels"
      :key="`${m.platform}:${m.name}`"
      type="button"
      class="group flex flex-col gap-2 rounded-xl border border-gray-200 bg-white px-4 py-3 text-left shadow-sm transition-all hover:-translate-y-0.5 hover:border-primary-300 hover:shadow-card focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-500/40 dark:border-dark-700/60 dark:bg-dark-800/60 dark:hover:border-primary-500/40"
      @click="emit('select', m)"
    >
      <!-- 行1:模型名 + 模式徽章 -->
      <div class="flex min-w-0 items-center gap-1.5">
        <span class="truncate font-medium text-gray-900 dark:text-white" :title="m.name">{{ m.name }}</span>
        <span
          v-if="platform && m.platform !== platform"
          :class="[
            'shrink-0 rounded-md px-1.5 py-0.5 text-[10px] font-medium',
            platformBadgeLightClass(m.platform)
          ]"
        >
          {{ platformLabel(m.platform) }}
        </span>
      </div>

      <!-- 行2:本站简价 -->
      <div class="flex items-baseline justify-between gap-2">
        <template v-if="billingMode(m) === BILLING_MODE_TOKEN">
          <div class="min-w-0">
            <template v-if="!m.pricing">
              <span class="text-xs text-gray-400 dark:text-dark-500">{{ t('modelPlaza.detail.noPricing') }}</span>
            </template>
            <template v-else>
              <div class="font-mono text-xs leading-5 text-gray-800 dark:text-gray-200">
                <span class="font-sans text-gray-400 dark:text-dark-500">{{ t('modelPlaza.table.input') }}</span>
                {{ paidPerMillion(basePaidInput(m)) }}
              </div>
              <div class="font-mono text-xs leading-5 text-gray-800 dark:text-gray-200">
                <span class="font-sans text-gray-400 dark:text-dark-500">{{ t('modelPlaza.table.output') }}</span>
                {{ paidPerMillion(basePaidOutput(m)) }}
              </div>
            </template>
          </div>
        </template>
        <template v-else>
          <div class="min-w-0">
            <span v-if="firstRequestPrice(m) != null" class="font-mono text-xs text-gray-800 dark:text-gray-200">
              {{ paidRequestPrice(m, firstRequestPrice(m)) }}
              <span class="font-sans text-gray-400 dark:text-dark-500">{{ perUnitSuffix(m) }}</span>
            </span>
            <span v-else class="text-xs text-gray-400 dark:text-dark-500">{{ t('modelPlaza.detail.noPricing') }}</span>
          </div>
        </template>
        <span class="shrink-0 font-mono text-[11px] font-semibold text-gray-500 dark:text-dark-400">
          ×{{ gridRate(m) }}
        </span>
      </div>

      <!-- 行3:能力徽章(阶梯/分时/独立倍率) -->
      <div v-if="hasBadges(m)" class="flex flex-wrap gap-1">
        <span
          v-if="tokenIntervals(m).length > 1 || officialIntervals(m).length > 1"
          class="rounded-md bg-gray-100 px-1.5 py-0.5 text-[10px] font-medium text-gray-500 dark:bg-dark-700/70 dark:text-dark-300"
          :title="tierHint(m)"
        >
          {{ t('modelPlaza.badges.longContext') }}
        </span>
        <span
          v-if="timePeriods(m).length"
          class="rounded-md bg-gray-100 px-1.5 py-0.5 text-[10px] font-medium text-gray-500 dark:bg-dark-700/70 dark:text-dark-300"
          :title="timePricingHint(m)"
        >
          {{ t('modelPlaza.badges.timePricing') }}
        </span>
        <span
          v-if="m.long_context_basis === 'marginal'"
          class="rounded-md bg-gray-100 px-1.5 py-0.5 text-[10px] font-medium text-gray-500 dark:bg-dark-700/70 dark:text-dark-300"
          :title="t('modelPlaza.table.tierHintMarginal')"
        >
          {{ t('modelPlaza.table.marginalBadge') }}
        </span>
      </div>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { platformBadgeLightClass, platformLabel } from '@/utils/platformColors'
import { BILLING_MODE_TOKEN, BILLING_MODE_IMAGE } from '@/constants/channel'
import type { PlazaModel } from '@/api/modelPlaza'
import {
  usePlazaPricing,
  plazaBillingMode,
  tokenIntervals,
  officialIntervals,
  requestIntervals,
  formatTimeWindow
} from './plazaPricing'

const props = defineProps<{
  models: PlazaModel[]
  /** 分组平台;跨平台模型(composite 分组)显示平台徽章。 */
  platform?: string
  rateMultiplier: number
  userRateMultiplier?: number | null
  imageRateIndependent?: boolean
  imageRateMultiplier?: number | null
  videoRateIndependent?: boolean
  videoRateMultiplier?: number | null
}>()

const emit = defineEmits<{ (e: 'select', model: PlazaModel): void }>()

const { t } = useI18n()
const {
  effectiveRate,
  billingMode,
  requestRate,
  paidPerMillion,
  paidRequestPrice
} = usePlazaPricing(props)

/** 与原表格一致的排序:token 模型在前,官方输出价降序,同价按名称降序。 */
const sortedModels = computed(() => {
  return [...props.models].sort((a, b) => {
    const ta = billingMode(a) === BILLING_MODE_TOKEN
    const tb = billingMode(b) === BILLING_MODE_TOKEN
    if (ta !== tb) return ta ? -1 : 1
    const pa = a.official_pricing?.output_price ?? null
    const pb = b.official_pricing?.output_price ?? null
    if (pa != null && pb != null && pa !== pb) return pb - pa
    if (pa != null && pb == null) return -1
    if (pa == null && pb != null) return 1
    return b.name.localeCompare(a.name)
  })
})

/** token 模型块展示首档实付价(与弹窗明细的基准档一致)。 */
function basePaidInput(m: PlazaModel): number | null {
  const ivs = tokenIntervals(m)
  if (ivs.length) return ivs[0].input_price
  return m.pricing?.input_price ?? null
}

function basePaidOutput(m: PlazaModel): number | null {
  const ivs = tokenIntervals(m)
  if (ivs.length) return ivs[0].output_price
  return m.pricing?.output_price ?? null
}

/** 按图/按次模型块展示最低档按次价。 */
function firstRequestPrice(m: PlazaModel): number | null {
  const ivs = requestIntervals(m)
  if (ivs.length) return ivs[0].per_request_price ?? null
  return m.pricing?.per_request_price ?? null
}

function gridRate(m: PlazaModel): string {
  const rate = billingMode(m) === BILLING_MODE_TOKEN ? effectiveRate.value : requestRate(m)
  return String(Math.round(rate * 1000) / 1000)
}

function perUnitSuffix(m: PlazaModel): string {
  return plazaBillingMode(m) === BILLING_MODE_IMAGE
    ? t('modelPlaza.table.perUnitImage')
    : t('modelPlaza.table.perUnitRequest')
}

function hasBadges(m: PlazaModel): boolean {
  return (
    tokenIntervals(m).length > 1 ||
    officialIntervals(m).length > 1 ||
    (m.time_pricing?.periods?.length ?? 0) > 0 ||
    m.long_context_basis === 'marginal'
  )
}

function timePeriods(m: PlazaModel) {
  return m.time_pricing?.periods ?? []
}

function tierHint(m: PlazaModel): string {
  return m.long_context_basis === 'marginal'
    ? t('modelPlaza.table.tierHintMarginal')
    : t('modelPlaza.table.tierHint')
}

function timePricingHint(m: PlazaModel): string {
  const periods = timePeriods(m).map(formatTimeWindow).join(' / ')
  const tz = m.time_pricing?.weekdays_only
    ? t('modelPlaza.table.timePricingRowHintWeekdays', { timezone: m.time_pricing?.timezone })
    : t('modelPlaza.table.timePricingRowHint', { timezone: m.time_pricing?.timezone })
  return `${tz} ${periods}`
}
</script>
