# Laima 开发指南

本指南将帮助你快速设置 Laima 的开发环境。

## 前置要求

- Go 1.22+
- Node.js 18+ and npm
- Docker and Docker Compose (可选)
- PostgreSQL 16+ (如果不使用 Docker)
- Redis 7+ (如果不使用 Docker)

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd laima
```

### 2. 启动依赖服务（使用 Docker）

项目包含一个 Docker Compose 配置，可以快速启动所有依赖服务：

```bash
cd backend
docker-compose up -d
```

这将启动以下服务：
- PostgreSQL (数据库)
- Redis (缓存和队列)
- MinIO (对象存储)
- Meilisearch (搜索服务)

### 3. 配置环境变量

后端配置：
```bash
cd backend
cp .env.example .env
# 根据需要编辑 .env 文件
```

前端配置（可选）：
```bash
cd frontend
# 创建 .env.local 文件来覆盖默认配置
echo "VITE_API_BASE_URL=http://localhost:8080" > .env.local
```

### 4. 启动后端服务

```bash
cd backend
go mod download
go run cmd/laima/main.go
```

后端服务将在 http://localhost:8080 启动。

### 5. 启动前端服务

```bash
cd frontend
npm install
npm run dev
```

前端服务将在 http://localhost:5173 启动。

## 项目结构

```
laima/
├── backend/              # 后端代码
│   ├── cmd/laima/       # 主程序入口
│   ├── internal/        # 内部模块
│   │   ├── repo/        # 仓库管理
│   │   ├── user/        # 用户管理
│   │   ├── pr/          # Pull Request
│   │   ├── ai/          # AI 功能
│   │   ├── cicd/        # CI/CD
│   │   ├── issue/       # Issue 管理
│   │   └── middleware/  # 中间件
│   └── scripts/         # 数据库脚本
├── frontend/            # 前端代码
│   ├── src/
│   │   ├── components/  # 组件
│   │   ├── views/      # 页面
│   │   ├── stores/     # 状态管理
│   │   └── services/   # API 服务
│   └── public/
└── docs/               # 文档
```

## API 文档

后端 API 将在启动后通过以下路径访问（如果已配置）：
- Swagger UI: http://localhost:8080/swagger/
- OpenAPI JSON: http://localhost:8080/swagger/doc.json

## 核心功能开发进度

### 已完成功能：
- [x] 项目基础架构
- [x] 用户认证系统（注册、登录、JWT）
- [x] 仓库管理基础功能（CRUD）
- [x] 前端基础组件和主题系统
- [x] 前端认证和仓库管理页面

### 待完成功能：
- [ ] Git 仓库实际存储和操作
- [ ] Pull Request 完整功能
- [ ] AI 代码审查功能
- [ ] CI/CD 流水线功能
- [ ] Issue 管理功能
- [ ] 权限系统
- [ ] 组织管理
- [ ] 等等...

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证。
