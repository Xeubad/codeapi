import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

const getSubscriptionSummary = vi.hoisted(() => vi.fn())
vi.mock('@/api/subscriptions', () => ({ getSubscriptionSummary }))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => key }),
  }
})

import UserDashboardPlanBanner from '../UserDashboardPlanBanner.vue'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'

const stubs = {
  RouterLink: { template: '<a><slot /></a>' },
  Icon: true,
}

function mountBanner() {
  const pinia = createPinia()
  setActivePinia(pinia)
  const authStore = useAuthStore()
  const appStore = useAppStore()
  authStore.user = {
    balance: 1,
    balance_notify_threshold: 5,
  } as any
  appStore.cachedPublicSettings = {
    subscription_enabled: true,
    payment_enabled: true,
    balance_low_notify_threshold: 0,
  } as any
  return mount(UserDashboardPlanBanner, {
    global: { plugins: [pinia], stubs },
  })
}

describe('UserDashboardPlanBanner', () => {
  beforeEach(() => {
    getSubscriptionSummary.mockReset()
  })

  it('shows the active plan, remaining days and quota progress', async () => {
    const expiresAt = new Date(Date.now() + 5 * 86400000).toISOString()
    getSubscriptionSummary.mockResolvedValue({
      active_count: 1,
      subscriptions: [{
        id: 1,
        group_id: 2,
        group_name: 'Pro',
        status: 'active',
        daily_used_usd: 2,
        daily_limit_usd: 10,
        expires_at: expiresAt,
      }],
    })

    const wrapper = mountBanner()
    await flushPromises()

    expect(wrapper.text()).toContain('dashboard.planBanner.activePlan')
    expect(wrapper.text()).toContain('dashboard.planBanner.daysLeft')
    expect(wrapper.text()).toContain('dashboard.planBanner.daily')
    expect(wrapper.text()).toContain('20%')
    expect(wrapper.text()).toContain('dashboard.planBanner.manage')
  })

  it('falls back to a low balance prompt without an active plan', async () => {
    getSubscriptionSummary.mockResolvedValue({ active_count: 0, subscriptions: [] })

    const wrapper = mountBanner()
    await flushPromises()

    expect(wrapper.text()).toContain('dashboard.planBanner.lowBalance')
    expect(wrapper.text()).toContain('dashboard.planBanner.recharge')
  })
})
