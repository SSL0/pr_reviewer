package service

import "errors"

var (
	ErrTeamExists       = errors.New("team already exists")
	ErrPRExists         = errors.New("PR id already exists")
	ErrResourceNotFound = errors.New("resource not found")
)
