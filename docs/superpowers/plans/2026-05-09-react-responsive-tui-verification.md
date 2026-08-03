# React Responsive TUI Verification

Date: 2026-05-10

## Commands Run

- `npm run test:react` — passed, 13 files / 44 tests.
- `npm run typecheck-react` — passed.
- `env LINUXLAB_USE_GO_DATA=1 LINUXLAB_BIN=./linuxlab npm run test:react` — passed, verifies React can load through the Go data boundary.
- `go test ./internal/runner ./internal/cli ./internal/tui -v` — passed.
- `go test ./... -v` — passed; Docker was available, one local-fallback test skipped by existing convention.
- `go build -buildvcs=false -o /tmp/linuxlab-build-check ./cmd/linuxlab` — passed.
- `go build -buildvcs=false -o ./linuxlab ./cmd/linuxlab` — passed.
- `/tmp/linuxlab-build-check doctor --json` — returned parseable JSON with Docker available, 278 challenges, and 78 references.
- `/tmp/linuxlab-build-check` — default Go/Bubble Tea TUI launched with no args and exited with `q`.
- `./linuxlab data dump --json` — returned the React data contract with categories, progress, and references.
- `./linuxlab challenge run ls-basic --json` — launched Docker mode, completed `ls -a /home/lab > /tmp/result.txt`, and returned `{"type":"result","passed":true,...}`.

## Terminal Size Smoke Checks

- `COLUMNS=60 LINES=20 LINUXLAB_BIN=./linuxlab npm run tui:react` — compact layout rendered top tabs, no sidebar, footer visible.
- `COLUMNS=80 LINES=24 LINUXLAB_BIN=./linuxlab npm run tui:react` — standard layout rendered sidebar and main workspace.
- `COLUMNS=100 LINES=30 LINUXLAB_BIN=./linuxlab npm run tui:react` — standard layout rendered with more vertical list budget.
- `COLUMNS=140 LINES=40 LINUXLAB_BIN=./linuxlab npm run tui:react` — wide layout rendered sidebar, main workspace, and inspector.

## Notes

- The React TUI now supports both the TypeScript YAML loader and the Go `data dump --json` loader. The Go loader is opt-in through `LINUXLAB_USE_GO_DATA=1`.
- Added an E2E guard that `challenge.LoadAll()` loads every `challenge.yaml`; this caught and fixed two invalid YAML escape sequences that made Go skip two challenges.
- `make build` still uses default Go VCS stamping. For deterministic verification in this workspace, use `go build -buildvcs=false -o /tmp/linuxlab-build-check ./cmd/linuxlab`.
