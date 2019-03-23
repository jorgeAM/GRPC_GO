package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

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
	//doClientStreaming(client)
	//doBiDiStreaming(client)
	doUnaryWithDeadLine(client, 5*time.Second) // Should complete
	doUnaryWithDeadLine(client, 1*time.Second) // Should timeout
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

func doBiDiStreaming(client greet.GreetServiceClient) {
	req := []*greet.GreetEveryoneRequest{
		&greet.GreetEveryoneRequest{
			Greeting: &greet.Greeting{
				FirstName: "Jorge",
				LastName:  "Alfaro",
			},
		},
		&greet.GreetEveryoneRequest{
			Greeting: &greet.Greeting{
				FirstName: "Basti",
				LastName:  "Murga",
			},
		},
		&greet.GreetEveryoneRequest{
			Greeting: &greet.Greeting{
				FirstName: "Pancho",
				LastName:  "Alfaro",
			},
		},
		&greet.GreetEveryoneRequest{
			Greeting: &greet.Greeting{
				FirstName: "Lili",
				LastName:  "Alfaro",
			},
		},
	}
	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreetEveryone RPC: %v", err)
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client
	go func() {
		for _, i := range req {
			fmt.Printf("Sending message: %v \n", i)
			stream.Send(i)
			time.Sleep(1000 * time.Millisecond)
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

			fmt.Printf("Received %v \n", res.GetResult())
		}

		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func doUnaryWithDeadLine(client greet.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting do Unary RPC with deadline")
	req := &greet.GreetWithDeadLineRequest{
		Greeting: &greet.Greeting{
			FirstName: "Jorge",
			LastName:  "Alfaro",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := client.GreetWithDeadLine(ctx, req)
	if err != nil {
		status, ok := status.FromError(err)
		if ok {
			if status.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("Unexpected error: %v", status)
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadLineRequest RPC: %v", err)
		}

		return
	}

	log.Fatalf("Response from GreetWithDeadLineRequest: %v", res.Result)
}
