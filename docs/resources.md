# Resources

`radcli` currently supports:

- campaigns
- ad groups
- ads

Campaigns now also support `create` and `update`.
Ad groups now also support `create` and `update`.

## Name or ID

For `get`, you can reference a resource by:

- ID
- exact name, case-insensitive

Examples:

```bash
rad campaign list
rad campaign get 1234567890
rad campaign get "Spring Launch"

rad adgroup get "Retargeting - US"
rad ad get "Horseflaps Ad 1"
```

If a name matches multiple resources, `radcli` will ask you to use the ID
instead.

## Campaigns

```bash
rad campaign list
rad campaign get <id-or-name>
rad campaign create --name <name> --objective <objective> --configured-status <status>
rad campaign update <id-or-name> [flags]
```

Notes:

- uses the currently selected ad account by default
- you can override with `--account-id <id-or-name>`
- `list` renders a compact table
- `get` renders JSON for deeper inspection
- `create` and `update` are documented in `docs/campaigns.md`

## Ad Groups

```bash
rad adgroup list
rad adgroup get <id-or-name>
rad adgroup create --campaign <id-or-name> --name <name> --configured-status <status>
rad adgroup update <id-or-name> [flags]
```

## Ads

```bash
rad ad list
rad ad get <id-or-name>
```

## Account Context

These commands use the selected ad account:

```bash
rad account use <ad-account-id-or-name>
```

You can override that per command:

```bash
rad campaign get "Spring Launch" --account-id "horseflaps"
```
