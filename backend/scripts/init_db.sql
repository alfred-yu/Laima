-- Laima 数据库初始化脚本
-- 版本: v1.0
-- 日期: 2026-04-22

-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(512),
    bio TEXT,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 组织表
CREATE TABLE IF NOT EXISTS organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255),
    description TEXT,
    owner_id INTEGER REFERENCES users(id),
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 仓库表
CREATE TABLE IF NOT EXISTS repositories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    full_path VARCHAR(512) UNIQUE NOT NULL,
    description TEXT,
    owner_type VARCHAR(50) NOT NULL,
    owner_id INTEGER NOT NULL,
    visibility VARCHAR(50) NOT NULL DEFAULT 'private',
    default_branch VARCHAR(255) NOT NULL DEFAULT 'main',
    size BIGINT NOT NULL DEFAULT 0,
    is_fork BOOLEAN NOT NULL DEFAULT false,
    fork_parent_id INTEGER REFERENCES repositories(id),
    is_mirror BOOLEAN NOT NULL DEFAULT false,
    mirror_url VARCHAR(512),
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT repo_owner_fk FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT repo_org_fk FOREIGN KEY (owner_id) REFERENCES organizations(id) ON DELETE CASCADE
);

-- 分支表
CREATE TABLE IF NOT EXISTS branches (
    id SERIAL PRIMARY KEY,
    repository_id INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    commit_sha VARCHAR(40) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(repository_id, name)
);

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    repository_id INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    commit_sha VARCHAR(40) NOT NULL,
    message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(repository_id, name)
);

-- Pull Request 表
CREATE TABLE IF NOT EXISTS pull_requests (
    id SERIAL PRIMARY KEY,
    number INTEGER NOT NULL,
    title VARCHAR(512) NOT NULL,
    description TEXT,
    repository_id INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
    author_id INTEGER REFERENCES users(id),
    source_repo_id INTEGER REFERENCES repositories(id),
    source_branch VARCHAR(255) NOT NULL,
    target_branch VARCHAR(255) NOT NULL,
    state VARCHAR(50) NOT NULL DEFAULT 'open',
    merge_state VARCHAR(50) NOT NULL DEFAULT 'checking',
    review_mode VARCHAR(50) NOT NULL DEFAULT 'standard',
    head_commit_sha VARCHAR(40) NOT NULL,
    base_commit_sha VARCHAR(40) NOT NULL,
    merge_commit_sha VARCHAR(40),
    merged_by INTEGER REFERENCES users(id),
    merged_at TIMESTAMP,
    closed_at TIMESTAMP,
    is_draft BOOLEAN NOT NULL DEFAULT false,
    ai_review_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(repository_id, number)
);

-- 审查表
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    pull_request_id INTEGER REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id INTEGER REFERENCES users(id),
    state VARCHAR(50) NOT NULL,
    score INTEGER,
    body TEXT,
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 审查评论表
CREATE TABLE IF NOT EXISTS review_comments (
    id SERIAL PRIMARY KEY,
    pull_request_id INTEGER REFERENCES pull_requests(id) ON DELETE CASCADE,
    review_id INTEGER REFERENCES reviews(id) ON DELETE SET NULL,
    author_id INTEGER REFERENCES users(id),
    type VARCHAR(50) NOT NULL DEFAULT 'human',
    path VARCHAR(512) NOT NULL,
    line INTEGER,
    diff_hunk TEXT,
    body TEXT NOT NULL,
    resolution VARCHAR(50),
    suggestion TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Issue 表
CREATE TABLE IF NOT EXISTS issues (
    id SERIAL PRIMARY KEY,
    number INTEGER NOT NULL,
    title VARCHAR(512) NOT NULL,
    description TEXT,
    repository_id INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
    author_id INTEGER REFERENCES users(id),
    assignee_id INTEGER REFERENCES users(id),
    state VARCHAR(50) NOT NULL DEFAULT 'open',
    milestone_id INTEGER,
    labels JSONB DEFAULT '[]',
    priority VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(repository_id, number)
);

-- 里程碑表
CREATE TABLE IF NOT EXISTS milestones (
    id SERIAL PRIMARY KEY,
    repository_id INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    state VARCHAR(50) NOT NULL DEFAULT 'open',
    due_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 流水线表
CREATE TABLE IF NOT EXISTS pipelines (
    id SERIAL PRIMARY KEY,
    repository_id INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
    commit_sha VARCHAR(40) NOT NULL,
    ref VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    trigger VARCHAR(50) NOT NULL,
    yaml_content TEXT,
    duration INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 任务表
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    pipeline_id INTEGER REFERENCES pipelines(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    stage VARCHAR(255),
    duration INTEGER,
    log TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 安全扫描表
CREATE TABLE IF NOT EXISTS security_scans (
    id SERIAL PRIMARY KEY,
    repository_id INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
    commit_sha VARCHAR(40) NOT NULL,
    scan_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    severity VARCHAR(50),
    findings JSONB DEFAULT '[]',
    duration INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 组织成员表
CREATE TABLE IF NOT EXISTS organization_members (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(organization_id, user_id)
);

-- 仓库成员表
CREATE TABLE IF NOT EXISTS repository_members (
    id SERIAL PRIMARY KEY,
    repository_id INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'developer',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(repository_id, user_id)
);

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    action VARCHAR(255) NOT NULL,
    target_type VARCHAR(50),
    target_id INTEGER,
    details JSONB DEFAULT '{}',
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_repositories_full_path ON repositories(full_path);
CREATE INDEX IF NOT EXISTS idx_repositories_owner ON repositories(owner_type, owner_id);
CREATE INDEX IF NOT EXISTS idx_pull_requests_repository ON pull_requests(repository_id, number);
CREATE INDEX IF NOT EXISTS idx_issues_repository ON issues(repository_id, number);
CREATE INDEX IF NOT EXISTS idx_pipelines_repository ON pipelines(repository_id);
CREATE INDEX IF NOT EXISTS idx_security_scans_repository ON security_scans(repository_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);

-- 创建超级管理员用户（默认密码：admin123）
INSERT INTO users (username, email, password_hash) 
VALUES ('admin', 'admin@laima.local', crypt('admin123', gen_salt('bf')))
ON CONFLICT (username) DO NOTHING;

-- 创建默认组织
INSERT INTO organizations (name, display_name, description, owner_id) 
VALUES ('laima', 'Laima', 'Laima 官方组织', (SELECT id FROM users WHERE username = 'admin'))
ON CONFLICT (name) DO NOTHING;

-- 创建默认仓库
INSERT INTO repositories (name, full_path, description, owner_type, owner_id, visibility) 
VALUES ('laima', 'laima/laima', 'Laima 主仓库', 'org', (SELECT id FROM organizations WHERE name = 'laima'), 'public')
ON CONFLICT (full_path) DO NOTHING;