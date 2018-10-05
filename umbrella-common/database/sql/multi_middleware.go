package sql

import (
	"context"
	"database/sql"
	"io"
)

type MultiDBMiddleware []DBMiddleware

func (mdm MultiDBMiddleware) Conn(ctx context.Context, next ConnFunc) (*sql.Conn, MiddlewareContext, error) {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.Wrapped(mdm[i])
	}

	conn, mctx, err := n(ctx)
	return conn, mctx, err
}

func (mdm MultiDBMiddleware) ExecContext(mctx MiddlewareContext, ctx context.Context, next ExecContextFunc, query string, args []interface{}) (sql.Result, error) {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.Wrapped(mdm[i])
	}
	return n(mctx, ctx, query, args)
}

func (mdm MultiDBMiddleware) QueryContext(mctx MiddlewareContext, ctx context.Context, next QueryContextFunc, query string, args []interface{}) (*sql.Rows, MiddlewareContext, error) {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.Wrapped(mdm[i])
	}
	rows, mctx, err := n(mctx, ctx, query, args)
	return rows, mctx, err
}

func (mdm MultiDBMiddleware) QueryRowContext(mctx MiddlewareContext, ctx context.Context, next QueryRowContextFunc, query string, args []interface{}) (*sql.Row, MiddlewareContext) {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.Wrapped(mdm[i])
	}
	return n(mctx, ctx, query, args)
}

func (mdm MultiDBMiddleware) CloseStmt(mctx MiddlewareContext, next CloseFunc) error {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.WrappedStmt(mdm[i])
	}
	return n(mctx)
}

func (mdm MultiDBMiddleware) BeginTx(mctx MiddlewareContext, ctx context.Context, opts *sql.TxOptions, next BeginTxFunc) (*sql.Tx, MiddlewareContext, error) {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.Wrapped(mdm[i])
	}

	return n(mctx, ctx, opts)
}

func (mdm MultiDBMiddleware) PingContext(mctx MiddlewareContext, ctx context.Context, next PingContextFunc) error {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.Wrapped(mdm[i])
	}

	return n(mctx, ctx)
}

func (mdm MultiDBMiddleware) PrepareContext(mctx MiddlewareContext, ctx context.Context, query string, next PrepareContextFunc) (*sql.Stmt, MiddlewareContext, error) {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.Wrapped(mdm[i])
	}

	return n(mctx, ctx, query)
}

func (mdm MultiDBMiddleware) CloseConn(mctx MiddlewareContext, next CloseFunc) error {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.WrappedConn(mdm[i])
	}
	return n(mctx)
}

type closerWrapper struct {
	io.Closer
	middleware DBMiddleware
}

func (cw *closerWrapper) Close() error {
	return cw.middleware.CloseDB(cw.Closer)
}

func wrapCloser(n io.Closer, m DBMiddleware) io.Closer {
	return &closerWrapper{
		Closer:     n,
		middleware: m,
	}
}

func (mdm MultiDBMiddleware) CloseDB(next io.Closer) error {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = wrapCloser(n, mdm[i])
	}

	return n.Close()
}

func (mdm MultiDBMiddleware) ScanRow(mctx MiddlewareContext, next ScanFunc, dest []interface{}) error {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.WrappedRow(mdm[i])
	}
	return n(mctx, dest)
}

func (mdm MultiDBMiddleware) ScanRows(mctx MiddlewareContext, next ScanFunc, dest []interface{}) error {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.WrappedRows(mdm[i])
	}
	return n(mctx, dest)
}

func (mdm MultiDBMiddleware) CloseRows(mctx MiddlewareContext, next CloseFunc) error {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.WrappedRows(mdm[i])
	}
	return n(mctx)
}

func (mdm MultiDBMiddleware) Commit(mctx MiddlewareContext, next CommitFunc) error {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.Wrapped(mdm[i])
	}

	return n(mctx)
}

func (mdm MultiDBMiddleware) Rollback(mctx MiddlewareContext, next RollbackFunc) error {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.Wrapped(mdm[i])
	}

	return n(mctx)
}

func (mdm MultiDBMiddleware) StmtContext(mctx MiddlewareContext, ctx context.Context, stmt *sql.Stmt, next StmtContextFunc) (*sql.Stmt, MiddlewareContext) {
	n := next
	for i := len(mdm) - 1; i >= 0; i-- {
		n = n.Wrapped(mdm[i])
	}

	return n(mctx, ctx, stmt)
}
