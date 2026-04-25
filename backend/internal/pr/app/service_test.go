package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"laima/internal/pr/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockPRRepository struct {
	mock.Mock
}

func (m *MockPRRepository) Create(pr *domain.PullRequest) error {
	args := m.Called(pr)
	return args.Error(0)
}

func (m *MockPRRepository) GetByID(id int) (*domain.PullRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PullRequest), args.Error(1)
}

func (m *MockPRRepository) GetByNumber(repoID, number int) (*domain.PullRequest, error) {
	args := m.Called(repoID, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PullRequest), args.Error(1)
}

func (m *MockPRRepository) Update(pr *domain.PullRequest) error {
	args := m.Called(pr)
	return args.Error(0)
}

func (m *MockPRRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPRRepository) List(filter *domain.PRFilter) ([]*domain.PullRequest, int64, error) {
	args := m.Called(filter)
	return args.Get(0).([]*domain.PullRequest), args.Get(1).(int64), args.Error(2)
}

type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) Create(review *domain.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockReviewRepository) GetByPRID(prID int) ([]*domain.Review, error) {
	args := m.Called(prID)
	return args.Get(0).([]*domain.Review), args.Error(1)
}

type prTestService struct {
	prRepo      *MockPRRepository
	reviewRepo  *MockReviewRepository
}

func NewTestPRService() *prTestService {
	return &prTestService{
		prRepo:     new(MockPRRepository),
		reviewRepo: new(MockReviewRepository),
	}
}

func (s *prTestService) CreatePR(ctx context.Context, req *domain.CreatePRRequest, authorID int) (*domain.PullRequest, error) {
	if req.Title == "" {
		return nil, errors.New("PR title is required")
	}

	if req.SourceBranch == "" || req.TargetBranch == "" {
		return nil, errors.New("source and target branches are required")
	}

	pr := &domain.PullRequest{
		ID:             1,
		Number:         1,
		Title:          req.Title,
		Description:    req.Description,
		RepositoryID:   req.RepositoryID,
		SourceRepoID:   req.SourceRepoID,
		SourceBranch:   req.SourceBranch,
		TargetBranch:   req.TargetBranch,
		AuthorID:       authorID,
		Status:         domain.PRStatusOpen,
		Mergeable:      true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.prRepo.Create(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *prTestService) GetPR(ctx context.Context, prID int) (*domain.PullRequest, error) {
	return s.prRepo.GetByID(prID)
}

func (s *prTestService) UpdatePR(ctx context.Context, prID int, req *domain.UpdatePRRequest) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetByID(prID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		pr.Title = req.Title
	}
	if req.Description != "" {
		pr.Description = req.Description
	}
	pr.UpdatedAt = time.Now()

	if err := s.prRepo.Update(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *prTestService) MergePR(ctx context.Context, prID int, userID int, mergeStrategy string) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetByID(prID)
	if err != nil {
		return nil, err
	}

	if pr.Status != domain.PRStatusOpen {
		return nil, errors.New("PR is not open")
	}

	if !pr.Mergeable {
		return nil, errors.New("PR is not mergeable")
	}

	pr.Status = domain.PRStatusMerged
	pr.MergedAt = time.Now()
	pr.MergedBy = userID
	pr.MergeStrategy = mergeStrategy
	pr.UpdatedAt = time.Now()

	if err := s.prRepo.Update(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *prTestService) ClosePR(ctx context.Context, prID int, userID int) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetByID(prID)
	if err != nil {
		return nil, err
	}

	if pr.Status != domain.PRStatusOpen {
		return nil, errors.New("PR is not open")
	}

	pr.Status = domain.PRStatusClosed
	pr.ClosedAt = time.Now()
	pr.ClosedBy = userID
	pr.UpdatedAt = time.Now()

	if err := s.prRepo.Update(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *prTestService) CreateReview(ctx context.Context, prID int, req *domain.ReviewRequest, reviewerID int) (*domain.Review, error) {
	pr, err := s.prRepo.GetByID(prID)
	if err != nil {
		return nil, err
	}

	if pr.Status != domain.PRStatusOpen {
		return nil, errors.New("cannot review a closed PR")
	}

	review := &domain.Review{
		ID:         1,
		PRID:       prID,
		ReviewerID: reviewerID,
		Body:       req.Body,
		State:      req.State,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *prTestService) GetReviews(ctx context.Context, prID int) ([]*domain.Review, error) {
	return s.reviewRepo.GetByPRID(prID)
}

func TestCreatePR(t *testing.T) {
	service := NewTestPRService()
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		req := &domain.CreatePRRequest{
			Title:        "Add new feature",
			Description:  "This PR adds a new feature",
			RepositoryID: 1,
			SourceRepoID: 1,
			SourceBranch: "feature/new-feature",
			TargetBranch: "main",
		}

		service.prRepo.On("Create", mock.AnythingOfType("*domain.PullRequest")).Return(nil)

		pr, err := service.CreatePR(ctx, req, 1)

		assert.NoError(t, err)
		assert.NotNil(t, pr)
		assert.Equal(t, "Add new feature", pr.Title)
		assert.Equal(t, "This PR adds a new feature", pr.Description)
		assert.Equal(t, domain.PRStatusOpen, pr.Status)
		assert.True(t, pr.Mergeable)
		assert.Equal(t, 1, pr.AuthorID)

		service.prRepo.AssertExpectations(t)
	})

	t.Run("missing title", func(t *testing.T) {
		req := &domain.CreatePRRequest{
			Description:  "This PR adds a new feature",
			RepositoryID: 1,
			SourceBranch: "feature/new-feature",
			TargetBranch: "main",
		}

		pr, err := service.CreatePR(ctx, req, 1)

		assert.Error(t, err)
		assert.Nil(t, pr)
		assert.Equal(t, "PR title is required", err.Error())
	})

	t.Run("missing branches", func(t *testing.T) {
		req := &domain.CreatePRRequest{
			Title:        "Add new feature",
			Description:  "This PR adds a new feature",
			RepositoryID: 1,
		}

		pr, err := service.CreatePR(ctx, req, 1)

		assert.Error(t, err)
		assert.Nil(t, pr)
		assert.Equal(t, "source and target branches are required", err.Error())
	})
}

func TestUpdatePR(t *testing.T) {
	service := NewTestPRService()
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		existingPR := &domain.PullRequest{
			ID:           1,
			Number:       1,
			Title:        "Old title",
			Description:  "Old description",
			RepositoryID: 1,
			AuthorID:     1,
			Status:       domain.PRStatusOpen,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		service.prRepo.On("GetByID", 1).Return(existingPR, nil)
		service.prRepo.On("Update", mock.AnythingOfType("*domain.PullRequest")).Return(nil)

		req := &domain.UpdatePRRequest{
			Title:       "New title",
			Description: "New description",
		}

		pr, err := service.UpdatePR(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, pr)
		assert.Equal(t, "New title", pr.Title)
		assert.Equal(t, "New description", pr.Description)

		service.prRepo.AssertExpectations(t)
	})

	t.Run("PR not found", func(t *testing.T) {
		service.prRepo.On("GetByID", 999).Return(nil, gorm.ErrRecordNotFound)

		req := &domain.UpdatePRRequest{
			Title: "New title",
		}

		pr, err := service.UpdatePR(ctx, 999, req)

		assert.Error(t, err)
		assert.Nil(t, pr)

		service.prRepo.AssertExpectations(t)
	})
}

func TestMergePR(t *testing.T) {
	service := NewTestPRService()
	ctx := context.Background()

	t.Run("successful merge", func(t *testing.T) {
		openPR := &domain.PullRequest{
			ID:           1,
			Number:       1,
			Title:        "Add new feature",
			RepositoryID: 1,
			AuthorID:     1,
			Status:       domain.PRStatusOpen,
			Mergeable:    true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		service.prRepo.On("GetByID", 1).Return(openPR, nil)
		service.prRepo.On("Update", mock.AnythingOfType("*domain.PullRequest")).Return(nil)

		pr, err := service.MergePR(ctx, 1, 2, "merge")

		assert.NoError(t, err)
		assert.NotNil(t, pr)
		assert.Equal(t, domain.PRStatusMerged, pr.Status)
		assert.NotNil(t, pr.MergedAt)
		assert.Equal(t, 2, pr.MergedBy)
		assert.Equal(t, "merge", pr.MergeStrategy)

		service.prRepo.AssertExpectations(t)
	})

	t.Run("PR not open", func(t *testing.T) {
		closedPR := &domain.PullRequest{
			ID:           1,
			Number:       1,
			Title:        "Add new feature",
			RepositoryID: 1,
			AuthorID:     1,
			Status:       domain.PRStatusClosed,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		service.prRepo.On("GetByID", 1).Return(closedPR, nil)

		pr, err := service.MergePR(ctx, 1, 2, "merge")

		assert.Error(t, err)
		assert.Nil(t, pr)
		assert.Equal(t, "PR is not open", err.Error())

		service.prRepo.AssertExpectations(t)
	})

	t.Run("PR not mergeable", func(t *testing.T) {
		unmergeablePR := &domain.PullRequest{
			ID:           1,
			Number:       1,
			Title:        "Add new feature",
			RepositoryID: 1,
			AuthorID:     1,
			Status:       domain.PRStatusOpen,
			Mergeable:    false,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		service.prRepo.On("GetByID", 1).Return(unmergeablePR, nil)

		pr, err := service.MergePR(ctx, 1, 2, "merge")

		assert.Error(t, err)
		assert.Nil(t, pr)
		assert.Equal(t, "PR is not mergeable", err.Error())

		service.prRepo.AssertExpectations(t)
	})
}

func TestClosePR(t *testing.T) {
	service := NewTestPRService()
	ctx := context.Background()

	t.Run("successful close", func(t *testing.T) {
		openPR := &domain.PullRequest{
			ID:           1,
			Number:       1,
			Title:        "Add new feature",
			RepositoryID: 1,
			AuthorID:     1,
			Status:       domain.PRStatusOpen,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		service.prRepo.On("GetByID", 1).Return(openPR, nil)
		service.prRepo.On("Update", mock.AnythingOfType("*domain.PullRequest")).Return(nil)

		pr, err := service.ClosePR(ctx, 1, 2)

		assert.NoError(t, err)
		assert.NotNil(t, pr)
		assert.Equal(t, domain.PRStatusClosed, pr.Status)
		assert.NotNil(t, pr.ClosedAt)
		assert.Equal(t, 2, pr.ClosedBy)

		service.prRepo.AssertExpectations(t)
	})

	t.Run("PR not open", func(t *testing.T) {
		mergedPR := &domain.PullRequest{
			ID:           1,
			Number:       1,
			Title:        "Add new feature",
			RepositoryID: 1,
			AuthorID:     1,
			Status:       domain.PRStatusMerged,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		service.prRepo.On("GetByID", 1).Return(mergedPR, nil)

		pr, err := service.ClosePR(ctx, 1, 2)

		assert.Error(t, err)
		assert.Nil(t, pr)
		assert.Equal(t, "PR is not open", err.Error())

		service.prRepo.AssertExpectations(t)
	})
}

func TestCreateReview(t *testing.T) {
	service := NewTestPRService()
	ctx := context.Background()

	t.Run("successful review", func(t *testing.T) {
		openPR := &domain.PullRequest{
			ID:           1,
			Number:       1,
			Title:        "Add new feature",
			RepositoryID: 1,
			AuthorID:     1,
			Status:       domain.PRStatusOpen,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		service.prRepo.On("GetByID", 1).Return(openPR, nil)
		service.reviewRepo.On("Create", mock.AnythingOfType("*domain.Review")).Return(nil)

		req := &domain.ReviewRequest{
			Body:  "Looks good to me!",
			State: domain.ReviewStateApproved,
		}

		review, err := service.CreateReview(ctx, 1, req, 2)

		assert.NoError(t, err)
		assert.NotNil(t, review)
		assert.Equal(t, "Looks good to me!", review.Body)
		assert.Equal(t, domain.ReviewStateApproved, review.State)
		assert.Equal(t, 2, review.ReviewerID)

		service.prRepo.AssertExpectations(t)
		service.reviewRepo.AssertExpectations(t)
	})

	t.Run("review closed PR", func(t *testing.T) {
		closedPR := &domain.PullRequest{
			ID:           1,
			Number:       1,
			Title:        "Add new feature",
			RepositoryID: 1,
			AuthorID:     1,
			Status:       domain.PRStatusClosed,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		service.prRepo.On("GetByID", 1).Return(closedPR, nil)

		req := &domain.ReviewRequest{
			Body:  "Looks good to me!",
			State: domain.ReviewStateApproved,
		}

		review, err := service.CreateReview(ctx, 1, req, 2)

		assert.Error(t, err)
		assert.Nil(t, review)
		assert.Equal(t, "cannot review a closed PR", err.Error())

		service.prRepo.AssertExpectations(t)
	})
}
