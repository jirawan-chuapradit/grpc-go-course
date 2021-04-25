package main

import (
	"context"
	"fmt"
	"github.com/jirawan-chuapradit/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {

	fmt.Println("Hello I'm a client.")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Created client: %f", c)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doUnaryWithDeadline(c, 5*time.Second) // should complete
	doUnaryWithDeadline(c, 1*time.Second) // should timeout
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a UnaryWithDeadline")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName: "Maarek",
		}	,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok :=   status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("timeout was hit! deadline was exceeded")
			}else {
				fmt.Printf("unexpected error: %v", statusErr)
			}
		}
		log.Fatalf("error while calling GreetWithDeadline RPC: %v", err)
		return
	}
	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC ... ")

	//	 we create a stream by invoke the client
	steam, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "user1",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "user2",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "user3",
			},
		},
	}

	waitc := make(chan struct{})
	// we send a bunch of message to the client (go routine)
	go func() {
		//	function to send a bunch of message
		for _, req := range requests {
			fmt.Printf("send message: %v\n",req)
			steam.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		steam.CloseSend()
	}()
	// we receive a bunch of message from the client go routine)
	go func() {
		//	function to send a bunch of message
		for {
			res, err := steam.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v", err)
				break
			}
			fmt.Printf("received: %v\n", res.GetResult())
		}
		close(waitc)
	}()
	// block util every thing is done.
	<-waitc
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
		fmt.Printf("Sending request: %v\n", req)
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
			LastName:  "lastname test",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
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

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "step",
			LastName:  "lastname",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v\n", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}
