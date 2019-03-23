package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/jorgeAM/udemy/basic/greet/greet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *server) LongGreet(stream greet.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked")
	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greet.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream  %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
}

func (s *server) GreetEveryone(stream greet.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked\n")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream  %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "!"
		res := &greet.GreetEveryoneResponse{
			Result: result,
		}

		if err := stream.Send(res); err != nil {
			log.Fatalf("Error while sending data to client  %v", err)
			return err
		}

	}
}

func (s *server) GreetWithDeadLine(ctx context.Context, req *greet.GreetWithDeadLineRequest) (*greet.GreetWithDeadLineResponse, error) {
	fmt.Printf("GreetWithDeadLine function was invoked with %v ", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "The client canceled the request.")
		}

		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greet.GreetWithDeadLineResponse{
		Result: result,
	}

	return res, nil
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
