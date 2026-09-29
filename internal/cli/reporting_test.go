package cli

import "testing"

func TestReportHasMorePages(t *testing.T) {
	cases := []struct {
		name    string
		payload map[string]any
		want    bool
	}{
		{"no pagination", map[string]any{"data": []any{}}, false},
		{"empty next_url", map[string]any{"pagination": map[string]any{"next_url": ""}}, false},
		{"next_url set", map[string]any{"pagination": map[string]any{"next_url": "https://ads-api.reddit.com/api/v3/x?page.token=abc"}}, true},
		{"nil payload", nil, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := reportHasMorePages(tc.payload); got != tc.want {
				t.Fatalf("reportHasMorePages() = %v, want %v", got, tc.want)
			}
		})
	}
}
