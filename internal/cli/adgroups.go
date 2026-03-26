package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/lloydhumphreys/radcli/internal/output"
)

var adGroupDefinition = assetDefinition{
	Command:            "adgroup",
	Label:              "ad group",
	CollectionEndpoint: "ad_groups",
	ItemEndpoint:       "ad_groups",
	ListColumns:        []string{"id", "name", "campaign_id", "configured_status", "effective_status", "bid_amount"},
}

func (a *App) runAdGroupCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(a.stdout, adGroupHelp)
		return err
	}

	switch args[0] {
	case "list":
		return a.runAssetListCommand(ctx, adGroupDefinition, args[1:])
	case "get":
		return a.runAssetGetCommand(ctx, adGroupDefinition, args[1:])
	case "create":
		return a.runAdGroupCreateCommand(ctx, args[1:])
	case "update":
		return a.runAdGroupUpdateCommand(ctx, args[1:])
	default:
		return fmt.Errorf("unknown adgroup command %q\n\n%s", args[0], adGroupHelp)
	}
}

func (a *App) runAdGroupCreateCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("adgroup create")
	accountInput := fs.String("account-id", "", "")
	campaignInput := fs.String("campaign", "", "")
	name := fs.String("name", "", "")
	configuredStatus := fs.String("configured-status", "", "")
	bidStrategy := fs.String("bid-strategy", "", "")
	bidType := fs.String("bid-type", "", "")
	bidValue := fs.String("bid-value", "", "")
	goalType := fs.String("goal-type", "", "")
	goalValue := fs.String("goal-value", "", "")
	optimizationGoal := fs.String("optimization-goal", "", "")
	optimizationStrategyType := fs.String("optimization-strategy-type", "", "")
	startTime := fs.String("start-time", "", "")
	endTime := fs.String("end-time", "", "")
	appID := fs.String("app-id", "", "")
	conversionPixelID := fs.String("conversion-pixel-id", "", "")
	viewThroughConversionType := fs.String("view-through-conversion-type", "", "")
	savedAudienceID := fs.String("saved-audience-id", "", "")
	productSetID := fs.String("product-set-id", "", "")
	shoppingType := fs.String("shopping-type", "", "")
	targetingInput := fs.String("targeting-json", "", "")
	scheduleInput := fs.String("schedule-json", "", "")
	shoppingTargetingInput := fs.String("shopping-targeting-json", "", "")
	dryRun := fs.Bool("dry-run", false, "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}

	if *campaignInput == "" || *name == "" || *configuredStatus == "" {
		return errors.New("adgroup create requires --campaign, --name, and --configured-status")
	}

	var campaignID string
	if *dryRun && looksLikeID(*campaignInput) {
		campaignID = *campaignInput
	} else {
		accountID, _, err := a.selectedAccountID(ctx, *accountInput)
		if err != nil {
			return err
		}
		campaignID, _, err = a.resolveAssetSelection(ctx, accountID, campaignDefinition, *campaignInput)
		if err != nil {
			return err
		}
	}

	data, err := adGroupWriteData(adGroupWriteOptions{
		CampaignID:                campaignID,
		Name:                      *name,
		ConfiguredStatus:          *configuredStatus,
		BidStrategy:               *bidStrategy,
		BidType:                   *bidType,
		BidValue:                  *bidValue,
		GoalType:                  *goalType,
		GoalValue:                 *goalValue,
		OptimizationGoal:          *optimizationGoal,
		OptimizationStrategyType:  *optimizationStrategyType,
		StartTime:                 *startTime,
		EndTime:                   *endTime,
		AppID:                     *appID,
		ConversionPixelID:         *conversionPixelID,
		ViewThroughConversionType: *viewThroughConversionType,
		SavedAudienceID:           *savedAudienceID,
		ProductSetID:              *productSetID,
		ShoppingType:              *shoppingType,
		TargetingInput:            *targetingInput,
		ScheduleInput:             *scheduleInput,
		ShoppingTargetingInput:    *shoppingTargetingInput,
		RequireIdentity:           true,
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

	response, err := a.api.RequestJSON(ctx, "POST", "/ad_accounts/"+accountID+"/ad_groups", nil, payload)
	if err != nil {
		return err
	}
	if !*jsonOut {
		if _, err := fmt.Fprintln(a.stdout, "Ad group created."); err != nil {
			return err
		}
	}
	return output.PrintJSON(a.stdout, dataOrSelf(response))
}

func (a *App) runAdGroupUpdateCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("adgroup update")
	accountInput := fs.String("account-id", "", "")
	campaignInput := fs.String("campaign", "", "")
	name := fs.String("name", "", "")
	configuredStatus := fs.String("configured-status", "", "")
	bidStrategy := fs.String("bid-strategy", "", "")
	bidType := fs.String("bid-type", "", "")
	bidValue := fs.String("bid-value", "", "")
	goalType := fs.String("goal-type", "", "")
	goalValue := fs.String("goal-value", "", "")
	optimizationGoal := fs.String("optimization-goal", "", "")
	optimizationStrategyType := fs.String("optimization-strategy-type", "", "")
	startTime := fs.String("start-time", "", "")
	endTime := fs.String("end-time", "", "")
	appID := fs.String("app-id", "", "")
	conversionPixelID := fs.String("conversion-pixel-id", "", "")
	viewThroughConversionType := fs.String("view-through-conversion-type", "", "")
	savedAudienceID := fs.String("saved-audience-id", "", "")
	productSetID := fs.String("product-set-id", "", "")
	shoppingType := fs.String("shopping-type", "", "")
	targetingInput := fs.String("targeting-json", "", "")
	scheduleInput := fs.String("schedule-json", "", "")
	shoppingTargetingInput := fs.String("shopping-targeting-json", "", "")
	dryRun := fs.Bool("dry-run", false, "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return errors.New("usage: rad adgroup update <id-or-name> [flags]")
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}
	adGroupID, adGroupName, err := a.resolveAssetSelection(ctx, accountID, adGroupDefinition, fs.Arg(0))
	if err != nil {
		return err
	}

	campaignID := ""
	if *campaignInput != "" {
		campaignID, _, err = a.resolveAssetSelection(ctx, accountID, campaignDefinition, *campaignInput)
		if err != nil {
			return err
		}
	}

	data, err := adGroupWriteData(adGroupWriteOptions{
		CampaignID:                campaignID,
		Name:                      *name,
		ConfiguredStatus:          *configuredStatus,
		BidStrategy:               *bidStrategy,
		BidType:                   *bidType,
		BidValue:                  *bidValue,
		GoalType:                  *goalType,
		GoalValue:                 *goalValue,
		OptimizationGoal:          *optimizationGoal,
		OptimizationStrategyType:  *optimizationStrategyType,
		StartTime:                 *startTime,
		EndTime:                   *endTime,
		AppID:                     *appID,
		ConversionPixelID:         *conversionPixelID,
		ViewThroughConversionType: *viewThroughConversionType,
		SavedAudienceID:           *savedAudienceID,
		ProductSetID:              *productSetID,
		ShoppingType:              *shoppingType,
		TargetingInput:            *targetingInput,
		ScheduleInput:             *scheduleInput,
		ShoppingTargetingInput:    *shoppingTargetingInput,
		RequireIdentity:           false,
	})
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("adgroup update requires at least one field to change")
	}

	payload := map[string]any{"data": data}
	if *dryRun {
		return output.PrintJSON(a.stdout, payload)
	}

	response, err := a.api.RequestJSON(ctx, "PATCH", "/ad_groups/"+adGroupID, nil, payload)
	if err != nil {
		return err
	}
	if !*jsonOut {
		if _, err := fmt.Fprintf(a.stdout, "Ad group updated: %s (%s)\n\n", adGroupName, adGroupID); err != nil {
			return err
		}
	}
	return output.PrintJSON(a.stdout, dataOrSelf(response))
}

type adGroupWriteOptions struct {
	CampaignID                string
	Name                      string
	ConfiguredStatus          string
	BidStrategy               string
	BidType                   string
	BidValue                  string
	GoalType                  string
	GoalValue                 string
	OptimizationGoal          string
	OptimizationStrategyType  string
	StartTime                 string
	EndTime                   string
	AppID                     string
	ConversionPixelID         string
	ViewThroughConversionType string
	SavedAudienceID           string
	ProductSetID              string
	ShoppingType              string
	TargetingInput            string
	ScheduleInput             string
	ShoppingTargetingInput    string
	RequireIdentity           bool
}

func adGroupWriteData(opts adGroupWriteOptions) (map[string]any, error) {
	data := map[string]any{}

	if opts.CampaignID != "" {
		data["campaign_id"] = opts.CampaignID
	}
	if opts.Name != "" {
		data["name"] = opts.Name
	}
	if opts.ConfiguredStatus != "" {
		data["configured_status"] = strings.ToUpper(opts.ConfiguredStatus)
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
	if opts.OptimizationGoal != "" {
		data["optimization_goal"] = strings.ToUpper(opts.OptimizationGoal)
	}
	if opts.OptimizationStrategyType != "" {
		data["optimization_strategy_type"] = strings.ToUpper(opts.OptimizationStrategyType)
	}
	if opts.StartTime != "" {
		data["start_time"] = opts.StartTime
	}
	if opts.EndTime != "" {
		data["end_time"] = opts.EndTime
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
	if opts.SavedAudienceID != "" {
		data["saved_audience_id"] = opts.SavedAudienceID
	}
	if opts.ProductSetID != "" {
		data["product_set_id"] = opts.ProductSetID
	}
	if opts.ShoppingType != "" {
		data["shopping_type"] = strings.ToUpper(opts.ShoppingType)
	}
	if opts.TargetingInput != "" {
		value, err := parseJSONInput(opts.TargetingInput, "--targeting-json")
		if err != nil {
			return nil, err
		}
		data["targeting"] = value
	}
	if opts.ScheduleInput != "" {
		value, err := parseJSONInput(opts.ScheduleInput, "--schedule-json")
		if err != nil {
			return nil, err
		}
		data["schedule"] = value
	}
	if opts.ShoppingTargetingInput != "" {
		value, err := parseJSONInput(opts.ShoppingTargetingInput, "--shopping-targeting-json")
		if err != nil {
			return nil, err
		}
		data["shopping_targeting"] = value
	}

	if (opts.StartTime == "") != (opts.EndTime == "") {
		return nil, errors.New("use both --start-time and --end-time together")
	}
	if opts.RequireIdentity {
		if opts.CampaignID == "" || opts.Name == "" || opts.ConfiguredStatus == "" {
			return nil, errors.New("adgroup create requires --campaign, --name, and --configured-status")
		}
	}

	return data, nil
}

func parseJSONInput(raw, flagName string) (any, error) {
	input := raw
	if strings.HasPrefix(raw, "@") {
		data, err := os.ReadFile(strings.TrimPrefix(raw, "@"))
		if err != nil {
			return nil, fmt.Errorf("could not read %s file %q: %w", flagName, strings.TrimPrefix(raw, "@"), err)
		}
		input = string(data)
	}

	var value any
	if err := json.Unmarshal([]byte(input), &value); err != nil {
		return nil, fmt.Errorf("invalid %s value: %w", flagName, err)
	}
	return value, nil
}

func looksLikeID(input string) bool {
	if input == "" {
		return false
	}
	if strings.HasPrefix(input, "t2_") || strings.HasPrefix(input, "a2_") {
		return true
	}
	if strings.Count(input, "-") == 4 && len(input) == 36 {
		return true
	}
	for _, r := range input {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

const adGroupHelp = `Usage:
  rad adgroup list [--account-id <id-or-name>] [--all] [--page-size N] [--json]
  rad adgroup get <id-or-name> [--account-id <id-or-name>] [--json]
  rad adgroup create --campaign <id-or-name> --name <name> --configured-status <status> [flags]
  rad adgroup update <id-or-name> [flags]

Common write flags:
  --bid-strategy <bidless|maximize_volume|target_cpx>
  --bid-type <cpc|cpm|cpv|cpv6>
  --bid-value <major-currency>
  --goal-type <daily_spend|lifetime_spend>
  --goal-value <major-currency>
  --optimization-goal <goal>
  --optimization-strategy-type <type>
  --start-time <rfc3339>
  --end-time <rfc3339>
  --app-id <store-app-id>
  --conversion-pixel-id <id>
  --view-through-conversion-type <type>
  --saved-audience-id <id>
  --product-set-id <id>
  --shopping-type <dynamic|static>
  --targeting-json <json-or-@file>
  --schedule-json <json-or-@file>
  --shopping-targeting-json <json-or-@file>
  --dry-run

Examples:
  rad adgroup create --campaign "Spring Launch" --name "US Retargeting" --configured-status PAUSED --bid-type CPC --bid-value 1.25 --dry-run
  rad adgroup update "US Retargeting" --configured-status ACTIVE
  rad adgroup update "US Retargeting" --targeting-json @targeting.json --dry-run`
