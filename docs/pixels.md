# Pixels

`radcli` now includes lookup helpers for Reddit pixels and conversion activity.

## Supported Commands

```bash
rad pixel list
rad pixel business-list
rad pixel events <pixel-id-or-name>
```

## Account-Level Listing

List pixels for the currently selected ad account:

```bash
rad pixel list
```

This is the fastest way to find a pixel ID to use with campaign or ad group
optimization settings.

## Business-Level Listing

List pixels for the currently selected business:

```bash
rad pixel business-list
rad pixel business-list --business-id YOUR_BUSINESS_ID
```

## Event Activity

Inspect the `last_fired_at` activity breakdown for a pixel by ID or exact name:

```bash
rad pixel events "Main Pixel"
rad pixel events 1234567890
```

This returns a simple table keyed by event name, which is useful for checking
whether a pixel has recent activity for events like `purchase`, `lead`,
`add_to_cart`, and `view_content`.

## Output

- table output is the default
- use `--json` for the raw API response
- list commands support `--all` and `--page-size`

## Notes

- `pixel list` uses the selected ad account by default
- `pixel business-list` uses the selected business by default
- `pixel events` resolves pixel names within the selected ad account
- `radcli` does not expose the raw `conversion_events` write endpoint yet, because this slice is intentionally read-only and lookup-focused
