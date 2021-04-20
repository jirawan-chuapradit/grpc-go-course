package main

import (
	"context"
	"fmt"
	"github.com/jirawan-chuapradit/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {}

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
	fmt.Println("Calculator test ... ")
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
