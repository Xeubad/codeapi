import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import GroupHubEntriesField from '../GroupHubEntriesField.vue'
import type { UpstreamEntry, UpstreamVendor } from '@/api/admin/upstreamHub'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => key })
  }
})

const vendors: UpstreamVendor[] = [
  { id: 1, name: 'Vendor A', concurrency_budget: 10, notes: '', base_urls: [], entry_count: 2, allocated_concurrency: 2, referenced_group_count: 0, over_budget: false, created_at: '', updated_at: '' }
]

const entries: UpstreamEntry[] = [
  {
    account_id: 11,
    vendor_id: 1,
    name: 'OpenAI-1',
    platform: 'openai',
    base_url: 'https://relay.example',
    has_api_key: true,
    concurrency: 1,
    priority: 50,
    schedulable: true,
    status: 'active',
    groups: [],
    models: ['gpt-5']
  },
  {
    account_id: 12,
    vendor_id: 1,
    name: 'Anthropic-1',
    platform: 'anthropic',
    base_url: 'https://relay.example',
    has_api_key: true,
    concurrency: 1,
    priority: 50,
    schedulable: true,
    status: 'active',
    groups: [],
    models: []
  },
  {
    account_id: 13,
    vendor_id: 1,
    name: 'OpenAI-2',
    platform: 'openai',
    base_url: 'https://relay.example',
    has_api_key: true,
    concurrency: 1,
    priority: 50,
    schedulable: false,
    status: 'active',
    groups: [],
    models: []
  }
]

const mountField = (platform: string, modelValue: number[] = []) =>
  mount(GroupHubEntriesField, {
    props: { platform, vendors, entries, modelValue }
  })

describe('GroupHubEntriesField', () => {
  it('lists only entries matching the group platform', () => {
    const wrapper = mountField('openai')
    expect(wrapper.text()).toContain('OpenAI-1')
    expect(wrapper.text()).toContain('OpenAI-2')
    expect(wrapper.text()).not.toContain('Anthropic-1')
  })

  it('lists entries from every platform for composite groups', () => {
    const wrapper = mountField('composite')
    expect(wrapper.text()).toContain('OpenAI-1')
    expect(wrapper.text()).toContain('Anthropic-1')
  })

  it('emits the next selection when toggling an entry', async () => {
    const wrapper = mountField('openai', [11])
    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    expect((checkboxes[0].element as HTMLInputElement).checked).toBe(true)

    await checkboxes[1].trigger('change')
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([[11, 13]])
  })

  it('flags entries without registered models and unschedulable entries', () => {
    const wrapper = mountField('openai')
    expect(wrapper.text()).toContain('admin.groups.hubRef.unschedulable')
  })
})
