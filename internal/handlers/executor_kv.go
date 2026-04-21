package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleSetExecutorKV(c echo.Context) error {
	bucket := c.Param("bucket")
	if bucket == "" {
		return wrapError(ErrRequiredFieldMissing, "bucket is required", nil, nil)
	}

	var req KVEntry
	if err := c.Bind(&req); err != nil {
		return wrapError(ErrInvalidInput, "could not decode request", err, nil)
	}

	if req.Key == "" {
		return wrapError(ErrRequiredFieldMissing, "key is required", nil, nil)
	}

	if err := h.co.SetExecutorKV(c.Request().Context(), bucket, req.Key, req.Value); err != nil {
		return wrapError(ErrOperationFailed, "could not set kv entry", err, nil)
	}

	return c.JSON(http.StatusOK, KVEntry{Key: req.Key, Value: req.Value})
}

func (h *Handler) HandleGetExecutorKV(c echo.Context) error {
	bucket := c.Param("bucket")
	key := c.Param("key")
	if bucket == "" || key == "" {
		return wrapError(ErrRequiredFieldMissing, "bucket and key are required", nil, nil)
	}

	value, err := h.co.GetExecutorKV(c.Request().Context(), bucket, key)
	if err != nil {
		return wrapError(ErrResourceNotFound, "kv entry not found", err, nil)
	}

	return c.JSON(http.StatusOK, KVEntry{Key: key, Value: value})
}

func (h *Handler) HandleListExecutorKV(c echo.Context) error {
	bucket := c.Param("bucket")
	if bucket == "" {
		return wrapError(ErrRequiredFieldMissing, "bucket is required", nil, nil)
	}

	entries, err := h.co.ListExecutorKV(c.Request().Context(), bucket)
	if err != nil {
		return wrapError(ErrOperationFailed, "could not list kv entries", err, nil)
	}

	result := make([]KVEntry, 0, len(entries))
	for k, v := range entries {
		result = append(result, KVEntry{Key: k, Value: v})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) HandleDeleteExecutorKV(c echo.Context) error {
	bucket := c.Param("bucket")
	key := c.Param("key")
	if bucket == "" || key == "" {
		return wrapError(ErrRequiredFieldMissing, "bucket and key are required", nil, nil)
	}

	if err := h.co.DeleteExecutorKV(c.Request().Context(), bucket, key); err != nil {
		return wrapError(ErrOperationFailed, "could not delete kv entry", err, nil)
	}

	return c.NoContent(http.StatusOK)
}
