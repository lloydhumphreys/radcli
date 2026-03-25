package releaseinfo

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	projectName = "radcli"
	binaryName  = "rad"
)

type Release struct {
	TagName string         `json:"tag_name"`
	Assets  []ReleaseAsset `json:"assets"`
}

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type tempReadCloser struct {
	*os.File
}

func (t *tempReadCloser) Close() error {
	name := t.Name()
	err := t.File.Close()
	_ = os.Remove(name)
	return err
}

func FetchLatestRelease(ctx context.Context, info Info) (*Release, error) {
	owner, repo, err := info.RepositoryParts()
	if err != nil {
		return nil, err
	}
	return fetchRelease(ctx, fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo))
}

func FetchRelease(ctx context.Context, info Info, version string) (*Release, error) {
	owner, repo, err := info.RepositoryParts()
	if err != nil {
		return nil, err
	}
	tag := version
	if tag != "" && !strings.HasPrefix(tag, "v") {
		tag = "v" + tag
	}
	return fetchRelease(ctx, fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag))
}

func fetchRelease(ctx context.Context, rawURL string) (*Release, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "radcli-self-update")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("github releases request failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	if release.TagName == "" {
		return nil, fmt.Errorf("github release response did not include a tag name")
	}
	return &release, nil
}

func FindAssetForCurrentPlatform(release *Release) (*ReleaseAsset, error) {
	suffix := archiveSuffixForCurrentPlatform()
	expected := fmt.Sprintf("%s_%s_%s_%s%s", projectName, versionFromTag(release.TagName), runtime.GOOS, runtime.GOARCH, suffix)
	for _, asset := range release.Assets {
		if asset.Name == expected {
			return &asset, nil
		}
	}
	return nil, fmt.Errorf("could not find a release asset for %s/%s named %q", runtime.GOOS, runtime.GOARCH, expected)
}

func versionFromTag(tag string) string {
	return strings.TrimPrefix(tag, "v")
}

func archiveSuffixForCurrentPlatform() string {
	if runtime.GOOS == "windows" {
		return ".zip"
	}
	return ".tar.gz"
}

func DownloadAndInstall(ctx context.Context, asset *ReleaseAsset) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.BrowserDownloadURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "radcli-self-update")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download failed (%d)", resp.StatusCode)
	}

	currentPath, err := os.Executable()
	if err != nil {
		return err
	}
	currentPath, err = filepath.EvalSymlinks(currentPath)
	if err != nil {
		return err
	}

	binary, err := extractBinary(asset.Name, resp.Body)
	if err != nil {
		return err
	}
	defer binary.Close()

	dir := filepath.Dir(currentPath)
	tempFile, err := os.CreateTemp(dir, "rad-update-*")
	if err != nil {
		return err
	}
	tempName := tempFile.Name()
	defer os.Remove(tempName)

	if _, err := io.Copy(tempFile, binary); err != nil {
		tempFile.Close()
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tempName, 0o755); err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		return fmt.Errorf("self-update is not supported on windows yet; please download the new release manually")
	}
	if err := os.Rename(tempName, currentPath); err != nil {
		return err
	}
	return nil
}

func extractBinary(assetName string, reader io.Reader) (io.ReadCloser, error) {
	if strings.HasSuffix(assetName, ".zip") {
		return extractBinaryFromZip(reader)
	}
	return extractBinaryFromTarGz(reader)
}

func extractBinaryFromTarGz(reader io.Reader) (io.ReadCloser, error) {
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			gz.Close()
			return nil, err
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		if filepath.Base(header.Name) != binaryName {
			continue
		}

		tempFile, err := os.CreateTemp("", "rad-binary-*")
		if err != nil {
			gz.Close()
			return nil, err
		}
		if _, err := io.Copy(tempFile, tr); err != nil {
			tempFile.Close()
			os.Remove(tempFile.Name())
			gz.Close()
			return nil, err
		}
		if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
			tempFile.Close()
			os.Remove(tempFile.Name())
			gz.Close()
			return nil, err
		}
		gz.Close()
		return &tempReadCloser{File: tempFile}, nil
	}

	gz.Close()
	return nil, fmt.Errorf("could not find %q inside the release archive", binaryName)
}

func extractBinaryFromZip(reader io.Reader) (io.ReadCloser, error) {
	tempArchive, err := os.CreateTemp("", "rad-archive-*.zip")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempArchive.Name())

	if _, err := io.Copy(tempArchive, reader); err != nil {
		tempArchive.Close()
		return nil, err
	}
	if err := tempArchive.Close(); err != nil {
		return nil, err
	}

	zr, err := zip.OpenReader(tempArchive.Name())
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	for _, file := range zr.File {
		if filepath.Base(file.Name) != binaryName+".exe" && filepath.Base(file.Name) != binaryName {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		tempFile, err := os.CreateTemp("", "rad-binary-*")
		if err != nil {
			rc.Close()
			return nil, err
		}
		if _, err := io.Copy(tempFile, rc); err != nil {
			rc.Close()
			tempFile.Close()
			os.Remove(tempFile.Name())
			return nil, err
		}
		rc.Close()
		if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
			tempFile.Close()
			os.Remove(tempFile.Name())
			return nil, err
		}
		return &tempReadCloser{File: tempFile}, nil
	}

	return nil, fmt.Errorf("could not find %q inside the release archive", binaryName)
}
