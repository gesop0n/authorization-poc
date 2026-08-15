package workspace

import (
	"errors"
	"testing"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
)

func TestReconstructRejectsInvalidRole(t *testing.T) {
	t.Parallel()

	owner, err := user.NewUser("owner")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}

	original, err := NewWorkspace("name", owner.ID())
	if err != nil {
		t.Fatalf("NewWorkspace() error = %v", err)
	}

	members := []WorkspaceMember{
		newWorkspaceMember(owner.ID(), WorkspaceRole("invalid")),
	}
	_, err = Reconstruct(original.ID(), original.Name(), members)
	if !errors.Is(err, ErrInvalidWorkspaceRole) {
		t.Fatalf("Reconstruct() error = %v, want %v", err, ErrInvalidWorkspaceRole)
	}
}
