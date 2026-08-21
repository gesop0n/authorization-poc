package user

type UserProfile struct {
	userID      UserID
	displayName string
}

func NewUserProfile(userID UserID, displayName string) UserProfile {
	return UserProfile{
		userID:      userID,
		displayName: displayName,
	}
}
