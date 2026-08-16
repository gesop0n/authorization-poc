package project

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
)

type Project struct {
	id          ProjectID
	name        string
	workspaceID workspace.WorkspaceID
	status      ProjectStatus
	members     []ProjectMember
}

func NewProject(name string, workspaceID workspace.WorkspaceID, adminUserID user.UserID) (*Project, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrProjectNameRequired
	}

	if workspaceID.IsZero() {
		return nil, workspace.ErrInvalidWorkspaceID
	}
	if adminUserID.IsZero() {
		return nil, ErrInvalidProjectMemberUserID
	}

	id, err := newProjectID()
	if err != nil {
		return nil, fmt.Errorf("create Project: %w", err)
	}

	return &Project{
		id:          id,
		name:        name,
		workspaceID: workspaceID,
		status:      ProjectStatusActive,
		members:     []ProjectMember{newProjectMember(adminUserID, ProjectRoleAdmin)},
	}, nil
}

func Reconstruct(id ProjectID, name string, workspaceID workspace.WorkspaceID, status ProjectStatus, members []ProjectMember) (*Project, error) {
	if id.IsZero() {
		return nil, ErrInvalidProjectID
	}

	if strings.TrimSpace(name) == "" {
		return nil, ErrProjectNameRequired
	}

	if workspaceID.IsZero() {
		return nil, workspace.ErrInvalidWorkspaceID
	}

	if !status.IsValid() {
		return nil, ErrInvalidProjectStatus
	}

	seen := make(map[string]struct{}, len(members))
	adminCount := 0
	for _, member := range members {
		if member.UserID().IsZero() {
			return nil, ErrInvalidProjectMemberUserID
		}
		if !member.Role().IsValid() {
			return nil, ErrInvalidProjectRole
		}
		key := member.UserID().String()
		if _, exists := seen[key]; exists {
			return nil, ErrProjectMemberAlreadyExists
		}
		seen[key] = struct{}{}
		if member.Role() == ProjectRoleAdmin {
			adminCount++
		}
	}
	if adminCount == 0 {
		return nil, ErrProjectMustHaveAdmin
	}

	return &Project{id: id, name: name, workspaceID: workspaceID, status: status, members: slices.Clone(members)}, nil
}

func (p *Project) ID() ProjectID {
	return p.id
}

func (p *Project) Name() string {
	return p.name
}

func (p *Project) WorkspaceID() workspace.WorkspaceID {
	return p.workspaceID
}

func (p *Project) Status() ProjectStatus {
	return p.status
}

func (p *Project) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrProjectNameRequired
	}

	p.name = name
	return nil
}

func (p *Project) Archive() {
	p.status = ProjectStatusArchived
}

func (p *Project) Restore() {
	p.status = ProjectStatusActive
}

// AddMember はProject内の不変条件を検証してメンバーを追加する。
// UserがこのProjectのWorkspaceに所属しているかは別集約の情報なので、呼び出し前に
// アプリケーション層で Workspace.HasMember を使って検証する。
func (p *Project) AddMember(userID user.UserID, role ProjectRole) error {
	if userID.IsZero() {
		return ErrInvalidProjectMemberUserID
	}
	if !role.IsValid() {
		return ErrInvalidProjectRole
	}
	if p.memberIndex(userID) >= 0 {
		return ErrProjectMemberAlreadyExists
	}
	p.members = append(p.members, newProjectMember(userID, role))
	return nil
}

func (p *Project) RemoveMember(userID user.UserID) error {
	if err := p.CanRemoveMember(userID); err != nil {
		return err
	}

	index := p.memberIndex(userID)
	p.members = append(p.members[:index:index], p.members[index+1:]...)
	return nil
}

// CanRemoveMember は、状態を変更せずにメンバーを削除できるか検証する。
func (p *Project) CanRemoveMember(userID user.UserID) error {
	index := p.memberIndex(userID)
	if index < 0 {
		return ErrProjectMemberNotFound
	}
	if p.members[index].Role() == ProjectRoleAdmin && p.adminCount() == 1 {
		return ErrProjectMustHaveAdmin
	}
	return nil
}

func (p *Project) ChangeMemberRole(userID user.UserID, role ProjectRole) error {
	if !role.IsValid() {
		return ErrInvalidProjectRole
	}
	index := p.memberIndex(userID)
	if index < 0 {
		return ErrProjectMemberNotFound
	}
	current := p.members[index]
	if current.Role() == role {
		return nil
	}
	if current.Role() == ProjectRoleAdmin && role != ProjectRoleAdmin && p.adminCount() == 1 {
		return ErrProjectMustHaveAdmin
	}
	p.members[index] = newProjectMember(userID, role)
	return nil
}

func (p *Project) HasMember(userID user.UserID) bool { return p.memberIndex(userID) >= 0 }

func (p *Project) Members() []ProjectMember { return slices.Clone(p.members) }

func (p *Project) memberIndex(userID user.UserID) int {
	for i, member := range p.members {
		if member.UserID().Equal(userID) {
			return i
		}
	}
	return -1
}

func (p *Project) adminCount() int {
	count := 0
	for _, member := range p.members {
		if member.Role() == ProjectRoleAdmin {
			count++
		}
	}
	return count
}
