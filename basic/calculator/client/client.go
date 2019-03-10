package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/jorgeAM/udemy/basic/calculator/pb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewCalculatorServiceClient(conn)
	//doUnary(client)
	//doServerStreaming(client)
	//doClientStreaming(client)
	//doBiDiStreaming(client)
	doErrorUnary(client)
}

func doUnary(client pb.CalculatorServiceClient) {
	fmt.Println("Starting do Unary RPC")
	req := &pb.AddRequest{
		Add: &pb.Add{
			A: 324,
			B: 120,
		},
	}

	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Fatalf("Response from Greet: %v", res.Result)
}

func doServerStreaming(client pb.CalculatorServiceClient) {
	fmt.Println("Starting do Server Streaming RPC")
	req := &pb.PrimeRequest{
		Number: 120,
	}
	stream, err := client.Prime(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Prime RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while response: %v", err)
		}

		fmt.Printf("Response from Prime: %v \n", res.Number)
	}
}

func doClientStreaming(client pb.CalculatorServiceClient) {
	fmt.Println("Starting do Client Streaming RPC")
	req := []*pb.AverageRequest{
		&pb.AverageRequest{
			Number: 1,
		},
		&pb.AverageRequest{
			Number: 2,
		},
		&pb.AverageRequest{
			Number: 3,
		},
		&pb.AverageRequest{
			Number: 4,
		},
	}
	stream, err := client.Average(context.Background())
	if err != nil {
		log.Fatalf("error while calling Average RPC: %v", err)
	}

	for _, i := range req {
		if err := stream.Send(i); err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("error while reading request: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while response: %v", err)
	}

	log.Printf("Response from Average: %v \n", res.Result)
}

func doBiDiStreaming(client pb.CalculatorServiceClient) {
	req := []*pb.MaxRequest{
		&pb.MaxRequest{
			Number: 1,
		},
		&pb.MaxRequest{
			Number: 5,
		},
		&pb.MaxRequest{
			Number: 3,
		},
		&pb.MaxRequest{
			Number: 6,
		},
		&pb.MaxRequest{
			Number: 2,
		},
		&pb.MaxRequest{
			Number: 20,
		},
	}
	stream, err := client.Max(context.Background())
	if err != nil {
		log.Fatalf("error while calling Max RPC: %v", err)
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client
	go func() {
		for _, i := range req {
			fmt.Printf("Sending message: %v \n", i)
			stream.Send(i)
		}

		stream.CloseSend()
	}()

	//we receive a bunch of message from the client
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while reciving %v", err)
				break
			}

			fmt.Printf("Received %v \n", res.GetNumber())
		}

		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func doErrorUnary(client pb.CalculatorServiceClient) {
	fmt.Println("Starting do Unary RPC")
	req := &pb.SquareRootRequest{
		Number: -144,
	}

	doErrorCall(client, req)
}

func doErrorCall(client pb.CalculatorServiceClient, req *pb.SquareRootRequest) {
	res, err := client.SquareRoot(context.Background(), req)
	if err != nil {
		status, ok := status.FromError(err)
		if ok {
			fmt.Println(status.Message())
			fmt.Println(status.Code())
			if status.Code() == codes.InvalidArgument {
				fmt.Println("We probably send a negative number!")
				return
			}
		} else {
			log.Fatalf("error while calling SquareRoot RPC: %v", err)
			return
		}
	}

	log.Fatalf("Response from SquareRoot: %v", res.GetNumberRoot())
}
