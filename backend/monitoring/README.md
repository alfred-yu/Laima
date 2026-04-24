# Laima 监控配置

## 1. 系统监控

### 1.1 Prometheus 配置

创建 `prometheus.yml` 文件：

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'laima'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: '/metrics'

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres:9187']

  - job_name: 'redis'
    static_configs:
      - targets: ['redis:9121']

  - job_name: 'minio'
    static_configs:
      - targets: ['minio:9000']
```

### 1.2 Grafana 仪表板

创建 Laima 系统仪表板，包含以下面板：

- **系统概览**：CPU、内存、磁盘使用情况
- **API 性能**：请求数、响应时间、错误率
- **数据库性能**：查询执行时间、连接数
- **Git 操作**：仓库创建、PR 操作、代码审查
- **CI/CD 状态**：流水线执行状态、成功率

## 2. 日志管理

### 2.1 日志配置

修改 `main.go` 文件，添加结构化日志：

```go
import (
	"github.com/rs/zerolog/log"
)

func main() {
	// 配置日志
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// 其他代码...
}
```

### 2.2 ELK 堆栈配置

使用 Elasticsearch、Logstash 和 Kibana 进行日志管理：

- **Elasticsearch**：存储日志数据
- **Logstash**：处理和转换日志
- **Kibana**：可视化日志和创建仪表板

## 3. 健康检查

### 3.1 API 健康检查

使用 `/health` 端点进行健康检查：

```bash
curl http://localhost:8080/health
```

### 3.2 服务健康检查

为每个服务添加健康检查配置，确保系统正常运行。

## 4. 告警配置

### 4.1 Prometheus 告警规则

创建 `alerts.yml` 文件：

```yaml
groups:
- name: laima_alerts
  rules:
  - alert: HighCPUUsage
    expr: cpu_usage_percent > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High CPU usage"
      description: "CPU usage is above 80% for 5 minutes"

  - alert: HighMemoryUsage
    expr: memory_usage_percent > 85
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High memory usage"
      description: "Memory usage is above 85% for 5 minutes"

  - alert: APIErrorRate
    expr: rate(api_errors_total[5m]) > 0.1
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High API error rate"
      description: "API error rate is above 10% for 5 minutes"
```

### 4.2 Alertmanager 配置

配置 Alertmanager 发送告警通知：

```yaml
global:
  resolve_timeout: 5m

route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1h
  receiver: 'email'

receivers:
- name: 'email'
  email_configs:
  - to: 'admin@example.com'
    from: 'alerts@example.com'
    smarthost: 'smtp.example.com:587'
    auth_username: 'alerts'
    auth_password: 'password'
```

## 5. 部署监控栈

使用 Docker Compose 部署监控栈：

```yaml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana

  alertmanager:
    image: prom/alertmanager:latest
    volumes:
      - ./alertmanager.yml:/etc/alertmanager/alertmanager.yml
    ports:
      - "9093:9093"

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.17.0
    environment:
      - discovery.type=single-node
    ports:
      - "9200:9200"
    volumes:
      - es_data:/usr/share/elasticsearch/data

  logstash:
    image: docker.elastic.co/logstash/logstash:7.17.0
    volumes:
      - ./logstash.conf:/etc/logstash/conf.d/logstash.conf
    ports:
      - "5044:5044"

  kibana:
    image: docker.elastic.co/kibana/kibana:7.17.0
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch

volumes:
  grafana_data:
  es_data:
```

## 6. 监控最佳实践

1. **设置合理的告警阈值**：根据系统实际情况调整告警阈值
2. **定期审查监控数据**：分析监控数据，识别潜在问题
3. **优化监控配置**：根据系统变化调整监控配置
4. **建立监控文档**：记录监控配置和告警处理流程
5. **培训团队成员**：确保团队成员了解监控系统和告警处理

## 7. 故障排查

### 7.1 常见问题

- **API 响应缓慢**：检查数据库性能、网络延迟
- **Git 操作失败**：检查文件系统权限、磁盘空间
- **CI/CD 流水线失败**：检查构建环境、依赖项
- **AI 审查超时**：检查 API 调用限制、网络连接

### 7.2 排查步骤

1. 查看监控仪表板，识别异常指标
2. 检查日志，查找错误信息
3. 分析系统资源使用情况
4. 检查依赖服务状态
5. 进行针对性测试，验证问题
