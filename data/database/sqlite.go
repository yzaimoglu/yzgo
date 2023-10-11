package data

import (
	"fmt"
	"os"
	"time"

	"github.com/gookit/slog"
	"github.com/yzaimoglu/yzgo/config"
	data "github.com/yzaimoglu/yzgo/data"
	data_const "github.com/yzaimoglu/yzgo/data/const"
	"github.com/yzaimoglu/yzgo/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteConfig struct {
	Database string
}

type SQLite struct {
	Config   *SQLiteConfig
	Database *gorm.DB
}

func NewSQLite() *SQLite {
ConnectionStart:
	if config.GetBoolean(data_const.EnvDBEnabled) {
		dbConfig := SQLiteConfig{
			Database: config.GetString(data_const.EnvDBDatabase),
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

		path := fmt.Sprintf("db/%s.db?journal_mode=WAL&busy_timeout=5000", dbConfig.Database)
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			slog.Errorf("Error connecting to SQLite: %v", err)
			slog.Errorf("Retrying in 10 seconds...")
			time.Sleep(10 * time.Second)
			config.Load()
			goto ConnectionStart
		}

		utils.NonChildExec(func() {
			slog.Infof("Connected to sqlite at path %s", path)
		})
		return &SQLite{
			Config:   &dbConfig,
			Database: db,
		}
	}
	return nil
}

func (db *SQLite) Init() error {
	return db.Database.AutoMigrate(&data.Placeholder{})
}

func (db *SQLite) Close() error {
	rawDb, err := db.Database.DB()
	if err != nil {
		return err
	}
	return rawDb.Close()
}

func (db *SQLite) Orm() *gorm.DB {
	return db.Database
}
