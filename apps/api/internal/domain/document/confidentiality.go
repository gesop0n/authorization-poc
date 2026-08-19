package document

type Confidentiality string

const (
	ConfidentialityPublic       Confidentiality = "public"
	ConfidentialityInternal     Confidentiality = "internal"
	ConfidentialityConfidential Confidentiality = "confidential"
	ConfidentialityPrivate      Confidentiality = "private"
)

func (c Confidentiality) IsValid() bool {
	switch c {
	case ConfidentialityPublic, ConfidentialityInternal, ConfidentialityConfidential, ConfidentialityPrivate:
		return true
	default:
		return false
	}
}
