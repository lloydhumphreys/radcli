# Profiles

`radcli` includes profile lookup commands so you can find the profile IDs needed
for post/creative management.

## Supported Commands

```bash
./bin/rad profile list
./bin/rad profile business-list
```

## Account-Level Listing

List profiles for the currently selected ad account:

```bash
./bin/rad profile list
```

## Business-Level Listing

List profiles for the currently selected business:

```bash
./bin/rad profile business-list
```

## Output

- table output is the default
- use `--json` for the raw API response
- list commands support `--all` and `--page-size`

## Notes

- `profile list` uses the selected ad account by default
- `profile business-list` uses the selected business by default
- profile IDs and names can then be used with `rad post`
