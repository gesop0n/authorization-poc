package document

type Confidentiality string

const (
	ConfidentialityPublic       Confidentiality = "public"
	ConfidentialityInternal     Confidentiality = "internal"
	ConfidentialityConfidential Confidentiality = "confidential"
)

func (c Confidentiality) IsValid() bool {
	switch c {
	case ConfidentialityPublic, ConfidentialityInternal, ConfidentialityConfidential:
		return true
	default:
		return false
	}
}
