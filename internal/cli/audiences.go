package cli

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"radcli/internal/output"
)

var savedAudienceDefinition = assetDefinition{
	Command:            "saved",
	Label:              "saved audience",
	CollectionEndpoint: "saved_audiences",
	ItemEndpoint:       "saved_audiences",
	ListColumns:        []string{"id", "name", "type", "status", "active_ad_groups_count", "size_range_lower", "size_range_upper", "updated_at"},
}

var customAudienceDefinition = assetDefinition{
	Command:            "custom",
	Label:              "custom audience",
	CollectionEndpoint: "custom_audiences",
	ItemEndpoint:       "custom_audiences",
	ListColumns:        []string{"id", "name", "type", "status", "delivery_status", "size_range_lower", "size_range_upper", "cost_partner", "cost_price", "cost_currency", "modified_at"},
}

func (a *App) runAudienceCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(a.stdout, audienceHelp)
		return err
	}

	switch args[0] {
	case "saved":
		return a.runSavedAudienceCommand(ctx, args[1:])
	case "custom":
		return a.runCustomAudienceCommand(ctx, args[1:])
	case "third-party":
		return a.runThirdPartyAudienceCommand(ctx, args[1:])
	default:
		return fmt.Errorf("unknown audience command %q\n\n%s", args[0], audienceHelp)
	}
}

func (a *App) runSavedAudienceCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(a.stdout, savedAudienceHelp)
		return err
	}

	switch args[0] {
	case "list":
		return a.runAssetListCommand(ctx, savedAudienceDefinition, args[1:])
	case "get":
		return a.runAssetGetCommand(ctx, savedAudienceDefinition, args[1:])
	case "create":
		return a.runSavedAudienceCreateCommand(ctx, args[1:])
	case "update":
		return a.runSavedAudienceUpdateCommand(ctx, args[1:])
	default:
		return fmt.Errorf("unknown audience saved command %q\n\n%s", args[0], savedAudienceHelp)
	}
}

func (a *App) runSavedAudienceCreateCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("audience saved create")
	accountInput := fs.String("account-id", "", "")
	name := fs.String("name", "", "")
	audienceType := fs.String("type", "", "")
	targetingInput := fs.String("targeting-json", "", "")
	dryRun := fs.Bool("dry-run", false, "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if *name == "" || *audienceType == "" || *targetingInput == "" {
		return errors.New("audience saved create requires --name, --type, and --targeting-json")
	}

	targeting, err := parseJSONInput(*targetingInput, "--targeting-json")
	if err != nil {
		return err
	}
	payload := map[string]any{
		"data": map[string]any{
			"name":      *name,
			"type":      *audienceType,
			"targeting": targeting,
		},
	}
	if *dryRun {
		return output.PrintJSON(a.stdout, payload)
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}
	response, err := a.api.RequestJSON(ctx, "POST", "/ad_accounts/"+accountID+"/saved_audiences", nil, payload)
	if err != nil {
		return err
	}
	if !*jsonOut {
		if _, err := fmt.Fprintln(a.stdout, "Saved audience created.\n"); err != nil {
			return err
		}
	}
	return output.PrintJSON(a.stdout, dataOrSelf(response))
}

func (a *App) runSavedAudienceUpdateCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("audience saved update")
	accountInput := fs.String("account-id", "", "")
	name := fs.String("name", "", "")
	audienceType := fs.String("type", "", "")
	targetingInput := fs.String("targeting-json", "", "")
	dryRun := fs.Bool("dry-run", false, "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return errors.New("usage: rad audience saved update <id-or-name> [flags]")
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}
	savedAudienceID, savedAudienceName, err := a.resolveAssetSelection(ctx, accountID, savedAudienceDefinition, fs.Arg(0))
	if err != nil {
		return err
	}

	data := map[string]any{}
	if *name != "" {
		data["name"] = *name
	}
	if *audienceType != "" {
		data["type"] = *audienceType
	}
	if *targetingInput != "" {
		targeting, err := parseJSONInput(*targetingInput, "--targeting-json")
		if err != nil {
			return err
		}
		data["targeting"] = targeting
	}
	if len(data) == 0 {
		return errors.New("audience saved update requires at least one field to change")
	}

	payload := map[string]any{"data": data}
	if *dryRun {
		return output.PrintJSON(a.stdout, payload)
	}
	response, err := a.api.RequestJSON(ctx, "PATCH", "/saved_audiences/"+savedAudienceID, nil, payload)
	if err != nil {
		return err
	}
	if !*jsonOut {
		if _, err := fmt.Fprintf(a.stdout, "Saved audience updated: %s (%s)\n\n", savedAudienceName, savedAudienceID); err != nil {
			return err
		}
	}
	return output.PrintJSON(a.stdout, dataOrSelf(response))
}

func (a *App) runCustomAudienceCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(a.stdout, customAudienceHelp)
		return err
	}

	switch args[0] {
	case "list":
		fs := newFlagSet("audience custom list")
		accountInput := fs.String("account-id", "", "")
		name := fs.String("name", "", "")
		all := fs.Bool("all", false, "")
		pageSize := fs.Int("page-size", 0, "")
		jsonOut := fs.Bool("json", false, "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}
		accountID, _, err := a.selectedAccountID(ctx, *accountInput)
		if err != nil {
			return err
		}
		query := url.Values{}
		if *name != "" {
			query.Set("name", *name)
		}
		if *pageSize > 0 {
			query.Set("page.size", fmt.Sprintf("%d", *pageSize))
		}
		payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/ad_accounts/"+accountID+"/custom_audiences", query, nil, *all)
		if err != nil {
			return err
		}
		rows := formatCustomAudienceRows(rowsFromPayload(payload))
		if *jsonOut {
			return output.PrintJSON(a.stdout, payload)
		}
		return output.PrintTable(a.stdout, rows, customAudienceDefinition.ListColumns)
	case "get":
		fs := newFlagSet("audience custom get")
		accountInput := fs.String("account-id", "", "")
		jsonOut := fs.Bool("json", false, "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}
		if fs.NArg() < 1 {
			return errors.New("usage: rad audience custom get <id-or-name>")
		}
		accountID, _, err := a.selectedAccountID(ctx, *accountInput)
		if err != nil {
			return err
		}
		id, name, err := a.resolveAssetSelection(ctx, accountID, customAudienceDefinition, fs.Arg(0))
		if err != nil {
			return err
		}
		payload, err := a.api.RequestJSON(ctx, "GET", "/custom_audiences/"+id, nil, nil)
		if err != nil {
			return err
		}
		if !*jsonOut {
			if _, err := fmt.Fprintf(a.stdout, "Custom audience %s (%s)\n\n", name, id); err != nil {
				return err
			}
		}
		return output.PrintJSON(a.stdout, dataOrSelf(payload))
	default:
		return fmt.Errorf("unknown audience custom command %q\n\n%s", args[0], customAudienceHelp)
	}
}

func (a *App) runThirdPartyAudienceCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(a.stdout, thirdPartyAudienceHelp)
		return err
	}
	if args[0] != "list" {
		return fmt.Errorf("unknown audience third-party command %q\n\n%s", args[0], thirdPartyAudienceHelp)
	}

	fs := newFlagSet("audience third-party list")
	all := fs.Bool("all", false, "")
	pageSize := fs.Int("page-size", 0, "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args[1:]); err != nil {
		return err
	}
	query := url.Values{}
	if *pageSize > 0 {
		query.Set("page.size", fmt.Sprintf("%d", *pageSize))
	}
	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/targeting/third_party_audiences", query, nil, *all)
	if err != nil {
		return err
	}
	rows := formatThirdPartyAudienceRows(rowsFromPayload(payload))
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rows, []string{"full_name", "name", "id", "size", "partner", "price", "currency"})
}

func formatCustomAudienceRows(rows []map[string]string) []map[string]string {
	for _, row := range rows {
		if row["cost"] == "" {
			continue
		}
		costMap := stringKeyValueMap(row["cost"])
		if costMap["price"] != "" {
			row["cost_price"] = formatMicroCurrency(costMap["price"])
		}
		row["cost_currency"] = costMap["currency"]
		row["cost_partner"] = costMap["partner"]
	}
	return rows
}

func formatThirdPartyAudienceRows(rows []map[string]string) []map[string]string {
	for _, row := range rows {
		if row["cost"] == "" {
			continue
		}
		costMap := stringKeyValueMap(row["cost"])
		row["partner"] = costMap["partner"]
		row["currency"] = costMap["currency"]
		if costMap["price"] != "" {
			row["price"] = formatMicroCurrency(costMap["price"])
		}
	}
	return rows
}

func stringKeyValueMap(raw string) map[string]string {
	out := map[string]string{}
	for _, part := range splitCommaSeparatedPairs(raw) {
		key, value, found := splitKeyValue(part)
		if found {
			out[key] = value
		}
	}
	return out
}

func splitCommaSeparatedPairs(raw string) []string {
	if raw == "" {
		return nil
	}
	return splitRespectingBrackets(raw, ',')
}

func splitKeyValue(raw string) (string, string, bool) {
	parts := splitRespectingBrackets(raw, '=')
	if len(parts) < 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func splitRespectingBrackets(raw string, sep rune) []string {
	parts := make([]string, 0)
	current := make([]rune, 0, len(raw))
	depth := 0
	for _, r := range raw {
		switch r {
		case '{', '[':
			depth++
		case '}', ']':
			if depth > 0 {
				depth--
			}
		}
		if r == sep && depth == 0 {
			parts = append(parts, string(current))
			current = current[:0]
			continue
		}
		current = append(current, r)
	}
	parts = append(parts, string(current))
	return parts
}

const audienceHelp = `Usage:
  rad audience saved list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
  rad audience saved get <id-or-name> [--account-id <id-or-name>] [--json]
  rad audience saved create --name <name> --type <type> --targeting-json <json-or-@file> [--account-id <id-or-name>] [--dry-run]
  rad audience saved update <id-or-name> [--name <name>] [--type <type>] [--targeting-json <json-or-@file>] [--account-id <id-or-name>] [--dry-run]
  rad audience custom list [--account-id <id-or-name>] [--name <text>] [--all] [--page-size N] [--json]
  rad audience custom get <id-or-name> [--account-id <id-or-name>] [--json]
  rad audience third-party list [--all] [--page-size N] [--json]`

const savedAudienceHelp = `Usage:
  rad audience saved list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
  rad audience saved get <id-or-name> [--account-id <id-or-name>] [--json]
  rad audience saved create --name <name> --type <type> --targeting-json <json-or-@file> [--account-id <id-or-name>] [--dry-run]
  rad audience saved update <id-or-name> [--name <name>] [--type <type>] [--targeting-json <json-or-@file>] [--account-id <id-or-name>] [--dry-run]`

const customAudienceHelp = `Usage:
  rad audience custom list [--account-id <id-or-name>] [--name <text>] [--all] [--page-size N] [--json]
  rad audience custom get <id-or-name> [--account-id <id-or-name>] [--json]`

const thirdPartyAudienceHelp = `Usage:
  rad audience third-party list [--all] [--page-size N] [--json]`
