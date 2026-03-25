# Ads

`radcli` now supports ad reads plus a first ad write path.

## Supported Commands

```bash
./bin/rad ad list
./bin/rad ad get <id-or-name>
./bin/rad ad create --ad-group <id-or-name> --name <name> --configured-status <status>
./bin/rad ad update <id-or-name> [flags]
```

## Name Or ID

For `get` and `update`, you can reference an ad by:

- ID
- exact name, case-insensitive

If a name matches multiple ads, `radcli` will ask you to use the ID instead.

## Create An Ad

The minimum required fields are:

- `--ad-group`
- `--name`
- `--configured-status`

Example:

```bash
./bin/rad ad create \
  --ad-group "US Retargeting" \
  --name "Spring Ad" \
  --configured-status PAUSED \
  --post-id t3_abcdef \
  --dry-run
```

## Update An Ad

You can update an ad by name or ID:

```bash
./bin/rad ad update "Spring Ad" --configured-status ACTIVE
./bin/rad ad update 123456789 --click-url https://example.com
```

## First Write Slice

The current write path supports the shared fields we can wrap safely today:

- `--post-id`
- `--click-url`
- `--click-url-query-parameter name=value`
- `--shopping-creative-json`

This is intentionally narrower than the full Ads API surface, but it covers the
stable shared fields without pretending we already have a perfect wrapper for
every creative variant.

## Shopping Creative JSON

For shopping ads, pass the creative payload as inline JSON or `@file`:

```bash
./bin/rad ad update "Catalog Ad" \
  --shopping-creative-json @shopping-creative.json \
  --dry-run
```

## Dry Runs

Use `--dry-run` to inspect the JSON request body without sending it to Reddit:

```bash
./bin/rad ad create \
  --ad-group 2145032584377720495 \
  --name "Spring Ad" \
  --configured-status PAUSED \
  --post-id t3_abcdef \
  --dry-run
```

If you use an ad group ID for `--ad-group`, dry runs can be built without
having to resolve the ad group name first.

## Query Parameters

`--click-url-query-parameter` accepts repeated `name=value` pairs:

```bash
./bin/rad ad update "Spring Ad" \
  --click-url https://example.com \
  --click-url-query-parameter utm_source=reddit \
  --click-url-query-parameter utm_medium={{AD_ID}}
```

## Notes

- write commands use the selected ad account by default
- override with `--account-id <id-or-name>` when needed
- `create` and `update` return JSON responses
- this slice does not expose every possible ad creative field yet
