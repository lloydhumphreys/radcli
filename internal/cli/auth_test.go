package cli

import "testing"

func TestParseAuthorizationInputURL(t *testing.T) {
	code, state, err := parseAuthorizationInput("https://example.com/callback?state=abc123&code=def456#_")
	if err != nil {
		t.Fatalf("parseAuthorizationInput() error = %v", err)
	}
	if code != "def456" {
		t.Fatalf("code = %q, want %q", code, "def456")
	}
	if state != "abc123" {
		t.Fatalf("state = %q, want %q", state, "abc123")
	}
}

func TestResolveAuthorizationCompletionRequiresStateWhenPending(t *testing.T) {
	_, _, err := resolveAuthorizationCompletion("def456", "", "", "abc123")
	if err == nil {
		t.Fatal("resolveAuthorizationCompletion() error = nil, want pending state error")
	}
}

func TestResolveAuthorizationCompletionAcceptsCallbackURLWhenPending(t *testing.T) {
	code, state, err := resolveAuthorizationCompletion("", "", "https://example.com/callback?state=abc123&code=def456#_", "abc123")
	if err != nil {
		t.Fatalf("resolveAuthorizationCompletion() error = %v", err)
	}
	if code != "def456" {
		t.Fatalf("code = %q, want %q", code, "def456")
	}
	if state != "abc123" {
		t.Fatalf("state = %q, want %q", state, "abc123")
	}
}

func TestResolveAuthorizationCompletionRejectsMismatchedState(t *testing.T) {
	_, _, err := resolveAuthorizationCompletion("def456", "wrong", "", "abc123")
	if err == nil {
		t.Fatal("resolveAuthorizationCompletion() error = nil, want state mismatch error")
	}
}
