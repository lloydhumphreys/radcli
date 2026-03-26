# radcli

An unofficial command-line interface for Reddit Ads. Fast, scriptable, no browser required.

## Install

**Homebrew:**

```bash
brew tap lloydhumphreys/radcli
brew install --cask radcli
rad version
```

**GitHub release installer:**

```bash
curl -fsSL https://raw.githubusercontent.com/lloydhumphreys/radcli/main/install.sh | bash
rad version
```

If `rad` is still not found after install, open a new shell or confirm your
`PATH` includes `$(brew --prefix)/bin`.

**From source:**

```bash
git clone https://github.com/lloydhumphreys/radcli.git
cd radcli
go build -o bin/rad ./cmd/rad
```

## Get started

### 1. Create a Reddit Ads app

Go to [Reddit Business Manager developer applications](https://ads.reddit.com/business/developer-applications)
and create a developer app. You'll need the client ID, client secret, and a
redirect URI.

### 2. Configure and log in

```bash
rad auth setup \
  --client-id YOUR_CLIENT_ID \
  --client-secret YOUR_CLIENT_SECRET \
  --redirect-uri https://yourdomain.com/oauth/callback

rad auth login --open
```

This opens Reddit in your browser. Approve the app and you'll be redirected to
your redirect URI with a `code` parameter in the URL, e.g.:

```
https://yourdomain.com/oauth/callback?state=abc123&code=def456#_
```

Copy the entire URL from your browser's address bar and paste it back into the
terminal. `rad` extracts the code automatically and exchanges it for an access
token. You can also paste just the code value if you prefer.

If the redirect URI doesn't resolve to a real server, that's fine — you only
need the URL from the address bar, not for the page to load.

### 3. Pick a business and ad account

```bash
rad business list
rad business use "My Business"
rad account list
rad account use "My Ad Account"
```

### 4. Start working

```bash
rad campaign list
rad campaign get "Spring Launch"
rad report campaign-summary --since 30d
```

## What can it do?

**Browse and manage your ad structure:**

```bash
rad campaign list
rad campaign create --name "Spring Launch" --objective CLICKS --configured-status PAUSED
rad adgroup create --campaign "Spring Launch" --name "US Traffic" --configured-status PAUSED --dry-run
rad ad create --ad-group "US Traffic" --name "Hero Ad" --configured-status ACTIVE --post-id t3_abc123
rad ad update --configured-status PAUSED "Hero Ad"
```

**Run reports:**

```bash
rad report campaign-summary --since 30d
rad report ad-summary --campaign "Spring Launch" --since 7d --csv --output report.csv
rad report run --from 2026-03-01T00:00:00Z --to 2026-03-08T00:00:00Z --field IMPRESSIONS --field CLICKS
```

**Inspect creatives:**

```bash
rad post get t3_abc123
rad post create --profile t2_xyz --type IMAGE --headline "My Ad" --content-json @content.json
```

**Find targeting options:**

```bash
rad targeting communities search --query "3d printing"
rad targeting interests list
rad targeting keywords suggest --keyword "filament"
```

**Check funding, pixels, and audiences:**

```bash
rad funding list
rad pixel list
rad pixel events "Main Pixel"
rad audience saved list
```

Every command supports `--json` for machine-readable output. Reports also
support `--csv`. Use `--dry-run` on any write command to preview the request
body before sending it. Use `rad update` to install the latest published
version in place.

## Docs

- [Examples with output](./docs/examples.md)
- [Command reference](./docs/commands.md)
- [Authentication flow](./docs/login.md)
- [Reports](./docs/reports.md)

## License

[MIT](./LICENSE)

## Disclaimer

`radcli` is an unofficial tool. It is not affiliated with, endorsed by, or
supported by Reddit, Inc. or the Reddit Ads team.
