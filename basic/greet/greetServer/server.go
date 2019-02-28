package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/jorgeAM/udemy/basic/greet/greet"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v ", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greet.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (s *server) GreetManyTimes(req *greet.GreetManyTimesRequest, stream greet.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v ", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greet.GreetManyTimesResponse{
			Result: "Hola " + firstName + " numbrer: " + strconv.Itoa(i),
		}
		if err := stream.Send(res); err != nil {
			log.Fatalf("something get wrong : %v", err)
			return err
		}

		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}

	grpcServe := grpc.NewServer()
	greet.RegisterGreetServiceServer(grpcServe, &server{})
	if err := grpcServe.Serve(lis); err != nil {
		log.Fatalf("Fallo al correr servidor: %v", err)
	}
}
