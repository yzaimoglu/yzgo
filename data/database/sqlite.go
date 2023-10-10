package data

import "database/sql"

type SQLiteConfig struct {
	Host      string
	Port      int
	User      string
	Password  string
	Namespace string
	Database  string
}

type SQLite struct {
	Config   *SurrealDBConfig
	Database *sql.DB
}

func NewSQLite() *SQLite {
	return &SQLite{}
}

func (db *SQLite) Init() error {
	return nil
}

func (db *SQLite) Close() error {
	return nil
}
