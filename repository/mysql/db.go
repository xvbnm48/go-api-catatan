package mysql

import (
	"database/sql"
	"fmt"
	"github.com/go-kit/kit/log/level"
	shv "github.com/xvbnm48/go-api-catatan/helper"
	_interface "github.com/xvbnm48/go-api-catatan/repository/interface"
	"time"
)

type dbReadWriter struct {
	db *sql.DB
}

func NewDbReadWriter(url, port, schema, user, password string) (_interface.ReadWriter, error) {
	mysqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", url, port, user, password, schema)
	db, err := sql.Open("mysql", mysqlConn)
	if err != nil {
		level.Error(shv.Logger).Log(fmt.Sprintf("%s", err.Error()))
		return nil, err
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &dbReadWriter{
		db: db,
	}, nil
}

func closeRows(rs *sql.Rows) {
	if rs != nil {
		if err := rs.Close(); err != nil {
			_ = level.Error(shv.Logger).Log(fmt.Sprintf("error while closing result set %+v", err.Error()))
		}
	}
}

func rollbackTx(tx *sql.Tx) {
	if tx == nil {
		return
	}

	if err := tx.Rollback(); err != nil {
		// _ = level.Error(shv.Logger).Log(fmt.Sprintf("error while rolling back transaction %+v", err.Error()))
	}
}

// Close is used for closing the sql connection
func (rw *dbReadWriter) Close() error {
	if rw.db != nil {
		if err := rw.db.Close(); err != nil {
			return err
		}
		rw.db = nil
	}

	return nil
}
func (rw *dbReadWriter) Begin() (*sql.Tx, error) {
	return rw.db.Begin()
}

func (rw *dbReadWriter) Commit(tx *sql.Tx) error {
	return tx.Commit()
}
func (rw *dbReadWriter) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}
