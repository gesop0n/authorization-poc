package document

import "errors"

var (
	ErrOnlyDocumentOwnerCanShare = errors.New("only document owner can share")
	ErrCannotShareWithOwner      = errors.New("cannot share document with its owner")
	ErrGranteeSuspended          = errors.New("grantee user is suspended")
)
