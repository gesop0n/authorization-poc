package projectusecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
	"github.com/gesop0n/authorization-poc/apps/api/internal/testutil"
	projectusecase "github.com/gesop0n/authorization-poc/apps/api/internal/usecase/project"
)

// workspaceに所属するユーザーをProjectへ追加し、
// 変更後のProjectがRepositoryへ保存されることを確認する。
func TestAddMemberUseCase(t *testing.T) {
	t.Parallel()
	ws, p, member := newTestEntities(t)
	workspaceRepository := &mockWorkspaceRepository{workspace: ws}
	projectRepository := &mockProjectRepository{project: p}
	uc := projectusecase.NewAddMemberUseCase(workspaceRepository, projectRepository)

	err := uc.Execute(context.Background(), projectusecase.AddMemberInput{
		WorkspaceID: ws.ID(),
		ProjectID:   p.ID(),
		UserID:      member.ID(),
		Role:        project.ProjectRoleEditor,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !p.HasMember(member.ID()) || projectRepository.saved != p {
		t.Fatal("Execute() did not add and save the project member")
	}
}

// Workspaceに所属していないユーザーのProject参加を拒否し、Projectを変更・保存しないことを確認する。
func TestAddMemberUseCaseRejectsNonWorkspaceMember(t *testing.T) {
	t.Parallel()
	ws, p, _ := newTestEntities(t)
	nonMember := testutil.NewUser(t, "non-member")
	projectRepository := &mockProjectRepository{project: p}
	uc := projectusecase.NewAddMemberUseCase(&mockWorkspaceRepository{workspace: ws}, projectRepository)

	err := uc.Execute(context.Background(), projectusecase.AddMemberInput{
		WorkspaceID: ws.ID(), ProjectID: p.ID(), UserID: nonMember.ID(), Role: project.ProjectRoleViewer,
	})
	if !errors.Is(err, projectusecase.ErrUserNotWorkspaceMember) {
		t.Fatalf("Execute() error = %v, want %v", err, projectusecase.ErrUserNotWorkspaceMember)
	}
	if p.HasMember(nonMember.ID()) || projectRepository.saved != nil {
		t.Fatal("Execute() changed or saved the project on rejection")
	}
}

// 指定したWorkspaceとProjectの対応関係が一致しない場合に参加を拒否し、Projectを変更・保存しないことを確認する。
func TestAddMemberUseCaseRejectsWorkspaceMismatch(t *testing.T) {
	t.Parallel()
	_, p, member := newTestEntities(t)
	otherOwner := testutil.NewUser(t, "other owner")
	otherWorkspace, err := workspace.NewWorkspace("other", otherOwner.ID())
	if err != nil {
		t.Fatal(err)
	}
	if err := otherWorkspace.AddMember(member.ID(), workspace.WorkspaceRoleMember); err != nil {
		t.Fatal(err)
	}
	projectRepository := &mockProjectRepository{project: p}
	uc := projectusecase.NewAddMemberUseCase(&mockWorkspaceRepository{workspace: otherWorkspace}, projectRepository)

	err = uc.Execute(context.Background(), projectusecase.AddMemberInput{
		WorkspaceID: otherWorkspace.ID(), ProjectID: p.ID(), UserID: member.ID(), Role: project.ProjectRoleViewer,
	})
	if !errors.Is(err, projectusecase.ErrProjectWorkspaceMismatch) {
		t.Fatalf("Execute() error = %v, want %v", err, projectusecase.ErrProjectWorkspaceMismatch)
	}
	if p.HasMember(member.ID()) || projectRepository.saved != nil {
		t.Fatal("Execute() changed or saved the project on mismatch")
	}
}

func newTestEntities(t *testing.T) (*workspace.Workspace, *project.Project, *user.User) {
	t.Helper()
	p, ws, _ := testutil.NewProjectWithWorkspace(t, "project", "workspace", "owner")
	member := testutil.NewUser(t, "member")
	if err := ws.AddMember(member.ID(), workspace.WorkspaceRoleMember); err != nil {
		t.Fatal(err)
	}
	return ws, p, member
}

type mockWorkspaceRepository struct {
	workspace *workspace.Workspace
	err       error
}

func (r *mockWorkspaceRepository) FindByID(context.Context, workspace.WorkspaceID) (*workspace.Workspace, error) {
	return r.workspace, r.err
}

func (r *mockWorkspaceRepository) Save(context.Context, *workspace.Workspace) error { return nil }

type mockProjectRepository struct {
	project *project.Project
	saved   *project.Project
	err     error
}

func (r *mockProjectRepository) FindByID(context.Context, project.ProjectID) (*project.Project, error) {
	return r.project, r.err
}

func (r *mockProjectRepository) FindByWorkspaceID(context.Context, workspace.WorkspaceID) ([]*project.Project, error) {
	return nil, nil
}

func (r *mockProjectRepository) Save(_ context.Context, p *project.Project) error {
	r.saved = p
	return nil
}
