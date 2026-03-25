# Distribution

`radcli` now has a release/distribution scaffold for:

- GitHub Releases
- Homebrew tap publishing
- in-place binary updates with `rad self-update`

## Versioning

The binary exposes build metadata:

```bash
./bin/rad version
./bin/rad version --json
```

Release builds inject:

- version
- git commit
- build date
- GitHub repository coordinates

## Self Update

Installed binaries can update themselves from GitHub Releases:

```bash
rad self-update --check
rad self-update
rad self-update --version v0.1.0
rad self-update --repo owner/repo
```

Notes:

- `rad self-update` uses the embedded release repository from the build
- you can override that with `RADCLI_UPDATE_REPOSITORY=owner/repo`
- current implementation supports in-place replacement on Unix-like systems
- Windows currently falls back to manual download

## GoReleaser

Release automation is configured in:

```text
.goreleaser.yaml
```

It builds `rad` for:

- macOS
- Linux
- Windows

And publishes:

- versioned archives
- checksums
- GitHub Releases
- a Homebrew cask update, if tap settings are configured

## GitHub Actions

Workflows:

- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`

`ci.yml` runs tests and a build on pushes and pull requests.

`release.yml` runs on version tags like:

```bash
git tag v0.1.0
git push origin v0.1.0
```

## Required Release Settings

The release workflow expects:

### Default GitHub Release Publishing

No extra setup is required for publishing the GitHub Release itself when the
workflow runs in the main repository.

### Homebrew Tap Publishing

Create a separate tap repository, for example:

```text
OWNER/homebrew-tap
```

Then set this repository variable:

- `RADCLI_HOMEBREW_TAP_OWNER`
- `RADCLI_HOMEBREW_TAP_NAME`

And this repository secret:

- `HOMEBREW_TAP_GITHUB_TOKEN`

That token needs content write access to the tap repository.

## Homebrew Install And Upgrade

Once releases are flowing and the tap repo is configured, users can install with:

```bash
brew tap OWNER/TAP_REPO
brew install --cask radcli
```

And upgrade with:

```bash
brew upgrade --cask radcli
```

## Suggested First Distribution Rollout

1. Create the GitHub repo for `radcli`.
2. Add `origin`.
3. Push `main`.
4. Create the Homebrew tap repo.
5. Set the repository variables and secret.
6. Tag a release:

```bash
git tag v0.1.0
git push origin main --tags
```

7. Verify:

- GitHub Release artifacts exist
- checksums are attached
- the tap repo received the cask update
- `rad version` shows release metadata
- `rad self-update --check` sees the latest release
