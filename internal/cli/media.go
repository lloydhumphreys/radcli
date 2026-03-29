package cli

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"

	"github.com/lloydhumphreys/radcli/internal/output"
)

func (a *App) runMediaCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintln(a.stdout, mediaHelp)
		return err
	}

	switch args[0] {
	case "upload":
		return a.runMediaUploadCommand(ctx, args[1:])
	default:
		return fmt.Errorf("unknown media command %q\n\n%s", args[0], mediaHelp)
	}
}

func (a *App) runMediaUploadCommand(ctx context.Context, args []string) error {
	fs := newFlagSet("media upload")
	filePath := fs.String("file", "", "")
	mimeType := fs.String("mime-type", "", "")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	if *filePath == "" {
		return errors.New("media upload requires --file")
	}

	info, err := os.Stat(*filePath)
	if err != nil {
		return fmt.Errorf("cannot access file %q: %w", *filePath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%q is a directory, not a file", *filePath)
	}

	detectedMIME := *mimeType
	if detectedMIME == "" {
		detectedMIME = mime.TypeByExtension(filepath.Ext(*filePath))
	}
	if detectedMIME == "" {
		return errors.New("could not detect MIME type from file extension; use --mime-type to specify it")
	}

	filename := filepath.Base(*filePath)

	if _, err := fmt.Fprintf(a.stdout, "Requesting upload lease for %s (%s)...\n", filename, detectedMIME); err != nil {
		return err
	}

	lease, err := a.api.RequestMediaUploadLease(ctx, filename, detectedMIME)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(a.stdout, "Uploading to Reddit...\n"); err != nil {
		return err
	}

	mediaURL, err := a.api.UploadMediaToS3(ctx, lease, *filePath)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(a.stdout, "Upload complete.\n\n"); err != nil {
		return err
	}

	result := map[string]any{
		"media_url": mediaURL,
		"asset_id":  lease.AssetID,
		"filename":  filename,
		"mime_type": detectedMIME,
	}

	if *jsonOut {
		return output.PrintJSON(a.stdout, result)
	}

	if _, err := fmt.Fprintf(a.stdout, "media_url: %s\n", mediaURL); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(a.stdout, "asset_id:  %s\n\n", lease.AssetID); err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.stdout, `Use this media_url in post creation:
  rad post create --profile <profile> --type IMAGE --headline <headline> \
    --content-json '[{"media_url":"` + mediaURL + `","destination_url":"https://example.com"}]'`)
	return err
}

const mediaHelp = `Usage:
  rad media upload --file <path> [--mime-type <type>] [--json]

Upload a local image or video file to Reddit and get back a hosted URL
for use with 'rad post create --content-json'.

NOTE: This command requires the 'submit' OAuth scope, which is not
available to ads-only applications. If you get a 403 error, use a
publicly accessible URL as the media_url instead:

  rad post create --profile <profile> --type IMAGE --headline <headline> \
    --content-json '[{"media_url":"https://example.com/image.jpg","destination_url":"https://example.com"}]'

Reddit will fetch and proxy the image at post creation time.

Examples:
  rad media upload --file hero.jpg
  rad media upload --file promo.mp4 --mime-type video/mp4
  rad media upload --file banner.png --json`
