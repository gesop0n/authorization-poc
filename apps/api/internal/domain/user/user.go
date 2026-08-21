package user

import (
	"fmt"
	"strings"
)

type User struct {
	id      UserID
	profile UserProfile
	status  UserStatus
}

func NewUser(name string) (*User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrUsernameRequired
	}

	id, err := newUserID()
	if err != nil {
		return nil, fmt.Errorf("create User: %w", err)
	}

	return &User{
		id:      id,
		profile: NewUserProfile(id, name),
		status:  UserStatusActive,
	}, nil
}

func Reconstruct(id UserID, name string, status UserStatus) (*User, error) {
	if id.IsZero() {
		return nil, ErrInvalidUserID
	}

	if strings.TrimSpace(name) == "" {
		return nil, ErrUsernameRequired
	}

	if !status.IsValid() {
		return nil, ErrInvalidUserStatus
	}

	return &User{id: id, profile: NewUserProfile(id, name), status: status}, nil
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Suspend() {
	u.status = UserStatusSuspended
}

func (u *User) Reactivate() {
	u.status = UserStatusActive
}

func (u *User) DisplayName() string {
	return u.profile.displayName
}

func (u *User) Status() UserStatus {
	return u.status
}
