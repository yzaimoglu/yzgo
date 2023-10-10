package internal

import (
	"os"
	"syscall"

	"github.com/yzaimoglu/yzgo/utils"
)

func NewNoSQLSyscallHandler() func(signal os.Signal, backend *NoSQLBackend) {
	return func(signal os.Signal, backend *NoSQLBackend) {
		utils.NonChildExec(func() {
			if signal == syscall.SIGTERM {
				backend.Logger.Info("Got terminate signal.")
				backend.Logger.Info("Terminating the program...")
				backend.Stop()
				os.Exit(0)
			} else if signal == syscall.SIGINT {
				backend.Logger.Info("Got interrupt signal.")
				backend.Logger.Info("Terminating the program...")
				backend.Stop()
				os.Exit(0)
			}
		})
	}
}

func NewSQLSyscallHandler() func(signal os.Signal, backend *SQLBackend) {
	return func(signal os.Signal, backend *SQLBackend) {
		utils.NonChildExec(func() {
			if signal == syscall.SIGTERM {
				backend.Logger.Info("Got terminate signal.")
				backend.Logger.Info("Terminating the program...")
				backend.Stop()
				os.Exit(0)
			} else if signal == syscall.SIGINT {
				backend.Logger.Info("Got interrupt signal.")
				backend.Logger.Info("Terminating the program...")
				backend.Stop()
				os.Exit(0)
			}
		})
	}
}
