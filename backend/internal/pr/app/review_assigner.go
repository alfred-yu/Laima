package app

import (
	"errors"

	"laima/internal/pr/domain"
)

var (
	ErrReviewerNotFound = errors.New("reviewer not found")
	ErrReviewRequired  = errors.New("review is required")
)

type ReviewAssigner interface {
	AssignReviewers(prID int, filePath string) ([]int, error)
	CheckReviewRequirements(prID int) (bool, error)
	GetRequiredReviewers(prID int) ([]int, error)
}

type reviewAssigner struct {
	codeownersService CODEOWNERSService
}

func NewReviewAssigner(codeownersService CODEOWNERSService) ReviewAssigner {
	return &reviewAssigner{
		codeownersService: codeownersService,
	}
}

func (a *reviewAssigner) AssignReviewers(prID int, filePath string) ([]int, error) {
	repoPath := "/path/to/repo"

	owners, err := a.codeownersService.GetOwners(repoPath, filePath)
	if err != nil {
		return nil, err
	}

	reviewerIDs := []int{}
	for _, owner := range owners {
		reviewerID := a.getReviewerID(owner)
		if reviewerID > 0 {
			reviewerIDs = append(reviewerIDs, reviewerID)
		}
	}

	return reviewerIDs, nil
}

func (a *reviewAssigner) CheckReviewRequirements(prID int) (bool, error) {
	requiredReviewers, err := a.GetRequiredReviewers(prID)
	if err != nil {
		return false, err
	}

	return len(requiredReviewers) > 0, nil
}

func (a *reviewAssigner) GetRequiredReviewers(prID int) ([]int, error) {
	pr := &domain.PullRequest{
		ID:           prID,
		RepositoryID: 1,
		Title:        "Test PR",
		Description:  "Test PR description",
	}

	files := []string{
		"src/main.go",
		"src/utils.go",
	}

	reviewerMap := make(map[int]bool)
	for _, file := range files {
		reviewers, err := a.AssignReviewers(prID, file)
		if err != nil {
			continue
		}

		for _, reviewer := range reviewers {
			reviewerMap[reviewer] = true
		}
	}

	requiredReviewers := []int{}
	for reviewer := range reviewerMap {
		requiredReviewers = append(requiredReviewers, reviewer)
	}

	return requiredReviewers, nil
}

func (a *reviewAssigner) getReviewerID(owner string) int {
	ownerMap := map[string]int{
		"@alice": 1,
		"@bob":   2,
		"@charlie": 3,
	}

	return ownerMap[owner]
}
