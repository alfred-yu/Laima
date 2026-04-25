package domain

import "time"

type SignaturePolicy struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	RepositoryID    int64     `json:"repository_id" gorm:"not null;unique;index"`
	RequireSignedCommits bool   `json:"require_signed_commits" gorm:"default:false"`
	RequireSignedMergeRequests bool `json:"require_signed_merge_requests" gorm:"default:false"`
	AllowUnsignedCommits bool   `json:"allow_unsigned_commits" gorm:"default:true"`
	EnforceSignedTags bool      `json:"enforce_signed_tags" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

type SignaturePolicyRequest struct {
	RequireSignedCommits     bool `json:"require_signed_commits"`
	RequireSignedMergeRequests bool `json:"require_signed_merge_requests"`
	AllowUnsignedCommits     bool `json:"allow_unsigned_commits"`
	EnforceSignedTags        bool `json:"enforce_signed_tags"`
}

type SignaturePolicyResponse struct {
	ID              int       `json:"id"`
	RepositoryID    int64     `json:"repository_id"`
	RequireSignedCommits bool   `json:"require_signed_commits"`
	RequireSignedMergeRequests bool `json:"require_signed_merge_requests"`
	AllowUnsignedCommits bool   `json:"allow_unsigned_commits"`
	EnforceSignedTags bool      `json:"enforce_signed_tags"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
