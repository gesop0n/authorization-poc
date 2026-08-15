package projectmembership

import (
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
)

type ProjectMembership struct {
	projectID   project.ProjectID
	userID      user.UserID
	projectRole ProjectRole
}
