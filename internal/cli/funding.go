package cli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"radcli/internal/output"
)

func (a *App) runFundingCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(a.stdout, fundingHelp)
		return err
	}

	switch args[0] {
	case "list":
		return a.runFundingListCommand(ctx, args[1:])
	case "business-list":
		return a.runFundingBusinessListCommand(ctx, args[1:])
	case "allocations":
		return a.runFundingAllocationsCommand(ctx, args[1:])
	default:
		return fmt.Errorf("unknown funding command %q\n\n%s", args[0], fundingHelp)
	}
}

func (a *App) runFundingListCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("funding list")
	accountInput := fs.String("account-id", "", "")
	search := fs.String("search", "", "")
	startTime := fs.String("start-time", "", "")
	endTime := fs.String("end-time", "", "")
	mode := fs.String("mode", "", "")
	selectable := fs.String("selectable", "", "")
	all := fs.Bool("all", false, "")
	pageSize := fs.Int("page-size", 0, "")
	jsonOut := fs.Bool("json", false, "")
	var ids stringList
	var types stringList
	fs.Var(&ids, "funding-instrument-id", "")
	fs.Var(&types, "type", "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}

	query := url.Values{}
	if *search != "" {
		query.Set("search", *search)
	}
	if *startTime != "" {
		query.Set("start_time", *startTime)
	}
	if *endTime != "" {
		query.Set("end_time", *endTime)
	}
	if *mode != "" {
		query.Set("mode", *mode)
	}
	if *selectable != "" {
		query.Set("is_selectable", *selectable)
	}
	if len(ids) > 0 {
		query.Set("funding_instrument_ids", joinCSV([]string(ids)))
	}
	if len(types) > 0 {
		query.Set("types", joinCSV([]string(types)))
	}
	if *pageSize > 0 {
		query.Set("page.size", fmt.Sprintf("%d", *pageSize))
	}

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/ad_accounts/"+accountID+"/funding_instruments", query, nil, *all)
	if err != nil {
		return err
	}

	rows := formatFundingRows(rowsFromPayload(payload))
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rows, fundingColumns())
}

func (a *App) runFundingBusinessListCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("funding business-list")
	businessInput := fs.String("business-id", "", "")
	search := fs.String("search", "", "")
	mode := fs.String("mode", "", "")
	all := fs.Bool("all", false, "")
	pageSize := fs.Int("page-size", 0, "")
	jsonOut := fs.Bool("json", false, "")
	var ids stringList
	fs.Var(&ids, "funding-instrument-id", "")
	fs.Var(&ids, "id", "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}

	businessID, _, err := a.selectedBusinessID(ctx, *businessInput)
	if err != nil {
		return err
	}

	query := url.Values{}
	if *search != "" {
		query.Set("search", *search)
	}
	if *pageSize > 0 {
		query.Set("page.size", fmt.Sprintf("%d", *pageSize))
	}

	data := map[string]any{}
	if len(ids) > 0 {
		data["funding_instrument_ids"] = []string(ids)
	}
	if *mode != "" {
		data["mode"] = *mode
	}

	payload, err := a.api.RequestPaginatedJSON(ctx, "POST", "/businesses/"+businessID+"/funding_instruments/query", query, map[string]any{"data": data}, *all)
	if err != nil {
		return err
	}

	rows := formatFundingRows(rowsFromPayload(payload))
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rows, fundingColumns())
}

func (a *App) runFundingAllocationsCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("funding allocations")
	accountInput := fs.String("account-id", "", "")
	all := fs.Bool("all", false, "")
	pageSize := fs.Int("page-size", 0, "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return errors.New("usage: rad funding allocations <funding-instrument-id-or-name> [flags]")
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}
	fundingID, _, err := a.resolveFundingSelection(ctx, accountID, fs.Arg(0))
	if err != nil {
		return err
	}

	query := url.Values{}
	if *pageSize > 0 {
		query.Set("page.size", fmt.Sprintf("%d", *pageSize))
	}

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/funding_instruments/"+fundingID+"/allocations", query, nil, *all)
	if err != nil {
		return err
	}

	rows := formatFundingRows(rowsFromPayload(payload))
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rows, fundingColumns())
}

func (a *App) resolveFundingSelection(ctx context.Context, accountID, input string) (string, string, error) {
	query := url.Values{}
	query.Set("page.size", "1000")

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/ad_accounts/"+accountID+"/funding_instruments", query, nil, true)
	if err != nil {
		return "", "", err
	}

	rows := rowsFromPayload(payload)
	if len(rows) == 0 {
		return "", "", errors.New("no funding instruments returned for the selected ad account")
	}

	for _, row := range rows {
		if row["id"] == input {
			return row["id"], row["name"], nil
		}
	}

	matches := make([]map[string]string, 0)
	for _, row := range rows {
		if strings.EqualFold(row["name"], input) {
			matches = append(matches, row)
		}
	}

	if len(matches) == 1 {
		return matches[0]["id"], matches[0]["name"], nil
	}
	if len(matches) > 1 {
		ids := make([]string, 0, len(matches))
		for _, match := range matches {
			ids = append(ids, match["id"])
		}
		return "", "", fmt.Errorf("funding instrument name %q matched multiple funding instruments: %s. use the ID instead", input, strings.Join(ids, ", "))
	}

	return "", "", fmt.Errorf("could not find funding instrument %q in the selected account. run `rad funding list` and use the ID from the `id` column", input)
}

func formatFundingRows(rows []map[string]string) []map[string]string {
	for _, row := range rows {
		for _, key := range []string{"credit_limit", "billable_amount"} {
			if value := row[key]; value != "" {
				row[key] = formatMicroCurrency(value)
			}
		}
	}
	return rows
}

func fundingColumns() []string {
	return []string{"name", "id", "currency", "credit_limit", "billable_amount", "is_servable", "authorize_status", "invoice_group_status", "start_time", "end_time", "reasons_not_servable"}
}

const fundingHelp = `Usage:
  rad funding list [--account-id <id-or-name>] [--search <text>] [--type <type> ...] [--funding-instrument-id <id> ...] [--start-time <rfc3339>] [--end-time <rfc3339>] [--mode <mode>] [--selectable <true|false>] [--all] [--page-size N] [--json]
  rad funding business-list [--business-id <id-or-name>] [--search <text>] [--funding-instrument-id <id> ...] [--mode <mode>] [--all] [--page-size N] [--json]
  rad funding allocations <funding-instrument-id-or-name> [--account-id <id-or-name>] [--all] [--page-size N] [--json]

Examples:
  rad funding list
  rad funding list --search amex
  rad funding business-list --search invoice
  rad funding allocations "Primary Card"`
