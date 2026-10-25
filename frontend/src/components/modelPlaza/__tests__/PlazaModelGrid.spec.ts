import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import PlazaModelGrid from '../PlazaModelGrid.vue'
import type { PlazaModel } from '@/api/modelPlaza'

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

function tokenModel(overrides: Partial<PlazaModel> = {}): PlazaModel {
  return {
    name: 'gpt-5.5',
    platform: 'openai',
    pricing: {
      billing_mode: 'token',
      input_price: 5e-6,
      output_price: 3e-5,
      cache_write_price: 5e-7,
      cache_write_1h_price: null,
      cache_read_price: 5e-8,
      image_input_price: null,
      image_output_price: null,
      per_request_price: null,
      intervals: []
    },
    official_pricing: {
      input_price: 5e-6,
      output_price: 3e-5,
      cache_write_price: 5e-7,
      cache_read_price: 5e-8
    },
    ...overrides
  }
}

function mountGrid(models: PlazaModel[], props: Record<string, unknown> = {}) {
  return mount(PlazaModelGrid, {
    props: {
      models,
      platform: 'openai',
      rateMultiplier: 2,
      ...props
    }
  })
}

describe('PlazaModelGrid', () => {
  it('块上展示按倍率合成的实付输入/输出价', () => {
    const wrapper = mountGrid([tokenModel()])
    const text = wrapper.text()
    // input 5e-6 ×2 = 10 ($/1M), output 3e-5 ×2 = 60 ($/1M)
    expect(text).toContain('10')
    expect(text).toContain('60')
    expect(text).toContain('×2')
  })

  it('专属倍率覆盖分组倍率', () => {
    const wrapper = mountGrid([tokenModel()], { userRateMultiplier: 0.5 })
    const text = wrapper.text()
    // input 5e-6 ×0.5 = 2.5 ($/1M)
    expect(text).toContain('×0.5')
    expect(text).toContain('2.5')
  })

  it('点击块时派发 select 事件携带模型', async () => {
    const wrapper = mountGrid([tokenModel()])
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('select')).toHaveLength(1)
    expect(wrapper.emitted('select')![0][0]).toMatchObject({ name: 'gpt-5.5' })
  })

  it('长上下文阶梯模型带阶梯徽章', () => {
    const m = tokenModel({
      pricing: {
        ...tokenModel().pricing!,
        intervals: [
          { min_tokens: 0, max_tokens: 272000, tier_label: '', input_price: 5e-6, output_price: 3e-5, cache_write_price: null, cache_read_price: null, per_request_price: null },
          { min_tokens: 272000, max_tokens: null, tier_label: '', input_price: 1e-5, output_price: 4.5e-5, cache_write_price: null, cache_read_price: null, per_request_price: null }
        ]
      }
    })
    const wrapper = mountGrid([m])
    expect(wrapper.text()).toContain('modelPlaza.badges.longContext')
  })

  it('未配置定价的模型显示占位文案', () => {
    const wrapper = mountGrid([tokenModel({ pricing: null })])
    expect(wrapper.text()).toContain('modelPlaza.detail.noPricing')
  })

  it('按次模型展示按次价而非 token 价', () => {
    const m = tokenModel({
      name: 'img-model',
      pricing: {
        billing_mode: 'per_request',
        input_price: null,
        output_price: null,
        cache_write_price: null,
        cache_write_1h_price: null,
        cache_read_price: null,
        image_input_price: null,
        image_output_price: null,
        per_request_price: 0.04,
        intervals: []
      }
    })
    const wrapper = mountGrid([m])
    expect(wrapper.text()).toContain('0.08') // 0.04 ×2
    expect(wrapper.text()).toContain('modelPlaza.table.perUnitRequest')
  })
})
