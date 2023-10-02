package internal

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gookit/slog"
	"github.com/stripe/stripe-go/v74"
	"github.com/yzaimoglu/yzgo/config"
	data_const "github.com/yzaimoglu/yzgo/data/const"
	data_db "github.com/yzaimoglu/yzgo/data/database"
	"github.com/yzaimoglu/yzgo/utils"
)

type Backend struct {
	Fiber    *fiber.App
	Database data_db.Database
	Logger   *Logger
	HTTP     *HTTPResponder
	Syscall  func(signal os.Signal, backend *Backend)
}

func NewBackend() *Backend {
	config.Load()
	return &Backend{
		Logger:  NewLogger(),
		HTTP:    NewHTTPResponder(),
		Syscall: NewSyscallHandler(),
	}
}

func (backend *Backend) Start() {
	prefork := false
	if config.GetBoolean("PREFORK") {
		prefork = true
	}

	backend.Fiber = fiber.New(fiber.Config{
		Prefork:      prefork,
		ServerHeader: config.GetString("SERVER_HEADER"),
		AppName:      config.GetString("APP_NAME"),
	})

	// Security and compression middleware
	if config.GetBoolean("HELMET_ENABLED") {
		backend.Fiber.Use(helmet.New(helmet.ConfigDefault))
	}
	if config.GetBoolean("COMPRESS_ENABLED") {
		backend.Fiber.Use(compress.New())
	}
	if config.GetBoolean("RECOVER_ENABLED") {
		backend.Fiber.Use(recover.New())
	}
	if config.GetBoolean("CSRF_ENABLED") {
		backend.Fiber.Use(csrf.New())
	}

	// Logger middleware
	file, err := os.OpenFile("./logs/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		backend.Logger.Fatal("Error opening file: %v", err)
	}
	backend.Fiber.Use(logger.New(logger.Config{
		Format:     "[${time}] [http] [${status}] ${method} ${path} - ${ip} | ${latency}\n",
		TimeFormat: "02.01.2006T15:04:05",
		TimeZone:   "Europe/Berlin",
		Output:     file,
	}))
	backend.Logger.SetFile(file)

	if config.GetBoolean("METRICS_ENABLED") {
		metrics := backend.Fiber.Group("/api/v1/metrics")
		metrics.Use(basicauth.New(basicauth.Config{
			Users: map[string]string{
				config.GetString("METRICS_USER"): config.GetString("METRICS_PASSWORD"),
			},
		}))
		metrics.Get("/", monitor.New())
	}

	if config.GetBoolean("CORS_ENABLED") {
		backend.Fiber.Use(cors.New(cors.Config{
			AllowCredentials: config.GetBoolean("CORS_ALLOW_CREDENTIALS"),
			AllowOrigins:     config.GetString("CORS_ALLOW_ORIGINS"),
			AllowMethods:     config.GetString("CORS_ALLOW_METHODS"),
		}))
	}

	utils.StripeEnabledExec(func() {
		stripe.Key = config.GetString("STRIPE_SECRET")
	})

	switch config.GetString("DB") {
	case data_const.DBTypeSurreal:
		backend.Database = data_db.NewSurrealDB()
	case data_const.DBTypeMySQL:
		backend.Logger.Fatal("MySQL is not supported yet:")
		return
	case data_const.DBTypePostgres:
		backend.Logger.Fatal("Postgres is not supported yet:")
		return
	case data_const.DBTypeSQLite:
		backend.Logger.Fatal("SQLite is not supported yet:")
		return
	default:
		backend.Logger.Fatal("Unknown database type: %s", config.GetString("DB"))
		return
	}

	// TODO: controller setup

	go backend.StartServer()
	backend.SysCallSetup()
}

func (backend *Backend) StartServer() {
	backend.Database.Init()
	backend.Logger.Debug("Server starting on port %d.", config.GetInteger("PORT"))
	if err := backend.Fiber.Listen(fmt.Sprintf(":%d", config.GetInteger("PORT"))); err != nil {
		slog.Fatalf("Error starting server: %v", err)
	}
}

func (backend *Backend) Stop() {
	backend.Fiber.Shutdown()
	backend.Logger.File.Close()
	backend.Database.Close()
	slog.Flush()
}

// SysCallSetup sets up the system call handler.
// It should be called after the webservers start at the very end.
func (backend *Backend) SysCallSetup() {
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
