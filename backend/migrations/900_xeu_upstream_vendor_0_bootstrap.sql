-- 全新数据库的排序补丁：文件名字典序下 900_xeu_upstream_vendor_base_urls.sql
-- 会先于 900_xeu_upstream_vendors.sql 执行，导致 ALTER TABLE 找不到表。
-- 这里提前用同一份定义幂等建表，让 base_urls 迁移有表可改；随后
-- 900_xeu_upstream_vendors.sql 的 CREATE TABLE IF NOT EXISTS 自然成为空操作。
-- 表结构以 900_xeu_upstream_vendors.sql 为准，本文件不再追加新列。
CREATE TABLE IF NOT EXISTS upstream_vendors (
    id                 BIGSERIAL PRIMARY KEY,
    name               VARCHAR(100) NOT NULL,
    concurrency_budget INTEGER NOT NULL DEFAULT 0,
    notes              TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_upstream_vendors_name_active
    ON upstream_vendors (name)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_upstream_vendors_deleted_at
    ON upstream_vendors (deleted_at);
