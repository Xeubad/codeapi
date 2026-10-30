<template>
  <div class="card flex h-full flex-col p-5">
    <div class="flex items-start justify-between gap-3">
      <div>
        <h2 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.dashboard.health.title') }}</h2>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.health.subtitle') }}</p>
      </div>
      <RouterLink to="/admin/ops" class="text-xs font-medium text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
        {{ t('admin.dashboard.health.viewOps') }} →
      </RouterLink>
    </div>

    <div v-if="loading" class="flex flex-1 items-center justify-center py-10">
      <LoadingSpinner size="md" />
    </div>
    <div v-else class="mt-4 flex flex-1 flex-col justify-between gap-4">
      <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800/50">
          <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.health.healthScore') }}</p>
          <p class="mt-1 text-xl font-bold tabular-nums" :class="healthScoreClass">{{ overview?.health_score ?? '--' }}</p>
        </div>
        <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800/50">
          <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.health.errorRate') }}</p>
          <p class="mt-1 text-xl font-bold tabular-nums" :class="errorRateClass">{{ errorRateText }}</p>
        </div>
        <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800/50">
          <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.health.p95') }}</p>
          <p class="mt-1 text-xl font-bold tabular-nums text-gray-900 dark:text-white">{{ p95Text }}</p>
        </div>
        <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800/50">
          <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.health.qps') }}</p>
          <p class="mt-1 text-xl font-bold tabular-nums text-gray-900 dark:text-white">{{ qpsText }}</p>
        </div>
      </div>

      <div>
        <div class="mb-2 flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
          <span>{{ t('admin.dashboard.health.accounts') }}</span>
          <span>{{ stats?.total_accounts ?? 0 }}</span>
        </div>
        <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
          <div class="flex items-center justify-between rounded-lg border border-gray-100 px-3 py-2 text-xs dark:border-dark-700">
            <span class="text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.health.normal') }}</span>
            <span class="font-semibold tabular-nums text-emerald-600 dark:text-emerald-400">{{ stats?.normal_accounts ?? 0 }}</span>
          </div>
          <div class="flex items-center justify-between rounded-lg border border-gray-100 px-3 py-2 text-xs dark:border-dark-700">
            <span class="text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.health.error') }}</span>
            <span class="font-semibold tabular-nums text-red-500 dark:text-red-400">{{ stats?.error_accounts ?? 0 }}</span>
          </div>
          <div class="flex items-center justify-between rounded-lg border border-gray-100 px-3 py-2 text-xs dark:border-dark-700">
            <span class="text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.health.rateLimited') }}</span>
            <span class="font-semibold tabular-nums text-amber-600 dark:text-amber-400">{{ stats?.ratelimit_accounts ?? 0 }}</span>
          </div>
          <div class="flex items-center justify-between rounded-lg border border-gray-100 px-3 py-2 text-xs dark:border-dark-700">
            <span class="text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.health.overloaded') }}</span>
            <span class="font-semibold tabular-nums text-orange-600 dark:text-orange-400">{{ stats?.overload_accounts ?? 0 }}</span>
          </div>
        </div>
      </div>

      <p v-if="error" class="text-xs text-amber-600 dark:text-amber-400">{{ t('admin.dashboard.health.unavailable') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { DashboardStats } from '@/types'
import type { OpsDashboardOverview } from '@/api/admin/ops'

const props = defineProps<{
  stats: DashboardStats | null
  overview: OpsDashboardOverview | null
  loading: boolean
  error: boolean
}>()

const { t } = useI18n()

const errorRateText = computed(() => {
  if (!props.overview) return '--'
  return `${((props.overview.error_rate ?? 0) * 100).toFixed(2)}%`
})

const p95Text = computed(() => {
  const value = props.overview?.duration?.p95_ms
  if (value === null || value === undefined) return '--'
  return value >= 1000 ? `${(value / 1000).toFixed(2)}s` : `${Math.round(value)}ms`
})

const qpsText = computed(() => {
  const value = props.overview?.qps?.current
  return value === undefined || value === null ? '--' : value.toFixed(2)
})

const healthScoreClass = computed(() => {
  const score = props.overview?.health_score
  if (score === undefined || score === null) return 'text-gray-900 dark:text-white'
  if (score >= 90) return 'text-emerald-600 dark:text-emerald-400'
  if (score >= 70) return 'text-amber-600 dark:text-amber-400'
  return 'text-red-500 dark:text-red-400'
})

const errorRateClass = computed(() => {
  const rate = (props.overview?.error_rate ?? 0) * 100
  if (rate >= 5) return 'text-red-500 dark:text-red-400'
  if (rate >= 1) return 'text-amber-600 dark:text-amber-400'
  return 'text-gray-900 dark:text-white'
})
</script>
