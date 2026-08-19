package documentgrant

import "github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"

// CanShareWith は、Documentの機密区分と共有先のWorkspace所属状況から
// 個別共有が可能かを判定する。
func CanShareWith(confidentiality document.Confidentiality, isWorkspaceMember bool) error {
	switch confidentiality {
	case document.ConfidentialityPublic:
		return nil
	case document.ConfidentialityInternal, document.ConfidentialityConfidential:
		if !isWorkspaceMember {
			return ErrExternalSharingProhibited
		}
		return nil
	case document.ConfidentialityPrivate:
		return ErrDocumentNotShareable
	default:
		return document.ErrInvalidConfidentiality
	}
}
