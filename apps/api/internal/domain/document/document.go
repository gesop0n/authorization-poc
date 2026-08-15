package document

import (
	"fmt"
	"strings"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
)

type Document struct {
	id          DocumentID
	projectID   project.ProjectID
	ownerUserID user.UserID
	title       string
	content     string
}

func NewDocument(projectID project.ProjectID, ownerUserID user.UserID, title string, content string) (*Document, error) {
	if projectID.IsZero() {
		return nil, project.ErrInvalidProjectID
	}

	if ownerUserID.IsZero() {
		return nil, user.ErrInvalidUserID
	}

	if strings.TrimSpace(title) == "" {
		return nil, ErrDocumentTitleRequired
	}

	id, err := newDocumentID()
	if err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}

	return &Document{
		id:          id,
		projectID:   projectID,
		ownerUserID: ownerUserID,
		title:       title,
		content:     content,
	}, nil
}

func (d *Document) ID() DocumentID {
	return d.id
}
