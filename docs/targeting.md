# Targeting

`radcli` now includes lookup-oriented targeting helpers so you can assemble ad
group targeting payloads without digging through the website.

## Supported Commands

```bash
./bin/rad targeting communities search --query <text>
./bin/rad targeting communities list [--name <community> ...]
./bin/rad targeting communities suggest [--name <community> ...] [--website-url <url>]
./bin/rad targeting interests list
./bin/rad targeting geolocations search [--postal-code <code>] [--city <name>] [--country <code>]
./bin/rad targeting devices list
./bin/rad targeting carriers list
./bin/rad targeting keywords suggest --keyword <term> [--keyword <term> ...]
./bin/rad targeting keywords validate --keyword <term> [--keyword <term> ...]
```

## Communities

Search by free text:

```bash
./bin/rad targeting communities search --query gaming
```

Look up exact communities:

```bash
./bin/rad targeting communities list --name games --name pcmasterrace
```

Get suggestions from a seed community or website:

```bash
./bin/rad targeting communities suggest --name games
./bin/rad targeting communities suggest --website-url https://example.com
```

## Interests

List the currently available interest taxonomy:

```bash
./bin/rad targeting interests list
```

## Geolocations

Search by postal code, city, or country:

```bash
./bin/rad targeting geolocations search --postal-code 94107
./bin/rad targeting geolocations search --city Copenhagen --country DK
./bin/rad targeting geolocations search --country US
```

## Devices And Carriers

```bash
./bin/rad targeting devices list
./bin/rad targeting carriers list
```

Both commands support `--all`, `--page-size`, and `--json`.

## Keywords

Get keyword suggestions from one or more seed keywords:

```bash
./bin/rad targeting keywords suggest --keyword fishing --keyword reels
```

Validate keywords for brand safety:

```bash
./bin/rad targeting keywords validate --keyword fishing --keyword reels
```

## Output

- table output is the default
- use `--json` for the raw API response
- paginated commands support `--all` and `--page-size`

## Notes

- this is a lookup-first surface, not a full targeting mutation DSL
- the main goal is to help build `--targeting-json` payloads for `rad adgroup create` and `rad adgroup update`
