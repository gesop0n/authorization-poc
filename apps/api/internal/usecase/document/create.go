package document

import (
	"context"

	domaindocument "github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	domainproject "github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
)

type CreateUseCase struct {
	projectRepository  domainproject.ProjectRepository
	documentRepository domaindocument.DocumentRepository
}

type CreateInput struct {
	ProjectID       domainproject.ProjectID
	OwnerUserID     user.UserID
	Title           string
	Content         string
	Confidentiality domaindocument.Confidentiality
}

func NewCreateUseCase(
	projectRepository domainproject.ProjectRepository,
	documentRepository domaindocument.DocumentRepository,
) *CreateUseCase {
	return &CreateUseCase{
		projectRepository:  projectRepository,
		documentRepository: documentRepository,
	}
}

func (uc *CreateUseCase) Execute(ctx context.Context, input CreateInput) (*domaindocument.Document, error) {
	p, err := uc.projectRepository.FindByID(ctx, input.ProjectID)
	if err != nil {
		return nil, err
	}
	if p.Status() == domainproject.ProjectStatusArchived {
		return nil, domainproject.ErrProjectArchived
	}

	doc, err := domaindocument.NewDocument(
		p.ID(),
		input.OwnerUserID,
		input.Title,
		input.Content,
		input.Confidentiality,
	)
	if err != nil {
		return nil, err
	}
	if err := uc.documentRepository.Save(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}
