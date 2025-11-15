package service

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrTeamNotFound = errors.New("team not found")
	ErrTeamExists   = errors.New("team already exists")
)
