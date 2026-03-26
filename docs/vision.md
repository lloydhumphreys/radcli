# radcli Vision

`radcli` makes Reddit Ads scriptable, predictable, and fast.

The goal is to replace the highest-friction parts of the Reddit Ads website with a CLI that helps an operator:

- inspect campaigns, ad groups, ads, and reports quickly
- automate repetitive setup and reporting work
- avoid clicking around the web UI
- feed info to an Agent to make choices about how to improve

## Product Principles

- **Fast local UX.** Every command should be composable, scriptable, and quiet by default.
- **Safe by default.** Write commands support dry runs, and error messages tell you what to do next.
- **Operator-friendly.** Wrap the raw API in good defaults instead of mirroring every field blindly.
- **Stable context.** Business and ad account selection is first-class CLI state. Set it once, use it everywhere.

## Architecture

### CLI Layer

Human-friendly commands and flags. Parses input, applies defaults from local config, renders tables, JSON, and CSV.

### Config and Session Layer

Local machine state at `~/.config/radcli/config.json` for app credentials, tokens, default business, and default ad account.

### Reddit Transport Layer

Shared HTTP client handling OAuth bearer tokens, token refresh, user agent, pagination, and normalized API errors.

### Resource Commands

Higher-level command groups for real workflows:

- auth, business, account
- campaign, adgroup, ad
- post, profile
- report
- targeting, audience, pixel
- funding

## API Facts Driving the Design

- OAuth 2.0 authorization-code flow is required.
- Developer apps live in Business Manager.
- Core surfaces: businesses, ad accounts, campaigns, ad groups, ads, posts, reports, targeting, audiences, pixels, and product catalogs.
- Pagination follows Reddit-provided `next_url` values directly.
