<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <LoadingSpinner />
      </div>

      <template v-else-if="stats">
        <!-- KPI Cards -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
          <MetricCard
            :label="t('admin.dashboard.todayRequests')"
            icon="chart"
            tone="green"
            :delta-current="stats.today_requests"
            :delta-previous="previousDay?.requests ?? null"
          >
            <template #value>
              <p class="text-3xl font-bold tabular-nums text-gray-900 dark:text-white">{{ stats.today_requests }}</p>
            </template>
            <template #footer>
              <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
                <span>{{ t('admin.dashboard.totalRequests') }}</span>
                <span class="font-semibold tabular-nums text-gray-900 dark:text-white">{{ formatNumber(stats.total_requests) }}</span>
              </div>
            </template>
          </MetricCard>

          <MetricCard
            :label="t('admin.dashboard.todayTokens')"
            icon="cube"
            tone="amber"
            :delta-current="stats.today_tokens"
            :delta-previous="previousDay?.tokens ?? null"
          >
            <template #value>
              <p class="text-3xl font-bold tabular-nums text-gray-900 dark:text-white">{{ formatTokens(stats.today_tokens) }}</p>
            </template>
            <template #meta>
              <p class="text-[11px] leading-5 text-gray-500 dark:text-gray-400">
                {{ t('admin.dashboard.input') }} {{ formatTokens(stats.today_input_tokens) }} ·
                {{ t('admin.dashboard.output') }} {{ formatTokens(stats.today_output_tokens) }} ·
                {{ t('admin.dashboard.cacheToday') }} {{ formatTokens(stats.today_cache_creation_tokens + stats.today_cache_read_tokens) }}
              </p>
            </template>
            <template #footer>
              <div class="flex flex-wrap items-center justify-between gap-x-3 gap-y-1 text-xs text-gray-500 dark:text-gray-400">
                <span>{{ t('admin.dashboard.totalTokens') }} <span class="font-semibold tabular-nums text-gray-900 dark:text-white">{{ formatTokens(stats.total_tokens) }}</span></span>
                <span>{{ t('dashboard.cacheHitRate') }} <span class="font-semibold tabular-nums text-gray-900 dark:text-white">{{ cacheHitRateToday.toFixed(1) }}%</span></span>
              </div>
            </template>
          </MetricCard>

          <MetricCard
            :label="t('admin.dashboard.todayCost')"
            icon="dollar"
            tone="purple"
            :delta-current="stats.today_actual_cost"
            :delta-previous="previousDay?.cost ?? null"
            :delta-positive-is-good="false"
          >
            <template #value>
              <p class="text-3xl font-bold tabular-nums text-purple-600 dark:text-purple-400">${{ formatCost(stats.today_actual_cost) }}</p>
            </template>
            <template #meta>
              <p class="text-[11px] leading-5 text-gray-500 dark:text-gray-400">
                {{ t('admin.dashboard.standard') }} ${{ formatCost(stats.today_cost) }} ·
                {{ t('admin.dashboard.accountCost') }} ${{ formatCost(stats.today_account_cost) }}
              </p>
            </template>
            <template #footer>
              <div class="flex flex-wrap items-center justify-between gap-x-3 gap-y-1 text-xs text-gray-500 dark:text-gray-400">
                <span>{{ t('admin.dashboard.totalCost') }} <span class="font-semibold tabular-nums text-gray-900 dark:text-white">${{ formatCost(stats.total_actual_cost) }}</span></span>
                <span>{{ t('admin.dashboard.standard') }} <span class="tabular-nums text-gray-700 dark:text-gray-300">${{ formatCost(stats.total_cost) }}</span></span>
              </div>
            </template>
          </MetricCard>

          <MetricCard
            :label="t('admin.dashboard.newUsersToday')"
            icon="userPlus"
            tone="emerald"
          >
            <template #value>
              <p class="text-3xl font-bold tabular-nums text-emerald-600 dark:text-emerald-400">+{{ stats.today_new_users }}</p>
            </template>
            <template #footer>
              <div class="flex flex-wrap items-center justify-between gap-x-3 gap-y-1 text-xs text-gray-500 dark:text-gray-400">
                <span>{{ t('admin.dashboard.totalUsers') }} <span class="font-semibold tabular-nums text-gray-900 dark:text-white">{{ formatNumber(stats.total_users) }}</span></span>
                <span>{{ t('admin.dashboard.activeUsers') }} <span class="font-semibold tabular-nums text-gray-900 dark:text-white">{{ formatNumber(stats.active_users) }}</span></span>
              </div>
            </template>
          </MetricCard>

          <MetricCard
            :label="t('admin.dashboard.apiKeys')"
            icon="key"
            tone="blue"
          >
            <template #value>
              <p class="text-3xl font-bold tabular-nums text-gray-900 dark:text-white">{{ stats.total_api_keys }}</p>
            </template>
            <template #meta>
              <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ stats.active_api_keys }} {{ t('admin.dashboard.activeApiKeys') }}</p>
            </template>
            <template #footer>
              <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
                <span>{{ t('admin.dashboard.activeRate') }}</span>
                <span class="font-semibold tabular-nums text-gray-900 dark:text-white">{{ activeKeyRate }}%</span>
              </div>
              <div class="mt-2 h-1.5 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
                <div class="h-full rounded-full bg-blue-500" :style="{ width: `${activeKeyRate}%` }"></div>
              </div>
            </template>
          </MetricCard>

          <MetricCard
            :label="t('admin.dashboard.accounts')"
            icon="server"
            tone="indigo"
          >
            <template #value>
              <p class="text-3xl font-bold tabular-nums text-gray-900 dark:text-white">{{ stats.total_accounts }}</p>
            </template>
            <template #meta>
              <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ stats.normal_accounts }} {{ t('admin.dashboard.ok') }}</p>
            </template>
            <template #footer>
              <div class="flex flex-wrap gap-x-3 gap-y-1 text-[11px] text-gray-500 dark:text-gray-400">
                <span>{{ t('admin.dashboard.health.normal') }} <span class="font-semibold text-emerald-600 dark:text-emerald-400">{{ stats.normal_accounts }}</span></span>
                <span>{{ t('admin.dashboard.health.error') }} <span class="font-semibold text-red-500 dark:text-red-400">{{ stats.error_accounts }}</span></span>
                <span>{{ t('admin.dashboard.health.rateLimited') }} <span class="font-semibold text-amber-600 dark:text-amber-400">{{ stats.ratelimit_accounts }}</span></span>
                <span>{{ t('admin.dashboard.health.overloaded') }} <span class="font-semibold text-orange-600 dark:text-orange-400">{{ stats.overload_accounts }}</span></span>
              </div>
            </template>
          </MetricCard>

          <MetricCard
            :label="t('admin.dashboard.performance')"
            icon="bolt"
            tone="violet"
          >
            <template #value>
              <p class="flex items-baseline gap-1.5 text-3xl font-bold tabular-nums text-gray-900 dark:text-white">
                {{ formatTokens(stats.rpm) }}<span class="text-xs font-medium text-gray-500 dark:text-gray-400">RPM</span>
              </p>
            </template>
            <template #meta>
              <p class="text-[11px] text-gray-500 dark:text-gray-400">TPM <span class="font-semibold tabular-nums text-violet-600 dark:text-violet-400">{{ formatTokens(stats.tpm) }}</span></p>
            </template>
            <template #footer>
              <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
                <span>{{ t('admin.dashboard.realtime') }}</span>
                <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
              </div>
            </template>
          </MetricCard>

          <MetricCard
            :label="t('admin.dashboard.avgResponse')"
            icon="clock"
            tone="rose"
          >
            <template #value>
              <p class="text-3xl font-bold tabular-nums text-gray-900 dark:text-white">{{ formatDuration(stats.average_duration_ms) }}</p>
            </template>
            <template #footer>
              <div class="flex flex-wrap items-center justify-between gap-x-3 gap-y-1 text-xs text-gray-500 dark:text-gray-400">
                <span>{{ stats.active_users }} {{ t('admin.dashboard.activeUsers') }}</span>
                <span>{{ t('admin.dashboard.currentHourActive') }} {{ stats.hourly_active_users }}</span>
              </div>
            </template>
          </MetricCard>
        </div>

        <!-- Operations Health -->
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <div class="lg:col-span-2">
            <SystemHealthCard
              :stats="stats"
              :overview="opsOverview"
              :loading="opsLoading"
              :error="opsError"
            />
          </div>
          <div class="lg:col-span-1">
            <AlertEventsCard
              :events="opsAlerts"
              :loading="alertsLoading"
              :error="alertsError"
            />
          </div>
        </div>

        <!-- Revenue & Top Spenders -->
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <div v-if="paymentEnabled" class="lg:col-span-1">
            <RevenueOverviewCard
              :stats="paymentStats"
              :loading="paymentLoading"
              :error="paymentError"
              :today-cost="stats.today_actual_cost"
              :total-cost="stats.total_actual_cost"
            />
          </div>
          <div :class="paymentEnabled ? 'lg:col-span-2' : 'lg:col-span-3'">
            <TopSpendersCard
              :items="rankingItems"
              :loading="rankingLoading"
              :error="rankingError"
              :total-actual-cost="rankingTotalActualCost"
              :total-requests="rankingTotalRequests"
              :total-tokens="rankingTotalTokens"
              @select="goToUserUsage"
            />
          </div>
        </div>

        <!-- Charts Section -->
        <div class="space-y-6">
          <!-- Date Range Filter -->
          <div class="card p-4">
            <div class="flex flex-wrap items-center gap-4">
              <div class="flex items-center gap-2">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-300"
                  >{{ t('admin.dashboard.timeRange') }}:</span
                >
                <DateRangePicker
                  v-model:start-date="startDate"
                  v-model:end-date="endDate"
                  @change="onDateRangeChange"
                />
              </div>
              <button @click="loadDashboardStats" :disabled="chartsLoading" class="btn btn-secondary">
                {{ t('common.refresh') }}
              </button>
              <RouterLink to="/admin/usage" class="text-xs font-medium text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
                {{ t('dashboard.viewAllUsage') }} →
              </RouterLink>
              <div class="ml-auto flex items-center gap-2">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-300"
                  >{{ t('admin.dashboard.granularity') }}:</span
                >
                <div class="w-28">
                  <Select
                    v-model="granularity"
                    :options="granularityOptions"
                    @change="loadChartData"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- Charts Grid -->
          <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
            <ModelDistributionChart
              :model-stats="modelStats"
              :enable-ranking-view="true"
              :ranking-items="rankingItems"
              :ranking-total-actual-cost="rankingTotalActualCost"
              :ranking-total-requests="rankingTotalRequests"
              :ranking-total-tokens="rankingTotalTokens"
              :loading="chartsLoading"
              :ranking-loading="rankingLoading"
              :ranking-error="rankingError"
              :start-date="startDate"
              :end-date="endDate"
              @ranking-click="goToUserUsage"
            />
            <TokenUsageTrend :trend-data="trendData" :loading="chartsLoading" />
          </div>
        </div>

        <!-- System Announcements -->
        <AnnouncementTimeline :items="adminAnnouncements" :loading="announcementLoading" :limit="5">
          <template #action>
            <RouterLink to="/admin/announcements" class="text-xs font-medium text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
              {{ t('dashboard.announcements.manage') }} →
            </RouterLink>
          </template>
        </AnnouncementTimeline>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
import { adminAPI } from '@/api/admin'
import type {
  Announcement,
  DashboardStats,
  TrendDataPoint,
  ModelStat,
  UserSpendingRankingItem
} from '@/types'
import type { AlertEvent, OpsDashboardOverview } from '@/api/admin/ops'
import type { DashboardStats as PaymentDashboardStats } from '@/types/payment'
import AppLayout from '@/components/layout/AppLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import Select from '@/components/common/Select.vue'
import ModelDistributionChart from '@/components/charts/ModelDistributionChart.vue'
import TokenUsageTrend from '@/components/charts/TokenUsageTrend.vue'
import AnnouncementTimeline from '@/components/common/AnnouncementTimeline.vue'
import MetricCard from '@/components/dashboard/MetricCard.vue'
import SystemHealthCard from '@/components/admin/dashboard/SystemHealthCard.vue'
import AlertEventsCard from '@/components/admin/dashboard/AlertEventsCard.vue'
import RevenueOverviewCard from '@/components/admin/dashboard/RevenueOverviewCard.vue'
import TopSpendersCard from '@/components/admin/dashboard/TopSpendersCard.vue'
import { adminPaymentAPI } from '@/api/admin/payment'
import { computeCacheHitRate, previousDaySnapshot } from '@/utils/dashboardMetrics'

import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
  Filler
)

const appStore = useAppStore()
const router = useRouter()
const stats = ref<DashboardStats | null>(null)
const loading = ref(false)
const chartsLoading = ref(false)
const announcementLoading = ref(false)
const rankingLoading = ref(false)
const rankingError = ref(false)
const opsLoading = ref(false)
const opsError = ref(false)
const alertsLoading = ref(false)
const alertsError = ref(false)
const paymentLoading = ref(false)
const paymentError = ref(false)

// Chart data
const trendData = ref<TrendDataPoint[]>([])
const modelStats = ref<ModelStat[]>([])
const comparisonTrend = ref<TrendDataPoint[]>([])
const adminAnnouncements = ref<Announcement[]>([])
const rankingItems = ref<UserSpendingRankingItem[]>([])
const rankingTotalActualCost = ref(0)
const rankingTotalRequests = ref(0)
const rankingTotalTokens = ref(0)
const opsOverview = ref<OpsDashboardOverview | null>(null)
const opsAlerts = ref<AlertEvent[]>([])
const paymentStats = ref<PaymentDashboardStats | null>(null)
let chartLoadSeq = 0
let rankingLoadSeq = 0
const rankingLimit = 12

const previousDay = computed(() => previousDaySnapshot(comparisonTrend.value))
const paymentEnabled = computed(() => appStore.cachedPublicSettings?.payment_enabled === true)
const activeKeyRate = computed(() => {
  const total = stats.value?.total_api_keys ?? 0
  if (total <= 0) return 0
  return Math.round(((stats.value?.active_api_keys ?? 0) / total) * 100)
})
const cacheHitRateToday = computed(() => computeCacheHitRate(
  stats.value?.today_input_tokens ?? 0,
  stats.value?.today_cache_read_tokens ?? 0,
  stats.value?.today_cache_creation_tokens ?? 0,
))

// Helper function to format date in local timezone
const formatLocalDate = (date: Date): string => {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

const getLast24HoursRangeDates = (): { start: string; end: string } => {
  const end = new Date()
  const start = new Date(end.getTime() - 24 * 60 * 60 * 1000)
  return {
    start: formatLocalDate(start),
    end: formatLocalDate(end)
  }
}

// Date range
const granularity = ref<'day' | 'hour'>('hour')
const defaultRange = getLast24HoursRangeDates()
const startDate = ref(defaultRange.start)
const endDate = ref(defaultRange.end)

// Granularity options for Select component
const granularityOptions = computed(() => [
  { value: 'day', label: t('admin.dashboard.day') },
  { value: 'hour', label: t('admin.dashboard.hour') }
])

// Format helpers
const formatTokens = (value: number | undefined): string => {
  if (value === undefined || value === null) return '0'
  if (value >= 1_000_000_000) {
    return `${(value / 1_000_000_000).toFixed(2)}B`
  } else if (value >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(2)}M`
  } else if (value >= 1_000) {
    return `${(value / 1_000).toFixed(2)}K`
  }
  return value.toLocaleString()
}

const toFiniteNumber = (value: unknown): number => {
  const numberValue = Number(value)
  return Number.isFinite(numberValue) ? numberValue : 0
}

const formatNumber = (value: number | null | undefined): string => {
  return toFiniteNumber(value).toLocaleString()
}

const formatCost = (value: number | null | undefined): string => {
  const safeValue = toFiniteNumber(value)
  if (safeValue >= 1000) {
    return (safeValue / 1000).toFixed(2) + 'K'
  } else if (safeValue >= 1) {
    return safeValue.toFixed(2)
  } else if (safeValue >= 0.01) {
    return safeValue.toFixed(3)
  }
  return safeValue.toFixed(4)
}

const formatDuration = (ms: number): string => {
  if (ms >= 1000) {
    return `${(ms / 1000).toFixed(2)}s`
  }
  return `${Math.round(ms)}ms`
}

const goToUserUsage = (item: UserSpendingRankingItem) => {
  void router.push({
    path: '/admin/usage',
    query: {
      user_id: String(item.user_id),
      start_date: startDate.value,
      end_date: endDate.value
    }
  })
}

// Date range change handler
const onDateRangeChange = (range: {
  startDate: string
  endDate: string
  preset: string | null
}) => {
  // Auto-select granularity based on date range
  const start = new Date(range.startDate)
  const end = new Date(range.endDate)
  const daysDiff = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24))

  // If range is 1 day, use hourly granularity
  if (daysDiff <= 1) {
    granularity.value = 'hour'
  } else {
    granularity.value = 'day'
  }

  loadChartData()
}

// Load data
const loadDashboardSnapshot = async (includeStats: boolean) => {
  const currentSeq = ++chartLoadSeq
  if (includeStats && !stats.value) {
    loading.value = true
  }
  chartsLoading.value = true
  try {
    const response = await adminAPI.dashboard.getSnapshotV2({
      start_date: startDate.value,
      end_date: endDate.value,
      granularity: granularity.value,
      include_stats: includeStats,
      include_trend: true,
      include_model_stats: true,
      include_group_stats: false,
      include_users_trend: false
    })
    if (currentSeq !== chartLoadSeq) return
    if (includeStats && response.stats) {
      stats.value = response.stats
    }
    trendData.value = response.trend || []
    modelStats.value = response.models || []
  } catch (error) {
    if (currentSeq !== chartLoadSeq) return
    appStore.showError(t('admin.dashboard.failedToLoad'))
    console.error('Error loading dashboard snapshot:', error)
  } finally {
    if (currentSeq === chartLoadSeq) {
      loading.value = false
      chartsLoading.value = false
    }
  }
}

const loadComparison = async () => {
  try {
    const now = new Date()
    const yesterday = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 1)
    const response = await adminAPI.dashboard.getUsageTrend({
      start_date: formatLocalDate(yesterday),
      end_date: formatLocalDate(now),
      granularity: 'day'
    })
    comparisonTrend.value = response.trend || []
  } catch (error) {
    console.error('Error loading day-over-day trend:', error)
    comparisonTrend.value = []
  }
}

const loadOpsHealth = async () => {
  opsLoading.value = true
  opsError.value = false
  try {
    opsOverview.value = await adminAPI.ops.getDashboardOverview({ time_range: '24h' })
  } catch (error) {
    console.error('Error loading ops overview:', error)
    opsOverview.value = null
    opsError.value = true
  } finally {
    opsLoading.value = false
  }
}

const loadOpsAlerts = async () => {
  alertsLoading.value = true
  alertsError.value = false
  try {
    opsAlerts.value = await adminAPI.ops.listAlertEvents({ status: 'firing', limit: 5 })
  } catch (error) {
    console.error('Error loading ops alerts:', error)
    opsAlerts.value = []
    alertsError.value = true
  } finally {
    alertsLoading.value = false
  }
}

const loadRevenue = async () => {
  if (!paymentEnabled.value) {
    paymentStats.value = null
    paymentError.value = false
    return
  }
  paymentLoading.value = true
  paymentError.value = false
  try {
    const response = await adminPaymentAPI.getDashboard(30)
    paymentStats.value = response.data
  } catch (error) {
    console.error('Error loading payment dashboard:', error)
    paymentStats.value = null
    paymentError.value = true
  } finally {
    paymentLoading.value = false
  }
}

const loadAnnouncements = async () => {
  announcementLoading.value = true
  try {
    const response = await adminAPI.announcements.list(1, 20, { status: 'active' })
    const now = Date.now()
    adminAnnouncements.value = response.items.filter((item) => {
      const startsAt = item.starts_at ? new Date(item.starts_at).getTime() : null
      const endsAt = item.ends_at ? new Date(item.ends_at).getTime() : null
      return (startsAt === null || startsAt <= now) && (endsAt === null || endsAt > now)
    })
  } catch (error) {
    console.error('Error loading announcements:', error)
    adminAnnouncements.value = []
  } finally {
    announcementLoading.value = false
  }
}

const loadUserSpendingRanking = async () => {
  const currentSeq = ++rankingLoadSeq
  rankingLoading.value = true
  rankingError.value = false
  try {
    const response = await adminAPI.dashboard.getUserSpendingRanking({
      start_date: startDate.value,
      end_date: endDate.value,
      limit: rankingLimit
    })
    if (currentSeq !== rankingLoadSeq) return
    rankingItems.value = response.ranking || []
    rankingTotalActualCost.value = response.total_actual_cost || 0
    rankingTotalRequests.value = response.total_requests || 0
    rankingTotalTokens.value = response.total_tokens || 0
  } catch (error) {
    if (currentSeq !== rankingLoadSeq) return
    console.error('Error loading user spending ranking:', error)
    rankingItems.value = []
    rankingTotalActualCost.value = 0
    rankingTotalRequests.value = 0
    rankingTotalTokens.value = 0
    rankingError.value = true
  } finally {
    if (currentSeq === rankingLoadSeq) {
      rankingLoading.value = false
    }
  }
}

const loadDashboardStats = async () => {
  await Promise.all([
    loadDashboardSnapshot(true),
    loadComparison(),
    loadAnnouncements(),
    loadUserSpendingRanking(),
    loadOpsHealth(),
    loadOpsAlerts(),
    loadRevenue()
  ])
}

const loadChartData = async () => {
  await Promise.all([
    loadDashboardSnapshot(false),
    loadUserSpendingRanking()
  ])
}

onMounted(() => {
  loadDashboardStats()
})
</script>

<style scoped>
</style>
