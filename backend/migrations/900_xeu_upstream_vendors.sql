-- 上游聚合：厂家表。同厂家多把 key 共享账号级并发预算；条目归属存 accounts.extra.upstream_vendor_id。
CREATE TABLE IF NOT EXISTS upstream_vendors (
    id                 BIGSERIAL PRIMARY KEY,
    name               VARCHAR(100) NOT NULL,
    concurrency_budget INTEGER NOT NULL DEFAULT 0,
    notes              TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ
);

-- 软删除后可重用名称（与 groups 的部分唯一索引同型）。
CREATE UNIQUE INDEX IF NOT EXISTS idx_upstream_vendors_name_active
    ON upstream_vendors (name)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_upstream_vendors_deleted_at
    ON upstream_vendors (deleted_at);
