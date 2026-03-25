# Command Reference

This is the current supported command surface for `radcli`.

## Global

```bash
./bin/rad help
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

## Campaigns

```bash
./bin/rad campaign list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad campaign get <id-or-name> [--account-id <id-or-name>] [--json]
```

## Ad Groups

```bash
./bin/rad adgroup list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad adgroup get <id-or-name> [--account-id <id-or-name>] [--json]
```

## Ads

```bash
./bin/rad ad list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
./bin/rad ad get <id-or-name> [--account-id <id-or-name>] [--json]
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
./bin/rad campaign list
./bin/rad campaign get <campaign-name>
./bin/rad report campaign-summary --since 30d
./bin/rad report campaign-summary --campaign <campaign-name>
./bin/rad report campaign-summary --since 30d --csv
./bin/rad report campaign-summary --since 30d --csv --output campaign-summary-30d.csv
```
