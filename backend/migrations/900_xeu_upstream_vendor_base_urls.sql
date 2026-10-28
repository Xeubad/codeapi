-- 上游厂家候选接入地址：允许一个厂家配置多个 Base URL。
-- 条目保存时会逐个测速并只把最快的一个写进条目凭证，运行时不再切换，
-- 避免请求途中因自动择路导致的行为漂移。
ALTER TABLE upstream_vendors
    ADD COLUMN IF NOT EXISTS base_urls JSONB NOT NULL DEFAULT '[]'::jsonb;
