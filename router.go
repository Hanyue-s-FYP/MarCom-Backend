package main

import (
	"fmt"
	"net/http"

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
	r.HandleFunc("POST /products", utils.MakeHttpHandler(product.CreateProduct))
	r.HandleFunc("PUT /products", utils.MakeHttpHandler(product.UpdateProduct))
	r.HandleFunc("DELETE /products/{id}", utils.MakeHttpHandler(product.DeleteProduct))
}
