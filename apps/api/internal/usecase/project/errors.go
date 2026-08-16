package projectusecase

import "errors"

var (
	// ワークスペースとプロジェクトが対応しない
	ErrProjectWorkspaceMismatch = errors.New("project does not belong to workspace")

	// ユーザーがワークスペースメンバーではない
	ErrUserNotWorkspaceMember = errors.New("user is not a workspace member")
)
