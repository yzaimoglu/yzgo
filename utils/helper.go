package utils

import "github.com/yzaimoglu/yzgo/config"

func StripeEnabledExec(function func()) {
	if config.GetBoolean("STRIPE_ENABLED") {
		function()
	}
}
