<template>
  <div class="cicd-detail">
    <h1>流水线 #{{ pipeline.id }}</h1>
    <div class="pipeline-meta">
      <span class="pipeline-repo">{{ pipeline.repo }}</span>
      <span class="pipeline-author">由 {{ pipeline.author }} 触发</span>
      <span class="pipeline-time">{{ pipeline.time }}</span>
      <span class="pipeline-branch">{{ pipeline.branch }}</span>
      <span class="pipeline-status" :class="pipeline.status">{{ pipeline.status }}</span>
    </div>
    <div class="pipeline-content">
      <div class="pipeline-jobs">
        <h2>任务</h2>
        <div class="jobs-list">
          <div v-for="job in pipeline.jobs" :key="job.id" class="job-item">
            <div class="job-info">
              <h3 class="job-name">{{ job.name }}</h3>
              <div class="job-meta">
                <span class="job-status" :class="job.status">{{ job.status }}</span>
                <span class="job-duration">{{ job.duration }}</span>
              </div>
            </div>
            <div class="job-logs">
              <button class="btn">查看日志</button>
            </div>
          </div>
        </div>
      </div>
      <div class="pipeline-logs">
        <h2>日志</h2>
        <div class="logs-content">
          <pre>{{ pipeline.logs }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const pipeline = ref({
  id: 1,
  repo: 'user1/laima',
  author: 'user1',
  time: '2小时前',
  branch: 'main',
  status: 'running',
  jobs: [
    {
      id: 1,
      name: '构建',
      status: 'success',
      duration: '1m 30s'
    },
    {
      id: 2,
      name: '测试',
      status: 'running',
      duration: '2m 15s'
    },
    {
      id: 3,
      name: '部署',
      status: 'pending',
      duration: '0s'
    }
  ],
  logs: `[2024-01-01 10:00:00] 开始构建
[2024-01-01 10:00:01] 安装依赖
[2024-01-01 10:00:10] 构建完成
[2024-01-01 10:00:10] 开始测试
[2024-01-01 10:00:15] 运行单元测试
[2024-01-01 10:02:30] 测试中...`
})
</script>

<style scoped>
.cicd-detail {
  padding: 20px;
}

.pipeline-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin: 16px 0;
  font-size: 14px;
  color: var(--text-secondary);
  transition: color 0.3s;
}

.pipeline-status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.pipeline-status.running {
  background: var(--warning-color);
  color: #fff;
}

.pipeline-status.success {
  background: var(--success-color);
  color: #fff;
}

.pipeline-status.failed {
  background: var(--danger-color);
  color: #fff;
}

.pipeline-content {
  margin-top: 24px;
}

.pipeline-jobs {
  margin-bottom: 24px;
}

.pipeline-jobs h2 {
  margin: 0 0 16px 0;
  font-size: 18px;
  color: var(--text-primary);
  transition: color 0.3s;
}

.jobs-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.job-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: background-color 0.3s, box-shadow 0.3s;
}

.job-info {
  flex: 1;
}

.job-name {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: var(--text-primary);
  transition: color 0.3s;
}

.job-meta {
  display: flex;
  gap: 12px;
  font-size: 14px;
  color: var(--text-secondary);
  transition: color 0.3s;
}

.job-status {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.job-status.success {
  background: var(--success-color);
  color: #fff;
}

.job-status.running {
  background: var(--warning-color);
  color: #fff;
}

.job-status.pending {
  background: var(--border-color);
  color: var(--text-secondary);
  transition: background-color 0.3s, color 0.3s;
}

.job-status.failed {
  background: var(--danger-color);
  color: #fff;
}

.btn {
  padding: 6px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
}

.btn:hover {
  background: var(--bg-secondary);
}

.pipeline-logs h2 {
  margin: 0 0 16px 0;
  font-size: 18px;
  color: var(--text-primary);
  transition: color 0.3s;
}

.logs-content {
  background: var(--bg-primary);
  border-radius: 8px;
  padding: 16px;
  overflow: auto;
  max-height: 400px;
  transition: background-color 0.3s;
}

.logs-content pre {
  margin: 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  color: var(--text-primary);
  transition: color 0.3s;
}
</style>