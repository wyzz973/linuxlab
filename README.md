# LinuxLab

在真实终端环境中掌握 Linux 命令、Vim 操作、Shell 脚本和运维技能的交互式 CLI 学习工具。

```
╭─ LinuxLab ─────────────────────────────╮
│                                         │
│   ▸ 开始练习                            │
│     能力图谱                            │
│     薄弱推荐                            │
│     命令速查                            │
│                                         │
│   ↑↓/jk 选择  Enter 确认  q 退出       │
╰─────────────────────────────────────────╯
```

## 特性

- **278 道挑战题** — 覆盖 Linux 基础、Vim、Shell 脚本、运维实战、容器部署 5 大模块
- **真实终端练习** — 在 Docker 容器中执行真命令，在真实 Vim 中编辑，不是模拟器
- **能力图谱** — 按知识领域追踪掌握程度，自动推荐薄弱环节
- **渐进式提示** — 每题 3 级提示，从概念到具体到接近答案，鼓励独立思考
- **命令速查** — 内置 78 个常用命令的速查手册，支持模糊搜索
- **降级模式** — 没有 Docker 时自动切换到本地沙盒，基础题照常练习
- **单文件分发** — 4.9MB 单二进制，无运行时依赖

## 快速开始

```bash
# 构建
make build

# 运行
./linuxlab

# 或者一步到位
make run
```

### 前置要求

- Go 1.22+（仅构建时需要）
- Docker（推荐，用于沙盒隔离；没有也可以用降级模式）

## 题库总览

| 模块 | 题数 | 涵盖内容 |
|------|------|----------|
| **Linux 基础** | 98 | 文件操作、内容查看、权限管理、用户管理、磁盘管理、管道重定向、进程管理、网络基础、系统信息、包管理 |
| **Vim 操作** | 36 | 移动导航、编辑操作、搜索替换、多窗口、宏操作、实用技巧 |
| **Shell 脚本** | 53 | 变量、字符串、数组、运算符、流程控制、函数、IO 重定向、参数处理 |
| **运维实战** | 50 | 网络排查、日志分析、系统监控、systemd、磁盘运维、安全、自动化 |
| **容器部署** | 41 | Docker 基础、Dockerfile、Docker Compose、Docker 网络、数据卷 |

## 工作流程

```
选择模块 → 选择关卡 → 阅读任务说明 → 进入真实终端
→ 完成任务 → exit 退出 → 自动检测结果 → 记录进度
```

1. 启动 `linuxlab`，在 TUI 主界面选择学习模块
2. 浏览关卡列表，选择一道题
3. 阅读任务描述，按 `h` 查看渐进式提示
4. 按 `Enter` 进入真实终端（Docker 容器或 Vim）
5. 完成任务后 `exit`（或 `:wq`）退出
6. 程序自动检测结果，记录到能力图谱

## 项目结构

```
cmd/linuxlab/main.go          # 入口
internal/
  challenge/                   # 题目类型定义和 YAML 加载器
  verify/                      # 验证引擎（6 种类型）
  sandbox/                     # 沙盒（Docker / Compose / 本地降级）
  progress/                    # 进度存储 + 能力图谱计算
  reference/                   # 命令速查数据和搜索
  tui/                         # Bubbletea TUI 界面（8 个屏幕）
challenges/                    # YAML 题库（目录制）
  linux-basics/                #   98 题
  vim/                         #   36 题
  shell-scripting/             #   53 题
  ops/                         #   50 题
  containers/                  #   41 题
references/
  commands.yaml                # 78 个命令速查数据
```

## 题目格式

每道题是一个目录，包含标准化文件：

```
challenges/linux-basics/find-large-files/
├── challenge.yaml    # 元数据：标题、难度、标签、描述、提示、验证规则
├── init.sh           # 环境初始化脚本
├── check.sh          # 验证脚本（exit 0 = 通过）
└── solution.sh       # 参考答案
```

### challenge.yaml 示例

```yaml
id: find-large-files
title: "查找大文件"
difficulty: 2
category: linux-basics
subcategory: file-operations
tags: [find, du, sort]
description: |
  /var/log 目录下有多个日志文件。
  找出其中最大的 3 个文件，将路径写入 /tmp/result.txt。
hints:
  - level: 1
    text: "试试 du 命令查看文件大小"
  - level: 2
    text: "du -a 配合 sort -rn 按大小降序"
  - level: 3
    text: "du -a /var/log | sort -rn | head -3 | awk '{print $2}'"
verify:
  - type: file_content
    path: /tmp/result.txt
    expect: |
      /var/log/auth.log
      /var/log/app.log
      /var/log/sys.log
```

## 验证类型

| 类型 | 说明 |
|------|------|
| `file_content` | 检查文件内容是否匹配 |
| `file_exists` | 检查文件/目录是否存在 |
| `command_output` | 执行命令，检查输出 |
| `exit_code` | 执行命令，检查返回码 |
| `permissions` | 检查文件权限（八进制） |
| `script` | 运行 check.sh 自定义检测 |

## 沙盒模式

| 模式 | 触发条件 | 说明 |
|------|----------|------|
| **Docker** | Docker 可用 | 每题启动独立容器，真实隔离 |
| **Compose** | 题目有 `compose_file` | 多服务场景（web + db 等） |
| **本地** | Docker 不可用 | 临时目录中执行，基础题可用 |

## 开发

```bash
# 运行所有测试
make test

# 运行指定包测试
go test ./internal/verify/ -v
go test ./internal/sandbox/ -v -timeout 60s

# 构建
make build
```

### 贡献题目

1. 在对应模块目录下创建新目录
2. 编写 `challenge.yaml`、`init.sh`、`check.sh`、`solution.sh`
3. 确保 `solution.sh` 在 `init.sh` 之后能通过 `check.sh`
4. 运行 `go test ./... -v` 确认 E2E 测试通过

## 技术栈

- **语言:** Go 1.22+
- **TUI:** [Bubbletea](https://github.com/charmbracelet/bubbletea) / [Lipgloss](https://github.com/charmbracelet/lipgloss) / [Bubbles](https://github.com/charmbracelet/bubbles)
- **沙盒:** Docker SDK for Go
- **数据:** YAML 题库 + 本地 JSON 进度（`~/.linuxlab/progress.json`）

## License

MIT
