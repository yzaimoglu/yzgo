package main

import (
	"github.com/yzaimoglu/yzgo/config"
	data_const "github.com/yzaimoglu/yzgo/data/const"
	"github.com/yzaimoglu/yzgo/http/controllers"
	"github.com/yzaimoglu/yzgo/internal"
)

func main() {
	config.Load()
	switch config.GetString(data_const.EnvDBType) {
	case data_const.DBTypeNoSQL:
		backend := internal.NewNoSQLBackend()
		nosqlController := controllers.NewNoSQLController(backend)
		nosqlController.Start()
	case data_const.DBTypeSQL:
		backend := internal.NewSQLBackend()
		sqlController := controllers.NewSQLController(backend)
		sqlController.Start()
	default:
		panic("Unknown database type.")
	}
}
