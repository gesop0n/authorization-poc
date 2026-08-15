package workspace

import "fmt"

// WorkspaceRole は、ワークスペース内のロール。
type WorkspaceRole string

const (
	// WorkspaceRoleOwner は、所有者ロール。
	WorkspaceRoleOwner WorkspaceRole = "owner"
	// WorkspaceRoleAdmin は、管理者ロール。
	WorkspaceRoleAdmin WorkspaceRole = "admin"
	// WorkspaceRoleMember は、一般メンバーロール。
	WorkspaceRoleMember WorkspaceRole = "member"
)

// NewWorkspaceRole は、検証済みのロールを生成する。
func NewWorkspaceRole(value string) (WorkspaceRole, error) {
	role := WorkspaceRole(value)

	switch role {
	case WorkspaceRoleOwner, WorkspaceRoleAdmin, WorkspaceRoleMember:
		return role, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidWorkspaceRole, value)
	}
}

// IsValid は、ロールが有効かを返す。
func (role WorkspaceRole) IsValid() bool {
	switch role {
	case WorkspaceRoleOwner, WorkspaceRoleAdmin, WorkspaceRoleMember:
		return true
	default:
		return false
	}
}
