-- 上游聚合：分组对账号的"本组调用"开关。
-- false = 暂不调用（仅本分组候选池剔除），不影响账号本身在其它分组/全局的调度状态。
ALTER TABLE account_groups
    ADD COLUMN IF NOT EXISTS enabled BOOLEAN NOT NULL DEFAULT TRUE;

CREATE INDEX IF NOT EXISTS idx_account_groups_group_enabled
    ON account_groups (group_id, enabled);
