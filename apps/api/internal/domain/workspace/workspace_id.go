package workspace

import (
	"fmt"

	"github.com/google/uuid"
)

type WorkspaceID struct {
	value uuid.UUID
}

func newWorkspaceID() (WorkspaceID, error) {
	value, err := uuid.NewV7()
	if err != nil {
		return WorkspaceID{}, fmt.Errorf("generate workspace ID: %w", err)
	}

	return WorkspaceID{value: value}, nil
}

func parseWorkspaceID(value string) (WorkspaceID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return WorkspaceID{}, fmt.Errorf("%w: %q", ErrInvalidWorkspaceID, value)
	}

	if parsed.Version() != 7 {
		return WorkspaceID{}, fmt.Errorf("%w: %q", ErrInvalidWorkspaceID, value)
	}

	return WorkspaceID{value: parsed}, nil
}

func (id WorkspaceID) String() string {
	return id.value.String()
}

func (id WorkspaceID) Equal(other WorkspaceID) bool {
	return id.value == other.value
}

func (id WorkspaceID) IsZero() bool {
	return id.value == uuid.Nil
}
