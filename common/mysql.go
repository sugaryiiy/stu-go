package common

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// MySQLConfig holds connection settings for a MySQL database.
type MySQLConfig struct {
	DSN             string
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	PingTimeout     time.Duration
}

// OpenMySQL creates a MySQL connection using the provided configuration.
func OpenMySQL(cfg MySQLConfig) (*xorm.Engine, error) {
	db, err := xorm.NewEngine("mysql", cfg.DSN)
	if err != nil {
		return nil, err
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = time.Hour
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 25
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 25
	}
	if cfg.PingTimeout == 0 {
		cfg.PingTimeout = 5 * time.Second
	}
	//日志打印SQL
	db.ShowSQL(true)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
