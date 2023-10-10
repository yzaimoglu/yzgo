package data

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gookit/slog"
	"github.com/yzaimoglu/yzgo/config"
	data "github.com/yzaimoglu/yzgo/data/const"
	"github.com/yzaimoglu/yzgo/utils"

	_ "github.com/glebarez/go-sqlite"
)

type SQLiteConfig struct {
	Database string
}

type SQLite struct {
	Config   *SQLiteConfig
	Database *sql.DB
}

func NewSQLite() *SQLite {
ConnectionStart:
	if config.GetBoolean(data.EnvDBEnabled) {
		dbConfig := SQLiteConfig{
			Database: config.GetString(data.EnvDBDatabase),
		}

		if _, err := os.Stat("db"); os.IsNotExist(err) {
			if err := os.MkdirAll("db", 0755); err != nil {
				slog.Errorf("Error creating db directory: %v", err)
				slog.Errorf("Retrying in 10 seconds...")
				time.Sleep(10 * time.Second)
				config.Load()
				goto ConnectionStart
			}
		}

		path := fmt.Sprintf("db/%s.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", dbConfig.Database)
		db, err := sql.Open("sqlite", path)
		if err != nil {
			slog.Errorf("Error connecting to SQLite: %v", err)
			slog.Errorf("Retrying in 10 seconds...")
			time.Sleep(10 * time.Second)
			config.Load()
			goto ConnectionStart
		}

		utils.NonChildExec(func() {
			var version string
			versionQuery := db.QueryRow("select sqlite_version()")
			if err := versionQuery.Scan(&version); err != nil {
				slog.Errorf("Error connecting to SQLite: %v", err)
				panic(err)
			}
			slog.Infof("Connected to sqlite v%s at path %s", version, path)
		})
		return &SQLite{
			Config:   &dbConfig,
			Database: db,
		}
	}
	return nil
}

func (db *SQLite) Init() error {
	return nil
}

func (db *SQLite) Close() error {
	return nil
}
