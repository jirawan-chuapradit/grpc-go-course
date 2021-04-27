package main

import (
	"context"
	"fmt"
	"github.com/jirawan-chuapradit/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
)

type server struct {}

func (server) SquareRoot(ctx context.Context, request *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")

	number := request.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument,fmt.Sprintf("Recieved a negative number"))
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func (s server) FindMaximum(maximumServer calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a streaming request\n")

	maximum :=int32(0)

	for {
		req, err := maximumServer.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
			return err
		}

		number := req.GetNumber()
		if number > maximum {
			maximum = number
			err = maximumServer.Send(&calculatorpb.FindMaximumResponse{
				Maximum: maximum,
			})
			if err != nil {
				log.Fatalf("error while send data to client stream: %v", err)
				return err
			}
		}
	}

}

func (server) ComputeAverage(averageServer calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with a streaming request\n")
	var nLists []int64
	for {
		req, err := averageServer.Recv()
		if err == io.EOF {
			var total int64
			for _, n := range nLists {
				total += n
			}
			average := float32(total) / float32(len(nLists))
			return averageServer.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}

		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		nLists = append(nLists, req.Number)
	}

}

func (server) PrimeNumberDecomposition(request *calculatorpb.PrimeNumberManyTimesRequest, decompositionServer calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Prime number decomposition function was invoked with %v\n", request)
	k := 2;
	n := request.GetNumber()
	for {
		if n <= 1 {
			break
		}
		
		if math.Mod(float64(n),float64(k)) == 0 {
			result := k
			res := &calculatorpb.PrimeNumberManyTimesResponse{
				Result: int64(result),
			}
			decompositionServer.Send(res)
			n = n/int32(k)
		}else {
			k = k + 1
		}
	}
	return nil
}

func (*server) Calculator(ctx context.Context, request *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Calculator function was invoked with %v\n", request)
	num1 := request.GetNumber1()
	num2 := request.GetNumber2()
	sum := num1 + num2
	res := &calculatorpb.CalculatorResponse{
		Sum: sum,
	}
	return res ,nil
}

func main () {
	fmt.Println("Calculator server start ... ")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s:= grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
