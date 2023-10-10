package controllers

import (
	"github.com/yzaimoglu/yzgo/http/api"
	"github.com/yzaimoglu/yzgo/internal"
)

type Controller interface {
	Start()
}

type NoSQLController struct {
	Backend *internal.NoSQLBackend
}

type SQLController struct {
	Backend *internal.SQLBackend
}

func NewSQLController(backend *internal.SQLBackend) *SQLController {
	return &SQLController{
		Backend: backend,
	}
}

func NewNoSQLController(backend *internal.NoSQLBackend) *NoSQLController {
	return &NoSQLController{
		Backend: backend,
	}
}

func (controller *SQLController) Start() {
	backendFiber, backendLogger := internal.BackendSetup()
	controller.Backend.Fiber = backendFiber
	controller.Backend.Logger = backendLogger

	api.V1(controller.Backend.Fiber)
	controller.Backend.Start()
}

func (controller *NoSQLController) Start() {
	backendFiber, backendLogger := internal.BackendSetup()
	controller.Backend.Fiber = backendFiber
	controller.Backend.Logger = backendLogger

	api.V1(controller.Backend.Fiber)
	controller.Backend.Start()
}
