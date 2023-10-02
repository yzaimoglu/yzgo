package internal

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
	"github.com/yzaimoglu/yzgo/config"
	"github.com/yzaimoglu/yzgo/utils"
)

type Logger struct {
	Info    func(format string, args ...any)
	Debug   func(format string, args ...any)
	Warning func(format string, args ...any)
	Error   func(format string, args ...any)
	Fatal   func(format string, args ...any)
	File    *os.File
}

func NewLogger() *Logger {
	utils.NonChildExec(func() {
		config.ConfigLogger()
	})
	return &Logger{
		Info: slog.Infof,
		Debug: func(format string, args ...any) {
			if utils.Debug() && !fiber.IsChild() {
				slog.Debugf(format, args...)
			}
		},
		Warning: slog.Warnf,
		Error:   slog.Errorf,
		Fatal:   slog.Fatalf,
	}
}

func (logger *Logger) SetFile(file *os.File) {
	logger.File = file
}
