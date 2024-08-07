package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Hanyue-s-FYP/Marcom-Backend/middleware"
	simulation_pb "github.com/Hanyue-s-FYP/Marcom-Backend/proto_simulation"
	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

func Main() {
	router := http.NewServeMux()
	config := utils.NewConfig(".env.development")

	SetupRouter(router)

	middlewares := middleware.Use(
		middleware.Auth,
		middleware.RequestLogger,
		middleware.Cors,
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

// to test grpc
func main() {
    utils.NewConfig(".env.development")
    var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
        wg.Add(1)

		go utils.UseSimulationClient(func(client simulation_pb.SimulationServiceClient) {
			slog.Info("Sending simulation data to simulation server")
			startSimulationResp, err := client.StartSimulation(context.Background(), &simulation_pb.SimulationRequest{
				Id:      int32(i),
				EnvDesc: "Testing 123",
				Agents: []*simulation_pb.Agent{
					{Id: 1, Name: "Test Agent", Desc: "Test Agent", Attrs: []*simulation_pb.AgentAttribute{
						{Key: "Test", Value: "Ing"},
					},
					},
				},
				Products: []*simulation_pb.Product{
					{Id: 1, Name: "Test Product", Desc: "Test Product", Price: 123.0, Cost: 100.0},
					{Id: 1, Name: "Test Product 2", Desc: "Test Product", Price: 123.0, Cost: 100.0},
				},
                TotalCycles: 5,
			})
			if err != nil {
				slog.Error(fmt.Sprintf("failed to send simulation data to grpc server: %v", err))
			}
			slog.Info(fmt.Sprintf("Start Simulation Response: %s", startSimulationResp.Message))
            wg.Done()
		})

        go utils.UseSimulationClient(func(client simulation_pb.SimulationServiceClient) {
            slog.Info("Requesting StreamUpdates from grpc server")
            stream, err := client.StreamUpdates(context.Background(), &simulation_pb.StreamRequest{SimulationId: int32(i)})
            if err != nil {
                slog.Error(fmt.Sprintf("failed to call StreamUpdates on simulation server: %v", err))
            }
            for {
                simulationEvent, err := stream.Recv()
                if err == io.EOF {
                    break
                }
                if err != nil {
                    slog.Error(fmt.Sprintf("%v.StreamUpdates(_) = _, %v", client, err))
                }
                slog.Info(fmt.Sprintf("Got event from simulation server: %v", simulationEvent))
            }
        })
	}
    wg.Wait()
}
