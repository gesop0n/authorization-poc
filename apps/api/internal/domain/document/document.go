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

func NewDocument(
	projectID project.ProjectID,
	ownerUserID user.UserID,
	title string,
	content string,
) (*Document, error) {
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

func Reconstruct(
	id DocumentID,
	projectID project.ProjectID,
	ownerUserID user.UserID,
	title string,
	content string,
) (*Document, error) {
	if id.IsZero() {
		return nil, ErrInvalidDocumentID
	}
	if projectID.IsZero() {
		return nil, project.ErrInvalidProjectID
	}
	if ownerUserID.IsZero() {
		return nil, user.ErrInvalidUserID
	}
	if strings.TrimSpace(title) == "" {
		return nil, ErrDocumentTitleRequired
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

func (d *Document) ProjectID() project.ProjectID { return d.projectID }

func (d *Document) OwnerUserID() user.UserID { return d.ownerUserID }

func (d *Document) Title() string { return d.title }

func (d *Document) Content() string { return d.content }

func (d *Document) ChangeTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrDocumentTitleRequired
	}
	d.title = title
	return nil
}

func (d *Document) ChangeContent(content string) { d.content = content }

func (d *Document) ChangeOwner(ownerUserID user.UserID) error {
	if ownerUserID.IsZero() {
		return user.ErrInvalidUserID
	}
	d.ownerUserID = ownerUserID
	return nil
}
