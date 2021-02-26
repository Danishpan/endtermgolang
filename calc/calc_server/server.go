package calc_server

import (
	"EndtermPracticeTask/calc/calcpb"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type Server struct {
	calcpb.UnimplementedCalculatorServiceServer
}

func (s *Server) PrimeNumberDecomposition(req *calcpb.PrimeNumberDecompositionRequest, stream calcpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v \n", req)

	n := req.GetX()

	// Get the number of 2s that divide n
	for n%2 == 0 {
		stream.Send(&calcpb.PrimeNumberDecompositionResponse{
			Result: 2,
		})
		n = n / 2
	}

	// n must be odd at this point. so we can skip one element
	// (note i = i + 2)
	for i := int32(3); i*i <= n; i = i + 2 {
		// while i divides n, append i and divide n
		for n%i == 0 {
			stream.Send(&calcpb.PrimeNumberDecompositionResponse{
				Result: i,
			})
			n = n / i
		}
	}

	// This condition is to handle the case when n is a prime number
	// greater than 2
	if n > 2 {
		stream.Send(&calcpb.PrimeNumberDecompositionResponse{
			Result: n,
		})
	}

	return nil
}

func (s *Server) ComputeAverage(stream calcpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with a streaming request\n")

	sum := int32(0)
	currentNumber := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Result: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		currentNumber = req.GetX()
		sum += currentNumber
		count++
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4000")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	calcpb.RegisterCalculatorServiceServer(s, &Server{})
	log.Println("Server is running on port:4000")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
