package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Hanyue-s-FYP/Marcom-Backend/middleware"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World\n")
	})

    middlewares := middleware.Use(
		middleware.RequestLogger,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewares(router),
	}

	fmt.Println("Starting to listen on port :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(fmt.Sprintf("Failed to start and listen to port 8080: %v", err))
	}
}
