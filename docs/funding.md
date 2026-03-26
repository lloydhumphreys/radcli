# Funding Instruments

`radcli` now includes funding instrument lookup commands so you can find the
right funding source for campaigns and inspect related allocations.

## Supported Commands

```bash
rad funding list
rad funding business-list
rad funding allocations <funding-instrument-id-or-name>
```

## Account-Level Listing

List funding instruments for the currently selected ad account:

```bash
rad funding list
```

You can filter by search text, type, explicit IDs, time window, mode, and
selectability:

```bash
rad funding list --search amex
rad funding list --type CREDIT_CARD --selectable true
rad funding list --funding-instrument-id 123 --funding-instrument-id 456
```

## Business-Level Query

Query funding instruments at the business level:

```bash
rad funding business-list
rad funding business-list --search invoice
rad funding business-list --business-id YOUR_BUSINESS_ID
```

This is useful when you want to inspect a broader funding pool before attaching
instruments to a specific ad account workflow.

## Allocations

Inspect allocations for a funding instrument by ID or exact name:

```bash
rad funding allocations "Primary Card"
rad funding allocations 1234567890
```

## Output

- table output is the default
- use `--json` for the raw API response
- paginated commands support `--all` and `--page-size`
- `credit_limit` and `billable_amount` are rendered in major currency units in table output

## Notes

- `funding list` uses the selected ad account by default
- `funding business-list` uses the selected business by default
- `funding allocations` resolves funding instrument names within the selected ad account
