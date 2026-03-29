# Reports

`radcli` supports both raw report requests and friendlier summary presets.

## Discover Available Fields

You can ask the CLI for the current report field list from Reddit's official v3
OpenAPI metadata:

```bash
rad report fields
rad report fields --match conversion
rad report fields --match ctr
```

## Raw Reports

Use `report run` when you want direct control over fields and breakdowns.

```bash
rad report run \
  --from 2026-03-01T00:00:00Z \
  --to 2026-03-08T00:00:00Z \
  --field IMPRESSIONS \
  --field CLICKS \
  --field SPEND \
  --breakdown CAMPAIGN_ID
```

Raw reports also support CSV output:

```bash
rad report run \
  --from 2026-03-01T00:00:00Z \
  --to 2026-03-08T00:00:00Z \
  --field IMPRESSIONS \
  --field CLICKS \
  --field SPEND \
  --breakdown CAMPAIGN_ID \
  --csv
```

You can also write report output directly to a file:

```bash
rad report run \
  --from 2026-03-01T00:00:00Z \
  --to 2026-03-08T00:00:00Z \
  --field IMPRESSIONS \
  --field CLICKS \
  --field SPEND \
  --breakdown CAMPAIGN_ID \
  --csv \
  --output reports/campaigns.csv
```

## Summary Presets

The first presets are:

- `campaign-summary`
- `adgroup-summary`
- `ad-summary`

These default to a 7-day window and return a practical metric set:

- impressions
- clicks
- spend
- CTR
- CPC
- eCPM

You can add more metrics later with extra `--field` flags once you know which
combinations your account accepts.

Preset table output is enriched with human-readable names where possible:

- campaign summaries add `campaign_name`
- ad group summaries add `ad_group_name` and `campaign_name`
- ad summaries add `ad_name`, `ad_group_name`, and `campaign_name`

### Campaign Summary

```bash
rad report campaign-summary
rad report campaign-summary --since 30d
rad report campaign-summary --daily
rad report campaign-summary --campaign "Spring Launch"
rad report campaign-summary --campaign-id 123456789
rad report campaign-summary --since 30d --csv
rad report campaign-summary --since 30d --csv --output campaign-summary-30d.csv
```

### Ad Group Summary

```bash
rad report adgroup-summary
rad report adgroup-summary --since 14d --daily
rad report adgroup-summary --campaign "Spring Launch"
rad report adgroup-summary --adgroup "Retargeting"
rad report adgroup-summary --campaign-id 123 --adgroup-id 456
```

### Ad Summary

```bash
rad report ad-summary
rad report ad-summary --since 7d
rad report ad-summary --campaign "Spring Launch" --adgroup "Retargeting"
rad report ad-summary --ad "Winner Variant"
```

## Filtering Summary Reports

Preset reports can be filtered after enrichment by either human-readable name or
ID:

- `--campaign` or `--campaign-id`
- `--adgroup` or `--adgroup-id`
- `--ad` or `--ad-id`

Support by preset:

- `campaign-summary`: campaign filters only
- `adgroup-summary`: campaign and ad group filters
- `ad-summary`: campaign, ad group, and ad filters

Examples:

```bash
rad report campaign-summary --campaign "Spring Launch" --since 30d
rad report adgroup-summary --adgroup "Retargeting" --since 7d
rad report ad-summary --ad-id 2145032584377720495 --since 14d
```

## Time Windows

Presets support:

- `--since 7d`
- `--since 30d`
- `--since 2w`
- `--since 168h`
- or explicit `--from` and `--to`

If you use `--from`, you should also use `--to`.

Reddit report windows require hourly granularity.

That means timestamps should look like:

```text
2026-03-01T00:00:00Z
```

`radcli` now rounds preset windows down to the hour automatically. If you pass
explicit `--from` and `--to` values, the CLI normalizes them to hourly UTC as
well.

## Time Zones

The API defaults to UTC unless you pass `--time-zone-id`. Use `rad timezone list`
to discover valid values.

Example:

```bash
rad timezone list
rad report campaign-summary --since 7d --time-zone-id America/Los_Angeles
```

## Extra Fields

You can extend a preset by adding more fields:

```bash
rad report campaign-summary --field CONVERSIONS
```

## Output

- preset commands render tables by default
- preset tables are enriched for readability
- use `--json` if you want the raw API response without enrichment
- use `--csv` if you want spreadsheet-friendly output
- use `--output FILE` if you want the CLI to write directly to disk
- raw reports can be rendered with either `--json` or `--csv`
- currency metrics like `spend`, `cpc`, and `ecpm` are rendered in major units for table and CSV output, while `--json` preserves Reddit's raw values
