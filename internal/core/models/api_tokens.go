package models

type APIToken struct {
	UUID       string `json:"uuid"`
	Label      string `json:"label"`
	LastUsedAt string `json:"last_used_at,omitempty"`
	CreatedAt  string `json:"created_at"`
}

type APITokenCreated struct {
	APIToken
	Token string `json:"token"`
}
