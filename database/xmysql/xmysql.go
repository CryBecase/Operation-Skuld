package xmysql

import (
	"database/sql"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"skuld/xorm"
	"skuld/xsql"
)

func NewDB(c Config) xsql.DBItf {
	db, err := sql.Open("mysql", c.Source)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxIdleTime(c.IdleTime)
	db.SetMaxIdleConns(c.Idle)
	db.SetMaxOpenConns(c.Open)

	return xsql.New(db)
}

func NewORM(c Config) xorm.ORMItf {
	var err error
	orm, err := gorm.Open(mysql.Open(c.Source), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				Colorful:                  false,
				IgnoreRecordNotFoundError: true,
				LogLevel:                  logger.Warn,
			}),
	})
	if err != nil {
		panic(err)
	}
	rawdb, err := orm.DB()
	if err != nil {
		panic(err)
	}

	rawdb.SetConnMaxIdleTime(c.IdleTime)
	rawdb.SetMaxIdleConns(c.Idle)
	rawdb.SetMaxOpenConns(c.Open)

	return xorm.New(orm)
}
