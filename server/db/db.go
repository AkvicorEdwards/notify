package db

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DEFAULT *MySQL

func SetDEFAULT(c *MySQL) {
	DEFAULT = c
}

func Open(m *MySQL) (*gorm.DB, error) {
	return gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=%s&loc=%s",
		m.User,
		m.Password,
		m.Host,
		m.DBName,
		m.Charset,
		m.ParseTime,
		m.Loc,
	))
}

func Close(db *gorm.DB) {
	_ = db.Close()
}

func NewRecord(value interface{}) bool {
	db, err := Open(DEFAULT)
	if err != nil {
		return false
	}
	defer Close(db)
	return db.NewRecord(value)
}

func Exec(sql string, values ...interface{}) (*gorm.DB, error) {
	db, err := Open(DEFAULT)
	if err != nil {
		return nil, err
	}
	defer Close(db)
	gdb := db.Exec(sql, values...)
	return gdb, gdb.Error
}

func Row(st interface{}, sql string, values ...interface{}) (*gorm.DB, error) {
	db, err := Open(DEFAULT)
	if err != nil {
		return nil, err
	}
	defer Close(db)
	gdb := db.Raw(sql, values...).Scan(st)
	return gdb, gdb.Error
}

func Query(sql string, values ...interface{}) (*sql.Rows, error) {
	db, err := Open(DEFAULT)
	if err != nil {
		return nil, err
	}
	defer Close(db)
	gdb := db.Raw(sql, values...)
	gdbs, err := gdb.Rows()

	return gdbs, nil
}