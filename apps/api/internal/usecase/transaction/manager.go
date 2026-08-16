package transaction

import "context"

// Manager は、複数Repositoryへの操作を同じトランザクション内で実行する。
type Manager interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
}
