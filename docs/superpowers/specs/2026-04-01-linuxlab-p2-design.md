# LinuxLab P2 — 设计规格

> P2 阶段：容器与部署题库、docker-compose 多服务沙盒、降级模式（无 Docker 时的基础练习支持）。

## 1. P0 + P1 现状

| 模块 | 状态 |
|------|------|
| 题库 | ✅ 237 题（Linux 98 + Vim 36 + Shell 53 + Ops 50） |
| 验证系统 | ✅ 6 种类型（file_content/file_exists/command_output/exit_code/permissions/script） |
| 命令速查 | ✅ 78 个命令 |
| TUI | ✅ 8 个屏幕（menu/modules/challenges/detail/skillmap/recommend/reference/result） |
| 沙盒 | ✅ Docker 单容器 + Vim runner |

## 2. P2 新增功能

### 2.1 容器与部署题库

参考 https://www.runoob.com/docker/docker-tutorial.html 全部章节。

子类别：

| 子类别 | 内容 | 最少题数 |
|--------|------|----------|
| docker-basics | docker run/pull/images/ps/stop/rm/exec/logs | 10 |
| dockerfile | FROM/RUN/COPY/CMD/ENTRYPOINT/EXPOSE/ENV/WORKDIR/多阶段构建 | 8 |
| docker-compose | docker-compose up/down/ps/logs、YAML 配置、多服务编排 | 7 |
| docker-network | bridge/host/none 网络模式、自定义网络、容器间通信 | 5 |
| docker-volume | 数据卷挂载、bind mount、volume 管理 | 5 |

**总计：≥35 题**

所有题目放在 `challenges/containers/` 下。

### 2.2 docker-compose 多服务沙盒

当前沙盒只支持单容器。P2 扩展支持 docker-compose 多服务场景。

#### 设计

在 challenge.yaml 中新增 `compose_file` 字段：

```yaml
id: compose-web-db
category: containers
subcategory: docker-compose
compose_file: docker-compose.yaml  # 相对于 challenge.Dir
description: |
  使用 docker-compose 启动一个 web + db 的服务...
```

challenge 目录下放置 `docker-compose.yaml`：

```yaml
version: "3"
services:
  web:
    image: nginx:alpine
    ports: ["8080:80"]
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_PASSWORD: test
```

#### 沙盒引擎扩展

```go
// internal/sandbox/compose.go
type ComposeSandbox struct {
    dir string  // docker-compose.yaml 所在目录
}

func NewComposeSandbox(ctx context.Context, dir string) (*ComposeSandbox, error)
func (s *ComposeSandbox) Up(ctx context.Context) error        // docker-compose up -d
func (s *ComposeSandbox) Down(ctx context.Context) error      // docker-compose down
func (s *ComposeSandbox) Exec(ctx context.Context, service, command string) (string, int, error)
func (s *ComposeSandbox) InteractiveShellArgs(service string) []string
```

### 2.3 降级模式（无 Docker 时的基础练习）

当 Docker 不可用时，允许用户在本地临时目录中完成基础练习。

#### 设计

```go
// internal/sandbox/local.go
type LocalSandbox struct {
    workDir string  // 临时目录
}

func NewLocalSandbox() (*LocalSandbox, error)  // 创建 os.MkdirTemp
func (s *LocalSandbox) Exec(ctx context.Context, command string) (string, int, error)  // 在 workDir 中执行
func (s *LocalSandbox) Destroy() error  // 清理临时目录
func (s *LocalSandbox) InteractiveShellArgs() []string  // 返回 bash -c "cd workDir && bash"
```

#### 降级策略

在 `internal/sandbox/sandbox.go` 中定义统一接口：

```go
type Sandbox interface {
    Exec(ctx context.Context, command string) (string, int, error)
    Destroy(ctx context.Context) error
    InteractiveShellArgs() []string
}

func NewSandbox(ctx context.Context, ch *challenge.Challenge) (Sandbox, error) {
    if ch.ComposeFile != "" {
        return NewComposeSandbox(ctx, ch.Dir)
    }
    if dockerAvailable() {
        return NewDockerSandbox(ctx, "")
    }
    return NewLocalSandbox()
}
```

#### 限制

降级模式下：
- 文件操作、管道、grep/awk/sed 等基础题可正常完成
- 进程管理、网络、systemd 等需要 root 或隔离环境的题标记为 `requires_docker: true`
- TUI 显示 "⚠ 需要 Docker" 提示

### 2.4 Challenge 类型扩展

在 `challenge.Challenge` 中新增字段：

```go
ComposeFile    string `yaml:"compose_file,omitempty"`
RequiresDocker bool   `yaml:"requires_docker,omitempty"`
```

## 3. 文件结构变更

```
新增文件：
├── internal/sandbox/sandbox.go           # Sandbox 接口
├── internal/sandbox/sandbox_test.go
├── internal/sandbox/compose.go           # docker-compose 沙盒
├── internal/sandbox/compose_test.go
├── internal/sandbox/local.go             # 本地降级沙盒
├── internal/sandbox/local_test.go
├── challenges/containers/                # 35+ 容器题目

修改文件：
├── internal/challenge/types.go           # 新增 ComposeFile、RequiresDocker 字段
├── internal/challenge/types_test.go      # 测试新字段
├── internal/tui/app.go                   # 使用 Sandbox 接口启动挑战
├── internal/tui/detail.go               # 显示 Docker 要求提示
├── internal/tui/modules.go              # categoryLabels 添加 containers
└── e2e_test.go                          # 更新断言
```

## 4. 实施顺序

### P2a — 沙盒重构 + 降级模式（代码）
1. Challenge 类型扩展（TDD）
2. Sandbox 接口定义（TDD）
3. LocalSandbox 本地降级沙盒（TDD）
4. ComposeSandbox docker-compose 沙盒（TDD）
5. NewSandbox 工厂函数（TDD）
6. TUI 集成 Sandbox 接口（TDD）

### P2b — 容器题库（内容）
7. Docker 基础题库（10+ 题）
8. Dockerfile 题库（8+ 题）
9. Docker Compose 题库（7+ 题）
10. Docker 网络题库（5+ 题）
11. Docker 数据卷题库（5+ 题）
12. E2E 测试更新
