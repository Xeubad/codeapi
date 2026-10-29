import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import AccountActionMenu from '../AccountActionMenu.vue'
import type { Account } from '@/types'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    })
  }
})

function makeAccount(overrides: Partial<Account>): Account {
  return {
    id: 1,
    name: 'test-account',
    platform: 'openai',
    type: 'apikey',
    proxy_id: null,
    concurrency: 3,
    priority: 50,
    status: 'active',
    error_message: null,
    last_used_at: null,
    expires_at: null,
    auto_pause_on_expired: false,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    schedulable: true,
    rate_limited_at: null,
    rate_limit_reset_at: null,
    overload_until: null,
    temp_unschedulable_until: null,
    temp_unschedulable_reason: null,
    session_window_start: null,
    session_window_end: null,
    session_window_status: null,
    credentials: { base_url: 'https://relay.example' },
    credentials_status: { has_api_key: true },
    extra: {},
    ...overrides
  }
}

const anchorRect = new DOMRect(100, 100, 24, 24)
const getBodyButtons = () => Array.from(document.body.querySelectorAll('button'))
const findAttachButton = () => getBodyButtons().find((button) => button.textContent?.includes('admin.accounts.hub.attach'))

function mountMenu(account: Account) {
  return mount(AccountActionMenu, {
    props: { show: true, account, anchorRect },
    attachTo: document.body
  })
}

describe('AccountActionMenu — 纳入上游聚合', () => {
  it('shows the action for an eligible relay account and emits the account', async () => {
    const account = makeAccount({})
    const wrapper = mountMenu(account)

    const attachButton = findAttachButton()
    expect(attachButton).toBeTruthy()
    await attachButton!.click()
    expect(wrapper.emitted('attach-to-hub')?.[0]).toEqual([account])

    wrapper.unmount()
  })

  it('hides the action for accounts already owned by a vendor', () => {
    const wrapper = mountMenu(makeAccount({ extra: { upstream_vendor_id: 3 } }))
    expect(findAttachButton()).toBeFalsy()
    wrapper.unmount()
  })

  it('hides the action when the account has no custom base url', () => {
    const wrapper = mountMenu(makeAccount({ credentials: {} }))
    expect(findAttachButton()).toBeFalsy()
    wrapper.unmount()
  })
})
