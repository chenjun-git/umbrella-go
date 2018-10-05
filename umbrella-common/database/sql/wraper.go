package sql

import (
	"context"
	"database/sql"
	"io"
)

type CloseFunc func(MiddlewareContext) error

func (f CloseFunc) Close() error {
	return f(context.Background())
}

func fromCloser(closer io.Closer) CloseFunc {
	return func(mctx MiddlewareContext) error {
		return closer.Close()
	}
}

func (f CloseFunc) WrappedConn(w DBMiddleware) CloseFunc {
	return CloseFunc(func(mctx MiddlewareContext) error {
		return w.CloseConn(mctx, f)
	})
}

func (f CloseFunc) WrappedStmt(w DBMiddleware) CloseFunc {
	return CloseFunc(func(mctx MiddlewareContext) error {
		return w.CloseStmt(mctx, f)
	})
}

func (f CloseFunc) WrappedRows(w DBMiddleware) CloseFunc {
	return CloseFunc(func(mctx MiddlewareContext) error {
		return w.CloseRows(mctx, f)
	})
}

type ExecContextFunc func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (sql.Result, error)

func fromExecContext(ec func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)) ExecContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (sql.Result, error) {
		return ec(ctx, query, args...)
	}
}

func (f ExecContextFunc) Wrapped(w DBMiddleware) ExecContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (sql.Result, error) {
		return w.ExecContext(mctx, ctx, f, query, args)
	}
}

type QueryContextFunc func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (*sql.Rows, MiddlewareContext, error)

func fromQueryContext(qc func(context.Context, string, ...interface{}) (*sql.Rows, error)) QueryContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (*sql.Rows, MiddlewareContext, error) {
		rows, err := qc(ctx, query, args...)
		return rows, mctx, err
	}
}

func (f QueryContextFunc) Wrapped(w DBMiddleware) QueryContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (*sql.Rows, MiddlewareContext, error) {
		return w.QueryContext(mctx, ctx, f, query, args)
	}
}

type QueryRowContextFunc func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (*sql.Row, MiddlewareContext)

func fromQueryRowContext(qrc func(ctx context.Context, query string, args ...interface{}) *sql.Row) QueryRowContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (*sql.Row, MiddlewareContext) {
		return qrc(ctx, query, args...), mctx
	}
}

func (f QueryRowContextFunc) Wrapped(w DBMiddleware) QueryRowContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (*sql.Row, MiddlewareContext) {
		return w.QueryRowContext(mctx, ctx, f, query, args)
	}
}

type BeginTxFunc func(mctx MiddlewareContext, ctx context.Context, opts *sql.TxOptions) (*sql.Tx, MiddlewareContext, error)

func fromBeginTx(bt func(context.Context, *sql.TxOptions) (*sql.Tx, error)) BeginTxFunc {
	return func(mctx MiddlewareContext, ctx context.Context, opts *sql.TxOptions) (*sql.Tx, MiddlewareContext, error) {
		tx, err := bt(ctx, opts)
		return tx, mctx, err
	}
}

func (f BeginTxFunc) Wrapped(w DBMiddleware) BeginTxFunc {
	return func(mctx MiddlewareContext, ctx context.Context, opts *sql.TxOptions) (*sql.Tx, MiddlewareContext, error) {
		return w.BeginTx(mctx, ctx, opts, f)
	}
}

type PingContextFunc func(mctx MiddlewareContext, ctx context.Context) error

func fromPingContext(pc func(ctx context.Context) error) PingContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context) error {
		return pc(ctx)
	}
}

func (f PingContextFunc) Wrapped(w DBMiddleware) PingContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context) error {
		return w.PingContext(mctx, ctx, f)
	}
}

type PrepareContextFunc func(mctx MiddlewareContext, ctx context.Context, query string) (*sql.Stmt, MiddlewareContext, error)

func fromPrepareContext(pc func(ctx context.Context, query string) (*sql.Stmt, error)) PrepareContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, query string) (*sql.Stmt, MiddlewareContext, error) {
		s, err := pc(ctx, query)
		return s, mctx, err
	}
}

func (f PrepareContextFunc) Wrapped(w DBMiddleware) PrepareContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, query string) (*sql.Stmt, MiddlewareContext, error) {
		return w.PrepareContext(mctx, ctx, query, f)
	}
}

type CommitFunc func(MiddlewareContext) error

func fromCommit(c func() error) CommitFunc {
	return func(mctx MiddlewareContext) error {
		return c()
	}
}

func (f CommitFunc) Wrapped(w DBMiddleware) CommitFunc {
	return CommitFunc(func(mctx MiddlewareContext) error {
		return w.Commit(mctx, f)
	})
}

type RollbackFunc func(MiddlewareContext) error

func fromRollback(r func() error) RollbackFunc {
	return func(mctx MiddlewareContext) error {
		return r()
	}
}

func (f RollbackFunc) Wrapped(w DBMiddleware) RollbackFunc {
	return RollbackFunc(func(mctx MiddlewareContext) error {
		return w.Rollback(mctx, f)
	})
}

type StmtContextFunc func(mctx MiddlewareContext, ctx context.Context, stmt *sql.Stmt) (*sql.Stmt, MiddlewareContext)

func fromStmtContext(sc func(ctx context.Context, stmt *sql.Stmt) *sql.Stmt) StmtContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, stmt *sql.Stmt) (*sql.Stmt, MiddlewareContext) {
		return sc(ctx, stmt), mctx
	}
}

func (f StmtContextFunc) Wrapped(w DBMiddleware) StmtContextFunc {
	return StmtContextFunc(func(mctx MiddlewareContext, ctx context.Context, stmt *sql.Stmt) (*sql.Stmt, MiddlewareContext) {
		return w.StmtContext(mctx, ctx, stmt, f)
	})
}

type ScanFunc func(mctx MiddlewareContext, dest []interface{}) error

func fromScan(s func(dest ...interface{}) error) ScanFunc {
	return func(mctx MiddlewareContext, dest []interface{}) error {
		return s(dest...)
	}
}

func (f ScanFunc) WrappedRow(w DBMiddleware) ScanFunc {
	return ScanFunc(func(mctx MiddlewareContext, dest []interface{}) error {
		return w.ScanRow(mctx, f, dest)
	})
}

func (f ScanFunc) WrappedRows(w DBMiddleware) ScanFunc {
	return ScanFunc(func(mctx MiddlewareContext, dest []interface{}) error {
		return w.ScanRows(mctx, f, dest)
	})
}

type ConnFunc func(ctx context.Context) (*sql.Conn, MiddlewareContext, error)

func fromConn(conn func(context.Context) (*sql.Conn, error)) ConnFunc {
	return func(ctx context.Context) (*sql.Conn, MiddlewareContext, error) {
		c, err := conn(ctx)
		return c, context.Background(), err
	}
}

func (f ConnFunc) Wrapped(w DBMiddleware) ConnFunc {
	return ConnFunc(func(ctx context.Context) (*sql.Conn, MiddlewareContext, error) {
		return w.Conn(ctx, f)
	})
}
