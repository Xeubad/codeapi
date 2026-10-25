package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
)

// models.dev（https://github.com/anomalyco/models.dev，数据文件 https://models.dev/api.json）
// 目录适配层：把「厂商 → 模型 → cost」两层结构、USD/百万token 计价的原始数据拍平成
// LiteLLM 目录格式（扁平、USD/token），供现有解析/计费管线原样消费。
//
// 设计约束（详见 docs/WORK_LOG_CN.md 2026-10-25）：
//   - 只做格式转换，不重写计费链路；转换失败时回退 LiteLLM 解析路径；
//   - 一方厂商价格优先，同名模型转售商仅在无一方覆盖时补入；
//   - 转售商补入会降低计费覆盖面但官方价语义更纯，因此目录条目记录来源厂商，
//     litellm_provider 换算回平台同步所需的 LiteLLM provider 名。

// modelsDevProvider 优先级：一方官方厂商在前。同名模型被多个厂商收录时，
// 序列中靠前的厂商胜出；都不在序列内时按字典序取首个，保证结果确定。
var modelsDevFirstPartyPriority = []string{
	"anthropic", "openai", "google", "xai", "deepseek", "moonshotai",
	"alibaba", "minimax", "zhipuai", "meta", "mistral", "cohere",
	"amazon", "nvidia", "ai21", "perplexity",
}

// modelsDevProviderToLiteLLM 把 models.dev 厂商 id 换算成 LiteLLM provider 名，
// 供 ListModelNamesByProvider（管理端按平台同步目录模型）继续工作。
// 未知厂商原样返回小写 id。
var modelsDevProviderToLiteLLM = map[string]string{
	"google":     "gemini",
	"moonshotai": "moonshot",
	"zhipuai":    "zhipu",
}

// modelsDevPricingTier models.dev cost.tiers 单档长上下文阶梯。
type modelsDevPricingTier struct {
	Input      *float64              `json:"input"`
	Output     *float64              `json:"output"`
	CacheRead  *float64              `json:"cache_read"`
	CacheWrite *float64              `json:"cache_write"`
	Tier       modelsDevTierBoundary `json:"tier"`
}

type modelsDevTierBoundary struct {
	Type string `json:"type"`
	Size int64  `json:"size"`
}

// modelsDevCost models.dev 的 cost 对象，单位 USD/百万token。
// 未知字段（input_audio/output_audio/reasoning 等）暂不映射。
type modelsDevCost struct {
	Input      *float64               `json:"input"`
	Output     *float64               `json:"output"`
	CacheRead  *float64               `json:"cache_read"`
	CacheWrite *float64               `json:"cache_write"`
	Tiers      []modelsDevPricingTier `json:"tiers"`
}

// modelsDevModelEntry 模型条目；modalities.output 含 image 表示图片输出模型
// （output 即图片 token 价）。
type modelsDevModelEntry struct {
	ID         string              `json:"id"`
	Modalities modelsDevModalities `json:"modalities"`
	Cost       *modelsDevCost      `json:"cost"`
}

type modelsDevProviderEntry struct {
	ID     string                         `json:"id"`
	Name   string                         `json:"name"`
	Models map[string]modelsDevModelEntry `json:"models"`
}

// modelsDevCandidate 记录某模型条目在目录中的来源厂商，供一方优先去重。
type modelsDevCandidate struct {
	providerID string
	entry      modelsDevModelEntry
}

// detectModelsDevFormat 判断 body 是否为 models.dev 目录格式。
// LiteLLM 目录是扁平的「模型名 → 对象」，其条目不携带 id/models 字段；
// models.dev 顶层是「厂商 id → {id, name, models}」。抽样至多 8 个顶层键，
// 全部含非空 models 映射即判定为 models.dev 格式。
func detectModelsDevFormat(body []byte) bool {
	var top map[string]modelsDevProviderEntry
	if err := json.Unmarshal(body, &top); err != nil {
		return false
	}
	if len(top) == 0 {
		return false
	}
	checked := 0
	for _, prov := range top {
		if prov.Models == nil {
			return false
		}
		checked++
		if checked >= 8 {
			break
		}
	}
	return true
}

// convertModelsDevCatalog 把 models.dev 目录转换为 LiteLLM 扁平目录 JSON。
// 返回值可直接喂给 parsePricingData（与 LiteLLM 目录同一套有效性过滤、
// override 叠加与 above_XXXk 折算逻辑）。
//
// 单位换算：models.dev 全部价格字段为 USD/百万token，除以 1e6 得 USD/token。
func convertModelsDevCatalog(body []byte) ([]byte, error) {
	var top map[string]modelsDevProviderEntry
	if err := json.Unmarshal(body, &top); err != nil {
		return nil, fmt.Errorf("parse models.dev catalog: %w", err)
	}

	models := make(map[string][]modelsDevCandidate)
	for providerID, prov := range top {
		for modelID, entry := range prov.Models {
			if strings.TrimSpace(modelID) == "" {
				continue
			}
			models[modelID] = append(models[modelID], modelsDevCandidate{providerID: providerID, entry: entry})
		}
	}

	out := make(map[string]map[string]any, len(models))
	for modelID, candidates := range models {
		chosen := pickModelsDevCandidate(candidates)
		entry := buildLiteLLMEntryFromModelsDev(chosen.providerID, chosen.entry)
		if entry != nil {
			out[modelID] = entry
		}
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("models.dev catalog produced no usable entries")
	}
	return json.Marshal(out)
}

// pickModelsDevCandidate 按一方优先规则选出模型条目的胜出厂商。
func pickModelsDevCandidate(candidates []modelsDevCandidate) modelsDevCandidate {
	if len(candidates) == 1 {
		return candidates[0]
	}
	for _, want := range modelsDevFirstPartyPriority {
		for _, c := range candidates {
			if c.providerID == want {
				return c
			}
		}
	}
	// 无一方厂商收录（如 openrouter 风格的 deepseek/deepseek-chat）：取字典序
	// 首个，保证同名模型在数据更新后仍稳定归到同一厂商。
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.providerID < best.providerID {
			best = c
		}
	}
	return best
}

// buildLiteLLMEntryFromModelsDev 生成单条 LiteLLM 目录条目（JSON 字段形态）。
// 无任何 token/图片价格时返回 nil（parsePricingData 的有效性过滤同口径）。
func buildLiteLLMEntryFromModelsDev(providerID string, entry modelsDevModelEntry) map[string]any {
	if entry.Cost == nil {
		return nil
	}
	cost := entry.Cost
	out := make(map[string]any, 8)

	litellmProvider := strings.ToLower(providerID)
	if mapped, ok := modelsDevProviderToLiteLLM[litellmProvider]; ok {
		litellmProvider = mapped
	}
	out["litellm_provider"] = litellmProvider

	// 图片输出模型：models.dev 的 output 是图片 token 价（如 gpt-image-2
	// $30/M），映射到 LiteLLM output_cost_per_image_token；文本价缺失时
	// 不再填 input/output token 字段，避免把图片价当文本价计费。
	isImageModel := false
	for _, modality := range entry.Modalities.Output {
		if strings.EqualFold(strings.TrimSpace(modality), "image") {
			isImageModel = true
			break
		}
	}

	perToken := func(v *float64) any {
		if v == nil || *v <= 0 {
			return nil
		}
		return *v / 1e6
	}

	setIfNotNil := func(key string, v *float64) {
		if converted := perToken(v); converted != nil {
			out[key] = converted
		}
	}

	setIfNotNil("input_cost_per_token", cost.Input)
	if isImageModel {
		// 图片输出模型：models.dev 的 output 是图片输出 token 价（OpenAI 官方
		// 文档口径），文本输出价目录未提供。若把它映射到 output_cost_per_token，
		// token 计费路径会把图片价当文本价用，造成超收。
		setIfNotNil("output_cost_per_image_token", cost.Output)
		out["mode"] = "image_generation"
	} else {
		setIfNotNil("output_cost_per_token", cost.Output)
	}
	setIfNotNil("cache_creation_input_token_cost", cost.CacheWrite)
	setIfNotNil("cache_read_input_token_cost", cost.CacheRead)

	// 长上下文阶梯：取首个 context 型 tier 的绝对价折算成 LiteLLM 的
	// *_above_<N>k_tokens 字段名（N 为千 token 数），交由
	// deriveLongContextFromAboveTierFields 统一折算阈值+倍率
	// （cache 侧按 LiteLLM 约定跟随输入倍率，不单独映射）。
	for _, tier := range cost.Tiers {
		if tier.Tier.Type != "context" || tier.Tier.Size <= 0 {
			continue
		}
		if tier.Tier.Size%1000 != 0 {
			// LiteLLM 字段名只支持千粒度阈值；非整千数据跳过并告警。
			logger.LegacyPrintf("service.pricing", "[Pricing] models.dev tier size %d is not a multiple of 1000; skipped", tier.Tier.Size)
			break
		}
		thousands := tier.Tier.Size / 1000
		if converted := perToken(tier.Input); converted != nil {
			out[fmt.Sprintf("input_cost_per_token_above_%dk_tokens", thousands)] = converted
		}
		if converted := perToken(tier.Output); converted != nil {
			out[fmt.Sprintf("output_cost_per_token_above_%dk_tokens", thousands)] = converted
		}
		break
	}

	if _, ok := out["input_cost_per_token"]; !ok {
		if _, ok := out["output_cost_per_token"]; !ok {
			if _, ok := out["output_cost_per_image_token"]; !ok {
				return nil
			}
		}
	}
	return out
}
