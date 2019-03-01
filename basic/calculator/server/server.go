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

func (s *server) Prime(req *pb.PrimeRequest, stream pb.CalculatorService_PrimeServer) error {
	fmt.Printf("Prime function was invoked with %v ", req)
	number := req.GetNumber()
	var i uint32 = 2
	for number > i {
		if number%i == 0 {
			number = number / i
			res := &pb.PrimeResponse{
				Number: i,
			}
			if err := stream.Send(res); err != nil {
				log.Fatalf("something get wrong : %v", err)
				return err
			}
		} else {
			i++
		}
	}

	return nil
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
