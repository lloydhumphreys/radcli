package cli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lloydhumphreys/radcli/internal/output"
)

type reportPreset struct {
	Name         string
	Breakdown    string
	DefaultSince string
	Fields       []string
	TableColumns []string
	Enrich       func(context.Context, *App, string, []map[string]string) ([]map[string]string, error)
}

type reportFilters struct {
	Campaign string
	AdGroup  string
	Ad       string
}

var reportPresets = map[string]reportPreset{
	"campaign-summary": {
		Name:         "campaign-summary",
		Breakdown:    "CAMPAIGN_ID",
		DefaultSince: "7d",
		Fields:       []string{"IMPRESSIONS", "CLICKS", "CTR", "SPEND", "CPC", "ECPM"},
		TableColumns: []string{"campaign_name", "campaign_id", "date", "impressions", "clicks", "ctr", "spend", "cpc", "ecpm"},
		Enrich:       enrichCampaignRows,
	},
	"adgroup-summary": {
		Name:         "adgroup-summary",
		Breakdown:    "AD_GROUP_ID",
		DefaultSince: "7d",
		Fields:       []string{"IMPRESSIONS", "CLICKS", "CTR", "SPEND", "CPC", "ECPM"},
		TableColumns: []string{"ad_group_name", "ad_group_id", "campaign_name", "campaign_id", "date", "impressions", "clicks", "ctr", "spend", "cpc", "ecpm"},
		Enrich:       enrichAdGroupRows,
	},
	"ad-summary": {
		Name:         "ad-summary",
		Breakdown:    "AD_ID",
		DefaultSince: "7d",
		Fields:       []string{"IMPRESSIONS", "CLICKS", "CTR", "SPEND", "CPC", "ECPM"},
		TableColumns: []string{"ad_name", "ad_id", "ad_group_name", "ad_group_id", "campaign_name", "campaign_id", "date", "impressions", "clicks", "ctr", "spend", "cpc", "ecpm"},
		Enrich:       enrichAdRows,
	},
}

func (a *App) runReportCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintln(a.stdout, reportHelp)
		return err
	}

	switch args[0] {
	case "fields":
		return a.runReportFields(ctx, args[1:])
	case "run":
		return a.runRawReport(ctx, args[1:])
	case "campaign-summary", "adgroup-summary", "ad-summary":
		return a.runPresetReport(ctx, reportPresets[args[0]], args[1:])
	default:
		return fmt.Errorf("unknown report command %q\n\n%s", args[0], reportHelp)
	}
}

func (a *App) runReportFields(ctx context.Context, args []string) error {
	fs := newFlagSet("report fields")
	match := fs.String("match", "", "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}

	fields, err := a.api.ReportFields(ctx)
	if err != nil {
		return err
	}

	if *match != "" {
		filtered := make([]string, 0, len(fields))
		matchUpper := strings.ToUpper(*match)
		for _, field := range fields {
			if strings.Contains(strings.ToUpper(field), matchUpper) {
				filtered = append(filtered, field)
			}
		}
		fields = filtered
	}

	if *jsonOut {
		return output.PrintJSON(a.stdout, fields)
	}
	for _, field := range fields {
		if _, err := fmt.Fprintln(a.stdout, field); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) runRawReport(ctx context.Context, args []string) error {
	fs := newFlagSet("report run")
	from := fs.String("from", "", "")
	to := fs.String("to", "", "")
	accountInput := fs.String("account-id", "", "")
	timeZoneID := fs.String("time-zone-id", "", "")
	all := fs.Bool("all", false, "")
	pageSize := fs.Int("page-size", 0, "")
	jsonOut := fs.Bool("json", false, "")
	csvOut := fs.Bool("csv", false, "")
	outputPath := fs.String("output", "", "")
	var fields stringList
	var breakdowns stringList
	fs.Var(&fields, "field", "")
	fs.Var(&breakdowns, "breakdown", "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if *jsonOut && *csvOut {
		return errors.New("use only one of --json or --csv")
	}

	if *from == "" || *to == "" {
		return errors.New("report run requires --from and --to")
	}
	if len(fields) == 0 {
		return errors.New("report run requires at least one --field")
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}

	data := map[string]any{
		"starts_at": *from,
		"ends_at":   *to,
		"fields":    uniqueStrings([]string(fields)),
	}
	if len(breakdowns) > 0 {
		data["breakdowns"] = uniqueStrings([]string(breakdowns))
	}
	if *timeZoneID != "" {
		data["time_zone_id"] = *timeZoneID
	}

	payload, err := a.requestReport(ctx, accountID, data, *pageSize, *all)
	if err != nil {
		return err
	}
	rows := formatReportRows(rowsFromPayload(payload))
	return a.writeReportOutput(*outputPath, *jsonOut, *csvOut, payload, rows, nil)
}

func (a *App) runPresetReport(ctx context.Context, preset reportPreset, args []string) error {
	fs := newFlagSet("report " + preset.Name)
	from := fs.String("from", "", "")
	to := fs.String("to", "", "")
	since := fs.String("since", preset.DefaultSince, "")
	accountInput := fs.String("account-id", "", "")
	timeZoneID := fs.String("time-zone-id", "", "")
	all := fs.Bool("all", false, "")
	pageSize := fs.Int("page-size", 0, "")
	jsonOut := fs.Bool("json", false, "")
	csvOut := fs.Bool("csv", false, "")
	outputPath := fs.String("output", "", "")
	daily := fs.Bool("daily", false, "")
	campaign := fs.String("campaign", "", "")
	campaignID := fs.String("campaign-id", "", "")
	adGroup := fs.String("adgroup", "", "")
	adGroupID := fs.String("adgroup-id", "", "")
	ad := fs.String("ad", "", "")
	adID := fs.String("ad-id", "", "")
	var extraFields stringList
	fs.Var(&extraFields, "field", "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if *jsonOut && *csvOut {
		return errors.New("use only one of --json or --csv")
	}

	startsAt, endsAt, err := resolveReportWindow(*from, *to, *since)
	if err != nil {
		return err
	}
	filters, err := resolveReportFilters(
		preset,
		*campaign, *campaignID,
		*adGroup, *adGroupID,
		*ad, *adID,
	)
	if err != nil {
		return err
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}

	breakdowns := []string{preset.Breakdown}
	if *daily {
		breakdowns = append(breakdowns, "DATE")
	}

	data := map[string]any{
		"starts_at":  startsAt,
		"ends_at":    endsAt,
		"fields":     uniqueStrings(append(append([]string{}, preset.Fields...), []string(extraFields)...)),
		"breakdowns": breakdowns,
	}
	if *timeZoneID != "" {
		data["time_zone_id"] = *timeZoneID
	}

	payload, err := a.requestReport(ctx, accountID, data, *pageSize, *all)
	if err != nil {
		return err
	}
	rows := formatReportRows(rowsFromPayload(payload))
	if preset.Enrich != nil {
		rows, err = preset.Enrich(ctx, a, accountID, rows)
		if err != nil {
			return err
		}
	}
	rows = filterReportRows(rows, filters)
	return a.writeReportOutput(*outputPath, *jsonOut, *csvOut, payload, rows, preset.TableColumns)
}

func (a *App) requestReport(ctx context.Context, accountID string, data map[string]any, pageSize int, fetchAll bool) (map[string]any, error) {
	query := url.Values{}
	if pageSize > 0 {
		query.Set("page.size", fmt.Sprintf("%d", pageSize))
	}

	return a.api.RequestPaginatedJSON(
		ctx,
		"POST",
		"/ad_accounts/"+accountID+"/reports",
		query,
		map[string]any{"data": data},
		fetchAll,
	)
}

func resolveReportWindow(from, to, since string) (string, string, error) {
	if (from == "") != (to == "") {
		return "", "", errors.New("use both --from and --to together")
	}
	if from != "" && to != "" {
		return hourlyTimestampStringPair(from, to)
	}

	duration, err := parseSinceDuration(since)
	if err != nil {
		return "", "", err
	}

	end := time.Now().UTC().Truncate(time.Hour)
	start := end.Add(-duration)
	return start.Format(time.RFC3339), end.Format(time.RFC3339), nil
}

func hourlyTimestampStringPair(from, to string) (string, string, error) {
	start, err := hourlyTimestampString(from)
	if err != nil {
		return "", "", err
	}
	end, err := hourlyTimestampString(to)
	if err != nil {
		return "", "", err
	}
	return start, end, nil
}

func hourlyTimestampString(raw string) (string, error) {
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return "", fmt.Errorf("invalid timestamp %q: use RFC3339 like 2026-03-01T00:00:00Z", raw)
	}
	return parsed.UTC().Truncate(time.Hour).Format(time.RFC3339), nil
}

func parseSinceDuration(input string) (time.Duration, error) {
	if input == "" {
		return 0, errors.New("missing --since value")
	}

	switch {
	case strings.HasSuffix(input, "d"):
		value, err := strconv.Atoi(strings.TrimSuffix(input, "d"))
		if err != nil || value <= 0 {
			return 0, fmt.Errorf("invalid day duration %q", input)
		}
		return time.Duration(value) * 24 * time.Hour, nil
	case strings.HasSuffix(input, "w"):
		value, err := strconv.Atoi(strings.TrimSuffix(input, "w"))
		if err != nil || value <= 0 {
			return 0, fmt.Errorf("invalid week duration %q", input)
		}
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	default:
		duration, err := time.ParseDuration(input)
		if err != nil || duration <= 0 {
			return 0, fmt.Errorf("invalid --since value %q. use values like 7d, 30d, 2w, or 168h", input)
		}
		return duration, nil
	}
}

func uniqueStrings(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func (a *App) writeReportOutput(path string, jsonOut, csvOut bool, payload map[string]any, rows []map[string]string, preferred []string) error {
	writer := a.stdout
	if path != "" {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = file
	}

	switch {
	case jsonOut:
		return output.PrintJSON(writer, payload)
	case csvOut:
		return output.PrintCSV(writer, rows, preferred)
	default:
		return output.PrintTable(writer, rows, preferred)
	}
}

func enrichCampaignRows(ctx context.Context, app *App, accountID string, rows []map[string]string) ([]map[string]string, error) {
	campaigns, err := app.fetchAssetNamesByID(ctx, accountID, "campaigns")
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if campaignID := firstNonEmpty(row["campaign_id"], row["campaignid"]); campaignID != "" {
			row["campaign_id"] = campaignID
			if name, ok := campaigns[campaignID]; ok {
				row["campaign_name"] = name
			}
		}
	}
	return rows, nil
}

func enrichAdGroupRows(ctx context.Context, app *App, accountID string, rows []map[string]string) ([]map[string]string, error) {
	adGroups, err := app.fetchAssetRowsByID(ctx, accountID, "ad_groups")
	if err != nil {
		return nil, err
	}
	campaigns, err := app.fetchAssetNamesByID(ctx, accountID, "campaigns")
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		adGroupID := firstNonEmpty(row["ad_group_id"], row["adgroup_id"])
		if adGroupID == "" {
			continue
		}
		row["ad_group_id"] = adGroupID
		if adGroup, ok := adGroups[adGroupID]; ok {
			if name := adGroup["name"]; name != "" {
				row["ad_group_name"] = name
			}
			campaignID := adGroup["campaign_id"]
			if campaignID != "" {
				row["campaign_id"] = campaignID
				if campaignName, ok := campaigns[campaignID]; ok {
					row["campaign_name"] = campaignName
				}
			}
		}
	}
	return rows, nil
}

func enrichAdRows(ctx context.Context, app *App, accountID string, rows []map[string]string) ([]map[string]string, error) {
	ads, err := app.fetchAssetRowsByID(ctx, accountID, "ads")
	if err != nil {
		return nil, err
	}
	adGroups, err := app.fetchAssetRowsByID(ctx, accountID, "ad_groups")
	if err != nil {
		return nil, err
	}
	campaigns, err := app.fetchAssetNamesByID(ctx, accountID, "campaigns")
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		adID := row["ad_id"]
		if adID == "" {
			continue
		}
		if ad, ok := ads[adID]; ok {
			if name := ad["name"]; name != "" {
				row["ad_name"] = name
			}
			adGroupID := ad["ad_group_id"]
			if adGroupID != "" {
				row["ad_group_id"] = adGroupID
				if adGroup, ok := adGroups[adGroupID]; ok {
					if adGroupName := adGroup["name"]; adGroupName != "" {
						row["ad_group_name"] = adGroupName
					}
					campaignID := adGroup["campaign_id"]
					if campaignID != "" {
						row["campaign_id"] = campaignID
						if campaignName, ok := campaigns[campaignID]; ok {
							row["campaign_name"] = campaignName
						}
					}
				}
			}
		}
	}
	return rows, nil
}

func (a *App) fetchAssetNamesByID(ctx context.Context, accountID, endpoint string) (map[string]string, error) {
	rows, err := a.fetchAssetRowsByID(ctx, accountID, endpoint)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for id, row := range rows {
		out[id] = row["name"]
	}
	return out, nil
}

func (a *App) fetchAssetRowsByID(ctx context.Context, accountID, endpoint string) (map[string]map[string]string, error) {
	query := url.Values{}
	query.Set("page.size", "1000")

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/ad_accounts/"+accountID+"/"+endpoint, query, nil, true)
	if err != nil {
		return nil, err
	}

	rows := rowsFromPayload(payload)
	out := make(map[string]map[string]string, len(rows))
	for _, row := range rows {
		if id := row["id"]; id != "" {
			out[id] = row
		}
	}
	return out, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func resolveReportFilters(preset reportPreset, campaign, campaignID, adGroup, adGroupID, ad, adID string) (reportFilters, error) {
	filters := reportFilters{}

	var err error
	filters.Campaign, err = mergeFilterValue("--campaign", campaign, "--campaign-id", campaignID)
	if err != nil {
		return reportFilters{}, err
	}
	filters.AdGroup, err = mergeFilterValue("--adgroup", adGroup, "--adgroup-id", adGroupID)
	if err != nil {
		return reportFilters{}, err
	}
	filters.Ad, err = mergeFilterValue("--ad", ad, "--ad-id", adID)
	if err != nil {
		return reportFilters{}, err
	}

	switch preset.Name {
	case "campaign-summary":
		if filters.AdGroup != "" || filters.Ad != "" {
			return reportFilters{}, errors.New("campaign-summary supports only --campaign or --campaign-id")
		}
	case "adgroup-summary":
		if filters.Ad != "" {
			return reportFilters{}, errors.New("adgroup-summary supports campaign and ad group filters, but not --ad")
		}
	}

	return filters, nil
}

func mergeFilterValue(flagName, value, aliasName, aliasValue string) (string, error) {
	switch {
	case value == "":
		return aliasValue, nil
	case aliasValue == "":
		return value, nil
	case strings.EqualFold(value, aliasValue):
		return value, nil
	default:
		return "", fmt.Errorf("%s and %s refer to different values; use only one", flagName, aliasName)
	}
}

func formatReportRows(rows []map[string]string) []map[string]string {
	for _, row := range rows {
		for key, value := range row {
			row[key] = formatReportValue(key, value)
		}
	}
	return rows
}

func formatReportValue(key, value string) string {
	if value == "" {
		return value
	}

	switch strings.ToLower(key) {
	case "spend", "cpc", "ecpm":
		return formatMicroCurrency(value)
	default:
		return formatPlainNumber(value)
	}
}

func formatMicroCurrency(value string) string {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return value
	}

	return strconv.FormatFloat(parsed/1_000_000, 'f', 2, 64)
}

func formatPlainNumber(value string) string {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return value
	}

	if !strings.ContainsAny(value, "eE") {
		return value
	}

	return strconv.FormatFloat(parsed, 'f', -1, 64)
}

func filterReportRows(rows []map[string]string, filters reportFilters) []map[string]string {
	if filters == (reportFilters{}) {
		return rows
	}

	filtered := make([]map[string]string, 0, len(rows))
	for _, row := range rows {
		if !matchesReportFilter(row, filters.Campaign, "campaign_name", "campaign_id", "campaignid") {
			continue
		}
		if !matchesReportFilter(row, filters.AdGroup, "ad_group_name", "ad_group_id", "adgroup_id") {
			continue
		}
		if !matchesReportFilter(row, filters.Ad, "ad_name", "ad_id") {
			continue
		}
		filtered = append(filtered, row)
	}
	return filtered
}

func matchesReportFilter(row map[string]string, needle string, keys ...string) bool {
	if needle == "" {
		return true
	}

	for _, key := range keys {
		value := row[key]
		switch {
		case value == "":
			continue
		case value == needle:
			return true
		case strings.EqualFold(value, needle):
			return true
		}
	}

	return false
}
