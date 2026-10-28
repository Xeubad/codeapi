import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

import AggregationView from '../AggregationView.vue'

const {
  listVendors,
  listAllEntries,
  setEntryGroups,
  setEntryGroupEnabled,
  getAllIncludingInactive,
  listCompositeRoutes,
  createCompositeRoute,
  getStats,
  createGroup
} = vi.hoisted(() => ({
  listVendors: vi.fn(),
  listAllEntries: vi.fn(),
  setEntryGroups: vi.fn(),
  setEntryGroupEnabled: vi.fn(),
  getAllIncludingInactive: vi.fn(),
  listCompositeRoutes: vi.fn(),
  createCompositeRoute: vi.fn(),
  getStats: vi.fn(),
  createGroup: vi.fn()
}))

vi.mock('@/api/admin/upstreamHub', () => ({
  default: {
    listVendors,
    listAllEntries,
    setEntryGroups,
    setEntryGroupEnabled
  }
}))

vi.mock('@/api/admin/groups', () => ({
  default: {
    getAllIncludingInactive,
    listCompositeRoutes,
    createCompositeRoute,
    updateCompositeRoute: vi.fn(),
    deleteCompositeRoute: vi.fn(),
    getStats,
    create: createGroup
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

const replace = vi.fn()
vi.mock('vue-router', () => ({
  useRoute: () => ({ query: {} }),
  useRouter: () => ({ replace, push: vi.fn() })
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
        BaseDialog: {
          props: ['show', 'title'],
          template: '<div v-if="show"><slot /><slot name="footer" /></div>'
        },
        Select: { props: ['modelValue', 'options'], template: '<select class="select-stub" />' },
        Toggle: { props: ['modelValue'], template: '<button class="toggle-stub" />' }
      }
    }
  })

describe('AggregationView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    getAllIncludingInactive.mockResolvedValue([
      { id: 3, name: 'mock-composite', platform: 'composite', rpm_limit: 60, status: 'active' },
      { id: 4, name: 'plain-group', platform: 'openai', rpm_limit: 0, status: 'active' }
    ])
    listVendors.mockResolvedValue([
      { id: 2, name: 'MockVendor', concurrency_budget: 10, allocated_concurrency: 5 }
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
    createGroup.mockResolvedValue({ id: 7, name: 'new-unit', platform: 'composite' })
  })

  it('lists only composite groups as scheduling units with attached counts', async () => {
    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.text()).toContain('mock-composite')
    expect(wrapper.text()).not.toContain('plain-group')
  })

  it('opens a unit and detaches an attached entry', async () => {
    const wrapper = mountView()
    await flushPromises()

    const manageButton = wrapper.findAll('button').find((button) => button.text() === 'admin.aggregation.manage')
    expect(manageButton).toBeTruthy()
    await manageButton!.trigger('click')
    await flushPromises()
    expect(replace).toHaveBeenCalled()

    await wrapper.find('input[type="checkbox"]').trigger('change')
    await flushPromises()
    expect(setEntryGroups).toHaveBeenCalledWith(9, [])
  })

  it('creates a scheduling unit as a composite group', async () => {
    const wrapper = mountView()
    await flushPromises()

    const openButton = wrapper
      .findAll('button')
      .find((button) => button.text() === 'admin.aggregation.createUnit')
    await openButton!.trigger('click')
    await flushPromises()

    await wrapper.find('input[placeholder="admin.aggregation.unitNamePlaceholder"]').setValue('new-unit')
    const submitButton = wrapper
      .findAll('button')
      .filter((button) => button.text() === 'admin.aggregation.createUnit')
      .at(-1)
    await submitButton!.trigger('click')
    await flushPromises()

    expect(createGroup).toHaveBeenCalledWith(
      expect.objectContaining({ name: 'new-unit', platform: 'composite' })
    )
  })
})
