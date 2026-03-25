# radcli Docs

This folder is the start of the public-facing docs for `radcli`.

## Start Here

- [commands.md](./commands.md): current command reference
- [login.md](./login.md): Reddit Ads authentication flow
- [resources.md](./resources.md): campaigns, ad groups, and ads
- [reports.md](./reports.md): report commands, presets, fields, and output formats

## Current Status

`radcli` is already useful for:

- logging into Reddit Ads
- selecting a business and ad account
- listing and inspecting campaigns, ad groups, and ads
- running raw reports
- running enriched summary reports
- filtering summary reports by campaign, ad group, and ad
- exporting report output as JSON or CSV

## Build

```bash
cd /Users/lloyd/Code/radcli
env GOCACHE=$PWD/.gocache go build -o bin/rad ./cmd/rad
```
