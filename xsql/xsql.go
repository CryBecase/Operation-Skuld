package xsql

import (
	"context"
	"database/sql"
	"errors"

	"skuld/xerr"
)

type Rows struct {
	*sql.Rows
}

func (r *Rows) Scan(desc ...interface{}) error {
	if err := r.Rows.Scan(desc...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return xerr.NoData
		}
		return err
	}
	return nil
}

type Row struct {
	*sql.Row
}

func (r *Row) Scan(desc ...interface{}) error {
	if err := r.Row.Scan(desc...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return xerr.NoData
		}
		return err
	}
	return nil
}

type Core interface {
	Exec(ctx context.Context, query string, args ...interface{}) (rst sql.Result, err error)
	Query(ctx context.Context, query string, args ...interface{}) (rows *Rows, err error)
	QueryRow(ctx context.Context, query string, args ...interface{}) (row *Row, err error)
}

type DBItf interface {
	Core
	Tx(ctx context.Context, f func(context.Context) error) error
}

type TxItf interface {
	Core
	Rollback() error
	Commit() error
	Tx(f func() error) error
}

type DB struct {
	db *sql.DB
}

type txkey struct{}

func New(db *sql.DB) *DB {
	return &DB{db: db}
}

func (d *DB) tx(ctx context.Context) TxItf {
	if v, ok := ctx.Value(txkey{}).(TxItf); ok {
		return v
	}
	return nil
}

func (d *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if tx := d.tx(ctx); tx != nil {
		return tx.Exec(ctx, query, args...)
	}

	return d.db.ExecContext(ctx, query, args...)
}

func (d *DB) Query(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	if tx := d.tx(ctx); tx != nil {
		return tx.Query(ctx, query, args...)
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{rows}, nil
}

func (d *DB) QueryRow(ctx context.Context, query string, args ...interface{}) (*Row, error) {
	if tx := d.tx(ctx); tx != nil {
		return tx.QueryRow(ctx, query, args...)
	}

	row := d.db.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return nil, err
	}
	return &Row{row}, nil
}

func (d *DB) Tx(ctx context.Context, f func(context.Context) error) error {
	tx := d.tx(ctx)
	if tx != nil {
		return f(ctx)
	}

	rawTx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil
	}
	tx = &Tx{tx: rawTx}
	return tx.Tx(func() error {
		return f(context.WithValue(ctx, txkey{}, tx))
	})
}

type Tx struct {
	tx *sql.Tx
}

func (t *Tx) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *Tx) Query(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{rows}, nil
}

func (t *Tx) QueryRow(ctx context.Context, query string, args ...interface{}) (*Row, error) {
	row := t.tx.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return nil, err
	}
	return &Row{row}, nil
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Tx(f func() error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			_ = t.Rollback()
			panic(p)
		} else if err != nil {
			_ = t.Rollback()
		} else {
			err = t.Commit()
		}
	}()

	err = f()
	return
}
