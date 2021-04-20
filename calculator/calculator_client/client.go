package main

import (
	"context"
	"github.com/jirawan-chuapradit/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
)

func main()  {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := calculatorpb.NewCalculatorServiceClient(conn)

	doUnary(client)
}

func doUnary(client calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.CalculatorRequest{
		Number1: 10,
		Number2: 3,
	}

	res, err := client.Calculator(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calculator RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Sum)
}