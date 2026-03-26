# Posts

`radcli` now includes a profile-scoped post surface for creative management.

## Supported Commands

```bash
rad post list --profile <id-or-name>
rad post get <post-id>
rad post create --profile <id-or-name> --type <type> --headline <headline>
rad post update <post-id> --allow-comments <true|false>
```

## Profiles First

Posts are created under Reddit profiles, so it usually starts with:

```bash
rad profile list
```

Then use the profile name or ID with `rad post`.

## List And Inspect Posts

```bash
rad post list --profile "brand_profile"
rad post list --profile "brand_profile" --type IMAGE
rad post get t3_abcdef
```

## Create Posts

The minimum required fields are:

- `--profile`
- `--type`
- `--headline`

Examples:

```bash
rad post create \
  --profile "brand_profile" \
  --type IMAGE \
  --headline "Spring Launch" \
  --content-json @content.json \
  --dry-run

rad post create \
  --profile "brand_profile" \
  --type TEXT \
  --headline "Spring Launch" \
  --body "This is the post body" \
  --allow-comments true \
  --dry-run
```

## Content JSON

For image, video, and carousel posts, pass the `content` array as inline JSON or
`@file`:

```bash
rad post create \
  --profile "brand_profile" \
  --type IMAGE \
  --headline "Spring Launch" \
  --content-json '[{"media_url":"https://example.com/image.png","destination_url":"https://example.com"}]' \
  --dry-run
```

Content items can include fields like:

- `media_url`
- `destination_url`
- `display_url`
- `caption`
- `call_to_action`

## Updates

Reddit’s current v3 patch schema for posts is narrow. In this first slice,
`radcli` exposes what the schema clearly supports:

```bash
rad post update t3_abcdef --allow-comments false
```

So this is not a full “edit every post field” command yet. It is currently a
comment-toggle update path.

## Output

- table output is the default for `list`
- use `--json` for raw API payloads
- `--dry-run` prints the request body without sending it

## Notes

- `post list` and `post create` use the selected ad account by default for profile resolution
- if you already know the profile ID, you can pass it directly
