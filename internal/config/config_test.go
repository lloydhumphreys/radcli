package config

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadAndSaveRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	store, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expectedPath := filepath.Join(home, ".config", "radcli", "config.json")
	if store.Path != expectedPath {
		t.Fatalf("store.Path = %q, want %q", store.Path, expectedPath)
	}

	store.Config.App = &AppCredentials{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURI:  "https://example.com/callback",
		Scopes:       []string{"adsread", "adsedit"},
		UserAgent:    "radcli-test",
	}
	store.Config.Auth = AuthState{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		TokenType:    "bearer",
		Scope:        "adsread adsedit",
		ExpiresAt:    "2026-03-26T12:00:00Z",
		PendingState: "state",
	}
	store.Config.Defaults = DefaultSelections{
		BusinessID:  "business-id",
		AdAccountID: "account-id",
	}

	if err := store.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	reloaded, err := Load()
	if err != nil {
		t.Fatalf("Load() after Save() error = %v", err)
	}

	if !reflect.DeepEqual(reloaded.Config, store.Config) {
		t.Fatalf("reloaded config mismatch\n got: %#v\nwant: %#v", reloaded.Config, store.Config)
	}
}

func TestSanitizedMapRedactsSecrets(t *testing.T) {
	store := &Store{
		Config: AppConfig{
			Version: 1,
			App: &AppCredentials{
				ClientID:     "client-id",
				ClientSecret: "very-secret",
				RedirectURI:  "https://example.com/callback",
				Scopes:       []string{"adsread"},
				UserAgent:    "radcli-test",
			},
			Auth: AuthState{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresAt:    "2026-03-26T12:00:00Z",
				Scope:        "adsread",
			},
			Defaults: DefaultSelections{
				BusinessID:  "business-id",
				AdAccountID: "account-id",
			},
		},
	}

	got := store.SanitizedMap()

	appMap, ok := got["app"].(map[string]any)
	if !ok {
		t.Fatalf("sanitized app config missing or wrong type: %#v", got["app"])
	}
	if _, exists := appMap["client_secret"]; exists {
		t.Fatalf("sanitized app config unexpectedly included client_secret")
	}
	if appMap["client_secret_saved"] != true {
		t.Fatalf("client_secret_saved = %#v, want true", appMap["client_secret_saved"])
	}

	authMap, ok := got["auth"].(map[string]any)
	if !ok {
		t.Fatalf("sanitized auth config missing or wrong type: %#v", got["auth"])
	}
	if authMap["has_access_token"] != true {
		t.Fatalf("has_access_token = %#v, want true", authMap["has_access_token"])
	}
	if authMap["has_refresh_token"] != true {
		t.Fatalf("has_refresh_token = %#v, want true", authMap["has_refresh_token"])
	}
	if _, exists := authMap["access_token"]; exists {
		t.Fatalf("sanitized auth config unexpectedly included access_token")
	}
	if _, exists := authMap["refresh_token"]; exists {
		t.Fatalf("sanitized auth config unexpectedly included refresh_token")
	}
}
