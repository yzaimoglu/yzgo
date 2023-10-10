package data

import (
	"fmt"
	"time"

	"github.com/gookit/slog"
	"github.com/rqlite/gorqlite"
	"github.com/yzaimoglu/yzgo/config"
	data "github.com/yzaimoglu/yzgo/data/const"
	"github.com/yzaimoglu/yzgo/utils"
)

type RQLiteConfig struct {
	Host           string
	Port           int
	User           string
	Password       string
	Namespace      string
	Database       string
	Authentication bool
}

type RQLite struct {
	Config   *RQLiteConfig
	Database *gorqlite.Connection
}

func NewRQLite() *RQLite {
ConnectionStart:
	if config.GetBoolean(data.EnvDBEnabled) {
		dbConfig := RQLiteConfig{
			Host:           config.GetString(data.EnvDBHost),
			Port:           config.GetInteger(data.EnvDBPort),
			User:           config.GetString(data.EnvDBUser),
			Password:       config.GetString(data.EnvDBPassword),
			Namespace:      config.GetString(data.EnvDBNamespace),
			Database:       config.GetString(data.EnvDBDatabase),
			Authentication: config.GetBoolean(data.EnvDBAuthEnabled),
		}
		var url string
		if dbConfig.Authentication {
			url = fmt.Sprintf("http://%s:%s@%s:%d/", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port)
		} else {
			url = fmt.Sprintf("http://%s:%d/", dbConfig.Host, dbConfig.Port)
		}

		db, err := gorqlite.Open(url)
		if err != nil {
			slog.Errorf("Error connecting to RQLite: %v", err)
			slog.Errorf("Retrying in 10 seconds...")
			time.Sleep(10 * time.Second)
			config.Load()
			goto ConnectionStart
		}

		utils.NonChildExec(func() {
			slog.Infof("Connected to rqlite at %s:%d", dbConfig.Host, dbConfig.Port)
		})
		return &RQLite{
			Config:   &dbConfig,
			Database: db,
		}
	}
	return nil
}

func (db *RQLite) Init() error {
	return nil
}

func (db *RQLite) Close() error {
	db.Database.Close()
	return nil
}
