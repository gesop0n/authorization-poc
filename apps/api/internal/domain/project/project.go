package project

import (
	"fmt"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
)

type Project struct {
	id          ProjectID
	name        string
	workspaceID workspace.WorkspaceID
}

func NewProject(name string, workspaceID workspace.WorkspaceID) (*Project, error) {
	id, err := newProjectID()
	if err != nil {
		return nil, fmt.Errorf("create Project: %w", err)
	}

	return &Project{
		id:          id,
		name:        name,
		workspaceID: workspaceID,
	}, nil
}

func Reconstruct(id ProjectID, name string, workspaceID workspace.WorkspaceID) (*Project, error) {
	if id.IsZero() {
		return nil, ErrInvalidProjectID
	}
	return &Project{id: id, name: name, workspaceID: workspaceID}, nil
}

func (p *Project) ID() ProjectID {
	return p.id
}

func (p *Project) Name() string {
	return p.name
}

func (p *Project) WorkspaceID() workspace.WorkspaceID {
	return p.workspaceID
}
