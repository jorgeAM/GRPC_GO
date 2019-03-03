package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/jorgeAM/udemy/basic/greet/greet"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()
	client := greet.NewGreetServiceClient(conn)
	//doUnary(client)
	//doServerStreaming(client)
	doClientStreaming(client)
}

func doUnary(client greet.GreetServiceClient) {
	fmt.Println("Starting do Unary RPC")
	req := &greet.GreetRequest{
		Greeting: &greet.Greeting{
			FirstName: "Jorge",
			LastName:  "Alfaro",
		},
	}
	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Fatalf("Response from Greet: %v", res.Result)
}

func doServerStreaming(client greet.GreetServiceClient) {
	fmt.Println("Starting do Server Streaming RPC")
	req := &greet.GreetManyTimesRequest{
		Greeting: &greet.Greeting{
			FirstName: "Jorge",
			LastName:  "Alfaro",
		},
	}
	stream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while response: %v", err)
		}

		fmt.Printf("Response from GreetManyTimes: %v \n", res.Result)
	}
}

func doClientStreaming(client greet.GreetServiceClient) {
	req := []*greet.LongGreetRequest{
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Jorge",
				LastName:  "Alfaro",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Basti",
				LastName:  "Murga",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Pancho",
				LastName:  "Alfaro",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Lili",
				LastName:  "Alfaro",
			},
		},
	}
	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet RPC: %v", err)
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

	log.Printf("Response from LongGreet: %v \n", res.Result)
}
