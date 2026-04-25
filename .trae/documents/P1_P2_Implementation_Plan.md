# P1 和 P2 任务实施计划

> **版本**：v1.0  
> **日期**：2026 年 4 月 24 日  
> **基于**：Laima 待完善功能实施计划

---

## 一、任务概览

### 1.1 P1 任务（中优先级）

| 任务 | 预计工作量 | 实施周期 |
|------|-----------|---------|
| 测试覆盖率提升 | 10 人天 | 2 周 |
| 高级功能实现 | 15 人天 | 3 周 |

### 1.2 P2 任务（低优先级）

| 任务 | 预计工作量 | 实施周期 |
|------|-----------|---------|
| 前端组件库完善 | 8 人天 | 1.5 周 |

### 1.3 实施路线图

```
Week 1-2           Week 3-5           Week 6-7
  │                 │                 │
  ▼                 ▼                 ▼
┌────────┐      ┌────────┐      ┌────────┐
│ 测试   │      │ 高级   │      │ 前端   │
│ 覆盖率 │─────▶│ 功能   │─────▶│ 组件库 │
│ 提升   │      │ 实现   │      │ 完善   │
└────────┘      └────────┘      └────────┘
```

---

## 二、P1 任务详细计划

### 2.1 测试覆盖率提升

#### 2.1.1 目标

- 单元测试覆盖率 ≥ 80%
- 集成测试覆盖核心流程
- E2E 测试覆盖主要功能

#### 2.1.2 实施步骤

**阶段一：单元测试（5 人天）**

**步骤 1：仓库模块测试**
- 文件：`backend/internal/repo/app/service_test.go`
- 测试内容：
  - `CreateRepository` - 创建仓库
  - `GetRepository` - 获取仓库
  - `UpdateRepository` - 更新仓库
  - `DeleteRepository` - 删除仓库
  - `ForkRepository` - Fork 仓库
  - `ListBranches` - 列出分支
  - `SearchCode` - 代码搜索
- Mock 依赖：数据库、Git 服务、缓存

**步骤 2：PR 模块测试**
- 文件：`backend/internal/pr/app/service_test.go`
- 测试内容：
  - `CreatePullRequest` - 创建 PR
  - `UpdatePullRequest` - 更新 PR
  - `MergePullRequest` - 合并 PR
  - `ClosePullRequest` - 关闭 PR
  - `AddReview` - 添加审查
  - `AddComment` - 添加评论
- Mock 依赖：数据库、Git 服务

**步骤 3：AI 模块测试**
- 文件：`backend/internal/ai/app/service_test.go`
- 测试内容：
  - `TriggerReview` - 触发审查
  - `GetReviewResult` - 获取审查结果
  - `PromptTemplateManager` - Prompt 模板管理
  - `ResultAggregator` - 结果聚合
- Mock 依赖：LLM Provider、数据库

**步骤 4：CI/CD 模块测试**
- 文件：`backend/internal/cicd/app/service_test.go`
- 测试内容：
  - `CreatePipeline` - 创建流水线
  - `UpdatePipelineStatus` - 更新流水线状态
  - `RunnerService` - Runner 管理
  - `JobScheduler` - 任务调度
- Mock 依赖：数据库、Redis

**步骤 5：安全模块测试**
- 文件：`backend/internal/security/app/service_test.go`
- 测试内容：
  - `TriggerScan` - 触发扫描
  - `GetScanResult` - 获取扫描结果
  - `GetDASTConfig` - DAST 配置
  - `GetContainerScanConfig` - 容器扫描配置
- Mock 依赖：数据库、扫描工具

**阶段二：集成测试（3 人天）**

**步骤 6：API 集成测试**
- 文件：`backend/tests/integration/api_test.go`
- 测试内容：
  - 用户注册/登录流程
  - 仓库创建 → 提交代码 → PR 创建 → 合并流程
  - CI/CD 流水线触发和执行
  - AI 审查触发和结果获取
- 使用测试数据库和测试环境

**步骤 7：业务流程集成测试**
- 文件：`backend/tests/integration/workflow_test.go`
- 测试内容：
  - 完整的 PR 生命周期（创建 → 审查 → 修改 → 合并）
  - CI/CD 流水线完整流程（触发 → 执行 → 结果）
  - 安全扫描完整流程（触发 → 扫描 → 报告）

**阶段三：E2E 测试（2 人天）**

**步骤 8：前端 E2E 测试**
- 工具：Playwright
- 文件：`frontend/tests/e2e/`
- 测试内容：
  - 用户登录流程
  - 仓库创建和浏览
  - PR 创建和合并
  - Issue 创建和管理
  - CI/CD 流水线查看

#### 2.1.3 验收标准

- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 所有单元测试通过
- [ ] 集成测试覆盖核心流程
- [ ] E2E 测试通过率 100%

---

### 2.2 高级功能实现

#### 2.2.1 目标

- 实现仓库镜像功能
- 实现 GPG 签名验证
- 实现 CODEOWNERS 功能
- 实现 Cherry-pick/Revert 功能

#### 2.2.2 实施步骤

**阶段一：仓库镜像（5 人天）**

**步骤 1：设计镜像数据模型**
- 文件：`backend/internal/repo/domain/mirror.go`
- 内容：
  - `Mirror` 模型（ID、源 URL、同步状态、最后同步时间）
  - `MirrorConfig` 配置结构体
  - 镜像同步策略枚举

**步骤 2：实现镜像同步服务**
- 文件：`backend/internal/repo/app/mirror_service.go`
- 功能：
  - `CreateMirror` - 创建镜像仓库
  - `SyncMirror` - 同步镜像
  - `GetMirrorStatus` - 获取镜像状态
  - `UpdateMirrorConfig` - 更新镜像配置
- 使用定时任务进行自动同步
- 支持增量同步

**步骤 3：实现镜像管理 API**
- 文件：`backend/internal/repo/api/mirror_handler.go`
- API 端点：
  - `POST /api/v1/repos/{owner}/{repo}/mirror` - 创建镜像
  - `POST /api/v1/repos/{owner}/{repo}/mirror/sync` - 触发同步
  - `GET /api/v1/repos/{owner}/{repo}/mirror/status` - 获取状态
  - `PUT /api/v1/repos/{owner}/{repo}/mirror` - 更新配置

**阶段二：GPG 签名（4 人天）**

**步骤 4：实现 GPG 密钥管理**
- 文件：`backend/internal/user/app/gpg_service.go`
- 功能：
  - `AddGPGKey` - 添加 GPG 密钥
  - `VerifyGPGKey` - 验证 GPG 密钥
  - `ListGPGKeys` - 列出用户的 GPG 密钥
  - `DeleteGPGKey` - 删除 GPG 密钥
- 使用 `golang.org/x/crypto/openpgp` 包

**步骤 5：实现提交签名验证**
- 文件：`backend/internal/git/signature.go`
- 功能：
  - `VerifyCommitSignature` - 验证提交签名
  - `GetCommitSignature` - 获取提交签名信息
  - 标记未签名提交

**步骤 6：实现签名策略**
- 文件：`backend/internal/repo/domain/signature_policy.go`
- 功能：
  - 定义签名策略模型
  - 强制签名配置
  - 签名验证规则
  - PR 合并时检查签名

**阶段三：CODEOWNERS（3 人天）**

**步骤 7：实现 CODEOWNERS 解析**
- 文件：`backend/internal/repo/app/codeowners.go`
- 功能：
  - `ParseCODEOWNERS` - 解析 CODEOWNERS 文件
  - `MatchPath` - 匹配文件路径
  - `GetOwners` - 获取文件所有者
- 支持通配符和正则表达式

**步骤 8：实现自动审查分配**
- 文件：`backend/internal/pr/app/review_assigner.go`
- 功能：
  - `AssignReviewers` - 根据 CODEOWNERS 分配审查者
  - `CheckReviewRequirements` - 检查审查要求
  - `GetRequiredReviewers` - 获取必需审查者

**阶段四：Cherry-pick/Revert（3 人天）**

**步骤 9：实现 Cherry-pick 功能**
- 文件：`backend/internal/pr/app/cherry_pick.go`
- 功能：
  - `CherryPick` - 选择提交
  - `CreateCherryPickBranch` - 创建新分支
  - `ResolveConflicts` - 处理冲突
  - `CreateCherryPickPR` - 创建 Cherry-pick PR

**步骤 10：实现 Revert 功能**
- 文件：`backend/internal/pr/app/revert.go`
- 功能：
  - `RevertCommit` - Revert 提交
  - `CreateRevertPR` - 创建 Revert PR
  - 保留原始提交信息

#### 2.2.3 验收标准

- [ ] 仓库镜像可正常同步
- [ ] GPG 签名可正常验证
- [ ] CODEOWNERS 可正确分配审查者
- [ ] Cherry-pick/Revert 功能正常

---

## 三、P2 任务详细计划

### 3.1 前端组件库完善

#### 3.1.1 目标

- 完善基础 UI 组件（≥10 个）
- 实现业务组件（≥8 个）
- 优化用户体验

#### 3.1.2 实施步骤

**阶段一：基础组件（3 人天）**

**步骤 1：完善现有组件**
- 文件：`frontend/src/components/ui/`
- 组件列表：
  - `Button.vue` - 按钮组件（添加更多变体）
  - `Input.vue` - 输入框组件
  - `Select.vue` - 选择器组件
  - `Textarea.vue` - 文本域组件
- 添加变体：primary、secondary、danger、ghost
- 添加尺寸：small、medium、large

**步骤 2：新增基础组件**
- 文件：`frontend/src/components/ui/`
- 组件列表：
  - `Modal.vue` - 模态框组件
  - `Dropdown.vue` - 下拉菜单组件
  - `Tabs.vue` - 标签页组件
  - `Table.vue` - 表格组件
  - `Pagination.vue` - 分页组件
  - `Badge.vue` - 徽章组件
  - `Tooltip.vue` - 工具提示组件
  - `Alert.vue` - 警告提示组件

**阶段二：业务组件（3 人天）**

**步骤 3：代码相关组件**
- 文件：`frontend/src/components/code/`
- 组件列表：
  - `CodeEditor.vue` - 代码编辑器（集成 Monaco Editor）
  - `DiffViewer.vue` - 差异查看器
  - `FileTree.vue` - 文件树组件
  - `CodeHighlight.vue` - 代码高亮组件

**步骤 4：Git 相关组件**
- 文件：`frontend/src/components/git/`
- 组件列表：
  - `CommitList.vue` - 提交列表
  - `BranchSelector.vue` - 分支选择器
  - `TagList.vue` - 标签列表
  - `CommitInfo.vue` - 提交信息卡片

**步骤 5：PR 相关组件**
- 文件：`frontend/src/components/pr/`
- 组件列表：
  - `ReviewPanel.vue` - 审查面板
  - `CommentThread.vue` - 评论线程
  - `MergeStatus.vue` - 合并状态
  - `PRStatusBadge.vue` - PR 状态徽章

**阶段三：用户体验优化（2 人天）**

**步骤 6：加载状态**
- 文件：`frontend/src/components/ui/`
- 组件列表：
  - `Skeleton.vue` - 骨架屏
  - `LoadingSpinner.vue` - 加载旋转器
  - `ProgressBar.vue` - 进度条

**步骤 7：反馈提示**
- 文件：`frontend/src/components/ui/`
- 组件列表：
  - `Toast.vue` - 消息提示
  - `EmptyState.vue` - 空状态
- 功能：
  - 全局 Toast 服务
  - 空状态插图和提示

**步骤 8：响应式设计**
- 优化移动端适配
- 优化平板适配
- 优化导航栏响应式布局

#### 3.1.3 验收标准

- [ ] 基础组件 ≥ 10 个
- [ ] 业务组件 ≥ 8 个
- [ ] 响应式设计完成
- [ ] 组件文档完善

---

## 四、实施顺序

### 4.1 推荐实施顺序

```
第 1 周：单元测试（仓库模块、PR 模块）
第 2 周：单元测试（AI 模块、CI/CD 模块、安全模块）+ 集成测试
第 3 周：仓库镜像 + GPG 签名
第 4 周：CODEOWNERS + Cherry-pick/Revert
第 5 周：E2E 测试 + 高级功能测试
第 6 周：前端基础组件
第 7 周：前端业务组件 + 用户体验优化
```

### 4.2 依赖关系

```
单元测试 → 集成测试 → E2E 测试
仓库镜像 → 无依赖
GPG 签名 → 无依赖
CODEOWNERS → PR 模块
Cherry-pick/Revert → PR 模块
前端组件库 → 无依赖
```

---

## 五、资源需求

### 5.1 人力需求

| 角色 | 人数 | 工作内容 |
|------|------|---------|
| 后端开发 | 2 人 | 测试编写、高级功能实现 |
| 前端开发 | 1 人 | 组件库开发、用户体验优化 |
| 测试工程师 | 1 人 | 集成测试、E2E 测试 |

### 5.2 技术资源

| 资源 | 用途 |
|------|------|
| 测试数据库 | 集成测试 |
| Playwright | E2E 测试 |
| Monaco Editor | 代码编辑器组件 |

---

## 六、风险管理

### 6.1 技术风险

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|---------|
| 测试覆盖率目标过高 | 低 | 中 | 分阶段提升，优先核心模块 |
| GPG 签名实现复杂 | 中 | 中 | 参考成熟方案，逐步实现 |
| Monaco Editor 集成困难 | 低 | 低 | 使用官方文档和示例 |

### 6.2 进度风险

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|---------|
| 测试编写耗时 | 中 | 中 | 使用 Mock 简化依赖 |
| 组件库开发量大 | 中 | 低 | 优先核心组件 |

---

## 七、验收标准

### 7.1 功能验收

- [ ] 所有功能模块按计划完成
- [ ] 功能测试通过率 100%
- [ ] 性能指标达标

### 7.2 质量验收

- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 集成测试覆盖核心流程
- [ ] E2E 测试通过率 100%
- [ ] 代码审查通过

### 7.3 文档验收

- [ ] API 文档更新
- [ ] 组件文档完善
- [ ] 测试报告生成

---

## 八、后续规划

### 8.1 v0.5 版本（当前计划）

- ✅ P0 任务（AI 审查引擎、CI/CD Runner）
- 🔄 P1 任务（测试覆盖率、高级功能）
- 🔄 P2 任务（前端组件库）

### 8.2 v0.6 版本（下一阶段）

- 容器镜像仓库
- 包管理功能
- Wiki 功能完善
- Pages 功能

### 8.3 v1.0 版本（正式发布）

- 性能优化
- 安全加固
- 文档完善
- 社区建设

---

> **文档版本**：v1.0  
> **最后更新**：2026 年 4 月 24 日
