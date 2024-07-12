package main

import (
	"fmt"
	"net/http"

	"github.com/Hanyue-s-FYP/Marcom-Backend/modules/user"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

func SetupRouter(r *http.ServeMux) {
	// testing routes
	r.HandleFunc("GET /auth_test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to auth test %s\n", r.Header.Get("userId"))
	})

	r.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

    // Auth routes
    r.HandleFunc("POST /login", utils.MakeHttpHandler(user.Login))
}
