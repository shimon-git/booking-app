package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbconn = &DB{}

const maxOpenDbConn = 10
const maxIdledbConn = 5
const maxDbLifeTime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDB(dsn)
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdledbConn)
	d.SetConnMaxLifetime(maxDbLifeTime)
	dbconn.SQL = d

	if err := TestDB(d); err != nil {
		return nil, err
	}
	return dbconn, nil
}

func TestDB(d *sql.DB) error {
	if err := d.Ping(); err != nil {
		return err
	}
	return nil
}

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = TestDB(db); err != nil {
		return nil, err
	}
	return db, nil
}
