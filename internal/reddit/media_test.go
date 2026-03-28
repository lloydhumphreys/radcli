package reddit

import (
	"testing"
)

func TestParseUploadLease(t *testing.T) {
	body := []byte(`{
		"args": {
			"action": "//reddit-uploaded-media.s3-accelerate.amazonaws.com/rte_images/abc123",
			"fields": [
				{"name": "x-amz-security-token", "value": "tok123"},
				{"name": "key", "value": "rte_images/abc123"},
				{"name": "x-amz-algorithm", "value": "AWS4-HMAC-SHA256"},
				{"name": "policy", "value": "eyJ..."}
			]
		},
		"asset": {
			"asset_id": "asset-abc-123",
			"websocket_url": "wss://ws-082c59a638dda.wss.redditmedia.com/abc123"
		}
	}`)

	lease, err := parseUploadLease(body)
	if err != nil {
		t.Fatalf("parseUploadLease() error = %v", err)
	}

	if lease.UploadURL != "https://reddit-uploaded-media.s3-accelerate.amazonaws.com/rte_images/abc123" {
		t.Errorf("UploadURL = %q", lease.UploadURL)
	}
	if lease.AssetID != "asset-abc-123" {
		t.Errorf("AssetID = %q", lease.AssetID)
	}
	if len(lease.Fields) != 4 {
		t.Errorf("Fields count = %d, want 4", len(lease.Fields))
	}
	if lease.Fields["key"] != "rte_images/abc123" {
		t.Errorf("Fields[key] = %q", lease.Fields["key"])
	}
}

func TestParseUploadLeaseNoProtocol(t *testing.T) {
	body := []byte(`{
		"args": {
			"action": "https://reddit-uploaded-media.s3-accelerate.amazonaws.com/path",
			"fields": []
		},
		"asset": {"asset_id": "id1"}
	}`)

	lease, err := parseUploadLease(body)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if lease.UploadURL != "https://reddit-uploaded-media.s3-accelerate.amazonaws.com/path" {
		t.Errorf("UploadURL = %q", lease.UploadURL)
	}
}

func TestParseS3Location(t *testing.T) {
	body := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<PostResponse>
  <Location>https%3A%2F%2Freddit-uploaded-media.s3-accelerate.amazonaws.com%2Frte_images%2Fabc123</Location>
  <Bucket>reddit-uploaded-media</Bucket>
  <Key>rte_images/abc123</Key>
  <ETag>"abc123"</ETag>
</PostResponse>`)

	loc, err := parseS3Location(body)
	if err != nil {
		t.Fatalf("parseS3Location() error = %v", err)
	}
	if loc != "https://reddit-uploaded-media.s3-accelerate.amazonaws.com/rte_images/abc123" {
		t.Errorf("location = %q", loc)
	}
}

func TestParseS3LocationRaw(t *testing.T) {
	body := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<PostResponse>
  <Location>https://reddit-uploaded-media.s3-accelerate.amazonaws.com/rte_images/abc123</Location>
</PostResponse>`)

	loc, err := parseS3Location(body)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if loc != "https://reddit-uploaded-media.s3-accelerate.amazonaws.com/rte_images/abc123" {
		t.Errorf("location = %q", loc)
	}
}

func TestTruncate(t *testing.T) {
	if got := truncate("hello", 10); got != "hello" {
		t.Errorf("truncate short = %q", got)
	}
	if got := truncate("hello world", 5); got != "hello..." {
		t.Errorf("truncate long = %q", got)
	}
}
