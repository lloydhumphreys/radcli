# Audiences

`radcli` now includes an audience surface for saved audiences, custom audiences,
and third-party audience lookup.

## Supported Commands

```bash
./bin/rad audience saved list
./bin/rad audience saved get <id-or-name>
./bin/rad audience saved create --name <name> --type <type> --targeting-json <json-or-@file>
./bin/rad audience saved update <id-or-name> [flags]

./bin/rad audience custom list
./bin/rad audience custom get <id-or-name>

./bin/rad audience third-party list
```

## Saved Audiences

List and inspect saved audiences:

```bash
./bin/rad audience saved list
./bin/rad audience saved get "US Retargeting Audience"
```

Create a saved audience from a targeting payload:

```bash
./bin/rad audience saved create \
  --name "US Retargeting Audience" \
  --type RETARGETING \
  --targeting-json @targeting.json \
  --dry-run
```

Update a saved audience:

```bash
./bin/rad audience saved update "US Retargeting Audience" \
  --name "US Retargeting Audience v2" \
  --targeting-json @targeting.json \
  --dry-run
```

Notes:

- saved audience writes use `--targeting-json` so we can stay flexible while the targeting shape evolves
- `--dry-run` prints the request body without sending it

## Custom Audiences

List and inspect custom audiences:

```bash
./bin/rad audience custom list
./bin/rad audience custom list --name purchasers
./bin/rad audience custom get "Purchasers 180d"
```

The list view surfaces parsed cost fields when Reddit returns them.

## Third-Party Audiences

Lookup third-party audience segments available through Reddit:

```bash
./bin/rad audience third-party list
```

The table includes:

- audience name
- full audience path
- partner
- price
- currency
- size

## Output

- table output is the default
- use `--json` for raw API payloads
- paginated commands support `--all` and `--page-size`

## Notes

- saved and custom audience commands use the selected ad account by default
- third-party audience lookup does not require an ad account argument
- `radcli` does not expose custom audience upload/mutation flows yet; this slice is intentionally focused on safe lookup plus saved-audience management
