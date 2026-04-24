package app

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"laima/internal/cicd/domain"
	prapp "laima/internal/pr/app"

	"gorm.io/gorm"
)

// CICDService CI/CD服务接口
type CICDService interface {
	// 流水线管理
	CreatePipeline(ctx context.Context, req *domain.PipelineRequest) (*domain.Pipeline, error)
	GetPipeline(ctx context.Context, pipelineID int) (*domain.Pipeline, error)
	ListPipelines(ctx context.Context, filter *domain.PipelineFilter) ([]*domain.Pipeline, int64, error)
	CancelPipeline(ctx context.Context, pipelineID int) (*domain.Pipeline, error)

	// 任务管理
	GetJobs(ctx context.Context, pipelineID int) ([]*domain.Job, error)
	GetJob(ctx context.Context, jobID int) (*domain.Job, error)
	UpdateJobStatus(ctx context.Context, jobID int, status string) (*domain.Job, error)
	AddJobLog(ctx context.Context, jobID int, log string) error

	// 流水线解析
	ParsePipelineYAML(ctx context.Context, yamlContent string) ([]*domain.Job, error)

	// 集成触发
	TriggerPipelineForPR(ctx context.Context, prID int) (*domain.Pipeline, error)
	TriggerPipelineForPush(ctx context.Context, repoID int, commitSHA string, ref string) (*domain.Pipeline, error)
}

// cicdService CI/CD服务实现
type cicdService struct {
	db    *gorm.DB
	prSvc prapp.PRService
}

// NewCICDService 创建CI/CD服务实例
func NewCICDService(db *gorm.DB, prSvc prapp.PRService) CICDService {
	return &cicdService{db: db, prSvc: prSvc}
}

// CreatePipeline 创建流水线
func (s *cicdService) CreatePipeline(ctx context.Context, req *domain.PipelineRequest) (*domain.Pipeline, error) {
	// 创建流水线记录
	pipeline := &domain.Pipeline{
		RepositoryID: req.RepositoryID,
		PRID:         req.PRID,
		CommitSHA:    req.CommitSHA,
		Ref:          req.Ref,
		Status:       domain.PipelineStatusPending,
		Trigger:      req.Trigger,
	}

	if err := s.db.Create(pipeline).Error; err != nil {
		return nil, err
	}

	var jobs []*domain.Job

	// 检查是否有YAML配置
	if pipeline.YAMLContent != "" {
		// 解析YAML配置
		parsedJobs, err := s.ParsePipelineYAML(ctx, pipeline.YAMLContent)
		if err != nil {
			// 解析失败，使用默认任务
			jobs = s.getDefaultJobs(pipeline.ID)
		} else {
			// 使用解析的任务
			for _, job := range parsedJobs {
				job.PipelineID = pipeline.ID
			}
			jobs = parsedJobs
		}
	} else {
		// 使用默认任务
		jobs = s.getDefaultJobs(pipeline.ID)
	}

	// 保存任务
	for _, job := range jobs {
		if err := s.db.Create(job).Error; err != nil {
			return nil, err
		}
	}

	// 异步触发流水线执行
	go func() {
		ctx := context.Background()
		for _, job := range jobs {
			// 更新任务状态为运行中
			_, _ = s.UpdateJobStatus(ctx, job.ID, domain.JobStatusRunning)
			
			// 模拟任务执行
			s.executeJob(ctx, job)
		}
	}()

	return pipeline, nil
}

// executeJob 执行单个任务
func (s *cicdService) executeJob(ctx context.Context, job *domain.Job) {
	// 模拟任务执行过程
	log := "Starting job: " + job.Name + "\n"
	
	// 模拟执行脚本
	simulatedScript := []string{
		"echo 'Hello from '" + job.Name,
		"echo 'Executing stage: '" + job.Stage,
		"sleep 1", // 模拟执行时间
		"echo 'Job completed successfully'",
	}
	
	for _, script := range simulatedScript {
		log += "$ " + script + "\n"
		// 模拟命令输出
		log += "Output: " + script + " executed\n"
	}
	
	// 添加日志
	_ = s.AddJobLog(ctx, job.ID, log)
	
	// 随机模拟成功或失败
	// 这里简单起见，总是成功
	_, _ = s.UpdateJobStatus(ctx, job.ID, domain.JobStatusSuccess)
}

// getDefaultJobs 获取默认任务列表
func (s *cicdService) getDefaultJobs(pipelineID int) []*domain.Job {
	return []*domain.Job{
		{
			PipelineID: pipelineID,
			Name:       "build",
			Status:     domain.JobStatusPending,
			Stage:      "build",
		},
		{
			PipelineID: pipelineID,
			Name:       "test",
			Status:     domain.JobStatusPending,
			Stage:      "test",
		},
	}
}

// GetPipeline 根据ID获取流水线
func (s *cicdService) GetPipeline(ctx context.Context, pipelineID int) (*domain.Pipeline, error) {
	var pipeline domain.Pipeline
	result := s.db.First(&pipeline, pipelineID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pipeline, nil
}

// ListPipelines 列出流水线
func (s *cicdService) ListPipelines(ctx context.Context, filter *domain.PipelineFilter) ([]*domain.Pipeline, int64, error) {
	var pipelines []*domain.Pipeline
	var total int64

	query := s.db.Model(&domain.Pipeline{})

	// 应用过滤条件
	if filter.RepositoryID > 0 {
		query = query.Where("repository_id = ?", filter.RepositoryID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Ref != "" {
		query = query.Where("ref = ?", filter.Ref)
	}
	if filter.CommitSHA != "" {
		query = query.Where("commit_sha = ?", filter.CommitSHA)
	}

	// 计算总数
	query.Count(&total)

	// 分页
	offset := (filter.Page - 1) * filter.PerPage
	result := query.Offset(offset).Limit(filter.PerPage).Order("created_at DESC").Find(&pipelines)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return pipelines, total, nil
}

// CancelPipeline 取消流水线
func (s *cicdService) CancelPipeline(ctx context.Context, pipelineID int) (*domain.Pipeline, error) {
	var pipeline domain.Pipeline
	if err := s.db.First(&pipeline, pipelineID).Error; err != nil {
		return nil, err
	}

	// 更新流水线状态
	pipeline.Status = domain.PipelineStatusCanceled
	if err := s.db.Save(&pipeline).Error; err != nil {
		return nil, err
	}

	// 更新相关任务状态
	s.db.Model(&domain.Job{}).Where("pipeline_id = ?", pipelineID).Updates(map[string]interface{}{
		"status": domain.JobStatusCanceled,
	})

	return &pipeline, nil
}

// GetJobs 获取流水线的任务列表
func (s *cicdService) GetJobs(ctx context.Context, pipelineID int) ([]*domain.Job, error) {
	var jobs []*domain.Job
	result := s.db.Where("pipeline_id = ?", pipelineID).Order("created_at ASC").Find(&jobs)
	if result.Error != nil {
		return nil, result.Error
	}
	return jobs, nil
}

// GetJob 根据ID获取任务
func (s *cicdService) GetJob(ctx context.Context, jobID int) (*domain.Job, error) {
	var job domain.Job
	result := s.db.First(&job, jobID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &job, nil
}

// UpdateJobStatus 更新任务状态
func (s *cicdService) UpdateJobStatus(ctx context.Context, jobID int, status string) (*domain.Job, error) {
	var job domain.Job
	if err := s.db.First(&job, jobID).Error; err != nil {
		return nil, err
	}

	job.Status = status
	if err := s.db.Save(&job).Error; err != nil {
		return nil, err
	}

	// 更新流水线状态
	s.updatePipelineStatus(ctx, job.PipelineID)

	return &job, nil
}

// AddJobLog 添加任务日志
func (s *cicdService) AddJobLog(ctx context.Context, jobID int, log string) error {
	var job domain.Job
	if err := s.db.First(&job, jobID).Error; err != nil {
		return err
	}

	job.Log += log
	return s.db.Save(&job).Error
}

// PipelineYAML 表示CI/CD配置文件的结构
type PipelineYAML struct {
	Stages []string          `yaml:"stages"`
	Jobs   map[string]JobYAML `yaml:"jobs"`
}

// JobYAML 表示单个任务的配置
type JobYAML struct {
	Stage    string   `yaml:"stage"`
	Script   []string `yaml:"script"`
	Artifacts []string `yaml:"artifacts"`
	Only      []string `yaml:"only"`
	Except    []string `yaml:"except"`
}

// ParsePipelineYAML 解析流水线YAML配置
func (s *cicdService) ParsePipelineYAML(ctx context.Context, yamlContent string) ([]*domain.Job, error) {
	// 解析YAML文件
	var pipeline PipelineYAML
	if err := yaml.Unmarshal([]byte(yamlContent), &pipeline); err != nil {
		return nil, fmt.Errorf("解析YAML失败: %w", err)
	}

	// 构建任务列表
	var jobs []*domain.Job
	for jobName, jobConfig := range pipeline.Jobs {
		job := &domain.Job{
			Name:   jobName,
			Status: domain.JobStatusPending,
			Stage:  jobConfig.Stage,
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// TriggerPipelineForPR 为PR触发流水线
func (s *cicdService) TriggerPipelineForPR(ctx context.Context, prID int) (*domain.Pipeline, error) {
	// 获取PR信息
	pr, err := s.prSvc.GetPR(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("获取PR信息失败: %w", err)
	}

	// 构建流水线请求
	req := &domain.PipelineRequest{
		RepositoryID: pr.RepositoryID,
		CommitSHA:    pr.HeadCommitSHA,
		Ref:          fmt.Sprintf("refs/heads/%s", pr.SourceBranch),
		Trigger:      "pr",
		PRID:         prID,
	}

	return s.CreatePipeline(ctx, req)
}

// TriggerPipelineForPush 为推送触发流水线
func (s *cicdService) TriggerPipelineForPush(ctx context.Context, repoID int, commitSHA string, ref string) (*domain.Pipeline, error) {
	req := &domain.PipelineRequest{
		RepositoryID: repoID,
		CommitSHA:    commitSHA,
		Ref:          ref,
		Trigger:      "push",
	}
	return s.CreatePipeline(ctx, req)
}

// updatePipelineStatus 更新流水线状态
func (s *cicdService) updatePipelineStatus(ctx context.Context, pipelineID int) error {
	var jobs []*domain.Job
	result := s.db.Where("pipeline_id = ?", pipelineID).Find(&jobs)
	if result.Error != nil {
		return result.Error
	}

	hasRunning := false
	hasFailed := false
	hasPending := false
	allSuccess := true

	for _, job := range jobs {
		switch job.Status {
		case domain.JobStatusRunning:
			hasRunning = true
		case domain.JobStatusFailed:
			hasFailed = true
		case domain.JobStatusPending:
			hasPending = true
		case domain.JobStatusSuccess:
			// nothing
		default:
			allSuccess = false
		}
	}

	var status string
	switch {
	case hasFailed:
		status = domain.PipelineStatusFailed
	case hasRunning:
		status = domain.PipelineStatusRunning
	case hasPending:
		status = domain.PipelineStatusPending
	case allSuccess:
		status = domain.PipelineStatusSuccess
	default:
		status = domain.PipelineStatusPending
	}

	return s.db.Model(&domain.Pipeline{}).Where("id = ?", pipelineID).Update("status", status).Error
}
