import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

import type { Announcement, DashboardStats, TrendDataPoint } from '@/types'
import DashboardView from '../DashboardView.vue'

const {
  getSnapshotV2,
  getUsageTrend,
  getUserSpendingRanking,
  listAnnouncements,
  getDashboardOverview,
  listAlertEvents
} = vi.hoisted(() => ({
  getSnapshotV2: vi.fn(),
  getUsageTrend: vi.fn(),
  getUserSpendingRanking: vi.fn(),
  listAnnouncements: vi.fn(),
  getDashboardOverview: vi.fn(),
  listAlertEvents: vi.fn()
}))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    dashboard: {
      getSnapshotV2,
      getUsageTrend,
      getUserSpendingRanking
    },
    announcements: {
      list: listAnnouncements
    },
    ops: {
      getDashboardOverview,
      listAlertEvents
    }
  }
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError: vi.fn(),
    cachedPublicSettings: null
  })
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  })
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    })
  }
})

const formatLocalDate = (date: Date): string => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const createDashboardStats = (): DashboardStats => ({
  total_users: 0,
  today_new_users: 0,
  active_users: 0,
  hourly_active_users: 0,
  stats_updated_at: '',
  stats_stale: false,
  total_api_keys: 0,
  active_api_keys: 0,
  total_accounts: 0,
  normal_accounts: 0,
  error_accounts: 0,
  ratelimit_accounts: 0,
  overload_accounts: 0,
  total_requests: 0,
  total_input_tokens: 0,
  total_output_tokens: 0,
  total_cache_creation_tokens: 0,
  total_cache_read_tokens: 0,
  total_tokens: 0,
  total_cost: 0,
  total_actual_cost: 0,
  total_account_cost: 0,
  today_requests: 0,
  today_input_tokens: 0,
  today_output_tokens: 0,
  today_cache_creation_tokens: 0,
  today_cache_read_tokens: 0,
  today_tokens: 0,
  today_cost: 0,
  today_actual_cost: 0,
  today_account_cost: 0,
  average_duration_ms: 0,
  uptime: 0,
  rpm: 0,
  tpm: 0
})

const createTrendPoint = (date: string, tokens: number): TrendDataPoint => ({
  date,
  requests: 1,
  input_tokens: tokens,
  output_tokens: 0,
  cache_creation_tokens: 0,
  cache_read_tokens: 0,
  total_tokens: tokens,
  cost: 0,
  actual_cost: 0
})

const createAnnouncement = (over: Partial<Announcement> & Pick<Announcement, 'id' | 'title'>): Announcement => ({
  content: 'Body',
  status: 'active',
  notify_mode: 'silent',
  targeting: {},
  created_at: '2026-10-01T09:00:00Z',
  updated_at: '2026-10-01T09:00:00Z',
  ...over
})

const globalStubs = {
  AppLayout: { template: '<div><slot /></div>' },
  LoadingSpinner: true,
  Icon: true,
  DateRangePicker: true,
  Select: true,
  ModelDistributionChart: true,
  TokenUsageTrend: true,
  RouterLink: { template: '<a><slot /></a>' }
}

describe('admin DashboardView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())

    getSnapshotV2.mockReset()
    getUsageTrend.mockReset()
    getUserSpendingRanking.mockReset()
    listAnnouncements.mockReset()
    getDashboardOverview.mockReset()
    listAlertEvents.mockReset()

    getSnapshotV2.mockResolvedValue({
      stats: createDashboardStats(),
      trend: [],
      models: []
    })
    getUsageTrend.mockResolvedValue({
      trend: [],
      start_date: '',
      end_date: '',
      granularity: 'day'
    })
    getUserSpendingRanking.mockResolvedValue({
      ranking: [],
      total_actual_cost: 0,
      total_requests: 0,
      total_tokens: 0,
      start_date: '',
      end_date: ''
    })
    listAnnouncements.mockResolvedValue({
      items: [],
      total: 0,
      page: 1,
      page_size: 20,
      pages: 0
    })
    getDashboardOverview.mockResolvedValue({
      health_score: 98,
      error_rate: 0.01,
      duration: { p95_ms: 120 },
      qps: { current: 2.5 }
    })
    listAlertEvents.mockResolvedValue([])
  })

  it('loads ops health, alerts, ranking and only live announcements', async () => {
    getUsageTrend.mockResolvedValue({
      trend: [createTrendPoint('2026-10-01', 1000)],
      start_date: '2026-10-14',
      end_date: '2026-10-15',
      granularity: 'day'
    })
    getUserSpendingRanking.mockResolvedValue({
      ranking: [{ user_id: 7, email: 'top@example.com', username: 'top', actual_cost: 12, requests: 3, tokens: 900 }],
      total_actual_cost: 12,
      total_requests: 3,
      total_tokens: 900,
      start_date: '2026-10-01',
      end_date: '2026-10-15'
    })
    listAnnouncements.mockResolvedValue({
      items: [
        createAnnouncement({ id: 1, title: 'Live announcement' }),
        createAnnouncement({ id: 2, title: 'Scheduled announcement', starts_at: '2099-01-01T00:00:00Z' }),
        createAnnouncement({ id: 3, title: 'Ended announcement', ends_at: '2020-01-01T00:00:00Z' })
      ],
      total: 3,
      page: 1,
      page_size: 20,
      pages: 1
    })

    const wrapper = mount(DashboardView, {
      global: { stubs: globalStubs }
    })
    await flushPromises()

    expect(getUsageTrend).toHaveBeenCalledWith(expect.objectContaining({ granularity: 'day' }))
    expect(getDashboardOverview).toHaveBeenCalledWith({ time_range: '24h' })
    expect(listAlertEvents).toHaveBeenCalledWith({ status: 'firing', limit: 5 })
    expect(listAnnouncements).toHaveBeenCalledWith(1, 20, { status: 'active' })
    expect(wrapper.findAll('[data-month]')).toHaveLength(0)
    expect(wrapper.find('[data-announcement-id="1"]').exists()).toBe(true)
    expect(wrapper.find('[data-announcement-id="2"]').exists()).toBe(false)
    expect(wrapper.find('[data-announcement-id="3"]').exists()).toBe(false)
    expect(wrapper.find('[data-user-id="7"]').exists()).toBe(true)
  })

  it('uses last 24 hours as default dashboard range', async () => {
    mount(DashboardView, {
      global: { stubs: globalStubs }
    })

    await flushPromises()

    const now = new Date()
    const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000)

    expect(getSnapshotV2).toHaveBeenCalledTimes(1)
    expect(getSnapshotV2).toHaveBeenCalledWith(expect.objectContaining({
      start_date: formatLocalDate(yesterday),
      end_date: formatLocalDate(now),
      granularity: 'hour'
    }))
  })
})
