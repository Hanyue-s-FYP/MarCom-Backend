package main

import (
	"fmt"
	"net/http"

	"github.com/Hanyue-s-FYP/Marcom-Backend/middleware"
)

func main() {
	router := http.NewServeMux()
    
    SetupRouter(router)

	middlewares := middleware.Use(
		middleware.RequestLogger,
        middleware.Auth,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewares(router),
	}

	fmt.Println("Starting to listen on port :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Failed to start and listen to port 8080: %v\n", err)
		panic(err) // cant even start listen d what else to do lol
	}
}
