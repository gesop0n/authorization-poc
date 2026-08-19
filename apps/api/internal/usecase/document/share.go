package document

import (
	"context"

	domaindocument "github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/documentgrant"
	domainproject "github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
)

type ShareUseCase struct {
	documentRepository      domaindocument.DocumentRepository
	projectRepository       domainproject.ProjectRepository
	workspaceRepository     workspace.WorkspaceRepository
	userRepository          user.UserRepository
	documentGrantRepository documentgrant.DocumentGrantRepository
}

type ShareInput struct {
	ActorUserID   user.UserID
	DocumentID    domaindocument.DocumentID
	GranteeUserID user.UserID
	Role          documentgrant.GrantRole
}

func NewShareUseCase(
	documentRepository domaindocument.DocumentRepository,
	projectRepository domainproject.ProjectRepository,
	workspaceRepository workspace.WorkspaceRepository,
	userRepository user.UserRepository,
	documentGrantRepository documentgrant.DocumentGrantRepository,
) *ShareUseCase {
	return &ShareUseCase{
		documentRepository:      documentRepository,
		projectRepository:       projectRepository,
		workspaceRepository:     workspaceRepository,
		userRepository:          userRepository,
		documentGrantRepository: documentGrantRepository,
	}
}

func (uc *ShareUseCase) Execute(ctx context.Context, input ShareInput) error {
	doc, err := uc.documentRepository.FindByID(ctx, input.DocumentID)
	if err != nil {
		return err
	}
	if doc.Status() == domaindocument.DocumentStatusArchived {
		return domaindocument.ErrDocumentArchived
	}

	p, err := uc.projectRepository.FindByID(ctx, doc.ProjectID())
	if err != nil {
		return err
	}
	if p.Status() == domainproject.ProjectStatusArchived {
		return domainproject.ErrProjectArchived
	}

	if !doc.OwnerUserID().Equal(input.ActorUserID) {
		return ErrOnlyDocumentOwnerCanShare
	}
	if doc.OwnerUserID().Equal(input.GranteeUserID) {
		return ErrCannotShareWithOwner
	}

	grantee, err := uc.userRepository.FindByID(ctx, input.GranteeUserID)
	if err != nil {
		return err
	}
	if grantee.Status() == user.UserStatusSuspended {
		return ErrGranteeSuspended
	}

	ws, err := uc.workspaceRepository.FindByID(ctx, p.WorkspaceID())
	if err != nil {
		return err
	}
	if err := documentgrant.CanShareWith(doc.Confidentiality(), ws.HasMember(input.GranteeUserID)); err != nil {
		return err
	}

	grants, err := uc.documentGrantRepository.FindByDocumentID(ctx, doc.ID())
	if err != nil {
		return err
	}
	if err := grants.Grant(doc.ID(), input.GranteeUserID, input.Role, input.ActorUserID); err != nil {
		return err
	}
	return uc.documentGrantRepository.Save(ctx, grants)
}
