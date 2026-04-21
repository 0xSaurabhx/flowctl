-- name: UpsertExecutorKV :one
INSERT INTO executor_kv_store (bucket, key, value)
VALUES ($1, $2, $3)
ON CONFLICT (bucket, key) DO UPDATE
    SET value = EXCLUDED.value,
        updated_at = NOW()
RETURNING id, uuid, bucket, key, value, created_at, updated_at;

-- name: GetExecutorKV :one
SELECT id, uuid, bucket, key, value, created_at, updated_at
FROM executor_kv_store
WHERE bucket = $1 AND key = $2;

-- name: ListExecutorKVByBucket :many
SELECT id, uuid, bucket, key, value, created_at, updated_at
FROM executor_kv_store
WHERE bucket = $1
ORDER BY key;

-- name: DeleteExecutorKV :exec
DELETE FROM executor_kv_store
WHERE bucket = $1 AND key = $2;

-- name: DeleteExecutorKVByBucket :exec
DELETE FROM executor_kv_store
WHERE bucket = $1;
