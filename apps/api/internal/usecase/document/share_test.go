package document_test

import (
	"context"
	"errors"
	"testing"

	domaindocument "github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/documentgrant"
	domainproject "github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
	"github.com/gesop0n/authorization-poc/apps/api/internal/testutil"
	documentusecase "github.com/gesop0n/authorization-poc/apps/api/internal/usecase/document"
)

func TestShareUseCaseSharesWithWorkspaceMember(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityInternal, true)

	err := fixture.useCase().Execute(context.Background(), fixture.input())
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	assertGrantSaved(t, fixture.grantRepository, fixture.document, fixture.grantee)
}

func TestShareUseCaseSharesPublicDocumentWithExternalUser(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityPublic, false)

	err := fixture.useCase().Execute(context.Background(), fixture.input())
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	assertGrantSaved(t, fixture.grantRepository, fixture.document, fixture.grantee)
}

func TestShareUseCaseRejectsExternalUserForInternalDocument(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityInternal, false)

	err := fixture.useCase().Execute(context.Background(), fixture.input())
	assertRejectedWithoutGrantSave(t, err, documentgrant.ErrExternalSharingProhibited, fixture.grantRepository)
}

func TestShareUseCaseRejectsExternalUserForConfidentialDocument(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityConfidential, false)

	err := fixture.useCase().Execute(context.Background(), fixture.input())
	assertRejectedWithoutGrantSave(t, err, documentgrant.ErrExternalSharingProhibited, fixture.grantRepository)
}

func TestShareUseCaseRejectsPrivateDocumentForWorkspaceMember(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityPrivate, true)

	err := fixture.useCase().Execute(context.Background(), fixture.input())
	assertRejectedWithoutGrantSave(t, err, documentgrant.ErrDocumentNotShareable, fixture.grantRepository)
}

func TestShareUseCaseRejectsNonOwner(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityPublic, true)
	input := fixture.input()
	input.ActorUserID = fixture.grantee.ID()

	err := fixture.useCase().Execute(context.Background(), input)
	assertRejectedWithoutGrantSave(t, err, documentusecase.ErrOnlyDocumentOwnerCanShare, fixture.grantRepository)
}

func TestShareUseCaseRejectsSharingWithOwner(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityPublic, true)
	input := fixture.input()
	input.GranteeUserID = fixture.owner.ID()

	err := fixture.useCase().Execute(context.Background(), input)
	assertRejectedWithoutGrantSave(t, err, documentusecase.ErrCannotShareWithOwner, fixture.grantRepository)
}

func TestShareUseCaseRejectsSuspendedGrantee(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityPublic, true)
	fixture.grantee.Suspend()

	err := fixture.useCase().Execute(context.Background(), fixture.input())
	assertRejectedWithoutGrantSave(t, err, documentusecase.ErrGranteeSuspended, fixture.grantRepository)
}

func TestShareUseCaseRejectsArchivedProject(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityPublic, true)
	fixture.project.Archive()

	err := fixture.useCase().Execute(context.Background(), fixture.input())
	assertRejectedWithoutGrantSave(t, err, domainproject.ErrProjectArchived, fixture.grantRepository)
}

func TestShareUseCaseRejectsArchivedDocument(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityPublic, true)
	fixture.document.Archive()

	err := fixture.useCase().Execute(context.Background(), fixture.input())
	assertRejectedWithoutGrantSave(t, err, domaindocument.ErrDocumentArchived, fixture.grantRepository)
}

func TestShareUseCaseRejectsDuplicateGrant(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityPublic, true)
	if err := fixture.grantRepository.grants.Grant(
		fixture.document.ID(), fixture.grantee.ID(), documentgrant.GrantRoleViewer, fixture.owner.ID(),
	); err != nil {
		t.Fatal(err)
	}

	err := fixture.useCase().Execute(context.Background(), fixture.input())
	assertRejectedWithoutGrantSave(t, err, documentgrant.ErrDocumentGrantAlreadyExists, fixture.grantRepository)
}

func TestShareUseCaseRejectsUnknownGrantee(t *testing.T) {
	t.Parallel()
	fixture := newShareFixture(t, domaindocument.ConfidentialityPublic, true)
	wantErr := errors.New("user not found")
	fixture.userRepository.err = wantErr

	err := fixture.useCase().Execute(context.Background(), fixture.input())
	assertRejectedWithoutGrantSave(t, err, wantErr, fixture.grantRepository)
}

type shareFixture struct {
	document            *domaindocument.Document
	project             *domainproject.Project
	workspace           *workspace.Workspace
	owner               *user.User
	grantee             *user.User
	documentRepository  *mockDocumentRepository
	projectRepository   *mockProjectRepository
	workspaceRepository *shareWorkspaceRepository
	userRepository      *shareUserRepository
	grantRepository     *shareGrantRepository
}

func newShareFixture(t *testing.T, confidentiality domaindocument.Confidentiality, internal bool) *shareFixture {
	t.Helper()
	p, ws, owner := testutil.NewProjectWithWorkspace(t, "project", "workspace", "owner")
	grantee := testutil.NewUser(t, "grantee")
	if internal {
		if err := ws.AddMember(grantee.ID(), workspace.WorkspaceRoleMember); err != nil {
			t.Fatal(err)
		}
	}
	doc, err := domaindocument.NewDocument(p.ID(), owner.ID(), "title", "content", confidentiality)
	if err != nil {
		t.Fatal(err)
	}
	grants, err := documentgrant.NewDocumentGrants(nil)
	if err != nil {
		t.Fatal(err)
	}
	return &shareFixture{
		document:            doc,
		project:             p,
		workspace:           ws,
		owner:               owner,
		grantee:             grantee,
		documentRepository:  &mockDocumentRepository{document: doc},
		projectRepository:   &mockProjectRepository{project: p},
		workspaceRepository: &shareWorkspaceRepository{workspace: ws},
		userRepository:      &shareUserRepository{user: grantee},
		grantRepository:     &shareGrantRepository{grants: grants},
	}
}

func (f *shareFixture) useCase() *documentusecase.ShareUseCase {
	return documentusecase.NewShareUseCase(
		f.documentRepository,
		f.projectRepository,
		f.workspaceRepository,
		f.userRepository,
		f.grantRepository,
	)
}

func (f *shareFixture) input() documentusecase.ShareInput {
	return documentusecase.ShareInput{
		ActorUserID:   f.owner.ID(),
		DocumentID:    f.document.ID(),
		GranteeUserID: f.grantee.ID(),
		Role:          documentgrant.GrantRoleViewer,
	}
}

func assertGrantSaved(t *testing.T, repository *shareGrantRepository, doc *domaindocument.Document, grantee *user.User) {
	t.Helper()
	if repository.saved == nil {
		t.Fatal("DocumentGrantRepository.Save() was not called")
	}
	items := repository.saved.Items()
	if len(items) != 1 || !items[0].DocumentID().Equal(doc.ID()) || !items[0].GranteeUserID().Equal(grantee.ID()) {
		t.Fatal("unexpected DocumentGrant was saved")
	}
}

func assertRejectedWithoutGrantSave(t *testing.T, err, wantErr error, repository *shareGrantRepository) {
	t.Helper()
	if !errors.Is(err, wantErr) {
		t.Fatalf("Execute() error = %v, want %v", err, wantErr)
	}
	if repository.saved != nil {
		t.Fatal("Execute() saved grants on rejection")
	}
}

type shareWorkspaceRepository struct {
	workspace *workspace.Workspace
	err       error
}

func (r *shareWorkspaceRepository) FindByID(context.Context, workspace.WorkspaceID) (*workspace.Workspace, error) {
	return r.workspace, r.err
}

func (r *shareWorkspaceRepository) Save(context.Context, *workspace.Workspace) error { return nil }

type shareUserRepository struct {
	user *user.User
	err  error
}

func (r *shareUserRepository) FindByID(context.Context, user.UserID) (*user.User, error) {
	return r.user, r.err
}

type shareGrantRepository struct {
	grants documentgrant.DocumentGrants
	saved  *documentgrant.DocumentGrants
	err    error
}

func (r *shareGrantRepository) FindByDocumentID(context.Context, domaindocument.DocumentID) (documentgrant.DocumentGrants, error) {
	return r.grants, r.err
}

func (r *shareGrantRepository) Save(_ context.Context, grants documentgrant.DocumentGrants) error {
	r.saved = &grants
	return r.err
}
