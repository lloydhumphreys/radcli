package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lloydhumphreys/radcli/internal/cli"
	"github.com/lloydhumphreys/radcli/internal/releaseinfo"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	repoOwner = ""
	repoName  = ""
)

func main() {
	ctx := context.Background()
	app, err := cli.New(os.Stdin, os.Stdout, os.Stderr, releaseinfo.Info{
		Version:   version,
		Commit:    commit,
		Date:      date,
		RepoOwner: repoOwner,
		RepoName:  repoName,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
