# LinuxLab — 设计规格

> 一个命令行交互式学习工具，帮助用户在真实终端环境中掌握 Linux 命令、Vim 操作和运维技能。

## 1. 项目概述

**目标用户**：希望系统性提升 Linux / Vim / 运维能力的开发者或运维工程师（自用）。

**核心理念**：在真实终端中练习，而非模拟器。学到的是真正的肌肉记忆。

**技术栈**：Go + Charm 生态（Bubbletea / Lipgloss / Bubbles / Glamour）

**分发方式**：单个二进制文件 + 外部 YAML 题库目录

## 2. 架构

```
linuxlab (单个二进制)
│
├── TUI 主界面 (Bubbletea)
│   ├── 主菜单（开始练习 / 能力图谱 / 薄弱推荐 / 命令速查）
│   ├── 模块选择（显示题量和完成进度）
│   ├── 关卡列表（已通过 ✓ / 当前 ▸ / 未完成 ○，带难度星级）
│   └── 能力图谱（分组概览 + 详情展开）
│
├── 练习引擎
│   ├── Linux 命令挑战 → 进入 Docker 容器的真实 shell
│   ├── Vim 挑战 → 启动真实 Vim，预设文件和任务
│   └── 结果检测器 → 分层验证系统
│
├── 题库加载器 → 读取目录制 YAML 题目
│
└── 进度管理 → 本地 ~/.linuxlab/progress.json
```

## 3. 学习模块

共 5 个模块，分阶段交付：

| 阶段 | 模块 | 内容范围 |
|------|------|----------|
| P0（首期） | Linux 基础命令 | 文件操作、权限、管道、重定向、进程管理 |
| P0（首期） | Vim 操作 | 移动、编辑、搜索替换、多窗口、宏 |
| P1 | Shell 脚本 | 变量、循环、条件、函数、正则 |
| P1 | 运维实战 | 网络排查、日志分析、系统监控、systemd、磁盘管理 |
| P2 | 容器与部署 | Docker、docker-compose、基础 k8s 操作 |

## 4. 题库结构

采用「目录制」，每道题一个目录，包含标准化文件：

```
challenges/
└── linux-basics/
    └── find-large-files/
        ├── challenge.yaml    # 元数据：标题、难度、标签、描述、渐进提示
        ├── init.sh           # 环境初始化脚本（在沙盒中执行）
        ├── check.sh          # 验证脚本（exit 0 = 通过）
        ├── solution.sh       # 参考答案（同时用于自动化测试）
        └── files/            # 挑战所需的预设文件
```

### challenge.yaml 结构

```yaml
id: find-large-files
title: "查找大文件"
difficulty: 2                # 1-5
category: linux-basics
subcategory: file-operations
tags: [find, du, sort]

description: |
  /var/log 目录下有很多日志文件。
  找出其中最大的 3 个文件，并将结果输出到 /tmp/result.txt
  格式：每行一个文件路径，按大小降序排列。

hints:
  - level: 1
    text: "试试 du 命令查看文件大小"
  - level: 2
    text: "du -a 可以列出所有文件，配合 sort -n 排序"
  - level: 3
    text: "du -a /var/log | sort -rn | head -3 | awk '{print $2}' > /tmp/result.txt"

verify:
  - type: file_content
    path: /tmp/result.txt
    expect: |
      /var/log/auth.log
      /var/log/app.log
      /var/log/sys.log
```

### Vim 题目示例

```yaml
id: vim-batch-replace
title: "批量替换文本"
difficulty: 2
category: vim
subcategory: editing
tags: [substitute, regex]

description: |
  打开文件后，将所有的 "foo" 替换为 "bar"，然后保存退出。

setup_files:
  - path: /tmp/vim-challenge.txt
    content: |
      foo is great
      I love foo
      foo and foo together

verify:
  - type: file_content
    path: /tmp/vim-challenge.txt
    expect: |
      bar is great
      I love bar
      bar and bar together
```

### 验证类型

分层验证系统，支持以下检测类型：

| 类型 | 说明 | 适用场景 |
|------|------|----------|
| `file_content` | 检查文件内容是否匹配 | 文件操作、Vim 编辑 |
| `file_exists` | 检查文件/目录是否存在 | 创建目录、移动文件 |
| `command_output` | 执行检测命令，检查输出 | 管道、进程管理 |
| `exit_code` | 执行检测命令，检查返回码 | 脚本编写 |
| `permissions` | 检查文件权限 | 权限管理 |
| `script` | 运行 check.sh 自定义检测 | 复杂运维场景 |

每道题可组合多个验证条件，全部通过才算完成。

## 5. 沙盒机制

### 混合模式

- **Docker 容器模式（默认）**：Linux 命令和运维题目在 Docker 容器中执行，真实且安全
- **降级模式**：未安装 Docker 时，基础文件操作题使用 chroot 或临时目录隔离
- **Vim 题目**：不需要沙盒，在本地启动真实 Vim，操作临时文件

### Docker 沙盒细节

- 基础镜像：`ubuntu:22.04`，预装常用运维工具
- 每题启动一个独立容器，`init.sh` 在容器内执行环境初始化
- 用户通过 `docker exec -it` 进入容器 shell
- 输入 `exit` 退出后，程序在容器内执行 `check.sh` 验证
- 验证完成后销毁容器
- 资源限制：CPU、内存、磁盘配额，防止误操作影响宿主机

## 6. TUI 界面

### 主菜单

```
╭─ LinuxLab ─────────────────────────────╮
│                                         │
│   ▸ 开始练习                            │
│     能力图谱                            │
│     薄弱推荐                            │
│     命令速查                            │
│                                         │
│   ↑↓ 选择  Enter 确认  q 退出          │
╰─────────────────────────────────────────╯
```

### 核心交互流程

1. 主菜单 → 选择「开始练习」
2. 模块选择 → 选择学习模块（显示各模块题量和完成进度）
3. 关卡列表 → 选择具体题目（✓ 已通过 / ▸ 当前 / ○ 未完成，带难度星级）
4. 任务说明 → 展示题目描述、难度、标签
5. `Enter` 进入真实终端（Docker shell 或 Vim）
6. 完成任务后 `exit`（或 `:wq`），程序自动检测结果
7. 显示通过/失败，可查看参考答案
8. 继续下一题或返回列表

### 快捷键

- `Enter` — 确认/进入
- `↑↓` 或 `j/k` — 上下导航
- `h` — 查看提示（渐进式，每按一次解锁下一级）
- `s` — 查看参考答案（仅在完成或放弃后可用）
- `q` / `Esc` — 返回上一级
- `?` — 帮助

## 7. 能力图谱

### 展示方式：分组概览 + 详情展开

```
╭─ 能力图谱 ──────────────────────────────────╮
│                                              │
│  ▸ Linux 基础    ████████████████░░  68%     │
│    ├ 文件操作       8/12  ██████░░           │
│    ├ 权限管理       3/8   ███░░░░░           │
│    └ 管道重定向    10/10  ████████  ✓        │
│                                              │
│  ▸ Vim            ██████████░░░░░░  48%     │
│  ▾ 运维实战       ███░░░░░░░░░░░░░  15%     │
│    ├ 网络排查       0/15  ░░░░░░░░  ← 推荐   │
│    └ 系统服务       3/10  ███░░░░░           │
│                                              │
│  整体: 34/85 (40%)  今日: +5                 │
╰──────────────────────────────────────────────╯
```

### 数据模型

```json
{
  "skills": {
    "linux-basics.file-operations": {
      "total": 12,
      "passed": 8,
      "score": 0.67
    }
  },
  "challenges": {
    "find-large-files": {
      "status": "passed",
      "attempts": 2,
      "hints_used": 1,
      "last_attempt": "2026-03-31"
    }
  }
}
```

### 评分规则

- 基础分：通过 = 1.0，未通过 = 0.0
- 提示惩罚：每使用一级提示，该题得分 × 0.8
- 子类别得分 = 该类下所有题的平均分
- 模块得分 = 该模块下所有子类别的平均分
- 薄弱推荐：自动推荐得分最低的子类别中未通过的题目

## 8. 命令速查

内置命令速查功能，借鉴 tldr-pages 的 `{{placeholder}}` 格式：

```
╭─ 命令速查: find ─────────────────────────────╮
│                                               │
│  find - 搜索文件和目录                        │
│                                               │
│  按名称查找:                                  │
│    find {{路径}} -name "{{文件名}}"            │
│                                               │
│  按大小查找（大于 100MB）:                     │
│    find {{路径}} -size +100M                   │
│                                               │
│  按修改时间查找（7天内）:                      │
│    find {{路径}} -mtime -7                     │
│                                               │
│  相关题目: file-ops-03, file-ops-07            │
╰───────────────────────────────────────────────╯
```

速查内容与题目关联，方便学完即练。

## 9. 数据存储

所有数据存储在本地，无需网络：

```
~/.linuxlab/
├── progress.json    # 学习进度和能力数据
└── config.yaml      # 用户配置（Docker 镜像、Vim 路径等）
```

## 10. 分期交付计划

### P0 — 最小可用版本

- TUI 主界面（主菜单、模块选择、关卡列表、任务说明）
- Docker 沙盒引擎
- Linux 基础命令题库（15-20 题）
- Vim 操作题库（10-15 题）
- 基础验证系统（file_content、file_exists、command_output）
- 进度记录
- 能力图谱

### P1 — 扩展内容

- Shell 脚本模块题库
- 运维实战模块题库
- 高级验证类型（permissions、script）
- 命令速查功能
- 薄弱推荐算法

### P2 — 高级特性

- 容器与部署模块题库
- docker-compose 多服务场景沙盒
- 降级模式（无 Docker 时的基础练习支持）
