CREATE TABLE IF NOT EXISTS executor_kv_store (
    id         SERIAL PRIMARY KEY,
    uuid       UUID NOT NULL DEFAULT uuid_generate_v4(),
    bucket     VARCHAR(255) NOT NULL,
    key        VARCHAR(255) NOT NULL,
    value      TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(bucket, key)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_executor_kv_uuid ON executor_kv_store(uuid);
CREATE INDEX IF NOT EXISTS idx_executor_kv_bucket ON executor_kv_store(bucket);
