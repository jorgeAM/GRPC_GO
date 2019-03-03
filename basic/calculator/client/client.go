package main

import (
	"context"
	"fmt"
	"io"
	"log"

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
	doClientStreaming(client)
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
