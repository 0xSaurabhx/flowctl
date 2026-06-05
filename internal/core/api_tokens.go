package core

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cvhariharan/flowctl/internal/core/models"
	"github.com/cvhariharan/flowctl/internal/repo"
	"github.com/google/uuid"
)

const (
	APITokenPrefix     = "fctl_pat_"
	apiTokenByteLength = 32
)

func hashAPIToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (c *Core) CreateAPIToken(ctx context.Context, userUUID, label string) (string, models.APIToken, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return "", models.APIToken{}, fmt.Errorf("label is required")
	}

	uid, err := uuid.Parse(userUUID)
	if err != nil {
		return "", models.APIToken{}, fmt.Errorf("user ID should be a UUID: %w", err)
	}

	buf := make([]byte, apiTokenByteLength)
	if _, err := rand.Read(buf); err != nil {
		return "", models.APIToken{}, fmt.Errorf("failed to generate token: %w", err)
	}
	raw := APITokenPrefix + base64.RawURLEncoding.EncodeToString(buf)

	row, err := c.store.CreateAPIToken(ctx, repo.CreateAPITokenParams{
		UserUuid:  uid,
		Label:     label,
		TokenHash: hashAPIToken(raw),
	})
	if err != nil {
		return "", models.APIToken{}, fmt.Errorf("could not create api token: %w", err)
	}

	return raw, repoAPITokenToModel(row), nil
}

func (c *Core) ListAPITokens(ctx context.Context, userUUID string) ([]models.APIToken, error) {
	uid, err := uuid.Parse(userUUID)
	if err != nil {
		return nil, fmt.Errorf("user ID should be a UUID: %w", err)
	}

	rows, err := c.store.ListAPITokensByUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("could not list api tokens: %w", err)
	}

	out := make([]models.APIToken, 0, len(rows))
	for _, r := range rows {
		out = append(out, repoAPITokenToModel(r))
	}
	return out, nil
}

func (c *Core) RevokeAPIToken(ctx context.Context, userUUID, tokenUUID string) error {
	uid, err := uuid.Parse(userUUID)
	if err != nil {
		return fmt.Errorf("user ID should be a UUID: %w", err)
	}
	tid, err := uuid.Parse(tokenUUID)
	if err != nil {
		return fmt.Errorf("token ID should be a UUID: %w", err)
	}

	if err := c.store.DeleteAPIToken(ctx, repo.DeleteAPITokenParams{
		TokenUuid: tid,
		UserUuid:  uid,
	}); err != nil {
		return fmt.Errorf("could not delete api token: %w", err)
	}
	return nil
}

// AuthenticateAPIToken validates a raw `fctl_pat_*` token and returns the owning
// user's UserInfo. The raw token is never logged or persisted.
func (c *Core) AuthenticateAPIToken(ctx context.Context, raw string) (models.UserInfo, error) {
	if !strings.HasPrefix(raw, APITokenPrefix) {
		return models.UserInfo{}, fmt.Errorf("invalid token prefix")
	}

	row, err := c.store.GetAPITokenByHash(ctx, hashAPIToken(raw))
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("invalid api token")
	}

	uwg, err := c.GetUserWithUUIDWithGroups(ctx, row.UserUuid.String())
	if err != nil {
		return models.UserInfo{}, fmt.Errorf("could not load token owner: %w", err)
	}

	// Best-effort last-used update; never block the request on it.
	tokenID := row.ID
	go func() {
		_ = c.store.TouchAPITokenLastUsed(context.Background(), tokenID)
	}()

	return uwg.ToUserInfo(), nil
}

func repoAPITokenToModel(r repo.ApiToken) models.APIToken {
	m := models.APIToken{
		UUID:      r.Uuid.String(),
		Label:     r.Label,
		CreatedAt: r.CreatedAt.Format(TimeFormat),
	}
	if r.LastUsedAt.Valid {
		m.LastUsedAt = r.LastUsedAt.Time.Format(TimeFormat)
	}
	return m
}
