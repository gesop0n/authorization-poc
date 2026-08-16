package project_test

import (
	"errors"
	"testing"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
	"github.com/gesop0n/authorization-poc/apps/api/internal/testutil"
)

func TestNewProject(t *testing.T) {
	t.Parallel()
	p, ws, admin := newTestProject(t)
	if p.ID().IsZero() || p.Name() != "name" || !p.WorkspaceID().Equal(ws.ID()) || p.Status() != project.ProjectStatusActive {
		t.Fatal("NewProject() returned unexpected values")
	}
	members := p.Members()
	if len(members) != 1 || !members[0].UserID().Equal(admin.ID()) || members[0].Role() != project.ProjectRoleAdmin {
		t.Fatal("creator must be the initial project admin")
	}
}

func TestNewProjectRejectsInvalidValues(t *testing.T) {
	t.Parallel()
	ws, admin := newTestWorkspace(t)
	tests := []struct {
		name, projectName string
		workspaceID       workspace.WorkspaceID
		adminID           user.UserID
		wantErr           error
	}{
		{"名前が空", "", ws.ID(), admin.ID(), project.ErrProjectNameRequired},
		{"名前が空白のみ", "  ", ws.ID(), admin.ID(), project.ErrProjectNameRequired},
		{"Workspace IDがゼロ値", "name", workspace.WorkspaceID{}, admin.ID(), workspace.ErrInvalidWorkspaceID},
		{"Admin IDがゼロ値", "name", ws.ID(), user.UserID{}, project.ErrInvalidProjectMemberUserID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := project.NewProject(tt.projectName, tt.workspaceID, tt.adminID)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewProject() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestReconstructProject(t *testing.T) {
	t.Parallel()
	original, ws, _ := newTestProject(t)
	p, err := project.Reconstruct(original.ID(), "restored", ws.ID(), project.ProjectStatusArchived, original.Members())
	if err != nil {
		t.Fatalf("Reconstruct() error = %v", err)
	}
	if p.Name() != "restored" || p.Status() != project.ProjectStatusArchived || len(p.Members()) != 1 {
		t.Fatal("Reconstruct() did not restore all values")
	}
}

func TestReconstructProjectRejectsInvalidValues(t *testing.T) {
	t.Parallel()
	original, ws, _ := newTestProject(t)
	tests := []struct {
		name, projectName string
		id                project.ProjectID
		workspaceID       workspace.WorkspaceID
		status            project.ProjectStatus
		members           []project.ProjectMember
		wantErr           error
	}{
		{"IDがゼロ値", "name", project.ProjectID{}, ws.ID(), project.ProjectStatusActive, original.Members(), project.ErrInvalidProjectID},
		{"名前が空白のみ", "  ", original.ID(), ws.ID(), project.ProjectStatusActive, original.Members(), project.ErrProjectNameRequired},
		{"Workspace IDがゼロ値", "name", original.ID(), workspace.WorkspaceID{}, project.ProjectStatusActive, original.Members(), workspace.ErrInvalidWorkspaceID},
		{"Statusが不正", "name", original.ID(), ws.ID(), project.ProjectStatus("invalid"), original.Members(), project.ErrInvalidProjectStatus},
		{"Adminが不在", "name", original.ID(), ws.ID(), project.ProjectStatusActive, nil, project.ErrProjectMustHaveAdmin},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := project.Reconstruct(tt.id, tt.projectName, tt.workspaceID, tt.status, tt.members)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Reconstruct() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectMemberLifecycle(t *testing.T) {
	t.Parallel()
	p, _, _ := newTestProject(t)
	member := newTestUser(t, "member")
	if err := p.AddMember(member.ID(), project.ProjectRoleViewer); err != nil {
		t.Fatalf("AddMember() error = %v", err)
	}
	if err := p.ChangeMemberRole(member.ID(), project.ProjectRoleEditor); err != nil {
		t.Fatalf("ChangeMemberRole() error = %v", err)
	}
	if !p.HasMember(member.ID()) || p.Members()[1].Role() != project.ProjectRoleEditor {
		t.Fatal("ChangeMemberRole() did not change the role")
	}
	if err := p.RemoveMember(member.ID()); err != nil {
		t.Fatalf("RemoveMember() error = %v", err)
	}
	if p.HasMember(member.ID()) {
		t.Fatal("RemoveMember() did not remove the member")
	}
}

func TestProjectRejectsDuplicateMember(t *testing.T) {
	t.Parallel()
	p, _, admin := newTestProject(t)
	before := p.Members()
	err := p.AddMember(admin.ID(), project.ProjectRoleViewer)
	if !errors.Is(err, project.ErrProjectMemberAlreadyExists) || len(p.Members()) != len(before) {
		t.Fatal("AddMember() must reject a duplicate without changing state")
	}
}

func TestProjectPreservesLastAdmin(t *testing.T) {
	t.Parallel()
	p, _, admin := newTestProject(t)
	if err := p.ChangeMemberRole(admin.ID(), project.ProjectRoleEditor); !errors.Is(err, project.ErrProjectMustHaveAdmin) {
		t.Fatalf("ChangeMemberRole() error = %v, want %v", err, project.ErrProjectMustHaveAdmin)
	}
	if err := p.RemoveMember(admin.ID()); !errors.Is(err, project.ErrProjectMustHaveAdmin) {
		t.Fatalf("RemoveMember() error = %v, want %v", err, project.ErrProjectMustHaveAdmin)
	}
	if len(p.Members()) != 1 || p.Members()[0].Role() != project.ProjectRoleAdmin {
		t.Fatal("last admin changed on error")
	}
}

func TestProjectMembersReturnsCopy(t *testing.T) {
	t.Parallel()
	p, _, admin := newTestProject(t)
	members := p.Members()
	members[0] = project.ProjectMember{}
	if !p.Members()[0].UserID().Equal(admin.ID()) {
		t.Fatal("Members() exposed internal state")
	}
}

func TestProjectRenameArchiveAndRestore(t *testing.T) {
	t.Parallel()
	p, _, _ := newTestProject(t)
	if err := p.Rename("renamed"); err != nil || p.Name() != "renamed" {
		t.Fatalf("Rename() = %v, name = %q", err, p.Name())
	}
	if err := p.Rename("  "); !errors.Is(err, project.ErrProjectNameRequired) || p.Name() != "renamed" {
		t.Fatal("Rename() must reject a blank name without changing state")
	}
	p.Archive()
	if p.Status() != project.ProjectStatusArchived {
		t.Fatal("Archive() did not archive the project")
	}
	p.Restore()
	if p.Status() != project.ProjectStatusActive {
		t.Fatal("Restore() did not restore the project")
	}
}

func newTestProject(t *testing.T) (*project.Project, *workspace.Workspace, *user.User) {
	t.Helper()
	return testutil.NewProjectWithWorkspace(t, "name", "workspace", "owner")
}

func newTestWorkspace(t *testing.T) (*workspace.Workspace, *user.User) {
	t.Helper()
	return testutil.NewWorkspaceWithOwner(t, "workspace", "owner")
}

func newTestUser(t *testing.T, name string) *user.User {
	t.Helper()
	return testutil.NewUser(t, name)
}
