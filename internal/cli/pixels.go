package cli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/lloydhumphreys/radcli/internal/output"
)

func (a *App) runPixelCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintln(a.stdout, pixelHelp)
		return err
	}

	switch args[0] {
	case "list":
		return a.runPixelListCommand(ctx, args[1:])
	case "business-list":
		return a.runPixelBusinessListCommand(ctx, args[1:])
	case "events", "last-fired-at":
		return a.runPixelEventsCommand(ctx, args[1:])
	default:
		return fmt.Errorf("unknown pixel command %q\n\n%s", args[0], pixelHelp)
	}
}

func (a *App) runPixelListCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("pixel list")
	accountInput := fs.String("account-id", "", "")
	all := fs.Bool("all", false, "")
	pageSize := fs.Int("page-size", 0, "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}

	query := url.Values{}
	if *pageSize > 0 {
		query.Set("page.size", fmt.Sprintf("%d", *pageSize))
	}

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/ad_accounts/"+accountID+"/pixels", query, nil, *all)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rowsFromPayload(payload), pixelColumns())
}

func (a *App) runPixelBusinessListCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("pixel business-list")
	businessInput := fs.String("business-id", "", "")
	all := fs.Bool("all", false, "")
	pageSize := fs.Int("page-size", 0, "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}

	businessID, _, err := a.selectedBusinessID(ctx, *businessInput)
	if err != nil {
		return err
	}

	query := url.Values{}
	if *pageSize > 0 {
		query.Set("page.size", fmt.Sprintf("%d", *pageSize))
	}

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/businesses/"+businessID+"/pixels", query, nil, *all)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rowsFromPayload(payload), pixelColumns())
}

func (a *App) runPixelEventsCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("pixel events")
	accountInput := fs.String("account-id", "", "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return errors.New("usage: rad pixel events <pixel-id-or-name> [--account-id <id-or-name>] [--json]")
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}
	pixelID, pixelName, err := a.resolvePixelSelection(ctx, accountID, fs.Arg(0))
	if err != nil {
		return err
	}

	payload, err := a.api.RequestJSON(ctx, "GET", "/pixels/"+pixelID+"/last_fired_at", nil, nil)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	if _, err := fmt.Fprintf(a.stdout, "Pixel %s (%s)\n\n", pixelName, pixelID); err != nil {
		return err
	}
	return output.PrintTable(a.stdout, pixelEventRows(payload), []string{"event", "last_fired_at"})
}

func (a *App) resolvePixelSelection(ctx context.Context, accountID, input string) (string, string, error) {
	query := url.Values{}
	query.Set("page.size", "1000")

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/ad_accounts/"+accountID+"/pixels", query, nil, true)
	if err != nil {
		return "", "", err
	}

	rows := rowsFromPayload(payload)
	if len(rows) == 0 {
		return "", "", errors.New("no pixels returned for the selected ad account")
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
		return "", "", fmt.Errorf("pixel name %q matched multiple pixels: %s. use the ID instead", input, strings.Join(ids, ", "))
	}

	return "", "", fmt.Errorf("could not find pixel %q in the selected account. run `rad pixel list` and use the ID from the `id` column", input)
}

func pixelColumns() []string {
	return []string{"name", "id", "business_id", "created_at", "created_by", "modified_at", "modified_by", "automatic_matching_config"}
}

func pixelEventRows(payload map[string]any) []map[string]string {
	data, _ := payload["data"].(map[string]any)
	if len(data) == 0 {
		return nil
	}

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	rows := make([]map[string]string, 0, len(keys))
	for _, key := range keys {
		rows = append(rows, map[string]string{
			"event":         key,
			"last_fired_at": stringify(data[key]),
		})
	}
	return rows
}

const pixelHelp = `Usage:
  rad pixel list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
  rad pixel business-list [--business-id <id-or-name>] [--all] [--page-size N] [--json]
  rad pixel events <pixel-id-or-name> [--account-id <id-or-name>] [--json]

Examples:
  rad pixel list
  rad pixel business-list
  rad pixel events "Main Pixel"
  rad pixel events 1234567890`
