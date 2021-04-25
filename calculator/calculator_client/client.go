package main

import (
	"context"
	"fmt"
	"github.com/jirawan-chuapradit/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//doClientStreaming(client)
	//doBiDiStreaming(client)
	doErrorUnary(client)
}

func doErrorUnary(client calculatorpb.CalculatorServiceClient) {
	// correct call
	correctRequest := &calculatorpb.SquareRootRequest{
		Number: -10,
	}
	res, err := client.SquareRoot(context.Background(),correctRequest)
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			log.Fatalf("big error calling SquareRoot: %v\n", err)
			return
		}else {
			fmt.Println(status.Message())
			fmt.Println(status.Code())
			if status.Code() == codes.InvalidArgument {
				fmt.Println("we probably sent a negative number!")
				return
			}
		}
	}

	fmt.Printf("result of square root: %v\n", res.GetNumberRoot())

}



func doBiDiStreaming(client calculatorpb.CalculatorServiceClient) {
	stream, err := client.FindMaximum(context.Background())
	if err !=nil {
		log.Fatalf("error while opening stream and calling maximum: %v", err)
	}

	waitc := make(chan struct{})

	// send go routine
	go func() {
		number := []int32{4,7,2,19,4,6,30}
		for _, req := range number{
			 stream.Send(&calculatorpb.FindMaximumRequest{
				Number: req,
			})
			time.Sleep(1000*time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("error while recieve data from server stream: %v", err)
				break
			}

			maximum := res.GetMaximum()
			fmt.Printf("recieve a new maximum: %v\n", maximum)
		}
		close(waitc)
	}()

	<- waitc
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