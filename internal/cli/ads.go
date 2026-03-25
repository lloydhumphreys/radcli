package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"radcli/internal/output"
)

var adDefinition = assetDefinition{
	Command:            "ad",
	Label:              "ad",
	CollectionEndpoint: "ads",
	ItemEndpoint:       "ads",
	ListColumns:        []string{"id", "name", "campaign_id", "ad_group_id", "configured_status", "effective_status"},
}

func (a *App) runAdCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(a.stdout, adHelp)
		return err
	}

	switch args[0] {
	case "list":
		return a.runAssetListCommand(ctx, adDefinition, args[1:])
	case "get":
		return a.runAssetGetCommand(ctx, adDefinition, args[1:])
	case "create":
		return a.runAdCreateCommand(ctx, args[1:])
	case "update":
		return a.runAdUpdateCommand(ctx, args[1:])
	default:
		return fmt.Errorf("unknown ad command %q\n\n%s", args[0], adHelp)
	}
}

func (a *App) runAdCreateCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("ad create")
	accountInput := fs.String("account-id", "", "")
	adGroupInput := fs.String("ad-group", "", "")
	name := fs.String("name", "", "")
	configuredStatus := fs.String("configured-status", "", "")
	postID := fs.String("post-id", "", "")
	clickURL := fs.String("click-url", "", "")
	shoppingCreativeInput := fs.String("shopping-creative-json", "", "")
	dryRun := fs.Bool("dry-run", false, "")
	jsonOut := fs.Bool("json", false, "")
	var queryParams stringList
	fs.Var(&queryParams, "click-url-query-parameter", "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if *adGroupInput == "" || *name == "" || *configuredStatus == "" {
		return errors.New("ad create requires --ad-group, --name, and --configured-status")
	}

	var adGroupID string
	if *dryRun && looksLikeID(*adGroupInput) {
		adGroupID = *adGroupInput
	} else {
		accountID, _, err := a.selectedAccountID(ctx, *accountInput)
		if err != nil {
			return err
		}
		adGroupID, _, err = a.resolveAssetSelection(ctx, accountID, adGroupDefinition, *adGroupInput)
		if err != nil {
			return err
		}
	}

	data, err := adWriteData(adWriteOptions{
		AdGroupID:            adGroupID,
		Name:                 *name,
		ConfiguredStatus:     *configuredStatus,
		PostID:               *postID,
		ClickURL:             *clickURL,
		QueryParams:          []string(queryParams),
		ShoppingCreativeJSON: *shoppingCreativeInput,
		RequireIdentity:      true,
	})
	if err != nil {
		return err
	}

	payload := map[string]any{"data": data}
	if *dryRun {
		return output.PrintJSON(a.stdout, payload)
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}
	response, err := a.api.RequestJSON(ctx, "POST", "/ad_accounts/"+accountID+"/ads", nil, payload)
	if err != nil {
		return err
	}
	if !*jsonOut {
		if _, err := fmt.Fprintln(a.stdout, "Ad created."); err != nil {
			return err
		}
	}
	return output.PrintJSON(a.stdout, dataOrSelf(response))
}

func (a *App) runAdUpdateCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("ad update")
	accountInput := fs.String("account-id", "", "")
	adGroupInput := fs.String("ad-group", "", "")
	name := fs.String("name", "", "")
	configuredStatus := fs.String("configured-status", "", "")
	postID := fs.String("post-id", "", "")
	clickURL := fs.String("click-url", "", "")
	shoppingCreativeInput := fs.String("shopping-creative-json", "", "")
	dryRun := fs.Bool("dry-run", false, "")
	jsonOut := fs.Bool("json", false, "")
	var queryParams stringList
	fs.Var(&queryParams, "click-url-query-parameter", "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return errors.New("usage: rad ad update <id-or-name> [flags]")
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}
	adID, adName, err := a.resolveAssetSelection(ctx, accountID, adDefinition, fs.Arg(0))
	if err != nil {
		return err
	}

	adGroupID := ""
	if *adGroupInput != "" {
		adGroupID, _, err = a.resolveAssetSelection(ctx, accountID, adGroupDefinition, *adGroupInput)
		if err != nil {
			return err
		}
	}

	data, err := adWriteData(adWriteOptions{
		AdGroupID:            adGroupID,
		Name:                 *name,
		ConfiguredStatus:     *configuredStatus,
		PostID:               *postID,
		ClickURL:             *clickURL,
		QueryParams:          []string(queryParams),
		ShoppingCreativeJSON: *shoppingCreativeInput,
		RequireIdentity:      false,
	})
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("ad update requires at least one field to change")
	}

	payload := map[string]any{"data": data}
	if *dryRun {
		return output.PrintJSON(a.stdout, payload)
	}

	response, err := a.api.RequestJSON(ctx, "PATCH", "/ads/"+adID, nil, payload)
	if err != nil {
		return err
	}
	if !*jsonOut {
		if _, err := fmt.Fprintf(a.stdout, "Ad updated: %s (%s)\n\n", adName, adID); err != nil {
			return err
		}
	}
	return output.PrintJSON(a.stdout, dataOrSelf(response))
}

type adWriteOptions struct {
	AdGroupID            string
	Name                 string
	ConfiguredStatus     string
	PostID               string
	ClickURL             string
	QueryParams          []string
	ShoppingCreativeJSON string
	RequireIdentity      bool
}

func adWriteData(opts adWriteOptions) (map[string]any, error) {
	data := map[string]any{}

	if opts.AdGroupID != "" {
		data["ad_group_id"] = opts.AdGroupID
	}
	if opts.Name != "" {
		data["name"] = opts.Name
	}
	if opts.ConfiguredStatus != "" {
		data["configured_status"] = strings.ToUpper(opts.ConfiguredStatus)
	}
	if opts.PostID != "" {
		data["post_id"] = opts.PostID
	}
	if opts.ClickURL != "" {
		data["click_url"] = opts.ClickURL
	}
	if len(opts.QueryParams) > 0 {
		queryParams, err := parseNamedValues(opts.QueryParams, "--click-url-query-parameter")
		if err != nil {
			return nil, err
		}
		if len(queryParams) > 0 {
			data["click_url_query_parameters"] = queryParams
		}
	}
	if opts.ShoppingCreativeJSON != "" {
		value, err := parseJSONInput(opts.ShoppingCreativeJSON, "--shopping-creative-json")
		if err != nil {
			return nil, err
		}
		data["shopping_creative"] = value
	}

	if opts.RequireIdentity {
		if opts.AdGroupID == "" || opts.Name == "" || opts.ConfiguredStatus == "" {
			return nil, errors.New("ad create requires --ad-group, --name, and --configured-status")
		}
	}

	return data, nil
}

func parseNamedValues(values []string, flagName string) ([]map[string]string, error) {
	out := make([]map[string]string, 0, len(values))
	for _, raw := range values {
		key, value, found := splitKeyValue(raw)
		if !found || key == "" {
			return nil, fmt.Errorf("invalid %s value %q: use name=value", flagName, raw)
		}
		out = append(out, map[string]string{
			"name":  key,
			"value": value,
		})
	}
	return out, nil
}

const adHelp = `Usage:
  rad ad list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
  rad ad get <id-or-name> [--account-id <id-or-name>] [--json]
  rad ad create --ad-group <id-or-name> --name <name> --configured-status <status> [--post-id <id>] [--click-url <url>] [--click-url-query-parameter <name=value> ...] [--shopping-creative-json <json-or-@file>] [--account-id <id-or-name>] [--dry-run]
  rad ad update <id-or-name> [--ad-group <id-or-name>] [--name <name>] [--configured-status <status>] [--post-id <id>] [--click-url <url>] [--click-url-query-parameter <name=value> ...] [--shopping-creative-json <json-or-@file>] [--account-id <id-or-name>] [--dry-run]

Examples:
  rad ad create --ad-group "US Retargeting" --name "Spring Ad" --configured-status PAUSED --post-id t3_abcdef --dry-run
  rad ad update "Spring Ad" --click-url https://example.com --click-url-query-parameter utm_source=reddit
  rad ad update "Catalog Ad" --shopping-creative-json @shopping-creative.json --dry-run`
