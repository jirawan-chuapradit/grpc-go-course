package main

import (
	"context"
	"fmt"
	"github.com/jirawan-chuapradit/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main () {

	fmt.Println("Hello I'm a client.")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer  cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Created client: %f", c)
	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client RPC ... ")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "user1",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "user2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "user3",
			},
		},

	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending request: %v\n",req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v\n", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC ...")
	req := &greetpb.GreetManytimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "firtname test",
			LastName: "lastname test",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil{
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream.
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doUnary (c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "step",
			LastName: "lastname",
		},
	}
	res, err := c.Greet(context.Background(),req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v\n", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}