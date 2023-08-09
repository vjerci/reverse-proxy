package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/vjerci/reverse-proxy/internal/app"
)

func main() {
	app, err := app.Build()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic(errors.New("port not set"))
	}

	log.Printf("starting proxy on port %s", port)

	err = http.ListenAndServe(port, app)

	if err != nil {
		if errors.Is(http.ErrServerClosed, err) {
			log.Print("got shutdown request... closing server, zZzzzZZZzzzz")
			return
		}

		panic(err)
	}
}
