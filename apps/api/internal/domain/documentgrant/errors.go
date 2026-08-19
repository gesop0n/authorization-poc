package documentgrant

import "errors"

var (
	ErrInvalidGrantRole           = errors.New("invalid document grant role")
	ErrDocumentGrantAlreadyExists = errors.New("document grant already exists")
	ErrDocumentGrantNotFound      = errors.New("document grant not found")
	ErrDocumentGrantRevoked       = errors.New("document grant is revoked")
)
