package releaseinfo

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"runtime"
	"testing"
)

func TestRepositoryUsesEnvOverride(t *testing.T) {
	t.Setenv("RADCLI_UPDATE_REPOSITORY", " override/repo ")

	info := Info{RepoOwner: "embedded", RepoName: "repo"}
	if got := info.Repository(); got != "override/repo" {
		t.Fatalf("Repository() = %q, want %q", got, "override/repo")
	}
}

func TestRepositoryPartsFromEmbeddedMetadata(t *testing.T) {
	info := Info{RepoOwner: "lloydhumphreys", RepoName: "radcli"}

	owner, repo, err := info.RepositoryParts()
	if err != nil {
		t.Fatalf("RepositoryParts() error = %v", err)
	}
	if owner != "lloydhumphreys" || repo != "radcli" {
		t.Fatalf("RepositoryParts() = (%q, %q), want (%q, %q)", owner, repo, "lloydhumphreys", "radcli")
	}
}

func TestRepositoryPartsRejectsInvalidOverride(t *testing.T) {
	t.Setenv("RADCLI_UPDATE_REPOSITORY", "bad/repo/value")

	_, _, err := (Info{}).RepositoryParts()
	if err == nil {
		t.Fatal("RepositoryParts() error = nil, want invalid repository error")
	}
}

func TestFindAssetForCurrentPlatform(t *testing.T) {
	expectedName := fmt.Sprintf(
		"%s_%s_%s_%s%s",
		projectName,
		"1.2.3",
		runtime.GOOS,
		runtime.GOARCH,
		archiveSuffixForCurrentPlatform(),
	)
	release := &Release{
		TagName: "v1.2.3",
		Assets: []ReleaseAsset{
			{Name: "something-else.tar.gz"},
			{Name: expectedName, BrowserDownloadURL: "https://example.com/" + expectedName},
		},
	}

	asset, err := FindAssetForCurrentPlatform(release)
	if err != nil {
		t.Fatalf("FindAssetForCurrentPlatform() error = %v", err)
	}
	if asset.Name != expectedName {
		t.Fatalf("asset.Name = %q, want %q", asset.Name, expectedName)
	}
}

func TestExtractBinaryFromTarGz(t *testing.T) {
	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tw := tar.NewWriter(gz)

	payload := []byte("tar binary")
	header := &tar.Header{
		Name: "rad",
		Mode: 0o755,
		Size: int64(len(payload)),
	}
	if err := tw.WriteHeader(header); err != nil {
		t.Fatalf("WriteHeader() error = %v", err)
	}
	if _, err := tw.Write(payload); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("tar close error = %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("gzip close error = %v", err)
	}

	rc, err := extractBinary("radcli_1.2.3_darwin_arm64.tar.gz", bytes.NewReader(archive.Bytes()))
	if err != nil {
		t.Fatalf("extractBinary() error = %v", err)
	}
	defer rc.Close()

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("extracted binary = %q, want %q", string(got), string(payload))
	}
}

func TestExtractBinaryFromZip(t *testing.T) {
	var archive bytes.Buffer
	zw := zip.NewWriter(&archive)

	file, err := zw.Create("rad")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	payload := []byte("zip binary")
	if _, err := file.Write(payload); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zip close error = %v", err)
	}

	rc, err := extractBinary("radcli_1.2.3_windows_arm64.zip", bytes.NewReader(archive.Bytes()))
	if err != nil {
		t.Fatalf("extractBinary() error = %v", err)
	}
	defer rc.Close()

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("extracted binary = %q, want %q", string(got), string(payload))
	}
}
