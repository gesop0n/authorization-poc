package document

import "errors"

var (
	ErrOnlyDocumentOwner = errors.New("only document owner can operate the document")
)
