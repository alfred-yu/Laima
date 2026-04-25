# Laima（莱码）待完善功能实施计划

> **版本**：v1.0  
> **日期**：2026 年 4 月 24 日  
> **基于**：《Laima（莱码）软件设计与实现一致性分析报告》

---

## 一、实施计划总览

### 1.1 优先级划分

| 优先级 | 功能模块 | 预计工作量 | 实施周期 |
|--------|---------|-----------|---------|
| **P0**（高优先级） | AI 审查引擎完善 | 15 人天 | 3 周 |
| **P0**（高优先级） | CI/CD Runner 管理 | 12 人天 | 2.5 周 |
| **P1**（中优先级） | 测试覆盖率提升 | 10 人天 | 2 周 |
| **P1**（中优先级） | 高级功能实现 | 15 人天 | 3 周 |
| **P2**（低优先级） | 前端组件库完善 | 8 人天 | 1.5 周 |

### 1.2 实施路线图

```
Week 1-3          Week 4-6          Week 7-9          Week 10-12
  │                 │                 │                 │
  ▼                 ▼                 ▼                 ▼
┌────────┐      ┌────────┐      ┌────────┐      ┌────────┐
│ AI审查  │      │ CI/CD  │      │ 测试   │      │ 高级   │
│ 引擎    │─────▶│ Runner │─────▶│ 覆盖率 │─────▶│ 功能   │
│ 完善    │      │ 管理   │      │ 提升   │      │ 实现   │
└────────┘      └────────┘      └────────┘      └────────┘
```

---

## 二、详细实施计划

### 2.1 AI 审查引擎完善（P0）

#### 2.1.1 目标

- 实现多 LLM 模型支持（OpenAI、Ollama、Qwen、DeepSeek）
- 完善 Prompt 工程模板
- 实现 AI 审查结果优化

#### 2.1.2 实施步骤

**阶段一：多模型支持（5 人天）**

1. **设计 LLM Provider 抽象层**
   - 文件：`backend/internal/ai/infrastructure/llm_provider.go`
   - 定义统一的 LLM Provider 接口
   - 实现工厂模式创建不同 Provider

2. **实现 OpenAI Provider**
   - 文件：`backend/internal/ai/infrastructure/openai_provider.go`
   - 集成 OpenAI API
   - 实现流式响应处理

3. **实现 Ollama Provider**
   - 文件：`backend/internal/ai/infrastructure/ollama_provider.go`
   - 集成本地 Ollama 服务
   - 支持离线模式

4. **实现 Qwen Provider**
   - 文件：`backend/internal/ai/infrastructure/qwen_provider.go`
   - 集成通义千问 API
   - 支持国产模型

5. **实现 DeepSeek Provider**
   - 文件：`backend/internal/ai/infrastructure/deepseek_provider.go`
   - 集成 DeepSeek API

**阶段二：Prompt 工程优化（5 人天）**

1. **创建 Prompt 模板管理器**
   - 文件：`backend/internal/ai/domain/prompt_template.go`
   - 支持模板 CRUD
   - 支持变量替换

2. **实现场景化 Prompt 模板**
   - Bug 检测模板
   - 安全检测模板
   - 性能分析模板
   - 代码风格检查模板

3. **实现上下文构建器**
   - 文件：`backend/internal/ai/app/context_builder.go`
   - 提取项目结构
   - 提取历史审查记录
   - 提取团队编码规范

**阶段三：结果优化（5 人天）**

1. **实现结果聚合器**
   - 文件：`backend/internal/ai/app/result_aggregator.go`
   - 去重逻辑
   - 优先级排序
   - 置信度计算

2. **实现修复建议生成器**
   - 文件：`backend/internal/ai/app/fix_generator.go`
   - 生成一键修复补丁
   - 支持多种修复策略

3. **实现审查摘要生成**
   - 文件：`backend/internal/ai/app/summary_generator.go`
   - 生成 PR 变更摘要
   - 生成审查报告

#### 2.1.3 验收标准

- [ ] 支持至少 4 种 LLM 模型
- [ ] Prompt 模板可配置
- [ ] AI 审查准确率 ≥ 70%
- [ ] 审查响应时间 ≤ 30s（中等 PR）

---

### 2.2 CI/CD Runner 管理（P0）

#### 2.2.1 目标

- 实现 Runner 注册和管理
- 实现任务调度和分发
- 实现任务执行和日志收集

#### 2.2.2 实施步骤

**阶段一：Runner 管理（4 人天）**

1. **完善 Runner 数据模型**
   - 文件：`backend/internal/cicd/domain/model.go`
   - 添加 Runner 状态字段
   - 添加 Runner 标签

2. **实现 Runner 注册 API**
   - 文件：`backend/internal/cicd/api/runner_handler.go`
   - `POST /api/v1/runners/register`
   - 生成 Runner Token
   - 返回配置信息

3. **实现 Runner 心跳机制**
   - 文件：`backend/internal/cicd/app/runner_service.go`
   - Runner 定期上报状态
   - 自动标记离线 Runner

4. **实现 Runner 管理 API**
   - 列出 Runner
   - 删除 Runner
   - 更新 Runner 配置

**阶段二：任务调度（4 人天）**

1. **实现任务队列**
   - 文件：`backend/internal/cicd/infrastructure/job_queue.go`
   - 基于 Redis 实现
   - 支持优先级队列

2. **实现任务调度器**
   - 文件：`backend/internal/cicd/app/scheduler.go`
   - 根据 Runner 标签匹配任务
   - 实现负载均衡
   - 支持任务重试

3. **实现任务分发机制**
   - Runner 拉取任务
   - 任务状态更新
   - 任务超时处理

**阶段三：任务执行（4 人天）**

1. **实现任务执行器**
   - 文件：`backend/internal/cicd/app/executor.go`
   - 支持 Docker 执行器
   - 支持 Shell 执行器
   - 支持 Kubernetes 执行器

2. **实现日志收集**
   - 文件：`backend/internal/cicd/infrastructure/log_collector.go`
   - 实时日志流
   - 日志存储
   - 日志查询

3. **实现制品管理**
   - 文件：`backend/internal/cicd/infrastructure/artifact_manager.go`
   - 制品上传
   - 制品下载
   - 制品清理

#### 2.2.3 验收标准

- [ ] Runner 可正常注册和心跳
- [ ] 任务可正确调度和执行
- [ ] 日志可实时查看
- [ ] 支持至少 2 种执行器

---

### 2.3 测试覆盖率提升（P1）

#### 2.3.1 目标

- 单元测试覆盖率 ≥ 80%
- 集成测试覆盖核心流程
- E2E 测试覆盖主要功能

#### 2.3.2 实施步骤

**阶段一：单元测试（5 人天）**

1. **仓库模块测试**
   - 文件：`backend/internal/repo/app/service_test.go`
   - 测试所有 Service 方法
   - Mock 数据库和 Git 服务

2. **PR 模块测试**
   - 文件：`backend/internal/pr/app/service_test.go`
   - 测试 PR 生命周期
   - 测试合并逻辑

3. **AI 模块测试**
   - 文件：`backend/internal/ai/app/service_test.go`
   - Mock LLM Provider
   - 测试审查流程

4. **CI/CD 模块测试**
   - 文件：`backend/internal/cicd/app/service_test.go`
   - 测试流水线创建
   - 测试任务调度

5. **安全模块测试**
   - 文件：`backend/internal/security/app/service_test.go`
   - 测试扫描流程
   - 测试结果聚合

**阶段二：集成测试（3 人天）**

1. **API 集成测试**
   - 文件：`backend/tests/integration/api_test.go`
   - 测试完整 API 流程
   - 使用测试数据库

2. **业务流程集成测试**
   - 文件：`backend/tests/integration/workflow_test.go`
   - 测试 PR 创建到合并流程
   - 测试 CI/CD 触发流程

**阶段三：E2E 测试（2 人天）**

1. **前端 E2E 测试**
   - 使用 Playwright/Cypress
   - 测试用户登录流程
   - 测试仓库创建流程
   - 测试 PR 流程

#### 2.3.3 验收标准

- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 集成测试覆盖核心流程
- [ ] E2E 测试通过率 100%

---

### 2.4 高级功能实现（P1）

#### 2.4.1 目标

- 实现仓库镜像功能
- 实现 GPG 签名验证
- 实现 CODEOWNERS 功能
- 实现 Cherry-pick/Revert 功能

#### 2.4.2 实施步骤

**阶段一：仓库镜像（5 人天）**

1. **设计镜像数据模型**
   - 文件：`backend/internal/repo/domain/mirror.go`
   - 镜像配置
   - 同步状态

2. **实现镜像同步服务**
   - 文件：`backend/internal/repo/app/mirror_service.go`
   - 定时同步
   - 增量同步
   - 错误处理

3. **实现镜像管理 API**
   - 文件：`backend/internal/repo/api/mirror_handler.go`
   - 创建镜像仓库
   - 触发同步
   - 查看同步状态

**阶段二：GPG 签名（4 人天）**

1. **实现 GPG 密钥管理**
   - 文件：`backend/internal/user/app/gpg_service.go`
   - 添加 GPG 密钥
   - 验证 GPG 密钥

2. **实现提交签名验证**
   - 文件：`backend/internal/git/signature.go`
   - 验证提交签名
   - 标记未签名提交

3. **实现签名策略**
   - 文件：`backend/internal/repo/domain/signature_policy.go`
   - 强制签名配置
   - 签名验证规则

**阶段三：CODEOWNERS（3 人天）**

1. **实现 CODEOWNERS 解析**
   - 文件：`backend/internal/repo/app/codeowners.go`
   - 解析 CODEOWNERS 文件
   - 匹配文件路径

2. **实现自动审查分配**
   - 文件：`backend/internal/pr/app/review_assigner.go`
   - 根据 CODEOWNERS 分配审查者
   - 检查审查要求

**阶段四：Cherry-pick/Revert（3 人天）**

1. **实现 Cherry-pick 功能**
   - 文件：`backend/internal/pr/app/cherry_pick.go`
   - 选择提交
   - 创建新分支
   - 处理冲突

2. **实现 Revert 功能**
   - 文件：`backend/internal/pr/app/revert.go`
   - 创建 Revert PR
   - 保留原始信息

#### 2.4.3 验收标准

- [ ] 仓库镜像可正常同步
- [ ] GPG 签名可正常验证
- [ ] CODEOWNERS 可正确分配审查者
- [ ] Cherry-pick/Revert 功能正常

---

### 2.5 前端组件库完善（P2）

#### 2.5.1 目标

- 完善基础 UI 组件
- 实现业务组件
- 优化用户体验

#### 2.5.2 实施步骤

**阶段一：基础组件（3 人天）**

1. **完善现有组件**
   - Button、Input、Select、Textarea
   - 添加更多变体
   - 优化样式

2. **新增基础组件**
   - Modal（模态框）
   - Dropdown（下拉菜单）
   - Tabs（标签页）
   - Table（表格）
   - Pagination（分页）

**阶段二：业务组件（3 人天）**

1. **代码相关组件**
   - CodeEditor（代码编辑器）
   - DiffViewer（差异查看器）
   - FileTree（文件树）

2. **Git 相关组件**
   - CommitList（提交列表）
   - BranchSelector（分支选择器）
   - TagList（标签列表）

3. **PR 相关组件**
   - ReviewPanel（审查面板）
   - CommentThread（评论线程）
   - MergeStatus（合并状态）

**阶段三：用户体验优化（2 人天）**

1. **加载状态**
   - Skeleton（骨架屏）
   - Loading Spinner
   - Progress Bar

2. **反馈提示**
   - Toast（消息提示）
   - Tooltip（工具提示）
   - Empty State（空状态）

3. **响应式设计**
   - 移动端适配
   - 平板适配

#### 2.5.3 验收标准

- [ ] 基础组件 ≥ 10 个
- [ ] 业务组件 ≥ 8 个
- [ ] 响应式设计完成

---

## 三、资源需求

### 3.1 人力需求

| 角色 | 人数 | 工作内容 |
|------|------|---------|
| 后端开发 | 2-3 人 | AI 引擎、CI/CD、高级功能 |
| 前端开发 | 1-2 人 | 组件库、用户体验优化 |
| 测试工程师 | 1 人 | 测试用例编写、自动化测试 |
| DevOps | 1 人 | CI/CD Runner 部署、环境配置 |

### 3.2 技术资源

| 资源 | 用途 |
|------|------|
| OpenAI API Key | AI 审查测试 |
| Ollama 服务器 | 本地 LLM 测试 |
| 测试服务器 | Runner 部署测试 |
| CI/CD 环境 | 自动化测试 |

---

## 四、风险管理

### 4.1 技术风险

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|---------|
| LLM API 不稳定 | 中 | 高 | 实现降级机制，支持本地模型 |
| Runner 调度复杂 | 中 | 中 | 参考成熟方案，逐步迭代 |
| 测试覆盖率目标过高 | 低 | 中 | 分阶段提升，优先核心模块 |

### 4.2 进度风险

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|---------|
| 人员变动 | 中 | 高 | 文档完善，知识共享 |
| 需求变更 | 中 | 中 | 敏捷开发，快速响应 |
| 技术难点攻关 | 高 | 中 | 预留缓冲时间，寻求外部支持 |

---

## 五、验收标准

### 5.1 功能验收

- [ ] 所有功能模块按计划完成
- [ ] 功能测试通过率 100%
- [ ] 性能指标达标

### 5.2 质量验收

- [ ] 代码审查通过
- [ ] 测试覆盖率达标
- [ ] 无严重 Bug

### 5.3 文档验收

- [ ] API 文档更新
- [ ] 用户手册更新
- [ ] 部署文档更新

---

## 六、后续规划

### 6.1 v0.5 版本（当前计划）

- AI 审查引擎完善
- CI/CD Runner 管理
- 测试覆盖率提升
- 高级功能实现

### 6.2 v0.6 版本（下一阶段）

- 容器镜像仓库
- 包管理功能
- Wiki 功能完善
- Pages 功能

### 6.3 v1.0 版本（正式发布）

- 性能优化
- 安全加固
- 文档完善
- 社区建设

---

> **文档版本**：v1.0  
> **最后更新**：2026 年 4 月 24 日
