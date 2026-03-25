package cli

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"radcli/internal/output"
)

var campaignDefinition = assetDefinition{
	Command:            "campaign",
	Label:              "campaign",
	CollectionEndpoint: "campaigns",
	ItemEndpoint:       "campaigns",
	ListColumns:        []string{"id", "name", "configured_status", "effective_status", "objective", "budget_amount"},
}

func (a *App) runCampaignCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(a.stdout, campaignHelp)
		return err
	}

	switch args[0] {
	case "list":
		return a.runAssetListCommand(ctx, campaignDefinition, args[1:])
	case "get":
		return a.runAssetGetCommand(ctx, campaignDefinition, args[1:])
	case "create":
		return a.runCampaignCreateCommand(ctx, args[1:])
	case "update":
		return a.runCampaignUpdateCommand(ctx, args[1:])
	default:
		return fmt.Errorf("unknown campaign command %q\n\n%s", args[0], campaignHelp)
	}
}

func (a *App) runCampaignCreateCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("campaign create")
	accountInput := fs.String("account-id", "", "")
	name := fs.String("name", "", "")
	objective := fs.String("objective", "", "")
	configuredStatus := fs.String("configured-status", "", "")
	fundingInstrumentID := fs.String("funding-instrument-id", "", "")
	invoiceLabel := fs.String("invoice-label", "", "")
	cbo := fs.String("campaign-budget-optimization", "", "")
	goalType := fs.String("goal-type", "", "")
	goalValue := fs.String("goal-value", "", "")
	spendCap := fs.String("spend-cap", "", "")
	startTime := fs.String("start-time", "", "")
	endTime := fs.String("end-time", "", "")
	bidStrategy := fs.String("bid-strategy", "", "")
	bidType := fs.String("bid-type", "", "")
	bidValue := fs.String("bid-value", "", "")
	appID := fs.String("app-id", "", "")
	conversionPixelID := fs.String("conversion-pixel-id", "", "")
	viewThroughConversionType := fs.String("view-through-conversion-type", "", "")
	dryRun := fs.Bool("dry-run", false, "")
	jsonOut := fs.Bool("json", false, "")
	var specialAdCategories stringList
	fs.Var(&specialAdCategories, "special-ad-category", "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}

	if *name == "" || *objective == "" || *configuredStatus == "" {
		return errors.New("campaign create requires --name, --objective, and --configured-status")
	}

	data, err := campaignWriteData(campaignWriteOptions{
		Name:                       *name,
		Objective:                  *objective,
		ConfiguredStatus:           *configuredStatus,
		FundingInstrumentID:        *fundingInstrumentID,
		InvoiceLabel:               *invoiceLabel,
		SpecialAdCategories:        []string(specialAdCategories),
		CampaignBudgetOptimization: *cbo,
		GoalType:                   *goalType,
		GoalValue:                  *goalValue,
		SpendCap:                   *spendCap,
		StartTime:                  *startTime,
		EndTime:                    *endTime,
		BidStrategy:                *bidStrategy,
		BidType:                    *bidType,
		BidValue:                   *bidValue,
		AppID:                      *appID,
		ConversionPixelID:          *conversionPixelID,
		ViewThroughConversionType:  *viewThroughConversionType,
		RequireIdentity:            true,
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

	response, err := a.api.RequestJSON(ctx, "POST", "/ad_accounts/"+accountID+"/campaigns", nil, payload)
	if err != nil {
		return err
	}
	if !*jsonOut {
		if _, err := fmt.Fprintln(a.stdout, "Campaign created.\n"); err != nil {
			return err
		}
	}
	return output.PrintJSON(a.stdout, dataOrSelf(response))
}

func (a *App) runCampaignUpdateCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("campaign update")
	accountInput := fs.String("account-id", "", "")
	name := fs.String("name", "", "")
	objective := fs.String("objective", "", "")
	configuredStatus := fs.String("configured-status", "", "")
	fundingInstrumentID := fs.String("funding-instrument-id", "", "")
	invoiceLabel := fs.String("invoice-label", "", "")
	cbo := fs.String("campaign-budget-optimization", "", "")
	goalType := fs.String("goal-type", "", "")
	goalValue := fs.String("goal-value", "", "")
	spendCap := fs.String("spend-cap", "", "")
	startTime := fs.String("start-time", "", "")
	endTime := fs.String("end-time", "", "")
	bidStrategy := fs.String("bid-strategy", "", "")
	bidType := fs.String("bid-type", "", "")
	bidValue := fs.String("bid-value", "", "")
	appID := fs.String("app-id", "", "")
	conversionPixelID := fs.String("conversion-pixel-id", "", "")
	viewThroughConversionType := fs.String("view-through-conversion-type", "", "")
	dryRun := fs.Bool("dry-run", false, "")
	jsonOut := fs.Bool("json", false, "")
	var specialAdCategories stringList
	fs.Var(&specialAdCategories, "special-ad-category", "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return errors.New("usage: rad campaign update <id-or-name> [flags]")
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}
	campaignID, campaignName, err := a.resolveAssetSelection(ctx, accountID, campaignDefinition, fs.Arg(0))
	if err != nil {
		return err
	}

	data, err := campaignWriteData(campaignWriteOptions{
		Name:                       *name,
		Objective:                  *objective,
		ConfiguredStatus:           *configuredStatus,
		FundingInstrumentID:        *fundingInstrumentID,
		InvoiceLabel:               *invoiceLabel,
		SpecialAdCategories:        []string(specialAdCategories),
		CampaignBudgetOptimization: *cbo,
		GoalType:                   *goalType,
		GoalValue:                  *goalValue,
		SpendCap:                   *spendCap,
		StartTime:                  *startTime,
		EndTime:                    *endTime,
		BidStrategy:                *bidStrategy,
		BidType:                    *bidType,
		BidValue:                   *bidValue,
		AppID:                      *appID,
		ConversionPixelID:          *conversionPixelID,
		ViewThroughConversionType:  *viewThroughConversionType,
		RequireIdentity:            false,
	})
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("campaign update requires at least one field to change")
	}

	payload := map[string]any{"data": data}
	if *dryRun {
		return output.PrintJSON(a.stdout, payload)
	}

	response, err := a.api.RequestJSON(ctx, "PATCH", "/campaigns/"+campaignID, nil, payload)
	if err != nil {
		return err
	}
	if !*jsonOut {
		if _, err := fmt.Fprintf(a.stdout, "Campaign updated: %s (%s)\n\n", campaignName, campaignID); err != nil {
			return err
		}
	}
	return output.PrintJSON(a.stdout, dataOrSelf(response))
}

type campaignWriteOptions struct {
	Name                       string
	Objective                  string
	ConfiguredStatus           string
	FundingInstrumentID        string
	InvoiceLabel               string
	SpecialAdCategories        []string
	CampaignBudgetOptimization string
	GoalType                   string
	GoalValue                  string
	SpendCap                   string
	StartTime                  string
	EndTime                    string
	BidStrategy                string
	BidType                    string
	BidValue                   string
	AppID                      string
	ConversionPixelID          string
	ViewThroughConversionType  string
	RequireIdentity            bool
}

func campaignWriteData(opts campaignWriteOptions) (map[string]any, error) {
	data := map[string]any{}

	if opts.Name != "" {
		data["name"] = opts.Name
	}
	if opts.Objective != "" {
		data["objective"] = strings.ToUpper(opts.Objective)
	}
	if opts.ConfiguredStatus != "" {
		data["configured_status"] = strings.ToUpper(opts.ConfiguredStatus)
	}
	if opts.FundingInstrumentID != "" {
		data["funding_instrument_id"] = opts.FundingInstrumentID
	}
	if opts.InvoiceLabel != "" {
		data["invoice_label"] = opts.InvoiceLabel
	}
	if len(opts.SpecialAdCategories) > 0 {
		categories := make([]string, 0, len(opts.SpecialAdCategories))
		for _, category := range opts.SpecialAdCategories {
			category = strings.TrimSpace(category)
			if category == "" {
				continue
			}
			categories = append(categories, strings.ToUpper(category))
		}
		if len(categories) > 0 {
			data["special_ad_categories"] = categories
		}
	}
	if opts.CampaignBudgetOptimization != "" {
		value, err := strconv.ParseBool(opts.CampaignBudgetOptimization)
		if err != nil {
			return nil, fmt.Errorf("invalid --campaign-budget-optimization value %q: use true or false", opts.CampaignBudgetOptimization)
		}
		data["is_campaign_budget_optimization"] = value
	}
	if opts.GoalType != "" {
		data["goal_type"] = strings.ToUpper(opts.GoalType)
	}
	if opts.GoalValue != "" {
		value, err := parseMoneyToMicros(opts.GoalValue, "--goal-value")
		if err != nil {
			return nil, err
		}
		data["goal_value"] = value
	}
	if opts.SpendCap != "" {
		value, err := parseMoneyToMicros(opts.SpendCap, "--spend-cap")
		if err != nil {
			return nil, err
		}
		data["spend_cap"] = value
	}
	if opts.StartTime != "" {
		data["start_time"] = opts.StartTime
	}
	if opts.EndTime != "" {
		data["end_time"] = opts.EndTime
	}
	if opts.BidStrategy != "" {
		data["bid_strategy"] = strings.ToUpper(opts.BidStrategy)
	}
	if opts.BidType != "" {
		data["bid_type"] = strings.ToUpper(opts.BidType)
	}
	if opts.BidValue != "" {
		value, err := parseMoneyToMicros(opts.BidValue, "--bid-value")
		if err != nil {
			return nil, err
		}
		data["bid_value"] = value
	}
	if opts.AppID != "" {
		data["app_id"] = opts.AppID
	}
	if opts.ConversionPixelID != "" {
		data["conversion_pixel_id"] = opts.ConversionPixelID
	}
	if opts.ViewThroughConversionType != "" {
		data["view_through_conversion_type"] = strings.ToUpper(opts.ViewThroughConversionType)
	}

	if (opts.StartTime == "") != (opts.EndTime == "") {
		return nil, errors.New("use both --start-time and --end-time together")
	}
	if opts.RequireIdentity {
		if opts.Name == "" || opts.Objective == "" || opts.ConfiguredStatus == "" {
			return nil, errors.New("campaign create requires --name, --objective, and --configured-status")
		}
	}

	return data, nil
}

func parseMoneyToMicros(raw, flagName string) (int64, error) {
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value %q", flagName, raw)
	}
	return int64(math.Round(value * 1_000_000)), nil
}

const campaignHelp = `Usage:
  rad campaign list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
  rad campaign get <id-or-name> [--account-id <id-or-name>] [--json]
  rad campaign create --name <name> --objective <objective> --configured-status <status> [flags]
  rad campaign update <id-or-name> [flags]

Common write flags:
  --funding-instrument-id <id>
  --invoice-label <text>
  --special-ad-category <category>
  --campaign-budget-optimization <true|false>
  --goal-type <daily_spend|lifetime_spend>
  --goal-value <major-currency>
  --spend-cap <major-currency>
  --start-time <rfc3339>
  --end-time <rfc3339>
  --bid-strategy <bidless|maximize_volume|target_cpx>
  --bid-type <cpc|cpm|cpv6>
  --bid-value <major-currency>
  --app-id <store-app-id>
  --conversion-pixel-id <id>
  --view-through-conversion-type <seven_day_clicks|seven_day_clicks_one_day_view>
  --dry-run

Examples:
  rad campaign create --name "Spring Launch" --objective CLICKS --configured-status PAUSED --spend-cap 250
  rad campaign update "Spring Launch" --configured-status ACTIVE
  rad campaign update "Spring Launch" --invoice-label "client-2026-q2" --dry-run`
