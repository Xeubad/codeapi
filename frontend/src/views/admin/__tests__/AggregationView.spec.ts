import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

import AggregationView from '../AggregationView.vue'

const {
  listVendors,
  listAllEntries,
  setEntryGroups,
  setEntryGroupEnabled,
  getAllIncludingInactive,
  listAggregationUnits,
  updateGroup,
  deleteGroup,
  listCompositeRoutes,
  createCompositeRoute,
  updateCompositeRoute,
  deleteCompositeRoute,
  getGroupApiKeys,
  createGroup
} = vi.hoisted(() => ({
  listVendors: vi.fn(),
  listAllEntries: vi.fn(),
  setEntryGroups: vi.fn(),
  setEntryGroupEnabled: vi.fn(),
  getAllIncludingInactive: vi.fn(),
  listAggregationUnits: vi.fn(),
  updateGroup: vi.fn(),
  deleteGroup: vi.fn(),
  listCompositeRoutes: vi.fn(),
  createCompositeRoute: vi.fn(),
  updateCompositeRoute: vi.fn(),
  deleteCompositeRoute: vi.fn(),
  getGroupApiKeys: vi.fn(),
  createGroup: vi.fn()
}))

vi.mock('@/api/admin/upstreamHub', () => ({
  default: {
    listVendors,
    listAllEntries,
    listAggregationUnits,
    setEntryGroups,
    setEntryGroupEnabled
  }
}))

vi.mock('@/api/admin/groups', () => ({
  default: {
    getAllIncludingInactive,
    update: updateGroup,
    delete: deleteGroup,
    listCompositeRoutes,
    createCompositeRoute,
    updateCompositeRoute,
    deleteCompositeRoute,
    getGroupApiKeys,
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

const routeQuery: { value: { query: Record<string, string> } } = { value: { query: {} } }
const replace = vi.fn()
vi.mock('vue-router', () => ({
  useRoute: () => routeQuery.value,
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
          props: ['show', 'title', 'width'],
          template: '<div v-if="show" class="dialog-stub"><slot /><slot name="footer" /></div>'
        },
        ConfirmDialog: {
          props: ['show', 'title', 'message'],
          emits: ['confirm', 'cancel'],
          template: '<button v-if="show" class="confirm-stub" @click="$emit(\'confirm\')">confirm</button>'
        },
        Select: {
          props: ['modelValue', 'options'],
          emits: ['update:modelValue'],
          template:
            '<select class="select-stub" :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)">' +
            '<option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option>' +
            '</select>'
        },
        Toggle: {
          props: ['modelValue', 'disabled'],
          emits: ['update:modelValue'],
          template:
            '<button class="toggle-stub" :disabled="disabled" @click="$emit(\'update:modelValue\', !modelValue)" />'
        }
      }
    }
  })

const activeRoute = {
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

function fixtures() {
  listAggregationUnits.mockResolvedValue([
    {
      id: 3,
      name: 'mock-composite',
      description: '',
      status: 'active',
      rpm_limit: 60,
      entry_count: 1,
      route_count: 1
    }
  ])
  listVendors.mockResolvedValue([
    { id: 2, name: 'MockVendor', concurrency_budget: 10, allocated_concurrency: 5, base_urls: ['https://mock.example'], over_budget: false }
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
      groups: [{ group_id: 3, group_name: 'mock-composite', platform: 'composite', priority: 50, enabled: true }],
      models: ['mock-deepseek-chat']
    },
    {
      account_id: 10,
      vendor_id: 2,
      name: 'Mock-OpenAI',
      platform: 'openai',
      base_url: 'http://127.0.0.1:18081/v1',
      has_api_key: true,
      concurrency: 5,
      priority: 50,
      schedulable: false,
      status: 'active',
      groups: [],
      models: []
    }
  ])
  listCompositeRoutes.mockResolvedValue([activeRoute])
  getGroupApiKeys.mockResolvedValue({
    items: [
      {
        id: 1,
        user_id: 1,
        key: 'sk-1234567890abcdef',
        name: 'customer-key',
        status: 'active',
        last_used_at: null,
        user: { id: 1, email: 'admin@admin.local' }
      }
    ],
    total: 1,
    page: 1,
    page_size: 100,
    pages: 1
  })
  setEntryGroups.mockResolvedValue([])
  setEntryGroupEnabled.mockResolvedValue(undefined)
  createCompositeRoute.mockResolvedValue({ id: 2 })
  updateCompositeRoute.mockResolvedValue(activeRoute)
  deleteCompositeRoute.mockResolvedValue({ message: 'ok' })
  createGroup.mockResolvedValue({ id: 7, name: 'new-unit', platform: 'composite' })
  updateGroup.mockResolvedValue({ id: 3, name: 'renamed-unit' })
  deleteGroup.mockResolvedValue({ message: 'ok' })
}

async function openUnit(wrapper: ReturnType<typeof mountView>) {
  const manageButton = wrapper.findAll('button').find((button) => button.text() === 'admin.aggregation.manage')
  expect(manageButton).toBeTruthy()
  await manageButton!.trigger('click')
  await flushPromises()
}

async function openTab(wrapper: ReturnType<typeof mountView>, key: string) {
  const tabButton = wrapper.findAll('button').find((button) => button.text() === key)
  expect(tabButton).toBeTruthy()
  await tabButton!.trigger('click')
  await flushPromises()
}

describe('AggregationView', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    routeQuery.value = { query: {} }
    fixtures()
  })

  it('lists only composite groups as scheduling units with attached counts', async () => {
    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.text()).toContain('mock-composite')
    expect(listAggregationUnits).toHaveBeenCalled()
  })

  it('manages a unit in a dialog and detaches an attached entry', async () => {
    const wrapper = mountView()
    await flushPromises()
    await openUnit(wrapper)

    expect(wrapper.find('.dialog-stub').exists()).toBe(true)
    expect(wrapper.text()).toContain('admin.aggregation.supply')

    const checkbox = wrapper.find('input[type="checkbox"]')
    await checkbox.trigger('change')
    await flushPromises()
    expect(setEntryGroups).toHaveBeenCalledWith(9, [])
  })

  it('attaches an entry and keeps the enable switch disabled until attached', async () => {
    const wrapper = mountView()
    await flushPromises()
    await openUnit(wrapper)

    const toggles = wrapper.findAll('.toggle-stub')
    expect(toggles).toHaveLength(2)
    const disabledToggle = toggles.find((toggle) => toggle.attributes('disabled') !== undefined)
    expect(disabledToggle).toBeTruthy()

    const attachCheckbox = wrapper.find('input[aria-label="Mock-OpenAI"]')
    await attachCheckbox.trigger('change')
    await flushPromises()
    expect(setEntryGroups).toHaveBeenCalledWith(10, [3])
  })

  it('shows real customer keys and active counts in the supply chain', async () => {
    const wrapper = mountView()
    await flushPromises()
    await openUnit(wrapper)
    await openTab(wrapper, 'admin.aggregation.tabChain')

    expect(getGroupApiKeys).toHaveBeenCalledWith(3, 1, 100)
    expect(wrapper.text()).toContain('customer-key')
    expect(wrapper.text()).toContain('admin@admin.local')
    expect(wrapper.text()).toContain('sk-123...cdef')
  })

  it('creates a route with the configured match type and endpoint', async () => {
    const wrapper = mountView()
    await flushPromises()
    await openUnit(wrapper)
    await openTab(wrapper, 'admin.aggregation.tabRoutes')

    const publicModelInput = wrapper.find('input[placeholder="admin.aggregation.publicModelPlaceholder"]')
    await publicModelInput.setValue('public-alias')
    const upstreamInput = wrapper.find('input[placeholder="admin.aggregation.upstreamModelPlaceholder"]')
    await upstreamInput.setValue('mock-deepseek-chat')
    const priorityInput = wrapper.find('input[type="number"]')
    await priorityInput.setValue('20')

    const addButton = wrapper
      .findAll('button')
      .filter((button) => button.text() === 'admin.aggregation.addRoute')
      .at(-1)
    await addButton!.trigger('click')
    await flushPromises()

    expect(createCompositeRoute).toHaveBeenCalledWith(
      3,
      expect.objectContaining({
        public_model: 'public-alias',
        target_platform: 'openai',
        upstream_model: 'mock-deepseek-chat',
        match_type: 'exact',
        endpoint: 'any',
        priority: 20
      })
    )
  })

  it('edits an existing route through the same form', async () => {
    const wrapper = mountView()
    await flushPromises()
    await openUnit(wrapper)
    await openTab(wrapper, 'admin.aggregation.tabRoutes')

    const editButton = wrapper.findAll('button').find((button) => button.text() === 'admin.aggregation.edit')
    await editButton!.trigger('click')
    await flushPromises()

    const publicModelInput = wrapper.find('input[placeholder="admin.aggregation.publicModelPlaceholder"]')
    expect((publicModelInput.element as HTMLInputElement).value).toBe('mock-public-model')

    const saveButton = wrapper.findAll('button').find((button) => button.text() === 'admin.aggregation.save')
    await saveButton!.trigger('click')
    await flushPromises()

    expect(updateCompositeRoute).toHaveBeenCalledWith(
      3,
      1,
      expect.objectContaining({ public_model: 'mock-public-model', upstream_model: 'mock-deepseek-chat' })
    )
  })

  it('opens the unit dialog from the ?unit= deep link', async () => {
    routeQuery.value = { query: { unit: '3' } }
    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.find('.dialog-stub').exists()).toBe(true)
  })

  it('renames the current unit through the groups API', async () => {
    const wrapper = mountView()
    await flushPromises()
    await openUnit(wrapper)

    const renameButton = wrapper.findAll('button').find((button) => button.text() === 'admin.aggregation.renameUnit')
    await renameButton!.trigger('click')
    await flushPromises()

    const input = wrapper.find('[data-testid="rename-unit-input"]')
    await input.setValue('renamed-unit')
    const saveButton = wrapper.findAll('button').find((button) => button.text() === 'admin.aggregation.save')
    await saveButton!.trigger('click')
    await flushPromises()

    expect(updateGroup).toHaveBeenCalledWith(3, { name: 'renamed-unit' })
  })

  it('deletes the current unit after confirmation', async () => {
    const wrapper = mountView()
    await flushPromises()
    await openUnit(wrapper)

    const deleteButton = wrapper.findAll('button').find((button) => button.text() === 'admin.aggregation.deleteUnit')
    await deleteButton!.trigger('click')
    await flushPromises()

    const confirmButton = wrapper.find('.confirm-stub')
    expect(confirmButton.exists()).toBe(true)
    await confirmButton.trigger('click')
    await flushPromises()

    expect(deleteGroup).toHaveBeenCalledWith(3)
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
