package project

import "fmt"

type ProjectRole string

const (
	ProjectRoleAdmin  ProjectRole = "admin"
	ProjectRoleEditor ProjectRole = "editor"
	ProjectRoleViewer ProjectRole = "viewer"
)

func NewProjectRole(value string) (ProjectRole, error) {
	role := ProjectRole(value)
	if !role.IsValid() {
		return "", fmt.Errorf("%w: %q", ErrInvalidProjectRole, value)
	}
	return role, nil
}

func (role ProjectRole) IsValid() bool {
	switch role {
	case ProjectRoleAdmin, ProjectRoleEditor, ProjectRoleViewer:
		return true
	default:
		return false
	}
}
