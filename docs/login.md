# Login

This is the currently working manual login flow for `radcli`.

## Prerequisites

You need a Reddit Ads developer application created in Business Manager.

Use a real redirect URI and make sure the value saved in Reddit matches exactly
what you pass to `rad auth setup`.

Example:

```text
https://radcli.lloydhumphreys.com/callback
```

Important:

- `redirect_uri` must match exactly
- scheme must match: `https://` vs `http://`
- host must match
- path must match
- trailing slash behavior must match
- if Reddit says `invalid redirect_uri parameter`, the saved app redirect URI and
  the CLI redirect URI do not exactly match

## 1. Build the CLI

```bash
cd /Users/lloyd/Code/radcli
env GOCACHE=$PWD/.gocache go build -o bin/rad ./cmd/rad
```

## 2. Save app credentials

```bash
./bin/rad auth setup \
  --client-id YOUR_CLIENT_ID \
  --client-secret YOUR_CLIENT_SECRET \
  --redirect-uri https://radcli.lloydhumphreys.com/callback \
  --scope adsread \
  --scope adsedit \
  --scope adsconversions \
  --scope history \
  --scope read \
  --user-agent 'macos:com.lloyd.radcli:v0.1.0 (by /u/YOUR_REDDIT_USERNAME)'
```

## 3. Start login

```bash
./bin/rad auth login
```

This prints a Reddit authorization URL and then waits for you to paste either:

- the full callback URL
- just the `code` value

Open it in your browser and approve the app.

## 4. Copy the code from the redirect URL

After approval, Reddit redirects to your configured callback URL with query
parameters that look like this:

```text
https://radcli.lloydhumphreys.com/callback?state=abc123&code=Ve49meZ9oh3lZFY35CEizm3jTsfJeA#_
```

You can paste the full callback URL back into `radcli`, or copy only the `code`
value:

```text
Ve49meZ9oh3lZFY35CEizm3jTsfJeA
```

Important:

- if the redirect URL ends with `#_`, the CLI strips it automatically
- if you are copying manually, do not include `#_` in the code
- in other words, do not pass `Ve49meZ9oh3lZFY35CEizm3jTsfJeA#_`
- pass only `Ve49meZ9oh3lZFY35CEizm3jTsfJeA`

## 5. Complete login

If you pasted the code or callback URL directly into `rad auth login`, you are
already done.

If you want the older manual flow, this still works:

```bash
./bin/rad auth complete --code YOUR_CODE
```

## 6. Verify access

```bash
./bin/rad auth whoami
./bin/rad business list
```

## Troubleshooting

### `invalid redirect_uri parameter`

This means the redirect URI in the auth request does not exactly match the one
saved in the Reddit developer app.

Check both:

```bash
./bin/rad config show
```

and the Redirect URL in the Reddit app settings.

### `code` copied from redirect does not work

Make sure you removed the trailing `#_` if Reddit appended it to the callback
URL, or just paste the full callback URL into `rad auth login` and let the CLI
handle it.
