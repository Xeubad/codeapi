import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import type { TrendDataPoint } from '@/types'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => key }),
  }
})

import TokenUsageHeatmap from '../TokenUsageHeatmap.vue'

function point(date: string, totalTokens: number, requests = 1): TrendDataPoint {
  return {
    date,
    requests,
    input_tokens: totalTokens,
    output_tokens: 0,
    cache_creation_tokens: 0,
    cache_read_tokens: 0,
    total_tokens: totalTokens,
    cost: 0,
    actual_cost: 0,
  }
}

const trend = [
  point('2026-09-01', 200),
  point('2026-09-02', 400),
  point('2026-09-03', 600),
  point('2026-09-04', 800),
  point('2026-10-01', 100),
  point('2026-10-02', 300),
  point('2026-10-03', 500),
  point('2026-10-04', 700),
]

describe('TokenUsageHeatmap', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date(2026, 9, 15, 12, 0, 0))
  })

  afterEach(() => {
    vi.useRealTimers()
    document.body.innerHTML = ''
  })

  it('renders one row per month, newest first, and fills missing days with zero', () => {
    const wrapper = mount(TokenUsageHeatmap, {
      props: { trend, loading: false, defaultView: 'month' },
    })

    const rows = wrapper.findAll('[data-month]')
    expect(rows).toHaveLength(12)
    expect(rows[0].attributes('data-month')).toBe('2026-10')
    expect(rows[1].attributes('data-month')).toBe('2026-09')
    expect(rows[11].attributes('data-month')).toBe('2025-11')
    expect(rows[0].findAll('[data-date]')).toHaveLength(31)
    expect(rows[0].attributes('data-total')).toBe('1600')
  })

  it('maps token volume to intensity levels and stamps raw values on cells', () => {
    const wrapper = mount(TokenUsageHeatmap, {
      props: { trend, loading: false, defaultView: 'month' },
    })

    expect(wrapper.get('[data-date="2026-10-01"]').attributes('data-level')).toBe('1')
    expect(wrapper.get('[data-date="2026-10-03"]').attributes('data-level')).toBe('2')
    expect(wrapper.get('[data-date="2026-10-04"]').attributes('data-level')).toBe('3')
    expect(wrapper.get('[data-date="2026-09-04"]').attributes('data-level')).toBe('4')
    expect(wrapper.get('[data-date="2026-10-05"]').attributes('data-tokens')).toBe('0')
    expect(wrapper.get('[data-date="2026-10-05"]').attributes('data-level')).toBe('0')
    expect(wrapper.get('[data-date="2026-10-15"]').attributes('data-future')).toBeUndefined()
    expect(wrapper.get('[data-date="2026-10-16"]').attributes('data-future')).toBe('true')
    expect(wrapper.get('[data-date="2026-10-16"]').classes()).not.toContain('bg-gray-100')
  })

  it('renders the year view by default and can switch to the month view', async () => {
    const wrapper = mount(TokenUsageHeatmap, {
      props: { trend, loading: false },
    })

    expect(wrapper.attributes('data-heatmap-view')).toBe('year')
    expect(wrapper.findAll('[data-month]')).toHaveLength(0)
    expect(wrapper.findAll('[data-date]').length).toBeGreaterThanOrEqual(340)

    const monthButton = wrapper.findAll('button')
      .find(button => button.text() === 'dashboard.tokenHeatmap.viewMonth')
    expect(monthButton).toBeTruthy()
    await monthButton!.trigger('click')

    expect(wrapper.attributes('data-heatmap-view')).toBe('month')
    expect(wrapper.findAll('[data-month]')).toHaveLength(12)
  })

  it('shows a detailed tooltip when hovering a day', async () => {
    const wrapper = mount(TokenUsageHeatmap, {
      props: { trend, loading: false, defaultView: 'month' },
    })

    await wrapper.get('[data-date="2026-10-04"]').trigger('mouseenter')

    expect(document.body.textContent).toContain('dashboard.tokenHeatmap.tooltip.requests')
    expect(document.body.textContent).toContain('dashboard.tokenHeatmap.tooltip.cache5m')
    expect(document.body.textContent).toContain('dashboard.tokenHeatmap.tooltip.cache1h')
    expect(document.body.textContent).toContain('700')
    expect(document.body.textContent).toContain('$0.0000')
    wrapper.unmount()
  })

  it('shows an empty hint when there is no token usage', () => {
    const wrapper = mount(TokenUsageHeatmap, {
      props: { trend: [], loading: false, defaultView: 'month' },
    })

    expect(wrapper.text()).toContain('dashboard.tokenHeatmap.empty')
    expect(wrapper.findAll('[data-month]')).toHaveLength(12)
  })

  it('hides the grid while loading', () => {
    const wrapper = mount(TokenUsageHeatmap, {
      props: { trend, loading: true, defaultView: 'month' },
    })

    expect(wrapper.findAll('[data-month]')).toHaveLength(0)
  })
})
