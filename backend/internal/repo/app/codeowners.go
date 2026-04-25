package app

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ErrCODEOWNERSNotFound = errors.New("CODEOWNERS file not found")
	ErrCODEOWNERSInvalid = errors.New("CODEOWNERS file is invalid")
)

type CODEOWNERS struct {
	Rules []*CODEOWNERRule
}

type CODEOWNERRule struct {
	Pattern  string
	Owners   []string
	Line     int
	Comment  string
}

type CODEOWNERSService interface {
	ParseCODEOWNERS(repoPath string) (*CODEOWNERS, error)
	MatchPath(owners *CODEOWNERS, filePath string) []string
	GetOwners(repoPath string, filePath string) ([]string, error)
}

type codeownersService struct {}

func NewCODEOWNERSService() CODEOWNERSService {
	return &codeownersService{}
}

func (s *codeownersService) ParseCODEOWNERS(repoPath string) (*CODEOWNERS, error) {
	locations := []string{
		filepath.Join(repoPath, "CODEOWNERS"),
		filepath.Join(repoPath, ".github", "CODEOWNERS"),
		filepath.Join(repoPath, "docs", "CODEOWNERS"),
	}

	var codeownersFile string
	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			codeownersFile = loc
			break
		}
	}

	if codeownersFile == "" {
		return nil, ErrCODEOWNERSNotFound
	}

	file, err := os.Open(codeownersFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rules := []*CODEOWNERRule{}
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		commentIndex := strings.Index(line, "#")
		var comment string
		if commentIndex != -1 {
			comment = strings.TrimSpace(line[commentIndex+1:])
			line = strings.TrimSpace(line[:commentIndex])
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		pattern := parts[0]
		owners := parts[1:]

		rules = append(rules, &CODEOWNERRule{
			Pattern: pattern,
			Owners:  owners,
			Line:    lineNum,
			Comment: comment,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &CODEOWNERS{Rules: rules}, nil
}

func (s *codeownersService) MatchPath(owners *CODEOWNERS, filePath string) []string {
	var matchedOwners []string
	bestMatch := -1

	for i, rule := range owners.Rules {
		if s.matchesPattern(rule.Pattern, filePath) {
			if i > bestMatch {
				bestMatch = i
				matchedOwners = rule.Owners
			}
		}
	}

	return matchedOwners
}

func (s *codeownersService) GetOwners(repoPath string, filePath string) ([]string, error) {
	owners, err := s.ParseCODEOWNERS(repoPath)
	if err != nil {
		if err == ErrCODEOWNERSNotFound {
			return []string{}, nil
		}
		return nil, err
	}

	return s.MatchPath(owners, filePath), nil
}

func (s *codeownersService) matchesPattern(pattern, path string) bool {
	if pattern == "*" {
		return true
	}

	if pattern == path {
		return true
	}

	if strings.HasSuffix(pattern, "/") {
		return strings.HasPrefix(path, pattern)
	}

	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(path, pattern[1:])
	}

	if strings.Contains(pattern, "*") {
		regexPattern := s.globToRegex(pattern)
		match, _ := regexp.MatchString(regexPattern, path)
		return match
	}

	return false
}

func (s *codeownersService) globToRegex(pattern string) string {
	regexPattern := regexp.QuoteMeta(pattern)
	regexPattern = strings.ReplaceAll(regexPattern, "\\*", ".*")
	return "^" + regexPattern + "$"
}
