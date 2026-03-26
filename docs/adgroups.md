# Ad Groups

`radcli` now supports ad group reads and the first ad group write path.

## Supported Commands

```bash
rad adgroup list
rad adgroup get <id-or-name>
rad adgroup create --campaign <id-or-name> --name <name> --configured-status <status>
rad adgroup update <id-or-name> [flags]
```

## Name Or ID

For `get` and `update`, you can reference an ad group by:

- ID
- exact name, case-insensitive

If a name matches multiple ad groups, `radcli` will ask you to use the ID
instead.

## Create An Ad Group

The minimum required fields are:

- `--campaign`
- `--name`
- `--configured-status`

Example:

```bash
rad adgroup create \
  --campaign "Spring Launch" \
  --name "US Retargeting" \
  --configured-status PAUSED \
  --bid-type CPC \
  --bid-value 1.25
```

## Update An Ad Group

You can update an ad group by name or ID:

```bash
rad adgroup update "US Retargeting" --configured-status ACTIVE
rad adgroup update 123456789 --goal-value 25
```

## Dry Runs

Use `--dry-run` to inspect the JSON request body without sending it to Reddit:

```bash
rad adgroup create \
  --campaign 2145032584377720495 \
  --name "US Retargeting" \
  --configured-status PAUSED \
  --targeting-json @targeting.json \
  --dry-run
```

For `--dry-run`, using a campaign ID avoids the need to resolve the campaign by
name first.

## Write Flags

The first write cut supports:

- `--bid-strategy`
- `--bid-type`
- `--bid-value`
- `--goal-type`
- `--goal-value`
- `--optimization-goal`
- `--optimization-strategy-type`
- `--start-time`
- `--end-time`
- `--app-id`
- `--conversion-pixel-id`
- `--view-through-conversion-type`
- `--saved-audience-id`
- `--product-set-id`
- `--shopping-type`
- `--targeting-json`
- `--schedule-json`
- `--shopping-targeting-json`

## Money Inputs

`radcli` accepts major currency units for write flags like:

- `--bid-value`
- `--goal-value`

Example:

```bash
rad adgroup create \
  --campaign "Spring Launch" \
  --name "US Retargeting" \
  --configured-status PAUSED \
  --bid-value 1.25
```

That means `1.25` is sent to the API as `1250000` micros.

## JSON Inputs

For nested payloads, `radcli` accepts either inline JSON or `@file` syntax:

```bash
rad adgroup update "US Retargeting" \
  --targeting-json '{"locations":["US"]}' \
  --dry-run

rad adgroup update "US Retargeting" \
  --targeting-json @targeting.json \
  --schedule-json @schedule.json \
  --dry-run
```

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
