package sql

import (
	"database/sql"
)

type Row struct {
	*sql.Row
	context    MiddlewareContext
	middleware DBMiddleware
}

func WrapRow(row *sql.Row, ctx MiddlewareContext, wrappers ...DBMiddleware) *Row {
	return &Row{
		Row:        row,
		context:    ctx,
		middleware: ChainDBMiddlewares(wrappers...),
	}
}

func UnwrapRow(row *Row) *sql.Row {
	return row.Row
}

func (r *Row) Scan(dest ...interface{}) error {
	return r.middleware.ScanRow(r.context, fromScan(r.Row.Scan), dest)
}

type Rows struct {
	*sql.Rows
	context    MiddlewareContext
	middleware DBMiddleware
}

func WrapRows(rows *sql.Rows, ctx MiddlewareContext, wrappers ...DBMiddleware) *Rows {
	return &Rows{
		Rows:       rows,
		context:    ctx,
		middleware: ChainDBMiddlewares(wrappers...),
	}
}

func UnwrapRows(rows *Rows) *sql.Rows {
	return rows.Rows
}

func (rs *Rows) Scan(dest ...interface{}) error {
	return rs.middleware.ScanRows(rs.context, fromScan(rs.Rows.Scan), dest)
}

func (rs *Rows) Close() error {
	return rs.middleware.CloseRows(rs.context, fromCloser(rs.Rows))
}
