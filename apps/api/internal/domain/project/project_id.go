package project

import (
	"fmt"

	"github.com/google/uuid"
)

type ProjectID struct {
	value uuid.UUID
}

func newProjectID() (ProjectID, error) {
	value, err := uuid.NewV7()
	if err != nil {
		return ProjectID{}, fmt.Errorf("generate project ID: %w", err)
	}

	return ProjectID{value: value}, nil
}

func parseProjectID(value string) (ProjectID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return ProjectID{}, fmt.Errorf("%w: %q", ErrInvalidProjectID, value)
	}

	if parsed.Version() != 7 {
		return ProjectID{}, fmt.Errorf("%w: %q", ErrInvalidProjectID, value)
	}

	return ProjectID{value: parsed}, nil
}

func (id ProjectID) String() string {
	return id.value.String()
}

func (id ProjectID) Equal(other ProjectID) bool {
	return id.value == other.value
}

func (id ProjectID) IsZero() bool {
	return id.value == uuid.Nil
}
