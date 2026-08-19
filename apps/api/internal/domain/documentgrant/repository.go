package documentgrant

import (
	"context"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
)

type DocumentGrantRepository interface {
	FindByDocumentID(ctx context.Context, documentID document.DocumentID) (DocumentGrants, error)
	Save(ctx context.Context, grants DocumentGrants) error
}
