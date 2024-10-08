package main

import (
	"fmt"
	"net/http"

	"github.com/Hanyue-s-FYP/Marcom-Backend/modules/agent"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules/dashboard"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules/environment"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules/product"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules/simulation"
	"github.com/Hanyue-s-FYP/Marcom-Backend/modules/user"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

func SetupRouter(r *http.ServeMux) {
	r.HandleFunc("OPTIONS /*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // web (or axios) requires preflight to response with status ok for some reason
	})

	// testing routes
	r.HandleFunc("GET /auth-test", func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseJSON(w, &struct{ Message string }{Message: fmt.Sprintf("Welcome to auth test %s\n", r.Header.Get("userId"))}, 200)
	})

	r.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseJSON(w, &struct{ Message string }{Message: "world"}, 200)
	})

	// Auth routes
	r.HandleFunc("POST /login", utils.MakeHttpHandler(user.Login))
	r.HandleFunc("POST /register-business", utils.MakeHttpHandler(user.RegisterBusiness, 201))
	r.HandleFunc("GET /get-me", utils.MakeHttpHandler(user.GetMe))
	r.HandleFunc("GET /check-username/{username}", utils.MakeHttpHandler(user.CheckUserWithUsername))
	r.HandleFunc("POST /forget-password", utils.MakeHttpHandler(user.ForgetPassword))
	r.HandleFunc("POST /reset-password", utils.MakeHttpHandler(user.ResetPassword))

	// Business routes
	r.HandleFunc("GET /business/{id}", utils.MakeHttpHandler(user.GetBusiness))
	r.HandleFunc("PUT /business", utils.MakeHttpHandler(user.UpdateBusiness))

	// Dashboard routes
	r.HandleFunc("GET /dashboard/products", utils.MakeHttpHandler(dashboard.GetDashboardProduct))
	r.HandleFunc("GET /dashboard/agents", utils.MakeHttpHandler(dashboard.GetDashboardAgent))
	r.HandleFunc("GET /dashboard/environments", utils.MakeHttpHandler(dashboard.GetDashboardEnvironment))

	// Product routes
	r.HandleFunc("GET /products", utils.MakeHttpHandler(product.GetAllProducts))
	r.HandleFunc("GET /products/{id}", utils.MakeHttpHandler(product.GetProduct))
	r.HandleFunc("GET /business-products/{id}", utils.MakeHttpHandler(product.GetAllProductsByBusiness))
	r.HandleFunc("POST /products", utils.MakeHttpHandler(product.CreateProduct))
	r.HandleFunc("PUT /products", utils.MakeHttpHandler(product.UpdateProduct))
	r.HandleFunc("DELETE /products/{id}", utils.MakeHttpHandler(product.DeleteProduct))
	r.HandleFunc("GET /products/report/{id}", utils.MakeHttpHandler(product.GetProductCompetitorReport))

	// Agent routes
	r.HandleFunc("GET /agents", utils.MakeHttpHandler(agent.GetAllAgents))
	r.HandleFunc("GET /agents/{id}", utils.MakeHttpHandler(agent.GetAgent))
	r.HandleFunc("GET /business-agents/{id}", utils.MakeHttpHandler(agent.GetAllAgentsByBusiness))
	r.HandleFunc("POST /agents", utils.MakeHttpHandler(agent.CreateAgent))
	r.HandleFunc("PUT /agents", utils.MakeHttpHandler(agent.UpdateAgent))
	r.HandleFunc("DELETE /agents/{id}", utils.MakeHttpHandler(agent.DeleteAgent))

	// Environment routes
	r.HandleFunc("GET /environments", utils.MakeHttpHandler(environment.GetAllEnvironments))
	r.HandleFunc("GET /environments/{id}", utils.MakeHttpHandler(environment.GetEnvironment))
	r.HandleFunc("GET /environments/has-product/{id}", utils.MakeHttpHandler(environment.GetSimplifiedEnvironmentsWithProduct))
	r.HandleFunc("GET /environments/has-agent/{id}", utils.MakeHttpHandler(environment.GetSimplifiedEnvironmentsWithAgent))
	r.HandleFunc("GET /business-environments/{id}", utils.MakeHttpHandler(environment.GetAllEnvironmentsByBusiness))
	r.HandleFunc("POST /environments", utils.MakeHttpHandler(environment.CreateEnvironment))
	r.HandleFunc("PUT /environments", utils.MakeHttpHandler(environment.UpdateEnvironment))
	r.HandleFunc("DELETE /environments/{id}", utils.MakeHttpHandler(environment.DeleteEnvironment))

	// Simulation routes
	r.HandleFunc("GET /simulations", utils.MakeHttpHandler(simulation.GetAllSimulations))
	r.HandleFunc("GET /simulations/{id}", utils.MakeHttpHandler(simulation.GetSimulation))
	r.HandleFunc("GET /simulations/of-environment/{id}", utils.MakeHttpHandler(simulation.GetSimulationsOfEnvID))
	r.HandleFunc("GET /business-simulations/{id}", utils.MakeHttpHandler(simulation.GetSimulationsByBusinessID))
	r.HandleFunc("POST /simulations", utils.MakeHttpHandler(simulation.CreateSimulation))
	r.HandleFunc("PUT /simulations", utils.MakeHttpHandler(simulation.UpdateSimulation))
	r.HandleFunc("DELETE /simulations/{id}", utils.MakeHttpHandler(simulation.DeleteSimulation))
	r.HandleFunc("GET /simulations/cycles/{id}", utils.MakeHttpHandler(simulation.GetSimulationCyclesBySimID))
	r.HandleFunc("GET /simulations/cycle/{id}", utils.MakeHttpHandler(simulation.GetSimulationCycleByCycleID)) // get specific cycle
	r.HandleFunc("GET /simulations/start/{id}", utils.MakeHttpHandler(simulation.StartSimulation))
	r.HandleFunc("GET /simulations/pause/{id}", utils.MakeHttpHandler(simulation.PauseSimulation))
	r.HandleFunc("GET /simulations/listen-for-event/{id}", simulation.ListenToSimulationUpdates)

	// handle images (load the path from config)
	imgPrefix := fmt.Sprintf("/%s/", utils.GetConfig().IMG_FOLDER)
	r.Handle(imgPrefix, http.StripPrefix(imgPrefix, http.FileServer(http.Dir(utils.GetConfig().IMG_FOLDER))))
}
