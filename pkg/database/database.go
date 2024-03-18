package database

import (
	"fmt"

	"github.com/arif-x/sqlx-mysql-boilerplate/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct{ *sqlx.DB }

var defaultDB = &DB{}

func (db *DB) connect(cfg *config.DB) (err error) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db.DB, err = sqlx.Open("mysql", dbURI)
	if err != nil {
		return err
	}

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return fmt.Errorf("can't sent ping to database, %w", err)
	}

	return nil
}

// GetDB returns db instance
func GetDB() *DB {
	return defaultDB
}

// ConnectDB sets the db client of database using default configuration
func ConnectDB() error {
	return defaultDB.connect(config.DBCfg())
}
