# Laima API 文档

## 1. 认证相关 API

### 1.1 登录
- **URL**: `/api/auth/login`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **响应**:
  ```json
  {
    "token": "string",
    "user": {
      "id": 1,
      "username": "string",
      "email": "string"
    }
  }
  ```

### 1.2 注册
- **URL**: `/api/auth/register`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "username": "string",
    "email": "string",
    "password": "string"
  }
  ```
- **响应**:
  ```json
  {
    "token": "string",
    "user": {
      "id": 1,
      "username": "string",
      "email": "string"
    }
  }
  ```

### 1.3 获取当前用户
- **URL**: `/api/users/me`
- **方法**: `GET`
- **认证**: 需要 JWT Token
- **响应**:
  ```json
  {
    "id": 1,
    "username": "string",
    "email": "string"
  }
  ```

## 2. 仓库相关 API

### 2.1 列出仓库
- **URL**: `/api/v1/repos`
- **方法**: `GET`
- **查询参数**:
  - `q`: 搜索关键词
  - `visibility`: 可见性 (public/private)
  - `page`: 页码
  - `per_page`: 每页数量
- **响应**:
  ```json
  {
    "items": [
      {
        "id": 1,
        "name": "string",
        "description": "string",
        "visibility": "public",
        "default_branch": "main",
        "full_path": "user/repo"
      }
    ]
  }
  ```

### 2.2 创建仓库
- **URL**: `/api/v1/repos`
- **方法**: `POST`
- **认证**: 需要 JWT Token
- **请求体**:
  ```json
  {
    "name": "string",
    "description": "string",
    "visibility": "private",
    "auto_init": false
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "name": "string",
    "description": "string",
    "visibility": "private",
    "default_branch": "main"
  }
  ```

### 2.3 获取仓库详情
- **URL**: `/api/v1/repos/{owner}/{repo}`
- **方法**: `GET`
- **响应**:
  ```json
  {
    "id": 1,
    "name": "string",
    "description": "string",
    "visibility": "public",
    "default_branch": "main"
  }
  ```

## 3. Pull Request 相关 API

### 3.1 列出 PR
- **URL**: `/api/v1/prs`
- **方法**: `GET`
- **查询参数**:
  - `repository_id`: 仓库 ID
  - `state`: 状态 (open/closed)
  - `page`: 页码
  - `per_page`: 每页数量
- **响应**:
  ```json
  {
    "items": [
      {
        "id": 1,
        "title": "string",
        "description": "string",
        "state": "open",
        "repository_id": 1,
        "source_branch": "feature",
        "target_branch": "main"
      }
    ]
  }
  ```

### 3.2 创建 PR
- **URL**: `/api/v1/prs`
- **方法**: `POST`
- **认证**: 需要 JWT Token
- **请求体**:
  ```json
  {
    "repository_id": 1,
    "title": "string",
    "description": "string",
    "source_branch": "feature",
    "target_branch": "main"
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "title": "string",
    "description": "string",
    "state": "open",
    "repository_id": 1,
    "source_branch": "feature",
    "target_branch": "main"
  }
  ```

### 3.3 合并 PR
- **URL**: `/api/v1/prs/{id}/merge`
- **方法**: `POST`
- **认证**: 需要 JWT Token
- **响应**:
  ```json
  {
    "success": true,
    "message": "PR merged successfully"
  }
  ```

## 4. Issue 相关 API

### 4.1 列出 Issue
- **URL**: `/api/v1/issues`
- **方法**: `GET`
- **查询参数**:
  - `repository_id`: 仓库 ID
  - `state`: 状态 (open/closed)
  - `page`: 页码
  - `per_page`: 每页数量
- **响应**:
  ```json
  {
    "items": [
      {
        "id": 1,
        "title": "string",
        "description": "string",
        "state": "open",
        "repository_id": 1
      }
    ]
  }
  ```

### 4.2 创建 Issue
- **URL**: `/api/v1/issues`
- **方法**: `POST`
- **认证**: 需要 JWT Token
- **请求体**:
  ```json
  {
    "repository_id": 1,
    "title": "string",
    "description": "string"
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "title": "string",
    "description": "string",
    "state": "open",
    "repository_id": 1
  }
  ```

## 5. CI/CD 相关 API

### 5.1 列出流水线
- **URL**: `/api/v1/cicd`
- **方法**: `GET`
- **查询参数**:
  - `repository_id`: 仓库 ID
  - `status`: 状态 (pending/running/success/failed)
  - `page`: 页码
  - `per_page`: 每页数量
- **响应**:
  ```json
  {
    "items": [
      {
        "id": 1,
        "repository_id": 1,
        "status": "success",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
  ```

### 5.2 获取流水线详情
- **URL**: `/api/v1/cicd/{id}`
- **方法**: `GET`
- **响应**:
  ```json
  {
    "id": 1,
    "repository_id": 1,
    "status": "success",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:01:00Z"
  }
  ```

## 6. AI 审查相关 API

### 6.1 触发 AI 审查
- **URL**: `/api/v1/ai/review`
- **方法**: `POST`
- **认证**: 需要 JWT Token
- **请求体**:
  ```json
  {
    "pr_id": 1
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "pr_id": 1,
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z"
  }
  ```

### 6.2 获取 AI 审查结果
- **URL**: `/api/v1/ai/review/{id}`
- **方法**: `GET`
- **响应**:
  ```json
  {
    "id": 1,
    "pr_id": 1,
    "status": "completed",
    "result": {
      "issues": [
        {
          "line": 10,
          "message": "Potential bug here"
        }
      ]
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:01:00Z"
  }
  ```

## 7. 组织相关 API

### 7.1 创建组织
- **URL**: `/api/orgs`
- **方法**: `POST`
- **认证**: 需要 JWT Token
- **请求体**:
  ```json
  {
    "name": "string",
    "description": "string"
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "name": "string",
    "description": "string"
  }
  ```

### 7.2 获取组织详情
- **URL**: `/api/orgs/{id}`
- **方法**: `GET`
- **响应**:
  ```json
  {
    "id": 1,
    "name": "string",
    "description": "string"
  }
  ```
