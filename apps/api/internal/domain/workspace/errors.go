package workspace

import "errors"

var (
	// ワークスペースIDが不正
	ErrInvalidWorkspaceID = errors.New("invalid workspace ID")

	// ワークスペースロールが不正な
	ErrInvalidWorkspaceRole = errors.New("invalid workspace role")

	// 名前が未入力
	ErrWorkspaceNameRequired = errors.New("workspace name is required")

	// メンバーのユーザーIDが不正
	ErrInvalidWorkspaceMemberUserID = errors.New("invalid workspace member user ID")

	// メンバーが重複
	ErrWorkspaceMemberAlreadyExists = errors.New("workspace member already exists")

	// メンバーが存在しない
	ErrWorkspaceMemberNotFound = errors.New("workspace member not found")

	// Ownerが不在
	ErrWorkspaceMustHaveOwner = errors.New("workspace must have at least one owner")
)
