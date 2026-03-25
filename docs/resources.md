# Resources

`radcli` currently supports:

- campaigns
- ad groups
- ads

Campaigns now also support `create` and `update`.

## Name or ID

For `get`, you can reference a resource by:

- ID
- exact name, case-insensitive

Examples:

```bash
./bin/rad campaign list
./bin/rad campaign get 1234567890
./bin/rad campaign get "Spring Launch"

./bin/rad adgroup get "Retargeting - US"
./bin/rad ad get "Horseflaps Ad 1"
```

If a name matches multiple resources, `radcli` will ask you to use the ID
instead.

## Campaigns

```bash
./bin/rad campaign list
./bin/rad campaign get <id-or-name>
./bin/rad campaign create --name <name> --objective <objective> --configured-status <status>
./bin/rad campaign update <id-or-name> [flags]
```

Notes:

- uses the currently selected ad account by default
- you can override with `--account-id <id-or-name>`
- `list` renders a compact table
- `get` renders JSON for deeper inspection
- `create` and `update` are documented in `docs/campaigns.md`

## Ad Groups

```bash
./bin/rad adgroup list
./bin/rad adgroup get <id-or-name>
```

## Ads

```bash
./bin/rad ad list
./bin/rad ad get <id-or-name>
```

## Account Context

These commands use the selected ad account:

```bash
./bin/rad account use <ad-account-id-or-name>
```

You can override that per command:

```bash
./bin/rad campaign get "Spring Launch" --account-id "horseflaps"
```
