import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

import AggregationView from '../AggregationView.vue'

const { listVendors, listAllEntries, setEntryGroups, listCompositeRoutes, createCompositeRoute, getStats } =
  vi.hoisted(() => ({
    listVendors: vi.fn(),
    listAllEntries: vi.fn(),
    setEntryGroups: vi.fn(),
    listCompositeRoutes: vi.fn(),
    createCompositeRoute: vi.fn(),
    getStats: vi.fn()
  }))

vi.mock('@/api/admin/upstreamHub', () => ({
  default: {
    listVendors,
    listAllEntries,
    setEntryGroups,
    setEntryGroupEnabled: vi.fn()
  }
}))

vi.mock('@/api/admin/groups', () => ({
  default: {
    getAllIncludingInactive: vi.fn(async () => [
      { id: 3, name: 'mock-composite', platform: 'composite', rpm_limit: 60, status: 'active' }
    ]),
    listCompositeRoutes,
    createCompositeRoute,
    updateCompositeRoute: vi.fn(),
    deleteCompositeRoute: vi.fn(),
    getStats
  }
}))

vi.mock('@/constants/platforms', () => ({
  CONCRETE_PLATFORM_OPTIONS: [
    { value: 'openai', label: 'OpenAI' },
    { value: 'deepseek', label: 'DeepSeek' }
  ]
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ showError: vi.fn(), showSuccess: vi.fn() })
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => key })
  }
})

const mountView = () =>
  mount(AggregationView, {
    global: {
      stubs: {
        AppLayout: { template: '<div><slot /></div>' },
        EmptyState: true,
        Select: { props: ['modelValue', 'options'], template: '<select class="select-stub" />' },
        Toggle: { props: ['modelValue'], template: '<button class="toggle-stub" />' }
      }
    }
  })

describe('AggregationView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    listVendors.mockResolvedValue([
      {
        id: 2,
        name: 'MockVendor',
        concurrency_budget: 10,
        allocated_concurrency: 5,
        entry_count: 1
      }
    ])
    listAllEntries.mockResolvedValue([
      {
        account_id: 9,
        vendor_id: 2,
        name: 'Mock-DeepSeek',
        platform: 'deepseek',
        base_url: 'http://127.0.0.1:18081/v1',
        has_api_key: true,
        concurrency: 5,
        priority: 50,
        schedulable: true,
        status: 'active',
        groups: [
          { group_id: 3, group_name: 'mock-composite', platform: 'composite', priority: 50, enabled: true }
        ],
        models: ['mock-deepseek-chat']
      }
    ])
    listCompositeRoutes.mockResolvedValue([
      {
        id: 1,
        group_id: 3,
        public_model: 'mock-public-model',
        match_type: 'exact',
        target_platform: 'deepseek',
        upstream_model: 'mock-deepseek-chat',
        endpoint: 'any',
        priority: 100,
        enabled: true,
        notes: ''
      }
    ])
    getStats.mockResolvedValue({ total_api_keys: 1, active_api_keys: 1, total_cost: 0, total_requests: 0 })
    setEntryGroups.mockResolvedValue([])
    createCompositeRoute.mockResolvedValue({ id: 2 })
  })

  it('renders vendors, entries and routes for the selected composite group', async () => {
    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.text()).toContain('MockVendor')
    expect(wrapper.text()).toContain('Mock-DeepSeek')
    expect(wrapper.text()).toContain('mock-public-model')
    expect(getStats).toHaveBeenCalledWith(3)
  })

  it('detaches an attached entry from the composite group', async () => {
    const wrapper = mountView()
    await flushPromises()

    await wrapper.find('input[type="checkbox"]').trigger('change')
    await flushPromises()

    expect(setEntryGroups).toHaveBeenCalledWith(9, [])
  })

  it('creates an exact-match composite route', async () => {
    const wrapper = mountView()
    await flushPromises()

    const publicModelInput = wrapper.find('input[placeholder="admin.aggregation.publicModelPlaceholder"]')
    await publicModelInput.setValue('mock-new-model')
    const addButton = wrapper
      .findAll('button')
      .find((button) => button.text() === 'admin.aggregation.addRoute')
    expect(addButton).toBeTruthy()
    await addButton!.trigger('click')
    await flushPromises()

    expect(createCompositeRoute).toHaveBeenCalledWith(
      3,
      expect.objectContaining({ public_model: 'mock-new-model', match_type: 'exact', target_platform: 'openai' })
    )
  })
})
