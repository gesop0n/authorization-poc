package projectmembership

import (
	"errors"
	"fmt"
)

type ProjectRole string

var ErrInvalidProjectRole = errors.New("Invalid Project Role")

const (
	ProjectRoleAdmin  ProjectRole = "admin"
	ProjectRoleEditor ProjectRole = "editor"
	ProjectRoleViewer ProjectRole = "viwer"
)

func NewProjectRole(value string) (ProjectRole, error) {
	role := ProjectRole(value)

	switch role {
	case ProjectRoleAdmin, ProjectRoleEditor, ProjectRoleViewer:
		return role, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidProjectRole, value)

	}
}
