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
		ProjectID:       p.ID(),
		OwnerUserID:     owner.ID(),
		Title:           "title",
		Content:         "content",
		Confidentiality: domaindocument.ConfidentialityInternal,
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
		ProjectID:       p.ID(),
		OwnerUserID:     owner.ID(),
		Title:           "title",
		Content:         "content",
		Confidentiality: domaindocument.ConfidentialityInternal,
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
	doc, p, _ := testutil.NewDocumentWithProject(t)
	documentRepository := &mockDocumentRepository{document: doc}
	uc := documentusecase.NewEditUseCase(&mockProjectRepository{project: p}, documentRepository)

	err := uc.Execute(context.Background(), documentusecase.EditInput{
		DocumentID: doc.ID(),
		Title:      "changed",
		Content:    "changed content",
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
	doc, p, _ := testutil.NewDocumentWithProject(t)
	p.Archive()
	documentRepository := &mockDocumentRepository{document: doc}
	uc := documentusecase.NewEditUseCase(&mockProjectRepository{project: p}, documentRepository)

	err := uc.Execute(context.Background(), documentusecase.EditInput{
		DocumentID: doc.ID(), Title: "changed", Content: "changed content",
	})
	if !errors.Is(err, domainproject.ErrProjectArchived) {
		t.Fatalf("Execute() error = %v, want %v", err, domainproject.ErrProjectArchived)
	}
	assertDocumentUnchangedAndUnsaved(t, doc, documentRepository)
}

func TestEditUseCaseRejectsArchivedDocument(t *testing.T) {
	t.Parallel()
	doc, p, _ := testutil.NewDocumentWithProject(t)
	doc.Archive()
	documentRepository := &mockDocumentRepository{document: doc}
	uc := documentusecase.NewEditUseCase(&mockProjectRepository{project: p}, documentRepository)

	err := uc.Execute(context.Background(), documentusecase.EditInput{
		DocumentID: doc.ID(), Title: "changed", Content: "changed content",
	})
	if !errors.Is(err, domaindocument.ErrDocumentArchived) {
		t.Fatalf("Execute() error = %v, want %v", err, domaindocument.ErrDocumentArchived)
	}
	assertDocumentUnchangedAndUnsaved(t, doc, documentRepository)
}

func TestEditUseCaseRejectsInvalidTitleWithoutPartialChange(t *testing.T) {
	t.Parallel()
	doc, p, _ := testutil.NewDocumentWithProject(t)
	documentRepository := &mockDocumentRepository{document: doc}
	uc := documentusecase.NewEditUseCase(&mockProjectRepository{project: p}, documentRepository)

	err := uc.Execute(context.Background(), documentusecase.EditInput{
		DocumentID: doc.ID(), Title: "  ", Content: "changed content",
	})
	if !errors.Is(err, domaindocument.ErrDocumentTitleRequired) {
		t.Fatalf("Execute() error = %v, want %v", err, domaindocument.ErrDocumentTitleRequired)
	}
	assertDocumentUnchangedAndUnsaved(t, doc, documentRepository)
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
	document *domaindocument.Document
	saved    *domaindocument.Document
	err      error
}

func (r *mockDocumentRepository) FindByID(context.Context, domaindocument.DocumentID) (*domaindocument.Document, error) {
	return r.document, r.err
}

func (r *mockDocumentRepository) Save(_ context.Context, doc *domaindocument.Document) error {
	r.saved = doc
	return r.err
}
