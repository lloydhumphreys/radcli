# radcli Docs

This folder is the start of the public-facing docs for `radcli`.

## Start Here

- [commands.md](./commands.md): current command reference
- [api-coverage.md](./api-coverage.md): API/docs coverage checklist
- [campaigns.md](./campaigns.md): campaign list, get, create, and update
- [adgroups.md](./adgroups.md): ad group list, get, create, and update
- [targeting.md](./targeting.md): targeting lookup helpers
- [funding.md](./funding.md): funding instrument lookup helpers
- [pixels.md](./pixels.md): pixel and conversion activity lookup helpers
- [audiences.md](./audiences.md): saved, custom, and third-party audiences
- [ads.md](./ads.md): ad list, get, create, and update
- [profiles.md](./profiles.md): profile lookup helpers
- [posts.md](./posts.md): post/creative list, get, create, and update
- [login.md](./login.md): Reddit Ads authentication flow
- [resources.md](./resources.md): campaigns, ad groups, and ads
- [reports.md](./reports.md): report commands, presets, fields, and output formats

## Current Status

`radcli` is already useful for:

- logging into Reddit Ads
- selecting a business and ad account
- listing and inspecting campaigns, ad groups, and ads
- creating and updating campaigns
- creating and updating ad groups
- looking up targeting entities
- looking up funding instruments
- looking up pixels and event activity
- managing saved audiences and looking up custom audiences
- creating and updating ads
- managing profiles and posts/creatives
- running raw reports
- running enriched summary reports
- filtering summary reports by campaign, ad group, and ad
- exporting report output as JSON or CSV

## Build

```bash
cd /Users/lloyd/Code/radcli
env GOCACHE=$PWD/.gocache go build -o bin/rad ./cmd/rad
```
