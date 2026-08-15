package user

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
)

func (status UserStatus) IsValid() bool {
	switch status {
	case UserStatusActive, UserStatusSuspended:
		return true
	default:
		return false
	}
}
