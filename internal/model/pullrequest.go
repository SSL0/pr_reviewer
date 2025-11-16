package model

import "time"

type PullRequestStatus string

const (
	PullRequestOpen   PullRequestStatus = "OPEN"
	PullRequestMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	ID        string            `db:"id"`
	Name      string            `db:"name"`
	AuthorID  string            `db:"author_id"`
	Status    PullRequestStatus `db:"status"`
	CreatedAt time.Time         `db:"created_at"`
	MergedAt  *time.Time        `db:"merged_at"`
}

type PullRequestReviewers struct {
	ID            string `db:"id"`
	PullRequestID string `db:"pull_request_id"`
	ReviewerID    string `db:"reviewer_id"`
}
