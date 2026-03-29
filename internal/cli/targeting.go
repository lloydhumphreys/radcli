package cli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/lloydhumphreys/radcli/internal/output"
)

func (a *App) runTargetingCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintln(a.stdout, targetingHelp)
		return err
	}

	switch args[0] {
	case "communities":
		return a.runTargetingCommunities(ctx, args[1:])
	case "interests":
		return a.runTargetingInterests(ctx, args[1:])
	case "geolocations":
		return a.runTargetingGeolocations(ctx, args[1:])
	case "devices":
		return a.runTargetingDevices(ctx, args[1:])
	case "carriers":
		return a.runTargetingCarriers(ctx, args[1:])
	case "keywords":
		return a.runTargetingKeywords(ctx, args[1:])
	default:
		return fmt.Errorf("unknown targeting command %q\n\n%s", args[0], targetingHelp)
	}
}

func (a *App) runTargetingCommunities(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintln(a.stdout, targetingCommunitiesHelp)
		return err
	}

	switch args[0] {
	case "search":
		fs := newFlagSet("targeting communities search")
		queryText := fs.String("query", "", "")
		all := fs.Bool("all", false, "")
		pageSize := fs.Int("page-size", 0, "")
		jsonOut := fs.Bool("json", false, "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}
		if *queryText == "" {
			return errors.New("targeting communities search requires --query")
		}

		query := url.Values{}
		query.Set("query", *queryText)
		if *pageSize > 0 {
			query.Set("page.size", fmt.Sprintf("%d", *pageSize))
		}

		payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/targeting/communities/search", query, nil, *all)
		if err != nil {
			return err
		}
		if *jsonOut {
			return output.PrintJSON(a.stdout, payload)
		}
		return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"name", "id", "subscriber_count", "categories", "description"})
	case "list":
		fs := newFlagSet("targeting communities list")
		all := fs.Bool("all", false, "")
		pageSize := fs.Int("page-size", 0, "")
		jsonOut := fs.Bool("json", false, "")
		var names stringList
		fs.Var(&names, "name", "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}

		query := url.Values{}
		if len(names) > 0 {
			query.Set("names", joinCSV([]string(names)))
		}
		if *pageSize > 0 {
			query.Set("page.size", fmt.Sprintf("%d", *pageSize))
		}

		payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/targeting/communities", query, nil, *all)
		if err != nil {
			return err
		}
		if *jsonOut {
			return output.PrintJSON(a.stdout, payload)
		}
		return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"name", "id", "subscriber_count", "categories", "description"})
	case "suggest":
		fs := newFlagSet("targeting communities suggest")
		websiteURL := fs.String("website-url", "", "")
		all := fs.Bool("all", false, "")
		pageSize := fs.Int("page-size", 0, "")
		jsonOut := fs.Bool("json", false, "")
		var names stringList
		fs.Var(&names, "name", "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}
		if len(names) == 0 && *websiteURL == "" {
			return errors.New("targeting communities suggest requires --name and/or --website-url")
		}

		query := url.Values{}
		if len(names) > 0 {
			query.Set("names", joinCSV([]string(names)))
		}
		if *websiteURL != "" {
			query.Set("website_url", *websiteURL)
		}
		if *pageSize > 0 {
			query.Set("page.size", fmt.Sprintf("%d", *pageSize))
		}

		payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/targeting/communities/suggestions", query, nil, *all)
		if err != nil {
			return err
		}
		if *jsonOut {
			return output.PrintJSON(a.stdout, payload)
		}
		return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"name", "id", "subscriber_count", "categories", "description"})
	default:
		return fmt.Errorf("unknown targeting communities command %q\n\n%s", args[0], targetingCommunitiesHelp)
	}
}

func (a *App) runTargetingInterests(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintln(a.stdout, targetingInterestsHelp)
		return err
	}
	if args[0] != "list" {
		return fmt.Errorf("unknown targeting interests command %q\n\n%s", args[0], targetingInterestsHelp)
	}

	fs := newFlagSet("targeting interests list")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args[1:]); err != nil {
		return err
	}

	payload, err := a.api.RequestJSON(ctx, "GET", "/targeting/interests", nil, nil)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"category", "name", "id"})
}

func (a *App) runTargetingGeolocations(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintln(a.stdout, targetingGeolocationsHelp)
		return err
	}
	if args[0] != "search" {
		return fmt.Errorf("unknown targeting geolocations command %q\n\n%s", args[0], targetingGeolocationsHelp)
	}

	fs := newFlagSet("targeting geolocations search")
	postalCode := fs.String("postal-code", "", "")
	city := fs.String("city", "", "")
	country := fs.String("country", "", "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args[1:]); err != nil {
		return err
	}
	if *postalCode == "" && *city == "" && *country == "" {
		return errors.New("targeting geolocations search requires at least one of --postal-code, --city, or --country")
	}

	query := url.Values{}
	if *postalCode != "" {
		query.Set("postal_code", *postalCode)
	}
	if *city != "" {
		query.Set("cities_search", *city)
	}
	if *country != "" {
		query.Set("country", *country)
	}

	payload, err := a.api.RequestJSON(ctx, "GET", "/targeting/geolocations", query, nil)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"country", "region", "city", "postal_code", "dma", "id", "name"})
}

func (a *App) runTargetingDevices(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintln(a.stdout, targetingDevicesHelp)
		return err
	}
	if args[0] != "list" {
		return fmt.Errorf("unknown targeting devices command %q\n\n%s", args[0], targetingDevicesHelp)
	}

	fs := newFlagSet("targeting devices list")
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

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/targeting/devices", query, nil, *all)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"make", "model"})
}

func (a *App) runTargetingCarriers(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintln(a.stdout, targetingCarriersHelp)
		return err
	}
	if args[0] != "list" {
		return fmt.Errorf("unknown targeting carriers command %q\n\n%s", args[0], targetingCarriersHelp)
	}

	fs := newFlagSet("targeting carriers list")
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

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/targeting/carriers", query, nil, *all)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"country_code", "name", "id"})
}

func (a *App) runTargetingKeywords(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintln(a.stdout, targetingKeywordsHelp)
		return err
	}

	switch args[0] {
	case "suggest":
		fs := newFlagSet("targeting keywords suggest")
		jsonOut := fs.Bool("json", false, "")
		var keywords stringList
		fs.Var(&keywords, "keyword", "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}
		if len(keywords) == 0 {
			return errors.New("targeting keywords suggest requires at least one --keyword")
		}

		payload, err := a.api.RequestJSON(ctx, "POST", "/targeting/keyword_suggestions", nil, map[string]any{
			"data": map[string]any{
				"seed_keywords": []string(keywords),
			},
		})
		if err != nil {
			return err
		}

		suggestions := nestedStringRows(payload, "data", "keyword_suggestions")
		if *jsonOut {
			return output.PrintJSON(a.stdout, payload)
		}
		return output.PrintTable(a.stdout, suggestions, []string{"keyword"})
	case "validate":
		fs := newFlagSet("targeting keywords validate")
		jsonOut := fs.Bool("json", false, "")
		var keywords stringList
		fs.Var(&keywords, "keyword", "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}
		if len(keywords) == 0 {
			return errors.New("targeting keywords validate requires at least one --keyword")
		}

		payload, err := a.api.RequestJSON(ctx, "POST", "/targeting/keyword_validations", nil, map[string]any{
			"data": map[string]any{
				"keywords": []string(keywords),
			},
		})
		if err != nil {
			return err
		}
		if *jsonOut {
			return output.PrintJSON(a.stdout, payload)
		}
		return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"keyword", "is_brand_safe"})
	default:
		return fmt.Errorf("unknown targeting keywords command %q\n\n%s", args[0], targetingKeywordsHelp)
	}
}

func nestedStringRows(payload map[string]any, objectKey, arrayKey string) []map[string]string {
	object, _ := payload[objectKey].(map[string]any)
	items, _ := object[arrayKey].([]any)
	rows := make([]map[string]string, 0, len(items))
	for _, item := range items {
		if value, ok := item.(string); ok {
			rows = append(rows, map[string]string{"keyword": value})
		}
	}
	return rows
}

func joinCSV(values []string) string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" {
			out = append(out, value)
		}
	}
	return strings.Join(out, ",")
}

const targetingHelp = `Usage:
  rad targeting communities search --query <text> [--all] [--page-size N] [--json]
  rad targeting communities list [--name <community> ...] [--all] [--page-size N] [--json]
  rad targeting communities suggest [--name <community> ...] [--website-url <url>] [--all] [--page-size N] [--json]
  rad targeting interests list [--json]
  rad targeting geolocations search [--postal-code <code>] [--city <name>] [--country <code>] [--json]
  rad targeting devices list [--all] [--page-size N] [--json]
  rad targeting carriers list [--all] [--page-size N] [--json]
  rad targeting keywords suggest --keyword <term> [--keyword <term> ...] [--json]
  rad targeting keywords validate --keyword <term> [--keyword <term> ...] [--json]`

const targetingCommunitiesHelp = `Usage:
  rad targeting communities search --query <text> [--all] [--page-size N] [--json]
  rad targeting communities list [--name <community> ...] [--all] [--page-size N] [--json]
  rad targeting communities suggest [--name <community> ...] [--website-url <url>] [--all] [--page-size N] [--json]`

const targetingInterestsHelp = `Usage:
  rad targeting interests list [--json]`

const targetingGeolocationsHelp = `Usage:
  rad targeting geolocations search [--postal-code <code>] [--city <name>] [--country <code>] [--json]`

const targetingDevicesHelp = `Usage:
  rad targeting devices list [--all] [--page-size N] [--json]`

const targetingCarriersHelp = `Usage:
  rad targeting carriers list [--all] [--page-size N] [--json]`

const targetingKeywordsHelp = `Usage:
  rad targeting keywords suggest --keyword <term> [--keyword <term> ...] [--json]
  rad targeting keywords validate --keyword <term> [--keyword <term> ...] [--json]`
