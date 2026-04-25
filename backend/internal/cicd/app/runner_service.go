package app

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"laima/internal/cicd/domain"
)

var (
	ErrRunnerNotFound      = errors.New("runner not found")
	ErrRunnerAlreadyExists = errors.New("runner already exists")
	ErrInvalidToken        = errors.New("invalid runner token")
	ErrRunnerBusy          = errors.New("runner is busy")
)

type RunnerService interface {
	RegisterRunner(req *domain.RunnerRegistration) (*domain.RunnerResponse, error)
	Heartbeat(req *domain.RunnerHeartbeat) error
	RequestJob(req *domain.RunnerJobRequest) (*domain.Job, error)
	UpdateJobStatus(req *domain.RunnerJobUpdate) error
	ListRunners() ([]*domain.RunnerResponse, error)
	GetRunner(id int) (*domain.RunnerResponse, error)
	DeleteRunner(id int) error
}

type runnerService struct {
	runners map[int]*domain.Runner
	jobs    map[int]*domain.Job
}

func NewRunnerService() RunnerService {
	return &runnerService{
		runners: make(map[int]*domain.Runner),
		jobs:    make(map[int]*domain.Job),
	}
}

func (s *runnerService) RegisterRunner(req *domain.RunnerRegistration) (*domain.RunnerResponse, error) {
	for _, runner := range s.runners {
		if runner.Name == req.Name {
			return nil, ErrRunnerAlreadyExists
		}
	}

	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	if req.MaxJobs == 0 {
		req.MaxJobs = 2
	}

	runner := &domain.Runner{
		ID:          len(s.runners) + 1,
		Name:        req.Name,
		Token:       token,
		Description: req.Description,
		Status:      domain.RunnerStatusOffline,
		Executor:    req.Executor,
		Tags:        req.Tags,
		MaxJobs:     req.MaxJobs,
		ActiveJobs:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.runners[runner.ID] = runner

	return toRunnerResponse(runner), nil
}

func (s *runnerService) Heartbeat(req *domain.RunnerHeartbeat) error {
	runner, exists := s.runners[req.RunnerID]
	if !exists {
		return ErrRunnerNotFound
	}

	if runner.Token != req.Token {
		return ErrInvalidToken
	}

	runner.LastSeenAt = time.Now()
	runner.Version = req.Version

	if req.Status != "" {
		runner.Status = req.Status
	} else {
		if runner.ActiveJobs > 0 {
			runner.Status = domain.RunnerStatusBusy
		} else {
			runner.Status = domain.RunnerStatusOnline
		}
	}

	runner.UpdatedAt = time.Now()

	return nil
}

func (s *runnerService) RequestJob(req *domain.RunnerJobRequest) (*domain.Job, error) {
	runner, exists := s.runners[req.RunnerID]
	if !exists {
		return nil, ErrRunnerNotFound
	}

	if runner.Token != req.Token {
		return nil, ErrInvalidToken
	}

	if runner.ActiveJobs >= runner.MaxJobs {
		return nil, ErrRunnerBusy
	}

	for _, job := range s.jobs {
		if job.Status == domain.JobStatusPending {
			job.Status = domain.JobStatusRunning
			runner.ActiveJobs++
			runner.Status = domain.RunnerStatusBusy
			runner.UpdatedAt = time.Now()
			return job, nil
		}
	}

	return nil, nil
}

func (s *runnerService) UpdateJobStatus(req *domain.RunnerJobUpdate) error {
	runner, exists := s.runners[req.RunnerID]
	if !exists {
		return ErrRunnerNotFound
	}

	if runner.Token != req.Token {
		return ErrInvalidToken
	}

	job, exists := s.jobs[req.JobID]
	if !exists {
		return errors.New("job not found")
	}

	job.Status = req.Status
	job.Log = req.Log
	job.UpdatedAt = time.Now()

	if req.Status == domain.JobStatusSuccess || 
	   req.Status == domain.JobStatusFailed || 
	   req.Status == domain.JobStatusCanceled {
		runner.ActiveJobs--
		if runner.ActiveJobs <= 0 {
			runner.ActiveJobs = 0
			runner.Status = domain.RunnerStatusOnline
		}
		runner.UpdatedAt = time.Now()
	}

	return nil
}

func (s *runnerService) ListRunners() ([]*domain.RunnerResponse, error) {
	var runners []*domain.RunnerResponse
	for _, runner := range s.runners {
		runners = append(runners, toRunnerResponse(runner))
	}
	return runners, nil
}

func (s *runnerService) GetRunner(id int) (*domain.RunnerResponse, error) {
	runner, exists := s.runners[id]
	if !exists {
		return nil, ErrRunnerNotFound
	}
	return toRunnerResponse(runner), nil
}

func (s *runnerService) DeleteRunner(id int) error {
	if _, exists := s.runners[id]; !exists {
		return ErrRunnerNotFound
	}
	delete(s.runners, id)
	return nil
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func toRunnerResponse(runner *domain.Runner) *domain.RunnerResponse {
	return &domain.RunnerResponse{
		ID:          runner.ID,
		Name:        runner.Name,
		Description: runner.Description,
		Status:      runner.Status,
		Executor:    runner.Executor,
		Tags:        runner.Tags,
		MaxJobs:     runner.MaxJobs,
		ActiveJobs:  runner.ActiveJobs,
		Version:     runner.Version,
		IPAddress:   runner.IPAddress,
		LastSeenAt:  runner.LastSeenAt,
		CreatedAt:   runner.CreatedAt,
		UpdatedAt:   runner.UpdatedAt,
	}
}
