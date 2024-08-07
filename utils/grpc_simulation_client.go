package utils

import (
	"fmt"
	"log/slog"

	simulation_pb "github.com/Hanyue-s-FYP/Marcom-Backend/proto_simulation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// creates the grpc simulation service client and calls callback with the client
func UseSimulationClient(callback func(client simulation_pb.SimulationServiceClient)) {
    config := GetConfig()

	conn, err := grpc.NewClient(config.GRPC_SIMULATION_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(fmt.Sprintf("fail to dial: %v", err))
	}
	defer conn.Close()

	client := simulation_pb.NewSimulationServiceClient(conn)
    callback(client)
}
