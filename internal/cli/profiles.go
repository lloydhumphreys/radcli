package cli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"radcli/internal/output"
)

func (a *App) runProfileCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(a.stdout, profileHelp)
		return err
	}

	switch args[0] {
	case "list":
		return a.runProfileListCommand(ctx, args[1:])
	case "business-list":
		return a.runProfileBusinessListCommand(ctx, args[1:])
	default:
		return fmt.Errorf("unknown profile command %q\n\n%s", args[0], profileHelp)
	}
}

func (a *App) runProfileListCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("profile list")
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
	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/ad_accounts/"+accountID+"/profiles", query, nil, *all)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"name", "id", "reddit_user_id", "business_id", "modified_at"})
}

func (a *App) runProfileBusinessListCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("profile business-list")
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
	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/businesses/"+businessID+"/profiles", query, nil, *all)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"name", "id", "reddit_user_id", "business_id", "modified_at"})
}

func (a *App) resolveProfileSelection(ctx context.Context, accountID, input string) (string, string, error) {
	query := url.Values{}
	query.Set("page.size", "1000")

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/ad_accounts/"+accountID+"/profiles", query, nil, true)
	if err != nil {
		return "", "", err
	}

	rows := rowsFromPayload(payload)
	if len(rows) == 0 {
		return "", "", errors.New("no profiles returned for the selected ad account")
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
		return "", "", fmt.Errorf("profile name %q matched multiple profiles: %s. use the ID instead", input, strings.Join(ids, ", "))
	}

	return "", "", fmt.Errorf("could not find profile %q in the selected account. run `rad profile list` and use the ID from the `id` column", input)
}

const profileHelp = `Usage:
  rad profile list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
  rad profile business-list [--business-id <id-or-name>] [--all] [--page-size N] [--json]

Examples:
  rad profile list
  rad profile business-list`
