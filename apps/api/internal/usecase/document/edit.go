package document

import (
	"context"

	domaindocument "github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	domainproject "github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
)

type EditUseCase struct {
	projectRepository  domainproject.ProjectRepository
	documentRepository domaindocument.DocumentRepository
}

type EditInput struct {
	DocumentID  domaindocument.DocumentID
	ActorUserID user.UserID
	Title       string
	Content     string
}

func NewEditUseCase(
	projectRepository domainproject.ProjectRepository,
	documentRepository domaindocument.DocumentRepository,
) *EditUseCase {
	return &EditUseCase{
		projectRepository:  projectRepository,
		documentRepository: documentRepository,
	}
}

func (uc *EditUseCase) Execute(ctx context.Context, input EditInput) error {
	doc, err := uc.documentRepository.FindByID(ctx, input.DocumentID)
	if err != nil {
		return err
	}
	if !doc.OwnerUserID().Equal(input.ActorUserID) {
		return ErrOnlyDocumentOwner
	}

	p, err := uc.projectRepository.FindByID(ctx, doc.ProjectID())
	if err != nil {
		return err
	}
	if p.Status() == domainproject.ProjectStatusArchived {
		return domainproject.ErrProjectArchived
	}

	// Titleを先に検証し、失敗時にContentだけが変更されることを防ぐ。
	if err := doc.ChangeTitle(input.Title); err != nil {
		return err
	}
	doc.ChangeContent(input.Content)
	return uc.documentRepository.Save(ctx, doc)
}
