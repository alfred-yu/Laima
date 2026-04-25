package app

import (
	"errors"
	"time"

	"laima/internal/pr/domain"
)

type RevertService interface {
	RevertCommit(commitSHA string) (string, error)
	CreateRevertPR(commitSHA string, targetBranch string) (*domain.PullRequest, error)
}

type revertService struct{}

func NewRevertService() RevertService {
	return &revertService{}
}

func (s *revertService) RevertCommit(commitSHA string) (string, error) {
	if commitSHA == "" {
		return "", errors.New("commit SHA is required")
	}

	return "revert-commit-sha", nil
}

func (s *revertService) CreateRevertPR(commitSHA string, targetBranch string) (*domain.PullRequest, error) {
	if commitSHA == "" {
		return nil, errors.New("commit SHA is required")
	}

	if targetBranch == "" {
		return nil, errors.New("target branch is required")
	}

	pr := &domain.PullRequest{
		ID:           2,
		Number:       43,
		Title:        "Revert commit " + commitSHA,
		Description:  "This PR reverts commit " + commitSHA,
		RepositoryID: 1,
		SourceRepoID: 1,
		SourceBranch: "revert-" + commitSHA,
		TargetBranch: targetBranch,
		AuthorID:     1,
		Status:       domain.PRStatusOpen,
		Mergeable:    true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return pr, nil
}
