package config

import (
	"os"
	"strconv"
	"time"
)

type DB struct {
	Host            string
	Port            int
	SslMode         string
	Name            string
	User            string
	Password        string
	Debug           bool
	MaxOpenConn     int
	MaxIdleConn     int
	MaxConnLifetime time.Duration
}

var db = &DB{}

func DBCfg() *DB {
	return db
}

func LoadDBCfg() {
	db.Host = os.Getenv("DB_HOST")
	db.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	db.User = os.Getenv("DB_USER")
	db.Password = os.Getenv("DB_PASSWORD")
	db.Name = os.Getenv("DB_NAME")
	db.SslMode = os.Getenv("DB_SSL_MODE")
	db.Debug, _ = strconv.ParseBool(os.Getenv("DB_DEBUG"))
}
