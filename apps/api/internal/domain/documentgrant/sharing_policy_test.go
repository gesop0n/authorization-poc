package documentgrant_test

import (
	"errors"
	"testing"

	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/document"
	"github.com/gesop0n/authorization-poc/apps/api/internal/domain/documentgrant"
)

func TestCanShareWith(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		confidentiality   document.Confidentiality
		isWorkspaceMember bool
		wantErr           error
	}{
		{"publicを内部Userに共有できる", document.ConfidentialityPublic, true, nil},
		{"publicを外部Userに共有できる", document.ConfidentialityPublic, false, nil},
		{"internalを内部Userに共有できる", document.ConfidentialityInternal, true, nil},
		{"internalの外部共有を拒否", document.ConfidentialityInternal, false, documentgrant.ErrExternalSharingProhibited},
		{"confidentialを内部Userに共有できる", document.ConfidentialityConfidential, true, nil},
		{"confidentialの外部共有を拒否", document.ConfidentialityConfidential, false, documentgrant.ErrExternalSharingProhibited},
		{"privateの内部共有を拒否", document.ConfidentialityPrivate, true, documentgrant.ErrDocumentNotShareable},
		{"privateの外部共有を拒否", document.ConfidentialityPrivate, false, documentgrant.ErrDocumentNotShareable},
		{"不正な機密区分を拒否", document.Confidentiality("invalid"), true, document.ErrInvalidConfidentiality},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := documentgrant.CanShareWith(tt.confidentiality, tt.isWorkspaceMember)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("CanShareWith() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
