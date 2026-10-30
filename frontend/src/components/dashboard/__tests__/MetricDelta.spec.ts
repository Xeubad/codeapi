import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => key }),
  }
})

import MetricDelta from '../MetricDelta.vue'

describe('MetricDelta', () => {
  it('shows a relative percentage change', () => {
    const wrapper = mount(MetricDelta, {
      props: { current: 120, previous: 100 },
      global: { stubs: { Icon: true } },
    })
    expect(wrapper.text()).toContain('+20.0%')
    expect(wrapper.text()).toContain('dashboard.metrics.vsYesterday')
    expect(wrapper.find('span').classes()).toContain('text-emerald-600')
  })

  it('supports inverted sentiment for cost metrics', () => {
    const wrapper = mount(MetricDelta, {
      props: { current: 120, previous: 100, positiveIsGood: false },
      global: { stubs: { Icon: true } },
    })
    expect(wrapper.find('span').classes()).toContain('text-red-500')
  })

  it('shows percentage-point changes for rates', () => {
    const wrapper = mount(MetricDelta, {
      props: { current: 45, previous: 40, mode: 'points' },
      global: { stubs: { Icon: true } },
    })
    expect(wrapper.text()).toContain('+5.0')
    expect(wrapper.text()).toContain('dashboard.metrics.pointsShort')
  })

  it('falls back when there is no previous value', () => {
    const wrapper = mount(MetricDelta, {
      props: { current: 10, previous: null },
      global: { stubs: { Icon: true } },
    })
    expect(wrapper.text()).toBe('dashboard.metrics.noComparison')
  })
})
