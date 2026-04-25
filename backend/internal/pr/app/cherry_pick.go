package app

import (
	"errors"
	"time"

	"laima/internal/pr/domain"
)

var (
	ErrCommitNotFound     = errors.New("commit not found")
	ErrBranchExists       = errors.New("branch already exists")
	ErrCherryPickFailed   = errors.New("cherry pick failed")
	ErrRevertFailed       = errors.New("revert failed")
)

type CherryPickService interface {
	CherryPick(commitSHA string, targetBranch string) (string, error)
	CreateCherryPickBranch(branchName string, baseBranch string) (string, error)
	ResolveConflicts(branch string) error
	CreateCherryPickPR(title string, sourceBranch string, targetBranch string) (*domain.PullRequest, error)
}

type RevertService interface {
	RevertCommit(commitSHA string) (string, error)
	CreateRevertPR(commitSHA string, targetBranch string) (*domain.PullRequest, error)
}

type cherryPickService struct {}

func NewCherryPickService() CherryPickService {
	return &cherryPickService{}
}

func (s *cherryPickService) CherryPick(commitSHA string, targetBranch string) (string, error) {
	if commitSHA == "" {
		return "", errors.New("commit SHA is required")
	}

	if targetBranch == "" {
		return "", errors.New("target branch is required")
	}

	return "new-commit-sha", nil
}

func (s *cherryPickService) CreateCherryPickBranch(branchName string, baseBranch string) (string, error) {
	if branchName == "" {
		return "", errors.New("branch name is required")
	}

	if baseBranch == "" {
		return "", errors.New("base branch is required")
	}

	return branchName, nil
}

func (s *cherryPickService) ResolveConflicts(branch string) error {
	if branch == "" {
		return errors.New("branch is required")
	}

	return nil
}

func (s *cherryPickService) CreateCherryPickPR(title string, sourceBranch string, targetBranch string) (*domain.PullRequest, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	if sourceBranch == "" || targetBranch == "" {
		return nil, errors.New("source and target branches are required")
	}

	pr := &domain.PullRequest{
		ID:           1,
		Number:       42,
		Title:        title,
		Description:  "Cherry-pick PR",
		RepositoryID: 1,
		SourceRepoID: 1,
		SourceBranch: sourceBranch,
		TargetBranch: targetBranch,
		AuthorID:     1,
		Status:       domain.PRStatusOpen,
		Mergeable:    true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return pr, nil
}

type revertService struct {}

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
