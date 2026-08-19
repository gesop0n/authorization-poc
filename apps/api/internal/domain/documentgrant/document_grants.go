package documentgrant

import (
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
)

// DocumentGrants は、同一Document・Userへの重複付与を防ぐコレクション。
// 永続化時にも同じ組み合わせへ一意制約を設定する。
type DocumentGrants struct {
	items []DocumentGrant
}

func NewDocumentGrants(items []DocumentGrant) (DocumentGrants, error) {
	grants := DocumentGrants{}
	for _, item := range items {
		if grants.index(item.DocumentID(), item.GranteeUserID()) >= 0 {
			return DocumentGrants{}, ErrDocumentGrantAlreadyExists
		}
		grants.items = append(grants.items, item)
	}
	return grants, nil
}

func (g *DocumentGrants) Grant(
	documentID document.DocumentID,
	granteeUserID user.UserID,
	role GrantRole,
	grantedBy user.UserID,
) error {
	if g.index(documentID, granteeUserID) >= 0 {
		return ErrDocumentGrantAlreadyExists
	}
	grant, err := NewDocumentGrant(documentID, granteeUserID, role, grantedBy)
	if err != nil {
		return err
	}
	g.items = append(g.items, grant)
	return nil
}

func (g *DocumentGrants) ChangeRole(documentID document.DocumentID, granteeUserID user.UserID, role GrantRole) error {
	index := g.index(documentID, granteeUserID)
	if index < 0 {
		return ErrDocumentGrantNotFound
	}
	return g.items[index].ChangeRole(role)
}

func (g *DocumentGrants) Revoke(documentID document.DocumentID, granteeUserID user.UserID) error {
	index := g.index(documentID, granteeUserID)
	if index < 0 {
		return ErrDocumentGrantNotFound
	}
	return g.items[index].Revoke()
}

func (g DocumentGrants) Items() []DocumentGrant {
	return append([]DocumentGrant(nil), g.items...)
}

func (g DocumentGrants) index(documentID document.DocumentID, granteeUserID user.UserID) int {
	for i := range g.items {
		if g.items[i].sameTarget(documentID, granteeUserID) {
			return i
		}
	}
	return -1
}
