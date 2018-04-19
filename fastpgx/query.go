

package runtime

import (
	"context"
	"errors"
	"github.com/jackc/pgx"
)

const (
	existingTXError = "tx allready in progress"
	noTXError       = "no tx in progress"
)

type Queryer interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	Exec(sql string, arguments ...interface{}) error
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Begin() error
	Commit() error
	Rollback() error
}

type queryer struct {
	pool *pgx.ConnPool
	tx   *pgx.Tx
}

func (p* queryer) Exec(sql string, arguments ...interface{}) (err error) {
	if p.tx != nil {
		_, err = p.tx.Exec(sql, arguments)
	}
	_, err = p.pool.Exec(sql, arguments)
	return
}

func (p *queryer) Query(sql string, args ...interface{}) (*pgx.Rows, error) {
	if p.tx != nil {
		return p.tx.Query(sql, args)
	}
	return p.pool.Query(sql, args)
}

func (p *queryer) QueryRow(sql string, args ...interface{}) *pgx.Row {
	if p.tx != nil {
		return p.tx.QueryRow(sql, args)
	}
	return p.pool.QueryRow(sql, args)
}

func (q *queryer) Begin() (err error) {
	if q.tx != nil {
		return errors.New(existingTXError)
	}
	q.tx, err = q.pool.Begin()
	return err
}

func (q *queryer) Commit() (err error) {
	if q.tx == nil {
		return errors.New(noTXError)
	}
	err = q.tx.Commit()
	q.tx = nil
	return
}

func (q *queryer) Rollback() (err error) {
	if q.tx == nil {
		return errors.New(noTXError)
	}
	err = q.tx.Rollback()
	q.tx = nil
	return
}

func NewContextWithPool(parent context.Context, pool *pgx.ConnPool) (ctx context.Context, q Queryer) {
	q = &queryer{pool: pool}
	ctx = context.WithValue(parent, queryer{}, q)
	return
}

func QueryerFromContext(context context.Context) Queryer {
	return context.Value(queryer{}).(*queryer)
}


