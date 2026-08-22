package document_test

import (
	"errors"
	"testing"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	"github.com/gesop0n/authorization-poc/apps/api/internal/testutil"
)

func TestNewDocument(t *testing.T) {
	t.Parallel()
	d, p, owner := newTestDocument(t)
	if d.ID().IsZero() || !d.ProjectID().Equal(p.ID()) || !d.OwnerUserID().Equal(owner.ID()) {
		t.Fatal("NewDocument() returned unexpected IDs")
	}
	if d.Title() != "title" || d.Content() != "content" {
		t.Fatal("NewDocument() returned unexpected values")
	}
}

func TestNewDocumentRejectsInvalidValues(t *testing.T) {
	t.Parallel()
	_, p, owner := newTestDocument(t)
	tests := []struct {
		name      string
		projectID project.ProjectID
		ownerID   user.UserID
		title     string
		wantErr   error
	}{
		{"Project IDがゼロ値", project.ProjectID{}, owner.ID(), "title", project.ErrInvalidProjectID},
		{"Owner IDがゼロ値", p.ID(), user.UserID{}, "title", user.ErrInvalidUserID},
		{"タイトルが空白のみ", p.ID(), owner.ID(), "  ", document.ErrDocumentTitleRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := document.NewDocument(tt.projectID, tt.ownerID, tt.title, "")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewDocument() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestReconstructDocument(t *testing.T) {
	t.Parallel()
	original, p, owner := newTestDocument(t)
	d, err := document.Reconstruct(original.ID(), p.ID(), owner.ID(), "restored", "body")
	if err != nil {
		t.Fatalf("Reconstruct() error = %v", err)
	}
	if !d.ID().Equal(original.ID()) || d.Title() != "restored" || d.Content() != "body" {
		t.Fatal("Reconstruct() did not restore all values")
	}
}

func TestReconstructDocumentRejectsInvalidValues(t *testing.T) {
	t.Parallel()
	original, p, owner := newTestDocument(t)
	tests := []struct {
		name      string
		id        document.DocumentID
		projectID project.ProjectID
		ownerID   user.UserID
		title     string
		wantErr   error
	}{
		{"Document IDがゼロ値", document.DocumentID{}, p.ID(), owner.ID(), "title", document.ErrInvalidDocumentID},
		{"Project IDがゼロ値", original.ID(), project.ProjectID{}, owner.ID(), "title", project.ErrInvalidProjectID},
		{"Owner IDがゼロ値", original.ID(), p.ID(), user.UserID{}, "title", user.ErrInvalidUserID},
		{"タイトルが空", original.ID(), p.ID(), owner.ID(), "", document.ErrDocumentTitleRequired},
		{"タイトルが空白のみ", original.ID(), p.ID(), owner.ID(), "  ", document.ErrDocumentTitleRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := document.Reconstruct(tt.id, tt.projectID, tt.ownerID, tt.title, "content")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Reconstruct() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDocumentChanges(t *testing.T) {
	t.Parallel()
	d, _, _ := newTestDocument(t)
	newOwner := testutil.NewUser(t, "new owner")
	if err := d.ChangeTitle("changed"); err != nil {
		t.Fatal(err)
	}
	d.ChangeContent("changed content")
	if err := d.ChangeOwner(newOwner.ID()); err != nil {
		t.Fatal(err)
	}
	if d.Title() != "changed" || d.Content() != "changed content" || !d.OwnerUserID().Equal(newOwner.ID()) {
		t.Fatal("change methods did not update all values")
	}
}

func TestDocumentRejectsInvalidChangesWithoutChangingState(t *testing.T) {
	t.Parallel()
	d, _, owner := newTestDocument(t)
	if err := d.ChangeTitle("  "); !errors.Is(err, document.ErrDocumentTitleRequired) || d.Title() != "title" {
		t.Fatal("ChangeTitle() must reject a blank title without changing state")
	}
	if err := d.ChangeOwner(user.UserID{}); !errors.Is(err, user.ErrInvalidUserID) || !d.OwnerUserID().Equal(owner.ID()) {
		t.Fatal("ChangeOwner() must reject a zero ID without changing state")
	}
}

func newTestDocument(t *testing.T) (*document.Document, *project.Project, *user.User) {
	t.Helper()
	return testutil.NewDocumentWithProject(t)
}
