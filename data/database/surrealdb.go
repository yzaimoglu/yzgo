package data

import (
	"errors"
	"fmt"
	"time"

	"github.com/gookit/slog"
	"github.com/surrealdb/surrealdb.go"
	"github.com/yzaimoglu/yzgo/config"
	"github.com/yzaimoglu/yzgo/utils"
)

const (
	Limit = 50
)

type SurrealDBConfig struct {
	Host      string
	Port      int
	User      string
	Password  string
	Namespace string
	Database  string
}

type SurrealDB struct {
	Config   *SurrealDBConfig
	Database *surrealdb.DB
}

type SurrealDBResponse struct {
	Result []interface{} `json:"result"`
	Status string        `json:"status"`
	Time   string        `json:"time"`
}

func NewSurrealDB() *SurrealDB {
ConnectionStart:
	if config.GetBoolean("DB_ENABLED") {
		dbConfig := SurrealDBConfig{
			Host:      config.GetString("DB_HOST"),
			Port:      config.GetInteger("DB_PORT"),
			User:      config.GetString("DB_USER"),
			Password:  config.GetString("DB_PASSWORD"),
			Namespace: config.GetString("DB_NAMESPACE"),
			Database:  config.GetString("DB_DATABASE"),
		}
		url := fmt.Sprintf("ws://%s:%d/rpc", dbConfig.Host, dbConfig.Port)
		db, err := surrealdb.New(url)
		if err != nil {
			slog.Errorf("Error connecting to SurrealDB: %v", err)
			slog.Errorf("Retrying in 10 seconds...")
			time.Sleep(10 * time.Second)
			config.Load()
			goto ConnectionStart
		}

		if _, err := db.Signin(map[string]interface{}{
			"user": dbConfig.User,
			"pass": dbConfig.Password,
		}); err != nil {
			slog.Errorf("Error connecting to SurrealDB: %v", err)
			slog.Errorf("Retrying in 10 seconds...")
			time.Sleep(10 * time.Second)
			config.Load()
			goto ConnectionStart
		}

		if _, err = db.Use(dbConfig.Namespace, dbConfig.Database); err != nil {
			slog.Errorf("Error connecting to SurrealDB: %v", err)
			slog.Errorf("Retrying in 10 seconds...")
			time.Sleep(10 * time.Second)
			config.Load()
			goto ConnectionStart
		}

		surrealdb := &SurrealDB{
			Config:   &dbConfig,
			Database: db,
		}

		surrealdb.Init()
		return surrealdb
	}
	return nil
}

func GetSurrealDBResponse(rawData interface{}) (map[string]interface{}, error) {
	var data []interface{}
	var ok bool
	if data, ok = rawData.([]interface{}); !ok {
		return map[string]interface{}{}, fmt.Errorf("failed raw unmarshaling to interface slice: %s", "InvalidResponse")
	}

	var responseObj map[string]interface{}
	if responseObj, ok = data[0].(map[string]interface{}); !ok {
		return map[string]interface{}{}, fmt.Errorf("failed mapping to response object: %s", "InvalidResponse")
	}

	var status string
	if status, ok = responseObj["status"].(string); !ok {
		return map[string]interface{}{}, fmt.Errorf("failed retrieving status: %s", "InvalidResponse")
	}
	if status != "OK" {
		return map[string]interface{}{}, fmt.Errorf("status was not ok: %s", "ErrQuery")
	}

	return responseObj, nil
}

func GetActionResultMap(rawData interface{}) (map[string]interface{}, error) {
	var data []interface{} = rawData.([]interface{})
	var responseObj map[string]interface{}
	var ok bool
	if responseObj, ok = data[0].(map[string]interface{}); !ok {
		return map[string]interface{}{}, fmt.Errorf("failed mapping to response object: %s", "InvalidResponse")
	}

	return responseObj, nil
}

func GetPatchResultMap(rawData interface{}) (map[string]interface{}, error) {
	data, ok := rawData.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, fmt.Errorf("failed raw unmarshaling to interface slice: %s", "InvalidResponse")
	}

	return data, nil
}

func GetSingleResultMap(rawData interface{}) (map[string]interface{}, error) {
	responseObj, err := GetSurrealDBResponse(rawData)
	if err != nil {
		return map[string]interface{}{}, err
	}

	result := responseObj["result"]
	if len(result.([]interface{})) == 0 {
		return map[string]interface{}{}, nil
	}

	resultObj := result.([]interface{})
	if len(resultObj) == 0 {
		return map[string]interface{}{}, nil
	}

	return resultObj[0].(map[string]interface{}), nil
}

func GetMultipleResultMap(rawData interface{}) ([]interface{}, error) {
	responseObj, err := GetSurrealDBResponse(rawData)
	if err != nil {
		return []interface{}{}, err
	}

	result := responseObj["result"]
	if len(result.([]interface{})) == 0 {
		return []interface{}{}, nil
	}

	resultObj := result.([]interface{})
	if len(resultObj) == 0 {
		return []interface{}{}, nil
	}

	return resultObj, nil
}

func (db *SurrealDB) Query(query string, vars interface{}) (interface{}, error) {
	result, err := db.Database.Query(query, vars)
	if err != nil {
		return nil, err
	}
	return GetSurrealDBResponse(result)
}

func (db *SurrealDB) QueryWithoutResponse(query string, vars interface{}) (interface{}, error) {
	result, err := db.Database.Query(query, vars)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db *SurrealDB) Create(table string, data interface{}) (interface{}, error) {
	result, err := db.Database.Create(table, data)
	if err != nil {
		return nil, err
	}

	mappedResult, err := GetActionResultMap(result)
	if err != nil {
		return nil, err
	}

	return mappedResult, nil
}

func (db *SurrealDB) Patch(id string, data interface{}) (interface{}, error) {
	// Needs to be changed to Merge
	result, err := db.Database.Change(id, data)
	if err != nil {
		return nil, err
	}

	mappedResult, err := GetPatchResultMap(result)
	if err != nil {
		return nil, err
	}

	return mappedResult, nil
}

func (db *SurrealDB) Delete(id string) error {
	if _, err := db.Database.Delete(id); err != nil {
		return err
	}
	return nil
}

func (db *SurrealDB) GetSingle(query string, vars interface{}) (interface{}, error) {
	result, err := db.Database.Query(query, vars)
	if err != nil {
		return nil, err
	}

	data, err := GetSingleResultMap(result)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("no results found")
	}

	return data, err
}

func (db *SurrealDB) GetMultiple(query string, vars interface{}) ([]interface{}, error) {
	result, err := db.Database.Query(query, vars)
	if err != nil {
		return nil, err
	}

	data, err := GetMultipleResultMap(result)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []interface{}{}, nil
	}

	return data, err
}

func (db *SurrealDB) Init() error {
	utils.NonChildExec(func() {
		_, err := db.QueryWithoutResponse(`
			BEGIN TRANSACTION;

			COMMIT TRANSACTION;
		`, nil)
		if err != nil {
			return
		}
	})
	return nil
}

func (db *SurrealDB) Close() error {
	db.Database.Close()
	return nil
}
