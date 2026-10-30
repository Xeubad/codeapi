<template>
  <div class="card flex h-full flex-col p-5">
    <div class="flex items-start justify-between gap-3">
      <div>
        <h2 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.dashboard.topSpenders.title') }}</h2>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.topSpenders.subtitle') }}</p>
      </div>
      <RouterLink to="/admin/usage" class="text-xs font-medium text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
        {{ t('admin.dashboard.topSpenders.viewAll') }} →
      </RouterLink>
    </div>

    <div v-if="loading" class="flex flex-1 items-center justify-center py-10">
      <LoadingSpinner size="md" />
    </div>
    <div v-else-if="error" class="flex flex-1 items-center justify-center py-10 text-sm text-gray-500 dark:text-gray-400">
      {{ t('admin.dashboard.topSpenders.unavailable') }}
    </div>
    <div v-else-if="items.length === 0" class="flex flex-1 items-center justify-center py-10 text-sm text-gray-500 dark:text-gray-400">
      {{ t('admin.dashboard.topSpenders.empty') }}
    </div>
    <template v-else>
      <div class="mt-4 flex-1 overflow-x-auto">
        <table class="w-full text-xs">
          <thead>
            <tr class="text-gray-500 dark:text-gray-400">
              <th class="pb-2 text-left">#</th>
              <th class="pb-2 text-left">{{ t('admin.dashboard.topSpenders.user') }}</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.topSpenders.requests') }}</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.topSpenders.tokens') }}</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.topSpenders.spend') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(item, index) in items.slice(0, 5)"
              :key="item.user_id"
              class="cursor-pointer border-t border-gray-100 transition-colors hover:bg-gray-50 dark:border-dark-700 dark:hover:bg-dark-700/40"
              :data-user-id="item.user_id"
              @click="emit('select', item)"
            >
              <td class="py-2 font-semibold tabular-nums text-gray-400 dark:text-dark-500">{{ index + 1 }}</td>
              <td class="max-w-[220px] truncate py-2 font-medium text-gray-900 dark:text-white" :title="userLabel(item)">
                {{ userLabel(item) }}
              </td>
              <td class="py-2 text-right tabular-nums text-gray-600 dark:text-gray-400">{{ formatNumberLocaleString(item.requests) }}</td>
              <td class="py-2 text-right tabular-nums text-gray-600 dark:text-gray-400">{{ formatTokensK(item.tokens) }}</td>
              <td class="py-2 text-right tabular-nums text-purple-600 dark:text-purple-400">${{ formatCostFixed(item.actual_cost) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="mt-auto flex flex-wrap items-center justify-between gap-2 border-t border-gray-100 pt-3 text-[11px] text-gray-500 dark:border-dark-700 dark:text-gray-400">
        <span>{{ t('admin.dashboard.topSpenders.total') }}</span>
        <div class="flex flex-wrap items-center gap-3">
          <span>{{ t('admin.dashboard.topSpenders.requests') }} {{ formatNumberLocaleString(totalRequests ?? 0) }}</span>
          <span>{{ t('admin.dashboard.topSpenders.tokens') }} {{ formatTokensK(totalTokens ?? 0) }}</span>
          <span class="font-medium text-purple-600 dark:text-purple-400">${{ formatCostFixed(totalActualCost ?? 0) }}</span>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { UserSpendingRankingItem } from '@/types'
import { formatCostFixed, formatNumberLocaleString, formatTokensK } from '@/utils/format'

defineProps<{
  items: UserSpendingRankingItem[]
  loading: boolean
  error: boolean
  totalActualCost?: number
  totalRequests?: number
  totalTokens?: number
}>()

const emit = defineEmits<{ select: [item: UserSpendingRankingItem] }>()
const { t } = useI18n()

const userLabel = (item: UserSpendingRankingItem) => item.username?.trim() || item.email?.trim() || `#${item.user_id}`
</script>
