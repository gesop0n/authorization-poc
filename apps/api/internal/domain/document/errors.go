package document

import "errors"

var (
	// ドキュメントIDが不正
	ErrInvalidDocumentID = errors.New("invalid document ID")

	// ドキュメントのタイトルは必須
	ErrDocumentTitleRequired = errors.New("document title is required")
)
