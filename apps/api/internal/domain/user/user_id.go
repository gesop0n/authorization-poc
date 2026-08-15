package user

import (
	"fmt"

	"github.com/google/uuid"
)

type UserID struct {
	value uuid.UUID
}

func newUserID() (UserID, error) {
	value, err := uuid.NewV7()
	if err != nil {
		return UserID{}, fmt.Errorf("generate user ID: %w", err)
	}

	return UserID{value: value}, nil
}

func parseUserID(value string) (UserID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return UserID{}, fmt.Errorf("%w: %q", ErrInvalidUserID, value)
	}

	if parsed.Version() != 7 {
		return UserID{}, fmt.Errorf("%w: %q", ErrInvalidUserID, value)
	}

	return UserID{value: parsed}, nil
}

func (id UserID) String() string {
	return id.value.String()
}

func (id UserID) Equal(other UserID) bool {
	return id.value == other.value
}

func (id UserID) IsZero() bool {
	return id.value == uuid.Nil
}
