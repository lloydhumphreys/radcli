package cli

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"radcli/internal/config"
	"radcli/internal/output"
	"radcli/internal/reddit"
)

type App struct {
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	store       *config.Store
	api         *reddit.Client
	interactive bool
}

func New(stdin io.Reader, stdout, stderr io.Writer) (*App, error) {
	store, err := config.Load()
	if err != nil {
		return nil, err
	}

	interactive := false
	if file, ok := stdin.(*os.File); ok {
		if stat, err := file.Stat(); err == nil {
			interactive = (stat.Mode() & os.ModeCharDevice) != 0
		}
	}

	return &App{
		stdin:       stdin,
		stdout:      stdout,
		stderr:      stderr,
		store:       store,
		api:         reddit.New(store),
		interactive: interactive,
	}, nil
}

func (a *App) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		_, err := fmt.Fprintln(a.stdout, rootHelp)
		return err
	}

	switch args[0] {
	case "help", "--help", "-h":
		_, err := fmt.Fprintln(a.stdout, rootHelp)
		return err
	case "auth":
		return a.runAuth(ctx, args[1:])
	case "config":
		return a.runConfig(args[1:])
	case "business":
		return a.runBusiness(ctx, args[1:])
	case "account":
		return a.runAccount(ctx, args[1:])
	case "funding":
		return a.runFundingCommand(ctx, args[1:])
	case "campaign":
		return a.runCampaignCommand(ctx, args[1:])
	case "adgroup":
		return a.runAdGroupCommand(ctx, args[1:])
	case "ad":
		return a.runAssetCommand(ctx, assetDefinition{
			Command:            "ad",
			Label:              "ad",
			CollectionEndpoint: "ads",
			ItemEndpoint:       "ads",
			ListColumns:        []string{"id", "name", "campaign_id", "ad_group_id", "configured_status", "effective_status"},
		}, args[1:])
	case "targeting":
		return a.runTargetingCommand(ctx, args[1:])
	case "report":
		return a.runReportCommand(ctx, args[1:])
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], rootHelp)
	}
}

func (a *App) runConfig(args []string) error {
	if len(args) == 0 || args[0] == "show" {
		return output.PrintJSON(a.stdout, a.store.SanitizedMap())
	}
	return fmt.Errorf("unknown config command %q", args[0])
}

func (a *App) runAuth(ctx context.Context, args []string) error {
	if len(args) == 0 {
		_, err := fmt.Fprintln(a.stdout, authHelp)
		return err
	}

	switch args[0] {
	case "help":
		_, err := fmt.Fprintln(a.stdout, authHelp)
		return err
	case "setup":
		fs := newFlagSet("auth setup")
		clientID := fs.String("client-id", "", "")
		clientSecret := fs.String("client-secret", "", "")
		redirectURI := fs.String("redirect-uri", "", "")
		userAgent := fs.String("user-agent", "macos:com.lloyd.radcli:v0.1.0 (by /u/unknown)", "")
		var scopes stringList
		fs.Var(&scopes, "scope", "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}
		if *clientID == "" || *clientSecret == "" || *redirectURI == "" {
			return errors.New("auth setup requires --client-id, --client-secret, and --redirect-uri")
		}
		if len(scopes) == 0 {
			scopes = []string{"adsread", "adsedit", "adsconversions", "history", "read"}
		}
		a.store.Config.App = &config.AppCredentials{
			ClientID:     *clientID,
			ClientSecret: *clientSecret,
			RedirectURI:  *redirectURI,
			Scopes:       normalizeScopes(scopes),
			UserAgent:    *userAgent,
		}
		if err := a.store.Save(); err != nil {
			return err
		}
		_, err := fmt.Fprintf(a.stdout, "Saved app credentials to %s.\n", a.store.Path)
		return err
	case "login":
		fs := newFlagSet("auth login")
		openBrowser := fs.Bool("open", false, "")
		noWait := fs.Bool("no-wait", false, "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}
		rawURL, expectedState, err := a.api.BuildAuthorizationURL()
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(a.stdout, "Open this URL in your browser and approve the app:\n\n%s\n", rawURL); err != nil {
			return err
		}
		if *openBrowser {
			if err := reddit.OpenBrowser(rawURL); err != nil {
				return err
			}
		}
		if *noWait || !a.interactive {
			_, err := fmt.Fprintln(a.stdout, "\nAfter the redirect, copy the `code` query parameter and run:\n  rad auth complete --code <code>")
			return err
		}
		return a.finishInteractiveLogin(ctx, expectedState)
	case "complete":
		fs := newFlagSet("auth complete")
		code := fs.String("code", "", "")
		state := fs.String("state", "", "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}
		if *code == "" {
			return errors.New("auth complete requires --code")
		}
		if *state != "" && a.store.Config.Auth.PendingState != "" && *state != a.store.Config.Auth.PendingState {
			return errors.New("provided --state does not match the pending login state")
		}
		if _, err := a.api.ExchangeAuthorizationCode(ctx, *code); err != nil {
			return err
		}
		_, err := fmt.Fprintln(a.stdout, "Authentication complete. Run `rad auth whoami` or `rad business list` next.")
		return err
	case "whoami":
		fs := newFlagSet("auth whoami")
		jsonOut := fs.Bool("json", false, "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}
		payload, err := a.api.RequestJSON(ctx, "GET", "/me", nil, nil)
		if err != nil {
			return err
		}
		if *jsonOut {
			return output.PrintJSON(a.stdout, payload)
		}
		return output.PrintJSON(a.stdout, dataOrSelf(payload))
	case "logout":
		a.store.Config.Auth = config.AuthState{}
		if err := a.store.Save(); err != nil {
			return err
		}
		_, err := fmt.Fprintln(a.stdout, "Cleared saved auth session.")
		return err
	default:
		return fmt.Errorf("unknown auth command %q\n\n%s", args[0], authHelp)
	}
}

func (a *App) runBusiness(ctx context.Context, args []string) error {
	if len(args) == 0 {
		_, err := fmt.Fprintln(a.stdout, businessHelp)
		return err
	}
	switch args[0] {
	case "help":
		_, err := fmt.Fprintln(a.stdout, businessHelp)
		return err
	case "list":
		fs := newFlagSet("business list")
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
		payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/me/businesses", query, nil, *all)
		if err != nil {
			return err
		}
		if *jsonOut {
			return output.PrintJSON(a.stdout, payload)
		}
		return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"id", "name", "role", "configured_status", "effective_status"})
	case "use":
		if len(args) < 2 {
			return errors.New("usage: rad business use <business-id-or-name>")
		}
		id, name, err := a.resolveBusinessSelection(ctx, args[1])
		if err != nil {
			return err
		}
		a.store.Config.Defaults.BusinessID = id
		if err := a.store.Save(); err != nil {
			return err
		}
		_, err = fmt.Fprintf(a.stdout, "Default business set to %s (%s).\n", name, id)
		return err
	default:
		return fmt.Errorf("unknown business command %q\n\n%s", args[0], businessHelp)
	}
}

func (a *App) runAccount(ctx context.Context, args []string) error {
	if len(args) == 0 {
		_, err := fmt.Fprintln(a.stdout, accountHelp)
		return err
	}
	switch args[0] {
	case "help":
		_, err := fmt.Fprintln(a.stdout, accountHelp)
		return err
	case "list":
		fs := newFlagSet("account list")
		businessID := fs.String("business-id", "", "")
		all := fs.Bool("all", false, "")
		pageSize := fs.Int("page-size", 0, "")
		jsonOut := fs.Bool("json", false, "")
		if err := parseFlags(fs, args[1:]); err != nil {
			return err
		}
		if *businessID == "" {
			*businessID = a.store.Config.Defaults.BusinessID
		}
		if *businessID == "" {
			return errors.New("no business selected. use `rad business use <business-id-or-name>` or pass --business-id")
		}
		resolvedBusinessID, _, err := a.resolveBusinessSelection(ctx, *businessID)
		if err != nil {
			return err
		}
		*businessID = resolvedBusinessID
		query := url.Values{}
		if *pageSize > 0 {
			query.Set("page.size", fmt.Sprintf("%d", *pageSize))
		}
		payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/businesses/"+*businessID+"/ad_accounts", query, nil, *all)
		if err != nil {
			return err
		}
		if *jsonOut {
			return output.PrintJSON(a.stdout, payload)
		}
		return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"id", "name", "currency", "configured_status", "effective_status"})
	case "use":
		if len(args) < 2 {
			return errors.New("usage: rad account use <ad-account-id-or-name>")
		}
		businessID := a.store.Config.Defaults.BusinessID
		if businessID == "" {
			return errors.New("no business selected. use `rad business use <business-id-or-name>` first")
		}
		businessID, _, err := a.resolveBusinessSelection(ctx, businessID)
		if err != nil {
			return err
		}
		id, name, err := a.resolveAccountSelection(ctx, businessID, args[1])
		if err != nil {
			return err
		}
		a.store.Config.Defaults.AdAccountID = id
		if err := a.store.Save(); err != nil {
			return err
		}
		_, err = fmt.Fprintf(a.stdout, "Default account set to %s (%s).\n", name, id)
		return err
	default:
		return fmt.Errorf("unknown account command %q\n\n%s", args[0], accountHelp)
	}
}

func rowsFromPayload(payload map[string]any) []map[string]string {
	data, ok := payload["data"]
	if !ok {
		return nil
	}

	if list, ok := data.([]any); ok {
		rows := make([]map[string]string, 0, len(list))
		for _, item := range list {
			if row, ok := item.(map[string]any); ok {
				rows = append(rows, stringifyMap(row))
			}
		}
		return rows
	}

	if object, ok := data.(map[string]any); ok {
		if metrics, ok := object["metrics"].([]any); ok {
			rows := make([]map[string]string, 0, len(metrics))
			for _, item := range metrics {
				if row, ok := item.(map[string]any); ok {
					rows = append(rows, stringifyMap(row))
				}
			}
			return rows
		}
		return []map[string]string{stringifyMap(object)}
	}

	return nil
}

func stringifyMap(value map[string]any) map[string]string {
	out := make(map[string]string, len(value))
	for key, item := range value {
		out[key] = stringify(item)
	}
	return out
}

func stringify(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return fmt.Sprintf("%t", v)
	case []any:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			parts = append(parts, stringify(item))
		}
		return strings.Join(parts, ",")
	case map[string]any:
		parts := make([]string, 0, len(v))
		for key, item := range v {
			parts = append(parts, key+"="+stringify(item))
		}
		return strings.Join(parts, ",")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func dataOrSelf(payload map[string]any) any {
	if data, ok := payload["data"]; ok {
		return data
	}
	return payload
}

func normalizeScopes(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		for _, piece := range strings.Split(value, ",") {
			piece = strings.TrimSpace(piece)
			if piece == "" {
				continue
			}
			if _, ok := seen[piece]; ok {
				continue
			}
			seen[piece] = struct{}{}
			out = append(out, piece)
		}
	}
	return out
}

func (a *App) resolveBusinessSelection(ctx context.Context, input string) (string, string, error) {
	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/me/businesses", nil, nil, true)
	if err != nil {
		return "", "", err
	}

	rows := rowsFromPayload(payload)
	if len(rows) == 0 {
		return "", "", errors.New("no businesses returned for the current user")
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
		return "", "", fmt.Errorf("business name %q matched multiple businesses: %s. use the business ID instead", input, strings.Join(ids, ", "))
	}

	return "", "", fmt.Errorf("could not find business %q. run `rad business list` and use the UUID from the `id` column", input)
}

func (a *App) resolveAccountSelection(ctx context.Context, businessID, input string) (string, string, error) {
	query := url.Values{}
	query.Set("page.size", "1000")

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/businesses/"+businessID+"/ad_accounts", query, nil, true)
	if err != nil {
		return "", "", err
	}

	rows := rowsFromPayload(payload)
	if len(rows) == 0 {
		return "", "", errors.New("no ad accounts returned for the selected business")
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
		return "", "", fmt.Errorf("account name %q matched multiple ad accounts: %s. use the account ID instead", input, strings.Join(ids, ", "))
	}

	return "", "", fmt.Errorf("could not find ad account %q in the selected business. run `rad account list` and use the UUID from the `id` column", input)
}

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}

func parseFlags(fs *flag.FlagSet, args []string) error {
	return fs.Parse(interspersedArgs(fs, args))
}

func interspersedArgs(fs *flag.FlagSet, args []string) []string {
	flags := make([]string, 0, len(args))
	positionals := make([]string, 0, len(args))

	for index := 0; index < len(args); index++ {
		arg := args[index]
		switch {
		case arg == "--":
			positionals = append(positionals, args[index+1:]...)
			index = len(args)
		case arg == "-" || !strings.HasPrefix(arg, "-"):
			positionals = append(positionals, arg)
		default:
			flags = append(flags, arg)
			if flagConsumesValue(fs, arg) && index+1 < len(args) {
				index++
				flags = append(flags, args[index])
			}
		}
	}

	return append(flags, positionals...)
}

func flagConsumesValue(fs *flag.FlagSet, arg string) bool {
	name := strings.TrimLeft(arg, "-")
	if name == "" {
		return false
	}
	if cut := strings.IndexByte(name, '='); cut >= 0 {
		return false
	}

	defined := fs.Lookup(name)
	if defined == nil {
		return false
	}

	type boolFlag interface {
		IsBoolFlag() bool
	}

	if v, ok := defined.Value.(boolFlag); ok && v.IsBoolFlag() {
		return false
	}

	return true
}

func (a *App) finishInteractiveLogin(ctx context.Context, expectedState string) error {
	if _, err := fmt.Fprintln(a.stdout, "\nPaste the full callback URL or just the code, then press Enter."); err != nil {
		return err
	}
	if _, err := fmt.Fprint(a.stdout, "> "); err != nil {
		return err
	}

	reader := bufio.NewReader(a.stdin)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	line = strings.TrimSpace(line)
	if line == "" {
		_, err := fmt.Fprintln(a.stdout, "No code entered. You can finish later with `rad auth complete --code <code>`.")
		return err
	}

	code, state, err := parseAuthorizationInput(line)
	if err != nil {
		return err
	}
	if state != "" && expectedState != "" && state != expectedState {
		return errors.New("the pasted callback URL state does not match the pending login state")
	}

	if _, err := a.api.ExchangeAuthorizationCode(ctx, code); err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.stdout, "Authentication complete. Run `rad auth whoami` or `rad business list` next.")
	return err
}

func parseAuthorizationInput(input string) (string, string, error) {
	input = strings.TrimSpace(input)
	input = strings.TrimPrefix(input, "code=")

	if strings.Contains(input, "://") {
		parsed, err := url.Parse(input)
		if err != nil {
			return "", "", fmt.Errorf("could not parse pasted callback URL: %w", err)
		}
		code := strings.TrimSuffix(parsed.Query().Get("code"), "#_")
		if code == "" {
			return "", "", errors.New("no `code` query parameter found in pasted callback URL")
		}
		return code, parsed.Query().Get("state"), nil
	}

	code := strings.TrimSuffix(input, "#_")
	if code == "" {
		return "", "", errors.New("no authorization code provided")
	}
	return code, "", nil
}

type stringList []string

func (s *stringList) String() string {
	return strings.Join(*s, ",")
}

func (s *stringList) Set(value string) error {
	*s = append(*s, value)
	return nil
}

const rootHelp = `radcli: Reddit Ads from the terminal

Commands:
  auth      Configure and authenticate with Reddit Ads
  config    Show local configuration
  business  List and select businesses
  account   List and select ad accounts
  funding   Look up funding instruments
  campaign  List campaigns
  adgroup   List ad groups
  ad        List ads
  targeting Look up targeting entities
  report    Run reports

Examples:
  rad auth setup --client-id <id> --client-secret <secret> --redirect-uri https://example.com/oauth/callback
  rad auth login
  rad auth complete --code <code>
  rad business list
  rad business use <business-id-or-name>
  rad account list
  rad account use <ad-account-id-or-name>
  rad funding list
  rad campaign list
  rad campaign get <id-or-name>
  rad campaign create --name <name> --objective <objective> --configured-status PAUSED
  rad campaign update <id-or-name> --name <name>
  rad adgroup create --campaign <campaign> --name <name> --configured-status PAUSED --dry-run
  rad targeting communities search --query gaming
  rad report campaign-summary --since 7d
  rad report run --from 2026-03-01T00:00:00Z --to 2026-03-08T00:00:00Z --field IMPRESSIONS --field CLICKS`

const authHelp = `Usage:
  rad auth setup --client-id <id> --client-secret <secret> --redirect-uri <uri> [--scope adsread] [--user-agent ua]
  rad auth login [--open] [--no-wait]
  rad auth complete --code <code> [--state <state>]
  rad auth whoami [--json]
  rad auth logout`

const businessHelp = `Usage:
  rad business list [--all] [--page-size N] [--json]
  rad business use <business-id-or-name>`

const accountHelp = `Usage:
  rad account list [--business-id <id>] [--all] [--page-size N] [--json]
  rad account use <ad-account-id-or-name>`

const reportHelp = `Usage:
  rad report fields [--match TEXT] [--json]

  rad report run --from <iso8601> --to <iso8601> --field <FIELD> [--field <FIELD> ...]
                 [--breakdown <BREAKDOWN> ...] [--account-id <id-or-name>] [--time-zone-id <tz>]
                 [--all] [--page-size N] [--json|--csv] [--output FILE]

  rad report campaign-summary [--since 7d] [--daily] [--account-id <id-or-name>] [--campaign <id-or-name>] [--field FIELD] [--json|--csv] [--output FILE]
  rad report adgroup-summary [--since 7d] [--daily] [--account-id <id-or-name>] [--campaign <id-or-name>] [--adgroup <id-or-name>] [--field FIELD] [--json|--csv] [--output FILE]
  rad report ad-summary [--since 7d] [--daily] [--account-id <id-or-name>] [--campaign <id-or-name>] [--adgroup <id-or-name>] [--ad <id-or-name>] [--field FIELD] [--json|--csv] [--output FILE]`
