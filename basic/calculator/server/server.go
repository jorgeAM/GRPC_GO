package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	for number >= i {
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

func (s *server) Average(stream pb.CalculatorService_AverageServer) error {
	var result, count float32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.AverageResponse{
				Result: result / count,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream  %v", err)
			return err
		}

		result += req.GetNumber()
		count++
	}
}

func (s *server) Max(stream pb.CalculatorService_MaxServer) error {
	var aux uint32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream  %v", err)
			return err
		}

		number := req.GetNumber()
		if number >= aux {
			aux = number
			res := &pb.MaxResponse{
				Number: number,
			}
			if err := stream.Send(res); err != nil {
				log.Fatalf("Error while sending data to client  %v", err)
				return err
			}
		}
	}
}

func (s *server) SquareRoot(ctx context.Context, req *pb.SquareRootRequest) (*pb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}

	return &pb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listening: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
