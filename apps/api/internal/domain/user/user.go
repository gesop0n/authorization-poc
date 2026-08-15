package user

import "fmt"

type User struct {
	id      UserID
	profile UserProfile
}

func NewUser(name string) (*User, error) {
	id, err := newUserID()
	if err != nil {
		return nil, fmt.Errorf("create User: %w", err)
	}

	return &User{
		id:      id,
		profile: NewUserProfile(id, name),
	}, nil
}

func Reconstruct(id UserID) (*User, error) {
	if id.IsZero() {
		return nil, ErrInvalidUserID
	}
	return &User{id: id}, nil
}

func (w *User) ID() UserID {
	return w.id
}
