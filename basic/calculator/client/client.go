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
	doServerStreaming(client)
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
