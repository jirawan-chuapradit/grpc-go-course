package main

import (
	"context"
	"fmt"
	"github.com/jirawan-chuapradit/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"math"
	"net"
)

type server struct {}

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
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
