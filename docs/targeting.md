# Targeting

`radcli` now includes lookup-oriented targeting helpers so you can assemble ad
group targeting payloads without digging through the website.

## Supported Commands

```bash
rad targeting communities search --query <text>
rad targeting communities list [--name <community> ...]
rad targeting communities suggest [--name <community> ...] [--website-url <url>]
rad targeting interests list
rad targeting geolocations search [--postal-code <code>] [--city <name>] [--country <code>]
rad targeting devices list
rad targeting carriers list
rad targeting keywords suggest --keyword <term> [--keyword <term> ...]
rad targeting keywords validate --keyword <term> [--keyword <term> ...]
```

## Communities

Search by free text:

```bash
rad targeting communities search --query gaming
```

Look up exact communities:

```bash
rad targeting communities list --name games --name pcmasterrace
```

Get suggestions from a seed community or website:

```bash
rad targeting communities suggest --name games
rad targeting communities suggest --website-url https://example.com
```

## Interests

List the currently available interest taxonomy:

```bash
rad targeting interests list
```

## Geolocations

Search by postal code, city, or country:

```bash
rad targeting geolocations search --postal-code 94107
rad targeting geolocations search --city Copenhagen --country DK
rad targeting geolocations search --country US
```

## Devices And Carriers

```bash
rad targeting devices list
rad targeting carriers list
```

Both commands support `--all`, `--page-size`, and `--json`.

## Keywords

Get keyword suggestions from one or more seed keywords:

```bash
rad targeting keywords suggest --keyword fishing --keyword reels
```

Validate keywords for brand safety:

```bash
rad targeting keywords validate --keyword fishing --keyword reels
```

## Output

- table output is the default
- use `--json` for the raw API response
- paginated commands support `--all` and `--page-size`

## Notes

- this is a lookup-first surface, not a full targeting mutation DSL
- the main goal is to help build `--targeting-json` payloads for `rad adgroup create` and `rad adgroup update`
