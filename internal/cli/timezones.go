package cli

import (
	"context"
	"fmt"

	"github.com/lloydhumphreys/radcli/internal/output"
)

func (a *App) runTimezoneCommand(ctx context.Context, args []string) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		_, err := fmt.Fprintln(a.stdout, timezoneHelp)
		return err
	}
	if args[0] != "list" {
		return fmt.Errorf("unknown timezone command %q\n\n%s", args[0], timezoneHelp)
	}

	fs := newFlagSet("timezone list")
	jsonOut := fs.Bool("json", false, "")
	if err := parseFlags(fs, args[1:]); err != nil {
		return err
	}

	payload, err := a.api.RequestJSON(ctx, "GET", "/targeting/timezones", nil, nil)
	if err != nil {
		return err
	}
	if *jsonOut {
		return output.PrintJSON(a.stdout, payload)
	}
	return output.PrintTable(a.stdout, rowsFromPayload(payload), []string{"id", "name", "country_code"})
}

const timezoneHelp = `Usage:
  rad timezone list [--json]`
