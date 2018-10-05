package sql

import "io"

var (
	defaultDBMiddleware = &DefaultDBMiddleware{}
)

type DefaultDBMiddleware struct{}

func (dm *DefaultDBMiddleware) ExecContext(mctx MiddlewareContext, ctx context.Context, next ExecContextFunc, query string, args []interface{}) (sql.Result, error) {
	return next(mctx, ctx, query, args)
}

func (dm *DefaultDBMiddleware) QueryContext(mctx MiddlewareContext, ctx context.Context, next QueryContextFunc, query string, args []interface{}) (*sql.Rows, MiddlewareContext, error) {
	return next(mctx, ctx, query, args)
}

func (dm *DefaultDBMiddleware) QueryRowContext(mctx MiddlewareContext, ctx context.Context, next QueryRowContextFunc, query string, args []interface{}) (*sql.Row, MiddlewareContext) {
	return next(mctx, ctx, query, args)
}

func (dm *DefaultDBMiddleware) Commit(mctx MiddlewareContext, next CommitFunc) error {
	return next(mctx)
}

func (dm *DefaultDBMiddleware) Rollback(mctx MiddlewareContext, next RollbackFunc) error {
	return next(mctx)
}

func (dm *DefaultDBMiddleware) PrepareContext(mctx MiddlewareContext, ctx context.Context, query string, next PrepareContextFunc) (*sql.Stmt, MiddlewareContext, error) {
	return next(mctx, ctx, query)
}

func (dm *DefaultDBMiddleware) StmtContext(mctx MiddlewareContext, ctx context.Context, stmt *sql.Stmt, next StmtContextFunc) (*sql.Stmt, MiddlewareContext) {
	return next(mctx, ctx, stmt)
}

func (dm *DefaultDBMiddleware) CloseStmt(mctx MiddlewareContext, next CloseFunc) error {
	return next(mctx)
}

func (dm *DefaultDBMiddleware) BeginTx(mctx MiddlewareContext, ctx context.Context, opts *sql.TxOptions, next BeginTxFunc) (*sql.Tx, MiddlewareContext, error) {
	return next(mctx, ctx, opts)
}

func (dm *DefaultDBMiddleware) PingContext(mctx MiddlewareContext, ctx context.Context, next PingContextFunc) error {
	return next(mctx, ctx)
}

func (dm *DefaultDBMiddleware) CloseConn(mctx MiddlewareContext, next CloseFunc) error {
	return next(mctx)
}

func (dm *DefaultDBMiddleware) Conn(ctx context.Context, next ConnFunc) (*sql.Conn, MiddlewareContext, error) {
	return next(ctx)
}

func (dm *DefaultDBMiddleware) CloseDB(next io.Closer) error {
	return next.Close()
}

func (dm *DefaultDBMiddleware) ScanRow(mctx MiddlewareContext, next ScanFunc, dest []interface{}) error {
	return next(mctx, dest)
}

func (dm *DefaultDBMiddleware) ScanRows(mctx MiddlewareContext, next ScanFunc, dest []interface{}) error {
	return next(mctx, dest)
}

func (dm *DefaultDBMiddleware) CloseRows(mctx MiddlewareContext, next CloseFunc) error {
	return next(mctx)
}
