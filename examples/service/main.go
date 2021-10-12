package main

import (
	"net/http"

	"github.com/danibix95/miabase"
	"github.com/danibix95/miabase/pkg/response"
)

func main() {
	service := miabase.NewService()

	service.Plugin.Get("/greet", func(rw http.ResponseWriter, r *http.Request) {
		response.JSON(rw, map[string]string{"message": "welcome"})
	})

	service.Start()
}
