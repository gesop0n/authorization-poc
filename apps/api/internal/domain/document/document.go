package document

import (
	"fmt"
	"strings"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/project"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/user"
)

type Document struct {
	id              DocumentID
	projectID       project.ProjectID
	ownerUserID     user.UserID
	title           string
	content         string
	confidentiality Confidentiality
	status          DocumentStatus
}

func NewDocument(
	projectID project.ProjectID,
	ownerUserID user.UserID,
	title string,
	content string,
	confidentiality Confidentiality,
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

	if !confidentiality.IsValid() {
		return nil, ErrInvalidConfidentiality
	}

	id, err := newDocumentID()
	if err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}

	return &Document{
		id:              id,
		projectID:       projectID,
		ownerUserID:     ownerUserID,
		title:           title,
		content:         content,
		confidentiality: confidentiality,
		status:          DocumentStatusActive,
	}, nil
}

func Reconstruct(
	id DocumentID,
	projectID project.ProjectID,
	ownerUserID user.UserID,
	title string,
	content string,
	confidentiality Confidentiality,
	status DocumentStatus,
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
	if !confidentiality.IsValid() {
		return nil, ErrInvalidConfidentiality
	}
	if !status.IsValid() {
		return nil, ErrInvalidDocumentStatus
	}

	return &Document{
		id:              id,
		projectID:       projectID,
		ownerUserID:     ownerUserID,
		title:           title,
		content:         content,
		confidentiality: confidentiality,
		status:          status,
	}, nil
}

func (d *Document) ID() DocumentID {
	return d.id
}

func (d *Document) ProjectID() project.ProjectID { return d.projectID }

func (d *Document) OwnerUserID() user.UserID { return d.ownerUserID }

func (d *Document) Title() string { return d.title }

func (d *Document) Content() string { return d.content }

func (d *Document) Confidentiality() Confidentiality { return d.confidentiality }

func (d *Document) Status() DocumentStatus { return d.status }

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

func (d *Document) ChangeConfidentiality(confidentiality Confidentiality) error {
	if !confidentiality.IsValid() {
		return ErrInvalidConfidentiality
	}
	d.confidentiality = confidentiality
	return nil
}

func (d *Document) Archive() { d.status = DocumentStatusArchived }

func (d *Document) Restore() { d.status = DocumentStatusActive }
