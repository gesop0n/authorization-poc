package document_test

import (
	"context"
	"errors"
	"testing"

	domaindocument "github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	domainproject "github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
	"github.com/gesop0n/authorization-poc/apps/api/internal/testutil"
	documentusecase "github.com/gesop0n/authorization-poc/apps/api/internal/usecase/document"
)

func TestCreateUseCase(t *testing.T) {
	t.Parallel()
	p, _, owner := testutil.NewProjectWithWorkspace(t, "project", "workspace", "owner")
	documentRepository := &mockDocumentRepository{}
	uc := documentusecase.NewCreateUseCase(&mockProjectRepository{project: p}, documentRepository)

	doc, err := uc.Execute(context.Background(), documentusecase.CreateInput{
		ProjectID:   p.ID(),
		OwnerUserID: owner.ID(),
		Title:       "title",
		Content:     "content",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if doc == nil || documentRepository.saved != doc || !doc.ProjectID().Equal(p.ID()) {
		t.Fatal("Execute() did not create and save a document")
	}
}

func TestCreateUseCaseRejectsArchivedProject(t *testing.T) {
	t.Parallel()
	p, _, owner := testutil.NewProjectWithWorkspace(t, "project", "workspace", "owner")
	p.Archive()
	documentRepository := &mockDocumentRepository{}
	uc := documentusecase.NewCreateUseCase(&mockProjectRepository{project: p}, documentRepository)

	doc, err := uc.Execute(context.Background(), documentusecase.CreateInput{
		ProjectID:   p.ID(),
		OwnerUserID: owner.ID(),
		Title:       "title",
		Content:     "content",
	})
	if !errors.Is(err, domainproject.ErrProjectArchived) {
		t.Fatalf("Execute() error = %v, want %v", err, domainproject.ErrProjectArchived)
	}
	if doc != nil || documentRepository.saved != nil {
		t.Fatal("Execute() created or saved a document in an archived project")
	}
}

func TestEditUseCase(t *testing.T) {
	t.Parallel()
	doc, p, owner := testutil.NewDocumentWithProject(t)
	documentRepository := &mockDocumentRepository{document: doc}
	uc := documentusecase.NewEditUseCase(&mockProjectRepository{project: p}, documentRepository)

	err := uc.Execute(context.Background(), documentusecase.EditInput{
		DocumentID:  doc.ID(),
		ActorUserID: owner.ID(),
		Title:       "changed",
		Content:     "changed content",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if doc.Title() != "changed" || doc.Content() != "changed content" || documentRepository.saved != doc {
		t.Fatal("Execute() did not edit and save the document")
	}
}

func TestEditUseCaseRejectsArchivedProject(t *testing.T) {
	t.Parallel()
	doc, p, owner := testutil.NewDocumentWithProject(t)
	p.Archive()
	documentRepository := &mockDocumentRepository{document: doc}
	uc := documentusecase.NewEditUseCase(&mockProjectRepository{project: p}, documentRepository)

	err := uc.Execute(context.Background(), documentusecase.EditInput{
		DocumentID: doc.ID(), ActorUserID: owner.ID(), Title: "changed", Content: "changed content",
	})
	if !errors.Is(err, domainproject.ErrProjectArchived) {
		t.Fatalf("Execute() error = %v, want %v", err, domainproject.ErrProjectArchived)
	}
	assertDocumentUnchangedAndUnsaved(t, doc, documentRepository)
}

func TestEditUseCaseRejectsNonOwner(t *testing.T) {
	t.Parallel()
	doc, p, _ := testutil.NewDocumentWithProject(t)
	nonOwner := testutil.NewUser(t, "non-owner")
	documentRepository := &mockDocumentRepository{document: doc}
	uc := documentusecase.NewEditUseCase(&mockProjectRepository{project: p}, documentRepository)

	err := uc.Execute(context.Background(), documentusecase.EditInput{
		DocumentID: doc.ID(), ActorUserID: nonOwner.ID(), Title: "changed", Content: "changed content",
	})
	if !errors.Is(err, documentusecase.ErrOnlyDocumentOwner) {
		t.Fatalf("Execute() error = %v, want %v", err, documentusecase.ErrOnlyDocumentOwner)
	}
	assertDocumentUnchangedAndUnsaved(t, doc, documentRepository)
}

func TestEditUseCaseRejectsInvalidTitleWithoutPartialChange(t *testing.T) {
	t.Parallel()
	doc, p, owner := testutil.NewDocumentWithProject(t)
	documentRepository := &mockDocumentRepository{document: doc}
	uc := documentusecase.NewEditUseCase(&mockProjectRepository{project: p}, documentRepository)

	err := uc.Execute(context.Background(), documentusecase.EditInput{
		DocumentID: doc.ID(), ActorUserID: owner.ID(), Title: "  ", Content: "changed content",
	})
	if !errors.Is(err, domaindocument.ErrDocumentTitleRequired) {
		t.Fatalf("Execute() error = %v, want %v", err, domaindocument.ErrDocumentTitleRequired)
	}
	assertDocumentUnchangedAndUnsaved(t, doc, documentRepository)
}

func TestDeleteUseCase(t *testing.T) {
	t.Parallel()
	doc, _, owner := testutil.NewDocumentWithProject(t)
	repository := &mockDocumentRepository{document: doc}
	uc := documentusecase.NewDeleteUseCase(repository)

	err := uc.Execute(context.Background(), documentusecase.DeleteInput{
		DocumentID: doc.ID(), ActorUserID: owner.ID(),
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !repository.deletedID.Equal(doc.ID()) {
		t.Fatal("Execute() did not delete the document")
	}
}

func TestDeleteUseCaseRejectsNonOwner(t *testing.T) {
	t.Parallel()
	doc, _, _ := testutil.NewDocumentWithProject(t)
	nonOwner := testutil.NewUser(t, "non-owner")
	repository := &mockDocumentRepository{document: doc}
	uc := documentusecase.NewDeleteUseCase(repository)

	err := uc.Execute(context.Background(), documentusecase.DeleteInput{
		DocumentID: doc.ID(), ActorUserID: nonOwner.ID(),
	})
	if !errors.Is(err, documentusecase.ErrOnlyDocumentOwner) {
		t.Fatalf("Execute() error = %v, want %v", err, documentusecase.ErrOnlyDocumentOwner)
	}
	if !repository.deletedID.IsZero() {
		t.Fatal("Execute() deleted a document owned by another user")
	}
}

func assertDocumentUnchangedAndUnsaved(t *testing.T, doc *domaindocument.Document, repository *mockDocumentRepository) {
	t.Helper()
	if doc.Title() != "title" || doc.Content() != "content" || repository.saved != nil {
		t.Fatal("Execute() changed or saved the document on rejection")
	}
}

type mockProjectRepository struct {
	project *domainproject.Project
	err     error
}

func (r *mockProjectRepository) FindByID(context.Context, domainproject.ProjectID) (*domainproject.Project, error) {
	return r.project, r.err
}

func (r *mockProjectRepository) FindByWorkspaceID(context.Context, workspace.WorkspaceID) ([]*domainproject.Project, error) {
	return nil, nil
}

func (r *mockProjectRepository) Save(context.Context, *domainproject.Project) error { return nil }

type mockDocumentRepository struct {
	document  *domaindocument.Document
	saved     *domaindocument.Document
	deletedID domaindocument.DocumentID
	err       error
}

func (r *mockDocumentRepository) FindByID(context.Context, domaindocument.DocumentID) (*domaindocument.Document, error) {
	return r.document, r.err
}

func (r *mockDocumentRepository) Save(_ context.Context, doc *domaindocument.Document) error {
	r.saved = doc
	return r.err
}

func (r *mockDocumentRepository) Delete(_ context.Context, id domaindocument.DocumentID) error {
	r.deletedID = id
	return r.err
}
