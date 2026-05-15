# Releasing

Releases are driven by [release-kit](https://github.com/pike00/release-kit) locally,
with [GoReleaser](https://goreleaser.com) handling cross-platform binary builds in CI.

## Release flow

```
just release patch   # or minor | major
```

This runs `release-kit cut`, which:

1. Computes the next version (semver bump from the latest tag)
2. Drafts release notes via LLM (deepseek-v4-pro-cloud via local LiteLLM proxy)
3. Opens `$EDITOR` for review and editing
4. Updates `CHANGELOG.md` via git-cliff and commits it
5. Creates the GitHub release with the reviewed notes
6. Pushes the tag to `origin`

Pushing the tag triggers the `Release` GitHub Actions workflow, which runs
GoReleaser to:

- Build `coldkey` for Linux/macOS/Windows (amd64 + arm64)
- Attach `.tar.gz`/`.zip` archives and `checksums.txt` to the existing GitHub release
- Update the Homebrew cask in `pike00/homebrew-tap`
- Build and push the Docker image to `ghcr.io/pike00/coldkey`

GoReleaser is configured with `mode: keep-existing` so it never overwrites the
LLM-drafted release notes.

## Other recipes

| Recipe | What it does |
|---|---|
| `just changelog-preview` | Preview the unreleased CHANGELOG section (no writes) |
| `just changelog-backfill` | Regenerate `CHANGELOG.md` from full git history |
| `just notes-dry-run` | Draft LLM release notes for HEAD without cutting a release |
| `just release-yes patch` | Skip the confirm prompt (scripted use) |

## Tooling

- **release-kit** — `uv tool install --editable ~/projects/release-kit`
- **git-cliff** — invoked via `uvx git-cliff@latest` (no install required)
- **GoReleaser** — runs in CI; not needed locally

## Config files

| File | Purpose |
|---|---|
| `release-kit.toml` | release-kit config: project name, LiteLLM model, voice notes for LLM |
| `cliff.toml` | git-cliff config (symlink to `~/projects/_templates/cliff.toml`) |
| `release.just` | Shared just recipes (symlink to `~/projects/_templates/release.just`) |
| `.goreleaser.yml` | GoReleaser config: builds, archives, Homebrew cask, Docker |
