package document

type DocumentStatus string

const (
	DocumentStatusActive   DocumentStatus = "active"
	DocumentStatusArchived DocumentStatus = "archived"
)

func (status DocumentStatus) IsValid() bool {
	switch status {
	case DocumentStatusActive, DocumentStatusArchived:
		return true
	default:
		return false
	}
}
