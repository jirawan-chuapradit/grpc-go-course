package main

import (
	"context"
	"fmt"
	"github.com/jirawan-chuapradit/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {}

func (*server) GreetManyTimes(request *greetpb.GreetManytimesRequest, timesServer greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", request)

	firstName := request.GetGreeting().GetFirstName()
	for i := 0; i< 10 ; i++ {
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
	greetpb.RegisterGreetServiceServer(s,&server{})

	if err:= s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
