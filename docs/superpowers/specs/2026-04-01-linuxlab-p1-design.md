# LinuxLab P1 — 设计规格

> P1 阶段：扩展题库（Shell 脚本 + 运维实战）、高级验证类型、命令速查、薄弱推荐算法。

## 1. P0 现状

| 模块 | 状态 |
|------|------|
| TUI 主界面 | ✅ 完成 |
| Docker 沙盒引擎 | ✅ 完成 |
| Linux 基础题库 | ✅ 98 题 |
| Vim 题库 | ✅ 36 题 |
| 基础验证系统 (file_content/file_exists/command_output/exit_code) | ✅ 完成 |
| 进度记录 | ✅ 完成 |
| 能力图谱 | ✅ 完成 |

## 2. P1 新增功能

### 2.1 高级验证类型

新增两种验证类型，补充到 `internal/verify/` 中：

#### `permissions` 验证器

检查文件权限是否匹配期望值。

```yaml
verify:
  - type: permissions
    path: /tmp/deploy.sh
    expect: "750"
```

实现逻辑：
- 在容器/沙盒内执行 `stat -c '%a' <path>` 获取八进制权限
- 与 expect 字段比较
- 注意：本地测试时也使用 `stat`，macOS 的 `stat` 语法不同（`-f '%Lp'`），需要兼容处理

#### `script` 验证器

运行 check.sh 自定义脚本，检查退出码。

```yaml
verify:
  - type: script
    path: check.sh      # 相对于 challenge.Dir 的路径
```

实现逻辑：
- 读取 challenge.Dir 下的 check.sh
- 在沙盒内执行，exit 0 表示通过
- 捕获 stdout/stderr 作为 Message

### 2.2 命令速查功能

在 TUI 主菜单的「命令速查」选项中，提供一个交互式命令搜索界面。

#### 数据结构

```yaml
# references/commands.yaml
commands:
  - name: find
    brief: "搜索文件和目录"
    examples:
      - desc: "按名称查找"
        cmd: "find {{路径}} -name '{{文件名}}'"
      - desc: "按大小查找（大于 100MB）"
        cmd: "find {{路径}} -size +100M"
      - desc: "按修改时间查找（7天内）"
        cmd: "find {{路径}} -mtime -7"
    related_challenges: [find-name, find-size, find-mtime]
```

#### TUI 界面

- 搜索输入框（支持模糊匹配）
- 命令列表（匹配结果）
- 选中命令后展示详情（brief + examples + related challenges）
- 快捷键：`/` 聚焦搜索框，`Enter` 查看详情，`Esc` 返回

#### 实现模块

```
internal/reference/
├── types.go         # CommandRef, Example 结构体
├── loader.go        # 从 YAML 加载命令参考
├── loader_test.go
├── search.go        # 模糊搜索
└── search_test.go

internal/tui/
├── reference.go     # 命令速查 TUI 屏幕
└── reference_test.go

references/
└── commands.yaml    # 命令参考数据
```

### 2.3 薄弱推荐算法增强

当前 `RecommendWeakest` 只返回一个最弱子类别中的第一道未通过的题。P1 增强：

- 返回多个推荐（最多 5 个），来自不同的薄弱子类别
- 考虑 ScoreWithHints（使用提示的题目权重更低）
- 在 TUI 中新增「薄弱推荐」屏幕（当前主菜单的 "recommend" 选项），展示推荐列表

#### 实现

```go
// internal/progress/skillmap.go 新增
func RecommendMultiple(store *Store, challenges []*challenge.Challenge, limit int) []*challenge.Challenge
```

```
internal/tui/
├── recommend.go       # 薄弱推荐 TUI 屏幕
└── recommend_test.go
```

### 2.4 Shell 脚本模块题库

参考 https://www.runoob.com/linux/linux-shell.html 全部章节。

子类别：

| 子类别 | 内容 | 最少题数 |
|--------|------|----------|
| variables | 变量定义、赋值、引用、只读、删除、作用域 | 8 |
| strings | 字符串拼接、长度、截取、查找、替换 | 6 |
| arrays | 数组定义、遍历、长度、关联数组 | 5 |
| operators | 算术、关系、布尔、字符串、文件测试 | 6 |
| flow-control | if/elif/else、for、while、until、case、break/continue | 10 |
| functions | 函数定义、参数、返回值、局部变量 | 5 |
| io-redirect | 输入输出重定向、here document、/dev/null | 5 |
| parameters | $0-$9、$#、$*、$@、$?、$$、shift | 5 |

**总计：≥50 题**

所有题目放在 `challenges/shell-scripting/` 下。

### 2.5 运维实战模块题库

参考 https://www.runoob.com/linux/linux-command-manual.html 中的系统管理、网络通讯、磁盘维护等分类。

子类别：

| 子类别 | 内容 | 最少题数 |
|--------|------|----------|
| network-troubleshoot | tcpdump、ss、netstat、iptables、route、traceroute | 8 |
| log-analysis | journalctl、grep/awk/sed 分析日志、logrotate | 6 |
| system-monitoring | top/htop、vmstat、iostat、sar、dstat | 6 |
| systemd | systemctl start/stop/enable/disable/status、unit 文件编写、journal | 8 |
| disk-ops | LVM 基础、RAID 概念、fsck、dd、fdisk/parted | 5 |
| security | ssh 密钥管理、firewall/ufw、fail2ban 概念、文件完整性检查 | 5 |
| automation | cron 高级、systemd timer、at/batch、脚本自动化 | 5 |

**总计：≥43 题**

所有题目放在 `challenges/ops/` 下。

## 3. 文件结构变更

```
新增文件：
├── internal/verify/permissions.go        # permissions 验证器
├── internal/verify/permissions_test.go
├── internal/verify/script.go             # script 验证器
├── internal/verify/script_test.go
├── internal/reference/types.go           # 命令参考类型
├── internal/reference/loader.go          # YAML 加载
├── internal/reference/loader_test.go
├── internal/reference/search.go          # 模糊搜索
├── internal/reference/search_test.go
├── internal/tui/reference.go             # 命令速查屏幕
├── internal/tui/reference_test.go
├── internal/tui/recommend.go             # 薄弱推荐屏幕
├── internal/tui/recommend_test.go
├── references/commands.yaml              # 命令参考数据
├── challenges/shell-scripting/           # 50+ Shell 题目
└── challenges/ops/                       # 43+ 运维题目

修改文件：
├── internal/verify/verifier.go           # 注册 permissions、script
├── internal/progress/skillmap.go         # RecommendMultiple
├── internal/progress/skillmap_test.go    # 新测试
├── internal/tui/app.go                   # 路由 reference、recommend 屏幕
├── internal/tui/app_test.go              # 新测试
└── internal/tui/modules.go              # categoryLabels 添加 shell-scripting、ops
```

## 4. 分期实施顺序

P1 内部分 3 个小阶段：

### P1a — 高级验证 + 算法增强（代码）
1. permissions 验证器（TDD）
2. script 验证器（TDD）
3. RecommendMultiple 算法（TDD）
4. 薄弱推荐 TUI 屏幕（TDD）

### P1b — 命令速查（代码 + 数据）
5. reference 数据类型和加载器（TDD）
6. 模糊搜索（TDD）
7. 命令速查 TUI 屏幕（TDD）
8. commands.yaml 数据文件
9. app.go 路由集成

### P1c — 题库扩展（内容）
10. Shell 脚本题库（50+ 题）
11. 运维实战题库（43+ 题）
12. E2E 测试更新
