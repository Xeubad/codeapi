<template>
  <section class="rounded-2xl border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-800" aria-labelledby="resource-analysis-title">
    <div class="flex flex-wrap items-start justify-between gap-3">
      <div>
        <h2 id="resource-analysis-title" class="text-base font-semibold text-gray-900 dark:text-white">{{ t('admin.ops.resourceAnalysis.title') }}</h2>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.ops.resourceAnalysis.description') }}</p>
      </div>
      <span class="rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600 dark:bg-dark-700 dark:text-gray-300">{{ t('admin.ops.resourceAnalysis.observation') }}</span>
    </div>

    <p v-if="!metrics" class="mt-4 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.ops.resourceAnalysis.noData') }}</p>
    <template v-else>
      <p class="mt-4 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.ops.resourceAnalysis.sampleTime', { time: sampleTime }) }}</p>
      <p v-if="stale" class="mt-2 text-sm text-amber-700 dark:text-amber-400">{{ t('admin.ops.resourceAnalysis.stale') }}</p>
      <div class="mt-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        <div v-for="item in readings" :key="item.label" class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
          <div class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ item.label }}</div>
          <div class="mt-2 text-lg font-semibold tabular-nums text-gray-900 dark:text-white">{{ item.value }}</div>
          <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ item.detail }}</div>
        </div>
      </div>
      <ul class="mt-4 space-y-2 text-sm text-gray-600 dark:text-gray-300">
        <li v-if="metrics.db_ok === false" class="text-red-600 dark:text-red-400">{{ t('admin.ops.resourceAnalysis.dbUnavailable') }}</li>
        <li v-if="metrics.redis_ok === false" class="text-red-600 dark:text-red-400">{{ t('admin.ops.resourceAnalysis.redisUnavailable') }}</li>
        <li>{{ t('admin.ops.resourceAnalysis.redisScope') }}</li>
        <li>{{ t('admin.ops.resourceAnalysis.limits') }}</li>
      </ul>
    </template>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { OpsSystemMetricsSnapshot } from '@/api/admin/ops'

const props = defineProps<{ metrics: OpsSystemMetricsSnapshot | null | undefined }>()
const { t } = useI18n()
const sampleTime = computed(() => {
  const date = props.metrics?.created_at ? new Date(props.metrics.created_at) : null
  return date && Number.isFinite(date.getTime()) ? date.toLocaleString() : t('admin.ops.noData')
})
const stale = computed(() => {
  const time = props.metrics?.created_at ? Date.parse(props.metrics.created_at) : NaN
  return !Number.isFinite(time) || Date.now() - time > 10 * 60 * 1000
})
const percent = (value: number | null | undefined) => typeof value === 'number' && Number.isFinite(value) ? `${value.toFixed(1)}%` : '—'
const readings = computed(() => {
  const m = props.metrics
  const db = m?.db_conn_active != null && m.db_conn_idle != null ? m.db_conn_active + m.db_conn_idle : null
  const count = (value: number | null | undefined) => value == null ? '—' : String(value)
  return [
    { label: 'CPU', value: percent(m?.cpu_usage_percent), detail: t('admin.ops.resourceAnalysis.cpuScope') },
    { label: t('admin.ops.memory'), value: percent(m?.memory_usage_percent), detail: m?.memory_used_mb != null && m.memory_total_mb != null ? `${m.memory_used_mb} / ${m.memory_total_mb} MiB` : t('admin.ops.noData') },
    { label: t('admin.ops.db'), value: `${count(db)} / ${count(m?.db_max_open_conns)}`, detail: t('admin.ops.resourceAnalysis.dbScope') },
    { label: 'Redis', value: `${count(m?.redis_conn_total)} / ${count(m?.redis_pool_size)}`, detail: t('admin.ops.resourceAnalysis.redisBase') },
  ]
})
</script>
