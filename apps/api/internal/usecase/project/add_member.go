package projectusecase

import (
	"context"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
)

type AddMemberUseCase struct {
	workspaceRepository workspace.WorkspaceRepository
	projectRepository   project.ProjectRepository
}

type AddMemberInput struct {
	WorkspaceID workspace.WorkspaceID
	ProjectID   project.ProjectID
	UserID      user.UserID
	Role        project.ProjectRole
}

func NewAddMemberUseCase(
	workspaceRepository workspace.WorkspaceRepository,
	projectRepository project.ProjectRepository,
) *AddMemberUseCase {
	return &AddMemberUseCase{
		workspaceRepository: workspaceRepository,
		projectRepository:   projectRepository,
	}
}

func (uc *AddMemberUseCase) Execute(ctx context.Context, input AddMemberInput) error {
	ws, err := uc.workspaceRepository.FindByID(ctx, input.WorkspaceID)
	if err != nil {
		return err
	}

	p, err := uc.projectRepository.FindByID(ctx, input.ProjectID)
	if err != nil {
		return err
	}

	if !p.WorkspaceID().Equal(ws.ID()) {
		return ErrProjectWorkspaceMismatch
	}

	if !ws.HasMember(input.UserID) {
		return ErrUserNotWorkspaceMember
	}

	if err := p.AddMember(input.UserID, input.Role); err != nil {
		return err
	}

	return uc.projectRepository.Save(ctx, p)
}
