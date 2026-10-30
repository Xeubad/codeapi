<template>
  <div class="card p-5">
    <div class="flex flex-col gap-5 xl:flex-row xl:items-center">
      <div class="flex min-w-0 items-center gap-4 xl:w-64 xl:shrink-0">
        <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-500 to-primary-600 text-white shadow-lg shadow-indigo-500/20">
          <Icon name="database" size="md" />
        </div>
        <div class="min-w-0">
          <div class="flex flex-wrap items-center gap-2">
            <h2 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('dashboard.cumulative.title') }}</h2>
            <span class="badge badge-gray">{{ t('dashboard.cumulative.since') }}</span>
          </div>
          <p class="mt-1 text-2xl font-bold tabular-nums text-gray-900 dark:text-white">{{ formatTokensK(stats.total_tokens) }}</p>
          <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('dashboard.cumulative.today', { value: formatTokensK(stats.today_tokens) }) }}</p>
        </div>
      </div>

      <div class="min-w-0 flex-1 xl:px-2">
        <div class="flex items-center justify-between text-[11px] text-gray-500 dark:text-gray-400">
          <span>{{ t('dashboard.cumulative.composition') }}</span>
          <span>{{ formatNumberLocaleString(stats.total_tokens) }} tokens</span>
        </div>
        <div class="mt-2 flex h-2 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
          <div
            v-for="segment in segments"
            :key="segment.key"
            class="h-full"
            :class="segment.color"
            :style="{ width: `${segment.percent}%` }"
            :title="`${segment.label}: ${segment.percent}%`"
          ></div>
        </div>
        <div class="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-[11px] text-gray-500 dark:text-gray-400">
          <span v-for="segment in segments" :key="segment.key" class="inline-flex items-center gap-1.5">
            <span class="h-2 w-2 rounded-full" :class="segment.color"></span>
            {{ segment.label }} {{ segment.percent }}%
          </span>
        </div>
      </div>

      <div class="grid grid-cols-2 gap-6 xl:w-72 xl:shrink-0">
        <div>
          <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('dashboard.cumulative.requests') }}</p>
          <p class="mt-1 text-lg font-bold tabular-nums text-gray-900 dark:text-white">{{ formatNumberLocaleString(stats.total_requests) }}</p>
          <p class="mt-0.5 text-[11px] text-gray-500 dark:text-gray-400">{{ t('dashboard.cumulative.today', { value: formatNumberLocaleString(stats.today_requests) }) }}</p>
        </div>
        <div>
          <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('dashboard.cumulative.spend') }}</p>
          <p class="mt-1 text-lg font-bold tabular-nums text-purple-600 dark:text-purple-400">${{ formatCostFixed(stats.total_actual_cost) }}</p>
          <p class="mt-0.5 text-[11px] text-gray-500 dark:text-gray-400">{{ t('dashboard.cumulative.standard', { value: `$${formatCostFixed(stats.total_cost)}` }) }}</p>
        </div>
      </div>

      <RouterLink to="/usage" class="btn btn-secondary shrink-0 self-start xl:self-auto">
        {{ t('dashboard.cumulative.viewDetails') }}
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import type { UserDashboardStats } from '@/api/usage'
import { formatCostFixed, formatNumberLocaleString, formatTokensK } from '@/utils/format'

const props = defineProps<{ stats: UserDashboardStats }>()
const { t } = useI18n()

const segments = computed(() => {
  const items = [
    { key: 'input', label: t('dashboard.cumulative.input'), value: props.stats.total_input_tokens, color: 'bg-blue-500' },
    { key: 'output', label: t('dashboard.cumulative.output'), value: props.stats.total_output_tokens, color: 'bg-emerald-500' },
    { key: 'cacheCreation', label: t('dashboard.cumulative.cacheCreation'), value: props.stats.total_cache_creation_tokens, color: 'bg-amber-500' },
    { key: 'cacheRead', label: t('dashboard.cumulative.cacheRead'), value: props.stats.total_cache_read_tokens, color: 'bg-cyan-500' },
  ]
  const componentTotal = items.reduce((sum, item) => sum + item.value, 0)
  return items.map(item => ({
    ...item,
    percent: componentTotal > 0 ? Math.round((item.value / componentTotal) * 1000) / 10 : 0,
  }))
})
</script>
