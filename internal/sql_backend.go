package internal

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gookit/slog"
	"github.com/yzaimoglu/yzgo/config"
	data_const "github.com/yzaimoglu/yzgo/data/const"
	data_db "github.com/yzaimoglu/yzgo/data/database"
)

func NewSQLBackend() *SQLBackend {
	return &SQLBackend{
		Logger:  NewLogger(),
		HTTP:    NewHTTPResponder(),
		Syscall: NewSQLSyscallHandler(),
	}
}

func (backend *SQLBackend) Start() {
	switch config.GetString(data_const.EnvDB) {
	case data_const.DBTypeSQLite:
		backend.Database = data_db.NewSQLite()
	case data_const.DBTypeMySQL:
		backend.Logger.Fatal("MySQL is not supported yet.")
		return
	case data_const.DBTypePostgres:
		backend.Logger.Fatal("PostgreSQL is not supported yet.")
		return
	// case data_const.DBTypeRQLite:
	// 	backend.Database = data_db.NewRQLite()
	default:
		backend.Logger.Fatal("Unknown database %s for database type %s", config.GetString(data_const.EnvDB), config.GetString(data_const.EnvDBType))
		return
	}

	go backend.StartServer()
	backend.SysCallSetup()
}

func (backend *SQLBackend) StartServer() {
	backend.Database.Init()
	backend.Logger.Debug("Server starting on port %d.", config.GetInteger(data_const.EnvPort))
	if err := backend.Fiber.Listen(fmt.Sprintf(":%d", config.GetInteger(data_const.EnvPort))); err != nil {
		slog.Fatalf("Error starting server: %v", err)
	}
}

func (backend *SQLBackend) Stop() {
	backend.Fiber.Shutdown()
	backend.Logger.File.Close()
	backend.Database.Close()
	slog.Flush()
}

func (backend *SQLBackend) SysCallSetup() {
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)
	exitchnl := make(chan int)

	go func() {
		for {
			s := <-sigchnl
			backend.Syscall(s, backend)
		}
	}()

	exitcode := <-exitchnl
	os.Exit(exitcode)
}
