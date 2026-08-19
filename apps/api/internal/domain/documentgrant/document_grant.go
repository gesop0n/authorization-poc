package documentgrant

import (
	"time"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
)

// DocumentGrant は、特定のDocumentに対してUserへ個別付与した権限を表す。
// 取り消し後も監査可能なようにレコードを残す。
type DocumentGrant struct {
	documentID    document.DocumentID
	granteeUserID user.UserID
	role          GrantRole
	grantedBy     user.UserID
	grantedAt     time.Time
	revokedAt     *time.Time
}

func NewDocumentGrant(
	documentID document.DocumentID,
	granteeUserID user.UserID,
	role GrantRole,
	grantedBy user.UserID,
) (DocumentGrant, error) {
	if documentID.IsZero() {
		return DocumentGrant{}, document.ErrInvalidDocumentID
	}
	if granteeUserID.IsZero() || grantedBy.IsZero() {
		return DocumentGrant{}, user.ErrInvalidUserID
	}
	if !role.IsValid() {
		return DocumentGrant{}, ErrInvalidGrantRole
	}

	return DocumentGrant{
		documentID:    documentID,
		granteeUserID: granteeUserID,
		role:          role,
		grantedBy:     grantedBy,
		grantedAt:     time.Now(),
	}, nil
}

func (g *DocumentGrant) ChangeRole(role GrantRole) error {
	if g.IsRevoked() {
		return ErrDocumentGrantRevoked
	}
	if !role.IsValid() {
		return ErrInvalidGrantRole
	}
	g.role = role
	return nil
}

func (g *DocumentGrant) Revoke() error {
	if g.IsRevoked() {
		return ErrDocumentGrantRevoked
	}
	now := time.Now()
	g.revokedAt = &now
	return nil
}

func (g DocumentGrant) DocumentID() document.DocumentID { return g.documentID }
func (g DocumentGrant) GranteeUserID() user.UserID      { return g.granteeUserID }
func (g DocumentGrant) Role() GrantRole                 { return g.role }
func (g DocumentGrant) GrantedBy() user.UserID          { return g.grantedBy }
func (g DocumentGrant) GrantedAt() time.Time            { return g.grantedAt }
func (g DocumentGrant) IsRevoked() bool                 { return g.revokedAt != nil }

func (g DocumentGrant) RevokedAt() (time.Time, bool) {
	if g.revokedAt == nil {
		return time.Time{}, false
	}
	return *g.revokedAt, true
}

func (g DocumentGrant) sameTarget(documentID document.DocumentID, userID user.UserID) bool {
	return g.documentID.Equal(documentID) && g.granteeUserID.Equal(userID)
}
