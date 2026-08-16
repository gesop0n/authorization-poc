package workspace_test

import (
	"context"
	"errors"
	"testing"

	domainproject "github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	domainworkspace "github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
	"github.com/gesop0n/authorization-poc/apps/api/internal/testutil"
	"github.com/gesop0n/authorization-poc/apps/api/internal/usecase/workspace"
)

// WorkspaceメンバーをWorkspaceと参加中の全Projectから削除し、変更した集約を保存することを確認する。
func TestRemoveMemberUseCase(t *testing.T) {
	t.Parallel()
	ws, projects, member := newRemovalTestEntities(t)
	workspaceRepository := &mockWorkspaceRepository{workspace: ws}
	projectRepository := &mockProjectRepository{projects: projects}
	transactionManager := &mockTransactionManager{}
	uc := workspace.NewRemoveMemberUseCase(workspaceRepository, projectRepository, transactionManager)

	err := uc.Execute(context.Background(), workspace.RemoveMemberInput{WorkspaceID: ws.ID(), UserID: member.ID()})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if ws.HasMember(member.ID()) || projects[0].HasMember(member.ID()) {
		t.Fatal("Execute() did not remove the member")
	}
	if workspaceRepository.saved != ws || len(projectRepository.saved) != 1 || projectRepository.saved[0] != projects[0] {
		t.Fatal("Execute() did not save the changed aggregates")
	}
	if transactionManager.callCount != 1 {
		t.Fatalf("WithinTransaction() calls = %d, want 1", transactionManager.callCount)
	}
}

// Projectに参加していないWorkspaceメンバーを脱退させる場合、Projectを保存しないことを確認する。
func TestRemoveMemberUseCaseDoesNotSaveUnchangedProject(t *testing.T) {
	t.Parallel()
	ws, projects, member := newRemovalTestEntities(t)
	if err := projects[0].RemoveMember(member.ID()); err != nil {
		t.Fatal(err)
	}
	projectRepository := &mockProjectRepository{projects: projects}
	uc := workspace.NewRemoveMemberUseCase(
		&mockWorkspaceRepository{workspace: ws}, projectRepository, &mockTransactionManager{},
	)

	err := uc.Execute(context.Background(), workspace.RemoveMemberInput{WorkspaceID: ws.ID(), UserID: member.ID()})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(projectRepository.saved) != 0 {
		t.Fatal("Execute() saved an unchanged project")
	}
}

// 最後のWorkspace Ownerは脱退を拒否し、どの集約も変更・保存しないことを確認する。
func TestRemoveMemberUseCasePreservesLastWorkspaceOwner(t *testing.T) {
	t.Parallel()
	p, ws, owner := testutil.NewProjectWithWorkspace(t, "project", "workspace", "owner")
	workspaceRepository := &mockWorkspaceRepository{workspace: ws}
	projectRepository := &mockProjectRepository{projects: []*domainproject.Project{p}}
	uc := workspace.NewRemoveMemberUseCase(workspaceRepository, projectRepository, &mockTransactionManager{})

	err := uc.Execute(context.Background(), workspace.RemoveMemberInput{WorkspaceID: ws.ID(), UserID: owner.ID()})
	if !errors.Is(err, domainworkspace.ErrWorkspaceMustHaveOwner) {
		t.Fatalf("Execute() error = %v, want %v", err, domainworkspace.ErrWorkspaceMustHaveOwner)
	}
	assertNoRemovalOrSave(t, ws, p, owner.ID(), workspaceRepository, projectRepository)
}

// 最後のProject AdminはWorkspace脱退を拒否し、どの集約も変更・保存しないことを確認する。
func TestRemoveMemberUseCasePreservesLastProjectAdmin(t *testing.T) {
	t.Parallel()
	p, ws, owner := testutil.NewProjectWithWorkspace(t, "project", "workspace", "owner")
	secondOwner := testutil.NewUser(t, "second owner")
	if err := ws.AddMember(secondOwner.ID(), domainworkspace.WorkspaceRoleOwner); err != nil {
		t.Fatal(err)
	}
	workspaceRepository := &mockWorkspaceRepository{workspace: ws}
	projectRepository := &mockProjectRepository{projects: []*domainproject.Project{p}}
	uc := workspace.NewRemoveMemberUseCase(workspaceRepository, projectRepository, &mockTransactionManager{})

	err := uc.Execute(context.Background(), workspace.RemoveMemberInput{WorkspaceID: ws.ID(), UserID: owner.ID()})
	if !errors.Is(err, domainproject.ErrProjectMustHaveAdmin) {
		t.Fatalf("Execute() error = %v, want %v", err, domainproject.ErrProjectMustHaveAdmin)
	}
	assertNoRemovalOrSave(t, ws, p, owner.ID(), workspaceRepository, projectRepository)
}

// Repositoryが別WorkspaceのProjectを返した場合は脱退を拒否することを確認する。
func TestRemoveMemberUseCaseRejectsWorkspaceMismatch(t *testing.T) {
	t.Parallel()
	ws, _, member := newRemovalTestEntities(t)
	otherProject, _, _ := testutil.NewProjectWithWorkspace(t, "other project", "other workspace", "other owner")
	workspaceRepository := &mockWorkspaceRepository{workspace: ws}
	projectRepository := &mockProjectRepository{projects: []*domainproject.Project{otherProject}}
	uc := workspace.NewRemoveMemberUseCase(workspaceRepository, projectRepository, &mockTransactionManager{})

	err := uc.Execute(context.Background(), workspace.RemoveMemberInput{WorkspaceID: ws.ID(), UserID: member.ID()})
	if !errors.Is(err, workspace.ErrProjectWorkspaceMismatch) {
		t.Fatalf("Execute() error = %v, want %v", err, workspace.ErrProjectWorkspaceMismatch)
	}
	if !ws.HasMember(member.ID()) || workspaceRepository.saved != nil || len(projectRepository.saved) != 0 {
		t.Fatal("Execute() changed or saved aggregates on mismatch")
	}
}

func newRemovalTestEntities(t *testing.T) (*domainworkspace.Workspace, []*domainproject.Project, *user.User) {
	t.Helper()
	p, ws, _ := testutil.NewProjectWithWorkspace(t, "project", "workspace", "owner")
	member := testutil.NewUser(t, "member")
	if err := ws.AddMember(member.ID(), domainworkspace.WorkspaceRoleMember); err != nil {
		t.Fatal(err)
	}
	if err := p.AddMember(member.ID(), domainproject.ProjectRoleEditor); err != nil {
		t.Fatal(err)
	}
	return ws, []*domainproject.Project{p}, member
}

func assertNoRemovalOrSave(
	t *testing.T,
	ws *domainworkspace.Workspace,
	p *domainproject.Project,
	userID user.UserID,
	workspaceRepository *mockWorkspaceRepository,
	projectRepository *mockProjectRepository,
) {
	t.Helper()
	if !ws.HasMember(userID) || !p.HasMember(userID) {
		t.Fatal("Execute() changed aggregates on rejection")
	}
	if workspaceRepository.saved != nil || len(projectRepository.saved) != 0 {
		t.Fatal("Execute() saved aggregates on rejection")
	}
}

type mockWorkspaceRepository struct {
	workspace *domainworkspace.Workspace
	saved     *domainworkspace.Workspace
}

func (r *mockWorkspaceRepository) FindByID(context.Context, domainworkspace.WorkspaceID) (*domainworkspace.Workspace, error) {
	return r.workspace, nil
}

func (r *mockWorkspaceRepository) Save(_ context.Context, ws *domainworkspace.Workspace) error {
	r.saved = ws
	return nil
}

type mockProjectRepository struct {
	projects []*domainproject.Project
	saved    []*domainproject.Project
}

func (r *mockProjectRepository) FindByID(context.Context, domainproject.ProjectID) (*domainproject.Project, error) {
	return nil, nil
}

func (r *mockProjectRepository) FindByWorkspaceID(context.Context, domainworkspace.WorkspaceID) ([]*domainproject.Project, error) {
	return r.projects, nil
}

func (r *mockProjectRepository) Save(_ context.Context, p *domainproject.Project) error {
	r.saved = append(r.saved, p)
	return nil
}

type mockTransactionManager struct {
	callCount int
}

func (m *mockTransactionManager) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	m.callCount++
	return fn(ctx)
}
