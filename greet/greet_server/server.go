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
	"net"
	"strconv"
	"time"
)

type server struct{}

func (*server) GreetWithDeadline(ctx context.Context, request *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked with a streaming request\n")
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// the client canceled the request
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.DeadlineExceeded, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := request.GetGreeting().GetFirstName()
	result := "Hello "+firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetEveryone(everyoneServer greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with a streaming request\n")
	for {
		req, err := everyoneServer.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error while reading client streaming: %v", err)
			return err
		}

		firstname := req.GetGreeting().GetFirstName()
		result := "Hello " + firstname + "! "
		err = everyoneServer.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if err != nil {
			log.Fatalf("error while sending data to client streaming: %v", err)
			return err
		}
	}
}

func (*server) LongGreet(greetServer greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")

	result := ""
	for {
		req, err := greetServer.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return greetServer.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "

	}
	return nil
}

func (*server) GreetManyTimes(request *greetpb.GreetManytimesRequest, timesServer greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", request)

	firstName := request.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " >>> number: " + strconv.Itoa(i)
		res := &greetpb.GreetManyTImesResponse{
			Result: result,
		}
		timesServer.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) Greet(ctx context.Context, request *greetpb.GreetRequest) (*greetpb.GreetingResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", request)
	firstName := request.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetingResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
