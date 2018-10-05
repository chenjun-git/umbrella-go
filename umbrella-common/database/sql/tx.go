package sql

import (
	"context"
	"database/sql"
)

type Tx struct {
	*sql.Tx
	context    MiddlewareContext
	middleware DBMiddleware
}

func WrapTx(tx *sql.Tx, mctx MiddlewareContext, wrappers ...DBMiddleware) *Tx {
	return &Tx{
		Tx:         tx,
		context:    mctx,
		middleware: ChainDBMiddlewares(wrappers...),
	}
}

func UnwrapTx(tx *Tx) *sql.Tx {
	return tx.Tx
}

func (tx *Tx) Commit() error {
	return tx.middleware.Commit(tx.context, fromCommit(tx.Tx.Commit))
}

func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.ExecContext(context.Background(), query, args...)
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return tx.middleware.ExecContext(tx.context, ctx, fromExecContext(tx.Tx.ExecContext), query, args)
}

func (tx *Tx) Prepare(query string) (*Stmt, error) {
	return tx.PrepareContext(context.Background(), query)
}

func (tx *Tx) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, mctx, err := tx.middleware.PrepareContext(tx.context, ctx, query, fromPrepareContext(tx.Tx.PrepareContext))
	if err != nil {
		return nil, err
	}

	return WrapStmt(stmt, query, mctx, tx.middleware), nil
}

func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error) {
	return tx.QueryContext(context.Background(), query, args...)
}

func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	rows, mctx, err := tx.middleware.QueryContext(tx.context, ctx, fromQueryContext(tx.Tx.QueryContext), query, args)
	if err != nil {
		return nil, err
	}
	return WrapRows(rows, mctx, tx.middleware), nil
}

func (tx *Tx) QueryRow(query string, args ...interface{}) *Row {
	return tx.QueryRowContext(context.Background(), query, args...)
}

func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	row, mctx := tx.middleware.QueryRowContext(tx.context, ctx, fromQueryRowContext(tx.Tx.QueryRowContext), query, args)
	return WrapRow(row, mctx, tx.middleware)
}

func (tx *Tx) Rollback() error {
	return tx.middleware.Rollback(tx.context, fromRollback(tx.Tx.Rollback))
}

func (tx *Tx) Stmt(stmt *Stmt) *Stmt {
	return tx.StmtContext(context.Background(), stmt)
}

func (tx *Tx) StmtContext(ctx context.Context, stmt *Stmt) *Stmt {
	s, mctx := tx.middleware.StmtContext(tx.context, ctx, UnwrapStmt(stmt), fromStmtContext(tx.Tx.StmtContext))
	return WrapStmt(s, stmt.query, mctx, tx.middleware)
}
