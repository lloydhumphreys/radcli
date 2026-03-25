package cli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"radcli/internal/output"
)

func (a *App) runPostCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_, err := fmt.Fprintln(a.stdout, postHelp)
		return err
	}

	switch args[0] {
	case "list":
		return a.runPostListCommand(ctx, args[1:])
	case "get":
		return a.runPostGetCommand(ctx, args[1:])
	case "create":
		return a.runPostCreateCommand(ctx, args[1:])
	case "update":
		return a.runPostUpdateCommand(ctx, args[1:])
	default:
		return fmt.Errorf("unknown post command %q\n\n%s", args[0], postHelp)
	}
}

func (a *App) runPostListCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("post list")
	accountInput := fs.String("account-id", "", "")
	profileInput := fs.String("profile", "", "")
	postType := fs.String("type", "", "")
	source := fs.String("source", "", "")
	all := fs.Bool("all", false, "")
	pageSize := fs.Int("page-size", 0, "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if *profileInput == "" {
		return errors.New("post list requires --profile")
	}

	accountID, _, err := a.selectedAccountID(ctx, *accountInput)
	if err != nil {
		return err
	}
	profileID, _, err := a.resolveProfileInput(ctx, accountID, *profileInput)
	if err != nil {
		return err
	}

	query := url.Values{}
	if *postType != "" {
		query.Set("type", strings.ToUpper(*postType))
	}
	if *source != "" {
		query.Set("source", *source)
	}
	if *pageSize > 0 {
		query.Set("page.size", fmt.Sprintf("%d", *pageSize))
	}

	payload, err := a.api.RequestPaginatedJSON(ctx, "GET", "/profiles/"+profileID+"/posts", query, nil, *all)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"id", "type", "headline", "body", "allow_comments", "profile_id", "created_at", "post_url"})
}

func (a *App) runPostGetCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("post get")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return errors.New("usage: rad post get <post-id> [--json]")
	}

	payload, err := a.api.RequestJSON(ctx, "GET", "/posts/"+fs.Arg(0), nil, nil)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintJSON(a.stdout, dataOrSelf(payload))
}

func (a *App) runPostCreateCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("post create")
	accountInput := fs.String("account-id", "", "")
	profileInput := fs.String("profile", "", "")
	postType := fs.String("type", "", "")
	headline := fs.String("headline", "", "")
	body := fs.String("body", "", "")
	contentInput := fs.String("content-json", "", "")
	allowComments := fs.String("allow-comments", "", "")
	isRichText := fs.String("is-richtext", "", "")
	thumbnailURL := fs.String("thumbnail-url", "", "")
	dryRun := fs.Bool("dry-run", false, "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if *profileInput == "" || *postType == "" || *headline == "" {
		return errors.New("post create requires --profile, --type, and --headline")
	}

	var profileID string
	if *dryRun && looksLikeID(*profileInput) {
		profileID = *profileInput
	} else {
		accountID, _, err := a.selectedAccountID(ctx, *accountInput)
		if err != nil {
			return err
		}
		profileID, _, err = a.resolveProfileInput(ctx, accountID, *profileInput)
		if err != nil {
			return err
		}
	}

	data, err := postWriteData(postWriteOptions{
		Type:          *postType,
		Headline:      *headline,
		Body:          *body,
		ContentJSON:   *contentInput,
		AllowComments: *allowComments,
		IsRichText:    *isRichText,
		ThumbnailURL:  *thumbnailURL,
		CreateMode:    true,
	})
	if err != nil {
		return err
	}

	payload := map[string]any{"data": data}
	if *dryRun {
		return output.PrintJSON(a.stdout, payload)
	}

	response, err := a.api.RequestJSON(ctx, "POST", "/profiles/"+profileID+"/posts", nil, payload)
	if err != nil {
		return err
	}
	if !*jsonOut {
		if _, err := fmt.Fprintln(a.stdout, "Post created.\n"); err != nil {
			return err
		}
	}
	return output.PrintJSON(a.stdout, dataOrSelf(response))
}

func (a *App) runPostUpdateCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("post update")
	allowComments := fs.String("allow-comments", "", "")
	dryRun := fs.Bool("dry-run", false, "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return errors.New("usage: rad post update <post-id> --allow-comments <true|false> [--dry-run]")
	}
	if *allowComments == "" {
		return errors.New("post update currently supports only --allow-comments")
	}

	data, err := postWriteData(postWriteOptions{
		AllowComments: *allowComments,
		CreateMode:    false,
	})
	if err != nil {
		return err
	}
	payload := map[string]any{"data": data}
	if *dryRun {
		return output.PrintJSON(a.stdout, payload)
	}

	response, err := a.api.RequestJSON(ctx, "PATCH", "/posts/"+fs.Arg(0), nil, payload)
	if err != nil {
		return err
	}
	if !*jsonOut {
		if _, err := fmt.Fprintf(a.stdout, "Post updated: %s\n\n", fs.Arg(0)); err != nil {
			return err
		}
	}
	return output.PrintJSON(a.stdout, dataOrSelf(response))
}

type postWriteOptions struct {
	Type          string
	Headline      string
	Body          string
	ContentJSON   string
	AllowComments string
	IsRichText    string
	ThumbnailURL  string
	CreateMode    bool
}

func postWriteData(opts postWriteOptions) (map[string]any, error) {
	data := map[string]any{}

	if opts.CreateMode {
		if opts.Type == "" || opts.Headline == "" {
			return nil, errors.New("post create requires --type and --headline")
		}
		data["type"] = strings.ToUpper(opts.Type)
		data["headline"] = opts.Headline
		if opts.Body != "" {
			data["body"] = opts.Body
		}
		if opts.ContentJSON != "" {
			content, err := parseJSONInput(opts.ContentJSON, "--content-json")
			if err != nil {
				return nil, err
			}
			data["content"] = content
		}
		if opts.IsRichText != "" {
			value, err := parseBoolStringFlag(opts.IsRichText, "--is-richtext")
			if err != nil {
				return nil, err
			}
			data["is_richtext"] = value
		}
		if opts.ThumbnailURL != "" {
			data["thumbnail_url"] = opts.ThumbnailURL
		}
	}

	if opts.AllowComments != "" {
		value, err := parseBoolStringFlag(opts.AllowComments, "--allow-comments")
		if err != nil {
			return nil, err
		}
		data["allow_comments"] = value
	}

	return data, nil
}

func parseBoolStringFlag(raw, flagName string) (bool, error) {
	switch strings.ToLower(raw) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("invalid %s value %q: use true or false", flagName, raw)
	}
}

func (a *App) resolveProfileInput(ctx context.Context, accountID, input string) (string, string, error) {
	if looksLikeID(input) {
		return input, "", nil
	}
	return a.resolveProfileSelection(ctx, accountID, input)
}

const postHelp = `Usage:
  rad post list --profile <id-or-name> [--type <type>] [--source <source>] [--account-id <id-or-name>] [--all] [--page-size N] [--json]
  rad post get <post-id> [--json]
  rad post create --profile <id-or-name> --type <IMAGE|VIDEO|TEXT|CAROUSEL> --headline <headline> [--body <text>] [--content-json <json-or-@file>] [--allow-comments <true|false>] [--is-richtext <true|false>] [--thumbnail-url <url>] [--account-id <id-or-name>] [--dry-run]
  rad post update <post-id> --allow-comments <true|false> [--dry-run] [--json]

Examples:
  rad post list --profile "brand_profile"
  rad post create --profile "brand_profile" --type IMAGE --headline "Spring Launch" --content-json @content.json --dry-run
  rad post update t3_abcdef --allow-comments false`
