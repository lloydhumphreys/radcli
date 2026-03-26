# Command Reference

This is the current supported command surface for `radcli`.

## Global

```bash
rad help
rad version [--json]
rad update [--check] [--version <tag>] [--repo <owner/repo>]
```

## Auth

```bash
rad auth setup --client-id <id> --client-secret <secret> --redirect-uri <uri> [--scope SCOPE] [--user-agent UA]
rad auth login [--open] [--no-wait]
rad auth complete --code <code> [--state <state>]
rad auth whoami [--json]
rad auth logout
```

## Config

```bash
rad config show
```

## Businesses

```bash
rad business list [--all] [--page-size N] [--json]
rad business use <business-id-or-name>
```

## Accounts

```bash
rad account list [--business-id <id-or-name>] [--all] [--page-size N] [--json]
rad account use <ad-account-id-or-name>
```

## Funding

```bash
rad funding list [--account-id <id-or-name>] [--search <text>] [--type <type> ...] [--funding-instrument-id <id> ...] [--start-time <rfc3339>] [--end-time <rfc3339>] [--mode <mode>] [--selectable <true|false>] [--all] [--page-size N] [--json]
rad funding business-list [--business-id <id-or-name>] [--search <text>] [--funding-instrument-id <id> ...] [--mode <mode>] [--all] [--page-size N] [--json]
rad funding allocations <funding-instrument-id-or-name> [--account-id <id-or-name>] [--all] [--page-size N] [--json]
```

## Pixels

```bash
rad pixel list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
rad pixel business-list [--business-id <id-or-name>] [--all] [--page-size N] [--json]
rad pixel events <pixel-id-or-name> [--account-id <id-or-name>] [--json]
```

## Audiences

```bash
rad audience saved list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
rad audience saved get <id-or-name> [--account-id <id-or-name>] [--json]
rad audience saved create --name <name> --type <type> --targeting-json <json-or-@file> [--account-id <id-or-name>] [--dry-run]
rad audience saved update <id-or-name> [--name <name>] [--type <type>] [--targeting-json <json-or-@file>] [--account-id <id-or-name>] [--dry-run]
rad audience custom list [--account-id <id-or-name>] [--name <text>] [--all] [--page-size N] [--json]
rad audience custom get <id-or-name> [--account-id <id-or-name>] [--json]
rad audience third-party list [--all] [--page-size N] [--json]
```

## Profiles

```bash
rad profile list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
rad profile business-list [--business-id <id-or-name>] [--all] [--page-size N] [--json]
```

## Posts

```bash
rad post list --profile <id-or-name> [--type <type>] [--source <source>] [--account-id <id-or-name>] [--all] [--page-size N] [--json]
rad post get <post-id> [--json]
rad post create --profile <id-or-name> --type <IMAGE|VIDEO|TEXT|CAROUSEL> --headline <headline> [--body <text>] [--content-json <json-or-@file>] [--allow-comments <true|false>] [--is-richtext <true|false>] [--thumbnail-url <url>] [--account-id <id-or-name>] [--dry-run]
rad post update <post-id> --allow-comments <true|false> [--dry-run] [--json]
```

## Campaigns

```bash
rad campaign list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
rad campaign get <id-or-name> [--account-id <id-or-name>] [--json]
rad campaign create --name <name> --objective <objective> --configured-status <status> [--account-id <id-or-name>] [--funding-instrument-id <id>] [--invoice-label <text>] [--special-ad-category <category>] [--campaign-budget-optimization <true|false>] [--goal-type <daily_spend|lifetime_spend>] [--goal-value <major-currency>] [--spend-cap <major-currency>] [--start-time <rfc3339>] [--end-time <rfc3339>] [--bid-strategy <bidless|maximize_volume|target_cpx>] [--bid-type <cpc|cpm|cpv6>] [--bid-value <major-currency>] [--app-id <id>] [--conversion-pixel-id <id>] [--view-through-conversion-type <type>] [--dry-run]
rad campaign update <id-or-name> [--account-id <id-or-name>] [--name <name>] [--objective <objective>] [--configured-status <status>] [--funding-instrument-id <id>] [--invoice-label <text>] [--special-ad-category <category>] [--campaign-budget-optimization <true|false>] [--goal-type <daily_spend|lifetime_spend>] [--goal-value <major-currency>] [--spend-cap <major-currency>] [--start-time <rfc3339>] [--end-time <rfc3339>] [--bid-strategy <bidless|maximize_volume|target_cpx>] [--bid-type <cpc|cpm|cpv6>] [--bid-value <major-currency>] [--app-id <id>] [--conversion-pixel-id <id>] [--view-through-conversion-type <type>] [--dry-run]
```

## Ad Groups

```bash
rad adgroup list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
rad adgroup get <id-or-name> [--account-id <id-or-name>] [--json]
rad adgroup create --campaign <id-or-name> --name <name> --configured-status <status> [--account-id <id-or-name>] [--bid-strategy <bidless|maximize_volume|target_cpx>] [--bid-type <cpc|cpm|cpv|cpv6>] [--bid-value <major-currency>] [--goal-type <daily_spend|lifetime_spend>] [--goal-value <major-currency>] [--optimization-goal <goal>] [--optimization-strategy-type <type>] [--start-time <rfc3339>] [--end-time <rfc3339>] [--app-id <id>] [--conversion-pixel-id <id>] [--view-through-conversion-type <type>] [--saved-audience-id <id>] [--product-set-id <id>] [--shopping-type <dynamic|static>] [--targeting-json <json-or-@file>] [--schedule-json <json-or-@file>] [--shopping-targeting-json <json-or-@file>] [--dry-run]
rad adgroup update <id-or-name> [--account-id <id-or-name>] [--campaign <id-or-name>] [--name <name>] [--configured-status <status>] [--bid-strategy <bidless|maximize_volume|target_cpx>] [--bid-type <cpc|cpm|cpv|cpv6>] [--bid-value <major-currency>] [--goal-type <daily_spend|lifetime_spend>] [--goal-value <major-currency>] [--optimization-goal <goal>] [--optimization-strategy-type <type>] [--start-time <rfc3339>] [--end-time <rfc3339>] [--app-id <id>] [--conversion-pixel-id <id>] [--view-through-conversion-type <type>] [--saved-audience-id <id>] [--product-set-id <id>] [--shopping-type <dynamic|static>] [--targeting-json <json-or-@file>] [--schedule-json <json-or-@file>] [--shopping-targeting-json <json-or-@file>] [--dry-run]
```

## Ads

```bash
rad ad list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
rad ad get <id-or-name> [--account-id <id-or-name>] [--json]
rad ad create --ad-group <id-or-name> --name <name> --configured-status <status> [--post-id <id>] [--click-url <url>] [--click-url-query-parameter <name=value> ...] [--shopping-creative-json <json-or-@file>] [--account-id <id-or-name>] [--dry-run]
rad ad update <id-or-name> [--ad-group <id-or-name>] [--name <name>] [--configured-status <status>] [--post-id <id>] [--click-url <url>] [--click-url-query-parameter <name=value> ...] [--shopping-creative-json <json-or-@file>] [--account-id <id-or-name>] [--dry-run]
```

## Targeting

```bash
rad targeting communities search --query <text> [--all] [--page-size N] [--json]
rad targeting communities list [--name <community> ...] [--all] [--page-size N] [--json]
rad targeting communities suggest [--name <community> ...] [--website-url <url>] [--all] [--page-size N] [--json]
rad targeting interests list [--json]
rad targeting geolocations search [--postal-code <code>] [--city <name>] [--country <code>] [--json]
rad targeting devices list [--all] [--page-size N] [--json]
rad targeting carriers list [--all] [--page-size N] [--json]
rad targeting keywords suggest --keyword <term> [--keyword <term> ...] [--json]
rad targeting keywords validate --keyword <term> [--keyword <term> ...] [--json]
```

## Reports

```bash
rad report fields [--match TEXT] [--json]
```

```bash
rad report run \
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
rad report campaign-summary [--since 7d] [--daily] [--account-id <id-or-name>] [--campaign <id-or-name>] [--campaign-id <id>] [--field FIELD] [--json|--csv] [--output FILE]
rad report adgroup-summary [--since 7d] [--daily] [--account-id <id-or-name>] [--campaign <id-or-name>] [--campaign-id <id>] [--adgroup <id-or-name>] [--adgroup-id <id>] [--field FIELD] [--json|--csv] [--output FILE]
rad report ad-summary [--since 7d] [--daily] [--account-id <id-or-name>] [--campaign <id-or-name>] [--campaign-id <id>] [--adgroup <id-or-name>] [--adgroup-id <id>] [--ad <id-or-name>] [--ad-id <id>] [--field FIELD] [--json|--csv] [--output FILE]
```

## Output Modes

Common patterns:

- default output is human-readable
- `--json` returns JSON
- report commands also support `--csv`
- report commands support `--output FILE` for writing output directly

## Good First Commands

```bash
rad auth whoami
rad business list
rad business use <business-name>
rad account list
rad account use <account-name>
rad funding list
rad pixel list
rad audience saved list
rad profile list
rad post create --profile "brand_profile" --type IMAGE --headline "Spring Launch" --content-json @content.json --dry-run
rad campaign list
rad ad create --ad-group "US Retargeting" --name "Spring Ad" --configured-status PAUSED --dry-run
rad campaign get <campaign-name>
rad campaign create --name "Spring Launch" --objective CLICKS --configured-status PAUSED --dry-run
rad campaign update "Spring Launch" --configured-status ACTIVE
rad adgroup create --campaign "Spring Launch" --name "US Retargeting" --configured-status PAUSED --dry-run
rad adgroup update "US Retargeting" --configured-status ACTIVE
rad targeting communities search --query gaming
rad targeting keywords suggest --keyword fishing --keyword reels
rad report campaign-summary --since 30d
rad report campaign-summary --campaign <campaign-name>
rad report campaign-summary --since 30d --csv
rad report campaign-summary --since 30d --csv --output campaign-summary-30d.csv
```
