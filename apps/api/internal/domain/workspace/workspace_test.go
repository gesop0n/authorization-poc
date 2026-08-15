package workspace_test

import (
	"errors"
	"testing"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
)

func TestNewWorkspace(t *testing.T) {
	t.Parallel()

	owner, err := user.NewUser("owner")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	ws, err := workspace.NewWorkspace("name", owner.ID())
	if err != nil {
		t.Fatalf("NewWorkspace() error = %v", err)
	}

	if ws.ID().IsZero() {
		t.Fatal("workspace ID must not be zero")
	}
	if ws.Name() != "name" {
		t.Fatalf("Name() = %q, want %q", ws.Name(), "name")
	}

	members := ws.Members()
	if len(members) != 1 {
		t.Fatalf("len(Members()) = %d, want 1", len(members))
	}
	if !members[0].UserID().Equal(owner.ID()) || members[0].Role() != workspace.WorkspaceRoleOwner {
		t.Fatal("creator must be the initial owner")
	}
}

func TestNewWorkspaceRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	owner := newTestUser(t, "owner")
	tests := []struct {
		name      string
		workspace string
		ownerID   user.UserID
		wantErr   error
	}{
		{"名前が空", "", owner.ID(), workspace.ErrWorkspaceNameRequired},
		{"名前が空白のみ", "  ", owner.ID(), workspace.ErrWorkspaceNameRequired},
		{"OwnerのIDがゼロ値", "name", user.UserID{}, workspace.ErrInvalidWorkspaceMemberUserID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := workspace.NewWorkspace(tt.workspace, tt.ownerID)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewWorkspace() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestReconstruct(t *testing.T) {
	t.Parallel()

	owner := newTestUser(t, "owner")
	member := newTestUser(t, "member")
	original, err := workspace.NewWorkspace("name", owner.ID())
	if err != nil {
		t.Fatalf("NewWorkspace() error = %v", err)
	}

	workspaceMember, err := workspace.NewWorkspaceMember(member.ID(), workspace.WorkspaceRoleMember)
	if err != nil {
		t.Fatalf("NewWorkspaceMember() error = %v", err)
	}

	members := append(original.Members(), workspaceMember)
	reconstructed, err := workspace.Reconstruct(original.ID(), original.Name(), members)
	if err != nil {
		t.Fatalf("Reconstruct() error = %v", err)
	}

	if got := len(reconstructed.Members()); got != 2 {
		t.Fatalf("len(Members()) = %d, want 2", got)
	}
}

func TestReconstructRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	owner := newTestUser(t, "owner")
	member := newTestUser(t, "member")
	original, err := workspace.NewWorkspace("name", owner.ID())
	if err != nil {
		t.Fatalf("NewWorkspace() error = %v", err)
	}

	ownerMember := original.Members()[0]
	regularMember, err := workspace.NewWorkspaceMember(member.ID(), workspace.WorkspaceRoleMember)
	if err != nil {
		t.Fatalf("NewWorkspaceMember() error = %v", err)
	}

	tests := []struct {
		name      string
		id        workspace.WorkspaceID
		workspace string
		members   []workspace.WorkspaceMember
		wantErr   error
	}{
		{
			name:      "IDがゼロ値",
			workspace: "name",
			members:   []workspace.WorkspaceMember{ownerMember},
			wantErr:   workspace.ErrInvalidWorkspaceID,
		},
		{
			name:      "名前が空白のみ",
			id:        original.ID(),
			workspace: "   ",
			members:   []workspace.WorkspaceMember{ownerMember},
			wantErr:   workspace.ErrWorkspaceNameRequired,
		},
		{
			name:      "メンバーのユーザーIDがゼロ値",
			id:        original.ID(),
			workspace: "name",
			members:   []workspace.WorkspaceMember{{}},
			wantErr:   workspace.ErrInvalidWorkspaceMemberUserID,
		},
		{
			name:      "ユーザーIDが重複",
			id:        original.ID(),
			workspace: "name",
			members:   []workspace.WorkspaceMember{ownerMember, ownerMember},
			wantErr:   workspace.ErrWorkspaceMemberAlreadyExists,
		},
		{
			name:      "Ownerが存在しない",
			id:        original.ID(),
			workspace: "name",
			members:   []workspace.WorkspaceMember{regularMember},
			wantErr:   workspace.ErrWorkspaceMustHaveOwner,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := workspace.Reconstruct(tt.id, tt.workspace, tt.members)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Reconstruct() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewWorkspaceMemberRejectsInvalidRole(t *testing.T) {
	t.Parallel()

	member := newTestUser(t, "member")
	_, err := workspace.NewWorkspaceMember(member.ID(), workspace.WorkspaceRole("invalid"))
	if !errors.Is(err, workspace.ErrInvalidWorkspaceRole) {
		t.Fatalf("NewWorkspaceMember() error = %v, want %v", err, workspace.ErrInvalidWorkspaceRole)
	}
}

func TestNewWorkspaceMember(t *testing.T) {
	t.Parallel()

	member := newTestUser(t, "member")
	got, err := workspace.NewWorkspaceMember(member.ID(), workspace.WorkspaceRoleAdmin)
	if err != nil {
		t.Fatalf("NewWorkspaceMember() error = %v", err)
	}
	if !got.UserID().Equal(member.ID()) || got.Role() != workspace.WorkspaceRoleAdmin {
		t.Fatal("NewWorkspaceMember() did not preserve values")
	}
}

func TestNewWorkspaceMemberRejectsZeroUserID(t *testing.T) {
	t.Parallel()

	_, err := workspace.NewWorkspaceMember(user.UserID{}, workspace.WorkspaceRoleMember)
	if !errors.Is(err, workspace.ErrInvalidWorkspaceMemberUserID) {
		t.Fatalf("NewWorkspaceMember() error = %v, want %v", err, workspace.ErrInvalidWorkspaceMemberUserID)
	}
}

func TestWorkspaceAddMember(t *testing.T) {
	t.Parallel()

	ws, _ := newTestWorkspace(t)
	member := newTestUser(t, "member")
	if err := ws.AddMember(member.ID(), workspace.WorkspaceRoleAdmin); err != nil {
		t.Fatalf("AddMember() error = %v", err)
	}

	members := ws.Members()
	if len(members) != 2 {
		t.Fatalf("len(Members()) = %d, want 2", len(members))
	}
	if !members[1].UserID().Equal(member.ID()) || members[1].Role() != workspace.WorkspaceRoleAdmin {
		t.Fatal("added member has unexpected values")
	}
}

func TestWorkspaceAddMemberRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(t *testing.T, ws *workspace.Workspace) user.UserID
		role    workspace.WorkspaceRole
		wantErr error
	}{
		{
			name:    "ユーザーIDがゼロ値",
			prepare: func(*testing.T, *workspace.Workspace) user.UserID { return user.UserID{} },
			role:    workspace.WorkspaceRoleMember,
			wantErr: workspace.ErrInvalidWorkspaceMemberUserID,
		},
		{
			name: "ロールが不正",
			prepare: func(t *testing.T, _ *workspace.Workspace) user.UserID {
				return newTestUser(t, "member").ID()
			},
			role:    workspace.WorkspaceRole("invalid"),
			wantErr: workspace.ErrInvalidWorkspaceRole,
		},
		{
			name: "メンバーが重複",
			prepare: func(t *testing.T, ws *workspace.Workspace) user.UserID {
				member := newTestUser(t, "member")
				if err := ws.AddMember(member.ID(), workspace.WorkspaceRoleMember); err != nil {
					t.Fatalf("AddMember() setup error = %v", err)
				}
				return member.ID()
			},
			role:    workspace.WorkspaceRoleAdmin,
			wantErr: workspace.ErrWorkspaceMemberAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ws, _ := newTestWorkspace(t)
			userID := tt.prepare(t, ws)
			before := ws.Members()

			err := ws.AddMember(userID, tt.role)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("AddMember() error = %v, want %v", err, tt.wantErr)
			}
			if len(ws.Members()) != len(before) {
				t.Fatal("AddMember() changed state on error")
			}
		})
	}
}

func TestWorkspaceRemoveMember(t *testing.T) {
	t.Parallel()

	ws, _ := newTestWorkspace(t)
	member := newTestUser(t, "member")
	if err := ws.AddMember(member.ID(), workspace.WorkspaceRoleMember); err != nil {
		t.Fatalf("AddMember() error = %v", err)
	}
	if err := ws.RemoveMember(member.ID()); err != nil {
		t.Fatalf("RemoveMember() error = %v", err)
	}
	if got := len(ws.Members()); got != 1 {
		t.Fatalf("len(Members()) = %d, want 1", got)
	}
}

func TestWorkspaceRemoveMemberPreservesOwner(t *testing.T) {
	t.Parallel()

	ws, owner := newTestWorkspace(t)
	err := ws.RemoveMember(owner.ID())
	if !errors.Is(err, workspace.ErrWorkspaceMustHaveOwner) {
		t.Fatalf("RemoveMember() error = %v, want %v", err, workspace.ErrWorkspaceMustHaveOwner)
	}
	if got := len(ws.Members()); got != 1 {
		t.Fatalf("len(Members()) = %d, want 1", got)
	}
}

func TestWorkspaceRemoveMemberAllowsMultipleOwners(t *testing.T) {
	t.Parallel()

	ws, owner := newTestWorkspace(t)
	secondOwner := newTestUser(t, "second owner")
	if err := ws.AddMember(secondOwner.ID(), workspace.WorkspaceRoleOwner); err != nil {
		t.Fatalf("AddMember() error = %v", err)
	}
	if err := ws.RemoveMember(owner.ID()); err != nil {
		t.Fatalf("RemoveMember() error = %v", err)
	}
}

func TestWorkspaceRemoveMemberReturnsNotFound(t *testing.T) {
	t.Parallel()

	ws, _ := newTestWorkspace(t)
	unknown := newTestUser(t, "unknown")
	err := ws.RemoveMember(unknown.ID())
	if !errors.Is(err, workspace.ErrWorkspaceMemberNotFound) {
		t.Fatalf("RemoveMember() error = %v, want %v", err, workspace.ErrWorkspaceMemberNotFound)
	}
}

func TestWorkspaceChangeMemberRole(t *testing.T) {
	t.Parallel()

	ws, _ := newTestWorkspace(t)
	member := newTestUser(t, "member")
	if err := ws.AddMember(member.ID(), workspace.WorkspaceRoleMember); err != nil {
		t.Fatalf("AddMember() error = %v", err)
	}
	if err := ws.ChangeMemberRole(member.ID(), workspace.WorkspaceRoleAdmin); err != nil {
		t.Fatalf("ChangeMemberRole() error = %v", err)
	}
	if got := ws.Members()[1].Role(); got != workspace.WorkspaceRoleAdmin {
		t.Fatalf("Role() = %q, want %q", got, workspace.WorkspaceRoleAdmin)
	}
}

func TestWorkspaceChangeMemberRoleRejectsInvalidOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  func(t *testing.T, ws *workspace.Workspace, owner *user.User) user.UserID
		role    workspace.WorkspaceRole
		wantErr error
	}{
		{
			name:    "最後のOwnerを降格",
			userID:  func(_ *testing.T, _ *workspace.Workspace, owner *user.User) user.UserID { return owner.ID() },
			role:    workspace.WorkspaceRoleAdmin,
			wantErr: workspace.ErrWorkspaceMustHaveOwner,
		},
		{
			name: "存在しないメンバー",
			userID: func(t *testing.T, _ *workspace.Workspace, _ *user.User) user.UserID {
				return newTestUser(t, "unknown").ID()
			},
			role:    workspace.WorkspaceRoleMember,
			wantErr: workspace.ErrWorkspaceMemberNotFound,
		},
		{
			name:    "不正なロール",
			userID:  func(_ *testing.T, _ *workspace.Workspace, owner *user.User) user.UserID { return owner.ID() },
			role:    workspace.WorkspaceRole("invalid"),
			wantErr: workspace.ErrInvalidWorkspaceRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ws, owner := newTestWorkspace(t)
			before := ws.Members()
			err := ws.ChangeMemberRole(tt.userID(t, ws, owner), tt.role)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ChangeMemberRole() error = %v, want %v", err, tt.wantErr)
			}
			if ws.Members()[0].Role() != before[0].Role() {
				t.Fatal("ChangeMemberRole() changed state on error")
			}
		})
	}
}

func TestWorkspaceChangeMemberRoleAllowsSameRole(t *testing.T) {
	t.Parallel()

	ws, owner := newTestWorkspace(t)
	if err := ws.ChangeMemberRole(owner.ID(), workspace.WorkspaceRoleOwner); err != nil {
		t.Fatalf("ChangeMemberRole() error = %v", err)
	}
}

func TestWorkspaceMembersReturnsCopy(t *testing.T) {
	t.Parallel()

	ws, owner := newTestWorkspace(t)
	members := ws.Members()
	members[0] = workspace.WorkspaceMember{}

	if !ws.Members()[0].UserID().Equal(owner.ID()) {
		t.Fatal("Members() exposed internal state")
	}
}

func TestReconstructCopiesMembers(t *testing.T) {
	t.Parallel()

	original, _ := newTestWorkspace(t)
	members := original.Members()
	reconstructed, err := workspace.Reconstruct(original.ID(), original.Name(), members)
	if err != nil {
		t.Fatalf("Reconstruct() error = %v", err)
	}
	members[0] = workspace.WorkspaceMember{}

	if reconstructed.Members()[0].UserID().IsZero() {
		t.Fatal("Reconstruct() retained the input slice")
	}
}

func TestNewWorkspaceRole(t *testing.T) {
	t.Parallel()

	for _, role := range []workspace.WorkspaceRole{
		workspace.WorkspaceRoleOwner,
		workspace.WorkspaceRoleAdmin,
		workspace.WorkspaceRoleMember,
	} {
		role := role
		t.Run(string(role), func(t *testing.T) {
			t.Parallel()
			got, err := workspace.NewWorkspaceRole(string(role))
			if err != nil || got != role {
				t.Fatalf("NewWorkspaceRole() = %q, %v; want %q, nil", got, err, role)
			}
		})
	}
}

func TestNewWorkspaceRoleRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	_, err := workspace.NewWorkspaceRole("invalid")
	if !errors.Is(err, workspace.ErrInvalidWorkspaceRole) {
		t.Fatalf("NewWorkspaceRole() error = %v, want %v", err, workspace.ErrInvalidWorkspaceRole)
	}
}

func newTestWorkspace(t *testing.T) (*workspace.Workspace, *user.User) {
	t.Helper()

	owner := newTestUser(t, "owner")
	ws, err := workspace.NewWorkspace("name", owner.ID())
	if err != nil {
		t.Fatalf("NewWorkspace() error = %v", err)
	}

	return ws, owner
}

func newTestUser(t *testing.T, name string) *user.User {
	t.Helper()

	u, err := user.NewUser(name)
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	return u
}
