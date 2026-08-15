package workspace

import "github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"

// WorkspaceMember は、ワークスペースのメンバー。
type WorkspaceMember struct {
	userID user.UserID
	role   WorkspaceRole
}

// NewWorkspaceMember は、検証済みのメンバーを生成する。
func NewWorkspaceMember(userID user.UserID, role WorkspaceRole) (WorkspaceMember, error) {
	if userID.IsZero() {
		return WorkspaceMember{}, ErrInvalidWorkspaceMemberUserID
	}

	if !role.IsValid() {
		return WorkspaceMember{}, ErrInvalidWorkspaceRole
	}

	return newWorkspaceMember(userID, role), nil
}

func newWorkspaceMember(userID user.UserID, role WorkspaceRole) WorkspaceMember {
	return WorkspaceMember{userID: userID, role: role}
}

// UserID は、ユーザーIDを返す。
func (wm WorkspaceMember) UserID() user.UserID {
	return wm.userID
}

// Role は、ロールを返す。
func (wm WorkspaceMember) Role() WorkspaceRole {
	return wm.role
}
