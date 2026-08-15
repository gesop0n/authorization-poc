package user_test

import (
	"errors"
	"testing"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	u, err := user.NewUser("name")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}
	if u.ID().IsZero() || u.Name() != "name" || u.Status() != user.UserStatusActive {
		t.Fatal("NewUser() returned an unexpected user")
	}
}

func TestNewUserRejectsBlankName(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"", "   "} {
		if _, err := user.NewUser(name); !errors.Is(err, user.ErrUsernameRequired) {
			t.Fatalf("NewUser(%q) error = %v, want %v", name, err, user.ErrUsernameRequired)
		}
	}
}

func TestReconstructUser(t *testing.T) {
	t.Parallel()

	original := newTestUser(t)
	u, err := user.Reconstruct(original.ID(), "restored", user.UserStatusSuspended)
	if err != nil {
		t.Fatalf("Reconstruct() error = %v", err)
	}
	if u.Name() != "restored" || u.Status() != user.UserStatusSuspended {
		t.Fatal("Reconstruct() did not restore all values")
	}
}

func TestReconstructUserRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	id := newTestUser(t).ID()
	tests := []struct {
		name     string
		id       user.UserID
		username string
		status   user.UserStatus
		wantErr  error
	}{
		{"IDがゼロ値", user.UserID{}, "name", user.UserStatusActive, user.ErrInvalidUserID},
		{"名前が空白のみ", id, "  ", user.UserStatusActive, user.ErrUsernameRequired},
		{"ステータスが不正", id, "name", user.UserStatus("invalid"), user.ErrInvalidUserStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := user.Reconstruct(tt.id, tt.username, tt.status)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Reconstruct() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserSuspendAndReactivate(t *testing.T) {
	t.Parallel()

	u := newTestUser(t)
	u.Suspend()
	if u.Status() != user.UserStatusSuspended {
		t.Fatalf("Status() = %q, want %q", u.Status(), user.UserStatusSuspended)
	}
	u.Reactivate()
	if u.Status() != user.UserStatusActive {
		t.Fatalf("Status() = %q, want %q", u.Status(), user.UserStatusActive)
	}
}

func newTestUser(t *testing.T) *user.User {
	t.Helper()
	u, err := user.NewUser("name")
	if err != nil {
		t.Fatalf("NewUser() error = %v", err)
	}
	return u
}
