package repository

import "errors"

var (
	ErrTeamExists   = errors.New("team already exists")
	ErrPRExists     = errors.New("PR id already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrTeamNotFound = errors.New("team not found")
	ErrPRNotFound   = errors.New("PR not found")
	ErrNoCanditate  = errors.New("no active replacement candidate in team")
	ErrNotAssigned  = errors.New("reviewer is not assigned to this PR")
	ErrPRMerged     = errors.New("cannot reassign on merged PR")
)
