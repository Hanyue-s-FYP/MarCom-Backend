package main

import (
	"fmt"
	"net/http"

	"github.com/Hanyue-s-FYP/Marcom-Backend/modules/agent"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules/product"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules/user"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

func SetupRouter(r *http.ServeMux) {
	// testing routes
	r.HandleFunc("GET /auth-test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to auth test %s\n", r.Header.Get("userId"))
	})

	r.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

	// Auth routes
	r.HandleFunc("POST /login", utils.MakeHttpHandler(user.Login))
	r.HandleFunc("POST /register-business", utils.MakeHttpHandler(user.RegisterBusiness, 201))

	// Product routes
	r.HandleFunc("GET /products", utils.MakeHttpHandler(product.GetAllProduct))
	r.HandleFunc("GET /products/{id}", utils.MakeHttpHandler(product.GetProduct))
    r.HandleFunc("GET /business-products", utils.MakeHttpHandler(product.GetAllProductByBusiness)) // need this otherwise business won't be able to retrieve as their id is in the header already
    r.HandleFunc("GET /business-products/{id}", utils.MakeHttpHandler(product.GetAllProductByBusiness))
	r.HandleFunc("POST /products", utils.MakeHttpHandler(product.CreateProduct))
	r.HandleFunc("PUT /products", utils.MakeHttpHandler(product.UpdateProduct))
	r.HandleFunc("DELETE /products/{id}", utils.MakeHttpHandler(product.DeleteProduct))

	// Agent routes
	r.HandleFunc("GET /agents", utils.MakeHttpHandler(agent.GetAllAgent))
	r.HandleFunc("GET /agents/{id}", utils.MakeHttpHandler(agent.GetAgent))
    r.HandleFunc("GET /business-agents", utils.MakeHttpHandler(agent.GetAllAgentByBusiness)) 
    r.HandleFunc("GET /business-agents/{id}", utils.MakeHttpHandler(agent.GetAllAgentByBusiness))
	r.HandleFunc("POST /agents", utils.MakeHttpHandler(agent.CreateAgent))
	r.HandleFunc("PUT /agents", utils.MakeHttpHandler(agent.UpdateAgent))
	r.HandleFunc("DELETE /agents/{id}", utils.MakeHttpHandler(agent.DeleteAgent))
}
