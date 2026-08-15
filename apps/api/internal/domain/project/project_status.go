package project

type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusArchived ProjectStatus = "archived"
)

func (status ProjectStatus) IsValid() bool {
	switch status {
	case ProjectStatusActive, ProjectStatusArchived:
		return true
	default:
		return false
	}
}
