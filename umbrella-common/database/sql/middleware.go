package sql

import (
	"context"
	"database/sql"
	"io"
)

type MiddlewareContext context.Context

type DBMiddleware interface {
	Conn(ctx context.Context, next ConnFunc) (*sql.Conn, MiddlewareContext, error)

	ExecContext(mctx MiddlewareContext, ctx context.Context, next ExecContextFunc, query string, args []interface{}) (sql.Result, error)
	QueryContext(mctx MiddlewareContext, ctx context.Context, next QueryContextFunc, query string, args []interface{}) (*sql.Rows, MiddlewareContext, error)
	QueryRowContext(mctx MiddlewareContext, ctx context.Context, next QueryRowContextFunc, query string, args []interface{}) (*sql.Row, MiddlewareContext)

	BeginTx(mctx MiddlewareContext, ctx context.Context, opts *sql.TxOptions, next BeginTxFunc) (*sql.Tx, MiddlewareContext, error)
	PingContext(mctx MiddlewareContext, ctx context.Context, next PingContextFunc) error

	Commit(mctx MiddlewareContext, next CommitFunc) error
	Rollback(mctx MiddlewareContext, next RollbackFunc) error
	PrepareContext(mctx MiddlewareContext, ctx context.Context, query string, next PrepareContextFunc) (*sql.Stmt, MiddlewareContext, error)
	StmtContext(mctx MiddlewareContext, ctx context.Context, stmt *sql.Stmt, next StmtContextFunc) (*sql.Stmt, MiddlewareContext)

	ScanRow(mctx MiddlewareContext, next ScanFunc, dest []interface{}) error
	ScanRows(mctx MiddlewareContext, next ScanFunc, dest []interface{}) error

	CloseStmt(mctx MiddlewareContext, next CloseFunc) error
	CloseConn(mctx MiddlewareContext, next CloseFunc) error
	CloseRows(mctx MiddlewareContext, next CloseFunc) error
	CloseDB(next io.Closer) error
}

func ChainDBMiddlewares(middlewares ...DBMiddleware) DBMiddleware {
	ms := removeNilDBMiddleware(middlewares)
	switch len(ms) {
	case 0:
		return defaultDBMiddleware
	case 1:
		return ms[0]
	default:
		return MultiDBMiddleware(ms)
	}
}

func removeNilDBMiddleware(middlewares []DBMiddleware) []DBMiddleware {
	var result []DBMiddleware
	for _, m := range middlewares {
		if m != nil {
			result = append(result, m)
		}
	}

	return result
}
