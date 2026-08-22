package document

import (
	"context"

	domaindocument "github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
)

type DeleteUseCase struct {
	documentRepository domaindocument.DocumentRepository
}

type DeleteInput struct {
	DocumentID  domaindocument.DocumentID
	ActorUserID user.UserID
}

func NewDeleteUseCase(documentRepository domaindocument.DocumentRepository) *DeleteUseCase {
	return &DeleteUseCase{documentRepository: documentRepository}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, input DeleteInput) error {
	doc, err := uc.documentRepository.FindByID(ctx, input.DocumentID)
	if err != nil {
		return err
	}
	if !doc.OwnerUserID().Equal(input.ActorUserID) {
		return ErrOnlyDocumentOwner
	}
	return uc.documentRepository.Delete(ctx, doc.ID())
}
