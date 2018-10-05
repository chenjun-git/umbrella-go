package sql

import (
	"context"
	"database/sql"
)

type Conn struct {
	*sql.Conn
	context    MiddlewareContext
	middleware DBMiddleware
}

func WrapConn(conn *sql.Conn, mctx MiddlewareContext, wrappers ...DBMiddleware) *Conn {
	return &Conn{
		Conn:       conn,
		context:    mctx,
		middleware: ChainDBMiddlewares(wrappers...),
	}
}

func (c *Conn) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, mctx, err := c.middleware.BeginTx(c.context, ctx, opts, fromBeginTx(c.Conn.BeginTx))
	if err != nil {
		return nil, err
	}
	return WrapTx(tx, mctx, c.middleware), nil
}

func (c *Conn) Close() error {
	return c.middleware.CloseConn(c.context, fromCloser(c.Conn))
}

func (c *Conn) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.middleware.ExecContext(c.context, ctx, fromExecContext(c.Conn.ExecContext), query, args)
}

func (c *Conn) PingContext(ctx context.Context) error {
	return c.middleware.PingContext(c.context, ctx, fromPingContext(c.Conn.PingContext))
}

func (c *Conn) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, mctx, err := c.middleware.PrepareContext(c.context, ctx, query, fromPrepareContext(c.Conn.PrepareContext))
	if err != nil {
		return nil, err
	}
	return WrapStmt(stmt, query, mctx, c.middleware), nil
}

func (c *Conn) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	rows, mctx, err := c.middleware.QueryContext(c.context, ctx, fromQueryContext(c.Conn.QueryContext), query, args)
	if err != nil {
		return nil, err
	}
	return WrapRows(rows, mctx, c.middleware), nil
}

func (c *Conn) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	row, mctx := c.middleware.QueryRowContext(c.context, ctx, fromQueryRowContext(c.Conn.QueryRowContext), query, args)

	return WrapRow(row, mctx, c.middleware)
}
