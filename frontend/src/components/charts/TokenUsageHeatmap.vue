<template>
  <div class="card" :data-heatmap-view="view">
    <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('dashboard.tokenHeatmap.title') }}</h2>
      <div class="flex items-center gap-2">
        <span class="badge badge-gray">{{ t('dashboard.tokenHeatmap.range', { months }) }}</span>
        <div class="inline-flex rounded-lg border border-gray-200 bg-gray-50 p-0.5 dark:border-dark-700 dark:bg-dark-800" role="group">
          <button
            v-for="option in viewOptions"
            :key="option.value"
            type="button"
            class="rounded-md px-2.5 py-1 text-xs font-medium transition-colors"
            :class="view === option.value
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            :aria-pressed="view === option.value"
            @click="view = option.value"
          >
            {{ option.label }}
          </button>
        </div>
      </div>
    </div>
    <div class="p-6">
      <div v-if="loading" class="flex items-center justify-center py-12">
        <LoadingSpinner size="lg" />
      </div>
      <template v-else>
        <div v-if="view === 'year'" class="overflow-x-auto">
          <div class="min-w-[760px]">
            <div class="flex gap-2">
              <div class="flex w-9 shrink-0 flex-col gap-1 pt-5">
                <span
                  v-for="(label, index) in weekdayLabels"
                  :key="index"
                  class="h-3.5 text-right text-[10px] leading-3.5 text-gray-400 dark:text-dark-500"
                >{{ label }}</span>
              </div>
              <div class="flex flex-1 gap-1">
                <div
                  v-for="column in yearColumns"
                  :key="column.key"
                  class="flex min-w-0 flex-1 flex-col gap-1"
                >
                  <span class="h-4 overflow-visible whitespace-nowrap text-[10px] leading-4 text-gray-400 dark:text-dark-500">{{ column.label }}</span>
                  <span
                    v-for="cell in column.cells"
                    :key="cell.date"
                    class="h-3.5 rounded-[2px] transition-colors"
                    :class="cell.inRange ? LEVEL_CLASSES[cell.level] : 'bg-transparent'"
                    :aria-label="cellTitle(cell)"
                    :data-date="cell.date"
                    :data-level="cell.level"
                    :data-tokens="cell.tokens"
                    :data-future="cell.future || undefined"
                    @mouseenter="showTooltip(cell, $event)"
                    @mouseleave="hideTooltip"
                  ></span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="overflow-x-auto">
          <div class="min-w-[720px]">
            <div class="flex items-center gap-3">
              <span class="w-20 shrink-0"></span>
              <div class="heatmap-grid flex-1">
                <span
                  v-for="day in DAYS_PER_ROW"
                  :key="day"
                  class="text-center text-[10px] leading-4 text-gray-400 dark:text-dark-500"
                >{{ dayTick(day) }}</span>
              </div>
              <span class="w-16 shrink-0"></span>
            </div>
            <div
              v-for="month in monthRows"
              :key="month.key"
              class="mt-1 flex items-center gap-3"
              :data-month="month.key"
              :data-total="month.total"
            >
              <span class="w-20 shrink-0 text-xs font-medium text-gray-600 dark:text-gray-300">{{ month.label }}</span>
              <div class="heatmap-grid flex-1">
                <span
                  v-for="cell in month.cells"
                  :key="cell.date"
                  class="h-5 rounded-[3px] transition-colors"
                  :class="cell.future ? 'bg-transparent' : LEVEL_CLASSES[cell.level]"
                  :aria-label="cellTitle(cell)"
                  :data-date="cell.date"
                  :data-level="cell.level"
                  :data-tokens="cell.tokens"
                  :data-future="cell.future || undefined"
                  @mouseenter="showTooltip(cell, $event)"
                  @mouseleave="hideTooltip"
                ></span>
              </div>
              <span
                class="w-16 shrink-0 text-right text-xs tabular-nums text-gray-500 dark:text-dark-400"
                :title="t('dashboard.tokenHeatmap.monthTotal', { tokens: formatNumberLocaleString(month.total) })"
              >{{ formatTokensK(month.total) }}</span>
            </div>
          </div>
        </div>

        <div class="mt-4 flex flex-wrap items-center justify-between gap-2 text-xs text-gray-500 dark:text-dark-400">
          <span>{{ t('dashboard.tokenHeatmap.hint') }}</span>
          <div class="flex items-center gap-3">
            <span class="tabular-nums">{{ t('dashboard.tokenHeatmap.total', { tokens: formatTokensK(yearTotal) }) }}</span>
            <div class="flex items-center gap-1">
              <span>{{ t('dashboard.tokenHeatmap.less') }}</span>
              <span
                v-for="level in 4"
                :key="level"
                class="h-3 w-3 rounded-[2px]"
                :class="LEVEL_CLASSES[level]"
              ></span>
              <span>{{ t('dashboard.tokenHeatmap.more') }}</span>
            </div>
          </div>
        </div>
        <p v-if="!hasData" class="mt-3 text-center text-sm text-gray-500 dark:text-gray-400">
          {{ t('dashboard.tokenHeatmap.empty') }}
        </p>
      </template>
    </div>

    <Teleport to="body">
      <div
        v-if="tooltip"
        class="pointer-events-none fixed z-[200] w-64 rounded-xl border border-gray-200 bg-white p-3 text-xs shadow-2xl dark:border-dark-700 dark:bg-dark-800"
        :style="{ left: `${tooltip.x}px`, top: `${tooltip.y}px` }"
      >
        <p class="font-semibold text-gray-900 dark:text-white">{{ tooltipTitle }}</p>
        <div class="mt-2 space-y-1">
          <div v-for="row in tooltipRows" :key="row.key" class="flex items-center justify-between gap-4">
            <span class="text-gray-500 dark:text-gray-400">{{ row.label }}</span>
            <span class="tabular-nums text-gray-900 dark:text-white">{{ row.value }}</span>
          </div>
        </div>
        <div class="mt-2 space-y-1 border-t border-gray-100 pt-2 dark:border-dark-700">
          <div class="flex items-center justify-between gap-4">
            <span class="text-gray-500 dark:text-gray-400">{{ t('dashboard.tokenHeatmap.tooltip.total') }}</span>
            <span class="font-medium tabular-nums text-gray-900 dark:text-white">{{ tooltipTotal }}</span>
          </div>
          <div class="flex items-center justify-between gap-4">
            <span class="text-gray-500 dark:text-gray-400">{{ t('dashboard.tokenHeatmap.tooltip.cost') }}</span>
            <span class="font-medium tabular-nums text-red-500 dark:text-red-400">{{ tooltipCost }}</span>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { TrendDataPoint } from '@/types'
import {
  formatCostFixed,
  formatDate,
  formatDateLocalInput,
  formatNumberLocaleString,
  formatTokensK,
} from '@/utils/format'

type HeatmapView = 'year' | 'month'

interface HeatmapCell {
  date: string
  tokens: number
  requests: number
  level: number
  future: boolean
  inRange: boolean
}

const props = withDefaults(defineProps<{
  trend: TrendDataPoint[]
  loading: boolean
  months?: number
  defaultView?: HeatmapView
}>(), {
  months: 12,
  defaultView: 'year',
})

const { t } = useI18n()

const DAYS_PER_ROW = 31
const DAY_TICKS = new Set([1, 5, 10, 15, 20, 25, 31])
const LEVEL_CLASSES = [
  'bg-gray-100 dark:bg-dark-700/60',
  'bg-primary-100 dark:bg-primary-900/60',
  'bg-primary-200 dark:bg-primary-800',
  'bg-primary-300 dark:bg-primary-700',
  'bg-primary-400 dark:bg-primary-600',
]

const view = ref<HeatmapView>(props.defaultView)
const tooltip = ref<{ date: string; x: number; y: number } | null>(null)
const viewOptions = computed(() => [
  { value: 'year' as const, label: t('dashboard.tokenHeatmap.viewYear') },
  { value: 'month' as const, label: t('dashboard.tokenHeatmap.viewMonth') },
])

const dayTick = (day: number) => DAY_TICKS.has(day) ? String(day) : ''

const pointsByDate = computed(() => {
  const map = new Map<string, TrendDataPoint>()
  for (const point of props.trend) {
    if (!point.date) continue
    map.set(point.date.slice(0, 10), point)
  }
  return map
})

const thresholds = computed(() => {
  const values = props.trend
    .map(point => point.total_tokens)
    .filter(tokens => tokens > 0)
    .sort((a, b) => a - b)
  if (values.length === 0) return []
  const pick = (ratio: number) => values[Math.min(values.length - 1, Math.floor(ratio * values.length))]
  return [pick(0.25), pick(0.5), pick(0.75), values[values.length - 1]]
})

const levelFor = (tokens: number) => {
  if (tokens <= 0) return 0
  const [q1, q2, q3, max] = thresholds.value
  if (q1 === undefined) return 4
  if (tokens >= max) return 4
  if (tokens <= q1) return 1
  if (tokens <= q2) return 2
  if (tokens <= q3) return 3
  return 4
}

const buildCell = (date: string, future: boolean, inRange: boolean): HeatmapCell => {
  const point = pointsByDate.value.get(date)
  const tokens = point?.total_tokens ?? 0
  return {
    date,
    tokens,
    requests: point?.requests ?? 0,
    level: future || !inRange ? 0 : levelFor(tokens),
    future,
    inRange,
  }
}

const monthRows = computed(() => {
  const now = new Date()
  const todayKey = formatDateLocalInput(now)
  const rows = []
  for (let offset = 0; offset < props.months; offset++) {
    const first = new Date(now.getFullYear(), now.getMonth() - offset, 1)
    const year = first.getFullYear()
    const monthIndex = first.getMonth()
    const daysInMonth = new Date(year, monthIndex + 1, 0).getDate()
    const cells: HeatmapCell[] = []
    let total = 0
    for (let day = 1; day <= daysInMonth; day++) {
      const date = formatDateLocalInput(new Date(year, monthIndex, day))
      const future = date > todayKey
      const cell = buildCell(date, future, true)
      total += cell.tokens
      cells.push(cell)
    }
    rows.push({
      key: `${year}-${String(monthIndex + 1).padStart(2, '0')}`,
      label: formatDate(first, { year: 'numeric', month: 'short' }),
      cells,
      total,
    })
  }
  return rows
})

const axisWeekdays = computed(() => [0, 2, 4].map((offset) => {
  const monday = new Date(2024, 0, 1 + offset)
  return formatDate(monday, { weekday: 'short' })
}))

const weekdayLabels = computed(() => {
  const labels = Array.from({ length: 7 }, () => '')
  labels[0] = axisWeekdays.value[0] ?? ''
  labels[2] = axisWeekdays.value[1] ?? ''
  labels[4] = axisWeekdays.value[2] ?? ''
  return labels
})

const yearColumns = computed(() => {
  const now = new Date()
  const start = new Date(now.getFullYear(), now.getMonth() - (props.months - 1), 1)
  const end = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const alignedStart = new Date(start)
  alignedStart.setDate(alignedStart.getDate() - ((alignedStart.getDay() + 6) % 7))
  const alignedEnd = new Date(end)
  alignedEnd.setDate(alignedEnd.getDate() + (6 - ((alignedEnd.getDay() + 6) % 7)))
  const startKey = formatDateLocalInput(start)
  const endKey = formatDateLocalInput(end)

  const columns = []
  const cursor = new Date(alignedStart)
  let weekIndex = 0
  while (cursor <= alignedEnd) {
    const cells: HeatmapCell[] = []
    let label = ''
    for (let i = 0; i < 7; i++) {
      const date = new Date(cursor.getFullYear(), cursor.getMonth(), cursor.getDate() + i)
      const key = formatDateLocalInput(date)
      const inRange = key >= startKey && key <= endKey
      if (date.getDate() === 1 && key >= startKey && key <= endKey) {
        label = formatDate(date, { month: 'short' })
      }
      cells.push(buildCell(key, key > endKey, inRange))
    }
    columns.push({ key: formatDateLocalInput(cursor), label, cells })
    cursor.setDate(cursor.getDate() + 7)
    weekIndex += 1
    if (weekIndex > 60) break
  }
  return columns
})

const yearTotal = computed(() => {
  let total = 0
  for (const column of yearColumns.value) {
    for (const cell of column.cells) {
      if (cell.inRange) total += cell.tokens
    }
  }
  return total
})

const hasData = computed(() => props.trend.some(point => point.total_tokens > 0))

const hoveredPoint = computed(() => tooltip.value ? pointsByDate.value.get(tooltip.value.date) : undefined)

const tooltipTitle = computed(() => {
  if (!tooltip.value) return ''
  return formatDate(new Date(`${tooltip.value.date}T00:00:00`), {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'short',
  })
})

const tooltipRows = computed(() => {
  const point = hoveredPoint.value
  const number = (value: number | undefined) => formatNumberLocaleString(value ?? 0)
  return [
    { key: 'requests', label: t('dashboard.tokenHeatmap.tooltip.requests'), value: number(point?.requests) },
    { key: 'input', label: t('dashboard.tokenHeatmap.tooltip.input'), value: number(point?.input_tokens) },
    { key: 'output', label: t('dashboard.tokenHeatmap.tooltip.output'), value: number(point?.output_tokens) },
    { key: 'cache5m', label: t('dashboard.tokenHeatmap.tooltip.cache5m'), value: number(point?.cache_creation_5m_tokens) },
    { key: 'cache1h', label: t('dashboard.tokenHeatmap.tooltip.cache1h'), value: number(point?.cache_creation_1h_tokens) },
    { key: 'cacheRead', label: t('dashboard.tokenHeatmap.tooltip.cacheRead'), value: number(point?.cache_read_tokens) },
  ]
})

const tooltipTotal = computed(() => formatNumberLocaleString(hoveredPoint.value?.total_tokens ?? 0))
const tooltipCost = computed(() => `$${formatCostFixed(hoveredPoint.value?.actual_cost ?? 0)}`)

const showTooltip = (cell: HeatmapCell, event: MouseEvent) => {
  if (!cell.inRange || cell.future) return
  const target = event.currentTarget as HTMLElement | null
  if (!target) return
  const rect = target.getBoundingClientRect()
  const width = 256
  const height = 300
  const gap = 10
  let x = rect.right + gap
  if (x + width > window.innerWidth - 8) x = rect.left - width - gap
  x = Math.max(8, Math.min(x, window.innerWidth - width - 8))
  const y = Math.max(8, Math.min(rect.top - 8, window.innerHeight - height - 8))
  tooltip.value = { date: cell.date, x, y }
}

const hideTooltip = () => {
  tooltip.value = null
}

watch(view, hideTooltip)

const cellTitle = (cell: HeatmapCell) => t('dashboard.tokenHeatmap.cellTooltip', {
  date: cell.date,
  tokens: formatNumberLocaleString(cell.tokens),
  requests: formatNumberLocaleString(cell.requests),
})
</script>

<style scoped>
.heatmap-grid {
  display: grid;
  grid-template-columns: repeat(31, minmax(0, 1fr));
  gap: 4px;
}
</style>
