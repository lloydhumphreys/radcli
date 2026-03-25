# Command Reference

This is the current supported command surface for `radcli`.

## Global

```bash
./bin/rad help
./bin/rad version [--json]
./bin/rad self-update [--check] [--version <tag>] [--repo <owner/repo>]
```

## Auth

```bash
./bin/rad auth setup --client-id <id> --client-secret <secret> --redirect-uri <uri> [--scope SCOPE] [--user-agent UA]
./bin/rad auth login [--open] [--no-wait]
./bin/rad auth complete --code <code> [--state <state>]
./bin/rad auth whoami [--json]
./bin/rad auth logout
```

## Config

```bash
./bin/rad config show
```

## Businesses

```bash
./bin/rad business list [--all] [--page-size N] [--json]
./bin/rad business use <business-id-or-name>
```

## Accounts

```bash
./bin/rad account list [--business-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad account use <ad-account-id-or-name>
```

## Funding

```bash
./bin/rad funding list [--account-id <id-or-name>] [--search <text>] [--type <type> ...] [--funding-instrument-id <id> ...] [--start-time <rfc3339>] [--end-time <rfc3339>] [--mode <mode>] [--selectable <true|false>] [--all] [--page-size N] [--json]
./bin/rad funding business-list [--business-id <id-or-name>] [--search <text>] [--funding-instrument-id <id> ...] [--mode <mode>] [--all] [--page-size N] [--json]
./bin/rad funding allocations <funding-instrument-id-or-name> [--account-id <id-or-name>] [--all] [--page-size N] [--json]
```

## Pixels

```bash
./bin/rad pixel list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad pixel business-list [--business-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad pixel events <pixel-id-or-name> [--account-id <id-or-name>] [--json]
```

## Audiences

```bash
./bin/rad audience saved list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad audience saved get <id-or-name> [--account-id <id-or-name>] [--json]
./bin/rad audience saved create --name <name> --type <type> --targeting-json <json-or-@file> [--account-id <id-or-name>] [--dry-run]
./bin/rad audience saved update <id-or-name> [--name <name>] [--type <type>] [--targeting-json <json-or-@file>] [--account-id <id-or-name>] [--dry-run]
./bin/rad audience custom list [--account-id <id-or-name>] [--name <text>] [--all] [--page-size N] [--json]
./bin/rad audience custom get <id-or-name> [--account-id <id-or-name>] [--json]
./bin/rad audience third-party list [--all] [--page-size N] [--json]
```

## Profiles

```bash
./bin/rad profile list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad profile business-list [--business-id <id-or-name>] [--all] [--page-size N] [--json]
```

## Posts

```bash
./bin/rad post list --profile <id-or-name> [--type <type>] [--source <source>] [--account-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad post get <post-id> [--json]
./bin/rad post create --profile <id-or-name> --type <IMAGE|VIDEO|TEXT|CAROUSEL> --headline <headline> [--body <text>] [--content-json <json-or-@file>] [--allow-comments <true|false>] [--is-richtext <true|false>] [--thumbnail-url <url>] [--account-id <id-or-name>] [--dry-run]
./bin/rad post update <post-id> --allow-comments <true|false> [--dry-run] [--json]
```

## Campaigns

```bash
./bin/rad campaign list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad campaign get <id-or-name> [--account-id <id-or-name>] [--json]
./bin/rad campaign create --name <name> --objective <objective> --configured-status <status> [--account-id <id-or-name>] [--funding-instrument-id <id>] [--invoice-label <text>] [--special-ad-category <category>] [--campaign-budget-optimization <true|false>] [--goal-type <daily_spend|lifetime_spend>] [--goal-value <major-currency>] [--spend-cap <major-currency>] [--start-time <rfc3339>] [--end-time <rfc3339>] [--bid-strategy <bidless|maximize_volume|target_cpx>] [--bid-type <cpc|cpm|cpv6>] [--bid-value <major-currency>] [--app-id <id>] [--conversion-pixel-id <id>] [--view-through-conversion-type <type>] [--dry-run]
./bin/rad campaign update <id-or-name> [--account-id <id-or-name>] [--name <name>] [--objective <objective>] [--configured-status <status>] [--funding-instrument-id <id>] [--invoice-label <text>] [--special-ad-category <category>] [--campaign-budget-optimization <true|false>] [--goal-type <daily_spend|lifetime_spend>] [--goal-value <major-currency>] [--spend-cap <major-currency>] [--start-time <rfc3339>] [--end-time <rfc3339>] [--bid-strategy <bidless|maximize_volume|target_cpx>] [--bid-type <cpc|cpm|cpv6>] [--bid-value <major-currency>] [--app-id <id>] [--conversion-pixel-id <id>] [--view-through-conversion-type <type>] [--dry-run]
```

## Ad Groups

```bash
./bin/rad adgroup list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad adgroup get <id-or-name> [--account-id <id-or-name>] [--json]
./bin/rad adgroup create --campaign <id-or-name> --name <name> --configured-status <status> [--account-id <id-or-name>] [--bid-strategy <bidless|maximize_volume|target_cpx>] [--bid-type <cpc|cpm|cpv|cpv6>] [--bid-value <major-currency>] [--goal-type <daily_spend|lifetime_spend>] [--goal-value <major-currency>] [--optimization-goal <goal>] [--optimization-strategy-type <type>] [--start-time <rfc3339>] [--end-time <rfc3339>] [--app-id <id>] [--conversion-pixel-id <id>] [--view-through-conversion-type <type>] [--saved-audience-id <id>] [--product-set-id <id>] [--shopping-type <dynamic|static>] [--targeting-json <json-or-@file>] [--schedule-json <json-or-@file>] [--shopping-targeting-json <json-or-@file>] [--dry-run]
./bin/rad adgroup update <id-or-name> [--account-id <id-or-name>] [--campaign <id-or-name>] [--name <name>] [--configured-status <status>] [--bid-strategy <bidless|maximize_volume|target_cpx>] [--bid-type <cpc|cpm|cpv|cpv6>] [--bid-value <major-currency>] [--goal-type <daily_spend|lifetime_spend>] [--goal-value <major-currency>] [--optimization-goal <goal>] [--optimization-strategy-type <type>] [--start-time <rfc3339>] [--end-time <rfc3339>] [--app-id <id>] [--conversion-pixel-id <id>] [--view-through-conversion-type <type>] [--saved-audience-id <id>] [--product-set-id <id>] [--shopping-type <dynamic|static>] [--targeting-json <json-or-@file>] [--schedule-json <json-or-@file>] [--shopping-targeting-json <json-or-@file>] [--dry-run]
```

## Ads

```bash
./bin/rad ad list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad ad get <id-or-name> [--account-id <id-or-name>] [--json]
./bin/rad ad create --ad-group <id-or-name> --name <name> --configured-status <status> [--post-id <id>] [--click-url <url>] [--click-url-query-parameter <name=value> ...] [--shopping-creative-json <json-or-@file>] [--account-id <id-or-name>] [--dry-run]
./bin/rad ad update <id-or-name> [--ad-group <id-or-name>] [--name <name>] [--configured-status <status>] [--post-id <id>] [--click-url <url>] [--click-url-query-parameter <name=value> ...] [--shopping-creative-json <json-or-@file>] [--account-id <id-or-name>] [--dry-run]
```

## Targeting

```bash
./bin/rad targeting communities search --query <text> [--all] [--page-size N] [--json]
./bin/rad targeting communities list [--name <community> ...] [--all] [--page-size N] [--json]
./bin/rad targeting communities suggest [--name <community> ...] [--website-url <url>] [--all] [--page-size N] [--json]
./bin/rad targeting interests list [--json]
./bin/rad targeting geolocations search [--postal-code <code>] [--city <name>] [--country <code>] [--json]
./bin/rad targeting devices list [--all] [--page-size N] [--json]
./bin/rad targeting carriers list [--all] [--page-size N] [--json]
./bin/rad targeting keywords suggest --keyword <term> [--keyword <term> ...] [--json]
./bin/rad targeting keywords validate --keyword <term> [--keyword <term> ...] [--json]
```

## Reports

```bash
./bin/rad report fields [--match TEXT] [--json]
```

```bash
./bin/rad report run \
  --from <iso8601> \
  --to <iso8601> \
  --field <FIELD> \
  [--field <FIELD> ...] \
  [--breakdown <BREAKDOWN> ...] \
  [--account-id <id-or-name>] \
  [--time-zone-id <tz>] \
  [--all] \
  [--page-size N] \
  [--json|--csv] \
  [--output FILE]
```

```bash
./bin/rad report campaign-summary [--since 7d] [--daily] [--account-id <id-or-name>] [--campaign <id-or-name>] [--campaign-id <id>] [--field FIELD] [--json|--csv] [--output FILE]
./bin/rad report adgroup-summary [--since 7d] [--daily] [--account-id <id-or-name>] [--campaign <id-or-name>] [--campaign-id <id>] [--adgroup <id-or-name>] [--adgroup-id <id>] [--field FIELD] [--json|--csv] [--output FILE]
./bin/rad report ad-summary [--since 7d] [--daily] [--account-id <id-or-name>] [--campaign <id-or-name>] [--campaign-id <id>] [--adgroup <id-or-name>] [--adgroup-id <id>] [--ad <id-or-name>] [--ad-id <id>] [--field FIELD] [--json|--csv] [--output FILE]
```

## Output Modes

Common patterns:

- default output is human-readable
- `--json` returns JSON
- report commands also support `--csv`
- report commands support `--output FILE` for writing output directly

## Good First Commands

```bash
./bin/rad auth whoami
./bin/rad business list
./bin/rad business use <business-name>
./bin/rad account list
./bin/rad account use <account-name>
./bin/rad funding list
./bin/rad pixel list
./bin/rad audience saved list
./bin/rad profile list
./bin/rad post create --profile "brand_profile" --type IMAGE --headline "Spring Launch" --content-json @content.json --dry-run
./bin/rad campaign list
./bin/rad ad create --ad-group "US Retargeting" --name "Spring Ad" --configured-status PAUSED --dry-run
./bin/rad campaign get <campaign-name>
./bin/rad campaign create --name "Spring Launch" --objective CLICKS --configured-status PAUSED --dry-run
./bin/rad campaign update "Spring Launch" --configured-status ACTIVE
./bin/rad adgroup create --campaign "Spring Launch" --name "US Retargeting" --configured-status PAUSED --dry-run
./bin/rad adgroup update "US Retargeting" --configured-status ACTIVE
./bin/rad targeting communities search --query gaming
./bin/rad targeting keywords suggest --keyword fishing --keyword reels
./bin/rad report campaign-summary --since 30d
./bin/rad report campaign-summary --campaign <campaign-name>
./bin/rad report campaign-summary --since 30d --csv
./bin/rad report campaign-summary --since 30d --csv --output campaign-summary-30d.csv
```
