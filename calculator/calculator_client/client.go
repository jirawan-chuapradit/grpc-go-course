package main

import (
	"context"
	"fmt"
	"github.com/jirawan-chuapradit/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main()  {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := calculatorpb.NewCalculatorServiceClient(conn)

	//doUnary(client)
	//doServerStreaming(client)
	doClientStreaming(client)
}

func doClientStreaming(client calculatorpb.CalculatorServiceClient) {

	nLists := []int64{1,2,3,4}
	steam, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v\n", err)
	}

	for _, n := range nLists {
		req := &calculatorpb.ComputeAverageRequest{
			Number: n,
		}
		err := steam.Send(req)
		if err != nil {
			log.Fatalf("error while send request to ComputeAverage: %v\n", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := steam.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v\n", err)
	}

	fmt.Printf("ComputeAverage Response: %v\n", res)

}

func doServerStreaming(client calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberManyTimesRequest{
		Number: 120,
	}

	resStream, err := client.PrimeNumberDecomposition(context.Background(),req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}else if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
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