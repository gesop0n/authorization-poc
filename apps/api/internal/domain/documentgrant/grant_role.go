package documentgrant

// GrantRole は、Document単位で付与する権限。
type GrantRole string

const (
	GrantRoleViewer GrantRole = "viewer"
	GrantRoleEditor GrantRole = "editor"
)

func (role GrantRole) IsValid() bool {
	switch role {
	case GrantRoleViewer, GrantRoleEditor:
		return true
	default:
		return false
	}
}
