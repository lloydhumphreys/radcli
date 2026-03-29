package cli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/lloydhumphreys/radcli/internal/output"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type assetDefinition struct {
	Command            string
	Label              string
	CollectionEndpoint string
	ItemEndpoint       string
	ListColumns        []string
}

func (a *App) runAssetCommand(ctx context.Context, def assetDefinition, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintf(a.stdout, "Usage:\n  rad %s list [--account-id <id-or-name>] [--all] [--page-size N] [--json]\n  rad %s get <id-or-name> [--account-id <id-or-name>] [--json]\n", def.Command, def.Command)
		return err
	}

	switch args[0] {
	case "list":
		return a.runAssetListCommand(ctx, def, args[1:])
	case "get":
		return a.runAssetGetCommand(ctx, def, args[1:])
	default:
		return fmt.Errorf("unknown %s command %q", def.Command, args[0])
	}
}

func (a *App) runAssetListCommand(ctx context.Context, def assetDefinition, args []string) error {
	fs := newFlagSet(def.Command + " list")
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

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/ad_accounts/"+accountID+"/"+def.CollectionEndpoint, query, nil, *all)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rowsFromPayload(payload), def.ListColumns)
}

func (a *App) runAssetGetCommand(ctx context.Context, def assetDefinition, args []string) error {
	fs := newFlagSet(def.Command + " get")
	accountInput := fs.String("account-id", "", "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: rad %s get <id-or-name>", def.Command)
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}

	id, name, err := a.resolveAssetSelection(ctx, accountID, def, fs.Arg(0))
	if err != nil {
		return err
	}

	payload, err := a.api.RequestJSON(ctx, "GET", "/"+def.ItemEndpoint+"/"+id, nil, nil)
	if err != nil {
		return err
	}

	if !*jsonOut {
		if _, err := fmt.Fprintf(a.stdout, "%s %s (%s)\n\n", cases.Title(language.English).String(def.Label), name, id); err != nil {
			return err
		}
	}
	return output.PrintJSON(a.stdout, dataOrSelf(payload))
}

func (a *App) selectedBusinessID(ctx context.Context, input string) (string, string, error) {
	if input == "" {
		input = a.store.Config.Defaults.BusinessID
	}
	if input == "" {
		return "", "", errors.New("no business selected. use `rad business use <business-id-or-name>` or pass --business-id")
	}
	return a.resolveBusinessSelection(ctx, input)
}

func (a *App) selectedAccountID(ctx context.Context, input string) (string, string, error) {
	if input == "" {
		input = a.store.Config.Defaults.AdAccountID
	}
	if input == "" {
		return "", "", errors.New("no ad account selected. use `rad account use <ad-account-id-or-name>` or pass --account-id")
	}

	businessID, _, err := a.selectedBusinessID(ctx, "")
	if err != nil {
		return "", "", err
	}
	return a.resolveAccountSelection(ctx, businessID, input)
}

func (a *App) resolveAssetSelection(ctx context.Context, accountID string, def assetDefinition, input string) (string, string, error) {
	query := url.Values{}
	query.Set("page.size", "1000")

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/ad_accounts/"+accountID+"/"+def.CollectionEndpoint, query, nil, true)
	if err != nil {
		return "", "", err
	}

	rows := rowsFromPayload(payload)
	if len(rows) == 0 {
		return "", "", fmt.Errorf("no %ss returned for the selected account", def.Label)
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
		return "", "", fmt.Errorf("%s name %q matched multiple %ss: %s. use the ID instead", def.Label, input, def.Label, strings.Join(ids, ", "))
	}

	return "", "", fmt.Errorf("could not find %s %q in the selected account. run `rad %s list` and use the ID from the `id` column", def.Label, input, def.Command)
}
