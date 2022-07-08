package xorm

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"skuld/xerr"
)

type Core interface {
	Builder

	Debug(ctx context.Context) Core

	Create(ctx context.Context, value interface{}) error
	CreateInBatches(ctx context.Context, value interface{}, batchSize int) error
	Save(ctx context.Context, value interface{}) error
	First(ctx context.Context, dest interface{}, conds ...interface{}) error
	Find(ctx context.Context, dest interface{}, conds ...interface{}) error
	Update(ctx context.Context, column string, value interface{}) error
	UpdateAffected(ctx context.Context, column string, value interface{}) (int64, error)
	Updates(ctx context.Context, values interface{}) error
	UpdateColumn(ctx context.Context, column string, value interface{}) error
	UpdateColumns(ctx context.Context, values interface{}) error
	Delete(ctx context.Context, value interface{}, conds ...interface{}) error
	Count(ctx context.Context, count *int64) error
	Exec(ctx context.Context, sql string, values ...interface{}) error
}

type Builder interface {
	Model(ctx context.Context, value interface{}) Core
	Table(ctx context.Context, name string, args ...interface{}) Core
	Distinct(ctx context.Context, args ...interface{}) Core
	Select(ctx context.Context, query interface{}, args ...interface{}) Core
	Omit(ctx context.Context, columns ...string) Core
	Where(ctx context.Context, query interface{}, args ...interface{}) Core
	Not(ctx context.Context, query interface{}, args ...interface{}) Core
	Or(ctx context.Context, query interface{}, args ...interface{}) Core
	Joins(ctx context.Context, query string, args ...interface{}) Core
	Group(ctx context.Context, name string) Core
	Having(ctx context.Context, query interface{}, args ...interface{}) Core
	Order(ctx context.Context, value interface{}) Core
	Limit(ctx context.Context, limit int) Core
	Offset(ctx context.Context, offset int) Core
	Unscoped(ctx context.Context) Core
	Raw(ctx context.Context, sql string, values ...interface{}) Core
}

type ORMItf interface {
	Core
	Tx(ctx context.Context, f func(context.Context) error) error
}

type TxORMItf interface {
	Core
	Rollback() error
	Commit() error
	Tx(f func() error) error
}

type ORM struct {
	orm *gorm.DB
}

type txkey struct{}

func New(orm *gorm.DB) ORMItf {
	return &ORM{orm: orm}
}

func handleErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return xerr.NoData
	}
	return err
}

func (x *ORM) tx(ctx context.Context) TxORMItf {
	if v, ok := ctx.Value(txkey{}).(TxORMItf); ok {
		return v
	}
	return nil
}

func (x *ORM) Tx(ctx context.Context, f func(context.Context) error) error {
	if tx := x.tx(ctx); tx != nil {
		return f(ctx)
	}

	tx := &TxORM{orm: x.orm.Begin()}
	return tx.Tx(func() error {
		return f(context.WithValue(ctx, txkey{}, tx))
	})
}

func (x *ORM) Debug(ctx context.Context) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Debug(ctx)
	}

	return &ORM{orm: x.orm.Debug()}
}

func (x *ORM) Create(ctx context.Context, value interface{}) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.Create(ctx, value)
	}

	return handleErr(x.orm.WithContext(ctx).Create(value).Error)
}

func (x *ORM) CreateInBatches(ctx context.Context, value interface{}, batchSize int) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.CreateInBatches(ctx, value, batchSize)
	}

	return handleErr(x.orm.WithContext(ctx).CreateInBatches(value, batchSize).Error)
}

func (x *ORM) Save(ctx context.Context, value interface{}) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.Save(ctx, value)
	}

	return handleErr(x.orm.WithContext(ctx).Save(value).Error)
}

func (x *ORM) First(ctx context.Context, dest interface{}, conds ...interface{}) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.First(ctx, dest, conds...)
	}

	return handleErr(x.orm.WithContext(ctx).First(dest, conds...).Error)
}

func (x *ORM) Find(ctx context.Context, dest interface{}, conds ...interface{}) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.Find(ctx, dest, conds...)
	}

	return handleErr(x.orm.WithContext(ctx).Find(dest, conds...).Error)
}

func (x *ORM) Update(ctx context.Context, column string, value interface{}) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.Update(ctx, column, value)
	}

	return handleErr(x.orm.WithContext(ctx).Update(column, value).Error)
}

func (x *ORM) UpdateAffected(ctx context.Context, column string, value interface{}) (int64, error) {
	if tx := x.tx(ctx); tx != nil {
		return tx.UpdateAffected(ctx, column, value)
	}

	rst := x.orm.WithContext(ctx).Update(column, value)
	return rst.RowsAffected, handleErr(rst.Error)
}

func (x *ORM) Updates(ctx context.Context, values interface{}) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.Updates(ctx, values)
	}

	return handleErr(x.orm.WithContext(ctx).Updates(values).Error)
}

func (x *ORM) UpdateColumn(ctx context.Context, column string, value interface{}) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.UpdateColumn(ctx, column, value)
	}

	return handleErr(x.orm.WithContext(ctx).UpdateColumn(column, value).Error)
}

func (x *ORM) UpdateColumns(ctx context.Context, values interface{}) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.UpdateColumns(ctx, values)
	}

	return handleErr(x.orm.WithContext(ctx).UpdateColumns(values).Error)
}

func (x *ORM) Delete(ctx context.Context, value interface{}, conds ...interface{}) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.Delete(ctx, value, conds...)
	}

	return handleErr(x.orm.WithContext(ctx).Delete(value, conds...).Error)
}

func (x *ORM) Count(ctx context.Context, count *int64) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.Count(ctx, count)
	}

	return handleErr(x.orm.WithContext(ctx).Count(count).Error)
}

func (x *ORM) Exec(ctx context.Context, sql string, values ...interface{}) error {
	if tx := x.tx(ctx); tx != nil {
		return tx.Exec(ctx, sql, values...)
	}

	return handleErr(x.orm.WithContext(ctx).Exec(sql, values...).Error)
}

func (x *ORM) Model(ctx context.Context, value interface{}) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Model(ctx, value)
	}

	return &ORM{orm: x.orm.Model(value)}
}

func (x *ORM) Table(ctx context.Context, name string, args ...interface{}) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Table(ctx, name, args...)
	}

	return &ORM{orm: x.orm.Table(name, args...)}
}

func (x *ORM) Distinct(ctx context.Context, args ...interface{}) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Distinct(ctx, args...)
	}

	return &ORM{orm: x.orm.Distinct(args...)}
}

func (x *ORM) Select(ctx context.Context, query interface{}, args ...interface{}) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Select(ctx, query, args...)
	}

	return &ORM{orm: x.orm.Select(query, args...)}
}

func (x *ORM) Omit(ctx context.Context, columns ...string) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Omit(ctx, columns...)
	}

	return &ORM{orm: x.orm.Omit(columns...)}
}

func (x *ORM) Where(ctx context.Context, query interface{}, args ...interface{}) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Where(ctx, query, args...)
	}

	return &ORM{orm: x.orm.Where(query, args...)}
}

func (x *ORM) Not(ctx context.Context, query interface{}, args ...interface{}) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Not(ctx, query, args...)
	}

	return &ORM{orm: x.orm.Not(query, args...)}
}

func (x *ORM) Or(ctx context.Context, query interface{}, args ...interface{}) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Or(ctx, query, args...)
	}

	return &ORM{orm: x.orm.Or(query, args...)}
}

func (x *ORM) Joins(ctx context.Context, query string, args ...interface{}) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Joins(ctx, query, args...)
	}

	return &ORM{orm: x.orm.Joins(query, args...)}
}

func (x *ORM) Group(ctx context.Context, name string) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Group(ctx, name)
	}

	return &ORM{orm: x.orm.Group(name)}
}

func (x *ORM) Having(ctx context.Context, query interface{}, args ...interface{}) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Having(ctx, query, args...)
	}

	return &ORM{orm: x.orm.Having(query, args...)}
}

func (x *ORM) Order(ctx context.Context, value interface{}) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Order(ctx, value)
	}

	return &ORM{orm: x.orm.Order(value)}
}

func (x *ORM) Limit(ctx context.Context, limit int) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Limit(ctx, limit)
	}

	return &ORM{orm: x.orm.Limit(limit)}
}

func (x *ORM) Offset(ctx context.Context, offset int) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Offset(ctx, offset)
	}

	return &ORM{orm: x.orm.Offset(offset)}
}

func (x *ORM) Unscoped(ctx context.Context) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Unscoped(ctx)
	}

	return &ORM{orm: x.orm.Unscoped()}
}

func (x *ORM) Raw(ctx context.Context, sql string, values ...interface{}) Core {
	if tx := x.tx(ctx); tx != nil {
		return tx.Raw(ctx, sql, values...)
	}

	return &ORM{x.orm.WithContext(ctx).Raw(sql, values...)}
}

type TxORM struct {
	orm *gorm.DB
}

func (x *TxORM) Debug(ctx context.Context) Core {
	return &TxORM{orm: x.orm.Debug()}
}

func (x *TxORM) Rollback() error {
	return x.orm.Rollback().Error
}

func (x *TxORM) Commit() error {
	return x.orm.Commit().Error
}

func (x *TxORM) Tx(f func() error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			_ = x.Rollback()
			panic(p)
		} else if err != nil {
			_ = x.Rollback()
		} else {
			err = x.Commit()
		}
	}()

	err = f()
	return
}

func (x *TxORM) Create(ctx context.Context, value interface{}) error {
	return handleErr(x.orm.WithContext(ctx).Create(value).Error)
}

func (x *TxORM) CreateInBatches(ctx context.Context, value interface{}, batchSize int) error {
	return handleErr(x.orm.WithContext(ctx).CreateInBatches(value, batchSize).Error)
}

func (x *TxORM) Save(ctx context.Context, value interface{}) error {
	return handleErr(x.orm.WithContext(ctx).Save(value).Error)
}

func (x *TxORM) First(ctx context.Context, dest interface{}, conds ...interface{}) error {
	return handleErr(x.orm.WithContext(ctx).First(dest, conds...).Error)
}

func (x *TxORM) Find(ctx context.Context, dest interface{}, conds ...interface{}) error {
	return handleErr(x.orm.WithContext(ctx).Find(dest, conds...).Error)
}

func (x *TxORM) Update(ctx context.Context, column string, value interface{}) error {
	return handleErr(x.orm.WithContext(ctx).Update(column, value).Error)
}

func (x *TxORM) UpdateAffected(ctx context.Context, column string, value interface{}) (int64, error) {
	rst := x.orm.WithContext(ctx).Update(column, value)
	return rst.RowsAffected, handleErr(rst.Error)
}

func (x *TxORM) Updates(ctx context.Context, values interface{}) error {
	return handleErr(x.orm.WithContext(ctx).Updates(values).Error)
}

func (x *TxORM) UpdateColumn(ctx context.Context, column string, value interface{}) error {
	return handleErr(x.orm.WithContext(ctx).UpdateColumn(column, value).Error)
}

func (x *TxORM) UpdateColumns(ctx context.Context, values interface{}) error {
	return handleErr(x.orm.WithContext(ctx).UpdateColumns(values).Error)
}

func (x *TxORM) Delete(ctx context.Context, value interface{}, conds ...interface{}) error {
	return handleErr(x.orm.WithContext(ctx).Delete(value, conds...).Error)
}

func (x *TxORM) Count(ctx context.Context, count *int64) error {
	return handleErr(x.orm.WithContext(ctx).Count(count).Error)
}

func (x *TxORM) Exec(ctx context.Context, sql string, values ...interface{}) error {
	return handleErr(x.orm.WithContext(ctx).Exec(sql, values...).Error)
}

func (x *TxORM) Model(ctx context.Context, value interface{}) Core {
	return &TxORM{orm: x.orm.Model(value)}
}

func (x *TxORM) Table(ctx context.Context, name string, args ...interface{}) Core {
	return &TxORM{orm: x.orm.Table(name, args...)}
}

func (x *TxORM) Distinct(ctx context.Context, args ...interface{}) Core {
	return &TxORM{orm: x.orm.Distinct(args...)}
}

func (x *TxORM) Select(ctx context.Context, query interface{}, args ...interface{}) Core {
	return &TxORM{orm: x.orm.Select(query, args...)}
}

func (x *TxORM) Omit(ctx context.Context, columns ...string) Core {
	return &TxORM{orm: x.orm.Omit(columns...)}
}

func (x *TxORM) Where(ctx context.Context, query interface{}, args ...interface{}) Core {
	return &TxORM{orm: x.orm.Where(query, args...)}
}

func (x *TxORM) Not(ctx context.Context, query interface{}, args ...interface{}) Core {
	return &TxORM{orm: x.orm.Not(query, args...)}
}

func (x *TxORM) Or(ctx context.Context, query interface{}, args ...interface{}) Core {
	return &TxORM{orm: x.orm.Or(query, args...)}
}

func (x *TxORM) Joins(ctx context.Context, query string, args ...interface{}) Core {
	return &TxORM{orm: x.orm.Joins(query, args...)}
}

func (x *TxORM) Group(ctx context.Context, name string) Core {
	return &TxORM{orm: x.orm.Group(name)}
}

func (x *TxORM) Having(ctx context.Context, query interface{}, args ...interface{}) Core {
	return &TxORM{orm: x.orm.Having(query, args...)}
}

func (x *TxORM) Order(ctx context.Context, value interface{}) Core {
	return &TxORM{orm: x.orm.Order(value)}
}

func (x *TxORM) Limit(ctx context.Context, limit int) Core {
	return &TxORM{orm: x.orm.Limit(limit)}
}

func (x *TxORM) Offset(ctx context.Context, offset int) Core {
	return &TxORM{orm: x.orm.Offset(offset)}
}

func (x *TxORM) Unscoped(ctx context.Context) Core {
	return &TxORM{orm: x.orm.Unscoped()}
}

func (x *TxORM) Raw(ctx context.Context, sql string, values ...interface{}) Core {
	return &TxORM{orm: x.orm.WithContext(ctx).Raw(sql, values...)}
}
