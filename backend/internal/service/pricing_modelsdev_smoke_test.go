package service

import (
	"encoding/json"
	"os"
	"testing"
)

// TestModelsDevCatalogSmoke 用真实 models.dev api.json（存在时）跑完整解析管线，
// 验证单位换算、一方优先去重、图片模型映射与阶梯折算。CI 无该文件时跳过。
func TestModelsDevCatalogSmoke(t *testing.T) {
	body, err := os.ReadFile("/tmp/modelsdev/api.json")
	if err != nil {
		t.Skip("local api.json snapshot not available")
	}

	s := &PricingService{cfg: nil}
	data, err := s.parsePricingData(body)
	if err != nil {
		t.Fatalf("parsePricingData: %v", err)
	}
	t.Logf("parsed %d models from models.dev catalog", len(data))
	if len(data) < 1000 {
		t.Fatalf("expected a broad catalog, got %d models", len(data))
	}

	// 单位换算：claude-opus-5-5 官方 $4/M input → 4e-6 USD/token
	opus, ok := data["claude-opus-5-5"]
	if !ok {
		t.Fatalf("claude-opus-5-5 missing from converted catalog")
	}
	if !nearlyEqual(opus.InputCostPerToken, 4e-6) {
		t.Errorf("claude-opus-5-5 input: got %g, want 4e-6", opus.InputCostPerToken)
	}
	if !nearlyEqual(opus.OutputCostPerToken, 2e-5) {
		t.Errorf("claude-opus-5-5 output: got %g, want 2e-5", opus.OutputCostPerToken)
	}
	if !nearlyEqual(opus.CacheCreationInputTokenCost, 5e-6) {
		t.Errorf("claude-opus-5-5 cache_write: got %g, want 5e-6", opus.CacheCreationInputTokenCost)
	}

	// 一方优先：gpt-5.5 出现在 6+ 个厂商下，必须取 openai 的 $5/$30
	gpt55, ok := data["gpt-5.5"]
	if !ok {
		t.Fatalf("gpt-5.5 missing from converted catalog")
	}
	if !nearlyEqual(gpt55.InputCostPerToken, 5e-6) || !nearlyEqual(gpt55.OutputCostPerToken, 3e-5) {
		t.Errorf("gpt-5.5: got %g/%g, want 5e-6/3e-5", gpt55.InputCostPerToken, gpt55.OutputCostPerToken)
	}

	// 长上下文阶梯折算：gpt-5.5 的 context tier @272k → 阈值 272000、倍率 2x
	if gpt55.LongContextInputTokenThreshold != 272000 {
		t.Errorf("gpt-5.5 long-context threshold: got %d, want 272000", gpt55.LongContextInputTokenThreshold)
	}
	if !nearlyEqual(gpt55.LongContextInputCostMultiplier, 2.0) || !nearlyEqual(gpt55.LongContextOutputCostMultiplier, 1.5) {
		t.Errorf("gpt-5.5 long-context multipliers: got %g/%g, want 2.0/1.5",
			gpt55.LongContextInputCostMultiplier, gpt55.LongContextOutputCostMultiplier)
	}

	// 图片输出模型：gpt-image-2 的 output 是图片 token 价 $30/M；文本输出价
	// models.dev 未提供，不得落到 output_cost_per_token。
	img, ok := data["gpt-image-2"]
	if !ok {
		t.Fatalf("gpt-image-2 missing from converted catalog")
	}
	if !nearlyEqual(img.OutputCostPerImageToken, 3e-5) {
		t.Errorf("gpt-image-2 image token output: got %g, want 3e-5", img.OutputCostPerImageToken)
	}
	if img.OutputCostPerToken != 0 {
		t.Errorf("gpt-image-2 must not carry a text output price, got %g", img.OutputCostPerToken)
	}
	if !nearlyEqual(img.InputCostPerToken, 5e-6) {
		t.Errorf("gpt-image-2 text input: got %g, want 5e-6", img.InputCostPerToken)
	}
	if img.TokenPricingAbsent {
		t.Errorf("gpt-image-2 has a text input price, TokenPricingAbsent must be false")
	}

	// provider 名换算：google → gemini，供管理端按平台同步继续工作
	if p := data["gemini-3.1-pro-preview"]; p == nil || p.LiteLLMProvider != "gemini" {
		if p == nil {
			t.Errorf("gemini-3.1-pro-preview missing")
		} else {
			t.Errorf("gemini-3.1-pro-preview provider: got %q, want gemini", p.LiteLLMProvider)
		}
	}
}

// TestProviderQuotesSurviveFallbackMerge 验证 models.dev 索引在 fallback/override
// 合并层（LiteLLM 格式重入 parsePricingData）后不被清空。
func TestProviderQuotesSurviveFallbackMerge(t *testing.T) {
	s := &PricingService{}
	modelsDevBody := []byte(`{
		"openai": {"models": {"gpt-5.5": {"modalities": {"output": ["text"]}, "cost": {"input": 5, "output": 30}}}}
	}`)
	if !s.hasProviderQuotes() {
		if _, err := s.parsePricingData(modelsDevBody); err != nil {
			t.Fatalf("parse models.dev body: %v", err)
		}
	}
	if !s.hasProviderQuotes() {
		t.Fatalf("provider quotes index must be populated after models.dev parse")
	}
	// fallback 文件是 LiteLLM 扁平格式,重入 parsePricingData 不得清空索引
	litellmBody := []byte(`{"claude-opus-5-5": {"input_cost_per_token": 4e-06, "output_cost_per_token": 2e-05, "litellm_provider": "anthropic"}}`)
	if _, err := s.parsePricingData(litellmBody); err != nil {
		t.Fatalf("parse litellm body: %v", err)
	}
	if !s.hasProviderQuotes() {
		t.Fatalf("provider quotes index must survive litellm-format re-parse (fallback merge)")
	}
	quotes := s.ListModelProviderQuotes("gpt-5.5")
	if len(quotes) != 1 || quotes[0].ProviderID != "openai" {
		t.Fatalf("gpt-5.5 quotes: got %+v", quotes)
	}
}

// TestLiteLLMCatalogNotDetectedAsModelsDev 验证旧 LiteLLM 目录（扁平结构）
// 不会被误判为 models.dev 格式而走转换路径。
func TestLiteLLMCatalogNotDetectedAsModelsDev(t *testing.T) {
	litellmBody := []byte(`{
		"sample_spec": {"input_cost_per_token": 0, "litellm_provider": "sample"},
		"claude-opus-5-5": {"input_cost_per_token": 4e-06, "output_cost_per_token": 2e-05, "litellm_provider": "anthropic"}
	}`)
	if detectModelsDevFormat(litellmBody) {
		t.Fatalf("LiteLLM flat catalog must not be detected as models.dev format")
	}

	// 空对象、数组与损坏 JSON 一律不触发转换。
	if detectModelsDevFormat([]byte(`{}`)) {
		t.Errorf("empty object must not be detected as models.dev")
	}
	if detectModelsDevFormat([]byte(`[]`)) {
		t.Errorf("array must not be detected as models.dev")
	}
	if detectModelsDevFormat([]byte(`not json`)) {
		t.Errorf("invalid JSON must not be detected as models.dev")
	}

	// 带空 models 映射的顶层对象视为 LiteLLM 目录（条目模型各自解析）。
	if detectModelsDevFormat([]byte(`{"claude-opus-5-5": {"input_cost_per_token": 4e-06}}`)) {
		t.Errorf("flat entry without models field must not be detected as models.dev")
	}
}

// TestModelsDevConversionIsDeterministic 验证同名模型无一方厂商时按字典序取首个，
// 且转换结果可直接序列化。
func TestModelsDevConversionIsDeterministic(t *testing.T) {
	body := []byte(`{
		"zeta-aggregator": {"models": {"shared-model": {"modalities": {"output": ["text"]}, "cost": {"input": 2, "output": 6}}}},
		"alpha-aggregator": {"models": {"shared-model": {"modalities": {"output": ["text"]}, "cost": {"input": 1, "output": 3}}}}
	}`)
	out, quotes, err := convertModelsDevCatalog(body)
	if err != nil {
		t.Fatalf("convertModelsDevCatalog: %v", err)
	}
	if q, ok := quotes["shared-model"]; !ok || len(q.Quotes) != 2 {
		t.Fatalf("shared-model should have 2 provider quotes, got %+v", q)
	}
	if q := quotes["shared-model"]; q.Quotes[0].ProviderID != "alpha-aggregator" || len(q.Quotes) != 2 {
		t.Fatalf("winner must be first in quotes: %+v", q.Quotes)
	}
	var flat map[string]json.RawMessage
	if err := json.Unmarshal(out, &flat); err != nil {
		t.Fatalf("converted output is not flat JSON: %v", err)
	}
	var entry struct {
		Input float64 `json:"input_cost_per_token"`
	}
	if err := json.Unmarshal(flat["shared-model"], &entry); err != nil {
		t.Fatalf("shared-model entry: %v", err)
	}
	// 无一方厂商收录时取字典序首个（alpha-aggregator）
	if entry.Input != 1e-6 {
		t.Errorf("shared-model input: got %g, want 1e-6 (alpha-aggregator wins)", entry.Input)
	}
}

func nearlyEqual(got, want float64) bool {
	if got == want {
		return true
	}
	diff := got - want
	if diff < 0 {
		diff = -diff
	}
	return diff/1e-12 < 1e4
}
