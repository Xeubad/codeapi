<template>
  <BaseDialog
    :show="show"
    :title="model?.name ?? ''"
    width="wide"
    close-on-click-outside
    @close="emit('close')"
  >
    <div v-if="model" class="space-y-5">
      <!-- 头部徽章:平台/计费模式/思考倍率 -->
      <div class="flex flex-wrap items-center gap-1.5">
        <span
          v-if="model.platform"
          :class="[
            'inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium',
            platformBadgeLightClass(model.platform)
          ]"
        >
          {{ platformLabel(model.platform) }}
        </span>
        <span
          v-if="billingMode(model) !== BILLING_MODE_TOKEN"
          class="rounded-md bg-gray-100 px-2 py-0.5 text-xs font-medium text-gray-500 dark:bg-dark-700/70 dark:text-dark-300"
        >
          {{ billingModeLabel(model) }}
        </span>
        <span
          v-for="([effort, multiplier]) in reasoningEffortMultipliers(model)"
          :key="effort"
          class="rounded-md bg-amber-50 px-1.5 py-0.5 text-[10px] font-medium text-amber-700 dark:bg-amber-900/20 dark:text-amber-300"
          :title="t('modelPlaza.table.reasoningMultiplierHint', { effort, multiplier })"
        >
          {{ t('modelPlaza.table.reasoningMultiplierBadge', { effort, multiplier }) }}
        </span>
        <span
          v-if="model.long_context_basis === 'marginal'"
          class="rounded-md bg-gray-100 px-1.5 py-0.5 text-[10px] font-medium text-gray-500 dark:bg-dark-700/70 dark:text-dark-300"
          :title="t('modelPlaza.table.tierHintMarginal')"
        >
          {{ t('modelPlaza.table.marginalBadge') }}
        </span>
      </div>

      <!-- 本站实收价 -->
      <section>
        <h4 class="mb-2 flex items-center gap-1.5 text-sm font-semibold text-gray-900 dark:text-white">
          {{ t('modelPlaza.dialog.paidTitle') }}
          <span class="rounded-md bg-primary-50 px-1.5 py-0.5 font-mono text-[11px] font-semibold text-primary-600 dark:bg-primary-500/10 dark:text-primary-300">
            ×{{ rateDisplay }}
          </span>
          <span class="font-normal text-gray-400 dark:text-dark-500">{{ t('modelPlaza.table.unitPerMillion') }}</span>
        </h4>
        <div
          v-if="model.pricing"
          class="overflow-x-auto rounded-xl border border-gray-200 dark:border-dark-700/60"
        >
          <table class="w-full text-sm tabular-nums">
            <thead>
              <tr class="bg-gray-50 text-left text-xs font-medium uppercase tracking-wide text-gray-500 dark:bg-dark-800/70 dark:text-dark-400">
                <th class="px-3 py-2">{{ t('modelPlaza.dialog.tier') }}</th>
                <th class="px-3 py-2">{{ t('modelPlaza.table.input') }}</th>
                <th class="px-3 py-2">{{ t('modelPlaza.table.output') }}</th>
                <th class="px-3 py-2">{{ t('modelPlaza.table.cacheWrite') }}</th>
                <th class="px-3 py-2">{{ t('modelPlaza.table.cacheRead') }}</th>
              </tr>
            </thead>
            <tbody>
              <!-- token 计费:无阶梯一行,有阶梯每档一行 -->
              <template v-if="billingMode(model) === BILLING_MODE_TOKEN">
                <tr
                  v-for="(row, idx) in paidRows"
                  :key="`paid-${idx}`"
                  class="border-t border-gray-100 font-mono dark:border-dark-700/60"
                >
                  <td class="px-3 py-2 font-sans text-xs text-gray-500 dark:text-dark-400" :title="tierHint">
                    {{ row.label || t('modelPlaza.dialog.baseTier') }}
                  </td>
                  <td class="px-3 py-2 text-gray-900 dark:text-gray-100">{{ paidPerMillion(row.input, row.period) }}</td>
                  <td class="px-3 py-2 text-gray-900 dark:text-gray-100">{{ paidPerMillion(row.output, row.period) }}</td>
                  <td class="px-3 py-2 text-gray-700 dark:text-gray-300">
                    <template v-if="row.cacheWrite != null">
                      {{ paidPerMillion(row.cacheWrite, row.period) }}
                      <span v-if="row.cacheWrite1h != null" class="text-xs text-gray-400 dark:text-dark-500">
                        (1h {{ paidPerMillion(row.cacheWrite1h, row.period) }})
                      </span>
                    </template>
                    <span v-else class="text-gray-400 dark:text-dark-500">-</span>
                  </td>
                  <td class="px-3 py-2 text-gray-700 dark:text-gray-300">{{ row.cacheRead != null ? paidPerMillion(row.cacheRead, row.period) : '-' }}</td>
                </tr>
              </template>
              <!-- 按图/按次:合并为档位价列表 -->
              <template v-else>
                <tr
                  v-for="(iv, idx) in requestRows"
                  :key="`paid-req-${idx}`"
                  class="border-t border-gray-100 dark:border-dark-700/60"
                >
                  <td class="px-3 py-2 font-sans text-xs text-gray-500 dark:text-dark-400">{{ tierLabel(iv) }}</td>
                  <td colspan="4" class="px-3 py-2 font-mono text-gray-900 dark:text-gray-100">
                    {{ paidRequestPrice(model, iv.per_request_price) }}
                    <span class="text-xs text-gray-400 dark:text-dark-500">{{ perUnitSuffix }}</span>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
        <p v-else class="text-sm text-gray-400 dark:text-dark-500">{{ t('modelPlaza.detail.noPricing') }}</p>

        <!-- 分时倍率时段 -->
        <div
          v-if="timePeriods.length"
          class="mt-2.5 rounded-xl border border-amber-200 bg-amber-50/60 px-3.5 py-2.5 text-xs text-amber-700 dark:border-amber-500/20 dark:bg-amber-500/5 dark:text-amber-300"
        >
          <p class="font-medium">{{ t('modelPlaza.dialog.timePricingTitle') }}</p>
          <p
            v-for="(p, idx) in timePeriods"
            :key="idx"
            class="mt-1 font-mono"
            :title="timePricingRowHint(model)"
          >
            {{ formatTimeWindow(p) }} ×{{ p.multiplier }} →
            {{ t('modelPlaza.dialog.effectiveRate') }} ×{{ periodRate(p) }}
          </p>
          <p v-if="peakWindow" class="mt-1 font-sans">{{ peakNote }}</p>
        </div>
      </section>

      <!-- 官方参考价 -->
      <section>
        <h4 class="mb-2 flex items-center gap-1.5 text-sm font-semibold text-gray-900 dark:text-white">
          {{ t('modelPlaza.dialog.officialTitle') }}
          <span class="font-normal text-gray-400 dark:text-dark-500">{{ t('modelPlaza.table.unitPerMillion') }}</span>
        </h4>
        <div
          v-if="model.official_pricing"
          class="overflow-x-auto rounded-xl border border-gray-200 dark:border-dark-700/60"
        >
          <table class="w-full text-sm tabular-nums">
            <thead>
              <tr class="bg-gray-50 text-left text-xs font-medium uppercase tracking-wide text-gray-500 dark:bg-dark-800/70 dark:text-dark-400">
                <th class="px-3 py-2">{{ t('modelPlaza.dialog.tier') }}</th>
                <th class="px-3 py-2">{{ t('modelPlaza.table.input') }}</th>
                <th class="px-3 py-2">{{ t('modelPlaza.table.output') }}</th>
                <th class="px-3 py-2">{{ t('modelPlaza.table.cacheWrite') }}</th>
                <th class="px-3 py-2">{{ t('modelPlaza.table.cacheRead') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(row, idx) in officialRows"
                :key="`official-${idx}`"
                class="border-t border-gray-100 font-mono text-gray-600 dark:border-dark-700/60 dark:text-dark-300"
              >
                <td class="px-3 py-2 font-sans text-xs text-gray-400 dark:text-dark-500" :title="tierHint">
                  {{ row.label || t('modelPlaza.dialog.baseTier') }}
                </td>
                <td class="px-3 py-2">{{ official(row.input) }}</td>
                <td class="px-3 py-2">{{ official(row.output) }}</td>
                <td class="px-3 py-2">
                  <template v-if="row.cacheWrite != null">
                    {{ official(row.cacheWrite) }}
                    <span v-if="row.cacheWrite1h != null" class="text-xs text-gray-400 dark:text-dark-500">
                      (1h {{ official(row.cacheWrite1h) }})
                    </span>
                  </template>
                  <span v-else class="text-gray-400 dark:text-dark-500">-</span>
                </td>
                <td class="px-3 py-2">{{ row.cacheRead != null ? official(row.cacheRead) : '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-else class="text-sm text-gray-400 dark:text-dark-500">{{ t('modelPlaza.dialog.noOfficial') }}</p>
      </section>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { platformBadgeLightClass, platformLabel } from '@/utils/platformColors'
import { BILLING_MODE_TOKEN, BILLING_MODE_IMAGE, BILLING_MODE_VIDEO } from '@/constants/channel'
import type { PlazaModel, PlazaTimePricingPeriod } from '@/api/modelPlaza'
import {
  usePlazaPricing,
  plazaBillingMode,
  tokenIntervals,
  officialIntervals,
  requestIntervals,
  tierLabel,
  formatTimeWindow,
  reasoningEffortMultipliers
} from './plazaPricing'

const props = defineProps<{
  show: boolean
  model: PlazaModel | null
  rateMultiplier: number
  userRateMultiplier?: number | null
  imageRateIndependent?: boolean
  imageRateMultiplier?: number | null
  videoRateIndependent?: boolean
  videoRateMultiplier?: number | null
  peakWindow?: string
  peakRateMultiplier?: number | null
}>()

const emit = defineEmits<{ (e: 'close'): void }>()

const { t } = useI18n()
const {
  effectiveRate,
  billingMode,
  periodRate,
  paidPerMillion,
  paidRequestPrice,
  official
} = usePlazaPricing(props)

interface PricingRow {
  label: string
  period: PlazaTimePricingPeriod | null
  input: number | null
  output: number | null
  cacheWrite: number | null
  cacheWrite1h: number | null
  cacheRead: number | null
}

/** 实收区行:基础档/阶梯档 + 分时时段展开。 */
const paidRows = computed<PricingRow[]>(() => {
  const m = props.model
  if (!m || !m.pricing || billingMode(m) !== BILLING_MODE_TOKEN) return []
  const rows: PricingRow[] = []
  const ivs = tokenIntervals(m)
  if (ivs.length) {
    for (const iv of ivs) {
      rows.push({
        label: tierLabel(iv),
        period: null,
        input: iv.input_price,
        output: iv.output_price,
        cacheWrite: iv.cache_write_price ?? null,
        cacheWrite1h: iv.cache_write_1h_price ?? null,
        cacheRead: iv.cache_read_price ?? null
      })
    }
  } else {
    rows.push({
      label: '',
      period: null,
      input: m.pricing.input_price ?? null,
      output: m.pricing.output_price ?? null,
      cacheWrite: m.pricing.cache_write_price ?? null,
      cacheWrite1h: m.pricing.cache_write_1h_price ?? null,
      cacheRead: m.pricing.cache_read_price ?? null
    })
  }
  for (const p of m.time_pricing?.periods ?? []) {
    rows.push({
      label: formatTimeWindow(p),
      period: p,
      input: m.pricing.input_price ?? null,
      output: m.pricing.output_price ?? null,
      cacheWrite: m.pricing.cache_write_price ?? null,
      cacheWrite1h: m.pricing.cache_write_1h_price ?? null,
      cacheRead: m.pricing.cache_read_price ?? null
    })
  }
  return rows
})

const requestRows = computed(() => (props.model ? requestIntervals(props.model) : []))

/** 官方区行:官方阶梯或单档。 */
const officialRows = computed<PricingRow[]>(() => {
  const m = props.model
  if (!m) return []
  const o = m.official_pricing
  if (!o) return []
  const ivs = officialIntervals(m)
  if (ivs.length) {
    return ivs.map(iv => ({
      label: tierLabel(iv),
      period: null,
      input: iv.input_price ?? null,
      output: iv.output_price ?? null,
      cacheWrite: iv.cache_write_price ?? null,
      cacheWrite1h: iv.cache_write_1h_price ?? null,
      cacheRead: iv.cache_read_price ?? null
    }))
  }
  return [{
    label: '',
    period: null,
    input: o.input_price,
    output: o.output_price,
    cacheWrite: o.cache_write_price,
    cacheWrite1h: o.cache_write_1h_price ?? null,
    cacheRead: o.cache_read_price
  }]
})

const timePeriods = computed(() => props.model?.time_pricing?.periods ?? [])

const rateDisplay = computed(() => {
  const m = props.model
  if (m && billingMode(m) !== BILLING_MODE_TOKEN) {
    const rate = plazaBillingMode(m) === BILLING_MODE_VIDEO && props.videoRateIndependent === true
      ? (props.videoRateMultiplier ?? 1)
      : plazaBillingMode(m) === BILLING_MODE_IMAGE && props.imageRateIndependent === true
        ? (props.imageRateMultiplier ?? 1)
        : effectiveRate.value
    return String(Math.round(rate * 1000) / 1000)
  }
  return String(Math.round(effectiveRate.value * 1000) / 1000)
})

const tierHint = computed(() =>
  props.model?.long_context_basis === 'marginal'
    ? t('modelPlaza.table.tierHintMarginal')
    : t('modelPlaza.table.tierHint')
)

const peakNote = computed(() =>
  props.peakWindow
    ? t('modelPlaza.table.timePricingRowHintPeak', {
        window: props.peakWindow,
        multiplier: props.peakRateMultiplier ?? 1
      })
    : ''
)

function billingModeLabel(m: PlazaModel): string {
  return plazaBillingMode(m) === BILLING_MODE_IMAGE
    ? t('modelPlaza.table.perImage')
    : t('modelPlaza.table.perRequest')
}

function perUnitSuffixFn(m: PlazaModel): string {
  return plazaBillingMode(m) === BILLING_MODE_IMAGE
    ? t('modelPlaza.table.perUnitImage')
    : t('modelPlaza.table.perUnitRequest')
}
const perUnitSuffix = computed(() => (props.model ? perUnitSuffixFn(props.model) : ''))

function timePricingRowHint(m: PlazaModel): string {
  const key = m.time_pricing?.weekdays_only
    ? 'modelPlaza.table.timePricingRowHintWeekdays'
    : 'modelPlaza.table.timePricingRowHint'
  return t(key, { timezone: m.time_pricing?.timezone })
}
</script>
