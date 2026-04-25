package app

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"laima/internal/ai/domain"
)

type FindingSource string

const (
	FindingSourceBugDetection  FindingSource = "bug_detection"
	FindingSourceSecurity      FindingSource = "security"
	FindingSourcePerformance   FindingSource = "performance"
	FindingSourceStyle         FindingSource = "style"
	FindingSourceBestPractice  FindingSource = "best_practice"
)

type AggregatedFinding struct {
	ID           string         `json:"id"`
	Path         string         `json:"path"`
	Line         int            `json:"line"`
	Message      string         `json:"message"`
	Severity     string         `json:"severity"`
	Category     string         `json:"category"`
	Suggestion   string         `json:"suggestion"`
	FixPatch     string         `json:"fix_patch,omitempty"`
	Confidence   float64        `json:"confidence"`
	Sources      []FindingSource `json:"sources"`
	Occurrences  int            `json:"occurrences"`
}

type ResultAggregator interface {
	Aggregate(findings []*domain.AIReviewIssue) ([]*AggregatedFinding, error)
	Deduplicate(findings []*domain.AIReviewIssue) []*domain.AIReviewIssue
	SortByPriority(findings []*AggregatedFinding) []*AggregatedFinding
	GenerateSummary(findings []*AggregatedFinding) string
}

type resultAggregator struct {
	severityOrder map[string]int
}

func NewResultAggregator() ResultAggregator {
	return &resultAggregator{
		severityOrder: map[string]int{
			domain.AIReviewSeverityCritical: 4,
			domain.AIReviewSeverityHigh:     3,
			domain.AIReviewSeverityMedium:   2,
			domain.AIReviewSeverityLow:      1,
		},
	}
}

func (a *resultAggregator) Aggregate(findings []*domain.AIReviewIssue) ([]*AggregatedFinding, error) {
	deduplicated := a.Deduplicate(findings)
	
	grouped := make(map[string][]*domain.AIReviewIssue)
	for _, f := range deduplicated {
		key := a.generateFindingKey(f)
		grouped[key] = append(grouped[key], f)
	}

	var result []*AggregatedFinding
	for key, group := range grouped {
		if len(group) == 0 {
			continue
		}

		primary := group[0]
		sources := a.extractSources(group)
		
		avgConfidence := 0.0
		for _, f := range group {
			avgConfidence += f.Confidence
		}
		avgConfidence /= float64(len(group))

		aggregated := &AggregatedFinding{
			ID:          key,
			Path:        primary.Path,
			Line:        primary.Line,
			Message:     primary.Description,
			Severity:    primary.Severity,
			Category:    primary.Category,
			Suggestion:  primary.Suggestion,
			Confidence:  avgConfidence,
			Sources:     sources,
			Occurrences: len(group),
		}

		result = append(result, aggregated)
	}

	return a.SortByPriority(result), nil
}

func (a *resultAggregator) Deduplicate(findings []*domain.AIReviewIssue) []*domain.AIReviewIssue {
	seen := make(map[string]bool)
	var result []*domain.AIReviewIssue

	for _, f := range findings {
		key := a.generateFindingKey(f)
		if !seen[key] {
			seen[key] = true
			result = append(result, f)
		}
	}

	return result
}

func (a *resultAggregator) SortByPriority(findings []*AggregatedFinding) []*AggregatedFinding {
	sort.Slice(findings, func(i, j int) bool {
		severityI := a.severityOrder[findings[i].Severity]
		severityJ := a.severityOrder[findings[j].Severity]
		
		if severityI != severityJ {
			return severityI > severityJ
		}
		
		if findings[i].Confidence != findings[j].Confidence {
			return findings[i].Confidence > findings[j].Confidence
		}
		
		if findings[i].Occurrences != findings[j].Occurrences {
			return findings[i].Occurrences > findings[j].Occurrences
		}
		
		return findings[i].Path < findings[j].Path
	})

	return findings
}

func (a *resultAggregator) GenerateSummary(findings []*AggregatedFinding) string {
	if len(findings) == 0 {
		return "代码审查完成，未发现问题。"
	}

	severityCounts := make(map[string]int)
	categoryCounts := make(map[string]int)

	for _, f := range findings {
		severityCounts[f.Severity]++
		categoryCounts[f.Category]++
	}

	var summaryParts []string

	if critical := severityCounts[domain.AIReviewSeverityCritical]; critical > 0 {
		summaryParts = append(summaryParts, formatCount(critical, "个严重问题"))
	}
	if high := severityCounts[domain.AIReviewSeverityHigh]; high > 0 {
		summaryParts = append(summaryParts, formatCount(high, "个高危问题"))
	}
	if medium := severityCounts[domain.AIReviewSeverityMedium]; medium > 0 {
		summaryParts = append(summaryParts, formatCount(medium, "个中危问题"))
	}
	if low := severityCounts[domain.AIReviewSeverityLow]; low > 0 {
		summaryParts = append(summaryParts, formatCount(low, "个低危问题"))
	}

	summary := "本次代码审查发现 " + strings.Join(summaryParts, "、") + "。"

	topFindings := findings
	if len(findings) > 3 {
		topFindings = findings[:3]
	}

	if len(topFindings) > 0 {
		summary += "\n\n主要问题："
		for i, f := range topFindings {
			summary += fmt.Sprintf("\n%d. [%s] %s (%s:%d)", 
				i+1, 
				strings.ToUpper(f.Severity),
				f.Message,
				f.Path,
				f.Line,
			)
		}
	}

	return summary
}

func (a *resultAggregator) generateFindingKey(f *domain.AIReviewIssue) string {
	data := fmt.Sprintf("%s:%d:%s", f.Path, f.Line, f.Category)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:16]
}

func (a *resultAggregator) extractSources(findings []*domain.AIReviewIssue) []FindingSource {
	sourceMap := make(map[FindingSource]bool)
	for _, f := range findings {
		source := FindingSource(f.Category)
		sourceMap[source] = true
	}

	var sources []FindingSource
	for source := range sourceMap {
		sources = append(sources, source)
	}

	return sources
}

func formatCount(count int, label string) string {
	return fmt.Sprintf("%d%s", count, label)
}
