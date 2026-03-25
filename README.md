# radcli

`radcli` is a command-line interface for Reddit Ads.

The goal is simple: make Reddit Ads fast, scriptable, and predictable without
having to live in the website.

Today `radcli` is already useful for:

- authenticating with Reddit Ads
- selecting a default business and ad account
- listing and inspecting campaigns, ad groups, and ads
- running raw reports
- running enriched summary reports
- filtering reports by campaign, ad group, and ad
- exporting report output as tables, JSON, or CSV

## Why

The Reddit Ads web UI is powerful, but it can also be slow and hard to navigate
for repeated operator tasks. `radcli` is meant to be a daily-driver tool for:

- advertisers
- agencies
- internal growth teams
- anyone who wants Reddit Ads to behave more like infrastructure

## Build

```bash
cd /Users/lloyd/Code/radcli
env GOCACHE=$PWD/.gocache go build -o bin/rad ./cmd/rad
```

## Quick Start

1. Create a Reddit Ads developer app in Business Manager.
2. Configure the CLI with your app credentials.
3. Log in once.
4. Pick a business and ad account.
5. Start listing assets and running reports.

Example:

```bash
./bin/rad auth setup \
  --client-id YOUR_CLIENT_ID \
  --client-secret YOUR_CLIENT_SECRET \
  --redirect-uri https://YOURDOMAIN.com/oauth/reddit-ads/callback \
  --scope adsread \
  --scope adsedit \
  --scope adsconversions \
  --scope history \
  --scope read \
  --user-agent 'macos:com.example.radcli:v0.1.0 (by /u/YOUR_USERNAME)'

./bin/rad auth login
./bin/rad auth whoami
./bin/rad business list
./bin/rad business use YOUR_BUSINESS_NAME
./bin/rad account list
./bin/rad account use YOUR_ACCOUNT_NAME
./bin/rad campaign list
./bin/rad report campaign-summary --since 30d
```

## Example Commands

```bash
./bin/rad campaign get "Spring Launch"
./bin/rad adgroup get "Retargeting"
./bin/rad ad get "Winner Variant"

./bin/rad report campaign-summary --since 30d
./bin/rad report campaign-summary --campaign "Spring Launch" --since 30d
./bin/rad report adgroup-summary --adgroup "Retargeting" --since 14d
./bin/rad report ad-summary --ad "Winner Variant" --since 7d

./bin/rad report campaign-summary --since 30d --csv --output campaign-summary-30d.csv
```

## Docs

The deeper docs live in [`docs/`](./docs):

- [`docs/index.md`](./docs/index.md)
- [`docs/commands.md`](./docs/commands.md)
- [`docs/login.md`](./docs/login.md)
- [`docs/resources.md`](./docs/resources.md)
- [`docs/reports.md`](./docs/reports.md)
- [`plan.md`](./plan.md)

## Current Scope

Implemented command groups:

- `auth`
- `config`
- `business`
- `account`
- `campaign`
- `adgroup`
- `ad`
- `report`

## Next

The next major milestone is the write path:

- campaign create/update
- ad group create/update
- dry-run support for write commands
- more operator-friendly report and targeting workflows
