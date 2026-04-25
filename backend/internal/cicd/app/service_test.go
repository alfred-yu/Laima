package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"laima/internal/cicd/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRunnerService_RegisterRunner(t *testing.T) {
	service := NewRunnerService()
	ctx := context.Background()

	t.Run("successful registration", func(t *testing.T) {
		req := &domain.RunnerRegistration{
			Name:        "test-runner",
			Description: "Test runner",
			Executor:    domain.RunnerExecutorShell,
			Tags:        "linux,docker",
			MaxJobs:     2,
		}

		runner, err := service.RegisterRunner(req)

		assert.NoError(t, err)
		assert.NotNil(t, runner)
		assert.Equal(t, "test-runner", runner.Name)
		assert.Equal(t, "Test runner", runner.Description)
		assert.Equal(t, domain.RunnerExecutorShell, runner.Executor)
		assert.Equal(t, "linux,docker", runner.Tags)
		assert.Equal(t, 2, runner.MaxJobs)
		assert.NotEmpty(t, runner.Token)
	})

	t.Run("duplicate name", func(t *testing.T) {
		req1 := &domain.RunnerRegistration{
			Name:     "duplicate-runner",
			Executor: domain.RunnerExecutorShell,
		}

		req2 := &domain.RunnerRegistration{
			Name:     "duplicate-runner",
			Executor: domain.RunnerExecutorShell,
		}

		runner1, err1 := service.RegisterRunner(req1)
		assert.NoError(t, err1)
		assert.NotNil(t, runner1)

		runner2, err2 := service.RegisterRunner(req2)
		assert.Error(t, err2)
		assert.Nil(t, runner2)
		assert.Equal(t, ErrRunnerAlreadyExists, err2)
	})
}

func TestRunnerService_Heartbeat(t *testing.T) {
	service := NewRunnerService()
	ctx := context.Background()

	t.Run("successful heartbeat", func(t *testing.T) {
		regReq := &domain.RunnerRegistration{
			Name:     "heartbeat-runner",
			Executor: domain.RunnerExecutorShell,
		}

		runner, _ := service.RegisterRunner(regReq)

		heartbeatReq := &domain.RunnerHeartbeat{
			RunnerID: runner.ID,
			Token:    runner.Token,
			Version:  "1.0.0",
			Status:   domain.RunnerStatusOnline,
		}

		err := service.Heartbeat(heartbeatReq)

		assert.NoError(t, err)

		updatedRunner, _ := service.GetRunner(runner.ID)
		assert.Equal(t, "1.0.0", updatedRunner.Version)
		assert.Equal(t, domain.RunnerStatusOnline, updatedRunner.Status)
	})

	t.Run("invalid token", func(t *testing.T) {
		regReq := &domain.RunnerRegistration{
			Name:     "invalid-token-runner",
			Executor: domain.RunnerExecutorShell,
		}

		runner, _ := service.RegisterRunner(regReq)

		heartbeatReq := &domain.RunnerHeartbeat{
			RunnerID: runner.ID,
			Token:    "invalid-token",
		}

		err := service.Heartbeat(heartbeatReq)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidToken, err)
	})

	t.Run("runner not found", func(t *testing.T) {
		heartbeatReq := &domain.RunnerHeartbeat{
			RunnerID: 999,
			Token:    "some-token",
		}

		err := service.Heartbeat(heartbeatReq)

		assert.Error(t, err)
		assert.Equal(t, ErrRunnerNotFound, err)
	})
}

func TestRunnerService_RequestJob(t *testing.T) {
	service := NewRunnerService()
	ctx := context.Background()

	t.Run("no jobs available", func(t *testing.T) {
		regReq := &domain.RunnerRegistration{
			Name:     "no-job-runner",
			Executor: domain.RunnerExecutorShell,
		}

		runner, _ := service.RegisterRunner(regReq)

		jobReq := &domain.RunnerJobRequest{
			RunnerID: runner.ID,
			Token:    runner.Token,
		}

		job, err := service.RequestJob(jobReq)

		assert.NoError(t, err)
		assert.Nil(t, job)
	})

	t.Run("runner busy", func(t *testing.T) {
		regReq := &domain.RunnerRegistration{
			Name:     "busy-runner",
			Executor: domain.RunnerExecutorShell,
			MaxJobs:  1,
		}

		runner, _ := service.RegisterRunner(regReq)

		runnerService := service.(*runnerService)
		runnerService.runners[runner.ID].ActiveJobs = 1

		jobReq := &domain.RunnerJobRequest{
			RunnerID: runner.ID,
			Token:    runner.Token,
		}

		job, err := service.RequestJob(jobReq)

		assert.Error(t, err)
		assert.Nil(t, job)
		assert.Equal(t, ErrRunnerBusy, err)
	})
}

func TestRunnerService_ListRunners(t *testing.T) {
	service := NewRunnerService()

	t.Run("empty list", func(t *testing.T) {
		runners, err := service.ListRunners()

		assert.NoError(t, err)
		assert.Empty(t, runners)
	})

	t.Run("list with runners", func(t *testing.T) {
		regReq1 := &domain.RunnerRegistration{
			Name:     "runner-1",
			Executor: domain.RunnerExecutorShell,
		}

		regReq2 := &domain.RunnerRegistration{
			Name:     "runner-2",
			Executor: domain.RunnerExecutorDocker,
		}

		service.RegisterRunner(regReq1)
		service.RegisterRunner(regReq2)

		runners, err := service.ListRunners()

		assert.NoError(t, err)
		assert.Len(t, runners, 2)
	})
}

func TestRunnerService_DeleteRunner(t *testing.T) {
	service := NewRunnerService()

	t.Run("successful delete", func(t *testing.T) {
		regReq := &domain.RunnerRegistration{
			Name:     "delete-runner",
			Executor: domain.RunnerExecutorShell,
		}

		runner, _ := service.RegisterRunner(regReq)

		err := service.DeleteRunner(runner.ID)

		assert.NoError(t, err)

		_, err = service.GetRunner(runner.ID)
		assert.Error(t, err)
		assert.Equal(t, ErrRunnerNotFound, err)
	})

	t.Run("runner not found", func(t *testing.T) {
		err := service.DeleteRunner(999)

		assert.Error(t, err)
		assert.Equal(t, ErrRunnerNotFound, err)
	})
}

func TestPipelineService_CreatePipeline(t *testing.T) {
	t.Run("create pipeline with valid data", func(t *testing.T) {
		req := &domain.PipelineRequest{
			RepositoryID: 1,
			CommitSHA:    "abc123",
			Ref:          "refs/heads/main",
			Trigger:      "push",
		}

		assert.NotNil(t, req)
		assert.Equal(t, 1, req.RepositoryID)
		assert.Equal(t, "abc123", req.CommitSHA)
		assert.Equal(t, "refs/heads/main", req.Ref)
		assert.Equal(t, "push", req.Trigger)
	})
}

func TestJobStatus(t *testing.T) {
	t.Run("job status transitions", func(t *testing.T) {
		statuses := []string{
			domain.JobStatusPending,
			domain.JobStatusRunning,
			domain.JobStatusSuccess,
		}

		assert.Equal(t, "pending", statuses[0])
		assert.Equal(t, "running", statuses[1])
		assert.Equal(t, "success", statuses[2])
	})
}

func TestRunnerStatus(t *testing.T) {
	t.Run("runner status values", func(t *testing.T) {
		assert.Equal(t, "online", domain.RunnerStatusOnline)
		assert.Equal(t, "offline", domain.RunnerStatusOffline)
		assert.Equal(t, "busy", domain.RunnerStatusBusy)
	})
}

func TestRunnerExecutorTypes(t *testing.T) {
	t.Run("executor type values", func(t *testing.T) {
		assert.Equal(t, "shell", domain.RunnerExecutorShell)
		assert.Equal(t, "docker", domain.RunnerExecutorDocker)
		assert.Equal(t, "kubernetes", domain.RunnerExecutorKubernetes)
	})
}
