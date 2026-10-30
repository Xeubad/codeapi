import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

const { list, markRead } = vi.hoisted(() => ({ list: vi.fn(), markRead: vi.fn() }))
vi.mock('@/api', () => ({ announcementsAPI: { list, markRead } }))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => key }),
  }
})

import UserDashboardAnnouncements from '../UserDashboardAnnouncements.vue'
import type { UserAnnouncement } from '@/types'

function notice(id: number, title: string, createdAt: string, content = 'body'): UserAnnouncement {
  return {
    id,
    title,
    content,
    notify_mode: 'silent',
    created_at: createdAt,
    updated_at: createdAt,
  }
}

function mountComponent() {
  const pinia = createPinia()
  setActivePinia(pinia)
  return mount(UserDashboardAnnouncements, { global: { plugins: [pinia] } })
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

afterEach(() => {
  vi.restoreAllMocks()
})

describe('UserDashboardAnnouncements', () => {
  it('groups announcements by date with the newest date first', async () => {
    list.mockResolvedValue([
      notice(1, 'Older', '2026-10-01T09:00:00'),
      notice(2, 'Newest', '2026-10-29T09:00:00', '## Newest body\n\n<script>window.__xss = true</script>'),
      notice(3, 'Same day later', '2026-10-01T18:00:00'),
    ])

    const wrapper = mountComponent()
    await flushPromises()

    expect(list).toHaveBeenCalledWith(false)
    const groups = wrapper.findAll('[data-date]')
    expect(groups.map(group => group.attributes('data-date'))).toEqual(['2026-10-29', '2026-10-01'])
    expect(groups[1].findAll('[data-announcement-id]').map(node => node.attributes('data-announcement-id')))
      .toEqual(['3', '1'])

    const newestContent = wrapper.get('[data-announcement-id="2"] .markdown-body').html()
    expect(newestContent).toContain('<h2>Newest body</h2>')
    expect(newestContent).not.toContain('<script')
  })

  it('expands an announcement and marks it as read', async () => {
    markRead.mockResolvedValue({ message: 'ok' })
    list.mockResolvedValue([
      notice(2, 'Newest', '2026-10-29T09:00:00', 'Longer body'),
    ])

    const wrapper = mountComponent()
    await flushPromises()

    const article = wrapper.get('[data-announcement-id="2"]')
    expect(article.get('.markdown-body').classes()).toContain('line-clamp-2')
    expect(wrapper.text()).toContain('dashboard.announcements.unreadCount')

    const expand = article.findAll('button')
      .find(button => button.text() === 'dashboard.announcements.expand')
    expect(expand).toBeTruthy()

    await expand!.trigger('click')
    await flushPromises()

    expect(markRead).toHaveBeenCalledWith(2)
    expect(article.get('.markdown-body').classes()).not.toContain('line-clamp-2')
    expect(article.findAll('button').some(button => button.text() === 'dashboard.announcements.collapse')).toBe(true)
  })

  it('shows the empty state when there are no announcements', async () => {
    list.mockResolvedValue([])

    const wrapper = mountComponent()
    await flushPromises()

    expect(wrapper.findAll('[data-date]')).toHaveLength(0)
    expect(wrapper.text()).toContain('dashboard.announcements.empty')
  })
})
