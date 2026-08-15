package project

import (
	"errors"
	"testing"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
)

func TestReconstructRejectsInvalidMembers(t *testing.T) {
	t.Parallel()
	p, ws, admin := newInternalTestProject(t)
	tests := []struct {
		name    string
		members []ProjectMember
		wantErr error
	}{
		{"User IDがゼロ値", []ProjectMember{newProjectMember(user.UserID{}, ProjectRoleAdmin)}, ErrInvalidProjectMemberUserID},
		{"Roleが不正", []ProjectMember{newProjectMember(admin.ID(), ProjectRole("invalid"))}, ErrInvalidProjectRole},
		{"Memberが重複", []ProjectMember{newProjectMember(admin.ID(), ProjectRoleAdmin), newProjectMember(admin.ID(), ProjectRoleAdmin)}, ErrProjectMemberAlreadyExists},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := Reconstruct(p.ID(), p.Name(), ws.ID(), p.Status(), tt.members)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Reconstruct() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func newInternalTestProject(t *testing.T) (*Project, *workspace.Workspace, *user.User) {
	t.Helper()
	admin, err := user.NewUser("admin")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}
	ws, err := workspace.NewWorkspace("workspace", admin.ID())
	if err != nil {
		t.Fatalf("NewWorkspace() error = %v", err)
	}
	p, err := NewProject("project", ws.ID(), admin.ID())
	if err != nil {
		t.Fatalf("NewProject() error = %v", err)
	}
	return p, ws, admin
}
