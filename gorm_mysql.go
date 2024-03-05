package gm

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	gormDB *gorm.DB
)

func InitGorm(host string, port int64, user, password, db_name string) (err error) {
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, db_name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm.Open err, msg: %v", err)
	}

	if gormDB == nil {
		gormDB = db
	}

	return
}

func NewGormSession() *gorm.DB {
	return gormDB.Session(&gorm.Session{})
}

func BeginGormTx() *gorm.DB {
	return gormDB.Begin()
}

func CommitGormTx(tx *gorm.DB) error {
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("tx.Commit %v", err)
	}
	return nil
}
