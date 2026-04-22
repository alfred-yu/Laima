# Laima（莱码）软件架构设计与系统功能设计

> **版本**：v1.0
> **日期**：2026 年 4 月 22 日
> **基于**：《代码托管平台竞品分析报告》、《Laima（莱码）产品方案》

---

## 目录

1. [架构设计原则](#1-架构设计原则)
2. [宏观架构设计](#2-宏观架构设计)
3. [系统功能设计](#3-系统功能设计)
4. [详细模块设计](#4-详细模块设计)
5. [数据库详细设计](#5-数据库详细设计)
6. [API 接口规范](#6-api-接口规范)
7. [核心业务流程设计](#7-核心业务流程设计)
8. [非功能性架构设计](#8-非功能性架构设计)

---

## 1. 架构设计原则

### 1.1 核心原则

| 原则 | 描述 | 落地策略 |
|------|------|----------|
| **轻量优先** | 最小化资源占用，512MB RAM 可运行核心功能 | 单体架构优先，Go 编译为单一二进制，按需加载模块 |
| **模块化** | 功能模块高内聚低耦合，可独立开发和测试 | Go interface 抽象模块边界，依赖注入 |
| **可扩展** | 支持水平扩展和功能插件化 | Runner/Worker 分布式架构，插件系统 |
| **安全内建** | 安全不是附加层，而是贯穿所有层级的设计 | 输入验证、RBAC、审计日志、加密存储 |
| **AI 原生** | AI 能力深度集成到核心流程中，而非外挂 | AI 审查引擎作为一等公民，参与 PR 生命周期 |
| **兼容性** | 兼容 GitHub API 格式和 CI/CD 语法 | API 兼容层、YAML 语法转换层 |

### 1.2 架构约束

| 约束 | 描述 |
|------|------|
| 单一二进制部署 | 核心服务编译为一个 Go 二进制文件，无需额外运行时 |
| 最小资源预算 | 核心服务（代码托管 + PR + Issue）≤ 512MB RAM |
| 无外部 Git 依赖 | 使用 go-git 库处理 Git 操作，不依赖系统安装的 Git |
| 数据库唯一 | 仅依赖 PostgreSQL 作为关系型数据库（Redis 为可选加速层） |
| 配置即代码 | 所有 CI/CD、审查规则、安全策略通过 YAML 配置文件管理 |

---

## 2. 宏观架构设计

### 2.1 系统分层架构

Laima 采用**改良的分层单体架构**（Modular Monolith），在保持单一进程部署优势的同时，通过清晰的模块边界实现高内聚低耦合。

```
┌─────────────────────────────────────────────────────────────────────┐
│                         接入层 (Access Layer)                        │
│                                                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │  HTTP Server  │  │  SSH Server  │  │  WebSocket   │              │
│  │  (:8080/:443) │  │  (:2222)     │  │  Server      │              │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘              │
│         │                 │                 │                        │
│  ┌──────▼─────────────────▼─────────────────▼───────┐              │
│  │              Router / Middleware                   │              │
│  │  认证 · 鉴权 · 限流 · CORS · 审计日志 · IP白名单   │              │
│  └──────────────────────┬───────────────────────────┘              │
└─────────────────────────┼───────────────────────────────────────────┘
                          │
┌─────────────────────────┼───────────────────────────────────────────┐
│                         应用层 (Application Layer)                   │
│                         │                                           │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │
│  │ Repo     │ │ PullReq  │ │ Issue    │ │ Pipeline │ │ User     │ │
│  │ Service  │ │ Service  │ │ Service  │ │ Service  │ │ Service  │ │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘ │
│       │            │            │            │            │        │
│  ┌────▼─────┐ ┌────▼─────┐ ┌────▼─────┐ ┌────▼─────┐ ┌────▼─────┐ │
│  │ AIReview │ │ Security │ │ Registry │ │ Wiki     │ │ Org      │ │
│  │ Service  │ │ Service  │ │ Service  │ │ Service  │ │ Service  │ │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘ │
└─────────────────────────┬───────────────────────────────────────────┘
                          │
┌─────────────────────────┼───────────────────────────────────────────┐
│                         领域层 (Domain Layer)                        │
│                         │                                           │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │
│  │ Repo     │ │ PullReq  │ │ Issue    │ │ Pipeline │ │ User     │ │
│  │ Domain   │ │ Domain   │ │ Domain   │ │ Domain   │ │ Domain   │ │
│  │ Model    │ │ Model    │ │ Model    │ │ Model    │ │ Model    │ │
│  │ + Rules  │ │ + Rules  │ │ + Rules  │ │ + Rules  │ │ + Rules  │ │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘ │
└─────────────────────────┬───────────────────────────────────────────┘
                          │
┌─────────────────────────┼───────────────────────────────────────────┐
│                         基础设施层 (Infrastructure Layer)             │
│                         │                                           │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │
│  │PostgreSQL│ │  Redis   │ │  MinIO   │ │Meili-   │ │ Git      │ │
│  │  Store   │ │  Store   │ │  Store   │ │ search  │ │ Store   │ │
│  │          │ │          │ │          │ │          │ │(go-git)  │ │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘ │
│                                                                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐              │
│  │  Queue   │ │  Cache   │ │  Mail    │ │  Logger  │              │
│  │(Redis    │ │ Manager  │ │ Sender   │ │ & Audit │              │
│  │ Stream)  │ │          │ │          │ │          │              │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘              │
└─────────────────────────────────────────────────────────────────────┘
                          │
┌─────────────────────────┼───────────────────────────────────────────┐
│                    外部服务层 (External Services)                    │
│                                                                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐              │
│  │  LLM     │ │  Runner  │ │  OIDC/   │ │  LDAP    │              │
│  │ Engine   │ │  Pool    │ │  SAML    │ │  Server  │              │
│  │(Ollama/  │ │(CI/CD    │ │ Provider │ │          │              │
│  │ OpenAI)  │ │ Workers) │ │          │ │          │              │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘              │
└─────────────────────────────────────────────────────────────────────┘
```

### 2.2 分层职责说明

| 层级 | 职责 | 关键组件 |
|------|------|----------|
| **接入层** | 协议处理、请求路由、认证鉴权、限流防护 | HTTP Server (Gin/Echo)、SSH Server、WebSocket Server、Middleware 链 |
| **应用层** | 业务编排、用例实现、事务管理、事件发布 | 各业务 Service（RepoService、PullReqService 等） |
| **领域层** | 领域模型、业务规则、不变量约束 | Domain Model、Domain Rules、Domain Events |
| **基础设施层** | 数据持久化、缓存、搜索、消息队列、外部通信 | Repository 实现、Cache Manager、Queue、Mail Sender |
| **外部服务层** | 外部系统对接 | LLM Engine、Runner Pool、OIDC/SAML Provider、LDAP |

### 2.3 模块划分

系统划分为 **10 个核心模块**，每个模块遵循相同的内部结构：

```
internal/
├── repo/                    # 仓库管理模块
│   ├── domain/              # 领域模型和规则
│   │   ├── model.go         # 仓库实体
│   │   ├── rules.go         # 业务规则（命名规范、权限校验等）
│   │   └── events.go        # 领域事件
│   ├── app/                 # 应用服务
│   │   ├── service.go       # 仓库应用服务
│   │   └── dto.go           # 数据传输对象
│   ├── infra/               # 基础设施
│   │   ├── repo_pg.go       # PostgreSQL 仓库实现
│   │   ├── repo_cache.go    # 缓存实现
│   │   └── git_ops.go       # Git 操作封装
│   ├── api/                 # API 层
│   │   ├── handler.go       # HTTP Handler
│   │   ├── router.go        # 路由注册
│   │   └── request.go       # 请求/响应结构体
│   └── module.go            # 模块注册（依赖注入）
│
├── pullreq/                 # Pull Request 模块
├── issue/                   # Issue 追踪模块
├── pipeline/                # CI/CD 流水线模块
├── user/                    # 用户管理模块
├── org/                     # 组织管理模块
├── aireview/                # AI 审查引擎模块
├── security/                # 安全扫描模块
├── registry/                # 容器/包仓库模块
├── wiki/                    # Wiki 模块
│
├── pkg/                     # 公共包
│   ├── auth/                # 认证鉴权
│   ├── middleware/           # 中间件
│   ├── git/                 # Git 操作封装（go-git）
│   ├── llm/                 # LLM 客户端封装
│   ├── queue/               # 消息队列封装
│   ├── cache/               # 缓存封装
│   ├── search/              # 搜索封装
│   ├── storage/             # 对象存储封装
│   ├── crypto/              # 加密工具
│   └── notification/        # 通知服务
│
└── cmd/
    └── laima/               # 主程序入口
        └── main.go
```

### 2.4 模块依赖关系

```
                    ┌─────────┐
                    │  user   │
                    └────┬────┘
                         │
              ┌──────────┼──────────┐
              │          │          │
        ┌─────▼────┐ ┌──▼───┐ ┌────▼─────┐
        │   org    │ │ repo │ │  pullreq │
        └────┬─────┘ └──┬───┘ └────┬─────┘
             │          │          │
             │     ┌────┼────┐     │
             │     │    │    │     │
        ┌────▼─────▼┐  │ ┌──▼─────▼────┐
        │   issue   │  │ │  aireview   │
        └───────────┘  │ └──────┬──────┘
                       │        │
                 ┌─────▼──┐ ┌──▼────────┐
                 │pipeline│ │  security │
                 └────────┘ └───────────┘
                       │
                 ┌─────▼──────┐
                 │  registry  │
                 └────────────┘
```

**依赖规则**：
- **单向依赖**：模块间只允许单向依赖，禁止循环依赖
- **领域隔离**：领域层不依赖任何外部模块
- **接口隔离**：模块间通过 interface 交互，不直接引用具体实现
- **事件驱动**：跨模块通信优先使用领域事件（Domain Event），降低耦合

### 2.5 全局数据流

```
用户请求
    │
    ▼
┌──────────────┐
│  接入层       │  认证 → 鉴权 → 限流 → 路由
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  应用层       │  参数校验 → 业务编排 → 事务管理
└──────┬───────┘
       │
       ├──────────────────────┐
       ▼                      ▼
┌──────────────┐      ┌──────────────┐
│  领域层       │      │  事件发布     │  领域事件 → 消息队列
│  业务规则     │      │  (异步)      │
│  状态变更     │      └──────┬───────┘
└──────┬───────┘             │
       │                     ▼
       ▼              ┌──────────────┐
┌──────────────┐      │  事件消费者   │  AI 审查 / 安全扫描 / 通知
│  基础设施层    │      │  (Worker)   │
│  数据持久化    │      └──────────────┘
│  缓存更新     │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  响应返回     │
└──────────────┘
```

---

## 3. 系统功能设计

### 3.1 功能模块全景图

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Laima 系统功能全景                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    代码托管 (repo)                            │   │
│  │  仓库CRUD · 分支管理 · 标签管理 · 代码浏览 · 代码搜索         │   │
│  │  Git LFS · 仓库镜像 · GPG签名 · 仓库模板 · 仓库Fork          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌──────────────────────┐  ┌──────────────────────────────────┐   │
│  │  Pull Request (pr)   │  │  AI 审查引擎 (aireview)          │   │
│  │  PR CRUD · 内联评论   │  │  Bug检测 · 安全检测 · 性能分析    │   │
│  │  审批规则 · 合并策略  │  │  风格检查 · 修复建议 · 审查摘要   │   │
│  │  Draft PR · CODEOWNERS│  │  多模型支持 · 自定义规则          │   │
│  │  冲突解决 · CherryPick│  │  审查模式 (宽松/标准/严格)       │   │
│  └──────────────────────┘  └──────────────────────────────────┘   │
│                                                                     │
│  ┌──────────────────────┐  ┌──────────────────────────────────┐   │
│  │  CI/CD (pipeline)    │  │  安全扫描 (security)              │   │
│  │  流水线定义 · 多阶段  │  │  SAST · 依赖扫描 · 密钥检测       │   │
│  │  Runner调度 · 矩阵   │  │  许可证合规 · DAST · 容器扫描     │   │
│  │  缓存 · 密钥管理      │  │  安全仪表盘 · 合规报告            │   │
│  │  制品管理 · 定时触发  │  │                                  │   │
│  └──────────────────────┘  └──────────────────────────────────┘   │
│                                                                     │
│  ┌──────────────────────┐  ┌──────────────────────────────────┐   │
│  │  Issue (issue)       │  │  仓库管理 (registry)              │   │
│  │  Issue CRUD · 标签   │  │  Container Registry · Package     │   │
│  │  里程碑 · 看板       │  │  Registry · 镜像签名 · 版本管理    │   │
│  │  时间追踪 · 依赖     │  │                                  │   │
│  └──────────────────────┘  └──────────────────────────────────┘   │
│                                                                     │
│  ┌──────────────────────┐  ┌──────────────────────────────────┐   │
│  │  用户/组织 (user/org) │  │  Wiki (wiki)                     │   │
│  │  用户管理 · 组织管理  │  │  项目Wiki · Pages · Mermaid      │   │
│  │  团队管理 · RBAC     │  │                                  │   │
│  │  SSO/SAML · LDAP     │  │                                  │   │
│  │  审计日志 · IP白名单  │  │                                  │   │
│  └──────────────────────┘  └──────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    通知与集成 (notification)                   │   │
│  │  Webhook · 邮件通知 · 机器人通知 · API · Marketplace         │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 3.2 模块功能清单与接口概要

#### 3.2.1 仓库管理模块 (repo)

| 功能 | 描述 | 核心接口 |
|------|------|----------|
| 创建仓库 | 创建空仓库或从模板创建 | `CreateRepo(ctx, req) → (*Repo, error)` |
| Fork 仓库 | Fork 他人仓库到自己的命名空间 | `ForkRepo(ctx, repoID, targetNS) → (*Repo, error)` |
| 导入仓库 | 从 GitHub/GitLab/Gitea 导入 | `ImportRepo(ctx, req) → (*ImportTask, error)` |
| 删除仓库 | 软删除（30天回收期） | `DeleteRepo(ctx, repoID) error` |
| 仓库镜像 | Push/Pull 镜像配置 | `SyncMirror(ctx, repoID) error` |
| 分支管理 | 创建/删除/保护分支 | `CreateBranch()`, `ProtectBranch()` |
| 标签管理 | 创建/删除标签 | `CreateTag()`, `DeleteTag()` |
| 代码浏览 | 文件树、文件内容、Blame | `GetTree()`, `GetBlob()`, `GetBlame()` |
| 代码搜索 | 全文搜索、正则搜索 | `SearchCode(ctx, query) → ([]Result, error)` |
| 仓库统计 | Star、Fork、Watch 计数 | `GetRepoStats(ctx, repoID) → (*Stats, error)` |

#### 3.2.2 Pull Request 模块 (pullreq)

| 功能 | 描述 | 核心接口 |
|------|------|----------|
| 创建 PR | 从分支创建 PR | `CreatePR(ctx, req) → (*PR, error)` |
| 更新 PR | 更新标题、描述、分支 | `UpdatePR(ctx, prID, req) error` |
| 提交审查 | 提交审查意见和评分 | `SubmitReview(ctx, prID, req) error` |
| 内联评论 | 行级/文件级评论 | `CreateComment(ctx, req) error` |
| 建议式修改 | 提供代码修改建议 | `CreateSuggestion(ctx, req) error` |
| 应用修复 | 应用 AI 或人工的修复建议 | `ApplySuggestion(ctx, commentID) error` |
| 合并 PR | 执行合并操作 | `MergePR(ctx, prID, strategy) error` |
| 关闭 PR | 关闭（不合并） | `ClosePR(ctx, prID) error` |
| 转换 Draft | Draft ↔ Ready 转换 | `ConvertDraft(ctx, prID) error` |
| 检查合并状态 | 检查是否满足合并条件 | `CheckMergeability(ctx, prID) → (*MergeCheck, error)` |

#### 3.2.3 AI 审查引擎模块 (aireview)

| 功能 | 描述 | 核心接口 |
|------|------|----------|
| 触发审查 | PR 创建/更新时自动触发 | `TriggerReview(ctx, prID) error` |
| 变更提取 | 提取 diff、上下文、文件信息 | `ExtractChanges(prID) → (*Changes, error)` |
| Bug 检测 | LLM 检测逻辑错误 | `DetectBugs(changes) → ([]Finding, error)` |
| 安全检测 | LLM + 规则检测安全问题 | `DetectSecurity(changes) → ([]Finding, error)` |
| 性能分析 | LLM 分析性能问题 | `AnalyzePerformance(changes) → ([]Finding, error)` |
| 风格检查 | 规则引擎检查代码风格 | `CheckStyle(changes) → ([]Finding, error)` |
| 生成摘要 | 生成 PR 变更摘要 | `GenerateSummary(changes) → (*Summary, error)` |
| 生成修复 | 生成一键修复补丁 | `GenerateFix(finding) → (*Patch, error)` |
| 发布审查意见 | 将结果以 PR 评论发布 | `PublishReview(prID, findings) error` |
| 管理审查规则 | CRUD 自定义审查规则 | `ManageRules(ctx, req) error` |

#### 3.2.4 CI/CD 流水线模块 (pipeline)

| 功能 | 描述 | 核心接口 |
|------|------|----------|
| 解析流水线 | 解析 .laima-ci.yml | `ParsePipeline(yaml) → (*Pipeline, error)` |
| 创建流水线 | 创建流水线实例 | `CreatePipeline(ctx, req) → (*Pipeline, error)` |
| 调度任务 | DAG 调度 Job 到 Runner | `ScheduleJobs(pipelineID) error` |
| Runner 注册 | Runner 注册和心跳 | `RegisterRunner(ctx, req) → (*Runner, error)` |
| 任务执行 | Runner 拉取和执行任务 | `AcquireJob(runnerID) → (*Job, error)` |
| 上报结果 | Runner 上报任务结果 | `ReportResult(jobID, result) error` |
| 流水线重试 | 重试失败的任务 | `RetryJob(jobID) error` |
| 取消流水线 | 取消运行中的流水线 | `CancelPipeline(pipelineID) error` |
| 制品管理 | 上传/下载构建制品 | `UploadArtifact()`, `DownloadArtifact()` |

#### 3.2.5 安全扫描模块 (security)

| 功能 | 描述 | 核心接口 |
|------|------|----------|
| SAST 扫描 | 静态应用安全测试 | `RunSAST(repoID, commitSHA) → (*ScanResult, error)` |
| 依赖扫描 | 检测依赖漏洞 (CVE) | `RunDependencyScan(repoID, commitSHA) → (*ScanResult, error)` |
| 密钥检测 | 检测硬编码密钥 | `RunSecretDetection(repoID, commitSHA) → (*ScanResult, error)` |
| 许可证检查 | 检查开源许可证 | `RunLicenseCheck(repoID, commitSHA) → (*ScanResult, error)` |
| DAST 扫描 | 动态应用安全测试 | `RunDAST(repoID, targetURL) → (*ScanResult, error)` |
| 容器扫描 | 镜像漏洞扫描 | `RunContainerScan(imageRef) → (*ScanResult, error)` |
| 聚合结果 | 聚合所有扫描结果 | `AggregateResults(repoID) → (*SecurityReport, error)` |
| 生成报告 | 生成合规报告 | `GenerateReport(repoID, format) → ([]byte, error)` |

---

## 4. 详细模块设计

### 4.1 仓库管理模块详细设计

#### 4.1.1 领域模型

```go
// internal/repo/domain/model.go

// Repository 仓库实体
type Repository struct {
    ID              int64          `json:"id"`
    Name            string         `json:"name"`              // 仓库名称
    FullName        string         `json:"full_name"`         // 完整路径 (org/repo)
    Description     string         `json:"description"`       // 描述
    OwnerID         int64          `json:"owner_id"`          // 所有者 ID
    OwnerType       OwnerType      `json:"owner_type"`        // 所有者类型 (user/org)
    Visibility      Visibility     `json:"visibility"`        // 可见性 (public/private/internal)
    DefaultBranch   string         `json:"default_branch"`    // 默认分支
    Size            int64          `json:"size"`              // 仓库大小 (bytes)
    IsFork          bool           `json:"is_fork"`           // 是否为 Fork
    ForkParentID    *int64         `json:"fork_parent_id"`    // Fork 源仓库 ID
    IsMirror        bool           `json:"is_mirror"`         // 是否为镜像
    MirrorURL       string         `json:"mirror_url"`        // 镜像源 URL
    IsTemplate      bool           `json:"is_template"`       // 是否为模板
    Topics          []string       `json:"topics"`            // 主题标签
    WebURL          string         `json:"web_url"`           // Web 访问 URL
    SSHURL          string         `json:"ssh_url"`           // SSH 克隆 URL
    HTTPURL         string         `json:"http_url"`          // HTTP 克隆 URL
    Settings        RepoSettings   `json:"settings"`          // 仓库设置
    Stats           RepoStats      `json:"stats"`             // 统计信息
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
}

type OwnerType   string
const (
    OwnerTypeUser OwnerType = "user"
    OwnerTypeOrg  OwnerType = "org"
)

type Visibility  string
const (
    VisibilityPublic   Visibility = "public"
    VisibilityPrivate  Visibility = "private"
    VisibilityInternal Visibility = "internal"
)

// RepoSettings 仓库设置
type RepoSettings struct {
    MergeStrategy      MergeStrategy `json:"merge_strategy"`       // 默认合并策略
    ReviewMode         ReviewMode    `json:"review_mode"`          // 审查模式
    RequirePR          bool          `json:"require_pr"`           // 是否强制 PR
    RequireSignedCommit bool         `json:"require_signed"`       // 是否要求签名提交
    EnableWiki         bool          `json:"enable_wiki"`          // 是否启用 Wiki
    EnableIssues       bool          `json:"enable_issues"`        // 是否启用 Issue
    EnablePipeline     bool          `json:"enable_pipeline"`      // 是否启用 CI/CD
    EnableAIReview     bool          `json:"enable_ai_review"`     // 是否启用 AI 审查
    AIReviewRules      []string      `json:"ai_review_rules"`      // AI 审查规则 ID 列表
    BranchProtection   []BranchProtection `json:"branch_protection"` // 分支保护规则
}

// BranchProtection 分支保护规则
type BranchProtection struct {
    BranchPattern     string   `json:"branch_pattern"`      // 分支匹配模式
    RequirePR         bool     `json:"require_pr"`          // 需要 PR
    RequiredApprovals int      `json:"required_approvals"`  // 需要审批人数
    RequireCI         bool     `json:"require_ci"`          // 需要 CI 通过
    RequireSigned     bool     `json:"require_signed"`      // 需要签名提交
    AllowedRoles      []string `json:"allowed_roles"`       // 允许直接 push 的角色
    CODEOWNERS        bool     `json:"codeowners"`          // 启用 CODEOWNERS
}

type MergeStrategy string
const (
    MergeStrategyMerge  MergeStrategy = "merge"
    MergeStrategySquash MergeStrategy = "squash"
    MergeStrategyRebase MergeStrategy = "rebase"
)

type ReviewMode string
const (
    ReviewModeRelaxed  ReviewMode = "relaxed"
    ReviewModeStandard ReviewMode = "standard"
    ReviewModeStrict   ReviewMode = "strict"
)
```

#### 4.1.2 应用服务接口

```go
// internal/repo/app/service.go

// RepoService 仓库应用服务接口
type RepoService interface {
    // 仓库 CRUD
    CreateRepo(ctx context.Context, req *CreateRepoRequest) (*Repository, error)
    GetRepo(ctx context.Context, repoID int64) (*Repository, error)
    GetRepoByPath(ctx context.Context, fullPath string) (*Repository, error)
    UpdateRepo(ctx context.Context, repoID int64, req *UpdateRepoRequest) (*Repository, error)
    DeleteRepo(ctx context.Context, repoID int64) error
    ListRepos(ctx context.Context, filter *RepoFilter) ([]*Repository, int, error)

    // Fork & 导入
    ForkRepo(ctx context.Context, repoID int64, targetNamespace string) (*Repository, error)
    ImportRepo(ctx context.Context, req *ImportRepoRequest) (*ImportTask, error)

    // 分支操作
    CreateBranch(ctx context.Context, repoID int64, req *CreateBranchRequest) (*Branch, error)
    DeleteBranch(ctx context.Context, repoID int64, branch string) error
    ListBranches(ctx context.Context, repoID int64) ([]*Branch, error)
    ProtectBranch(ctx context.Context, repoID int64, rule *BranchProtection) error

    // 标签操作
    CreateTag(ctx context.Context, repoID int64, req *CreateTagRequest) (*Tag, error)
    DeleteTag(ctx context.Context, repoID int64, tagName string) error
    ListTags(ctx context.Context, repoID int64) ([]*Tag, error)

    // 代码浏览
    GetTree(ctx context.Context, repoID int64, ref, path string) (*Tree, error)
    GetBlob(ctx context.Context, repoID int64, ref, path string) (*Blob, error)
    GetBlame(ctx context.Context, repoID int64, ref, path string) ([]*BlameLine, error)
    GetRawFile(ctx context.Context, repoID int64, ref, path string) ([]byte, error)

    // 代码搜索
    SearchCode(ctx context.Context, query *SearchQuery) ([]*SearchResult, int, error)

    // 统计
    StarRepo(ctx context.Context, repoID int64) error
    UnstarRepo(ctx context.Context, repoID int64) error
    WatchRepo(ctx context.Context, repoID int64) error
}
```

#### 4.1.3 Git 操作封装

```go
// pkg/git/operations.go

// GitOperations Git 操作接口
type GitOperations interface {
    // 仓库操作
    InitRepository(path string, bare bool) error
    CloneRepository(src, dst string, opts *CloneOptions) error
    GetRepository(path string) (*git.Repository, error)

    // 引用操作
    GetBranches(repo *git.Repository) ([]*plumbing.Reference, error)
    GetTags(repo *git.Repository) ([]*plumbing.Reference, error)
    CreateBranch(repo *git.Repository, name, base string) error
    DeleteBranch(repo *git.Repository, name string) error
    CreateTag(repo *git.Repository, name, message string) error

    // 对象操作
    GetTree(repo *git.Repository, ref string) (*object.Tree, error)
    GetBlob(repo *git.Repository, hash plumbing.Hash) (object.Blob, error)
    GetCommit(repo *git.Repository, ref string) (*object.Commit, error)
    GetBlame(repo *git.Repository, path, ref string) (*git.BlameResult, error)

    // Diff 操作
    GetDiff(repo *git.Repository, baseRef, headRef string) ([]*DiffFile, error)
    GetCommitDiff(repo *git.Repository, commit *object.Commit) ([]*DiffFile, error)

    // 合并操作
    MergeBase(repo *git.Repository, ref1, ref2 string) (plumbing.Hash, error)
    CanMerge(repo *git.Repository, base, ours, theirs string) (*MergeStatus, error)
    Merge(repo *git.Repository, base, ours, theirs string, strategy MergeStrategy) error

    // LFS 操作
    LFSUpload(repo *git.Repository, oid string, size int64, reader io.Reader) error
    LFSDownload(repo *git.Repository, oid string) (io.ReadCloser, error)
}
```

### 4.2 Pull Request 模块详细设计

#### 4.2.1 领域模型

```go
// internal/pullreq/domain/model.go

// PullRequest PR 实体
type PullRequest struct {
    ID              int64          `json:"id"`
    Number          int64          `json:"number"`            // PR 编号（仓库内自增）
    Title           string         `json:"title"`
    Description     string         `json:"description"`       // 支持 Markdown
    RepositoryID    int64          `json:"repository_id"`
    AuthorID        int64          `json:"author_id"`
    SourceRepoID    int64          `json:"source_repo_id"`    // 源仓库（Fork PR 时不同）
    SourceBranch    string         `json:"source_branch"`
    TargetBranch    string         `json:"target_branch"`
    State           PRState        `json:"state"`             // open/draft/merged/closed
    MergeState      MergeState     `json:"merge_state"`       // mergeable/unstable/conflict/blocked
    MergeStrategy   *MergeStrategy `json:"merge_strategy"`    // 合并策略（nil=使用仓库默认）
    ReviewMode      ReviewMode     `json:"review_mode"`       // 审查模式
    HeadCommitSHA   string         `json:"head_commit_sha"`
    BaseCommitSHA   string         `json:"base_commit_sha"`
    MergeCommitSHA  *string        `json:"merge_commit_sha"`  // 合并后的 commit SHA
    MergedBy        *int64         `json:"merged_by"`
    MergedAt        *time.Time     `json:"merged_at"`
    ClosedAt        *time.Time     `json:"closed_at"`
    Draft           bool           `json:"is_draft"`
    Labels          []string       `json:"labels"`
    MilestoneID     *int64         `json:"milestone_id"`
    RelatedIssues   []int64        `json:"related_issues"`    // 关联的 Issue ID
    AIReviewStatus  AIReviewStatus `json:"ai_review_status"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
}

type PRState string
const (
    PRStateOpen   PRState = "open"
    PRStateDraft  PRState = "draft"
    PRStateMerged PRState = "merged"
    PRStateClosed PRState = "closed"
)

type MergeState string
const (
    MergeStateMergeable  MergeState = "mergeable"
    MergeStateUnstable   MergeState = "unstable"   // CI 未通过
    MergeStateConflict   MergeState = "conflict"   // 有冲突
    MergeStateBlocked    MergeState = "blocked"    // 审查未满足
    MergeStateChecking   MergeState = "checking"   // 检查中
)

type AIReviewStatus string
const (
    AIReviewPending   AIReviewStatus = "pending"
    AIReviewRunning   AIReviewStatus = "running"
    AIReviewCompleted AIReviewStatus = "completed"
    AIReviewFailed    AIReviewStatus = "failed"
    AIReviewSkipped   AIReviewStatus = "skipped"
)

// Review 审查记录
type Review struct {
    ID          int64       `json:"id"`
    PRID        int64       `json:"pr_id"`
    ReviewerID  int64       `json:"reviewer_id"`
    State       ReviewState `json:"state"`        // approved/changes_requested/commented/pending
    Score       int         `json:"score"`         // 严格模式: +2/+1/0/-1/-2
    Body        string      `json:"body"`
    SubmittedAt time.Time   `json:"submitted_at"`
}

type ReviewState string
const (
    ReviewApproved         ReviewState = "approved"
    ReviewChangesRequested ReviewState = "changes_requested"
    ReviewCommented        ReviewState = "commented"
    ReviewPending          ReviewState = "pending"
)

// ReviewComment 审查评论
type ReviewComment struct {
    ID          int64     `json:"id"`
    PRID        int64     `json:"pr_id"`
    ReviewID    *int64    `json:"review_id"`      // 关联的 Review
    AuthorID    int64     `json:"author_id"`
    Type        CommentType `json:"type"`          // human/ai
    Path        string    `json:"path"`           // 文件路径
    Line        *int      `json:"line"`           // 行号（nil=文件级评论）
    DiffHunk    string    `json:"diff_hunk"`      // 关联的 Diff 片段
    Body        string    `json:"body"`
    Resolution  *string   `json:"resolution"`     // resolved/unresolved/nil
    Suggestion  *string   `json:"suggestion"`     // 建议式修改内容
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CommentType string
const (
    CommentTypeHuman CommentType = "human"
    CommentTypeAI   CommentType = "ai"
)

// MergeCheckResult 合并检查结果
type MergeCheckResult struct {
    CanMerge       bool              `json:"can_merge"`
    MergeState     MergeState        `json:"merge_state"`
    Checks         map[string]*Check `json:"checks"`  // 各项检查结果
    BlockReasons   []string          `json:"block_reasons"`
}

type Check struct {
    Name     string `json:"name"`
    Status   string `json:"status"`  // passed/failed/pending/skipped
    Required bool   `json:"required"`
    Message  string `json:"message"`
}
```

#### 4.2.2 应用服务接口

```go
// internal/pullreq/app/service.go

type PullReqService interface {
    // PR 生命周期
    CreatePR(ctx context.Context, req *CreatePRRequest) (*PullRequest, error)
    GetPR(ctx context.Context, prID int64) (*PullRequest, error)
    GetPRByNumber(ctx context.Context, repoID, number int64) (*PullRequest, error)
    UpdatePR(ctx context.Context, prID int64, req *UpdatePRRequest) (*PullRequest, error)
    ListPRs(ctx context.Context, filter *PRFilter) ([]*PullRequest, int, error)

    // 审查操作
    SubmitReview(ctx context.Context, prID int64, req *SubmitReviewRequest) (*Review, error)
    CreateComment(ctx context.Context, req *CreateCommentRequest) (*ReviewComment, error)
    ResolveComment(ctx context.Context, commentID int64, resolved bool) error
    ApplySuggestion(ctx context.Context, commentID int64) error

    // 合并操作
    CheckMergeability(ctx context.Context, prID int64) (*MergeCheckResult, error)
    MergePR(ctx context.Context, prID int64, req *MergeRequest) error
    ClosePR(ctx context.Context, prID int64, reason string) error
    ReopenPR(ctx context.Context, prID int64) error

    // Draft 操作
    ConvertToDraft(ctx context.Context, prID int64) error
    ConvertFromDraft(ctx context.Context, prID int64) error

    // Diff 操作
    GetDiff(ctx context.Context, prID int64) ([]*DiffFile, error)
    GetCommits(ctx context.Context, prID int64) ([]*PRCommit, error)
}
```

### 4.3 AI 审查引擎模块详细设计

#### 4.3.1 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                      AI 审查引擎架构                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Review Orchestrator                    │   │
│  │  (审查编排器 - 协调整个审查流程)                          │   │
│  └──────────────┬──────────────────────────────┬────────────┘   │
│                 │                              │                 │
│  ┌──────────────▼──────────┐    ┌─────────────▼─────────────┐  │
│  │   Change Analyzer       │    │   Context Builder          │  │
│  │   (变更分析器)           │    │   (上下文构建器)            │  │
│  │   · 提取 diff           │    │   · 项目结构               │  │
│  │   · 文件分类            │    │   · 历史审查记录            │  │
│  │   · 语言检测            │    │   · 团队编码规范            │  │
│  │   · 变更规模评估        │    │   · 关联文件内容            │  │
│  └──────────────┬──────────┘    └─────────────┬─────────────┘  │
│                 │                              │                 │
│                 └──────────────┬───────────────┘                 │
│                                │                                 │
│  ┌─────────────────────────────▼─────────────────────────────┐  │
│  │                    Analyzer Pool                           │  │
│  │  (分析器池 - 并行执行多个分析任务)                          │  │
│  │                                                            │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐    │  │
│  │  │ Bug      │ │ Security │ │ Perf     │ │ Style    │    │  │
│  │  │ Analyzer │ │ Analyzer │ │ Analyzer │ │ Analyzer │    │  │
│  │  │ (LLM)    │ │(Rule+LLM)│ │ (LLM)    │ │ (Rule)   │    │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘    │  │
│  └─────────────────────────────┬─────────────────────────────┘  │
│                                │                                 │
│  ┌─────────────────────────────▼─────────────────────────────┐  │
│  │                   Result Aggregator                        │  │
│  │  (结果聚合器 - 去重、排序、优先级分级)                       │  │
│  │  · 去重（同一问题多个分析器发现）                            │  │
│  │  · 优先级排序（critical > high > medium > low）             │  │
│  │  · 生成审查摘要                                            │  │
│  │  · 生成修复建议                                            │  │
│  └─────────────────────────────┬─────────────────────────────┘  │
│                                │                                 │
│  ┌─────────────────────────────▼─────────────────────────────┐  │
│  │                   Review Publisher                         │  │
│  │  (审查发布器 - 将结果发布为 PR 评论)                         │  │
│  │  · 创建 AI 审查摘要评论                                    │  │
│  │  · 创建内联审查评论                                        │  │
│  │  · 更新 PR AI 审查状态                                     │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    LLM Provider                            │  │
│  │  (LLM 提供者抽象层 - 统一接口，多后端)                      │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐    │  │
│  │  │ OpenAI   │ │ Ollama   │ │ Qwen     │ │ DeepSeek │    │  │
│  │  │ Provider │ │ Provider │ │ Provider │ │ Provider │    │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘    │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

#### 4.3.2 核心接口设计

```go
// internal/aireview/app/service.go

// AIReviewService AI 审查服务接口
type AIReviewService interface {
    // 触发审查
    TriggerReview(ctx context.Context, prID int64) error

    // 获取审查状态
    GetReviewStatus(ctx context.Context, prID int64) (*AIReviewStatus, error)

    // 获取审查结果
    GetReviewResult(ctx context.Context, prID int64) (*AIReviewResult, error)

    // 手动重新触发
    RetryReview(ctx context.Context, prID int64) error

    // 管理审查规则
    CreateRule(ctx context.Context, req *CreateAIRuleRequest) (*AIRule, error)
    UpdateRule(ctx context.Context, ruleID int64, req *UpdateAIRuleRequest) error
    DeleteRule(ctx context.Context, ruleID int64) error
    ListRules(ctx context.Context, repoID int64) ([]*AIRule, error)
}

// Analyzer 分析器接口
type Analyzer interface {
    Name() string
    Analyze(ctx context.Context, req *AnalyzeRequest) ([]*Finding, error)
    SupportedLanguages() []string
}

// LLMProvider LLM 提供者接口
type LLMProvider interface {
    Name() string
    ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    StreamChatCompletion(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error)
    IsAvailable(ctx context.Context) bool
}

// Finding 审查发现
type Finding struct {
    ID            string         `json:"id"`
    Analyzer      string         `json:"analyzer"`       // 来源分析器
    Severity      Severity       `json:"severity"`       // 严重性
    Category      FindingCategory `json:"category"`      // 分类
    Path          string         `json:"path"`           // 文件路径
    Line          *int           `json:"line"`           // 行号
    Message       string         `json:"message"`        // 问题描述
    Suggestion    *string        `json:"suggestion"`     // 修复建议
    FixPatch      *string        `json:"fix_patch"`      // 修复补丁
    Confidence    float64        `json:"confidence"`     // 置信度 (0-1)
    RuleID        *string        `json:"rule_id"`        // 关联的规则 ID
}

type Severity string
const (
    SeverityCritical Severity = "critical"
    SeverityHigh     Severity = "high"
    SeverityMedium   Severity = "medium"
    SeverityLow      Severity = "low"
    SeverityInfo     Severity = "info"
)

type FindingCategory string
const (
    CategoryBug       FindingCategory = "bug"
    CategorySecurity  FindingCategory = "security"
    CategoryPerformance FindingCategory = "performance"
    CategoryStyle     FindingCategory = "style"
    CategoryBestPractice FindingCategory = "best_practice"
)
```

#### 4.3.3 Prompt 工程设计

```go
// internal/aireview/prompts/templates.go

// PromptTemplate 审查 Prompt 模板
var PromptTemplates = map[string]string{

    // Bug 检测 Prompt
    "bug_detection": `你是一位经验丰富的代码审查专家。请审查以下代码变更，检测潜在的 Bug。

## 项目信息
- 语言: {{.Language}}
- 框架: {{.Framework}}
- 项目描述: {{.ProjectDescription}}

## 代码变更
{{.Diff}}

## 审查要求
1. 重点关注：空指针、边界条件、竞态条件、资源泄漏、类型错误
2. 不要报告代码风格问题（由专门的分析器处理）
3. 每个发现必须包含：文件路径、行号、问题描述、修复建议
4. 仅报告你确信是 Bug 的问题，不要猜测

## 输出格式 (JSON)
{
  "findings": [
    {
      "path": "文件路径",
      "line": 行号,
      "message": "问题描述",
      "suggestion": "修复建议",
      "confidence": 0.0-1.0
    }
  ]
}`,

    // 安全检测 Prompt
    "security_detection": `你是一位安全审计专家。请审查以下代码变更，检测安全漏洞。

## 代码变更
{{.Diff}}

## 审查要求
1. 重点关注：SQL 注入、XSS、CSRF、路径遍历、硬编码密钥、不安全的反序列化
2. 参考 OWASP Top 10 进行检查
3. 每个发现必须包含 CWE 编号（如适用）

## 输出格式 (JSON)
{
  "findings": [
    {
      "path": "文件路径",
      "line": 行号,
      "message": "问题描述",
      "cwe": "CWE-XXX",
      "suggestion": "修复建议",
      "confidence": 0.0-1.0
    }
  ]
}`,

    // 审查摘要 Prompt
    "summary_generation": `请为以下 Pull Request 生成简洁的变更摘要。

## PR 标题
{{.PRTitle}}

## PR 描述
{{.PRDescription}}

## 变更文件列表
{{.ChangedFiles}}

## 变更统计
- 新增 {{.Additions}} 行
- 删除 {{.Deletions}} 行
- 涉及 {{.FileCount}} 个文件

## 要求
用 2-3 句话概括本次变更的目的和影响。使用中文。`,
}
```

### 4.4 CI/CD 流水线模块详细设计

#### 4.4.1 领域模型

```go
// internal/pipeline/domain/model.go

// Pipeline 流水线
type Pipeline struct {
    ID           int64        `json:"id"`
    RepoID       int64        `json:"repo_id"`
    CommitSHA    string       `json:"commit_sha"`
    Ref          string       `json:"ref"`           // 分支/标签引用
    Trigger      TriggerType  `json:"trigger"`       // push/pr/manual/schedule
    TriggeredBy  int64        `json:"triggered_by"`
    Status       PipelineStatus `json:"status"`
    YAMLContent  string       `json:"yaml_content"`
    Jobs         []*Job       `json:"jobs"`
    Variables    []Variable   `json:"variables"`
    StartedAt    *time.Time   `json:"started_at"`
    FinishedAt   *time.Time   `json:"finished_at"`
    Duration     *int64       `json:"duration"`      // 持续时间（秒）
    CreatedAt    time.Time    `json:"created_at"`
}

type PipelineStatus string
const (
    PipelinePending   PipelineStatus = "pending"
    PipelineRunning   PipelineStatus = "running"
    PipelineSuccess   PipelineStatus = "success"
    PipelineFailed    PipelineStatus = "failed"
    PipelineCanceled  PipelineStatus = "canceled"
    PipelineSkipped   PipelineStatus = "skipped"
)

// Job 流水线任务
type Job struct {
    ID           int64       `json:"id"`
    PipelineID   int64       `json:"pipeline_id"`
    Name         string      `json:"name"`
    Stage        string      `json:"stage"`
    Status       JobStatus   `json:"status"`
    DependsOn    []string    `json:"depends_on"`    // 依赖的 Job 名称
    RunnerID     *int64      `json:"runner_id"`     // 执行的 Runner
    RunnerLabels []string    `json:"runner_labels"` // 要求的 Runner 标签
    Steps        []Step      `json:"steps"`
    Log          string      `json:"log"`           // 任务日志
    StartedAt    *time.Time  `json:"started_at"`
    FinishedAt   *time.Time  `json:"finished_at"`
    Duration     *int64      `json:"duration"`
    Artifacts    []Artifact  `json:"artifacts"`
}

type JobStatus string
const (
    JobPending   JobStatus = "pending"
    JobQueued    JobStatus = "queued"
    JobRunning   JobStatus = "running"
    JobSuccess   JobStatus = "success"
    JobFailed    JobStatus = "failed"
    JobCanceled  JobStatus = "canceled"
)

// Step 任务步骤
type Step struct {
    Name    string `json:"name"`
    Uses    string `json:"uses"`     // Action 引用
    Run     string `json:"run"`      // Shell 命令
    Env     map[string]string `json:"env"`
    With    map[string]interface{} `json:"with"`
    Status  StepStatus `json:"status"`
}

// Runner CI/CD Runner
type Runner struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    RepoID      *int64    `json:"repo_id"`       // 仓库级 Runner（nil=全局）
    OrgID       *int64    `json:"org_id"`        // 组织级 Runner
    Labels      []string  `json:"labels"`
    Status      RunnerStatus `json:"status"`
    LastPingAt  time.Time `json:"last_ping_at"`
    Version     string    `json:"version"`
    CreatedAt   time.Time `json:"created_at"`
}

type RunnerStatus string
const (
    RunnerOnline  RunnerStatus = "online"
    RunnerOffline RunnerStatus = "offline"
)
```

#### 4.4.2 流水线 YAML 语法（兼容 GitHub Actions）

```yaml
# .laima-ci.yml 示例
name: Build and Test

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  GO_VERSION: '1.22'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      - run: go build ./...
      - run: go test -race -coverprofile=coverage.out ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: golangci/golangci-lint-action@v4

  security-scan:
    needs: [build]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: laima security scan --type sast

  deploy:
    needs: [build, lint, security-scan]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v4
      - run: docker build -t myapp:${{ github.sha }} .
      - run: docker push registry.laima.cloud/myapp:${{ github.sha }}
```

### 4.5 安全扫描模块详细设计

```go
// internal/security/domain/model.go

// ScanResult 扫描结果
type ScanResult struct {
    ID          int64          `json:"id"`
    RepoID      int64          `json:"repo_id"`
    CommitSHA   string         `json:"commit_sha"`
    Branch      string         `json:"branch"`
    ScanType    ScanType       `json:"scan_type"`
    Status      ScanStatus     `json:"status"`
    Severity    Severity       `json:"severity"`       // 最高严重性
    Findings    []ScanFinding  `json:"findings"`
    Summary     string         `json:"summary"`
    StartedAt   time.Time      `json:"started_at"`
    FinishedAt  time.Time      `json:"finished_at"`
    Duration    int64          `json:"duration"`       // 毫秒
}

type ScanType string
const (
    ScanTypeSAST          ScanType = "sast"
    ScanTypeDependency    ScanType = "dependency"
    ScanTypeSecret        ScanType = "secret"
    ScanTypeLicense       ScanType = "license"
    ScanTypeDAST          ScanType = "dast"
    ScanTypeContainer     ScanType = "container"
)

// ScanFinding 扫描发现
type ScanFinding struct {
    ID          string  `json:"id"`
    RuleID      string  `json:"rule_id"`       // 规则 ID (如 CWE-79)
    Severity    Severity `json:"severity"`
    Category    string  `json:"category"`      // 分类
    Path        string  `json:"path"`          // 文件路径
    Line        *int    `json:"line"`          // 行号
    Message     string  `json:"message"`       // 问题描述
    Remediation string  `json:"remediation"`   // 修复建议
    Reference   string  `json:"reference"`     // 参考链接
    Confidence  float64 `json:"confidence"`    // 置信度
    CVE         *string `json:"cve"`           // CVE 编号（依赖扫描）
    Package     *string `json:"package"`       // 包名（依赖扫描）
    Version     *string `json:"version"`       // 版本（依赖扫描）
    License     *string `json:"license"`       // 许可证（许可证检查）
}

// Scanner 扫描器接口
type Scanner interface {
    Name() string
    ScanType() ScanType
    Scan(ctx context.Context, req *ScanRequest) (*ScanResult, error)
    IsAvailable(ctx context.Context) bool
}
```

---

## 5. 数据库详细设计

### 5.1 ER 关系图

```
┌──────────┐       ┌──────────────┐       ┌──────────────┐
│  users   │──1:N──│ repositories │──1:N──│ pull_requests │
└────┬─────┘       └──────┬───────┘       └──────┬───────┘
     │                    │                      │
     │              ┌─────┼──────┐         ┌────┼────┐
     │              │     │      │         │    │    │
     ▼              ▼     ▼      ▼         ▼    ▼    ▼
┌──────────┐ ┌──────┐┌────┐┌──────┐ ┌────────┐┌────┐┌──────────┐
│orgs      │ │issues││wiki││labels│ │reviews ││comm│ │ai_review │
│          │ │      ││    ││      │ │        ││ents│ │_results  │
└────┬─────┘ └──┬───┘└────┘└──────┘ └────────┘└────┘ └──────────┘
     │          │
     ▼          ▼
┌──────────┐ ┌──────────┐
│ teams    │ │milestones│
│          │ │          │
└──────────┘ └──────────┘

┌──────────────┐       ┌──────────────┐
│ repositories │──1:N──│  pipelines   │──1:N──│ pipeline_jobs │
└──────────────┘       └──────────────┘       └───────────────┘

┌──────────────┐       ┌──────────────┐
│ repositories │──1:N──│security_scans│──1:N──│scan_findings  │
└──────────────┘       └──────────────┘       └───────────────┘

┌──────────────┐       ┌──────────────┐
│ repositories │──1:N──│   runners    │
└──────────────┘       └──────────────┘
```

### 5.2 完整表结构设计

#### 5.2.1 用户与组织

```sql
-- 用户表
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(39) NOT NULL UNIQUE,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    display_name    VARCHAR(255),
    avatar_url      VARCHAR(500),
    bio             TEXT,
    status          VARCHAR(20) DEFAULT 'active',  -- active/disabled/suspended
    language        VARCHAR(10) DEFAULT 'zh-CN',
    theme           VARCHAR(10) DEFAULT 'dark',     -- dark/light
    settings        JSONB DEFAULT '{}',
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- 组织表
CREATE TABLE organizations (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(39) NOT NULL UNIQUE,
    display_name    VARCHAR(255),
    description     TEXT,
    avatar_url      VARCHAR(500),
    visibility      VARCHAR(10) DEFAULT 'public',  -- public/private
    owner_id        BIGINT NOT NULL REFERENCES users(id),
    settings        JSONB DEFAULT '{}',
    max_repos       INT DEFAULT -1,                -- -1 = 无限制
    max_members     INT DEFAULT -1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 组织成员表
CREATE TABLE org_members (
    id              BIGSERIAL PRIMARY KEY,
    org_id          BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(20) NOT NULL DEFAULT 'member',  -- owner/admin/member
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, user_id)
);

-- 团队表
CREATE TABLE teams (
    id              BIGSERIAL PRIMARY KEY,
    org_id          BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    permission      VARCHAR(20) NOT NULL DEFAULT 'read',  -- read/write/admin
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(org_id, name)
);

-- 团队成员表
CREATE TABLE team_members (
    team_id         BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(20) NOT NULL DEFAULT 'member',
    PRIMARY KEY(team_id, user_id)
);

-- 仓库权限表（仓库级 RBAC）
CREATE TABLE repo_collaborators (
    repo_id         BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(20) NOT NULL,  -- owner/admin/maintainer/developer/reporter/guest
    PRIMARY KEY(repo_id, user_id)
);
```

#### 5.2.2 仓库

```sql
-- 仓库表
CREATE TABLE repositories (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    full_path       VARCHAR(255) NOT NULL UNIQUE,  -- owner/repo
    description     TEXT DEFAULT '',
    owner_id        BIGINT NOT NULL,
    owner_type      VARCHAR(10) NOT NULL DEFAULT 'user',  -- user/org
    visibility      VARCHAR(10) NOT NULL DEFAULT 'private',
    default_branch  VARCHAR(255) NOT NULL DEFAULT 'main',
    size            BIGINT DEFAULT 0,
    is_fork         BOOLEAN DEFAULT FALSE,
    fork_parent_id  BIGINT REFERENCES repositories(id),
    is_mirror       BOOLEAN DEFAULT FALSE,
    mirror_url      VARCHAR(500),
    mirror_sync_at  TIMESTAMPTZ,
    is_template     BOOLEAN DEFAULT FALSE,
    topics          TEXT[],
    web_url         VARCHAR(500),
    ssh_url         VARCHAR(500),
    http_url        VARCHAR(500),
    settings        JSONB DEFAULT '{}',
    stars_count     INT DEFAULT 0,
    forks_count     INT DEFAULT 0,
    watches_count   INT DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_repos_owner ON repositories(owner_id, owner_type);
CREATE INDEX idx_repos_visibility ON repositories(visibility);
CREATE INDEX idx_repos_fork ON repositories(fork_parent_id) WHERE is_fork = TRUE;

-- 分支保护规则表
CREATE TABLE branch_protections (
    id              BIGSERIAL PRIMARY KEY,
    repo_id         BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    branch_pattern  VARCHAR(255) NOT NULL,     -- glob 模式，如 "main", "release/*"
    require_pr      BOOLEAN DEFAULT TRUE,
    required_approvals INT DEFAULT 1,
    require_ci      BOOLEAN DEFAULT TRUE,
    require_signed  BOOLEAN DEFAULT FALSE,
    allowed_roles   TEXT[],
    enable_codeowners BOOLEAN DEFAULT FALSE,
    review_mode     VARCHAR(10) DEFAULT 'standard',  -- relaxed/standard/strict
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 仓库 Star 表
CREATE TABLE repo_stars (
    repo_id         BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(repo_id, user_id)
);

-- 仓库 Watch 表
CREATE TABLE repo_watches (
    repo_id         BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mode            VARCHAR(10) DEFAULT 'watching',  -- watching/ignoring
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(repo_id, user_id)
);
```

#### 5.2.3 Pull Request

```sql
-- Pull Request 表
CREATE TABLE pull_requests (
    id              BIGSERIAL PRIMARY KEY,
    number          BIGINT NOT NULL,               -- 仓库内自增编号
    repo_id         BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT DEFAULT '',
    author_id       BIGINT NOT NULL REFERENCES users(id),
    source_repo_id  BIGINT NOT NULL REFERENCES repositories(id),
    source_branch   VARCHAR(255) NOT NULL,
    target_branch   VARCHAR(255) NOT NULL,
    state           VARCHAR(10) NOT NULL DEFAULT 'open',  -- open/draft/merged/closed
    merge_state     VARCHAR(15) DEFAULT 'checking',  -- mergeable/unstable/conflict/blocked/checking
    merge_strategy  VARCHAR(10),                    -- merge/squash/rebase (NULL=使用仓库默认)
    head_commit_sha VARCHAR(40) NOT NULL,
    base_commit_sha VARCHAR(40) NOT NULL,
    merge_commit_sha VARCHAR(40),
    merged_by       BIGINT REFERENCES users(id),
    merged_at       TIMESTAMPTZ,
    closed_at       TIMESTAMPTZ,
    is_draft        BOOLEAN DEFAULT FALSE,
    labels          TEXT[],
    milestone_id    BIGINT,
    ai_review_status VARCHAR(15) DEFAULT 'pending',  -- pending/running/completed/failed/skipped
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(repo_id, number)
);

CREATE INDEX idx_prs_repo_state ON pull_requests(repo_id, state);
CREATE INDEX idx_prs_author ON pull_requests(author_id);
CREATE INDEX idx_prs_target ON pull_requests(repo_id, target_branch);

-- PR 审查记录表
CREATE TABLE reviews (
    id              BIGSERIAL PRIMARY KEY,
    pr_id           BIGINT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id     BIGINT NOT NULL REFERENCES users(id),
    state           VARCHAR(20) NOT NULL,  -- approved/changes_requested/commented/pending
    score           INT DEFAULT 0,          -- 严格模式: +2/+1/0/-1/-2
    body            TEXT DEFAULT '',
    submitted_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reviews_pr ON reviews(pr_id);

-- PR 审查评论表
CREATE TABLE review_comments (
    id              BIGSERIAL PRIMARY KEY,
    pr_id           BIGINT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    review_id       BIGINT REFERENCES reviews(id) ON DELETE SET NULL,
    author_id       BIGINT NOT NULL REFERENCES users(id),
    comment_type    VARCHAR(10) NOT NULL DEFAULT 'human',  -- human/ai
    path            VARCHAR(500) NOT NULL,
    line            INT,
    diff_hunk       TEXT,
    body            TEXT NOT NULL,
    resolution      VARCHAR(10),            -- resolved/unresolved
    suggestion      TEXT,                   -- 建议式修改内容
    in_reply_to     BIGINT REFERENCES review_comments(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_pr ON review_comments(pr_id);
CREATE INDEX idx_comments_path ON review_comments(pr_id, path);

-- PR 关联 Issue 表
CREATE TABLE pr_issue_refs (
    pr_id           BIGINT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    issue_id        BIGINT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    ref_type        VARCHAR(10) NOT NULL DEFAULT 'closes',  -- closes/relates/mentions
    PRIMARY KEY(pr_id, issue_id)
);
```

#### 5.2.4 AI 审查

```sql
-- AI 审查结果表
CREATE TABLE ai_review_results (
    id              BIGSERIAL PRIMARY KEY,
    pr_id           BIGINT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    status          VARCHAR(15) NOT NULL,  -- pending/running/completed/failed
    summary         TEXT,                   -- AI 生成的审查摘要
    total_findings  INT DEFAULT 0,
    critical_count  INT DEFAULT 0,
    high_count      INT DEFAULT 0,
    medium_count    INT DEFAULT 0,
    low_count       INT DEFAULT 0,
    model_used      VARCHAR(50),            -- 使用的 LLM 模型
    duration_ms     INT,
    triggered_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ
);

-- AI 审查发现表
CREATE TABLE ai_review_findings (
    id              BIGSERIAL PRIMARY KEY,
    review_result_id BIGINT NOT NULL REFERENCES ai_review_results(id) ON DELETE CASCADE,
    analyzer        VARCHAR(50) NOT NULL,   -- bug_detection/security/performance/style
    severity        VARCHAR(10) NOT NULL,
    category        VARCHAR(30) NOT NULL,
    path            VARCHAR(500) NOT NULL,
    line            INT,
    message         TEXT NOT NULL,
    suggestion      TEXT,
    fix_patch       TEXT,
    confidence      FLOAT DEFAULT 0.0,
    rule_id         BIGINT REFERENCES ai_review_rules(id),
    comment_id      BIGINT REFERENCES review_comments(id),  -- 关联的 PR 评论
    status          VARCHAR(10) DEFAULT 'open',  -- open/applied/dismissed
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_findings_review ON ai_review_findings(review_result_id);
CREATE INDEX idx_findings_severity ON ai_review_findings(severity);

-- AI 审查规则表
CREATE TABLE ai_review_rules (
    id              BIGSERIAL PRIMARY KEY,
    repo_id         BIGINT REFERENCES repositories(id) ON DELETE CASCADE,  -- NULL=全局规则
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    category        VARCHAR(30) NOT NULL,   -- bug/security/performance/style/custom
    language        VARCHAR(50),            -- 编程语言（NULL=所有语言）
    prompt_template TEXT NOT NULL,           -- 自定义 Prompt 模板
    severity        VARCHAR(10) DEFAULT 'medium',
    enabled         BOOLEAN DEFAULT TRUE,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### 5.2.5 CI/CD 流水线

```sql
-- 流水线表
CREATE TABLE pipelines (
    id              BIGSERIAL PRIMARY KEY,
    repo_id         BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    commit_sha      VARCHAR(40) NOT NULL,
    ref             VARCHAR(255) NOT NULL,
    trigger_type    VARCHAR(10) NOT NULL,    -- push/pr/manual/schedule
    triggered_by    BIGINT REFERENCES users(id),
    status          VARCHAR(10) NOT NULL DEFAULT 'pending',
    yaml_content    TEXT,
    variables       JSONB DEFAULT '{}',
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    duration        INT,                     -- 秒
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pipelines_repo ON pipelines(repo_id, created_at DESC);
CREATE INDEX idx_pipelines_commit ON pipelines(repo_id, commit_sha);

-- 流水线任务表
CREATE TABLE pipeline_jobs (
    id              BIGSERIAL PRIMARY KEY,
    pipeline_id     BIGINT NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    stage           VARCHAR(100) NOT NULL,
    status          VARCHAR(10) NOT NULL DEFAULT 'pending',
    depends_on      TEXT[],
    runner_id       BIGINT REFERENCES runners(id),
    runner_labels   TEXT[],
    steps           JSONB DEFAULT '[]',
    log             TEXT,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    duration        INT,
    exit_code       INT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_jobs_pipeline ON pipeline_jobs(pipeline_id);

-- Runner 表
CREATE TABLE runners (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,
    description     TEXT,
    repo_id         BIGINT REFERENCES repositories(id) ON DELETE CASCADE,
    org_id          BIGINT REFERENCES organizations(id) ON DELETE CASCADE,
    global          BOOLEAN DEFAULT FALSE,   -- 全局 Runner
    token           VARCHAR(64) NOT NULL UNIQUE,
    labels          TEXT[],
    status          VARCHAR(10) NOT NULL DEFAULT 'offline',
    last_ping_at    TIMESTAMPTZ,
    version         VARCHAR(50),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 流水线制品表
CREATE TABLE pipeline_artifacts (
    id              BIGSERIAL PRIMARY KEY,
    job_id          BIGINT NOT NULL REFERENCES pipeline_jobs(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    path            VARCHAR(500) NOT NULL,
    size            BIGINT,
    checksum        VARCHAR(64),
    storage_path    VARCHAR(500),            -- MinIO/S3 路径
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### 5.2.6 安全扫描

```sql
-- 安全扫描结果表
CREATE TABLE security_scans (
    id              BIGSERIAL PRIMARY KEY,
    repo_id         BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    commit_sha      VARCHAR(40) NOT NULL,
    branch          VARCHAR(255),
    scan_type       VARCHAR(20) NOT NULL,    -- sast/dependency/secret/license/dast/container
    status          VARCHAR(10) NOT NULL,    -- pending/running/passed/failed
    severity        VARCHAR(10),             -- 最高严重性
    findings_count  INT DEFAULT 0,
    summary         TEXT,
    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at     TIMESTAMPTZ,
    duration_ms     INT
);

CREATE INDEX idx_scans_repo ON security_scans(repo_id, created_at DESC);

-- 安全扫描发现表
CREATE TABLE security_scan_findings (
    id              BIGSERIAL PRIMARY KEY,
    scan_id         BIGINT NOT NULL REFERENCES security_scans(id) ON DELETE CASCADE,
    rule_id         VARCHAR(100) NOT NULL,
    severity        VARCHAR(10) NOT NULL,
    category        VARCHAR(50),
    path            VARCHAR(500),
    line            INT,
    message         TEXT NOT NULL,
    remediation     TEXT,
    reference       VARCHAR(500),
    confidence      FLOAT DEFAULT 0.0,
    cve             VARCHAR(20),             -- CVE 编号
    package_name    VARCHAR(255),            -- 包名
    package_version VARCHAR(100),            -- 版本
    license_type    VARCHAR(100),            -- 许可证
    status          VARCHAR(10) DEFAULT 'open',  -- open/fixed/dismissed
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_findings_scan ON security_scan_findings(scan_id);
CREATE INDEX idx_findings_severity ON security_scan_findings(severity);
CREATE INDEX idx_findings_cve ON security_scan_findings(cve) WHERE cve IS NOT NULL;
```

#### 5.2.7 Issue 追踪

```sql
-- Issue 表
CREATE TABLE issues (
    id              BIGSERIAL PRIMARY KEY,
    number          BIGINT NOT NULL,
    repo_id         BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT DEFAULT '',
    author_id       BIGINT NOT NULL REFERENCES users(id),
    assignee_id     BIGINT REFERENCES users(id),
    state           VARCHAR(10) NOT NULL DEFAULT 'open',
    priority        VARCHAR(10),             -- critical/high/medium/low
    labels          TEXT[],
    milestone_id    BIGINT,
    type            VARCHAR(10) DEFAULT 'issue',  -- issue/bug/feature/task
    due_date        DATE,
    time_estimate   INT,                     -- 预估工时（分钟）
    time_spent      INT,                     -- 实际工时（分钟）
    closed_by       BIGINT REFERENCES users(id),
    closed_at       TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(repo_id, number)
);

CREATE INDEX idx_issues_repo_state ON issues(repo_id, state);
CREATE INDEX idx_issues_assignee ON issues(assignee_id);

-- Issue 评论表
CREATE TABLE issue_comments (
    id              BIGSERIAL PRIMARY KEY,
    issue_id        BIGINT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    author_id       BIGINT NOT NULL REFERENCES users(id),
    body            TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 标签表
CREATE TABLE labels (
    id              BIGSERIAL PRIMARY KEY,
    repo_id         BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    color           VARCHAR(7) NOT NULL DEFAULT '#4F7DF7',
    description     TEXT DEFAULT '',
    UNIQUE(repo_id, name)
);

-- 里程碑表
CREATE TABLE milestones (
    id              BIGSERIAL PRIMARY KEY,
    repo_id         BIGINT NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT DEFAULT '',
    state           VARCHAR(10) NOT NULL DEFAULT 'open',
    due_date        DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### 5.2.8 审计日志

```sql
-- 审计日志表
CREATE TABLE audit_logs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT REFERENCES users(id),
    action          VARCHAR(50) NOT NULL,     -- repo.create / pr.merge / user.login 等
    resource_type   VARCHAR(50) NOT NULL,     -- repo / pr / issue / user / org
    resource_id     BIGINT,
    resource_name   VARCHAR(255),
    ip_address      INET,
    user_agent      TEXT,
    request_id      VARCHAR(36),              -- 请求追踪 ID
    details         JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 按时间分区（每月一个分区，便于归档和清理）
CREATE INDEX idx_audit_user ON audit_logs(user_id, created_at DESC);
CREATE INDEX idx_audit_action ON audit_logs(action, created_at DESC);
CREATE INDEX idx_audit_resource ON audit_logs(resource_type, resource_id);
```

---

## 6. API 接口规范

### 6.1 API 设计规范

#### 通用约定

| 规范 | 描述 |
|------|------|
| **协议** | HTTPS（生产环境强制） |
| **格式** | JSON（Content-Type: application/json） |
| **版本** | URL 路径版本化 `/api/v1/` |
| **认证** | Bearer Token（Personal Access Token / OAuth2 Token） |
| **分页** | `page` + `per_page` 查询参数，默认 page=1, per_page=20，最大 100 |
| **排序** | `sort` 查询参数，如 `sort=created_at&direction=desc` |
| **错误格式** | 统一 JSON 错误响应 |
| **限流** | 认证用户 1000 req/min，匿名用户 100 req/min |

#### 统一响应格式

```json
// 成功响应
{
    "id": 1,
    "name": "my-repo",
    "full_name": "myorg/my-repo",
    "created_at": "2026-04-22T10:00:00Z"
}

// 列表响应
{
    "items": [...],
    "total": 42,
    "page": 1,
    "per_page": 20
}

// 错误响应
{
    "error": {
        "code": "REPO_NOT_FOUND",
        "message": "Repository not found",
        "details": null
    }
}
```

#### 统一错误码

| HTTP 状态码 | 错误码 | 描述 |
|-------------|--------|------|
| 400 | VALIDATION_ERROR | 请求参数校验失败 |
| 401 | UNAUTHORIZED | 未认证 |
| 403 | FORBIDDEN | 无权限 |
| 404 | NOT_FOUND | 资源不存在 |
| 409 | CONFLICT | 资源冲突 |
| 422 | UNPROCESSABLE | 业务逻辑错误 |
| 429 | RATE_LIMITED | 请求频率超限 |
| 500 | INTERNAL_ERROR | 服务器内部错误 |

### 6.2 核心 API 端点详细设计

#### 6.2.1 仓库 API

```
POST   /api/v1/repos                          # 创建仓库
GET    /api/v1/repos/{owner}/{repo}            # 获取仓库详情
PATCH  /api/v1/repos/{owner}/{repo}            # 更新仓库
DELETE /api/v1/repos/{owner}/{repo}            # 删除仓库
GET    /api/v1/repos/{owner}/{repo}/branches   # 列出分支
POST   /api/v1/repos/{owner}/{repo}/branches   # 创建分支
DELETE /api/v1/repos/{owner}/{repo}/branches/{branch}  # 删除分支
GET    /api/v1/repos/{owner}/{repo}/tags       # 列出标签
POST   /api/v1/repos/{owner}/{repo}/tags       # 创建标签
GET    /api/v1/repos/{owner}/{repo}/tree/{ref} # 获取文件树
GET    /api/v1/repos/{owner}/{repo}/blob/{ref}/{path}  # 获取文件内容
GET    /api/v1/repos/{owner}/{repo}/raw/{ref}/{path}    # 获取原始文件
GET    /api/v1/repos/{owner}/{repo}/blame/{ref}/{path}  # 获取 Blame
POST   /api/v1/repos/{owner}/{repo}/forks      # Fork 仓库
POST   /api/v1/repos/{owner}/{repo}/mirrors/sync  # 触发镜像同步
GET    /api/v1/repos/{owner}/{repo}/star       # Star 仓库
DELETE /api/v1/repos/{owner}/{repo}/star       # 取消 Star
```

**创建仓库请求/响应示例**：

```json
// POST /api/v1/repos
// Request
{
    "name": "my-project",
    "description": "A new project",
    "visibility": "private",
    "default_branch": "main",
    "initialize": true,
    "template": {
        "readme": true,
        "gitignore": "Go",
        "license": "MIT"
    }
}

// Response 201
{
    "id": 42,
    "name": "my-project",
    "full_name": "myorg/my-project",
    "description": "A new project",
    "visibility": "private",
    "default_branch": "main",
    "web_url": "https://laima.cloud/myorg/my-project",
    "ssh_url": "git@laima.cloud:myorg/my-project.git",
    "http_url": "https://laima.cloud/myorg/my-project.git",
    "created_at": "2026-04-22T10:00:00Z"
}
```

#### 6.2.2 Pull Request API

```
POST   /api/v1/repos/{owner}/{repo}/pulls              # 创建 PR
GET    /api/v1/repos/{owner}/{repo}/pulls              # 列出 PR
GET    /api/v1/repos/{owner}/{repo}/pulls/{number}     # 获取 PR 详情
PATCH  /api/v1/repos/{owner}/{repo}/pulls/{number}     # 更新 PR
POST   /api/v1/repos/{owner}/{repo}/pulls/{number}/close   # 关闭 PR
POST   /api/v1/repos/{owner}/{repo}/pulls/{number}/reopen  # 重新打开 PR
POST   /api/v1/repos/{owner}/{repo}/pulls/{number}/merge   # 合并 PR
GET    /api/v1/repos/{owner}/{repo}/pulls/{number}/mergeability  # 检查合并状态
GET    /api/v1/repos/{owner}/{repo}/pulls/{number}/diff      # 获取 Diff
GET    /api/v1/repos/{owner}/{repo}/pulls/{number}/commits   # 获取关联提交
POST   /api/v1/repos/{owner}/{repo}/pulls/{number}/reviews   # 提交审查
GET    /api/v1/repos/{owner}/{repo}/pulls/{number}/reviews   # 获取审查列表
POST   /api/v1/repos/{owner}/{repo}/pulls/{number}/comments  # 创建评论
PATCH  /api/v1/repos/{owner}/{repo}/pulls/comments/{id}/resolve  # 解决评论
POST   /api/v1/repos/{owner}/{repo}/pulls/comments/{id}/apply    # 应用修复建议
```

**创建 PR 请求/响应示例**：

```json
// POST /api/v1/repos/{owner}/{repo}/pulls
// Request
{
    "title": "feat: add user authentication",
    "description": "## 变更内容\n- 新增登录/注册接口\n- JWT Token 管理\n- 密码加密存储",
    "source_branch": "feature/auth",
    "target_branch": "main",
    "is_draft": false,
    "labels": ["feature", "security"],
    "related_issues": [12, 15]
}

// Response 201
{
    "id": 1,
    "number": 42,
    "title": "feat: add user authentication",
    "state": "open",
    "merge_state": "checking",
    "source_branch": "feature/auth",
    "target_branch": "main",
    "author": {
        "id": 1,
        "username": "alice",
        "avatar_url": "https://..."
    },
    "ai_review_status": "pending",
    "created_at": "2026-04-22T10:00:00Z"
}
```

#### 6.2.3 AI 审查 API

```
GET    /api/v1/repos/{owner}/{repo}/pulls/{number}/ai-review          # 获取 AI 审查状态和结果
POST   /api/v1/repos/{owner}/{repo}/pulls/{number}/ai-review/trigger  # 手动触发 AI 审查
GET    /api/v1/repos/{owner}/{repo}/ai-review/rules                    # 获取 AI 审查规则
POST   /api/v1/repos/{owner}/{repo}/ai-review/rules                    # 创建 AI 审查规则
PATCH  /api/v1/repos/{owner}/{repo}/ai-review/rules/{id}              # 更新规则
DELETE /api/v1/repos/{owner}/{repo}/ai-review/rules/{id}              # 删除规则
POST   /api/v1/repos/{owner}/{repo}/pulls/{number}/ai-review/findings/{id}/dismiss  # 忽略发现
POST   /api/v1/repos/{owner}/{repo}/pulls/{number}/ai-review/findings/{id}/apply    # 应用修复
```

**AI 审查结果响应示例**：

```json
// GET /api/v1/repos/{owner}/{repo}/pulls/42/ai-review
// Response 200
{
    "id": 1,
    "pr_id": 42,
    "status": "completed",
    "summary": "本次变更新增了用户认证模块。发现 2 个安全问题和 1 个性能建议。",
    "total_findings": 3,
    "severity_breakdown": {
        "critical": 0,
        "high": 2,
        "medium": 1,
        "low": 0
    },
    "model_used": "qwen-coder-32b",
    "duration_ms": 12500,
    "findings": [
        {
            "id": 1,
            "analyzer": "security",
            "severity": "high",
            "category": "security",
            "path": "auth/login.go",
            "line": 45,
            "message": "使用 MD5 进行密码哈希不安全，容易被彩虹表攻击",
            "suggestion": "建议使用 bcrypt 或 argon2 进行密码哈希",
            "fix_patch": "--- a/auth/login.go\n+++ b/auth/login.go\n@@ -42,3 +42,4 @@\n-   password := md5.Sum(rawPassword)\n+   password, err := bcrypt.GenerateFromPassword(\n+       rawPassword, bcrypt.DefaultCost)",
            "confidence": 0.95
        }
    ]
}
```

#### 6.2.4 CI/CD API

```
GET    /api/v1/repos/{owner}/{repo}/pipelines                  # 列出流水线
GET    /api/v1/repos/{owner}/{repo}/pipelines/{id}             # 获取流水线详情
POST   /api/v1/repos/{owner}/{repo}/pipelines/trigger          # 手动触发流水线
POST   /api/v1/repos/{owner}/{repo}/pipelines/{id}/cancel      # 取消流水线
POST   /api/v1/repos/{owner}/{repo}/pipelines/{id}/retry       # 重试流水线
GET    /api/v1/repos/{owner}/{repo}/pipelines/{id}/jobs        # 获取任务列表
GET    /api/v1/repos/{owner}/{repo}/pipelines/{id}/jobs/{jid}/log  # 获取任务日志
POST   /api/v1/repos/{owner}/{repo}/pipelines/{id}/artifacts   # 获取制品列表
GET    /api/v1/runners                                         # 列出 Runner
POST   /api/v1/runners                                         # 注册 Runner
DELETE /api/v1/runners/{id}                                    # 删除 Runner
```

#### 6.2.5 安全扫描 API

```
GET    /api/v1/repos/{owner}/{repo}/security/scans             # 列出扫描记录
GET    /api/v1/repos/{owner}/{repo}/security/scans/{id}        # 获取扫描详情
POST   /api/v1/repos/{owner}/{repo}/security/scan/trigger      # 触发扫描
GET    /api/v1/repos/{owner}/{repo}/security/findings          # 列出所有发现
PATCH  /api/v1/repos/{owner}/{repo}/security/findings/{id}     # 更新发现状态
GET    /api/v1/repos/{owner}/{repo}/security/dashboard         # 安全仪表盘
GET    /api/v1/repos/{owner}/{repo}/security/report            # 生成合规报告
```

#### 6.2.6 Issue API

```
POST   /api/v1/repos/{owner}/{repo}/issues                     # 创建 Issue
GET    /api/v1/repos/{owner}/{repo}/issues                     # 列出 Issue
GET    /api/v1/repos/{owner}/{repo}/issues/{number}            # 获取 Issue 详情
PATCH  /api/v1/repos/{owner}/{repo}/issues/{number}            # 更新 Issue
POST   /api/v1/repos/{owner}/{repo}/issues/{number}/close      # 关闭 Issue
POST   /api/v1/repos/{owner}/{repo}/issues/{number}/comments   # 创建评论
GET    /api/v1/repos/{owner}/{repo}/labels                     # 列出标签
POST   /api/v1/repos/{owner}/{repo}/labels                     # 创建标签
GET    /api/v1/repos/{owner}/{repo}/milestones                 # 列出里程碑
POST   /api/v1/repos/{owner}/{repo}/milestones                 # 创建里程碑
```

---

## 7. 核心业务流程设计

### 7.1 Pull Request 生命周期（状态机）

```
                    创建 PR
                      │
                      ▼
               ┌─────────────┐
               │   DRAFT     │◄──────────────┐
               │  (草稿模式)  │               │
               └──────┬──────┘               │
                      │ Ready for Review      │ Convert to Draft
                      ▼                      │
               ┌─────────────┐               │
               │    OPEN     │───────────────┘
               │  (待审查)    │
               └──────┬──────┘
                      │
          ┌───────────┼───────────┐
          │           │           │
     Close PR    Merge PR    Reopen PR
          │           │           │
          ▼           ▼           │
   ┌───────────┐ ┌───────────┐   │
   │  CLOSED   │ │  MERGED   │   │
   │ (已关闭)   │ │ (已合并)   │   │
   └───────────┘ └───────────┘   │
                      │           │
                      │           │
                      └───────────┘

PR 状态子机 (merge_state):
               ┌───────────┐
               │ CHECKING  │  正在检查合并状态
               └─────┬─────┘
                     │
          ┌──────────┼──────────┬──────────┐
          ▼          ▼          ▼          ▼
    ┌──────────┐┌──────────┐┌──────────┐┌──────────┐
    │MERGEABLE ││UNSTABLE  ││CONFLICT  ││ BLOCKED  │
    │(可合并)   ││(CI未通过) ││(有冲突)   ││(审查未满足)│
    └──────────┘└──────────┘└──────────┘└──────────┘
```

**状态转换规则**：

| 当前状态 | 事件 | 目标状态 | 前置条件 |
|----------|------|----------|----------|
| - | 创建 PR（is_draft=true） | DRAFT | 用户有仓库写权限 |
| - | 创建 PR（is_draft=false） | OPEN | 用户有仓库写权限 |
| DRAFT | Mark Ready | OPEN | - |
| OPEN | Close | CLOSED | 用户有仓库写权限 |
| OPEN | Merge | MERGED | 满足合并条件（见下表） |
| CLOSED | Reopen | OPEN | 用户有仓库写权限 |
| OPEN | Convert to Draft | DRAFT | PR 作者 |

**合并条件检查（CheckMergeability）**：

| 检查项 | 宽松模式 | 标准模式 | 严格模式 |
|--------|---------|---------|---------|
| 无冲突 | ✅ | ✅ | ✅ |
| CI 通过 | 可选 | ✅ | ✅ |
| 至少 1 人审批 | 可选 | ✅ | - |
| 至少 2 人 +2 审批 | - | - | ✅ |
| CODEOWNERS 审批 | 可选 | 可选 | ✅ |
| AI 审查无 critical | 可选 | 可选 | ✅ |
| 签名提交 | 可选 | 可选 | ✅ |

### 7.2 Pull Request 创建与审查流程（时序图）

```
用户A          API网关       PR Service     Git Store      事件总线       AI审查引擎      CI引擎      通知服务
 │              │              │              │              │              │              │           │
 │─ 创建PR ────▶│              │              │              │              │              │           │
 │              │─ CreatePR()─▶│              │              │              │              │           │
 │              │              │─ 校验权限     │              │              │              │           │
 │              │              │─ 检查分支保护  │              │              │              │           │
 │              │              │─ 计算Diff ──▶│              │              │              │           │
 │              │              │◀─ Diff结果 ──│              │              │              │           │
 │              │              │─ 保存PR ────▶│(DB)          │              │              │           │
 │              │              │              │              │              │              │           │
 │              │              │─ 发布事件 ──────────────────▶│              │              │           │
 │              │              │              │              │─ pr.created ──▶│              │           │
 │              │              │              │              │              │              │           │
 │              │              │              │              │              │─ 触发AI审查  │           │
 │              │              │              │              │              │  (异步)       │           │
 │              │              │              │              │              │              │           │
 │              │              │              │              │              │              │─ 触发CI ──▶│
 │              │              │              │              │              │              │  (异步)    │
 │              │              │              │              │              │              │           │
 │              │              │              │              │              │─ 提取变更     │           │
 │              │              │              │              │              │─ 并行分析     │           │
 │              │              │              │              │              │  (Bug/安全/   │           │
 │              │              │              │              │              │   性能/风格)   │           │
 │              │              │              │              │              │              │           │
 │              │              │              │              │              │─ 聚合结果     │           │
 │              │              │              │              │              │─ 发布PR评论 ─▶│(DB)       │
 │              │              │              │              │              │              │           │
 │              │              │              │              │              │              │─ 通知 ────▶│
 │              │              │              │              │              │              │           │
 │◀─ PR创建成功─│              │              │              │              │              │           │
 │              │              │              │              │              │              │           │
 │              │              │              │              │              │              │           │
 │─ 提交审查 ──▶│              │              │              │              │              │           │
 │              │─ SubmitReview()────────────▶│              │              │              │           │
 │              │              │─ 校验权限     │              │              │              │           │
 │              │              │─ 检查评分规则  │              │              │              │           │
 │              │              │─ 保存Review ─▶│(DB)          │              │              │           │
 │              │              │─ 更新merge_state             │              │              │           │
 │              │              │─ 发布事件 ──────────────────▶│              │              │           │
 │              │              │              │              │              │              │─ 通知 ────▶│
 │              │              │              │              │              │              │           │
 │◀─ 审查成功 ──│              │              │              │              │              │           │
 │              │              │              │              │              │              │           │
 │              │              │              │              │              │              │           │
 │─ 合并PR ────▶│              │              │              │              │              │           │
 │              │─ MergePR()──▶│              │              │              │              │           │
 │              │              │─ CheckMergeability()         │              │              │           │
 │              │              │  ├─ 检查冲突 ──────────────▶│              │              │           │
 │              │              │  ├─ 检查审批 ──────────────▶│(DB)          │              │           │
 │              │              │  ├─ 检查CI状态 ────────────▶│              │              │           │
 │              │              │  └─ 检查AI审查 ───────────▶│(DB)          │              │           │
 │              │              │              │              │              │              │           │
 │              │              │─ 执行Git合并 ──────────────▶│              │              │           │
 │              │              │◀─ 合并完成 ─────────────────│              │              │           │
 │              │              │─ 更新PR状态 ──────────────▶│(DB)          │              │              │
 │              │              │─ 关闭关联Issue ────────────▶│(DB)          │              │           │
 │              │              │─ 删除源分支（可选）────────▶│              │              │           │
 │              │              │─ 发布事件 ─────────────────▶│              │              │           │
 │              │              │              │              │              │              │─ 通知 ────▶│
 │              │              │              │              │              │              │           │
 │◀─ 合并成功 ──│              │              │              │              │              │           │
```

### 7.3 AI 审查引擎流程（时序图）

```
事件总线       审查编排器       变更分析器      上下文构建器     分析器池       结果聚合器      审查发布器
   │              │              │              │              │              │              │
   │─ pr.created ─▶│              │              │              │              │              │
   │              │              │              │              │              │              │
   │              │─ ExtractChanges()──────────▶│              │              │              │
   │              │              │              │─ 获取Diff     │              │              │
   │              │              │              │─ 分类文件     │              │              │
   │              │              │              │─ 检测语言     │              │              │
   │              │◀─ Changes ──────────────────│              │              │              │
   │              │              │              │              │              │              │
   │              │─ BuildContext()────────────────────────▶│              │              │
   │              │              │              │              │              │              │
   │              │              │              │              │              │              │
   │              │─ Analyze() ──────────────────────────────▶│              │              │
   │              │              │              │              │              │              │
   │              │              │              │  ┌───────────┼───────────┐  │              │
   │              │              │              │  ▼           ▼           ▼  │              │
   │              │              │              │ Bug检测    安全检测    性能分析│              │
   │              │              │              │ (LLM)     (规则+LLM)  (LLM) │              │
   │              │              │              │  │           │           │  │              │
   │              │              │              │  ▼           ▼           ▼  │              │
   │              │              │              │ Findings   Findings   Findings│             │
   │              │              │              │  │           │           │  │              │
   │              │              │              │  └───────────┼───────────┘  │              │
   │              │              │              │              │              │              │
   │              │◀─ AllFindings ────────────────────────────│              │              │
   │              │              │              │              │              │              │
   │              │─ Aggregate()────────────────────────────────────────────▶│              │
   │              │              │              │              │              │─ 去重        │
   │              │              │              │              │              │─ 排序        │
   │              │              │              │              │              │─ 生成摘要    │
   │              │              │              │              │              │─ 生成修复    │
   │              │◀─ AggregatedResult ──────────────────────────────────────│              │
   │              │              │              │              │              │              │
   │              │─ Publish()─────────────────────────────────────────────────────────────▶│
   │              │              │              │              │              │              │─ 创建摘要评论
   │              │              │              │              │              │              │─ 创建内联评论
   │              │              │              │              │              │              │─ 更新PR状态
   │              │◀─ Published ───────────────────────────────────────────────────────────│
   │              │              │              │              │              │              │
   │◀─ ai_review.completed │         │              │              │              │              │
```

### 7.4 CI/CD 流水线执行流程（时序图）

```
事件/用户      API网关      Pipeline Service    YAML解析器     调度器      Runner      任务执行器
   │             │              │                 │            │           │            │
   │─ push/PR ──▶│              │                 │            │           │            │
   │             │─ 触发流水线 ─▶│                 │            │           │            │
   │             │              │─ 获取 .laima-ci.yml           │           │            │
   │             │              │────────────────▶│            │           │            │
   │             │              │                 │─ 解析YAML   │           │            │
   │             │              │                 │─ 构建DAG    │           │            │
   │             │              │◀─ Pipeline DAG ─│            │           │            │
   │             │              │                 │            │           │            │
   │             │              │─ 创建Pipeline ─▶│(DB)        │           │            │
   │             │              │─ 创建Jobs ────▶│(DB)        │           │            │
   │             │              │                 │            │           │            │
   │             │              │─ Schedule() ──────────────▶│           │            │
   │             │              │                 │            │           │            │
   │             │              │                 │            │─ 查找可用Runner          │
   │             │              │                 │            │─ 分配Job ─────────────▶│
   │             │              │                 │            │           │            │
   │             │              │                 │            │           │─ 拉取代码   │
   │             │              │                 │            │           │─ 执行Steps │
   │             │              │                 │            │           │            │
   │             │              │                 │            │           │  ┌─────────┐│
   │             │              │                 │            │           │  │ Step 1  ││
   │             │              │                 │            │           │  │ Step 2  ││
   │             │              │                 │            │           │  │ Step 3  ││
   │             │              │                 │            │           │  └─────────┘│
   │             │              │                 │            │           │            │
   │             │              │                 │            │◀─ 上报结果 ─│            │
   │             │              │                 │            │           │            │
   │             │              │                 │            │─ 更新Job状态            │
   │             │              │                 │            │─ 检查后续Job           │
   │             │              │                 │            │  (DAG依赖)              │
   │             │              │                 │            │           │            │
   │             │              │                 │            │─ 分配下一批Job ────────▶│
   │             │              │                 │            │           │            │
   │             │              │                 │            │◀─ 全部完成 ─│            │
   │             │              │                 │            │           │            │
   │             │              │◀─ 更新Pipeline状态           │           │            │
   │             │              │─ 发布 pipeline.completed 事件│           │            │
   │             │              │                 │            │           │            │
   │◀─ 流水线结果 │              │                 │            │           │            │
```

### 7.5 安全扫描触发流程

```
触发源         安全服务         扫描调度器        SAST引擎      依赖扫描      密钥检测      结果聚合
   │             │               │               │            │            │            │
   │─ push事件 ─▶│               │               │            │            │            │
   │             │─ 检查仓库配置  │               │            │            │            │
   │             │  (是否启用扫描)│               │            │            │            │
   │             │               │               │            │            │            │
   │             │─ 调度扫描 ─────────────────▶│            │            │            │
   │             │               │               │            │            │            │
   │             │               │  ┌────────────┼────────────┼────────────┐│            │
   │             │               │  ▼            ▼            ▼            ││            │
   │             │               │ SAST       依赖扫描      密钥检测      ││            │
   │             │               │ (Semgrep)  (OSV DB)     (Gitleaks)   ││            │
   │             │               │  │            │            │            ││            │
   │             │               │  ▼            ▼            ▼            ││            │
   │             │               │ Findings   Findings    Findings      ││            │
   │             │               │  │            │            │            ││            │
   │             │               │  └────────────┼────────────┘            ││            │
   │             │               │               │                       ││            │
   │             │◀─ 聚合结果 ──────────────────────────────────────────▶│            │
   │             │               │               │                       ││            │
   │             │─ 保存结果 ────▶│(DB)           │                       ││            │
   │             │─ 创建Issue（高危）────────────▶│(DB)                   ││            │
   │             │- 发送通知     │               │                       ││            │
   │             │               │               │                       ││            │
   │◀─ 扫描完成  │               │               │                       ││            │
```

### 7.6 事件驱动架构

系统采用**领域事件**实现模块间松耦合通信：

```
┌─────────────────────────────────────────────────────────────────────┐
│                       领域事件清单                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  仓库事件 (repo.*):                                                 │
│    repo.created          仓库创建                                   │
│    repo.updated          仓库更新                                   │
│    repo.deleted          仓库删除                                   │
│    repo.forked           仓库 Fork                                  │
│    repo.mirror_synced    镜像同步完成                                │
│    repo.starred          仓库 Star                                  │
│                                                                     │
│  PR 事件 (pr.*):                                                    │
│    pr.created            PR 创建                                    │
│    pr.updated            PR 更新                                    │
│    pr.closed             PR 关闭                                    │
│    pr.reopened           PR 重新打开                                 │
│    pr.merged             PR 合并                                    │
│    pr.review_submitted   审查提交                                   │
│    pr.comment_created    评论创建                                    │
│    pr.suggestion_applied 修复建议已应用                              │
│                                                                     │
│  AI 审查事件 (ai_review.*):                                         │
│    ai_review.triggered  AI 审查触发                                 │
│    ai_review.completed  AI 审查完成                                 │
│    ai_review.failed     AI 审查失败                                 │
│                                                                     │
│  流水线事件 (pipeline.*):                                           │
│    pipeline.created     流水线创建                                  │
│    pipeline.started     流水线开始                                  │
│    pipeline.completed   流水线完成                                  │
│    pipeline.failed      流水线失败                                  │
│    job.started          任务开始                                    │
│    job.completed        任务完成                                    │
│                                                                     │
│  安全事件 (security.*):                                             │
│    security.scan_triggered  安全扫描触发                             │
│    security.scan_completed 安全扫描完成                             │
│    security.vulnerability_found 发现新漏洞                          │
│                                                                     │
│  用户事件 (user.*):                                                 │
│    user.registered       用户注册                                   │
│    user.login            用户登录                                   │
│    user.updated          用户信息更新                                │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**事件消费者映射**：

| 事件 | 消费者 | 动作 |
|------|--------|------|
| `pr.created` | AI 审查引擎 | 触发 AI 审查 |
| `pr.created` | CI/CD 引擎 | 触发流水线 |
| `pr.created` | 安全扫描 | 触发安全扫描 |
| `pr.created` | 通知服务 | 发送通知给 CODEOWNERS |
| `pr.merged` | Issue 服务 | 自动关闭关联 Issue |
| `pr.merged` | 通知服务 | 发送合并通知 |
| `pr.review_submitted` | PR 服务 | 更新 merge_state |
| `ai_review.completed` | PR 服务 | 更新 AI 审查状态 |
| `ai_review.completed` | 通知服务 | 发送审查结果通知 |
| `pipeline.completed` | PR 服务 | 更新 CI 状态检查 |
| `security.scan_completed` | 通知服务 | 发送安全扫描结果通知 |
| `security.vulnerability_found` | Issue 服务 | 自动创建 Issue |

---

## 8. 非功能性架构设计

### 8.1 性能设计

| 场景 | 目标 | 策略 |
|------|------|------|
| API 响应时间 | P99 ≤ 200ms | Redis 缓存热点数据、数据库索引优化、查询分页 |
| Git Clone 速度 | 与原生 Git 相当 | go-git 优化、HTTP 智能协议、浅克隆支持 |
| 代码搜索 | 10万行代码 < 500ms | Meilisearch 全文索引、增量索引更新 |
| AI 审查响应 | 中等 PR < 30s | 增量分析、并行分析器、本地 LLM 优先 |
| 大文件处理 | 支持 5GB LFS 文件 | 流式传输、分块上传、MinIO 对象存储 |
| 并发 PR 审查 | 支持 100+ 并发 | Worker 池、LLM 请求队列、限流控制 |

### 8.2 缓存策略

| 缓存对象 | 缓存层 | TTL | 失效策略 |
|----------|--------|-----|----------|
| 仓库详情 | Redis | 5 min | 仓库更新时主动失效 |
| 用户信息 | Redis | 30 min | 用户更新时主动失效 |
| 分支列表 | Redis | 2 min | push 事件触发失效 |
| PR 详情 | Redis | 1 min | PR 更新事件触发失效 |
| 代码搜索结果 | Meilisearch | - | push 事件触发增量索引 |
| API 响应 | HTTP Cache | 10 min | ETag + Last-Modified |
| 静态资源 | CDN | 1 year | 内容哈希文件名 |

### 8.3 安全架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                         安全架构分层                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  网络安全层                                                   │   │
│  │  TLS 1.3 · DDoS 防护 · WAF · IP 白名单 · 速率限制            │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  认证层                                                       │   │
│  │  JWT Token · Session 管理 · OAuth2 · SAML · LDAP · 2FA       │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  鉴权层                                                       │   │
│  │  RBAC (6级角色) · ABAC (属性策略) · 分支级权限 · CODEOWNERS   │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  输入安全层                                                   │   │
│  │  参数校验 · SQL 注入防护 · XSS 防护 · CSRF Token · 文件校验  │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  数据安全层                                                   │   │
│  │  密码 bcrypt · 敏感字段 AES-256 · 数据库加密 · 传输加密      │   │
│  │  密钥轮换 · HSM 支持 · 安全擦除                               │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  审计层                                                       │   │
│  │  全操作审计日志 · 不可篡改 · 审计日志导出 · 合规报告          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 8.4 可观测性设计

| 维度 | 工具 | 内容 |
|------|------|------|
| **日志** | 结构化日志 (zap) | 请求日志、错误日志、审计日志 |
| **指标** | Prometheus | API QPS/延迟、Git 操作计数、CI/CD 成功率、AI 审查耗时 |
| **追踪** | OpenTelemetry | 分布式追踪、请求链路分析 |
| **健康检查** | /healthz, /readyz | 存活探针、就绪探针 |
| **仪表盘** | Grafana | 系统概览、业务指标、告警面板 |

### 8.5 配置管理

```yaml
# config.yaml 示例
server:
  http:
    host: 0.0.0.0
    port: 8080
    tls:
      enabled: true
      cert: /etc/laima/cert.pem
      key: /etc/laima/key.pem
  ssh:
    host: 0.0.0.0
    port: 2222
    server_key: /etc/laima/ssh_host_key

database:
  postgres:
    host: localhost
    port: 5432
    name: laima
    user: laima
    password: ${DB_PASSWORD}       # 环境变量引用
    ssl_mode: disable
    max_open_conns: 50
    max_idle_conns: 10

cache:
  redis:
    host: localhost
    port: 6379
    password: ${REDIS_PASSWORD}
    db: 0
  # Redis 不可用时降级为内存缓存
  fallback_memory: true

search:
  meilisearch:
    host: http://localhost:7700
    api_key: ${MEILI_API_KEY}

storage:
  minio:
    endpoint: localhost:9000
    access_key: ${MINIO_ACCESS_KEY}
    secret_key: ${MINIO_SECRET_KEY}
    bucket: laima
    use_ssl: false

ai:
  default_provider: ollama          # ollama / openai / qwen / deepseek
  ollama:
    base_url: http://localhost:11434
    model: qwen-coder:32b
    timeout: 120s
  openai:
    api_key: ${OPENAI_API_KEY}
    model: gpt-4o
    timeout: 60s
  max_concurrent_reviews: 5         # 最大并发审查数
  review_timeout: 300s              # 单次审查超时

security:
  scanner:
    sast:
      enabled: true
      engine: semgrep
    dependency:
      enabled: true
      osv_url: https://api.osv.dev
    secret:
      enabled: true
      engine: gitleaks
    license:
      enabled: true

runner:
  token_expire: 30d                # Runner Token 过期时间
  max_concurrent_jobs: 10           # 单 Runner 最大并发任务数
```

---

> **文档版本**：v1.0
> **最后更新**：2026 年 4 月 22 日
> **关联文档**：《代码托管平台竞品分析报告》、《Laima（莱码）产品方案》
