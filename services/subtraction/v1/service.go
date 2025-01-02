package v1

import (
	"context"

	subtractionpb "example/gen/go/subtraction/v1"
	factorialpb "example/gen/go/x/factorial"
	"example/pkg/factorial"
)

var _ subtractionpb.SubtractionServiceServer = (*server)(nil)

type server struct{}

func New() *server {
	return &server{}
}

func (s *server) Subtract(ctx context.Context, request *subtractionpb.SubtractRequest) (*subtractionpb.SubtractResponse, error) {
	return &subtractionpb.SubtractResponse{
		Rez: request.GetA() - request.GetB(),
	}, nil
}

func (s *server) Factorial(ctx context.Context, request *factorialpb.FactorialRequest) (*factorialpb.FactorialResponse, error) {
	return &factorialpb.FactorialResponse{
		Res: int64(factorial.Calculate(int(request.GetNum()))),
	}, nil
}
