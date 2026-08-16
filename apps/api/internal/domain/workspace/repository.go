package workspace

import "context"

type WorkspaceRepository interface {
	FindByID(ctx context.Context, id WorkspaceID) (*Workspace, error)
	Save(ctx context.Context, workspace *Workspace) error
}
