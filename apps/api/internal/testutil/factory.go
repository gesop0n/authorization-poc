package testutil

import (
	"testing"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/workspace"
)

// 指定した名前の有効なUserを生成する。
func NewUser(t testing.TB, name string) *user.User {
	t.Helper()
	u, err := user.NewUser(name)
	if err != nil {
		t.Fatalf("user.NewUser() error = %v", err)
	}
	return u
}

// 有効なOwnerを持つWorkspaceを生成する。
func NewWorkspaceWithOwner(t testing.TB, workspaceName, ownerName string) (*workspace.Workspace, *user.User) {
	t.Helper()
	owner := NewUser(t, ownerName)
	ws, err := workspace.NewWorkspace(workspaceName, owner.ID())
	if err != nil {
		t.Fatalf("workspace.NewWorkspace() error = %v", err)
	}
	return ws, owner
}

// 有効なWorkspaceとAdminを持つProjectを生成する。
func NewProjectWithWorkspace(t testing.TB, projectName, workspaceName, ownerName string) (*project.Project, *workspace.Workspace, *user.User) {
	t.Helper()
	ws, owner := NewWorkspaceWithOwner(t, workspaceName, ownerName)
	p, err := project.NewProject(projectName, ws.ID(), owner.ID())
	if err != nil {
		t.Fatalf("project.NewProject() error = %v", err)
	}
	return p, ws, owner
}

// 有効なProjectとOwnerを持つDocumentを生成する。
func NewDocumentWithProject(t testing.TB) (*document.Document, *project.Project, *user.User) {
	t.Helper()
	p, _, owner := NewProjectWithWorkspace(t, "project", "workspace", "owner")
	d, err := document.NewDocument(p.ID(), owner.ID(), "title", "content", document.ConfidentialityInternal)
	if err != nil {
		t.Fatalf("document.NewDocument() error = %v", err)
	}
	return d, p, owner
}
