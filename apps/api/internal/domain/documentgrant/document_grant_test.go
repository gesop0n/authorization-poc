package documentgrant_test

import (
	"errors"
	"testing"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/documentgrant"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	"github.com/gesop0n/authorization-poc/apps/api/internal/testutil"
)

func TestNewDocumentGrant(t *testing.T) {
	t.Parallel()
	d, _, owner := testutil.NewDocumentWithProject(t)
	grantee := newTestUser(t, "grantee")

	grant, err := documentgrant.NewDocumentGrant(d.ID(), grantee.ID(), documentgrant.GrantRoleViewer, owner.ID())
	if err != nil {
		t.Fatalf("NewDocumentGrant() error = %v", err)
	}
	if !grant.DocumentID().Equal(d.ID()) || !grant.GranteeUserID().Equal(grantee.ID()) || grant.Role() != documentgrant.GrantRoleViewer || !grant.GrantedBy().Equal(owner.ID()) {
		t.Fatal("NewDocumentGrant() returned unexpected values")
	}
	if grant.GrantedAt().IsZero() || grant.IsRevoked() {
		t.Fatal("NewDocumentGrant() returned an invalid state")
	}
}

func TestNewDocumentGrantRejectsInvalidValues(t *testing.T) {
	t.Parallel()
	d, _, owner := testutil.NewDocumentWithProject(t)
	grantee := newTestUser(t, "grantee")

	tests := []struct {
		name       string
		documentID document.DocumentID
		granteeID  user.UserID
		role       documentgrant.GrantRole
		grantedBy  user.UserID
		wantErr    error
	}{
		{"Document IDがゼロ値", document.DocumentID{}, grantee.ID(), documentgrant.GrantRoleViewer, owner.ID(), document.ErrInvalidDocumentID},
		{"共有先User IDがゼロ値", d.ID(), user.UserID{}, documentgrant.GrantRoleViewer, owner.ID(), user.ErrInvalidUserID},
		{"付与者User IDがゼロ値", d.ID(), grantee.ID(), documentgrant.GrantRoleViewer, user.UserID{}, user.ErrInvalidUserID},
		{"Roleが不正", d.ID(), grantee.ID(), documentgrant.GrantRole("invalid"), owner.ID(), documentgrant.ErrInvalidGrantRole},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := documentgrant.NewDocumentGrant(tt.documentID, tt.granteeID, tt.role, tt.grantedBy)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewDocumentGrant() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDocumentGrantChangeRoleAndRevoke(t *testing.T) {
	t.Parallel()
	d, _, owner := testutil.NewDocumentWithProject(t)
	grantee := newTestUser(t, "grantee")
	grant, err := documentgrant.NewDocumentGrant(d.ID(), grantee.ID(), documentgrant.GrantRoleViewer, owner.ID())
	if err != nil {
		t.Fatal(err)
	}

	if err := grant.ChangeRole(documentgrant.GrantRoleEditor); err != nil {
		t.Fatalf("ChangeRole() error = %v", err)
	}
	if grant.Role() != documentgrant.GrantRoleEditor {
		t.Fatal("ChangeRole() did not change role")
	}
	if err := grant.Revoke(); err != nil {
		t.Fatalf("Revoke() error = %v", err)
	}
	if _, ok := grant.RevokedAt(); !grant.IsRevoked() || !ok {
		t.Fatal("Revoke() did not revoke grant")
	}
}

func TestDocumentGrantRejectsInvalidOrRevokedChanges(t *testing.T) {
	t.Parallel()
	d, _, owner := testutil.NewDocumentWithProject(t)
	grantee := newTestUser(t, "grantee")
	grant, err := documentgrant.NewDocumentGrant(d.ID(), grantee.ID(), documentgrant.GrantRoleViewer, owner.ID())
	if err != nil {
		t.Fatal(err)
	}

	if err := grant.ChangeRole(documentgrant.GrantRole("invalid")); !errors.Is(err, documentgrant.ErrInvalidGrantRole) || grant.Role() != documentgrant.GrantRoleViewer {
		t.Fatal("ChangeRole() must reject an invalid role without changing state")
	}
	if err := grant.Revoke(); err != nil {
		t.Fatal(err)
	}
	if err := grant.ChangeRole(documentgrant.GrantRoleEditor); !errors.Is(err, documentgrant.ErrDocumentGrantRevoked) {
		t.Fatalf("ChangeRole() error = %v, want %v", err, documentgrant.ErrDocumentGrantRevoked)
	}
	if err := grant.Revoke(); !errors.Is(err, documentgrant.ErrDocumentGrantRevoked) {
		t.Fatalf("Revoke() error = %v, want %v", err, documentgrant.ErrDocumentGrantRevoked)
	}
}

func TestDocumentGrantsRejectsDuplicateTarget(t *testing.T) {
	t.Parallel()
	d, _, owner := testutil.NewDocumentWithProject(t)
	grantee := newTestUser(t, "grantee")
	grants, err := documentgrant.NewDocumentGrants(nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := grants.Grant(d.ID(), grantee.ID(), documentgrant.GrantRoleViewer, owner.ID()); err != nil {
		t.Fatalf("Grant() error = %v", err)
	}
	if err := grants.Grant(d.ID(), grantee.ID(), documentgrant.GrantRoleEditor, owner.ID()); !errors.Is(err, documentgrant.ErrDocumentGrantAlreadyExists) {
		t.Fatalf("duplicate Grant() error = %v, want %v", err, documentgrant.ErrDocumentGrantAlreadyExists)
	}
	if len(grants.Items()) != 1 || grants.Items()[0].Role() != documentgrant.GrantRoleViewer {
		t.Fatal("duplicate Grant() changed collection")
	}
}

func TestDocumentGrantsChangesAndRevokesExistingGrant(t *testing.T) {
	t.Parallel()
	d, _, owner := testutil.NewDocumentWithProject(t)
	grantee := newTestUser(t, "grantee")
	grants, _ := documentgrant.NewDocumentGrants(nil)
	if err := grants.Grant(d.ID(), grantee.ID(), documentgrant.GrantRoleViewer, owner.ID()); err != nil {
		t.Fatal(err)
	}

	if err := grants.ChangeRole(d.ID(), grantee.ID(), documentgrant.GrantRoleEditor); err != nil {
		t.Fatalf("ChangeRole() error = %v", err)
	}
	if err := grants.Revoke(d.ID(), grantee.ID()); err != nil {
		t.Fatalf("Revoke() error = %v", err)
	}
	items := grants.Items()
	if items[0].Role() != documentgrant.GrantRoleEditor || !items[0].IsRevoked() {
		t.Fatal("collection did not update the grant")
	}
}

func TestDocumentGrantsRejectsUnknownTarget(t *testing.T) {
	t.Parallel()
	d, _, _ := testutil.NewDocumentWithProject(t)
	grantee := newTestUser(t, "grantee")
	grants, _ := documentgrant.NewDocumentGrants(nil)

	if err := grants.ChangeRole(d.ID(), grantee.ID(), documentgrant.GrantRoleEditor); !errors.Is(err, documentgrant.ErrDocumentGrantNotFound) {
		t.Fatalf("ChangeRole() error = %v, want %v", err, documentgrant.ErrDocumentGrantNotFound)
	}
	if err := grants.Revoke(d.ID(), grantee.ID()); !errors.Is(err, documentgrant.ErrDocumentGrantNotFound) {
		t.Fatalf("Revoke() error = %v, want %v", err, documentgrant.ErrDocumentGrantNotFound)
	}
}

func newTestUser(t *testing.T, name string) *user.User {
	t.Helper()
	u, err := user.NewUser(name)
	if err != nil {
		t.Fatal(err)
	}
	return u
}
