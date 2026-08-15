package document

import (
	"fmt"

	"github.com/google/uuid"
)

type DocumentID struct {
	value uuid.UUID
}

func newDocumentID() (DocumentID, error) {
	value, err := uuid.NewV7()
	if err != nil {
		return DocumentID{}, fmt.Errorf("generate document ID: %w", err)
	}

	return DocumentID{value: value}, nil
}

func parseDocumentID(value string) (DocumentID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return DocumentID{}, fmt.Errorf("%w: %q", ErrInvalidDocumentID, value)
	}

	if parsed.Version() != 7 {
		return DocumentID{}, fmt.Errorf("%w: %q", ErrInvalidDocumentID, value)
	}

	return DocumentID{value: parsed}, nil
}

func (id DocumentID) String() string {
	return id.value.String()
}

func (id DocumentID) Equal(other DocumentID) bool {
	return id.value == other.value
}

func (id DocumentID) IsZero() bool {
	return id.value == uuid.Nil
}
