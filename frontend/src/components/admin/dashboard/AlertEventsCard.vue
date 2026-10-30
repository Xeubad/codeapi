<template>
  <div class="card flex h-full flex-col p-5">
    <div class="flex items-start justify-between gap-3">
      <div>
        <h2 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.dashboard.alerts.title') }}</h2>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.alerts.subtitle') }}</p>
      </div>
      <RouterLink to="/admin/ops" class="text-xs font-medium text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
        {{ t('admin.dashboard.alerts.viewAll') }} →
      </RouterLink>
    </div>

    <div v-if="loading" class="flex flex-1 items-center justify-center py-10">
      <LoadingSpinner size="md" />
    </div>
    <div v-else-if="error" class="flex flex-1 items-center justify-center py-10 text-sm text-gray-500 dark:text-gray-400">
      {{ t('admin.dashboard.alerts.unavailable') }}
    </div>
    <div v-else-if="events.length === 0" class="flex flex-1 flex-col items-center justify-center py-10 text-center">
      <div class="mb-2 flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-100 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400">
        <Icon name="check" size="md" />
      </div>
      <p class="text-sm text-gray-600 dark:text-gray-300">{{ t('admin.dashboard.alerts.empty') }}</p>
    </div>
    <ul v-else class="mt-4 flex flex-1 flex-col justify-center gap-2">
      <li
        v-for="event in events"
        :key="event.id"
        class="rounded-lg border border-gray-100 p-3 dark:border-dark-700"
      >
        <div class="flex items-start gap-2">
          <span class="mt-1.5 h-2 w-2 shrink-0 rounded-full" :class="severityDot(event.severity)"></span>
          <div class="min-w-0 flex-1">
            <p class="truncate text-xs font-medium text-gray-900 dark:text-white" :title="event.title || event.description">
              {{ event.title || event.description || `#${event.rule_id}` }}
            </p>
            <div class="mt-1 flex flex-wrap items-center gap-2 text-[11px] text-gray-500 dark:text-gray-400">
              <span class="rounded px-1.5 py-0.5 font-medium uppercase" :class="severityBadge(event.severity)">{{ event.severity }}</span>
              <span>{{ formatRelativeTime(event.fired_at) }}</span>
            </div>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import Icon from '@/components/icons/Icon.vue'
import type { AlertEvent } from '@/api/admin/ops'
import { formatRelativeTime } from '@/utils/format'

defineProps<{
  events: AlertEvent[]
  loading: boolean
  error: boolean
}>()

const { t } = useI18n()

const severityDot = (severity: string) => {
  if (severity === 'critical') return 'bg-red-500'
  if (severity === 'warning') return 'bg-amber-500'
  if (severity === 'info') return 'bg-blue-500'
  return 'bg-gray-400'
}

const severityBadge = (severity: string) => {
  if (severity === 'critical') return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
  if (severity === 'warning') return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
  if (severity === 'info') return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
  return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-gray-300'
}
</script>
