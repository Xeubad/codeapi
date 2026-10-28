-- 上游聚合条目改用 type='apikey'。
-- 网关各转发路径按 account.type 分支取凭据（OpenAIGatewayService.GetAccessToken、
-- GatewayService.GetAccessToken 等 20+ 处 switch），'upstream' 只有 antigravity relay
-- 有专门处理，openai/anthropic/gemini/国产供应商都会落到 default 分支报
-- "unsupported account type"。apikey + 自定义 base_url 才是平台全覆盖的透传形态。
-- 条目归属只认 extra.upstream_vendor_id，与 type 无关，故此处仅改 type。
-- 只转换本 fork 上游聚合创建的条目，不触碰平台自身的 antigravity relay upstream 账号。
UPDATE accounts
SET type       = 'apikey',
    updated_at = NOW()
WHERE type = 'upstream'
  AND deleted_at IS NULL
  AND extra ->> 'upstream_vendor_id' IS NOT NULL;
