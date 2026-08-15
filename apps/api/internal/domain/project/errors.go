package project

import "errors"

var (
	ErrInvalidProjectID           = errors.New("Invalid Project ID")
	ErrProjectNameRequired        = errors.New("project name is required")
	ErrInvalidProjectStatus       = errors.New("invalid project status")
	ErrInvalidProjectMemberUserID = errors.New("invalid project member user ID")
	ErrInvalidProjectRole         = errors.New("invalid project role")
	ErrProjectMemberAlreadyExists = errors.New("project member already exists")
	ErrProjectMemberNotFound      = errors.New("project member not found")
	ErrProjectMustHaveAdmin       = errors.New("project must have at least one admin")
)
