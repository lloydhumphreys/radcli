package reddit

import "testing"

func TestBrowserCommandForOS(t *testing.T) {
	tests := []struct {
		name     string
		goos     string
		wantName string
		wantArgs []string
	}{
		{
			name:     "darwin",
			goos:     "darwin",
			wantName: "/usr/bin/open",
			wantArgs: []string{"https://example.com"},
		},
		{
			name:     "linux",
			goos:     "linux",
			wantName: "xdg-open",
			wantArgs: []string{"https://example.com"},
		},
		{
			name:     "windows",
			goos:     "windows",
			wantName: "rundll32",
			wantArgs: []string{"url.dll,FileProtocolHandler", "https://example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, args, err := browserCommandForOS(tt.goos, "https://example.com")
			if err != nil {
				t.Fatalf("browserCommandForOS() error = %v", err)
			}
			if name != tt.wantName {
				t.Fatalf("name = %q, want %q", name, tt.wantName)
			}
			if len(args) != len(tt.wantArgs) {
				t.Fatalf("len(args) = %d, want %d", len(args), len(tt.wantArgs))
			}
			for i := range args {
				if args[i] != tt.wantArgs[i] {
					t.Fatalf("args[%d] = %q, want %q", i, args[i], tt.wantArgs[i])
				}
			}
		})
	}
}

func TestBrowserCommandForUnsupportedOS(t *testing.T) {
	_, _, err := browserCommandForOS("plan9", "https://example.com")
	if err == nil {
		t.Fatal("browserCommandForOS() error = nil, want unsupported OS error")
	}
}
