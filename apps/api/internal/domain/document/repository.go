package document

import "context"

type DocumentRepository interface {
	FindByID(ctx context.Context, id DocumentID) (*Document, error)
	Save(ctx context.Context, document *Document) error
}
