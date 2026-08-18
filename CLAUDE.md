# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
make build              # Build binary → ./linuxlab
make test               # go test ./... -v
make run                # Build + run

# Single package test
go test ./internal/challenge/ -v
go test ./internal/sandbox/ -v -timeout 60s

# Batch-validate all challenges in Docker containers
go run scripts/validate_challenges.go              # all categories
go run scripts/validate_challenges.go linux-basics  # single category
```

## Architecture

**Flow:** `cmd/linuxlab/main.go` loads challenges + progress + references → creates `AppModel` → Bubbletea event loop.

**TUI (internal/tui/):** Root `AppModel` owns a `screenID` enum and delegates to sub-models (menu, modules, challenges, detail, skillmap, recommend, reference). Screen transitions happen via `Update()` returning new screen state. Terminal handoff to Docker/Vim uses `tea.ExecProcess`.

**Sandbox (internal/sandbox/):** `Sandbox` interface with three backends: `DockerSandbox` (per-challenge container via Docker SDK), `ComposeSandbox` (multi-service via docker-compose), `LocalSandbox` (degraded fallback). `NewSandbox()` auto-selects based on Docker availability and `compose_file` field. Verification runs INSIDE the container via `sb.Exec()`, not on the host.

**Challenge (internal/challenge/):** YAML loader reads `challenges/<category>/<id>/challenge.yaml`. Each challenge dir also has `init.sh`, `check.sh`, `solution.sh`. The `SetupFiles` field injects files into the container before the challenge starts (used by Vim challenges).

**Verify (internal/verify/):** Registry pattern — 6 verifier types registered by name. `RunAll()` iterates rules and dispatches. Most challenges use `type: script` pointing to `check.sh`.

**Progress (internal/progress/):** JSON store at `~/.linuxlab/progress.json`. `BuildSkillMap()` aggregates by subcategory. `ScoreWithHints()` penalizes hint usage.

## Key Conventions

- **Language:** Chinese for all user-facing strings (TUI, challenge descriptions, error messages). English for code identifiers and commit messages.
- **Challenge scripts:** `init.sh` must not run slow `apt-get` unconditionally — use `if ! command -v <tool>` guards. Scripts that use network tools (tcpdump, iptables) must have fallbacks for containers without NET_ADMIN/NET_RAW capabilities.
- **TUI borders:** `contentBox` in `styles.go` builds border lines manually (not via lipgloss border rendering) to avoid ANSI escape corruption when injecting titles.
- **Docker tests:** Call `t.Skip("docker not available")` when Docker is not running.
- **Test isolation:** Use `t.TempDir()`, never hardcoded temp paths.

## Challenge Validation

The validation script (`scripts/validate_challenges.go`) creates a fresh Docker container per challenge, runs init.sh → solution.sh → verify rules. The `containers` category (41 题) runs in host mode matching runner semantics (LocalSandbox/ComposeSandbox): init/solution/check execute on the host in the challenge dir, and containers started by the solution are removed afterwards. Exit code -1 from check.sh typically means timeout/killed (blocking command or slow apt-get).

## Docs

- Design spec: `docs/superpowers/specs/2026-03-31-linuxlab-design.md`
- Implementation plan: `docs/superpowers/plans/2026-03-31-linuxlab-p0.md`
