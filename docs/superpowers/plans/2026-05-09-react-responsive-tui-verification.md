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

## 2026-08-17 接入收尾验证 (React TUI ↔ Go 引擎)

### 新增能力

- **渐进式提示**：详情页 `h` 逐条解锁提示，CTA 显示「已用 N/M（影响得分）」；解锁数通过 `--hints N` 透传 `challenge run`，与 Go TUI 的 `ScoreWithHints` 扣分语义对齐；结果页 `r` 重试保留已用提示数。
- **进度即时刷新**：挑战结束后乐观更新内存进度（`applyRunResultToData`），随后后台执行 `linuxlab data dump --json` 与 Go 权威数据对账（技能图谱、跨会话进度一并刷新）。
- **后端状态徽标**：启动时并行执行 `doctor --json`，Header 显示「Docker 就绪 / 本地模式」。
- **Compose 宿主语义修复**（PTY 冒烟测试发现）：
  - `NewComposeSandbox` 将 compose 文件解析为绝对路径，修复相对挑战目录下路径被 `cmd.Dir` 二次拼接导致全部 compose 挑战无法启动的 bug；
  - Compose 挑战的 init/check 脚本改为在宿主的题目目录内执行（容器内没有 docker CLI、也没有宿主 /tmp 文件）；
  - Compose 交互 shell 改为宿主 shell（工作目录 = 题目目录），用户可执行 `docker compose up -d` 等宿主命令。

### 本次命令与结果

- `go test ./... -timeout 300s` — 全部通过（新增：`TestHintsFromArgs`、`TestComposeSandbox_RelativeDirDoesNotDoublePath`、`TestRunHostShellUsesChallengeDir`、`TestVerifyInSandboxScriptRuleUsesExecFn` 等）。
- `npm run test:react` — 61 个测试通过；`npm run typecheck-react` — 通过。
- `LINUXLAB_USE_GO_DATA=1 LINUXLAB_BIN=./linuxlab npm run test:react` — 通过（真实 Go 数据边界）。
- `./linuxlab doctor --json` → `{"docker":true,"challenges":278,"references":78}`。
- `./linuxlab challenge run compose-basics --json --hints 2` → setup/handoff(mode=compose)/result 事件，`hintsUsed:2` 正确透传。
- **PTY 全链路冒烟**（`scripts/react_tui_smoke.py`）：真实终端下 菜单 → 模块 → 详情 → `h` 解锁提示 → Enter 移交 Compose 宿主 shell → `docker compose up -d; exit` → 结果页「挑战通过 · ✓ 检查 1: 脚本检测通过」，Header 进度 1/278 → 2/278 即时刷新。全部通过。

### 已知限制

- Compose 挑战的宿主 shell 依赖本机 docker CLI 可用；`docker compose up -d` 的首次镜像拉取耗时可能超过 shell 移交等待（冒烟脚本 30s 上限内完成）。
- React TUI 仍是预览层，默认入口保持 Go/Bubble Tea TUI。

## 2026-08-17 全量批量验证 (278/278)

`scripts/validate_challenges.go` 扩展为覆盖全部 5 个类别：

- **containers 类别（41 题）**：新增宿主验证模式（与 runner 的 LocalSandbox/ComposeSandbox 语义一致），init/solution/check 在宿主的题目目录执行，验证后自动清理新增容器。顺带修复：
  - `dockerfile-workdir` solution 生成的 app.sh 缺 shebang（`exec format error`）
  - `network-ls-create` 子网 172.28.0.0/16 与本机已有网络冲突 → 改用 192.168.100.0/24
  - compose 挑战 init/solution 以文件方式执行（修复 `sh -c` 下 `$0=/bin/sh` 导致 `cd "$(dirname "$0")"` 落到 /bin）
- **ops 类别（50 题）**：修复 4 个 init.sh 无条件 apt（dns-dig-query / ssh-key-management / systemctl-start-stop / iostat-disk-io 加 guard，traceroute-hop 合并两次 apt）；验证单题超时 120s → 300s。
- **验证稳定性**：新增预构建沙盒镜像 `linuxlab/sandbox:22.04`（`Dockerfile.sandbox` 预装全部 init 工具并保留 apt lists），apt 在验证容器内变为 no-op，消除镜像延迟导致的间歇性超时。
- **curl-download**：YAML verify 规则与 check.sh 不一致（check.sh 接受 page.html 或 result.txt）→ 改用 `type: script`。

最终结果：**278/278 通过，0 失败，0 跳过**（此前 containers 类别从未纳入批量验证）。
