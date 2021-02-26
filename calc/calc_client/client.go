package calc_client

import (
	"EndtermPracticeTask/calc/calcpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := calcpb.NewCalculatorServiceClient(conn)

	doPrimeNumberDecomposition(c)
	doComputeAverage(c)
}

func doPrimeNumberDecomposition(c calcpb.CalculatorServiceClient) {
	req := &calcpb.PrimeNumberDecompositionRequest{
		X: 120,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reciving from PrimeNumberDecomposition RPC %v", err)
		}
		fmt.Printf("%v \n", res.GetResult())
	}
}

func doComputeAverage(c calcpb.CalculatorServiceClient) {
	requests := []int32{1, 2, 3, 4, 5}

	ctx := context.Background()
	stream, err := c.ComputeAverage(ctx)
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending number: %v\n", req)
		stream.Send(&calcpb.ComputeAverageRequest{
			X: req,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}

	fmt.Printf("Average is: %v\n", res.GetResult())
}
