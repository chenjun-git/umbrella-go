package sql

import (
	"context"
	"database/sql"
)

type Stmt struct {
	*sql.Stmt
	query      string
	context    MiddlewareContext
	middleware DBMiddleware
}

type isPreparedKey struct{}

func IsPreparedStatement(mctx MiddlewareContext) bool {
	return mctx.Value(isPreparedKey{}) != nil
}

func WrapStmt(stmt *sql.Stmt, query string, mctx MiddlewareContext, wrappers ...DBMiddleware) *Stmt {
	mctx = context.WithValue(mctx, isPreparedKey{}, true)
	return &Stmt{
		Stmt:       stmt,
		query:      query,
		context:    mctx,
		middleware: ChainDBMiddlewares(wrappers...),
	}
}

func UnwrapStmt(stmt *Stmt) *sql.Stmt {
	return stmt.Stmt
}

func (s *Stmt) QueryStr() string {
	return s.query
}

func (s *Stmt) Close() error {
	return s.middleware.CloseStmt(s.context, fromCloser(s.Stmt))
}

func (s *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	return s.ExecContext(context.Background(), args...)
}

func fromStmtExecContext(ec func(ctx context.Context, args ...interface{}) (sql.Result, error)) ExecContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (sql.Result, error) {
		return ec(ctx, args...)
	}
}

func fromStmtQueryContext(qc func(context.Context, ...interface{}) (*sql.Rows, error)) QueryContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (*sql.Rows, MiddlewareContext, error) {
		rows, err := qc(ctx, args...)
		return rows, mctx, err
	}
}

func fromStmtQueryRowContext(qrc func(ctx context.Context, args ...interface{}) *sql.Row) QueryRowContextFunc {
	return func(mctx MiddlewareContext, ctx context.Context, query string, args []interface{}) (*sql.Row, MiddlewareContext) {
		return qrc(ctx, args...), mctx
	}
}

func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	return s.middleware.ExecContext(s.context, ctx, fromStmtExecContext(s.Stmt.ExecContext), s.query, args)
}

func (s *Stmt) Query(args ...interface{}) (*Rows, error) {
	return s.QueryContext(context.Background(), args...)
}

func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*Rows, error) {
	rows, mctx, err := s.middleware.QueryContext(s.context, ctx, fromStmtQueryContext(s.Stmt.QueryContext), s.query, args)
	if err != nil {
		return nil, err
	}
	return WrapRows(rows, mctx, s.middleware), nil
}

func (s *Stmt) QueryRow(args ...interface{}) *Row {
	return s.QueryRowContext(context.Background(), args...)
}

func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *Row {
	row, mctx := s.middleware.QueryRowContext(s.context, ctx, fromStmtQueryRowContext(s.Stmt.QueryRowContext), s.query, args)
	return WrapRow(row, mctx, s.middleware)
}
