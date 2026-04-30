package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cvhariharan/flowctl/internal/repo"
)

func (c *Core) SetExecutorKV(ctx context.Context, bucket, key, value string) error {
	if bucket == "" {
		return errors.New("bucket is required")
	}
	if key == "" {
		return errors.New("key is required")
	}
	_, err := c.store.UpsertExecutorKV(ctx, repo.UpsertExecutorKVParams{
		Bucket: bucket,
		Key:    key,
		Value:  value,
	})
	return err
}

func (c *Core) GetExecutorKV(ctx context.Context, bucket, key string) (string, error) {
	row, err := c.store.GetExecutorKV(ctx, repo.GetExecutorKVParams{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("key %q not found in bucket %q", key, bucket)
		}
		return "", err
	}
	return row.Value, nil
}

func (c *Core) ListExecutorKV(ctx context.Context, bucket string) (map[string]string, error) {
	rows, err := c.store.ListExecutorKVByBucket(ctx, bucket)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(rows))
	for _, r := range rows {
		result[r.Key] = r.Value
	}
	return result, nil
}

func (c *Core) DeleteExecutorKV(ctx context.Context, bucket, key string) error {
	return c.store.DeleteExecutorKV(ctx, repo.DeleteExecutorKVParams{
		Bucket: bucket,
		Key:    key,
	})
}
