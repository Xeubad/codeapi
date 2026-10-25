import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import PlazaModelDialog from '../PlazaModelDialog.vue'
import type { PlazaModel } from '@/api/modelPlaza'

const adminState = { isAdmin: true }
vi.mock('@/stores/auth', () => ({
  useAuthStore: () => adminState
}))

const quoteMock = vi.fn(async () => ({
  model: 'claude-opus-5-5',
  quotes: [
    { provider_id: 'anthropic', provider_name: 'Anthropic', input_price: 4, output_price: 20, cache_read_price: 0.2, cache_write_price: 5 },
    { provider_id: '302ai', provider_name: '302.AI', input_price: 4, output_price: 22, cache_read_price: null, cache_write_price: null },
    { provider_id: 'azure', provider_name: 'Azure', input_price: 3, output_price: 18, cache_read_price: 0.2, cache_write_price: 5 }
  ]
}))
vi.mock('@/api/admin/channels', () => ({
  getModelProviderQuotes: (...args: unknown[]) => quoteMock(...(args as []))
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, unknown>) => {
        if (params) {
          return Object.values(params).join(' ')
        }
        return key
      }
    })
  }
})

const tokenModel: PlazaModel = {
  name: 'claude-opus-5-5',
  platform: 'anthropic',
  pricing: {
    billing_mode: 'token',
    input_price: 4e-6,
    output_price: 2e-5,
    cache_write_price: 5e-6,
    cache_write_1h_price: 8e-6,
    cache_read_price: 2e-7,
    image_input_price: null,
    image_output_price: null,
    per_request_price: null,
    intervals: []
  },
  official_pricing: {
    input_price: 4e-6,
    output_price: 2e-5,
    cache_write_price: 5e-6,
    cache_read_price: 2e-7
  }
}

function mountDialog(model: PlazaModel | null, props: Record<string, unknown> = {}) {
  return mount(PlazaModelDialog, {
    props: {
      show: model != null,
      model,
      rateMultiplier: 1.5,
      ...props
    },
    global: {
      stubs: {
        Teleport: true
      }
    }
  })
}

describe('PlazaModelDialog', () => {
  it('实付区按倍率合成、官方区不乘倍率', () => {
    const wrapper = mountDialog(tokenModel)
    const text = wrapper.text()
    // 实付: 4e-6×1.5=6, 2e-5×1.5=30 ($/1M)
    expect(text).toContain('×1.5')
    expect(text).toContain('6')
    expect(text).toContain('30')
  })

  it('缓存写入的 1h 价单独标注', () => {
    const wrapper = mountDialog(tokenModel)
    const text = wrapper.text()
    // 实付 5m: 5e-6×1.5=7.5; 1h: 8e-6×1.5=12
    expect(text).toContain('7.5')
    expect(text).toContain('12')
  })

  it('官方参考价不乘倍率', () => {
    const wrapper = mountDialog(tokenModel)
    const text = wrapper.text()
    // 官方 4 / 20 / 5 ($/1M) 直接出现
    expect(text).toContain('4')
    expect(text).toContain('20')
  })

  it('无官方价时显示占位文案', () => {
    const wrapper = mountDialog({ ...tokenModel, official_pricing: null })
    expect(wrapper.text()).toContain('modelPlaza.dialog.noOfficial')
  })

  it('阶梯模型按档位逐行展开实付价', () => {
    const m: PlazaModel = {
      ...tokenModel,
      pricing: {
        ...tokenModel.pricing!,
        intervals: [
          { min_tokens: 0, max_tokens: 272000, tier_label: '', input_price: 5e-6, output_price: 3e-5, cache_write_price: null, cache_read_price: null, per_request_price: null },
          { min_tokens: 272000, max_tokens: null, tier_label: '', input_price: 1e-5, output_price: 4.5e-5, cache_write_price: null, cache_read_price: null, per_request_price: null }
        ]
      },
      long_context_basis: 'whole_request'
    }
    const wrapper = mountDialog(m)
    const text = wrapper.text()
    // 首档 5e-6×1.5=7.5; 高档 1e-5×1.5=15
    expect(text).toContain('7.5')
    expect(text).toContain('15')
  })

  it('分时倍率模型展示时段行与生效倍率', () => {
    const m: PlazaModel = {
      ...tokenModel,
      time_pricing: {
        timezone: 'Asia/Shanghai',
        periods: [{ start_time: '00:30:00', end_time: '08:30:00', multiplier: 0.5 }]
      }
    }
    const wrapper = mountDialog(m)
    const text = wrapper.text()
    expect(text).toContain('00:30–08:30')
    expect(text).toContain('×0.75') // 1.5 × 0.5
  })

  it('管理员登录时展示厂商比价区,胜出厂商排最前', async () => {
    const wrapper = mountDialog(tokenModel)
    await vi.waitFor(() => {
      expect(wrapper.text()).toContain('modelPlaza.dialog.providersTitle')
    })
    const text = wrapper.text()
    expect(quoteMock).toHaveBeenCalledWith('claude-opus-5-5')
    expect(text).toContain('Anthropic')
    expect(text).toContain('modelPlaza.dialog.adoptedProvider')
    // 最低输入价厂商(azure 3)标最低价徽章
    expect(text).toContain('Azure')
    expect(text).toContain('modelPlaza.dialog.cheapest')
  })

  it('非管理员不展示厂商比价区', async () => {
    adminState.isAdmin = false
    quoteMock.mockClear()
    try {
      const wrapper = mountDialog(tokenModel)
      await vi.waitFor(() => {
        expect(wrapper.text()).not.toContain('modelPlaza.dialog.providersTitle')
      })
      expect(quoteMock).not.toHaveBeenCalled()
    } finally {
      adminState.isAdmin = true
    }
  })
})
