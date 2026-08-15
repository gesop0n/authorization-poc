package user

import "errors"

var (
	// 不正なユーザーID
	ErrInvalidUserID = errors.New("Invalid User ID")

	// ユーザー名は必須
	ErrUsernameRequired = errors.New("username is required")

	// ユーザーステータスが不正
	ErrInvalidUserStatus = errors.New("invalid user status")
)
