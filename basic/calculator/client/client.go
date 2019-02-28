package main

import (
	"context"
	"fmt"
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
	doUnary(client)
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
