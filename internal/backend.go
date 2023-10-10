package internal

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/stripe/stripe-go/v74"
	"github.com/yzaimoglu/yzgo/config"
	data_const "github.com/yzaimoglu/yzgo/data/const"
	data_db "github.com/yzaimoglu/yzgo/data/database"
	"github.com/yzaimoglu/yzgo/utils"
)

type Backend interface {
	Start()
	Stop()
	SysCallSetup()
}

type NoSQLBackend struct {
	Fiber    *fiber.App
	Database data_db.NoSQLDatabase
	Logger   *Logger
	HTTP     *HTTPResponder
	Syscall  func(signal os.Signal, backend *NoSQLBackend)
}

type SQLBackend struct {
	Fiber    *fiber.App
	Database data_db.SQLDatabase
	Logger   *Logger
	HTTP     *HTTPResponder
	Syscall  func(signal os.Signal, backend *SQLBackend)
}

func BackendSetup() (*fiber.App, *Logger) {
	prefork := false
	if config.GetBoolean(data_const.EnvPrefork) {
		prefork = true
	}

	fiberBackend := fiber.New(fiber.Config{
		Prefork:      prefork,
		ServerHeader: config.GetString(data_const.EnvServerHeader),
		AppName:      config.GetString(data_const.EnvAppName),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(404).JSON(&fiber.Map{
				"status":  404,
				"message": "Not Found",
			})
		},
	})

	loggerBackend := NewLogger()

	// Security and compression middleware
	if config.GetBoolean(data_const.EnvHelmetEnabled) {
		fiberBackend.Use(helmet.New(helmet.ConfigDefault))
	}
	if config.GetBoolean(data_const.EnvCompressEnabled) {
		fiberBackend.Use(compress.New())
	}
	if config.GetBoolean(data_const.EnvRecoverEnabled) {
		fiberBackend.Use(recover.New())
	}
	if config.GetBoolean(data_const.EnvCsrfEnabled) {
		fiberBackend.Use(csrf.New())
	}

	if config.GetBoolean(data_const.EnvLogEnabled) {
		file, err := os.OpenFile(fmt.Sprintf("./%s/%s", config.GetString(data_const.EnvLogPath), config.GetString(data_const.EnvLogFile)), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			loggerBackend.Fatal("Error opening file: %v", err)
		}
		fiberBackend.Use(logger.New(logger.Config{
			Format:     config.GetString(data_const.EnvLogFormat),
			TimeFormat: config.GetString(data_const.EnvTimeFormat),
			TimeZone:   config.GetString(data_const.EnvTimezone),
			Output:     file,
		}))
		loggerBackend.SetFile(file)
	}

	if config.GetBoolean(data_const.EnvMetricsEnabled) {
		metrics := fiberBackend.Group("/api/v1/metrics")
		metrics.Use(basicauth.New(basicauth.Config{
			Users: map[string]string{
				config.GetString(data_const.EnvMetricsUser): config.GetString(data_const.EnvMetricsPassword),
			},
		}))
		metrics.Get("/", monitor.New())
	}

	if config.GetBoolean(data_const.EnvCorsEnabled) {
		fiberBackend.Use(cors.New(cors.Config{
			AllowCredentials: config.GetBoolean(data_const.EnvCorsAllowCredentials),
			AllowOrigins:     config.GetString(data_const.EnvCorsAllowOrigins),
			AllowMethods:     config.GetString(data_const.EnvCorsAllowMethods),
		}))
	}

	utils.StripeEnabledExec(func() {
		stripe.Key = config.GetString(data_const.EnvStripeSecret)
	})

	return fiberBackend, loggerBackend
}
