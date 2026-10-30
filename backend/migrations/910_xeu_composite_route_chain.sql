-- 归一模型回退链：允许同一对外模型按优先级配置多条（跨平台）路由。
-- 旧唯一索引 (group_id, endpoint, match_type, public_model) 会拒绝第二行；
-- 改为非唯一索引，resolver 按 (match strength, endpoint, prefix, priority, id)
-- 取链，并在选路时挑选第一个有可用账号的平台。
DROP INDEX IF EXISTS idx_composite_model_routes_unique_active;

CREATE INDEX IF NOT EXISTS idx_composite_model_routes_chain
    ON composite_model_routes (group_id, endpoint, match_type, public_model, priority, id)
    WHERE deleted_at IS NULL;
