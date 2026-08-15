package workspace

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
)

// ワークスペース集約のルート。
type Workspace struct {
	id      WorkspaceID
	name    string
	members []WorkspaceMember
}

// Owner付きのワークスペースを生成する。
func NewWorkspace(name string, ownerID user.UserID) (*Workspace, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrWorkspaceNameRequired
	}

	if ownerID.IsZero() {
		return nil, ErrInvalidWorkspaceMemberUserID
	}

	id, err := newWorkspaceID()
	if err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}

	return &Workspace{
		id:   id,
		name: name,
		members: []WorkspaceMember{
			newWorkspaceMember(ownerID, WorkspaceRoleOwner),
		},
	}, nil
}

// 永続化された値からワークスペースを再構築する。
func Reconstruct(id WorkspaceID, name string, members []WorkspaceMember) (*Workspace, error) {
	if id.IsZero() {
		return nil, ErrInvalidWorkspaceID
	}

	if strings.TrimSpace(name) == "" {
		return nil, ErrWorkspaceNameRequired
	}

	seen := make(map[string]struct{}, len(members))
	ownerCount := 0

	for _, member := range members {
		if member.UserID().IsZero() {
			return nil, ErrInvalidWorkspaceMemberUserID
		}

		if !member.Role().IsValid() {
			return nil, ErrInvalidWorkspaceRole
		}

		userID := member.UserID().String()
		if _, exists := seen[userID]; exists {
			return nil, ErrWorkspaceMemberAlreadyExists
		}
		seen[userID] = struct{}{}

		if member.Role() == WorkspaceRoleOwner {
			ownerCount++
		}
	}

	// Ownerが１人以上いるか確認
	if ownerCount == 0 {
		return nil, ErrWorkspaceMustHaveOwner
	}

	return &Workspace{
		id:      id,
		name:    name,
		members: slices.Clone(members),
	}, nil
}

func (w *Workspace) ID() WorkspaceID {
	return w.id
}

func (w *Workspace) Name() string {
	return w.name
}

// ワークスペースにメンバーを追加する。
func (w *Workspace) AddMember(userID user.UserID, role WorkspaceRole) error {
	if userID.IsZero() {
		return ErrInvalidWorkspaceMemberUserID
	}

	if !role.IsValid() {
		return ErrInvalidWorkspaceRole
	}

	if w.hasMember(userID) {
		return ErrWorkspaceMemberAlreadyExists
	}

	w.members = append(w.members, newWorkspaceMember(userID, role))

	return nil
}

// ワークスペースからメンバーを削除する。
func (w *Workspace) RemoveMember(userID user.UserID) error {
	index := w.memberIndex(userID)
	if index == -1 {
		return ErrWorkspaceMemberNotFound
	}

	target := w.members[index]

	if target.Role() == WorkspaceRoleOwner && w.ownerCount() == 1 {
		return ErrWorkspaceMustHaveOwner
	}

	w.members = append(w.members[:index:index], w.members[index+1:]...)

	return nil
}

// ワークスペースのメンバーのロールを変更する。
func (w *Workspace) ChangeMemberRole(userID user.UserID, role WorkspaceRole) error {
	if !role.IsValid() {
		return ErrInvalidWorkspaceRole
	}

	index := w.memberIndex(userID)
	if index == -1 {
		return ErrWorkspaceMemberNotFound
	}

	currentRole := w.members[index].Role()

	if currentRole == role {
		return nil
	}

	if currentRole == WorkspaceRoleOwner && role != WorkspaceRoleOwner && w.ownerCount() == 1 {
		return ErrWorkspaceMustHaveOwner
	}

	w.members[index] = newWorkspaceMember(userID, role)

	return nil
}

// メンバーのコピーを返す。
func (w *Workspace) Members() []WorkspaceMember {
	return slices.Clone(w.members)
}

// HasMember は、ユーザーがワークスペースに所属しているかを返す。
func (w *Workspace) HasMember(userID user.UserID) bool {
	return w.hasMember(userID)
}

func (w *Workspace) hasMember(userID user.UserID) bool {
	for _, member := range w.members {
		if member.UserID().Equal(userID) {
			return true
		}
	}

	return false
}

func (w *Workspace) memberIndex(userID user.UserID) int {
	for i, member := range w.members {
		if member.UserID().Equal(userID) {
			return i
		}
	}

	return -1
}

func (w *Workspace) ownerCount() int {
	count := 0

	for _, member := range w.members {
		if member.Role() == WorkspaceRoleOwner {
			count++
		}
	}

	return count
}
