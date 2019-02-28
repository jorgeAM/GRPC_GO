package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/jorgeAM/udemy/basic/calculator/pb"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	fmt.Printf("Greet function was invoked with %v ", req)
	result := req.GetAdd().GetA() + req.GetAdd().GetB()
	res := &pb.AddResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listening: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
