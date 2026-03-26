package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lloydhumphreys/radcli/internal/output"
	"github.com/lloydhumphreys/radcli/internal/releaseinfo"
)

func (a *App) runVersion(args []string) error {
	fs := newFlagSet("version")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}

	payload := map[string]any{
		"version":    a.release.NormalizedVersion(),
		"commit":     a.release.Commit,
		"date":       a.release.Date,
		"repository": a.release.Repository(),
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}

	_, err := fmt.Fprintf(
		a.stdout,
		"radcli %s\ncommit: %s\ndate: %s\nrepository: %s\n",
		a.release.NormalizedVersion(),
		emptyFallback(a.release.Commit, "unknown"),
		emptyFallback(a.release.Date, "unknown"),
		emptyFallback(a.release.Repository(), "unconfigured"),
	)
	return err
}

func (a *App) runSelfUpdate(ctx context.Context, args []string) error {
	fs := newFlagSet("self-update")
	version := fs.String("version", "", "")
	repo := fs.String("repo", "", "")
	check := fs.Bool("check", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}

	info := a.release
	if *repo != "" {
		if err := validateRepoOverride(*repo); err != nil {
			return err
		}
		parts := splitRepo(*repo)
		info.RepoOwner = parts[0]
		info.RepoName = parts[1]
	}

	var release *releaseinfo.Release
	var err error
	switch {
	case *version != "":
		release, err = releaseinfo.FetchRelease(ctx, info, *version)
	default:
		release, err = releaseinfo.FetchLatestRelease(ctx, info)
	}
	if err != nil {
		return err
	}

	if *check {
		_, err := fmt.Fprintf(a.stdout, "Latest release: %s\n", release.TagName)
		return err
	}

	if sameReleaseVersion(release.TagName, a.release.NormalizedVersion()) {
		_, err := fmt.Fprintf(a.stdout, "Already on %s.\n", release.TagName)
		return err
	}

	asset, err := releaseinfo.FindAssetForCurrentPlatform(release)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(a.stdout, "Updating to %s...\n", release.TagName); err != nil {
		return err
	}
	if err := releaseinfo.DownloadAndInstall(ctx, asset); err != nil {
		return err
	}
	_, err = fmt.Fprintf(a.stdout, "Updated radcli to %s. Restart the command to use the new version.\n", release.TagName)
	return err
}

func validateRepoOverride(repo string) error {
	parts := splitRepo(repo)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return errors.New("invalid --repo value: use owner/repo")
	}
	return nil
}

func splitRepo(repo string) []string {
	return splitRespectingBrackets(repo, '/')
}

func emptyFallback(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func sameReleaseVersion(tag, version string) bool {
	tag = strings.TrimPrefix(tag, "v")
	version = strings.TrimPrefix(version, "v")
	return tag != "" && tag == version
}
