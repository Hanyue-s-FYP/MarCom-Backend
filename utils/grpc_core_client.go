package utils

import (
	"fmt"
	"log/slog"

	core_pb "github.com/Hanyue-s-FYP/Marcom-Backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// creates the grpc simulation service client and calls callback with the client
func UseCoreGRPCClient(callback func(client core_pb.MarcomServiceClient)) {
	config := GetConfig()

	conn, err := grpc.NewClient(config.GRPC_CORE_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(fmt.Sprintf("fail to dial: %v", err))
	}
	defer conn.Close()

	client := core_pb.NewMarcomServiceClient(conn)
	callback(client)
}
