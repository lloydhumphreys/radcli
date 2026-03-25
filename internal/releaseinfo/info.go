package releaseinfo

import (
	"fmt"
	"os"
	"strings"
)

type Info struct {
	Version   string
	Commit    string
	Date      string
	RepoOwner string
	RepoName  string
}

func (i Info) NormalizedVersion() string {
	if i.Version == "" {
		return "dev"
	}
	return i.Version
}

func (i Info) Repository() string {
	if repo := strings.TrimSpace(os.Getenv("RADCLI_UPDATE_REPOSITORY")); repo != "" {
		return repo
	}
	if i.RepoOwner == "" || i.RepoName == "" {
		return ""
	}
	return i.RepoOwner + "/" + i.RepoName
}

func (i Info) RepositoryParts() (string, string, error) {
	repo := i.Repository()
	if repo == "" {
		return "", "", fmt.Errorf("release repository is not configured. set RADCLI_UPDATE_REPOSITORY or build with repo metadata")
	}
	parts := strings.Split(repo, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid release repository %q: use owner/repo", repo)
	}
	return parts[0], parts[1], nil
}
