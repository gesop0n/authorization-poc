package workspace

import (
	"context"

	domainproject "github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	domainworkspace "github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
	"github.com/gesop0n/authorization-poc/apps/api/internal/usecase/transaction"
)

type RemoveMemberUseCase struct {
	workspaceRepository domainworkspace.WorkspaceRepository
	projectRepository   domainproject.ProjectRepository
	transactionManager  transaction.Manager
}

type RemoveMemberInput struct {
	WorkspaceID domainworkspace.WorkspaceID
	UserID      user.UserID
}

func NewRemoveMemberUseCase(
	workspaceRepository domainworkspace.WorkspaceRepository,
	projectRepository domainproject.ProjectRepository,
	transactionManager transaction.Manager,
) *RemoveMemberUseCase {
	return &RemoveMemberUseCase{
		workspaceRepository: workspaceRepository,
		projectRepository:   projectRepository,
		transactionManager:  transactionManager,
	}
}

// Execute は、ユーザーをWorkspaceと配下の全Projectから一括で脱退させる。
func (uc *RemoveMemberUseCase) Execute(ctx context.Context, input RemoveMemberInput) error {
	return uc.transactionManager.WithinTransaction(ctx, func(ctx context.Context) error {
		ws, err := uc.workspaceRepository.FindByID(ctx, input.WorkspaceID)
		if err != nil {
			return err
		}

		projects, err := uc.projectRepository.FindByWorkspaceID(ctx, input.WorkspaceID)
		if err != nil {
			return err
		}

		if err := validateRemoval(ws, projects, input.UserID); err != nil {
			return err
		}

		for _, p := range projects {
			if !p.HasMember(input.UserID) {
				continue
			}
			if err := p.RemoveMember(input.UserID); err != nil {
				return err
			}
			if err := uc.projectRepository.Save(ctx, p); err != nil {
				return err
			}
		}

		if err := ws.RemoveMember(input.UserID); err != nil {
			return err
		}
		return uc.workspaceRepository.Save(ctx, ws)
	})
}

func validateRemoval(ws *domainworkspace.Workspace, projects []*domainproject.Project, userID user.UserID) error {
	if err := ws.CanRemoveMember(userID); err != nil {
		return err
	}
	for _, p := range projects {
		if !p.WorkspaceID().Equal(ws.ID()) {
			return ErrProjectWorkspaceMismatch
		}
		if p.HasMember(userID) {
			if err := p.CanRemoveMember(userID); err != nil {
				return err
			}
		}
	}
	return nil
}
