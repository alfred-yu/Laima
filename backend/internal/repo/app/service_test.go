package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"laima/internal/repo/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepoRepository struct {
	mock.Mock
}

func (m *MockRepoRepository) Create(repo *domain.Repository) error {
	args := m.Called(repo)
	return args.Error(0)
}

func (m *MockRepoRepository) GetByID(id int64) (*domain.Repository, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Repository), args.Error(1)
}

func (m *MockRepoRepository) GetByFullPath(fullPath string) (*domain.Repository, error) {
	args := m.Called(fullPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Repository), args.Error(1)
}

func (m *MockRepoRepository) Update(repo *domain.Repository) error {
	args := m.Called(repo)
	return args.Error(0)
}

func (m *MockRepoRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepoRepository) List(filter *RepoFilter) ([]*domain.Repository, int64, error) {
	args := m.Called(filter)
	return args.Get(0).([]*domain.Repository), args.Get(1).(int64), args.Error(2)
}

type MockGitService struct {
	mock.Mock
}

func (m *MockGitService) InitRepository(repoPath string) error {
	args := m.Called(repoPath)
	return args.Error(0)
}

func (m *MockGitService) CreateBranch(repoPath, branchName, sourceRef string) error {
	args := m.Called(repoPath, branchName, sourceRef)
	return args.Error(0)
}

func (m *MockGitService) ListBranches(repoPath string) ([]*domain.Branch, error) {
	args := m.Called(repoPath)
	return args.Get(0).([]*domain.Branch), args.Error(1)
}

func (m *MockGitService) CloneRepository(cloneURL, repoPath string) error {
	args := m.Called(cloneURL, repoPath)
	return args.Error(0)
}

type repoService struct {
	repoRepo   *MockRepoRepository
	gitService *MockGitService
}

func NewTestRepoService() *repoService {
	return &repoService{
		repoRepo:   new(MockRepoRepository),
		gitService: new(MockGitService),
	}
}

func (s *repoService) CreateRepo(ctx context.Context, req *CreateRepoRequest) (*domain.Repository, error) {
	if req.Name == "" {
		return nil, errors.New("repository name is required")
	}

	repo := &domain.Repository{
		Name:          req.Name,
		FullPath:      req.Name,
		Description:   req.Description,
		OwnerType:     req.OwnerType,
		OwnerID:       req.OwnerID,
		Visibility:    domain.Visibility(req.Visibility),
		DefaultBranch: "main",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repoRepo.Create(repo); err != nil {
		return nil, err
	}

	repoPath := req.Name
	if err := s.gitService.InitRepository(repoPath); err != nil {
		return nil, err
	}

	return repo, nil
}

func (s *repoService) GetRepo(ctx context.Context, repoID int64) (*domain.Repository, error) {
	return s.repoRepo.GetByID(repoID)
}

func (s *repoService) UpdateRepo(ctx context.Context, repoID int64, req *UpdateRepoRequest) (*domain.Repository, error) {
	repo, err := s.repoRepo.GetByID(repoID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		repo.Name = req.Name
	}
	if req.Description != "" {
		repo.Description = req.Description
	}
	if req.Visibility != "" {
		repo.Visibility = domain.Visibility(req.Visibility)
	}
	repo.UpdatedAt = time.Now()

	if err := s.repoRepo.Update(repo); err != nil {
		return nil, err
	}

	return repo, nil
}

func (s *repoService) DeleteRepo(ctx context.Context, repoID int64) error {
	return s.repoRepo.Delete(repoID)
}

func (s *repoService) ForkRepo(ctx context.Context, repoID int64, targetNamespace string) (*domain.Repository, error) {
	originalRepo, err := s.repoRepo.GetByID(repoID)
	if err != nil {
		return nil, err
	}

	forkedRepo := &domain.Repository{
		Name:          originalRepo.Name,
		FullPath:      targetNamespace + "/" + originalRepo.Name,
		Description:   originalRepo.Description,
		OwnerType:     domain.OwnerTypeUser,
		OwnerID:       1,
		Visibility:    originalRepo.Visibility,
		DefaultBranch: originalRepo.DefaultBranch,
		IsFork:        true,
		ForkParentID:  &repoID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repoRepo.Create(forkedRepo); err != nil {
		return nil, err
	}

	return forkedRepo, nil
}

func (s *repoService) ListBranches(ctx context.Context, repoID int64) ([]*domain.Branch, error) {
	repo, err := s.repoRepo.GetByID(repoID)
	if err != nil {
		return nil, err
	}

	return s.gitService.ListBranches(repo.FullPath)
}

func TestCreateRepo(t *testing.T) {
	service := NewTestRepoService()
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		req := &CreateRepoRequest{
			Name:        "test-repo",
			Description: "Test repository",
			OwnerType:   domain.OwnerTypeUser,
			OwnerID:     1,
			Visibility:  "public",
		}

		service.repoRepo.On("Create", mock.AnythingOfType("*domain.Repository")).Return(nil)
		service.gitService.On("InitRepository", "test-repo").Return(nil)

		repo, err := service.CreateRepo(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, repo)
		assert.Equal(t, "test-repo", repo.Name)
		assert.Equal(t, "Test repository", repo.Description)
		assert.Equal(t, domain.OwnerTypeUser, repo.OwnerType)
		assert.Equal(t, domain.Visibility("public"), repo.Visibility)
		assert.Equal(t, "main", repo.DefaultBranch)

		service.repoRepo.AssertExpectations(t)
		service.gitService.AssertExpectations(t)
	})

	t.Run("empty name", func(t *testing.T) {
		req := &CreateRepoRequest{
			Description: "Test repository",
			OwnerType:   domain.OwnerTypeUser,
			OwnerID:     1,
		}

		repo, err := service.CreateRepo(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, repo)
		assert.Equal(t, "repository name is required", err.Error())
	})

	t.Run("database error", func(t *testing.T) {
		req := &CreateRepoRequest{
			Name:        "test-repo",
			Description: "Test repository",
			OwnerType:   domain.OwnerTypeUser,
			OwnerID:     1,
		}

		service.repoRepo.On("Create", mock.AnythingOfType("*domain.Repository")).
			Return(errors.New("database error"))

		repo, err := service.CreateRepo(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, repo)
		assert.Equal(t, "database error", err.Error())

		service.repoRepo.AssertExpectations(t)
	})
}

func TestGetRepo(t *testing.T) {
	service := NewTestRepoService()
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		expectedRepo := &domain.Repository{
			ID:          1,
			Name:        "test-repo",
			FullPath:    "test-repo",
			Description: "Test repository",
			OwnerType:   domain.OwnerTypeUser,
			OwnerID:     1,
			Visibility:  domain.VisibilityPublic,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		service.repoRepo.On("GetByID", int64(1)).Return(expectedRepo, nil)

		repo, err := service.GetRepo(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, repo)
		assert.Equal(t, int64(1), repo.ID)
		assert.Equal(t, "test-repo", repo.Name)

		service.repoRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		service.repoRepo.On("GetByID", int64(999)).Return(nil, errors.New("repository not found"))

		repo, err := service.GetRepo(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, repo)
		assert.Equal(t, "repository not found", err.Error())

		service.repoRepo.AssertExpectations(t)
	})
}

func TestUpdateRepo(t *testing.T) {
	service := NewTestRepoService()
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		existingRepo := &domain.Repository{
			ID:          1,
			Name:        "old-name",
			Description: "Old description",
			OwnerType:   domain.OwnerTypeUser,
			OwnerID:     1,
			Visibility:  domain.VisibilityPublic,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		service.repoRepo.On("GetByID", int64(1)).Return(existingRepo, nil)
		service.repoRepo.On("Update", mock.AnythingOfType("*domain.Repository")).Return(nil)

		req := &UpdateRepoRequest{
			Name:        "new-name",
			Description: "New description",
			Visibility:  "private",
		}

		repo, err := service.UpdateRepo(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, repo)
		assert.Equal(t, "new-name", repo.Name)
		assert.Equal(t, "New description", repo.Description)
		assert.Equal(t, domain.Visibility("private"), repo.Visibility)

		service.repoRepo.AssertExpectations(t)
	})

	t.Run("repo not found", func(t *testing.T) {
		service.repoRepo.On("GetByID", int64(999)).Return(nil, errors.New("repository not found"))

		req := &UpdateRepoRequest{
			Name: "new-name",
		}

		repo, err := service.UpdateRepo(ctx, 999, req)

		assert.Error(t, err)
		assert.Nil(t, repo)

		service.repoRepo.AssertExpectations(t)
	})
}

func TestDeleteRepo(t *testing.T) {
	service := NewTestRepoService()
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		service.repoRepo.On("Delete", int64(1)).Return(nil)

		err := service.DeleteRepo(ctx, 1)

		assert.NoError(t, err)
		service.repoRepo.AssertExpectations(t)
	})

	t.Run("delete error", func(t *testing.T) {
		service.repoRepo.On("Delete", int64(999)).Return(errors.New("repository not found"))

		err := service.DeleteRepo(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, "repository not found", err.Error())

		service.repoRepo.AssertExpectations(t)
	})
}

func TestForkRepo(t *testing.T) {
	service := NewTestRepoService()
	ctx := context.Background()

	t.Run("successful fork", func(t *testing.T) {
		originalRepo := &domain.Repository{
			ID:          1,
			Name:        "original-repo",
			FullPath:    "original-repo",
			Description: "Original repository",
			OwnerType:   domain.OwnerTypeUser,
			OwnerID:     1,
			Visibility:  domain.VisibilityPublic,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		service.repoRepo.On("GetByID", int64(1)).Return(originalRepo, nil)
		service.repoRepo.On("Create", mock.AnythingOfType("*domain.Repository")).Return(nil)

		forkedRepo, err := service.ForkRepo(ctx, 1, "my-namespace")

		assert.NoError(t, err)
		assert.NotNil(t, forkedRepo)
		assert.Equal(t, "original-repo", forkedRepo.Name)
		assert.Equal(t, "my-namespace/original-repo", forkedRepo.FullPath)
		assert.True(t, forkedRepo.IsFork)
		assert.Equal(t, int64(1), *forkedRepo.ForkParentID)

		service.repoRepo.AssertExpectations(t)
	})

	t.Run("original repo not found", func(t *testing.T) {
		service.repoRepo.On("GetByID", int64(999)).Return(nil, errors.New("repository not found"))

		forkedRepo, err := service.ForkRepo(ctx, 999, "my-namespace")

		assert.Error(t, err)
		assert.Nil(t, forkedRepo)

		service.repoRepo.AssertExpectations(t)
	})
}

func TestListBranches(t *testing.T) {
	service := NewTestRepoService()
	ctx := context.Background()

	t.Run("successful list", func(t *testing.T) {
		repo := &domain.Repository{
			ID:        1,
			Name:      "test-repo",
			FullPath:  "test-repo",
			OwnerType: domain.OwnerTypeUser,
			OwnerID:   1,
		}

		expectedBranches := []*domain.Branch{
			{Name: "main", IsDefault: true},
			{Name: "develop", IsDefault: false},
			{Name: "feature/test", IsDefault: false},
		}

		service.repoRepo.On("GetByID", int64(1)).Return(repo, nil)
		service.gitService.On("ListBranches", "test-repo").Return(expectedBranches, nil)

		branches, err := service.ListBranches(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, branches)
		assert.Len(t, branches, 3)
		assert.Equal(t, "main", branches[0].Name)
		assert.True(t, branches[0].IsDefault)

		service.repoRepo.AssertExpectations(t)
		service.gitService.AssertExpectations(t)
	})

	t.Run("repo not found", func(t *testing.T) {
		service.repoRepo.On("GetByID", int64(999)).Return(nil, errors.New("repository not found"))

		branches, err := service.ListBranches(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, branches)

		service.repoRepo.AssertExpectations(t)
	})
}
