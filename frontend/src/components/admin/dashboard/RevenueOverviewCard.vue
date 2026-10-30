<template>
  <div class="card flex h-full flex-col p-5">
    <div class="flex items-start justify-between gap-3">
      <div>
        <h2 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.dashboard.revenue.title') }}</h2>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.revenue.subtitle') }}</p>
      </div>
      <RouterLink to="/admin/orders/dashboard" class="text-xs font-medium text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
        {{ t('admin.dashboard.revenue.viewDetails') }} →
      </RouterLink>
    </div>

    <div v-if="loading" class="flex flex-1 items-center justify-center py-10">
      <LoadingSpinner size="md" />
    </div>
    <div v-else-if="error" class="flex flex-1 items-center justify-center py-10 text-sm text-gray-500 dark:text-gray-400">
      {{ t('admin.dashboard.revenue.unavailable') }}
    </div>
    <div v-else class="mt-4 flex flex-1 flex-col gap-3">
      <div class="grid grid-cols-2 gap-3">
        <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800/50">
          <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.revenue.todayRevenue') }}</p>
          <p class="mt-1 truncate text-lg font-bold tabular-nums text-emerald-600 dark:text-emerald-400" :title="amountsText(stats?.today_amount)">
            {{ amountsText(stats?.today_amount) }}
          </p>
          <p class="mt-0.5 text-[11px] text-gray-500 dark:text-gray-400">
            {{ t('admin.dashboard.revenue.orders', { count: stats?.today_count ?? 0 }) }}
          </p>
        </div>
        <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-800/50">
          <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.revenue.totalRevenue') }}</p>
          <p class="mt-1 truncate text-lg font-bold tabular-nums text-gray-900 dark:text-white" :title="amountsText(stats?.total_amount)">
            {{ amountsText(stats?.total_amount) }}
          </p>
          <p class="mt-0.5 text-[11px] text-gray-500 dark:text-gray-400">
            {{ t('admin.dashboard.revenue.orders', { count: stats?.total_count ?? 0 }) }}
          </p>
        </div>
      </div>
      <div class="mt-auto grid grid-cols-2 gap-3 border-t border-gray-100 pt-3 dark:border-dark-700">
        <div>
          <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.revenue.todayCost') }}</p>
          <p class="mt-0.5 text-sm font-semibold tabular-nums text-purple-600 dark:text-purple-400">${{ formatCost(todayCost) }}</p>
        </div>
        <div>
          <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.revenue.totalCost') }}</p>
          <p class="mt-0.5 text-sm font-semibold tabular-nums text-gray-900 dark:text-white">${{ formatCost(totalCost) }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { CurrencyAmounts, DashboardStats as PaymentDashboardStats } from '@/types/payment'

defineProps<{
  stats: PaymentDashboardStats | null
  loading: boolean
  error: boolean
  todayCost: number
  totalCost: number
}>()

const { t } = useI18n()

function amountsText(amounts?: CurrencyAmounts | null): string {
  if (!amounts) return '—'
  const entries = Object.entries(amounts).filter(([, value]) => Number.isFinite(value))
  if (entries.length === 0) return '—'
  return entries
    .sort(([left], [right]) => left.localeCompare(right))
    .map(([currency, value]) => new Intl.NumberFormat(undefined, { style: 'currency', currency }).format(value))
    .join(' · ')
}

function formatCost(value: number): string {
  if (!Number.isFinite(value)) return '0.0000'
  return value.toFixed(4)
}
</script>
