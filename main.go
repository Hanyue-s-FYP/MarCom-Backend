package main

import (
	"fmt"
	"net/http"

	"github.com/Hanyue-s-FYP/Marcom-Backend/middleware"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

func main() {
	router := http.NewServeMux()
	config := utils.NewConfig(".env.development")

	SetupRouter(router)

	middlewares := middleware.Use(
		middleware.Cors,
		middleware.RequestLogger,
		middleware.Auth,
	)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", config.PORT),
		Handler: middlewares(router),
	}

	fmt.Printf("Starting to listen on port :%s\n", config.PORT)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Failed to start and listen to port %s: %v\n", config.PORT, err)
		panic(err) // cant even start listen d what else to do lol
	}
}
