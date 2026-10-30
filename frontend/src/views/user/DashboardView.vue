<template>
  <AppLayout>
    <div class="space-y-6">
      <div v-if="loading" class="flex items-center justify-center py-12"><LoadingSpinner /></div>
      <template v-else-if="stats">
        <UserDashboardPlanBanner />
        <UserDashboardStats :stats="stats" :balance="user?.balance || 0" :is-simple="authStore.isSimpleMode" :platform-quotas="platformQuotas" :previous-day="previousDay" />

        <div v-if="stats.total_requests === 0" class="card flex flex-col gap-4 p-6 sm:flex-row sm:items-center sm:justify-between">
          <div class="flex items-start gap-4">
            <div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-primary-100 text-primary-600 dark:bg-primary-900/30 dark:text-primary-400">
              <Icon name="key" size="lg" />
            </div>
            <div>
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('dashboard.gettingStarted.title') }}</h2>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('dashboard.gettingStarted.description') }}</p>
            </div>
          </div>
          <div class="flex flex-wrap gap-2 sm:justify-end">
            <RouterLink to="/keys" class="btn btn-primary">{{ t('dashboard.gettingStarted.createKey') }}</RouterLink>
            <a v-if="appStore.docUrl" :href="appStore.docUrl" target="_blank" rel="noopener noreferrer" class="btn btn-secondary">{{ t('dashboard.gettingStarted.docs') }}</a>
            <RouterLink v-if="modelPlazaEnabled" to="/model-plaza" class="btn btn-secondary">{{ t('dashboard.gettingStarted.modelPlaza') }}</RouterLink>
          </div>
        </div>

        <UserDashboardCumulativeUsage v-if="stats.total_requests > 0" :stats="stats" />

        <UserDashboardCharts v-model:startDate="startDate" v-model:endDate="endDate" v-model:granularity="granularity" :loading="loadingCharts" :trend="trendData" :models="modelStats" @dateRangeChange="loadCharts" @granularityChange="loadCharts" @refresh="refreshAll" />
        <TokenUsageHeatmap v-if="stats.total_requests > 0" :trend="heatmapTrend" :loading="loadingHeatmap" />
        <div class="grid grid-cols-1 gap-6 lg:h-[30rem] lg:min-h-0 lg:grid-cols-3">
          <div class="h-full min-h-0 overflow-hidden lg:col-span-2"><UserDashboardAnnouncements /></div>
          <div class="h-full min-h-0 overflow-hidden lg:col-span-1"><UserDashboardQuickActions /></div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'; import { useAuthStore } from '@/stores/auth'; import { useAppStore } from '@/stores/app'; import { useI18n } from 'vue-i18n'; import { usageAPI, type UserDashboardStats as UserStatsType } from '@/api/usage'
import AppLayout from '@/components/layout/AppLayout.vue'; import LoadingSpinner from '@/components/common/LoadingSpinner.vue'; import Icon from '@/components/icons/Icon.vue'
import UserDashboardStats from '@/components/user/dashboard/UserDashboardStats.vue'; import UserDashboardCharts from '@/components/user/dashboard/UserDashboardCharts.vue'
import TokenUsageHeatmap from '@/components/charts/TokenUsageHeatmap.vue'; import UserDashboardAnnouncements from '@/components/user/dashboard/UserDashboardAnnouncements.vue'
import UserDashboardQuickActions from '@/components/user/dashboard/UserDashboardQuickActions.vue'; import UserDashboardPlanBanner from '@/components/user/dashboard/UserDashboardPlanBanner.vue'
import UserDashboardCumulativeUsage from '@/components/user/dashboard/UserDashboardCumulativeUsage.vue'
import type { TrendDataPoint, ModelStat, PlatformQuotaItem } from '@/types'
import { getMyPlatformQuotas } from '@/api/user'
import { formatDateLocalInput } from '@/utils/format'
import { previousDaySnapshot } from '@/utils/dashboardMetrics'

const { t } = useI18n(); const authStore = useAuthStore(); const appStore = useAppStore(); const user = computed(() => authStore.user)
const modelPlazaEnabled = computed(() => appStore.cachedPublicSettings?.model_plaza_enabled !== false)
const stats = ref<UserStatsType | null>(null); const loading = ref(false); const loadingCharts = ref(false); const loadingHeatmap = ref(false)
const trendData = ref<TrendDataPoint[]>([]); const modelStats = ref<ModelStat[]>([]); const heatmapTrend = ref<TrendDataPoint[]>([])
const platformQuotas = ref<PlatformQuotaItem[] | null>(null)
const previousDay = computed(() => previousDaySnapshot(heatmapTrend.value))

const startDate = ref(formatDateLocalInput(new Date(Date.now() - 6 * 86400000))); const endDate = ref(formatDateLocalInput(new Date())); const granularity = ref('day')

const loadStats = async () => { loading.value = true; try { await authStore.refreshUser(); stats.value = await usageAPI.getDashboardStats() } catch (error) { console.error('Failed to load dashboard stats:', error) } finally { loading.value = false } }
const loadCharts = async () => { loadingCharts.value = true; try { const res = await Promise.all([usageAPI.getDashboardTrend({ start_date: startDate.value, end_date: endDate.value, granularity: granularity.value as any }), usageAPI.getDashboardModels({ start_date: startDate.value, end_date: endDate.value })]); trendData.value = res[0].trend || []; modelStats.value = res[1].models || [] } catch (error) { console.error('Failed to load charts:', error) } finally { loadingCharts.value = false } }
const loadHeatmap = async () => { loadingHeatmap.value = true; try { const now = new Date(); const start = new Date(now.getFullYear(), now.getMonth() - 11, 1); const res = await usageAPI.getDashboardTrend({ start_date: formatDateLocalInput(start), end_date: formatDateLocalInput(now), granularity: 'day' }); heatmapTrend.value = res.trend || [] } catch (error) { console.error('Failed to load token heatmap:', error) } finally { loadingHeatmap.value = false } }
const loadPlatformQuotas = async () => { try { const data = await getMyPlatformQuotas(); platformQuotas.value = data.platform_quotas ?? [] } catch (error) { console.warn('Failed to load platform quotas:', error); platformQuotas.value = [] } }
const refreshAll = () => { loadStats(); loadCharts(); loadHeatmap(); loadPlatformQuotas() }

onMounted(() => { refreshAll() })
</script>
