<template>
  <div
    v-if="visible"
    class="card flex flex-col gap-4 p-4 sm:flex-row sm:items-center sm:justify-between"
    :class="activePlan ? 'border-primary-100 dark:border-primary-900/40' : 'border-amber-200/70 dark:border-amber-900/40'"
  >
    <div class="flex min-w-0 items-start gap-3">
      <div
        class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl"
        :class="activePlan ? 'bg-primary-100 text-primary-600 dark:bg-primary-900/30 dark:text-primary-400' : 'bg-amber-100 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400'"
      >
        <Icon :name="activePlan ? 'creditCard' : 'dollar'" size="md" />
      </div>
      <div class="min-w-0">
        <p class="text-sm font-semibold text-gray-900 dark:text-white">
          {{ activePlan ? t('dashboard.planBanner.activePlan', { name: activePlan.group_name || `#${activePlan.group_id}` }) : t('dashboard.planBanner.lowBalance') }}
        </p>
        <p class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">{{ subtitle }}</p>
        <div v-if="activePlan && progressItems.length > 0" class="mt-3 flex flex-wrap gap-x-6 gap-y-2">
          <div v-for="item in progressItems" :key="item.key" class="w-32">
            <div class="flex items-center justify-between text-[11px] text-gray-500 dark:text-gray-400">
              <span>{{ item.label }}</span>
              <span class="tabular-nums">{{ item.percent }}%</span>
            </div>
            <div class="mt-1 h-1.5 overflow-hidden rounded-full bg-gray-200 dark:bg-dark-700">
              <div class="h-full rounded-full" :class="barClass(item.percent)" :style="{ width: `${item.percent}%` }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <RouterLink :to="activePlan ? '/subscriptions' : '/purchase'" class="btn btn-secondary shrink-0">
      {{ activePlan ? t('dashboard.planBanner.manage') : t('dashboard.planBanner.recharge') }}
    </RouterLink>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { getSubscriptionSummary, type SubscriptionSummary, type SubscriptionSummaryItem } from '@/api/subscriptions'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()
const summary = ref<SubscriptionSummary | null>(null)

const subscriptionEnabled = computed(() => appStore.cachedPublicSettings?.subscription_enabled !== false)
const paymentEnabled = computed(() => appStore.cachedPublicSettings?.payment_enabled === true)
const user = computed(() => authStore.user)

const activePlan = computed<SubscriptionSummaryItem | null>(() => {
  const items = (summary.value?.subscriptions ?? []).filter(item => item.status === 'active')
  if (items.length === 0) return null
  return [...items].sort((a, b) => {
    const aTime = a.expires_at ? new Date(a.expires_at).getTime() : Number.MAX_SAFE_INTEGER
    const bTime = b.expires_at ? new Date(b.expires_at).getTime() : Number.MAX_SAFE_INTEGER
    return aTime - bTime
  })[0]
})

const lowBalanceThreshold = computed(() => {
  const userThreshold = user.value?.balance_notify_threshold ?? 0
  const globalThreshold = appStore.cachedPublicSettings?.balance_low_notify_threshold ?? 0
  return userThreshold > 0 ? userThreshold : globalThreshold
})

const lowBalance = computed(() => {
  if (!paymentEnabled.value) return false
  const balance = user.value?.balance ?? 0
  return balance <= 0 || (lowBalanceThreshold.value > 0 && balance <= lowBalanceThreshold.value)
})

const visible = computed(() => Boolean(activePlan.value) || lowBalance.value)

const daysLeft = computed(() => {
  const expiresAt = activePlan.value?.expires_at
  if (!expiresAt) return null
  const time = new Date(expiresAt).getTime()
  if (!Number.isFinite(time)) return null
  return Math.max(0, Math.ceil((time - Date.now()) / 86400000))
})

const subtitle = computed(() => {
  if (activePlan.value) {
    const parts = [daysLeft.value === null ? t('dashboard.planBanner.noExpiry') : t('dashboard.planBanner.daysLeft', { days: daysLeft.value })]
    const extra = (summary.value?.active_count ?? 0) - 1
    if (extra > 0) parts.push(t('dashboard.planBanner.extraSubscriptions', { count: extra }))
    return parts.join(' · ')
  }
  return t('dashboard.planBanner.lowBalanceDescription', { balance: (user.value?.balance ?? 0).toFixed(2) })
})

const progressItems = computed(() => {
  const plan = activePlan.value
  if (!plan) return []
  return [
    { key: 'daily', label: t('dashboard.planBanner.daily'), used: plan.daily_used_usd ?? 0, limit: plan.daily_limit_usd ?? 0 },
    { key: 'weekly', label: t('dashboard.planBanner.weekly'), used: plan.weekly_used_usd ?? 0, limit: plan.weekly_limit_usd ?? 0 },
    { key: 'monthly', label: t('dashboard.planBanner.monthly'), used: plan.monthly_used_usd ?? 0, limit: plan.monthly_limit_usd ?? 0 },
  ]
    .filter(item => item.limit > 0)
    .map(item => ({
      key: item.key,
      label: item.label,
      percent: Math.min(100, Math.max(0, Math.round((item.used / item.limit) * 100))),
    }))
})

const barClass = (percent: number) => {
  if (percent >= 90) return 'bg-red-500'
  if (percent >= 70) return 'bg-amber-500'
  return 'bg-primary-500'
}

onMounted(async () => {
  if (!subscriptionEnabled.value) return
  try {
    summary.value = await getSubscriptionSummary()
  } catch (error) {
    console.warn('Failed to load subscription summary:', error)
  }
})
</script>
