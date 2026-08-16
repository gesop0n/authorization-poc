package project

import (
	"context"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
)

type ProjectRepository interface {
	FindByID(ctx context.Context, id ProjectID) (*Project, error)
	FindByWorkspaceID(ctx context.Context, workspaceID workspace.WorkspaceID) ([]*Project, error)
	Save(ctx context.Context, project *Project) error
}
