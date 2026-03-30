# LinuxLab

Interactive CLI tool for mastering Linux commands, Vim, and ops skills in real terminal environments.

## Tech Stack

- **Language:** Go 1.22+
- **TUI:** Bubbletea / Lipgloss / Bubbles / Glamour (Charm ecosystem, v1 API)
- **Sandbox:** Docker SDK for Go (`github.com/docker/docker/client`)
- **YAML:** `gopkg.in/yaml.v3`
- **Data:** Local JSON (`~/.linuxlab/progress.json`)

## Project Structure

```
cmd/linuxlab/main.go          # Entry point
internal/
  challenge/                   # Challenge types and YAML loader
  verify/                      # Verification engine (file, command, composite)
  sandbox/                     # Docker sandbox + Vim runner
  progress/                    # Progress store + skill map calculation
  tui/                         # Bubbletea screens (menu, modules, challenges, detail, skillmap)
challenges/                    # YAML challenge directories (linux-basics/, vim/)
```

## Development

```bash
# Build
make build

# Run all tests
make test

# Run specific package tests
go test ./internal/challenge/ -v
go test ./internal/verify/ -v
go test ./internal/sandbox/ -v -timeout 60s
go test ./internal/progress/ -v
go test ./internal/tui/ -v

# Run the app
make run
```

## Workflow Rules

- **TDD:** Always write failing test first, then implement, then verify pass
- **Commit per feature:** Each completed feature point gets its own git commit
- **Tests must pass before commit:** Never commit with failing tests
- **Run `go test ./... -v` before committing** to ensure nothing is broken

## Code Conventions

- Go standard project layout: `cmd/` for entrypoints, `internal/` for private packages
- Use `t.TempDir()` for test isolation — no hardcoded temp paths
- Docker tests should call `t.Skip("docker not available")` when Docker is not running
- Chinese for user-facing strings (TUI labels, error messages, challenge descriptions)
- English for code identifiers, comments, and commit messages

## Key Design Decisions

- **Directory-per-challenge:** Each challenge is a directory with `challenge.yaml`, `init.sh`, `check.sh`, `solution.sh`
- **Layered verification:** file_content → file_exists → command_output → exit_code → script
- **Mixed sandbox:** Docker containers for Linux challenges, real Vim for Vim challenges
- **tea.ExecProcess:** Used to hand terminal control to `docker exec -it` and `vim`
- **Screen switching:** Root AppModel delegates to sub-models based on screenID enum

## Docs

- Design spec: `docs/superpowers/specs/2026-03-31-linuxlab-design.md`
- Implementation plan: `docs/superpowers/plans/2026-03-31-linuxlab-p0.md`
