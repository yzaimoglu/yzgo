package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yzaimoglu/yzgo/config"
)

func Debug() bool {
	return config.GetBoolean("DEBUG")
}

func DebugExec(function func()) {
	if Debug() {
		function()
	}
}

func Child() bool {
	return fiber.IsChild()
}

func NonChildExec(function func()) {
	if !Child() {
		function()
	}
}
