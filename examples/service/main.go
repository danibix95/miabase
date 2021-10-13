package main

import (
	"net/http"

	"github.com/danibix95/miabase"
	"github.com/danibix95/miabase/pkg/response"
	"github.com/mia-platform/configlib"
)

type Environment struct {
	LogLevel string
	HTTPPort int
}

var envVariablesConfig = []configlib.EnvConfig{
	{
		Key:          "LOG_LEVEL",
		Variable:     "LogLevel",
		DefaultValue: "info",
	},
	{
		Key:          "HTTP_PORT",
		Variable:     "HTTPPort",
		DefaultValue: "3000",
	},
}

func main() {
	var env Environment
	miabase.LoadEnv(envVariablesConfig, &env)

	service := miabase.NewService("example", "v0.0.1", env.LogLevel)

	service.Plugin.Get("/greet", func(rw http.ResponseWriter, r *http.Request) {
		response.JSON(rw, map[string]string{"message": "welcome"})
	})

	service.Start(env.HTTPPort)
}
