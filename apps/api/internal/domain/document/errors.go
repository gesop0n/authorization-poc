package document

import "errors"

var (
	// ドキュメントIDが不正
	ErrInvalidDocumentID = errors.New("invalid document ID")

	// ドキュメントのタイトルは必須
	ErrDocumentTitleRequired = errors.New("document title is required")

	// 不正な Confidentiality
	ErrInvalidConfidentiality = errors.New("invalid confidentiality")

	// ドキュメントの状態が不正
	ErrInvalidDocumentStatus = errors.New("invalid document status")

	// ドキュメントはアーカイブ済み
	ErrDocumentArchived = errors.New("document is archived")
)
