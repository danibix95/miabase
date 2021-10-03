package main

import (
	"io"
	"net/http"

	"github.com/danibix95/miabase"
)

func main() {
	service := miabase.NewService()

	service.Plugin.Get("/greet", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		io.WriteString(rw, `{"msg": "welcome"}`)
	})

	service.Start()
}
