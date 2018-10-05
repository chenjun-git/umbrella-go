package sql

import (
	"context"
	"database/sql"
)

type (
	RawDB   = sql.DB
	RawConn = sql.Conn
	RawStmt = sql.Stmt
	RawTx   = sql.Tx
	RawRow  = sql.Row
	RawRows = sql.Rows

	Result         = sql.Result
	TxOptions      = sql.TxOptions
	ColumnType     = sql.ColumnType
	DBStats        = sql.DBStats
	IsolationLevel = sql.IsolationLevel
	NamedArg       = sql.NamedArg
	NullBool       = sql.NullBool
	NullFloat64    = sql.NullFloat64
	NullInt64      = sql.NullInt64
	NullString     = sql.NullString
	Out            = sql.Out
	RawBytes       = sql.RawBytes
	Scanner        = sql.Scanner
)

var (
	Named     = sql.Named
	Drivers   = sql.Drivers
	Register  = sql.Register
	ErrNoRows = sql.ErrNoRows
)

const (
	LevelDefault         = sql.LevelDefault
	LevelReadUncommitted = sql.LevelReadUncommitted
	LevelReadCommitted   = sql.LevelReadCommitted
	LevelWriteCommitted  = sql.LevelWriteCommitted
	LevelRepeatableRead  = sql.LevelRepeatableRead
	LevelSnapshot        = sql.LevelSnapshot
	LevelSerializable    = sql.LevelSerializable
	LevelLinearizable    = sql.LevelLinearizable
)

type DB struct {
	*sql.DB
	middleware DBMiddleware
}

func Open(driverName, dataSourceName string, middlewares ...DBMiddleware) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return WrapDB(db, middlewares...), nil
}

func WrapDB(db *sql.DB, wrappers ...DBMiddleware) *DB {
	return &DB{
		DB:         db,
		middleware: ChainDBMiddlewares(wrappers...),
	}
}

func UnwrapDB(db *DB) *sql.DB {
	return db.DB
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, mctx, err := db.middleware.BeginTx(context.Background(), ctx, opts, fromBeginTx(db.DB.BeginTx))
	if err != nil {
		return nil, err
	}

	return WrapTx(tx, mctx, db.middleware), nil
}

func (db *DB) Begin() (*Tx, error) {
	return db.BeginTx(context.Background(), nil)
}

func (db *DB) Close() error {
	return db.middleware.CloseDB(db.DB)
}

func (db *DB) Conn(ctx context.Context) (*Conn, error) {
	conn, mctx, err := db.middleware.Conn(ctx, fromConn(db.DB.Conn))
	if err != nil {
		return nil, err
	}

	return WrapConn(conn, mctx, db.middleware), nil
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.middleware.ExecContext(context.Background(), ctx, fromExecContext(db.DB.ExecContext), query, args)
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.ExecContext(context.Background(), query, args...)
}

func (db *DB) PingContext(ctx context.Context) error {
	return db.middleware.PingContext(context.Background(), ctx, fromPingContext(db.DB.PingContext))
}

func (db *DB) Ping() error {
	return db.PingContext(context.Background())
}

func (db *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, mctx, err := db.middleware.PrepareContext(context.Background(), ctx, query, fromPrepareContext(db.DB.PrepareContext))
	if err != nil {
		return nil, err
	}
	return WrapStmt(stmt, query, mctx, db.middleware), nil
}

func (db *DB) Prepare(query string) (*Stmt, error) {
	return db.PrepareContext(context.Background(), query)
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	rows, mctx, err := db.middleware.QueryContext(context.Background(), ctx, fromQueryContext(db.DB.QueryContext), query, args)
	if err != nil {
		return nil, err
	}
	return WrapRows(rows, mctx, db.middleware), nil
}

func (db *DB) Query(query string, args ...interface{}) (*Rows, error) {
	return db.QueryContext(context.Background(), query, args...)
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	row, mctx := db.middleware.QueryRowContext(context.Background(), ctx, fromQueryRowContext(db.DB.QueryRowContext), query, args)
	return WrapRow(row, mctx, db.middleware)
}

func (db *DB) QueryRow(query string, args ...interface{}) *Row {
	return db.QueryRowContext(context.Background(), query, args...)
}
