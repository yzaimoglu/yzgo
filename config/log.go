package config

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func ConfigLogger() {
	standardHandler := handler.MustRotateFile("./logs/system.log", 60*60*12)
	errorHandler := handler.MustRotateFile("./logs/system_err.log", 60*60*12, handler.WithLogLevels(slog.DangerLevels))

	slog.Infof("Setting up the Logger...")
	slog.PushHandler(standardHandler)
	slog.PushHandler(errorHandler)
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})
	slog.Infof("Logger setup successfully.")
}
