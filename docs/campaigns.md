# Campaigns

`radcli` now supports campaign reads and the first campaign write path.

## Supported Commands

```bash
rad campaign list
rad campaign get <id-or-name>
rad campaign create --name <name> --objective <objective> --configured-status <status>
rad campaign update <id-or-name> [flags]
```

## Name Or ID

For `get` and `update`, you can reference a campaign by:

- ID
- exact name, case-insensitive

If a name matches multiple campaigns, `radcli` will ask you to use the ID
instead.

## Create A Campaign

The minimum required fields are:

- `--name`
- `--objective`
- `--configured-status`

Example:

```bash
rad campaign create \
  --name "Spring Launch" \
  --objective CLICKS \
  --configured-status PAUSED
```

## Update A Campaign

You can update a campaign by name or ID:

```bash
rad campaign update "Spring Launch" --configured-status ACTIVE
rad campaign update 123456789 --invoice-label "client-2026-q2"
```

## Dry Runs

Use `--dry-run` to inspect the JSON request body without sending it to Reddit:

```bash
rad campaign create \
  --name "Spring Launch" \
  --objective CLICKS \
  --configured-status PAUSED \
  --spend-cap 250 \
  --dry-run
```

## Write Flags

The first write cut supports:

- `--funding-instrument-id`
- `--invoice-label`
- `--special-ad-category`
- `--campaign-budget-optimization true|false`
- `--goal-type daily_spend|lifetime_spend`
- `--goal-value`
- `--spend-cap`
- `--start-time`
- `--end-time`
- `--bid-strategy bidless|maximize_volume|target_cpx`
- `--bid-type cpc|cpm|cpv6`
- `--bid-value`
- `--app-id`
- `--conversion-pixel-id`
- `--view-through-conversion-type seven_day_clicks|seven_day_clicks_one_day_view`

## Money Inputs

`radcli` accepts major currency units for write flags like:

- `--spend-cap`
- `--goal-value`
- `--bid-value`

Example:

```bash
rad campaign create \
  --name "Spring Launch" \
  --objective CLICKS \
  --configured-status PAUSED \
  --spend-cap 250
```

That means `250` is sent to the API as `250000000` micros.

## Time Inputs

If you pass one of `--start-time` or `--end-time`, you should pass both.

Use RFC3339 timestamps:

```text
2026-04-01T00:00:00Z
```

## Notes

- write commands use the selected ad account by default
- override with `--account-id <id-or-name>` when needed
- `create` and `update` both return JSON responses
- `--json` is optional today because write commands already print JSON by default
