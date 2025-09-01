package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	controlpb "example.com/remote/proto"
	svc "example.com/remote/server/control"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	controlpb.RegisterDeviceRegistryServer(grpcServer, &svc.Server{Store: &svc.MemoryStore{}})
	log.Println("control server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
