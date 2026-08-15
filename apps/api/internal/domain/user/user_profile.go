package user

type UserProfile struct {
	userID UserID
	name   string
}

func NewUserProfile(userID UserID, name string) UserProfile {
	return UserProfile{
		userID: userID,
		name:   name,
	}
}
