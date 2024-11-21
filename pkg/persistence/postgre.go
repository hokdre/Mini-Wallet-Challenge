package persistence

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq" // Example for PostgreSQL; use the driver for your DB
)

var DB *sql.DB
var onceDB sync.Once

type Config struct {
	Host        string
	Port        string
	Username    string
	Password    string
	DB          string
	SSLMode     string
	MaxIdleConn int
	MaxOpenConn int
}

func OpenPostgreDB(cfg Config) (*sql.DB, error) {
	var err error
	onceDB.Do(func() {
		db, errDB := sql.Open("postgres", fmt.Sprintf(`
		host=%s
		port=%s
		user=%s
		password=%s
		dbname=%s
		sslmode=%s
		`,
			cfg.Host,
			cfg.Port,
			cfg.Username,
			cfg.Password,
			cfg.DB,
			cfg.SSLMode,
		))
		DB = db
		err = errDB
	})
	if err != nil {
		return nil, err
	}

	return DB, nil
}
