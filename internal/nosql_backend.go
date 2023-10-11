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

func NewNoSQLBackend() *NoSQLBackend {
	return &NoSQLBackend{
		Logger:  NewLogger(),
		HTTP:    NewHTTPResponder(),
		Syscall: NewNoSQLSyscallHandler(),
	}
}

func (backend *NoSQLBackend) Start() {
	switch config.GetString(data_const.EnvDB) {
	case data_const.DBTypeSurreal:
		backend.Database = data_db.NewSurrealDB()
	case data_const.DBTypeMongo:
		backend.Logger.Fatal("MongoDB is not supported yet.")
	case data_const.DBTypeArango:
		backend.Logger.Fatal("ArangoDB is not supported yet.")
	default:
		backend.Logger.Fatal("Unknown database %s for database type %s", config.GetString(data_const.EnvDB), config.GetString(data_const.EnvDBType))
		return
	}

	go backend.StartServer()
	backend.SysCallSetup()
}

func (backend *NoSQLBackend) StartServer() {
	backend.Database.Init()
	backend.Logger.Debug("Server starting on port %d.", config.GetInteger(data_const.EnvPort))
	if err := backend.Fiber.Listen(fmt.Sprintf(":%d", config.GetInteger(data_const.EnvPort))); err != nil {
		slog.Fatalf("Error starting server: %v", err)
	}
}

func (backend *NoSQLBackend) Stop() {
	backend.Fiber.Shutdown()
	backend.Logger.File.Close()
	backend.Database.Close()
	slog.Flush()
}

func (backend *NoSQLBackend) SysCallSetup() {
	// Wait for a signal to exit.
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
