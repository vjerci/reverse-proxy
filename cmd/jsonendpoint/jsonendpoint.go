package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		panic(errors.New("port not set"))
	}

	log.Printf("starting json endpoint on port %s", port)
	err := http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("jsonendpoint serving response for method %s url %s", r.Method, r.URL.String())

		jsonData := struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			IsAdmin   bool   `json:"is_admin"`
			Metrics   int    `json:"metrics"`
		}{
			FirstName: "mark",
			LastName:  "jonson",
			IsAdmin:   true,
			Metrics:   100,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(&jsonData)
		if err != nil {
			panic(fmt.Errorf("failed to encode json %s", err))
		}

	}))

	if err != nil {
		if errors.Is(http.ErrServerClosed, err) {
			log.Print("got shutdown request... closing json endpoint, zZzzzZZZzzzz")
			return
		}
	}
}
