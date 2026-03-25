package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type AppConfig struct {
	Version  int               `json:"version"`
	App      *AppCredentials   `json:"app,omitempty"`
	Auth     AuthState         `json:"auth"`
	Defaults DefaultSelections `json:"defaults"`
}

type AppCredentials struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURI  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
	UserAgent    string   `json:"user_agent"`
}

type AuthState struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	Scope        string `json:"scope,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	PendingState string `json:"pending_state,omitempty"`
}

type DefaultSelections struct {
	BusinessID  string `json:"business_id,omitempty"`
	AdAccountID string `json:"ad_account_id,omitempty"`
}

type Store struct {
	Path   string
	Config AppConfig
}

func Load() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".config", "radcli", "config.json")
	store := &Store{
		Path: path,
		Config: AppConfig{
			Version:  1,
			Auth:     AuthState{},
			Defaults: DefaultSelections{},
		},
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &store.Config); err != nil {
		return nil, err
	}
	if store.Config.Version == 0 {
		store.Config.Version = 1
	}
	return store, nil
}

func (s *Store) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.Config, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(s.Path, data, 0o600)
}

func (s *Store) SanitizedMap() map[string]any {
	out := map[string]any{
		"version": s.Config.Version,
		"auth": map[string]any{
			"has_access_token":  s.Config.Auth.AccessToken != "",
			"has_refresh_token": s.Config.Auth.RefreshToken != "",
			"expires_at":        emptyToNil(s.Config.Auth.ExpiresAt),
			"scope":             emptyToNil(s.Config.Auth.Scope),
		},
		"defaults": map[string]any{
			"business_id":   emptyToNil(s.Config.Defaults.BusinessID),
			"ad_account_id": emptyToNil(s.Config.Defaults.AdAccountID),
		},
	}

	if s.Config.App != nil {
		out["app"] = map[string]any{
			"client_id":           s.Config.App.ClientID,
			"redirect_uri":        s.Config.App.RedirectURI,
			"scopes":              s.Config.App.Scopes,
			"user_agent":          s.Config.App.UserAgent,
			"client_secret_saved": s.Config.App.ClientSecret != "",
		}
	} else {
		out["app"] = nil
	}

	return out
}

func emptyToNil(value string) any {
	if value == "" {
		return nil
	}
	return value
}
