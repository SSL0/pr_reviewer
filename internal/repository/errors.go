package repository

import "errors"

var (
	ErrTeamExists   = errors.New("team already exists")
	ErrPRExists     = errors.New("PR id already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrPRNotFound   = errors.New("PR not found")
)
