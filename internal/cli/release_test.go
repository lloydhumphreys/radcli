package cli

import "testing"

func TestValidateRepoOverride(t *testing.T) {
	if err := validateRepoOverride("lloydhumphreys/radcli"); err != nil {
		t.Fatalf("validateRepoOverride() error = %v", err)
	}
	if err := validateRepoOverride("radcli"); err == nil {
		t.Fatal("validateRepoOverride() error = nil, want invalid repo error")
	}
}

func TestSameReleaseVersion(t *testing.T) {
	if !sameReleaseVersion("v1.2.3", "1.2.3") {
		t.Fatal("sameReleaseVersion(v-prefixed tag, plain version) = false, want true")
	}
	if sameReleaseVersion("v1.2.3", "1.2.4") {
		t.Fatal("sameReleaseVersion() = true for mismatched versions, want false")
	}
}
